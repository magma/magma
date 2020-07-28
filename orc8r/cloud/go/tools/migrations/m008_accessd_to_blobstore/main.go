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

	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/tools/migrations"
	"magma/orc8r/lib/go/protos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Accessd-specific
	placeholderNID     = "placeholder_network"
	oldAccessdTable    = "access_control"
	newAccessdTable    = "access_control_blobstore"
	accessdDefaultType = "access_control"
)

var (
	verify bool
)

// main runs the accessd_to_blobstore migration.
// Migration without SQL error will result in `SUCCESS` printed to stderr.
//
// Optional `-verify` flag logs values from accessd service, post-migration,
// allowing manual evaluation.
func main() {
	flag.BoolVar(&verify, "verify", false, "partially verify successful migration via accessd RPC")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")
	migrations.MigrateNetworkAgnosticServiceToBlobstore(
		placeholderNID, accessdDefaultType, oldAccessdTable, newAccessdTable)
	glog.Info("END MIGRATION")

	if verify {
		err := manuallyVerifyAccessdMigration()
		if err != nil {
			glog.Errorf("accessd verification failed: %s", err)
		}
	}

}

// manuallyVerifyAccessdMigration attempts to log several (key, value) pairs from
// the accessd service to allow for manual verification of successful migration. This
// ensures the accessd service is able to understand newly-migrated data.
//
// While not necessary, migrators may choose to have the relevant values
// (attained via the accessd client_api) printed to manually compare with the
// values from the old and new tables.
func manuallyVerifyAccessdMigration() error {
	// NOTE: we hard-code the accessd server location and bypass accessd/client_api.go
	// to avoid using application code as much as possible (registry etc.)
	accessdSrvAddr := "localhost:9091"
	conn, err := grpc.Dial(accessdSrvAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to connect to accessd server")
	}
	client := accessprotos.NewAccessControlManagerClient(conn)
	ctx := context.Background()

	operators, err := client.ListOperators(ctx, &protos.Void{})
	if err != nil {
		return errors.Wrap(err, "failed to list operators")
	}
	glog.Infof("[manually verify] number of operators: %d", len(operators.List))

	if len(operators.List) == 0 {
		glog.Warning("no operators found")
		return nil
	}

	operator := operators.List[0]
	acl, err := client.GetOperatorACL(ctx, operator)
	if err != nil {
		return errors.Wrapf(err, "failed to get operator ACL for operator: %+v", operator)
	}
	glog.Infof("[manually verify] operator-acl pair: {operator: %+v, acl: %+v}", operator, acl)

	return nil
}
