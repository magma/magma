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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/encoding/np_interface"
	"magma/lte/cloud/go/services/nprobe/encoding/x2_interface"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
)

// Encoder represtents the exported message
type Encoder interface {
	// Encode encodes this element by writing Len() bytes to dst.
	Encode() ([]byte, error)
}

func MakeField(event models.Event, format string, opID, xID string, correlationID int64) (Encoder, error) {
	switch format {
	case nprobe.IRIRecord:
		record, err := x2_interface.ConstructEpsIRIMessage(event, opID, xID, correlationID)
		if err != nil {
			glog.Errorf("Failed to construct IRI record %v\n", event)
			return &X2Encoder{}, err
		}
		return &X2Encoder{record: record}, nil
	case nprobe.NProbeRecord:
		record, err := np_interface.ConstructNProbeRecord(event, xID, correlationID)
		if err != nil {
			glog.Errorf("Failed to construct NProbe record %v\n", event)
			return &NProbeEncoder{}, err
		}
		return &NProbeEncoder{record: record}, nil
	}
	return nil, fmt.Errorf("Unsupported encoding format %s", format)
}

// X2Encoder represents an encoder for Json format used by Lawful Interception
type X2Encoder struct {
	record x2_interface.EpsIRIMessage
}

func (e *X2Encoder) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e.record)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NProbeEncoder represents an encoder for Json format used by Network Probe
type NProbeEncoder struct {
	record np_interface.NProbeMessage
}

func (e *NProbeEncoder) Encode() ([]byte, error) {
	return json.Marshal(e.record)
}
