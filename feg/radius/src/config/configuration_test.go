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

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRadiusConfig(t *testing.T) {
	conf, err := Read("../radius.config.json")
	require.Nil(t, err)
	require.NotNil(t, conf)
	require.NotNil(t, conf.Server)
	require.Equal(t, conf.Server.LoadBalance, LoadBalanceConfig{})
}

func TestLoadLBConfig(t *testing.T) {
	conf, err := Read("./samples/lb.config.json")
	require.Nil(t, err)
	require.NotNil(t, conf)
	require.NotNil(t, conf.Server)
	require.NotEmpty(t, conf.Server.LoadBalance.ServiceTiers)
	require.NotEmpty(t, conf.Server.LoadBalance.LiveTier)
	require.NotEmpty(t, conf.Server.LoadBalance.Canaries)
}
