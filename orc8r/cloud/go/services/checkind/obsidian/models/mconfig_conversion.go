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

	"magma/orc8r/cloud/go/obsidian/models"
	"magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
)

var formatsRegistry = strfmt.NewFormats()

// Config conversion

// GatewayStatus's FromMconfig fills in models.GatewayStatus from
// passed protos.GatewayStatus
func (mstatus *GatewayStatus) FromMconfig(pm proto.Message) error {
	pstatus, ok := pm.(*protos.GatewayStatus)
	if !ok {
		return fmt.Errorf(
			"Invalid Source Type %s, *protos.GatewayStatus expected",
			reflect.TypeOf(pm))
	}
	if pstatus != nil && mstatus != nil {
		mstatus.CheckinTime = pstatus.Time
		if pstatus.Checkin == nil {
			mstatus.SystemStatus = nil
			mstatus.HardwareID = ""
			mstatus.Version = ""
			mstatus.Meta = map[string]string{}
			mstatus.VpnIP = ""
			mstatus.KernelVersion = ""
		} else {
			mstatus.HardwareID = pstatus.Checkin.GatewayId
			mstatus.Version = pstatus.Checkin.MagmaPkgVersion
			if pstatus.Checkin.SystemStatus == nil {
				mstatus.SystemStatus = nil
			} else {
				if mstatus.SystemStatus == nil {
					mstatus.SystemStatus = new(SystemStatus)
				}
				protos.FillIn(pstatus.Checkin.SystemStatus, mstatus.SystemStatus)
			}
			mstatus.Meta = map[string]string{}
			if pstatus.Checkin.Status != nil && pstatus.Checkin.Status.Meta != nil {
				for key, value := range pstatus.Checkin.Status.Meta {
					mstatus.Meta[key] = value
				}
			}
			mstatus.VpnIP = pstatus.Checkin.VpnIp
			mstatus.KernelVersion = pstatus.Checkin.KernelVersion
			mstatus.KernelVersionsInstalled = pstatus.Checkin.KernelVersionsInstalled
		}
	}
	return nil
}

// GatewayStatus is read only (GET) property, ToMconfig should not be used and
// doesn't need to be fully implemented
func (_ *GatewayStatus) ToMconfig(_ proto.Message) error {
	return fmt.Errorf("Not Implemented, GatewayStatus cannot be 'set'")
}

// Verify validates given GatewayStatus
func (mstatus *GatewayStatus) Verify() error {
	if mstatus == nil {
		return fmt.Errorf("Nil GatewayStatus pointer")
	}
	err := mstatus.Validate(formatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("GatewayStatus Validation Error: %s", err)
	}
	return err
}
