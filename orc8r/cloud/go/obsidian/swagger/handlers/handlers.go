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

package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/spec"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
)

// UIInfo contains the templating variables injected into the Swagger
// UI template.
type UIInfo struct {
	// URL of the underlying Swagger spec
	URL string
	// Services list
	Services []string
	// SelectedService in the sidebar
	SelectedService string
}

// RegisterSwaggerHandlers registers routes for Swagger specs and
// the Swagger UI.
func RegisterSwaggerHandlers(e *echo.Echo) error {
	trailSlashMiddleware := middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})

	err := registerSpecHandlers(e, trailSlashMiddleware)
	if err != nil {
		return err
	}

	err = registerUIHandlers(e, trailSlashMiddleware)
	if err != nil {
		return err
	}

	// Redirect requests for apidocs/v1/ to swagger/v1/ui/
	e.GET(obsidian.StaticURLPrefixLegacy+"/v1", nil, trailSlashMiddleware)
	e.GET(obsidian.StaticURLPrefixLegacy+"/v1/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, obsidian.StaticURLPrefix+"/v1/ui/")
	})

	return nil
}

// GetCombinedSpecHandler returns a routing handler which creates and serves
// the combined Swagger Spec.
func GetCombinedSpecHandler(yamlCommon string) echo.HandlerFunc {
	return func(c echo.Context) error {
		combined, err := swagger.GetCombinedSpec(yamlCommon)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.String(http.StatusOK, combined)
	}
}

// GetSpecHandler returns a routing handler which serves a standalone raw YAML
// spec of a particular service.
func GetSpecHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		service, ok, err := getServiceName(c.Param("service"))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if !ok {
			return obsidian.HttpError(errors.New("service not found"), http.StatusNotFound)
		}

		yamlSpec, err := swagger.GetServiceSpec(service)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, yamlSpec)
	}
}

// GetUIHandler returns a routing handler which serves the UI of a service.
func GetUIHandler(tmpl *template.Template) echo.HandlerFunc {
	return func(c echo.Context) error {
		service, ok, err := getServiceName(c.Param("service"))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if !ok {
			return obsidian.HttpError(errors.New("service not found"), http.StatusNotFound)
		}

		services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		sort.Strings(services)

		uiInfo := UIInfo{
			URL:             obsidian.StaticURLPrefix + "/v1/spec/" + service,
			Services:        services,
			SelectedService: service,
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, uiInfo)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

// registerSpecHandlers registers routes for Swagger specs.
func registerSpecHandlers(e *echo.Echo, trailSlashMiddleware echo.MiddlewareFunc) error {
	yamlCommon, err := spec.GetDefaultLoader().GetCommonSpec()
	if err != nil {
		return err
	}

	if obsidian.EnableDynamicSwaggerSpecs {
		e.GET(obsidian.StaticURLPrefix+"/v1/spec/", GetCombinedSpecHandler(yamlCommon))
	} else {
		// If dynamic Swagger spec is off, serve static spec
		route := obsidian.StaticURLPrefix + "/v1/spec/"
		file := obsidian.StaticFolder + obsidian.StaticURLPrefix + "/v1/spec/swagger.yml"
		e.File(route, file)
	}
	e.GET(obsidian.StaticURLPrefix+"/v1/spec", nil, trailSlashMiddleware)

	e.GET(obsidian.StaticURLPrefix+"/v1/spec/:service", GetSpecHandler())

	return nil
}

// registerUIHandlers registers routes for the Swagger UI.
func registerUIHandlers(e *echo.Echo, trailSlashMiddleware echo.MiddlewareFunc) error {
	tmpl, err := template.ParseFiles(obsidian.StaticFolder + "/swagger/v1/ui/index.html")
	if err != nil {
		return errors.Wrap(err, "retrieve Swagger template")
	}

	e.GET(obsidian.StaticURLPrefix+"/v1/ui/", GetUIHandler(tmpl))
	e.GET(obsidian.StaticURLPrefix+"/v1/ui", nil, trailSlashMiddleware)

	e.GET(obsidian.StaticURLPrefix+"/v1/ui/:service", GetUIHandler(tmpl))

	return nil
}

// getServiceName returns if the service is registered with a Swagger spec.
func getServiceName(service string) (string, bool, error) {
	// If no service is provided, default to empty string to serve combined spec
	if service == "" {
		// Explicitly not returning user-generated string as a
		// defense-in-depth measure against XSS attacks
		return "", true, nil
	}

	services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
	if err != nil {
		return "", false, err
	}

	for _, s := range services {
		// Explicitly not returning user-generated string as a
		// defense-in-depth measure against XSS attacks
		if service == s {
			return strings.ToLower(s), true, nil
		}
	}

	return "", false, nil
}
