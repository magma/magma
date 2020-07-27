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

package adaptruckus

import (
	"encoding/binary"

	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/ruckus"

	"go.uber.org/zap"
)

// Init module interface implementation
func Init(loggert *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	return nil, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	found := false
	values, err := ruckus.RuckusTCAcctCtrs_Gets(r.Packet)
	if err == nil && values != nil {
		for _, ctrs := range values {
			for _, ctr := range ctrs {
				if ctr.RuckusAcctCtrsTCName == "internet" {
					c.Logger.Debug("found ruckus internet quota, converting to plain radius")

					r.Set(
						rfc2866.AcctInputOctets_Type,
						radius.Attribute(toNetowkOrder(ctr.RuckusAcctCtrsInputOctets)),
					)

					r.Set(
						rfc2866.AcctInputPackets_Type,
						radius.Attribute(toNetowkOrder(ctr.RuckusAcctCtrsInputPackets)),
					)

					r.Set(
						rfc2866.AcctOutputOctets_Type,
						radius.Attribute(toNetowkOrder(ctr.RuckusAcctCtrsOutputOctets)),
					)

					r.Set(
						rfc2866.AcctOutputPackets_Type,
						radius.Attribute(toNetowkOrder(ctr.RuckusAcctCtrsOutputPackets)),
					)
					found = true
					break
				}
			}
		}
	}

	if !found {
		c.Logger.Debug("could not find Ruckus Acct attributes - skipping")
	}

	return next(c, r)
}

func toNetowkOrder(value uint64) []byte {
	var result = make([]byte, 8)
	binary.BigEndian.PutUint64(result, value)
	return result
}
