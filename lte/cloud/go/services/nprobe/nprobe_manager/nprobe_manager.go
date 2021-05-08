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
	"context"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/encoding"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	eventdC "magma/orc8r/cloud/go/services/eventd/eventd_client"
	eventdM "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
)

const (
	LteNetwork = "lte"
	querySize  = 50
)

// NProbeManager provides the main functionality for the nprobe
// service. It collects ES events, encode records and export
// them to a remote collector server.
type NProbeManager struct {
	ElasticClient *elastic.Client
	OperatorID    uint32
}

// NewNProbeManager creates and returns a new nprobe manager
func NewNProbeManager(config nprobe.Config) (*NProbeManager, error) {
	client, err := eventdC.GetElasticClient()
	if err != nil {
		return nil, err
	}
	return &NProbeManager{
		ElasticClient: client,
		OperatorID:    config.OperatorID,
	}, nil
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

func getEvents(networkID, targetID string, start_time *time.Time, client *elastic.Client) ([]eventdM.Event, error) {
	queryParams := eventdC.MultiStreamEventQueryParams{
		NetworkID: networkID,
		Streams:   nprobe.GetESStreams(),
		Events:    nprobe.GetESEventTypes(),
		Tags:      []string{targetID, targetID[4:]},
		Start:     start_time,
		Size:      querySize,
	}
	return eventdC.GetMultiStreamEvents(context.Background(), queryParams, client)
}

func (np *NProbeManager) processNProbeTask(networkID string, task *models.NetworkProbeTask) error {
	// TBD - get the latest state of the task, collect latest events
	// then process them to create iri records.
	targetID := task.TaskDetails.TargetID
	timeSinceReported := time.Time(task.TaskDetails.Timestamp)
	events, err := getEvents(networkID, targetID, &timeSinceReported, np.ElasticClient)
	if err != nil {
		glog.Errorf("Failed to collect events for targetID %s: %s\n", targetID, err)
		return err
	}

	// TBD -  retrieve seq nbr from subscriber state and export records
	var seqNbr uint32 = 0
	for _, event := range events {
		_, err := encoding.MakeRecord(&event, task, np.OperatorID, seqNbr)
		if err != nil {
			glog.Errorf("Failed to collect events for targetID %s: %s\n", targetID, err)
			continue
		}
		seqNbr++
	}
	return nil
}

// ProcessNProbeTasks runs in loop and retrieves all nprobe tasks and process them.
// For each task, it collects latest events, creates the corresponding IRI record then
// export them to a remote destination.
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
				return err
			}
		}
	}
	return nil
}
