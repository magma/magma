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
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	models2 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/obsidian"
	magmadModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/protos"
)

const (
	CommandRootV1           = ManageGatewayPath + "/command"
	RebootGatewayV1         = CommandRootV1 + "/reboot"
	RestartServicesV1       = CommandRootV1 + "/restart_services"
	GatewayPingV1           = CommandRootV1 + "/ping"
	GatewayGenericCommandV1 = CommandRootV1 + "/generic"
	TailGatewayLogsV1       = CommandRootV1 + "/tail_logs"
)

func rebootGateway(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	err := magmad.GatewayReboot(c.Request().Context(), networkID, gatewayID)
	if err != nil {
		if err == merrors.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func restartServices(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	var services []string
	err := c.Bind(&services)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = magmad.GatewayRestartServices(c.Request().Context(), networkID, gatewayID, services)
	if err != nil {
		if err == merrors.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func gatewayPing(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	pingRequest := magmadModels.PingRequest{}
	err := c.Bind(&pingRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	response, err := magmad.GatewayPing(c.Request().Context(), networkID, gatewayID, pingRequest.Packets, pingRequest.Hosts)
	if err != nil {
		if err == merrors.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var pingResponse magmadModels.PingResponse
	for _, ping := range response.Pings {
		pingResult := &magmadModels.PingResult{
			HostOrIP:           &ping.HostOrIp,
			NumPackets:         &ping.NumPackets,
			PacketsTransmitted: ping.PacketsTransmitted,
			PacketsReceived:    ping.PacketsReceived,
			AvgResponseMs:      ping.AvgResponseMs,
			Error:              ping.Error,
		}
		pingResponse.Pings = append(pingResponse.Pings, pingResult)
	}
	return c.JSON(http.StatusOK, &pingResponse)
}

func gatewayGenericCommand(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	request := magmadModels.GenericCommandParams{}
	err := c.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	params, err := models2.JSONMapToProtobufStruct(request.Params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	genericCommandParams := protos.GenericCommandParams{
		Command: *request.Command,
		Params:  params,
	}

	response, err := magmad.GatewayGenericCommand(c.Request().Context(), networkID, gatewayID, &genericCommandParams)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.InvalidArgument {
			return echo.NewHTTPError(http.StatusNotFound, st.Message())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	resp, err := models2.ProtobufStructToJSONMap(response.Response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	genericCommandResponse := magmadModels.GenericCommandResponse{
		Response: resp,
	}
	return c.JSON(http.StatusOK, genericCommandResponse)
}

func tailGatewayLogs(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	request := magmadModels.TailLogsRequest{}
	err := c.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stream, err := magmad.TailGatewayLogs(c.Request().Context(), networkID, gatewayID, request.Service)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	go func() {
		<-c.Request().Context().Done()
	}()
	// https://echo.labstack.com/cookbook/streaming-response
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	c.Response().Header().Set(echo.HeaderXContentTypeOptions, "nosniff")
	c.Response().WriteHeader(http.StatusOK)
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if _, err := c.Response().Write([]byte(line.Line)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Response().Flush()
	}

	return c.NoContent(http.StatusNoContent)
}
