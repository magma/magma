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
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"magma/orc8r/cloud/go/clock"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const ServiceName = "CERTIFIER"

// Utility function to get a RPC connection to the certifier service
func getCertifierClient() (certprotos.CertifierClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		return nil, merrors.NewInitError(err, ServiceName)
	}

	return certprotos.NewCertifierClient(conn), err
}

// Get the certificate for the requested CA
func GetCACert(ctx context.Context, getCAReq *certprotos.GetCARequest) (*protos.CACert, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	ca, err := client.GetCA(ctx, getCAReq)
	if err != nil {
		glog.Errorf("Failed to get CA: %s", err)
		return nil, err
	}

	return ca, nil
}

// Return a signed certificate given CSR
func SignCSR(ctx context.Context, csr *protos.CSR) (*protos.Certificate, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	cert, err := client.SignAddCertificate(ctx, csr)
	if err != nil {
		glog.Errorf("Failed to sign CSR: %s", err)
		return nil, err
	}
	return cert, nil
}

// Add an existing Certificate & associate it with operator
func AddCertificate(ctx context.Context, oper *protos.Identity, certDer []byte) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.AddCertificate(ctx, &certprotos.AddCertRequest{Id: oper, CertDer: certDer})
	return err
}

// Get the CertificateInfo {Identity, NotAfter} of an SN
func GetIdentity(ctx context.Context, sn *protos.Certificate_SN) (*certprotos.CertificateInfo, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}

	certInfo, err := client.GetIdentity(ctx, sn)
	if err != nil {
		glog.Errorf("Failed to get identity with SN: %s, %s", sn.Sn, err)
		return nil, err
	}
	return certInfo, nil
}

// GetCertificateIdentity returns CertificateInfo of Certificate with the given
// Serial Number String. It's a simple wrapper for GetIdentity
func GetCertificateIdentity(ctx context.Context, serialNum string) (*certprotos.CertificateInfo, error) {
	return GetIdentity(ctx, &protos.Certificate_SN{Sn: serialNum})
}

// GetVerifiedCertificateIdentity returns CertificateInfo of Certificate with
// the given Serial Number String and verifies its validity
func GetVerifiedCertificateIdentity(ctx context.Context, serialNum string) (*protos.Identity, error) {
	certInfo, err := GetIdentity(ctx, &protos.Certificate_SN{Sn: serialNum})
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
func ListCertificates(ctx context.Context) ([]string, error) {
	client, err := getCertifierClient()
	if err != nil {
		return []string{}, err
	}
	slist, err := client.ListCertificates(ctx, &protos.Void{})
	if err != nil || slist == nil {
		return []string{}, err
	}
	return slist.Sns, err
}

// Finds & returns Serial Numbers of all Certificates associated with the
// given Identity
func FindCertificates(ctx context.Context, id *protos.Identity) ([]string, error) {
	client, err := getCertifierClient()
	if err != nil {
		return []string{}, err
	}
	slist, err := client.FindCertificates(ctx, id)
	if err != nil || slist == nil {
		return []string{}, err
	}
	return slist.Sns, err
}

// GetAll returns all Certificates Records
func GetAll(ctx context.Context) (map[string]*certprotos.CertificateInfo, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	certMap, err := client.GetAll(ctx, &protos.Void{})
	if err != nil || certMap == nil {
		return nil, err
	}
	return certMap.GetCertificates(), err
}

// Revoke Certificate and delete record of given SN
func RevokeCertificate(ctx context.Context, sn *protos.Certificate_SN) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}

	glog.V(2).Infof("Certifier: revoking certificate with SN: %s", sn.Sn)

	_, err = client.RevokeCertificate(ctx, sn)
	if err != nil {
		glog.Errorf("Failed to revoke certificate with SN: %s, %s", sn.Sn, err)
		return err
	}
	return nil
}

func RevokeCertificateSN(ctx context.Context, sn string) error {
	return RevokeCertificate(ctx, &protos.Certificate_SN{Sn: sn})
}

// Let certifier to remove expired certificates
func CollectGarbage(ctx context.Context) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}

	_, err = client.CollectGarbage(ctx, &protos.Void{})
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

// GetPolicyDecision makes a policy decision when a user attempts to access a resource
func GetPolicyDecision(ctx context.Context, getPDReq *certprotos.GetPolicyDecisionRequest) (*certprotos.GetPolicyDecisionResponse, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	pd, err := client.GetPolicyDecision(ctx, getPDReq)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

// CreateUser creates a new user with the specified password and policy
func CreateUser(ctx context.Context, user *certprotos.User) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.CreateUser(ctx, &certprotos.CreateUserRequest{User: user})
	return err
}

// ListUsers lists all users and their tokens in the database
func ListUsers(ctx context.Context) ([]*certprotos.User, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	users, err := client.ListUsers(ctx, &certprotos.ListUsersRequest{})
	return users.Users, err
}

func GetUser(ctx context.Context, username string) (*certprotos.User, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	user, err := client.GetUser(ctx, &certprotos.GetUserRequest{User: &certprotos.User{Username: username}})
	return user.User, err
}

func UpdateUser(ctx context.Context, user *certprotos.User) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.UpdateUser(ctx, &certprotos.UpdateUserRequest{User: user})
	return err
}

func DeleteUser(ctx context.Context, user *certprotos.User) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteUser(ctx, &certprotos.DeleteUserRequest{User: user})
	return err
}

func ListUserTokens(ctx context.Context, user *certprotos.User) (*certprotos.ListUserTokensResponse, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	tokens, err := client.ListUserTokens(ctx, &certprotos.ListUserTokensRequest{User: user})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func AddUserToken(ctx context.Context, req *certprotos.AddUserTokenRequest) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.AddUserToken(ctx, req)
	return err
}

func DeleteUserToken(ctx context.Context, req *certprotos.DeleteUserTokenRequest) error {
	client, err := getCertifierClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteUserToken(ctx, req)
	return err
}

func Login(ctx context.Context, req *certprotos.LoginRequest) (*certprotos.LoginResponse, error) {
	client, err := getCertifierClient()
	if err != nil {
		return nil, err
	}
	res, err := client.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
