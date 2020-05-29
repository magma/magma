/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"fmt"

	"magma/orc8r/cloud/go/services/state/indexer"

	"github.com/stretchr/testify/mock"
)

// NewMockIndexer returns a do-nothing test indexer with specified elements.
// 	- id		-- GetID return
//	- version	-- GetVersion return
//	- subs		-- GetSubscriptions return
//	- prepare	-- write PrepareReindex args to chan when called
//	- complete	-- write CompleteReindex args to chan when called
//	- index		-- write Index args to chan when called
func NewMockIndexer(id string, version indexer.Version, subs []indexer.Subscription, prepare, complete, index chan mock.Arguments) *Indexer {
	idx := &Indexer{}

	idx.On("GetID").Return(id)
	idx.On("GetVersion").Return(version)
	idx.On("GetSubscriptions").Return(subs)
	idx.On("PrepareReindex", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if prepare != nil {
			prepare <- args
		}
	}).Return(nil)
	idx.On("CompleteReindex", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if complete != nil {
			complete <- args
		}
	}).Return(nil)
	idx.On("Index", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if index != nil {
			index <- args
		}
	}).Return(nil, nil)

	return idx
}

func (_m *Indexer) String() string {
	return fmt.Sprintf("{id: %s, version: %d}", _m.GetID(), _m.GetVersion())
}
