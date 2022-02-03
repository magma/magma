package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"magma/dp/cloud/go/active_mode_controller/config"
)

const (
	dialTimeoutKey               = "DIAL_TIMEOUT_SEC"
	heartbeatSendTimeoutKey      = "HEARTBEAT_SEND_TIMEOUT_SEC"
	requestTimeoutKey            = "REQUEST_TIMEOUT_SEC"
	requestProcessingIntervalKey = "REQUEST_PROCESSING_INTERVAL_SEC"
	pollingIntervalKey           = "POLLING_INTERVAL_SEC"
	grpcServiceKey               = "GRPC_SERVICE"
	grpcPortKey                  = "GRPC_PORT"
	cbsdInactivityTimeoutKey     = "CBSD_INACTIVITY_TIMEOUT_SEC"
)

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}

type ConfigTestSuite struct {
	suite.Suite
	env map[string]string
}

func (s *ConfigTestSuite) SetupTest() {
	s.env = map[string]string{
		dialTimeoutKey:               os.Getenv(dialTimeoutKey),
		heartbeatSendTimeoutKey:      os.Getenv(heartbeatSendTimeoutKey),
		requestTimeoutKey:            os.Getenv(requestTimeoutKey),
		requestProcessingIntervalKey: os.Getenv(requestProcessingIntervalKey),
		pollingIntervalKey:           os.Getenv(pollingIntervalKey),
		grpcServiceKey:               os.Getenv(grpcServiceKey),
		grpcPortKey:                  os.Getenv(grpcPortKey),
		cbsdInactivityTimeoutKey:     os.Getenv(cbsdInactivityTimeoutKey),
	}
}

func (s *ConfigTestSuite) TearDownTest() {
	for key, value := range s.env {
		err := os.Setenv(key, value)
		s.NoError(err)
	}
}

func (s *ConfigTestSuite) TestReadDefaultValues() {
	actual, err := config.Read()
	s.NoError(err)

	expected := &config.Config{
		DialTimeout:               time.Second * 60,
		HeartbeatSendTimeout:      time.Second * 10,
		RequestTimeout:            time.Second * 5,
		RequestProcessingInterval: time.Second * 10,
		PollingInterval:           time.Second * 10,
		GrpcService:               "domain-proxy-radio-controller",
		GrpcPort:                  50053,
		CbsdInactivityTimeout:     time.Hour * 4,
	}
	s.Equal(expected, actual)
}

func (s *ConfigTestSuite) TestReadFromEnv() {
	s.NoError(os.Setenv(dialTimeoutKey, "1"))
	s.NoError(os.Setenv(heartbeatSendTimeoutKey, "4"))
	s.NoError(os.Setenv(requestTimeoutKey, "2"))
	s.NoError(os.Setenv(requestProcessingIntervalKey, "5"))
	s.NoError(os.Setenv(pollingIntervalKey, "3"))
	s.NoError(os.Setenv(grpcServiceKey, "some_grpc_service"))
	s.NoError(os.Setenv(grpcPortKey, "1234"))
	s.NoError(os.Setenv(cbsdInactivityTimeoutKey, "6"))

	actual, err := config.Read()
	s.NoError(err)

	expected := &config.Config{
		DialTimeout:               time.Second * 1,
		HeartbeatSendTimeout:      time.Second * 4,
		RequestTimeout:            time.Second * 2,
		RequestProcessingInterval: time.Second * 5,
		PollingInterval:           time.Second * 3,
		GrpcService:               "some_grpc_service",
		GrpcPort:                  1234,
		CbsdInactivityTimeout:     time.Second * 6,
	}
	s.Equal(expected, actual)
}
