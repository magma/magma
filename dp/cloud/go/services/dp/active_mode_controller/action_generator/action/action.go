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

package action

import (
	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

type Delete struct {
	Id int64
}

func (d *Delete) Do(runner sq.BaseRunner, manager storage.AmcManager) error {
	cbsd := &storage.DBCbsd{Id: db.MakeInt(d.Id)}
	return manager.DeleteCbsd(runner, cbsd)
}

type Update struct {
	Data *storage.DBCbsd
	Mask db.FieldMask
}

func (u *Update) Do(runner sq.BaseRunner, manager storage.AmcManager) error {
	return manager.UpdateCbsd(runner, u.Data, u.Mask)
}

type Request struct {
	Data *storage.MutableRequest
}

func (r *Request) Do(runner sq.BaseRunner, manager storage.AmcManager) error {
	return manager.CreateRequest(runner, r.Data)
}
