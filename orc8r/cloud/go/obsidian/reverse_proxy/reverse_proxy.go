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
	"fmt"
	"net/url"
	"strings"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// AddReverseProxyPaths adds reverse proxying from the echo server to every
// service that has registered obsidian handlers. The proxying is based off of
// the path prefixes that have been registered in the services' annotations.
func AddReverseProxyPaths(server *echo.Echo) (*echo.Echo, error) {
	pathPrefixesByAddr, err := getEchoServerAddressToPathPrefixes()
	if err != nil {
		return nil, err
	}
	// Echo does not enforce a one to one mapping of path to group.
	// To ensure that every group registers a unique path, track the
	// paths that have been registered.
	registeredPaths := map[string]string{}
	for addr, prefixes := range pathPrefixesByAddr {
		for _, prefix := range prefixes {
			if registeredAddr, exists := registeredPaths[prefix]; exists {
				glog.Errorf("path prefix '%s' was added to multiple addresses: %s, %s", prefix, registeredAddr, addr.String())
				continue
			}
			target := []*middleware.ProxyTarget{
				{
					URL: addr,
				},
			}
			g := server.Group(prefix)
			g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(target)))
			registeredPaths[prefix] = addr.String()
		}
	}
	return server, nil
}

func getEchoServerAddressToPathPrefixes() (map[*url.URL][]string, error) {
	pathPrefixesByAddr := map[*url.URL][]string{}
	services := registry.FindServices(orc8r.ObsidianHandlersLabel)
	for _, srv := range services {
		prefixAnnotation, err := registry.GetAnnotation(srv, orc8r.ObsidianHandlersPathPrefixesAnnotation)
		if err != nil {
			return map[*url.URL][]string{}, err
		}
		trimmedPrefixAnnotation := strings.Trim(prefixAnnotation, "\n")
		strippedPrefixAnnotation := strings.ReplaceAll(trimmedPrefixAnnotation, " ", "")
		pathPrefixes := strings.Split(strippedPrefixAnnotation, orc8r.AnnotationListSeparator)
		echoServerAddress, err := getEchoServerAddressForService(srv)
		if err != nil {
			return map[*url.URL][]string{}, err
		}
		pathPrefixesByAddr[echoServerAddress] = pathPrefixes
	}
	return pathPrefixesByAddr, nil
}

func getEchoServerAddressForService(service string) (*url.URL, error) {
	echoPort, err := registry.GetEchoServerPort(service)
	if err != nil {
		return nil, err
	}
	serviceAddr, err := registry.GetServiceAddress(service)
	if err != nil || len(serviceAddr) == 0 {
		return nil, err
	}
	splitServiceAddr := strings.Split(serviceAddr, ":")
	rawUrl := fmt.Sprintf("http://%s:%d", splitServiceAddr[0], echoPort)
	return url.Parse(rawUrl)
}
