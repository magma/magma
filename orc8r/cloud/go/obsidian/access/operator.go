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

package access

import (
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

// getOperator returns Identity of request's Operator (client).
// If either the request is missing TLS certificate headers or the certificate's
// SN is not found by Certifier or one of certificate & its identity checks fail
// - nil will be returned & the corresponding error logged
func getOperator(req *http.Request, decorate logDecorator) (*protos.Identity, error) {
	// Get Certificate SN header value
	// TBD: to optimize - use map directly
	csn := req.Header.Get(CLIENT_CERT_SN_KEY)
	if len(csn) == 0 {
		glog.V(1).Info(decorate("Missing REST client certificate"))
		return nil, fmt.Errorf("missing client certificate")
	}
	certInfo, err := certifier.GetCertificateIdentity(csn)
	if err != nil {
		if _, ok := err.(errors.ClientInitError); ok {
			glog.Error(decorate("Certificate SN '%s' lookup error '%s'", csn, err))
			return nil, err
		}
		glog.V(1).Info(decorate("Certificate SN '%s' lookup error '%s'", csn, err))
		return nil, fmt.Errorf("unknown client certificate SN: %s, err: %v", csn, err)
	}
	if certInfo == nil {
		glog.V(1).Info(decorate("No certificate info for SN: %s", csn))
		return nil, fmt.Errorf("unregistered client certificate, SN: %s", csn)
	}
	// Check if certificate time is not expired/not active yet
	err = certifier.VerifyDateRange(certInfo)
	if err != nil {
		glog.V(1).Info(decorate("Certificate validation error '%s' for SN: %s", err, csn))
		return nil, fmt.Errorf("certificate validation error: %s", err)
	}
	opId := certInfo.Id
	if opId == nil {
		glog.Error(decorate("Nil identity for certificate SN: %s", csn))
		return nil, fmt.Errorf("Internal server error (identity)")
	}
	// Check if it's operator identity
	if !identity.IsOperator(opId) {
		glog.V(1).Info(decorate("Identity '%s' of CSN '%s' is not an operator", opId.HashString(), csn))
		return nil, fmt.Errorf("identity must be for an operator")
	}

	return certInfo.Id, nil
}
