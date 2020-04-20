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

package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// returns 'public part' of the passed asymmetric encryption algo key 'priv'
func PublicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

// returns PEM key block structure with ASN.1 DER encoded private key 'priv'
func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			err = fmt.Errorf("Unable to marshal ECDSA private key: %v\n", err)
			return nil, err
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		err := fmt.Errorf("Key type %v is not supported.", reflect.TypeOf(priv))
		return nil, err
	}
}

// write private key 'priv' into 'keyFile' in ASN.1 DER format
// type of 'priv' can be *rsa.PrivateKey or *ecdsa.PrivateKey
func WriteKey(keyFile string, priv interface{}) error {
	keyOut, err := os.OpenFile(
		keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open %s for writing: %s", keyFile, err)
	}

	keyPEM, err := pemBlockForKey(priv)
	if err != nil {
		return err
	}

	pem.Encode(keyOut, keyPEM)
	keyOut.Close()
	return nil
}

// read and parse private key from 'keyFile', return the 'priv' key
// in the form of either *rsa.PrivateKey or *ecdsa.PrivateKey
func ReadKey(keyFile string) (priv interface{}, err error) {
	byteKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	pemKey, _ := pem.Decode(byteKey)
	if pemKey == nil {
		return nil, fmt.Errorf("Failed to find PEM block in file %s", keyFile)
	}

	switch pemKey.Type {
	case "RSA PRIVATE KEY":
		priv, err = x509.ParsePKCS1PrivateKey(pemKey.Bytes)
	case "EC PRIVATE KEY":
		priv, err = x509.ParseECPrivateKey(pemKey.Bytes)
	default:
		err = fmt.Errorf("Key type %s is not supported.", pemKey.Type)
		priv = nil
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse key from file %s: %s", keyFile, err)
	}
	return
}

// return a private key generated using rsa or ecdsa.
// if param ecdsaCurve is empty string, then rsa with 'rsaBits' will be used;
// otherwise corresponding ecdsa algorithm will be used and 'rsaBits' is ignored.
func GenerateKey(ecdsaCurve string, rsaBits int) (priv interface{}, err error) {
	switch ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		err = fmt.Errorf("Unrecognized elliptic curve: %q", ecdsaCurve)
	}
	return
}
