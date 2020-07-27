/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	certifierTestUtils "magma/orc8r/lib/go/security/csr"
)

func StartTestService(t *testing.T) {
	caMap := map[protos.CertType]*servicers.CAInfo{}

	bootstrapCert, bootstrapKey, err := certifierTestUtils.CreateSignedCertAndPrivKey(time.Hour * 24 * 10)
	if err != nil {
		t.Fatalf("Failed to create bootstrap certifier certificate: %s", err)
	} else {
		caMap[protos.CertType_DEFAULT] = &servicers.CAInfo{Cert: bootstrapCert, PrivKey: bootstrapKey}
	}

	vpnCert, vpnKey, err := certifierTestUtils.CreateSignedCertAndPrivKey(time.Hour * 24 * 10)
	if err != nil {
		t.Fatalf("Failed to create VPN certifier certificate: %s", err)
	} else {
		caMap[protos.CertType_VPN] = &servicers.CAInfo{Cert: vpnCert, PrivKey: vpnKey}
	}
	store := test_utils.NewSQLBlobstore(t, storage.CertifierTableBlobstore)
	certStore := storage.NewCertifierBlobstore(store)
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
