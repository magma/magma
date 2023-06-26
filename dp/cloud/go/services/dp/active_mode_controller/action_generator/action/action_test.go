/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action_test

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/action"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestRequest(t *testing.T) {
	m := &stubAmcManager{}
	data := &storage.MutableRequest{
		Request: &storage.DBRequest{
			CbsdId:  db.MakeInt(123),
			Payload: `{"some":"request"}`,
		},
		RequestType: &storage.DBRequestType{
			Name: db.MakeString("some type"),
		},
	}

	a := action.Request{Data: data}
	require.NoError(t, a.Do(nil, m))

	assert.Equal(t, m.action, createRequest)
	assert.Equal(t, m.request, data)
}

func TestUpdate(t *testing.T) {
	m := &stubAmcManager{}
	data := &storage.DBCbsd{
		AvailableFrequencies: []uint32{1, 2, 3, 4},
	}
	mask := db.NewIncludeMask("available_frequencies")

	a := action.Update{Data: data, Mask: mask}
	require.NoError(t, a.Do(nil, m))

	assert.Equal(t, m.action, updateCbsd)
	assert.Equal(t, m.cbsd, data)
	assert.Equal(t, m.mask, mask)
}

func TestDelete(t *testing.T) {
	m := &stubAmcManager{}

	a := action.Delete{Id: 123}
	require.NoError(t, a.Do(nil, m))

	assert.Equal(t, m.action, deleteCbsd)
	assert.Equal(t, &storage.DBCbsd{Id: db.MakeInt(123)}, m.cbsd)
}

type actionType uint8

const (
	createRequest actionType = iota
	deleteCbsd
	updateCbsd
)

type stubAmcManager struct {
	action  actionType
	request *storage.MutableRequest
	cbsd    *storage.DBCbsd
	mask    db.FieldMask
}

func (s *stubAmcManager) GetState(_ squirrel.BaseRunner) ([]*storage.DetailedCbsd, error) {
	return nil, nil
}

func (s *stubAmcManager) CreateRequest(_ squirrel.BaseRunner, request *storage.MutableRequest) error {
	s.action = createRequest
	s.request = request
	return nil
}

func (s *stubAmcManager) DeleteCbsd(_ squirrel.BaseRunner, cbsd *storage.DBCbsd) error {
	s.action = deleteCbsd
	s.cbsd = cbsd
	return nil
}

func (s *stubAmcManager) UpdateCbsd(_ squirrel.BaseRunner, cbsd *storage.DBCbsd, mask db.FieldMask) error {
	s.action = updateCbsd
	s.cbsd = cbsd
	s.mask = mask
	return nil
}
