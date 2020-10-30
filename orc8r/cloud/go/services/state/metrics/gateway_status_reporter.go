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
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
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
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return err
	}
	for _, networkID := range networks {
		gateways, _, err := configurator.LoadEntities(
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
			status, err := wrappers.GetGatewayStatus(networkID, gatewayEntity.PhysicalID)
			if err != nil {
				if err != errors.ErrNotFound {
					glog.Errorf("Error getting gateway state for network:%v, gateway:%v, %v", networkID, gatewayID, err)
				}
				continue
			}
			// last check in more than 5 minutes ago
			if (time.Now().Unix() - int64(status.CheckinTime)/1000) > 60*5 {
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
