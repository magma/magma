// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/protos/magma/sctpd/mock_sctpd"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSctpdDownlinkClient := mock_sctpd.NewMockSctpdDownlinkClient(ctrl)
	mockSctpdUplinkClient := mock_sctpd.NewMockSctpdUplinkClient(ctrl)

	router := NewRouter(mockSctpdDownlinkClient, mockSctpdUplinkClient)
	assert.Equal(t, mockSctpdDownlinkClient, router.SctpdDownlinkClient())
	assert.Equal(t, mockSctpdUplinkClient, router.SctpdUplinkClient())
}

func TestRouter_SctpdDownlinkClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSctpdDownlinkClient := mock_sctpd.NewMockSctpdDownlinkClient(ctrl)

	router := NewRouter(mockSctpdDownlinkClient, nil)
	assert.Equal(t, mockSctpdDownlinkClient, router.SctpdDownlinkClient())
}

func TestRouter_SctpdUplinkClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSctpdUplinkClient := mock_sctpd.NewMockSctpdUplinkClient(ctrl)

	router := NewRouter(nil, mockSctpdUplinkClient)
	assert.Equal(t, mockSctpdUplinkClient, router.SctpdUplinkClient())
}
