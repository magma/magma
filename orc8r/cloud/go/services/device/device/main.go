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
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/services/device/servicers"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, device.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating device service %s", err)
	}
	db, err := sqorc.Open(storage2.GetSQLDriver(), storage2.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	store := blobstore.NewEntStorage(device.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Failed to initialize device database: %s", err)
	}
	// Add servicers to the service
	deviceServicer, err := servicers.NewDeviceServicer(store)
	if err != nil {
		glog.Fatalf("Failed to instantiate the device servicer: %v", deviceServicer)
	}
	protos.RegisterDeviceServer(srv.GrpcServer, deviceServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Failed to start device service: %v", err)
	}
}
