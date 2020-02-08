/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/security/key"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
)

func CreateCSR(validTime time.Duration, cn, idCn string) (*protos.CSR, error) {
	return createCSR(
		validTime,
		cn,
		&protos.Identity{Value: &protos.Identity_Operator{Operator: idCn}})
}

func CreateCSRForId(validTime time.Duration, id *protos.Identity) (*protos.CSR, error) {
	cn := id.ToCommonName()
	if cn == nil {
		t := "UNDEFINED"
		cn = &t
	}
	return createCSR(validTime, *cn, id)
}

func createCSR(validTime time.Duration, cn string, id *protos.Identity) (*protos.CSR, error) {
	priv, err := key.GenerateKey("", 2048)
	if err != nil {
		return nil, fmt.Errorf("Failed to create key: %s", err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"FB"},
			OrganizationalUnit: []string{"FB Inc."},
			CommonName:         cn,
		},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	csr := &protos.CSR{
		Id:        id,
		ValidTime: ptypes.DurationProto(validTime),
		CsrDer:    csrDER,
	}

	return csr, nil
}

func CreateSignedCertAndPrivKey(validTime time.Duration) (*x509.Certificate, interface{}, error) {
	priv, err := key.GenerateKey("", 2048)
	notBefore := clock.Now().UTC()
	notAfter := notBefore.Add(validTime)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	ski := make([]byte, 32)
	rand.Read(ski)

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		DNSNames:              []string{"cloud.magma.facebook.com"},
		IsCA:                  true, // MUST be CA to sign & verify client certs
		SubjectKeyId:          ski,
		BasicConstraintsValid: true,
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"FB TEST CA"},
			OrganizationalUnit: []string{"FB TEST CA"},
			CommonName:         "",
		},
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	certDER, err := x509.CreateCertificate(
		rand.Reader, &template, &template, key.PublicKey(priv), priv)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create certificate: %s", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse certificate: %s", err)
	}
	return cert, priv, nil
}
