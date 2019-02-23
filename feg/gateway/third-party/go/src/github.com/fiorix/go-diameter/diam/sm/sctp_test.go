// +build go1.8
// +build linux,!386

// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"testing"
)

func TestHandleCER_HandshakeMetadataSCTP(t *testing.T) {
	testHandleCER_HandshakeMetadata(t, "sctp")
}

func testClient_Handshake_CustomIP_SCTP(t *testing.T) {
	testClient_Handshake_CustomIP(t, "sctp")
}

// TestStateMachineSCTP establishes a connection with a test SCTP server and
// sends a Re-Auth-Request message to ensure the handshake was
// completed and that the RAR handler has context from the peer.
func TestStateMachineSCTP(t *testing.T) {
	testStateMachine(t, "sctp")
}
