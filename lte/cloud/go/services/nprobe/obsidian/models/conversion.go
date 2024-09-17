/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"magma/lte/cloud/go/lte"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
)

func (m *NetworkProbeTask) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:      lte.NetworkProbeTaskEntityType,
		Key:       string(m.TaskID),
		NewConfig: m.TaskDetails,
	}
}

func (m *NetworkProbeTask) FromBackendModels(ent configurator.NetworkEntity) *NetworkProbeTask {
	m.TaskID = NetworkProbeTaskID(ent.Key)
	m.TaskDetails = ent.Config.(*NetworkProbeTaskDetails)
	return m
}

func (m *NetworkProbeDestination) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:      lte.NetworkProbeDestinationEntityType,
		Key:       string(m.DestinationID),
		NewConfig: m.DestinationDetails,
	}
}

func (m *NetworkProbeDestination) FromBackendModels(ent configurator.NetworkEntity) *NetworkProbeDestination {
	m.DestinationID = NetworkProbeDestinationID(ent.Key)
	m.DestinationDetails = ent.Config.(*NetworkProbeDestinationDetails)
	return m
}

func ToMConfigNProbeTask(task *NetworkProbeTask) *lte_mconfig.NProbeTask {
	return &lte_mconfig.NProbeTask{
		TaskId:        string(task.TaskID),
		DomainId:      task.TaskDetails.DomainID,
		TargetId:      task.TaskDetails.TargetID,
		TargetType:    task.TaskDetails.TargetType,
		DeliveryType:  task.TaskDetails.DeliveryType,
		CorrelationId: task.TaskDetails.CorrelationID,
	}
}
