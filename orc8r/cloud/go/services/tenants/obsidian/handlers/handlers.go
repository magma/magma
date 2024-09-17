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

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"

	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/services/tenants/obsidian/models"
	"magma/orc8r/cloud/go/services/tenants/protos"
	"magma/orc8r/lib/go/merrors"
)

const (
	TenantRootPath  = obsidian.V1Root + "tenants"
	TenantInfoURL   = TenantRootPath + obsidian.UrlSep + ":tenant_id"
	ControlProxyURL = TenantInfoURL + obsidian.UrlSep + "control_proxy"
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
		{
			Path:        ControlProxyURL,
			Methods:     obsidian.GET,
			HandlerFunc: GetControlProxyHandler,
		},
		{
			Path:        ControlProxyURL,
			Methods:     obsidian.PUT,
			HandlerFunc: CreateOrUpdateControlProxyHandler,
		},
	}
}

func GetTenantsHandler(c echo.Context) error {
	tenants, err := tenants.GetAllTenants(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error decoding request: %v", err))
	}
	if tenantInfo.ID == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "must provide tenant ID")
	}

	_, err = tenants.CreateTenant(c.Request().Context(), *tenantInfo.ID, &protos.Tenant{
		Name:     tenantInfo.Name,
		Networks: tenantInfo.Networks,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error creating tenant: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func GetTenantHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}

	tenantInfo, err := tenants.GetTenant(c.Request().Context(), tenantID)
	if err != nil {
		return mapErr(err, fmt.Errorf("tenant %d does not exist", tenantID), err)
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error decoding request: %v", err))
	}

	err = tenants.SetTenant(c.Request().Context(), tenantID, &tenantInfo)
	if err != nil {
		return mapErr(err, fmt.Errorf("tenant %d does not exist", tenantID), fmt.Errorf("error setting tenant info: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteTenantHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}

	err := tenants.DeleteTenant(c.Request().Context(), tenantID)
	if err != nil {
		return mapErr(err, fmt.Errorf("Tenant %d does not exist", tenantID), err)
	}

	return c.NoContent(http.StatusNoContent)
}

func GetControlProxyHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}

	err := validateTenantExists(c, tenantID)
	if err != nil {
		return err
	}

	controlProxy, err := tenants.GetControlProxy(c.Request().Context(), tenantID)
	if err != nil {
		return mapErr(err, fmt.Errorf("control proxy for tenantID %d does not exist", tenantID), err)
	}

	return c.JSON(http.StatusOK, models.ControlProxy{ControlProxy: &controlProxy.ControlProxy})
}

func CreateOrUpdateControlProxyHandler(c echo.Context) error {
	tenantID, terr := obsidian.GetTenantID(c)
	if terr != nil {
		return terr
	}

	err := validateTenantExists(c, tenantID)
	if err != nil {
		return err
	}
	var req = protos.CreateOrUpdateControlProxyRequest{}
	req.Id = tenantID
	data := &models.ControlProxy{}
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error decoding request: %v", err))
	}
	if err := data.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req.ControlProxy = *data.ControlProxy

	err = tenants.CreateOrUpdateControlProxy(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error setting control_proxy contents: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func mapErr(err error, notFoundErr error, nonNilErr error) error {
	switch {
	case err == merrors.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, notFoundErr.Error())
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, nonNilErr.Error())
	}
	return nil
}

func validateTenantExists(c echo.Context, tenantID int64) error {
	_, err := tenants.GetTenant(c.Request().Context(), tenantID)
	if err != nil {
		return mapErr(err, fmt.Errorf("tenant %d does not exist", tenantID), err)
	}
	return nil
}
