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

package np_interface

import (
	"time"

	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
)

func ConstructNProbeRecord(event models.Event, correlationID, xID string) (NProbeMessage, error) {
	generalizedTime, err := time.Parse(time.RFC3339Nano, event.Timestamp)
	if err != nil {
		glog.Errorf("Failed to parse timestamp %s.", event.Timestamp)
		return NProbeMessage{}, err
	}

	record := NProbeMessage{
		XID:           xID,
		Timestamp:     generalizedTime,
		MatchedTarget: event.Tag,
		CorrelationID: correlationID,
	}
	return record, nil
}
