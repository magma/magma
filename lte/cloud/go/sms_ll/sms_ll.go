package sms_ll

import (
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/tpdu"
	"time"
	"fmt"
)

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
func GenerateSmsDelivers(message string, from_num string, timestamp time.Time, references []uint8) ([][]byte, error) {
	tpdus := createTpdus(message, from_num, timestamp)
	if len(references) != len(tpdus) {
		return nil, fmt.Errorf("Insufficient references for generated TPDU (have %d, need %d)", len(references), len(tpdus))
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

// Decodes an SMS-DELIVERY-REPORT message.
// Inputs:
//	input: A byte array representing a fully encoded SMS delivered from a UE
// Outputs:
//	- uint8: Reference number representing the SMS
//	- bool: true if the message was successfully delivered
//	- string: descriptive delivery status (only present for failures)
//	- error: if the message received is not an SMS-DELIVERY-REPORT.
func Decode(input []byte) (uint8, bool, string, error) {
	// A message is a delivery report iff we receive a CP-DATA(RP-ACK(TPDU)). We can ignore everything else.
	cpm := new(cpMessage)
	err := cpm.unmarshalBinary(input)
	if err != nil {
		return 0, false, "", err
	}

	// Valid CPM, check type
	if cpm.messageType != CpData {
		return 0, false, "", fmt.Errorf("Not a CP-DATA message: %x", cpm.messageType)
	}

	rpm := new(rpMessage)
	err = rpm.unmarshalBinary(cpm.rpdu)
	if err != nil {
		return 0, false, "", err
	}

	// Valid RPM, must be of type RP-ERROR or RP-ACK
	msg_type, _ := rpm.msgType()
	switch msg_type {
	case RpAck: // success! get reference
		return uint8(rpm.reference), true, "", nil
	case RpError: // failure, get cause
		return uint8(rpm.reference), false, rpm.cause.cause_str, nil
	default:
		return 0, false, "", fmt.Errorf("RP-DATA message, ignoring")
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
		mti: mti,
		reference: reference,
		userData: ud,
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
