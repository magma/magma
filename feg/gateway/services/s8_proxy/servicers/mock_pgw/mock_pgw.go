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

package mock_pgw

import (
	"context"
	"fmt"
	"net"
	"time"

	"magma/feg/gateway/gtp"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

const (
	dummyUserPlanePgwIP = "10.0.0.1"
)

// mockPgw is just a wrapper around gtp.Client
type mockPgw struct {
	*gtp.Client
}

func NewStarted(ctx context.Context, sgwAddrStr, pgwAddrsStr string) (*mockPgw, error) {
	mPgw := New()
	err := mPgw.Start(ctx, sgwAddrStr, pgwAddrsStr)
	if err != nil {
		return nil, err
	}
	return mPgw, nil
}

func New() *mockPgw {
	return &mockPgw{}
}

func (mPgw *mockPgw) Start(ctx context.Context, sgwAddrStr, pgwAddrsStr string) error {
	pgwAddrs, err := net.ResolveUDPAddr("udp", pgwAddrsStr)
	if err != nil {
		return fmt.Errorf("Failed to get mock PGW IP: %s", err)
	}

	sgwAddrs, err := net.ResolveUDPAddr("udp", sgwAddrStr)
	if err != nil {
		return fmt.Errorf("Failed to get SGW IP: %s", err)
	}

	// start listening on the specified IP:Port.
	mPgw.Client, err = gtp.NewRunningClient(ctx, pgwAddrs, sgwAddrs, gtpv2.IFTypeS5S8PGWGTPC)
	if err != nil {
		return fmt.Errorf("Failed to get SGW IP: %s", err)
	}
	// Better handle wait for start of service to be ready
	time.Sleep(time.Millisecond * 20)
	time.Sleep(time.Millisecond * 20)

	// register handlers for ALL the message you expect remote endpoint to send.
	mPgw.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest: getHandleCreateSessionRequest(),
		message.MsgTypeDeleteSessionRequest: getHandleDeleteSessionRequest(),
	})
	return nil
}
