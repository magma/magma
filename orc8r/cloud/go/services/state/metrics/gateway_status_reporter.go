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

package metrics

import (
	"context"
	"time"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/lib/go/merrors"
)

const (
	// metricsGracePeriodSeconds is the grace period for gateway health check.
	// If a gateway's last check-in was longer ago than this value, it is
	// considered unhealthy and reported as status=0 in Prometheus metrics.
	//
	// This value should match the API's grace period calculation in
	// orc8r/cloud/go/services/orchestrator/obsidian/models/conversion.go:
	//   graceFactor (10) * defaultCheckinInterval (60s) = 600 seconds
	//
	// NOTE: This is a simplified constant that assumes the default checkin_interval.
	// For deployments with custom checkin_interval, this may not exactly match
	// the API's dynamic calculation. See Issue #15725 for details.
	//
	// Previously this was hardcoded as 60*5=300 seconds, which caused status
	// fluctuation in the NMS UI when a gateway was between 300-600 seconds
	// since last check-in.
	metricsGracePeriodSeconds = 10 * 60 // 600 seconds
)

func PeriodicallyReportGatewayStatus(dur time.Duration) {
	for range time.Tick(dur) {
		err := reportGatewayStatus()
		if err != nil {
			glog.Errorf("err in reportGatewayStatus: %v\n", err)
		}
	}
}

func reportGatewayStatus() error {
	networks, err := configurator.ListNetworkIDs(context.Background())
	if err != nil {
		return err
	}
	for _, networkID := range networks {
		gateways, _, err := configurator.LoadEntities(
			context.Background(),
			networkID,
			swag.String(orc8r.MagmadGatewayType),
			nil,
			nil,
			nil,
			configurator.EntityLoadCriteria{},
			serdes.Entity,
		)
		if err != nil {
			glog.Errorf("error getting gateways for network %v: %v\n", networkID, err)
			continue
		}
		numUpGateways := 0
		for _, gatewayEntity := range gateways {
			gatewayID := gatewayEntity.Key
			status, err := wrappers.GetGatewayStatus(context.Background(), networkID, gatewayEntity.PhysicalID)
			if err != nil {
				if err != merrors.ErrNotFound {
					glog.Errorf("Error getting gateway state for network:%v, gateway:%v, %v", networkID, gatewayID, err)
				}
				continue
			}
			// Check if last check-in was too long ago
			// Use >= to match the API behavior in conversion.go which uses time.Now().Before()
			// At the exact boundary, both API and metrics should report unhealthy
			if (time.Now().Unix() - int64(status.CheckinTime)/1000) >= metricsGracePeriodSeconds {
				gwCheckinStatus.WithLabelValues(networkID, gatewayID).Set(0)
			} else {
				gwCheckinStatus.WithLabelValues(networkID, gatewayID).Set(1)
				numUpGateways += 1
			}

			// report mconfig age
			if status.PlatformInfo != nil && status.PlatformInfo.ConfigInfo != nil {
				mconfigCreatedAt := status.PlatformInfo.ConfigInfo.MconfigCreatedAt
				if mconfigCreatedAt != 0 {
					gwMconfigAge.WithLabelValues(networkID, gatewayID).Set(float64(status.CheckinTime/1000 - mconfigCreatedAt))
				}
			} else {
				glog.Errorf("Status for networkID %s, gatewayID %s is missing the MconfigCreatedAt field", networkID, gatewayID)
			}
		}
		upGwCount.WithLabelValues(networkID).Set(float64(numUpGateways))
		totalGwCount.WithLabelValues(networkID).Set(float64(len(gateways)))
	}
	return nil
}
