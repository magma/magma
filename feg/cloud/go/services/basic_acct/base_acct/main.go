/*
Copyright 2021 The Magma Authors.

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

	"github.com/golang/glog"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/basic_acct"
	"magma/feg/cloud/go/services/basic_acct/servicers"
	"magma/orc8r/cloud/go/service"
)

func main() {
	flag.Parse()

	// Create the service
	srv, err := service.NewOrchestratorService(feg.ModuleName, basic_acct.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating Basic Accounting service: %s", err)
	}

	bas := servicers.NewBaseAcctService()
	cfg := bas.GetConfig()
	if cfg == nil {
		cfg = &servicers.Config{}
	}
	glog.Infof("Starting %s service with configs: %+v", basic_acct.ServiceName, *cfg)
	protos.RegisterAccountingServer(srv.GrpcServer, bas)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running Basic Accounting service: %s", err)
	}
}
