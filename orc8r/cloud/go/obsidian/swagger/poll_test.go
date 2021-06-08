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
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	spec "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

func Test_GetCombinedSwaggerSpecs(t *testing.T) {
	commonTag := spec.TagDefinition{Name: "Tag Common"}
	commonSpec := spec.Spec{Tags: []spec.TagDefinition{commonTag}}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Success with no registered servicers
	expectedSpec := spec.Spec{
		Tags: []spec.TagDefinition{commonTag},
	}
	expectedYaml := marshalToYAML(t, expectedSpec)

	combined, err := swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)

	// Success with registered servicers
	tags := []spec.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_spec_service1", "test_spec_service2", "test_spec_service3"}

	expectedSpec = spec.Spec{
		Tags: []spec.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}
	expectedYaml = marshalToYAML(t, expectedSpec)

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	for i, s := range services {
		registerServicer(t, s, tags[i], spec.TagDefinition{})
	}

	combined, err = swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)

	// Success even with merge warnings (duplicate tag)
	serviceDuplicate := "test_spec_service_dup"
	registerServicer(t, serviceDuplicate, tags[0], spec.TagDefinition{})

	combined, err = swagger.GetCombinedSpec(yamlCommon)
	assert.NoError(t, err)

	assert.Equal(t, expectedYaml, combined)
}

func Test_GetServiceSpec(t *testing.T) {
	// Fail with empty service
	_, err := swagger.GetServiceSpec("")
	assert.Error(t, err)

	// Fail with invalid service
	_, err = swagger.GetServiceSpec("invalid_test_spec_service")
	assert.Error(t, err)

	// Success with valid service
	tag := spec.TagDefinition{Name: "Tag 1"}
	testService := "test_spec_service"
	expected := spec.Spec{
		Tags: []spec.TagDefinition{tag},
	}
	expectedYaml := marshalToYAML(t, expected)

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, "test_spec_service", spec.TagDefinition{}, tag)

	combined, err := swagger.GetServiceSpec(testService)
	assert.NoError(t, err)
	assert.Equal(t, expectedYaml, combined)
}

func registerServicer(
	t *testing.T,
	service string,
	partialTag spec.TagDefinition,
	standaloneTag spec.TagDefinition,
) {
	labels := map[string]string{
		orc8r.SwaggerSpecLabel: "true",
	}

	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, service, labels, nil)
	partialSpec := spec.Spec{Tags: []spec.TagDefinition{partialTag}}
	standaloneSpec := spec.Spec{Tags: []spec.TagDefinition{standaloneTag}}

	partialYamlSpec := marshalToYAML(t, partialSpec)
	standaloneYamlSpec := marshalToYAML(t, standaloneSpec)
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicer(partialYamlSpec, standaloneYamlSpec))

	go srv.RunTest(lis)
}

// marshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func marshalToYAML(t *testing.T, spec spec.Spec) string {
	data, err := spec.MarshalBinary()
	assert.NoError(t, err)
	return string(data)
}
