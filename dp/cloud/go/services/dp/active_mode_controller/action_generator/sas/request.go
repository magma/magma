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

import (
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

const (
	Registration    = "registrationRequest"
	SpectrumInquiry = "spectrumInquiryRequest"
	Grant           = "grantRequest"
	Heartbeat       = "heartbeatRequest"
	Relinquishment  = "relinquishmentRequest"
	Deregistration  = "deregistrationRequest"
)

func makeRequest(requestType string, data any) *storage.MutableRequest {
	return &storage.MutableRequest{
		Request: &storage.DBRequest{
			Payload: data,
		},
		RequestType: &storage.DBRequestType{
			Name: db.MakeString(requestType),
		},
	}
}
