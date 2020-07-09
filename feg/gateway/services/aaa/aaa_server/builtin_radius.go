//
// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

// +build with_builtin_radius

// Package main
package main

import (
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa/radius"
)

func startBuiltInRadius(cfg *mconfig.AAAConfig) {
	go radius.New(cfg.GetRadiusConfig()).StartAuth()
}
