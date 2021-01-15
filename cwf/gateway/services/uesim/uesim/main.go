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

// This starts the user equipment (ue) service.
package main

import (
	"flag"

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func createUeSimServer(store blobstore.BlobStorageFactory) (protos.UESimServer, error) {
	config, err := servicers.GetUESimConfig()
	if err != nil {
		glog.Fatalf("Error getting UESim Config : %s ", err)
		return nil, err
	}
	if servicers.GetBypassHssFlag(config) == true {
		return servicers.NewUESimServerHssLess(store)
	} else {
		return servicers.NewUESimServer(store)
	}
}

func main() {
	flag.Parse()
	glog.Info("Starting UESim service")

	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.UeSim)
	if err != nil {
		glog.Fatalf("Error creating UeSim service: %s", err)
	}

	store, err := test_utils.NewSQLBlobstoreForServices("uesim_main_blobstore")
	if err != nil {
		glog.Fatalf("Error creating in-memory blobstore: %+v", err)
	}
	servicer, err := createUeSimServer(store)

	protos.RegisterUESimServer(srv.GrpcServer, servicer)
	if err != nil {
		glog.Fatalf("Error creating UE server: %s", err)
	}
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running UE service: %s", err)
	}
}
