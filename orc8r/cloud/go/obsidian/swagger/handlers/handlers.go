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
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
)

// RegisterSwaggerHandlers registers routes for Swagger specs and
// the Swagger UI.
func RegisterSwaggerHandlers(e *echo.Echo) error {
	trailSlashMiddleware := middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})

	err := registerSwaggerSpecHandlers(e, trailSlashMiddleware)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	err = registerSwaggerUIHandlers(e, trailSlashMiddleware)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Redirect requests for apidocs/v1/ to swagger/v1/ui/
	e.GET(obsidian.StaticURLPrefixLegacy+"/v1", nil, trailSlashMiddleware)
	e.GET(obsidian.StaticURLPrefixLegacy+"/v1/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, obsidian.StaticURLPrefix+"/ui/")
	})

	return nil
}

// GetCombinedSpecHandler returns a routing handler which creates
// and serves the combined Swagger Spec.
func GetCombinedSpecHandler(yamlCommon string) echo.HandlerFunc {
	return func(c echo.Context) error {
		combined, err := swagger.GetCombinedSpec(yamlCommon)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.String(http.StatusOK, combined)
	}
}

// GetSpecHandler returns a routing handler which creates and
// serves the raw YAML spec of a singular service.
func GetSpecHandler(yamlCommon string) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestedService := c.Param("service")

		service, ok, err := getServiceName(requestedService)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if !ok {
			return obsidian.HttpError(errors.New("service not found"), http.StatusNotFound)
		}

		combined, err := swagger.GetCombinedSpecFromService(yamlCommon, service)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, combined)
	}
}

// SpecInfo defines the URL to be injected into
// the Swagger UI template.
type SpecInfo struct {
	URL string
}

// GetUIHandler returns a routing handler which serves
// the UI of a service.
func GetUIHandler(tmpl *template.Template, enableDynamicSwaggerSpecs bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestedService := c.Param("service")

		service, ok, err := getServiceName(requestedService)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if !ok {
			return obsidian.HttpError(errors.New("service not found"), http.StatusNotFound)
		}

		// If runtime Swagger spec is off, serve static spec
		if service == "" && !enableDynamicSwaggerSpecs {
			service = "swagger.yml"
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, SpecInfo{URL: service})
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

// registerSwaggerSpecHandlers registers routes for Swagger Specs.
func registerSwaggerSpecHandlers(e *echo.Echo, trailSlashMiddleware echo.MiddlewareFunc) error {
	yamlCommon, err := swagger.GetCommonSpec()
	if err != nil {
		return err
	}

	e.GET(obsidian.StaticURLPrefix+"/spec", nil, trailSlashMiddleware)
	e.GET(obsidian.StaticURLPrefix+"/spec/", GetCombinedSpecHandler(yamlCommon))
	e.GET(obsidian.StaticURLPrefix+"/spec/:service", GetSpecHandler(yamlCommon))

	return nil
}

// registerSwaggerUIHandlers registers routes for the Swagger UI.
func registerSwaggerUIHandlers(e *echo.Echo, trailSlashMiddleware echo.MiddlewareFunc) error {
	tmpl, err := template.ParseFiles(obsidian.StaticFolder + "/swagger/v1/ui/index.html")
	if err != nil {
		return errors.Wrap(err, "retrieve Swagger template")
	}

	e.GET(obsidian.StaticURLPrefix+"/ui", nil, trailSlashMiddleware)
	e.GET(obsidian.StaticURLPrefix+"/ui/", GetUIHandler(tmpl, obsidian.EnableDynamicSwaggerSpecs))
	e.GET(obsidian.StaticURLPrefix+"/ui/:service", GetUIHandler(tmpl, obsidian.EnableDynamicSwaggerSpecs))

	return nil
}

// getServiceName returns if the service is a registered
// service with a Swagger spec.
func getServiceName(service string) (string, bool, error) {
	// If no service is provided, default to empty string to serve combined spec.
	if service == "" {
		// Explicitly not returning user-generated string as a
		// defense-in-depth measure against XSS attacks.
		return "", true, nil
	}

	services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
	if err != nil {
		return "", false, err
	}

	for _, s := range services {
		// Explicitly not returning user-generated string as a
		// defense-in-depth measure against XSS attacks.
		if service == s {
			return strings.ToLower(s), true, nil
		}
	}

	return "", false, nil
}
