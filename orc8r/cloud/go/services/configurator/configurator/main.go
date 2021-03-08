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

/*
Configurator is a dedicated Magma Cloud service which maintains configurations
and meta data for the network and network entity structures.
*/

package main

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/servicers"
	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

const (
	maxEntityLoadSizeConfigKey = "maxEntityLoadSize"
)

func bar() int {
	var luckyNumber []int
	return luckyNumber[42]
}

func foo() int {
	return bar()
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f6a54d1a20134c258b1e0b227d4d0982@o529355.ingest.sentry.io/5667116",
	})
	if err != nil {
		glog.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()
	foo()
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, configurator.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}
	db, err := sqorc.Open(storage2.SQLDriver, storage2.DatabaseSource)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	maxEntityLoadSize, err := srv.Config.GetInt(maxEntityLoadSizeConfigKey)
	if err != nil {
		glog.Fatalf("Failed to load '%s' from config: %s", maxEntityLoadSizeConfigKey, err)
	}
	factory := storage.NewSQLConfiguratorStorageFactory(db, &storage2.UUIDGenerator{}, sqorc.GetSqlBuilder(), uint32(maxEntityLoadSize))
	err = factory.InitializeServiceStorage()
	if err != nil {
		glog.Fatalf("Failed to initialize configurator database: %s", err)
	}

	nbServicer, err := servicers.NewNorthboundConfiguratorServicer(factory)
	if err != nil {
		glog.Fatalf("Failed to instantiate the user-facing configurator servicer: %v", nbServicer)
	}
	protos.RegisterNorthboundConfiguratorServer(srv.GrpcServer, nbServicer)

	sbServicer, err := servicers.NewSouthboundConfiguratorServicer(factory)
	if err != nil {
		glog.Fatalf("Failed to instantiate the device-facing configurator servicer: %v", sbServicer)
	}
	protos.RegisterSouthboundConfiguratorServer(srv.GrpcServer, sbServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Failed to start configurator service: %v", err)
	}
}
