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

package servicers

import (
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/services/state/protos"
)

func validateGetStatesRequest(req *protos.GetStatesRequest) error {
	if !funk.IsEmpty(req.Ids) && funk.IsEmpty(req.NetworkID) {
		return errors.New("network ID must be non-empty for non-empty state IDs")
	}
	return nil
}
