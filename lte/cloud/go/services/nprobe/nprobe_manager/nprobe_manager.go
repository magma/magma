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
  "fmt"

  "magma/lte/cloud/go/lte"
  "magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/exporter"
	"magma/lte/cloud/go/services/nprobe/collector"
  "magma/lte/cloud/go/services/nprobe/obsidian/models"
  "magma/orc8r/cloud/go/services/configurator"

  "github.com/golang/glog"
)

const LteNetwork = "lte"

// NetworkProbeManager provides the main functionality for the nprobe
// service. It collects ES events, encode IRIs records and export
// them to LIMS.
type NetworkProbeManager struct {
	Collector               *collector.EventsCollector
	Exporter                *exporter.RecordExporter
	MaxEventsCollectRetries uint32
	MaxRecordsExportRetries uint32
}

// NewNetworkProbeManager creates and returns a new interceptManager
func NewNetworkProbeManager(
	collector *collector.EventsCollector,
	exporter *exporter.RecordExporter,
	config nprobe.Config,
) *NetworkProbeManager {
	return &NetworkProbeManager{
		Collector:               collector,
		Exporter:                exporter,
		MaxEventsCollectRetries: config.MaxEventsCollectRetries,
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

// CollectAndProcessEvents collects events, process them and
// sends IRIs reports to an external LIMS
func (im *NetworkProbeManager) CollectAndProcessEvents() error {
  networks, err := configurator.ListNetworksOfType(LteNetwork)
  if err != nil {
    glog.Errorf("Error retrieving lte network list: %s", err)
    return err
  }

  for _, networkID := range networks {
    tasks, err := getNetworkProbeTasks(networkID)
    if err != nil {
      glog.Errorf("Error retrieving network probe tasks for network %s: %s", networkID, err)
      continue
    }

    for _, task := range tasks {
	    // TBD
	    // Encode records
	    // Export records to Remote server
      _, err := im.Collector.GetMultiStreamsEvents(networkID, "", []string{task.TaskDetails.TargetID})
      if err != nil {
        fmt.Printf("Error while retrieving events for subscriber %v: %s\n", task.TaskDetails.TargetID, err)
        return err
      }
    }
  }
	return nil
}
