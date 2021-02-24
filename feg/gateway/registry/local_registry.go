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

package registry

import (
	"log"

	"google.golang.org/grpc"

	"magma/orc8r/lib/go/registry"
	platform_registry "magma/orc8r/lib/go/registry"
)

const (
	ModuleName = "feg"

	CONTROL_PROXY    = "CONTROL_PROXY"
	S6A_PROXY        = "S6A_PROXY"
	S8_PROXY         = "S8_PROXY"
	SESSION_PROXY    = "SESSION_PROXY"
	SWX_PROXY        = "SWX_PROXY"
	HLR_PROXY        = "HLR_PROXY"
	HEALTH           = "HEALTH"
	CSFB             = "CSFB"
	FEG_HELLO        = "FEG_HELLO"
	AAA_SERVER       = "AAA_SERVER"
	ENVOY_CONTROLLER = "ENVOY_CONTROLLER"
	EAP              = "EAP"
	EAP_SIM          = "EAP_SIM"
	EAP_AKA          = "EAP_AKA"
	RADIUSD          = "RADIUSD"
	RADIUS           = "RADIUS"
	REDIS            = "REDIS"
	PIPELINED        = "PIPELINED"
	MOCK_VLR         = "MOCK_VLR"
	MOCK_OCS         = "MOCK_OCS"
	MOCK_OCS2        = "MOCK_OCS2"
	MOCK_PCRF        = "MOCK_PCRF"
	MOCK_PCRF2       = "MOCK_PCRF2"
	MOCK_HSS         = "HSS"

	SESSION_MANAGER = "SESSIOND"
)

// Add a new service.
// If the service already exists, overwrites the service config.
func AddService(serviceType, host string, port int) {
	fegRegistry.AddService(platform_registry.ServiceLocation{Name: serviceType, Host: host, Port: port})
}

// Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return fegRegistry.GetServiceAddress(service)
}

// Provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	return fegRegistry.GetConnection(service)
}

func addLocalService(serviceType string, port int) {
	AddService(serviceType, "localhost", port)
}

var fegRegistry = Get()

func init() {

	// Add default Local Service Locations
	addLocalService(REDIS, 6380)

	addLocalService(FEG_HELLO, 9093)
	addLocalService(SESSION_PROXY, 9097)
	addLocalService(S6A_PROXY, 9098)
	addLocalService(S8_PROXY, 9099)
	addLocalService(CSFB, 9101)
	addLocalService(HEALTH, 9107)

	addLocalService(RADIUS, 9108)
	addLocalService(EAP, 9109)
	addLocalService(AAA_SERVER, 9109)
	addLocalService(EAP_SIM, 9118)
	addLocalService(EAP_AKA, 9123)
	addLocalService(SWX_PROXY, 9110)
	addLocalService(RADIUSD, 9115)
	addLocalService(HLR_PROXY, 9116)
	addLocalService(PIPELINED, 9117)
	addLocalService(ENVOY_CONTROLLER, 9118)

	addLocalService(MOCK_OCS, 9201)
	addLocalService(MOCK_PCRF, 9202)
	addLocalService(MOCK_OCS2, 9205)
	addLocalService(MOCK_PCRF2, 9206)
	addLocalService(MOCK_VLR, 9203)
	addLocalService(MOCK_HSS, 9204)

	// Overwrite/Add from /etc/magma/service_registry.yml if it exists
	// moduleName is "" since all feg configs lie in /etc/magma without a module name
	locations, err := registry.LoadServiceRegistryConfig("")
	if err != nil {
		log.Printf("Error loading FeG service_registry.yml: %v", err)
	} else if len(locations) > 0 {
		fegRegistry.AddServices(locations...)
	}
}
