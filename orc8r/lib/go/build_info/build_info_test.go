/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// package build_info_test
package build_info_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ldflags = "-ldflags=-X 'magma/orc8r/lib/go/build_info.buildBranch=test_branch'" +
		" -X 'magma/orc8r/lib/go/build_info.buildCommitHash=abcdef123456'" +
		" -X 'magma/orc8r/lib/go/build_info.buildCommitDate=Tue Jun 29 13:26:43 2021 -0700'" +
		" -X 'magma/orc8r/lib/go/build_info.buildTag=test_tag'"
	expected = "\nBuild Info:\n-----------\n\tCommit Branch: test_branch\n\tCommit Tag:    test_tag\n\t" +
		"Commit Hash:   abcdef123456\n\tCommit Date:   Tue Jun 29 13:26:43 2021 -0700\n\tBuild  Date:   UNDEFINED\n"
)

func TestBuildInfoAPI(t *testing.T) {
	cmd := exec.Command("go", "run", ldflags, "magma/orc8r/lib/go/build_info/test_build_info")
	t.Log(cmd.String())
	out, err := cmd.Output()
	assert.NoError(t, err)
	assert.Equal(t, expected, string(out))
}
