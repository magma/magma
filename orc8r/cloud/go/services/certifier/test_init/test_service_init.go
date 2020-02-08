/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/servicers"
	"magma/orc8r/cloud/go/services/certifier/storage"
	certifier_test_utils "magma/orc8r/cloud/go/services/certifier/test_utils"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func StartTestService(t *testing.T) {
	caMap := map[protos.CertType]*servicers.CAInfo{}

	bootstrapCert, bootstrapKey, err := certifier_test_utils.CreateSignedCertAndPrivKey(
		time.Duration(time.Hour * 24 * 10))
	if err != nil {
		t.Fatalf("Failed to create bootstrap certifier certificate: %s", err)
	} else {
		caMap[protos.CertType_DEFAULT] = &servicers.CAInfo{bootstrapCert, bootstrapKey}
	}

	vpnCert, vpnKey, err := certifier_test_utils.CreateSignedCertAndPrivKey(
		time.Duration(time.Hour * 24 * 10))
	if err != nil {
		t.Fatalf("Failed to create VPN certifier certificate: %s", err)
	} else {
		caMap[protos.CertType_VPN] = &servicers.CAInfo{vpnCert, vpnKey}
	}
	ds := test_utils.GetMockDatastoreInstance()
	certStore := storage.NewCertifierDatastore(ds)
	certServer, err := servicers.NewCertifierServer(certStore, caMap)
	if err != nil {
		t.Fatalf("Failed to create certifier server: %s", err)
	}
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, certifier.ServiceName)
	certprotos.RegisterCertifierServer(
		srv.GrpcServer,
		certServer,
	)
	go srv.RunTest(lis)
}
