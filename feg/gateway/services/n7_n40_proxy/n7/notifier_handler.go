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

package n7

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"

	n7_client "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/lte/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
)

// notify_handler implements the following
// - Creates a HTTP server to receive notifications from PCF
// - Handles SmPolicyUpdateNotification from PCF convert the message to proto and sends it as
//   PolicyReAuth to SessionProxyResponder(feg_relay)
// - Handles SmPolicyTerminateNotification from PCF convert the message to prot and sends it as
//   AbortSession to feg_relay.

const (
	EncodedSessionId = "encodedSessionId"
)

// RegisterHandlers registers the UpdateNotification and TerminateNotification handlers.
// The notification url is of the form
//
//	{notifyRoot}/{encodedSessionId}/{operation}
//
// Example:
//
//	http://magma-feg.magma.com/sm-policy-control/v1/MTIzNDU2Nzg5MDsxMjM0NQo=/update
//
// where
//
//	notifyRoot = http://magma-feg.magma.com/sm-policy-control/v1
//	encodedSessionId = MTIzNDU2Nzg5MDsxMjM0NQo= (Session-Id is urlencoded)
//	operation = update
//
// This notification url is send to PCF in the SmPolicyCreate request
func (c *N7Client) registerHandlers() error {
	urlDef, err := url.ParseRequestURI(c.NotifyServer.NotifierCfg.NotifyApiRoot)
	if err != nil {
		return fmt.Errorf("error parsing notify api root - %s", err)
	}
	basePath := path.Join(urlDef.Path, fmt.Sprintf(":%s", EncodedSessionId))
	// Register the handlers
	c.NotifyServer.Server.POST(path.Join(basePath, "update"), c.postSmPolicyUpdateNotification)
	c.NotifyServer.Server.POST(path.Join(basePath, "terminate"), c.postSmPolicyTerminateNotification)
	return nil
}

// postSmPolicyUpdateNotification handles SM policy update notification request from PCF.
func (c *N7Client) postSmPolicyUpdateNotification(ctx echo.Context) error {
	client, err := relay.GetSessionProxyResponderClient(c.CloudRegistry)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	sessionId, imsi, err := getSessionIdAndIMSI(ctx.Param(EncodedSessionId))
	if err != nil {
		glog.Errorf("postSmPolicyUpdateNotification unable to fetch session-id for UpdateNotify - %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	policyReauthProto := GetPolicyReauthRequestProto(sessionId, imsi, smUpdateNotify.SmPolicyDecision)
	ans, err := client.PolicyReAuth(context.Background(), policyReauthProto)
	if err != nil {
		glog.Errorf("Error relaying N7 policy reauth request to gateway: %s", err)
		return fmt.Errorf("internal server error")
	}

	if !IsReAuthSuccess(ans) {
		// Build the error report for the failed rules
		partialSuccessReport := BuildPartialSuccessReportN7(ans)
		return ctx.JSON(http.StatusOK, partialSuccessReport)
	}
	return ctx.NoContent(http.StatusOK)
}

// postSmPolicyTerminateNotification handles SM policy Terminate notification request from PCF
func (c *N7Client) postSmPolicyTerminateNotification(ctx echo.Context) error {
	client, err := relay.GetAbortSessionResponderClient(c.CloudRegistry)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessionId, imsi, err := getSessionIdAndIMSI(ctx.Param(EncodedSessionId))
	if err != nil {
		glog.Errorf("postSmPolicyTerminateNotification unable to fetch session-id for TerminateNotify - %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
