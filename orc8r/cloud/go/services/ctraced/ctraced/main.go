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
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/ctraced"
	"magma/orc8r/cloud/go/services/ctraced/obsidian/handlers"
	"magma/orc8r/cloud/go/services/ctraced/servicers"
	ctraced_storage "magma/orc8r/cloud/go/services/ctraced/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, ctraced.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating ctraced service: %+v", err)
	}

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	fact := blobstore.NewSQLBlobStorageFactory(ctraced.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing ctraced table: %+v", err)
	}
	ctracedBlobstore := ctraced_storage.NewCtracedBlobstore(fact)

	// Init gRPC servicer
	protos.RegisterCallTraceControllerServer(srv.GrpcServer, servicers.NewCallTraceServicer(ctracedBlobstore))
	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(ctraced.ServiceName))

	gwClient := handlers.NewGwCtracedClient()
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers(gwClient, ctracedBlobstore))

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running ctraced service: %+v", err)
	}
}
