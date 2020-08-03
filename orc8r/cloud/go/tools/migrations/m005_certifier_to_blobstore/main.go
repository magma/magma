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
	"context"
	"flag"

	certifierprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/tools/migrations"
	"magma/orc8r/lib/go/protos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Certifier-specific
	placeholderNID       = "placeholder_network"
	oldCertInfoTable     = "certificate_info_db"
	newCertInfoTable     = "certificate_info_blobstore"
	certifierStorageType = "certificate_info"
)

var (
	verify bool
)

// main runs the certifier_to_blobstore migration.
// Migration without SQL error will result in `SUCCESS` printed to stderr.
//
// Optional `-verify` flag logs values from certifier service, post-migration,
// allowing manual evaluation.
func main() {
	flag.BoolVar(&verify, "verify", false, "partially verify successful migration via certifier RPC")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")
	migrations.MigrateNetworkAgnosticServiceToBlobstore(
		placeholderNID, certifierStorageType, oldCertInfoTable, newCertInfoTable,
	)
	glog.Info("END MIGRATION")

	if verify {
		err := manuallyVerifyCertifierMigration()
		if err != nil {
			glog.Errorf("certifier verification failed: %s", err)
		}
	}

}

// manuallyVerifyCertifierMigration attempts to log several (key, value) pairs from
// the certifier service to allow for manual verification of successful migration. This
// ensures the certifier service is able to understand newly-migrated data.
//
// While not necessary, migrators may choose to have the relevant values
// (attained via the certifier client_api) printed to manually compare with the
// values from the old and new tables.
func manuallyVerifyCertifierMigration() error {
	// NOTE: we hard-code the certifier server location and bypass certifier/client_api.go
	// to avoid using application code as much as possible (registry etc.)
	certSrvAddr := "localhost:9086"
	conn, err := grpc.Dial(certSrvAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to connect to certifier server")
	}
	client := certifierprotos.NewCertifierClient(conn)

	snsProto, err := client.ListCertificates(context.Background(), &protos.Void{})
	if err != nil {
		return errors.Wrap(err, "failed to list certificates")
	}
	glog.Infof("[manually verify] serial number count: %d", len(snsProto.GetSns()))

	certInfos, err := client.GetAll(context.Background(), &protos.Void{})
	if err != nil {
		return errors.Wrap(err, "failed to get all certificates")
	}
	maxToPrint := 5
	i := 0
	for sn, info := range certInfos.Certificates {
		if i++; i > maxToPrint {
			break
		}
		glog.Infof("[manually verify] key-value pair: {key: %s, value: %+v}", sn, info)
	}
	return nil
}
