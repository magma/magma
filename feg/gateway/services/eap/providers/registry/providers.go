/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package registry defines API to register and fing EAP providers
package registry

import (
	aka_provider "magma/feg/gateway/services/eap/providers/aka/provider"
)

func init() {
	Register(aka_provider.New())
}
