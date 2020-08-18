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

package pipelined

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/golang/glog"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	lte_protos "magma/lte/cloud/go/protos"
)

type pipelinedClient struct {
	lte_protos.PipelinedClient
}

// getPipelinedClient is a utility function to get a RPC connection to the
// pipelined client
func getPipelinedClient() (*pipelinedClient, error) {
	conn, err := registry.GetConnection(registry.PIPELINED)
	if err != nil {
		err = fmt.Errorf("Local SessionManager client initialization error: %v", err)
		glog.Error(err)
		return nil, err
	}
	return &pipelinedClient{lte_protos.NewPipelinedClient(conn)}, err
}

func AddUeMacFlow(sid *lte_protos.SubscriberID, aaaCtx *protos.Context) error {
	flowRequest := createFlowRequest(sid, aaaCtx)

	cli, err := getPipelinedClient()
	if err != nil {
		return err
	}

	response, err := cli.AddUEMacFlow(context.Background(), flowRequest)
	if err != nil {
		return err
	}
	if response.Result != lte_protos.FlowResponse_SUCCESS {
		return fmt.Errorf("Could not activate mac flow for subscriber %s", flowRequest.GetSid())
	}
	return nil
}

func DeleteUeMacFlow(sid *lte_protos.SubscriberID, aaaCtx *protos.Context) error {
	flowRequest := createFlowRequest(sid, aaaCtx)

	cli, err := getPipelinedClient()
	if err != nil {
		return err
	}

	response, err := cli.DeleteUEMacFlow(context.Background(), flowRequest)
	if err != nil {
		return err
	}
	if response.Result != lte_protos.FlowResponse_SUCCESS {
		return fmt.Errorf("Could not delte mac flow for subscriber %s", flowRequest.GetSid())
	}
	return nil
}

func createFlowRequest(sid *lte_protos.SubscriberID, aaaCtx *protos.Context) *lte_protos.UEMacFlowRequest {
	parsedApMacAddr, parsedApName := getApMacAndApName(aaaCtx.GetApn())
	return &lte_protos.UEMacFlowRequest{
		Sid:       sid,
		MacAddr:   aaaCtx.GetMacAddr(),
		Msisdn:    aaaCtx.GetMsisdn(),
		ApMacAddr: parsedApMacAddr,
		ApName:    parsedApName,
	}
}

// getApMacAndApName splits apn from AP (mac:name) into mac and name
// Example - 1C-B9-C4-36-04-F0:Wifi-Offload-hotspot20
// if apn is not valid it returns  apnName as apn and blank mac address
func getApMacAndApName(apn string) (apMAC string, apName string) {
	apMAC, apName, err := parseCompositeAPN(apn)
	if err != nil {
		apMAC = ""
		apName = apn
	}
	return
}

func parseCompositeAPN(apn string) (apMAC string, apName string, err error) {
	splitAPmacAndName := strings.Split(apn, ":")
	if len(splitAPmacAndName) != 2 {
		err = fmt.Errorf("Invalid composite APN: %v", err)
		return
	}
	_, err = net.ParseMAC(splitAPmacAndName[0])
	if err != nil {
		err = fmt.Errorf("Invalid AP MAC Address: %v", err)
		return
	}
	apMAC, apName = splitAPmacAndName[0], splitAPmacAndName[1]
	return
}
