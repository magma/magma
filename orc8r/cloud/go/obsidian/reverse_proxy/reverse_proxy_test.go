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
	"net/http"
	"testing"

	"magma/orc8r/cloud/go/obsidian/access/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
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
	go srv1.RunTest(lis1)
	go srv2.RunTest(lis2)

	e, err := startTestServer()
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
	go srv1.RunTest(lis1)
	go srv2.RunTest(lis2)

	e, err := startTestServer()
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

func startTestServer() (*echo.Echo, error) {
	e := echo.New()
	e, err := AddReverseProxyPaths(e)
	if err != nil {
		return nil, err
	}
	go func() {
		e.Start("")
	}()
	return e, nil
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
