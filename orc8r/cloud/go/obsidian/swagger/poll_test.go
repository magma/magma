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
	"magma/orc8r/cloud/go/obsidian/swagger/mswagger"
	"magma/orc8r/cloud/go/obsidian/swagger/mswagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

func Test_GetCommonSpec(t *testing.T) {
	specPath := "/etc/magma/swagger/specs"
	commonSpecDir := "/etc/magma/swagger/specs/common"
	commonSpecFilePath := "/etc/magma/swagger/specs/common/swagger-common.yml"

	os.RemoveAll(specPath)
	defer os.RemoveAll(specPath)

	err := os.MkdirAll(commonSpecDir, os.ModePerm)
	assert.NoError(t, err)

	commonSpec := swagger.Spec{
		Tags: []swagger.TagDefinition{{Name: "Tag Common"}},
	}
	yamlCommon := marshalToYAML(t, commonSpec)

	err = ioutil.WriteFile(commonSpecFilePath, []byte(yamlCommon), 0644)
	assert.NoError(t, err)

	actual, err := swagger.GetCommonSpec()
	assert.NoError(t, err)
	assert.Equal(t, yamlCommon, actual)
}

func Test_GetCombinedSwaggerSpecs(t *testing.T) {
	commonTag := swagger.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger.Spec{Tags: []swagger.TagDefinition{commonTag}}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Success with no registered servicers
	expectedSpec := swagger.Spec{
		Tags: []swagger.TagDefinition{commonTag},
	}
	expectedYaml := marshalToYAML(t, expectedSpec)

	combined, err := swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)

	// Success with registered servicers
	tags := []swagger.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_spec_service1", "test_spec_service2", "test_spec_service3"}

	expectedSpec = swagger.Spec{
		Tags: []swagger.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}
	expectedYaml = marshalToYAML(t, expectedSpec)

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	for i, s := range services {
		registerServicer(t, s, tags[i])
	}

	combined, err = swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)

	// Success even with merge warnings (duplicate tag)
	serviceDuplicate := "test_spec_service_dup"
	registerServicer(t, serviceDuplicate, tags[0])

	combined, err = swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)
}

func Test_GetCombinedSpecFromService(t *testing.T) {
	commonTag := swagger.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger.Spec{Tags: []swagger.TagDefinition{commonTag}}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Fail with empty service
	_, err := swagger.GetServiceSpec(yamlCommon, "")
	assert.Error(t, err)

	// Fail with invalid service
	_, err = swagger.GetServiceSpec(yamlCommon, "invalid_test_spec_service")
	assert.Error(t, err)

	// Success with valid service
	tag := swagger.TagDefinition{Name: "Tag 1"}
	testService := "test_spec_service"
	expected := swagger.Spec{
		Tags: []swagger.TagDefinition{tag, commonTag},
	}
	expectedYaml := marshalToYAML(t, expected)

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, "test_spec_service", tag)

	combined, err := swagger.GetServiceSpec(yamlCommon, testService)
	assert.NoError(t, err)
	assert.Equal(t, expectedYaml, combined)
}

func registerServicer(t *testing.T, service string, tag swagger.TagDefinition) {
	labels := map[string]string{
		orc8r.SwaggerSpecLabel: "true",
	}

	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, service, labels, nil)
	spec := swagger.Spec{Tags: []swagger.TagDefinition{tag}}

	yamlSpec := marshalToYAML(t, spec)
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, mswagger.NewSpecServicer(yamlSpec))

	go srv.RunTest(lis)
}

// marshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func marshalToYAML(t *testing.T, spec swagger.Spec) string {
	data, err := spec.MarshalBinary()
	assert.NoError(t, err)
	return string(data)
}
