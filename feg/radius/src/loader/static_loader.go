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

package loader

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/filters"
	filtlballocate "fbc/cwf/radius/filters/lballocate"
	filtlbcanary "fbc/cwf/radius/filters/lbcanary"
	"fbc/cwf/radius/modules"
	modadaptruckus "fbc/cwf/radius/modules/adaptruckus"
	modmsisdn "fbc/cwf/radius/modules/addmsisdn"
	modalwaysaccept "fbc/cwf/radius/modules/alwaysaccept"
	modan "fbc/cwf/radius/modules/analytics"
	modcoadynamic "fbc/cwf/radius/modules/coadynamic"
	modcoafixed "fbc/cwf/radius/modules/coafixedip"
	modcoanas "fbc/cwf/radius/modules/coanas"
	modeap "fbc/cwf/radius/modules/eap"
	modlbserve "fbc/cwf/radius/modules/lbserve"
	modmagmaacct "fbc/cwf/radius/modules/magmaacct"
	ofpanalytics "fbc/cwf/radius/modules/ofpanalytics"
	modproxy "fbc/cwf/radius/modules/proxy"
	modloopback "fbc/cwf/radius/modules/testloopback"
	testsessionstorage "fbc/cwf/radius/modules/testsessionstorage"
	modxwfv3 "fbc/cwf/radius/modules/xwfv3"
	"fbc/lib/go/radius"
	"fmt"

	"go.uber.org/zap"
)

// FilterNameMap a map from the filter-name to the API functions
type FilterNameMap map[string]func() filters.Filter

// ModuleNameMap a map from the module-name to the API functions
type ModuleNameMap map[string]func() modules.Module

// StaticLoader a loader based on a predefined set of supported plugins
type StaticLoader struct {
	logger *zap.Logger

	// the mapping of a filter-name to the API's it provides
	filters FilterNameMap

	// the mapping of a module-name to the API's it provides
	modules ModuleNameMap
}

// CWFModuleMap the available CWF modules with their names, for use by the configuration file
var CWFModuleMap = ModuleNameMap{
	"addmsisdn":          func() modules.Module { return NewModule(modmsisdn.Init, modmsisdn.Handle) },
	"analytics":          func() modules.Module { return NewModule(modan.Init, modan.Handle) },
	"eap":                func() modules.Module { return NewModule(modeap.Init, modeap.Handle) },
	"lbserve":            func() modules.Module { return NewModule(modlbserve.Init, modlbserve.Handle) },
	"proxy":              func() modules.Module { return NewModule(modproxy.Init, modproxy.Handle) },
	"ofpanalytics":       func() modules.Module { return NewModule(ofpanalytics.Init, ofpanalytics.Handle) },
	"xwfv3":              func() modules.Module { return NewModule(modxwfv3.Init, modxwfv3.Handle) },
	"testloopback":       func() modules.Module { return NewModule(modloopback.Init, modloopback.Handle) },
	"coafixedip":         func() modules.Module { return NewModule(modcoafixed.Init, modcoafixed.Handle) },
	"coanas":             func() modules.Module { return NewModule(modcoanas.Init, modcoanas.Handle) },
	"coadynamic":         func() modules.Module { return NewModule(modcoadynamic.Init, modcoadynamic.Handle) },
	"adaptruckus":        func() modules.Module { return NewModule(modadaptruckus.Init, modadaptruckus.Handle) },
	"alwaysaccept":       func() modules.Module { return NewModule(modalwaysaccept.Init, modalwaysaccept.Handle) },
	"magmaacct":          func() modules.Module { return NewModule(modmagmaacct.Init, modmagmaacct.Handle) },
	"testsessionstorage": func() modules.Module { return NewModule(testsessionstorage.Init, testsessionstorage.Handle) },
}

var CWFFilterMap = FilterNameMap{
	"lballocate": func() filters.Filter { return NewFilter(filtlballocate.Init, filtlballocate.Process) },
	"lbcanary":   func() filters.Filter { return NewFilter(filtlbcanary.Init, filtlbcanary.Process) },
}

// NewStaticLoader create a loader that loads from file system
func NewStaticLoader(logger *zap.Logger) Loader {
	return StaticLoader{logger: logger, modules: CWFModuleMap, filters: CWFFilterMap}
}

// LoadFilter returns a module invocation interface
func (l StaticLoader) LoadFilter(name string) (filters.Filter, error) {
	logger := l.logger.With(zap.String("fiter_name", name))
	logger.Info("creating filter")
	if filt, ok := l.filters[name]; ok {
		return filt(), nil
	}
	return nil, fmt.Errorf("failed to create filter %s", name)
}

// LoadModule returns a module invocation interface
func (l StaticLoader) LoadModule(name string) (modules.Module, error) {
	logger := l.logger.With(zap.String("module_name", name))
	logger.Info("creating module")
	if mod, ok := l.modules[name]; ok {
		return mod(), nil
	}
	return nil, fmt.Errorf("failed to create module %s", name)
}

// filter filters.Filter instatiation
type filter struct {
	init    filters.FilterInitFunc
	process filters.FilterProcessFunc
}

func (f filter) Init(config *config.ServerConfig) error {
	return f.init(config)
}

func (f filter) Process(c *modules.RequestContext, l string, r *radius.Request) error {
	return f.process(c, l, r)
}

// NewFilter create a new filter interface
func NewFilter(init filters.FilterInitFunc, process filters.FilterProcessFunc) filters.Filter {
	return filter{
		init:    init,
		process: process,
	}
}

// module modules.Module instatiation
type module struct {
	init   modules.ModuleInitFunc
	handle modules.ModuleHandleFunc
}

func (m module) Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return m.init(logger, config)
}

func (m module) Handle(x modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	return m.handle(x, c, r, next)
}

// NewModule create a new module interface
func NewModule(init modules.ModuleInitFunc, handle modules.ModuleHandleFunc) modules.Module {
	return module{
		init:   init,
		handle: handle,
	}
}
