/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test_utils

import (
	"crypto/x509"
	"sync"
	"testing"
	"time"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/lib/go/security/cert"
	"magma/orc8r/lib/go/security/csr"

	"github.com/stretchr/testify/assert"
)

var (
	once sync.Once
)

// StartMockGwAccessControl starts certifier & adds a Gw Identities for
// Gateways with hwGwIds.
// Returns a list of corresponding Certificate Serial Numbers.
func StartMockGwAccessControl(t *testing.T, hwGwIds []string) []string {
	// TODO(hcgatewood): the first-best solution here is having callers start
	// the certifier test service on their own. But this sync.Once solution is a
	// serviceable stopgap for now.
	once.Do(func() { test_init.StartTestService(t) })

	result := make([]string, len(hwGwIds))
	for idx, hwId := range hwGwIds {
		// create and sign Gw's csr
		csrMsg, err := csr.CreateCSRForId(
			time.Hour*4, identity.NewGateway(hwId, "", ""))
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
