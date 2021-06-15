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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/exporter"
	manager "magma/lte/cloud/go/services/nprobe/nprobe_manager"
	"magma/lte/cloud/go/services/nprobe/obsidian/handlers"
	np_storage "magma/lte/cloud/go/services/nprobe/storage"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(lte.ModuleName, nprobe.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %v", err)
	}

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	fact := blobstore.NewSQLBlobStorageFactory(nprobe.NProbeTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing nprobe table: %+v", err)
	}
	nprobeBlobstore := np_storage.NewNProbeBlobstore(fact)

	// Attach handlers
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers(nprobeBlobstore))
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(nprobe.ServiceName))

	serviceConfig := nprobe.GetServiceConfig()
	tlsConfig, err := exporter.NewTlsConfig(
		serviceConfig.ExporterCrtFile,
		serviceConfig.ExporterKeyFile,
		serviceConfig.SkipVerifyServer,
	)
	if err != nil {
		glog.Fatalf("Failed to create new TlsConfig: %v", err)
	}

	// Init records exporter
	recordExporter := exporter.NewRecordExporter(serviceConfig.DeliveryFunctionAddr, tlsConfig)
	nProbeManager, err := manager.NewNProbeManager(serviceConfig, nprobeBlobstore, recordExporter)
	if err != nil {
		glog.Fatalf("Failed to create new NProbeManager: %v", err)
	}

	// Run LI service in Loop
	go func() {
		for {
			err := nProbeManager.ProcessNProbeTasks()
			if err != nil {
				glog.Errorf("Failed to process tasks: %v", err)
				<-time.After(time.Duration(serviceConfig.BackOffIntervalSecs) * time.Second)
			}
			<-time.After(time.Duration(serviceConfig.UpdateIntervalSecs) * time.Second)
		}
	}()

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %v", err)
	}
}
