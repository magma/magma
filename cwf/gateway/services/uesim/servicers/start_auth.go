/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fbc/lib/go/radius"
	"magma/feg/gateway/services/eap"
)

// todo use a config to assign this value
const (
	EapIdentityRequestPacket = "\x01\x00\x00\x05\x01"
)

// CreateEAPIdentityRequest simulates starting the EAP-AKA authentication by sending a UE an
// EAP Identity Request packet.
func (srv *UESimServer) CreateEAPIdentityRequest(imsi string) (*radius.Packet, error) {
	ue, err := getUE(srv.store, imsi)
	if err != nil {
		return nil, err
	}

	eapReponse, err := srv.HandleEap(ue, eap.Packet(EapIdentityRequestPacket))
	if err != nil {
		return nil, err
	}

	// Set packet Identifier to 0.
	return srv.EapToRadius(eapReponse, imsi, 0)
}
