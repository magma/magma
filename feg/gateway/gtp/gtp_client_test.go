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

package gtp

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-gtp/gtpv2"
)

const (
	sgwAddrs = "127.0.0.1:0"
	pgwAddrs = "127.0.0.1:0" //port 0 means go will choose. Selected port will be injected on getDefaultConfig
	IMSI1    = "123456789012345"
)

func TestEcho(t *testing.T) {
	// run GTP server (PGW)
	pgwCli := startGTPServer(t)
	actualServerIPAndPort := pgwCli.LocalAddr().String()

	// run GTP client (SGW) and send echo message.
	_, err := NewConnectedAutoClient(context.Background(), actualServerIPAndPort, gtpv2.IFTypeS5S8SGWGTPC)
	assert.NoError(t, err)
	// if no error service was started and echo was received properly
}

func startGTPServer(t *testing.T) *Client {
	pgwUPDAddr, err := net.ResolveUDPAddr("udp", pgwAddrs)
	assert.NoError(t, err)
	sgwUPDAddr, err := net.ResolveUDPAddr("udp", sgwAddrs)
	assert.NoError(t, err)

	pgwConn, err := NewRunningClient(context.Background(), pgwUPDAddr, sgwUPDAddr, gtpv2.IFTypeS5S8PGWGTPC)

	// Better handle wait for start of service
	time.Sleep(time.Millisecond * 20)
	time.Sleep(time.Millisecond * 20)
	return pgwConn
}
