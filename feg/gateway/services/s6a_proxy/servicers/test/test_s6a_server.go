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
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/servicers"
)

const (
	TEST_PLMN_ID = "\x00\xF1\x10"
	TEST_IMSI    = "001010000000001"
	TEST_IMSI_2  = "001030000000001"
	VENDOR_3GPP  = diameter.Vendor3GPP
)

// StartTestS6aServer starts a new Test S6a Server on given network & address
func StartTestS6aServer(network, addr string, useStaticResp bool) error {
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

	if useStaticResp {
		mux.HandleIdx(
			diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: true},
			testHandleAIRStatic(settings))
	} else {
		mux.HandleIdx(
			diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: true},
			testHandleAIR(settings))
	}

	// expected ULRFlags = 290 (100100010) where the fifth 1 is DualRegistration_5GIndicator : true
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: true},
		testHandleULR(settings, 290))

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
		_, err = testSendAIA(
			c, a, int(req.RequestedEUTRANAuthInfo.NumVectors), int(req.RequestedUtranGeranAuthInfo.NumVectors))
		if err != nil {
			fmt.Printf("Failed to send AIA: %s", err.Error())
		}
	}
}

func testSendAIA(w io.Writer, m *diam.Message, eutranVectors, utranVectors int) (n int64, err error) {
	if eutranVectors > 5 {
		eutranVectors = 5
	}
	vectorAvps := []*diam.AVP{}
	for i := 0; i < eutranVectors; i++ {
		vectorAvps = append(vectorAvps,
			diam.NewAVP(avp.EUTRANVector, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.ItemNumber, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(i)),
					diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z"+strconv.Itoa(14+i))),
					diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
					diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
					diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
				},
			}))
	}
	if utranVectors > 0 {
		for i := 0; i < utranVectors; i++ {
			vectorAvps = append(vectorAvps,
				diam.NewAVP(avp.UTRANVector, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.ItemNumber, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(i)),
						diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(
							[]byte{57, 22, 40, 33, 82, 189, 193, 89, 219, 31, 18, 64, 95, 197, 50, 240 + byte(i)})),
						diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(
							[]byte{155, 36, 12, 2, 227, 43, 246, 254})),
						diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(
							[]byte{188, 167, 68, 25, 19, 11, 128, 0, 228, 20, 201, 246, 253, 57, 224, 99})),
						diam.NewAVP(avp.ConfidentialityKey, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(
							[]byte{235, 74, 254, 58, 73, 108, 112, 173, 61, 24, 169, 176, 219, 233, 85, 180})),
						diam.NewAVP(avp.IntegrityKey, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString(
							[]byte{8, 114, 43, 29, 82, 150, 220, 38, 242, 123, 82, 108, 116, 174, 27, 212})),
					},
				}))
		}
	}
	m.NewAVP(avp.AuthenticationInfo, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{AVP: vectorAvps})
	return m.WriteTo(w)
}

// Test AIA message
// Version: 0x01
// Length: 772
// Flags: 0x40, Proxyable
// Command Code: 318 3GPP-Authentication-Information
// ApplicationId: 3GPP S6a/S6d (16777251)
// Hop-by-Hop Identifier: 0x38cf1263
// End-to-End Identifier: 0x2c31ced8
// [Request In: 1]
// [Response Time: 0.000000000 seconds]
// AVP: Session-Id(263) l=33 f=-M- val=magma;1536276998_947F557F
// AVP: Vendor-Specific-Application-Id(260) l=32 f=-M-
// AVP: Auth-Session-State(277) l=12 f=-M- val=NO_STATE_MAINTAINED (1)
// AVP: Origin-Host(264) l=48 f=-M- val=hw-hss.epc.mnc001.mcc001.3gppnetwork.org
// AVP: Origin-Realm(296) l=41 f=-M- val=epc.mnc001.mcc001.3gppnetwork.org
// AVP: Result-Code(268) l=12 f=-M- val=DIAMETER_SUCCESS (2001)
// AVP: Authentication-Info(1413) l=456 f=VM- vnd=TGPP
// AVP Code: 1413 Authentication-Info
// AVP Flags: 0xc0, Vendor-Specific: Set, Mandatory: Set
// AVP Length: 456
// AVP Vendor Id: 3GPP (10415)
// Authentication-Info: 00000586c0000094000028af0000058bc0000010000028af...
// AVP: E-UTRAN-Vector(1414) l=148 f=VM- vnd=TGPP
// AVP: E-UTRAN-Vector(1414) l=148 f=VM- vnd=TGPP
// AVP: E-UTRAN-Vector(1414) l=148 f=VM- vnd=TGPP
// AVP: Supported-Features(628) l=56 f=V-- vnd=TGPP
// AVP: Supported-Features(628) l=56 f=V-- vnd=TGPP
var staticAIA = []byte("\x01\x00\x03\x04\x40\x00\x01\x3e\x01\x00\x00\x23\x38\xcf\x12\x63" +
	"\x2c\x31\xce\xd8\x00\x00\x01\x07\x40\x00\x00\x21\x6d\x61\x67\x6d" +
	"\x61\x3b\x31\x35\x33\x36\x32\x37\x36\x39\x39\x38\x5f\x39\x34\x37" +
	"\x46\x35\x35\x37\x46\x00\x00\x00\x00\x00\x01\x04\x40\x00\x00\x20" +
	"\x00\x00\x01\x0a\x40\x00\x00\x0c\x00\x00\x28\xaf\x00\x00\x01\x02" +
	"\x40\x00\x00\x0c\x01\x00\x00\x23\x00\x00\x01\x15\x40\x00\x00\x0c" +
	"\x00\x00\x00\x01\x00\x00\x01\x08\x40\x00\x00\x30\x68\x77\x2d\x68" +
	"\x73\x73\x2e\x65\x70\x63\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63" +
	"\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b" +
	"\x2e\x6f\x72\x67\x00\x00\x01\x28\x40\x00\x00\x29\x65\x70\x63\x2e" +
	"\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67" +
	"\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00\x00\x00" +
	"\x00\x00\x01\x0c\x40\x00\x00\x0c\x00\x00\x07\xd1\x00\x00\x05\x85" +
	"\xc0\x00\x01\xc8\x00\x00\x28\xaf\x00\x00\x05\x86\xc0\x00\x00\x94" +
	"\x00\x00\x28\xaf\x00\x00\x05\x8b\xc0\x00\x00\x10\x00\x00\x28\xaf" +
	"\x00\x00\x00\x00\x00\x00\x05\xa7\xc0\x00\x00\x1c\x00\x00\x28\xaf" +
	"\x15\x9a\xbf\x21\xca\xe2\xbf\x0a\xdb\xcb\xf1\x47\xef\x87\x74\x9d" +
	"\x00\x00\x05\xa8\xc0\x00\x00\x14\x00\x00\x28\xaf\x66\x6a\x9b\x73" +
	"\x6d\x41\x7a\x41\x00\x00\x05\xa9\xc0\x00\x00\x1c\x00\x00\x28\xaf" +
	"\x6c\x3c\x22\x75\x44\x1a\x80\x00\x7e\x59\x79\x64\x47\x61\x43\x9a" +
	"\x00\x00\x05\xaa\xc0\x00\x00\x2c\x00\x00\x28\xaf\x3b\x6d\xd0\x9b" +
	"\x14\x6b\xfb\x53\x05\x1a\xa7\x5d\xe5\xe8\x93\xdb\x43\xd1\x00\xe4" +
	"\x10\x48\x7d\x75\xcb\x26\x99\xe0\xe7\x80\xbb\x9e\x00\x00\x05\x86" +
	"\xc0\x00\x00\x94\x00\x00\x28\xaf\x00\x00\x05\x8b\xc0\x00\x00\x10" +
	"\x00\x00\x28\xaf\x00\x00\x00\x01\x00\x00\x05\xa7\xc0\x00\x00\x1c" +
	"\x00\x00\x28\xaf\xd0\xb3\x82\xe2\xec\x53\xe3\xa6\xaf\xd8\x1c\xa4" +
	"\x57\x92\xd8\xa6\x00\x00\x05\xa8\xc0\x00\x00\x14\x00\x00\x28\xaf" +
	"\xf9\xf8\xcd\xcb\xc6\x50\x59\x47\x00\x00\x05\xa9\xc0\x00\x00\x1c" +
	"\x00\x00\x28\xaf\x63\x82\xb8\x54\x48\x59\x80\x00\xf5\xaf\x37\xa5" +
	"\xe9\x6d\x76\x58\x00\x00\x05\xaa\xc0\x00\x00\x2c\x00\x00\x28\xaf" +
	"\x6d\xb8\x62\xd0\xd8\x54\x79\x51\x07\xc1\xbb\x97\xad\x14\x5d\xc0" +
	"\x68\xd7\x6e\xa5\xbe\x1f\x7d\x80\x1a\x96\xe1\x48\xf5\x04\x89\x96" +
	"\x00\x00\x05\x86\xc0\x00\x00\x94\x00\x00\x28\xaf\x00\x00\x05\x8b" +
	"\xc0\x00\x00\x10\x00\x00\x28\xaf\x00\x00\x00\x02\x00\x00\x05\xa7" +
	"\xc0\x00\x00\x1c\x00\x00\x28\xaf\x97\xce\x80\x32\x83\x47\x13\xdd" +
	"\x31\xf7\x06\xb3\x77\xd6\x1a\x92\x00\x00\x05\xa8\xc0\x00\x00\x14" +
	"\x00\x00\x28\xaf\x37\x2d\x89\xbc\x7c\x6b\x31\x9c\x00\x00\x05\xa9" +
	"\xc0\x00\x00\x1c\x00\x00\x28\xaf\x8a\x9c\x9d\x2b\xdf\xf7\x80\x00" +
	"\x77\xc5\xa4\x1a\xd3\xbc\x21\x29\x00\x00\x05\xaa\xc0\x00\x00\x2c" +
	"\x00\x00\x28\xaf\x74\x60\x79\x2b\x8d\x5e\xb1\x62\xfd\x88\x28\xc2" +
	"\x1a\x3b\xa0\xc5\x6e\x06\xed\xbf\x5b\x20\x54\x72\x50\x06\x36\xc5" +
	"\xfa\xd9\x0b\x84\x00\x00\x02\x74\x80\x00\x00\x38\x00\x00\x28\xaf" +
	"\x00\x00\x01\x0a\x40\x00\x00\x0c\x00\x00\x28\xaf\x00\x00\x02\x75" +
	"\x80\x00\x00\x10\x00\x00\x28\xaf\x00\x00\x00\x01\x00\x00\x02\x76" +
	"\x80\x00\x00\x10\x00\x00\x28\xaf\xbf\xff\xff\xff\x00\x00\x02\x74" +
	"\x80\x00\x00\x38\x00\x00\x28\xaf\x00\x00\x01\x0a\x40\x00\x00\x0c" +
	"\x00\x00\x28\xaf\x00\x00\x02\x75\x80\x00\x00\x10\x00\x00\x28\xaf" +
	"\x00\x00\x00\x02\x00\x00\x02\x76\x80\x00\x00\x10\x00\x00\x28\xaf" +
	"\x00\x01\x00\x00")

func testHandleAIRStatic(settings *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req servicers.AIR

		a, err := diam.ReadMessage(bytes.NewReader(staticAIA), dict.Default)
		if err != nil {
			fmt.Printf("Failed to read static AIA: %v", err)
			panic(err)
		}
		err = m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("AIR Unmarshal for message: %s failed: %s\n", m, err)
		} else {
			a.Header.ApplicationID = m.Header.ApplicationID
			a.Header.EndToEndID = m.Header.EndToEndID
			a.Header.HopByHopID = m.Header.HopByHopID
			sidAVP, err := a.FindAVP(avp.SessionID, 0)
			if err != nil {
				fmt.Printf("SessionID is not found in AIA: %s\n", a)
				a.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID)
			} else if sidAVP != nil {
				fmt.Printf("Setting AIA SessionID to: %s\n", req.SessionID)
				a.Header.MessageLength -= uint32(sidAVP.Len())
				sidAVP.Data = req.SessionID
				a.Header.MessageLength += uint32(sidAVP.Len())
			}
		}
		aavp, _ := a.FindAVP(avp.OriginHost, 0)
		a.Header.MessageLength -= uint32(aavp.Len())
		aavp.Data = settings.OriginHost
		a.Header.MessageLength += uint32(aavp.Len())
		aavp, _ = a.FindAVP(avp.OriginRealm, 0)
		a.Header.MessageLength -= uint32(aavp.Len())
		aavp.Data = settings.OriginRealm
		a.Header.MessageLength += uint32(aavp.Len())
		_, err = a.WriteTo(c)
		if err != nil {
			fmt.Printf("AIA Send Error: %v\n", err)
		}
	}
}

// S6a UL
func testHandleULR(settings *sm.Settings, expectedULRFlags uint32) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var req servicers.ULR
		var code uint32

		err := m.Unmarshal(&req)
		if err != nil {
			fmt.Printf("ULR Unmarshal for message: %s failed: %s", m, err)
			code = diam.UnableToComply
		} else if uint32(req.ULRFlags) != expectedULRFlags {
			// Flags needs to exist for this test
			fmt.Printf("error: ULRFlags (%d) doesnt match with the expected ULRFlags (%d)\n", req.ULRFlags, expectedULRFlags)
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

		// Add Feature-List-ID if exists
		if len(req.SupportedFeatures) > 0 {
			for _, suportedFeature := range req.SupportedFeatures {
				if suportedFeature.FeatureListID == 1 {
					a.NewAVP(avp.SupportedFeatures, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
							diam.NewAVP(avp.FeatureListID, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1)),
							diam.NewAVP(avp.FeatureList, avp.Vbit, diameter.Vendor3GPP,
								datatype.Unsigned32(suportedFeature.FeatureList)),
						},
					})
				}
				if suportedFeature.FeatureListID == 2 {
					a.NewAVP(avp.SupportedFeatures, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
							diam.NewAVP(avp.FeatureListID, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(2)),
							diam.NewAVP(avp.FeatureList, avp.Vbit, diameter.Vendor3GPP,
								datatype.Unsigned32(suportedFeature.FeatureList)),
						},
					})
				}
			}
		}

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
			diam.NewAVP(avp.RegionalSubscriptionZoneCode, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString([]byte{155, 36, 12, 2, 227, 43, 246, 254})),
			diam.NewAVP(avp.RegionalSubscriptionZoneCode, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString([]byte{1, 1, 0, 1})),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(
						avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(math.MaxUint32)),
					diam.NewAVP(
						avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(math.MaxUint32)),
					diam.NewAVP(
						avp.ExtendedMaxRequestedBWDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(500)),
					diam.NewAVP(
						avp.ExtendedMaxRequestedBWUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(600)),
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
							diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, VENDOR_3GPP, &diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(50)),
									diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(60)),
								},
							}),
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
