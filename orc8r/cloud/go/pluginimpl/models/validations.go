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
	"net"
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

func validateNetworkDNSRecordsConfig(records []*DNSConfigRecord) error {
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

func validateNetworkDNSConfigRecordsItems(config *DNSConfigRecord) error {
	if config == nil {
		return errors.New("NetworkDNSconfig Records Item is nil.")
	}

	if err := ValidateNetworkDNSConfigARecord(config.ARecord); err != nil {
		return err
	}
	if err := ValidateNetworkDNSConfigAaaaRecord(config.AaaaRecord); err != nil {
		return err
	}

	if err := ValidateNetworkDNSConfigDomain(config.Domain); err != nil {
		return err
	}

	if err := ValidateNetworkDNSConfigCname(config.CnameRecord); err != nil {
		return err
	}

	return nil
}

func ValidateNetworkDNSConfigARecord(ARecord []string) error {
	if ARecord == nil {
		return nil
	}
	for _, record := range ARecord {
		if net.ParseIP(record).To4() == nil {
			return errors.New("ARecord must be in the form of an IpV4 address.")
		}
	}
	return nil
}

func ValidateNetworkDNSConfigAaaaRecord(AaaaRecord []string) error {
	if AaaaRecord == nil {
		return nil
	}
	for _, record := range AaaaRecord {
		if net.ParseIP(record).To16() == nil {
			return errors.New("AaaaRecord must be in the form of an IpV6 address.")
		}
	}
	return nil
}

func ValidateNetworkDNSConfigDomain(domain string) error {
	// TODO: Figure out how to validate a string is a domain
	if domain == "" {
		return errors.New("Domain cannot be empty string.")
	}
	return nil
}

func ValidateNetworkDNSConfigCname(CnameRecord []string) error {
	if CnameRecord == nil {
		return nil
	}
	for _, record := range CnameRecord {
		if err := ValidateNetworkDNSConfigDomain(record); err != nil {
			return err
		}
	}
	return nil
}
