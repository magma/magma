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

type FieldMask interface {
	ShouldInclude(string) bool
}

type includeMask map[string]bool

func NewIncludeMask(columns ...string) includeMask {
	m := make(includeMask, len(columns))
	for _, column := range columns {
		m[column] = true
	}
	return m
}

func (m includeMask) ShouldInclude(column string) bool {
	return m[column]
}

type excludeMask map[string]bool

func NewExcludeMask(columns ...string) excludeMask {
	m := make(excludeMask, len(columns))
	for _, column := range columns {
		m[column] = true
	}
	return m
}

func (m excludeMask) ShouldInclude(column string) bool {
	return !m[column]
}
