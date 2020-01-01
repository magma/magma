// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package main

// Base Diameter SCTP client example.

import (
	"flag"
	"io"
	"log"
	"net"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/ishidawataru/sctp"
)

var addr = flag.String("addr", "192.168.60.145:3868", "server address in form of ip:port to connect to")
var laddr = flag.String("laddr", "", "local address in form of [ip]:port to bind Dailer to")
var host = flag.String("diam_host", "magma-oai.openair4G.eur", "diameter identity host")
var realm = flag.String("diam_realm", "openair4G.eur", "diameter identity realm")
var vendorID = flag.Uint("vendor", 10415, "Vendor ID")
var appID = flag.Uint("app", 16777251, "AuthApplicationID")
var wait = flag.Int("wait", 10, "Time to wait for completion")

var sctpLAdds, localAddr *sctp.SCTPAddr
var originStateID = datatype.Unsigned32(time.Now().Unix())

func main() {
	flag.Parse()
	if len(*addr) == 0 {
		flag.Usage()
	}

	errc := make(chan error, 1)

	cmux := diam.NewServeMux()

	cmux.Handle("CEA", handleCEA(errc))
	cmux.Handle("DWR", handleDWR(errc))

	var err error
	if len(*laddr) > 0 {
		localAddr, err = sctp.ResolveSCTPAddr("sctp", *laddr)
		if err != nil {
			log.Fatalf("Invalid Local Address '%s': %v", *laddr, err)
		}
	}
	log.Printf("Connecting to SCTP server at %s", *addr)
	cli, err := diam.DialExt("sctp", *addr, cmux, nil, 0, localAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	sctpConn := cli.Connection().(*sctp.SCTPConn)
	sctpLAdds, err = sctpConn.SCTPLocalAddr(0)
	if err != nil {
		log.Fatal(err)
		return
	}
	sctpRAdds, err := sctpConn.SCTPRemoteAddr(0)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("SCTP Association %s -> %s", sctpLAdds, sctpRAdds)
	primaryAddr, err := sctpConn.SCTPGetPrimaryPeerAddr()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Primary peer address: %s", primaryAddr)
	log.Printf("Sending CER")
	_, err = sendCER(cli)
	if err != nil {
		log.Fatal(err)
		return
	}

	select {
	case err := <-errc:
		log.Fatal(err)
	case <-time.After(time.Second * time.Duration(*wait)):
		log.Printf("Completed\n")
	}
}

func sendCER(w io.Writer) (n int64, err error) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString(*host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString(*realm))
	m.NewAVP(
		avp.HostIPAddress,
		avp.Mbit,
		0,
		datatype.Address(net.ParseIP(sctpLAdds.IPAddrs[len(sctpLAdds.IPAddrs)-1].String())))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, originStateID)
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(*vendorID)),
		},
	})
	return m.WriteTo(w)
}

func handleCEA(errc chan error) diam.HandlerFunc {
	type CEA struct {
		OriginHost        string `avp:"Origin-Host"`
		OriginRealm       string `avp:"Origin-Realm"`
		VendorID          int    `avp:"Vendor-Id"`
		ProductName       string `avp:"Product-Name"`
		OriginStateID     int    `avp:"Origin-State-Id"`
		AcctApplicationID int    `avp:"Acct-Application-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var resp CEA
		err := m.Unmarshal(&resp)
		if err != nil {
			errc <- err
			return
		}
	}
}

func handleDWR(errc chan error) diam.HandlerFunc {
	// DWR is a Device-Watchdog-Request message.
	// See RFC 6733 section 5.5.1 for details.
	type DWR struct {
		OriginHost    datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm   datatype.DiameterIdentity `avp:"Origin-Realm"`
		OriginStateID datatype.Unsigned32       `avp:"Origin-State-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var dwr DWR
		err := m.Unmarshal(&dwr)
		if err != nil {
			errc <- err
			return
		}
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString(*host))
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString(*realm))
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, originStateID)
		_, err = a.WriteTo(c)
		if err != nil {
			errc <- err
		}
	}
}
