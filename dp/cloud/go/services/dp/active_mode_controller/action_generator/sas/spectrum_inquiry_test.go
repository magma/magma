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

package sas_test

import (
	"testing"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestSpectrumInquiryRequestGenerator(t *testing.T) {
	cbsd := &storage.DBCbsd{CbsdId: db.MakeString("some_id")}
	data := &storage.DetailedCbsd{Cbsd: cbsd}
	g := &sas.SpectrumInquiryRequestGenerator{}
	actual := g.GenerateRequests(data)
	expected := []*request{{
		requestType: "spectrumInquiryRequest",
		data: `{
	"cbsdId": "some_id",
	"inquiredSpectrum": [{
		"lowFrequency": 3550000000,
		"highFrequency": 3700000000
	}]
}`,
	}}
	assertRequestsEqual(t, expected, actual)
}
