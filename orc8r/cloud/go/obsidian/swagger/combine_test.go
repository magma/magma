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
	"io/ioutil"
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	swagger_lib "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_GetCommonSpec(t *testing.T) {
	specPath := "/etc/magma/configs/orc8r/swagger_specs"
	commonSpecDir := "/etc/magma/configs/orc8r/swagger_specs/common"
	commonSpecFilePath := "/etc/magma/configs/orc8r/swagger_specs/common/swagger-common.yml"

	os.RemoveAll(specPath)
	defer os.RemoveAll(specPath)

	err := os.MkdirAll(commonSpecDir, os.ModePerm)
	assert.NoError(t, err)

	commonSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{{Name: "Tag Common"}},
	}
	yamlCommon := marshalToYAML(t, commonSpec)

	err = ioutil.WriteFile(commonSpecFilePath, []byte(yamlCommon), 0644)
	assert.NoError(t, err)

	actual, err := swagger.GetCommonSpec()
	assert.NoError(t, err)
	assert.Equal(t, yamlCommon, actual)
}

func Test_GetCombinedSwaggerSpecs(t *testing.T) {
	commonTag := swagger_lib.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger_lib.Spec{Tags: []swagger_lib.TagDefinition{commonTag}}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Success with no registered servicers
	expectedSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{commonTag},
	}
	expectedYaml := marshalToYAML(t, expectedSpec)

	combined, err := swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)

	// Success with registered servicers
	tags := []swagger_lib.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_service1", "test_service2", "test_service3"}

	expectedSpec = swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}
	expectedYaml = marshalToYAML(t, expectedSpec)

	setup(t, services, tags)

	combined, err = swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)
}

func setup(t *testing.T, services []string, tags []swagger_lib.TagDefinition) {
	labels := map[string]string{
		orc8r.SpecServicerLabel: "true",
	}

	for i, s := range services {
		srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, s, labels, nil)
		spec := swagger_lib.Spec{Tags: []swagger_lib.TagDefinition{tags[i]}}

		yamlSpec := marshalToYAML(t, spec)
		protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicer(yamlSpec))

		go srv.RunTest(lis)
	}
}

// marshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func marshalToYAML(t *testing.T, spec swagger_lib.Spec) string {
	data, err := yaml.Marshal(&spec)
	assert.NoError(t, err)
	return string(data)
}
