/*
 *  Copyright 2020 The Magma Authors.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package sms_ll

import (
	"fmt"
)

func smsRpError(field string) error {
	return fmt.Errorf("smsrp: %s", field)
}

// Used for both RP-Destination and RP-Originator Address, since it's the
// same layout (TS 24.011 8.2.5.1 and 8.2.5.2). Note that while 8.2.5.1-2
// suggest that the first octet is an IEI, 7.3.1 notes that this is a Type
// 4 LV IE, which means there's no IEI present -- just a length and values.
type rpAddressElement struct {
	length     byte // of the number in half-octets!
	numberInfo byte // octet 3
	number     []byte
}

// The length field of an RP Address Element is the number of half-octets in
// the number. This converts to a byte length of the number itself.
func (rpadde rpAddressElement) getNumberOctets() int {
	l := int(rpadde.length)
	if l%2 == 0 {
		return l / 2
	} else {
		return (l + 1) / 2 // last half-octet is padded with 0xf
	}
}

func (rpadde rpAddressElement) marshalBinary() []byte {
	if rpadde.length == 0x0 {
		return []byte{rpadde.length}
	}

	b := []byte{rpadde.length, rpadde.numberInfo}
	b = append(b, rpadde.number...)
	return b
}

// Decode an address element. Returns the length of the address element if present.
func (rpadde *rpAddressElement) unmarshalBinary(input []byte) (int, error) {
	// Empty addresses will be one byte long with a zero value length
	if len(input) == 1 {
		if input[0] != 0x0 {
			return -1, smsRpError("Invalid RP Address of length 1")
		} else {
			rpadde.length = input[0]
			return 1, nil
		}
	} else if len(input) < 3 { // if it's not zero length, we must have at least 3 octets
		return -1, smsRpError("Invalid RP Address")
	}

	rpadde.length = input[0]
	rpadde.numberInfo = input[1]

	num_bytes := rpadde.getNumberOctets()
	rpadde.number = make([]byte, num_bytes)
	copy(rpadde.number, input[2:num_bytes+2])

	return num_bytes + 2, nil
}

// The RP-Address-Element refers to the number of the SMSC. In our case, we
// generally don't use an SMSC, so we just set the number to 11.
func newFakeRpAddressElement() rpAddressElement {
	return rpAddressElement{
		numberInfo: 0xb9,         // network specific number, private numbering plan
		number:     []byte{0x11}, // 1 1
		length:     0x2,
	}
}

// RP-User data element (TS 24.011 8.2.5.3)
type rpUserElement struct {
	iei    byte // Not present for RP-DATA
	length byte
	tpdu   []byte
}

func createRpUserElement(data []byte) (rpUserElement, error) {
	if len(data) > 232 { // TS24.011 8.2.5.3
		return rpUserElement{}, smsRpError("UserData-Element too long (>232 bytes)")
	}

	return rpUserElement{
		iei:    RpUdeIei,
		length: byte(len(data)),
		tpdu:   data,
	}, nil

}

func (rpue rpUserElement) marshalBinary(msgType byte) []byte {
	rpu_len := len(rpue.tpdu) + 1
	b := make([]byte, 0, rpu_len)
	if msgType == RpAck || msgType == RpError { // these start with IEI
		b = append(b, rpue.iei)
	}
	b = append(b, rpue.length)
	b = append(b, rpue.tpdu...)
	return b
}

func (rpue *rpUserElement) unmarshalBinary(msgType byte, input []byte) int {
	idx := 0
	if msgType == RpAck || msgType == RpError { // these start with IEI
		rpue.iei = input[idx]
		idx++
	}
	rpue.length = input[idx]
	idx++

	end := idx + int(rpue.length)
	rpue.tpdu = make([]byte, rpue.length)
	copy(rpue.tpdu, input[idx:end])
	return end
}

// RP-Cause element (TS 24.011 8.2.5.4)
type rpCauseElement struct {
	iei        byte // Never serialized
	length     byte
	cause      byte
	diagnostic byte   // Optional
	causeStr   string // Derived
}

func (rpce rpCauseElement) marshalBinary() []byte {
	b := []byte{rpce.length, rpce.cause}
	if int(rpce.length) == 2 {
		b = append(b, rpce.diagnostic)
	}
	return b
}

func (rpce *rpCauseElement) unmarshalBinary(input []byte) (int, error) {
	if cs, ok := RpCauseStr[input[1]]; ok {
		rpce.cause = input[1]
		rpce.causeStr = cs
	} else {
		return 0, smsRpError(fmt.Sprintf("Invalid cause: %x", rpce.cause))
	}

	rpce.length = input[0]
	if int(rpce.length) == 2 {
		rpce.diagnostic = input[2]
		return 3, nil
	}
	return 2, nil
}

type rpMessage struct {
	mti       byte
	reference byte

	// Mandantory. If UE->Network, must be length 0
	originatorAddress rpAddressElement

	// Mandantory. If Network->UE, must be length 0
	destinationAddress rpAddressElement

	// Mandantory for RP-DATA, includes tpdu. If RP-ACK or RP-ERROR, must
	// include IEI.
	userData rpUserElement

	// Mandantory for RP-ERROR
	cause rpCauseElement
}

func (rpm rpMessage) marshalBinary() []byte {
	b := []byte{rpm.mti, rpm.reference}

	rpmt, _ := rpm.msgType()

	switch rpmt {
	case RpData:
		b = append(b, rpm.originatorAddress.marshalBinary()...)
		b = append(b, rpm.destinationAddress.marshalBinary()...)
		b = append(b, rpm.userData.marshalBinary(RpData)...)
	case RpError:
		b = append(b, rpm.cause.marshalBinary()...)
	case RpAck:
		// Do nothing
	}

	return b
}

func (rpm *rpMessage) unmarshalBinary(input []byte) error {
	if len(input) < 2 {
		return smsRpError("SMS-RP Message too short")
	}

	idx := 0
	rpm.mti = input[idx]
	idx++ // 1
	rpm.reference = input[idx]
	idx++ // 2

	rpmt, err := rpm.msgType()
	if err != nil {
		return err
	}

	switch rpmt {
	case RpData:
		// The next two IEs should be adddresses in this case. So, get the lengths and pass to unmarshal
		n, _ := rpm.originatorAddress.unmarshalBinary(input[idx:])
		if rpm.direction() == RpMo && n != 1 {
			return smsRpError("SMS-RP-DATA is MO, but OA length != 1")
		}
		idx += n
		n, _ = rpm.destinationAddress.unmarshalBinary(input[idx:])
		if rpm.direction() == RpMt && n != 1 {
			return smsRpError("SMS-RP-DATA is MT, but DA length != 1")
		}
		idx += n

		rpm.userData.unmarshalBinary(RpData, input[idx:])
	case RpAck:
		// RP-ACK and RP-ERROR may optionally contain an RP-User-Data
		// element (TS24.001 7.3.3). If this is the case, it will be a
		// TLV IE, with the first octet starting with the RP-User-Data
		// IE ID (0x41).
		if len(input) > 2 && input[idx] == RpUdeIei {
			rpm.userData.unmarshalBinary(RpAck, input[idx:])
		}
	case RpError:
		// Do nothing
		_, err := rpm.cause.unmarshalBinary(input[idx:])
		if err != nil {
			return err
		}
		// TODO: Add support for optional UserData TLV element.
	default:
		return smsRpError(fmt.Sprintf("Invalid RP-SMS MTI: 0x%08b", rpm.mti))
	}

	return nil
}

// If MTI is even, message is UE->Network
func (rpm rpMessage) direction() byte {
	if rpm.mti&0x1 == 0 {
		return RpMo
	}
	return RpMt
}

func (rpm rpMessage) msgType() (byte, error) {
	switch rpm.mti {
	case RpMtiMoData, RpMtiMtData:
		return RpData, nil
	case RpMtiMoErr, RpMtiMtErr:
		return RpError, nil
	case RpMtiMoAck, RpMtiMtAck:
		return RpAck, nil
	default:
		return RpInvalid, smsRpError(fmt.Sprintf("Invalid RP-SMS MTI: 0x%08b", rpm.mti))
	}
}

const (
	RpData = iota
	RpError
	RpAck
	RpInvalid // error type
	RpMo      // UE/MS->Network
	RpMt      // Network->UE/MS
)

const (
	// RP Message fields

	//RP-MTI (TS24.011 8.2.2) is technically defined as a 3 bit field in
	//the low order bits of the first octet of the RPDU. However, the five
	//high order bits are defined to always be 0, so here we treat these
	//fields as a full octet.
	RpMtiMoData = 0x0
	RpMtiMoAck  = 0x2
	RpMtiMoErr  = 0x4
	RpMtiMoSmma = 0x6
	RpMtiMtData = 0x1
	RpMtiMtAck  = 0x3
	RpMtiMtErr  = 0x5

	RpUdeIei   = 0x41
	RpCauseIei = 0x42
)

const (
	// RP Cause types (TS24.011 Table 8.4)
	RpCauseUnassigned             = 0x1
	RpCauseOpBarred               = 0x8
	RpCauseCallBarred             = 0xa
	RpCauseReserved               = 0xb
	RpCauseSmTransferRejected     = 0x15
	RpCauseMemExceeded            = 0x16
	RpCauseDestOutOfOrder         = 0x1b
	RpCauseUnidentifiedSub        = 0x1c
	RpCauseFacilityRejected       = 0x1d
	RpCauseUnknownSub             = 0x1e
	RpCauseNetOutOfOrder          = 0x26
	RpCauseTempFailure            = 0x29
	RpCauseCongestion             = 0x2a
	RpCauseResourceUnavailable    = 0x2f
	RpCauseRequestedFacNotSub     = 0x32
	RpCauseRequestedFacNotImpl    = 0x45
	RpCauseInvalidSmTransRef      = 0x51
	RpCauseSemIncorrectMessage    = 0x5f
	RpCauseInvalidMandantoryInfo  = 0x60
	RpCauseMsgTypeNotImpl         = 0x61
	RpCauseMsgTypeNotCompatible   = 0x62
	RpCauseInfoElementNonexistant = 0x63
	RpCauseProtocolError          = 0x6f
	RpCauseInterworking           = 0x7f
)

var RpCauseStr = map[byte]string{
	RpCauseUnassigned:             "Unassigned (unallocated) number",
	RpCauseOpBarred:               "Operator determined barring",
	RpCauseCallBarred:             "Call barred",
	RpCauseReserved:               "Reserved",
	RpCauseSmTransferRejected:     "Short message transfer rejected",
	RpCauseMemExceeded:            "Memory capacity exceeded",
	RpCauseDestOutOfOrder:         "Destination out of order",
	RpCauseUnidentifiedSub:        "Unidentified subscriber",
	RpCauseFacilityRejected:       "Facility rejected",
	RpCauseUnknownSub:             "Unknown subscriber",
	RpCauseNetOutOfOrder:          "Network out of order",
	RpCauseTempFailure:            "Temporary failure",
	RpCauseCongestion:             "Congestion",
	RpCauseResourceUnavailable:    "Resources unavailable, unspecified",
	RpCauseRequestedFacNotSub:     "Requested facility not subscribed",
	RpCauseRequestedFacNotImpl:    "Requested facility not implemented",
	RpCauseInvalidSmTransRef:      "Invalid short message transfer reference value",
	RpCauseSemIncorrectMessage:    "Semantically incorrect message",
	RpCauseInvalidMandantoryInfo:  "Invalid mandantory information",
	RpCauseMsgTypeNotImpl:         "Message type not non-existent or not implemented",
	RpCauseMsgTypeNotCompatible:   "Message not compatible with short message protocol state",
	RpCauseInfoElementNonexistant: "Information element non-existent or not implemented",
	RpCauseProtocolError:          "Protocol error, unspecified",
	RpCauseInterworking:           "Interworking, unspecified",
}
