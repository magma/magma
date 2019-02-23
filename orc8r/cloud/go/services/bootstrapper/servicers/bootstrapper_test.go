/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/security/key"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers"
	certifier_test_init "magma/orc8r/cloud/go/services/certifier/test_init"
	certifier_test_utils "magma/orc8r/cloud/go/services/certifier/test_utils"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func testWithECHO(
	t *testing.T, networkId string, srv *servicers.BootstrapperServer, ctx context.Context) {

	testAgHwId := "test_ag_echo"

	_, err := magmad.RegisterGateway(
		networkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW echo",
			Key:  &protos.ChallengeKey{KeyType: protos.ChallengeKey_ECHO},
		})
	assert.NoError(t, err)

	// check challenge type
	challenge, err := srv.GetChallenge(ctx, &protos.AccessGatewayID{Id: testAgHwId})
	assert.NoError(t, err)
	assert.Equal(t, challenge.KeyType, protos.ChallengeKey_ECHO)

	// create response
	response := &protos.Response_EchoResponse{
		EchoResponse: &protos.Response_Echo{Response: challenge.Challenge},
	}
	csr, err := certifier_test_utils.CreateCSR(time.Duration(time.Hour*24*10), "cn", "cn")
	assert.NoError(t, err)
	resp := protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       csr,
	}
	cert, err := srv.RequestSign(ctx, &resp)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func testWithRSA(
	t *testing.T, networkId string, srv *servicers.BootstrapperServer, ctx context.Context) {

	testAgHwId := "test_ag_rsa"
	privateKey, err := key.GenerateKey("", 1024)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)

	_, err = magmad.RegisterGateway(
		networkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW RSA",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_SOFTWARE_RSA_SHA256,
				Key:     marshaledPubKey},
		})
	assert.NoError(t, err)

	challenge, err := srv.GetChallenge(ctx, &protos.AccessGatewayID{Id: testAgHwId})
	assert.NoError(t, err)
	assert.Equal(t, challenge.KeyType, protos.ChallengeKey_SOFTWARE_RSA_SHA256)

	// sign challenge with private key
	hashed := sha256.Sum256(challenge.Challenge)
	signature, err := rsa.SignPKCS1v15(
		rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	assert.NoError(t, err)

	// create response
	response := &protos.Response_RsaResponse{
		RsaResponse: &protos.Response_RSA{Signature: signature},
	}
	csr, err := certifier_test_utils.CreateCSR(time.Duration(time.Hour*24*10), "cn", "cn")
	assert.NoError(t, err)
	resp := protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       csr,
	}
	cert, err := srv.RequestSign(ctx, &resp)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func testWithECDSA(
	t *testing.T, networkId string, srv *servicers.BootstrapperServer, ctx context.Context) {

	testAgHwId := "test_ag_ecdsa"
	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)

	_, err = magmad.RegisterGateway(
		networkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW ECDSA",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_SOFTWARE_ECDSA_SHA256,
				Key:     marshaledPubKey},
		})
	assert.NoError(t, err)

	challenge, err := srv.GetChallenge(ctx, &protos.AccessGatewayID{Id: testAgHwId})
	assert.NoError(t, err)
	assert.Equal(t, challenge.KeyType, protos.ChallengeKey_SOFTWARE_ECDSA_SHA256)

	hashed := sha256.Sum256(challenge.Challenge)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hashed[:])
	assert.NoError(t, err)

	// create response
	response := &protos.Response_EcdsaResponse{
		EcdsaResponse: &protos.Response_ECDSA{R: r.Bytes(), S: s.Bytes()},
	}
	csr, err := certifier_test_utils.CreateCSR(time.Duration(time.Hour*24*10), "cn", "cn")
	assert.NoError(t, err)
	resp := protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       csr,
	}
	cert, err := srv.RequestSign(ctx, &resp)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func testNegative(
	t *testing.T, networkId string, srv *servicers.BootstrapperServer, ctx context.Context) {

	testAgHwId := "test_ag_negative"
	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)

	_, err = magmad.RegisterGateway(
		networkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW ECDSA",
			Key:  &protos.ChallengeKey{KeyType: 10, Key: marshaledPubKey},
		})
	assert.NoError(t, err)
	// cannot get challenge because of unsupported key type
	_, err = srv.GetChallenge(ctx, &protos.AccessGatewayID{Id: testAgHwId})
	assert.Error(t, err)

	testAgHwId = "test_ag_negative2"
	_, err = magmad.RegisterGateway(
		networkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW ECDSA",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_SOFTWARE_ECDSA_SHA256,
				Key:     marshaledPubKey},
		})
	assert.NoError(t, err)

	challenge, err := srv.GetChallenge(ctx, &protos.AccessGatewayID{Id: testAgHwId})
	assert.NoError(t, err)

	// compute response
	hashed := sha256.Sum256(challenge.Challenge)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hashed[:])
	assert.NoError(t, err)

	csr, err := certifier_test_utils.CreateCSR(time.Duration(time.Hour*24*10), "cn", "cn")
	assert.NoError(t, err)

	// create response
	response := &protos.Response_EcdsaResponse{
		EcdsaResponse: &protos.Response_ECDSA{R: r.Bytes(), S: s.Bytes()},
	}

	// mess up challenge
	resp := protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: []byte("mess up challenge"),
		Response:  response,
		Csr:       csr,
	}
	_, err = srv.RequestSign(ctx, &resp)
	assert.Error(t, err)

	// mess up csr
	resp = protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       nil,
	}
	_, err = srv.RequestSign(ctx, &resp)
	assert.Error(t, err)

	// mess up response
	response = &protos.Response_EcdsaResponse{
		EcdsaResponse: &protos.Response_ECDSA{R: []byte("12344"), S: s.Bytes()},
	}
	resp = protos.Response{
		HwId:      &protos.AccessGatewayID{Id: testAgHwId},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       csr,
	}
	_, err = srv.RequestSign(ctx, &resp)
	assert.Error(t, err)

	// mess up hw_id
	resp = protos.Response{
		HwId:      &protos.AccessGatewayID{Id: "mess up hw_id"},
		Challenge: challenge.Challenge,
		Response:  response,
		Csr:       csr,
	}
	_, err = srv.RequestSign(ctx, &resp)
	assert.Error(t, err)
}

func TestBootstrapperServer(t *testing.T) {
	magmad_test_init.StartTestService(t)
	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"bootstrapper_test_network")
	assert.NoError(t, err)

	ctx := context.Background()

	// create bootstrapper with short key
	privateKey, err := key.GenerateKey("", 512)
	assert.NoError(t, err)
	_, err = servicers.NewBootstrapperServer(privateKey.(*rsa.PrivateKey))
	assert.Error(t, err)

	// create bootstrapper server
	privateKey, err = key.GenerateKey("", 2048)
	assert.NoError(t, err)
	srv, err := servicers.NewBootstrapperServer(privateKey.(*rsa.PrivateKey))

	// for signing csr
	certifier_test_init.StartTestService(t)

	testWithECHO(t, testNetworkId, srv, ctx)
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", "bla"))
	testWithRSA(t, testNetworkId, srv, ctx)
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-serial", ""))
	testWithECDSA(t, testNetworkId, srv, ctx)
	ctx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-magma-client-cert-cn", "bla"))
	testNegative(t, testNetworkId, srv, ctx)
}
