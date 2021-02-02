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

package swagger_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"

	"github.com/stretchr/testify/assert"
)

var (
	invalidPath      = "invalidPath"
	testFile         = "test.swagger.v1.yml"
	testFileContents = "test yaml spec"
	tmpDir           = "/etc/magma/configs/orc8r/swagger_specs/"
)

func TestSpecServicer_NewSpecServicerFromFile(t *testing.T) {
	req := &swagger_protos.GetSpecRequest{}

	err := os.Mkdir(tmpDir, os.ModePerm)
	assert.NoError(t, err)

	defer os.RemoveAll(tmpDir)

	tmpSpecPath := filepath.Join(tmpDir, testFile)
	err = ioutil.WriteFile(tmpSpecPath, []byte(testFileContents), 0644)
	assert.NoError(t, err)

	// Success
	servicer := swagger.NewSpecServicerFromFile("test")
	assert.NoError(t, err)

	res, err := servicer.GetSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, res.SwaggerSpec, testFileContents)
}

func TestSpecServicer_GetSpec(t *testing.T) {
	req := &swagger_protos.GetSpecRequest{}

	// Success
	servicer := swagger.NewSpecServicer(testFileContents)
	res, err := servicer.GetSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, res.SwaggerSpec, testFileContents)
}
