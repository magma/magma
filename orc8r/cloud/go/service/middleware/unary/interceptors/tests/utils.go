/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tests

import (
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/identity"
	security_cert "magma/orc8r/cloud/go/security/cert"
	"magma/orc8r/cloud/go/services/certifier"
	certifier_test_service "magma/orc8r/cloud/go/services/certifier/test_init"
	certifier_test_utils "magma/orc8r/cloud/go/services/certifier/test_utils"
)

// Starts certifier & adds a Gw Identities for Gateways with hwGwIds
// Returns a list of corresponding Certificate Serial Numbers
func StartMockGwAccessControl(t *testing.T, hwGwIds []string) []string {

	// Start services
	certifier_test_service.StartTestService(t)

	result := make([]string, len(hwGwIds))

	for idx, hwId := range hwGwIds {
		// create and sign Gw's csr
		csrMsg, err := certifier_test_utils.CreateCSRForId(
			time.Duration(time.Hour*4), identity.NewGateway(hwId, "", ""))
		assert.NoError(t, err)

		certMsg, err := certifier.SignCSR(csrMsg)
		assert.NoError(t, err, "Failed to sign Gateway's CSR")
		// get cert sn from cert
		gwCert, err := x509.ParseCertificates(certMsg.CertDer)
		assert.NoError(t, err, "Failed to parse Gateway's cert")

		certSerialNum := security_cert.SerialToString(gwCert[0].SerialNumber)
		t.Logf("Test Gateway Certificate SN: %s", certSerialNum)

		result[idx] = certSerialNum
	}
	return result
}
