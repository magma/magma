/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package npmanager

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/collector"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	eventd "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
)

const LteNetwork = "lte"

// NProbeManager provides the main functionality for the nprobe
// service. It collects ES events, encode records and export
// them to a remote collector server.
type NProbeManager struct {
	Collector *collector.EventsCollector
}

// NewNProbeManager creates and returns a new nprobe manager
func NewNProbeManager(
	collector *collector.EventsCollector,
	config nprobe.Config,
) *NProbeManager {
	return &NProbeManager{
		Collector: collector,
	}
}

func getNetworkProbeTasks(networkID string) (map[string]*models.NetworkProbeTask, error) {
	ents, _, err := configurator.LoadAllEntitiesOfType(
		networkID,
		lte.NetworkProbeTaskEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*models.NetworkProbeTask, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&models.NetworkProbeTask{}).FromBackendModels(ent)
	}
	return ret, nil
}

func (np *NProbeManager) collectEvents(
	networkID string,
	targetID string,
	timeSinceReported int64,
) ([]eventd.Event, error) {
	tags := []string{targetID, targetID[4:]}
	events, err := np.Collector.GetMultiStreamsEvents(networkID, timeSinceReported, tags)
	if err != nil {
		return []eventd.Event{}, err
	}
	return events, nil
}

func (np *NProbeManager) processNProbeTask(
	networkID string,
	task *models.NetworkProbeTask,
) error {
	// TBD - get the latest state of the task then
	// process the latest events to create records
	// and export them
	timeSinceReported := int64(0)
	targetID := task.TaskDetails.TargetID
	_, err := np.collectEvents(networkID, targetID, timeSinceReported)
	if err != nil {
		glog.Errorf("Failed to collect events for targetID %s: %s\n", targetID, err)
		return err
	}
	return nil
}

// ProcessNProbeTasks retrieves the list of all nprobe tasks and process them.
// It collects all required events, creates the corresponding records then
// export them to a remote server
func (np *NProbeManager) ProcessNProbeTasks() error {
	networks, err := configurator.ListNetworksOfType(LteNetwork)
	if err != nil {
		glog.Errorf("Failed to retrieve lte network list: %s", err)
		return err
	}

	for _, networkID := range networks {
		tasks, err := getNetworkProbeTasks(networkID)
		if err != nil {
			glog.Errorf("Failed to retrieve nprobe task for network %s: %s", networkID, err)
			continue
		}

		for _, task := range tasks {
			err = np.processNProbeTask(networkID, task)
			if err != nil {
				glog.Errorf("Failed to process events for targetID %s: %s\n", task.TaskDetails.TargetID, err)
				continue
			}
		}
	}
	return nil
}
