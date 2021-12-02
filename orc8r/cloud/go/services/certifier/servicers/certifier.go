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
	"net/http"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/identity"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/cert"
	unarylib "magma/orc8r/lib/go/service/middleware/unary"
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
	if snString == unarylib.ORC8R_CLIENT_CERT_VALUE {
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
	if snStr == unarylib.ORC8R_CLIENT_CERT_VALUE {
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
	var multiErr *errors.Multi
	count := 0
	for _, sn := range snList.Sns {
		certInfo, err := srv.getCertInfo(sn)
		if err != nil {
			multiErr.AddFmt(err, "'%s' get info error:", sn)
		}
		notAfter, _ := ptypes.Timestamp(certInfo.NotAfter)
		notAfter = notAfter.Add(CollectGarbageAfter)
		if time.Now().UTC().After(notAfter) {
			err = srv.store.DeleteCertInfo(sn)
			if err != nil {
				multiErr.AddFmt(err, "'%s' delete error:", sn)
			} else {
				count += 1
			}
		}
	}
	if multiErr.AsError() != nil {
		glog.Errorf("Failed to delete certificate[s]: %v", multiErr)
		return count, status.Error(codes.Internal, multiErr.Error())
	}
	return count, nil
}

// GetPolicyDecision makes a policy decision when a user attempts to access a resource
func (srv *CertifierServer) GetPolicyDecision(ctx context.Context, getPDReq *certprotos.GetPolicyDecisionRequest) (*certprotos.PolicyDecision, error) {
	username := getPDReq.Username
	user, err := srv.store.GetUser(username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user %s from database: %v", username, err)
	}

	err = isTokenWithUser(getPDReq.Token, user.Tokens)
	if err != nil {
		return nil, err
	}

	decision, err := srv.getPolicyDecisionFromTokenMany(user.Tokens, getPDReq)
	if err != nil {
		return nil, err
	}

	return decision, nil
}

// CreateUser creates a new user with the specified password and policy
func (srv *CertifierServer) CreateUser(ctx context.Context, user *certprotos.User) (*protos.Void, error) {
	// token, err := certifier.GenerateToken(certifier.Personal)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "failed to generate token while creating user")
	// }

	_, err := srv.store.GetUser(user.Username)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "error hashing password: %v", err)
	}

	user = &certprotos.User{
		Username: user.Username,
		Password: hashedPassword,
	}
	err = srv.store.PutUser(user.Username, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store user while creating user")
	}

	// policy := &certprotos.Policy{
	// 	Token:     token,
	// 	Resources: createUserReq.Policy.Resources,
	// }
	// err = srv.store.PutPolicy(token, policy)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to store policy while creating user")
	// }
	return &protos.Void{}, nil
}

func (srv *CertifierServer) ListUsers(ctx context.Context, void *protos.Void) (*certprotos.ListUsersResponse, error) {
	users, err := srv.store.ListUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users from db")
	}
	return &certprotos.ListUsersResponse{Users: users}, nil
}

func (srv *CertifierServer) GetUser(ctx context.Context, getUserReq *certprotos.GetUserRequest) (*certprotos.User, error) {
	user, err := srv.store.GetUser(getUserReq.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user %s", getUserReq.Username)
	}
	return user, nil
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

func (srv *CertifierServer) getPolicyDecisionFromTokenMany(tokens *certprotos.TokenList, getPDReq *certprotos.GetPolicyDecisionRequest) (*certprotos.PolicyDecision, error) {
	resource := getPDReq.Resource

	finalEffect := certprotos.Effect_DENY
	for _, t := range tokens.Tokens {
		effect, _ := srv.getPolicyDecisionFromToken(t, resource)
		switch effect {
		case certprotos.Effect_DENY:
			// Return early if there is a DENY in any of the permissions
			return &certprotos.PolicyDecision{Effect: certprotos.Effect_DENY}, nil
		case certprotos.Effect_ALLOW:
			finalEffect = certprotos.Effect_ALLOW
		default:
			continue
		}
	}
	// Return DENY if the policy unknown
	return &certprotos.PolicyDecision{Effect: finalEffect}, nil
}

func (srv *CertifierServer) getPolicyDecisionFromToken(token string, reqResource *certprotos.Resource) (certprotos.Effect, error) {
	policy, err := srv.store.GetPolicy(token)
	if err != nil {
		return certprotos.Effect_UNKNOWN, status.Errorf(codes.Internal, "failed to get policy from db, checking next policy: %v", err)
	}
	effect := certprotos.Effect_UNKNOWN
	for _, policyResource := range policy.Resources.Resources {
		effect = checkResource(reqResource, policyResource)
		if effect == certprotos.Effect_DENY {
			return certprotos.Effect_DENY, nil
		}

		effect = checkAction(reqResource, policyResource)
		if effect == certprotos.Effect_DENY {
			return certprotos.Effect_DENY, nil
		}
	}

	return effect, nil
}

func checkResource(reqResource, policyResource *certprotos.Resource) certprotos.Effect {
	switch reqResource.ResourceType {
	case certprotos.ResourceType_URI:
		return checkRequestURI(reqResource, policyResource)
	default:
		// TODO(christinewang5): implement for other req resource types
		return certprotos.Effect_UNKNOWN
	}
	return certprotos.Effect_UNKNOWN
}

// checkRequestURI checks if the requested resource is authorized by the policy
func checkRequestURI(reqResource *certprotos.Resource, policyResource *certprotos.Resource) certprotos.Effect {
	// The policy does not handle this case, so return UNKNOWN.
	if reqResource.ResourceType != certprotos.ResourceType_URI || policyResource.ResourceType != certprotos.ResourceType_URI {
		return certprotos.Effect_UNKNOWN
	}
	reqURI := reqResource.Resource
	policyURI := policyResource.Resource
	if ok, _ := doublestar.Match(reqURI, policyURI); ok {
		return policyResource.Effect
	}
	return certprotos.Effect_UNKNOWN
}

// checkAction checks if the requested action is authorized by the policy
func checkAction(reqResource *certprotos.Resource, policyResource *certprotos.Resource) certprotos.Effect {
	if policyResource.Action == certprotos.Action_WRITE {
		return policyResource.Effect
	}
	if policyResource.Action == certprotos.Action_READ && reqResource.Action == certprotos.Action_READ {
		return policyResource.Effect
	}
	return certprotos.Effect_UNKNOWN
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
