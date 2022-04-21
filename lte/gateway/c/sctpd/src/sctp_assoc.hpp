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

#include <stdint.h>

namespace magma {
namespace sctpd {

// Models the state of an SCTP association
class SctpAssoc {
 public:
  int sd;                  ///< Socket descriptor
  uint32_t ppid;           ///< Payload protocol Identifier
  uint16_t instreams;      ///< Number of input streams negotiated
  uint16_t outstreams;     ///< Number of output strams negotiated
  uint32_t assoc_id;       ///< SCTP association id
  uint32_t messages_recv;  ///< Number of messages received
  uint32_t messages_sent;  ///< Number of messages sent

  SctpAssoc();

  // Dump debug information about the SCTP assocation to the log
  void dump() const;
};

}  // namespace sctpd
}  // namespace magma
