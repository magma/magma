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

package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	dp_service "magma/dp/cloud/go/services/dp"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"
)

const (
	Dp                = "dp"
	DpPath            = obsidian.V1Root + Dp
	ManageNetworkPath = DpPath + obsidian.UrlSep + ":network_id"

	ManageCbsdsPath = ManageNetworkPath + obsidian.UrlSep + "cbsds"
	ManageCbsdPath  = ManageCbsdsPath + obsidian.UrlSep + ":cbsd_id"

	ManageLogsPath = ManageNetworkPath + obsidian.UrlSep + "logs"
)

const (
	baseWrongValMsg = "'%s' is not a proper value for %s"
	responseCode    = "response_code"
	beginTimestamp  = "begin"
	endTimestamp    = "end"
)

func GetHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: ManageCbsdsPath, Methods: obsidian.GET, HandlerFunc: listCbsds},
		{Path: ManageCbsdsPath, Methods: obsidian.POST, HandlerFunc: createCbsd},
		{Path: ManageCbsdPath, Methods: obsidian.GET, HandlerFunc: fetchCbsd},
		{Path: ManageCbsdPath, Methods: obsidian.DELETE, HandlerFunc: deleteCbsd},
		{Path: ManageCbsdPath, Methods: obsidian.PUT, HandlerFunc: updateCbsd},
		{Path: ManageLogsPath, Methods: obsidian.GET, HandlerFunc: listLogs},
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
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing Cbsd ID"))
}

func listLogs(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	client, err := getLogsManagerClient()
	if err != nil {
		return err
	}
	filter, err := getLogsFilter(c)
	if err != nil {
		return err
	}
	pagination, err := GetPagination(c)
	if err != nil {
		return err
	}
	req := protos.ListLogsRequest{
		NetworkId:  networkId,
		Filter:     filter,
		Pagination: pagination,
	}
	ctx := c.Request().Context()
	resp, err := client.ListLogs(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var payload []models.Message
	for _, l := range resp.Logs {
		payload = append(payload, *models.MessageFromBackend(l))
	}
	return c.JSON(http.StatusOK, payload)
}

func listCbsds(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	client, err := getCbsdManagerClient()
	if err != nil {
		return err
	}
	pagination, err := GetPagination(c)
	if err != nil {
		return err
	}
	req := protos.ListCbsdRequest{
		NetworkId:  networkId,
		Pagination: pagination,
	}
	ctx := c.Request().Context()
	cbsds, err := client.ListCbsds(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var payload []models.Cbsd
	for _, cd := range cbsds.Details {
		payload = append(payload, *models.CbsdFromBackend(cd))
	}
	return c.JSON(http.StatusOK, payload)
}

func fetchCbsd(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	client, err := getCbsdManagerClient()
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	req := protos.FetchCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	cbsd, err := client.FetchCbsd(ctx, &req)
	if err != nil {
		return getHttpError(err)
	}
	return c.JSON(http.StatusOK, models.CbsdFromBackend(cbsd.Details))
}

func createCbsd(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	payload := &models.MutableCbsd{}
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	client, err := getCbsdManagerClient()
	if err != nil {
		return err
	}
	data := models.CbsdToBackend(payload)
	req := protos.CreateCbsdRequest{NetworkId: networkId, Data: data}
	ctx := c.Request().Context()
	_, err = client.CreateCbsd(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

func deleteCbsd(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	client, err := getCbsdManagerClient()
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	req := protos.DeleteCbsdRequest{NetworkId: networkId, Id: int64(id)}
	ctx := c.Request().Context()
	_, ierr := client.DeleteCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateCbsd(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	payload := &models.MutableCbsd{}
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	cbsdId, nerr := getCbsdId(c)
	if nerr != nil {
		return nerr
	}
	client, err := getCbsdManagerClient()
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(cbsdId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	data := models.CbsdToBackend(payload)
	req := protos.UpdateCbsdRequest{NetworkId: networkId, Id: int64(id), Data: data}
	ctx := c.Request().Context()
	_, ierr := client.UpdateCbsd(ctx, &req)
	if ierr != nil {
		return getHttpError(ierr)
	}
	return c.NoContent(http.StatusNoContent)
}

func getHttpError(err error) error {
	switch s, _ := status.FromError(err); s.Code() {
	case codes.NotFound:
		return echo.NewHTTPError(http.StatusNotFound, err)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
}

func getLogsFilter(c echo.Context) (*protos.LogFilter, error) {
	var p string
	var respCode *wrapperspb.Int64Value
	var err error
	p = c.QueryParam(responseCode)
	if p != "" {
		rc, err := strconv.Atoi(p)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, p, responseCode)
		}
		respCode = wrapperspb.Int64(int64(rc))
	}
	p = c.QueryParam(beginTimestamp)
	beginTS, err := getTimeStamp(time.RFC3339, p, beginTimestamp)
	if err != nil {
		return nil, err
	}
	p = c.QueryParam(endTimestamp)
	endTS, err := getTimeStamp(time.RFC3339, p, endTimestamp)
	if err != nil {
		return nil, err
	}
	return &protos.LogFilter{
		From:                c.QueryParam("from"),
		To:                  c.QueryParam("to"),
		Name:                c.QueryParam("type"),
		SerialNumber:        c.QueryParam("serial_number"),
		FccId:               c.QueryParam("fcc_id"),
		ResponseCode:        respCode,
		BeginTimestampMilli: beginTS,
		EndTimestampMilli:   endTS,
	}, nil
}

func getTimeStamp(dateLayout string, p string, paramName string) (*wrapperspb.Int64Value, error) {
	if p == "" {
		return nil, nil
	}
	ts, err := time.Parse(dateLayout, p)
	if err != nil {
		return nil, newBadRequest(baseWrongValMsg, p, paramName)
	}
	return wrapperspb.Int64(ts.UnixNano() / int64(time.Millisecond)), nil
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

func getLogsManagerClient() (protos.LogFetcherClient, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
	}
	return protos.NewLogFetcherClient(conn), nil
}

func getConn() (*grpc.ClientConn, error) {
	conn, err := registry.GetConnection(dp_service.ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, dp_service.ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return conn, nil
}
