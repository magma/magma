// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/diamtest"
)

func TestCapabilitiesExchange(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", handleCER(errc, false))

	srv := diamtest.NewServer(smux, nil)
	defer srv.Close()

	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.HandleIdx(diam.CommandIndex{AppID: 0, Code: diam.CapabilitiesExchange, Request: false}, handleCEA(errc, wait))

	cli, err := diam.Dial(srv.Addr, cmux, nil)
	if err != nil {
		t.Fatal(err)
	}

	sendCER(cli)

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

func TestCapabilitiesExchangeTLS(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", handleCER(errc, true))

	srv := diamtest.NewUnstartedServer(smux, nil)
	tm := time.Second
	srv.Config.ReadTimeout = tm
	srv.Config.WriteTimeout = tm
	srv.TLS = &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS10,
	}
	srv.StartTLS()
	time.Sleep(time.Millisecond * 10) // let srv start
	defer srv.Close()
	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.Handle("CEA", handleCEA(errc, wait))

	cli, err := diam.DialTLS(srv.Addr, "", "", cmux, nil)
	if err != nil {
		t.Fatalf("diam.DialTLS Error: %v", err)
	}

	n, err := sendCER(cli)
	if err != nil {
		t.Fatalf("sendCER Error: %v", err)
	}
	if n <= 0 {
		t.Fatalf("sendCER: %d bytes sent", n)
	}

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(time.Second * 3):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

func sendCER(w io.Writer) (n int64, err error) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("cli"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1234))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1))
	return m.WriteTo(w)
}

func handleCER(errc chan error, useTLS bool) diam.HandlerFunc {
	type CER struct {
		OriginHost        string    `avp:"Origin-Host"`
		OriginRealm       string    `avp:"Origin-Realm"`
		VendorID          int       `avp:"Vendor-Id"`
		ProductName       string    `avp:"Product-Name"`
		OriginStateID     *diam.AVP `avp:"Origin-State-Id"`
		AcctApplicationID *diam.AVP `avp:"Acct-Application-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {

		if c.LocalAddr() == nil {
			errc <- fmt.Errorf("LocalAddr is nil")
		}
		if c.RemoteAddr() == nil {
			errc <- fmt.Errorf("LocalAddr is nil")
		}
		if useTLS && c.TLS() == nil {
			errc <- fmt.Errorf("TLS is nil")
		}
		if !useTLS && c.TLS() != nil {
			errc <- fmt.Errorf("TLS is supposed to be nil")
		}
		var req CER
		err := m.Unmarshal(&req)
		if err != nil {
			errc <- err
			return
		}
		if req.OriginHost != "cli" {
			errc <- fmt.Errorf("Unexpected OriginHost. Want cli, have %q", req.OriginHost)
			return
		}
		if req.OriginRealm != "localhost" {
			errc <- fmt.Errorf("Unexpected OriginRealm. Want localhost, have %q", req.OriginRealm)
			return
		}
		if req.VendorID != 99 {
			errc <- fmt.Errorf("Unexpected VendorID. Want 99, have %d", req.VendorID)
			return
		}
		if req.ProductName != "go-diameter" {
			errc <- fmt.Errorf("Unexpected ProductName. Want go-diameter, have %q", req.ProductName)
			return
		}
		a := m.Answer(diam.Success)
		_, err = sendCEA(c, a, req.OriginStateID, req.AcctApplicationID)
		if err != nil {
			errc <- err
		}
		c.(diam.CloseNotifier).CloseNotify()
		go func() {
			<-c.(diam.CloseNotifier).CloseNotify()
		}()
		//log.Println("Client", c.RemoteAddr(), "disconnected")
	}
}

func sendCEA(w io.Writer, m *diam.Message, OriginStateID, AcctApplicationID *diam.AVP) (n int64, err error) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.AddAVP(OriginStateID)
	m.AddAVP(AcctApplicationID)
	return m.WriteTo(w)
}

func handleCEA(errc chan error, wait chan struct{}) diam.HandlerFunc {
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
		if resp.OriginHost != "srv" {
			errc <- fmt.Errorf("Unexpected OriginHost. Want srv, have %q", resp.OriginHost)
			return
		}
		if resp.OriginRealm != "localhost" {
			errc <- fmt.Errorf("Unexpected OriginRealm. Want localhost, have %q", resp.OriginRealm)
			return
		}
		if resp.VendorID != 99 {
			errc <- fmt.Errorf("Unexpected VendorID. Want 99, have %d", resp.VendorID)
			return
		}
		if resp.ProductName != "go-diameter" {
			errc <- fmt.Errorf("Unexpected ProductName. Want go-diameter, have %q", resp.ProductName)
			return
		}
		if resp.OriginStateID != 1234 {
			errc <- fmt.Errorf("Unexpected OriginStateID. Want 1234, have %d", resp.OriginStateID)
			return
		}
		if resp.AcctApplicationID != 1 {
			errc <- fmt.Errorf("Unexpected AcctApplicationID. Want 1, have %d", resp.AcctApplicationID)
			return
		}
		// Initialize & start close notifier
		closeNotifyChan := c.(diam.CloseNotifier).CloseNotify()
		// Wait on close notify chan outside of main serve loop, closeNotifier routine is started by
		// liveSwitchReader.Read to avoid io.Pipe deadlock issue
		go func() {
			<-closeNotifyChan // wait on c.Close to complete
			select {          // close only if not already closed
			case <-wait:
			default:
				close(wait)
			}
		}()
		c.Close()
	}
}
