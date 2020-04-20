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
#include <stdio.h>
#include <stdlib.h>
#include <string>
#include <string.h>
#include "IMSIEncoder.h"

namespace openflow {

/**
 * Convert a IMSI string to a uint + length. IMSI strings can contain two
 * prefix zeros for test MCC and maximum fifteen digits. Bit 1 of the
 * compacted uint is always 1, so that we can match on it set. Bits 2-3
 * the compacted uint contain how many leading 0's are in the IMSI. For
 * example, if the IMSI is 001010000000013, the first bit is 0b1, the second
 * two bits would be 0b10 and the remaining bits would be 1010000000013 << 3
 */
uint64_t IMSIEncoder::compact_imsi(const std::string& imsi) {
  const char* imsi_cstr = imsi.c_str();
  uint32_t pad          = 0;
  for (; pad < strlen(imsi_cstr); pad++) {
    if (imsi_cstr[pad] != '0') {
      break;
    }
  }
  uint64_t compacted = strtoull(imsi_cstr + pad, NULL, 10);
  compacted          = compacted << 2 | (pad & 0x3);
  return compacted << 1 | 0x1;  // last bit signifies that IMSI is set
}

/**
 * Convert from the compacted uint back to a string, using the second two bits
 * to determine the padding
 */
std::string IMSIEncoder::expand_imsi(uint64_t compact) {
  // log(MAX_UINT_64) + 2 leading zeros + null
  const int buf_size = 19 + 2 + 1;
  char buf[buf_size];
  uint32_t pad           = (compact >> 1) & 0x3;
  unsigned long long num = compact >> 3;
  for (uint32_t i = 0; i < pad; i++) {
    buf[i] = '0';
  }
  snprintf(buf + pad, buf_size, "%llu", num);
  return std::string(buf);
}

}  // namespace openflow
