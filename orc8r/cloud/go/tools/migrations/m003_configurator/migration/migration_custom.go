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

package migration

import (
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func RunCustomPluginMigrations(sc *squirrel.StmtCache, builder sqorc.StatementBuilder) error {
	for _, plug := range allPlugins {
		glog.Infof("Running custom migrations for plugin %T", plug)
		// reload gateway metas every time in case a plugin changes gw meta
		migratedGatewayMetasByNetwork, err := reloadGatewayMetas(sc, builder)
		if err != nil {
			return errors.WithStack(err)
		}

		err = plug.RunCustomMigrations(sc, builder, migratedGatewayMetasByNetwork)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func reloadGatewayMetas(sc *squirrel.StmtCache, builder sqorc.StatementBuilder) (map[string]map[string]MigratedGatewayMeta, error) {
	rows, err := builder.Select(EntNidCol, EntKeyCol, EntPkCol, EntGidCol).
		From(EntityTable).
		Where(squirrel.Eq{EntTypeCol: "magmad_gateway"}).
		RunWith(sc).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to reload gateway meta info")
	}
	defer rows.Close()

	ret := map[string]map[string]MigratedGatewayMeta{}
	for rows.Next() {
		var nid, k, pk, gid string
		err := rows.Scan(&nid, &k, &pk, &gid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan gateway meta info")
		}

		m, found := ret[nid]
		if !found {
			m = map[string]MigratedGatewayMeta{}
			ret[nid] = m
		}
		m[k] = MigratedGatewayMeta{Pk: pk, GraphID: gid}
	}
	return ret, nil
}
