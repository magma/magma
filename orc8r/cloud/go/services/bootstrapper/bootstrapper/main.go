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
	"log"

	"github.com/getsentry/sentry-go"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/key"
)

var (
	keyFile = flag.String("cak", "bootstrapper.key.pem", "Bootstrapper's Private Key file")
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f6a54d1a20134c258b1e0b227d4d0982@o529355.ingest.sentry.io/5667116",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Create the service, flag will be parsed inside this function
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, bootstrapper.ServiceName)
	if err != nil {
		log.Fatalf("Error creating bootstrapper service: %s", err)
	}

	// Add servicers to the service
	privKey, err := key.ReadKey(*keyFile)
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}
	servicer, err := servicers.NewBootstrapperServer(privKey.(*rsa.PrivateKey))
	if err != nil {
		log.Fatalf("Failed to create bootstrapper servicer: %s", err)
	}
	protos.RegisterBootstrapperServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
