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

#include "lte/gateway/c/sctpd/src/sctp_desc.h"

#include "assert.h"

namespace magma {
namespace sctpd {

SctpDesc::SctpDesc(int sd) : _sd(sd) { assert(sd >= 0); }

void SctpDesc::addAssoc(const SctpAssoc& assoc) {
  _assocs[assoc.assoc_id] = assoc;
}

SctpAssoc& SctpDesc::getAssoc(uint32_t assoc_id) {
  return _assocs.at(assoc_id);  // throws std::out_of_range
}

int SctpDesc::delAssoc(uint32_t assoc_id) {
  auto num_removed = _assocs.erase(assoc_id);
  return num_removed == 1 ? 0 : -1;
}

AssocMap::const_iterator SctpDesc::begin() const { return _assocs.cbegin(); }

AssocMap::const_iterator SctpDesc::end() const { return _assocs.cend(); }

int SctpDesc::sd() const { return _sd; }

void SctpDesc::dump() const {
  for (auto const& kv : _assocs) {
    auto assoc = kv.second;
    assoc.dump();
  }
}

}  // namespace sctpd
}  // namespace magma
