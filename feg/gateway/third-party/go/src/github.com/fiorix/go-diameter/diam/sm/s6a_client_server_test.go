// Copyright 2013-2018 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package sm

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

const (
	VENDOR_3GPP          = 10415
	PLMN_ID              = "\x00\xF1\x10"
	TEST_IMSI            = "001010000000001"
	ULR_FLAGS            = 1<<1 | 1<<5
	CONCURENT_CLIENTS    = 128 // Number of clients (go routines) simultaneously using a single diameter connection
	TEST_TIMEOUT_SECONDS = 10
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestS6aClientServerTCP(t *testing.T) {
	testS6aClientServer("tcp", t)
}

var (
	sentAIRs, receivedAIRs,
	sentAIAs, receivedAIAs,
	sentULRs, receivedULRs,
	sentULAs, receivedULAs,
	sentCLRs, receivedCLRs,
	sentCLAs, receivedCLAs uint32
)

func testS6aClientServer(network string, t *testing.T) {

	resetTestStats()
	settings := &Settings{
		OriginHost:       datatype.DiameterIdentity("test.host"),
		OriginRealm:      datatype.DiameterIdentity("test.realm"),
		VendorID:         VENDOR_3GPP,
		ProductName:      "go-diameter-s6a",
		FirmwareRevision: 1,
	}

	results := make(chan error, CONCURENT_CLIENTS*2)

	// Create the state machine (mux) and set its message handlers.
	mux := New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: true},
		testHandleAIR(results, settings))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: true},
		testHandleULR(results, settings))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.CancelLocation, Request: false},
		testHandleCLA(results))

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, testHandleALL(results))

	// Print error reports.
	go testPrintErrors(mux.ErrorReports(), results)

	// Start Server
	go func() {
		results <- nil
		err := diam.ListenAndServeNetwork(network, "127.0.0.1:3868", mux, nil)
		if err != nil {
			results <- err
		}
	}()
	err := <-results
	time.Sleep(time.Millisecond * 10)

	// Initialize Client
	cfg := &Settings{
		OriginHost:       datatype.DiameterIdentity("test.host"),
		OriginRealm:      datatype.DiameterIdentity("test.realm"),
		VendorID:         VENDOR_3GPP,
		ProductName:      "go-diameter-s6a",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses:  []datatype.Address{datatype.Address(net.ParseIP("127.0.0.1"))},
	}

	// Create the state machine (it's a diam.ServeMux) and client.
	cmux := New(cfg)

	cli := &Client{
		Dict:               dict.Default,
		Handler:            cmux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   time.Second * 3,
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
				},
			}),
		},
	}

	cmux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: false},
		testHandleAIA(results, cfg))

	cmux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: false},
		testHandleULA(results))

	cmux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.CancelLocation, Request: true},
		testHandleCLR(results, cfg))

	cmux.HandleFunc("ALL", testHandleALL(results)) // Catch all.

	// Print error reports.
	go testPrintErrors(cmux.ErrorReports(), results)

	c, err := cli.DialNetwork(network, "127.0.0.1:3868")
	if err != nil {
		t.Fatal(err)
	}
	timeOut := time.NewTimer(time.Second * TEST_TIMEOUT_SECONDS)
	go func() {
		<-timeOut.C
		results <- fmt.Errorf("TestClientServer %s Timed Out", network)
	}()
	for i := 0; i < CONCURENT_CLIENTS; i++ {
		go func() {
			time.Sleep(time.Nanosecond * time.Duration(rand.Intn(int(time.Millisecond))))
			err := testSendAIR(c, cfg)
			if err != nil {
				results <- err
			} else {
				atomic.AddUint32(&sentAIRs, 1)
			}
		}()
	}
	for i := 0; i < CONCURENT_CLIENTS; i++ {
		err = <-results
		if err != nil {
			t.Error(err)
			for e := 0; e < len(results); e++ {
				err = <-results
				if err != nil {
					t.Error(err)
				}
			}
			break
		}
	}
	time.Sleep(time.Second)
	logStats(t)
}

func testHandleALL(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		results <- fmt.Errorf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
	}
}

// S6a AI
func testHandleAIR(results chan error, settings *Settings) diam.HandlerFunc {
	type RequestedEUTRANAuthInfo struct {
		NumVectors        datatype.Unsigned32  `avp:"Number-Of-Requested-Vectors"`
		ImmediateResponse datatype.Unsigned32  `avp:"Immediate-Response-Preferred"`
		ResyncInfo        datatype.OctetString `avp:"Re-synchronization-Info"`
	}

	type AIR struct {
		SessionID               datatype.UTF8String       `avp:"Session-Id"`
		OriginHost              datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm             datatype.DiameterIdentity `avp:"Origin-Realm"`
		AuthSessionState        datatype.UTF8String       `avp:"Auth-Session-State"`
		UserName                string                    `avp:"User-Name"`
		VisitedPLMNID           datatype.Unsigned32       `avp:"Visited-PLMN-Id"`
		RequestedEUTRANAuthInfo RequestedEUTRANAuthInfo   `avp:"Requested-EUTRAN-Authentication-Info"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var req AIR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			err = fmt.Errorf("Unmarshal failed: %s", err)
			code = diam.UnableToComply
			results <- err
		} else {
			code = diam.Success
			atomic.AddUint32(&receivedAIRs, 1)
		}

		a := m.Answer(code)
		// SessionID is required to be the AVP in position 1
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID))
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.OriginStateID, avp.Mbit, 0, settings.OriginStateID)
		_, err = testSendAIA(c, a)
		if err != nil {
			results <- fmt.Errorf("Failed to send AIA: %s", err.Error())
		} else {
			atomic.AddUint32(&sentAIAs, 1)
		}
	}
}

func testSendAIA(w io.Writer, m *diam.Message) (n int64, err error) {

	m.NewAVP(avp.AuthenticationInfo, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.EUTRANVector, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z\x19")),
					diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
					diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
					diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
				},
			}),
		},
	})

	return m.WriteTo(w)
}

func testHandleCLR(results chan error, settings *Settings) diam.HandlerFunc {
	type CLR struct {
		SessionID        string                    `avp:"Session-Id"`
		AuthSessionState int32                     `avp:"Auth-Session-State"`
		OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
		CancellationType int32                     `avp:"Cancellation-Type"`
		DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
		DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
		UserName         string                    `avp:"User-Name"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var code uint32
		var clr CLR
		err := m.Unmarshal(&clr)

		if err != nil {
			err = fmt.Errorf("CLR Unmarshal failed: %s", err)
			code = diam.UnableToComply
			results <- err
		} else {
			code = diam.Success
			atomic.AddUint32(&receivedCLRs, 1)
		}

		a := m.Answer(code)
		// SessionID is required to be the AVP in position 1
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(clr.SessionID)))
		a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(clr.AuthSessionState))
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		_, err = testSendCLA(settings, c, a)
		if err != nil {
			err := fmt.Errorf("Failed to send ULA: %s", err.Error())
			results <- err
		} else {
			atomic.AddUint32(&sentCLAs, 1)
		}
	}
}

// S6a UL
func testHandleULR(results chan error, settings *Settings) diam.HandlerFunc {

	type ULR struct {
		SessionID        datatype.UTF8String       `avp:"Session-Id"`
		OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
		AuthSessionState datatype.Unsigned32       `avp:"Auth-Session-State"`
		UserName         datatype.UTF8String       `avp:"User-Name"`
		VisitedPLMNID    datatype.Unsigned32       `avp:"Visited-PLMN-Id"`
		RATType          datatype.Unsigned32       `avp:"RAT-Type"`
		ULRFlags         datatype.Unsigned32       `avp:"ULR-Flags"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var req ULR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			err = fmt.Errorf("Unmarshal failed: %s", err)
			code = diam.UnableToComply
			results <- err
		} else {
			code = diam.Success
			atomic.AddUint32(&receivedULRs, 1)
		}

		a := m.Answer(code)
		// SessionID is required to be the AVP in position 1
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID))
		a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, req.AuthSessionState)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.OriginStateID, avp.Mbit, 0, settings.OriginStateID)
		_, err = testSendULA(settings, c, a)
		if err != nil {
			results <- fmt.Errorf("Failed to send ULA: %s", err.Error())
		} else {
			atomic.AddUint32(&sentULAs, 1)
		}
		// send cancel location request
		err = testSendCLR(c, settings)
		if err != nil {
			results <- fmt.Errorf("Failed to send CLR: %s\n", err.Error())
		} else {
			atomic.AddUint32(&sentCLRs, 1)
		}
	}
}

func testSendCLA(settings *Settings, w io.Writer, m *diam.Message) (n int64, err error) {
	return m.WriteTo(w)
}

func testSendCLR(c diam.Conn, cfg *Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.CancelLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(TEST_IMSI))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.CancellationType, avp.Mbit, 0, datatype.Enumerated(2))
	_, err := m.WriteTo(c)
	return err
}

func testSendULA(settings *Settings, w io.Writer, m *diam.Message) (n int64, err error) {
	m.NewAVP(avp.ULAFlags, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(1))
	m.NewAVP(avp.SubscriptionData, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MSISDN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("12345")),
			diam.NewAVP(avp.AccessRestrictionData, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(47)),
			diam.NewAVP(avp.SubscriberStatus, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
			diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(2)),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(
						avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
					diam.NewAVP(
						avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
				},
			}),
			diam.NewAVP(avp.APNConfigurationProfile, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
					diam.NewAVP(avp.AllAPNConfigurationsIncludedIndicator, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
					diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(2)),
					diam.NewAVP(avp.APNConfiguration, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
							diam.NewAVP(avp.PDNType, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
							diam.NewAVP(avp.ServiceSelection, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.UTF8String("oai.ipv4")),
							diam.NewAVP(avp.EPSSubscribedQoSProfile, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(avp.QoSClassIdentifier, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(9)),
									diam.NewAVP(avp.AllocationRetentionPriority, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
										AVP: []*diam.AVP{
											diam.NewAVP(avp.PriorityLevel, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(15)),
											diam.NewAVP(avp.PreemptionCapability, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(1)),
											diam.NewAVP(avp.PreemptionVulnerability, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(0)),
										},
									}),
								},
							}),
						},
					}),
					diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
							diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
						},
					}),
				},
			}),
		},
	})

	return m.WriteTo(w)
}

// Create & send Authentication-Information Request
func testSendAIR(c diam.Conn, cfg *Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(TEST_IMSI))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), datatype.OctetString(PLMN_ID))
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), datatype.Unsigned32(3)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), datatype.Unsigned32(0)),
		},
	})
	_, err := m.WriteTo(c)
	return err
}

// Create & send Update-Location Request
func testSendULR(c diam.Conn, cfg *Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(TEST_IMSI))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.RATType, avp.Mbit, uint32(VENDOR_3GPP), datatype.Enumerated(1004))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), datatype.Unsigned32(ULR_FLAGS))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, uint32(VENDOR_3GPP), datatype.OctetString(PLMN_ID))
	_, err := m.WriteTo(c)
	return err
}

type EUtranVector struct {
	RAND  datatype.OctetString `avp:"RAND"`
	XRES  datatype.OctetString `avp:"XRES"`
	AUTN  datatype.OctetString `avp:"AUTN"`
	KASME datatype.OctetString `avp:"KASME"`
}

type ExperimentalResult struct {
	VendorId               uint32 `avp:"Vendor-Id"`
	ExperimentalResultCode uint32 `avp:"Experimental-Result-Code"`
}

type AuthenticationInfo struct {
	EUtranVector EUtranVector `avp:"E-UTRAN-Vector"`
}

type AIA struct {
	SessionID          string                    `avp:"Session-Id"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	AIs                []AuthenticationInfo      `avp:"Authentication-Info"`
}

type AMBR struct {
	MaxRequestedBandwidthUL uint32 `avp:"Max-Requested-Bandwidth-UL"`
	MaxRequestedBandwidthDL uint32 `avp:"Max-Requested-Bandwidth-DL"`
}

type AllocationRetentionPriority struct {
	PriorityLevel           uint32 `avp:"Priority-Level"`
	PreemptionCapability    int32  `avp:"Pre-emption-Capability"`
	PreemptionVulnerability int32  `avp:"Pre-emption-Vulnerability"`
}

type EPSSubscribedQoSProfile struct {
	QoSClassIdentifier          int32                       `avp:"QoS-Class-Identifier"`
	AllocationRetentionPriority AllocationRetentionPriority `avp:"Allocation-Retention-Priority"`
}

type APNConfiguration struct {
	ContextIdentifier       uint32                  `avp:"Context-Identifier"`
	PDNType                 int32                   `avp:"PDN-Type"`
	ServiceSelection        string                  `avp:"Service-Selection"`
	EPSSubscribedQoSProfile EPSSubscribedQoSProfile `avp:"EPS-Subscribed-QoS-Profile"`
	AMBR                    AMBR                    `avp:"AMBR"`
}

type APNConfigurationProfile struct {
	ContextIdentifier                     uint32             `avp:"Context-Identifier"`
	AllAPNConfigurationsIncludedIndicator int32              `avp:"All-APN-Configurations-Included-Indicator"`
	APNConfigs                            []APNConfiguration `avp:"APN-Configuration"`
}

type SubscriptionData struct {
	MSISDN                        datatype.OctetString    `avp:"MSISDN"`
	AccessRestrictionData         uint32                  `avp:"Access-Restriction-Data"`
	SubscriberStatus              int32                   `avp:"Subscriber-Status"`
	NetworkAccessMode             int32                   `avp:"Network-Access-Mode"`
	AMBR                          AMBR                    `avp:"AMBR"`
	APNConfigurationProfile       APNConfigurationProfile `avp:"APN-Configuration-Profile"`
	SubscribedPeriodicRauTauTimer uint32                  `avp:"Subscribed-Periodic-RAU-TAU-Timer"`
}

type ULA struct {
	SessionID          string                    `avp:"Session-Id"`
	ULAFlags           uint32                    `avp:"ULA-Flags"`
	SubscriptionData   SubscriptionData          `avp:"Subscription-Data"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
}

type CLA struct {
	SessionID          string                    `avp:"Session-Id"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
}

func testHandleAIA(results chan error, cfg *Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		if m.Header.CommandCode != diam.AuthenticationInformation {
			results <- fmt.Errorf("Unexpected Command Code for AIA: %d", m.Header.CommandCode)
		} else {
			atomic.AddUint32(&receivedAIAs, 1)
			var req AIA
			err := m.Unmarshal(&req) // Make sure, we can unmarshal it
			if err != nil {
				err = fmt.Errorf("AIA Unmarshal failed: %s", err)
				results <- err
				return
			}

			err = testSendULR(c, cfg)
			if err != nil {
				results <- err
			} else {
				atomic.AddUint32(&sentULRs, 1)
			}
		}
	}
}

func testHandleULA(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		if m.Header.CommandCode != diam.UpdateLocation {
			results <- fmt.Errorf("Unexpected Command Code for ULA: %d", m.Header.CommandCode)
		} else {
			atomic.AddUint32(&receivedULAs, 1)
			var req ULA
			err := m.Unmarshal(&req) // Make sure, we can unmarshal it
			if err != nil {
				err = fmt.Errorf("ULA Unmarshal failed: %s", err)
				results <- err
				return
			}
			results <- nil
		}
	}
}

func testHandleCLA(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		if m.Header.CommandCode != diam.CancelLocation {
			results <- fmt.Errorf("Unexpected Command Code for CLA: %d", m.Header.CommandCode)
		} else {
			atomic.AddUint32(&receivedCLAs, 1)
			var cla CLA
			err := m.Unmarshal(&cla) // Make sure, we can unmarshal it
			if err != nil {
				err = fmt.Errorf("CLA Unmarshal failed: %s", err)
				results <- err
				return
			}
			results <- nil
		}
	}
}

func resetTestStats() {
	atomic.StoreUint32(&sentAIRs, 0)
	atomic.StoreUint32(&receivedAIRs, 0)
	atomic.StoreUint32(&sentAIAs, 0)
	atomic.StoreUint32(&receivedAIAs, 0)
	atomic.StoreUint32(&sentULRs, 0)
	atomic.StoreUint32(&receivedULRs, 0)
	atomic.StoreUint32(&sentULAs, 0)
	atomic.StoreUint32(&receivedULAs, 0)
}

func logStats(t *testing.T) {
	t.Logf(
		"AIRs sent/received: %d/%d; AIAs sent/received: %d/%d; ULRs sent/received: %d/%d; ULAs sent/received: %d/%d; "+
			"CLRs sent/received: %d/%d; CLAs sent/received: %d/%d",
		atomic.LoadUint32(&sentAIRs), atomic.LoadUint32(&receivedAIRs),
		atomic.LoadUint32(&sentAIAs), atomic.LoadUint32(&receivedAIAs),
		atomic.LoadUint32(&sentULRs), atomic.LoadUint32(&receivedULRs),
		atomic.LoadUint32(&sentULAs), atomic.LoadUint32(&receivedULAs),
		atomic.LoadUint32(&sentCLRs), atomic.LoadUint32(&receivedCLRs),
		atomic.LoadUint32(&sentCLAs), atomic.LoadUint32(&receivedCLAs),
	)
}

func testPrintErrors(ec <-chan *diam.ErrorReport, results chan error) {
	for err := range ec {
		results <- fmt.Errorf("Error: %v for Message: %s", err.Error, err.Message)
	}
}
