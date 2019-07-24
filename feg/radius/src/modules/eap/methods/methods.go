/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package methods

import (
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
)

// MethodConfig An abstract for configuration method
type MethodConfig interface{}

// HandlerResponse the response of an EAP Method handler
type HandlerResponse struct {
	// The EAP packet to be sent as response
	Packet *packet.Packet

	// RadiusCode The RADIUS response code that should be sent
	// (per EAP<->RADIUS binding RFC, RADIUS code is dependent on
	// the EAP response packet itself)
	RadiusCode radius.Code

	// NewProtocolState The new state of the protocol to persist
	NewProtocolState string

	// ExtraAttributes contains extra RADIUS attributes to be added to the response
	ExtraAttributes radius.Attributes
}

// EapMethod the interface between RADIUS server and EAP method
type EapMethod interface {
	// Handle an EAP packet
	Handle(c *modules.RequestContext, p *packet.Packet, eapState string, r *radius.Request) (*HandlerResponse, error)
}
