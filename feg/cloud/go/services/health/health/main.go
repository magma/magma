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

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/reporter"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

const (
	NETWORK_HEALTH_STATUS_REPORT_INTERVAL = time.Second * 60
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(feg.ModuleName, health.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Failed to connect to database: %+v", err)
	}
	store := blobstore.NewEntStorage(health.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing health database: %+v", err)
	}
	// Add servicers to the service
	healthServer, err := servicers.NewHealthServer(store)
	if err != nil {
		glog.Fatalf("Error creating health servicer: %+v", err)
	}
	protos.RegisterHealthServer(srv.GrpcServer, healthServer)

	// create a networkHealthStatusReporter to monitor and periodically log metrics
	// on if all gateways in a network are unhealthy
	healthStatusReporter := &reporter.NetworkHealthStatusReporter{}
	go healthStatusReporter.ReportHealthStatus(NETWORK_HEALTH_STATUS_REPORT_INTERVAL)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running health service: %+v", err)
	}
}
