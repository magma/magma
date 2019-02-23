// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
//
// If you'd like to test diameter over SSL, generate SSL certificates:
//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
//
// And start the server with `-cert_file cert.pem -key_file key.pem`.
//
// By default this server runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
)

func main() {
	addr := flag.String("addr", ":3868", "address in the form of ip:port to listen on")
	ppaddr := flag.String("pprof_addr", ":9000", "address in form of ip:port for the pprof server")
	host := flag.String("diam_host", "server", "diameter identity host")
	realm := flag.String("diam_realm", "go-diameter", "diameter identity realm")
	certFile := flag.String("cert_file", "", "tls certificate file (optional)")
	keyFile := flag.String("key_file", "", "tls key file (optional)")
	silent := flag.Bool("s", false, "silent mode, useful for benchmarks")
	flag.Parse()

	// Load our custom dictionary on top of the default one, which
	// always have the Base Protocol (RFC6733) and Credit Control
	// Application (RFC4006).
	err := dict.Default.Load(bytes.NewReader([]byte(helloDictionary)))
	if err != nil {
		log.Fatal(err)
	}

	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         13,
		ProductName:      "go-diameter",
		FirmwareRevision: 1,
	}

	// Create the state machine (mux) and set its message handlers.
	mux := sm.New(settings)
	mux.Handle("HMR", handleHMR(*silent))
	mux.Handle("ACR", handleACR(*silent))
	mux.HandleFunc("ALL", handleALL) // Catch all.

	// Print error reports.
	go printErrors(mux.ErrorReports())

	if len(*ppaddr) > 0 {
		go func() { log.Fatal(http.ListenAndServe(*ppaddr, nil)) }()
	}

	err = listen(*addr, *certFile, *keyFile, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.Println(err)
	}
}

func listen(addr, cert, key string, handler diam.Handler) error {
	// Start listening for connections.
	if len(cert) > 0 && len(key) > 0 {
		log.Println("Starting secure diameter server on", addr)
		return diam.ListenAndServeTLS(addr, cert, key, handler, nil)
	}
	log.Println("Starting diameter server on", addr)
	return diam.ListenAndServe(addr, handler, nil)
}

func handleHMR(silent bool) diam.HandlerFunc {
	type HelloRequest struct {
		SessionID        datatype.UTF8String       `avp:"Session-Id"`
		OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
		DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
		DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
		UserName         string                    `avp:"User-Name"`
	}
	return func(c diam.Conn, m *diam.Message) {
		if !silent {
			log.Printf("Received HMR from %s:\n%s", c.RemoteAddr(), m)
		}
		var hmr HelloRequest
		if err := m.Unmarshal(&hmr); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}
		a := m.Answer(diam.Success)
		a.NewAVP(avp.SessionID, avp.Mbit, 0, hmr.SessionID)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, hmr.DestinationHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, hmr.DestinationRealm)
		a.NewAVP(avp.DestinationRealm, avp.Mbit, 0, hmr.OriginRealm)
		a.NewAVP(avp.DestinationHost, avp.Mbit, 0, hmr.OriginHost)
		_, err := a.WriteTo(c)
		if err != nil {
			log.Printf("Failed to write message to %s: %s\n%s\n",
				c.RemoteAddr(), err, a)
			return
		}
		if !silent {
			log.Printf("Sent HMA to %s:\n%s", c.RemoteAddr(), a)
		}
	}
}

func handleACR(silent bool) diam.HandlerFunc {
	type AccountingRequest struct {
		SessionID              *diam.AVP                 `avp:"Session-Id"`
		OriginHost             *diam.AVP                 `avp:"Origin-Host"`
		OriginRealm            *diam.AVP                 `avp:"Origin-Realm"`
		DestinationRealm       datatype.DiameterIdentity `avp:"Destination-Realm"`
		AccountingRecordType   *diam.AVP                 `avp:"Accounting-Record-Type"`
		AccountingRecordNumber *diam.AVP                 `avp:"Accounting-Record-Number"`
		DestinationHost        datatype.DiameterIdentity `avp:"Destination-Host"`
	}
	return func(c diam.Conn, m *diam.Message) {
		if !silent {
			log.Printf("Received ACR from %s\n%s", c.RemoteAddr(), m)
		}
		var acr AccountingRequest
		if err := m.Unmarshal(&acr); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}
		a := m.Answer(diam.Success)
		a.InsertAVP(acr.SessionID)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, acr.DestinationHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, acr.DestinationRealm)
		a.AddAVP(acr.AccountingRecordType)
		a.AddAVP(acr.AccountingRecordNumber)
		_, err := a.WriteTo(c)
		if err != nil {
			log.Printf("Failed to write message to %s: %s\n%s\n",
				c.RemoteAddr(), err, a)
			return
		}
		if !silent {
			log.Printf("Sent ACA to %s:\n%s", c.RemoteAddr(), a)
		}
	}
}

func handleALL(c diam.Conn, m *diam.Message) {
	log.Printf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
}

// helloDictionary is our custom, example dictionary.
var helloDictionary = xml.Header + `
<diameter>
	<application id="999" type="acct">
		<command code="111" short="HM" name="Hello-Message">
			<request>
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
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
