/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"net"

	"github.com/go-openapi/strfmt"
)

func (m BaseNames) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m RuleNames) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

// validateIPBlocks parses and validates IP networks containing subnet masks.
// Returns an error in case any IP network in list is invalid.
func validateIPBlocks(ipBlocks []string) error {
	for _, ipBlock := range ipBlocks {
		_, _, err := net.ParseCIDR(ipBlock)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *PolicyRule) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *RatingGroup) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

// ValidateModel does standard swagger validation and any custom validation
func (m *MutableRatingGroup) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *NetworkSubscriberConfig) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
