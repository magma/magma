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
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
)

func (c *CallTrace) ToEntity() configurator.NetworkEntity {
	ret := configurator.NetworkEntity{
		Type:   orc8r.CallTraceEntityType,
		Key:    c.Config.TraceID,
		Config: c,
	}
	return ret
}

func (c *CallTrace) FromEntity(ent configurator.NetworkEntity) *CallTrace {
	return ent.Config.(*CallTrace)
}

func (m *CallTrace) FromBackendModels(ent configurator.NetworkEntity) error {
	if ent.Config == nil {
		return fmt.Errorf("could not convert entity to CallTrace; config was nil")
	}
	cfg, ok := ent.Config.(*CallTrace)
	if !ok {
		return fmt.Errorf("could not convert entity config type %T to CallTrace", ent.Config)
	}
	*m = *cfg
	return nil
}

func (c *MutableCallTrace) ToEntityUpdateCriteria(callTraceID string, callTrace CallTrace) configurator.EntityUpdateCriteria {
	update := configurator.EntityUpdateCriteria{
		Type:      orc8r.CallTraceEntityType,
		Key:       callTraceID,
		NewConfig: c.ToCallTrace(callTrace),
	}
	return update
}

func (c *MutableCallTrace) ToCallTrace(callTrace CallTrace) *CallTrace {
	callTrace.State.CallTraceEnding = *c.RequestedEnd
	callTrace.State.CallTraceAvailable = *c.RequestedEnd
	return &callTrace
}
