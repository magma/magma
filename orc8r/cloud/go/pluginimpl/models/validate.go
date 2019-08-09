/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"

	"github.com/go-openapi/strfmt"
)

const echoKeyType = "ECHO"
const ecdsaKeyType = "SOFTWARE_ECDSA_SHA256"

func (m *Network) ValidateModel() error {
	if m == nil {
		return errors.New("Network is nil.")
	}
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if err := m.DNS.ValidateModel(); err != nil {
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

func (m *MagmadGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDevice) ValidateModel() error {
	if err := m.Key.ValidateModel(); err != nil {
		return err
	}
	return m.Validate(strfmt.Default)
}

func (m *ChallengeKey) ValidateModel() error {
	switch m.KeyType {
	case echoKeyType:
		if m.Key != nil {
			return errors.New("ECHO mode should not have key value")
		}
		return nil
	case ecdsaKeyType:
		if m.Key == nil {
			return fmt.Errorf("No key supplied")
		}
		_, err := x509.ParsePKIXPublicKey([]byte(*m.Key))
		if err != nil {
			return fmt.Errorf("Failed to parse key: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("Unknown key type %s", m.KeyType)
	}
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
