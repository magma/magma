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
	swagger_lib "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// RegisterSpecHandlers registers routes for the Swagger UI.
func RegisterSpecHandlers(e *echo.Echo, yamlCommon string, tmpl *template.Template) {
	e.GET(obsidian.StaticURLPrefix+"/spec/:service", GetGenerateSpecHandler(yamlCommon))

	e.GET(obsidian.StaticURLPrefix+"/ui/:service", GetGenerateSpecUIHandler(tmpl))
	e.GET(obsidian.StaticURLPrefix+"/ui/", GetGenerateSpecUIHandler(tmpl))

	// Redirect requests for apidocs/v1/ to swagger/v1/ui/
	e.GET("/apidocs/v1/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, obsidian.StaticURLPrefix+"/ui/")
	})
}

// GetGenerateCombinedSpecHandler returns a routing handler which creates
// and serves the combined Swagger Spec.
func GetGenerateCombinedSpecHandler(yamlCommon string) echo.HandlerFunc {
	return func(c echo.Context) error {
		combined, err := swagger.GetCombinedSpec(yamlCommon)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.String(http.StatusOK, combined)
	}
}

// GetGenerateSpecHandler returns a routing handler which creates and
// serves the raw YAML spec of a service.
func GetGenerateSpecHandler(yamlCommon string) echo.HandlerFunc {
	return func(c echo.Context) error {
		service := c.Param("service")

		service, err := isValidService(service)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		remoteSpec := swagger.NewRemoteSpec(service)

		yamlSpec, err := remoteSpec.GetSpec()
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		combined, warnings, err := swagger_lib.Combine(yamlCommon, []string{yamlSpec})
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if warnings != nil {
			glog.Infof("Some Swagger spec traits were overwritten or unable to be read: %+v", warnings)
		}

		return c.String(http.StatusOK, combined)
	}
}

// GetGenerateSpecUIHandler returns a routing handler which serves
// the UI of a service.
func GetGenerateSpecUIHandler(tmpl *template.Template) echo.HandlerFunc {
	return func(c echo.Context) error {
		service := c.Param("service")

		service, err := isValidService(service)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, service)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		return c.HTML(http.StatusOK, tpl.String())
	}
}

// isValidService returns if the service is a registered
// service with a Swagger spec.
func isValidService(service string) (string, error) {
	services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
	if err != nil {
		return "", err
	}

	// If no service is provided, default to empty string to serve combined spec
	if service == "" {
		return "", nil
	}

	for _, s := range services {
		if service == s {
			return strings.ToLower(s), nil
		}
	}

	return "", errors.New("Service provided is not registered.")
}
