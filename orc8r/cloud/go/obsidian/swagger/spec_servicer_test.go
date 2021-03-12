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
	"magma/orc8r/cloud/go/obsidian/swagger/protos"

	"github.com/stretchr/testify/assert"
)

func TestSpecServicer_NewSpecServicerFromFile(t *testing.T) {
	testFile := "test_spec_servicer.swagger.v1.yml"
	testPartialFileContents := "test partial yaml spec"
	testStandaloneFileContents := "test standalone yaml spec"
	tmpDir := "/etc/magma/swagger/specs"
	partialDir := "/etc/magma/swagger/specs/partial"
	standaloneDir := "/etc/magma/swagger/specs/standalone"

	os.RemoveAll(tmpDir)
	defer os.RemoveAll(tmpDir)

	err := os.MkdirAll(partialDir, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(standaloneDir, os.ModePerm)
	assert.NoError(t, err)

	tmpPartialSpecPath := filepath.Join(partialDir, testFile)
	err = ioutil.WriteFile(tmpPartialSpecPath, []byte(testPartialFileContents), 0644)
	assert.NoError(t, err)

	tmpStandaloneSpecPath := filepath.Join(standaloneDir, testFile)
	err = ioutil.WriteFile(tmpStandaloneSpecPath, []byte(testStandaloneFileContents), 0644)
	assert.NoError(t, err)

	// Success
	servicer := swagger.NewSpecServicerFromFile("test_spec_servicer")

	partialReq := &protos.PartialSpecRequest{}
	partialRes, err := servicer.GetPartialSpec(context.Background(), partialReq)
	assert.NoError(t, err)

	assert.Equal(t, testPartialFileContents, partialRes.SwaggerSpec)

	standaloneReq := &protos.StandaloneSpecRequest{}
	standaloneRes, err := servicer.GetStandaloneSpec(context.Background(), standaloneReq)
	assert.NoError(t, err)

	assert.Equal(t, testStandaloneFileContents, standaloneRes.SwaggerSpec)
}

func TestSpecServicer_GetPartialSpec(t *testing.T) {
	testFileContents := "test partial yaml spec"

	// Success
	servicer := swagger.NewSpecServicer(testFileContents, "")

	req := &protos.PartialSpecRequest{}
	res, err := servicer.GetPartialSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, testFileContents, res.SwaggerSpec)
}

func TestSpecServicer_GetStandaloneSpec(t *testing.T) {
	testFileContents := "test standalone yaml spec"

	// Success
	servicer := swagger.NewSpecServicer("", testFileContents)

	req := &protos.StandaloneSpecRequest{}
	res, err := servicer.GetStandaloneSpec(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, testFileContents, res.SwaggerSpec)
}
