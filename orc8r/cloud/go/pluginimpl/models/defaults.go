/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import "github.com/go-openapi/swag"

func NewDefaultDNSConfig() *NetworkDNSConfig {
	return &NetworkDNSConfig{
		EnableCaching: swag.Bool(true),
		LocalTTL:      swag.Uint32(60),
	}
}

func NewDefaultFeaturesConfig() *NetworkFeatures {
	return &NetworkFeatures{Features: map[string]string{"foo": "bar"}}
}
