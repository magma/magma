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
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/registry"
)

// ReverseProxyHandler tracks registered paths to their associated proxy
// backends. This is used to dynamically update the proxy middleware based
// off of the service registry
type ReverseProxyHandler struct {
	proxyBackendsByPathPrefix map[string]*reverseProxyBackend
	certifierServiceConfig    *certifier.Config
}

type reverseProxyBackend struct {
	serverUrl *url.URL
	active    bool
}

// NewReverseProxyHandler initializes a ReverseProxyHandler
func NewReverseProxyHandler(config *certifier.Config) *ReverseProxyHandler {
	return &ReverseProxyHandler{
		proxyBackendsByPathPrefix: map[string]*reverseProxyBackend{},
		certifierServiceConfig:    config,
	}
}

// AddReverseProxyPaths adds reverse proxying from the echo server to every
// service that has registered obsidian handlers. The proxying is based off of
// the path prefixes that have been registered in the services' annotations.
func (r *ReverseProxyHandler) AddReverseProxyPaths(server *echo.Echo, pathPrefixesByAddr map[*url.URL][]string) (*echo.Echo, error) {
	// Echo does not enforce a one to one mapping of path to group.
	// To ensure that a path isn't registered to multiple backends,
	// track the paths and associated addresses that have been registered.
	activePrefixes := map[string]bool{}
	for addr, prefixes := range pathPrefixesByAddr {
		for _, prefix := range prefixes {
			activePrefixes[prefix] = true
			backend, exists := r.proxyBackendsByPathPrefix[prefix]
			if exists && backend.serverUrl.String() == addr.String() {
				continue
			} else if exists {
				glog.Errorf("path prefix '%s' was attempted to be added to address %s; prefix already registered to %s",
					prefix,
					addr.String(),
					backend.serverUrl.String(),
				)
				continue
			}
			target := []*middleware.ProxyTarget{
				{
					URL: addr,
				},
			}
			g := server.Group(prefix)
			if r.certifierServiceConfig != nil && r.certifierServiceConfig.UseToken {
				g.Use(access.TokenMiddleware)
			}
			g.Use(r.activeBackendMiddleware)
			g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(target)))
			r.proxyBackendsByPathPrefix[prefix] = &reverseProxyBackend{
				serverUrl: addr,
				active:    true,
			}
		}
	}
	r.updateInactiveBackends(activePrefixes)
	return server, nil
}

func (r *ReverseProxyHandler) activeBackendMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Echo does not allow dynamically removing routes.
		// If the backend is no longer active (i.e. the service no longer
		// exists), return a 404.
		path := c.Path()
		path = strings.TrimSuffix(path, "/*")
		backend, exists := r.proxyBackendsByPathPrefix[path]
		if exists && !backend.active {
			return echo.NewHTTPError(http.StatusNotFound, "service not found")
		}
		return next(c)
	}
}

func (r *ReverseProxyHandler) updateInactiveBackends(activePrefixes map[string]bool) {
	// proxyBackendsByPathPrefix contains the union of previously existing
	// proxy configuration and all new paths that have been added. In order to
	// figure out deactivated paths, we take the set difference of the unioned
	// state and the current snapshot.
	for prefix, backend := range r.proxyBackendsByPathPrefix {
		_, exists := activePrefixes[prefix]
		if !exists {
			backend.active = false
		}
	}
}

func GetEchoServerAddressToPathPrefixes() (map[*url.URL][]string, error) {
	pathPrefixesByAddr := map[*url.URL][]string{}
	services, err := registry.FindServices(orc8r.ObsidianHandlersLabel)
	if err != nil {
		return pathPrefixesByAddr, err
	}
	for _, srv := range services {
		pathPrefixes, err := registry.GetAnnotationList(srv, orc8r.ObsidianHandlersPathPrefixesAnnotation)
		if err != nil {
			return map[*url.URL][]string{}, err
		}
		echoServerAddress, err := getEchoServerAddressForService(srv)
		if err != nil {
			return map[*url.URL][]string{}, err
		}
		pathPrefixesByAddr[echoServerAddress] = pathPrefixes
	}
	return pathPrefixesByAddr, nil
}

func getEchoServerAddressForService(service string) (*url.URL, error) {
	httpServerAddr, err := registry.GetHttpServerAddress(service)
	if err != nil {
		return nil, err
	}
	rawUrl := fmt.Sprintf("http://%s", httpServerAddr)
	return url.Parse(rawUrl)
}
