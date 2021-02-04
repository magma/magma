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
	"flag"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/interceptd"
	"magma/orc8r/cloud/go/services/interceptd/collector"
	manager "magma/orc8r/cloud/go/services/interceptd/intercept_manager"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, interceptd.ServiceName)
	if err != nil {
		glog.Fatalf("Failed running Interceptd service: %v", err)
	}

	serviceConfig := interceptd.GetServiceConfig()
	eventsCollector, err := collector.NewEventsCollector()
	if err != nil {
		glog.Fatalf("Failed to create new EventsCollector: %v", err)
	}

	interceptManager := manager.NewInterceptManager(eventsCollector, serviceConfig)
	// Run LI service in Loop
	go func() {
		for {
			interceptManager.CollectAndProcessEvents()
			<-time.After(time.Duration(serviceConfig.UpdateIntervalSecs) * time.Second)
		}
	}()

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
