/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package alert_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient/mocks"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testNID      = "test"
	testRuleFile = `groups:
- name: test
  rules:
  - alert: test_rule_1
    expr: up == 0{tenantID="test"}
    for: 5s
    labels:
      severity: major
      tenantID: test
  - alert: test_rule_2
    expr: up == 1{tenantID="test"}
    for: 5s
    labels:
      severity: critical
      tenantID: test
    annotations:
      summary: A test rule`

	otherNID      = "other"
	otherRuleFile = `groups:
- name: other
  rules:
  - alert: other_rule_1
    expr: up == 0{tenantID="other"}
    for: 5s
    labels:
      severity: major
      tenantID: other
  - alert: test_rule_2
    expr: up == 1{tenantID="other"}
    for: 5s
    labels:
      severity: critical
      tenantID: other
    annotations:
      summary: A test rule`
)

var (
	fiveSeconds, _ = model.ParseDuration("5s")
	testRule1      = rulefmt.Rule{
		Alert:  "test_rule_1",
		Expr:   "up==0",
		For:    fiveSeconds,
		Labels: map[string]string{"severity": "major", "tenantID": testNID},
	}
	badRule = rulefmt.Rule{
		Alert: "bad_rule",
		Expr:  "malformed{.}",
	}
)

func TestClient_ValidateRule(t *testing.T) {
	client := newTestClient("tenantID")

	err := client.ValidateRule(sampleRule)
	assert.NoError(t, err)

	invalidRule := rulefmt.Rule{
		// Only one of Record/Alert can be set
		Record: "x",
		Alert:  "x",
	}
	err = client.ValidateRule(invalidRule)
	assert.Error(t, err)
}
func TestClient_RuleExists(t *testing.T) {
	client := newTestClient("tenantID")
	assert.True(t, client.RuleExists(testNID, "test_rule_1"))
	assert.True(t, client.RuleExists(testNID, "test_rule_2"))
	assert.False(t, client.RuleExists(testNID, "no_rule"))
	assert.False(t, client.RuleExists(testNID, "other_rule_1"))

	assert.True(t, client.RuleExists(otherNID, "other_rule_1"))
	assert.True(t, client.RuleExists(otherNID, "test_rule_2"))
	assert.False(t, client.RuleExists(otherNID, "no_rule"))
	assert.False(t, client.RuleExists(otherNID, "test_rule_1"))
}

func TestClient_WriteRule(t *testing.T) {
	client := newTestClient("tenantID")
	err := client.WriteRule(testNID, sampleRule)
	assert.NoError(t, err)
}

func TestClient_UpdateRule(t *testing.T) {
	client := newTestClient("tenantID")

	err := client.UpdateRule(testNID, testRule1)
	assert.NoError(t, err)

	// Returns error when updating non-existent rule
	err = client.UpdateRule(testNID, sampleRule)
	assert.Error(t, err)
}

func TestClient_ReadRules(t *testing.T) {
	client := newTestClient("tenantID")

	rules, err := client.ReadRules(testNID, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rules))
	assert.Equal(t, "test_rule_1", rules[0].Alert)
	assert.Equal(t, "test_rule_2", rules[1].Alert)

	rules, err = client.ReadRules(otherNID, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rules))
	assert.Equal(t, "other_rule_1", rules[0].Alert)
	assert.Equal(t, "test_rule_2", rules[1].Alert)

	rules, err = client.ReadRules(testNID, "test_rule_1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "test_rule_1", rules[0].Alert)

	rules, err = client.ReadRules(testNID, "no_rule")
	assert.Error(t, err)
	assert.Equal(t, 0, len(rules))
}

func TestClient_DeleteRule(t *testing.T) {
	client := newTestClient("tenantID")
	err := client.DeleteRule(testNID, "test_rule_1")
	assert.NoError(t, err)

	err = client.DeleteRule(testNID, "no_rule")
	assert.Error(t, err)
}

func TestClient_BulkUpdateRules(t *testing.T) {
	client := newTestClient("tenantID")
	results, err := client.BulkUpdateRules(testNID, []rulefmt.Rule{sampleRule, testRule1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results.Statuses))
	assert.Equal(t, 0, len(results.Errors))

	results, err = client.BulkUpdateRules(testNID, []rulefmt.Rule{badRule, sampleRule, testRule1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results.Statuses))
	assert.Equal(t, 1, len(results.Errors))
}

func newTestClient(multitenantLabel string) alert.PrometheusAlertClient {
	dClient := newHealthyDirClient("test")
	fileLocks, _ := alert.NewFileLocker(dClient)
	fsClient := &mocks.FSClient{}
	fsClient.On("Stat", mock.AnythingOfType("string")).Return(nil, nil)
	fsClient.On("ReadFile", "test_rules/test_rules.yml").Return([]byte(testRuleFile), nil)
	fsClient.On("ReadFile", "test_rules/other_rules.yml").Return([]byte(otherRuleFile), nil)
	fsClient.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	tenancy := alert.TenancyConfig{
		RestrictorLabel: multitenantLabel,
		RestrictQueries: true,
	}
	return alert.NewClient(fileLocks, "test_rules", "prometheus-host.com", fsClient, tenancy)
}
