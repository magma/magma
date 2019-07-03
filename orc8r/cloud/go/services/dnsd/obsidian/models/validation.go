/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"errors"

	"magma/orc8r/cloud/go/services/dnsd/protos"
)

func (m *NetworkDNSConfig) ValidateNetworkConfig() error {
	if m == nil {
		return errors.New("NetworkDNSconfig is nil.")
	}
	if err := validateNetworkDNSRecordsConfig(m.Records); err != nil {
		return err
	}
	return nil
}

func validateNetworkDNSRecordsConfig(records []*NetworkDNSConfigRecordsItems0) error {
	if records == nil {
		return nil
	}

	for _, item := range records {
		if err := validateNetworkDNSConfigRecordsItems(item); err != nil {
			return err
		}
	}
	return nil
}

func validateNetworkDNSConfigRecordsItems(config *NetworkDNSConfigRecordsItems0) error {
	if config == nil {
		return errors.New("NetworkDNSconfig Records Item is nil.")
	}

	if err := protos.ValidateNetworkDNSConfigARecord(config.ARecord); err != nil {
		return err
	}
	if err := protos.ValidateNetworkDNSConfigAaaaRecord(config.AaaaRecord); err != nil {
		return err
	}

	if err := protos.ValidateNetworkDNSConfigDomain(config.Domain); err != nil {
		return err
	}

	if err := protos.ValidateNetworkDNSConfigCname(config.CnameRecord); err != nil {
		return err
	}

	return nil
}
