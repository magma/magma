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

#include "DiameterCodes.h"

namespace magma {
bool DiameterCodeHandler::is_transient_failure(const uint32_t code) {
  return 4000 <= code && code < 5000;
}

// Diameter code of form 5xxx marks a permanent failure
bool DiameterCodeHandler::is_permanent_failure(const uint32_t code) {
  return 5000 <= code && code < 6000;
}
}  // namespace magma
