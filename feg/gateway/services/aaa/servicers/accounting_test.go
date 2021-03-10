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

package servicers_test

import (
	"context"
	"testing"
	"time"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/servicers"
	"magma/feg/gateway/services/aaa/store"
	"magma/feg/gateway/services/aaa/test/mock_pipelined"
	"magma/feg/gateway/services/aaa/test/mock_sessiond"

	"github.com/stretchr/testify/assert"
)

const (
	IMSI1      = "123456789012345"
	SESSIONID1 = "sessionid0001"
	SESSIONID2 = "sessionid9999"
	imsiPrefix = "IMSI"
)

func TestAccountingStartStop(t *testing.T) {
	// Run mock sessiond and pipelinedServer
	mockPipelined := mock_pipelined.NewRunningPipelined(t)
	mock_sessiond.NewRunningSessionManager(t)

	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	sessionTable := createSessionTableWithAuthenticatedUE(t, aaaCtx)
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	// Start session for UE
	_, err = accService.Start(context.Background(), aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertMacFlowInstall(t, mockPipelined)
	validateCoreSessionId(t, sessionTable, SESSIONID1)

	// Stop session for UE
	stopRequest := &protos.StopRequest{
		Cause: protos.StopRequest_USER_REQUEST,
		Ctx:   aaaCtx,
	}
	_, err = accService.Stop(context.Background(), stopRequest)
	assert.NoError(t, err)
}

func TestAccountingCreate(t *testing.T) {
	mockPipelined := mock_pipelined.NewRunningPipelined(t)
	mock_sessiond.NewRunningSessionManager(t)

	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	sessionTable := store.NewMemorySessionTable()
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	// Create session for UE
	_, err = accService.CreateSession(context.Background(), aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertMacFlowInstall(t, mockPipelined)
}

// TODO: implement radius client (protos.NewAuthorizationClient(conn) to be executed)
func testAccountingTerminate(t *testing.T) {
	mock_pipelined.NewRunningPipelined(t)
	mock_sessiond.NewRunningSessionManager(t)

	// create and authenticate a user on database
	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	sessionTable := createSessionTableWithAuthenticatedUE(t, aaaCtx)
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	// Terminate Session session for UE
	terminateRequest := &protos.TerminateSessionRequest{
		RadiusSessionId: aaaCtx.GetSessionId(),
		Imsi:            imsiPrefix + aaaCtx.GetImsi(),
	}
	_, err = accService.TerminateSession(context.Background(), terminateRequest)
	assert.NoError(t, err)
}

// TestAccountingCreateWithRecycle to trigger the recycle of the session we need to authenticate
// a ue and then try to create a session with the same imsi, but idfferent session id.
func TestAccountingCreateWithRecycle(t *testing.T) {
	mockPipelined := mock_pipelined.NewRunningPipelined(t)
	mock_sessiond.NewRunningSessionManager(t)

	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	sessionTable := createSessionTableWithAuthenticatedUE(t, aaaCtx)
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	// Start session for UE
	_, err = accService.Start(context.Background(), aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertMacFlowInstall(t, mockPipelined)
	validateCoreSessionId(t, sessionTable, SESSIONID1)

	// Create session same IMSI, different SessionId
	aaaCtx.SessionId = SESSIONID2
	_, err = accService.CreateSession(context.Background(), aaaCtx)
	assert.NoError(t, err)
}

// TestAccountingBadMACaddress uses a wrong mac address.
// This shouldn't cause an error.
func TestAccountingBadMACaddress(t *testing.T) {
	// Run mock sessiond and pipelined
	mock_pipelined.NewRunningPipelined(t)
	mock_sessiond.NewRunningSessionManager(t)

	// MAC address has a missing character at the end on start
	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	aaaCtx.Apn = "98-76-54-AA-BB-C:Wifi-Offload-hotspot20"
	sessionTable := createSessionTableWithAuthenticatedUE(t, aaaCtx)
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	_, err = accService.Start(context.Background(), aaaCtx)
	assert.NoError(t, err)
}

func TestAccountingSessiondErrors(t *testing.T) {
	mock_pipelined.NewRunningPipelined(t)
	sessiond := mock_sessiond.NewRunningSessionManager(t)

	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	sessionTable := createSessionTableWithAuthenticatedUE(t, aaaCtx)
	aaaConfig := getAAAConfig()
	accService, err := servicers.NewAccountingService(sessionTable, aaaConfig)
	assert.NoError(t, err)

	// Force error on sessiond on Start
	sessiond.ReturnErrors(true)
	_, err = accService.Start(context.Background(), aaaCtx)
	assert.Error(t, err)

	// Force error on sessiond on End
	sessiond.ReturnErrors(false)
	_, err = accService.Start(context.Background(), aaaCtx)
	assert.NoError(t, err)

	sessiond.ReturnErrors(true)
	stopRequest := &protos.StopRequest{
		Cause: protos.StopRequest_USER_REQUEST,
		Ctx:   aaaCtx,
	}
	_, err = accService.Stop(context.Background(), stopRequest)
	assert.Error(t, err)
}

func createSessionTableWithAuthenticatedUE(t *testing.T, aaaCtx *protos.Context) aaa.SessionTable {
	// Create a shared Session Table and add IMSI and sessionId
	sessionTable := store.NewMemorySessionTable()
	_, err := sessionTable.AddSession(aaaCtx, time.Minute*10, nil)
	assert.NoError(t, err)
	return sessionTable
}

func getAAAConfig() *mconfig.AAAConfig {
	return &mconfig.AAAConfig{
		IdleSessionTimeoutMs: 0,
		AccountingEnabled:    true,
		CreateSessionOnAuth:  false,
	}
}

func getAAAcontext(sessionId, IMSI string) *protos.Context {
	return &protos.Context{
		SessionId: sessionId,
		Imsi:      IMSI,
		Msisdn:    "0015551234567",
		Apn:       "98-76-54-AA-BB-CC:Wifi-Offload-hotspot20",
		MacAddr:   "12-34-AB-CD-EF-FF",
	}
}

func validateCoreSessionId(t *testing.T, sessionTable aaa.SessionTable, sessionId string) {
	session := sessionTable.GetSession(sessionId)
	assert.NotNil(t, session)
	aaaCtx := session.GetCtx()
	assert.NotNil(t, aaaCtx)
	coreSessionId := aaaCtx.GetAcctSessionId()
	assert.Contains(t, coreSessionId, "-")
}
