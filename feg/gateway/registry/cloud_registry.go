/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

import (
	"magma/gateway/service_registry"
)

// Get returns a singleton of gateway's registry
func Get() service_registry.GatewayRegistry {
	return service_registry.Get()
}
