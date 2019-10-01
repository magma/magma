/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"fmt"

	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2869"
	"magma/feg/gateway/services/eap"

	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	Auth            = "\x73\xea\x5e\xdf\x10\x25\x45\x3b\x21\x15\xdb\xc2\xa9\x8a\x7c\x99"
	CalledStationID = "98-DE-D0-84-B5-47:CWF-TP-LINK_B547_5G"
)

// HandleRadius routes the Radius packet to the UE with the specified imsi.
func (srv *UESimServer) HandleRadius(imsi string, p radius.Packet) (radius.Packet, error) {
	// todo Validate the packet. (Requires keeping state)

	// Extract EAP packet.
	eapMessage, err := packet.NewPacketFromRadius(&p)
	if err != nil {
		err = errors.Wrap(err, "Error extracting EAP message from Radius packet")
		return radius.Packet{}, err
	}
	eapBytes, err := eapMessage.Bytes()
	if err != nil {
		err = errors.Wrap(err, "Error converting EAP packet to bytes")
		return radius.Packet{}, err
	}

	// Get the specified UE from the blobstore.
	ue, err := getUE(srv.store, imsi)
	if err != nil {
		return radius.Packet{}, err
	}

	// Generate EAP response.
	eapRes, err := srv.HandleEap(ue, eap.Packet(eapBytes))
	if err != nil {
		return radius.Packet{}, err
	}

	// Wrap EAP response in Radius packet.
	res, err := srv.EapToRadius(eapRes, imsi, p.Identifier+1)
	if err != nil {
		return radius.Packet{}, err
	}

	return res, err
}

// EapToRadius puts an Eap packet payload in a Radius packet.
func (srv *UESimServer) EapToRadius(eapP eap.Packet, imsi string, identifier uint8) (radius.Packet, error) {
	radiusP := radius.New(radius.CodeAccessRequest, []byte(srv.cfg.radiusSecret))
	radiusP.Identifier = identifier

	// Hardcode in the auth.
	copy(radiusP.Authenticator[:], []byte(Auth)[:])
	radiusP.Attributes[rfc2865.UserName_Type] = []radius.Attribute{
		radius.Attribute([]byte(imsi + IdentityPostfix)),
	}
	// TODO: Fetch UE MAC addr and use as CallingStationID
	radiusP.Attributes[rfc2865.CallingStationID_Type] = []radius.Attribute{
		radius.Attribute([]byte(srv.cfg.brMac)),
	}
	radiusP.Attributes[rfc2865.CalledStationID_Type] = []radius.Attribute{
		radius.Attribute([]byte(CalledStationID)),
	}
	encoded, err := radiusP.Encode()
	if err != nil {
		return radius.Packet{}, errors.Wrap(err, "Error encoding Radius packet")
	}
	// Put EAP message in the EAP message Attribute.
	encoded = append(encoded, uint8(rfc2869.EAPMessage_Type))
	encoded = append(encoded, uint8(len(eapP)+2))
	encoded = append(encoded, eapP...)

	// Add Message-Authenticator Attribute.
	encoded = srv.addMessageAuthenticator(encoded)

	// Parse to Radius packet.
	res, err := radius.Parse(encoded, []byte(srv.cfg.radiusSecret))
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return *res, nil
}

// addMessageAuthenticator calculates and adds the Message-Authenticator
// Attribute to a RADIUS packet.
func (srv *UESimServer) addMessageAuthenticator(encoded []byte) []byte {
	// Calculate new size
	size := uint16(len(encoded)) + radius.MessageAuthenticatorAttrLength
	binary.BigEndian.PutUint16(encoded[2:4], uint16(size))

	// Append the empty Message-Authenticator Attribute to the packet
	encoded = append(
		encoded,
		uint8(rfc2869.MessageAuthenticator_Type),
		uint8(radius.MessageAuthenticatorAttrLength),
	)
	encoded = append(encoded, make([]byte, 16)...)

	// Calculate Message-Authenticator and overwrite.
	hash := hmac.New(md5.New, []byte(srv.cfg.radiusSecret))
	hash.Write(encoded)
	encoded = hash.Sum(encoded[:len(encoded)-16])

	return encoded
}
