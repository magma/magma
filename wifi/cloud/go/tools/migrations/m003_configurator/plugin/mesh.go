/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"regexp"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"
	"magma/wifi/cloud/go/tools/migrations/m003_configurator/plugin/types"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const meshTableName = "mesh"

func migrateMeshes(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, migratedGatewayMetas map[string]map[string]migration.MigratedGatewayMeta) error {
	for nid, metas := range migratedGatewayMetas {
		if err := migrateMeshesForNetwork(sc, builder, nid, metas); err != nil {
			return errors.Wrapf(err, "failed to migrate meshes for network %s", nid)
		}
	}
	return nil
}

func migrateMeshesForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, migratedGateways map[string]migration.MigratedGatewayMeta) error {
	// load meshes from mesh table
	meshes, err := migration.UnmarshalJSONPBProtosFromDatastore(sc, builder, networkID, meshTableName, &types.MeshNode{})
	if err != nil {
		return errors.WithStack(err)
	}
	if funk.IsEmpty(meshes) {
		return nil
	}

	// load mesh configs from config service
	rows, err := builder.Select(migration.ConfigKeyCol, migration.ConfigValCol).
		From(migration.GetLegacyTableName(networkID, migration.ConfigTable)).
		Where(squirrel.Eq{migration.ConfigTypeCol: "mesh"}).
		RunWith(sc).
		Query()
	if err != nil {
		return errors.Wrap(err, "failed to load mesh configs")
	}
	defer rows.Close()
	newMeshConfigs := map[string][]byte{}
	for rows.Next() {
		var k string
		var v []byte
		err := rows.Scan(&k, &v)
		if err != nil {
			return errors.Wrapf(err, "failed to scan mesh config")
		}

		oldConf := &types.LegacyMeshConfig{}
		err = migration.Unmarshal(v, oldConf)
		if err != nil {
			return errors.Wrapf(err, "could not marshal mesh config %s", k)
		}
		newConf := &types.MeshConfigs{}
		migration.FillIn(oldConf, newConf)
		newConfVal, err := newConf.MarshalBinary()
		if err != nil {
			return errors.Wrapf(err, "could not marshal migrated mesh config %s", k)
		}
		newMeshConfigs[k] = newConfVal
	}

	// parse out the gateway IDs so we know what assocs to create
	// from chat with Ilya, this is the format of soma gateway IDs:
	// e.g. likoni2_id_5ce28cf1a8b6
	// likoni2 = mesh
	// 5ce28cf1a8b6 = device MAC
	gids := []string{}
	gatewayMetasByMesh := map[string][]migration.MigratedGatewayMeta{}
	meshRe := regexp.MustCompile("^([0-9a-zA-Z_-]+)_id_.+$")
	for gwID, gatewayMeta := range migratedGateways {
		match := meshRe.FindStringSubmatch(gwID)
		if len(match) < 2 || match[1] == "" {
			log.Printf("gateway ID %s does not match mesh gateway ID format, skipping", gwID)
			log.Printf("match array: %v", gwID)
			continue
		}
		gatewayMetasByMesh[match[1]] = append(gatewayMetasByMesh[match[1]], gatewayMeta)
		gids = append(gids, gatewayMeta.GraphID)
	}

	meshInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntNameCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	assocInsertBuilder := builder.Insert(migration.EntityAssocTable).
		Columns(migration.AFrCol, migration.AToCol).
		RunWith(sc)
	hasAssocs := false
	meshGids := map[string]string{}
	for meshID, meshMsg := range meshes {
		meshNode := meshMsg.(*types.MeshNode)
		pk, gid := uuid.New().String(), uuid.New().String()
		meshCfg := newMeshConfigs[meshID]

		gids = append(gids, gid)
		meshGids[meshID] = gid

		meshInsertBuilder = meshInsertBuilder.Values(pk, networkID, "mesh", meshID, meshNode.Name, meshCfg, gid)
		for _, gwMeta := range gatewayMetasByMesh[meshID] {
			assocInsertBuilder = assocInsertBuilder.Values(pk, gwMeta.Pk)
			hasAssocs = true
		}
	}

	_, err = meshInsertBuilder.Exec()
	if err != nil {
		return errors.Wrapf(err, "failed to insert new meshes for network %s", networkID)
	}

	if hasAssocs {
		_, err = assocInsertBuilder.Exec()
		if err != nil {
			return errors.Wrapf(err, "failed to create mesh assocs for network %s", networkID)
		}
	}

	if funk.IsEmpty(funk.UniqString(gids)) {
		return nil
	}

	// use a union-find structure to figure out what the final components for
	// the merged graph should be.
	uf := newUnionFind(funk.UniqString(gids))
	for meshID, gwMetas := range gatewayMetasByMesh {
		meshGid := meshGids[meshID]
		for _, gwMeta := range gwMetas {
			uf.union(meshGid, gwMeta.GraphID)
		}
	}
	graphComponents := uf.getComponents()
	for _, component := range graphComponents {
		if len(component) <= 1 {
			continue
		}

		targetGid := component[0]
		_, err := builder.Update(migration.EntityTable).
			Set(migration.EntGidCol, targetGid).
			Where(squirrel.Eq{migration.EntGidCol: component[1:]}).
			RunWith(sc).
			Exec()
		if err != nil {
			return errors.Wrapf(err, "failed to update gateway graph IDs to mesh graph ID %s", targetGid)
		}
	}

	return nil
}
