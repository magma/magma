/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"
	"io"
	"reflect"

	"magma/cwf/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/lte/cloud/go/crypto"
	"magma/lte/cloud/go/services/eps_authentication/servicers"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	IND            = 0
	CheckcodeValue = "\x00\x00\x86\xe8\x20\x4d\xc6\xe1\xe3\xd8\x94\x44\x3c\x26" +
		"\xa7\xc6\x5d\xee\x3c\x42\xab\xf8"
)

// handleEapAka routes the EAP-AKA request to the UE with the specified imsi.
func (srv *UESimServer) handleEapAka(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
	switch aka.Subtype(req[eap.EapSubtype]) {
	case aka.SubtypeIdentity:
		return eapAkaIdentityRequest(ue, req)
	case aka.SubtypeChallenge:
		return eapAkaChallengeRequest(ue, srv.op, srv.amf, req)
	default:
		return nil, errors.Errorf("Unsupported Subtype: %d", req[eap.EapSubtype])
	}
}

// Given a UE and the EAP-AKA identity request, generates the EAP response.
func eapAkaIdentityRequest(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new attribute scanner")
	}

	var a eap.Attribute

	// Parse out attributes.
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_PERMANENT_ID_REQ, aka.AT_ANY_ID_REQ:
			// Create the response EAP packet with the identity attribute.
			p := eap.NewPacket(
				eap.ResponseCode,
				req.Identifier(),
				[]byte{aka.TYPE, byte(aka.SubtypeIdentity), 0, 0},
			)

			// Append Identity Attribute data to packet.
			id := []byte("\x30" + ue.Imsi + IdentityPostfix)
			p, err = p.Append(
				eap.NewAttribute(
					aka.AT_IDENTITY,
					append(
						[]byte{uint8(len(id) >> 8), uint8(len(id))}, // actual len of Identity
						id...,
					),
				),
			)
			if err != nil {
				return nil, errors.Wrap(err, "Error appending attribute to packet")
			}
			return p, nil
		default:
			glog.Info(fmt.Sprintf("Unexpected EAP-AKA Identity Request Attribute type %d", a.Type()))
		}
	}
	return nil, errors.Wrap(err, "Error while processing EAP-AKA Identity Request")
}

type challengeAttributes struct {
	rand eap.Attribute
	autn eap.Attribute
	mac  eap.Attribute
}

// Given an EAP packet, parses out the RAND, AUTN, and MAC.
func parseChallengeAttributes(req eap.Packet) (challengeAttributes, error) {
	attrs := challengeAttributes{}

	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		return attrs, errors.Wrap(err, "Error creating new attribute scanner")
	}
	var a eap.Attribute
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_RAND:
			attrs.rand = a
		case aka.AT_AUTN:
			attrs.autn = a
		case aka.AT_MAC:
			if len(a.Marshaled()) < aka.ATT_HDR_LEN+aka.MAC_LEN {
				return attrs, fmt.Errorf("Malformed AT_MAC")
			}
			attrs.mac = a
		default:
			glog.Info(fmt.Sprintf("Unexpected EAP-AKA Challenge Request Attribute type %d", a.Type()))
		}
	}
	return attrs, err
}

// Given a UE, the Op, the Amf, and the EAP challenge, generates the EAP response.
func eapAkaChallengeRequest(ue *protos.UEConfig, op []byte, amf []byte, req eap.Packet) (eap.Packet, error) {
	attrs, err := parseChallengeAttributes(req)
	if err != io.EOF {
		return nil, errors.Wrap(err, "Error while parsing attributes of request packet")
	}
	if attrs.rand == nil || attrs.autn == nil || attrs.mac == nil {
		return nil, errors.Errorf("Missing one or more expected attributes\nRAND: %s\nAUTN: %s\nMAC: %s\n", attrs.rand, attrs.autn, attrs.mac)
	}

	// Parse out RAND, expected AUTN, and expected MAC values.
	rand := attrs.rand.Marshaled()[aka.ATT_HDR_LEN:]
	expectedAutn := attrs.autn.Marshaled()[aka.ATT_HDR_LEN:]
	expectedMac := attrs.mac.Marshaled()[aka.ATT_HDR_LEN:]

	id := []byte("\x30" + ue.Imsi + IdentityPostfix)
	key := []byte(ue.AuthKey)

	// Calculate SQN using SEQ and IND
	sqn := servicers.SeqToSqn(ue.Seq, IND) // todo decide how to increment SEQ

	// Calculate Opc using key and Op, and verify that it matches the UE's Opc
	opc, err := crypto.GenerateOpc(key, op)
	if err != nil {
		return nil, fmt.Errorf("Error while calculating Opc")
	}
	if !reflect.DeepEqual(opc[:], ue.AuthOpc) {
		return nil, fmt.Errorf("Invalid Opc: Expected Opc: %x; Actual Opc: %x", opc[:], ue.AuthOpc)
	}

	// Calculate RES and other keys.
	milenage, err := crypto.NewMilenageCipher(amf)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating milenage cipher")
	}
	vector, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc[:], sqn)
	if err != nil {
		return nil, errors.Wrap(err, "Error calculating authentication vector")
	}

	// Make copy of packet and zero out MAC value.
	copyReq := make([]byte, len(req))
	copy(copyReq, req)
	copyAttrs, err := parseChallengeAttributes(eap.Packet(copyReq))
	if err != io.EOF {
		return nil, errors.Wrap(err, "Error while parsing attributes of copied request packet")
	}
	copyMacBytes := copyAttrs.mac.Marshaled()
	for i := aka.ATT_HDR_LEN; i < len(copyMacBytes); i++ {
		copyMacBytes[i] = 0
	}

	// Calculate and verify MAC.
	_, kAut, _, _ := aka.MakeAKAKeys(id, vector.IntegrityKey[:], vector.ConfidentialityKey[:])
	mac := aka.GenMac(copyReq, kAut)
	if !reflect.DeepEqual(expectedMac, mac) {
		return nil, fmt.Errorf("Invalid MAC: Expected MAC: %x; Actual MAC: %x", expectedMac, mac)
	}

	// Calculate and verify AUTN.
	if !reflect.DeepEqual(expectedAutn, vector.Autn[:]) {
		return nil, fmt.Errorf("Invalid AUTN: Expected AUTN: %x; Actual AUTN: %x", expectedAutn, vector.Autn[:])
	}

	// Create the response EAP packet.
	p := eap.NewPacket(eap.ResponseCode, req.Identifier(), []byte{aka.TYPE, byte(aka.SubtypeChallenge), 0, 0})

	// Add the RES attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_RES,
			append(
				[]byte{uint8(len(vector.Xres[:]) * 8 >> 8), uint8(len(vector.Xres[:]) * 8)},
				[]byte(vector.Xres[:])...,
			),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error appending attribute to packet")
	}

	// Add the CHECKCODE attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_CHECKCODE,
			[]byte(CheckcodeValue),
		),
	)

	atMacOffset := len(p) + aka.ATT_HDR_LEN

	// Add the empty MAC attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_MAC,
			append(make([]byte, 2+16)),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error appending attribute to packet")
	}

	// Calculate and Copy MAC into packet.
	mac = aka.GenMac(p, kAut)
	copy(p[atMacOffset:], mac)

	return p, nil
}
