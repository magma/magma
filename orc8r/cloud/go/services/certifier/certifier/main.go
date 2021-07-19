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
	"flag"
	"math/rand"
	"time"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/analytics"
	analytics_protos "magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/certifier"
	analytics_service "magma/orc8r/cloud/go/services/certifier/analytics"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/servicers"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/cert"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

var (
	bootstrapCACertFile = flag.String("cac", "server_cert.pem", "Signer CA's Certificate file")
	bootstrapCAKeyFile  = flag.String("cak", "server_cert.key.pem", "Signer CA's Private Key file")

	vpnCertFile = flag.String("vpnc", "vpn_ca.crt", "VPN CA's Certificate file")
	vpnKeyFile  = flag.String("vpnk", "vpn_ca.key", "VPN CA's Private Key file")

	gcHours = flag.Int64("gc-hours", 12, "Garbage Collection time interval (in hours)")
)

func main() {
	// Create the service, flag will be parsed inside this function
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, certifier.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	// Init storage
	db, err := sqorc.Open(storage2.GetSQLDriver(), storage2.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	fact := blobstore.NewEntStorage(storage.CertifierTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing certifier database: %s", err)
	}
	store := storage.NewCertifierBlobstore(fact)

	// Add servicers to the service
	caMap := map[protos.CertType]*servicers.CAInfo{}
	bootstrapCert, bootstrapPrivKey, err := cert.LoadCertAndPrivKey(*bootstrapCACertFile, *bootstrapCAKeyFile)
	if err != nil {
		glog.Infof("ERROR: Failed to load bootstrap CA cert and key: %v", err)
	} else {
		caMap[protos.CertType_DEFAULT] = &servicers.CAInfo{Cert: bootstrapCert, PrivKey: bootstrapPrivKey}
	}
	vpnCert, vpnPrivKey, vpnErr := cert.LoadCertAndPrivKey(*vpnCertFile, *vpnKeyFile)
	if vpnErr != nil {
		fmtstr := "ERROR: Failed to load VPN cert and key: %v"
		if err != nil {
			glog.Fatalf(fmtstr, vpnErr)
		} else {
			glog.Infof(fmtstr, vpnErr)
		}
	} else {
		caMap[protos.CertType_VPN] = &servicers.CAInfo{Cert: vpnCert, PrivKey: vpnPrivKey}
	}

	var serviceConfig certifier.Config
	_, _, err = config.GetStructuredServiceConfig(orc8r.ModuleName, certifier.ServiceName, &serviceConfig)
	if err != nil {
		glog.Fatalf("err %v failed parsing the config file: skipping CollectorServicer creation ", err)
	}
	collectorServicer := analytics.NewCollectorServicer(
		&serviceConfig.Analytics,
		analytics.GetPrometheusClient(),
		analytics_service.GetAnalyticsCalculations(&serviceConfig),
		nil,
	)
	analytics_protos.RegisterAnalyticsCollectorServer(srv.GrpcServer, collectorServicer)

	// Register servicer
	servicer, err := servicers.NewCertifierServer(store, caMap)
	if err != nil {
		glog.Fatalf("Failed to create certifier server: %s", err)
	}
	certprotos.RegisterCertifierServer(srv.GrpcServer, servicer)

	// Start Garbage Collector Ticker
	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			// wait for *gcHours +/- rand(1/20 of *gcHours)
			after := time.Hour * time.Duration(*gcHours)
			tenth := (after / 10) + 1 // +1 to make sure, it's not 0
			randomDelta := time.Duration(rand.Int63n(int64(tenth))) - tenth/2
			<-time.After(after + randomDelta)
			glog.Infof("removing stale certificates")
			count, err := servicer.CollectGarbageImpl(context.Background())
			if err != nil {
				glog.Errorf("error collecting garbage for certifier: %v", err)
			}
			glog.Infof("removed %d stale certificates", count)
		}
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
