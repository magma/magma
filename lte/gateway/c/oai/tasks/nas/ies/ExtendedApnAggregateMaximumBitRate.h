/*
   Copyright 2020 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#ifndef EXTENDED_APN_AGGREGATE_MAXIMUM_BIT_RATE_SEEN
#define EXTENDED_APN_AGGREGATE_MAXIMUM_BIT_RATE_SEEN

#include <stdint.h>

#define EXTENDED_APN_AGGREGATE_MAXIMUM_BIT_RATE_MINIMUM_LENGTH 4
#define EXTENDED_APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_LENGTH 8

typedef struct ExtendedApnAggregateMaximumBitRate_tag {
  uint8_t extendedapnambrfordownlinkunit;
  uint8_t extendedapnambrfordownlink;
  uint8_t extendedapnambrfordownlink_continued;
  uint8_t extendedapnambrforuplinkunit;
  uint8_t extendedapnambrforuplink;
  uint8_t extendedapnambrforuplink_continued;
} ExtendedApnAggregateMaximumBitRate;

int encode_extended_apn_aggregate_maximum_bit_rate(
    ExtendedApnAggregateMaximumBitRate* extendedapnaggregatemaximumbitrate,
    uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_extended_apn_aggregate_maximum_bit_rate(
    ExtendedApnAggregateMaximumBitRate* extendedapnaggregatemaximumbitrate,
    uint8_t iei, uint8_t* buffer, uint32_t len);

void extended_bit_rate_value(
    ExtendedApnAggregateMaximumBitRate* extended_apn_ambr, uint64_t ambr_dl,
    uint64_t ambr_ul);

#endif /* APN_AGGREGATE_MAXIMUM_BIT_RATE_SEEN */
