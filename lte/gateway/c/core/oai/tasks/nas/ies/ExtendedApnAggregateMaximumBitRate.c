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

#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "ExtendedApnAggregateMaximumBitRate.h"

//------------------------------------------------------------------------------
int decode_extended_apn_aggregate_maximum_bit_rate(
    ExtendedApnAggregateMaximumBitRate* extendedapnaggregatemaximumbitrate,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  extendedapnaggregatemaximumbitrate->extendedapnambrfordownlinkunit =
      *(buffer + decoded);
  decoded++;
  extendedapnaggregatemaximumbitrate->extendedapnambrfordownlink =
      *(buffer + decoded);
  decoded++;
  extendedapnaggregatemaximumbitrate->extendedapnambrfordownlink_continued =
      *(buffer + decoded);
  decoded++;
  extendedapnaggregatemaximumbitrate->extendedapnambrforuplinkunit =
      *(buffer + decoded);
  decoded++;
  extendedapnaggregatemaximumbitrate->extendedapnambrforuplink =
      *(buffer + decoded);
  decoded++;
  extendedapnaggregatemaximumbitrate->extendedapnambrforuplink_continued =
      *(buffer + decoded);
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_extended_apn_aggregate_maximum_bit_rate(
    ExtendedApnAggregateMaximumBitRate* extendedapnaggregatemaximumbitrate,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, EXTENDED_APN_AGGREGATE_MAXIMUM_BIT_RATE_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrfordownlinkunit;
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrfordownlink;
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrfordownlink_continued;
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrforuplinkunit;
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrforuplink;
  encoded++;
  *(buffer + encoded) =
      extendedapnaggregatemaximumbitrate->extendedapnambrforuplink_continued;
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void extended_bit_rate_value(
    ExtendedApnAggregateMaximumBitRate* extended_apn_ambr, uint64_t ambr_dl,
    uint64_t ambr_ul) {
  uint64_t enc_ambr_dl = ambr_dl / 4000000;  // ambr_dl is expected in bps
  uint64_t enc_ambr_ul = ambr_ul / 4000000;  // ambr_ul is expected in bps

  uint8_t unit_dl = 3;
  while (enc_ambr_dl & 0xffffffffffff0000) {
    enc_ambr_dl =
        enc_ambr_dl >> 2;  // fit in 16 bits, right shift by 2 for a step of 4
    unit_dl += 1;
  }
  extended_apn_ambr->extendedapnambrfordownlinkunit = unit_dl;

  extended_apn_ambr->extendedapnambrfordownlink = (enc_ambr_dl & 0x00ff);
  extended_apn_ambr->extendedapnambrfordownlink_continued = (enc_ambr_dl >> 8);

  uint8_t unit_ul = 3;
  while (enc_ambr_ul & 0xffffffffffff0000) {
    enc_ambr_ul = enc_ambr_ul >> 2;
    unit_ul += 1;
  }
  extended_apn_ambr->extendedapnambrforuplinkunit = unit_ul;

  extended_apn_ambr->extendedapnambrforuplink = (enc_ambr_ul & 0x00ff);
  extended_apn_ambr->extendedapnambrforuplink_continued = (enc_ambr_ul >> 8);
}