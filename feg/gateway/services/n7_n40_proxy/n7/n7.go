/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package n7

import (
	"fmt"
	"strings"

	"magma/feg/gateway/sbi"
	sbi_NpcfSMPolicyControl "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/gateway/service_registry"
)

type N7Client struct {
	*sbi.BaseClientWithNotifier
	sbi_NpcfSMPolicyControl.ClientWithResponsesInterface
	CloudRegistry service_registry.GatewayRegistry
}

// NewN7ClientWithHandlers creates a N7 client and adds N7 handlers
func NewN7ClientWithHandlers(cfg *N7Config, cloudReg service_registry.GatewayRegistry) (*N7Client, error) {
	// client creation to handle magma initiated request
	n7Options := sbi_NpcfSMPolicyControl.WithHTTPClient(cfg.ServerConfig.BuildHttpClient())
	serverString := cfg.ServerConfig.BuildServerString()
	cliWithResponses, err := sbi_NpcfSMPolicyControl.NewClientWithResponses(serverString, n7Options)
	if err != nil {
		return nil, fmt.Errorf("error creating NewClientWithResponses: %s", err)
	}
	n7Cli := NewN7Client(cfg, cliWithResponses, cloudReg)

	// add handlers to handle PCF initiated requests
	err = n7Cli.registerHandlers()
	if err != nil {
		return nil, fmt.Errorf("error registering handlers: %s", err)
	}
	n7Cli.NotifyServer.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting notification handler: %s", err)
	}
	return n7Cli, nil
}

// NewN7Client creates a N7 api client and sets the OAuth2 client credentials for authorizing requests
func NewN7Client(cfg *N7Config, cliWithResponses sbi_NpcfSMPolicyControl.ClientWithResponsesInterface, cloudReg service_registry.GatewayRegistry,
) *N7Client {
	return &N7Client{
		BaseClientWithNotifier:       sbi.NewBaseClientWithNotifyServer(cfg.ClientConfig, cfg.ServerConfig),
		ClientWithResponsesInterface: cliWithResponses,
		CloudRegistry:                cloudReg,
	}
}

func removeIMSIPrefix(imsi string) string {
	return strings.TrimPrefix(imsi, "IMSI")
}
