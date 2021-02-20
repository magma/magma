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
	"text/template"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_lib "magma/orc8r/cloud/go/swagger"

	"github.com/labstack/echo"
)

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
// the UI of a service
func GenerateSpecUIHandler(c echo.Context) error {
	service := c.Param("service")

	tmpl, err := template.ParseFiles("index.html")
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
