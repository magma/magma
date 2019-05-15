/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package alert

import (
	"magma/orc8r/cloud/go/services/metricsd/obsidian/security"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type File struct {
	RuleGroups []rulefmt.RuleGroup `yaml:"groups"`
}

func NewFile(networkID string) *File {
	return &File{
		RuleGroups: []rulefmt.RuleGroup{{
			Name: networkID,
		}},
	}
}

func (f *File) Rules() []rulefmt.Rule {
	return f.RuleGroups[0].Rules
}

func (f *File) GetRule(rulename string) *rulefmt.Rule {
	for _, rule := range f.RuleGroups[0].Rules {
		if rule.Alert == rulename {
			return &rule
		}
	}
	return nil
}

func (f *File) AddRule(rule rulefmt.Rule) {
	f.RuleGroups[0].Rules = append(f.RuleGroups[0].Rules, rule)
}

func SecureRule(rule rulefmt.Rule, networkID string) error {
	networkLabels := map[string]string{exporters.NetworkLabelInstance: networkID}
	restrictor := security.NewQueryRestrictor(networkLabels)

	restrictedExpression, err := restrictor.RestrictQuery(rule.Expr)
	if err != nil {
		return err
	}
	rule.Expr = restrictedExpression
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
