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

package migrations

import (
	"crypto/md5"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Duplicated from configurator
const (
	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"

	entPkCol   = "pk"
	entTypeCol = "type"

	aFrCol = "from_pk"
	aToCol = "to_pk"
)

// AssocDirection is the desired direction of the assoc edge between two types
type AssocDirection struct {
	FromType string
	ToType   string
}

type assocPair struct {
	fromPk string
	toPk   string
}

func MakePK() string {
	return uuid.New().String()
}

// MakeDeterministicPK returns a PK in the expected UUID format, but created
// deterministically from the (network, type, key) primary key.
func MakeDeterministicPK(network, typ, key string) string {
	var in []byte
	in = append(in, []byte(network)...)
	in = append(in, []byte(typ)...)
	in = append(in, []byte(key)...)

	// [16]byte is UUID alias
	// [16]byte is md5.Sum return type
	var id uuid.UUID = md5.Sum(in)
	// Copied from uuid.NewRandom
	id[6] = (id[6] & 0x0f) | 0x40 // Version 4
	id[8] = (id[8] & 0x3f) | 0x80 // Variant is 10
	return id.String()
}

func SetAssocDirections(builder squirrel.StatementBuilderType, edges []AssocDirection) error {
	/*
		SELECT from_pk,to_pk
		FROM cfg_entities AS a
		JOIN cfg_assocs ON a.pk=from_pk
		JOIN cfg_entities AS b ON to_pk=b.pk
		WHERE
			(a.type='edge0.FromType' AND b.type='edge0.ToType') OR
			(a.type='edge1.FromType' AND b.type='edge1.ToType') OR
			...
	*/

	tblA, tblB := "a", "b"
	pkColA, pkColB := makeCol(tblA, entPkCol), makeCol(tblB, entPkCol)
	typeColA, typeColB := makeCol(tblA, entTypeCol), makeCol(tblB, entTypeCol)

	var where squirrel.Or
	for _, assoc := range edges {
		where = append(where, squirrel.Eq{typeColA: assoc.ToType, typeColB: assoc.FromType})
	}

	rows, err := builder.
		Select(aFrCol, aToCol).
		From(fmt.Sprintf("%s AS %s", entityTable, tblA)).
		Join(fmt.Sprintf("%s ON %s=%s", entityAssocTable, pkColA, aFrCol)).
		Join(fmt.Sprintf("%s AS %s ON %s=%s", entityTable, tblB, aToCol, pkColB)).
		Where(where).
		Query()
	if err != nil {
		return errors.Wrap(err, "get existing assocs to flip")
	}

	var assocs []assocPair
	for rows.Next() {
		assoc := assocPair{}
		err = rows.Scan(&assoc.fromPk, &assoc.toPk)
		if err != nil {
			return errors.Wrap(err, "scan assocs")
		}
		assocs = append(assocs, assoc)
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "get existing assocs: SQL rows error")
	}

	glog.Infof("Flipping %d assocs", len(assocs))
	for _, assoc := range assocs {
		err := flipAssocDirection(builder, assoc)
		if err != nil {
			return err
		}
	}

	return nil
}

func flipAssocDirection(builder squirrel.StatementBuilderType, assoc assocPair) error {
	b := builder.
		Update(entityAssocTable).
		Set(aFrCol, assoc.toPk).
		Set(aToCol, assoc.fromPk).
		Where(squirrel.Eq{aFrCol: assoc.fromPk, aToCol: assoc.toPk})
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	_, err := b.Exec()
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	return nil
}

func makeCol(table, col string) string {
	return fmt.Sprintf("%s.%s", table, col)
}
