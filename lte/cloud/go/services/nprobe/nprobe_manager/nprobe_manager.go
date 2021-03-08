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
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/collector"
	"magma/lte/cloud/go/services/nprobe/encoding"
	"magma/lte/cloud/go/services/nprobe/exporter"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/lte/cloud/go/services/nprobe/storage"
	eventd_models "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
)

const LteNetwork = "lte"

// NProbeManager provides the main functionality for the nprobe
// service. It collects ES events, encode IRIs records and export
// them to LIMS.
type NProbeManager struct {
	Exporter                *exporter.RecordExporter
	Collector               *collector.EventsCollector
	StateStore              storage.NProbeStateLookup
	OperatorID              string
	MaxRecordsExportRetries uint32
}

// NewNProbeManager creates and returns a new interceptManager
func NewNProbeManager(
	collector *collector.EventsCollector,
	exporter *exporter.RecordExporter,
	stateStore storage.NProbeStateLookup,
	config nprobe.Config,
) *NProbeManager {
	return &NProbeManager{
		Collector:               collector,
		Exporter:                exporter,
		StateStore:              stateStore,
		MaxRecordsExportRetries: config.MaxRecordsExportRetries,
	}
}

func getNetworkProbeTasks(networkID string) (map[string]*models.NetworkProbeTask, error) {
	ents, err := configurator.LoadAllEntitiesOfType(
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

func (im *NProbeManager) getTimeOfLastExportedRecord(networkID, targetID string, timestamp strfmt.DateTime) (int64, error) {
	state, err := im.StateStore.GetNProbeState(networkID, targetID)
	if err != nil {
		return 0, err
	}

	if state == nil || state.LastExported == 0 {
		return time.Time(timestamp).UnixNano() / int64(time.Millisecond), nil
	}
	return state.LastExported, nil
}

func (im *NProbeManager) updateTimeOfLastExportedRecord(networkID, targetID, timestamp string) error {
	ptime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		glog.Errorf("Failed to parse timestamp %s for targetID %s: %s\n", timestamp, targetID, err)
		return err
	}
	state := &models.NetworkProbeState{
		LastExported: ptime.UnixNano() / int64(time.Millisecond),
	}
	err = im.StateStore.SetNProbeState(networkID, targetID, state)
	if err != nil {
		return err
	}
	return nil
}

func (im *NProbeManager) collectEvents(
	networkID string,
	targetID string,
	startTime int64,
) ([]eventd_models.Event, error) {
	tags := []string{targetID}
	events, err := im.Collector.GetMultiStreamsEvents(networkID, startTime, tags)
	if err != nil {
		return []eventd_models.Event{}, err
	}
	return events, nil
}

func (im *NProbeManager) processEvents(
	networkID string,
	task *models.NetworkProbeTask,
	events []eventd_models.Event,
) error {

	var err error
	timestamp := ""
	for _, event := range events {
		record, err := encoding.MakeRecord(&event, nprobe.IRIRecord,
			im.OperatorID, string(task.TaskID), task.TaskDetails.CorrelationID)
		if err != nil {
			glog.Errorf("Failed to build record for event %s: %s\n", event, err)
			continue
		}
		err = im.Exporter.SendRecord(record, im.MaxRecordsExportRetries)
		if err != nil {
			glog.Errorf("Failed to export record for targetID %s: %s\n", task.TaskDetails.TargetID, err)
			break
		}
		timestamp = event.Timestamp
	}

	if timestamp != "" {
		err = im.updateTimeOfLastExportedRecord(networkID, task.TaskDetails.TargetID, timestamp)
		if err != nil {
			glog.Errorf("Failed to update state for targetID %s: %s\n", task.TaskDetails.TargetID, err)
			return err
		}
	}
	return err
}

// CollectAndProcessEvents collects events, process them and
// sends IRIs reports to an external LIMS
func (im *NProbeManager) CollectAndProcessEvents() error {
	networks, err := configurator.ListNetworksOfType(LteNetwork)
	if err != nil {
		glog.Errorf("Failed to retrieve lte network list: %s", err)
		return err
	}

	for _, networkID := range networks {
		// retrieve configured tasks
		tasks, err := getNetworkProbeTasks(networkID)
		if err != nil {
			glog.Errorf("Failed to retrieve nprobe tasks for network %s: %s", networkID, err)
			continue
		}

		for _, task := range tasks {
			details := task.TaskDetails
			timeSinceReported, err := im.getTimeOfLastExportedRecord(networkID, details.TargetID, details.Timestamp)
			if err != nil {
				glog.Errorf("Failed to retrieve state for targetID %s: %s\n", details.TargetID, err)
				continue
			}

			events, err := im.collectEvents(networkID, details.TargetID, timeSinceReported)
			if err != nil {
				glog.Errorf("Failed to collect events for targetID %s: %s\n", details.TargetID, err)
				continue
			}

			err = im.processEvents(networkID, task, events)
			if err != nil {
				glog.Errorf("Failed to process events for targetID %s: %s\n", details.TargetID, err)
				return err
			}
		}
	}
	return nil
}
