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
	"magma/lte/cloud/go/lte"
	subscriberdb_storage "magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, subscriberdb_cache.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	digestStore := subscriberdb_storage.NewDigestLookup(db, sqorc.GetSqlBuilder())
	if err := digestStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing digest lookup storage: %+v", err)
	}

	serviceConfig := subscriberdb_cache.MustGetServiceConfig()
	glog.Infof("Subscriberdb_cache service config %+v", serviceConfig)

	go subscriberdb_cache.MonitorDigests(digestStore, serviceConfig)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %+v", err)
	}
}
