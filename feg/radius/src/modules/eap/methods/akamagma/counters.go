/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package akamagma

import (
	"fbc/cwf/radius/monitoring"
)

var (
	// MarshalProtocolState marshal eap state
	MarshalProtocolState = monitoring.NewOperation("eap_marshal_state")

	// UnmarshalProtocolState marshal eap state
	UnmarshalProtocolState = monitoring.NewOperation("eap_marshal_state")
)
