// +build go1.8
// +build linux,!386

// Copyright 2013-2018 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package sm

import "testing"

func TestS6aClientServerSCTP(t *testing.T) {
	testS6aClientServer("sctp", t)
}
