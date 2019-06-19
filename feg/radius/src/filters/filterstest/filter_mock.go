package filterstest

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"

	"fbc/lib/go/radius"

	"github.com/stretchr/testify/mock"
)

// MockFilter ...
type MockFilter struct {
	mock.Mock
}

// Init ...
func (m *MockFilter) Init(c *config.ServerConfig) error {
	args := m.Called(c)
	return args.Error(0)
}

// Process ...
func (m *MockFilter) Process(c *modules.RequestContext, l string, r *radius.Request) error {
	args := m.Called(c, l, r)
	err := args.Get(0).(error)
	return err
}
