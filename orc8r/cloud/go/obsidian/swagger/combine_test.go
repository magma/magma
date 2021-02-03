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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func Test_GetCombinedSwaggerSpecs(t *testing.T) {
	tmpDir := "/etc/magma/configs/orc8r/swagger_specs/"
	testdataDir := "testdata"
	commonSpecDir := "/etc/magma/configs/orc8r/swagger_specs/common"
	commonSpecFile := "swagger-common.yml"
	goldenFilePath := filepath.Join(testdataDir, "out.yml.golden")
	testServices := []string{"test1", "test2", "test3"}

	os.RemoveAll(tmpDir)
	defer os.RemoveAll(tmpDir)

	err := os.MkdirAll(commonSpecDir, os.ModePerm)
	assert.NoError(t, err)

	commonSpecPath := filepath.Join(testdataDir, commonSpecFile)
	commonSpecTempPath := filepath.Join(commonSpecDir, commonSpecFile)
	copyFileToTestImage(t, commonSpecPath, commonSpecTempPath)

	for _, service := range testServices {
		file := fmt.Sprintf("%s.swagger.v1.yml", service)
		path := filepath.Join(testdataDir, file)
		tmpSpecPath := filepath.Join(tmpDir, file)

		copyFileToTestImage(t, path, tmpSpecPath)
	}

	StartTestServiceInternal(t, testServices[0], testServices[1], testServices[2])

	combined, err := swagger.GetCombinedSwaggerSpecs()
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(goldenFilePath)
	assert.NoError(t, err)
	expected := string(data)

	assert.Equal(t, expected, combined)
}

func StartTestServiceInternal(
	t *testing.T,
	serviceOne string,
	serviceTwo string,
	serviceThree string,
) {
	labels := map[string]string{}

	labels[orc8r.SpecServicerLabel] = "true"

	srv1, lis1 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service1", labels, nil)
	srv2, lis2 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service2", labels, nil)
	srv3, lis3 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service3", labels, nil)

	protos.RegisterSwaggerSpecServer(srv1.GrpcServer, swagger.NewSpecServicerFromFile(serviceOne))
	protos.RegisterSwaggerSpecServer(srv2.GrpcServer, swagger.NewSpecServicerFromFile(serviceTwo))
	protos.RegisterSwaggerSpecServer(srv3.GrpcServer, swagger.NewSpecServicerFromFile(serviceThree))

	go srv1.RunTest(lis1)
	go srv2.RunTest(lis2)
	go srv3.RunTest(lis3)
}

func copyFileToTestImage(t *testing.T, src string, dst string) {
	data, err := ioutil.ReadFile(src)
	assert.NoError(t, err)
	err = ioutil.WriteFile(dst, data, 0644)
	assert.NoError(t, err)
}
