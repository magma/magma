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

	"github.com/go-openapi/strfmt"
)

func (n *Network) ValidateModel() error {
	if n == nil {
		return errors.New("Network is nil.")
	}
	if err := n.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := n.DNS.ValidateModel(); err != nil {
		return err
	}
	return nil
}

func (m *NetworkDNSConfig) ValidateModel() error {
	if err := validateNetworkDNSRecordsConfig(m.Records); err != nil {
		return err
	}
	return m.Validate(strfmt.Default)
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

	if err := validateNetworkDNSConfigARecord(config.ARecord); err != nil {
		return err
	}
	if err := validateNetworkDNSConfigAaaaRecord(config.AaaaRecord); err != nil {
		return err
	}

	return nil
}

func validateNetworkDNSConfigARecord(ARecord []string) error {
	for _, record := range ARecord {
		if net.ParseIP(record).To4() == nil {
			return errors.New("ARecord must be in the form of an IpV4 address.")
		}
	}
	return nil
}

func validateNetworkDNSConfigAaaaRecord(AaaaRecord []string) error {
	for _, record := range AaaaRecord {
		if net.ParseIP(record).To16() == nil {
			return errors.New("AaaaRecord must be in the form of an IpV6 address.")
		}
	}
	return nil
}
