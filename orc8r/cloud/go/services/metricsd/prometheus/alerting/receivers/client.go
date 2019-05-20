/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// Client provides methods to create and read receiver configurations
type Client struct {
	configPath string
	sync.RWMutex
}

func NewClient(configPath string) *Client {
	return &Client{
		configPath: configPath,
	}
}

// CreateReceiver writes a new receiver to the config file with the networkID
// prepended to the name so multiple networks can be supported
func (c *Client) CreateReceiver(rec *Receiver, networkID string) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	rec.Secure(networkID)
	conf.Receivers = append(conf.Receivers, rec)
	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetReceivers returns the receiver configs for the given networkID
func (c *Client) GetReceivers(networkID string) ([]Receiver, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return []Receiver{}, nil
	}

	recs := make([]Receiver, 0)
	for _, rec := range conf.Receivers {
		if strings.HasPrefix(rec.Name, receiverNetworkPrefix(networkID)) {
			rec.Unsecure(networkID)
			recs = append(recs, *rec)
		}
	}
	return recs, nil
}

func (c *Client) readConfigFile() (*Config, error) {
	configFile := Config{}
	file, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config files: %v", err)
	}
	err = yaml.Unmarshal(file, &configFile)
	return &configFile, err
}

func (c *Client) writeConfigFile(conf *Config) error {
	yamlFile, err := yaml.Marshal(conf)
	err = ioutil.WriteFile(c.configPath, yamlFile, 0660)
	if err != nil {
		return fmt.Errorf("error writing config file: %v\n", yamlFile)
	}
	return nil
}

func receiverNetworkPrefix(networkID string) string {
	return strings.Replace(networkID, "_", "", -1) + "_"
}
