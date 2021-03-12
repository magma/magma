/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers_test

import (
	"html/template"
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/handlers"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	spec "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_GetCombinedSpecHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	commonTag := spec.TagDefinition{Name: "Tag Common"}
	commonSpec := spec.Spec{
		Tags: []spec.TagDefinition{commonTag},
	}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Success with no registered servicers
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.GetCombinedSpecHandler(yamlCommon),
		ExpectedStatus: 200,
		ExpectedResult: commonSpec,
	}
	tests.RunUnitTest(t, e, tc)

	tags := []spec.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_spec_service1", "test_spec_service2", "test_spec_service3"}
	expectedSpec := spec.Spec{
		Tags: []spec.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	for i, s := range services {
		registerServicer(t, s, tags[i], spec.TagDefinition{})
	}

	// Success with registered servicers
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.GetCombinedSpecHandler(yamlCommon),
		ExpectedStatus: 200,
		ExpectedResult: expectedSpec,
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GetSpecHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	// Fail with invalid service name.
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{"invalid_test_spec_service"},
		Handler:        handlers.GetSpecHandler(),
		ExpectedStatus: 404,
		ExpectedError:  "service not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Success with valid service name.
	tag := spec.TagDefinition{Name: "Tag 1"}
	expected := spec.Spec{
		Tags: []spec.TagDefinition{tag},
	}

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, "test_spec_service1", spec.TagDefinition{}, tag)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{"test_spec_service1"},
		Handler:        handlers.GetSpecHandler(),
		ExpectedStatus: 200,
		ExpectedResult: expected,
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GetUIHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	tmpl, err := template.New("test_template.html").Parse("swagger_spec_url: {{.URL}}")
	assert.NoError(t, err)

	// Fail with invalid service name
	tc := tests.Test{
		Method:                 "GET",
		URL:                    testURLRoot,
		Payload:                nil,
		ParamNames:             []string{"service"},
		ParamValues:            []string{"fake_test_service"},
		Handler:                handlers.GetUIHandler(tmpl),
		ExpectedStatus:         404,
		ExpectedErrorSubstring: "service not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Success with empty service name as it should serve the
	// monolithic Swagger spec
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{""},
		Handler:        handlers.GetUIHandler(tmpl),
		ExpectedStatus: 200,
		ExpectedResult: tests.StringMarshaler("swagger_spec_url: /swagger/v1/spec/"),
	}
	tests.RunUnitTest(t, e, tc)

	// Success with valid service name

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, "test_spec_service2", spec.TagDefinition{}, spec.TagDefinition{})

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{"test_spec_service2"},
		Handler:        handlers.GetUIHandler(tmpl),
		ExpectedStatus: 200,
		ExpectedResult: tests.StringMarshaler("swagger_spec_url: /swagger/v1/spec/test_spec_service2"),
	}
	tests.RunUnitTest(t, e, tc)
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

func marshalToYAML(t *testing.T, spec spec.Spec) string {
	data, err := spec.MarshalBinary()
	assert.NoError(t, err)
	return string(data)
}
