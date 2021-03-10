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

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"magma/feg/cloud/go/protos"
	_ "magma/feg/gateway/registry"
	"magma/feg/gateway/services/csfb"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/test_utils"
)

type marshalFunc func() (decode.SGsMessageType, *any.Any, error)

var marshallerMap = map[string]marshalFunc{
	"AR":     marshalAlertRequest,
	"DUD":    marshalDownlinkUnitData,
	"EPSDA":  marshalEPSDetachAck,
	"IMSIDA": marshalIMSIDetachAck,
	"LUA":    marshalLocationUpdateAccept,
	"LUR":    marshalLocationUpdateReject,
	"MMIR":   marshalMMInformationRequest,
	"RR":     marshalReleaseRequest,
	"SAR":    marshalServiceAbortRequest,
	"VLRRA":  marshalVLRResetAck,
	"VLRRI":  marshalVLRResetIndication,
	"VLRS":   marshalVLRStatus,
}

func main() {
	// setting up flags of the CLI
	helpPtr := flag.Bool("help", false, "[optional] Display this help message")
	cmdPtr := flag.String("rpcCall", "", "[required] The RPC call on the service. "+
		"{AR|DUD|EPSDA|IMSIDA|LUA|LUR|MMIR|RR|SAR|VLRRA|VLRRI|VLRS}")

	// setting up helper message of the CLI
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("	gw_csfb_service_cli [-h] " +
			"-rpcCall={AR|DUD|EPSDA|IMSIDA|LUA|LUR|MMIR|RR|SAR|VLRRA|VLRRI|VLRS}" +
			" <IMSI if required>")
		fmt.Println("Flags: ")
		fmt.Printf("	%s: %s\n", "rpcCall", flag.Lookup("rpcCall").Usage)
		fmt.Printf("	%s: %s\n", "help   ", flag.Lookup("help").Usage)
	}

	// parse command line inputs
	flag.Parse()

	// print usage if requested or if required arguments are not provided
	if *helpPtr || *cmdPtr == "" {
		flag.Usage()
		os.Exit(0)
	}

	// handle commands, make corresponding rpc calls, and print results
	err := handleCommands(*cmdPtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func handleCommands(cmd string) error {

	if decoderFunc, ok := marshallerMap[cmd]; ok {
		if len(flag.Args()) == 0 &&
			cmd != "VLRRA" && cmd != "VLRRI" && cmd != "VLRS" {
			return fmt.Errorf("please add IMSI as the argument")
		} else if len(flag.Args()) != 0 &&
			(cmd == "VLRRA" || cmd == "VLRRI" || cmd == "VLRS") {
			return fmt.Errorf("Reset Ack, Reset Indication, and Status don't need argument")
		}

		msgType, marshalledMsg, err := decoderFunc()
		if err != nil {
			return fmt.Errorf("error marshaling SGs message to Any: %s", err)
		}
		return sendMessage(msgType, marshalledMsg)
	}
	flag.Usage()
	return fmt.Errorf("command %s is not supported", cmd)
}

func sendMessage(messageType decode.SGsMessageType, msg *any.Any) error {
	_, err := csfb.SendSGsMessageToGateway(messageType, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	fmt.Println("Successfully sent message")
	return nil
}

func marshalAlertRequest() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.AlertRequest{
		Imsi: flag.Arg(0),
	})
	return decode.SGsAPAlertRequest, marshalledMsg, err
}

func marshalDownlinkUnitData() (decode.SGsMessageType, *any.Any, error) {
	nasMessageContainer := test_utils.ConstructDefaultIE(
		decode.IEINASMessageContainer,
		5,
	)
	marshalledMsg, err := ptypes.MarshalAny(&protos.DownlinkUnitdata{
		Imsi:                flag.Arg(0),
		NasMessageContainer: nasMessageContainer[2:],
	})
	return decode.SGsAPDownlinkUnitdata, marshalledMsg, err
}

func marshalEPSDetachAck() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.EPSDetachAck{
		Imsi: flag.Arg(0),
	})
	return decode.SGsAPEPSDetachAck, marshalledMsg, err
}

func marshalIMSIDetachAck() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: flag.Arg(0),
	})
	return decode.SGsAPIMSIDetachAck, marshalledMsg, err
}

func marshalLocationUpdateAccept() (decode.SGsMessageType, *any.Any, error) {
	LAI := test_utils.ConstructDefaultLocationAreaIdentifier()
	marshalledMsg, err := ptypes.MarshalAny(&protos.LocationUpdateAccept{
		Imsi:                   flag.Arg(0),
		LocationAreaIdentifier: LAI[2:],
	})
	return decode.SGsAPLocationUpdateAccept, marshalledMsg, err
}

func marshalLocationUpdateReject() (decode.SGsMessageType, *any.Any, error) {
	var rejectCause []byte
	if len(flag.Args()) == 2 {
		rejectCauseCode, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			return decode.SGsAPLocationUpdateReject, nil, err
		}
		rejectCause = []byte{byte(rejectCauseCode)}
	} else {
		rejectCause = []byte{byte(0x11)}
	}
	marshalledMsg, err := ptypes.MarshalAny(&protos.LocationUpdateReject{
		Imsi:        flag.Arg(0),
		RejectCause: rejectCause,
	})
	return decode.SGsAPLocationUpdateReject, marshalledMsg, err
}

func marshalMMInformationRequest() (decode.SGsMessageType, *any.Any, error) {
	mmInfo := test_utils.ConstructDefaultMMInformation()
	marshalledMsg, err := ptypes.MarshalAny(&protos.MMInformationRequest{
		Imsi:          flag.Arg(0),
		MmInformation: mmInfo[2:],
	})
	return decode.SGsAPMMInformationRequest, marshalledMsg, err
}

func marshalReleaseRequest() (decode.SGsMessageType, *any.Any, error) {
	sgsCause := test_utils.ConstructDefaultIE(
		decode.IEISGsCause,
		1,
	)
	marshalledMsg, err := ptypes.MarshalAny(&protos.ReleaseRequest{
		Imsi:     flag.Arg(0),
		SgsCause: sgsCause[2:],
	})
	return decode.SGsAPReleaseRequest, marshalledMsg, err
}

func marshalServiceAbortRequest() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.ServiceAbortRequest{
		Imsi: flag.Arg(0),
	})
	return decode.SGsAPServiceAbortRequest, marshalledMsg, err
}

func marshalVLRResetAck() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.ResetAck{
		VlrName: "www.facebook.com",
	})
	return decode.SGsAPResetAck, marshalledMsg, err
}

func marshalVLRResetIndication() (decode.SGsMessageType, *any.Any, error) {
	marshalledMsg, err := ptypes.MarshalAny(&protos.ResetIndication{
		VlrName: "www.facebook.com",
	})
	return decode.SGsAPResetIndication, marshalledMsg, err
}

func marshalVLRStatus() (decode.SGsMessageType, *any.Any, error) {
	sgsCause := test_utils.ConstructDefaultIE(
		decode.IEISGsCause,
		1,
	)
	erroneousMsg := test_utils.ConstructDefaultIE(
		decode.IEIErroneousMessage,
		10,
	)
	marshalledMsg, err := ptypes.MarshalAny(&protos.Status{
		Imsi:             flag.Arg(0),
		SgsCause:         sgsCause[2:],
		ErroneousMessage: erroneousMsg[2:],
	})
	return decode.SGsAPStatus, marshalledMsg, err
}
