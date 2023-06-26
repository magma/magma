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
#include "lte/gateway/c/core/oai/test/spgw_task/state_creators.hpp"

#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/conversions.h"
}

#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/state_converter.hpp"

namespace magma {

gtpv1u_data_t make_gtpv1u_data(int fd0, int fd1u) {
  gtpv1u_data_t data;
  data.fd0 = fd0;
  data.fd1u = fd1u;
  return data;
}

spgw_state_t make_spgw_state(uint32_t gtpv1u_teid, int fd0, int fd1u) {
  spgw_state_t result;
  result.gtpv1u_teid = gtpv1u_teid;
  result.gtpv1u_data = make_gtpv1u_data(fd0, fd1u);
  return result;
}

}  // namespace magma
