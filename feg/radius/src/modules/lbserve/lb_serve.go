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

package lbserve

import (
	"context"
	"errors"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"go.uber.org/zap"
)

const errMissingRequiredUpstreamHostText = "session state required field upstream host is missing, unable to serve request"

var errMissingRequiredUpstreamHost = errors.New(errMissingRequiredUpstreamHostText)

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	state, err := c.SessionStorage.Get()
	if err != nil {
		c.Logger.Error(
			"Error loading session state, unable to serve request",
			zap.Error(err),
		)
		return nil, err
	}

	if state.UpstreamHost == "" {
		c.Logger.Error(errMissingRequiredUpstreamHostText)
		return nil, errMissingRequiredUpstreamHost
	}

	res, err := radius.Exchange(context.Background(), r.Packet, state.UpstreamHost)
	if err != nil {
		c.Logger.Error("LB Serve received failed response", zap.Error(err))
		return nil, err
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
	}, nil
}
