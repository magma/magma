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

// Package test provides common definitions and function for eap related tests
package mock_sessiond

import (
	"context"
	"fmt"
	"testing"

	"magma/feg/gateway/registry"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/test_utils"
	orc8rProtos "magma/orc8r/lib/go/protos"
)

// MockSessionManager test sessiond  implementation
type MockSessionManager struct {
	returnErrors bool
}

func NewRunningSessionManager(t *testing.T) *MockSessionManager {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SESSION_MANAGER)
	service := &MockSessionManager{
		returnErrors: false,
	}
	protos.RegisterLocalSessionManagerServer(srv.GrpcServer, service)
	go srv.RunTest(lis)
	return service
}

func (c *MockSessionManager) ReportRuleStats(ctx context.Context, in *protos.RuleRecordTable) (*orc8rProtos.Void, error) {
	out := new(orc8rProtos.Void)
	if c.returnErrors {
		return out, fmt.Errorf("CreateSession returnErrors enabled")
	}
	err := fmt.Errorf("ReportRuleStats not implemented on test sessionManager")
	return out, err
}

func (c *MockSessionManager) CreateSession(ctx context.Context, in *protos.LocalCreateSessionRequest) (*protos.LocalCreateSessionResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("CreateSession returnErrors enabled")
	}

	out := &protos.LocalCreateSessionResponse{
		SessionId: fmt.Sprintf("%s-12345678", in.CommonContext.Sid.Id),
	}
	return out, nil
}

func (c *MockSessionManager) EndSession(ctx context.Context, in *protos.LocalEndSessionRequest) (*protos.LocalEndSessionResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("CreateSession returnErrors enabled")
	}
	return &protos.LocalEndSessionResponse{}, nil
}

func (c *MockSessionManager) BindPolicy2Bearer(ctx context.Context, in *protos.PolicyBearerBindingRequest) (*protos.PolicyBearerBindingResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("BindPolicy2Bearer returnErrors enabled")
	}
	return &protos.PolicyBearerBindingResponse{}, nil
}

func (c *MockSessionManager) SetSessionRules(ctx context.Context, in *protos.SessionRules) (*orc8rProtos.Void, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("SetSessionRules returnErrors enabled")
	}
	return &orc8rProtos.Void{}, nil
}

func (c *MockSessionManager) UpdateTunnelIds(ctx context.Context, in *protos.UpdateTunnelIdsRequest) (*protos.UpdateTunnelIdsResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("UpdateTunnelIds returnErrors enabled")
	}
	return &protos.UpdateTunnelIdsResponse{}, nil
}

func (c *MockSessionManager) ReturnErrors(enable bool) {
	c.returnErrors = enable
}
