/*
 * Copyright 2020 The Magma Authors. *
 *
 * This source code is licensed under the BSD-style license found in the *
 * LICENSE file in the root directory of this source tree. *
 *
 * Unless required by applicable law or agreed to in writing, software *
 * distributed under the License is distributed on an "AS IS" BASIS, *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and *
 * limitations under the License. *
 */

package status_reporter

import (
	"fmt"
	"os"

	"magma/cwf/cloud/go/services/cwf/obsidian/models"
	"magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"
	"magma/feg/cloud/go/protos"
	"magma/gateway/redis"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("hacluster_status_reporter")

// StatusReporter is composed of redis clients that allow the operator to report
// HA cluster status and gateway health
type StatusReporter struct {
	redisHealthClient  *redis.RedisStateClient
	redisClusterClient *redis.RedisStateClient
}

const (
	redisEnvVar     = "REDIS_ADDR"
	healthStateType = "cwf_gateway_health"
	haPairStateType = "cwf_ha_pair_status"
)

func NewStatusReporter() *StatusReporter {
	redisAddr := os.Getenv(redisEnvVar)
	if len(redisAddr) == 0 {
		log.Error(fmt.Errorf("%s env variable not set", redisEnvVar), "Operator will not report HACluster status to Orchestrator!")
		return &StatusReporter{
			redisHealthClient:  nil,
			redisClusterClient: nil,
		}
	}
	gatewayHealthStateSerde := redis.NewJsonStateSerde(healthStateType, &models.CarrierWifiGatewayHealthStatus{})
	haPairStateSerde := redis.NewJsonStateSerde(haPairStateType, &models.CarrierWifiHaPairStatus{})
	return &StatusReporter{
		redisHealthClient:  redis.NewDefaultRedisStateClient(redisAddr, gatewayHealthStateSerde),
		redisClusterClient: redis.NewDefaultRedisStateClient(redisAddr, haPairStateSerde),
	}
}

// UpdateHAClusterStatus updates redis with the cluster's gateway health and
// active/standby status
func (r *StatusReporter) UpdateHAClusterStatus(
	status v1alpha1.HAClusterStatus,
	spec v1alpha1.HAClusterSpec,
	activeHealth *protos.HealthStatus,
	standbyHealth *protos.HealthStatus,
) {
	if r.redisHealthClient == nil {
		return
	}
	var activeGatewayID string
	var standbyGatewayID string
	for _, resource := range spec.GatewayResources {
		if resource.HelmReleaseName == status.Active {
			activeGatewayID = resource.GatewayID
		} else {
			standbyGatewayID = resource.GatewayID
		}
	}
	err := r.updateHaPairStatus(spec.HAPairID, activeGatewayID)
	if err != nil {
		log.Error(err, "")
	}
	err = r.updateGatewayHealthStatus(activeGatewayID, activeHealth)
	if err != nil {
		log.Error(err, "")
	}
	err = r.updateGatewayHealthStatus(standbyGatewayID, standbyHealth)
	if err != nil {
		log.Error(err, "")
	}
}

func (r *StatusReporter) updateHaPairStatus(pairID string, activeGateway string) error {
	haPairStatus := models.CarrierWifiHaPairStatus{
		ActiveGateway: activeGateway,
	}
	return r.redisClusterClient.Set(pairID, haPairStatus)
}

func (r *StatusReporter) updateGatewayHealthStatus(gatewayID string, healthStatus *protos.HealthStatus) error {
	if healthStatus == nil {
		return fmt.Errorf("Health status for gateway %s is nil", gatewayID)
	}
	healthValue := models.CarrierWifiGatewayHealthStatus{
		Status:      healthStatus.Health.String(),
		Description: healthStatus.HealthMessage,
	}
	return r.redisHealthClient.Set(gatewayID, healthValue)
}
