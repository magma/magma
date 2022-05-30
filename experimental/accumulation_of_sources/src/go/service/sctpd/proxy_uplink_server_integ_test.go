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
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/magma/magma/src/go/internal/testutil"
	pb "github.com/magma/magma/src/go/protos/magma/sctpd"
	"github.com/magma/magma/src/go/protos/magma/sctpd/mock_sctpd"
	"github.com/magma/magma/src/go/service"
	"github.com/magma/magma/src/go/service/sctpd"
)

func TestProxyUplinkServer_SendUl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.SendUlReq{
		AssocId: 1,
		Stream:  2,
		Payload: []byte{'3'},
		Ppid:    4,
	}
	res := &pb.SendUlRes{}

	mock_sctpd_uplink := mock_sctpd.NewMockSctpdUplinkClient(ctrl)
	mock_sctpd_uplink.EXPECT().
		SendUl(ctx, req).
		Return(res, nil)

	router := service.NewRouter(nil, mock_sctpd_uplink)

	logger, logBuffer := testutil.NewTestLogger()
	pus := sctpd.NewProxyUplinkServer(logger, router)

	got, err := pus.SendUl(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tSendUl\t{\"assoc_id\": 1, \"stream\": 2, \"ppid\": 4, \"payload_size\": 1}\n",
		logBuffer.String())
}

func TestProxyUplinkServer_NewAssoc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.NewAssocReq{
		AssocId:     1,
		Instreams:   2,
		Outstreams:  3,
		RanCpIpaddr: net.IPv4(1, 2, 3, 4),
		Ppid:        4,
	}
	res := &pb.NewAssocRes{}

	mock_sctpd_uplink := mock_sctpd.NewMockSctpdUplinkClient(ctrl)
	mock_sctpd_uplink.EXPECT().
		NewAssoc(ctx, req).
		Return(res, nil)

	router := service.NewRouter(nil, mock_sctpd_uplink)

	logger, logBuffer := testutil.NewTestLogger()
	pus := sctpd.NewProxyUplinkServer(logger, router)

	got, err := pus.NewAssoc(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tNewAssoc\t{\"assoc_id\": 1, \"instreams\": 2, \"outstreams\": 3, \"ppid\": 4, \"ran_cp_ipaddr\": \"1.2.3.4\"}\n",
		logBuffer.String())
}

func TestProxyUplinkServer_CloseAssoc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.CloseAssocReq{
		AssocId: 1,
		IsReset: true,
		Ppid:    2,
	}
	res := &pb.CloseAssocRes{}

	mock_sctpd_uplink := mock_sctpd.NewMockSctpdUplinkClient(ctrl)
	mock_sctpd_uplink.EXPECT().
		CloseAssoc(ctx, req).
		Return(res, nil)

	router := service.NewRouter(nil, mock_sctpd_uplink)

	logger, logBuffer := testutil.NewTestLogger()
	pus := sctpd.NewProxyUplinkServer(logger, router)

	got, err := pus.CloseAssoc(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res, got)
	assert.Equal(
		t,
		"DEBUG\tCloseAssoc\t{\"assoc_id\": 1, \"is_reset\": true, \"ppid\": 2}\n",
		logBuffer.String())
}
