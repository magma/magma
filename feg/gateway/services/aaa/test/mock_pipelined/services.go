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

package mock_pipelined

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/feg/gateway/registry"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/test_utils"
)

// mock_pipelined stores all the messages received and sent on the Actions slice
// You can trace back messages or just check the last message that was received to see
// what was the reply or what was the content of the message received by pipelined

const (
	AddUEMacFlowCommand    = 1
	UpdateIPFIXFlowCommand = 2
	DeleteUEMacFlowCommand = 3
)

// records all the requests, responses and errors done to the client
type Action struct {
	Request         *protos.UEMacFlowRequest
	ExecutedCommand int
	Response        *protos.FlowResponse
	Err             error
}

// SessionManager test sessiond implementation
type MockPipelined struct {
	*protos.UnimplementedPipelinedServer
	Actions []*Action
}

func NewRunningPipelined(t *testing.T) *MockPipelined {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.PIPELINED)
	service := &MockPipelined{
		Actions: make([]*Action, 0),
	}
	protos.RegisterPipelinedServer(srv.GrpcServer, service)
	go srv.RunTest(lis)
	return service
}

func (c *MockPipelined) AddUEMacFlow(ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {
	c.appendNewAction(AddUEMacFlowCommand)
	return defaultHandler(req, c.Actions)
}

func (c *MockPipelined) UpdateIPFIXFlow(ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {
	c.appendNewAction(UpdateIPFIXFlowCommand)
	return defaultHandler(req, c.Actions)
}

func (c *MockPipelined) DeleteUEMacFlow(ctx context.Context, req *protos.UEMacFlowRequest) (*protos.FlowResponse, error) {
	c.appendNewAction(DeleteUEMacFlowCommand)
	return defaultHandler(req, c.Actions)
}

func (c *MockPipelined) GetLastAction() (*Action, error) {
	if len(c.Actions) == 0 {
		return nil, fmt.Errorf("No actions have been logged")
	}
	return c.Actions[len(c.Actions)-1], nil
}

func (c *MockPipelined) ClearActions() {
	c.Actions = make([]*Action, 0)
}

func (c *MockPipelined) appendNewAction(actionCommand int) {
	action := &Action{ExecutedCommand: actionCommand}
	c.Actions = append(c.Actions, action)
}

// defaultHandler executes a default action that always replies with a goood answer
// unless there is an issue with the mac address.
func defaultHandler(req *protos.UEMacFlowRequest, actions []*Action) (*protos.FlowResponse, error) {
	action := actions[len(actions)-1]
	action.Request = req
	_, err := net.ParseMAC(req.GetMacAddr())
	if err != nil {
		action.Response = &protos.FlowResponse{Result: protos.FlowResponse_FAILURE}
		action.Err = nil
		return action.Response, action.Err
	}
	action.Response = &protos.FlowResponse{Result: protos.FlowResponse_SUCCESS}
	action.Err = nil
	return action.Response, action.Err
}

// helper assert functions for testing
func AssertMacFlowInstall(t *testing.T, pipelined *MockPipelined) {
	action, err := pipelined.GetLastAction()
	assert.NoError(t, err)
	assert.Equal(t, AddUEMacFlowCommand, action.ExecutedCommand)
}

func AssertIPIXFlowUpdate(t *testing.T, pipelined *MockPipelined) {
	action, err := pipelined.GetLastAction()
	assert.NoError(t, err)
	assert.Equal(t, UpdateIPFIXFlowCommand, action.ExecutedCommand)
}

func AssertIDeleteMacFlow(t *testing.T, pipelined *MockPipelined) {
	action, err := pipelined.GetLastAction()
	assert.NoError(t, err)
	assert.Equal(t, DeleteUEMacFlowCommand, action.ExecutedCommand)
}

func AssertReceivedApMacAndAddress(t *testing.T,
	pipelined *MockPipelined, expectedApMac string, expectedApName string) {
	action, err := pipelined.GetLastAction()
	assert.NoError(t, err)
	receivedMac := action.Request.GetApMacAddr()
	receivedApName := action.Request.GetApName()
	assert.Equal(t, expectedApMac, receivedMac)
	assert.Equal(t, expectedApName, receivedApName)
}
