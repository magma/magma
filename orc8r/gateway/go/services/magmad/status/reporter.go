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

// package status implements magmad status amd metrics collectors & reporters
package status

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang/glog"

	"magma/gateway/config"
	"magma/gateway/mconfig"
	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	mconfig_proto "magma/orc8r/lib/go/protos/mconfig"
)

const (
	MinCheckinIntervalSeconds     = 30
	DefaultCheckinIntervalSeconds = MinCheckinIntervalSeconds * 2
	serviceCollectDelay           = time.Second * 10
)

// StartReporter starts state collection & reporting loop
// StartReporter also initiates periodic services metrics collection. Since status & metrics collections are done using
// the same GRPC clients and goroutines, status collection timer is the main driver of collections intervals,
// status collection and reporting will be done according to magmad's config checkin interval settings.
// Metrics collection & reporting is also done during status collection/reporting window, but an
// additiona condition of time elapsed from last metrics collection/reporting is equal or longer then
// Metricsd.CollectInterval, Metricsd.SyncInterval respectively must be also satisfied.
// It's therefore preferable to set Metricsd.CollectInterval, Metricsd.SyncInterval to be multiples
// of mconf.CheckinInterval
// StartReporter never returns, it'll log errors if any and continue
func StartReporter() {
	var (
		timer          = time.NewTimer(time.Second * time.Duration(DefaultCheckinIntervalSeconds))
		metricsEnabled bool

		lastMetricsCollection, lastMetricsReporting time.Time
	)

	for {
		mdc := config.GetMagmadConfigs()

		fb303services := mdc.MagmaServices
		nonFb303Services := map[string]struct{}{}
		for _, s := range mdc.NonService303Services {
			s := strings.ToLower(s)
			nonFb303Services[s] = struct{}{}
		}

		metricsCollectInterval, metricsSyncInterval := mdc.Metricsd.CollectInterval, mdc.Metricsd.SyncInterval
		if metricsCollectInterval < MinCheckinIntervalSeconds {
			metricsCollectInterval = MinCheckinIntervalSeconds
		}
		if metricsSyncInterval < metricsCollectInterval {
			metricsSyncInterval = metricsCollectInterval
		}
		metricsServices := map[string]bool{}
		nextMetricsCollectionTime := lastMetricsCollection.Add(time.Second * time.Duration(metricsCollectInterval))
		newMetricsEnabled := mdc.Metricsd.QueueLength > 0
		if newMetricsEnabled != metricsEnabled {
			resetMetricsQueue()
			metricsEnabled = newMetricsEnabled
		}
		now := time.Now()
		// use !After() to check if it's <= now
		if metricsEnabled && !nextMetricsCollectionTime.After(now) {
			lastMetricsCollection = now
			for _, s := range mdc.Metricsd.Services {
				s := strings.ToLower(s)
				metricsServices[s] = true
			}
		}
		for _, fb303service := range fb303services {
			fb303service := strings.ToLower(fb303service)
			if _, nonFb303 := nonFb303Services[fb303service]; !nonFb303 {
				err := startServiceQuery(fb303service, metricsServices[fb303service], mdc.Metricsd.QueueLength)
				if err != nil {
					glog.Errorf("error querying service '%s' state: %v", fb303service, err)
				}
			}
		}
		time.Sleep(serviceCollectDelay)

		serviceStates := collect()
		stateConn, err := service_registry.Get().GetSharedCloudConnection(definitions.StateServiceName)
		if err != nil {
			glog.Errorf("failed to connect to state reporting service: %v", err)
		} else {
			res, err := protos.NewStateServiceClient(stateConn).ReportStates(context.Background(), serviceStates)
			if err != nil {
				glog.Errorf("ReportStates error: %v", err)
			} else if len(res.GetUnreportedStates()) > 0 {
				resStr, _ := json.Marshal(res.GetUnreportedStates())
				glog.Warningf("status unreported states: %s", resStr)
			} else {
				glog.V(1).Info("states report success")
			}
		}
		nextMetricsSyncTime := lastMetricsReporting.Add(time.Second * time.Duration(metricsSyncInterval))
		now = time.Now()
		// use !After() to check if it's <= now
		if metricsEnabled && !nextMetricsSyncTime.After(now) {
			lastMetricsReporting = now
			err = reportMetrics(mdc.Metricsd.QueueLength)
			if err != nil {
				glog.Errorf("metrics reporting error: %v", err)
			} else {
				glog.V(1).Info("metrics report success")
			}
		}
		<-timer.C // wait on timer for the remainder of intervalSeconds

		// update timer based on the latest configs
		intervalSeconds := DefaultCheckinIntervalSeconds
		mconf := &mconfig_proto.MagmaD{}
		err = mconfig.GetServiceConfigs(definitions.MagmadServiceName, mconf)
		if err == nil && mconf.CheckinInterval > 0 {
			intervalSeconds = int(mconf.CheckinInterval)
			if intervalSeconds < MinCheckinIntervalSeconds {
				intervalSeconds = MinCheckinIntervalSeconds
			}
		}
		timer.Reset(time.Second * time.Duration(intervalSeconds))
	}
}
