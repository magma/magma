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

package testloopback

import (
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"go.uber.org/zap"
)

/*
 * This module is used for testing - it reflects all attributes from the request into the response.
 * This allows a test to check the changes made by the modules within the chain
 * its also possible to add some new attributes if a test requires them & they dont interfere with generic testing
 */

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(_ modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	logger := c.Logger.With(zap.String("module_name", "testloopback"))

	// Create a response with all attributes copied from request
	resp := modules.Response{
		Code:       getResponseCode(r),
		Attributes: r.Packet.Attributes,
	}

	logger.Debug("generating dummy response", zap.Any("dummy response", resp))
	return &resp, nil
}

func getResponseCode(r *radius.Request) radius.Code {
	switch r.Code {
	case radius.CodeAccessRequest:
		return radius.CodeAccessAccept
	case radius.CodeAccountingRequest:
		return radius.CodeAccountingResponse
	case radius.CodeCoARequest:
		return radius.CodeCoAACK
	case radius.CodeDisconnectRequest:
		return radius.CodeCoAACK

	case radius.CodeAccessReject:
	case radius.CodeAccountingResponse:
	case radius.CodeAccessChallenge:
	case radius.CodeStatusServer:
	case radius.CodeStatusClient:
	case radius.CodeDisconnectACK:
	case radius.CodeDisconnectNAK:
	case radius.CodeCoAACK:
	case radius.CodeCoANAK:
	case radius.CodeAccessAccept:
	case radius.CodeReserved:
	}
	return r.Code
}
