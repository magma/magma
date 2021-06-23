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
#pragma once

#include <lte/protos/mconfig/mconfigs.pb.h>

namespace magma {
namespace lte {

#define LIAGENTD "liagentd"
#define LIAGENTD_VERSION "1.0"

uint64_t get_time_in_sec_since_epoch();
uint64_t time_difference_from_now(const uint64_t timestamp);

magma::mconfig::LIAgentD get_default_mconfig();
magma::mconfig::LIAgentD load_mconfig();

}  // namespace lte
}  // namespace magma
