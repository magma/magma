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
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"

	"github.com/labstack/echo"
)

// GenerateCombinedSpecHandler is a routing handler which creates and serves the
// combined Swagger Spec.
func GenerateCombinedSpecHandler(c echo.Context) error {
	yamlCommon, err := swagger.GetCommonSpec()
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}

	combined, err := swagger.GetCombinedSpec(yamlCommon)
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}

	return c.String(http.StatusOK, combined)
}
