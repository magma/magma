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

// Package collection provides functions used by the health manager to collect
// health related metrics for FeG services and the system
package collection

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/service_health"
)

// CollectServiceStats fills out the ServiceHealthStats proto for the provided service
// If the service cannot be reached, the service state is listed as UNAVAILABLE
func CollectServiceStats(serviceType string) *protos.ServiceHealthStats {
	healthStatus, err := service_health.GetHealthStatus(serviceType)
	if err != nil {
		if healthStatus != nil {
			return &protos.ServiceHealthStats{
				ServiceState:        protos.ServiceHealthStats_AVAILABLE,
				ServiceHealthStatus: healthStatus,
			}
		}
		return &protos.ServiceHealthStats{
			ServiceState: protos.ServiceHealthStats_UNAVAILABLE,
			ServiceHealthStatus: &protos.HealthStatus{
				Health:        protos.HealthStatus_UNHEALTHY,
				HealthMessage: "Service unavailable",
			},
		}
	}
	return &protos.ServiceHealthStats{
		ServiceState:        protos.ServiceHealthStats_AVAILABLE,
		ServiceHealthStatus: healthStatus,
	}
}
