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

#include "sctp_assoc.h"

#include <iostream>

#include "util.h"

namespace magma {
namespace sctpd {

SctpAssoc::SctpAssoc()
    : sd(0),
      ppid(0),
      instreams(0),
      outstreams(0),
      assoc_id(0),
      messages_recv(0),
      messages_sent(0) {}

void SctpAssoc::dump() const {
  MLOG(MDEBUG) << "SctpAssoc<id: " << std::to_string(this->assoc_id) << ">"
               << std::endl;
}

}  // namespace sctpd
}  // namespace magma
