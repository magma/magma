/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package alert

import (
	"fmt"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/restrictor"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type File struct {
	RuleGroups []rulefmt.RuleGroup `yaml:"groups"`
}

func NewFile(tenantID string) *File {
	return &File{
		RuleGroups: []rulefmt.RuleGroup{{
			Name: tenantID,
		}},
	}
}

// Rules returns the rule configs from this file
func (f *File) Rules() []rulefmt.Rule {
	return f.RuleGroups[0].Rules
}

// GetRule returns the specific rule by name. Nil if it isn't found
func (f *File) GetRule(rulename string) *rulefmt.Rule {
	for _, rule := range f.RuleGroups[0].Rules {
		if rule.Alert == rulename {
			return &rule
		}
	}
	return nil
}

// AddRule appends a new rule to the list of rules in this file
func (f *File) AddRule(rule rulefmt.Rule) {
	f.RuleGroups[0].Rules = append(f.RuleGroups[0].Rules, rule)
}

// ReplaceRule replaces an existing rule. Returns error if rule does not
// exist already
func (f *File) ReplaceRule(newRule rulefmt.Rule) error {
	ruleIdx := -1
	for idx, rule := range f.RuleGroups[0].Rules {
		if rule.Alert == newRule.Alert {
			ruleIdx = idx
		}
	}
	if ruleIdx < 0 {
		return fmt.Errorf("rule %s does not exist", newRule.Alert)
	}

	f.RuleGroups[0].Rules[ruleIdx] = newRule
	return nil
}

func (f *File) DeleteRule(name string) error {
	rules := f.RuleGroups[0].Rules
	for idx, rule := range rules {
		if rule.Alert == name {
			f.RuleGroups[0].Rules = append(rules[:idx], rules[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("alert with name %s not found", name)
}

// SecureRule attaches a label for tenantID to the given alert expression to
// to ensure that only metrics owned by this tenant can be alerted on
func SecureRule(matcherName, matcherValue string, rule *rulefmt.Rule) error {
	tenantLabels := map[string]string{matcherName: matcherValue}
	queryRestrictor := restrictor.NewQueryRestrictor(tenantLabels)

	restrictedExpression, err := queryRestrictor.RestrictQuery(rule.Expr)
	if err != nil {
		return err
	}
	rule.Expr = restrictedExpression
	if rule.Labels == nil {
		rule.Labels = make(map[string]string)
	}
	rule.Labels[matcherName] = matcherValue
	return nil
}

// RuleJSONWrapper Provides a struct to marshal/unmarshal into a rulefmt.Rule
// since rulefmt does not support json encoding
type RuleJSONWrapper struct {
	Record      string            `json:"record,omitempty"`
	Alert       string            `json:"alert,omitempty"`
	Expr        string            `json:"expr"`
	For         string            `json:"for,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (r *RuleJSONWrapper) ToRuleFmt() (rulefmt.Rule, error) {
	if r.Labels == nil {
		r.Labels = make(map[string]string)
	}
	if r.Annotations == nil {
		r.Annotations = make(map[string]string)
	}

	rule := rulefmt.Rule{
		Record:      r.Record,
		Alert:       r.Alert,
		Expr:        r.Expr,
		Labels:      r.Labels,
		Annotations: r.Annotations,
	}
	if r.For != "" {
		modelFor, err := model.ParseDuration(r.For)
		if err != nil {
			return rulefmt.Rule{}, err
		}
		rule.For = modelFor
	}
	return rule, nil
}
