// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client example. This is by no means a complete client.
//
// If you'd like to test diameter over SSL, make sure the server supports
// it and add -ssl to the command line. To use client certificates,
// run the client with -ssl -cert_file cert.pem -key_file key.pem.
//
// When the client connects, the underlying state machine (diam/sm package)
// performs the handshake (CER/CEA) and returns a connection. If the
// client is configured with watchdog, it automatically sends DWR and
// handles DWA in background.
//
// The -hello command line flag makes the client connect, handshake,
// send a hello message, and disconnect. This is to demonstrate how to
// use custom dictionaries.
//
// The -bench option turns the client into a benchmark tool to test
// the server. It uses ACR/ACA messages for this.
//
// By default this client runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr := flag.String("addr", "localhost:3868", "address in form of ip:port to connect to")
	ssl := flag.Bool("ssl", false, "connect to server using tls")
	host := flag.String("diam_host", "client", "diameter identity host")
	realm := flag.String("diam_realm", "go-diameter", "diameter identity realm")
	certFile := flag.String("cert_file", "", "tls client certificate file (optional)")
	keyFile := flag.String("key_file", "", "tls client key file (optional)")
	hello := flag.Bool("hello", false, "send a hello message, wait for the response and disconnect")
	bench := flag.Bool("bench", false, "benchmark the server by sending ACR messages")
	benchCli := flag.Int("bench_clients", 1, "number of client connections")
	benchMsgs := flag.Int("bench_msgs", 1000, "number of ACR messages to send")
	networkType := flag.String("network_type", "tcp", "protocol type tcp/sctp")

	flag.Parse()
	if len(*addr) == 0 {
		flag.Usage()
	}

	// Load our custom dictionary on top of the default one.
	err := dict.Default.Load(bytes.NewReader([]byte(helloDictionary)))
	if err != nil {
		log.Fatal(err)
	}

	cfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(net.ParseIP("127.0.0.1")),
		},
	}

	// Create the state machine (it's a diam.ServeMux) and client.
	mux := sm.New(cfg)

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   5 * time.Second,
		AcctApplicationID: []*diam.AVP{
			// Advertise that we want support accounting application with id 999
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(helloApplication)),
		},
		AuthApplicationID: []*diam.AVP{
			// Advertise support for credit control application
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)), // RFC 4006
		},
	}

	// Set message handlers.
	done := make(chan struct{}, 1000)
	mux.Handle("HMA", handleHMA(done))
	mux.Handle("ACA", handleACA(done))

	// Print error reports.
	go printErrors(mux.ErrorReports())

	connect := func() (diam.Conn, error) {
		return dial(cli, *addr, *certFile, *keyFile, *ssl, *networkType)
	}

	if *bench {
		cli.EnableWatchdog = false
		benchmark(connect, cfg, *benchCli, *benchMsgs, done)
		return
	}

	if *hello {
		c, err := connect()
		if err != nil {
			log.Fatal(err)
		}
		err = sendHMR(c, cfg)
		if err != nil {
			log.Fatal(err)
		}
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			log.Fatal("timeout: no hello answer received")
		}
		return
	}

	// Makes a persisent connection with back-off.
	log.Println("Use wireshark to see the messages, or try -hello")
	backoff := 1
	for {
		c, err := connect()
		if err != nil {
			log.Println(err)
			backoff *= 2
			if backoff > 20 {
				backoff = 20
			}
			time.Sleep(time.Duration(backoff) * time.Second)
			continue
		}
		log.Println("Client connected, handshake ok")
		backoff = 1
		<-c.(diam.CloseNotifier).CloseNotify()
		log.Println("Client disconnected")
	}
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.Println(err)
	}
}

func dial(cli *sm.Client, addr, cert, key string, ssl bool, networkType string) (diam.Conn, error) {
	if ssl {
		return cli.DialNetworkTLS(networkType, addr, cert, key, nil)
	}
	return cli.DialNetwork(networkType, addr)
}

func sendHMR(c diam.Conn, cfg *sm.Settings) error {
	// Get this client's metadata from the connection object,
	// which is set by the state machine after the handshake.
	// It contains the peer's Origin-Host and Realm from the
	// CER/CEA handshake. We use it to populate the AVPs below.
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(helloMessage, helloApplication, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("foobar"))
	log.Printf("Sending HMR to %s\n%s", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return err
}

func handleHMA(done chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received HMA from %s\n%s", c.RemoteAddr(), m)
		close(done)
	}
}

func handleACA(done chan struct{}) diam.HandlerFunc {
	ok := struct{}{}
	return func(c diam.Conn, m *diam.Message) {
		done <- ok
	}
}

type dialFunc func() (diam.Conn, error)

func benchmark(df dialFunc, cfg *sm.Settings, ncli, msgs int, done chan struct{}) {
	var err error
	c := make([]diam.Conn, ncli)
	log.Println("Connecting", ncli, "clients...")
	for i := 0; i < ncli; i++ {
		c[i], err = df() // Dial and do CER/CEA handshake.
		if err != nil {
			log.Fatal(err)
		}
		defer c[i].Close()
	}
	log.Println("Done. Sending messages...")
	start := time.Now()
	for _, cli := range c {
		go sendACR(cli, cfg, msgs)
	}
	count := 0
	total := ncli * msgs
wait:
	for {
		select {
		case <-done:
			count++
			if count == total {
				break wait
			}
		case <-time.After(time.Second):
			log.Fatal("Timeout waiting for messages.")
		}
	}
	elapsed := time.Since(start)
	total = total * 2 // req+resp
	log.Printf("%d messages in %s: %d/s", total, elapsed,
		int(float64(total)/elapsed.Seconds()))
}

var eventRecord = datatype.Unsigned32(1) // RFC 6733: EVENT_RECORD 1

func sendACR(c diam.Conn, cfg *sm.Settings, n int) {
	// Get this client's metadata from the connection object,
	// which is set by the state machine after the handshake.
	// It contains the peer's Origin-Host and Realm from the
	// CER/CEA handshake. We use it to populate the AVPs below.
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		log.Fatal("Client connection does not contain metadata")
	}
	var err error
	var m *diam.Message
	for i := 0; i < n; i++ {
		m = diam.NewRequest(diam.Accounting, 0, c.Dictionary())
		m.NewAVP(avp.SessionID, avp.Mbit, 0,
			datatype.UTF8String(strconv.Itoa(i)))
		m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
		m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
		m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
		m.NewAVP(avp.AccountingRecordType, avp.Mbit, 0, eventRecord)
		m.NewAVP(avp.AccountingRecordNumber, avp.Mbit, 0,
			datatype.Unsigned32(i))
		m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
		if _, err = m.WriteTo(c); err != nil {
			log.Fatal(err)
		}
	}
}

// Example dictionary.

const (
	helloApplication = 999 // Our custom app from the dictionary below.
	helloMessage     = 111
)

// helloDictionary is our custom, example dictionary.
var helloDictionary = xml.Header + `
<diameter>
	<application id="999" type="acct">
		<command code="111" short="HM" name="Hello-Message">
			<request>
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Destination-Host" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
			</request>
			<answer>
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Error-Message" required="false" max="1"/>
			</answer>
		</command>
	</application>
</diameter>
`
