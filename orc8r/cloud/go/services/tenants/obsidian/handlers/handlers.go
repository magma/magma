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
	"encoding/json"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/services/tenants/obsidian/models"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/labstack/echo"
)

const (
	TenantRootPath = obsidian.V1Root + "tenants"
	TenantInfoURL  = TenantRootPath + obsidian.UrlSep + ":tenant_id"
)

func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{
			Path:        TenantRootPath,
			Methods:     obsidian.GET,
			HandlerFunc: GetTenantsHandler,
		},
		{
			Path:        TenantRootPath,
			Methods:     obsidian.POST,
			HandlerFunc: CreateTenantHandler,
		},
		{
			Path:        TenantInfoURL,
			Methods:     obsidian.GET,
			HandlerFunc: GetTenantHandler,
		},
		{
			Path:        TenantInfoURL,
			Methods:     obsidian.PUT,
			HandlerFunc: SetTenantHandler,
		},
		{
			Path:        TenantInfoURL,
			Methods:     obsidian.DELETE,
			HandlerFunc: DeleteTenantHandler,
		},
	}
}

func GetTenantsHandler(c echo.Context) error {
	tenants, err := tenants.GetAllTenants()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	tenantsAndIDs := make([]models.Tenant, 0)
	for _, tenant := range tenants.Tenants {
		tenantsAndIDs = append(tenantsAndIDs, models.Tenant{
			ID:       &tenant.Id,
			Networks: tenant.Tenant.Networks,
			Name:     tenant.Tenant.Name})
	}
	return c.JSON(http.StatusOK, tenantsAndIDs)
}

func CreateTenantHandler(c echo.Context) error {
	var tenantInfo = models.Tenant{}
	err := json.NewDecoder(c.Request().Body).Decode(&tenantInfo)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding request: %v", err), http.StatusBadRequest)
	}
	if tenantInfo.ID == nil {
		return obsidian.HttpError(fmt.Errorf("Must provide tenant ID"), http.StatusBadRequest)
	}

	_, err = tenants.CreateTenant(*tenantInfo.ID, &protos.Tenant{
		Name:     tenantInfo.Name,
		Networks: tenantInfo.Networks,
	})
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Error creating tenant: %v", err), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetTenantHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}
	tenantInfo, err := tenants.GetTenant(tenantID)
	switch {
	case err == errors.ErrNotFound:
		return obsidian.HttpError(fmt.Errorf("Tenant %d does not exist", tenantID), http.StatusNotFound)
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, models.Tenant{ID: &tenantID, Name: tenantInfo.Name, Networks: tenantInfo.Networks})
}

func SetTenantHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}

	var tenantInfo = protos.Tenant{}
	err := json.NewDecoder(c.Request().Body).Decode(&tenantInfo)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding request: %v", err), http.StatusBadRequest)
	}

	err = tenants.SetTenant(tenantID, tenantInfo)
	switch {
	case err == errors.ErrNotFound:
		return obsidian.HttpError(fmt.Errorf("Tenant %d does not exist", tenantID), http.StatusNotFound)
	case err != nil:
		return obsidian.HttpError(fmt.Errorf("Error setting tenant info: %v", err), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteTenantHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}
	err := tenants.DeleteTenant(tenantID)
	switch {
	case err == errors.ErrNotFound:
		return obsidian.HttpError(fmt.Errorf("Tenant %d does not exist", tenantID), http.StatusNotFound)
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
