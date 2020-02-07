/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/cloud/go/services/certifier/protos"
)

// CertifierStorage provides storage functionality for mapping serial numbers to certificate information.
type CertifierStorage interface {
	// ListSerialNumbers returns all tracked serial numbers.
	ListSerialNumbers() ([]string, error)

	// GetCertInfo returns the certificate info associated with the serial number.
	// If not found, returns ErrNotFound from magma/orc8r/cloud/go/errors.
	GetCertInfo(serialNumber string) (*protos.CertificateInfo, error)

	// GetManyCertInfo maps the passed serial numbers to their associated certificate info.
	GetManyCertInfo(serialNumbers []string) (map[string]*protos.CertificateInfo, error)

	// GetAllCertInfo returns a map of all serial numbers to their associated certificate info.
	GetAllCertInfo() (map[string]*protos.CertificateInfo, error)

	// PutCertInfo associates certificate info with the passed serial number.
	PutCertInfo(serialNumber string, certInfo *protos.CertificateInfo) error

	// DeleteCertInfo removes the serial number and its certificate info.
	// Returns success even when nothing is deleted (i.e. serial number not found).
	DeleteCertInfo(serialNumber string) error
}
