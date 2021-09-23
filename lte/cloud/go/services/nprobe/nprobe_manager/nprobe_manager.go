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
	"fmt"
	"time"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/encoding"
	"magma/lte/cloud/go/services/nprobe/exporter"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/lte/cloud/go/services/nprobe/storage"
	"magma/orc8r/cloud/go/services/configurator"
	eventdC "magma/orc8r/cloud/go/services/eventd/eventd_client"
	eventdM "magma/orc8r/cloud/go/services/eventd/obsidian/models"
)

const (
	LteNetwork = "lte"
	querySize  = 50
)

// NProbeManager provides the main functionality for the nprobe
// service. It collects ES events, encode records and export
// them to a remote collector server.
type NProbeManager struct {
	ElasticClient    *elastic.Client
	Storage          storage.NProbeStorage
	Exporters        map[string]*exporter.RecordExporter
	MaxExportRetries uint32
}

// NewNProbeManager creates and returns a new nprobe manager
func NewNProbeManager(
	config nprobe.Config,
	storage storage.NProbeStorage,
) (*NProbeManager, error) {
	client, err := eventdC.GetElasticClient()
	if err != nil {
		return nil, err
	}
	exporters := make(map[string]*exporter.RecordExporter)
	return &NProbeManager{
		ElasticClient:    client,
		Storage:          storage,
		Exporters:        exporters,
		MaxExportRetries: config.MaxExportRetries,
	}, nil
}

// getNetworkProbeTasks retrieves the list of all tasks provisioned for a specific network
func getNetworkProbeTasks(networkID string) (map[string]*models.NetworkProbeTask, error) {
	ents, _, err := configurator.LoadAllEntitiesOfType(
		context.Background(),
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

// createNewRecordExporter retrieves the first destination provisioned
// for a specific network and creates a record exporter
func createNewRecordExporter(networkID string) (*exporter.RecordExporter, error) {
	ents, _, err := configurator.LoadAllEntitiesOfType(
		context.Background(),
		networkID,
		lte.NetworkProbeDestinationEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}

	if len(ents) == 0 {
		return nil, fmt.Errorf("Could not find a destination for network  %v", networkID)
	}

	destination := (&models.NetworkProbeDestination{}).FromBackendModels(ents[0])
	return exporter.NewRecordExporter(destination)
}

// getEvents retrieves all events since start_time from fluentd
func getEvents(
	networkID string,
	state *models.NetworkProbeData,
	client *elastic.Client,
) ([]eventdM.Event, error) {

	// build multi-stream es query
	targetID := state.TargetID
	startTime := time.Time(state.LastExported).Add(time.Millisecond * 1)
	queryParams := eventdC.MultiStreamEventQueryParams{
		NetworkID: networkID,
		Streams:   nprobe.GetESStreams(),
		Events:    nprobe.GetESEventTypes(),
		Tags:      []string{targetID, targetID[4:]},
		Start:     &startTime,
		Size:      querySize,
	}

	return eventdC.GetMultiStreamEvents(context.Background(), queryParams, client)
}

// updateRecordState updates nprobe state with last sequence number and timestamp
func (np *NProbeManager) updateRecordState(
	networkID, taskID string,
	state models.NetworkProbeData,
	timestamp string,
	sequenceNumber uint32,
) error {
	ptime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}

	// update state with last timestamp and sequence nbr
	state.LastExported = strfmt.DateTime(ptime)
	state.SequenceNumber = sequenceNumber
	return np.Storage.StoreNProbeData(networkID, taskID, state)
}

// processNProbeTask is the main function processing each task, managing state and exporting data
func (np *NProbeManager) processNProbeTask(networkID string, task *models.NetworkProbeTask) error {
	taskID := string(task.TaskID)
	state, err := np.Storage.GetNProbeData(networkID, taskID)
	if err != nil {
		glog.Errorf("Failed to get state for record %s: %v", taskID, err)
		return err
	}

	events, err := getEvents(networkID, state, np.ElasticClient)
	if err != nil {
		glog.Errorf("Failed to collect events for targetID %s: %s\n", state.TargetID, err)
		return err
	}

	var nerr error
	seq := state.SequenceNumber
	for _, event := range events {
		record, err := encoding.MakeRecord(&event, task, seq)
		if err != nil {
			glog.Errorf("Failed to build record from event %v: %s\n", event, err)
			continue
		}

		nerr = np.Exporters[networkID].SendMessageWithRetries(record, np.MaxExportRetries)
		if nerr != nil {
			glog.Errorf("Failed to export record for targetID %s: %s\n", state.TargetID, nerr)
			np.Exporters[networkID] = nil
			break
		}
		seq++
	}

	if seq > state.SequenceNumber {
		idx := seq - state.SequenceNumber - 1
		err = np.updateRecordState(networkID, taskID, *state, events[idx].Timestamp, seq)
		if err != nil {
			glog.Errorf("Failed to update state for targetID %s: %s\n", state.TargetID, err)
			return err
		}
	}
	return nerr
}

// ProcessNProbeTasks runs in loop, retrieves all nprobe tasks and process them.
// For each task, it collects latest events, creates the corresponding IRI record then
// export them to a remote destination.
func (np *NProbeManager) ProcessNProbeTasks() error {
	networks, err := configurator.ListNetworksOfType(context.Background(), LteNetwork)
	if err != nil {
		glog.Errorf("Failed to retrieve lte network list: %s", err)
		return err
	}

	for _, networkID := range networks {
		if np.Exporters[networkID] == nil {
			exporter, err := createNewRecordExporter(networkID)
			if err != nil || exporter == nil {
				glog.Infof("Could not create an exporter for network %s: %s", networkID, err)
				continue
			}
			np.Exporters[networkID] = exporter
		}

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
