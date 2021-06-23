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

// Package main implements dual purpose directory service which manages UE location records and provides RPCs
// to look them up: DirectoryLookupServer & GatewayDirectoryService
//
// GatewayDirectoryService RPC can be provided by a local Gateway/Device service as well as the cloud hosted service.
// UE locations are reported directly from the relevant device/gateway.
// Depending on a gateway/device capacity & resources, UE locations can be reported using GatewayDirectoryService RPC:
//   1) to a locally run directoryd service, backed up by a local DB and synchronized with the cloud by state service
//   2) to a cloud based directoryd service
//
// Gateways/devices with an excess capacity to host and run local GatewayDirectoryService service & database (currently
// python based & requiring 40-70MB extra RAM & storage) should prefer reporting path #1 due to its higher persistency,
// and recovery capabilities.
// Embedded or underpowered gateways/devices on another hand, have to use reporting path #2, it's far more efficient in
// terms of required HW resources (less then a few hundred KB of extra RAM & storage), but its availability and
// correctness provided on a best-effort basis. Such a compromise should be a reasonable tradeoff for enabling location
// based features on underpowered devices.
// Note that that paths #1 & #2 can coexist within a single deployment/network, cloud's state service is a common
// destination for UE locations in both cases.
package main

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	dstorage "magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, directoryd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating directory service: %s", err)
	}

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %s", err)
	}

	fact := blobstore.NewEntStorage(dstorage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing directory storage: %s", err)
	}

	store := dstorage.NewDirectorydBlobstore(fact)

	// Add servicers
	// Cloud lookup service
	servicer, err := servicers.NewDirectoryLookupServicer(store)
	if err != nil {
		glog.Fatalf("Error creating initializing directory servicer: %s", err)
	}
	protos.RegisterDirectoryLookupServer(srv.GrpcServer, servicer)
	protos.RegisterGatewayDirectoryServiceServer(srv.GrpcServer, servicers.NewDirectoryUpdateServicer())

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running directory service: %s", err)
	}
}
