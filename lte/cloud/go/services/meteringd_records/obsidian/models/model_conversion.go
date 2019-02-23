/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"fmt"
	"reflect"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/obsidian/models"
	"magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
)

// mconfig_converters provides model receiver convertors to/from protobuf messages
type ProtoConverter interface {
	FromProto(msg proto.Message) error
	ToProto(msg proto.Message) error
	Verify() error
}

// Default fmt registry implementation suggests that it's thread safe, but
// we need to monitor it
var sharedFRFormatsRegistry = strfmt.NewFormats()

// FlowRecords's FromProto fills in models.FlowRecordsstruct from
// passed protos.FlowRecord
func (flowRecord *FlowRecord) FromProto(pfrm proto.Message) error {
	flowRecordProto, ok := pfrm.(*lteprotos.FlowRecord)
	if !ok {
		return fmt.Errorf(
			"Invalid Source Type %s, *protos.FlowRecord expected",
			reflect.TypeOf(pfrm))
	}
	if flowRecord != nil {
		if flowRecordProto != nil {
			protos.FillIn(flowRecordProto, flowRecord)
			flowRecord.SubscriberID = SubscriberID(flowRecordProto.Sid)
			return flowRecord.Verify()
		}
	}
	return nil
}

// Verify validates given FlowRecord
func (flowRecord *FlowRecord) Verify() error {
	if flowRecord == nil {
		return fmt.Errorf("Nil FlowRecord pointer")
	}
	err := flowRecord.Validate(sharedFRFormatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("Flow Record Validation Error: %s", err)
	}
	return err
}
