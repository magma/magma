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

	"magma/cwf/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	IdentityPostfix = "\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63" +
		"\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e" +
		"\x6f\x72\x67"
)

// Handle routes the EAP request to the UE with the specified imsi.
func (srv *UESimServer) Handle(imsi string, req eap.Packet) (res eap.Packet, err error) {
	err = req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Error validating EAP packet")
	}

	// Get the specified UE from the blobstore.
	store, err := srv.store.StartTransaction()
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				glog.Errorf("Error while rolling back transaction: %s", err)
			}
		}
	}()

	blob, err := store.Get(networkIDPlaceholder, storage.TypeAndKey{Type: blobTypePlaceholder, Key: imsi})
	if err != nil {
		return
	}
	ue, err := blobToUE(blob)
	if err != nil {
		return
	}

	switch aka.Subtype(req[eap.EapSubtype]) {
	case aka.SubtypeIdentity:
		return identityRequest(ue, req)
	default:
		return nil, errors.Errorf("Unsupported Subtype: %d", req[eap.EapSubtype])
	}
}

// Given a UE and the EAP identity request, generates the EAP response.
func identityRequest(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
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
