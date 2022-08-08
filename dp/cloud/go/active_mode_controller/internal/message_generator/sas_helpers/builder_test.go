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

package sas_helpers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas_helpers"
)

func TestBuild(t *testing.T) {
	const someDeregistrationRequest = `{"cbsdId":"someId"}`
	const otherDeregistrationRequest = `{"cbsdId":"otherId"}`
	const someHeartbeatRequest = `{"cbsdId":"someId","grantId":"grantId"}`
	const someRegistrationRequest = `{"key":"value"}`
	requests := []*sas.Request{{
		Type: sas.Deregistration,
		Data: []byte(someDeregistrationRequest),
	}, {
		Type: sas.Heartbeat,
		Data: []byte(someHeartbeatRequest),
	}, {
		Type: sas.Deregistration,
		Data: []byte(otherDeregistrationRequest),
	}, {
		Type: sas.Registration,
		Data: []byte(someRegistrationRequest),
	}}
	actual := sas_helpers.Build(requests)
	expected := []string{
		fmt.Sprintf(`{"%s":[%s]}`, sas.Registration, someRegistrationRequest),
		fmt.Sprintf(`{"%s":[%s]}`, sas.Heartbeat, someHeartbeatRequest),
		fmt.Sprintf(`{"%s":[%s,%s]}`, sas.Deregistration,
			someDeregistrationRequest, otherDeregistrationRequest),
	}
	assert.Equal(t, expected, actual)
}

func TestSkipNil(t *testing.T) {
	requests := []*sas.Request{nil, nil}
	actual := sas_helpers.Build(requests)
	assert.Empty(t, actual)
}
