package sms_ll

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test cases:
// 1) Marshal a specific SMS-Deliver and ensure the binary matches what we expect
// 2) Marshal a long SMS and ensure binary matches what we expect
// 3) Given a delivery report, decode the TPDU

// Pick a consistent timestamp for tests

func TestEncodeSingleSMS(t *testing.T) {
	msg := "Here's a test."
	ts := time.Date(2020, 9, 14, 16, 30, 50, 12345, time.UTC)
	num := "18658675309"
	ref := []byte{7}
	expected := "790127010702b9110020240b918156685703f90000029041610305000ec8b2bc7c9a83c2207a794e7701"
	b, e := GenerateSmsDelivers(msg, num, ts, ref)
	if e != nil {
		t.Errorf("Error: %s", e)
	}

	if GetMessageCount(msg) != 1 || len(b) != 1 {
		t.Errorf("Incorrect number of PDUs generated")
	}
	if hex.EncodeToString(b[0]) != expected {
		t.Errorf("Incorrect PDU generated/wanted:\n%s\n%s", hex.EncodeToString(b[0]), expected)
	}

}

func TestEncodingMultipleSMS(t *testing.T) {
	msg := "Here's a test of a veeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeerrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrryyyyyyyyyyyyyyyyyy long message that's super super long."
	ts := time.Date(2020, 9, 14, 16, 30, 50, 12345, time.UTC)
	num := "18658675309"
	ref := []byte{1, 2}
	expected := []string{"1901a6010102b911009f640b918156685703f9000002904161030500a0050003010201906579f934078541f4f29c0e7a9b416190bd5c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c96cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2e97cbe572b95c2ecfe7f3f97c3e9fcfe7f3f97c3e9fcfe741ecb7fb0c6a97e7f3f0b90ca2a3c3", "290133010202b911002c640b918156685703f90000029041610305001c050003010202e8a739685e8797e5a0791d5e9683d86ff7d905"}

	b, e := GenerateSmsDelivers(msg, num, ts, ref)
	if e != nil {
		t.Errorf("Error: %s", e)
	}
	if GetMessageCount(msg) != 2 || len(b) != 2 {
		t.Errorf("Wrong number of PDUs generated")
	}
	for i := range b {
		if hex.EncodeToString(b[i]) != expected[i] {
			t.Errorf("Incorrect PDU generated/wanted:\n%s\n%s", hex.EncodeToString(b[i]), expected[i])
		}
	}
}

func TestDecode(t *testing.T) {
	type decodeResult struct {
		ref   uint8
		res   bool
		cause string
		err   bool // true if we get an error
	}

	tests := []struct {
		name  string
		input string
		want  decodeResult
	}{
		{name: "deliver-success", input: "d90106020141020000", want: decodeResult{1, true, "", false}},
		{name: "deliver-fail", input: "d9010404040160", want: decodeResult{4, false, "Invalid mandantory information", false}},
		{name: "giberish", input: "209491c192c912ca1010", want: decodeResult{0, false, "", true}},
		{name: "data", input: "790127010702b9110020240b918156685703f90000029041610305000ec8b2bc7c9a83c2207a794e7701", want: decodeResult{0, false, "", true}},
	}

	for _, tc := range tests {
		msg, _ := hex.DecodeString(tc.input)
		actual, err := Decode(msg)
		result := decodeResult{actual.Reference, actual.IsSuccessful, actual.ErrorMessage, err != nil}
		assert.Equal(t, tc.want, result)
	}
}

func TestPiecewiseDecodeDeliveryFailure(t *testing.T) {
	input := "d9010404010160"
	cp_hex, _ := hex.DecodeString(input)
	cpm := new(cpMessage)
	err := cpm.unmarshalBinary(cp_hex)
	if err != nil {
		t.Errorf("Failed to decode valid CP-DATA")
	}
	if cpm.messageType != CpData {
		t.Errorf("Failed to decode valid CP-DATA")
	}
	if int(cpm.length) != 4 || len(cpm.rpdu) != 4 {
		t.Errorf("CP-DATA length incorrect")
	}
	if cpm.cause != 0x0 {
		t.Errorf("CP-DATA has cause set to non-zero")
	}

	rpm := new(rpMessage)
	err = rpm.unmarshalBinary(cpm.rpdu)
	if err != nil {
		t.Errorf("Failed to decode valid RP-ERROR")
	}
	msg_type, err := rpm.msgType()
	if err != nil {
		t.Errorf("Failed to decode valid RP-ERROR")
	}
	if msg_type != RpError || rpm.cause.cause != RpCauseInvalidMandantoryInfo {
		t.Errorf("Failed to decode valid RP-ERROR")
	}
}

func TestPiecewiseDecodeDeliveryReport(t *testing.T) {
	input := "d90106020141020000"
	cp_hex, _ := hex.DecodeString(input)
	cpm := new(cpMessage)
	err := cpm.unmarshalBinary(cp_hex)
	if err != nil {
		t.Errorf("Failed to decode valid CP-DATA")
	}
	if cpm.messageType != CpData {
		t.Errorf("Failed to decode valid CP-DATA")
	}
	if int(cpm.length) != 6 && len(cpm.rpdu) != 6 {
		t.Errorf("CP-DATA length incorrect")
	}
	if cpm.cause != 0x0 {
		t.Errorf("CP-DATA has cause set to non-zero")
	}

	rpm := new(rpMessage)
	err = rpm.unmarshalBinary(cpm.rpdu)
	if err != nil {
		t.Errorf("Failed to decode valid RP-ACK")
	}
	msg_type, err := rpm.msgType()
	if err != nil {
		t.Errorf("Failed to decode valid RP-ACK")
	}
	if msg_type != RpAck {
		t.Errorf("Failed to decode valid RP-ACK")
	}
	if rpm.userData.iei != RpUdeIei {
		t.Errorf("Failed to decode valid RP-ACK User Data IEI")
	}
	if rpm.userData.length != byte(2) {
		t.Errorf("Failed to decode valid RP-ACK User Data Length")
	}
	if len(rpm.userData.tpdu) != int(rpm.userData.length) {
		t.Errorf("RP-ACK user data length doesn't match payload")
	}
}

func TestUnmarshalAddressElement(t *testing.T) {
	input := "0b911605935713f2"
	rpadde_hex, _ := hex.DecodeString(input)
	rpadde := new(rpAddressElement)
	l, err := rpadde.unmarshalBinary(rpadde_hex)
	if err != nil {
		t.Errorf("Failed to decode RP Address Element")
	}

	if l != 8 || rpadde.length != 0x0b {
		t.Errorf("RPAddressElement incorrect length")
	}
	if rpadde.numberInfo != 0x91 {
		t.Errorf("RPAddressElement incorrect number info")
	}
	num, _ := hex.DecodeString("1605935713f2")
	if !reflect.DeepEqual(rpadde.number, num) {
		t.Errorf("RPAddressElement incorrect number. Have:\n%s\nwant\n%s", hex.Dump(rpadde.number), hex.Dump(num))
	}
}
