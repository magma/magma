// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Trigger is the base interface for implementing a trigger
type Trigger interface {
	ID() TriggerID
	Description() string
	SupportedActionIDs() []ActionID
	SupportedFilters() []Filter
	Evaluate(Rule) (bool, error)
}
