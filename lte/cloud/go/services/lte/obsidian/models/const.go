/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

const (
	// NATAllocationMode NAT IP allocation mode
	NATAllocationMode = "NAT"
	// StaticAllocationMode Static IP allocation mode
	StaticAllocationMode = "STATIC"
	// DHCPPassthroughAllocationMode DHCP Passthrough (carrier wifi) IP allocation mode
	DHCPPassthroughAllocationMode = "DHCP_PASSTHROUGH"
	// DHCPBroadcastAllocationMode DHCP Broadcast IP allocation mode
	DHCPBroadcastAllocationMode = "DHCP_BROADCAST"
)
