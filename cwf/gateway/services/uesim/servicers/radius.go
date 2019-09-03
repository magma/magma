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
	"fbc/lib/go/radius/rfc2869"
	"magma/feg/gateway/services/eap"

	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	Auth       = "\x73\xea\x5e\xdf\x10\x25\x45\x3b\x21\x15\xdb\xc2\xa9\x8a\x7c\x99"
	Attributes = "\x01\x35\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30" +
		"\x39\x31\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d" +
		"\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72" +
		"\x6b\x2e\x6f\x72\x67\x04\x06\xc0\xa8\x00\x01\x05\x06\x00\x00\x00" +
		"\x00\x1e\x27\x39\x38\x2d\x44\x45\x2d\x44\x30\x2d\x38\x34\x2d\x42" +
		"\x35\x2d\x34\x37\x3a\x43\x57\x46\x2d\x54\x50\x2d\x4c\x49\x4e\x4b" +
		"\x5f\x42\x35\x34\x37\x5f\x35\x47\x1f\x13\x41\x43\x2d\x35\x46\x2d" +
		"\x33\x45\x2d\x31\x32\x2d\x38\x41\x2d\x42\x37\x0c\x06\x00\x00\x05" +
		"\x78\x3d\x06\x00\x00\x00\x13\x4d\x17\x43\x4f\x4e\x4e\x45\x43\x54" +
		"\x20\x30\x4d\x62\x70\x73\x20\x38\x30\x32\x2e\x31\x31\x67"
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
	res, err := srv.EapToRadius(eapRes, p.Identifier+1)
	if err != nil {
		return radius.Packet{}, err
	}

	return res, err
}

// EapToRadius puts an Eap packet payload in a Radius packet.
func (srv *UESimServer) EapToRadius(eapP eap.Packet, identifier uint8) (radius.Packet, error) {
	radiusP := radius.New(radius.CodeAccessRequest, []byte(srv.cfg.radiusSecret))
	radiusP.Identifier = identifier

	// Hardcode in the auth.
	copy(radiusP.Authenticator[:], []byte(Auth)[:])

	encoded, err := radiusP.Encode()
	if err != nil {
		return radius.Packet{}, errors.Wrap(err, "Error encoding Radius packet")
	}

	// Add Attributes.
	encoded = append(encoded, []byte(Attributes)...)

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
