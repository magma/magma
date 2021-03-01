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

package mswagger_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger/mswagger"
	"magma/orc8r/cloud/go/obsidian/swagger/mswagger/protos"

	"github.com/stretchr/testify/assert"
)

func TestSpecServicer_NewSpecServicerFromFile(t *testing.T) {
	testFile := "test_spec_servicer.swagger.v1.yml"
	testFileContents := "test yaml spec"
	tmpDir := "/etc/magma/swagger/specs"

	os.RemoveAll(tmpDir)
	defer os.RemoveAll(tmpDir)

	err := os.Mkdir(tmpDir, os.ModePerm)
	assert.NoError(t, err)

	tmpSpecPath := filepath.Join(tmpDir, testFile)
	err = ioutil.WriteFile(tmpSpecPath, []byte(testFileContents), 0644)
	assert.NoError(t, err)

	// Success
	servicer := mswagger.NewSpecServicerFromFile("test_spec_servicer")
	assert.NoError(t, err)

	req := &protos.GetSpecRequest{}
	res, err := servicer.GetSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, testFileContents, res.SwaggerSpec)
}

func TestSpecServicer_GetSpec(t *testing.T) {
	testFileContents := "test yaml spec"

	// Success
	servicer := mswagger.NewSpecServicer(testFileContents)

	req := &protos.GetSpecRequest{}
	res, err := servicer.GetSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, testFileContents, res.SwaggerSpec)
}
