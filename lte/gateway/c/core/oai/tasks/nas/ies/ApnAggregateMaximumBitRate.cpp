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

#include "lte/gateway/c/core/oai/tasks/nas/ies/ApnAggregateMaximumBitRate.hpp"

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/common/common_types.h"

//------------------------------------------------------------------------------
int decode_apn_aggregate_maximum_bit_rate(
    ApnAggregateMaximumBitRate* apnaggregatemaximumbitrate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  apnaggregatemaximumbitrate->apnambrfordownlink = *(buffer + decoded);
  decoded++;
  apnaggregatemaximumbitrate->apnambrforuplink = *(buffer + decoded);
  decoded++;

  if (ielen >= 4) {
    apnaggregatemaximumbitrate->extensions =
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apnaggregatemaximumbitrate->apnambrfordownlink_extended =
        *(buffer + decoded);
    decoded++;
    apnaggregatemaximumbitrate->apnambrforuplink_extended = *(buffer + decoded);
    decoded++;

    if (ielen >= 6) {
      apnaggregatemaximumbitrate->extensions =
          (APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT) |
          (APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION2_PRESENT);
      apnaggregatemaximumbitrate->apnambrfordownlink_extended2 =
          *(buffer + decoded);
      decoded++;
      apnaggregatemaximumbitrate->apnambrforuplink_extended2 =
          *(buffer + decoded);
      decoded++;
    }
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_apn_aggregate_maximum_bit_rate(
    ApnAggregateMaximumBitRate* apnaggregatemaximumbitrate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, APN_AGGREGATE_MAXIMUM_BIT_RATE_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = apnaggregatemaximumbitrate->apnambrfordownlink;
  encoded++;
  *(buffer + encoded) = apnaggregatemaximumbitrate->apnambrforuplink;
  encoded++;

  if (((apnaggregatemaximumbitrate->extensions) &
       (APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT)) |
      (APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION2_PRESENT)) {
    *(buffer + encoded) =
        apnaggregatemaximumbitrate->apnambrfordownlink_extended;
    encoded++;
    *(buffer + encoded) = apnaggregatemaximumbitrate->apnambrforuplink_extended;
    encoded++;
  }
  if (apnaggregatemaximumbitrate->extensions &
      APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION2_PRESENT) {
    *(buffer + encoded) =
        apnaggregatemaximumbitrate->apnambrfordownlink_extended2;
    encoded++;
    *(buffer + encoded) =
        apnaggregatemaximumbitrate->apnambrforuplink_extended2;
    encoded++;
  }

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

// Use 3GPP TS 24.008 figure 10.5.136A, table 10.5.154A
void bit_rate_value_to_eps_qos(ApnAggregateMaximumBitRate* apn_ambr,
                               uint64_t ambr_dl, uint64_t ambr_ul,
                               const apn_ambr_bitrate_unit_t br_unit) {
  uint64_t ambr_dl_kbps = ambr_dl;
  uint64_t ambr_ul_kbps = ambr_ul;
  if (br_unit == BPS) {
    ambr_dl_kbps = ambr_dl / 1000;  // ambr_dl is expected in bps
    ambr_ul_kbps = ambr_ul / 1000;  // ambr_ul is expected in bps
  }
  apn_ambr->apnambrforuplink_extended2 = 0;
  apn_ambr->apnambrforuplink_extended = 0;
  apn_ambr->apnambrfordownlink_extended2 = 0;
  apn_ambr->apnambrfordownlink_extended = 0;
  apn_ambr->apnambrforuplink = 0xff;
  apn_ambr->apnambrfordownlink = 0xff;

  if ((ambr_dl_kbps > 256000) && (ambr_dl_kbps <= 65280000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION2_PRESENT;
    apn_ambr->apnambrfordownlink_extended2 = (ambr_dl_kbps) / 256000;
    ambr_dl_kbps =
        ambr_dl_kbps - (apn_ambr->apnambrfordownlink_extended2) * 256000;
    apn_ambr->apnambrfordownlink = 0xff;
  }
  if ((ambr_dl_kbps > 8640) &&
      ((ambr_dl_kbps > 8600) && (ambr_dl_kbps <= 16000))) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrfordownlink_extended = (ambr_dl_kbps - 8600) / 100;

    ambr_dl_kbps =
        ambr_dl_kbps - (apn_ambr->apnambrfordownlink_extended) * 100 - 8600;
  } else if ((ambr_dl_kbps > 16000) && (ambr_dl_kbps <= 128000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrfordownlink_extended =
        ((ambr_dl_kbps - 16000) / 1000) + 74;
    ambr_dl_kbps = ambr_dl_kbps -
                   (apn_ambr->apnambrfordownlink_extended - 74) * 1000 - 16000;
  } else if ((ambr_dl_kbps > 128000) && (ambr_dl_kbps <= 256000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrfordownlink_extended =
        ((ambr_dl_kbps - 128000) / 2000) + 186;
    ambr_dl_kbps = ambr_dl_kbps -
                   (apn_ambr->apnambrfordownlink_extended - 186) * 2000 -
                   128000;
  }
  // As per 3GPP-24.301, If the network wants to indicate an APN-AMBR for
  // downlink higher than 8640 kbps, it shall set octet 3 to "11111110".
  if ((apn_ambr->apnambrfordownlink != 0xfe) && !(ambr_dl_kbps > 8640)) {
    if (ambr_dl_kbps == 0) {
      apn_ambr->apnambrfordownlink = 0xff;
    }
    if ((ambr_dl_kbps > 0) && (ambr_dl_kbps <= 63)) {
      apn_ambr->apnambrfordownlink = ambr_dl_kbps;
    } else if ((ambr_dl_kbps > 63) && (ambr_dl_kbps <= 575)) {
      apn_ambr->apnambrfordownlink = ((ambr_dl_kbps - 64) / 8) + 64;
      ambr_dl_kbps =
          ambr_dl_kbps - (apn_ambr->apnambrfordownlink_extended - 64) * 8 - 63;
    } else if ((ambr_dl_kbps > 575) && (ambr_dl_kbps <= 8640)) {
      apn_ambr->apnambrfordownlink = ((ambr_dl_kbps - 576) / 64) + 128;
      ambr_dl_kbps = ambr_dl_kbps -
                     (apn_ambr->apnambrfordownlink_extended - 128) * 64 - 575;
    }
  }
  if ((ambr_ul_kbps > 256000) && (ambr_ul_kbps <= 65280000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION2_PRESENT;

    apn_ambr->apnambrforuplink_extended2 = (ambr_ul_kbps) / 256000;
    ambr_ul_kbps =
        ambr_ul_kbps - (apn_ambr->apnambrforuplink_extended2) * 256000;
    apn_ambr->apnambrforuplink = 0xff;
  }
  if ((ambr_ul_kbps > 8640) &&
      ((ambr_ul_kbps > 8600) && (ambr_ul_kbps <= 16000))) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrforuplink_extended = (ambr_ul_kbps - 8600) / 100;
    ambr_ul_kbps =
        ambr_ul_kbps - (apn_ambr->apnambrforuplink_extended) * 100 - 8600;
  } else if ((ambr_ul_kbps > 16000) && (ambr_ul_kbps <= 128000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrforuplink_extended = ((ambr_ul_kbps - 16000) / 1000) + 74;
    ambr_ul_kbps = ambr_ul_kbps -
                   (apn_ambr->apnambrforuplink_extended - 74) * 1000 - 16000;
  } else if ((ambr_ul_kbps > 128000) && (ambr_ul_kbps <= 256000)) {
    apn_ambr->extensions |=
        APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_EXTENSION_PRESENT;
    apn_ambr->apnambrforuplink_extended =
        ((ambr_ul_kbps - 128000) / 2000) + 186;

    ambr_ul_kbps = ambr_ul_kbps -
                   (apn_ambr->apnambrforuplink_extended - 186) * 2000 - 128000;
  }
  // As per 3GPP-24.301, If the network wants to indicate an APN-AMBR for
  // downlink higher than 8640 kbps, it shall set octet 3 to "11111110".
  if ((apn_ambr->apnambrforuplink != 0xfe) && !(ambr_ul_kbps > 8640)) {
    if (ambr_ul_kbps == 0) {
      apn_ambr->apnambrforuplink = 0xff;
    }
    if ((ambr_ul_kbps > 0) && (ambr_ul_kbps <= 63)) {
      apn_ambr->apnambrforuplink = ambr_ul_kbps;
    } else if ((ambr_ul_kbps > 63) && (ambr_ul_kbps <= 575)) {
      apn_ambr->apnambrforuplink = ((ambr_ul_kbps - 64) / 8) + 64;
      ambr_ul_kbps =
          ambr_ul_kbps - (apn_ambr->apnambrforuplink_extended - 64) * 8 - 63;
    } else if ((ambr_ul_kbps > 575) && (ambr_ul_kbps <= 8640)) {
      apn_ambr->apnambrforuplink = ((ambr_ul_kbps - 576) / 64) + 128;
      ambr_ul_kbps =
          ambr_ul_kbps - (apn_ambr->apnambrforuplink_extended - 128) * 64 - 575;
    }
  }
}
