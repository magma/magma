/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package radius_test

import (
	"strings"
	"testing"

	"fbc/lib/go/radius"
)

func TestParseAttributes_invalid(t *testing.T) {
	tests := []struct {
		Wire  string
		Error string
	}{
		{"\x01", "short buffer"},

		{"\x01\xff", "invalid attribute length"},
		{"\x01\x01", "invalid attribute length"},
	}

	for _, test := range tests {
		attrs, err := radius.ParseAttributes([]byte(test.Wire))
		if len(attrs) != 0 {
			t.Errorf("(%#x): expected empty attrs, got %v", test.Wire, attrs)
		} else if err == nil {
			t.Errorf("(%#x): expected error, got none", test.Wire)
		} else if !strings.Contains(err.Error(), test.Error) {
			t.Errorf("(%#x): expected error %q, got %q", test.Wire, test.Error, err.Error())
		}
	}
}
