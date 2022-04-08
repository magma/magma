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

type FieldMap map[string]*Field

type Field struct {
	BaseType     BaseType
	Nullable     bool
	HasDefault   bool
	DefaultValue interface{}
	Unique       bool
}

func (f *Field) GetValue() interface{} {
	if f.Nullable {
		return f.BaseType.nullableValue()
	}
	return f.BaseType.baseValue()
}
