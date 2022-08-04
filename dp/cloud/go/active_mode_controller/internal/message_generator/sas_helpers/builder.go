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

package sas_helpers

import (
	"encoding/json"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
)

func Build(reqs []*sas.Request) []string {
	byType := [sas.RequestTypeCount][]json.RawMessage{}
	for _, r := range reqs {
		if r != nil {
			byType[r.Type] = append(byType[r.Type], r.Data)
		}
	}
	payloads := make([]string, 0, len(byType))
	for k, v := range byType {
		if len(v) != 0 {
			payloads = append(payloads, toRequest(sas.RequestType(k), v))
		}
	}
	return payloads
}

func toRequest(requestType sas.RequestType, reqs []json.RawMessage) string {
	data := map[string][]json.RawMessage{
		requestType.String(): reqs,
	}
	payload, _ := json.Marshal(data)
	return string(payload)
}
