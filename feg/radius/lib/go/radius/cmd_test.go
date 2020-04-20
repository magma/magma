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

package radius_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func Test_cmd_radiusdictgen(t *testing.T) {
	t.Parallel()

	name := tempFile(t)
	defer os.Remove(name)
	output, err := exec.Command("go", "build", "-o", name, "fbc/lib/go/radius/cmd/radius-dict-gen").CombinedOutput()
	if err != nil {
		t.Fatalf("%s\n", output)
	}
}

func Test_cmd_radserver(t *testing.T) {
	t.Parallel()

	name := tempFile(t)
	defer os.Remove(name)
	output, err := exec.Command("go", "build", "-o", name, "fbc/lib/go/radius/cmd/radserver").CombinedOutput()
	if err != nil {
		t.Fatalf("%s\n", output)
	}
}

func Test_cmd_radtest(t *testing.T) {
	t.Parallel()

	name := tempFile(t)
	defer os.Remove(name)
	output, err := exec.Command("go", "build", "-o", name, "fbc/lib/go/radius/cmd/radtest").CombinedOutput()
	if err != nil {
		t.Fatalf("%s\n", output)
	}
}

func tempFile(t *testing.T) string {
	f, err := ioutil.TempFile(os.TempDir(), "gobuild")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}
