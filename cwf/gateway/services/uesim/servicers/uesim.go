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
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"fbc/lib/go/radius"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	networkIDPlaceholder = "magma"
	blobTypePlaceholder  = "uesim"
	trafficMSS           = "1300"
	trafficSrvIP         = "192.168.129.42"
	trafficSrvSSHport    = "22"
	numRetries           = 10
	retryDelay           = 1000 * time.Millisecond
)

// UESimServer tracks all the UEs being simulated.
type UESimServer struct {
	store blobstore.BlobStorageFactory
	cfg   *UESimConfig
}

type UESimConfig struct {
	op                []byte
	amf               []byte
	radiusAuthAddress string
	radiusAcctAddress string
	radiusSecret      string
	brMac             string
	bypassHssAuth     bool
}

// NewUESimServer initializes a UESimServer with an empty store map.
// Output: a new UESimServer
func NewUESimServer(factory blobstore.BlobStorageFactory) (*UESimServer, error) {
	config, err := GetUESimConfig()
	if err != nil {
		return nil, err
	}
	return &UESimServer{
		store: factory,
		cfg:   config,
	}, nil
}

// AddUE tries to add this UE to the server.
// Input: The UE data which will be added.
func (srv *UESimServer) AddUE(ctx context.Context, ue *cwfprotos.UEConfig) (ret *protos.Void, err error) {
	ret = &protos.Void{}

	err = validateUEData(ue)
	if err != nil {
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	addUeToStore(srv.store, ue)
	return
}

// Authenticate triggers an authentication for the UE with the specified IMSI.
// Input: The IMSI of the UE to try to authenticate.
// Output: The resulting Radius packet returned by the Radius server.
func (srv *UESimServer) Authenticate(ctx context.Context, id *cwfprotos.AuthenticateRequest) (*cwfprotos.AuthenticateResponse, error) {
	eapIDResp, err := srv.CreateEAPIdentityRequest(id.GetImsi(), id.GetCalledStationID())
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaIDReq, err := radius.Exchange(context.Background(), eapIDResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaIDResp, err := srv.HandleRadius(id.GetImsi(), id.GetCalledStationID(), akaIDReq)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaChalReq, err := radius.Exchange(context.Background(), akaIDResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	akaChalResp, err := srv.HandleRadius(id.GetImsi(), id.GetCalledStationID(), akaChalReq)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	result, err := radius.Exchange(context.Background(), akaChalResp, srv.cfg.radiusAuthAddress)
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, err
	}

	resultBytes, err := result.Encode()
	if err != nil {
		return &cwfprotos.AuthenticateResponse{}, errors.Wrap(err, "Error encoding Radius packet")
	}
	radiusPacket := &cwfprotos.AuthenticateResponse{RadiusPacket: resultBytes}

	return radiusPacket, nil
}

func (srv *UESimServer) Disconnect(ctx context.Context, id *cwfprotos.DisconnectRequest) (*cwfprotos.DisconnectResponse, error) {
	radiusP, err := srv.MakeAccountingStopRequest(id.GetCalledStationID())
	if err != nil {
		return nil, errors.Wrap(err, "Error making Accounting Stop Radius message")
	}
	response, err := radius.Exchange(context.Background(), radiusP, srv.cfg.radiusAcctAddress)
	if err != nil {
		return nil, errors.Wrap(err, "Error exchanging Radius message")
	}
	encoded, err := response.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "Error encoding Radius packet")
	}
	return &cwfprotos.DisconnectResponse{RadiusPacket: encoded}, nil
}

// GenTraffic generates traffic using a remote iperf server. The command to be sent is configured using GenTrafficRequest
// Note that GenTrafficRequest have parameter that configures iperf client itself, and parameters that configure UESim
// Configuration parameters related to the UESim client itself (not iperf) are:
// - timeout: if different than 0 stops iperf externally after n seconds. Use it to avoid the test to hang on a unreachable server
// 	 If the test timesout it will be counted as an error. By default this is 0 (DISABLED)
// - disableServerReachabilityCheck: enables/disables the function to request the server to send the UE small packets to check if
//   the server is alive. By default this is ENABLED
func (srv *UESimServer) GenTraffic(ctx context.Context, req *cwfprotos.GenTrafficRequest) (*cwfprotos.GenTrafficResponse, error) {
	if req == nil {
		return &cwfprotos.GenTrafficResponse{}, fmt.Errorf("Nil GenTrafficRequest provided")
	}

	restartIperfServer(trafficSrvIP, trafficSrvSSHport)

	argList := []string{"--json", "-c", trafficSrvIP, "-M", trafficMSS}
	if req.Volume != nil {
		argList = append(argList, []string{"-n", req.Volume.Value}...)
	}

	if req.ReverseMode {
		argList = append(argList, "-R")
	}

	if req.Bitrate != nil {
		argList = append(argList, []string{"-b", req.Bitrate.Value}...)
	}

	if req.TimeInSecs != 0 {
		argList = append(argList, []string{"-t", strconv.FormatUint(req.TimeInSecs, 10)}...)
	}

	if req.ReportingIntervalInSecs != 0 {
		argList = append(argList, []string{"-i", strconv.FormatUint(req.ReportingIntervalInSecs, 10)}...)
	}
	output, err := executeIperfWithOptions(argList, req)
	return &cwfprotos.GenTrafficResponse{Output: output}, err
}

// Converts a blob back into a UE config
func blobToUE(blob blobstore.Blob) (*cwfprotos.UEConfig, error) {
	ue := &cwfprotos.UEConfig{}
	err := protos.Unmarshal(blob.Value, ue)
	if err != nil {
		return nil, err
	}
	return ue, nil
}

// getUE gets the UE with the specified IMSI from the blobstore.
func getUE(blobStoreFactory blobstore.BlobStorageFactory, imsi string) (ue *cwfprotos.UEConfig, err error) {
	store, err := blobStoreFactory.StartTransaction(nil)
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				glog.Errorf("Error while rolling back transaction: %s", err)
			}
		}
	}()

	blob, err := store.Get(networkIDPlaceholder, storage.TypeAndKey{Type: blobTypePlaceholder, Key: imsi})
	if err != nil {
		err = errors.Wrap(err, "Error getting UE with specified IMSI")
		return
	}
	ue, err = blobToUE(blob)
	return
}

// ConvertStorageErrorToGrpcStatus converts a UE error into a gRPC status error.
func ConvertStorageErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Unknown, err.Error())
}

// executeIperfWithOptions runs iperf with the timeout and server reachability options per req
func executeIperfWithOptions(argList []string, req *cwfprotos.GenTrafficRequest) ([]byte, error) {
	// server reachability option (Enabled by default)

	if req.DisableServerReachabilityCheck == false {
		// Check if server is reachable by requesting the server to send UE 10b of data
		reachable, err := checkIperfServerReachabilityWithRetries()
		if !reachable {
			return nil, fmt.Errorf("(%s) iperf server not reachable or didn't send traffic back to the UE."+
				"This may happen when traffic is requested before rules had time to be synched, %+v", trafficSrvIP, err)
		}
	}

	// timeout option
	if req.Timeout > 0 {
		return executeIperfWithTimeout(argList, req.Timeout)
	}
	return executeIperf(argList)
}

func checkIperfServerReachabilityWithRetries() (bool, error) {
	var (
		res bool
		err error
	)

	for i := 0; i < numRetries; i++ {
		res, err = checkIperfServerReachability()
		if res == true {
			break
		}
		glog.V(2).Infof("Iperf server was not reachable, trying one more time (%d out of %d)", i+1, numRetries)
		time.Sleep(retryDelay)
	}
	return res, err
}

// checkIperfServerReachability will request the server to send the UE a very small amount of data to check
// if the server is able to reach UE. This is useful to detect situations were we are able to send
// traffic from UE->server but not traffic UE<-server
func checkIperfServerReachability() (bool, error) {
	// iperf during 1s, reverse sending 10 bytes
	argList := []string{"1s", "iperf3", "--json", "-c", trafficSrvIP, "-R", "-n", "10", "-l", "2"}

	// run timeout command but ignore error since timeout always produce an error
	cmd := exec.Command("timeout", argList...)
	cmd.Dir = "/usr/bin"
	output, _ := cmd.Output()

	totalBytes, err := ExtractBytesReceived(output)
	if err != nil {
		return false, fmt.Errorf("Couldnt parse response from server reachability: %s", err)
	}
	if totalBytes == 0 {
		return false, nil
	}
	return true, nil
}

// executeIperfWithTimeout runs iperf with a maximum timeout. If timeout is reached, iperf will return
// error and any traffic it has logged
func executeIperfWithTimeout(argList []string, timeout uint32) ([]byte, error) {
	timeoutString := fmt.Sprintf("%ds", timeout)
	argsList2 := []string{timeoutString, "iperf3"}
	argsList2 = append(argsList2, argList...)
	return executeCommandWithRetries("timeout", argsList2)
}

func executeIperf(argList []string) ([]byte, error) {
	return executeCommandWithRetries("iperf3", argList)
}

// executeCommandWithRetries will retry a command if the error of that command contains
// a specific content (so far it will only retry in case of error `unable to receive control`
func executeCommandWithRetries(command string, argList []string) ([]byte, error) {
	var err error
	var res []byte

	for i := 0; i < numRetries; i++ {
		res, err = executeCommand(command, argList)
		if !isIperfErrorDueToControlMessage(err) {
			break
		}
		err_msg := "Retried IPERF command due to unable to receive control message"
		glog.Warning(err_msg)
		err = errors.Wrap(err, err_msg)
		time.Sleep(300 * time.Millisecond)
	}
	return res, err
}

func executeCommand(command string, argList []string) ([]byte, error) {
	glog.V(2).Info("Execute: ", command, argList)
	cmd := exec.Command(command, argList...)
	cmd.Dir = "/usr/bin"
	output, err := cmd.Output()
	if err != nil {
		newError := errors.Wrap(err, fmt.Sprintf(
			"error while executing \"%s %s\"\n output:\n%v",
			command, strings.Join(argList, " "), string(output)))
		glog.Error(newError)
		return output, newError
	}
	glog.V(5).Infof("Result:\n %s", PrettyPrintIperfResponse(output))
	return output, nil
}

func isIperfErrorDueToControlMessage(iperf_err error) bool {
	if iperf_err == nil {
		return false
	}
	return strings.Contains(iperf_err.Error(), "unable to receive control message")
}

// TODO: create a new file and structs to to parse and dump iperf message
// extractBytesReceived returns the amount of bytes sent by the Server to the UE
func ExtractBytesReceived(output []byte) (float64, error) {
	outputU, err := unmarshallIper3Response(output)
	if err != nil {
		return 0, err
	}
	endSection, found := outputU["end"].(map[string]interface{})
	if !found {
		return 0, fmt.Errorf("Couldn't parse iperf result 'end' section\n%s", PrettyPrintIperfResponse(output))
	}
	sumReceived_section, found := endSection["sum_received"].(map[string]interface{})
	if !found {
		return 0, fmt.Errorf("Couldn't parse iperf result 'sum_received' section\n%s", PrettyPrintIperfResponse(output))
	}
	bytes_param, found := sumReceived_section["bytes"].(float64)
	if !found {
		return 0, fmt.Errorf("Couldn't parse iperf result 'bytes'(float64)\n%s", PrettyPrintIperfResponse(output))
	}
	return bytes_param, nil
}

func ExtractIperfError(output []byte) (string, error) {
	outputU, err := unmarshallIper3Response(output)
	if err != nil {
		return "", err
	}
	iperfErrorMessage, found := outputU["error"].(string)
	if !found {
		return "", fmt.Errorf("Couldn't parse iperf result 'end' section\n%s", PrettyPrintIperfResponse(output))
	}
	return iperfErrorMessage, nil
}

func unmarshallIper3Response(output []byte) (map[string]interface{}, error) {
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(output, &jsonResponse)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

func PrettyPrintIperfResponse(input []byte) string {
	prettyOutput := &bytes.Buffer{}
	err := json.Indent(prettyOutput, input, "", "  ")
	if err != nil {
		return "Couldn't parse iperf3 response into JSON"
	}
	return prettyOutput.String()
}
