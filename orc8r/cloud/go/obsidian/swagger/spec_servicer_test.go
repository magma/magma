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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/obsidian/swagger/spec"
	"magma/orc8r/cloud/go/obsidian/swagger/spec/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSpecServicer_NewSpecServicerFromFile(t *testing.T) {
	// Missing default dir / spec files should not panic.
	servicer := swagger.NewSpecServicerFromFile("foo")
	assert.NotNil(t, servicer)
}

func TestSpecServicer_NewSpecServicerFromSpecs(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	partialDir := filepath.Join(d, "partial")
	standaloneDir := filepath.Join(d, "standalone")
	err = os.MkdirAll(partialDir, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(standaloneDir, os.ModePerm)
	assert.NoError(t, err)

	testFile := "test_spec_servicer.swagger.v1.yml"
	testPartialFileContents := "test partial yaml spec"
	testStandaloneFileContents := "test standalone yaml spec"

	tmpPartialSpecPath := filepath.Join(partialDir, testFile)
	err = ioutil.WriteFile(tmpPartialSpecPath, []byte(testPartialFileContents), 0644)
	assert.NoError(t, err)

	tmpStandaloneSpecPath := filepath.Join(standaloneDir, testFile)
	err = ioutil.WriteFile(tmpStandaloneSpecPath, []byte(testStandaloneFileContents), 0644)
	assert.NoError(t, err)

	// Success
	servicer := swagger.NewSpecServicerWithLoader(spec.NewFSLoader(d), "test_spec_servicer")

	partialReq := &protos.PartialSpecRequest{}
	partialRes, err := servicer.GetPartialSpec(context.Background(), partialReq)
	assert.NoError(t, err)

	assert.Equal(t, testPartialFileContents, partialRes.SwaggerSpec)

	standaloneReq := &protos.StandaloneSpecRequest{}
	standaloneRes, err := servicer.GetStandaloneSpec(context.Background(), standaloneReq)
	assert.NoError(t, err)

	assert.Equal(t, testStandaloneFileContents, standaloneRes.SwaggerSpec)
}

func TestSpecServicer_NewSpecServicerFromSpecsPartialErr(t *testing.T) {
	service := "a service"

	partialErrMock := &mocks.Loader{}
	partialErrMock.On("GetPartialSpec", service).Return("", errors.New("partial err"))
	standaloneSpec := "standalone spec"
	partialErrMock.On("GetStandaloneSpec", service).Return(standaloneSpec, nil)

	servicer := swagger.NewSpecServicerWithLoader(partialErrMock, service)
	assert.NotNil(t, servicer)

	partialRes, err := servicer.GetPartialSpec(context.Background(), &protos.PartialSpecRequest{})
	assert.NoError(t, err)
	assert.Empty(t, partialRes.SwaggerSpec)

	standaloneRes, err := servicer.GetStandaloneSpec(context.Background(), &protos.StandaloneSpecRequest{})
	assert.NoError(t, err)
	assert.Equal(t, standaloneSpec, standaloneRes.SwaggerSpec)
}

func TestSpecServicer_NewSpecServicerFromSpecsStandaloneErr(t *testing.T) {
	service := "another service"

	standaloneErrMock := &mocks.Loader{}
	partialSpec := "partial spec"
	standaloneErrMock.On("GetPartialSpec", service).Return(partialSpec, errors.New("partial err"))
	standaloneErrMock.On("GetStandaloneSpec", service).Return("", nil)

	servicer := swagger.NewSpecServicerWithLoader(standaloneErrMock, service)
	assert.NotNil(t, servicer)

	partialRes, err := servicer.GetPartialSpec(context.Background(), &protos.PartialSpecRequest{})
	assert.NoError(t, err)
	assert.Equal(t, partialSpec, partialRes.SwaggerSpec)

	standaloneRes, err := servicer.GetStandaloneSpec(context.Background(), &protos.StandaloneSpecRequest{})
	assert.NoError(t, err)
	assert.Empty(t, standaloneRes.SwaggerSpec)
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
