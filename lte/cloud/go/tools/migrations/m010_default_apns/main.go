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
	m010_default_apns ensures all subscribers are associated with an APN.

	Creates a default APN, then, for each subscriber without an APN, associates
	the subscriber with the default APN. Is idempotent.

	This migration touches a lot of SQL objects within a serializable
	transaction. If the migration fails due to serialization failure, consider
	retrying.

	The migration executable imports main library code in two ways
		- sqorc	-- statement builders
		- other	-- optional, manual verification
	If these imports break this executable, recourse to the following
		- sqorc	-- fix breakage, or change statement builders to manual SQL
				   string execs
		- other	-- fix breakage, or remove manual verification code

	In a single, serializable transaction, performs roughly the following,
	per network:

	Create default APN if not exist
		- Access point name			-- oai.ipv4
		- AMBR downlink				-- unlimited (200000000)
		- AMBR uplink				-- unlimited (100000000)
		- QoS class ID				-- 9
		- Priority level			-- 15
		- Preemption capability		-- true
		- Preemption vulnerability	-- false

	Associate subscribers without an APN to the default APN
		- Create subscriber->APN assoc to default APN if not exist
		- Merge subscriber and APN graph IDs as necessary

	Validate (optional)
		- Use the -verify flag
		- Make calls to configurator service to ensure all subscribers have
		  an associated APN
		- Print a subset of subscribers for manual inspection
*/

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"sort"

	"magma/lte/cloud/go/tools/migrations/m010_default_apns/types"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/tools/migrations"
	"magma/orc8r/lib/go/registry"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	apnEntType        = "apn"
	subscriberEntType = "subscriber"

	numManualVerificationsToLog = 5
)

// Duplicated from configurator
const (
	networksTable = "cfg_networks"

	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"

	nwIDCol = "id"

	entPkCol   = "pk"
	entNidCol  = "network_id"
	entTypeCol = "type"
	entKeyCol  = "\"key\""
	entGidCol  = "graph_id"
	entConfCol = "config"

	aFrCol = "from_pk"
	aToCol = "to_pk"

	internalNetworkID = "network_magma_internal"
)

var (
	dryRun bool
	verify bool
)

// ent identifies essential entity info.
type ent struct {
	pk      string
	typ     string
	graphID string
}

func init() {
	flag.BoolVar(&dryRun, "dry", false, "don't commit changes")
	flag.BoolVar(&verify, "verify", false, "verify successful migration")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
}

// main runs the default_apns migration.
// Migration without SQL error will result in `SUCCESS` printed to stderr.
//
// Optional `-verify` flag interfaces with configurator post-migration
// to verify successful migration.
func main() {
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}
	builder := sqorc.GetSqlBuilder()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		networks, err := getNetworks(tx, builder)
		if err != nil {
			return nil, err
		}

		for _, network := range networks {
			defaultAPN, err := createDefaultAPN(tx, builder, network)
			if err != nil {
				return nil, err
			}
			err = ensureAllSubscribersHaveAPN(tx, builder, network, defaultAPN)
			if err != nil {
				return nil, err
			}
		}

		if dryRun {
			return nil, errors.New("throwing error instead of committing because -dry was specified")
		}

		return nil, nil
	}
	_, err = migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		glog.Fatal(err)
	}

	if verify {
		err := verifyMigration(db, builder)
		if err != nil {
			glog.Fatalf("configurator verification failed: %s", err)
		}
	}

	glog.Info("SUCCESS")
	glog.Info("END MIGRATION")
}

func getNetworks(tx *sql.Tx, builder sqorc.StatementBuilder) ([]string, error) {
	rows, err := builder.
		Select(nwIDCol).From(networksTable).
		Where(squirrel.NotEq{nwIDCol: internalNetworkID}).
		RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, "getNetworks: select networks")
	}

	var networks []string
	for rows.Next() {
		var n string
		err = rows.Scan(&n)
		if err != nil {
			return nil, errors.Wrap(err, "getNetworks: scan network ID")
		}
		networks = append(networks, n)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "getNetworks: SQL rows error")
	}

	return networks, nil
}

// createDefaultAPN ensures the default APN exists for the network, creating
// it if necessary.
// Returns default APN.
// Inserted values:
//	- pk			-- new UUID
//	- network_id	-- passed network
//	- type			-- "apn"
//	- key			-- "oai.ipv4"
//	- graph_id		-- new UUID (different from pk)
//	- config		-- defaultAPNVal
func createDefaultAPN(tx *sql.Tx, builder sqorc.StatementBuilder, network string) (ent, error) {
	apn, err := getDefaultAPN(tx, builder, network)
	if err != nil {
		return ent{}, err
	}
	// Return early if default APN already exists
	if apn != nil {
		return *apn, nil
	}

	pk, gid := newUUID(), newUUID()
	b := builder.
		Insert(entityTable).
		Columns(entPkCol, entNidCol, entTypeCol, entKeyCol, entGidCol, entConfCol).
		Values(pk, network, apnEntType, types.DefaultAPNName, gid, types.DefaultAPNVal)
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	_, err = b.RunWith(tx).Exec()
	if err != nil {
		return ent{}, errors.Wrap(err, "createDefaultAPNIfNotExist: insert error")
	}

	newAPN := ent{pk: pk, typ: apnEntType, graphID: gid}
	return newAPN, nil
}

func getDefaultAPN(tx *sql.Tx, builder sqorc.StatementBuilder, network string) (*ent, error) {
	rows, err := builder.
		Select(entPkCol, entGidCol).From(entityTable).
		Where(
			squirrel.And{
				squirrel.Eq{entNidCol: network},
				squirrel.Eq{entKeyCol: types.DefaultAPNName},
			},
		).
		RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, "getDefaultAPN: select existing default APN")
	}

	var apns []*ent
	for rows.Next() {
		apn := &ent{}
		err = rows.Scan(&apn.pk, &apn.graphID)
		if err != nil {
			return nil, errors.Wrap(err, "getDefaultAPN: scan APN")
		}
		apns = append(apns, apn)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "getDefaultAPN: SQL rows error")
	}

	if len(apns) > 1 {
		return nil, fmt.Errorf("getDefaultAPN: found more than 1 default APN for network %s %v", network, apns)
	}
	if len(apns) == 0 {
		return nil, nil
	}

	return apns[0], nil
}

func ensureAllSubscribersHaveAPN(tx *sql.Tx, builder sqorc.StatementBuilder, network string, apn ent) error {
	subs, err := getSubscribersMissingAPN(tx, builder, network)
	if err != nil {
		return err
	}
	if len(subs) == 0 {
		return nil
	}

	subPKs := funk.Map(subs, func(s ent) string { return s.pk }).([]string)
	allGIDs := funk.Map(append(subs, apn), func(e ent) string { return e.graphID }).([]string)

	err = mapSubscribersToAPN(tx, builder, apn.pk, subPKs)
	if err != nil {
		return err
	}

	err = mergeGraphs(tx, builder, allGIDs)
	if err != nil {
		return err
	}

	return nil
}

func getSubscribersMissingAPN(tx *sql.Tx, builder sqorc.StatementBuilder, network string) ([]ent, error) {
	// SELECT a.graph_id,a.pk, b.pk,b.type
	// FROM cfg_entities AS a
	// LEFT JOIN cfg_assocs ON a.pk=from_pk
	// LEFT JOIN cfg_entities AS b ON to_pk=b.pk
	// WHERE a.network=network AND a.type='subscriber';
	tblA, tblB := "a", "b"
	nidColA := makeCol(tblA, entNidCol)
	gidColA := makeCol(tblA, entGidCol)
	pkColA, pkColB := makeCol(tblA, entPkCol), makeCol(tblB, entPkCol)
	typeColA, typeColB := makeCol(tblA, entTypeCol), makeCol(tblB, entTypeCol)

	rows, err := builder.
		Select(gidColA, pkColA, pkColB, typeColB).
		From(fmt.Sprintf("%s AS %s", entityTable, tblA)).
		LeftJoin(fmt.Sprintf("%s ON %s=%s", entityAssocTable, pkColA, aFrCol)).
		LeftJoin(fmt.Sprintf("%s AS %s ON %s=%s", entityTable, tblB, aToCol, pkColB)).
		Where(
			squirrel.And{
				squirrel.Eq{nidColA: network},
				squirrel.Eq{typeColA: subscriberEntType},
			},
		).
		RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, "getSubscribersMissingAPN: select subscribers")
	}

	assocs := map[ent][]ent{}
	for rows.Next() {
		a := ent{typ: subscriberEntType}
		pkB, typeB := &sql.NullString{}, &sql.NullString{}
		err = rows.Scan(&a.graphID, &a.pk, pkB, typeB)
		if err != nil {
			return nil, errors.Wrap(err, "getSubscribersMissingAPN: scan subscriber")
		}
		assocs[a] = append(assocs[a], ent{pk: pkB.String, typ: typeB.String})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "getSubscribersMissingAPN: SQL rows error")
	}

	var subsMissingAPN []ent
	for sub, ents := range assocs {
		hasAPN := false
		for _, ent := range ents {
			if ent.typ == apnEntType {
				hasAPN = true
			}
		}
		if !hasAPN {
			subsMissingAPN = append(subsMissingAPN, sub)
		}
	}

	return subsMissingAPN, nil
}

func mapSubscribersToAPN(tx *sql.Tx, builder sqorc.StatementBuilder, apnPK string, subscriberPKs []string) error {
	b := builder.Insert(entityAssocTable).Columns(aFrCol, aToCol)
	for _, subPK := range subscriberPKs {
		b = b.Values(subPK, apnPK)
	}
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	_, err := b.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "mapSubscribersToAPN: insert error")
	}

	return nil
}

// mergeGraphs merges the graphs for all passed graph IDs.
// Assumes gids is non-empty.
func mergeGraphs(tx *sql.Tx, builder sqorc.StatementBuilder, gids []string) error {
	gids = funk.UniqString(gids)
	sort.Strings(gids)
	mergedGID, oldGIDs := gids[0], gids[1:]

	b := builder.Update(entityTable).Set(entGidCol, mergedGID).Where(squirrel.Eq{entGidCol: oldGIDs})
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	_, err := b.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "mergeGraphs: update error")
	}

	return nil
}

// verifyMigration checks invariants reported by the configurator service.
//	- Ensures all subscribers have an APN assoc
//	- Ensure default APN and its subscribers have same graph ID
func verifyMigration(db *sql.DB, builder sqorc.StatementBuilder) error {
	registry.MustPopulateServices()

	serdes := serde.NewRegistry(
		configurator.NewNetworkEntityConfigSerde(apnEntType, &types.ApnConfiguration{}),
	)

	nids, err := configurator.ListNetworkIDs()
	if err != nil {
		return err
	}

	for _, nid := range nids {

		// All subscribers have an APN

		allSubs, _, err := configurator.LoadAllEntitiesOfType(
			nid, subscriberEntType,
			configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
			serdes,
		)
		if err != nil {
			return err
		}

		i := 0
		for _, sub := range allSubs {
			apns := funk.
				Chain(sub.Associations).
				Filter(func(tk storage.TypeAndKey) bool { return tk.Type == apnEntType }).
				Map(func(tk storage.TypeAndKey) string { return tk.Key }).
				Value().([]string)

			if len(apns) == 0 {
				return fmt.Errorf("subscriber %s has no APN assocs", sub.Key)
			}

			if i < numManualVerificationsToLog {
				glog.Infof("Subscriber %+v has APN assocs %+v", sub, apns)
			}
			i += 1
		}

		// Default APN and its subscribers have same graph ID

		defaultAPN, err := configurator.LoadEntity(
			nid, apnEntType, types.DefaultAPNName,
			configurator.EntityLoadCriteria{LoadAssocsToThis: true},
			serdes,
		)
		if err != nil {
			return err
		}
		subsFromAssocs := defaultAPN.ParentAssociations

		subsFromGraphID, err := getSubscribersInDefaultAPNGraph(db, builder, nid)
		if err != nil {
			return err
		}

		// The default APN's subscribers should be a subset of the subscribers
		// in the default APN's full graph
		for _, sub := range subsFromAssocs {
			if !funk.Contains(subsFromGraphID, sub.Key) {
				return fmt.Errorf("network %s has graph ID error: subscriber %v not in default APN's graph", nid, sub)
			}
		}
	}

	return nil
}

// getSubscribersInDefaultAPNGraph lists keys for all subscribers with the same
// graph ID as the default APN.
func getSubscribersInDefaultAPNGraph(db *sql.DB, builder sqorc.StatementBuilder, network string) ([]string, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		defaultAPN, err := getDefaultAPN(tx, builder, network)
		if err != nil {
			return nil, err
		}

		rows, err := builder.
			Select(entKeyCol).From(entityTable).
			Where(
				squirrel.And{
					squirrel.Eq{entTypeCol: subscriberEntType},
					squirrel.Eq{entGidCol: defaultAPN.graphID},
				},
			).
			RunWith(tx).Query()
		if err != nil {
			return nil, errors.Wrap(err, "getSubscribersInDefaultAPNGraph: select subscribers")
		}

		var subs []string
		for rows.Next() {
			key := ""
			err = rows.Scan(&key)
			if err != nil {
				return nil, errors.Wrap(err, "getSubscribersInDefaultAPNGraph: scan subscriber PK")
			}
			subs = append(subs, key)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "getSubscribersInDefaultAPNGraph: SQL rows error")
		}

		return subs, nil
	}
	ret, err := migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		return nil, err
	}
	subs := ret.([]string)

	return subs, nil
}

func newUUID() string {
	return uuid.New().String()
}

func makeCol(table, col string) string {
	return fmt.Sprintf("%s.%s", table, col)
}
