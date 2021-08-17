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
	"errors"
	"net/http"
	"os"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/labstack/echo"
)

func getVersionHandler(c echo.Context) error {
	version, ok := os.LookupEnv("VERSION_TAG")
	if !ok {
		return obsidian.HttpError(errors.New("Failed to get Orc8r version"), http.StatusInternalServerError)
	}
	chartVersion, ok := os.LookupEnv("HELM_VERSION_TAG")
	if !ok {
		return obsidian.HttpError(errors.New("Failed to get Helm chart version"), http.StatusInternalServerError)
	}
	versionInfo := models.VersionInfo{
		ContainerImageVersion: version,
		HelmChartVersion:      chartVersion,
	}
	return c.JSON(http.StatusOK, versionInfo)
}
