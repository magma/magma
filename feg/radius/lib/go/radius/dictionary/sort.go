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

package dictionary

import (
	"sort"
	"strconv"
)

func SortAttributes(attrs []*Attribute) {
	sort.Stable(sortAttributes(attrs))
}

type sortAttributes []*Attribute

func (s sortAttributes) Len() int { return len(s) }

func (s sortAttributes) Less(i, j int) bool {
	iOID, _ := strconv.Atoi(s[i].OID)
	jOID, _ := strconv.Atoi(s[j].OID)
	return iOID < jOID
}

func (s sortAttributes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func SortValues(values []*Value) {
	sort.Stable(sortValues(values))
}

type sortValues []*Value

func (s sortValues) Len() int           { return len(s) }
func (s sortValues) Less(i, j int) bool { return s[i].Number < s[j].Number }
func (s sortValues) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func SortVendors(vendors []*Vendor) {
	sort.Stable(sortVendors(vendors))
}

type sortVendors []*Vendor

func (s sortVendors) Len() int           { return len(s) }
func (s sortVendors) Less(i, j int) bool { return s[i].Number < s[j].Number }
func (s sortVendors) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
