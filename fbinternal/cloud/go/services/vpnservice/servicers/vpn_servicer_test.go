package servicers_test

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	fbprotos "magma/fbinternal/cloud/go/protos"
	"magma/fbinternal/cloud/go/services/vpnservice/servicers"
	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	certifier_test_init "magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/key"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	TestTaKey = `#
# 2048 bit OpenVPN static key
#
-----BEGIN OpenVPN Static key V1-----
abcdef01234567890123456789012345
abcdef01234567890123456789012345
-----END OpenVPN Static key V1-----
`
	TestTaKeyFile = "/tmp/vpn_ta.key"
)

func TestGetCA(t *testing.T) {
	srv, vpnCert := getVPNServicer(t)

	// check the CA comes back the same
	srvCA, err := srv.GetCA(context.Background(), &protos.Void{})
	assert.NoError(t, err)

	assert.Equal(t, srvCA.Cert, vpnCert.Raw)
}

func TestVPNServicer_RequestPSK(t *testing.T) {
	f, err := os.Create(TestTaKeyFile)
	assert.NoError(t, err)
	_, err = f.WriteString(TestTaKey)
	assert.NoError(t, err)
	f.Close()
	srv := servicers.NewVPNServicer(TestTaKeyFile)
	id := protos.NewGatewayIdentity("testHwId", "testNwId", "testLogicalId")

	key, err := srv.RequestPSK(id.NewContextWithIdentity(context.Background()), &protos.Void{})
	assert.NoError(t, err)
	fileInBytes, err := ioutil.ReadFile(TestTaKeyFile)
	assert.NoError(t, err)
	assert.Equal(t, fileInBytes, key.TaKey)
	err = os.Remove(TestTaKeyFile)
	assert.NoError(t, err)
	_, err = srv.RequestPSK(id.NewContextWithIdentity(context.Background()), &protos.Void{})
	assert.True(t, strings.Contains(err.Error(), "err loading PSK"))
}

func TestRequestSign(t *testing.T) {
	srv, vpnCert := getVPNServicer(t)

	// make private key and csr
	privKey, err := key.GenerateKey("", 2048)
	assert.NoError(t, err)

	csr, err := createCSRBytes(privKey)
	assert.NoError(t, err)

	signReq := &fbprotos.VPNCertRequest{
		Request: csr,
	}

	signCheckCert := func(id *protos.Identity) {
		testContext := id.NewContextWithIdentity(context.Background())
		signedCertMsg, err := srv.RequestCert(testContext, signReq)
		assert.NoError(t, err)
		verifySignedCert(t, vpnCert, signedCertMsg)
	}

	signCheckCert(protos.NewGatewayIdentity("test", "test", "test"))
	signCheckCert(protos.NewOperatorIdentity("test"))
}

// Starts up certifier and returns VPNServicer, VPN CA cert
func getVPNServicer(t *testing.T) (*servicers.VPNServicer, *x509.Certificate) {
	certifier_test_init.StartTestService(t)

	caMsg, err := certifier.GetCACert(&certprotos.GetCARequest{CertType: protos.CertType_VPN})
	assert.NoError(t, err)

	vpnCA, err := x509.ParseCertificate(caMsg.Cert)
	assert.NoError(t, err)

	return servicers.NewVPNServicer(TestTaKeyFile), vpnCA
}

// Make a certificate request with the given private key
func createCSRBytes(privKey interface{}) ([]byte, error) {
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"FB"},
			OrganizationalUnit: []string{"FB Inc."},
			CommonName:         "test",
		},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create csr: %s", err)
	}
	return csrDER, nil
}

// Perform various checks on certificate that we want signed with given CA
func verifySignedCert(t *testing.T, ca *x509.Certificate, certMsg *fbprotos.VPNCertificate) {
	// make sure we can get the identity
	_, err := certifier.GetCertificateIdentity(certMsg.Serial)
	assert.NoError(t, err)

	// deserialize the certificate
	cert, err := x509.ParseCertificate(certMsg.Cert)
	assert.NoError(t, err)

	// check the time is valid
	now := time.Now()
	assert.True(t, now.After(cert.NotBefore))
	assert.True(t, now.Before(cert.NotAfter))

	// check it's signed by CA
	clientCert, err := x509.ParseCertificate(certMsg.Cert)
	assert.NoError(t, err)

	caPool := x509.NewCertPool()
	caPool.AddCert(ca)
	opts := x509.VerifyOptions{
		Roots:         caPool,
		Intermediates: x509.NewCertPool(),
		// Make sure client cert has ExtKeyUsageClientAuth
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	_, err = clientCert.Verify(opts)
	assert.NoError(t, err)
}
