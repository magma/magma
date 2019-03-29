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
			mstatus.PlatformInfo = nil
			mstatus.MachineInfo = nil
			mstatus.Meta = map[string]string{}

			// Deprecated fields
			mstatus.HardwareID = ""
			mstatus.Version = ""
			mstatus.VpnIP = ""
			mstatus.KernelVersion = ""
		} else {
			mstatus.HardwareID = pstatus.Checkin.GatewayId
			if pstatus.Checkin.SystemStatus == nil {
				mstatus.SystemStatus = nil
			} else {
				if mstatus.SystemStatus == nil {
					mstatus.SystemStatus = new(SystemStatus)
				}
				mstatus.SystemStatus.fillSystemStatus(pstatus.Checkin.SystemStatus)
			}

			if mstatus.PlatformInfo == nil {
				mstatus.PlatformInfo = new(PlatformInfo)
			}
			if pstatus.Checkin.PlatformInfo == nil {
				// Fallback to using deprecated fields to fill platform info
				mstatus.PlatformInfo.VpnIP = pstatus.Checkin.VpnIp
				mstatus.PlatformInfo.Packages = []*Package{
					{
						Name:    "magma",
						Version: pstatus.Checkin.MagmaPkgVersion,
					},
				}
				mstatus.PlatformInfo.KernelVersion = pstatus.Checkin.KernelVersion
				mstatus.PlatformInfo.KernelVersionsInstalled = pstatus.Checkin.KernelVersionsInstalled
			} else {
				mstatus.PlatformInfo.fillPlatformInfo(pstatus.Checkin.PlatformInfo)
			}

			if pstatus.Checkin.MachineInfo == nil {
				mstatus.MachineInfo = nil
			} else {
				if mstatus.MachineInfo == nil {
					mstatus.MachineInfo = new(MachineInfo)
				}
				mstatus.MachineInfo.fillMachineInfo(pstatus.Checkin.MachineInfo)
			}
			mstatus.Meta = map[string]string{}
			if pstatus.Checkin.Status != nil && pstatus.Checkin.Status.Meta != nil {
				for key, value := range pstatus.Checkin.Status.Meta {
					mstatus.Meta[key] = value
				}
			}

			// Populate deprecated fields to support API backwards compatibility
			// TODO: Remove this and related tests when deprecated fields are no longer used
			mstatus.VpnIP = mstatus.PlatformInfo.VpnIP
			mstatus.KernelVersion = mstatus.PlatformInfo.KernelVersion
			mstatus.KernelVersionsInstalled = mstatus.PlatformInfo.KernelVersionsInstalled
			for _, mPackage := range mstatus.PlatformInfo.Packages {
				if mPackage.Name == "magma" {
					mstatus.Version = mPackage.Version
				}
			}
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

func (mSystemStatus *SystemStatus) fillSystemStatus(pSystemStatus *protos.SystemStatus) {
	protos.FillIn(pSystemStatus, mSystemStatus)

	if pSystemStatus.DiskPartitions != nil {
		mSystemStatus.DiskPartitions = []*DiskPartition{}
		for _, pDiskPartition := range pSystemStatus.DiskPartitions {
			mDiskPartition := &DiskPartition{}
			protos.FillIn(pDiskPartition, mDiskPartition)
			mSystemStatus.DiskPartitions = append(mSystemStatus.DiskPartitions, mDiskPartition)
		}
	}
}

func (mPlatformInfo *PlatformInfo) fillPlatformInfo(pPlatformInfo *protos.PlatformInfo) {
	protos.FillIn(pPlatformInfo, mPlatformInfo)

	if pPlatformInfo.Packages != nil {
		mPlatformInfo.Packages = []*Package{}
		for _, pPackage := range pPlatformInfo.Packages {
			mPackage := &Package{}
			protos.FillIn(pPackage, mPackage)
			mPlatformInfo.Packages = append(mPlatformInfo.Packages, mPackage)
		}
	}
}

func (mMachineInfo *MachineInfo) fillMachineInfo(pMachineInfo *protos.MachineInfo) {
	protos.FillIn(pMachineInfo, mMachineInfo)

	if pMachineInfo.CpuInfo != nil {
		mMachineInfo.CPUInfo = &MachineInfoCPUInfo{}
		protos.FillIn(pMachineInfo.CpuInfo, mMachineInfo.CPUInfo)
	}

	if pMachineInfo.NetworkInfo != nil {
		mMachineInfo.NetworkInfo = &MachineInfoNetworkInfo{}
		if pMachineInfo.NetworkInfo.NetworkInterfaces != nil {
			mMachineInfo.NetworkInfo.NetworkInterfaces = []*NetworkInterface{}
			for _, pNetworkInterfaces := range pMachineInfo.NetworkInfo.NetworkInterfaces {
				mNetworkInterface := &NetworkInterface{}
				protos.FillIn(pNetworkInterfaces, mNetworkInterface)
				switch pNetworkInterfaces.Status {
				case protos.NetworkInterface_UP:
					mNetworkInterface.Status = NetworkInterfaceStatusUP
				case protos.NetworkInterface_DOWN:
					mNetworkInterface.Status = NetworkInterfaceStatusDOWN
				default:
					mNetworkInterface.Status = NetworkInterfaceStatusUNKNOWN
				}
				mMachineInfo.NetworkInfo.NetworkInterfaces = append(mMachineInfo.NetworkInfo.NetworkInterfaces, mNetworkInterface)
			}
		}
		if pMachineInfo.NetworkInfo.RoutingTable != nil {
			mMachineInfo.NetworkInfo.RoutingTable = []*Route{}
			for _, pRoute := range pMachineInfo.NetworkInfo.RoutingTable {
				mRoute := &Route{}
				protos.FillIn(pRoute, mRoute)
				mMachineInfo.NetworkInfo.RoutingTable = append(mMachineInfo.NetworkInfo.RoutingTable, mRoute)
			}
		}
	}
}
