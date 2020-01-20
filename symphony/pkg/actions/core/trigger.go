// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"

	"github.com/pkg/errors"
)

// Trigger is the base interface for implementing a Trigger
type Trigger interface {
	ID() TriggerID
	Description() string
	SupportedActionIDs() []ActionID
	SupportedFilters() []Filter
}

// EvaluateTrigger evaluates the user-supplied rule for if this rule
// should be executed or not
func EvaluateTrigger(trigger Trigger, rule Rule, inputParams map[string]interface{}) (bool, error) {
	supportedfilters := supportedFiltersMap(trigger)

	for _, ruleFilter := range rule.RuleFilters {
		filter, ok := supportedfilters[ruleFilter.FilterID]
		if !ok {
			return false, fmt.Errorf("invalid filter id: %s", ruleFilter.FilterID)
		}
		isValid, err := filter.Evaluate(ruleFilter, inputParams)
		if err != nil {
			return false, errors.Wrap(err, "evaluating ruleFilter")
		}
		if !isValid {
			return false, nil
		}
	}
	return true, nil
}

func supportedFiltersMap(t Trigger) map[string]Filter {
	ret := make(map[string]Filter)
	for _, filter := range t.SupportedFilters() {
		ret[filter.FilterID()] = filter
	}
	return ret
}
