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
	"log"
	"net/http"
	"strings"
	"text/template"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/orc8r"
	swagger_lib "magma/orc8r/cloud/go/swagger"
	"magma/orc8r/lib/go/registry"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// RegisterSpecHandlers registers routes for the Swagger UI.
func RegisterSpecHandlers(e *echo.Echo) error {
	yamlCommon, err := swagger.GetCommonSpec()
	if err != nil {
		return err
	}

	handler := GetGenerateCombinedSpecHandler(yamlCommon)
	e.GET(obsidian.StaticURLPrefix+"/v1/swagger.yml", handler)
	e.GET("/swagger/v1/spec/", handler)

	handler = GetGenerateSpecHandler(yamlCommon)
	e.GET("/swagger/v1/spec/:service", handler)

	e.GET("/swagger/v1/ui/:service", GenerateSpecUIHandler)
	e.GET("/swagger/v1/ui/", GenerateSpecUIHandler)

	// Redirect requests for apidocs/v1/ to swagger/v1/ui/
	// e.GET("/apidocs/v1", )
	return nil
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
		if service == "" {
			err := errors.New("Service provided is not registered.")
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
			log.Printf("Some Swagger spec traits were overwritten or unable to be read: %+v", warnings)
		}

		return c.String(http.StatusOK, combined)
	}
}

// GenerateSpecUIHandler returns a routing handler which serves
// the UI of a service.
func GenerateSpecUIHandler(c echo.Context) error {
	service := c.Param("service")

	service, err := isValidService(service)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if service == "" {
		err := errors.New("Service provided is not registered.")
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	templatePath := "/var/opt/magma/static/apidocs/v1/indexTemplated.html"
	tmpl, err := template.New("indexTemplated.html").Funcs(template.FuncMap{
		"enableDynamicSwaggerSpecs": func() bool {
			return true
		},
	}).ParseFiles(templatePath)
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

// isValidService returns if the service is a registered
// service with a Swagger spec.
func isValidService(service string) (string, error) {
	services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
	if err != nil {
		return "", err
	}

	for _, s := range services {
		if service == s {
			return strings.ToLower(s), nil
		}
	}

	return "", nil
}
