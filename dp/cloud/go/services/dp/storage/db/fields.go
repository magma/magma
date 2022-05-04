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

package db

import "magma/orc8r/cloud/go/sqorc"

type Field struct {
	Name         string
	SqlType      sqorc.ColumnType
	Nullable     bool
	HasDefault   bool
	DefaultValue interface{}
	Unique       bool
	Relation     string
}
