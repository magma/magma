/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_utils

import (
	"crypto/x509"
	"testing"
	"time"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/lib/go/security/cert"
	"magma/orc8r/lib/go/security/csr"

	"github.com/stretchr/testify/assert"
)

// StartMockGwAccessControl starts certifier & adds a Gw Identities for
// Gateways with hwGwIds.
// Returns a list of corresponding Certificate Serial Numbers.
func StartMockGwAccessControl(t *testing.T, hwGwIds []string) []string {
	// Start services
	test_init.StartTestService(t)

	result := make([]string, len(hwGwIds))
	for idx, hwId := range hwGwIds {
		// create and sign Gw's csr
		csrMsg, err := csr.CreateCSRForId(
			time.Duration(time.Hour*4), identity.NewGateway(hwId, "", ""))
		assert.NoError(t, err)

		certMsg, err := certifier.SignCSR(csrMsg)
		assert.NoError(t, err, "Failed to sign Gateway's CSR")
		// get cert sn from cert
		gwCert, err := x509.ParseCertificates(certMsg.CertDer)
		assert.NoError(t, err, "Failed to parse Gateway's cert")

		certSerialNum := cert.SerialToString(gwCert[0].SerialNumber)
		t.Logf("Test Gateway Certificate SN: %s", certSerialNum)

		result[idx] = certSerialNum
	}
	return result
}
