/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package testlib

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"magma/orc8r/cloud/go/parallel"
)

var clusterConfig = flag.String("cluster-config", "", "Location of the cluster config file")

// Gateway gateway specific configuration of the cluster
type Gateway struct {
	ID         string `json:"gateway_id"`
	HardwareID string `json:"hardware_id"`
	HostName   string `json:"hostname"`
}

type Gateways []Gateway

func (g Gateways) Hostnames() []string {
	var hh []string
	for _, gw := range g {
		hh = append(hh, gw.HostName)
	}
	return hh
}

// ClusterInternalConfig cluster's internal configuration
type ClusterInternalConfig struct {
	BastionIP string `json:"bastion_ip"`
}

// Cluster configuration of the cluster
type Cluster struct {
	UUID           string                 `json:"uuid"`
	ClusterType    int                    `json:"cluster_type"`
	InternalConfig ClusterInternalConfig  `json:"internal_config"`
	Template       map[string]interface{} `json:"template"`
	Gateways       Gateways               `json:"gateways"`
}

func (c *Cluster) RunCmdOnGateways(cmd string) ([]string, error) {
	outs, err := parallel.MapString(c.Gateways.Hostnames(), parallel.DefaultNumWorkers, func(in parallel.In) (parallel.Out, error) {
		hostname := in.(string)
		out, err := runRemoteCommand(hostname, c.InternalConfig.BastionIP, []string{cmd})
		if err != nil {
			return nil, fmt.Errorf("run remote commands on gateway '%s'; out: %+v: %w", hostname, out, err)
		}
		return out, nil
	})
	return outs, err
}

// GetClusterInfo parses cluster config file and returns config instance
func GetClusterInfo() (*Cluster, error) {
	c, err := GetClusterFromFile()
	if err != nil {
		return nil, err
	}
	err = c.populateHWIDs()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GetClusterFromFile() (*Cluster, error) {
	if clusterConfig == nil {
		return nil, fmt.Errorf("Cluster config doesn't exist")
	}

	bytes, err := ioutil.ReadFile(*clusterConfig)
	if err != nil {
		return nil, err
	}
	cluster := Cluster{}
	_ = json.Unmarshal([]byte(bytes), &cluster)
	return &cluster, nil
}

// SetClusterInfo saves cluster information to cluster config
func SetClusterInfo(cluster *Cluster) error {
	if clusterConfig == nil {
		return fmt.Errorf("Cluster config doesn't exist")
	}
	file, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("Failed marshalling payload")
		return err
	}
	err = ioutil.WriteFile(*clusterConfig, file, 0644)
	return err
}

func (c *Cluster) populateHWIDs() error {
	hwids, err := parallel.MapString(c.Gateways.Hostnames(), parallel.DefaultNumWorkers, func(in parallel.In) (parallel.Out, error) {
		hostname := in.(string)
		hwid, err := fetchHardwareID(hostname, c.InternalConfig.BastionIP)
		if err != nil {
			return nil, err
		}
		return hwid, nil
	})
	if err != nil {
		return err
	}

	for i, hwid := range hwids {
		c.Gateways[i].HardwareID = hwid
	}

	return nil
}
