package loaderstest

import (
	"fbc/cwf/radius/filters"
	"fbc/cwf/radius/modules"

	"github.com/stretchr/testify/mock"
)

// MockLoader ...
type MockLoader struct {
	mock.Mock
}

// LoadFilter ...
func (l *MockLoader) LoadFilter(name string) (filters.Filter, error) {
	args := l.Called(name)
	f, ok := args.Get(0).(filters.Filter)
	if !ok {
		return nil, args.Error(1)
	}
	return f, args.Error(1)
}

// LoadModule ...
func (l *MockLoader) LoadModule(name string) (modules.Module, error) {
	args := l.Called(name)
	m, ok := args.Get(0).(modules.Module)
	if !ok {
		return nil, args.Error(1)
	}
	return m, args.Error(1)
}
