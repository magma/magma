/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/user"
	"testing"
	"time"

	"fbc/cwf/radius/config"
	"fbc/cwf/radius/filters/filterstest"
	"fbc/cwf/radius/loader"
	"fbc/cwf/radius/loader/loaderstest"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/modulestest"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type FullRADIUSSessiontWithAnalyticsModulesTestParam struct {
	CallingStationID string
	CalledStationID  string
	AcctSessionID    string
	FramedIPAddr     net.IP
	NasIdentifier    string
	AnalyticsModIdx  int                 // the index of the Analytics module within the module chain
	Config           config.ServerConfig // the server configuration
	Server           *Server             // the RADIUS server we created for the test
}

// TestAnalyticsModulesAuthenticate tests the Analytics module handling of the Authenticate RADIUS packet
func TestAnalyticsModulesAuthenticate(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	testParam := analyticsModuleTestEnvCreate(t, logger)
	testFullRADIUSSessiontWithAnalyticsModulesAuthenticate(t, logger, testParam)
	analyticsModuleTestEnvDestroy(testParam)
}

// TestAnalyticsModulesAccountingStart test the processing of Accounting-Start RADIUS packet, with/out processing of Auth-Request
func TestAnalyticsModulesAccountingStart(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	testParam := analyticsModuleTestEnvCreate(t, logger)
	// case 1: run the packet as if the Auth packet was dropped & not received. we should ignore the packet bcz we didnt create the XWFEntRadiusSession
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStart(t, logger, testParam, true)
	// case 2: create the session state we expect would be created by the AuthRequest & then fire the packet
	server := testParam.Server
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, testParam.AcctSessionID)
	stg := server.getSessionStateAPI(sessionID)
	stg.Set(session.State{RadiusSessionFBID: 123 /* non-zero value*/})
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStart(t, logger, testParam, false)
	analyticsModuleTestEnvDestroy(testParam)
}

// TestAnalyticsModulesAccountingUpdate test the processing of Accounting-Update RADIUS packet, with/out processing of Auth-Request
func TestAnalyticsModulesAccountingUpdate(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	testParam := analyticsModuleTestEnvCreate(t, logger)
	// case 1: run the packet as if the Auth packet was dropped & not received. we should ignore the packet bcz we didnt create the XWFEntRadiusSession
	testFullRADIUSSessiontWithAnalyticsModulesAccountingUpdate(t, logger, testParam, true)
	// case 2: create the session state we expect would be created by the AuthRequest & then fire the packet
	server := testParam.Server
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, "")
	stg := server.getSessionStateAPI(sessionID)
	stg.Set(session.State{RadiusSessionFBID: 123 /* non-zero value*/})
	testFullRADIUSSessiontWithAnalyticsModulesAccountingUpdate(t, logger, testParam, false)
	analyticsModuleTestEnvDestroy(testParam)
}

// TestAnalyticsModulesAccountingStop test the processing of Accounting-Stop RADIUS packet, with/out processing of Auth-Request
func TestAnalyticsModulesAccountingStop(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	testParam := analyticsModuleTestEnvCreate(t, logger)
	testParam.AcctSessionID = "1003/1.0.0.6/567" // Test With AP's Accounting Session ID
	// case 1: run the packet as if the Auth packet was dropped & not received. we should ignore the packet bcz we didnt create the XWFEntRadiusSession
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStop(t, logger, testParam, true)
	// case 2: create the session state we expect would be created by the AuthRequest & then fire the packet
	server := testParam.Server
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, testParam.AcctSessionID)
	stg := server.getSessionStateAPI(sessionID)
	stg.Set(session.State{RadiusSessionFBID: 123 /* non-zero value*/})
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStop(t, logger, testParam, false)
	analyticsModuleTestEnvDestroy(testParam)
}

// full session lifetime test of Analytics module.
func TestFullRADIUSSessiontWithAnalyticsModules(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	testParam := analyticsModuleTestEnvCreate(t, logger)

	// Act & Assert
	// step 1: Authorization, to establish the session
	testFullRADIUSSessiontWithAnalyticsModulesAuthenticate(t, logger, testParam)
	// step 2: Acct-start, no Gigawords defined
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStart(t, logger, testParam, false)
	// step 3: Acct-update, Gigawords defined as 0
	testFullRADIUSSessiontWithAnalyticsModulesAccountingUpdate(t, logger, testParam, false)
	// step 4: Acct-update, Gigawords defined as > 0
	testFullRADIUSSessiontWithAnalyticsModulesAccountingStop(t, logger, testParam, false)

	// Cleanup
	analyticsModuleTestEnvDestroy(testParam)
}

func TestRequestWithModules(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	moduleChain := []string{"auth"}
	moduleCount := []int{3}
	config := getConfigWithAuthListener(t, moduleChain, moduleCount, true)

	// Create server with one filter and 3 modules concatenated.
	// Test will make sure they are all called as expected

	mFilter1 := createMockFilterWithReturn(nil)
	fLoader := loaderstest.MockLoader{}
	fLoader.On("LoadFilter", "filter.auth.1").Return(mFilter1, nil)

	mModule1 := createMockHandlerWithReturn(&modules.Response{
		Code: radius.CodeAccessAccept,
		Attributes: map[radius.Type][]radius.Attribute{
			10: {[]byte("hello")},
		},
	}, nil)
	mModule2 := createMockHandlerWithReturn(&modules.Response{}, nil)
	mModule3 := createMockHandlerWithReturn(&modules.Response{}, nil)

	loader := loaderstest.MockLoader{}
	loader.On("LoadModule", "module.auth.1").Return(mModule1, nil)
	loader.On("LoadModule", "module.auth.2").Return(mModule2, nil)
	loader.On("LoadModule", "module.auth.3").Return(mModule3, nil)

	server, err := New(config, logger, &loader)
	assert.Equal(t, err, nil)
	isReady := server.StartAndWait()
	require.True(t, isReady, "failed to initialize the server")

	// Act
	packet := radius.New(radius.CodeAccessRequest, []byte(config.Secret))
	rfc2865.UserName_SetString(packet, "tim")
	rfc2865.UserPassword_SetString(packet, "12345")
	port := config.Listeners[0].Extra["Port"].(int)
	response, err := radius.Exchange(
		context.Background(),
		packet,
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		log.Fatal(err)
	}
	server.Stop()

	// Assert
	assert.Equal(t, response.Code, radius.CodeAccessAccept)
	loader.AssertExpectations(t)
	mModule1.AssertExpectations(t)
	mModule2.AssertExpectations(t)
	mModule3.AssertExpectations(t)
}

func TestModuleFailsToLoad(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	config := getConfigWithAuthListener(t, []string{"auth"}, []int{1}, true)

	// Create loaders
	loader := loaderstest.MockLoader{}
	defer loader.AssertExpectations(t)
	loader.On("LoadModule", "module.auth.1").Return(
		modules.Module(nil),
		errors.New("failed to load"),
	)

	// Act
	server, err := New(config, logger, &loader)

	// Assert
	assert.Nil(t, server)
	assert.Equal(t, err.Error(), "failed to load")
}

func TestModuleFailsToInit(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	config := getConfigWithAuthListener(t, []string{"auth"}, []int{1}, true)

	loader := loaderstest.MockLoader{}
	nModule := modulestest.MockModule{}
	nModule.On("Init", mock.Anything).
		Return(errors.New("failed to init")).Once()

	defer loader.AssertExpectations(t)
	loader.On("LoadModule", "module.auth.1").Return(&nModule, nil)

	// Act
	server, err := New(config, logger, &loader)

	// Assert
	assert.Nil(t, server)
	assert.Equal(t, err.Error(), "failed to init")
}

func TestModuleFailsToHandle(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	config := getConfigWithAuthListener(t, []string{"auth"}, []int{1}, true)

	loader := loaderstest.MockLoader{}

	mModule1 := createMockHandlerWithReturn(nil, errors.New("failed to handle"))

	loader.On("LoadModule", "module.auth.1").Return(mModule1, nil)

	server, err := New(config, logger, &loader)
	assert.Equal(t, err, nil)
	isReady := server.StartAndWait()
	require.True(t, isReady, "failed to initialize the server")

	// Act
	packet := radius.New(radius.CodeAccessRequest, []byte(config.Secret))
	client := radius.Client{
		Retry: 0,
	}

	go func() {
		port := config.Listeners[0].Extra["Port"].(int)
		client.Exchange(
			context.Background(),
			packet,
			fmt.Sprintf(":%d", port),
		)
	}()
	time.Sleep(time.Millisecond * 500)

	// Assert
	assert.NotNil(t, server)
	mModule1.AssertExpectations(t)
}

func TestFilterFailsToInit(t *testing.T) {
	// Arrange
	config := getConfigWithFilters(t, []string{"filter.1"})

	nFilter := filterstest.MockFilter{}
	nFilter.On("Init", mock.Anything).
		Return(errors.New("failed to init")).Once()

	loader := loaderstest.MockLoader{}
	loader.On("LoadFilter", "filter.1").Return(&nFilter, nil)
	defer loader.AssertExpectations(t)

	// Act
	server, err := New(config, zap.NewNop(), &loader)

	// Assert
	assert.Nil(t, server)
	assert.Equal(t, err.Error(), "failed to init")
}

func TestFilterFailsToProcess(t *testing.T) {
	// Arrange
	config := getConfigWithFilters(t, []string{"filter.1"})

	mFilter1 := createMockFilterWithReturn(errors.New("failed to handle"))

	loader := loaderstest.MockLoader{}
	loader.On("LoadFilter", "filter.1").Return(mFilter1, nil)

	mModule1 := createMockHandlerWithReturn(&modules.Response{}, nil)
	loader.On("LoadModule", "module.auth.1").Return(mModule1, nil)

	server, err := New(config, zap.NewNop(), &loader)
	assert.Equal(t, err, nil)
	isReady := server.StartAndWait()
	require.True(t, isReady, "failed to initialize the server")

	// Act
	packet := radius.New(radius.CodeAccessRequest, []byte(config.Secret))
	client := radius.Client{
		Retry: 0,
	}

	go func() {
		port := config.Listeners[0].Extra["Port"].(int)
		client.Exchange(
			context.Background(),
			packet,
			fmt.Sprintf(":%d", port),
		)
	}()
	time.Sleep(time.Millisecond * 500)

	// Assert
	assert.NotNil(t, server)
	mFilter1.AssertExpectations(t)
}

func TestDedup(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	moduleChain := []string{"auth"}
	moduleCount := []int{1}
	config := getConfigWithAuthListener(t, moduleChain, moduleCount, true)
	mModule1 := createMockHandlerWithReturn(nil, nil)

	loader := loaderstest.MockLoader{}
	loader.On("LoadModule", "module.auth.1").Return(mModule1, nil)

	server, err := New(config, logger, &loader)
	assert.Equal(t, err, nil)
	isReady := server.StartAndWait()
	require.True(t, isReady, "failed to initialize the server")
	packet := radius.New(radius.CodeAccessRequest, []byte(config.Secret))
	rfc2865.UserName_SetString(packet, "tim")
	rfc2865.UserPassword_SetString(packet, "12345")

	// Act (no response, package will be sent multiple times)
	radius.DefaultClient.Retry, _ = time.ParseDuration("10ms")
	deadline := time.Now().Add(time.Millisecond * 100)
	d, cancelFunc := context.WithDeadline(context.Background(), deadline)
	port := config.Listeners[0].Extra["Port"].(int)
	_, _ = radius.Exchange(
		d,
		packet,
		fmt.Sprintf(":%d", port),
	)
	server.Stop()
	cancelFunc()

	// Assert
	loader.AssertExpectations(t)
	mModule1.AssertExpectations(t)

	// This ASSERT is a bit tricky;
	// The test is timed to take 100 millisec, retrying packet every 10ms
	// This timing is expected to give us 9 retries (first attempt is not
	// counted) however, this integration test depends on timing, so we
	// ease the expected count.
	// IF THIS ASSERT FAILS EITHER THERE'S A JITTER IN THE RADIUS CLIENT
	// TIMER (which is an issue by itself) OR THE RETRY LOGIC BROKE
	assert.True(t, server.GetDroppedCount() > 5)
}

func getConfigWithFilters(t *testing.T, filterNames []string) config.ServerConfig {
	conf := getConfigWithAuthListener(t, []string{"auth"}, []int{1}, true)
	conf.Filters = filterNames
	return conf
}

// getConfigWithAuthListener create a config with a chain of modules as set by args:
// moduleChain the names of modules to chain - ordered is maintained
// moduleCount the numner of module instances from the name in the same index within 'moduleChain'
func getConfigWithAuthListener(t *testing.T, moduleChain []string, moduleCount []int, beutifyName bool) config.ServerConfig {
	require.Equal(t, len(moduleChain), len(moduleCount),
		"chain of module must have same number of elements in names & count")
	initialPort := 2000 + rand.Intn(5000)
	result := config.ServerConfig{
		Secret: `123456`,
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{},
			LiveTier:     config.TierRouting{},
			Canaries:     []config.Canary{},
		},
		SessionStorage: &config.SessionStorageConfig{
			StorageType: "memory",
		},
	}
	listenerCfg := config.ListenerConfig{
		Name: "listener.0",
		Extra: map[string]interface{}{
			"Port": initialPort,
		},
		Modules: []config.ModuleDescriptor{},
		Type:    "udp",
	}
	for mi, moduleName := range moduleChain {
		// Add module instances
		for j := 1; j <= moduleCount[mi]; j++ {
			modName := moduleChain[mi]
			if beutifyName {
				modName = fmt.Sprintf("module.%s.%d", moduleName, j)
			}
			listenerCfg.Modules = append(
				listenerCfg.Modules,
				config.ModuleDescriptor{
					Name:   modName,
					Config: make(modules.ModuleConfig),
				},
			)
		}
		result.Listeners = append(result.Listeners, listenerCfg)
	}

	return result
}

func createMockHandlerWithReturn(r *modules.Response, err error) *modulestest.MockModule {
	mModule := modulestest.MockModule{}
	mModule.On("Init", mock.Anything, mock.Anything).
		Return(nil).On(
		"Handle",
		mock.AnythingOfType("*modules.RequestContext"),
		mock.AnythingOfType("*radius.Request"),
		mock.AnythingOfType("modules.Middleware"),
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(*modules.RequestContext)
		req := args.Get(1).(*radius.Request)
		next := args.Get(2).(modules.Middleware)
		next(ctx, req)
	}).Return(r, err)
	return &mModule
}

func createMockFilterWithReturn(err error) *filterstest.MockFilter {
	mFilter := filterstest.MockFilter{}
	mFilter.On("Init", mock.Anything).Return(nil).On(
		"Process",
		mock.AnythingOfType("*modules.RequestContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*radius.Request"),
	).Return(err)
	return &mFilter
}

// generate the Authentication packet & verify processing, final state
func testFullRADIUSSessiontWithAnalyticsModulesAuthenticate(
	t *testing.T,
	logger *zap.Logger,
	testParam *FullRADIUSSessiontWithAnalyticsModulesTestParam,
) {
	// Arrange
	config := testParam.Config
	server := testParam.Server
	pkt := radius.New(radius.CodeAccessRequest, []byte(config.Secret))
	rfc2865.FramedIPAddress_Add(pkt, testParam.FramedIPAddr)
	rfc2865.NASIPAddress_Add(pkt, net.IP{1, 0, 0, 2})
	rfc2865.CalledStationID_SetString(pkt, testParam.CalledStationID)
	rfc2865.NASIdentifier_SetString(pkt, testParam.NasIdentifier)
	rfc2866.AcctSessionID_SetString(pkt, testParam.AcctSessionID)
	rfc2865.CallingStationID_SetString(pkt, testParam.CallingStationID)
	port := config.Listeners[testParam.AnalyticsModIdx].Extra["Port"].(int)

	// Act
	logger.Debug("sending RADIUS access-request packet", zap.Int("port", port))
	response, err := radius.Exchange(context.Background(), pkt, fmt.Sprintf(":%d", port))
	require.NoError(t, err, "failed to exchange RADIUS packet")
	logger.Debug("RADIUS auth response", zap.Any("response", response))

	// Assert
	// Verify that response has the attributes we sent (test-module generates that response)
	require.Equal(
		t,
		testParam.FramedIPAddr,
		rfc2865.FramedIPAddress_Get(response),
		"expecting response to hold all AVP's from request",
	)

	// Read the session state & verify the accounting counters were updated (if we decide to save them)
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, testParam.AcctSessionID)
	sessionState, err := server.getSessionState(sessionID)
	require.NoError(t, err, "cant find session state for pkt we just sent")
	require.NotEqual(t, sessionState.RadiusSessionFBID, 0)
}

// generate the Accounting-Start packet & verify processing, final state
func testFullRADIUSSessiontWithAnalyticsModulesAccountingStart(
	t *testing.T,
	logger *zap.Logger,
	testParam *FullRADIUSSessiontWithAnalyticsModulesTestParam,
	shouldFail bool,
) {
	// Arrange
	config := testParam.Config
	server := testParam.Server
	pkt := radius.New(radius.CodeAccountingRequest, []byte(config.Secret))
	rfc2866.AcctStatusType_Set(pkt, rfc2866.AcctStatusType_Value_Start)
	rfc2865.CalledStationID_SetString(pkt, testParam.CalledStationID)
	rfc2865.CallingStationID_SetString(pkt, testParam.CallingStationID)
	rfc2865.NASIdentifier_SetString(pkt, testParam.NasIdentifier)
	rfc2866.AcctInputOctets_Set(pkt, 1111)
	rfc2866.AcctOutputOctets_Set(pkt, 1112)
	port := config.Listeners[testParam.AnalyticsModIdx].Extra["Port"].(int)

	// Act
	logger.Debug("sending RADIUS accounting-request packet", zap.Int("port", port))
	response, err := radius.Exchange(context.Background(), pkt, fmt.Sprintf(":%d", port))
	require.NoError(t, err, "failed to exchange RADIUS packet")
	logger.Debug("RADIUS accounting-request response", zap.Any("response", response))

	// Assert
	// Verify that response has the attributes we sent (test-module generates that response)
	require.Equal(
		t,
		rfc2866.AcctStatusType_Value_Start,
		rfc2866.AcctStatusType_Get(response),
		"expecting response to hold all AVP's from request",
	)

	// Read the session state & verify the accounting counters were updated (if we decide to save them)
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, "")
	sessionState, err := server.getSessionState(sessionID)
	if shouldFail {
		require.Error(t, err, "expecting no session state for pkt we just sent")
		require.Nil(t, sessionState, 0)
	} else {
		require.NoError(t, err, "cant find session state for pkt we just sent")
		require.NotEqual(t, sessionState.RadiusSessionFBID, 0)
	}
}

// generate the Accounting-Intermediate-Update packet & verify processing, final state
func testFullRADIUSSessiontWithAnalyticsModulesAccountingUpdate(
	t *testing.T,
	logger *zap.Logger,
	testParam *FullRADIUSSessiontWithAnalyticsModulesTestParam,
	shouldFail bool,
) {
	// Arrange
	config := testParam.Config
	server := testParam.Server
	pkt := radius.New(radius.CodeAccountingRequest, []byte(config.Secret))
	rfc2866.AcctStatusType_Set(pkt, rfc2866.AcctStatusType_Value_InterimUpdate)
	rfc2865.CalledStationID_SetString(pkt, testParam.CalledStationID)
	rfc2865.CallingStationID_SetString(pkt, testParam.CallingStationID)
	if len(testParam.AcctSessionID) > 0 {
		rfc2866.AcctSessionID_SetString(pkt, testParam.AcctSessionID)
	}
	rfc2865.NASIdentifier_SetString(pkt, testParam.NasIdentifier)
	rfc2866.AcctInputOctets_Set(pkt, 1111)
	rfc2869.AcctInputGigawords_Set(pkt, 0)
	rfc2866.AcctOutputOctets_Set(pkt, 1112)
	rfc2869.AcctOutputGigawords_Set(pkt, 0)
	port := config.Listeners[testParam.AnalyticsModIdx].Extra["Port"].(int)

	// Act
	logger.Debug("sending RADIUS accounting-interim-update packet", zap.Int("port", port))
	response, err := radius.Exchange(context.Background(), pkt, fmt.Sprintf(":%d", port))
	require.NoError(t, err, "failed to exchange RADIUS packet")
	logger.Debug("RADIUS accounting-interim-update response", zap.Any("response", response))

	// Assert
	// Verify that response has the attributes we sent (test-module generates that response)
	require.Equal(
		t,
		rfc2866.AcctStatusType_Value_InterimUpdate,
		rfc2866.AcctStatusType_Get(response),
		"expecting response to hold all AVP's from request",
	)

	// Read the session state & verify the accounting counters were updated (if we decide to save them)
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, testParam.AcctSessionID)
	sessionState, err := server.getSessionState(sessionID)
	if shouldFail {
		require.Error(t, err, "expecting no session state for pkt we just sent")
		require.Nil(t, sessionState, 0)
	} else {
		require.NoError(t, err, "cant find session state for pkt we just sent")
		require.NotEqual(t, sessionState.RadiusSessionFBID, 0)
	}
}

// generate the Accounting-Stop packet & verify processing, final state
func testFullRADIUSSessiontWithAnalyticsModulesAccountingStop(
	t *testing.T,
	logger *zap.Logger,
	testParam *FullRADIUSSessiontWithAnalyticsModulesTestParam,
	shouldFail bool,
) {
	// Arrange
	config := testParam.Config
	server := testParam.Server
	pkt := radius.New(radius.CodeAccountingRequest, []byte(config.Secret))
	rfc2866.AcctStatusType_Set(pkt, rfc2866.AcctStatusType_Value_Stop)
	rfc2865.CalledStationID_SetString(pkt, testParam.CalledStationID)
	rfc2865.CallingStationID_SetString(pkt, testParam.CallingStationID)
	if len(testParam.AcctSessionID) > 0 {
		rfc2866.AcctSessionID_SetString(pkt, testParam.AcctSessionID)
	}
	rfc2865.NASIdentifier_SetString(pkt, testParam.NasIdentifier)
	rfc2866.AcctInputOctets_Set(pkt, 1111)
	rfc2869.AcctInputGigawords_Set(pkt, 1)
	rfc2866.AcctOutputOctets_Set(pkt, 1112)
	rfc2869.AcctOutputGigawords_Set(pkt, 2)
	port := config.Listeners[testParam.AnalyticsModIdx].Extra["Port"].(int)

	// Act
	logger.Debug("sending RADIUS accounting-stop packet", zap.Int("port", port))
	response, err := radius.Exchange(context.Background(), pkt, fmt.Sprintf(":%d", port))
	require.NoError(t, err, "failed to exchange RADIUS packet")
	logger.Debug("RADIUS accounting-stop response", zap.Any("response", response))

	// Assert
	// Verify that response has the attributes we sent (test-module generates that response)
	require.Equal(
		t,
		rfc2866.AcctStatusType_Value_Stop,
		rfc2866.AcctStatusType_Get(response),
		"expecting response to hold all AVP's from request",
	)

	// Read the session state & verify the accounting counters were updated (if we decide to save them)
	sessionID := getSessionIDStrings(server, testParam.CallingStationID, testParam.CalledStationID, testParam.AcctSessionID)
	sessionState, err := server.getSessionState(sessionID)
	if shouldFail {
		require.Error(t, err, "expecting no session state for pkt we just sent")
		require.Nil(t, sessionState, 0)
	} else {
		require.NoError(t, err, "cant find session state for pkt we just sent")
		require.NotEqual(t, sessionState.RadiusSessionFBID, 0)
	}
}

// setup the test env for Analytics module tests
func analyticsModuleTestEnvCreate(t *testing.T, logger *zap.Logger) *FullRADIUSSessiontWithAnalyticsModulesTestParam {
	// session key identification - must be identical in all RADIUS packets
	testParam := &FullRADIUSSessiontWithAnalyticsModulesTestParam{
		FramedIPAddr:     net.IP{1, 0, 0, 1},
		CallingStationID: "1.0.0.6",
		CalledStationID:  "1.0.0.3",
		NasIdentifier:    "1.0.0.4",
	}
	u, err := user.Current()
	require.NoError(t, err, "failed getting user")

	testParam.Config = getConfigWithAuthListener(t, []string{"analytics", "testloopback"}, []int{1, 1}, false)
	testParam.AnalyticsModIdx = 0
	// add "analytics" module config
	analyticsMod := &testParam.Config.Listeners[testParam.AnalyticsModIdx].Modules[0]
	analyticsMod.Config["AccessToken"] = "dummy token" // valid token not required bcz GraphQL calls are in dry-run mode
	analyticsMod.Config["GraphQLURL"] = fmt.Sprintf("https://graph.%s.sb.expresswifi.com/graphql", u.Username)
	analyticsMod.Config["DryRunGraphQL"] = true

	// Create server with Analytics module handler.
	mLoader := loader.NewStaticLoader(logger)
	testParam.Server, err = New(testParam.Config, logger, mLoader)
	require.NoError(t, err, "failed to create server")
	isReady := testParam.Server.StartAndWait()
	require.True(t, isReady, "failed to initialize the server")

	return testParam
}

// analyticsModuleTestEnvDestroy destroy the test env that was created for Analytics module tests
func analyticsModuleTestEnvDestroy(testParam *FullRADIUSSessiontWithAnalyticsModulesTestParam) {
	testParam.Server.Stop()
}

func getSessionIDStrings(server *Server, calling string, called string, acctSessionId string) string {
	r := radius.Request{
		Packet: &radius.Packet{
			Attributes: radius.Attributes{
				rfc2865.CallingStationID_Type: []radius.Attribute{radius.Attribute(calling)},
				rfc2865.CalledStationID_Type:  []radius.Attribute{radius.Attribute(called)},
				rfc2866.AcctSessionID_Type:    []radius.Attribute{radius.Attribute(acctSessionId)},
			},
		},
	}
	return server.GetSessionID(&r)
}
