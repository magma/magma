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
	"io/ioutil"
	"time"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/bootstrapper"
	bootstrapper_config "magma/orc8r/cloud/go/services/bootstrapper/config"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/key"
)

var (
	keyFilepath    = flag.String("cak", "bootstrapper.key.pem", "Bootstrapper's Private Key file")
	rootCAFilepath = "/var/opt/magma/certs/rootCA.pem"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, bootstrapper.ServiceName)
	if err != nil {
		glog.Fatalf("error creating service: %+v", err)
	}

	bootstrapperServicer := createBootstrapperServicer()
	cloudRegistrationServicer, registrationServicer := createRegistrationServicers(srv)

	protos.RegisterBootstrapperServer(srv.GrpcServer, bootstrapperServicer)
	protos.RegisterCloudRegistrationServer(srv.GrpcServer, cloudRegistrationServicer)
	protos.RegisterRegistrationServer(srv.GrpcServer, registrationServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("error running service: %+v", err)
	}
}

func createBootstrapperServicer() *servicers.BootstrapperServer {
	key, err := key.ReadKey(*keyFilepath)
	if err != nil {
		glog.Fatalf("error reading bootstrapper private key: %+v", err)
	}
	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		glog.Fatalf("error coercing bootstrapper private key to RSA private key; actual type: %T", key)
	}

	servicer, err := servicers.NewBootstrapperServer(rsaPrivateKey)
	if err != nil {
		glog.Fatalf("error creating bootstrapper server: %+v", err)
	}
	return servicer
}

func createRegistrationServicers(srv *service.OrchestratorService) (protos.CloudRegistrationServer, protos.RegistrationServer) {
	db, err := sqorc.Open(storage2.GetSQLDriver(), storage2.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("failed to connect to database: %+v", err)
	}
	factory := blobstore.NewSQLStoreFactory(bootstrapper.BlobstoreTableName, db, sqorc.GetSqlBuilder())
	err = factory.InitializeFactory()
	if err != nil {
		glog.Fatalf("error initializing bootstrapper database: %+v", err)
	}
	store := registration.NewBlobstoreStore(factory)

	rootCA, err := getRootCA()
	if err != nil {
		glog.Fatalf("failed to get rootCA: %+v", err)
	}

	timeoutDurationInMinutes := srv.Config.MustGetInt(bootstrapper_config.TokenTimeoutDurationInMinutes)
	timeout := time.Duration(timeoutDurationInMinutes) * time.Minute
	cloudRegistrationServicer, err := registration.NewCloudRegistrationServicer(store, rootCA, timeout)
	if err != nil {
		glog.Fatalf("error creating cloud registration servicer: %+v", err)
	}

	registrationServicer := registration.NewRegistrationServicer()

	return cloudRegistrationServicer, registrationServicer
}

func getRootCA() (string, error) {
	body, err := ioutil.ReadFile(rootCAFilepath)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
