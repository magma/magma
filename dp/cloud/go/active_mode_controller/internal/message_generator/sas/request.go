/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sas

import "encoding/json"

type Request struct {
	Data json.RawMessage
	Type RequestType
}
type RequestType uint8

const (
	Registration RequestType = iota
	SpectrumInquiry
	Grant
	Heartbeat
	Relinquishment
	Deregistration
	RequestTypeCount
)

func (r RequestType) String() string {
	var pref string
	switch r {
	case Registration:
		pref = "registration"
	case SpectrumInquiry:
		pref = "spectrumInquiry"
	case Grant:
		pref = "grant"
	case Heartbeat:
		pref = "heartbeat"
	case Relinquishment:
		pref = "relinquishment"
	case Deregistration:
		pref = "deregistration"
	}
	return pref + "Request"
}

func asRequest(requestType RequestType, data interface{}) *Request {
	b, _ := json.Marshal(data)
	return &Request{
		Type: requestType,
		Data: b,
	}
}
