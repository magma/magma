package config_test

import (
	"os"
	"testing"
	"time"

	"magma/dp/cloud/go/active_mode_controller/config"

	"github.com/stretchr/testify/suite"

	"magma/dp/cloud/go/active_mode_controller/config"
)

const (
	dialTimeoutKey     = "DIAL_TIMEOUT"
	requestTimeoutKey  = "REQUEST_TIMEOUT"
	pollingIntervalKey = "POLLING_INTERVAL"
	grpcServiceKey     = "GRPC_SERVICE"
	grpcPortKey        = "GRPC_PORT"
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
		dialTimeoutKey:     os.Getenv(dialTimeoutKey),
		requestTimeoutKey:  os.Getenv(requestTimeoutKey),
		pollingIntervalKey: os.Getenv(pollingIntervalKey),
		grpcServiceKey:     os.Getenv(grpcPortKey),
		grpcPortKey:        os.Getenv(grpcPortKey),
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
		DialTimeout:     time.Second * 60,
		RequestTimeout:  time.Second * 5,
		PollingInterval: time.Second * 30,
		GrpcService:     "domain-proxy-radio-controller",
		GrpcPort:        50053,
	}
	s.Equal(expected, actual)
}

func (s *ConfigTestSuite) TestReadFromEnv() {
	s.NoError(os.Setenv(dialTimeoutKey, "1s"))
	s.NoError(os.Setenv(requestTimeoutKey, "2s"))
	s.NoError(os.Setenv(pollingIntervalKey, "3s"))
	s.NoError(os.Setenv(grpcServiceKey, "some_grpc_service"))
	s.NoError(os.Setenv(grpcPortKey, "1234"))

	actual, err := config.Read()
	s.NoError(err)

	expected := &config.Config{
		DialTimeout:     time.Second * 1,
		RequestTimeout:  time.Second * 2,
		PollingInterval: time.Second * 3,
		GrpcService:     "some_grpc_service",
		GrpcPort:        1234,
	}
	s.Equal(expected, actual)
}
