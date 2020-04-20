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

package test_utils

import (
	"testing"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"

	"github.com/stretchr/testify/assert"
)

const (
	defaultServerProtocol = "tcp"
	defaultServerAddr     = "127.0.0.1:0"
	defaultServerHost     = "magma.com"
	defaultServerRealm    = "magma.com"
	DefaultMaxUlBitRate   = uint64(100000000)
	DefaultMaxDlBitRate   = uint64(200000000)
)

var (
	defaultLteAuthAmf = []byte("\x80\x00")
	defaultLteAuthOp  = []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18")
)

// NewTestHomeSubscriberServer creates a HSS with test users so its functionality
// can be tested.
func NewTestHomeSubscriberServer(t *testing.T) *servicers.HomeSubscriberServer {
	store := storage.NewMemorySubscriberStore()

	subs := GetTestSubscribers()
	for _, sub := range subs {
		err := store.AddSubscriber(sub)
		assert.NoError(t, err)
	}

	config := &mconfig.HSSConfig{
		Server: &mconfig.DiamServerConfig{
			Protocol:  defaultServerProtocol,
			Address:   defaultServerAddr,
			DestHost:  defaultServerHost,
			DestRealm: defaultServerRealm,
		},
		LteAuthAmf: defaultLteAuthAmf,
		LteAuthOp:  defaultLteAuthOp,
		DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
			MaxUlBitRate: DefaultMaxUlBitRate,
			MaxDlBitRate: DefaultMaxDlBitRate,
		},
		SubProfiles: make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
	}
	server, err := servicers.NewHomeSubscriberServer(store, config)
	assert.NoError(t, err)
	return server
}
