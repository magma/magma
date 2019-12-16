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
	"github.com/go-openapi/swag"
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
	set := make(map[string][]uint32)
	for _, peer := range m.AllowedGrePeers {
		for _, key := range set[string(peer.IP)] {
			if swag.Uint32Value(peer.Key) == key {
				return errors.New(fmt.Sprintf("Found duplicate peer %s with key %d", string(peer.IP), key))
			}
		}
		set[string(peer.IP)] = append(set[string(peer.IP)], swag.Uint32Value(peer.Key))
	}
	return nil
}

func (m *CwfSubscriberDirectoryRecord) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
