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

package reindex

import (
	"database/sql"
	"sort"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/sqorc"
)

type versioner struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewVersioner(db *sql.DB, builder sqorc.StatementBuilder) Versioner {
	return &versioner{db: db, builder: builder}
}

func (v *versioner) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := v.builder.CreateTable(versionTableName).
			IfNotExists().
			Column(idColVersions).Type(sqorc.ColumnTypeText).NotNull().PrimaryKey().EndColumn().
			Column(actualColVersions).Type(sqorc.ColumnTypeInt).Default(0).NotNull().EndColumn().
			Column(desiredColVersions).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize indexer versions table")
	}
	_, err := sqorc.ExecInTx(v.db, &sql.TxOptions{Isolation: sql.LevelRepeatableRead}, nil, txFn)
	return err
}

func (v *versioner) GetIndexerVersions() ([]*indexer.Versions, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return getIndexerVersions(v.builder, tx)
	}
	txRet, err := sqorc.ExecInTx(v.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.([]*indexer.Versions)
	return ret, nil
}

func (v *versioner) SetIndexerActualVersion(indexerID string, version indexer.Version) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return nil, setIndexerActualVersion(v.builder, tx, indexerID, version)
	}
	_, err := sqorc.ExecInTx(v.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	return err
}

func getIndexerVersions(builder sqorc.StatementBuilder, tx *sql.Tx) ([]*indexer.Versions, error) {
	old, err := getTrackedVersions(builder, tx)
	if err != nil {
		return nil, err
	}

	composed, err := getComposedVersions(old)
	if err != nil {
		return nil, err
	}
	if EqualVersions(composed, old) {
		return composed, nil
	}

	// Test hook after first db call so the tx has "officially" started by acquiring some locks
	TestHookGet()

	err = overwriteAllVersions(builder, tx, composed)
	if err != nil {
		return nil, err
	}

	return composed, nil
}

func setIndexerActualVersion(builder sqorc.StatementBuilder, tx *sql.Tx, indexerID string, version indexer.Version) error {
	_, err := builder.Update(versionTableName).
		Set(actualColVersions, version).
		Where(squirrel.Eq{idColVersions: indexerID}).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrapf(err, "update indexer actual version for %s to %d", indexerID, version)
	}
	return nil
}

func getTrackedVersions(builder sqorc.StatementBuilder, tx *sql.Tx) ([]*indexer.Versions, error) {
	var ret []*indexer.Versions

	rows, err := builder.Select(idColVersions, actualColVersions, desiredColVersions).
		From(versionTableName).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "get all indexer versions, select existing versions")
	}

	defer sqorc.CloseRowsLogOnError(rows, "GetAllIndexerVersions")

	var idVal string
	var actualVal, desiredVal int64 // int64 is driver's default int type, though these cols are actually int32 storing a uint32
	for rows.Next() {
		err = rows.Scan(&idVal, &actualVal, &desiredVal)
		if err != nil {
			return ret, errors.Wrap(err, "get all indexer versions, SQL row scan error")
		}
		v, err := newVersions(idVal, actualVal, desiredVal)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}

	err = rows.Err()
	if err != nil {
		return ret, errors.Wrap(err, "get all indexer versions, SQL rows error")
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].IndexerID < ret[j].IndexerID }) // make deterministic
	return ret, nil
}

func overwriteAllVersions(builder sqorc.StatementBuilder, tx *sql.Tx, versions []*indexer.Versions) error {
	_, err := builder.Delete(versionTableName).RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "overwrite all indexer versions, delete existing versions")
	}

	if len(versions) == 0 {
		return nil
	}

	b := builder.Insert(versionTableName).Columns(idColVersions, actualColVersions, desiredColVersions)
	for _, v := range versions {
		b = b.Values(v.IndexerID, v.Actual, v.Desired)
	}
	_, err = b.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrapf(err, "overwrite all indexer desired versions, insert new versions %+v", versions)
	}

	return nil
}

// getComposedVersions writes the composition of tracked (old) and local (new) indexers to store.
// Determining whether an indexer needs to be reindexed depends on three recorded version infos per indexer:
//	- new_desired	-- desired version from indexer registry
//	- old_desired	-- desired version from existing reindex jobs
//	- actual		-- actual version updated upon successful reindex job completion
func getComposedVersions(old []*indexer.Versions) ([]*indexer.Versions, error) {
	newv, err := getIndexerVersionsByID()
	if err != nil {
		return nil, err
	}
	composed := map[string]*indexer.Versions{}

	// Insert all old versions -- old_desired and actual values
	for _, v := range old {
		composed[v.IndexerID] = v
	}

	// Insert all new versions -- new_desired overwrite any existing old_desired
	for id, newDesired := range newv {
		if _, present := composed[id]; present {
			composed[id].Desired = newDesired
		} else {
			composed[id] = &indexer.Versions{IndexerID: id, Actual: 0, Desired: newDesired}
		}
	}

	ret := funk.Map(composed, func(k string, v *indexer.Versions) *indexer.Versions { return v }).([]*indexer.Versions)
	sort.Slice(ret, func(i, j int) bool { return ret[i].IndexerID < ret[j].IndexerID }) // make deterministic
	return ret, nil
}

// getIndexerVersionsByID returns a map of registered indexer IDs to their registered ("desired") versions.
func getIndexerVersionsByID() (map[string]indexer.Version, error) {
	indexers, err := indexer.GetIndexers()
	if err != nil {
		return nil, err
	}
	ret := map[string]indexer.Version{}
	for _, x := range indexers {
		ret[x.GetID()] = x.GetVersion()
	}
	return ret, nil
}

// newVersions returns a new indexer versions view.
// First checks the indexer versions fit in an indexer.Version.
func newVersions(indexerID string, actualVersion, desiredVersion int64) (*indexer.Versions, error) {
	actual, err := indexer.NewIndexerVersion(actualVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "new actual version for indexer %s", indexerID)
	}
	desired, err := indexer.NewIndexerVersion(desiredVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "new desired version for indexer %s", indexerID)
	}
	v := &indexer.Versions{
		IndexerID: indexerID,
		Actual:    actual,
		Desired:   desired,
	}
	return v, nil
}
