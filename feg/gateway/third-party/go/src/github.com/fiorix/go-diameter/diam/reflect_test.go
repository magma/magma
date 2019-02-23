// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

func TestUnmarshalAVP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type Data struct {
		OriginHost1 AVP  `avp:"Origin-Host"`
		OriginHost2 *AVP `avp:"Origin-Host"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v, ok := d.OriginHost1.Data.(datatype.DiameterIdentity); !ok {
		t.Fatalf("Unexpected value. Want datatype.DiameterIdentity, have %#v", d.OriginHost1.Data)
	} else if v != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", v)
	}
	if v, ok := d.OriginHost2.Data.(datatype.DiameterIdentity); !ok {
		t.Fatalf("Unexpected value. Want datatype.DiameterIdentity, have %#v", d.OriginHost2.Data)
	} else if v != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", v)
	}
}

func TestUnmarshalString(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type Data struct {
		OriginHost string `avp:"Origin-Host"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if d.OriginHost != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", d.OriginHost)
	}
}

func TestUnmarshalTimeDatatype(t *testing.T) {
	expectedTime := "2015-12-09 15:40:53 +0000 UTC"
	m, _ := ReadMessage(bytes.NewReader(testMessageWithVendorID), dict.Default)
	type Data struct {
		EventTimestamp datatype.Time `avp:"Event-Timestamp"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	timestamp := time.Time(d.EventTimestamp)
	if timestamp.UTC().String() != expectedTime {
		t.Fatalf("Unexpected value, want %s, have %s", expectedTime, timestamp.UTC().String())
	}
}

func TestUnmarshalTimeType(t *testing.T) {
	expectedTime := "2015-12-09 15:40:53 +0000 UTC"
	m, _ := ReadMessage(bytes.NewReader(testMessageWithVendorID), dict.Default)
	type Data struct {
		EventTimestamp time.Time `avp:"Event-Timestamp"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if d.EventTimestamp.UTC().String() != expectedTime {
		t.Fatalf("Unexpected value, want %s, have %s", expectedTime, d.EventTimestamp.UTC().String())
	}
}

func TestUnmarshalNetIP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type Data struct {
		HostIP1 *AVP   `avp:"Host-IP-Address"`
		HostIP2 net.IP `avp:"Host-IP-Address"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v := d.HostIP1.Data.(datatype.Address); net.IP(v).String() != "10.1.0.1" {
		t.Fatalf("Unexpected value. Want 10.1.0.1, have %s", v)
	}
	if v := d.HostIP2.String(); v != "10.1.0.1" {
		t.Fatalf("Unexpected value. Want 10.1.0.1, have %s", v)
	}
}

func TestUnmarshalInt(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type Data struct {
		VendorID1 *AVP `avp:"Vendor-Id"`
		VendorID2 int  `avp:"Vendor-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v := d.VendorID1.Data.(datatype.Unsigned32); v != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", v)
	}
	if d.VendorID2 != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", d.VendorID2)
	}
}

func TestUnmarshalSlice(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type Data struct {
		Vendors1 []*AVP `avp:"Supported-Vendor-Id"`
		Vendors2 []int  `avp:"Supported-Vendor-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if len(d.Vendors1) != 2 {
		t.Fatalf("Unexpected value. Want 2, have %d", len(d.Vendors1))
	}
	if v := d.Vendors1[0].Data.(datatype.Unsigned32); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := d.Vendors1[1].Data.(datatype.Unsigned32); v != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", v)
	}
	if len(d.Vendors2) != 2 {
		t.Fatalf("Unexpected value. Want 2, have %d", len(d.Vendors2))
	}
	if d.Vendors2[0] != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.Vendors2[0])
	}
	if d.Vendors2[1] != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", d.Vendors2[1])
	}
}

func TestUnmarshalGrouped(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type VSA struct {
		AuthAppID1 AVP  `avp:"Auth-Application-Id"`
		AuthAppID2 *AVP `avp:"Auth-Application-Id"`
		AuthAppID3 int  `avp:"Auth-Application-Id"`
		VendorID1  AVP  `avp:"Vendor-Id"`
		VendorID2  *AVP `avp:"Vendor-Id"`
		VendorID3  int  `avp:"Vendor-Id"`
	}
	type Data struct {
		VSA1 AVP  `avp:"Vendor-Specific-Application-Id"`
		VSA2 *AVP `avp:"Vendor-Specific-Application-Id"`
		VSA3 VSA  `avp:"Vendor-Specific-Application-Id"`
		VSA4 *VSA `avp:"Vendor-Specific-Application-Id"`
		VSA5 struct {
			AuthAppID1 AVP  `avp:"Auth-Application-Id"`
			AuthAppID2 *AVP `avp:"Auth-Application-Id"`
			AuthAppID3 int  `avp:"Auth-Application-Id"`
			VendorID1  AVP  `avp:"Vendor-Id"`
			VendorID2  *AVP `avp:"Vendor-Id"`
			VendorID3  int  `avp:"Vendor-Id"`
		} `avp:"Vendor-Specific-Application-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v, ok := d.VSA1.Data.(*GroupedAVP); !ok {
		t.Fatalf("Unexpected value. Want Grouped, have %v", d.VSA1)
	} else if len(v.AVP) != 2 { // There must be 2 AVPs in it.
		t.Fatalf("Unexpected value. Want 2, have %d", len(v.AVP))
	}
	if v, ok := d.VSA2.Data.(*GroupedAVP); !ok {
		t.Fatalf("Unexpected value. Want Grouped, have %s", d.VSA2)
	} else if len(v.AVP) != 2 { // There must be 2 AVPs in it.
		t.Fatalf("Unexpected value. Want 2, have %d", len(v.AVP))
	}
	if v := int(d.VSA3.AuthAppID1.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA3.AuthAppID2.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA3.AuthAppID3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA3.AuthAppID3)
	}
	if v := int(d.VSA3.VendorID1.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA3.VendorID2.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA3.VendorID3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA3.VendorID3)
	}
	if v := int(d.VSA4.AuthAppID1.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA4.AuthAppID2.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA4.AuthAppID3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA4.AuthAppID3)
	}
	if v := int(d.VSA4.VendorID1.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA4.VendorID2.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA4.VendorID3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA4.VendorID3)
	}
	if v := int(d.VSA5.AuthAppID1.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA5.AuthAppID2.Data.(datatype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA5.AuthAppID3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA5.AuthAppID3)
	}
	if v := int(d.VSA5.VendorID1.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA5.VendorID2.Data.(datatype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA5.VendorID3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA5.VendorID3)
	}
}

func TestUnmarshalGroupedSlice(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	type VSA struct {
		AuthAppID int `avp:"Auth-Application-Id"`
		VendorID  int `avp:"Vendor-Id"`
	}
	type Data struct {
		VSA1 []*VSA `avp:"Vendor-Specific-Application-Id"`
		VSA2 []*AVP `avp:"Vendor-Specific-Application-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if len(d.VSA1) != 1 {
		t.Fatalf("Unexpected value. Want 1, have %d", len(d.VSA1))
	}
	if len(d.VSA2) != 1 {
		t.Fatalf("Unexpected value. Want 1, have %d", len(d.VSA2))
	}
	if v, ok := d.VSA2[0].Data.(*GroupedAVP); !ok {
		t.Fatalf("Unexpected value. Want Grouped, have %s", d.VSA2)
	} else if len(v.AVP) != 2 { // There must be 2 AVPs in it.
		t.Fatalf("Unexpected value. Want 2, have %d", len(v.AVP))
	}
}

func TestUnmarshalCER(t *testing.T) {
	msg, err := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	type CER struct {
		OriginHost  string `avp:"Origin-Host"`
		OriginRealm string `avp:"Origin-Realm"`
		HostIP      net.IP `avp:"Host-IP-Address"`
		VendorID    int    `avp:"Vendor-Id"`
		ProductName string `avp:"Product-Name"`
		StateID     int    `avp:"Origin-State-Id"`
		Vendors     []int  `avp:"Supported-Vendor-Id"`
		AuthAppID   int    `avp:"Auth-Application-Id"`
		InbandSecID int    `avp:"Inband-Security-Id"`
		VSA         struct {
			AuthAppID int `avp:"Auth-Application-Id"`
			VendorID  int `avp:"Vendor-Id"`
		} `avp:"Vendor-Specific-Application-Id"`
		Firmware int `avp:"Firmware-Revision"`
	}
	var d CER
	if err := msg.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	switch {
	case d.OriginHost != "test":
		t.Fatalf("Unexpected Code. Want test, have %s", d.OriginHost)
	case d.OriginRealm != "localhost":
		t.Fatalf("Unexpected Code. Want localhost, have %s", d.OriginRealm)
	case d.HostIP.String() != "10.1.0.1":
		t.Fatalf("Unexpected Host-IP-Address. Want 10.1.0.1, have %s", d.HostIP)
	case d.VendorID != 13:
		t.Fatalf("Unexpected Host-Vendor-Id. Want 13, have %d", d.VendorID)
	case d.ProductName != "go-diameter":
		t.Fatalf("Unexpected Product-Name. Want go-diameter, have %s", d.ProductName)
	case d.StateID != 1397760650:
		t.Fatalf("Unexpected Origin-State-Id. Want 1397760650, have %d", d.StateID)
	case d.Vendors[0] != 10415:
		t.Fatalf("Unexpected Origin-State-Id. Want 10415, have %d", d.StateID)
	case d.Vendors[1] != 13:
		t.Fatalf("Unexpected Origin-State-Id. Want 13, have %d", d.StateID)
	case d.AuthAppID != 4:
		t.Fatalf("Unexpected Origin-State-Id. Want 4, have %d", d.AuthAppID)
	case d.InbandSecID != 0:
		t.Fatalf("Unexpected Origin-State-Id. Want 0, have %d", d.InbandSecID)
	case d.VSA.AuthAppID != 4:
		t.Fatalf("Unexpected Origin-State-Id. Want 4, have %d", d.VSA.AuthAppID)
	case d.VSA.VendorID != 10415:
		t.Fatalf("Unexpected Origin-State-Id. Want 10415, have %d", d.VSA.VendorID)
	case d.Firmware != 1:
		t.Fatalf("Unexpected Origin-State-Id. Want 1, have %d", d.Firmware)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	msg, err := ReadMessage(bytes.NewReader(testMessage), dict.Default)
	if err != nil {
		b.Fatal(err)
	}
	type CER struct {
		OriginHost  AVP    `avp:"Origin-Host"`
		OriginRealm *AVP   `avp:"Origin-Realm"`
		HostIP      net.IP `avp:"Host-IP-Address"`
		VendorID    int    `avp:"Vendor-Id"`
		ProductName string `avp:"Product-Name"`
		StateID     int    `avp:"Origin-State-Id"`
	}
	var cer CER
	for n := 0; n < b.N; n++ {
		msg.Unmarshal(&cer)
	}
}
