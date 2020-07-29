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

package dictionarygen

import "sort"

func sortExternalAttributes(e []*externalAttribute) {
	sort.Stable(sortableExternalAttributes(e))
}

type sortableExternalAttributes []*externalAttribute

func (s sortableExternalAttributes) Len() int           { return len(s) }
func (s sortableExternalAttributes) Less(i, j int) bool { return s[i].Attribute < s[j].Attribute }
func (s sortableExternalAttributes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
