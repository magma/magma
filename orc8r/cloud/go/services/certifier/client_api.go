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

package certifier

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/clock"
	certifierprotos "magma/orc8r/cloud/go/services/certifier/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
)

const ServiceName = "CERTIFIER"

// Utility function to get a RPC connection to the certifier service
func getCertifierClient() (certifierprotos.CertifierClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		return nil, merrors.NewInitError(err, ServiceName)
	}

	return certifierprotos.NewCertifierClient(conn), err
}

// Get the certificate for the requested CA
func GetCACert(getCAReq *certifierprotos.GetCARequest) (*protos.CACert, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	ca, err := client.GetCA(context.Background(), getCAReq)
	if err != nil {
		glog.Errorf("Failed to get CA: %s", err)
		return nil, err
	}

	return ca, nil
}

// Return a signed certificate given CSR
func SignCSR(csr *protos.CSR) (*protos.Certificate, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	cert, err := client.SignAddCertificate(context.Background(), csr)
	if err != nil {
		glog.Errorf("Failed to sign CSR: %s", err)
		return nil, err
	}
	return cert, nil
}

// Add an existing Certificate & associate it with operator
func AddCertificate(oper *protos.Identity, certDer []byte) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.AddCertificate(
		context.Background(), &certifierprotos.AddCertRequest{Id: oper, CertDer: certDer})
	return err
}

// Get the CertificateInfo {Identity, NotAfter} of an SN
func GetIdentity(sn *protos.Certificate_SN) (*certifierprotos.CertificateInfo, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	certInfo, err := client.GetIdentity(context.Background(), sn)
	if err != nil {
		glog.Errorf("Failed to get identity with SN: %s, %s", sn.Sn, err)
		return nil, err
	}
	return certInfo, nil
}

// GetCertificateIdentity returns CertificateInfo of Certificate with the given
// Serial Number String. It's a simple wrapper for GetIdentity
func GetCertificateIdentity(serialNum string) (*certifierprotos.CertificateInfo, error) {
	return GetIdentity(&protos.Certificate_SN{Sn: serialNum})
}

// GetVerifiedCertificateIdentity returns CertificateInfo of Certificate with
// the given Serial Number String and verifies its validity
func GetVerifiedCertificateIdentity(serialNum string) (*protos.Identity, error) {
	certInfo, err := GetIdentity(&protos.Certificate_SN{Sn: serialNum})
	if err != nil {
		glog.Errorf("Lookup error '%s' for Cert SN: %s", err, serialNum)
		return nil, err
	}
	if certInfo == nil {
		err = fmt.Errorf("Missing Certificate Info for Cert SN: %s", serialNum)
		glog.Error(err)
		return nil, err
	}
	// Check if certificate time is not expired/not active yet
	err = VerifyDateRange(certInfo)
	if err != nil {
		glog.Errorf(
			"Certificate Validation Error '%s' for Cert SN: %s", err, serialNum)
		return nil, err
	}
	if certInfo.Id == nil {
		err = fmt.Errorf("Missing Identity for Cert SN: %s", serialNum)
		glog.Error(err)
		return nil, err
	}
	return certInfo.Id, nil
}

// Returns serial numbers of all registered certificates
func ListCertificates() ([]string, error) {
	client, err := getCertifierClient()
	if err != nil {
		return []string{}, err
	}
	slist, err := client.ListCertificates(context.Background(), &protos.Void{})
	if err != nil || slist == nil {
		return []string{}, err
	}
	return slist.Sns, err
}

// Finds & returns Serial Numbers of all Certificates associated with the
// given Identity
func FindCertificates(id *protos.Identity) ([]string, error) {
	client, err := getCertifierClient()
	if err != nil {
		return []string{}, err
	}
	slist, err := client.FindCertificates(context.Background(), id)
	if err != nil || slist == nil {
		return []string{}, err
	}
	return slist.Sns, err
}

// GetAll returns all Certificates Records
func GetAll() (map[string]*certifierprotos.CertificateInfo, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	certMap, err := client.GetAll(context.Background(), &protos.Void{})
	if err != nil || certMap == nil {
		return nil, err
	}
	return certMap.GetCertificates(), err
}

// Revoke Certificate and delete record of given SN
func RevokeCertificate(sn *protos.Certificate_SN) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}

	glog.V(2).Infof("Certifier: revoking certificate with SN: %s", sn.Sn)

	_, err = client.RevokeCertificate(context.Background(), sn)
	if err != nil {
		glog.Errorf("Failed to revoke certificate with SN: %s, %s", sn.Sn, err)
		return err
	}
	return nil
}

func RevokeCertificateSN(sn string) error {
	return RevokeCertificate(&protos.Certificate_SN{Sn: sn})
}

// Let certifier to remove expired certificates
func CollectGarbage() error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}

	_, err = client.CollectGarbage(context.Background(), &protos.Void{})
	if err != nil {
		glog.Errorf("Failed to collect garbage: %v", err)
		return err
	}
	return nil
}

type CertDateRange interface {
	GetNotBefore() *timestamp.Timestamp
	GetNotAfter() *timestamp.Timestamp
}

// Check if certificate time is not expired/not active yet
func VerifyDateRange(certInfo CertDateRange) error {
	tm := clock.Now()
	notBefore, _ := ptypes.Timestamp(certInfo.GetNotBefore())
	notAfter, _ := ptypes.Timestamp(certInfo.GetNotAfter())
	if tm.After(notAfter) {
		return errors.New("Expired")
	}
	if tm.Before(notBefore) {
		return errors.New("Not yet valid")
	}
	return nil
}
