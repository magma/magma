/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package testcore

import "magma/orc8r/cloud/go/registry"

const (
	MockVLRServiceName  = "MOCK_VLR"
	MockOCSServiceName  = "MOCK_OCS"
	MockPCRFServiceName = "MOCK_PCRF"
	MockHSSServiceName  = "HSS"

	HSSServiceHost = "localhost"
	HSSServicePort = 9204
)

func init() {
	registry.AddService(registry.ServiceLocation{
		Name: MockHSSServiceName,
		Host: HSSServiceHost,
		Port: HSSServicePort,
	})
}
