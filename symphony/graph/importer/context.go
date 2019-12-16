// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
)

type ctxKey struct{}

type importContext struct {
	lowestHierarchyIndex        int
	indexToLocationTypeID       map[int]string
	equipmentTypeNameToID       map[string]string
	equipmentTypeIDToProperties map[string][]string
	propNameToIndex             map[string]int
	typeIDsToProperties         map[string][]string
}

func newImportContext(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKey{}, &importContext{
		lowestHierarchyIndex:        -1,
		indexToLocationTypeID:       make(map[int]string),
		equipmentTypeNameToID:       make(map[string]string),
		equipmentTypeIDToProperties: make(map[string][]string),
		propNameToIndex:             make(map[string]int),
		typeIDsToProperties:         make(map[string][]string),
	})
}

func getImportContext(ctx context.Context) *importContext {
	ld, _ := ctx.Value(ctxKey{}).(*importContext)
	return ld
}
