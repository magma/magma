package server

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
)

// TestLogger tests the error logging through the echo middleware.
// Run `go test -v ./... --args --logtostderr -v=<verbosity>` to see logs
func TestLogger(t *testing.T) {
	e := echo.New()
	obsidianHandlers := handlers.GetObsidianHandlers()
	// test without configurator service, should throw a 500 status error
	req := httptest.NewRequest(echo.GET, "/magma/v1/networks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.GET).HandlerFunc
	handlerFunc := Logger(listNetworks)
	handlerFunc(c)

	// set up
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)


	// test functional handler -- should not have logs from middleware
	req = httptest.NewRequest(echo.GET, "/magma/v1/networks", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	listNetworks = tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.GET).HandlerFunc
	handlerFunc = Logger(listNetworks)
	handlerFunc(c)

	// test dysfunctional handler -- should have logs
	networkId := "blah"
	req = httptest.NewRequest(echo.GET, fmt.Sprintf("/magma/v1/networks/%s", networkId), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues(networkId)
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.GET).HandlerFunc
	handlerFunc = Logger(getNetwork)
	handlerFunc(c)
}
