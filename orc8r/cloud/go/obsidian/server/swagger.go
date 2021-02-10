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

package server

import (
	"io/ioutil"
	"log"

	"magma/orc8r/cloud/go/obsidian/swagger"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// GenerateSwaggerSpec is a middleware function which creates and writes the
// combined Swagger Spec.
func GenerateSwaggerSpec(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Print("Generating Combined Swagger Spec")
		outFile := "/var/opt/magma/static/apidocs/v1/swagger.yml"

		yamlCommon, err := swagger.GetCommonSpec()
		if err != nil {
			log.Printf("An error occurred while retrieving the Swagger common spec %+v", err)
			return next(c)
		}

		combined, err := swagger.GetCombinedSpec(yamlCommon)
		if err != nil {
			log.Printf("An error occurred while producing the combined Swagger spec %+v", err)
			return next(c)
		}

		err = ioutil.WriteFile(outFile, []byte(combined), 0644)
		if err != nil {
			err = errors.Wrapf(err, "write combined spec to file")
			log.Printf("An error occurred while writing the combined spec to file %+v", err)
		}
		return next(c)
	}
}
