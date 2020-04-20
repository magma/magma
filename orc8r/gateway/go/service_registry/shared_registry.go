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
// package cloud_registry provides CloudRegistry interface for Go based gateways
package service_registry

import (
	"sync/atomic"

	"github.com/golang/glog"

	_ "magma/orc8r/lib/go/initflag"
	platform_registry "magma/orc8r/lib/go/registry"
)

// default service registry shared by all GW process services
var initializedRegistry atomic.Value

const serviceRegLoadErrorFmt = "Error loading Gateway service_registry.yml: %v"

// Get returns default service registry which can be shared by all GW process services
func Get() GatewayRegistry {
	reg, ok := initializedRegistry.Load().(GatewayRegistry)
	if (!ok) || reg == nil {
		reg = NewDefaultRegistry()
		// Overwrite/Add from /etc/magma/service_registry.yml if it exists
		// moduleName is "" since all feg configs lie in /etc/magma without a module name
		locations, err := platform_registry.LoadServiceRegistryConfig("")
		if err != nil {
			glog.Warningf(serviceRegLoadErrorFmt, err)
			// return registry, but don't store/cache it
			return reg
		}
		if len(locations) > 0 {
			reg.AddServices(locations...)
		}
		initializedRegistry.Store(reg)
	}
	return reg
}

// NewDefaultRegistry returns a new service registry populated by default services
func NewDefaultRegistry() GatewayRegistry {
	reg := platform_registry.Get()
	if _, err := reg.GetServiceAddress(platform_registry.ControlProxyServiceName); err != nil {
		addLocal(reg, platform_registry.ControlProxyServiceName, 5053)
	}
	return reg
}

// Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return Get().GetServiceAddress(service)
}

func addLocal(reg GatewayRegistry, serviceType string, port int) {
	reg.AddService(platform_registry.ServiceLocation{Name: serviceType, Host: "localhost", Port: port})
}
