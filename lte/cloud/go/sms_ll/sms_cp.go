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
package sms_ll

import (
	"fmt"
)

func smsCpError(field string) error {
	return fmt.Errorf("smscp: %s", field)
}

// Handles creation of SMS-CM messages (3GPP TS 24.011 7.2)

// CP Message represents
type cpMessage struct {
	// Contains Transaction ID and Protocol Disc
	firstOctet byte

	// CP-DATA, CP-ACK, or CP-ERROR
	messageType byte

	// Only present for CP-ERROR
	cause byte

	// Only present for CP-DATA
	length byte

	// Only present for CP-DATA
	rpdu []byte
}

func createCpMessage(txID byte, messageType byte, rpdu []byte) (cpMessage, error) {
	if int(txID) < 0 || int(txID) > 15 {
		return cpMessage{}, smsCpError(fmt.Sprintf("Transaction ID must be 0-15: %x", txID))
	}

	switch messageType {
	case CpData, CpError, CpAck:
		// do nothing
	default:
		return cpMessage{}, smsCpError(fmt.Sprintf("Invalid CP Message Type: 0x%x", messageType))
	}

	// Set txID and protocol disc
	fo := byte(0x0)
	fo |= (txID << 4)
	fo |= CpProtocolDisc

	rpduCopy := make([]byte, len(rpdu))
	copy(rpduCopy, rpdu)

	return cpMessage{
		firstOctet:  fo,
		messageType: messageType,
		length:      byte(len(rpduCopy)),
		rpdu:        rpduCopy,
	}, nil
}

func (cpm cpMessage) marshalBinary() []byte {
	b := []byte{cpm.firstOctet, cpm.messageType}

	switch cpm.messageType {
	case CpData:
		b = append(b, cpm.length)
		b = append(b, cpm.rpdu...)
	case CpError:
		b = append(b, cpm.cause)
	case CpAck:
		// No additional data, do nothing
	}
	return b
}

func (cpm *cpMessage) unmarshalBinary(input []byte) error {
	// must be at least two octets long
	if len(input) < 2 {
		return smsCpError("Message too short")
	}

	cpm.firstOctet = input[0]
	cpm.messageType = input[1]

	switch cpm.messageType {
	case CpData:
		if len(input) < 3 {
			return smsCpError(fmt.Sprintf("message too short for message type %x", CpData))
		}
		cpm.length = input[2]
		cpm.rpdu = make([]byte, len(input[3:]))
		copy(cpm.rpdu, input[3:])
	case CpError:
		if len(input) < 3 {
			return smsCpError(fmt.Sprintf("message too short for message type %x", CpError))
		}
		if _, ok := CpCauseStr[input[2]]; ok {
			cpm.cause = input[2]
		} else {
			return smsCpError(fmt.Sprintf("Invalid cause: %x", cpm.cause))
		}
	case CpAck:
		// Do nothing -- no more data
	default:
		return smsCpError(fmt.Sprintf("Invalid IE type: %x", cpm.messageType))
	}

	return nil
}

func (cpm cpMessage) GetTransactionId() byte {
	return cpm.firstOctet >> 4
}

func (cpm cpMessage) GetProtocolDisc() byte {
	return cpm.firstOctet & 0xf
}

const (
	// CP Message bit fields

	// Protocol discriminator (3GPP TS 24.007 11.2.3.1.1)
	// For SMS-related messages, this is always 0x9 (half-octet)
	CpProtocolDisc = 0x9

	// Message types 24.011 8.1.3
	CpData  = 0x1
	CpAck   = 0x4
	CpError = 0x10

	// IE Types
	CpIeiUser  = 0x1
	CpIeiCause = 0x2
)

const (
	// CP Cause error types (24.011 8.1.4.2, Table 8.2)
	CpCauseNetworkFailure               = 0x11
	CpCauseCongestion                   = 0x16
	CpCauseInvalidTi                    = 0x51
	CpCauseSemanticallyIncorrect        = 0x5f
	CpCauseInvalidMandantoryInformation = 0x60
	CpCauseMessageTypeNonexistant       = 0x61
	CpCauseMessageNotCompatible         = 0x62
	CpCauseInfoElementNonexistant       = 0x63
	CpCauseProtocolError                = 0x6f
)

var CpCauseStr = map[byte]string{
	CpCauseNetworkFailure:               "Network failure",
	CpCauseCongestion:                   "Congestion",
	CpCauseInvalidTi:                    "Invalid Transaction Identifier value",
	CpCauseSemanticallyIncorrect:        "Semantically incorrect message",
	CpCauseInvalidMandantoryInformation: "Invalid mandantory information",
	CpCauseMessageTypeNonexistant:       "Message type non-existent or not implemented",
	CpCauseMessageNotCompatible:         "Message not compatible with the short message protocol state",
	CpCauseInfoElementNonexistant:       "Information element non-existent or not implemented",
	CpCauseProtocolError:                "Protocol error, unspecified",
}
