// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"os"
	"testing"
)

var testDicts = []string{
	"./testdata/base.xml",
	"./testdata/credit_control.xml",
	"./testdata/network_access_server.xml",
	"./testdata/tgpp_ro_rf.xml",
	"./testdata/tgpp_s6a.xml",
	"./testdata/tgpp_swx.xml"}

func TestNewParser(t *testing.T) {
	for _, dict := range testDicts {
		p, err := NewParser(dict)
		if err != nil {
			t.Fatalf("Error Creating Parser from %s: %s", dict, err)
		}
		t.Log(p)
	}
}

func TestLoadFile(t *testing.T) {
	for _, dict := range testDicts {
		p, _ := NewParser()
		if err := p.LoadFile(dict); err != nil {
			t.Fatalf("Error Loading %s: %s", dict, err)
		}
	}
}

func TestLoad(t *testing.T) {
	for _, dict := range testDicts {
		f, err := os.Open(dict)
		if err != nil {
			t.Fatalf("Error Opening %s: %s", dict, err)
		}
		defer f.Close()
		p, _ := NewParser()
		if err = p.Load(f); err != nil {
			t.Fatalf("Error Loading Parsing %s: %s", dict, err)
		}
	}
}
