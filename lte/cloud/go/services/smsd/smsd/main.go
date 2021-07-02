/*
 *  Copyright 2020 The Magma Authors.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/smsd"
	"magma/lte/cloud/go/services/smsd/servicers"
	storage2 "magma/lte/cloud/go/services/smsd/storage"
	"magma/lte/cloud/go/sms_ll"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, smsd.ServiceName)
	if err != nil {
		glog.Fatalf("error creating smsd service: %v", err)
	}

	// Storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("error opening db conn: %v", err)
	}
	store := storage2.NewSQLSMSStorage(db, sqorc.GetSqlBuilder(), &storage2.DefaultSMSReferenceCounter{}, &storage.UUIDGenerator{})
	err = store.Init()
	if err != nil {
		glog.Fatalf("error initializing smsd storage: %s", err)
	}

	restServicer := servicers.NewRESTServicer(store)
	obsidian.AttachHandlers(srv.EchoServer, restServicer.GetHandlers())
	protos.RegisterSmsDServer(srv.GrpcServer, servicers.NewSMSDServicer(store, &sms_ll.DefaultSMSSerde{}))

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(smsd.ServiceName))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("error while running smsd service: %v", err)
	}
}
