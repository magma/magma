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

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"

	"testing"
)

var (
	invalidPath      = "invalidPath"
	testFile         = "test.swagger.v1.yml"
	testFileContents = "test yaml spec"
)

func TestSpecServicer_NewSpecServicer(t *testing.T) {
	// Setup
	request := &swagger_protos.GetSpecRequest{}

	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	tmpSpecPath := filepath.Join(dir, testFile)
	err = ioutil.WriteFile(tmpSpecPath, []byte(testFileContents), 0644)
	assert.NoError(t, err)

	// Failed Servicer Initialization
	_, err = swagger.NewSpecServicerWithPath(invalidPath)
	assert.Error(t, err)

	// Success
	specServicer, err := swagger.NewSpecServicerWithPath(tmpSpecPath)
	assert.NoError(t, err)

	response, err := specServicer.GetSpec(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, response.SwaggerSpec, testFileContents)
}

func TestSpecServicer_GetSpec(t *testing.T) {
	request := &swagger_protos.GetSpecRequest{}

	// Success
	specServicer := swagger.NewSpecServicer(testFileContents)
	response, err := specServicer.GetSpec(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, response.SwaggerSpec, testFileContents)
}
