package loader

import (
	"fbc/cwf/radius/filters"
	"fbc/cwf/radius/modules"
)

// Loader an interface for a Loader, which loads plugins
type Loader interface {
	LoadFilter(name string) (filters.Filter, error)
	LoadModule(name string) (modules.Module, error)
}
