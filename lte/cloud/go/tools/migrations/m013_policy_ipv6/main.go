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
	m013_policy_ipv6 updates the IPv6-related policy rules.

	If the migration fails due to serialization failure, consider retrying.

	This migration moves the old ipv4_src/ipv4_dst into new IP Address struct
	that support both ipv4 and ipv6.
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

	"magma/lte/cloud/go/tools/migrations/m013_policy_ipv6/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"
)

// Duplicated from configurator
const (
	entityTable = "cfg_entities"

	entConfCol = "config"
	entKeyCol  = "\"key\""
	entPkCol   = "pk"
	entTypeCol = "type"
)

const (
	policyEntType = "policy"
)

var (
	dryRun bool
)

type policyEnt struct {
	pk      string
	network string
	gid     string
	config  types.OldPolicyRuleConfig
}

func init() {
	flag.BoolVar(&dryRun, "dry", false, "don't commit changes")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
}

// Migration without SQL error will result in `SUCCESS` printed to stderr.
func main() {
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
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

	err := migratePolicyRules(tx, builder)
	if err != nil {
		return nil, err
	}

	if dryRun {
		return nil, errors.New("throwing error instead of committing because -dry was specified")
	}

	return nil, nil
}

// Changes the old ipv4_src/ipv4_dst strings to ip_src/ip_dst IPAddress structure
func migratePolicyRules(tx *sql.Tx, builder squirrel.StatementBuilderType) error {
	// Get all existing policy ents
	rows, err := builder.
		Select(entPkCol, entKeyCol, entConfCol).
		From(entityTable).
		Where(squirrel.Eq{entTypeCol: policyEntType}).
		RunWith(tx).Query()
	if err != nil {
		return errors.Wrap(err, "get policy ents")
	}
	oldByKey := map[string]policyEnt{}
	for rows.Next() {
		var key string
		old := policyEnt{}
		var confBytes []byte
		err = rows.Scan(&old.pk, &key, &confBytes)
		if err != nil {
			return errors.Wrap(err, "scan policy")
		}
		err = json.Unmarshal(confBytes, &old.config)
		if err != nil {
			return errors.Wrap(err, "unmarshal existing policy rule config")
		}
		oldByKey[key] = old
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "get existing assocs: SQL rows error")
	}

	// Convert policy configs to {new config}
	// Note we are explicitly updating the old config to drop the Qos field,
	// as we use that as the marker of idempotence.
	updateConfByKey := map[string][]byte{}

	for key, old := range oldByKey {
		c := old.config
		newConf := types.PolicyRuleConfig{
			AppName:        c.AppName,
			AppServiceType: c.AppServiceType,
			MonitoringKey:  c.MonitoringKey,
			Priority:       c.Priority,
			RatingGroup:    c.RatingGroup,
			Redirect:       c.Redirect,
			TrackingType:   c.TrackingType,
		}

		for _, flow_desc := range c.FlowList {
			fd := &types.FlowDescription{
				Action: flow_desc.Action,
				Match: &types.FlowMatch{
					Direction: flow_desc.Match.Direction,
					IPProto:   flow_desc.Match.IPProto,
					IPSrc:     get_ip_address(flow_desc.Match.IPV4Src),
					IPDst:     get_ip_address(flow_desc.Match.IPV4Dst),
					TCPDst:    flow_desc.Match.TCPDst,
					TCPSrc:    flow_desc.Match.TCPSrc,
					UDPDst:    flow_desc.Match.UDPDst,
					UDPSrc:    flow_desc.Match.UDPSrc,
				},
			}
			newConf.FlowList = append(newConf.FlowList, fd)
		}

		newBytes, err := json.Marshal(newConf)
		if err != nil {
			return errors.Wrap(err, "marshal updated policy config")
		}
		updateConfByKey[key] = newBytes
	}

	// Update policy configs
	for key := range updateConfByKey {
		old := oldByKey[key]
		updateConf := updateConfByKey[key]
		bu := builder.Update(entityTable).Set(entConfCol, updateConf).Where(squirrel.Eq{entPkCol: old.pk})
		sqlStr, args, _ := bu.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bu.RunWith(tx).Exec()
		if err != nil {
			return errors.Wrap(err, "error updating policy ent")
		}
	}

	return nil
}

func get_ip_address(ip string) *types.IPAddress {
	if ip == "" {
		return nil
	}
	return &types.IPAddress{
		Version: types.IPAddressVersionIPV4,
		Address: ip,
	}
}
