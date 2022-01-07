package app_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"magma/dp/cloud/go/active_mode_controller/config"
	"magma/dp/cloud/go/active_mode_controller/internal/app"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufferSize  = 16
	timeout     = time.Millisecond * 20
	currentTime = 1000
)

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, &AppTestSuite{})
}

type AppTestSuite struct {
	suite.Suite
	clock                *stubClock
	activeModeController *stubActiveModeControllerService
	radioController      *stubRadioControllerService
	appDone              chan error
	cancel               context.CancelFunc
	dialer               app.Dialer
	grpcServerDone       chan error
	grpcServer           *grpc.Server
}

func (s *AppTestSuite) SetupTest() {
	s.clock = &stubClock{ticker: make(chan time.Time, bufferSize)}
	s.activeModeController = &stubActiveModeControllerService{
		states: make(chan *active_mode.State, bufferSize),
		err:    make(chan error, bufferSize),
	}
	s.radioController = &stubRadioControllerService{
		requests: make(chan *requests.RequestPayload, bufferSize),
	}
	s.givenGrpcServer()
	s.givenAppRunning()
}

func (s *AppTestSuite) TearDownTest() {
	s.whenAppWasShutdown()
	s.thenAppWasShutdown()
	s.thenNoOtherRequestWasReceived()

	s.whenGrpcServerWasShutdown()
	s.thenGrpcServerWasShutdown()
}

func (s *AppTestSuite) TestGetStateAndSendRequests() {
	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))
}

func (s *AppTestSuite) TestFilterPendingRequests() {
	s.givenState(withPendingRequests(buildSomeState("some", "other"), "some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("other"))
}

func (s *AppTestSuite) TestCalculateHeartbeatDeadline() {
	const interval = 50 * time.Second
	const deadline = interval + 3*timeout
	now := s.clock.Now()
	timestamps := []time.Time{
		now.Add(-deadline + 2), now.Add(-deadline + 1),
		now.Add(-deadline), now.Add(-deadline - 1),
	}
	s.givenState(buildStateWithAuthorizedGrants("some", interval, timestamps...))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedHeartbeatRequests("some", "0", "1"))
}

func (s *AppTestSuite) TestAppWorkInALoop() {
	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))

	s.givenState(buildSomeState("other"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("other"))
}

func (s *AppTestSuite) TestContinueWhenFailedToGetState() {
	s.givenStateError(errors.New("some error"))
	s.whenTickerFired()

	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))
}

func (s *AppTestSuite) givenAppRunning() {
	s.appDone = make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	a := app.NewApp(
		app.WithDialer(s.dialer),
		app.WithClock(s.clock),
		app.WithConfig(&config.Config{
			DialTimeout:               timeout,
			HeartbeatSendTimeout:      timeout,
			RequestTimeout:            timeout,
			PollingInterval:           timeout,
			RequestProcessingInterval: timeout,
			GrpcService:               "",
			GrpcPort:                  0,
			CbsdInactivityTimeout:     timeout,
		}),
	)
	go func() {
		s.appDone <- a.Run(ctx)
	}()
}

func (s *AppTestSuite) givenGrpcServer() {
	listener := bufconn.Listen(bufferSize)
	s.grpcServer = grpc.NewServer()
	active_mode.RegisterActiveModeControllerServer(s.grpcServer, s.activeModeController)
	requests.RegisterRadioControllerServer(s.grpcServer, s.radioController)
	s.grpcServerDone = make(chan error)
	go func() {
		s.grpcServerDone <- s.grpcServer.Serve(listener)
	}()
	s.dialer = func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (s *AppTestSuite) givenState(state *active_mode.State) {
	s.activeModeController.states <- state
	s.activeModeController.err <- nil
}

func (s *AppTestSuite) givenStateError(err error) {
	s.activeModeController.states <- nil
	s.activeModeController.err <- err
}

func (s *AppTestSuite) whenAppWasShutdown() {
	s.cancel()
}

func (s *AppTestSuite) whenGrpcServerWasShutdown() {
	s.grpcServer.Stop()
}

func (s *AppTestSuite) whenTickerFired() {
	s.clock.ticker <- time.Time{}
}

func (s *AppTestSuite) thenAppWasShutdown() {
	select {
	case err := <-s.appDone:
		s.EqualError(err, context.Canceled.Error())
		return
	case <-time.After(timeout):
		s.Fail("failed to stop app")
	}
}

func (s *AppTestSuite) thenGrpcServerWasShutdown() {
	select {
	case err := <-s.grpcServerDone:
		s.NoError(err)
		return
	case <-time.After(timeout):
		s.Fail("failed to stop grpc server")
	}
}

func (s *AppTestSuite) thenRequestsWereEventuallyReceived(expectedRequests []*requests.RequestPayload) {
	timer := time.After(timeout)
	for _, expected := range expectedRequests {
		select {
		case actual := <-s.radioController.requests:
			s.JSONEq(expected.Payload, actual.Payload)
		case <-timer:
			s.Fail("Waiting for requests timed out")
		}
	}
}

func (s *AppTestSuite) thenNoOtherRequestWasReceived() {
	select {
	case actual := <-s.radioController.requests:
		s.Failf("Expected no more requests, got: %s", actual.Payload)
	default:
	}
}

func withPendingRequests(state *active_mode.State, name string) *active_mode.State {
	for _, cfg := range state.ActiveModeConfigs {
		if cfg.Cbsd.UserId == name {
			cfg.Cbsd.PendingRequests = []string{getExpectedSingleRequest(name)}
			break
		}
	}
	return state
}

func buildSomeState(names ...string) *active_mode.State {
	configs := make([]*active_mode.ActiveModeConfig, len(names))
	for i, name := range names {
		configs[i] = &active_mode.ActiveModeConfig{
			DesiredState: active_mode.CbsdState_Registered,
			Cbsd: &active_mode.Cbsd{
				UserId:            name,
				FccId:             name,
				SerialNumber:      name,
				State:             active_mode.CbsdState_Unregistered,
				LastSeenTimestamp: currentTime,
			},
		}
	}
	return &active_mode.State{
		ActiveModeConfigs: configs,
	}
}

func buildStateWithAuthorizedGrants(name string, interval time.Duration, timestamps ...time.Time) *active_mode.State {
	grants := make([]*active_mode.Grant, len(timestamps))
	for i, timestamp := range timestamps {
		grants[i] = &active_mode.Grant{
			Id:                     fmt.Sprintf("%d", i),
			State:                  active_mode.GrantState_Authorized,
			HeartbeatIntervalSec:   int64(interval / time.Second),
			LastHeartbeatTimestamp: timestamp.Unix(),
		}
	}
	configs := []*active_mode.ActiveModeConfig{{
		DesiredState: active_mode.CbsdState_Registered,
		Cbsd: &active_mode.Cbsd{
			Id:                name,
			State:             active_mode.CbsdState_Registered,
			Grants:            grants,
			LastSeenTimestamp: currentTime,
		},
	}}
	return &active_mode.State{
		ActiveModeConfigs: configs,
	}
}

func getExpectedRequests(name string) []*requests.RequestPayload {
	const template = `{"registrationRequest":[%s]}`
	request := fmt.Sprintf(template, getExpectedSingleRequest(name))
	return []*requests.RequestPayload{{Payload: request}}
}

func getExpectedSingleRequest(name string) string {
	const template = `{"userId":"%[1]s","fccId":"%[1]s","cbsdSerialNumber":"%[1]s"}`
	return fmt.Sprintf(template, name)
}

func getExpectedHeartbeatRequests(id string, grantIds ...string) []*requests.RequestPayload {
	const template = `{"heartbeatRequest":[%s]}`
	result := make([]*requests.RequestPayload, len(grantIds))
	for i, grantId := range grantIds {
		request := fmt.Sprintf(template, getExpectedHeartbeatRequest(id, grantId))
		result[i] = &requests.RequestPayload{Payload: request}
	}
	return result
}

func getExpectedHeartbeatRequest(id string, grantId string) string {
	const template = `{"cbsdId":"%s","grantId":"%s","operationState":"AUTHORIZED"}`
	return fmt.Sprintf(template, id, grantId)
}

type stubClock struct {
	ticker chan time.Time
}

func (s *stubClock) Now() time.Time {
	return time.Unix(currentTime, 0)
}

func (s *stubClock) Tick(_ time.Duration) *time.Ticker {
	return &time.Ticker{C: s.ticker}
}

type stubActiveModeControllerService struct {
	active_mode.UnimplementedActiveModeControllerServer
	states chan *active_mode.State
	err    chan error
}

func (s *stubActiveModeControllerService) GetState(_ context.Context, _ *active_mode.GetStateRequest) (*active_mode.State, error) {
	return <-s.states, <-s.err
}

type stubRadioControllerService struct {
	requests.UnimplementedRadioControllerServer
	requests chan *requests.RequestPayload
	err      error
}

func (s *stubRadioControllerService) UploadRequests(_ context.Context, in *requests.RequestPayload) (*requests.RequestDbIds, error) {
	s.requests <- in
	return &requests.RequestDbIds{}, s.err
}

func (s *stubRadioControllerService) GetResponse(_ context.Context, _ *requests.RequestDbId) (*requests.ResponsePayload, error) {
	return nil, errors.New("not implemented")
}
