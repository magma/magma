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
	"errors"
	"fmt"
	"time"

	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/tpdu"
)

// SMSSerde (SER-ializer/DE-serializer) is an interface that wraps the encode/
// decode functions in this package to facilitate unit testing of components
// that depend on this package's functionality.
type SMSSerde interface {
	EncodeMessage(message string, fromNum string, timestamp time.Time, references []uint8) ([][]byte, error)
	DecodeDelivery(input []byte) (SMSDeliveryReport, error)
}

// DefaultSMSSerde is the SMSSerde impl that's backed by the exported functions
// in this package.
type DefaultSMSSerde struct{}

func (d *DefaultSMSSerde) EncodeMessage(message string, fromNum string, timestamp time.Time, references []uint8) ([][]byte, error) {
	return GenerateSmsDelivers(message, fromNum, timestamp, references)
}

func (d *DefaultSMSSerde) DecodeDelivery(input []byte) (SMSDeliveryReport, error) {
	return Decode(input)
}

// Generate fully encoded SMS PDUs for delivery to a UE (MS). Will handle
// encoding and chunking of messages as appropriate. We first generate TPDUs,
// then RP-DATA headers, and finally CP-DATA headers, resulting in a set of
// byte arrays that can be directly delivered to a UE (MS).
// Inputs:
// 	message: A UTF-8 string representing the SMS.
//	from_num: A E.164 encoded source number.
//	timestamp: The sender timestamp for the SMS (generally, use current server time)
//	references: An array of references for the messages we'll generate. Must match number of PDUs generated.
// Outputs:
//	- Array of byte array representing the set of fully-encoded CP-DATA(RP-DATA(TPDU)) messages generated
//	- Error	(if any)
func GenerateSmsDelivers(message string, fromNum string, timestamp time.Time, references []uint8) ([][]byte, error) {
	tpdus := createTpdus(message, fromNum, timestamp)
	if len(references) != len(tpdus) {
		return nil, fmt.Errorf("insufficient references for generated TPDU (have %d, need %d)", len(references), len(tpdus))
	}

	output := make([][]byte, 0)

	for i := range tpdus {
		tp, err := tpdus[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		rpm, err := createRpDataMessage(RpMtiMtData, references[i], tp)
		if err != nil {
			return nil, err
		}

		marshaledRPM := rpm.marshalBinary()
		cp_id := references[i] & 0xf // low order bits become the cp_id
		cpm, err := createCpDataMessage(marshaledRPM, cp_id)
		if err != nil {
			return nil, err
		}

		marshaledCPM := cpm.marshalBinary()

		output = append(output, marshaledCPM)
	}

	return output, nil
}

// Return the number of SMS PDUs that will be generated for a given string,
// after taking TPDU encoding into account.
// Inputs:
// 	message: A UTF-8 string representing the SMS
// Outputs:
//	- An integer representing the number of SMS messages the input will
//	need to be split across after encoding.
func GetMessageCount(message string) int {
	// Number of SMS is determined only by the message content, not the
	// timestamp or number.
	return len(createTpdus(message, "123456", time.Now()))
}

// SMSDeliveryReport is a struct that wraps the decoded result of a
// SMS-DELIVERY-REPORT message.
// ErrorMessage field will be non-empty if IsSuccessful is false and the input
// was successfully decoded.
type SMSDeliveryReport struct {
	Reference    uint8
	IsSuccessful bool
	ErrorMessage string
}

// Decodes an SMS-DELIVERY-REPORT message.
// Inputs:
//	input: A byte array representing a fully encoded SMS delivered from a UE
// Outputs:
//	- uint8: Reference number representing the SMS
//	- bool: true if the message was successfully delivered
//	- string: descriptive delivery status (only present for failures)
//	- error: if the message received is not an SMS-DELIVERY-REPORT.
func Decode(input []byte) (SMSDeliveryReport, error) {
	ret := SMSDeliveryReport{}
	// A message is a delivery report iff we receive a CP-DATA(RP-ACK(TPDU)). We can ignore everything else.
	cpm := new(cpMessage)
	err := cpm.unmarshalBinary(input)
	if err != nil {
		return ret, err
	}

	// Valid CPM, check type
	if cpm.messageType != CpData {
		return ret, fmt.Errorf("not a CP-DATA message: %x", cpm.messageType)
	}

	rpm := new(rpMessage)
	err = rpm.unmarshalBinary(cpm.rpdu)
	if err != nil {
		return ret, err
	}

	// Valid RPM, must be of type RP-ERROR or RP-ACK
	msgType, _ := rpm.msgType()
	switch msgType {
	case RpAck: // success! get reference
		return SMSDeliveryReport{
			Reference:    rpm.reference,
			IsSuccessful: true,
			ErrorMessage: "",
		}, nil
	case RpError: // failure, get cause
		return SMSDeliveryReport{
			Reference:    rpm.reference,
			IsSuccessful: false,
			ErrorMessage: rpm.cause.causeStr,
		}, nil
	default:
		return ret, errors.New("RP-DATA message, ignoring")
	}
}

func createTpdus(message string, from_num string, timestamp time.Time) []tpdu.TPDU {
	tpdus, _ := sms.Encode([]byte(message), sms.AsDeliver, sms.From(from_num))
	for i := range tpdus {
		tpdus[i].FirstOctet |= tpdu.FoMMS // Android won't accept if this bit isn't set.
		tpdus[i].FirstOctet |= tpdu.FoSRI // Request a delivery report.
		tpdus[i].SCTS = tpdu.Timestamp{Time: timestamp}
	}

	return tpdus
}

// Helper for creating a RP-DATA message.
func createRpDataMessage(mti byte, reference byte, data []byte) (rpMessage, error) {
	msg_type, err := rpMessage{mti: mti}.msgType()
	if msg_type != RpData {
		return rpMessage{}, smsRpError(fmt.Sprintf("MTI isn't RP-DATA: %x", mti))
	} else if err != nil {
		return rpMessage{}, err
	}

	ud, err := createRpUserElement(data)
	if err != nil {
		return rpMessage{}, err
	}

	rpm := rpMessage{
		mti:       mti,
		reference: reference,
		userData:  ud,
	}

	if rpm.direction() == RpMt { // MT-SMS should have empty dest address
		rpm.originatorAddress, rpm.destinationAddress = newFakeRpAddressElement(), rpAddressElement{length: 0}
	} else { // MO-SMS should have empty orig address
		rpm.originatorAddress, rpm.destinationAddress = rpAddressElement{length: 0}, newFakeRpAddressElement()
	}

	return rpm, nil
}

func createCpDataMessage(rpdu []byte, txID byte) (cpMessage, error) {
	cpm, err := createCpMessage(txID, CpData, rpdu)
	if err != nil {
		return cpm, err
	}
	return cpm, nil
}
