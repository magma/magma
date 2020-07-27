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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEutranVector(t *testing.T) {
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	key := []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb")
	sqn := uint64(7351)
	opc := []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]")
	amf := []byte("\x80\x00")
	plmn := []byte("\x02\xf8\x59")

	milenage, err := NewMockMilenageCipher(amf, rand)
	assert.NoError(t, err)

	eutran, err := milenage.GenerateEutranVector(key, opc, sqn, plmn)
	assert.NoError(t, err)
	assert.Equal(t, rand, eutran.Rand[:])
	assert.Equal(t, []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"), eutran.Xres[:])
	assert.Equal(t, []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"), eutran.Autn[:])
	assert.Equal(t, []byte("\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2"), eutran.Kasme[:])
}

func TestGenerateSIPAuthVector(t *testing.T) {
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	key := []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb")
	sqn := uint64(7351)
	opc := []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]")
	amf := []byte("\x80\x00")

	milenage, err := NewMockMilenageCipher(amf, rand)
	assert.NoError(t, err)

	vector, err := milenage.GenerateSIPAuthVector(key, opc, sqn)
	assert.NoError(t, err)
	assert.Equal(t, rand, vector.Rand[:])
	assert.Equal(t, []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"), vector.Xres[:])
	assert.Equal(t, []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"), vector.Autn[:])
	assert.Equal(t, []byte{0xf0, 0x6e, 0x32, 0xf9, 0x13, 0xee, 0xfb, 0x49, 0xfb, 0x72, 0xf1, 0x9, 0xb3, 0xa5, 0xf3, 0xc8}, vector.ConfidentialityKey[:])
	assert.Equal(t, []byte{0xb0, 0x6a, 0x7b, 0x46, 0xf, 0x4f, 0x53, 0xc4, 0x16, 0x6b, 0xf4, 0xa2, 0xe0, 0xa0, 0xc2, 0x5c}, vector.IntegrityKey[:])
	assert.Equal(t, []byte{0x6f, 0xbf, 0xa3, 0x80, 0x3, 0xe0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, vector.AnonymityKey[:])
}

func TestGenerateSIPAuthVectorWithRand(t *testing.T) {
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	key := []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb")
	sqn := uint64(7351)
	opc := []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]")
	amf := []byte("\x80\x00")

	milenage, err := NewMilenageCipher(amf)
	assert.NoError(t, err)

	vector, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc, sqn)
	assert.NoError(t, err)
	assert.Equal(t, rand, vector.Rand[:])
	assert.Equal(t, []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"), vector.Xres[:])
	assert.Equal(t, []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"), vector.Autn[:])
	assert.Equal(t, []byte{0xf0, 0x6e, 0x32, 0xf9, 0x13, 0xee, 0xfb, 0x49, 0xfb, 0x72, 0xf1, 0x9, 0xb3, 0xa5, 0xf3, 0xc8}, vector.ConfidentialityKey[:])
	assert.Equal(t, []byte{0xb0, 0x6a, 0x7b, 0x46, 0xf, 0x4f, 0x53, 0xc4, 0x16, 0x6b, 0xf4, 0xa2, 0xe0, 0xa0, 0xc2, 0x5c}, vector.IntegrityKey[:])
	assert.Equal(t, []byte{0x6f, 0xbf, 0xa3, 0x80, 0x3, 0xe0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, vector.AnonymityKey[:])
}

func TestGenerateEutranVectorError(t *testing.T) {
	amf := []byte("\x80\x00")
	milenage, err := NewMilenageCipher(amf)
	assert.NoError(t, err)

	eutran, err := milenage.GenerateEutranVector(nil, nil, 0, nil)
	assert.Nil(t, eutran)
	assert.Error(t, err)
}

func TestValidateGenerateEutranVectorInputs(t *testing.T) {
	err := validateGenerateEutranVectorInputs(nil, nil, 0, nil)
	assert.Error(t, err)

	var key = make([]byte, ExpectedKeyBytes)
	var opc = make([]byte, ExpectedOpcBytes)
	var plmn = make([]byte, ExpectedPlmnBytes)
	err = validateGenerateEutranVectorInputs(key, opc, 0, plmn)
	assert.NoError(t, err)

	var badKey = make([]byte, ExpectedKeyBytes*2)
	err = validateGenerateEutranVectorInputs(badKey, opc, 0, plmn)
	assert.Error(t, err)

	var badOpc = make([]byte, ExpectedOpcBytes/2)
	err = validateGenerateEutranVectorInputs(key, badOpc, 0, plmn)
	assert.Error(t, err)

	var badPlmn = make([]byte, ExpectedPlmnBytes+1)
	err = validateGenerateEutranVectorInputs(key, opc, 0, badPlmn)
	assert.Error(t, err)

	err = validateGenerateEutranVectorInputs(key, opc, maxSqn+1, plmn)
	assert.Error(t, err)
}

func TestValidateGenerateSIPAuthVectorInputs(t *testing.T) {
	err := validateGenerateSIPAuthVectorInputs(nil, nil, 0)
	assert.Error(t, err)

	var key = make([]byte, ExpectedKeyBytes)
	var opc = make([]byte, ExpectedOpcBytes)
	err = validateGenerateSIPAuthVectorInputs(key, opc, 0)
	assert.NoError(t, err)

	var badKey = make([]byte, ExpectedKeyBytes*2)
	err = validateGenerateSIPAuthVectorInputs(badKey, opc, 0)
	assert.EqualError(t, err, "incorrect key size. Expected 16 bytes, but got 32 bytes")

	var badOpc = make([]byte, ExpectedOpcBytes/2)
	err = validateGenerateSIPAuthVectorInputs(key, badOpc, 0)
	assert.EqualError(t, err, "incorrect opc size. Expected 16 bytes, but got 8 bytes")

	err = validateGenerateSIPAuthVectorInputs(key, opc, maxSqn+1)
	assert.EqualError(t, err, "sequence number too large, expected a number which can fit in 48 bits. Got: 281474976710656")
}

func TestValidateGenerateSIPAuthVectorWithRandInputs(t *testing.T) {
	err := validateGenerateSIPAuthVectorWithRandInputs(nil, nil, nil, 0)
	assert.Error(t, err)

	var rand = make([]byte, RandChallengeBytes)
	var key = make([]byte, ExpectedKeyBytes)
	var opc = make([]byte, ExpectedOpcBytes)
	err = validateGenerateSIPAuthVectorWithRandInputs(rand, key, opc, 0)
	assert.NoError(t, err)

	var badRand = make([]byte, RandChallengeBytes*2)
	err = validateGenerateSIPAuthVectorWithRandInputs(badRand, key, opc, 0)
	assert.EqualError(t, err, "incorrect rand size. Expected 16 bytes, but got 32 bytes")

	var badKey = make([]byte, ExpectedKeyBytes*2)
	err = validateGenerateSIPAuthVectorWithRandInputs(rand, badKey, opc, 0)
	assert.EqualError(t, err, "incorrect key size. Expected 16 bytes, but got 32 bytes")

	var badOpc = make([]byte, ExpectedOpcBytes/2)
	err = validateGenerateSIPAuthVectorWithRandInputs(rand, key, badOpc, 0)
	assert.EqualError(t, err, "incorrect opc size. Expected 16 bytes, but got 8 bytes")

	err = validateGenerateSIPAuthVectorWithRandInputs(rand, key, opc, maxSqn+1)
	assert.EqualError(t, err, "sequence number too large, expected a number which can fit in 48 bits. Got: 281474976710656")
}

func TestNewMilenageError(t *testing.T) {
	amf := []byte("\x80")
	_, err := NewMilenageCipher(amf)
	assert.Error(t, err)

	amf = []byte("\x80\x80\x80")
	_, err = NewMilenageCipher(amf)
	assert.Error(t, err)

	_, err = NewMilenageCipher(nil)
	assert.Error(t, err)

	amf = []byte("\x80\x80")
	_, err = NewMilenageCipher(amf)
	assert.NoError(t, err)
}

// This is the test from set 1 from 3GPP 35.207 4.3
func TestF1_Set1(t *testing.T) {
	rand, err := hex.DecodeString("23553cbe9637a89d218ae64dae47bf35")
	assert.NoError(t, err)

	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	sqn := []byte("\xff\x9b\xb4\xd0\xb6\x07")
	amf := []byte("\xb9\xb9")
	opc := []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\x4a\x9f\xfa\xc3\x54\xdf\xaf\xb3"), macA)
	assert.Equal(t, []byte("\x01\xcf\xaf\x9e\xc4\xe8\x71\xe9"), macS)
}

// This is the test from set 2 from 3GPP 35.207 4.3
func TestF1_Set2(t *testing.T) {
	rand, err := hex.DecodeString("c00d603103dcee52c4478119494202e8")
	assert.NoError(t, err)
	key, err := hex.DecodeString("0396eb317b6d1c36f19c1c84cd6ffd16")
	assert.NoError(t, err)
	sqn, err := hex.DecodeString("fd8eef40df7d")
	assert.NoError(t, err)
	amf, err := hex.DecodeString("af17")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("53c15671c60a4b731c55b4a441c0bde2")
	assert.NoError(t, err)
	expectedMacA, err := hex.DecodeString("5df5b31807e258b0")
	assert.NoError(t, err)
	expectedMacS, err := hex.DecodeString("a8c016e51ef4a343")
	assert.NoError(t, err)

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, expectedMacA, macA)
	assert.Equal(t, expectedMacS, macS)
}

// This is the test from set 3 from 3GPP 35.207 4.3
func TestF1_Set3(t *testing.T) {
	rand, err := hex.DecodeString("9f7c8d021accf4db213ccff0c7f71a6a")
	assert.NoError(t, err)
	key, err := hex.DecodeString("fec86ba6eb707ed08905757b1bb44b8f")
	assert.NoError(t, err)
	sqn, err := hex.DecodeString("9d0277595ffc")
	assert.NoError(t, err)
	amf, err := hex.DecodeString("725c")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("1006020f0a478bf6b699f15c062e42b3")
	assert.NoError(t, err)
	expectedMacA, err := hex.DecodeString("9cabc3e99baf7281")
	assert.NoError(t, err)
	expectedMacS, err := hex.DecodeString("95814ba2b3044324")
	assert.NoError(t, err)

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, expectedMacA, macA)
	assert.Equal(t, expectedMacS, macS)
}

// This is the test from set 4 from 3GPP 35.207 4.3
func TestF1_Set4(t *testing.T) {
	rand, err := hex.DecodeString("ce83dbc54ac0274a157c17f80d017bd6")
	assert.NoError(t, err)
	key, err := hex.DecodeString("9e5944aea94b81165c82fbf9f32db751")
	assert.NoError(t, err)
	sqn, err := hex.DecodeString("0b604a81eca8")
	assert.NoError(t, err)
	amf, err := hex.DecodeString("9e09")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("a64a507ae1a2a98bb88eb4210135dc87")
	assert.NoError(t, err)
	expectedMacA, err := hex.DecodeString("74a58220cba84c49")
	assert.NoError(t, err)
	expectedMacS, err := hex.DecodeString("ac2cc74a96871837")
	assert.NoError(t, err)

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, expectedMacA, macA)
	assert.Equal(t, expectedMacS, macS)
}

// This is the test from set 5 from 3GPP 35.207 4.3
func TestF1_Set5(t *testing.T) {
	rand, err := hex.DecodeString("74b0cd6031a1c8339b2b6ce2b8c4a186")
	assert.NoError(t, err)
	key, err := hex.DecodeString("4ab1deb05ca6ceb051fc98e77d026a84")
	assert.NoError(t, err)
	sqn, err := hex.DecodeString("e880a1b580b6")
	assert.NoError(t, err)
	amf, err := hex.DecodeString("9f07")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("dcf07cbd51855290b92a07a9891e523e")
	assert.NoError(t, err)
	expectedMacA, err := hex.DecodeString("49e785dd12626ef2")
	assert.NoError(t, err)
	expectedMacS, err := hex.DecodeString("9e85790336bb3fa2")
	assert.NoError(t, err)

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, expectedMacA, macA)
	assert.Equal(t, expectedMacS, macS)
}

// This is the test from set 6 from 3GPP 35.207 4.3
func TestF1_Set6(t *testing.T) {
	rand, err := hex.DecodeString("ee6466bc96202c5a557abbeff8babf63")
	assert.NoError(t, err)
	key, err := hex.DecodeString("6c38a116ac280c454f59332ee35c8c4f")
	assert.NoError(t, err)
	sqn, err := hex.DecodeString("414b98222181")
	assert.NoError(t, err)
	amf, err := hex.DecodeString("4464")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("3803ef5363b947c6aaa225e58fae3934")
	assert.NoError(t, err)
	expectedMacA, err := hex.DecodeString("078adfb488241a57")
	assert.NoError(t, err)
	expectedMacS, err := hex.DecodeString("80246b8d0186bcf1")
	assert.NoError(t, err)

	macA, macS, err := f1(key, sqn, rand, opc, amf)
	assert.NoError(t, err)
	assert.Equal(t, expectedMacA, macA)
	assert.Equal(t, expectedMacS, macS)
}

//  This is the test from set 1 from 3GPP 35.207 5.3
func TestF2F5_Set1(t *testing.T) {
	rand := []byte("#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5")
	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	opc := []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\xa5\x42\x11\xd5\xe3\xba\x50\xbf"), xres)
	assert.Equal(t, []byte("\xaa\x68\x9c\x64\x83\x70"), ak)
}

//  This is the test from set 2 from 3GPP 35.207 5.3
func TestF2F5_Set2(t *testing.T) {
	rand, err := hex.DecodeString("c00d603103dcee52c4478119494202e8")
	assert.NoError(t, err)
	key, err := hex.DecodeString("0396eb317b6d1c36f19c1c84cd6ffd16")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("53c15671c60a4b731c55b4a441c0bde2")
	assert.NoError(t, err)
	expectedXres, err := hex.DecodeString("d3a628ed988620f0")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("c47783995f72")
	assert.NoError(t, err)

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedXres, xres)
	assert.Equal(t, expectedAk, ak)
}

//  This is the test from set 3 from 3GPP 35.207 5.3
func TestF2F5_Set3(t *testing.T) {
	rand, err := hex.DecodeString("9f7c8d021accf4db213ccff0c7f71a6a")
	assert.NoError(t, err)
	key, err := hex.DecodeString("fec86ba6eb707ed08905757b1bb44b8f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("1006020f0a478bf6b699f15c062e42b3")
	assert.NoError(t, err)
	expectedXres, err := hex.DecodeString("8011c48c0c214ed2")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("33484dc2136b")
	assert.NoError(t, err)

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedXres, xres)
	assert.Equal(t, expectedAk, ak)
}

//  This is the test from set 4 from 3GPP 35.207 5.3
func TestF2F5_Set4(t *testing.T) {
	rand, err := hex.DecodeString("ce83dbc54ac0274a157c17f80d017bd6")
	assert.NoError(t, err)
	key, err := hex.DecodeString("9e5944aea94b81165c82fbf9f32db751")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("a64a507ae1a2a98bb88eb4210135dc87")
	assert.NoError(t, err)
	expectedXres, err := hex.DecodeString("f365cd683cd92e96")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("f0b9c08ad02e")
	assert.NoError(t, err)

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedXres, xres)
	assert.Equal(t, expectedAk, ak)
}

//  This is the test from set 5 from 3GPP 35.207 5.3
func TestF2F5_Set5(t *testing.T) {
	rand, err := hex.DecodeString("74b0cd6031a1c8339b2b6ce2b8c4a186")
	assert.NoError(t, err)
	key, err := hex.DecodeString("4ab1deb05ca6ceb051fc98e77d026a84")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("dcf07cbd51855290b92a07a9891e523e")
	assert.NoError(t, err)
	expectedXres, err := hex.DecodeString("5860fc1bce351e7e")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("31e11a609118")
	assert.NoError(t, err)

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedXres, xres)
	assert.Equal(t, expectedAk, ak)
}

//  This is the test from set 6 from 3GPP 35.207 5.3
func TestF2F5_Set6(t *testing.T) {
	rand, err := hex.DecodeString("ee6466bc96202c5a557abbeff8babf63")
	assert.NoError(t, err)
	key, err := hex.DecodeString("6c38a116ac280c454f59332ee35c8c4f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("3803ef5363b947c6aaa225e58fae3934")
	assert.NoError(t, err)
	expectedXres, err := hex.DecodeString("16c8233f05a0ac28")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("45b0f69ab06c")
	assert.NoError(t, err)

	xres, ak, err := f2F5(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedXres, xres)
	assert.Equal(t, expectedAk, ak)
}

//  This is the test from set 1 from 3GPP 35.207 5.3
func TestF3_Set1(t *testing.T) {
	rand := []byte("#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5")
	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	opc := []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\xb4\x0b\xa9\xa3\xc5\x8b\x2a\x05\xbb\xf0\xd9\x87\xb2\x1b\xf8\xcb"), ck)
}

//  This is the test from set 2 from 3GPP 35.207 5.3
func TestF3_Set2(t *testing.T) {
	rand, err := hex.DecodeString("c00d603103dcee52c4478119494202e8")
	assert.NoError(t, err)
	key, err := hex.DecodeString("0396eb317b6d1c36f19c1c84cd6ffd16")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("53c15671c60a4b731c55b4a441c0bde2")
	assert.NoError(t, err)
	expectedCk, err := hex.DecodeString("58c433ff7a7082acd424220f2b67c556")
	assert.NoError(t, err)

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedCk, ck)
}

//  This is the test from set 3 from 3GPP 35.207 5.3
func TestF3_Set3(t *testing.T) {
	rand, err := hex.DecodeString("9f7c8d021accf4db213ccff0c7f71a6a")
	assert.NoError(t, err)
	key, err := hex.DecodeString("fec86ba6eb707ed08905757b1bb44b8f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("1006020f0a478bf6b699f15c062e42b3")
	assert.NoError(t, err)
	expectedCk, err := hex.DecodeString("5dbdbb2954e8f3cde665b046179a5098")
	assert.NoError(t, err)

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedCk, ck)
}

//  This is the test from set 4 from 3GPP 35.207 5.3
func TestF3_Set4(t *testing.T) {
	rand, err := hex.DecodeString("ce83dbc54ac0274a157c17f80d017bd6")
	assert.NoError(t, err)
	key, err := hex.DecodeString("9e5944aea94b81165c82fbf9f32db751")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("a64a507ae1a2a98bb88eb4210135dc87")
	assert.NoError(t, err)
	expectedCk, err := hex.DecodeString("e203edb3971574f5a94b0d61b816345d")
	assert.NoError(t, err)

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedCk, ck)
}

//  This is the test from set 5 from 3GPP 35.207 5.3
func TestF3_Set5(t *testing.T) {
	rand, err := hex.DecodeString("74b0cd6031a1c8339b2b6ce2b8c4a186")
	assert.NoError(t, err)
	key, err := hex.DecodeString("4ab1deb05ca6ceb051fc98e77d026a84")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("dcf07cbd51855290b92a07a9891e523e")
	assert.NoError(t, err)
	expectedCk, err := hex.DecodeString("7657766b373d1c2138f307e3de9242f9")
	assert.NoError(t, err)

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedCk, ck)
}

//  This is the test from set 6 from 3GPP 35.207 5.3
func TestF3_Set6(t *testing.T) {
	rand, err := hex.DecodeString("ee6466bc96202c5a557abbeff8babf63")
	assert.NoError(t, err)
	key, err := hex.DecodeString("6c38a116ac280c454f59332ee35c8c4f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("3803ef5363b947c6aaa225e58fae3934")
	assert.NoError(t, err)
	expectedCk, err := hex.DecodeString("3f8c7587fe8e4b233af676aede30ba3b")
	assert.NoError(t, err)

	ck, err := f3(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedCk, ck)
}

// This is the test from set 1 from 3GPP 35.207 6.3
func TestF4_Set1(t *testing.T) {
	rand := []byte("#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5")
	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	opc := []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\xf7\x69\xbc\xd7\x51\x04\x46\x04\x12\x76\x72\x71\x1c\x6d\x34\x41"), ik)
}

// This is the test from set 2 from 3GPP 35.207 6.3
func TestF4_Set2(t *testing.T) {
	rand, err := hex.DecodeString("c00d603103dcee52c4478119494202e8")
	assert.NoError(t, err)
	key, err := hex.DecodeString("0396eb317b6d1c36f19c1c84cd6ffd16")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("53c15671c60a4b731c55b4a441c0bde2")
	assert.NoError(t, err)
	expectedIk, err := hex.DecodeString("21a8c1f929702adb3e738488b9f5c5da")
	assert.NoError(t, err)

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedIk, ik)
}

// This is the test from set 3 from 3GPP 35.207 6.3
func TestF4_Set3(t *testing.T) {
	rand, err := hex.DecodeString("9f7c8d021accf4db213ccff0c7f71a6a")
	assert.NoError(t, err)
	key, err := hex.DecodeString("fec86ba6eb707ed08905757b1bb44b8f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("1006020f0a478bf6b699f15c062e42b3")
	assert.NoError(t, err)
	expectedIk, err := hex.DecodeString("59a92d3b476a0443487055cf88b2307b")
	assert.NoError(t, err)

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedIk, ik)
}

// This is the test from set 4 from 3GPP 35.207 6.3
func TestF4_Set4(t *testing.T) {
	rand, err := hex.DecodeString("ce83dbc54ac0274a157c17f80d017bd6")
	assert.NoError(t, err)
	key, err := hex.DecodeString("9e5944aea94b81165c82fbf9f32db751")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("a64a507ae1a2a98bb88eb4210135dc87")
	assert.NoError(t, err)
	expectedIk, err := hex.DecodeString("0c4524adeac041c4dd830d20854fc46b")
	assert.NoError(t, err)

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedIk, ik)
}

// This is the test from set 5 from 3GPP 35.207 6.3
func TestF4_Set5(t *testing.T) {
	rand, err := hex.DecodeString("74b0cd6031a1c8339b2b6ce2b8c4a186")
	assert.NoError(t, err)
	key, err := hex.DecodeString("4ab1deb05ca6ceb051fc98e77d026a84")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("dcf07cbd51855290b92a07a9891e523e")
	assert.NoError(t, err)
	expectedIk, err := hex.DecodeString("1c42e960d89b8fa99f2744e0708ccb53")
	assert.NoError(t, err)

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedIk, ik)
}

// This is the test from set 6 from 3GPP 35.207 6.3
func TestF4_Set6(t *testing.T) {
	rand, err := hex.DecodeString("ee6466bc96202c5a557abbeff8babf63")
	assert.NoError(t, err)
	key, err := hex.DecodeString("6c38a116ac280c454f59332ee35c8c4f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("3803ef5363b947c6aaa225e58fae3934")
	assert.NoError(t, err)
	expectedIk, err := hex.DecodeString("a7466cc1e6b2a1337d49d3b66e95d7b4")
	assert.NoError(t, err)

	ik, err := f4(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedIk, ik)
}

// This is the test from set 1 from 3GPP 35.207 6.3
func TestF5Star_Set1(t *testing.T) {
	rand := []byte("#U<\xbe\x967\xa8\x9d!\x8a\xe6M\xaeG\xbf5")
	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	opc := []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\x45\x1e\x8b\xec\xa4\x3b"), ak)
}

// This is the test from set 2 from 3GPP 35.207 6.3
func TestF5Star_Set2(t *testing.T) {
	rand, err := hex.DecodeString("c00d603103dcee52c4478119494202e8")
	assert.NoError(t, err)
	key, err := hex.DecodeString("0396eb317b6d1c36f19c1c84cd6ffd16")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("53c15671c60a4b731c55b4a441c0bde2")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("30f1197061c1")
	assert.NoError(t, err)

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedAk, ak)
}

// This is the test from set 3 from 3GPP 35.207 6.3
func TestF5Star_Set3(t *testing.T) {
	rand, err := hex.DecodeString("9f7c8d021accf4db213ccff0c7f71a6a")
	assert.NoError(t, err)
	key, err := hex.DecodeString("fec86ba6eb707ed08905757b1bb44b8f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("1006020f0a478bf6b699f15c062e42b3")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("deacdd848cc6")
	assert.NoError(t, err)

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedAk, ak)
}

// This is the test from set 4 from 3GPP 35.207 6.3
func TestF5Star_Set4(t *testing.T) {
	rand, err := hex.DecodeString("ce83dbc54ac0274a157c17f80d017bd6")
	assert.NoError(t, err)
	key, err := hex.DecodeString("9e5944aea94b81165c82fbf9f32db751")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("a64a507ae1a2a98bb88eb4210135dc87")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("6085a86c6f63")
	assert.NoError(t, err)

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedAk, ak)
}

// This is the test from set 5 from 3GPP 35.207 6.3
func TestF5Star_Set5(t *testing.T) {
	rand, err := hex.DecodeString("74b0cd6031a1c8339b2b6ce2b8c4a186")
	assert.NoError(t, err)
	key, err := hex.DecodeString("4ab1deb05ca6ceb051fc98e77d026a84")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("dcf07cbd51855290b92a07a9891e523e")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("fe2555e54aa9")
	assert.NoError(t, err)

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedAk, ak)
}

// This is the test from set 6 from 3GPP 35.207 6.3
func TestF5Star_Set6(t *testing.T) {
	rand, err := hex.DecodeString("ee6466bc96202c5a557abbeff8babf63")
	assert.NoError(t, err)
	key, err := hex.DecodeString("6c38a116ac280c454f59332ee35c8c4f")
	assert.NoError(t, err)
	opc, err := hex.DecodeString("3803ef5363b947c6aaa225e58fae3934")
	assert.NoError(t, err)
	expectedAk, err := hex.DecodeString("1f53cd2b1113")
	assert.NoError(t, err)

	ak, err := f5Star(key, rand, opc)
	assert.NoError(t, err)
	assert.Equal(t, expectedAk, ak)
}

func TestGenerateAutn(t *testing.T) {
	sqn := []byte("\x00\x01\x01\x00\xaa\x0a")
	ak := []byte("\x00\x00\x01\x01\xcc\x0c")
	macA := []byte("\x34\xdd\x01\x00\xde\x12\x76\x40")
	amf := []byte("\xab\xcd")
	autn := generateAutn(sqn, ak, macA, amf)
	assert.Equal(t, []byte("\x00\x01\x00\x01\x66\x06\xab\xcd\x34\xdd\x01\x00\xde\x12\x76\x40"), autn)
}

func TestGenerateKasme(t *testing.T) {
	ck := []byte("\xb4\x0b\xa9\xa3\xc5\x8b\x2a\x05\xbb\xf0\xd9\x87\xb2\x1b\xf8\xcb")
	ik := []byte("\xf7\x69\xbc\xd7\x51\x04\x46\x04\x12\x76\x72\x71\x1c\x6d\x34\x41")
	sqn := []byte("\x00\x01\x01\x00\xaa\x0a")
	ak := []byte("\x00\x00\x01\x01\xcc\x0c")
	plmn := []byte("\x02\xf8\x59")

	kasme, err := generateKasme(ck, ik, plmn, sqn, ak)
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0xd0, 0x87, 0x64, 0xb4, 0xfd, 0x45, 0xfa, 0x3b, 0x6, 0x59, 0x79,
			0x18, 0x94, 0xb6, 0x73, 0x8b, 0x23, 0x3, 0x24, 0x1e, 0x14, 0x7a,
			0xef, 0x80, 0x4f, 0x6c, 0x53, 0xd1, 0xf0, 0x3b, 0xd1, 0xea}, kasme)
}

func TestGenerateOpc(t *testing.T) {
	key := []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")
	op := []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18")

	opc, err := GenerateOpc(key, op)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf"), opc[:])
}

func TestGenerateOpc_InvalidInput(t *testing.T) {
	var key = make([]byte, ExpectedKeyBytes)
	var op = make([]byte, ExpectedOpBytes)
	var invalidKey = make([]byte, ExpectedKeyBytes-1)
	var invalidOp = make([]byte, ExpectedOpBytes+1)

	_, err := GenerateOpc(key, op)
	assert.NoError(t, err)

	_, err = GenerateOpc(invalidKey, op)
	assert.EqualError(t, err, "incorrect key size. Expected 16 bytes, but got 15 bytes")

	_, err = GenerateOpc(key, invalidOp)
	assert.EqualError(t, err, "incorrect op size. Expected 16 bytes, but got 17 bytes")
}

func TestGenerateResync(t *testing.T) {
	rand := []byte("\xcd\x14\xa7S\x97\x7f\xbcq\x8eb\xbd\xdbS]\x88\xf8")
	key := []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb")
	opc := []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]")
	amf := []byte{0, 0}
	auts := []byte{236, 25, 14, 177, 16, 88, 219, 95, 99, 96, 89, 31, 52, 234}

	milenage, err := NewMilenageCipher(amf)
	assert.NoError(t, err)

	sqn, macS, err := milenage.GenerateResync(auts, key, opc, rand)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), sqn)
	assert.Equal(t, []byte("\xdb_c`Y\x1f4\xea"), macS[:])
}

func TestGenerateResync_InvalidInput(t *testing.T) {
	rand := make([]byte, RandChallengeBytes)
	key := make([]byte, ExpectedKeyBytes)
	opc := make([]byte, ExpectedOpcBytes)
	auts := make([]byte, ExpectedAutsBytes)
	amf := make([]byte, ExpectedAmfBytes)

	invalidRand := make([]byte, RandChallengeBytes+1)
	invalidKey := make([]byte, ExpectedKeyBytes-1)
	invalidOpc := make([]byte, ExpectedOpcBytes*2)
	invalidAuts := make([]byte, ExpectedAutsBytes/2)

	milenage, err := NewMilenageCipher(amf)
	assert.NoError(t, err)

	_, _, err = milenage.GenerateResync(auts, key, opc, rand)
	assert.NoError(t, err)

	_, _, err = milenage.GenerateResync(auts, key, opc, invalidRand)
	assert.EqualError(t, err, "incorrect rand size. Expected 16 bytes, but got 17 bytes")

	_, _, err = milenage.GenerateResync(auts, invalidKey, opc, rand)
	assert.EqualError(t, err, "incorrect key size. Expected 16 bytes, but got 15 bytes")

	_, _, err = milenage.GenerateResync(auts, key, invalidOpc, rand)
	assert.EqualError(t, err, "incorrect opc size. Expected 16 bytes, but got 32 bytes")

	_, _, err = milenage.GenerateResync(invalidAuts, key, opc, rand)
	assert.EqualError(t, err, "incorrect auts size. Expected 14 bytes, but got 7 bytes")
}

func TestValidateGenerateResyncInputs(t *testing.T) {
	rand := make([]byte, RandChallengeBytes)
	key := make([]byte, ExpectedKeyBytes)
	opc := make([]byte, ExpectedOpcBytes)
	auts := make([]byte, ExpectedAutsBytes)

	invalidRand := make([]byte, RandChallengeBytes+1)
	invalidKey := make([]byte, ExpectedKeyBytes-1)
	invalidOpc := make([]byte, ExpectedOpcBytes*2)
	invalidAuts := make([]byte, ExpectedAutsBytes/2)

	err := validateGenerateResyncInputs(auts, key, opc, rand)
	assert.NoError(t, err)

	err = validateGenerateResyncInputs(auts, key, opc, invalidRand)
	assert.EqualError(t, err, "incorrect rand size. Expected 16 bytes, but got 17 bytes")

	err = validateGenerateResyncInputs(auts, invalidKey, opc, rand)
	assert.EqualError(t, err, "incorrect key size. Expected 16 bytes, but got 15 bytes")

	err = validateGenerateResyncInputs(auts, key, invalidOpc, rand)
	assert.EqualError(t, err, "incorrect opc size. Expected 16 bytes, but got 32 bytes")

	err = validateGenerateResyncInputs(invalidAuts, key, opc, rand)
	assert.EqualError(t, err, "incorrect auts size. Expected 14 bytes, but got 7 bytes")
}

func TestXor(t *testing.T) {
	a := []byte("\x00\x01\x01\x00\xaa")
	b := []byte("\x00\x00\x01\x01\xcc")
	assert.Equal(t, []byte("\x00\x01\x00\x01\x66"), xor(a, b))
}

func TestRotate(t *testing.T) {
	arr := []byte("\x00\x01\x02\x03")
	assert.Equal(t, arr, rotate(arr, 0))
	assert.Equal(t, []byte("\x01\x02\x03\x00"), rotate(arr, 1))
	assert.Equal(t, []byte("\x02\x03\x00\x01"), rotate(arr, 2))
	assert.Equal(t, []byte("\x03\x00\x01\x02"), rotate(arr, 3))
	assert.Equal(t, arr, rotate(arr, 4))
	assert.Equal(t, []byte("\x01\x02\x03\x00"), rotate(arr, 5))
}
