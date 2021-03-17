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
	"crypto/rsa"
	"flag"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/key"

	"github.com/golang/glog"
)

var (
	keyFilepath = flag.String("cak", "bootstrapper.key.pem", "Bootstrapper's Private Key file")
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, bootstrapper.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	key, err := key.ReadKey(*keyFilepath)
	if err != nil {
		glog.Fatalf("Error reading bootstrapper private key: %+v", err)
	}
	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		glog.Fatalf("Error coercing bootstrapper private key to RSA private key; actual type: %T", key)
	}

	servicer, err := servicers.NewBootstrapperServer(rsaPrivateKey)
	if err != nil {
		glog.Fatalf("Error creating bootstrapper servicer: %+v", err)
	}
	protos.RegisterBootstrapperServer(srv.GrpcServer, servicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %+v", err)
	}
}
