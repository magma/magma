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
	"io/ioutil"
	"net/http"
	"os"

	"magma/orc8r/cloud/go/obsidian/handlers"

	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v2"
)

const (
	rulesFilePostfix = "_rules.yml"
)

// Client provides thread-safe methods for writing, reading, and modifying
// alert configuration files
type Client struct {
	fileLocks *FileLocker
	rulesDir  string
}

func NewClient(rulesDir string) (*Client, error) {
	fileLocks, err := NewFileLocker(rulesDir)
	if err != nil {
		return nil, err
	}
	return &Client{
		fileLocks: fileLocks,
		rulesDir:  rulesDir,
	}, nil
}

// WriteAlert takes an alerting rule and writes it to the rules file for the
// given networkID
func (c *Client) WriteAlert(rule rulefmt.Rule, networkID string) error {
	errs := rule.Validate()
	if len(errs) != 0 {
		return handlers.HttpError(fmt.Errorf("invalid rule: %v", errs), http.StatusBadRequest)
	}
	filename := makeFilename(networkID, c.rulesDir)

	c.fileLocks.Lock(filename)
	defer c.fileLocks.Unlock(filename)

	ruleFile, err := c.initializeRuleFile(filename, networkID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	ruleFile.AddRule(rule)

	err = c.writeRuleFile(ruleFile, filename)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return nil
}

func (c *Client) ReadRules(ruleName string, networkID string) ([]rulefmt.Rule, error) {
	filename := makeFilename(networkID, c.rulesDir)
	c.fileLocks.RLock(filename)
	defer c.fileLocks.RUnlock(filename)

	ruleFile, err := c.readRuleFile(makeFilename(networkID, c.rulesDir))
	if err != nil {
		return []rulefmt.Rule{}, err
	}
	if ruleName == "" {
		return ruleFile.Rules(), nil
	}
	foundRule, err := ruleFile.GetRule(ruleName)
	if err != nil {
		return nil, err
	}
	return []rulefmt.Rule{*foundRule}, nil
}

func (c *Client) writeRuleFile(ruleFile *File, filename string) error {
	yamlFile, err := yaml.Marshal(ruleFile)
	err = ioutil.WriteFile(filename, yamlFile, 0660)
	if err != nil {
		return fmt.Errorf("error writing rules file: %v\n", yamlFile)
	}
	return nil
}

func (c *Client) initializeRuleFile(filename, networkID string) (*File, error) {
	if _, err := os.Stat(filename); err == nil {
		file, err := c.readRuleFile(filename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return NewFile(networkID), nil
}

func (c *Client) readRuleFile(requestedFile string) (*File, error) {
	ruleFile := File{}
	file, err := ioutil.ReadFile(requestedFile)
	if err != nil {
		return &File{}, fmt.Errorf("error reading rules files: %v", err)
	}
	err = yaml.Unmarshal(file, &ruleFile)
	return &ruleFile, err
}

func makeFilename(networkID, path string) string {
	return path + "/" + networkID + rulesFilePostfix
}
