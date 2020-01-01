// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// testMessage is used by the test cases below and also in reflect_test.go.
// The same testMessage is re-created programmatically in TestNewMessage.
//
// Capabilities-Exchange-Request (CER)
// {Code:257,Flags:0x80,Version:0x1,Length:204,ApplicationId:0,HopByHopId:0xa8cc407d,EndToEndId:0xa8c1b2b4}
//   Origin-Host {Code:264,Flags:0x40,Length:12,VendorId:0,Value:DiameterIdentity{test},Padding:0}
//   Origin-Realm {Code:296,Flags:0x40,Length:20,VendorId:0,Value:DiameterIdentity{localhost},Padding:3}
//   Host-IP-Address {Code:257,Flags:0x40,Length:16,VendorId:0,Value:Address{10.1.0.1},Padding:2}
//   Vendor-Id {Code:266,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{13}}
//   Product-Name {Code:269,Flags:0x0,Length:20,VendorId:0,Value:UTF8String{go-diameter},Padding:1}
//   Origin-State-Id {Code:278,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{1397760650}}
//   Supported-Vendor-Id {Code:265,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{10415}}
//   Supported-Vendor-Id {Code:265,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{13}}
//   Auth-Application-Id {Code:258,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{4}}
//   Inband-Security-Id {Code:299,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{0}}
//   Vendor-Specific-Application-Id {Code:260,Flags:0x40,Length:32,VendorId:0,Value:Grouped{
//     Auth-Application-Id {Code:258,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{4}},
//     Vendor-Id {Code:266,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{10415}},
//   }}
//   Firmware-Revision {Code:267,Flags:0x0,Length:12,VendorId:0,Value:Unsigned32{1}}
var testMessage = []byte{
	0x01, 0x00, 0x00, 0xcc,
	0x80, 0x00, 0x01, 0x01,
	0x00, 0x00, 0x00, 0x00,
	0xa8, 0xcc, 0x40, 0x7d,
	0xa8, 0xc1, 0xb2, 0xb4,
	0x00, 0x00, 0x01, 0x08,
	0x40, 0x00, 0x00, 0x0c,
	0x74, 0x65, 0x73, 0x74,
	0x00, 0x00, 0x01, 0x28,
	0x40, 0x00, 0x00, 0x11,
	0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x68, 0x6f, 0x73,
	0x74, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x01,
	0x40, 0x00, 0x00, 0x0e,
	0x00, 0x01, 0x0a, 0x01,
	0x00, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x0a,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x0d,
	0x00, 0x00, 0x01, 0x0d,
	0x00, 0x00, 0x00, 0x13,
	0x67, 0x6f, 0x2d, 0x64,
	0x69, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x00,
	0x00, 0x00, 0x01, 0x16,
	0x40, 0x00, 0x00, 0x0c,
	0x53, 0x50, 0x22, 0x8a,
	0x00, 0x00, 0x01, 0x09,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x28, 0xaf,
	0x00, 0x00, 0x01, 0x09,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x0d,
	0x00, 0x00, 0x01, 0x02,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x04,
	0x00, 0x00, 0x01, 0x2b,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x04,
	0x40, 0x00, 0x00, 0x20,
	0x00, 0x00, 0x01, 0x02,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x04,
	0x00, 0x00, 0x01, 0x0a,
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x28, 0xaf,
	0x00, 0x00, 0x01, 0x0b,
	0x00, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x01,
}

// CCR
//	Diameter Protocol
//		Version: 0x01
//		Length: 264
//		Flags: 0x80, Request
//		Command Code: 272 Credit-Control
//		ApplicationId: Diameter Credit Control Application (4)
//		Hop-by-Hop Identifier: 0x1df4b75e
//		End-to-End Identifier: 0xb1b50a4b
//	AVP: Session-Id(263) l=14 f=-M- val=123456
//	AVP: Origin-Host(264) l=14 f=-M- val=client
//	AVP: Origin-Realm(296) l=19 f=-M- val=go-diameter
//	AVP: Destination-Realm(283) l=13 f=-M- val=realm
//	AVP: Auth-Application-Id(258) l=12 f=-M- val=Diameter Credit Control Application (4)
//	AVP: Service-Context-Id(461) l=32 f=-M- val=302.220.8.32251@3gpp.org
//	AVP: CC-Request-Type(416) l=12 f=-M- val=INITIAL_REQUEST (1)
//	AVP: User-Name(1) l=23 f=-M- val=123456@test.com
//	AVP: Origin-State-Id(278) l=12 f=-M- val=1437497267
//	AVP: Event-Timestamp(55) l=12 f=-M- val=Dec  9, 2015 15:40:53.000000000 UTC
//	AVP: CC-Request-Number(415) l=12 f=-M- val=0
//	AVP: Service-Information(873) l=60 f=VM- vnd=TGPP
//		AVP: PS-Information(874) l=48 f=VM- vnd=TGPP
//			AVP: GGSN-Address(847) l=18 f=V-- vnd=TGPP val=127.0.0.1 (127.0.0.1)
//			AVP: 3GPP-RAT-Type(21) l=14 f=V-- vnd=TGPP val=3031
var testMessageWithVendorID = []byte{
	0x01, 0x00, 0x01, 0x08, 0x80, 0x00, 0x01, 0x10, 0x00, 0x00, 0x00, 0x04, 0x1d, 0xf4,
	0xb7, 0x5e, 0xb1, 0xb5, 0x0a, 0x4b, 0x00, 0x00, 0x01, 0x07, 0x40, 0x00, 0x00, 0x0e, 0x31, 0x32,
	0x33, 0x34, 0x35, 0x36, 0x00, 0x00, 0x00, 0x00, 0x01, 0x08, 0x40, 0x00, 0x00, 0x0e, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x6d, 0x65, 0x00, 0x00, 0x01, 0x28, 0x40, 0x00, 0x00, 0x13, 0x67, 0x6f,
	0x2d, 0x64, 0x69, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x00, 0x00, 0x00, 0x01, 0x1b, 0x40, 0x00,
	0x00, 0x0d, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x00, 0x01, 0x0d, 0x00, 0x00, 0x01, 0x02, 0x40, 0x00,
	0x00, 0x0c, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x01, 0xcd, 0x40, 0x00, 0x00, 0x20, 0x33, 0x30,
	0x32, 0x2e, 0x32, 0x32, 0x30, 0x2e, 0x38, 0x2e, 0x33, 0x32, 0x32, 0x35, 0x31, 0x40, 0x33, 0x67,
	0x70, 0x70, 0x2e, 0x6f, 0x72, 0x67, 0x00, 0x00, 0x01, 0xa0, 0x40, 0x00, 0x00, 0x0c, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x40, 0x00, 0x00, 0x17, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
	0x40, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x00, 0x01, 0x16, 0x40, 0x00,
	0x00, 0x0c, 0x55, 0xae, 0x77, 0xb3, 0x00, 0x00, 0x00, 0x37, 0x40, 0x00, 0x00, 0x0c, 0xda, 0x12,
	0xca, 0x05, 0x00, 0x00, 0x01, 0x9f, 0x40, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x03, 0x69, 0xc0, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x28, 0xaf, 0x00, 0x00, 0x03, 0x6a, 0xc0, 0x00,
	0x00, 0x30, 0x00, 0x00, 0x28, 0xaf, 0x00, 0x00, 0x03, 0x4f, 0x80, 0x00, 0x00, 0x12, 0x00, 0x00,
	0x28, 0xaf, 0x00, 0x01, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0x80, 0x00,
	0x00, 0x0e, 0x00, 0x00, 0x28, 0xaf, 0x30, 0x31, 0x00, 0x00,
}

func TestReadMessage(t *testing.T) {
	msg, err := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Message:\n%s", msg)
}

func TestReadMessageWithVendorID(t *testing.T) {
	msg, err := ReadMessage(bytes.NewReader(testMessageWithVendorID), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Message:\n%s", msg)
}

func TestNewMessage(t *testing.T) {
	want, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	m := NewMessage(CapabilitiesExchange, RequestFlag, 0, 0xa8cc407d, 0xa8c1b2b4, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("10.1.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	m.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("go-diameter"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1397760650))
	m.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(10415))
	m.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
			NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
		},
	})
	m.NewAVP(avp.FirmwareRevision, 0, 0, datatype.Unsigned32(1))
	if m.Len() != want.Len() {
		t.Fatalf("Unexpected message length.\nWant: %d\n%s\nHave: %d\n%s",
			want.Len(), want, m.Len(), m)
	}
	a, err := m.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := want.Serialize()
	if !bytes.Equal(a, b) {
		t.Fatalf("Unexpected message.\nWant:\n%s\n%s\nHave:\n%s\n%s",
			want, hex.Dump(b), m, hex.Dump(a))
	}
	t.Logf("%d bytes\n%s", len(a), m)
	t.Logf("Message:\n%s", hex.Dump(a))
}

func TestMessageFindAVP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	a, err := m.FindAVP(avp.OriginStateID, 0)
	if err != nil {
		t.Fatal(err)
	}
	a, err = m.FindAVP("Origin-State-Id", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)

	a, err = m.FindAVP("Vendor-Id", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)
	var avps []*AVP
	avps, err = m.FindAVPs("Supported-Vendor-Id", 0)
	if err != nil || len(avps) != 2 {
		t.Fatal(err)
	}
	t.Log(avps)
}

func TestMessageFindAVPsWithPath(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	if avps, err := m.FindAVPsWithPath(nil, 0); err != nil || len(avps) != len(m.AVP) {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
	if avps, err := m.FindAVPsWithPath([]interface{}{avp.VendorID}, 0); len(avps) != 1 {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
	if avps, err := m.FindAVPsWithPath([]interface{}{avp.VendorSpecificApplicationID}, 0); len(avps) != 1 {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
	if avps, err := m.FindAVPsWithPath([]interface{}{"Vendor-Specific-Application-Id", avp.VendorID}, 0); len(avps) != 1 {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
	if avps, err := m.FindAVPsWithPath([]interface{}{avp.VendorID}, 0); len(avps) != 1 {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
	if avps, err := m.FindAVPsWithPath([]interface{}{avp.VendorSpecificApplicationID, avp.OriginStateID}, 0); len(avps) != 0 {
		t.Errorf("Received nr of AVPs: %d, error: %v", len(avps), err)
	}
}

func TestMessageWriteTo(t *testing.T) {
	var mydictXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="4">
    <avp name="Service-Information" code="873" must="V,M" may="P" must-not="-" may-encrypt="N" vendor-id="10415">
      <data type="Grouped">
        <rule avp="IN-Information" required="false" max="1" />
      </data>
    </avp>
    <avp name="IN-Information" code="20300" must="V" may="P,M" must-not="-" may-encrypt="N" vendor-id="20300">
      <data type="Grouped">
        <rule avp="Charge-Flow-Type" required="false" max="1" />
        <rule avp="Calling-Vlr-Number" required="false" max="1" />
      </data>
    </avp>
    <avp name="Charge-Flow-Type" code="20339" must="V" may="P,M" must-not="-" may-encrypt="N" vendor-id="20300">
      <data type="Unsigned32" />
    </avp>
    <avp name="Calling-Vlr-Number" code="20302" must="V" may="P,M" must-not="-" may-encrypt="N" vendor-id="20300">
      <data type="UTF8String" />
    </avp>
  </application>
</diameter>`
	dict.Default.Load(bytes.NewReader([]byte(mydictXML)))
	m := NewRequest(CreditControl, 4, nil)
	m.NewAVP("Session-Id", avp.Mbit, 0, datatype.UTF8String("890f81bee22a0dfddc8b9037eb367781cea1f328"))
	m.NewAVP("Service-Information", avp.Mbit, 10415, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(20300, avp.Mbit, 20300, &GroupedAVP{ // IN-Information
				AVP: []*AVP{
					NewAVP(20339, avp.Mbit, 20300, datatype.Unsigned32(0)),  // Charge-Flow-Type
					NewAVP(20302, avp.Mbit, 20300, datatype.UTF8String("")), // Calling-Vlr-Number
				},
			}),
		}})
	if _, err := m.WriteTo(ioutil.Discard); err != nil {
		t.Error(err)
	}
}

func BenchmarkReadMessage(b *testing.B) {
	reader := bytes.NewReader(testMessage)
	for n := 0; n < b.N; n++ {
		ReadMessage(reader, dict.Default)
		reader.Seek(0, 0)
	}
}

func BenchmarkWriteMessage(b *testing.B) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	for n := 0; n < b.N; n++ {
		m.WriteTo(ioutil.Discard)
	}
}
