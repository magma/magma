/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicers

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"reflect"

	"magma/cwf/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/lte/cloud/go/crypto"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	defaultInd     = 0
	CheckcodeValue = "\x00\x00\x86\xe8\x20\x4d\xc6\xe1\xe3\xd8\x94\x44\x3c\x26" +
		"\xa7\xc6\x5d\xee\x3c\x42\xab\xf8"

	SqnLen    = 6
	MacAStart = 8

	// maxSeqDelta is the maximum allowed increase to SEQ.
	// eg. if x was the last accepted SEQ, then the next SEQ must
	// be greater than x and less than (x + maxSeqDelta) to be accepted.
	// See 3GPP TS 33.102 Appendix C.2.1.
	maxSeqDelta = 1 << 28
)

// handleEapAka routes the EAP-AKA request to the UE with the specified imsi.
func (srv *UESimServer) handleEapAka(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
	switch aka.Subtype(req[eap.EapSubtype]) {
	case aka.SubtypeIdentity:
		return srv.eapAkaIdentityRequest(ue, req)
	case aka.SubtypeChallenge:
		return srv.eapAkaChallengeRequest(ue, req)
	default:
		return nil, errors.Errorf("Unsupported Subtype: %d", req[eap.EapSubtype])
	}
}

// Given a UE and the EAP-AKA identity request, generates the EAP response.
func (srv *UESimServer) eapAkaIdentityRequest(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new attribute scanner")
	}

	var a eap.Attribute

	// Parse out attributes.
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_PERMANENT_ID_REQ, aka.AT_ANY_ID_REQ:
			// Create the response EAP packet with the identity attribute.
			p := eap.NewPacket(
				eap.ResponseCode,
				req.Identifier(),
				[]byte{aka.TYPE, byte(aka.SubtypeIdentity), 0, 0},
			)

			// Append Identity Attribute data to packet.
			id := []byte("\x30" + ue.GetImsi() + IdentityPostfix)
			p, err = p.Append(
				eap.NewAttribute(
					aka.AT_IDENTITY,
					append(
						[]byte{uint8(len(id) >> 8), uint8(len(id))}, // actual len of Identity
						id...,
					),
				),
			)
			if err != nil {
				return nil, errors.Wrap(err, "Error appending attribute to packet")
			}
			return p, nil
		default:
			glog.Info(fmt.Sprintf("Unexpected EAP-AKA Identity Request Attribute type %d", a.Type()))
		}
	}
	return nil, errors.Wrap(err, "Error while processing EAP-AKA Identity Request")
}

type challengeAttributes struct {
	rand eap.Attribute
	autn eap.Attribute
	mac  eap.Attribute
}

// Given a UE, the Op, the Amf, and the EAP challenge, generates the EAP response.
func (srv *UESimServer) eapAkaChallengeRequest(ue *protos.UEConfig, req eap.Packet) (eap.Packet, error) {
	attrs, err := parseChallengeAttributes(req)
	if err != io.EOF {
		return nil, errors.Wrap(err, "Error while parsing attributes of request packet")
	}
	if attrs.rand == nil || attrs.autn == nil || attrs.mac == nil {
		return nil, errors.Errorf("Missing one or more expected attributes\nRAND: %s\nAUTN: %s\nMAC: %s\n", attrs.rand, attrs.autn, attrs.mac)
	}

	// Parse out RAND, expected AUTN, and expected MAC values.
	rand := attrs.rand.Marshaled()[aka.ATT_HDR_LEN:]
	expectedAutn := attrs.autn.Marshaled()[aka.ATT_HDR_LEN:]
	expectedMac := attrs.mac.Marshaled()[aka.ATT_HDR_LEN:]

	id := []byte("\x30" + ue.GetImsi() + IdentityPostfix)
	key := []byte(ue.AuthKey)

	// Calculate SQN using SEQ and arbitrary IND
	sqn := servicers.SeqToSqn(ue.Seq, defaultInd)

	// Calculate Opc using key and Op, and verify that it matches the UE's Opc
	opc, err := crypto.GenerateOpc(key, srv.cfg.op)
	if err != nil {
		return nil, fmt.Errorf("Error while calculating Opc")
	}
	if !reflect.DeepEqual(opc[:], ue.AuthOpc) {
		return nil, fmt.Errorf("Invalid Opc: Expected Opc: %x; Actual Opc: %x", opc[:], ue.AuthOpc)
	}

	// Calculate RES and other keys.
	milenage, err := crypto.NewMilenageCipher(srv.cfg.amf)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating milenage cipher")
	}
	intermediateVec, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc[:], sqn)
	if err != nil {
		return nil, errors.Wrap(err, "Error calculating authentication vector")
	}
	// Make copy of packet and zero out MAC value.
	copyReq := make([]byte, len(req))
	copy(copyReq, req)
	copyAttrs, err := parseChallengeAttributes(eap.Packet(copyReq))
	if err != io.EOF {
		return nil, errors.Wrap(err, "Error while parsing attributes of copied request packet")
	}
	copyMacBytes := copyAttrs.mac.Marshaled()
	for i := aka.ATT_HDR_LEN; i < len(copyMacBytes); i++ {
		copyMacBytes[i] = 0
	}

	// Calculate and verify MAC.
	_, kAut, _, _ := aka.MakeAKAKeys(id, intermediateVec.IntegrityKey[:], intermediateVec.ConfidentialityKey[:])
	mac := aka.GenMac(copyReq, kAut)
	if !reflect.DeepEqual(expectedMac, mac) {
		return nil, fmt.Errorf("Invalid MAC: Expected MAC: %x; Actual MAC: %x", expectedMac, mac)
	}

	// Verify AUTN (MacA must be equal and SEQ in correct range)
	receivedSqn := extractSqnFromAutn(expectedAutn, intermediateVec.AnonymityKey[:])
	resultVec, err := milenage.GenerateSIPAuthVectorWithRand(rand, key, opc[:], receivedSqn)
	if err != nil {
		return nil, errors.Wrap(err, "Error calculating authentication vector")
	}
	if !reflect.DeepEqual(expectedAutn[MacAStart:], resultVec.Autn[MacAStart:]) {
		return nil, fmt.Errorf("Invalid MacA in AUTN: Received MacA %x; Calculated MacA: %x",
			expectedAutn[MacAStart:],
			resultVec.Autn[MacAStart:],
		)
	}
	seq, _ := servicers.SplitSqn(receivedSqn)
	isSeqValid := seq > ue.Seq && (seq-ue.GetSeq()) < maxSeqDelta
	if !isSeqValid {
		// TODO: Implement re-sync procedure
		// For now just return the error
		return nil, fmt.Errorf("Invalid SEQ received. HSS SEQ: %d, UE SEQ: %d", seq, ue.GetSeq())
	}

	// Update UE SEQ number
	ue.Seq = seq
	_, err = srv.AddUE(context.Background(), ue)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred while updating SEQ: %s", err)
	}

	// Create the response EAP packet.
	p := eap.NewPacket(eap.ResponseCode, req.Identifier(), []byte{aka.TYPE, byte(aka.SubtypeChallenge), 0, 0})

	// Add the RES attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_RES,
			append(
				[]byte{uint8(len(resultVec.Xres[:]) * 8 >> 8), uint8(len(resultVec.Xres[:]) * 8)},
				[]byte(resultVec.Xres[:])...,
			),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error appending attribute to packet")
	}

	// Add the CHECKCODE attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_CHECKCODE,
			[]byte(CheckcodeValue),
		),
	)

	atMacOffset := len(p) + aka.ATT_HDR_LEN

	// Add the empty MAC attribute.
	p, err = p.Append(
		eap.NewAttribute(
			aka.AT_MAC,
			append(make([]byte, 2+16)),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error appending attribute to packet")
	}

	// Calculate and Copy MAC into packet.
	mac = aka.GenMac(p, kAut)
	copy(p[atMacOffset:], mac)

	return p, nil
}

// Given an EAP packet, parses out the RAND, AUTN, and MAC.
func parseChallengeAttributes(req eap.Packet) (challengeAttributes, error) {
	attrs := challengeAttributes{}

	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		return attrs, errors.Wrap(err, "Error creating new attribute scanner")
	}
	var a eap.Attribute
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_RAND:
			attrs.rand = a
		case aka.AT_AUTN:
			attrs.autn = a
		case aka.AT_MAC:
			if len(a.Marshaled()) < aka.ATT_HDR_LEN+aka.MAC_LEN {
				return attrs, fmt.Errorf("Malformed AT_MAC")
			}
			attrs.mac = a
		default:
			glog.Info(fmt.Sprintf("Unexpected EAP-AKA Challenge Request Attribute type %d", a.Type()))
		}
	}
	return attrs, err
}

func extractSqnFromAutn(autn []byte, ak []byte) uint64 {
	sqn := xor(autn[:SqnLen], ak[:SqnLen])
	sqn64bits := make([]byte, 0, 8)
	sqn64bits = append(sqn64bits, []byte{0, 0}...)
	sqn64bits = append(sqn64bits, sqn...)
	return binary.BigEndian.Uint64(sqn64bits)
}

func xor(a, b []byte) []byte {
	n := len(a)
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}
