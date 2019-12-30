// +build go1.8
// +build linux,!386

package test

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/examples/s6a_proxy/service"
)

const (
	TEST_PLMN_ID = "\x00\xF1\x10"
	TEST_IMSI    = "001010000000001"
)

// StartTestS6aServer starts a new Test S6a Server on given network & address
func StartTestS6aServer(network, addr string) error {
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("magma-oai.openair4G.eur"),
		OriginRealm:      datatype.DiameterIdentity("openair4G.eur"),
		VendorID:         datatype.Unsigned32(service.VENDOR_3GPP),
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

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, testHandleALL(results))

	// Print error reports.
	go testPrintErrors(mux.ErrorReports())

	// Start S6a Diameter Server
	go func() {
		results <- nil
		err := diam.ListenAndServeNetwork(network, addr, mux, nil)
		if err != nil {
			fmt.Printf("StartTestS6aServer Error: %v\n", err)
			results <- err
		}
	}()
	err := <-results
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 10)
	return nil
}

func testHandleALL(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		results <- fmt.Errorf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
	}
}

// S6a AI
func testHandleAIR(settings *sm.Settings) diam.HandlerFunc {
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
			fmt.Printf("Unmarshal failed: %s", err)
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

func testSendAIA(w io.Writer, m *diam.Message, vectors int) (int64, error) {
	if vectors < 0 {
		vectors = 1
	}
	if vectors > 5 {
		vectors = 5
	}
	for ; vectors > 0; vectors-- {
		m.NewAVP(avp.AuthenticationInfo, avp.Mbit, service.VENDOR_3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.EUTRANVector, avp.Mbit, service.VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z"+strconv.Itoa(14+vectors))),
						diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
						diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
						diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
					},
				}),
			},
		})
	}
	n, err := m.WriteToStream(w.(diam.MultistreamWriter), m.MessageStream())
	return int64(n), err
}

// S6a UL
func testHandleULR(settings *sm.Settings) diam.HandlerFunc {

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
			fmt.Printf("Unmarshal failed: %s", err)
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

func testSendULA(settings *sm.Settings, w io.Writer, m *diam.Message) (int64, error) {
	m.NewAVP(avp.ULAFlags, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(1))
	m.NewAVP(avp.SubscriptionData, avp.Mbit, service.VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MSISDN, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.OctetString("12345")),
			diam.NewAVP(avp.AccessRestrictionData, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(47)),
			diam.NewAVP(avp.SubscriberStatus, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
			diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(2)),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(
						avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(500)),
					diam.NewAVP(
						avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(500)),
				},
			}),
			diam.NewAVP(avp.APNConfigurationProfile, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
					diam.NewAVP(avp.AllAPNConfigurationsIncludedIndicator, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
					diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(2)),
					diam.NewAVP(avp.APNConfiguration, avp.Mbit, service.VENDOR_3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
							diam.NewAVP(avp.PDNType, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
							diam.NewAVP(avp.ServiceSelection, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.UTF8String("oai.ipv4")),
							diam.NewAVP(avp.EPSSubscribedQoSProfile, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, &diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(avp.QoSClassIdentifier, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(9)),
									diam.NewAVP(avp.AllocationRetentionPriority, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, &diam.GroupedAVP{
										AVP: []*diam.AVP{
											diam.NewAVP(avp.PriorityLevel, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(15)),
											diam.NewAVP(avp.PreemptionCapability, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(1)),
											diam.NewAVP(avp.PreemptionVulnerability, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(0)),
										},
									}),
								},
							}),
						},
					}),
					diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(500)),
							diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, service.VENDOR_3GPP, datatype.Unsigned32(500)),
						},
					}),
				},
			}),
		},
	})

	n, err := m.WriteToStream(w.(diam.MultistreamWriter), m.MessageStream())
	return int64(n), err
}

func testPrintErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		fmt.Printf("Error: %v for Message: %s", err.Error, err.Message)
	}
}
