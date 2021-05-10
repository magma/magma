/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"magma/feg/gateway/diameter"
	swx "magma/feg/gateway/services/swx_proxy/servicers"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
)

const (
	BASE_IMSI   = "0010100000"
	VENDOR_3GPP = diameter.Vendor3GPP

	DefaultSIPAuthenticate  = "\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z"
	DefaultSIPAuthorization = "F\xf0\"\xb9%#\xf58"
	DefaultCK               = "\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_"
	DefaultIK               = "\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:"
)

// StartTestSwxServer starts a new Test Swx Server on given network & address
func StartTestSwxServer(network, addr string) (string, error) {
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("magma-oai.openair4G.eur"),
		OriginRealm:      datatype.DiameterIdentity("openair4G.eur"),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      "go-diameter-swx",
		FirmwareRevision: 1,
	}
	// Create the state machine (mux) and set its message handlers.
	errResults := make(chan error, 2)
	addrResult := make(chan string, 1)

	mux := sm.New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.MultimediaAuthentication, Request: true},
		testHandleMAR(settings))
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.ServerAssignment, Request: true},
		testHandleSAR(settings))

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, testHandleALL(errResults))

	// Print error reports.
	go testPrintErrors(mux.ErrorReports())

	// Start Swx Diameter Server
	go func() {
		errResults <- nil
		server := diam.Server{
			Network: network,
			Addr:    addr,
			Handler: mux,
		}
		lis, err := diam.MultistreamListen(network, addr)
		if err != nil {
			fmt.Printf("StartTestSwxServer Error: %v for address: %s\n", err, addr)
			errResults <- err
		}
		addrResult <- lis.Addr().String()
		server.Serve(lis)
	}()
	err := <-errResults
	serverAddr := <-addrResult
	if err != nil {
		return "", err
	}
	time.Sleep(time.Millisecond * 20)
	return serverAddr, nil
}

func testHandleALL(results chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		results <- fmt.Errorf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
	}
}

// Swx Multimedia Authentication Request
func testHandleMAR(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req swx.MAR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("MAR Unmarshal for message: %s failed: %s", m, err)
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
		_, err = testSendMAA(c, a, int(req.NumberAuthItems))
		if err != nil {
			fmt.Printf("Failed to send MAA: %s", err.Error())
		}
	}
}

// SWx Server Assignment Request
func testHandleSAR(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req swx.SAR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("SAR Unmarshal for message: %s failed: %s", m, err)
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
		a.NewAVP(avp.UserName, avp.Mbit, 0, req.UserName)
		_, err = testSendSAA(c, a)
		if err != nil {
			fmt.Printf("Failed to send SAA: %s", err.Error())
		}
	}
}

// Send Multimedia Authentication Answer
func testSendMAA(w io.Writer, m *diam.Message, vectors int) (n int64, err error) {
	m.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(vectors))
	for i := 0; i < vectors; i++ {
		m.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.UTF8String("EAP-AKA")),
				diam.NewAVP(avp.SIPAuthenticate, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(DefaultSIPAuthenticate+strconv.Itoa(14+i))),
				diam.NewAVP(avp.SIPAuthorization, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(DefaultSIPAuthorization)),
				diam.NewAVP(avp.ConfidentialityKey, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(DefaultCK)),
				diam.NewAVP(avp.IntegrityKey, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(DefaultIK)),
			},
		})
	}
	return m.WriteTo(w)
}

// Send Server Assignment Answer
func testSendSAA(w io.Writer, m *diam.Message) (n int64, err error) {
	m.NewAVP(avp.Non3GPPUserData, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)), // END_USER_E164
					diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.OctetString("12345")),
				},
			}),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(
						avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
					diam.NewAVP(
						avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
				},
			}),
			diam.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Enumerated(2002)),
			diam.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Enumerated(2003)),
			diam.NewAVP(avp.Non3GPPIPAccess, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Enumerated(swx.Non3GPPIPAccess_ENABLED)),
			diam.NewAVP(avp.Non3GPPIPAccessAPN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Enumerated(0)),
		},
	})
	return m.WriteTo(w)
}

func testPrintErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		fmt.Printf("Error: %v for Message: %s", err.Error, err.Message)
	}
}

// StartEmptyDiameterServer starts an empty server for testing
func StartEmptyDiameterServer(network, addr string) (string, error) {
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("magma-oai.openair4G.eur"),
		OriginRealm:      datatype.DiameterIdentity("openair4G.eur"),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      "go-diameter-swx",
		FirmwareRevision: 1,
	}
	// Create the state machine (mux) and set its message handlers.
	errResults := make(chan error, 2)
	addrResult := make(chan string, 1)

	mux := sm.New(settings)

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, testHandleALL(errResults))

	// Print error reports.
	go testPrintErrors(mux.ErrorReports())

	// Start Swx Diameter Server
	go func() {
		errResults <- nil
		server := diam.Server{
			Network: network,
			Addr:    addr,
			Handler: mux,
		}
		lis, err := diam.MultistreamListen(network, addr)
		if err != nil {
			fmt.Printf("StartEmptyDiameterServer Error: %v for address: %s\n", err, addr)
			errResults <- err
		}
		addrResult <- lis.Addr().String()
		server.Serve(lis)
	}()
	err := <-errResults
	serverAddr := <-addrResult
	if err != nil {
		return "", err
	}
	time.Sleep(time.Millisecond * 20)
	return serverAddr, nil
}
