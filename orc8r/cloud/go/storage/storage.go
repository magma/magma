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

// Package storage contains common definitions to be used across service
// storage interfaces
package storage

import (
	"fmt"
	"sort"

	"magma/orc8r/lib/go/definitions"

	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

type IsolationLevel int

// TxOptions specifies options for transactions
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

type TypeAndKey struct {
	Type string
	Key  string
}

func (tk TypeAndKey) String() string {
	return fmt.Sprintf("%s-%s", tk.Type, tk.Key)
}

func (tk TypeAndKey) IsLessThan(tkb TypeAndKey) bool {
	return tk.String() < tkb.String()
}

type TKs []TypeAndKey

// Filter returns the tks which match the passed type.
func (tks TKs) Filter(typ string) TKs {
	var filtered TKs
	for _, tk := range tks {
		if tk.Type == typ {
			filtered = append(filtered, tk)
		}
	}
	return filtered
}

// MultiFilter returns the tks which match any of the passed types.
func (tks TKs) MultiFilter(types ...string) TKs {
	var filtered TKs
	for _, tk := range tks {
		if funk.ContainsString(types, tk.Type) {
			filtered = append(filtered, tk)
		}
	}
	return filtered
}

// GetFirst returns the first TK with the passed type.
// Returns err only on tk not found.
func (tks TKs) GetFirst(typ string) (TypeAndKey, error) {
	for _, tk := range tks {
		if tk.Type == typ {
			return tk, nil
		}
	}
	return TypeAndKey{}, fmt.Errorf("no TK of type %s found in %v", typ, tks)
}

// Keys returns the keys of the TKs.
func (tks TKs) Keys() []string {
	var keys []string
	for _, tk := range tks {
		keys = append(keys, tk.Key)
	}
	return keys
}

func MakeTKs(typ string, keys []string) TKs {
	var tks TKs
	for _, key := range keys {
		tks = append(tks, TypeAndKey{Type: typ, Key: key})
	}
	return tks
}

// Difference returns (A-B, B-A) when called as A.Difference(B).
func (tks TKs) Difference(b TKs) (TKs, TKs) {
	a := tks

	aa := map[TypeAndKey]struct{}{}
	for _, tk := range a {
		aa[tk] = struct{}{}
	}
	bb := map[TypeAndKey]struct{}{}
	for _, tk := range b {
		bb[tk] = struct{}{}
	}

	var diffA, diffB TKs
	for tk := range aa {
		if _, inB := bb[tk]; !inB {
			diffA = append(diffA, tk)
		}
	}
	for tk := range bb {
		if _, inA := aa[tk]; !inA {
			diffB = append(diffB, tk)
		}
	}

	return diffA, diffB
}

func (tks TKs) Sort() {
	sort.Slice(tks, func(i, j int) bool {
		return tks[i].IsLessThan(tks[j])
	})
}

// IDGenerator is an interface which wraps the creation of unique IDs
type IDGenerator interface {
	// New returns a new unique ID
	New() string
}

// UUIDGenerator is an implementation of IDGenerator which uses uuidv4
type UUIDGenerator struct{}

func (*UUIDGenerator) New() string {
	return uuid.New().String()
}

func GetSQLDriver() string {
	return definitions.MustGetEnv("SQL_DRIVER")
}

func GetDatabaseSource() string {
	return definitions.MustGetEnv("DATABASE_SOURCE")
}
