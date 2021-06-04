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

package spec_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger/spec"

	"github.com/stretchr/testify/assert"
)

func TestSpecs_GetCommonSpec(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	want := []byte("common yaml content")
	commonSpecFilePath := filepath.Join(d, "common/swagger-common.yml")
	assert.NoError(t, os.MkdirAll(filepath.Dir(commonSpecFilePath), 0777))
	assert.NoError(t, ioutil.WriteFile(commonSpecFilePath, want, 0644))

	got, err := spec.NewFSLoader(d).GetCommonSpec()
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestSpecs_GetPartialSpec(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	want := []byte("partial yaml content")
	commonSpecFilePath := filepath.Join(d, "partial/lte.swagger.v1.yml")
	assert.NoError(t, os.MkdirAll(filepath.Dir(commonSpecFilePath), 0777))
	assert.NoError(t, ioutil.WriteFile(commonSpecFilePath, want, 0644))

	got, err := spec.NewFSLoader(d).GetPartialSpec("LTE")
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestSpecs_GetStandaloneSpec(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	want := []byte("standalone yaml content")
	commonSpecFilePath := filepath.Join(d, "standalone/feg.swagger.v1.yml")
	assert.NoError(t, os.MkdirAll(filepath.Dir(commonSpecFilePath), 0777))
	assert.NoError(t, ioutil.WriteFile(commonSpecFilePath, want, 0644))

	got, err := spec.NewFSLoader(d).GetStandaloneSpec("feg")
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
