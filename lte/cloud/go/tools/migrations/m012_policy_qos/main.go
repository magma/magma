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
	m012_policy_qos updates QoS-related configurator entities and assocs
	to support proposal #2102.

	If the migration fails due to serialization failure, consider retrying.

	This migration contains two logical updates
		- Swap three assoc directions in configurator
		- Migrate policy.flow_qos to a standalone qos_profile entity

    For the assoc swaps, enforce the following assoc edge directions,
	flipping existing edge directions when necessary
        - base_name -> policy_rule
        - base_name -> subscriber
        - subscriber -> policy_rule

	For the policy migration, convert the existing policy entities'
	config.flow_qos field to a standalone policy_qos_profile.
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

	"magma/lte/cloud/go/tools/migrations/m012_policy_qos/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"
)

// Duplicated from configurator
const (
	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"

	entConfCol = "config"
	entGidCol  = "graph_id"
	entKeyCol  = "\"key\""
	entNidCol  = "network_id"
	entPkCol   = "pk"
	entTypeCol = "type"

	aFrCol = "from_pk"
	aToCol = "to_pk"
)

const (
	basenameEntType   = "base_name"
	policyEntType     = "policy"
	subscriberEntType = "subscriber"
	qosProfileEntType = "policy_qos_profile"
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

// main runs the policy_qos migration.
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

	edges := []migrations.AssocDirection{
		{FromType: subscriberEntType, ToType: basenameEntType}, // subscriber -> base_name
		{FromType: subscriberEntType, ToType: policyEntType},   // subscriber -> policy
		{FromType: basenameEntType, ToType: policyEntType},     // base_name -> policy
	}

	err := migrations.SetAssocDirections(builder, edges)
	if err != nil {
		return nil, err
	}

	err = migratePolicyRules(tx, builder)
	if err != nil {
		return nil, err
	}

	if dryRun {
		return nil, errors.New("throwing error instead of committing because -dry was specified")
	}

	return nil, nil
}

// migratePolicyRules derives and creates a policy_qos_profile entity for each
// existing policy entity.
// Also removes the Qos field from existing policy entities.
func migratePolicyRules(tx *sql.Tx, builder squirrel.StatementBuilderType) error {
	// Get all existing policy ents
	rows, err := builder.
		Select(entPkCol, entKeyCol, entNidCol, entGidCol, entConfCol).
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
		err = rows.Scan(&old.pk, &key, &old.network, &old.gid, &confBytes)
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

	// Convert policy configs to {new config, policy_qos_profile ent}
	// Note we are explicitly updating the old config to drop the Qos field,
	// as we use that as the marker of idempotence.
	updateConfByKey := map[string][]byte{}
	createProfileByKey := map[string][]byte{}
	for key, old := range oldByKey {
		c := old.config
		newConf := types.PolicyRuleConfig{
			AppName:        c.AppName,
			AppServiceType: c.AppServiceType,
			FlowList:       c.FlowList,
			MonitoringKey:  c.MonitoringKey,
			Priority:       c.Priority,
			RatingGroup:    c.RatingGroup,
			Redirect:       c.Redirect,
			TrackingType:   c.TrackingType,
		}
		newBytes, err := json.Marshal(newConf)
		if err != nil {
			return errors.Wrap(err, "marshal updated policy config")
		}
		updateConfByKey[key] = newBytes

		profileConf := types.PolicyQosProfile{
			Arp:     nil, // previously unused
			ClassID: 0,   // previously unused
			Gbr:     nil, // previously unused
			ID:      key, // same key, different type
		}
		if c.Qos != nil {
			profileConf.MaxReqBwDl = c.Qos.MaxReqBwDl
			profileConf.MaxReqBwUl = c.Qos.MaxReqBwUl
		}
		profileBytes, err := json.Marshal(profileConf)
		if err != nil {
			return errors.Wrap(err, "marshal new policy_qos_profile config")
		}
		createProfileByKey[key] = profileBytes
	}

	// Update policy configs, insert new policy_qos_profile ents, and add
	// policy->policy_qos_profile edge
	for key := range updateConfByKey {
		old := oldByKey[key]

		if old.config.Qos == nil {
			continue
		}

		updateConf := updateConfByKey[key]
		bu := builder.Update(entityTable).Set(entConfCol, updateConf).Where(squirrel.Eq{entPkCol: old.pk})
		sqlStr, args, _ := bu.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bu.RunWith(tx).Exec()
		if err != nil {
			return errors.Wrap(err, "error updating policy ent")
		}

		qosProfilePK := migrations.MakePK()
		createConf := createProfileByKey[key]
		bi := builder.
			Insert(entityTable).
			Columns(entPkCol, entNidCol, entTypeCol, entKeyCol, entGidCol, entConfCol).
			Values(qosProfilePK, old.network, qosProfileEntType, key, old.gid, createConf)
		sqlStr, args, _ = bi.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bi.RunWith(tx).Exec()
		if err != nil {
			return errors.Wrap(err, "error inserting policy_qos_profile ent")
		}

		bi = builder.
			Insert(entityAssocTable).
			Columns(aFrCol, aToCol).
			Values(old.pk, qosProfilePK)
		sqlStr, args, _ = bi.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bi.RunWith(tx).Exec()
		if err != nil {
			return errors.Wrap(err, "error inserting policy->policy_qos_profile assoc")
		}
	}

	return nil
}
