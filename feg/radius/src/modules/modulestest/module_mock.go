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

package modulestest

import (
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockModule ...
type MockModule struct {
	mock.Mock
}

// Init ...
func (m *MockModule) Init(loggert *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	args := m.Called(config)
	return nil, args.Error(0)
}

// Handle ...
func (m *MockModule) Handle(mCtx modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	args := m.Called(c, r, next)
	res, ok := args.Get(0).(*modules.Response)
	if !ok {
		return nil, args.Error(1)
	}
	return res, args.Error(1)
}
