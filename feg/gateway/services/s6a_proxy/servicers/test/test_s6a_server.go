/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/servicers"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
)

const (
	TEST_PLMN_ID = "\x00\xF1\x10"
	TEST_IMSI    = "001010000000001"
	VENDOR_3GPP  = diameter.Vendor3GPP
)

// StartTestS6aServer starts a new Test S6a Server on given network & address
func StartTestS6aServer(network, addr string) error {
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("magma-oai.openair4G.eur"),
		OriginRealm:      datatype.DiameterIdentity("openair4G.eur"),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      "go-diameter-s6a",
		FirmwareRevision: 1,
	}
	// Create the state machine (mux) and set its message handlers.
	results := make(chan error, 2)
	mux := sm.New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: true},
		testHandleAIR(settings))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: true},
		testHandleULR(settings))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.PurgeUE, Request: true},
		testHandlePUR(settings))

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, testHandleALL(results))

	// Print error reports.
	go testPrintErrors(mux.ErrorReports())

	// Start S6a Diameter Server
	go func() {
		results <- nil
		err := diam.ListenAndServeNetwork(network, addr, mux, nil)
		if err != nil {
			fmt.Printf("StartTestS6aServer Error: %v for address: %s\n", err, addr)
			results <- err
		}
	}()
	err := <-results
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)
	return nil
}

func testHandleALL(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		results <- fmt.Errorf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
	}
}

// S6a AI
func testHandleAIR(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req servicers.AIR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("AIR Unmarshal for message: %s failed: %s", m, err)
			code = diam.UnableToComply
		} else {
			code = diam.Success
		}

		a := m.Answer(code)
		// SessionID is required to be the AVP in position 1
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID))
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.OriginStateID, avp.Mbit, 0, settings.OriginStateID)
		_, err = testSendAIA(c, a, int(req.RequestedEUTRANAuthInfo.NumVectors))
		if err != nil {
			fmt.Printf("Failed to send AIA: %s", err.Error())
		}
	}
}

func testSendAIA(w io.Writer, m *diam.Message, vectors int) (n int64, err error) {
	if vectors < 0 {
		vectors = 1
	}
	if vectors > 5 {
		vectors = 5
	}
	for ; vectors > 0; vectors-- {
		m.NewAVP(avp.AuthenticationInfo, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.EUTRANVector, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z"+strconv.Itoa(14+vectors))),
						diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
						diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
						diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
					},
				}),
			},
		})
	}
	return m.WriteTo(w)
}

// S6a UL
func testHandleULR(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req servicers.ULR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("ULR Unmarshal for message: %s failed: %s", m, err)
			code = diam.UnableToComply
		} else {
			code = diam.Success
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
			fmt.Printf("Failed to send ULA: %s", err.Error())
		}
	}
}

func testSendULA(settings *sm.Settings, w io.Writer, m *diam.Message) (n int64, err error) {
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

func testHandlePUR(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req servicers.PUR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("PUR Unmarshal for message: %s failed: %s", m, err)
			code = diam.UnableToComply
		} else {
			code = diam.Success
		}

		a := m.Answer(code)
		// SessionID is required to be the AVP in position 1
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(req.SessionID)))
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)

		_, err = testSendPUA(c, a)
		if err != nil {
			fmt.Printf("Failed to send PUA: %s", err.Error())
		}
	}
}

func testSendPUA(w io.Writer, m *diam.Message) (n int64, err error) {
	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
	})
	return m.WriteTo(w)
}

func testPrintErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		fmt.Printf("Error: %v for Message: %s", err.Error, err.Message)
	}
}
