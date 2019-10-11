/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package aruba_test

import (
	"testing"

	"fbc/lib/go/radius"
	. "fbc/lib/go/radius/vendors/aruba"
)

func TestLookup(t *testing.T) {
	p := radius.New(radius.CodeAccessRequest, []byte(`12345`))
	ArubaUserRole_SetString(p, "Admin")
	ArubaDeviceType_SetString(p, "Desktop")

	if dt := ArubaDeviceType_GetString(p); dt != "Desktop" {
		t.Fatalf("ArubaDeviceType = %v; expecting %v", dt, "Desktop")
	}
}
