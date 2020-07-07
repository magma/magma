/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

func NewDefaultSubscriberConfig() *NetworkSubscriberConfig {
	return &NetworkSubscriberConfig{
		NetworkWideBaseNames: []BaseName{"base1"},
		NetworkWideRuleNames: []string{"rule1"},
	}
}
