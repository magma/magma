// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/actions/core"
)

// DataLoader is an interface for querying data for the executor
type DataLoader interface {
	QueryRules(context.Context, core.TriggerID) ([]core.Rule, error)
}

// BasicDataLoader is a simple implementation for querying rules.
// In the real world, this will query a database
type BasicDataLoader struct {
	Rules []core.Rule
}

// QueryRules returns all rules matches the specified TriggerID
func (b BasicDataLoader) QueryRules(ctx context.Context, triggerID core.TriggerID) (ret []core.Rule, _err error) {
	for _, rule := range b.Rules {
		if rule.TriggerID == triggerID {
			ret = append(ret, rule)
		}
	}
	return
}
