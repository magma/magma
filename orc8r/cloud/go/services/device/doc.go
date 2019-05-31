/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package device contains the device service.
// The device service is a simple blob-storage service for tracking
// physical device.
package device

// SerdeDomain is the domain for all Serde implementations for the device
// service
const (
	SerdeDomain = "device"

	// ServiceName is the name of this service
	ServiceName = "DEVICE"

	// DBTableName is the name of the sql table used for this service
	DBTableName = "device"

	GatewayInfoType = "access_gateway_record"
)
