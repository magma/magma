/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestGatewayNonEpsConfigs_ValidateModel(t *testing.T) {
	testCases := []struct {
		cfg           *GatewayNonEpsConfigs
		expectedError string
	}{
		// Don't validate if service control is off
		{
			cfg: &GatewayNonEpsConfigs{
				NonEpsServiceControl: swag.Uint32(0),
			},
			expectedError: "",
		},

		// Validate if service control is on
		{
			cfg: &GatewayNonEpsConfigs{
				NonEpsServiceControl: swag.Uint32(1),
			},
			expectedError: "validation failure list:\n" +
				"arfcn_2g in body is required\n" +
				"csfb_mcc in body is required\n" +
				"csfb_mnc in body is required\n" +
				"csfb_rat in body is required\n" +
				"lac in body is required",
		},

		// Happy path
		{
			cfg: &GatewayNonEpsConfigs{
				NonEpsServiceControl: swag.Uint32(2),
				Lac:                  swag.Uint32(1),
				CsfbRat:              swag.Uint32(1),
				CsfbMcc:              "001",
				CsfbMnc:              "01",
				Arfcn2g:              []uint32{1},
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		err := tc.cfg.ValidateModel()
		if err == nil {
			assert.Equal(t, "", tc.expectedError)
		} else {
			assert.Equal(t, err.Error(), tc.expectedError)
		}
	}
}
