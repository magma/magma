/*
Copyright 2022 The Magma Authors.

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
	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/eps_authentication"
	"magma/lte/cloud/go/services/eps_authentication/servicers"
	eps_storage "magma/lte/cloud/go/services/eps_authentication/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
)

// eps_authentication service
func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, eps_authentication.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}
	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	stateStoreFactory := blobstore.NewSQLStoreFactory(eps_storage.EpsAuthStateStore, db, sqorc.GetSqlBuilder())
	if err := stateStoreFactory.InitializeFactory(); err != nil {
		glog.Fatalf("Error initializing EPS Authentication storage: %+v", err)
	}
	// Add servicers to the service
	store := eps_storage.NewSubscriberDBStorage(stateStoreFactory)
	servicer, err := servicers.NewEPSAuthServer(store)
	if err != nil {
		glog.Fatalf("EPS Auth Servicer Initialization Error: %s", err)
	}
	protos.RegisterS6AProxyServer(srv.GrpcServer, servicer) // EPS Auth server implements S6a Proxy API

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
