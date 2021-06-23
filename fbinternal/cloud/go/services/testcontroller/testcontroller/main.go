/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net/http"
	"time"

	"magma/fbinternal/cloud/go/fbinternal"
	"magma/fbinternal/cloud/go/serdes"
	"magma/fbinternal/cloud/go/services/testcontroller"
	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/handlers"
	"magma/fbinternal/cloud/go/services/testcontroller/protos"
	"magma/fbinternal/cloud/go/services/testcontroller/servicers"
	"magma/fbinternal/cloud/go/services/testcontroller/statemachines"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(fbinternal.ModuleName, testcontroller.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	db, err := sqorc.Open(storage2.GetSQLDriver(), storage2.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	defer db.Close()

	// Add node leasor servicer
	nodeStore := storage.NewSQLNodeLeasorStorage(db, &storage2.UUIDGenerator{}, sqorc.GetSqlBuilder())
	err = nodeStore.Init()
	if err != nil {
		glog.Fatalf("failed to initialize CI node storage: %s", err)
	}
	nodeServicer := servicers.NewNodeLeasorServicer(nodeStore)
	protos.RegisterNodeLeasorServer(srv.GrpcServer, nodeServicer)

	// Add testcontroller servicer
	e2eStore := storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
	err = e2eStore.Init()
	if err != nil {
		glog.Fatalf("failed to initialize testcontroller storage: %s", err)
	}
	e2eServicer := servicers.NewTestControllerServicer(e2eStore)
	protos.RegisterTestControllerServer(srv.GrpcServer, e2eServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(testcontroller.ServiceName))

	// Instantiate state machines, start test execution loop
	go func() {
		magmadClient := &statemachines.MagmadClient{}
		machines := map[string]statemachines.TestMachine{
			testcontroller.EnodedTestCaseType:       statemachines.NewEnodebdE2ETestStateMachine(e2eStore, http.DefaultClient, magmadClient),
			testcontroller.EnodedTestExcludeTraffic: statemachines.NewEnodebdE2ETestStateMachineNoTraffic(e2eStore, http.DefaultClient, magmadClient),
		}

		for {
			err := testcontroller.ExecuteNextTestCase(machines, e2eStore, serdes.TestController)
			if err != nil {
				glog.Error(err)
			}
			// 10 second sleep so we don't hot-loop
			time.Sleep(10 * time.Second)
		}
	}()
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers())
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Failed to run test controller service: %s", err)
	}
}
