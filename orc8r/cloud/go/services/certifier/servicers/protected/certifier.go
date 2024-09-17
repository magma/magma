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

package servicers

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/constants"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/security/cert"
)

var (
	NumTrialsForSn      int
	CollectGarbageAfter time.Duration // remove cert if expired for certain amount of time
)

func init() {
	NumTrialsForSn = 1
	CollectGarbageAfter = time.Hour * 24
}

type CAInfo struct {
	Cert    *x509.Certificate
	PrivKey interface{}
}

type CertifierServer struct {
	store storage.CertifierStorage
	CAs   map[protos.CertType]*CAInfo
}

func NewCertifierServer(store storage.CertifierStorage, CAs map[protos.CertType]*CAInfo) (srv *CertifierServer, err error) {
	srv = new(CertifierServer)
	srv.store = store
	if CAs == nil {
		return nil, fmt.Errorf("CA info not provided to certifier")
	}
	if len(CAs) == 0 {
		return nil, fmt.Errorf("No Certificates are provided to certifier")
	}
	srv.CAs = CAs
	return srv, nil
}

func (srv *CertifierServer) GetCA(ctx context.Context, getCAReqMsg *certprotos.GetCARequest) (*protos.CACert, error) {
	if getCAReqMsg == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid CA request")
	}

	ca, ok := srv.CAs[getCAReqMsg.CertType]
	if !ok {
		return nil, fmt.Errorf("no CA found for given CA type: %s", getCAReqMsg.CertType.String())
	}

	caCertMsg := &protos.CACert{Cert: ca.Cert.Raw}

	return caCertMsg, nil
}

func (srv *CertifierServer) SignAddCertificate(ctx context.Context, csrMsg *protos.CSR) (*protos.Certificate, error) {

	sn, err := generateSerialNumber(srv.store)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Error generating serial number: %s", err)
	}

	csr, err := parseAndCheckCSR(csrMsg.CsrDer)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Error parsing CSR: %s", err)
	}

	err = checkOrOverwriteCN(csr, csrMsg)
	if err != nil {
		return nil, err
	}

	validTime, err := ptypes.Duration(csrMsg.ValidTime)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Invalid requested certificate duration: %s", err)
	}

	certDER, notBefore, notAfter, err := srv.signCSR(csr, sn, csrMsg.CertType, validTime)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Error signing CSR: %s", err)
	}

	notBeforeProto, _ := ptypes.TimestampProto(notBefore)
	notAfterProto, _ := ptypes.TimestampProto(notAfter)

	// create CertificateInfo
	certInfo := &certprotos.CertificateInfo{
		Id:        csrMsg.Id,
		CertType:  csrMsg.CertType,
		NotBefore: notBeforeProto,
		NotAfter:  notAfterProto,
	}
	// add to table
	snString := cert.SerialToString(sn)
	// Ensure serial number is not the orc8r client reserved SN
	if snString == registry.ORC8R_CLIENT_CERT_VALUE {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Serial Number")
	}
	err = srv.store.PutCertInfo(snString, certInfo)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Error adding CertificateInfo: %s", err)
	}

	// create Certificate
	certMsg := protos.Certificate{
		Sn:        &protos.Certificate_SN{Sn: snString},
		NotBefore: notBeforeProto,
		NotAfter:  notAfterProto,
		CertDer:   certDER,
	}
	return &certMsg, nil
}

func (srv *CertifierServer) GetIdentity(
	ctx context.Context, snMsg *protos.Certificate_SN) (*certprotos.CertificateInfo, error) {

	var certSN string
	if snMsg != nil {
		certSN = strings.TrimLeft(snMsg.Sn, "0")
	}
	certInfo, err := srv.store.GetCertInfo(certSN)
	if err != nil {
		return &certprotos.CertificateInfo{}, status.Errorf(
			codes.NotFound, "Certificate with serial number '%s' is not found", certSN)
	}

	// check timestamp
	notBefore, _ := ptypes.Timestamp(certInfo.NotBefore)
	notAfter, _ := ptypes.Timestamp(certInfo.NotAfter)
	now := clock.Now().UTC()
	if now.After(notAfter) {
		return &certprotos.CertificateInfo{}, status.Errorf(codes.OutOfRange,
			"Certificate with serial number '%s' has expired", certSN)
	}
	if now.Before(notBefore) {
		return &certprotos.CertificateInfo{}, status.Errorf(codes.OutOfRange,
			"Certificate with serial number '%s' is not yet valid", certSN)
	}
	return certInfo, nil
}

func (srv *CertifierServer) RevokeCertificate(
	ctx context.Context, snMsg *protos.Certificate_SN) (*protos.Void, error) {

	var certSN string
	if snMsg != nil {
		certSN = strings.TrimLeft(snMsg.Sn, "0")
	}
	_, err := srv.store.GetCertInfo(certSN)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot find certificate with SN: %s", certSN)
	}
	err = srv.store.DeleteCertInfo(certSN)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to delete certificate: %s", err)
	}
	return &protos.Void{}, nil
}

func (srv *CertifierServer) AddCertificate(ctx context.Context, req *certprotos.AddCertRequest) (*protos.Void, error) {

	res := &protos.Void{}
	x509Cert, err := x509.ParseCertificate(req.CertDer)
	if err != nil {
		return res,
			status.Errorf(codes.InvalidArgument, "DER Parse Error: %s", err)
	}
	if x509Cert.SerialNumber == nil {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Serial Number")
	}
	snStr := cert.SerialToString(x509Cert.SerialNumber)
	// Ensure serial number is not the orc8r client reserved SN
	if snStr == registry.ORC8R_CLIENT_CERT_VALUE {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Serial Number")
	}
	// Verify that the certificate is signed by our CA
	if err = srv.verifyCert(x509Cert, req.CertType); err != nil {
		return res, status.Errorf(
			codes.InvalidArgument, "%s for Certificate SN %s", err, snStr)
	}
	// Check if a certificate with the same SN is already there
	_, err = srv.store.GetCertInfo(snStr)
	if err == nil {
		return res, status.Errorf(
			codes.AlreadyExists, "Certificate SN %s already exists", snStr)
	}
	// create CertificateInfo
	notBeforeProto, _ := ptypes.TimestampProto(x509Cert.NotBefore)
	notAfterProto, _ := ptypes.TimestampProto(x509Cert.NotAfter)
	certInfo := &certprotos.CertificateInfo{
		Id:        req.Id,
		CertType:  req.CertType,
		NotBefore: notBeforeProto,
		NotAfter:  notAfterProto,
	}
	// add to table
	err = srv.store.PutCertInfo(snStr, certInfo)
	if err != nil {
		return res,
			status.Errorf(codes.Internal, "Error adding CertificateInfo: %s", err)
	}
	return res, nil
}

// Finds & returns Serial Numbers of all Certificates associated with the
// given Identity
func (srv *CertifierServer) FindCertificates(ctx context.Context, id *protos.Identity) (*certprotos.SerialNumbers, error) {

	res := &certprotos.SerialNumbers{}
	if id != nil {
		idKey := id.HashString()
		snList, err := srv.ListCertificates(ctx, &protos.Void{})
		if err != nil {
			return res, err
		}
		for _, sn := range snList.Sns {
			certInfo, err := srv.getCertInfo(sn)
			if err != nil {
				return res, err
			}
			if certInfo != nil && certInfo.Id.HashString() == idKey {
				res.Sns = append(res.Sns, sn)
			}
		}
	}
	return res, nil
}

// Returns serial numbers of all certificates in the table
func (srv *CertifierServer) ListCertificates(ctx context.Context, void *protos.Void) (*certprotos.SerialNumbers, error) {
	res := &certprotos.SerialNumbers{}
	snList, err := srv.store.ListSerialNumbers()
	if err != nil {
		return res, status.Errorf(
			codes.Internal, "Failed to get certificate serial numbers: %s", err)
	}
	res.Sns = snList
	return res, nil
}

// GetAll returns all Certificates Records
func (srv *CertifierServer) GetAll(context.Context, *protos.Void) (*certprotos.CertificateInfoMap, error) {
	res := &certprotos.CertificateInfoMap{Certificates: map[string]*certprotos.CertificateInfo{}}
	certInfos, err := srv.store.GetAllCertInfo()
	if err != nil {
		return res, status.Errorf(codes.Internal, "Failed to get all certificates: %v", err)
	}
	res.Certificates = certInfos
	return res, nil
}

func (srv *CertifierServer) CollectGarbage(ctx context.Context, void *protos.Void) (*protos.Void, error) {
	count, err := srv.CollectGarbageImpl(ctx)
	glog.Infof("purged %d expired certificates", count)
	return &protos.Void{}, err
}

func (srv *CertifierServer) CollectGarbageImpl(ctx context.Context) (int, error) {
	snList, err := srv.ListCertificates(ctx, &protos.Void{})
	if err != nil {
		return 0, err
	}
	errs := &multierror.Error{}
	count := 0
	for _, sn := range snList.Sns {
		certInfo, err := srv.getCertInfo(sn)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("'%s' get info error: %v", sn, err))
		}
		notAfter, _ := ptypes.Timestamp(certInfo.NotAfter)
		notAfter = notAfter.Add(CollectGarbageAfter)
		if time.Now().UTC().After(notAfter) {
			err = srv.store.DeleteCertInfo(sn)
			if err != nil {
				errs = multierror.Append(errs, fmt.Errorf("'%s' delete error: %v", sn, err))
			} else {
				count += 1
			}
		}
	}
	if errs.ErrorOrNil() != nil {
		glog.Errorf("Failed to delete certificate[s]: %v", errs)
		return count, status.Error(codes.Internal, errs.Error())
	}
	return count, nil
}

// GetPolicyDecision makes a policy decision when a user attempts to access a resource.
// For conflicting policy decisions from multiple tokens (e.g. one policy is ALLOW and the other DENY), the DENY effect
// will take precedent.
// For resources that do not have any policies addressing it, the policy decision defaults to DENY as well.
func (srv *CertifierServer) GetPolicyDecision(ctx context.Context, getPDReq *certprotos.GetPolicyDecisionRequest) (*certprotos.GetPolicyDecisionResponse, error) {
	if err := certifier.ValidateToken(getPDReq.Token); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	username := getPDReq.Username
	user, err := srv.store.GetUser(username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user %s from database: %v", username, err)
	}

	err = isTokenWithUser(getPDReq.Token, user.Tokens)
	if err != nil {
		return nil, err
	}

	decision, err := srv.getPolicyDecisionFromTokenMany(ctx, user.Tokens, getPDReq)
	if err != nil {
		return nil, err
	}

	return decision, nil
}

// CreateUser creates a new user with the specified password and policy
func (srv *CertifierServer) CreateUser(ctx context.Context, req *certprotos.CreateUserRequest) (*certprotos.CreateUserResponse, error) {
	user, _ := srv.store.GetUser(req.User.Username)
	if user != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(req.User.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password: %v", err)
	}

	user = &certprotos.User{
		Username: req.User.Username,
		Password: hashedPassword,
	}
	err = srv.store.PutUser(user.Username, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store user while creating user")
	}
	return &certprotos.CreateUserResponse{}, nil
}

func (srv *CertifierServer) ListUsers(ctx context.Context, req *certprotos.ListUsersRequest) (*certprotos.ListUsersResponse, error) {
	users, err := srv.store.ListUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users from db")
	}
	return &certprotos.ListUsersResponse{Users: users}, nil
}

func (srv *CertifierServer) GetUser(ctx context.Context, req *certprotos.GetUserRequest) (*certprotos.GetUserResponse, error) {
	user, err := srv.store.GetUser(req.User.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user %s", req.User.Username)
	}
	return &certprotos.GetUserResponse{User: user}, nil
}

func (srv *CertifierServer) UpdateUser(ctx context.Context, req *certprotos.UpdateUserRequest) (*certprotos.UpdateUserResponse, error) {
	user, err := srv.store.GetUser(req.User.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(req.User.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password: %v", err)
	}

	newUser := &certprotos.User{
		Username: req.User.Username,
		Password: hashedPassword,
		Tokens:   req.User.Tokens,
	}
	err = srv.store.PutUser(user.Username, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating user")
	}
	return &certprotos.UpdateUserResponse{}, nil
}

func (srv *CertifierServer) DeleteUser(ctx context.Context, req *certprotos.DeleteUserRequest) (*certprotos.DeleteUserResponse, error) {
	userToDelete, err := srv.store.GetUser(req.User.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "error getting user")
	}
	if userToDelete != nil && userToDelete.Tokens != nil {
		for _, token := range userToDelete.Tokens.Tokens {
			err = srv.store.DeletePolicy(token)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "error deleting token")
			}
		}
	}
	err = srv.store.DeleteUser(req.User.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting user: %v", err)
	}

	return &certprotos.DeleteUserResponse{}, nil
}

func (srv *CertifierServer) ListUserTokens(ctx context.Context, req *certprotos.ListUserTokensRequest) (*certprotos.ListUserTokensResponse, error) {
	user, err := srv.store.GetUser(req.User.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting user for listing user tokens")
	}
	policies, err := srv.getPolicyFromTokenMany(user.Tokens)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting policy from token list: %v", err)
	}
	// Add empty policy list if nil for marshaling purposes
	if policies == nil {
		policies = new(certprotos.ListUserTokensResponse)
	}
	return policies, nil
}

func (srv *CertifierServer) AddUserToken(ctx context.Context, req *certprotos.AddUserTokenRequest) (*certprotos.AddUserTokenResponse, error) {
	user, err := srv.store.GetUser(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting user for adding token: %v", err)
	}
	token, err := certifier.GenerateToken(certifier.Personal)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating token: %v", err)
	}

	if user.Tokens == nil {
		user.Tokens = &certprotos.TokenList{Tokens: []string{token}}
	} else {
		user.Tokens.Tokens = append(user.Tokens.Tokens, token)
	}
	newUser := &certprotos.User{
		Username: user.Username,
		Password: user.Password,
		Tokens:   &certprotos.TokenList{Tokens: user.Tokens.Tokens},
	}
	err = srv.store.PutUser(user.Username, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error putting user: %v", err)
	}

	policy := &certprotos.PolicyList{
		Token:    token,
		Policies: req.Policies,
	}
	err = srv.store.PutPolicy(token, policy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error putting policy: %v", err)
	}
	return &certprotos.AddUserTokenResponse{}, nil
}

func (srv *CertifierServer) DeleteUserToken(ctx context.Context, req *certprotos.DeleteUserTokenRequest) (*certprotos.DeleteUserTokenResponse, error) {
	user, err := srv.store.GetUser(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error getting user to delete token")
	}
	newTokenList, err := srv.deleteTokenFromUser(user.Tokens, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting token from user: %v", err)
	}
	newUser := &certprotos.User{
		Username: user.Username,
		Password: user.Password,
		Tokens:   newTokenList,
	}
	err = srv.store.PutUser(user.Username, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating user while deleting token: %v", err)
	}
	err = srv.store.DeletePolicy(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting policy: %v", err)
	}
	return &certprotos.DeleteUserTokenResponse{}, nil
}

func (srv *CertifierServer) Login(ctx context.Context, req *certprotos.LoginRequest) (*certprotos.LoginResponse, error) {
	userRes, err := srv.GetUser(ctx, &certprotos.GetUserRequest{User: req.User})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	hashedPassword := userRes.User.Password
	err = bcrypt.CompareHashAndPassword(hashedPassword, req.User.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "wrong password")
	}
	listTokensRes, err := srv.ListUserTokens(ctx, &certprotos.ListUserTokensRequest{User: userRes.User})
	if err != nil {
		return nil, err
	}
	return &certprotos.LoginResponse{PolicyLists: listTokensRes.PolicyLists}, nil
}

func (srv *CertifierServer) deleteTokenFromUser(tokenList *certprotos.TokenList, reqToken string) (*certprotos.TokenList, error) {
	remove := -1
	for idx, token := range tokenList.Tokens {
		if token == reqToken {
			remove = idx
		}
	}
	newTokenList := append(tokenList.Tokens[:remove], tokenList.Tokens[remove+1:]...)
	return &certprotos.TokenList{Tokens: newTokenList}, nil
}

func generateSerialNumber(store storage.CertifierStorage) (sn *big.Int, err error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)

	for i := 0; i < NumTrialsForSn; i++ {
		sn, err = rand.Int(rand.Reader, limit)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate serial number: %s", err)
		}
		_, err := store.GetCertInfo(cert.SerialToString(sn))
		if err != nil {
			return sn, nil
		}
	}
	return nil, fmt.Errorf(
		"Failed to genearte serial number after %d trials.", NumTrialsForSn)
}

func parseAndCheckCSR(csrDER []byte) (*x509.CertificateRequest, error) {
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate request: %s", err)
	}

	err = csr.CheckSignature()
	if err != nil {
		return nil, fmt.Errorf("Failed to check certificate request signature: %s", err)
	}
	return csr, err
}

func (srv *CertifierServer) signCSR(
	csr *x509.CertificateRequest,
	sn *big.Int,
	certType protos.CertType,
	validTime time.Duration,
) ([]byte, time.Time, time.Time, error) {

	if srv.CAs == nil {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("CAInfo not found")
	}
	ca, ok := srv.CAs[certType]
	if !ok {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("No CA found for given cert type: %s", certType.String())
	}
	signingCert := ca.Cert
	signingKey := ca.PrivKey

	now := clock.Now().UTC()
	// Provide a cert from an hour ago to account for clock skews
	notBefore := now.Add(-1 * time.Hour)
	notAfter := now.Add(validTime)
	if notAfter.After(signingCert.NotAfter) {
		glog.Warningln("The requested time is longer than signing certificate valid time.")
		notAfter = signingCert.NotAfter
	}
	template := x509.Certificate{
		SerialNumber:          sn,
		Subject:               csr.Subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	clientCertDER, err := x509.CreateCertificate(
		rand.Reader, &template, signingCert, csr.PublicKey, signingKey)
	if err != nil {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("Failed to sign csr: %s", err)
	}

	return clientCertDER, notBefore, notAfter, nil
}

func checkOrOverwriteCN(csr *x509.CertificateRequest, csrMsg *protos.CSR) error {
	id := csrMsg.Id
	idCn := id.ToCommonName()
	if idCn == nil {
		return nil
	}
	if len(csr.Subject.CommonName) == 0 {
		csr.Subject.CommonName = *idCn
		return nil
	}

	if csr.Subject.CommonName != *idCn {
		return status.Errorf(
			codes.Aborted,
			"CN from CSR (%s) and CN in Identity (%s) do not match", csr.Subject.CommonName, *idCn)
	}

	if csrMsg.CertType == protos.CertType_VPN && identity.IsGateway(id) {
		// Use networkID & logicalID to identify the vpn client instead of hwID
		gw := id.GetGateway()
		csr.Subject.CommonName = gw.GetLogicalId()
	}

	return nil
}

func (srv *CertifierServer) getCertInfo(sn string) (*certprotos.CertificateInfo, error) {
	certInfo, err := srv.store.GetCertInfo(sn)
	if err != nil {
		return &certprotos.CertificateInfo{}, status.Errorf(codes.NotFound, "Failed to load certificate: %s", err)
	}
	return certInfo, nil
}

// Verify that the certificate is signed by our CA
func (srv *CertifierServer) verifyCert(clientCert *x509.Certificate, certType protos.CertType) error {
	// Check if CAInfo / cert exists for requested cert type
	if srv.CAs == nil {
		return fmt.Errorf("CAInfo not found")
	}
	ca, ok := srv.CAs[certType]
	if !ok {
		return fmt.Errorf("No CA found for given cert type: %s", certType.String())
	}

	caPool := x509.NewCertPool()
	caPool.AddCert(ca.Cert) // Use appropriate cert to check against
	opts := x509.VerifyOptions{
		Roots:         caPool,
		Intermediates: x509.NewCertPool(),
		// Make sure client cert has ExtKeyUsageClientAuth
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	if _, err := clientCert.Verify(opts); err != nil {
		return fmt.Errorf("Certificate Verification Failure: %s", err)
	}
	return nil
}

func (srv *CertifierServer) getPolicyFromTokenMany(tokens *certprotos.TokenList) (*certprotos.ListUserTokensResponse, error) {
	if tokens == nil {
		return nil, nil
	}
	ret := make([]*certprotos.PolicyList, len(tokens.Tokens))
	for i, token := range tokens.Tokens {
		policy, err := srv.getPolicyFromToken(token)
		if err != nil {
			return nil, err
		}
		ret[i] = policy
	}
	return &certprotos.ListUserTokensResponse{PolicyLists: ret}, nil
}

func (srv *CertifierServer) getPolicyFromToken(token string) (*certprotos.PolicyList, error) {
	policy, err := srv.store.GetPolicy(token)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (srv *CertifierServer) getPolicyDecisionFromTokenMany(ctx context.Context, tokens *certprotos.TokenList, getPDReq *certprotos.GetPolicyDecisionRequest) (*certprotos.GetPolicyDecisionResponse, error) {
	resource := getPDReq.Request

	finalEffect := certprotos.Effect_DENY
	for _, t := range tokens.Tokens {
		effect, err := srv.getPolicyDecisionFromToken(ctx, t, resource)
		// Return DENY if there are any unexpected errors
		if err != nil {
			return &certprotos.GetPolicyDecisionResponse{Effect: certprotos.Effect_DENY}, nil
		}
		switch effect {
		case certprotos.Effect_DENY:
			// Return early if there is a DENY in any of the permissions
			return &certprotos.GetPolicyDecisionResponse{Effect: certprotos.Effect_DENY}, nil
		case certprotos.Effect_ALLOW:
			finalEffect = certprotos.Effect_ALLOW
		default:
			continue
		}
	}
	// Return DENY if the policy unknown
	return &certprotos.GetPolicyDecisionResponse{Effect: finalEffect}, nil
}

func (srv *CertifierServer) getPolicyDecisionFromToken(ctx context.Context, token string, req *certprotos.Request) (certprotos.Effect, error) {
	policyList, err := srv.store.GetPolicy(token)
	if err != nil {
		return certprotos.Effect_DENY, status.Errorf(codes.Internal, "failed to get policyList from db %v", err)
	}
	effect := certprotos.Effect_UNKNOWN

	// Networks are registered with tenants, hence any tenant scoped policies
	// have additional policies for their networks.
	tenantNetworkResource, err := getTenantPolicyNetworkResourceMany(ctx, policyList)
	if err != nil {
		return effect, nil
	}
	policyList.Policies = append(policyList.Policies, tenantNetworkResource...)

	for _, policy := range policyList.Policies {
		actionEffect := getActionAuthorization(req, policy)
		resourceEffect := getResourceAuthorization(req, policy)
		// The effects from both action and resource should match
		// for the effect of the policyList to apply to the request.
		notBothUnknown := actionEffect != certprotos.Effect_UNKNOWN && resourceEffect != certprotos.Effect_UNKNOWN
		equal := actionEffect == resourceEffect
		if notBothUnknown && equal {
			effect = actionEffect
		}

		// actionEffect only checks the read/write permission of a requested resource,
		// but it may not apply to that specific resource (resourceEffect ensures
		// that the policy applies to that resource), thus both action and
		// resource effect need deny in order for the final effect to be DENY.
		if actionEffect == certprotos.Effect_DENY && resourceEffect == certprotos.Effect_DENY {
			return certprotos.Effect_DENY, nil
		}
	}

	return effect, nil
}

// getTenantPolicyNetworkResourceMany builds a set of network policies for each tenant policy
func getTenantPolicyNetworkResourceMany(ctx context.Context, policyList *certprotos.PolicyList) ([]*certprotos.Policy, error) {
	var networkResources []*certprotos.Policy
	for _, policy := range policyList.Policies {
		if t := policy.GetTenant(); t != nil {
			networks, err := getTenantPolicyNetworkResource(ctx, t)
			if err != nil {
				return nil, err
			}
			resource := &certprotos.Policy{
				Effect:   policy.Effect,
				Action:   policy.Action,
				Resource: &certprotos.Policy_Network{Network: &certprotos.NetworkResource{Networks: networks}},
			}
			networkResources = append(networkResources, resource)
		}
	}
	return networkResources, nil
}

// getTenantPolicyNetworkResource retrieves all the networks registered with the tenant
func getTenantPolicyNetworkResource(ctx context.Context, t *certprotos.TenantResource) ([]string, error) {
	var networks []string
	for _, i := range t.Tenants {
		tenant, err := tenants.GetTenant(ctx, i)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
		networks = append(networks, tenant.Networks...)
	}

	return networks, nil
}

func getRequestType(req *certprotos.Request) constants.ResourceType {
	reqType := constants.Path
	switch req.ResourceId.(type) {
	case *certprotos.Request_NetworkId:
		reqType = constants.NetworkID
	case *certprotos.Request_TenantId:
		reqType = constants.TenantID
	}
	return reqType
}

func getPolicyType(policy *certprotos.Policy) constants.ResourceType {
	policyType := constants.Path
	switch policy.Resource.(type) {
	case *certprotos.Policy_Network:
		policyType = constants.NetworkID
	case *certprotos.Policy_Tenant:
		policyType = constants.TenantID
	}
	return policyType
}

// doResourceTypesMatch checks if the policy type applies to the requested type
// for network and tenant resource types.
// Every request has a path but not necessarily network/tenant ids, so path policy
// types are always applied to every request, and they are exempted from being matched by type.
// Returns true if the resource type is either network or tenant and the resource types match
// Returns false otherwise
func doResourceTypesMatch(req *certprotos.Request, policy *certprotos.Policy) bool {
	reqType := getRequestType(req)
	policyType := getPolicyType(policy)
	return req.GetResourceId() == nil || policyType == constants.Path || reqType == policyType
}

// getActionAuthorization checks if the requested action to read/write is authorized by the policy
// Returns the effect if the policy allows/denies write access regardless of requested action,
// or if policy allows/denies read access and the requested action is read
// Otherwise, returns an unknown effect, which implies that the policy does not apply to this requested resource
func getActionAuthorization(req *certprotos.Request, policy *certprotos.Policy) certprotos.Effect {
	// The policy does not handle this case, so return UNKNOWN.
	if !doResourceTypesMatch(req, policy) {
		return certprotos.Effect_UNKNOWN
	}
	if policy.Action == certprotos.Action_WRITE {
		return policy.Effect
	}
	if policy.Action == certprotos.Action_READ && req.Action == certprotos.Action_READ {
		return policy.Effect
	}
	return certprotos.Effect_UNKNOWN
}

// getResourceAuthorization checks if the requested resource (either a path, network, or tenant) is authorized by the policy
// Returns the effect if the policy's path matches the requested path and if the policy's tenant/network ID matches that
// of the requested when applicable
// Otherwise, returns an unknown effect, which implies that the policy does not apply to this requested resource
func getResourceAuthorization(req *certprotos.Request, policy *certprotos.Policy) certprotos.Effect {
	// The policy does not handle this case, so return UNKNOWN.
	if !doResourceTypesMatch(req, policy) {
		return certprotos.Effect_UNKNOWN
	}

	finalEffect := certprotos.Effect_UNKNOWN
	if path := policy.GetPath(); path != nil {
		if ok, _ := doublestar.Match(path.Path, req.Resource); ok {
			finalEffect = policy.Effect
		}
	}

	reqType := getRequestType(req)
	policyType := getPolicyType(policy)
	if reqType == constants.NetworkID && policyType == constants.NetworkID {
		reqID := req.GetNetworkId()
		policyIDs := policy.GetNetwork().Networks
		for _, policyID := range policyIDs {
			if finalEffect != certprotos.Effect_DENY && reqID == policyID {
				return policy.Effect
			}
		}
	}

	if reqType == constants.TenantID && policyType == constants.TenantID {
		reqID := req.GetTenantId()
		policyIDs := policy.GetTenant().Tenants
		for _, policyID := range policyIDs {
			if finalEffect != certprotos.Effect_DENY && reqID == policyID {
				return policy.Effect
			}
		}
	}
	return finalEffect
}

func isTokenWithUser(token string, tokenList *certprotos.TokenList) error {
	flag := false
	for _, t := range tokenList.Tokens {
		if t == token {
			flag = true
		}
	}
	if !flag {
		return status.Errorf(codes.PermissionDenied, "token is not registered with user")
	}
	return nil
}
