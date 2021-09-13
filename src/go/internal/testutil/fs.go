// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

// MustTempDir returns a randomly generated temp dir from ioutil.TempDir and
// a cleanup function. If the test failed, the temp dir is not cleaned up to
// aid in debugging. Any encountered err results in a panic.
func MustTempDir() (string, func(*testing.T)) {
	td, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return td, func(t *testing.T) {
		if t.Failed() {
			return
		}
		if err := os.RemoveAll(td); err != nil {
			panic(err)
		}
	}
}
