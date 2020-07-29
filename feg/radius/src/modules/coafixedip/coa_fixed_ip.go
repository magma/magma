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

package coafixedip

import (
	"context"
	"net"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"

	"go.uber.org/zap"
)

// Config config has only one parameter which is the ip to forward the request
type Config struct {
	Target string
}

// ModuleCtx ...
type ModuleCtx struct {
	target string
}

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var coaConfig Config
	err := mapstructure.Decode(config, &coaConfig)
	if err != nil {
		return nil, err
	}

	if coaConfig.Target == "" {
		return nil, errors.New("coa module cannot be initialized with empty target value")
	}

	// Validating the correctness of Target
	var host string
	host, _, err = net.SplitHostPort(coaConfig.Target)
	if err != nil {
		return nil, err
	}

	if nil == net.ParseIP(host) {
		return nil, errors.Wrap(err, "Invalid ip address specified")
	}

	return ModuleCtx{target: coaConfig.Target}, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	mod := m.(ModuleCtx)
	c.Logger.Debug("Starting to handle coa radius request")
	requestCode := r.Code
	// Checking that we have received a coa request
	if requestCode != radius.CodeDisconnectRequest && requestCode != radius.CodeCoARequest {
		return next(c, r)
	}

	// Handling the coa request
	res, err := radius.Exchange(context.Background(), r.Packet, mod.target)
	if err != nil {
		return nil, err
	}

	b, err := res.Encode()
	if err != nil {
		c.Logger.Info("failed to serialize CoA response")
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
		Raw:        b,
	}, nil
}
