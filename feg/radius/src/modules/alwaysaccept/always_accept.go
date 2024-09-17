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

package alwaysaccept

import (
	"errors"

	"fbc/cwf/radius/modules"

	"go.uber.org/zap"
	"layeh.com/radius"
)

// Init module interface implementation
func Init(loggert *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	if r.Code != radius.CodeAccessRequest {
		return nil, errors.New("module cannot handle anything other than Access-Request messages")
	}

	return &modules.Response{
		Code:       radius.CodeAccessAccept,
		Attributes: radius.Attributes{},
	}, nil
}
