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

package service

import (
	"testing"

	"magma/orc8r/cloud/go/service/middleware/unary"
	platform_service "magma/orc8r/lib/go/service"

	"google.golang.org/grpc"
)

// NewTestService returns a new gRPC orchestrator service without
// loading Orchestrator plugins from disk. This should only be used in test
// contexts, where plugins are registered manually.
func NewTestService(t *testing.T, moduleName string, serviceType string) (*platform_service.Service, error) {
	if t == nil {
		panic("for tests only")
	}
	return platform_service.NewServiceWithOptions(moduleName, serviceType, grpc.UnaryInterceptor(unary.MiddlewareHandler))
}

// NewTestOrchestratorService returns a new gRPC orchestrator service without
// loading Orchestrator plugins from disk. This should only be used in test
// contexts, where plugins are registered manually.
func NewTestOrchestratorService(t *testing.T, moduleName string, serviceType string) (*OrchestratorService, error) {
	if t == nil {
		panic("for tests only")
	}
	platformService, err := platform_service.NewServiceWithOptions(moduleName, serviceType, grpc.UnaryInterceptor(unary.MiddlewareHandler))
	if err != nil {
		return nil, err
	}
	echoSrv, err := getEchoServerForOrchestratorService(serviceType)
	if err != nil {
		return nil, err
	}
	return &OrchestratorService{
		Service:    platformService,
		EchoServer: echoSrv,
	}, nil
}
