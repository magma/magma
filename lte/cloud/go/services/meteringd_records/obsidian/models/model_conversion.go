/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func (m *FlowRecord) FromProto(record *lteprotos.FlowRecord) *FlowRecord {
	m.SubscriberID = record.Sid
	m.BytesRx = swag.Uint64(record.BytesRx)
	m.BytesTx = swag.Uint64(record.BytesTx)
	m.PktsRx = swag.Uint64(record.PktsRx)
	m.PktsTx = swag.Uint64(record.PktsTx)
	return m
}

func (m *FlowRecord) Verify() error {
	return m.Validate(strfmt.Default)
}
