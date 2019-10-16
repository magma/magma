/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package machine

import (
	"net"
	"sort"
)

type (
	// Interfaces ...
	Interfaces []net.Interface
)

// DefaultMacAddress ...
const DefaultMacAddress string = "00:00:00:00:00:00"

func (ifs Interfaces) Len() int {
	return len(ifs)
}

func (ifs Interfaces) Swap(i, j int) {
	ifs[i], ifs[j] = ifs[j], ifs[i]
}

func (ifs Interfaces) Less(i, j int) bool {
	return ifs[i].Name < ifs[j].Name
}

// GetMachineMACAddressID gets a unique MAC address which identifies the machine.
// This means that on the same machine, the same MAC will be returned on every call
// to this function.
func GetMachineMACAddressID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return DefaultMacAddress
	}

	sort.Sort(Interfaces(interfaces))

	for _, intf := range interfaces {
		mac := intf.HardwareAddr.String()
		if mac != "" {
			return mac
		}
	}
	return DefaultMacAddress
}
