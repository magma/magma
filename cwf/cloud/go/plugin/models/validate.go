/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

func (m *CwfNetwork) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *NetworkCarrierWifiConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *GatewayCwfConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	set := make(map[string]int)
	for _, peer := range m.AllowedGrePeers {
		set[string(peer.IP)]++
		if set[string(peer.IP)] > 1 {
			return errors.New(fmt.Sprintf("Found duplicate peer %s", string(peer.IP)))
		}
	}
	return nil
}

func (m *CwfSubscriberDirectoryRecord) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
