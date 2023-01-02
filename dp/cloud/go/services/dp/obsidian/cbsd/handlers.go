/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cbsd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	dp_service "magma/dp/cloud/go/services/dp"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/lib/go/merrors"
	lib_protos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const (
	Dp                = "dp"
	DpPath            = obsidian.V1Root + Dp
	ManageNetworkPath = DpPath + obsidian.UrlSep + ":network_id"

	ManageCbsdsPath    = ManageNetworkPath + obsidian.UrlSep + "cbsds"
	ManageCbsdPath     = ManageCbsdsPath + obsidian.UrlSep + ":cbsd_id"
	DeregisterCbsdPath = ManageCbsdPath + obsidian.UrlSep + "deregister"
	RelinquishCbsdPath = ManageCbsdPath + obsidian.UrlSep + "relinquish"
)

const baseWrongValMsg = "'%s' is not a proper value for %s"

func GetHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: ManageCbsdsPath, Methods: obsidian.GET, HandlerFunc: withNetworkIdAndCbsdClient(listCbsds)},
		{Path: ManageCbsdsPath, Methods: obsidian.POST, HandlerFunc: withNetworkIdAndCbsdClient(createCbsd)},
		{Path: ManageCbsdPath, Methods: obsidian.GET, HandlerFunc: withNetworkIdAndCbsdClient(fetchCbsd)},
		{Path: ManageCbsdPath, Methods: obsidian.DELETE, HandlerFunc: withNetworkIdAndCbsdClient(deleteCbsd)},
		{Path: ManageCbsdPath, Methods: obsidian.PUT, HandlerFunc: withNetworkIdAndCbsdClient(updateCbsd)},
		{Path: DeregisterCbsdPath, Methods: obsidian.POST, HandlerFunc: withNetworkIdAndCbsdClient(deregisterCbsd)},
		{Path: RelinquishCbsdPath, Methods: obsidian.POST, HandlerFunc: withNetworkIdAndCbsdClient(relinquishCbsd)},
	}
}

func withNetworkIdAndCbsdClient(handler func(c echo.Context, networkId string, client protos.CbsdManagementClient) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkId, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		client, err := getCbsdManagerClient()
		if err != nil {
			return err
		}
		return handler(c, networkId, client)
	}
}

func getCbsdId(c echo.Context) (string, *echo.HTTPError) {
	id := c.Param("cbsd_id")
	if id == "" {
		return id, cbsdIdHTTPError()
	}
	return id, nil
}

func cbsdIdHTTPError() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, "missing Cbsd ID")
}

func listCbsds(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	pagination, err := GetPagination(c)
	if err != nil {
		return err
	}
	filter := GetCbsdFilter(c)
	req := protos.ListCbsdRequest{
		NetworkId:  networkId,
		Pagination: pagination,
		Filter:     filter,
	}
	ctx := c.Request().Context()
	cbsds, err := client.ListCbsds(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	payload := &models.PaginatedCbsds{
		Cbsds:      make([]*models.Cbsd, len(cbsds.Details)),
		TotalCount: cbsds.TotalCount,
	}
	for i, cd := range cbsds.Details {
		payload.Cbsds[i] = models.CbsdFromBackend(cd)
	}
	return c.JSON(http.StatusOK, payload)
}

func fetchCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req := protos.FetchCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	cbsd, err := client.FetchCbsd(ctx, &req)
	if err != nil {
		return getHttpError(err)
	}
	return c.JSON(http.StatusOK, models.CbsdFromBackend(cbsd.Details))
}

func createCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	payload := &models.MutableCbsd{}
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if err := payload.ValidateModel(ctx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data, err := models.CbsdToBackend(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req := protos.CreateCbsdRequest{NetworkId: networkId, Data: data}
	_, err = client.CreateCbsd(ctx, &req)
	if err != nil {
		return getHttpError(err)
	}
	return c.NoContent(http.StatusCreated)
}

func deleteCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	req := protos.DeleteCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	_, ierr := client.DeleteCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	payload := &models.MutableCbsd{}
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	data, err := models.CbsdToBackend(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req := protos.UpdateCbsdRequest{NetworkId: networkId, Id: int64(id), Data: data}
	ctx := c.Request().Context()
	_, ierr := client.UserUpdateCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func deregisterCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	req := protos.DeregisterCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	_, ierr := client.DeregisterCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func relinquishCbsd(c echo.Context, networkId string, client protos.CbsdManagementClient) error {
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	req := protos.RelinquishCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	_, ierr := client.RelinquishCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func GetCbsdFilter(c echo.Context) *protos.CbsdFilter {
	return &protos.CbsdFilter{
		SerialNumber: c.QueryParam("serial_number"),
	}
}

func getHttpError(err error) error {
	switch s, _ := status.FromError(err); s.Code() {
	case codes.NotFound:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case codes.AlreadyExists:
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func GetPagination(c echo.Context) (*protos.Pagination, error) {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")

	pagination := protos.Pagination{}
	if l != "" {
		limit, err := strconv.Atoi(l)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, l, "limit")
		}
		pagination.Limit = wrapperspb.Int64(int64(limit))
	}
	if o != "" {
		offset, err := strconv.Atoi(o)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, o, "offset")
		}
		if pagination.Limit == nil {
			return nil, newBadRequest("offset requires a limit")
		}
		pagination.Offset = wrapperspb.Int64(int64(offset))
	}
	return &pagination, nil
}

func newBadRequest(format string, a ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(format, a...))
}

func getCbsdManagerClient() (protos.CbsdManagementClient, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
	}
	return protos.NewCbsdManagementClient(conn), nil
}

func getConn() (*grpc.ClientConn, error) {
	conn, err := registry.GetConnection(dp_service.ServiceName, lib_protos.ServiceType_SOUTHBOUND)
	if err != nil {
		initErr := merrors.NewInitError(err, dp_service.ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return conn, nil
}
