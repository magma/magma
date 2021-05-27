/**
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
#include <vector>
#include "sgw_defs.h"
#include "spgw_types.h"

namespace magma {

gtpv1u_data_t make_gtpv1u_data(int fd0, int fd1u);

spgw_state_t make_spgw_state(uint32_t gtpv1u_teid, int fd0, int fd1u);

}  // namespace magma
