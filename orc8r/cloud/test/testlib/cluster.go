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
)

var clusterConfig = flag.String("cluster-config", "", "Location of the cluster config file")

//Gateway gateway specific configuration of the cluster
type Gateway struct {
	ID         string `json:"gateway_id"`
	HardwareID string `json:"hardware_id"`
	HostName   string `json:"hostname"`
}

//ClusterInternalConfig cluster's internal configuration
type ClusterInternalConfig struct {
	BastionIP string `json:"bastion_ip"`
}

//Cluster configuration of the cluster
type Cluster struct {
	UUID           string                 `json:"uuid"`
	ClusterType    int                    `json:"cluster_type"`
	InternalConfig ClusterInternalConfig  `json:"internal_config"`
	Template       map[string]interface{} `json:"template"`
	Gateways       []Gateway              `json:"gateways"`
}

// GetClusterInfo parses cluster config file and returns config instance
func GetClusterInfo() (*Cluster, error) {
	if clusterConfig == nil {
		return nil, fmt.Errorf("Cluster config doesn't exist")
	}

	bytes, err := ioutil.ReadFile(*clusterConfig)
	if err != nil {
		return nil, err
	}
	cluster := Cluster{}
	_ = json.Unmarshal([]byte(bytes), &cluster)

	// populate gateway information
	for i := range cluster.Gateways {
		hardwareID, err := fetchHardwareID(
			cluster.Gateways[i].HostName,
			cluster.InternalConfig.BastionIP)
		if err != nil {
			return nil, err
		}
		cluster.Gateways[i].HardwareID = hardwareID
	}
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
