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

package sctpd_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/magma/magma/internal/testutil"
	pb "github.com/magma/magma/protos/magma/sctpd"
	"github.com/magma/magma/protos/magma/sctpd/mock_sctpd"
	"github.com/magma/magma/service"
	"github.com/magma/magma/service/sctpd"
)

func TestProxyDownlinkServer_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.InitReq{
		UseIpv4:      true,
		UseIpv6:      true,
		Ipv4Addrs:    []string{"10.0.0.1"},
		Ipv6Addrs:    []string{"fe80::1"},
		Port:         1,
		Ppid:         2,
		ForceRestart: true,
	}
	res := &pb.InitRes{}

	mock_sctpd_downlink := mock_sctpd.NewMockSctpdDownlinkClient(ctrl)
	mock_sctpd_downlink.EXPECT().
		Init(ctx, req).
		Return(res, nil)

	router := service.NewRouter(mock_sctpd_downlink, nil)

	logger, logBuffer := testutil.NewTestLogger()
	pus := sctpd.NewProxyDownlinkServer(logger, router)

	got, err := pus.Init(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tInit\t{\"use_ipv4\": true, \"use_ipv6\": true, \"ipv4_addrs\": [\"10.0.0.1\"], \"ipv6_addrs\": [\"fe80::1\"], \"port\": 1, \"ppid\": 2, \"force_restart\": true}\n",
		logBuffer.String())
}

func TestProxyDownlinkServer_SendDl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.SendDlReq{
		AssocId: 1,
		Stream:  2,
		Payload: []byte{'3'},
		Ppid:    4,
	}
	res := &pb.SendDlRes{}

	mock_sctpd_downlink := mock_sctpd.NewMockSctpdDownlinkClient(ctrl)
	mock_sctpd_downlink.EXPECT().
		SendDl(ctx, req).
		Return(res, nil)

	router := service.NewRouter(mock_sctpd_downlink, nil)

	logger, logBuffer := testutil.NewTestLogger()
	pus := sctpd.NewProxyDownlinkServer(logger, router)

	got, err := pus.SendDl(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tSendDl\t{\"assoc_id\": 1, \"stream\": 2, \"ppid\": 4, \"payload_size\": 1}\n",
		logBuffer.String())
}
