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
	"strings"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/files"

	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v2"
)

const (
	rulesFilePostfix = "_rules.yml"
)

// PrometheusAlertClient provides thread-safe methods for writing, reading,
// and modifying alert configuration files
type PrometheusAlertClient interface {
	ValidateRule(rule rulefmt.Rule) error
	RuleExists(networkID, rulename string) bool
	WriteRule(networkID string, rule rulefmt.Rule) error
	UpdateRule(networkID string, rule rulefmt.Rule) error
	ReadRules(networkID, ruleName string) ([]rulefmt.Rule, error)
	DeleteRule(networkID, ruleName string) error
	BulkUpdateRules(networkID string, rules []rulefmt.Rule) (BulkUpdateResults, error)
}

type client struct {
	fileLocks *FileLocker
	rulesDir  string
	fsClient  files.FSClient
}

func NewClient(fileLocks *FileLocker, rulesDir string, fsClient files.FSClient) PrometheusAlertClient {
	return &client{
		fileLocks: fileLocks,
		rulesDir:  rulesDir,
		fsClient:  fsClient,
	}
}

// ValidateRule checks that a new alert rule is a valid specification
func (c *client) ValidateRule(rule rulefmt.Rule) error {
	errs := rule.Validate()
	if len(errs) != 0 {
		return fmt.Errorf("invalid rule: %v", errs)
	}
	return nil
}

func (c *client) RuleExists(networkID, rulename string) bool {
	filename := makeFilename(networkID, c.rulesDir)

	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.initializeRuleFile(networkID, filename)
	if err != nil {
		return false
	}
	return ruleFile.GetRule(rulename) != nil
}

// WriteRule takes an alerting rule and writes it to the rules file for the
// given networkID
func (c *client) WriteRule(networkID string, rule rulefmt.Rule) error {
	filename := makeFilename(networkID, c.rulesDir)

	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.initializeRuleFile(networkID, filename)
	if err != nil {
		return err
	}
	ruleFile.AddRule(rule)

	err = c.writeRuleFile(ruleFile, filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) UpdateRule(networkID string, rule rulefmt.Rule) error {
	filename := makeFilename(networkID, c.rulesDir)

	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.initializeRuleFile(networkID, filename)
	if err != nil {
		return err
	}

	err = SecureRule(networkID, &rule)
	if err != nil {
		return err
	}

	err = ruleFile.ReplaceRule(rule)
	if err != nil {
		return err
	}

	err = c.writeRuleFile(ruleFile, filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) ReadRules(networkID, ruleName string) ([]rulefmt.Rule, error) {
	filename := makeFilename(networkID, c.rulesDir)
	c.fileLocks.RLock(filename)
	defer c.fileLocks.RUnlock(filename)

	if !c.ruleFileExists(filename) {
		return []rulefmt.Rule{}, nil
	}

	ruleFile, err := c.readRuleFile(makeFilename(networkID, c.rulesDir))
	if err != nil {
		return []rulefmt.Rule{}, err
	}
	if ruleName == "" {
		return ruleFile.Rules(), nil
	}
	foundRule := ruleFile.GetRule(ruleName)
	if foundRule == nil {
		return nil, fmt.Errorf("rule %s not found", ruleName)
	}
	return []rulefmt.Rule{*foundRule}, nil
}

func (c *client) DeleteRule(networkID, ruleName string) error {
	filename := makeFilename(networkID, c.rulesDir)
	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.readRuleFile(filename)
	if err != nil {
		return err
	}

	err = ruleFile.DeleteRule(ruleName)
	if err != nil {
		return err
	}

	err = c.writeRuleFile(ruleFile, filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) BulkUpdateRules(networkID string, rules []rulefmt.Rule) (BulkUpdateResults, error) {
	filename := makeFilename(networkID, c.rulesDir)
	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.readRuleFile(filename)
	if err != nil {
		return BulkUpdateResults{}, err
	}

	results := NewBulkUpdateResults()
	for _, newRule := range rules {
		ruleName := newRule.Alert
		err := SecureRule(networkID, &newRule)
		if err != nil {
			results.Errors[ruleName] = err
			continue
		}
		if ruleFile.GetRule(ruleName) != nil {
			err := ruleFile.ReplaceRule(newRule)
			if err != nil {
				results.Errors[ruleName] = err
			} else {
				results.Statuses[ruleName] = "updated"
			}
		} else {
			ruleFile.AddRule(newRule)
			results.Statuses[ruleName] = "created"
		}
	}

	err = c.writeRuleFile(ruleFile, filename)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (c *client) writeRuleFile(ruleFile *File, filename string) error {
	yamlFile, err := yaml.Marshal(ruleFile)
	err = c.fsClient.WriteFile(filename, yamlFile, 0666)
	if err != nil {
		return fmt.Errorf("error writing rules file: %v\n", yamlFile)
	}
	return nil
}

func (c *client) initializeRuleFile(networkID, filename string) (*File, error) {
	if _, err := c.fsClient.Stat(filename); err == nil {
		file, err := c.readRuleFile(filename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return NewFile(networkID), nil
}

func (c *client) ruleFileExists(filename string) bool {
	_, err := c.fsClient.Stat(filename)
	return err == nil
}

func (c *client) readRuleFile(requestedFile string) (*File, error) {
	ruleFile := File{}
	file, err := c.fsClient.ReadFile(requestedFile)
	if err != nil {
		return &File{}, fmt.Errorf("error reading rules files: %v", err)
	}
	err = yaml.Unmarshal(file, &ruleFile)
	return &ruleFile, err
}

type BulkUpdateResults struct {
	Errors   map[string]error
	Statuses map[string]string
}

func NewBulkUpdateResults() BulkUpdateResults {
	return BulkUpdateResults{
		Errors:   make(map[string]error, 0),
		Statuses: make(map[string]string, 0),
	}
}

func (r BulkUpdateResults) String() string {
	str := strings.Builder{}
	if len(r.Errors) > 0 {
		str.WriteString("Errors: \n")
		for name, err := range r.Errors {
			str.WriteString(fmt.Sprintf("\t%s: %s\n", name, err))
		}
	}
	if len(r.Statuses) > 0 {
		str.WriteString("Statuses: \n")
		for name, status := range r.Statuses {
			str.WriteString(fmt.Sprintf("\t%s: %s\n", name, status))
		}
	}
	return str.String()
}

func makeFilename(networkID, path string) string {
	return path + "/" + networkID + rulesFilePostfix
}
