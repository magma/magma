/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <stdint.h>
#include <string>

namespace openflow {

class IMSIEncoder {
 public:
  /*
   * Convert a IMSI string to a uint + length. IMSI strings can contain two
   * prefix zeros for test MCC and maximum fifteen digits. The first 2 bits of
   * the compacted uint contain how many leading 0's are in the IMSI. For
   * example, if the IMSI is 001010000000013, the first two bits would be 10 and
   * the remaining bits would be 1010000000013 << 2
   * @param imsi - string representation of imsi
   * @return uint representation of imsi with padding amount at end
   */
  static uint64_t compact_imsi(const std::string& imsi);

  /*
   * Convert from the compacted uint back to a string, using the first two bits
   * to determine the padding
   * @param compact - compacted representation of imsi with padding at end
   * @return pointer to static buffer with imsi inside
   */
  static std::string expand_imsi(uint64_t compact);
};

}  // namespace openflow
