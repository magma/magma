/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package filters

import (
	"fmt"

	"go.uber.org/zap"
)

// Loader an interface for a Loader, which can load filters
type Loader interface {
	LoadFilter(name string) (Filter, error)
}

// FilterNameMap a map from the filter-name to the API functions
type FilterNameMap map[string]Filter

// StaticFilterLoader a filter loader based on a predefined set of supported filters
type StaticFilterLoader struct {
	logger *zap.Logger
	// the mapping of a filter-name to the API's it provides
	filters FilterNameMap
}

// NewStaticFilterLoader create a loader that loads from file system
func NewStaticFilterLoader(logger *zap.Logger, filt FilterNameMap) Loader {
	return StaticFilterLoader{logger: logger, filters: filt}
}

// LoadFilter returns a filter invocation interface
func (l StaticFilterLoader) LoadFilter(name string) (Filter, error) {
	logger := l.logger.With(zap.String("filter_name", name))
	logger.Info("creating filter")
	if filt, ok := l.filters[name]; ok {
		return filt, nil
	}
	logger.Error("failed to create filter")
	return nil, fmt.Errorf("failed to create filter %s", name)
}
