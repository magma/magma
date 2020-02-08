/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

import (
	"magma/gateway/cloud_registry"
)

// CloudRegistry interface for FeG scoped users, see cloud_registry.CloudRegistry
type CloudRegistry interface {
	cloud_registry.CloudRegistry
}

// NewCloudRegistry returns a new instance of gateway's cloud registry
func NewCloudRegistry() CloudRegistry {
	return cloud_registry.New()
}
