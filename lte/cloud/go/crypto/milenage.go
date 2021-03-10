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

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	// ExpectedKeyBytes is the number of bytes for the subscriber key.
	ExpectedKeyBytes = 16

	// ExpectedOpcBytes is the number of bytes for the operator variant algorithm configuration field.
	ExpectedOpcBytes = 16

	// ExpectedPlmnBytes is the number of bytes for the network identifier.
	ExpectedPlmnBytes = 3

	// ExpectedAmfBytes is the number of bytes for the authentication management field.
	ExpectedAmfBytes = 2

	// ExpectedOpBytes is the number of bytes for the operator variant configuration field.
	ExpectedOpBytes = 16

	// ExpectedAutsBytes is the number of bytes for the authentication token from the client key.
	ExpectedAutsBytes = 14

	// RandChallengeBytes is the number of bytes for the random challenge.
	RandChallengeBytes = 16

	// XresBytes is the number of bytes for the expected response.
	XresBytes = 8

	// AutnBytes is the number of bytes for the authentication token.
	AutnBytes = 16

	// KasmeBytes is the number of bytes for the base network authentication token.
	KasmeBytes = 32

	// ConfidentialityKeyBytes is the number of bytes for the confidentiality key.
	ConfidentialityKeyBytes = 16

	// IntegrityKeyBytes is the number of bytes for the integrity key.
	IntegrityKeyBytes = 16

	// AnonymityKeyBytes is the number of bytes for the anonymity key.
	AnonymityKeyBytes = 16

	// The highest valid sequence number (since sequence numbers are 48 bits).
	maxSqn      = (1 << 48) - 1
	sqnMaxBytes = 6
)

// MilenageCipher implements the milenage algorithm (3GPP TS 35.205, .206, .207, .208)
type MilenageCipher struct {
	// rng is a cryptographically secure random number generator
	rng cryptoRNG

	// amf is a 16 bit authentication management field
	amf [ExpectedAmfBytes]byte
}

// NewMilenageCipher instantiates the Milenage algo using crypto/rand for rng.
func NewMilenageCipher(amf []byte) (*MilenageCipher, error) {
	if len(amf) != ExpectedAmfBytes {
		return nil, fmt.Errorf("incorrect amf size. Expected 2 bytes, but got %v bytes", len(amf))
	}

	milenage := &MilenageCipher{rng: defaultCryptoRNG{}}
	copy(milenage.amf[:], amf)
	return milenage, nil
}

// EutranVector reprsents an E-UTRAN key vector.
type EutranVector struct {
	// Rand is a random challenge
	Rand [RandChallengeBytes]byte

	// Xres is the expected response
	Xres [XresBytes]byte

	// Autn is an authentication token
	Autn [AutnBytes]byte

	// Kasme is a base network authentication token
	Kasme [KasmeBytes]byte
}

// UtranVector represents a UTRAN key vector
type UtranVector struct {
	// Rand is a random challenge
	Rand [RandChallengeBytes]byte

	// Xres is the expected response
	Xres [XresBytes]byte

	// Autn is an authentication token
	Autn [AutnBytes]byte

	// Confidentialitykey is used to ensure the confidentiality of messages
	ConfidentialityKey [ConfidentialityKeyBytes]byte

	// IntegrityKey is used to ensure the integrity of messages
	IntegrityKey [IntegrityKeyBytes]byte
}

// SIPAuthVector represents the data encoded in a SIP auth data item.
type SIPAuthVector struct {
	// Rand is a random challenge
	Rand [RandChallengeBytes]byte

	// Xres is the expected response
	Xres [XresBytes]byte

	// Autn is an authentication token
	Autn [AutnBytes]byte

	// Confidentialitykey is used to ensure the confidentiality of messages
	ConfidentialityKey [ConfidentialityKeyBytes]byte

	// IntegrityKey is used to ensure the integrity of messages
	IntegrityKey [IntegrityKeyBytes]byte

	// AnonymityKey is used to ensure the anonymity of messages
	AnonymityKey [AnonymityKeyBytes]byte
}

// GenerateEutranVector creates an E-UTRAN key vector.
// Inputs:
//   key: 128 bit subscriber key
//   opc: 128 bit operator variant algorithm configuration field
//   sqn: 48 bit sequence number
//   plmn: 24 bit network identifier
//      Octet           Description
//         1      MCC digit 2 | MCC digit 1
//         2      MNC digit 3 | MCC digit 3
//         3      MNC digit 2 | MNC digit 1
// Outputs: An EutranVector or an error. The EutranVector is not nil if and only if err == nil.
func (milenage *MilenageCipher) GenerateEutranVector(key, opc []byte, sqn uint64, plmn []byte) (*EutranVector, error) {
	var randChallenge = make([]byte, RandChallengeBytes)
	_, err := milenage.rng.Read(randChallenge)
	if err != nil {
		return nil, err
	}
	return milenage.GenerateEutranVectorWithRand(key, opc, randChallenge, sqn, plmn)
}

// GenerateEutranVectorWithRand creates an E-UTRAN key vector.
// Inputs:
//   key:  128 bit subscriber key
//   opc:  128 bit operator variant algorithm configuration field
//   rand: 128 bit random challenge
//   sqn:  48 bit sequence number
//   plmn: 24 bit network identifier
//      Octet           Description
//         1      MCC digit 2 | MCC digit 1
//         2      MNC digit 3 | MCC digit 3
//         3      MNC digit 2 | MNC digit 1
// Outputs: An EutranVector or an error. The EutranVector is not nil if and only if err == nil.
func (milenage *MilenageCipher) GenerateEutranVectorWithRand(
	key, opc, rand []byte, sqn uint64, plmn []byte) (*EutranVector, error) {

	err := validateGenerateEutranVectorInputs(key, opc, sqn, plmn)
	if err != nil {
		return nil, err
	}
	vector, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc, sqn)
	if err != nil {
		return nil, err
	}
	sqnBytes := getSqnBytes(sqn)
	kasme, err := generateKasme(vector.ConfidentialityKey[:], vector.IntegrityKey[:], plmn, sqnBytes, vector.AnonymityKey[:])
	if err != nil {
		return nil, err
	}
	return newEutranVector(vector.Rand[:], vector.Xres[:], vector.Autn[:], kasme), nil
}

// GenerateUtranVector creates UTRAN auth vector
// Inputs:
//   key:  128 bit subscriber key
//   opc:  128 bit operator variant algorithm configuration field
//   sqn:  48 bit sequence number
// Outputs: A E-UTRAN & UTRAN auth vector or an error
func (milenage *MilenageCipher) GenerateUtranVector(key, opc []byte, sqn uint64) (*UtranVector, error) {
	var randChallenge = make([]byte, RandChallengeBytes)
	_, err := milenage.rng.Read(randChallenge)
	if err != nil {
		return nil, err
	}
	return milenage.GenerateUtranVectorWithRand(key, opc, randChallenge, sqn)
}

// GenerateUtranVectorWithRand creates UTRAN auth vector
// Inputs:
//   key:  128 bit subscriber key
//   opc:  128 bit operator variant algorithm configuration field
//   rand: 128 bit random challenge
//   sqn:  48 bit sequence number
// Outputs: A E-UTRAN & UTRAN auth vector or an error
func (milenage *MilenageCipher) GenerateUtranVectorWithRand(key, opc, rand []byte, sqn uint64) (*UtranVector, error) {
	vector, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc, sqn)
	if err != nil {
		return nil, err
	}
	return &UtranVector{
		Rand:               vector.Rand,
		Xres:               vector.Xres,
		Autn:               vector.Autn,
		ConfidentialityKey: vector.ConfidentialityKey,
		IntegrityKey:       vector.IntegrityKey,
	}, nil
}

// GenerateSIPAuthVector creates a SIP auth vector.
// Inputs:
//   key: 128 bit subscriber key
//   opc: 128 bit operator variant algorithm configuration field
//   sqn: 48 bit sequence number
// Outputs: A SIP auth vector or an error. The SIP auth vector is not nil if and only if err == nil.
func (milenage *MilenageCipher) GenerateSIPAuthVector(key []byte, opc []byte, sqn uint64) (*SIPAuthVector, error) {
	if err := validateGenerateSIPAuthVectorInputs(key, opc, sqn); err != nil {
		return nil, err
	}

	var randChallenge = make([]byte, RandChallengeBytes)
	_, err := milenage.rng.Read(randChallenge)
	if err != nil {
		return nil, err
	}
	return milenage.GenerateSIPAuthVectorWithRand(randChallenge, key, opc, sqn)
}

// GenerateSIPAuthVectorWithRand creates a SIP auth vector using a specific random challenge value.
// Inputs:
//   rand: 128 bit random challenge
//   key:  128 bit subscriber key
//   opc:  128 bit operator variant algorithm configuration field
//   sqn:  48 bit sequence number
// Outputs: A SIP auth vector or an error. The SIP auth vector is not nil if and only if err == nil.
func (milenage *MilenageCipher) GenerateSIPAuthVectorWithRand(rand []byte, key []byte, opc []byte, sqn uint64) (*SIPAuthVector, error) {
	if err := validateGenerateSIPAuthVectorWithRandInputs(rand, key, opc, sqn); err != nil {
		return nil, err
	}
	sqnBytes := getSqnBytes(sqn)

	macA, _, err := f1(key, sqnBytes, rand, opc, milenage.amf[:])
	if err != nil {
		return nil, err
	}

	xres, ak, err := f2F5(key, rand, opc)
	if err != nil {
		return nil, err
	}

	ck, err := f3(key, rand, opc)
	if err != nil {
		return nil, err
	}
	ik, err := f4(key, rand, opc)
	if err != nil {
		return nil, err
	}

	autn := generateAutn(sqnBytes, ak, macA, milenage.amf[:])
	return newSIPAuthVector(rand, xres, autn, ck, ik, ak), nil
}

// GenerateOpc returns the OP_c according to 3GPP 35.205 8.2
// Inputs:
//   key: 128 bit subscriber key
//   op: 128 bit operator variant configuration field
func GenerateOpc(key, op []byte) ([ExpectedOpcBytes]byte, error) {
	var opc [ExpectedOpcBytes]byte
	if len(key) != ExpectedKeyBytes {
		return opc, fmt.Errorf("incorrect key size. Expected %v bytes, but got %v bytes", ExpectedKeyBytes, len(key))
	}
	if len(op) != ExpectedOpBytes {
		return opc, fmt.Errorf("incorrect op size. Expected %v bytes, but got %v bytes", ExpectedOpBytes, len(op))
	}

	output, err := encrypt(key, op)
	if err != nil {
		return opc, err
	}
	copy(opc[:], xor(output, op))
	return opc, nil
}

// GenerateResync computes SQN_MS and MAC-S from AUTS for re-synchronization.
//    AUTS = SQN_MS ^ AK || f1*(SQN_MS || RAND || AMF*)
// Inputs:
//    auts: 112 bit authentication token from client key
//    opc: 128 bit operator variant algorithm configuration field
//    key: 128 bit subscriber key
//    rand: 128 bit random challenge
// Outputs: (sqnMs, macS) or an error
//	sqn_ms, 48 bit sequence number from client
//	mac_s, 64 bit resync authentication code
func (milenage *MilenageCipher) GenerateResync(auts, key, opc, rand []byte) (uint64, [8]byte, error) {
	var macS [8]byte
	err := validateGenerateResyncInputs(auts, key, opc, rand)
	if err != nil {
		return 0, macS, err
	}

	ak, err := f5Star(key, rand, opc)
	if err != nil {
		return 0, macS, err
	}
	sqnMs := xor(auts[:6], ak)
	sqnMsInt := uint64(sqnMs[5]) | uint64(sqnMs[4])<<8 | uint64(sqnMs[3])<<16 | uint64(sqnMs[2])<<24 |
		uint64(sqnMs[1])<<32 | uint64(sqnMs[0])<<40
	_, macSSlice, err := f1(key, sqnMs, rand, opc, milenage.amf[:])
	if err != nil {
		return 0, macS, err
	}
	copy(macS[:], macSSlice)
	return sqnMsInt, macS, nil
}

// validateGenerateResyncInputs ensures that each byte slice has the correct number of bytes.
// Output: An error if any of the arguments is invalid or nil otherwise.
func validateGenerateResyncInputs(auts, key, opc, rand []byte) error {
	if len(auts) != ExpectedAutsBytes {
		return fmt.Errorf("incorrect auts size. Expected %v bytes, but got %v bytes", ExpectedAutsBytes, len(auts))
	}
	if len(key) != ExpectedKeyBytes {
		return fmt.Errorf("incorrect key size. Expected %v bytes, but got %v bytes", ExpectedKeyBytes, len(key))
	}
	if len(opc) != ExpectedOpcBytes {
		return fmt.Errorf("incorrect opc size. Expected %v bytes, but got %v bytes", ExpectedOpcBytes, len(opc))
	}
	if len(rand) != RandChallengeBytes {
		return fmt.Errorf("incorrect rand size. Expected %v bytes, but got %v bytes", RandChallengeBytes, len(rand))
	}
	return nil
}

// validateGenerateEutranVectorInputs ensures that each argument has the required form.
// Each byte slice must be the correct number of bytes and sqn must fit within 48 bits.
// Output: An error if any of the arguments is invalid or nil otherwise.
func validateGenerateEutranVectorInputs(key []byte, opc []byte, sqn uint64, plmn []byte) error {
	if err := validateGenerateSIPAuthVectorInputs(key, opc, sqn); err != nil {
		return err
	}
	if len(plmn) != ExpectedPlmnBytes {
		return fmt.Errorf("incorrect plmn size. Expected 3 bytes, but got %v bytes", len(plmn))
	}
	return nil
}

// validateGenerateSIPAuthVectorInputs ensures that each argument has the required form.
// Each byte slice must be the correct number of bytes and sqn must fit within 48 bits.
// Output: An error if any of the arguments is invalid or nil otherwise.
func validateGenerateSIPAuthVectorInputs(key []byte, opc []byte, sqn uint64) error {
	if len(key) != ExpectedKeyBytes {
		return fmt.Errorf("incorrect key size. Expected %v bytes, but got %v bytes", ExpectedKeyBytes, len(key))
	}
	if len(opc) != ExpectedOpcBytes {
		return fmt.Errorf("incorrect opc size. Expected %v bytes, but got %v bytes", ExpectedOpcBytes, len(opc))
	}
	if sqn > maxSqn {
		return fmt.Errorf("sequence number too large, expected a number which can fit in 48 bits. Got: %v", sqn)
	}
	return nil
}

// validateGenerateSIPAuthVectorWithRandInputs ensures that each argument has the required form.
// Each byte slice must be the correct number of bytes and sqn must fit within 48 bits.
// Output: An error if any of the arguments is invalid or nil otherwise.
func validateGenerateSIPAuthVectorWithRandInputs(rand []byte, key []byte, opc []byte, sqn uint64) error {
	if len(rand) != RandChallengeBytes {
		return fmt.Errorf("incorrect rand size. Expected %v bytes, but got %v bytes", RandChallengeBytes, len(rand))
	}
	if len(key) != ExpectedKeyBytes {
		return fmt.Errorf("incorrect key size. Expected %v bytes, but got %v bytes", ExpectedKeyBytes, len(key))
	}
	if len(opc) != ExpectedOpcBytes {
		return fmt.Errorf("incorrect opc size. Expected %v bytes, but got %v bytes", ExpectedOpcBytes, len(opc))
	}
	if sqn > maxSqn {
		return fmt.Errorf("sequence number too large, expected a number which can fit in 48 bits. Got: %v", sqn)
	}
	return nil
}

// f1 and f1* implementation, the network authentication function and
// the re-synchronization message authentication function according to
// 3GPP 35.206 4.1
//
// Inputs:
//   key: 128 bit subscriber key
//   sqn: 48 bit sequence number
//   rand: 128 bit random challenge
//   opc: 128 bit computed from OP and subscriber key
//   amf: 16 bit authentication management field
// Outputs: (64 bit Network auth code, 64 bit Resync auth code) or an error
func f1(key, sqn, rand, opc, amf []byte) ([]byte, []byte, error) {
	// TEMP = E_K(RAND XOR OP_C)
	temp, err := encrypt(key, xor(rand, opc))
	if err != nil {
		return nil, nil, err
	}

	// IN1 = SQN || AMF || SQN || AMF
	var in1 = make([]byte, 0, ExpectedOpcBytes)
	in1 = append(in1, sqn...)
	in1 = append(in1, amf...)
	in1 = append(in1, in1...)

	const rotationBytes = 8 // Constant from 3GPP 35.206 4.1

	// OUT1 = E_K(TEMP XOR rotate(IN1 XOR OP_C, r1) XOR c1) XOR OP_C
	out1, err := encrypt(key, xor(temp, rotate(xor(in1, opc), rotationBytes)))
	if err != nil {
		return nil, nil, err
	}
	out1 = xor(out1, opc)

	// MAC-A = f1 = OUT1[0] .. OUT1[63]
	// MAC-S = f1* = OUT1[64] .. OUT1[127]
	return out1[:8], out1[8:], nil
}

// f2F5 implements f2 and f5, the compute anonymity key and response to
// challenge functions according to 3GPP 35.206 4.1
// Inputs:
//   key: 128 bit subscriber key
//   rand: 128 bit random challenge
//   opc: 128 bit computed from OP and subscriber key
//	Outputs:
//   (xres, ak) = (64 bit response to challenge, 48 bit anonymity key) or an error
func f2F5(key, rand, opc []byte) ([]byte, []byte, error) {
	var additionConstant = make([]byte, ExpectedOpcBytes) // Constant from 3GPP 35.206 4.1
	additionConstant[15] = 1

	// TEMP = E_K(RAND XOR OP_C)
	temp, err := encrypt(key, xor(rand, opc))
	if err != nil {
		return nil, nil, err
	}

	// OUT2 = E_K(rotate(TEMP XOR OP_C, r2) XOR c2) XOR OP_C
	out2, err := encrypt(key, xor(xor(temp, opc), additionConstant))
	if err != nil {
		return nil, nil, err
	}
	out2 = xor(out2, opc)

	// res = f2 = OUT2[64] ... OUT2[127]
	// ak = f5 = OUT2[0] ... OUT2[47]
	return out2[8:16], out2[:6], nil
}

// f3 implementation, the compute confidentiality key according
// to 3GPP 35.206 4.1
//
// Inputs:
//   key: 128 bit subscriber key
//   rand: 128 bit random challenge
//   opc: 128 bit computed from OP and subscriber key
// Outputs: 128 bit confidentiality key or an error
func f3(key, rand, opc []byte) ([]byte, error) {
	// Constants from 3GPP 35.206 4.1
	var additionConstant = make([]byte, ExpectedOpcBytes)
	additionConstant[15] = 2
	const rotationBytes = 4

	return f3F4Impl(key, rand, opc, additionConstant, rotationBytes)
}

// f4 implementation, the integrity key according
// to 3GPP 35.206 4.1
//
// Inputs:
//   key: 128 bit subscriber key
//   rand: 128 bit random challenge
//   opc: 128 bit computed from OP and subscriber key
// Outputs: 128 bit integrity key or an error
func f4(key, rand, opc []byte) ([]byte, error) {
	// Constants from 3GPP 35.206 4.1
	var additionConstant = make([]byte, ExpectedOpcBytes)
	additionConstant[15] = 4
	const rotationBytes = 8

	return f3F4Impl(key, rand, opc, additionConstant, rotationBytes)
}

// Implementation of f3 and f4, the compute confidentiality key according
// to 3GPP 35.206 4.1
// (f3 and f4 are the same except they use different addition and rotation constants)
//
// Inputs:
//   key (bytes): 128 bit subscriber key
//   rand (bytes): 128 bit random challenge
//   opc (bytes): 128 bit computed from OP and subscriber key
//   additionConstant (bytes): 128 bit fixed constant (defined by 3GPP 35.206 4.1)
//   rotationBytes (int): the number of bytes to shift by (defined by 3GPP 35.206 4.1)
// Outputs: 128 bit key or an error
func f3F4Impl(key, rand, opc, additionConstant []byte, rotationBytes int) ([]byte, error) {
	// TEMP = E_K(RAND XOR OP_C)
	temp, err := encrypt(key, xor(rand, opc))
	if err != nil {
		return nil, err
	}

	// OUT = E_K(rotate(TEMP XOR OP_C, r3) XOR c3) XOR OP_C
	out, err := encrypt(key, xor(rotate(xor(temp, opc), rotationBytes), additionConstant))
	if err != nil {
		return nil, err
	}
	return xor(out, opc), nil
}

// f5* implementation, the anonymity key according to 3GPP 35.206 4.1
// Inputs:
//    key: 128 bit subscriber key
//    rand: 128 bit random challenge
//    opc: 128 bit computed from OP and subscriber key
// Outputs: ak, 48 bit anonymity key or an error
func f5Star(key, rand, opc []byte) ([]byte, error) {
	// Constants from 3GPP 35.206 4.1
	var additionConstant = make([]byte, ExpectedOpcBytes)
	additionConstant[15] = 8
	const rotationBytes = 12

	// TEMP = E_K(RAND XOR OP_C)
	temp, err := encrypt(key, xor(rand, opc))
	if err != nil {
		return nil, err
	}

	// OUT = E_K(rotate(TEMP XOR OP_C, r5 XOR c5) XOR OP_C
	out, err := encrypt(key, xor(rotate(xor(temp, opc), rotationBytes), additionConstant))
	if err != nil {
		return nil, err
	}

	// ak = f5* = OUT[0] ... OUT[47]
	return xor(out, opc)[:6], nil
}

// generateAutn generates network authentication tokens as defined in 3GPP 25.205 7.2
//
// Inputs:
//   sqn: 48 bit sequence number
//   ak: 48 bit anonymity key
//   macA: 64 bit network authentication code
//   amf: 16 bit authentication management field
// Outputs: 128 bit authentication token
func generateAutn(sqn, ak, macA, amf []byte) []byte {
	autn := make([]byte, 0, AutnBytes)
	autn = append(autn, xor(sqn, ak)...)
	autn = append(autn, amf...)
	autn = append(autn, macA...)
	return autn
}

// generateKasme is the KASME derivation function (S_2) according to 3GPP 33.401 Annex A.2.
// This function creates an input string to a key derivation function.
//
// The input string to the KDF is composed of 2 input parameters P0, P1
// and their lengths L0, L1 a constant FC which identifies this algorithm.
//				S = FC || P0 || L0 || P1 || L1
// The FC = 0x10 and argument P0 is the 3 octets of the PLMN, and P1 is
// SQN XOR AK. The lengths are in bytes.
//
// The Kasme is computed by calling the key derivation function with S
// using key CK || IK
//
// Inputs:
//   ck: 128 bit confidentiality key
//   ik: 128 bit integrity key
//   plmn: 24 bit network identifier
//      Octet           Description
//         1      MCC digit 2 | MCC digit 1
//         2      MNC digit 3 | MCC digit 3
//         3      MNC digit 2 | MNC digit 1
//   sqn: 48 bit sequence number
//   ak: 48 bit anonymity key
// Outputs: 256 bit network base key or an error
func generateKasme(ck, ik, plmn, sqn, ak []byte) ([]byte, error) {
	const fc = 16 // identifies the algorithm
	const inputBytes = 14

	var msg = make([]byte, inputBytes)
	msg[0] = fc
	copy(msg[1:], plmn)
	msg[5] = ExpectedPlmnBytes
	copy(msg[6:], xor(sqn, ak))
	msg[13] = sqnMaxBytes
	key := append(ck, ik...)

	// 3GPP Key Derivation Function defined in TS 33.220 to be hmac-sha256
	hash := hmac.New(sha256.New, key)
	_, err := hash.Write(msg)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// getSqnBytes encodes sqn in a byte slice.
func getSqnBytes(sqn uint64) []byte {
	const uint64Bytes = 8
	sqnBytes := make([]byte, uint64Bytes)
	binary.BigEndian.PutUint64(sqnBytes, sqn)
	return sqnBytes[uint64Bytes-sqnMaxBytes:]
}

// encrypt implements the Rijndael (AES-128) cipher function used by Milenage
// Inputs:
//   key: 128 bit encryption key
//   buf: 128 bit buffer to encrypt
// Outputs: encrypted output or an error
func encrypt(key, buf []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(aesCipher, make([]byte, aes.BlockSize))
	output := make([]byte, len(buf))
	encrypter.CryptBlocks(output, buf)
	return output, nil
}

// xor xors the bytes in a and b.
// If len(b) > len(a), then this function will panic. Otherwise, this function
// will only run on the first len(a) bytes of each input slice.
// Inputs: The two byte arrays to be xor'd.
// Output: The xor'd result
func xor(a, b []byte) []byte {
	n := len(a)
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}

// rotate a byte array by a number (k >= 0) of bytes.
func rotate(arr []byte, k int) []byte {
	n := len(arr)
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = arr[(i+k)%n]
	}
	return dst
}

// newEutranVector creates an EutranVector by copying in the given slices.
func newEutranVector(rand, xres, autn, kasme []byte) *EutranVector {
	var eutran = &EutranVector{}
	copy(eutran.Rand[:], rand)
	copy(eutran.Xres[:], xres)
	copy(eutran.Autn[:], autn)
	copy(eutran.Kasme[:], kasme)
	return eutran
}

// newSIPAuthVector creates a SIP auth vector by copying in the given slices.
func newSIPAuthVector(rand, xres, autn, ck, ik, ak []byte) *SIPAuthVector {
	var vector = &SIPAuthVector{}
	copy(vector.Rand[:], rand)
	copy(vector.Xres[:], xres)
	copy(vector.Autn[:], autn)
	copy(vector.ConfidentialityKey[:], ck)
	copy(vector.IntegrityKey[:], ik)
	copy(vector.AnonymityKey[:], ak)
	return vector
}

// cryptoRNG allows reading random bytes
type cryptoRNG interface {
	Read([]byte) (int, error)
}

// defaultCryptoRNG is a type which forwards crypto/Rand's Read function.
type defaultCryptoRNG struct{}

func (defaultCryptoRNG) Read(b []byte) (int, error) {
	return rand.Read(b)
}
