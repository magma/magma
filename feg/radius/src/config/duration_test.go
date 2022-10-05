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
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TestWrapper struct {
	D Duration `json:"d"`
}

func TestMarshal(t *testing.T) {
	// Act
	d, err := json.Marshal(&TestWrapper{
		D: Duration{time.Second * 5},
	})

	// Assert
	require.Nil(t, err)
	require.True(t, strings.Contains(string(d), "\"d\":\"5s\""), string(d))
}

func TestUnmarshal(t *testing.T) {
	var d TestWrapper
	err := json.Unmarshal([]byte(`{"d": "1h"}`), &d)

	// Assert
	require.Nil(t, err)
	require.Equal(t, uint64(3600000000000), uint64(d.D.Duration))
}
