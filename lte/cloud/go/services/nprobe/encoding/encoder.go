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

package encoding

import (
	"encoding/json"

	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/encoding/np_interface"
	"magma/lte/cloud/go/services/nprobe/encoding/x2_interface"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	eventd_models "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
)

// Encoder represtents the exported message
type Encoder interface {
	// Encode encodes this element by writing Len() bytes to dst.
	Encode() ([]byte, error)
}

// X2Encoder represents an encoder for Json format used by Lawful Interception
type X2Encoder struct {
	record x2_interface.EpsIRIRecord
}

func (e *X2Encoder) Encode() ([]byte, error) {
	msg := e.record.Header.Serialize()
	msg = append(msg, e.record.Payload...)
	return msg, nil
}

// NProbeEncoder represents an encoder for Json format used by Network Probe
type NProbeEncoder struct {
	record np_interface.NProbeMessage
}

func (e *NProbeEncoder) Encode() ([]byte, error) {
	return json.Marshal(e.record)
}

func MakeRecord(event *eventd_models.Event, format, operatorID string, task *models.NetworkProbeTask, seqNbr uint64) ([]byte, error) {
	var encoder Encoder
	switch format {
	case nprobe.IRIRecord:
		record, err := x2_interface.MakeRecord(event, operatorID, task, seqNbr)
		if err != nil {
			glog.Errorf("Failed to construct IRI record %v\n", event)
			return []byte{}, err
		}
		encoder = &X2Encoder{record: record}
	case nprobe.NProbeRecord:
		record, err := np_interface.MakeRecord(event, operatorID, string(task.TaskID), task.TaskDetails.CorrelationID)
		if err != nil {
			glog.Errorf("Failed to construct NProbe record %v\n", event)
			return []byte{}, err
		}
		encoder = &NProbeEncoder{record: record}
	}
	return encoder.Encode()
}
