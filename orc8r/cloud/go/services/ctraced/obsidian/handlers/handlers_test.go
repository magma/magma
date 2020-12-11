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
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/ctraced/obsidian/handlers"
	traceModels "magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	"testing"

	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCtracedHandlersBasic(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetObsidianHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	listTraces := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/tracing", obsidian.GET).HandlerFunc
	createTrace := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/tracing", obsidian.POST).HandlerFunc
	getTrace := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/tracing/:trace_id", obsidian.GET).HandlerFunc
	updateTrace := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/tracing/:trace_id", obsidian.PUT).HandlerFunc
	deleteTrace := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/tracing/:trace_id", obsidian.DELETE).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/tracing?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listTraces,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*traceModels.CallTrace{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/tracing"
	tc.ExpectedResult = tests.JSONMarshaler(map[string]*traceModels.CallTrace{})
	tests.RunUnitTest(t, e, tc)

	testTraceCfg := &traceModels.CallTraceConfig{
		TraceID:   "CallTrace1",
		GatewayID: "test_gateway_id",
		Timeout:   300,
		TraceType: traceModels.CallTraceConfigTraceTypeGATEWAY,
	}
	testTrace := &traceModels.CallTrace{
		Config: testTraceCfg,
		State: &traceModels.CallTraceState{
			CallTraceAvailable: false,
			CallTraceEnding:    false,
		},
	}

	// Test create call trace
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/tracing",
		Payload:        testTraceCfg,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createTrace,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that policy rule was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/tracing?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listTraces,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*traceModels.CallTrace{
			"CallTrace1": testTrace,
		}),
	}
	tests.RunUnitTest(t, e, tc)

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/tracing/CallTrace1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "trace_id"},
		ParamValues:    []string{"n1", "CallTrace1"},
		Handler:        getTrace,
		ExpectedStatus: 200,
		ExpectedResult: testTrace,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update Rule Using URL based ID
	testMutableTrace := &traceModels.MutableCallTrace{
		RequestedEnd: swag.Bool(true),
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/tracing/CallTrace1",
		Payload:        testMutableTrace,
		ParamNames:     []string{"network_id", "trace_id"},
		ParamValues:    []string{"n1", "CallTrace1"},
		Handler:        updateTrace,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update results
	testTrace.State.CallTraceEnding = true
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/tracing/CallTrace1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "trace_id"},
		ParamValues:    []string{"n1", "CallTrace1"},
		Handler:        getTrace,
		ExpectedStatus: 200,
		ExpectedResult: testTrace,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete a rule
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/tracing/CallTrace1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "trace_id"},
		ParamValues:    []string{"n1", "CallTrace1"},
		Handler:        deleteTrace,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/tracing?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listTraces,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*traceModels.CallTrace{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/tracing"
	tc.ExpectedResult = tests.JSONMarshaler(map[string]*traceModels.CallTrace{})
	tests.RunUnitTest(t, e, tc)
}
