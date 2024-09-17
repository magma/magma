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

package parallel_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/parallel"
)

func TestDo(t *testing.T) {
	t.Run("echo", func(t *testing.T) {
		f := func(in parallel.In) (parallel.Out, error) { return in, nil }
		out, err := parallel.MapString(strings.Split("abcdefghijklmnopqrstuvwxyz", ""), parallel.DefaultNumWorkers, f)
		assert.NoError(t, err)
		assert.Equal(t, strings.Split("abcdefghijklmnopqrstuvwxyz", ""), out)
	})

	t.Run("double", func(t *testing.T) {
		f := func(in parallel.In) (parallel.Out, error) {
			s := in.(string)
			return s + s, nil
		}
		out, err := parallel.MapString(strings.Split("abcdefghijklmn", ""), parallel.DefaultNumWorkers, f)
		assert.NoError(t, err)
		assert.Equal(t, strings.Split("aa bb cc dd ee ff gg hh ii jj kk ll mm nn", " "), out)
	})
}
