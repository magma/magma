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
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ChallengeExpireTime = time.Minute * 5
const ChallengeLength = 512
const TimeLength = 8 // length of time encoded in byte array from int64
const MinKeyLength = 1024
const GatewayCertificateDuration = time.Hour * 97 // 4 days, lifetime of GW Certificate

type BootstrapperServer struct {
	privKey *rsa.PrivateKey
}

func NewBootstrapperServer(privKey *rsa.PrivateKey) (*BootstrapperServer, error) {
	srv := &BootstrapperServer{}
	if privKey.N.BitLen() < MinKeyLength {
		return nil, errorLogger(errors.Errorf("private key is too short: actual len (%d) is less than minimum len (%d)", privKey.N.BitLen(), MinKeyLength))
	}
	srv.privKey = privKey
	return srv, nil
}

// generate challenge in the format of [randomText : timestamp : signature]
func (srv *BootstrapperServer) GetChallenge(ctx context.Context, hwId *protos.AccessGatewayID) (*protos.Challenge, error) {
	var keyType protos.ChallengeKey_KeyType

	// case based on the env variable whether to use magmad or configurator
	var err error
	keyType, _, err = getChallengeKey(hwId.Id)
	if err != nil {
		return nil, err
	}

	if keyType != protos.ChallengeKey_ECHO &&
		keyType != protos.ChallengeKey_SOFTWARE_RSA_SHA256 &&
		keyType != protos.ChallengeKey_SOFTWARE_ECDSA_SHA256 {
		return nil, errorLogger(status.Errorf(codes.Aborted, "Unsupported key type: %s", keyType))
	}

	// generate random text
	randText, err := generateRandomText(ChallengeLength - TimeLength - srv.signatureLength())
	if err != nil {
		return nil, errorLogger(status.Errorf(codes.Aborted, "Failed to generate random text: %s", err))
	}

	// generate timestamp
	timeBytes := make([]byte, TimeLength)
	binary.BigEndian.PutUint64(timeBytes, uint64(time.Now().UTC().Unix()))

	// generate challenge
	challenge := append(randText, timeBytes...)
	signature, err := srv.sign(challenge)
	if err != nil {
		err = status.Errorf(codes.Aborted, "Failed to sign the challenge: %s", err)
		return nil, errorLogger(err)
	}
	challenge = append(challenge, signature...)

	return &protos.Challenge{KeyType: keyType, Challenge: challenge}, nil
}

// verify the response by client and return signed certificate if response is correct
func (srv *BootstrapperServer) RequestSign(ctx context.Context, resp *protos.Response) (*protos.Certificate, error) {
	hwId := resp.HwId.Id
	keyType, key, err := getChallengeKey(hwId)
	if err != nil {
		return nil, err
	}

	err = srv.verifyChallenge(resp.Challenge)
	if err != nil {
		return nil, errorLogger(status.Errorf(codes.Aborted, "Failed to verify challenge: %s", err))
	}

	// verify authentication / real response
	switch keyType {
	case protos.ChallengeKey_ECHO:
		err = verifyEcho(resp)
	case protos.ChallengeKey_SOFTWARE_RSA_SHA256:
		err = verifySoftwareRSASHA256(resp, key)
	case protos.ChallengeKey_SOFTWARE_ECDSA_SHA256:
		err = verifySoftwareECDSASHA256(resp, key)
	default:
		err = fmt.Errorf("Unsupported key type: %s", keyType)
	}
	if err != nil {
		return nil, errorLogger(status.Errorf(codes.Aborted, "Failed to verify response: %s", err))
	}

	// Ignore requested cert duration & overwrite it with our own if it's
	// longer than our default duration (allow shorter-lived certs)
	if resp.Csr != nil {
		reqValidDuration, err := ptypes.Duration(resp.Csr.ValidTime)
		if err != nil || reqValidDuration.Nanoseconds() > GatewayCertificateDuration.Nanoseconds() {
			resp.Csr.ValidTime = ptypes.DurationProto(GatewayCertificateDuration)
		}
	}
	cert, err := certifier.SignCSR(resp.Csr)
	if err != nil {
		return nil, errorLogger(status.Errorf(codes.Aborted, "Failed to sign csr: %s", err))
	}
	return cert, nil
}

// return the length of signature (number of bytes)
func (srv *BootstrapperServer) signatureLength() int {
	keyLength := srv.privKey.N.BitLen()
	return keyLength / 8
}

// authenticate text (challenge) with signature
func (srv *BootstrapperServer) sign(text []byte) ([]byte, error) {
	hashed := sha256.Sum256(text)
	return rsa.SignPKCS1v15(rand.Reader, srv.privKey, crypto.SHA256, hashed[:])
}

// verify the signature of text (challenge) to confirm the text is sent by server
func (srv *BootstrapperServer) verify(text, signature []byte) error {
	hashed := sha256.Sum256(text)
	return rsa.VerifyPKCS1v15(&srv.privKey.PublicKey, crypto.SHA256, hashed[:], signature)
}

// verify the challenge is the one sent by server and is not expired
// challenge[randomText : timestamp : signature]
func (srv *BootstrapperServer) verifyChallenge(challenge []byte) error {
	if n := len(challenge); n != ChallengeLength {
		return fmt.Errorf("Wrong length for challenge, expected %d, got %d", ChallengeLength, n)
	}

	// check time
	randLen := ChallengeLength - TimeLength - srv.signatureLength()
	timeBytes := challenge[randLen : randLen+TimeLength]
	issueTime := time.Unix(int64(binary.BigEndian.Uint64(timeBytes)), 0)
	expireTime := issueTime.Add(ChallengeExpireTime)
	now := time.Now().UTC()
	if issueTime.After(now) {
		return fmt.Errorf("Challenge is not valid yet")

	}
	if expireTime.Before(now) {
		return fmt.Errorf("Challenge has expired")
	}

	// verify signature
	signatureIdx := ChallengeLength - srv.signatureLength()
	err := srv.verify(challenge[:signatureIdx], challenge[signatureIdx:])
	if err != nil {
		return fmt.Errorf("Failed to veriry the authenticity of challenge")
	}
	return nil
}

// generate random byte slice given length
func generateRandomText(length int) ([]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("Incorrect length")
	}
	randText := make([]byte, length)
	_, err := rand.Read(randText)
	if err != nil {
		return nil, err
	}
	return randText, nil
}

// verify response with echo "encryption" method
func verifyEcho(resp *protos.Response) error {
	response := resp.GetEchoResponse() //.Response
	if response == nil {
		return fmt.Errorf("Wrong type of response, expected Echo")
	}
	if !bytes.Equal(response.Response, resp.Challenge) {
		return fmt.Errorf("Incorrect response")
	}
	return nil
}

// verify response with RSA signature and sha256 hash
func verifySoftwareRSASHA256(resp *protos.Response, key []byte) error {
	publicKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("Failed to parse RSA public key: %s", err)
	}

	response := resp.GetRsaResponse()
	if response == nil {
		return fmt.Errorf("Wrong type of response, expected RSA")
	}

	hashed := sha256.Sum256(resp.Challenge)
	err = rsa.VerifyPKCS1v15(
		publicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], response.Signature)
	return err
}

// verify response with ecdsa signature ahd sha256 hash
func verifySoftwareECDSASHA256(resp *protos.Response, key []byte) error {
	publicKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("Failed to parse ECDSA public key %s", err)
	}

	response := resp.GetEcdsaResponse()
	if response == nil {
		return fmt.Errorf("Wrong type of response, expected ECDSA")
	}

	var r, s big.Int
	r.SetBytes(response.R)
	s.SetBytes(response.S)
	hashed := sha256.Sum256(resp.Challenge)
	if !ecdsa.Verify(publicKey.(*ecdsa.PublicKey), hashed[:], &r, &s) {
		return fmt.Errorf("Wrong response")
	}
	return nil
}

func getChallengeKey(hwID string) (protos.ChallengeKey_KeyType, []byte, error) {
	var empty protos.ChallengeKey_KeyType
	entity, err := configurator.LoadEntityForPhysicalID(hwID, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return empty, nil, errorLogger(status.Errorf(codes.NotFound, "Gateway with hwid %s is not registered: %s", hwID, err))
	}
	iRecord, err := device.GetDevice(entity.NetworkID, orc8r.AccessGatewayRecordType, hwID, serdes.Device)
	if err != nil {
		return empty, nil, errorLogger(status.Errorf(codes.NotFound, "Failed to find gateway record: %s", err))
	}
	record, ok := iRecord.(*models.GatewayDevice)
	if !ok {
		return empty, nil, errorLogger(status.Errorf(codes.NotFound, "Failed to find gateway record"))
	}

	var key []byte
	keyType, ok := protos.ChallengeKey_KeyType_value[record.Key.KeyType]
	if !ok {
		return empty, nil, errorLogger(status.Errorf(codes.Aborted, "Unsupported key type: %v", keyType))
	}
	if record.Key.Key != nil {
		key = *record.Key.Key
	}
	return protos.ChallengeKey_KeyType(keyType), key, nil
}

func errorLogger(err error) error {
	log.Printf("Bootstrapper Error: %v", err)
	return err
}
