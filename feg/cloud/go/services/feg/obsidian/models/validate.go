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

	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
)

func (m FegNetworkID) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if !swag.IsZero(m) {
		exists, err := configurator.DoesNetworkExist(string(m))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to search for network %s", string(m)))
		}
		if !exists {
			return errors.New(fmt.Sprintf("Network: %s does not exist", string(m)))
		}
	}
	return nil
}

func (m *FederatedNetworkConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	nid := *m.FegNetworkID
	if !swag.IsZero(nid) {
		exists, err := configurator.DoesNetworkExist(string(nid))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to search for network %s", string(nid)))
		}
		if !exists {
			return errors.New(fmt.Sprintf("Network: %s does not exist", string(nid)))
		}
	}
	return nil
}

func (m *DiameterClientConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *DiameterServerConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *EapAkaTimeouts) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *GatewayFederationConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *NetworkFederationConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *SubscriptionProfile) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
