/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

func NewDefaultDNSConfig() *NetworkDNSConfig {
	return &NetworkDNSConfig{
		EnableCaching: bPtr(true),
		LocalTTL:      iPtr(60),
	}
}

func NewDefaultFeaturesConfig() *NetworkFeatures {
	return &NetworkFeatures{Features: map[string]string{"foo": "bar"}}
}

func bPtr(b bool) *bool {
	return &b
}

func iPtr(i int32) *int32 {
	return &i
}
