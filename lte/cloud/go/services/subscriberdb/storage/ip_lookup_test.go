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

package storage_test

import (
	"testing"

	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestIPLookup(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewIPLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("empty initially", func(t *testing.T) {
		got, err := s.GetIPs("n0", []string{})
		assert.NoError(t, err)
		assert.Empty(t, got)

		got, err = s.GetIPs("n0", []string{"ip0"})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("basic insert", func(t *testing.T) {
		err := s.SetIPs("n0", []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
			{Ip: "ipB", Imsi: "IMSI1", Apn: "apn1"},
			{Ip: "ipB", Imsi: "IMSI2", Apn: "apn2"},
		})
		assert.NoError(t, err)

		got, err := s.GetIPs("n0", []string{"ipA"})
		assert.NoError(t, err)
		want := []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
		}
		assert.Equal(t, want, got)

		got, err = s.GetIPs("n0", []string{"ipA", "ipB"})
		assert.NoError(t, err)
		want = []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
			{Ip: "ipB", Imsi: "IMSI1", Apn: "apn1"},
			{Ip: "ipB", Imsi: "IMSI2", Apn: "apn2"},
		}
		assert.Equal(t, want, got)
	})

	t.Run("upsert", func(t *testing.T) {
		err := s.SetIPs("n0", []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"}, // same pk
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn1"}, // different pk, non-unique imsi.apn (change APN)
			{Ip: "ipB", Imsi: "IMSI2", Apn: "apn1"}, // different pk, non-unique imsi.apn (change IMSI)
			{Ip: "ipC", Imsi: "IMSI2", Apn: "apn2"}, // different pk, overwrite IP
		})
		assert.NoError(t, err)

		got, err := s.GetIPs("n0", []string{"ipA"})
		assert.NoError(t, err)
		want := []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn1"},
		}
		assert.Equal(t, want, got)

		got, err = s.GetIPs("n0", []string{"ipA", "ipB", "ipC"})
		assert.NoError(t, err)
		want = []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn1"},
			{Ip: "ipB", Imsi: "IMSI2", Apn: "apn1"},
			{Ip: "ipC", Imsi: "IMSI2", Apn: "apn2"},
		}
		assert.Equal(t, want, got)
	})

	t.Run("additional network", func(t *testing.T) {
		err := s.SetIPs("n1", []*protos.IPMapping{
			{Ip: "ipZ", Imsi: "IMSI0", Apn: "apn0"},
		})
		assert.NoError(t, err)

		got, err := s.GetIPs("n0", []string{"ipA", "ipB", "ipC", "ipZ"})
		assert.NoError(t, err)
		want := []*protos.IPMapping{
			// Same as in previous sub-test
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn0"},
			{Ip: "ipA", Imsi: "IMSI1", Apn: "apn1"},
			{Ip: "ipB", Imsi: "IMSI2", Apn: "apn1"},
			{Ip: "ipC", Imsi: "IMSI2", Apn: "apn2"},
		}
		assert.Equal(t, want, got)

		got, err = s.GetIPs("n1", []string{"ipA", "ipB", "ipC", "ipZ"})
		assert.NoError(t, err)
		want = []*protos.IPMapping{
			{Ip: "ipZ", Imsi: "IMSI0", Apn: "apn0"},
		}
		assert.Equal(t, want, got)
	})
}
