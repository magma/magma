// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/diamtest"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm/smparser"
)

func TestClient_Dial_MissingStateMachine(t *testing.T) {
	cli := &Client{}
	_, err := cli.Dial("")
	if err != ErrMissingStateMachine {
		t.Fatal(err)
	}
}

func TestClient_Dial_InvalidAddress(t *testing.T) {
	cli := &Client{
		Handler: New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0,
				datatype.Unsigned32(0)),
		},
	}
	c, err := cli.Dial(":0")
	if err == nil {
		c.Close()
		t.Fatal("Invalid client address succeeded")
	}
}

func TestClient_DialTLS_InvalidAddress(t *testing.T) {
	cli := &Client{
		Handler: New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(0)),
		},
	}
	c, err := cli.DialTLS(":0", "", "")
	if err == nil {
		c.Close()
		t.Fatal("Invalid client address succeeded")
	}
}

func TestClient_Handshake(t *testing.T) {
	srv := diamtest.NewServer(New(serverSettings), dict.Default)
	defer srv.Close()
	cli := &Client{
		Handler: New(clientSettings),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, clientSettings.VendorID),
		},
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
				},
			}),
		},
	}
	c, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
}

func TestClient_Handshake_CustomIP_TCP(t *testing.T) {
	testClient_Handshake_CustomIP(t, "tcp")
}

func testClient_Handshake_CustomIP(t *testing.T, network string) {
	srv := diamtest.NewServerNetwork(network, New(serverSettings), dict.Default)
	defer srv.Close()
	cli := &Client{
		RetransmitInterval: time.Second * 3,
		Handler:            New(clientSettings2),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, clientSettings.VendorID),
		},
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
				},
			}),
		},
	}
	c, err := cli.DialNetwork(network, srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
}

func TestClient_Handshake_Notify(t *testing.T) {
	srv := diamtest.NewServer(New(serverSettings), dict.Default)
	defer srv.Close()
	cli := &Client{
		Handler: New(clientSettings),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, clientSettings.VendorID),
		},
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
				},
			}),
		},
	}
	handshakeOK := make(chan struct{})
	go func() {
		<-cli.Handler.HandshakeNotify()
		close(handshakeOK)
	}()
	c, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	select {
	case <-handshakeOK:
	case <-time.After(time.Second):
		t.Fatal("Handshake timed out")
	}
}

func TestClient_Handshake_FailParseCEA(t *testing.T) {
	mux := diam.NewServeMux()
	mux.HandleFunc("CER", func(c diam.Conn, m *diam.Message) {
		a := m.Answer(diam.Success)
		// Missing Origin-Host and other mandatory AVPs.
		a.WriteTo(c)
	})
	srv := diamtest.NewServer(mux, dict.Default)
	defer srv.Close()
	cli := &Client{
		Handler: New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
	}
	_, err := cli.Dial(srv.Addr)
	if err != smparser.ErrMissingOriginHost {
		t.Fatal(err)
	}
}

func TestClient_Handshake_FailedResultCode(t *testing.T) {
	mux := diam.NewServeMux()
	mux.HandleFunc("CER", func(c diam.Conn, m *diam.Message) {
		cer := new(smparser.CER)
		if _, err := cer.Parse(m, smparser.Server); err != nil {
			panic(err)
		}
		a := m.Answer(diam.NoCommonApplication)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
		if cer.OriginStateID != nil {
			a.AddAVP(cer.OriginStateID)
		}
		a.AddAVP(cer.AcctApplicationID[0]) // The one we send below.
		a.WriteTo(c)
	})
	srv := diamtest.NewServer(mux, dict.Default)
	defer srv.Close()
	cli := &Client{
		Handler: New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
	}
	_, err := cli.Dial(srv.Addr)
	if err == nil {
		t.Fatal("Unexpected CER worked")
	}
	e, ok := err.(*smparser.ErrFailedResultCode)
	if !ok {
		t.Fatal(err)
	}
	if !strings.Contains(e.Error(), "failed Result-Code AVP") {
		t.Fatal(e.Error())
	}
}

func TestClient_Handshake_RetransmitTimeout(t *testing.T) {
	mux := diam.NewServeMux()
	var retransmits uint32
	mux.HandleFunc("CER", func(c diam.Conn, m *diam.Message) {
		// Do nothing to force timeout.
		atomic.AddUint32(&retransmits, 1)
	})
	srv := diamtest.NewServer(mux, dict.Default)
	defer srv.Close()
	cli := &Client{
		Handler:            New(clientSettings),
		MaxRetransmits:     3,
		RetransmitInterval: time.Millisecond,
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
	}
	_, err := cli.Dial(srv.Addr)
	if err == nil {
		t.Fatal("Unexpected CER worked")
	}
	if err != ErrHandshakeTimeout {
		t.Fatal(err)
	}
	if n := atomic.LoadUint32(&retransmits); n != 4 {
		t.Fatalf("Unexpected # of retransmits. Want 4, have %d", n)
	}
}

func TestClient_Watchdog(t *testing.T) {
	srv := diamtest.NewServer(New(serverSettings), dict.Default)
	defer srv.Close()
	cli := &Client{
		EnableWatchdog:   true,
		WatchdogInterval: 100 * time.Millisecond,
		Handler:          New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
	}
	c, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	resp := make(chan struct{}, 1)
	dwa := handleDWA(cli.Handler, resp)
	cli.Handler.mux.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) {
		dwa(c, m)
	})
	select {
	case <-resp:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timeout waiting for DWA")
	}
}

func TestClient_Watchdog_Timeout(t *testing.T) {
	sm := New(serverSettings)
	var once sync.Once
	sm.mux.HandleIdx(baseDWRIdx, handshakeOK(func(c diam.Conn, m *diam.Message) {
		once.Do(func() { m.Answer(diam.UnableToComply).WriteTo(c) })
	}))
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	cli := &Client{
		MaxRetransmits:     3,
		RetransmitInterval: 50 * time.Millisecond,
		EnableWatchdog:     true,
		WatchdogInterval:   50 * time.Millisecond,
		Handler:            New(clientSettings),
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(3)),
		},
	}
	c, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	select {
	case <-c.(diam.CloseNotifier).CloseNotify():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for watchdog to disconnect client")
	}
}
