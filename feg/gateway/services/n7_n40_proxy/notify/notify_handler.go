/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package notify

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"

	"magma/feg/gateway/sbi"
	n7_client "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/gateway/service_registry"
	"magma/lte/cloud/go/protos"
)

// notify_handler implements the following
// - Creates a HTTP server to receive notifications from PCF
// - Handles SmPolicyUpdateNotification from PCF convert the message to proto and sends it as
//   PolicyReAuth to SessionProxyResponder(feg_relay)
// - Handles SmPolicyTerminateNotification from PCF convert the message to prot and sends it as
//   AbortSession to feg_relay.

const (
	UnableToDeliver  = 3002
	EncodedSessionId = "encodedSessionId"
	// Shutdown timeout is used as a timeout to exit the Echo Server Shutdown
	ShutdownTimeout = 10 * time.Second
)

type NotificationHandler struct {
	config        *n7.N7ClientConfig
	NotifyServer  *sbi.SbiServer
	cloudRegistry service_registry.GatewayRegistry
}

func NewStartedNotificationHandlerWithHandlers(
	config *n7.N7ClientConfig,
	cloudReg service_registry.GatewayRegistry,
) (*NotificationHandler, error) {
	handler := &NotificationHandler{
		config:        config,
		NotifyServer:  sbi.NewSbiServer(config.LocalAddr),
		cloudRegistry: cloudReg,
	}
	err := handler.registerHandlers()
	if err != nil {
		return nil, fmt.Errorf("error registering handlers: %s", err)
	}
	err = handler.NotifyServer.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting notification handler: %s", err)
	}
	return handler, err
}

// RegisterHandlers registers the UpdateNotification and TerminateNotification handlers.
// The notification url is of the form
//		{notifyRoot}/{encodedSessionId}/{operation}
// Example:
//		http://magma-feg.magma.com/sm-policy-control/v1/MTIzNDU2Nzg5MDsxMjM0NQo=/update
// where
// 		notifyRoot = http://magma-feg.magma.com/sm-policy-control/v1
//      encodedSessionId = MTIzNDU2Nzg5MDsxMjM0NQo= (Session-Id is urlencoded)
//      operation = update
// This notification url is send to PCF in the SmPolicyCreate request
func (handler *NotificationHandler) registerHandlers() error {
	urlDef, err := url.ParseRequestURI(handler.config.NotifyApiRoot)
	if err != nil {
		return fmt.Errorf("error parsing notify api root - %s", err)
	}
	basePath := path.Join(urlDef.Path, fmt.Sprintf(":%s", EncodedSessionId))
	// Register the handlers
	handler.NotifyServer.Server.POST(path.Join(basePath, "update"), handler.postSmPolicyUpdateNotification)
	handler.NotifyServer.Server.POST(path.Join(basePath, "terminate"), handler.postSmPolicyTerminateNotification)
	return nil
}

func (handler *NotificationHandler) Shutdown() {
	handler.NotifyServer.Shutdown(ShutdownTimeout)
}

// postSmPolicyUpdateNotification handles SM policy update notification request from PCF.
func (handler *NotificationHandler) postSmPolicyUpdateNotification(ctx echo.Context) error {
	client, err := relay.GetSessionProxyResponderClient(handler.cloudRegistry)
	if err != nil {
		glog.Errorf("postSmPolicyUpdateNotification failed to get SessionProxyResponderClient: %s", err)
		return fmt.Errorf("internal server error")
	}
	defer client.Close()

	var smUpdateNotify n7_client.SmPolicyNotification
	err = ctx.Bind(&smUpdateNotify)
	if err != nil {
		err = fmt.Errorf("invalid SmPolicyNotification received: %s", err)
		glog.Errorf("postSmPolicyUpdateNotification: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	sessionId, imsi, err := getSessionIdAndIMSI(ctx.Param(EncodedSessionId))
	if err != nil {
		glog.Errorf("postSmPolicyUpdateNotification unable to fetch session-id for UpdateNotify - %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	policyReauthProto := n7.GetPolicyReauthRequestProto(sessionId, imsi, smUpdateNotify.SmPolicyDecision)
	ans, err := client.PolicyReAuth(context.Background(), policyReauthProto)
	if err != nil {
		glog.Errorf("Error relaying N7 policy reauth request to gateway: %s", err)
		return fmt.Errorf("internal server error")
	}

	if !n7.IsReAuthSuccess(ans) {
		// Build the error report for the failed rules
		partialSuccessReport := n7.BuildPartialSuccessReportN7(ans)
		return ctx.JSON(http.StatusOK, partialSuccessReport)
	}
	return ctx.NoContent(http.StatusOK)
}

// postSmPolicyTerminateNotification handles SM policy Terminate notification request from PCF
func (handler *NotificationHandler) postSmPolicyTerminateNotification(ctx echo.Context) error {
	client, err := relay.GetAbortSessionResponderClient(handler.cloudRegistry)
	if err != nil {
		glog.Errorf("postSmPolicyTerminateNotification failed to get AbortSessionResponderClient: %s", err)
		return fmt.Errorf("internal server error")
	}
	defer client.Close()

	var smTermNotify n7_client.TerminationNotification
	err = ctx.Bind(&smTermNotify)
	if err != nil {
		err = fmt.Errorf("invalid TerminateNotification received: %s", err)
		glog.Errorf("postSmPolicyTerminateNotification: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sessionId, imsi, err := getSessionIdAndIMSI(ctx.Param(EncodedSessionId))
	if err != nil {
		glog.Errorf("postSmPolicyTerminateNotification unable to fetch session-id for TerminateNotify - %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ans, err := client.AbortSession(context.Background(), &protos.AbortSessionRequest{
		UserName:  imsi,
		SessionId: sessionId,
	})
	if err != nil {
		glog.Errorf("postSmPolicyTerminateNotification error relaying ASR to gateway: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error relaying to gateway")
	}

	switch ans.Code {
	case protos.AbortSessionResult_SESSION_NOT_FOUND:
		return echo.NewHTTPError(http.StatusNotFound, "Session not found")
	case protos.AbortSessionResult_USER_NOT_FOUND:
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	case protos.AbortSessionResult_GATEWAY_NOT_FOUND:
		return echo.NewHTTPError(http.StatusInternalServerError, "Gateway not found")
	case protos.AbortSessionResult_RADIUS_SERVER_ERROR:
		return echo.NewHTTPError(http.StatusInternalServerError, "Reaidus server error")
	case protos.AbortSessionResult_SESSION_REMOVED:
		return ctx.NoContent(http.StatusNoContent)
	default:
		return ctx.NoContent(http.StatusNoContent)
	}
}

func getSessionIdAndIMSI(encSessionId string) (sessionId string, imsi string, err error) {
	if len(encSessionId) == 0 {
		err = fmt.Errorf("encodedSessionId path parameter empty")
		return
	}
	sessionIdBytes, err := base64.URLEncoding.DecodeString(encSessionId)
	if err != nil {
		err = fmt.Errorf("invalid encodedSessionId path parameter, unable to decode session-id: %s", err)
		return
	}
	sessionId = string(sessionIdBytes)
	imsi, err = protos.GetIMSIwithPrefixFromSessionId(sessionId)
	if err != nil {
		err = fmt.Errorf("invalid session-id unable to decode imsi: %s", err)
		return
	}
	return
}
