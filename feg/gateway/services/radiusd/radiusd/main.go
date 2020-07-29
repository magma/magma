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

package main

import (
	"time"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/radiusd/collection"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.RADIUSD)
	if err != nil {
		glog.Fatalf("Error creating RADIUSD service: %s", err)
	}
	metricAggregateRegistry := collection.NewMetricAggregateRegistry()
	metricsRequester, err := collection.NewMetricsRequester()
	if err != nil {
		glog.Fatalf("Error getting metrics requester: %s", err)
	}

	radiusdCfg := collection.GetRadiusdConfig()
	interval := radiusdCfg.GetUpdateIntervalSecs()
	// Run Radius metrics collection Loop
	go func() {
		for {
			<-time.After(time.Duration(interval) * time.Second)
			prometheusText, err := metricsRequester.FetchMetrics()
			if err != nil {
				glog.Errorf("Error getting metrics from server: %s", err)
				metricsRequester.RefreshConfig()
				interval *= 2
				continue
			}
			metricFamilies, err := collection.ParsePrometheusText(prometheusText)
			if err != nil {
				glog.Errorf("Unable to parse prometheus text: %s", err)
				interval *= 2
				continue
			}
			metricAggregateRegistry.Update(metricFamilies)
			interval = radiusdCfg.GetUpdateIntervalSecs()
		}
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
