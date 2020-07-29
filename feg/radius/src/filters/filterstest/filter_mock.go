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

package filterstest

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"

	"fbc/lib/go/radius"

	"github.com/stretchr/testify/mock"
)

// MockFilter ...
type MockFilter struct {
	mock.Mock
}

// Init ...
func (m *MockFilter) Init(c *config.ServerConfig) error {
	args := m.Called(c)
	return args.Error(0)
}

// Process ...
func (m *MockFilter) Process(c *modules.RequestContext, l string, r *radius.Request) error {
	args := m.Called(c, l, r)
	err := args.Get(0).(error)
	return err
}
