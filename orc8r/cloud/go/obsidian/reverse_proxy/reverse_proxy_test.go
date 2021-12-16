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

package reverse_proxy

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/obsidian/access/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"
)

func init() {
	flag.Set(service.RunEchoServerFlag, "true")
}

func TestReverseProxy(t *testing.T) {
	pathPrefix1 := "/magma/v1/foo"
	pathPrefix2 := "/magma/v1/bar"
	pathPrefix3 := "/magma/v1/foo/:foo_id/baz"

	labels := map[string]string{
		orc8r.ObsidianHandlersLabel: "true",
	}
	annotations1 := map[string]string{
		orc8r.ObsidianHandlersPathPrefixesAnnotation: fmt.Sprintf("%s,%s", pathPrefix1, pathPrefix2),
	}
	annotations2 := map[string]string{
		orc8r.ObsidianHandlersPathPrefixesAnnotation: pathPrefix3,
	}
	srv1, lis1 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service1", labels, annotations1)
	srv1.EchoServer.GET(pathPrefix1, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	srv1.EchoServer.GET(pathPrefix1+"/:foo_id/rue", func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	srv1.EchoServer.GET(pathPrefix2, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	srv2, lis2 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service2", labels, annotations2)
	srv2.EchoServer.GET(pathPrefix3, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	startTestService(t, srv1, lis1)
	startTestService(t, srv2, lis2)

	handler := NewReverseProxyHandler(nil)
	e, err := startTestServer(handler)
	assert.NoError(t, err)

	listener := tests.WaitForTestServer(t, e)
	if listener == nil {
		return // WaitForTestServer should have 'logged' error already
	}
	urlPrefix := "http://" + listener.Addr().String()

	// Ensure both prefixes of service1 work
	s, err := sendRequest("GET", urlPrefix+pathPrefix1)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	s, err = sendRequest("GET", urlPrefix+pathPrefix2)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	// Ensure the most specific path gets used for proxying
	s, err = sendRequest("GET", urlPrefix+"/magma/v1/foo/foo1/baz")
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	s, err = sendRequest("GET", urlPrefix+"/magma/v1/foo/foo1/rue")
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	// Ensure unregistered path is not found
	s, err = sendRequest("GET", urlPrefix+"/magma/v1/nue")
	assert.NoError(t, err)
	assert.Equal(t, 404, s)

	// Test ReverseProxy properly handles dynamic changes to service registry

	// Update existing service annotation with new prefix
	newPrefix := "/magma/v1/newprefix"
	err = addPrefixesToExistingService("test_service2", newPrefix)
	assert.NoError(t, err)
	srv2.EchoServer.GET(newPrefix, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	// Add new service to registry
	pathPrefix4 := "/magma/v1/dynamic"
	annotations3 := map[string]string{
		orc8r.ObsidianHandlersPathPrefixesAnnotation: pathPrefix4,
	}
	srv3, lis3 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_service3", labels, annotations3)
	srv3.EchoServer.GET(pathPrefix4, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})

	pathPrefixesByAddr, err := GetEchoServerAddressToPathPrefixes()
	assert.NoError(t, err)
	_, err = handler.AddReverseProxyPaths(e, pathPrefixesByAddr)
	assert.NoError(t, err)

	startTestService(t, srv3, lis3)

	// Ensure added prefix to test_service2 works
	s, err = sendRequest("GET", urlPrefix+newPrefix)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	// Ensure new test_service3 can be proxied to
	s, err = sendRequest("GET", urlPrefix+pathPrefix4)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	// Remove services from the registry to test inactive backends
	registry.RemoveService("test_service3")
	registry.RemoveService("test_service1")
	pathPrefixesByAddr, err = GetEchoServerAddressToPathPrefixes()
	assert.NoError(t, err)
	_, err = handler.AddReverseProxyPaths(e, pathPrefixesByAddr)
	assert.NoError(t, err)

	// Test path without wildcard '*'
	s, err = sendRequest("GET", urlPrefix+pathPrefix4)
	assert.NoError(t, err)
	assert.Equal(t, 404, s)

	// Test path with wildcard '*'
	assert.NoError(t, err)
	s, err = sendRequest("GET", urlPrefix+pathPrefix1+"/foo1/rue")
	assert.NoError(t, err)
	assert.Equal(t, 404, s)

	// Ensure active backend still works
	s, err = sendRequest("GET", urlPrefix+"/magma/v1/foo/foo1/baz")
	assert.NoError(t, err)
	assert.Equal(t, 200, s)
}

func TestReverseProxyPathCollision(t *testing.T) {
	pathPrefix := "/magma/v1/foo"
	labels := map[string]string{
		orc8r.ObsidianHandlersLabel: "true",
	}
	annotations := map[string]string{
		orc8r.ObsidianHandlersPathPrefixesAnnotation: pathPrefix,
	}
	srv1, lis1 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "mock_server1", labels, annotations)
	srv1.EchoServer.GET(pathPrefix, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	srv2, lis2 := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "mock_server2", labels, annotations)
	srv2.EchoServer.GET(pathPrefix, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	startTestService(t, srv1, lis1)
	startTestService(t, srv2, lis2)

	handler := NewReverseProxyHandler(nil)
	e, err := startTestServer(handler)
	assert.NoError(t, err)

	listener := tests.WaitForTestServer(t, e)
	if listener == nil {
		return // WaitForTestServer should have 'logged' error already
	}
	urlPrefix := "http://" + listener.Addr().String()
	// Ensure the path still works properly as the server that the shared prefix
	// was registered with will be proxied to
	s, err := sendRequest("GET", urlPrefix+pathPrefix)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)
}

func startTestServer(handler *ReverseProxyHandler) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	pathPrefixesByAddr, err := GetEchoServerAddressToPathPrefixes()
	if err != nil {
		return nil, err
	}
	e, err = handler.AddReverseProxyPaths(e, pathPrefixesByAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		e.Start("")
	}()
	return e, nil
}

func startTestService(t *testing.T, srv *service.OrchestratorService, lis net.Listener) {
	go srv.RunTest(lis)
	tests.WaitForTestServer(t, srv.EchoServer)
}

func sendRequest(method string, url string) (int, error) {
	var body io.Reader = nil
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	var client = &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	return response.StatusCode, err
}

func addPrefixesToExistingService(serviceName string, newPrefixes string) error {
	port, err := registry.GetServicePort(serviceName)
	if err != nil {
		return err
	}
	echoPort, err := registry.GetEchoServerPort(serviceName)
	if err != nil {
		return err
	}
	existingPrefixes, err := registry.GetAnnotation(serviceName, orc8r.ObsidianHandlersPathPrefixesAnnotation)
	if err != nil {
		return err
	}
	updatedPrefixes := existingPrefixes + "," + newPrefixes
	obsidianLabels := map[string]string{
		orc8r.ObsidianHandlersLabel: "true",
	}
	newAnnotations := map[string]string{
		orc8r.ObsidianHandlersPathPrefixesAnnotation: updatedPrefixes,
	}
	registry.AddService(registry.ServiceLocation{
		Name:        serviceName,
		Host:        "localhost",
		Port:        port,
		EchoPort:    echoPort,
		Labels:      obsidianLabels,
		Annotations: newAnnotations,
	})
	return nil
}
