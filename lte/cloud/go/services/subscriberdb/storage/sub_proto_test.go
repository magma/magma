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
	"encoding/base64"
	"testing"

	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestSubProtoStore(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewSubProtoStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("initially empty", func(t *testing.T) {
		page, nextToken, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)
	})

	subProtos1 := []*lte_protos.SubscriberData{
		subProtoFromId("IMSI00001"),
		subProtoFromId("IMSI00002"),
		subProtoFromId("IMSI00003"),
	}
	subProtos2 := []*lte_protos.SubscriberData{
		subProtoFromId("IMSI00004"),
		subProtoFromId("IMSI00005"),
		subProtoFromId("IMSI00006"),
	}
	subProtos3 := []*lte_protos.SubscriberData{
		subProtoFromId("IMSI00007"),
		subProtoFromId("IMSI00008"),
		subProtoFromId("IMSI00009"),
	}
	subProtos4 := []*lte_protos.SubscriberData{
		subProtoFromId("IMSI00010"),
	}

	t.Run("insert into tmp table", func(t *testing.T) {
		err = s.InsertManyByNetwork("n0", subProtos1)
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n0", subProtos2)
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n0", subProtos3)
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n0", subProtos4)
		assert.NoError(t, err)

		// The actual sub protos table should still be empty
		page, nextToken, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)
	})

	t.Run("commit update and get by pages", func(t *testing.T) {
		err = s.CommitUpdateByNetwork("n0")
		assert.NoError(t, err)

		page1, nextToken, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assertEqualSubProtos(t, subProtos1, page1)
		expectedNextToken := getTokenByLastIncludedEntity(t, "IMSI00003")
		assert.Equal(t, expectedNextToken, nextToken)

		page2, nextToken, err := s.GetPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assertEqualSubProtos(t, subProtos2, page2)
		expectedNextToken = getTokenByLastIncludedEntity(t, "IMSI00006")
		assert.Equal(t, expectedNextToken, nextToken)

		page3, nextToken, err := s.GetPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assertEqualSubProtos(t, subProtos3, page3)
		expectedNextToken = getTokenByLastIncludedEntity(t, "IMSI00009")
		assert.Equal(t, expectedNextToken, nextToken)

		page4, nextToken, err := s.GetPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assertEqualSubProtos(t, subProtos4, page4)
		expectedNextToken = getTokenByLastIncludedEntity(t, "IMSI00010")
		assert.Equal(t, expectedNextToken, nextToken)

		// The next token should be empty when all pages have been fetched
		finalPage, nextToken, err := s.GetPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assert.Empty(t, nextToken)
		assert.Empty(t, finalPage)
	})

	t.Run("get by ids", func(t *testing.T) {
		// The queried protos are ordered in ascending order by their IDs
		ids := []string{"IMSI00006", "IMSI00003", "IMSI00002", "IMSI00001", "IMSI00000"}
		expectedSubProtos := []*lte_protos.SubscriberData{
			subProtoFromId("IMSI00001"), subProtoFromId("IMSI00002"),
			subProtoFromId("IMSI00003"), subProtoFromId("IMSI00006"),
		}

		subProtos, err := s.GetByIDs("n0", ids)
		assert.NoError(t, err)
		assertEqualSubProtos(t, expectedSubProtos, subProtos)

		// If no matching subscribers exist in store, return an empty array
		ids = []string{"IMSI99991", "IMSI99992", "IMSI99993"}
		subProtos, err = s.GetByIDs("n0", ids)
		assert.NoError(t, err)
		assert.Empty(t, subProtos)
	})

	t.Run("clear tmp table", func(t *testing.T) {
		// After the last CommitUpdate, the tmp table should currently be empty
		err := s.CommitUpdateByNetwork("n0")
		assert.NoError(t, err)
		page, nextToken, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, nextToken)
		assert.Empty(t, page)

		err = s.InsertManyByNetwork("n0", subProtos1)
		assert.NoError(t, err)
		err = s.ClearTmpTable()
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n0", subProtos2)
		assert.NoError(t, err)

		err = s.CommitUpdateByNetwork("n0")
		assert.NoError(t, err)
		// Since the tmp table was cleared halfway through, we'll only commit subProtos2 into the actual table
		page, nextToken, err = s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assertEqualSubProtos(t, subProtos2, page)
		expectedNextToken := getTokenByLastIncludedEntity(t, "IMSI00006")
		assert.Equal(t, expectedNextToken, nextToken)

		page, nextToken, err = s.GetPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.Empty(t, nextToken)
		assert.Empty(t, page)
	})

	t.Run("multiple network insert and get", func(t *testing.T) {
		err = s.ClearTmpTable()
		assert.NoError(t, err)
		err = s.CommitUpdateByNetwork("n0")
		assert.NoError(t, err)

		err = s.InsertManyByNetwork("n0", subProtos1)
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n1", subProtos2)
		assert.NoError(t, err)
		err = s.InsertManyByNetwork("n2", subProtos3)
		assert.NoError(t, err)

		err = s.CommitUpdateByNetwork("n0")
		assert.NoError(t, err)
		err = s.CommitUpdateByNetwork("n1")
		assert.NoError(t, err)
		err = s.CommitUpdateByNetwork("n2")
		assert.NoError(t, err)

		page1, nextToken1, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		expectedNextToken1 := getTokenByLastIncludedEntity(t, "IMSI00003")
		assert.Equal(t, expectedNextToken1, nextToken1)
		assertEqualSubProtos(t, subProtos1, page1)

		page2, nextToken2, err := s.GetPage("n1", "", 3)
		assert.NoError(t, err)
		expectedNextToken2 := getTokenByLastIncludedEntity(t, "IMSI00006")
		assert.Equal(t, expectedNextToken2, nextToken2)
		assertEqualSubProtos(t, subProtos2, page2)

		page3, nextToken3, err := s.GetPage("n2", "", 3)
		assert.NoError(t, err)
		expectedNextToken3 := getTokenByLastIncludedEntity(t, "IMSI00009")
		assert.Equal(t, expectedNextToken3, nextToken3)
		assertEqualSubProtos(t, subProtos3, page3)
	})

	t.Run("delete sub protos", func(t *testing.T) {
		err = s.DeleteSubProtos([]string{"n0", "n1", "n2"})
		assert.NoError(t, err)

		page, nextToken, err := s.GetPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)

		page, nextToken, err = s.GetPage("n1", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)

		page, nextToken, err = s.GetPage("n2", "", 3)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)
	})
}

func assertEqualSubProtos(t *testing.T, expected []*lte_protos.SubscriberData, got []*lte_protos.SubscriberData) {
	assert.Equal(t, len(expected), len(got))
	for i := range expected {
		assert.True(t, proto.Equal(expected[i], got[i]))
	}
}

func subProtoFromId(sid string) *lte_protos.SubscriberData {
	subProto := &lte_protos.SubscriberData{
		Sid: lte_protos.SidFromString(sid),
	}
	return subProto
}

func getTokenByLastIncludedEntity(t *testing.T, sid string) string {
	token := &configurator_storage.EntityPageToken{
		LastIncludedEntity: sid,
	}
	serialized, err := proto.Marshal(token)
	assert.NoError(t, err)

	encoded := base64.StdEncoding.EncodeToString(serialized)
	return encoded
}
