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

package main

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/golang/glog"

	"magma/feg/gateway/registry"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/service"
)

var addFlowLat = flag.Duration("add_lat", 0, "add, update & delete flow delay")
var activiteFlowLat = flag.Duration("activate_lat", 0, "Activate flows delay")

// Dummy Pipelined
type DummyPipelined struct {
	*protos.UnimplementedPipelinedServer
}

func NewPipelined() *DummyPipelined {
	return &DummyPipelined{}
}

func (c *DummyPipelined) AddUEMacFlow(
	ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {

	return defaultHandler(req)
}

func (c *DummyPipelined) UpdateIPFIXFlow(
	ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {

	return defaultHandler(req)
}

func (c *DummyPipelined) DeleteUEMacFlow(
	ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {

	return defaultHandler(req)
}

func (c *DummyPipelined) ActivateFlows(
	ctx context.Context, req *protos.ActivateFlowsRequest) (res *protos.ActivateFlowsResult, err error) {

	res = &protos.ActivateFlowsResult{
		PolicyResults: []*protos.RuleModResult{},
	}
	for _, policyRule := range req.GetPolicies() {
		res.PolicyResults = append(res.PolicyResults, &protos.RuleModResult{
			RuleId: policyRule.GetRule().GetId(),
			Result: protos.RuleModResult_SUCCESS,
		})
	}
	if activiteFlowLat != nil && *activiteFlowLat != 0 {
		time.Sleep(*activiteFlowLat)
	}
	return
}

// defaultHandler executes a default action that always replies with a goood answer
// unless there is an issue with the mac address.
func defaultHandler(req *protos.UEMacFlowRequest) (resp *protos.FlowResponse, err error) {
	_, err = net.ParseMAC(req.GetMacAddr())
	if err != nil {
		resp = &protos.FlowResponse{Result: protos.FlowResponse_FAILURE}
		return
	}
	resp = &protos.FlowResponse{Result: protos.FlowResponse_SUCCESS}
	if addFlowLat != nil && *addFlowLat != 0 {
		time.Sleep(*addFlowLat)
	}
	return
}

func main() {
	flag.Parse() // for glog

	// Dummy Pipelined service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.PIPELINED)
	if err != nil {
		glog.Fatalf("Error creating %s service: %v", registry.PIPELINED, err)
	}
	protos.RegisterPipelinedServer(srv.GrpcServer, NewPipelined())
	glog.Info("Starting Dummy Pipelined Service")
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running Dummy Pipelined: %v", err)
	}
}
