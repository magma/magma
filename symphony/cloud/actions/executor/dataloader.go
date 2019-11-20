// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/facebookincubator/symphony/cloud/actions/core"
)

// DataLoader is an interface for querying data for the executor
type DataLoader interface {
	QueryRules(core.TriggerID) []core.Rule
}

// BasicDataLoader is a simple implementation for querying rules.
// In the real world, this will query a database
type BasicDataLoader struct {
	Rules []core.Rule
}

// QueryRules returns all rules matches the specified TriggerID
func (b BasicDataLoader) QueryRules(triggerID core.TriggerID) (ret []core.Rule) {
	for _, rule := range b.Rules {
		if rule.TriggerID == triggerID {
			ret = append(ret, rule)
		}
	}
	return
}
