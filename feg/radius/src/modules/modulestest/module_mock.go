package modulestest

import (
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockModule ...
type MockModule struct {
	mock.Mock
}

// Init ...
func (m *MockModule) Init(loggert *zap.Logger, config modules.ModuleConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

// Handle ...
func (m *MockModule) Handle(c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	args := m.Called(c, r, next)
	res, ok := args.Get(0).(*modules.Response)
	if !ok {
		return nil, args.Error(1)
	}
	return res, args.Error(1)
}
