/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"errors"
	"net"
)

// See dns_service.proto for full documentation on all config protobuf
// fields and which ones are optional/required.
func ValidateNetworkConfig(config *NetworkDNSConfig) error {
	if config == nil {
		return errors.New("NetworkDNSconfig is nil.")
	}
	if err := validateNetworkDNSRecordsConfig(config.GetRecords()); err != nil {
		return err
	}
	return nil
}

func validateNetworkDNSRecordsConfig(records []*NetworkDNSConfigRecordsItems) error {
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

func validateNetworkDNSConfigRecordsItems(config *NetworkDNSConfigRecordsItems) error {
	if config == nil {
		return errors.New("NetworkDNSconfig Records Item is nil.")
	}

	if err := validateNetworkDNSConfigARecord(config.GetARecord()); err != nil {
		return err
	}
	if err := validateNetworkDNSConfigAaaaRecord(config.GetAaaaRecord()); err != nil {
		return err
	}

	if err := validateNetworkDNSConfigDomain(config.GetDomain()); err != nil {
		return err
	}

	if err := validateNetworkDNSConfigCname(config.GetCnameRecord()); err != nil {
		return err
	}

	return nil
}

func validateNetworkDNSConfigARecord(ARecord []string) error {
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

func validateNetworkDNSConfigAaaaRecord(AaaaRecord []string) error {
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

func validateNetworkDNSConfigDomain(domain string) error {
	// TODO: Figure out how to validate a string is a domain
	if domain == "" {
		return errors.New("Domain cannot be empty string.")
	}
	return nil
}

func validateNetworkDNSConfigCname(CnameRecord []string) error {
	if CnameRecord == nil {
		return nil
	}
	for _, record := range CnameRecord {
		if err := validateNetworkDNSConfigDomain(record); err != nil {
			return err
		}
	}
	return nil
}
