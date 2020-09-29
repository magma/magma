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

/*
	m013_relay_split updates lte network configs to split federation support
	between gx_gy and hss.

	This migration contains one logical updates
		- Migrate cellular_network_config to the new, relay-split version.
			- relay_enabled field removed
			- hss_relay_enabled field added
			- gx_gy_relay_enabled field added
*/

package main

import (
	"database/sql"
	"encoding/json"
	"flag"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"magma/lte/cloud/go/tools/migrations/m013_relay_split/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"
)

// Duplicated from configurator
const (
	networkConfigTable = "cfg_network_configs"

	nwcIDCol   = "network_id"
	nwcTypeCol = "type"
	nwcValCol  = "value"

	cellularNetworkConfigType = "cellular_network"
)

var (
	dryRun bool
)

func init() {
	flag.BoolVar(&dryRun, "dry", false, "don't commit changes")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
}

func main() {
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	glog.Infof("SQL_DRIVER: %s", dbDriver)
	glog.Infof("DATABASE_SOURCE: %s", dbSource)

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}

	_, err = migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, doMigration)
	if err != nil {
		glog.Fatalf("%+v", err)
	}

	glog.Info("SUCCESS")
	glog.Info("END MIGRATION")
}

func doMigration(tx *sql.Tx) (interface{}, error) {
	sc := squirrel.NewStmtCache(tx)
	defer func() { _ = sc.Clear() }()
	builder := sqorc.GetSqlBuilder().RunWith(sc)

	err := migrateNetworkEpcConfigs(tx, builder)
	if err != nil {
		return nil, err
	}

	if dryRun {
		return nil, errors.New("throwing error instead of committing because -dry was specified")
	}
	return nil, nil
}

// migrateNetworkEpcConfigs removes the relay_enabled field
// and generates the hss_relay_enabled and gx_gy_relay_enabled fields
func migrateNetworkEpcConfigs(tx *sql.Tx, builder squirrel.StatementBuilderType) error {
	// First read in all the old lte_network configs
	b := builder.
		Select(nwcIDCol, nwcValCol).
		From(networkConfigTable).
		Where(squirrel.Eq{nwcTypeCol: cellularNetworkConfigType})
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	rows, err := b.RunWith(tx).Query()
	if err != nil {
		return errors.Wrap(err, "get lte_network configs")
	}
	oldByNid := map[string]types.OldNetworkCellularConfigs{}
	for rows.Next() {
		var nid string
		var confBytes []byte
		err = rows.Scan(&nid, &confBytes)
		if err != nil {
			return errors.Wrap(err, "scan lte_network config")
		}
		old := types.OldNetworkCellularConfigs{}
		err = json.Unmarshal(confBytes, &old)
		if err != nil {
			return errors.Wrap(err, "unmarshal existing lte_network config")
		}
		oldByNid[nid] = old
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "get existing lte_network configs: SQL rows error")
	}

	// Convert network configs to new versions
	updateConfByNid := map[string][]byte{}
	for nid, old := range oldByNid {
		oldEpc := old.Epc
		newConf := types.NetworkCellularConfigs{
			Epc: &types.NetworkEpcConfigs{
				CloudSubscriberdbEnabled: oldEpc.CloudSubscriberdbEnabled,
				DefaultRuleID:            oldEpc.DefaultRuleID,
				GxGyRelayEnabled:         oldEpc.RelayEnabled,
				HssRelayEnabled:          oldEpc.RelayEnabled,
				LteAuthAmf:               oldEpc.LteAuthAmf,
				LteAuthOp:                oldEpc.LteAuthOp,
				Mcc:                      oldEpc.Mcc,
				Mnc:                      oldEpc.Mnc,
				Mobility:                 oldEpc.Mobility,
				NetworkServices:          oldEpc.NetworkServices,
				SubProfiles:              oldEpc.SubProfiles,
				Tac:                      oldEpc.Tac,
			},
			FegNetworkID: old.FegNetworkID,
			Ran:          old.Ran,
		}
		newBytes, err := json.Marshal(newConf)
		if err != nil {
			return errors.Wrap(err, "marshal updated lte_network config")
		}
		updateConfByNid[nid] = newBytes
	}

	// Update lte_network configs
	for nid, updateConf := range updateConfByNid {
		bu := builder.
			Update(networkConfigTable).
			Set(nwcValCol, updateConf).
			Where(squirrel.Eq{nwcIDCol: nid, nwcTypeCol: cellularNetworkConfigType})
		sqlStr, args, _ := bu.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bu.RunWith(tx).Exec()
		if err != nil {
			return errors.Wrap(err, "error updating lte_network config")
		}
	}

	return nil
}
