// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"bytes"
	"testing"
)

func TestApps(t *testing.T) {
	apps := Default.Apps()
	if len(apps) != 8 {
		t.Fatalf("Unexpected # of apps. Want 8, have %d", len(apps))
	}
	// Base protocol.
	if apps[0].ID != 0 {
		t.Fatalf("Unexpected app.ID. Want 0, have %d", apps[0].ID)
	}
	// Base accounting
	if apps[1].ID != 3 {
		t.Fatalf("Unexpected app.ID. Want 3, have %d", apps[1].ID)
	}
	// Credit-Control applications.
	if apps[2].ID != 4 {
		t.Fatalf("Unexpected app.ID. Want 4, have %d", apps[2].ID)
	}
	// 3GPP Gx Charging Control applications
	if apps[3].ID != 16777238 {
		t.Fatalf("Unexpected app.ID. Want 16777238, have %d", apps[3].ID)
	}
	// NASREQ applications
	if apps[4].ID != 1 {
		t.Fatalf("Unexpected app.ID. Want 1, have %d", apps[4].ID)
	}
	// 3GPP S6a applications
	if apps[6].ID != 16777251 {
		t.Fatalf("Unexpected app.ID. Want 16777251, have %d", apps[6].ID)
	}
	if apps[7].ID != 16777265 {
		t.Fatalf("Unexpected app.ID. Want 16777265, have %d", apps[7].ID)
	}
}

func TestApp(t *testing.T) {
	// Base protocol.
	if _, err := Default.App(0); err != nil {
		t.Fatal(err)
	}
	// Credit-Control applications.
	if _, err := Default.App(4); err != nil {
		t.Fatal(err)
	}
}

func findAVPCodeTest(t *testing.T, app uint32, codeStr string, vendor, expectedCode uint32) {
	if avp, err := Default.FindAVPWithVendor(app, codeStr, vendor); err != nil {
		t.Fatalf("FindAVP error: %v for app %d & %s AVP", err, app, codeStr)
	} else if avp.Code != expectedCode {
		t.Fatalf(
			"Unexpected code %d for %s AVP and %d vendor. Expected: %d",
			avp.Code, codeStr, vendor, expectedCode)
	}
}

func TestFindAVPWithVendor(t *testing.T) {
	var nokiaXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="43">
    <vendor id="94" name="Nokia" />
    <avp name="Session-Start-Indicator" code="5105" must="V" may="P,M" must-not="-" may-encrypt="N" vendor-id="94">
      <data type="UTF8String" />
    </avp>
  </application>
</diameter>`
	Default.Load(bytes.NewReader([]byte(nokiaXML)))
	if _, err := Default.FindAVPWithVendor(4, 999, UndefinedVendorID); err == nil {
		t.Error("Should get not found")
	}
	findAVPCodeTest(t, 4, "Session-Id", UndefinedVendorID, 263)
	findAVPCodeTest(t, 43, "Session-Start-Indicator", 94, 5105)
	findAVPCodeTest(t, 43, "Session-Start-Indicator", UndefinedVendorID, 5105)

	if _, err := Default.FindAVPWithVendor(4, "Session-Start-Indicator", 0); err == nil {
		t.Error("Should get not found")
	}
	findAVPCodeTest(t, 16777251, "Supported-Features", UndefinedVendorID, 628)

	// Test 'parent' AVP find - S6a app ID, tgpp_ro_rf dictionary
	findAVPCodeTest(t, 16777251, "GMLC-Address", UndefinedVendorID, 2405)

	if _, err := Default.FindAVPWithVendor(43, "User-Password", UndefinedVendorID); err == nil {
		t.Error("User-Password Should not be found for app 43")
	}
	findAVPCodeTest(t, 1, "User-Password", UndefinedVendorID, 2)
	findAVPCodeTest(t, 4, "User-Password", UndefinedVendorID, 2)
	findAVPCodeTest(t, 16777251, "User-Password", UndefinedVendorID, 2)
}

func TestFindAVP(t *testing.T) {
	if _, err := Default.FindAVP(999, 263); err != nil {
		t.Fatal(err)
	}
}

func TestScanAVP(t *testing.T) {
	if avp, err := Default.ScanAVP("Session-Id"); err != nil {
		t.Error(err)
	} else if avp.Code != 263 {
		t.Fatalf("Unexpected code %d for Session-Id AVP", avp.Code)
	}
}

func TestFindCommand(t *testing.T) {
	if cmd, err := Default.FindCommand(999, 257); err != nil {
		t.Error(err)
	} else if cmd.Short != "CE" {
		t.Fatalf("Unexpected command: %#v", cmd)
	}

	if cmd, err := Default.FindCommand(16777251, 316); err != nil {
		t.Error(err)
	} else if cmd.Short != "UL" {
		t.Fatalf("Unexpected command: %#v", cmd)
	}

	if cmd, err := Default.FindCommand(16777251, 318); err != nil {
		t.Error(err)
	} else if cmd.Short != "AI" {
		t.Fatalf("Unexpected command: %#v", cmd)
	}
}

func TestEnum(t *testing.T) {
	if item, err := Default.Enum(0, 274, 1); err != nil {
		t.Fatal(err)
	} else if item.Name != "AUTHENTICATE_ONLY" {
		t.Errorf(
			"Unexpected value %s, expected AUTHENTICATE_ONLY",
			item.Name,
		)
	}
}

func TestRule(t *testing.T) {
	if rule, err := Default.Rule(0, 284, "Proxy-Host"); err != nil {
		t.Fatal(err)
	} else if !rule.Required {
		t.Errorf("Unexpected rule %#v", rule)
	}
}

func BenchmarkFindAVPName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.FindAVP(0, "Session-Id")
	}
}

func BenchmarkFindAVPCode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.FindAVP(0, 263)
	}
}

func BenchmarkScanAVPName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.ScanAVP("Session-Id")
	}
}

func BenchmarkScanAVPCode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.ScanAVP(263)
	}
}
