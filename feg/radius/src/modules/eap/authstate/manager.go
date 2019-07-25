/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package authstate

import (
	"fbc/cwf/radius/modules/eap/packet"

	"fbc/lib/go/radius"
)

// Container A storable container for protocol state
type Container struct {
	LogCorrelationID uint64         `json:"correlation_id"` // Request correlation id
	EapType          packet.EAPType `json:"eap_type"`       // EAP type of the auth session
	ProtocolState    string         `json:"protocol_state"` // EAP-* Protocol-specific state
	RadiusSessionID  *string        `json:"session_id"`     // RADIUS Session ID
}

// Manager an interface for EAP state management storage
type Manager interface {
	Set(authReq *radius.Packet, eaptype packet.EAPType, state Container) error
	Get(authReq *radius.Packet, eaptype packet.EAPType) (*Container, error)
	Reset(authReq *radius.Packet, eapType packet.EAPType) error
}
