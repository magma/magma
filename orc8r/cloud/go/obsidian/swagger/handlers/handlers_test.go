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
	swagger_lib "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateCombinedSpecHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	commonTag := swagger_lib.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{commonTag},
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

	tags := []swagger_lib.TagDefinition{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
	}
	services := []string{"test_spec_service1", "test_spec_service2", "test_spec_service3"}
	expectedSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{tags[0], tags[1], tags[2], commonTag},
	}

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	for i, s := range services {
		registerServicer(t, s, tags[i])
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

func Test_GenerateSpecHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	commonTag := swagger_lib.TagDefinition{Name: "Tag Common"}
	commonSpec := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{commonTag},
	}
	yamlCommon := marshalToYAML(t, commonSpec)

	// Fail with invalid service name.
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{"invalid_test_spec_service"},
		Handler:        handlers.GetSpecHandler(yamlCommon),
		ExpectedStatus: 404,
		ExpectedError:  "service not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Success with valid service name.
	tag := swagger_lib.TagDefinition{Name: "Tag 1"}
	expected := swagger_lib.Spec{
		Tags: []swagger_lib.TagDefinition{tag, commonTag},
	}

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, "test_spec_service1", tag)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{"test_spec_service1"},
		Handler:        handlers.GetSpecHandler(yamlCommon),
		ExpectedStatus: 200,
		ExpectedResult: expected,
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GenerateSpecUIHandler(t *testing.T) {
	e := echo.New()
	testURLRoot := "/magma/v1"

	tmplStr := "Test Result is "
	tmpl, err := template.New("test_template.html").Parse(tmplStr + "{{.URL}}")
	assert.NoError(t, err)

	// Fail with invalid service name
	tc := tests.Test{
		Method:                 "GET",
		URL:                    testURLRoot,
		Payload:                nil,
		ParamNames:             []string{"service"},
		ParamValues:            []string{"fake_test_service"},
		Handler:                handlers.GetUIHandler(tmpl, true),
		ExpectedStatus:         404,
		ExpectedErrorSubstring: "service not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Success with empty service name as it should serve the static
	// monolithic Swagger spec
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{""},
		Handler:        handlers.GetUIHandler(tmpl, false),
		ExpectedStatus: 200,
		ExpectedResult: tests.StringMarshaler(tmplStr + "swagger.yml"),
	}
	tests.RunUnitTest(t, e, tc)

	// Success with empty service name as it should serve the dynamic
	// monolithic Swagger spec
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{""},
		Handler:        handlers.GetUIHandler(tmpl, true),
		ExpectedStatus: 200,
		ExpectedResult: tests.StringMarshaler(tmplStr),
	}
	tests.RunUnitTest(t, e, tc)

	// Success with valid service name
	testService2 := "test_spec_service2"

	// Clean up registry
	defer registry.RemoveServicesWithLabel(orc8r.SwaggerSpecLabel)

	registerServicer(t, testService2, swagger_lib.TagDefinition{})

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		ParamNames:     []string{"service"},
		ParamValues:    []string{testService2},
		Handler:        handlers.GetUIHandler(tmpl, true),
		ExpectedStatus: 200,
		ExpectedResult: tests.StringMarshaler(tmplStr + testService2),
	}
	tests.RunUnitTest(t, e, tc)
}

func registerServicer(t *testing.T, service string, tag swagger_lib.TagDefinition) {
	labels := map[string]string{
		orc8r.SwaggerSpecLabel: "true",
	}

	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, service, labels, nil)
	spec := swagger_lib.Spec{Tags: []swagger_lib.TagDefinition{tag}}

	yamlSpec := marshalToYAML(t, spec)
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicer(yamlSpec))

	go srv.RunTest(lis)
}

func marshalToYAML(t *testing.T, spec swagger_lib.Spec) string {
	data, err := spec.MarshalBinary()
	assert.NoError(t, err)
	return string(data)
}
