/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

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
