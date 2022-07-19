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

package containers

type MinQueue []*minQueueItem

func (m *MinQueue) Push(value int) {
	i := len(*m) - 1
	cnt := 1
	for ; i >= 0 && (*m)[i].value >= value; i-- {
		cnt += (*m)[i].count
	}
	item := &minQueueItem{value: value, count: cnt}
	*m = append((*m)[:i+1], item)
}

func (m *MinQueue) Pop() {
	(*m)[0].count--
	if (*m)[0].count == 0 {
		*m = (*m)[1:]
	}
}

func (m MinQueue) Top() int {
	return m[0].value
}

type minQueueItem struct {
	value int
	count int
}
