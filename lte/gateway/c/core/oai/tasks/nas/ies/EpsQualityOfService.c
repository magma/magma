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

#include "EpsQualityOfService.h"

#include <stdint.h>
#include <stdbool.h>
#include <string.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "common_types.h"
#include "common_defs.h"

//------------------------------------------------------------------------------
static int decode_eps_qos_bit_rates(
    EpsQoSBitRates* epsqosbitrates, const uint8_t* buffer) {
  int decoded = 0;

  epsqosbitrates->maxBitRateForUL = *(buffer + decoded);
  decoded++;
  epsqosbitrates->maxBitRateForDL = *(buffer + decoded);
  decoded++;
  epsqosbitrates->guarBitRateForUL = *(buffer + decoded);
  decoded++;
  epsqosbitrates->guarBitRateForDL = *(buffer + decoded);
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int decode_eps_quality_of_service(
    EpsQualityOfService* epsqualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  epsqualityofservice->qci = *(buffer + decoded);
  decoded++;

  if (ielen > 2 + (iei > 0) ? 1 : 0) {
    /*
     * bitRates is present
     */
    epsqualityofservice->bitRatesPresent = 1;
    decoded += decode_eps_qos_bit_rates(
        &epsqualityofservice->bitRates, buffer + decoded);
  } else {
    /*
     * bitRates is not present
     */
    epsqualityofservice->bitRatesPresent = 0;
  }

  if (ielen > 6 + (iei > 0) ? 1 : 0) {
    /*
     * bitRatesExt is present
     */
    epsqualityofservice->bitRatesExtPresent = 1;
    decoded += decode_eps_qos_bit_rates(
        &epsqualityofservice->bitRatesExt, buffer + decoded);
  } else {
    /*
     * bitRatesExt is not present
     */
    epsqualityofservice->bitRatesExtPresent = 0;
  }

  return decoded;
}

//------------------------------------------------------------------------------
static int encode_eps_qos_bit_rates(
    const EpsQoSBitRates* epsqosbitrates, uint8_t* buffer) {
  int encoded = 0;

  *(buffer + encoded) = epsqosbitrates->maxBitRateForUL;
  encoded++;
  *(buffer + encoded) = epsqosbitrates->maxBitRateForDL;
  encoded++;
  *(buffer + encoded) = epsqosbitrates->guarBitRateForUL;
  encoded++;
  *(buffer + encoded) = epsqosbitrates->guarBitRateForDL;
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
int encode_eps_quality_of_service(
    EpsQualityOfService* epsqualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, EPS_QUALITY_OF_SERVICE_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = epsqualityofservice->qci;
  encoded++;

  if (epsqualityofservice->bitRatesPresent) {
    encoded += encode_eps_qos_bit_rates(
        &epsqualityofservice->bitRates, buffer + encoded);
  }

  if (epsqualityofservice->bitRatesExtPresent) {
    encoded += encode_eps_qos_bit_rates(
        &epsqualityofservice->bitRatesExt, buffer + encoded);
  }

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

#define EPS_QOS_BIT_RATE_MAX 262144  // 256 Mbps
//------------------------------------------------------------------------------
int eps_qos_bit_rate_value(uint8_t br) {
  if (br < 0b00000001) {
    return (EPS_QOS_BIT_RATE_MAX);
  } else if ((br > 0b00000000) && (br < 0b01000000)) {
    return (br);
  } else if ((br > 0b00111111) && (br < 0b10000000)) {
    return (64 + (br - 0b01000000) * 8);
  } else if ((br > 0b01111111) && (br < 0b11111111)) {
    return (576 + (br - 0b10000000) * 64);
  } else {
    return (0);
  }
}

//------------------------------------------------------------------------------
int eps_qos_bit_rate_ext_value(uint8_t br) {
  if ((br > 0b00000000) && (br < 0b01001011)) {
    return (8600 + br * 100);
  } else if ((br > 0b01001010) && (br < 0b10111011)) {
    return (16384 + (br - 0b01001010) * 1024);
  } else if ((br > 0b10111010) && (br < 0b11111011)) {
    return (131072 + (br - 0b10111010) * 2048);
  } else {
    return (-1);
  }
}

//------------------------------------------------------------------------------
int qos_params_to_eps_qos(
    const qci_t qci, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    const bitrate_t gbr_dl, const bitrate_t gbr_ul,
    EpsQualityOfService* const eps_qos, bool is_default_bearer) {
  uint64_t mbr_dl_kbps = mbr_dl / 1000;  // mbr_dl is expected in bps
  uint64_t mbr_ul_kbps = mbr_ul / 1000;  // mbr_ul is expected in bps
  uint64_t gbr_dl_kbps = gbr_dl / 1000;  // mbr_dl is expected in bps
  uint64_t gbr_ul_kbps = gbr_ul / 1000;  // mbr_ul is expected in bps

  if (eps_qos) {
    memset(eps_qos, 0, sizeof(EpsQualityOfService));
    eps_qos->qci = qci;
    if (!is_default_bearer) {
      eps_qos->bitRatesPresent = 1;
      if (mbr_ul_kbps == 0) {
        eps_qos->bitRates.maxBitRateForUL = 0xff;
      } else if ((mbr_ul_kbps > 0) && (mbr_ul_kbps <= 63)) {
        eps_qos->bitRates.maxBitRateForUL = mbr_ul_kbps;
      } else if ((mbr_ul_kbps > 63) && (mbr_ul_kbps <= 575)) {
        eps_qos->bitRates.maxBitRateForUL = ((mbr_ul_kbps - 64) / 8) + 64;
      } else if ((mbr_ul_kbps > 575) && (mbr_ul_kbps <= 8640)) {
        eps_qos->bitRates.maxBitRateForUL = ((mbr_ul_kbps - 576) / 64) + 128;
      } else if (mbr_ul_kbps > 8640) {
        eps_qos->bitRates.maxBitRateForUL = 0xfe;
        eps_qos->bitRatesExtPresent       = 1;
        if ((mbr_ul_kbps >= 8600) && (mbr_ul_kbps <= 16000)) {
          eps_qos->bitRatesExt.maxBitRateForUL = (mbr_ul_kbps - 8600) / 100;
        } else if ((mbr_ul_kbps > 16000) && (mbr_ul_kbps <= 128000)) {
          eps_qos->bitRatesExt.maxBitRateForUL =
              ((mbr_ul_kbps - 16000) / 1000) + 74;
        } else if ((mbr_ul_kbps > 128000) && (mbr_ul_kbps <= 256000)) {
          eps_qos->bitRatesExt.maxBitRateForUL =
              ((mbr_ul_kbps - 128000) / 2000) + 186;
        }
      }
      if (mbr_dl_kbps == 0) {
        eps_qos->bitRates.maxBitRateForDL = 0xff;
      } else if ((mbr_dl_kbps > 0) && (mbr_dl_kbps <= 63)) {
        eps_qos->bitRates.maxBitRateForDL = mbr_dl_kbps;
      } else if ((mbr_dl_kbps > 63) && (mbr_dl_kbps <= 575)) {
        eps_qos->bitRates.maxBitRateForDL = ((mbr_dl_kbps - 64) / 8) + 64;
      } else if ((mbr_dl_kbps > 575) && (mbr_dl_kbps <= 8640)) {
        eps_qos->bitRates.maxBitRateForDL = ((mbr_dl_kbps - 576) / 64) + 128;
      } else if (mbr_dl_kbps > 8640) {
        eps_qos->bitRates.maxBitRateForDL = 0xfe;
        eps_qos->bitRatesExtPresent       = 1;
        if ((mbr_dl_kbps >= 8600) && (mbr_dl_kbps <= 16000)) {
          eps_qos->bitRatesExt.maxBitRateForDL = (mbr_dl_kbps - 8600) / 100;
        } else if ((mbr_dl_kbps > 16000) && (mbr_dl_kbps <= 128000)) {
          eps_qos->bitRatesExt.maxBitRateForDL =
              ((mbr_dl_kbps - 16000) / 1000) + 74;
        } else if ((mbr_dl_kbps > 128000) && (mbr_dl_kbps <= 256000)) {
          eps_qos->bitRatesExt.maxBitRateForDL =
              ((mbr_dl_kbps - 128000) / 2000) + 186;
        }
      }
      if (gbr_ul_kbps == 0) {
        eps_qos->bitRates.guarBitRateForUL = 0xff;
      } else if ((gbr_ul_kbps > 0) && (gbr_ul_kbps <= 63)) {
        eps_qos->bitRates.guarBitRateForUL = gbr_ul_kbps;
      } else if ((gbr_ul_kbps > 63) && (gbr_ul_kbps <= 575)) {
        eps_qos->bitRates.guarBitRateForUL = ((gbr_ul_kbps - 64) / 8) + 64;
      } else if ((gbr_ul_kbps > 575) && (gbr_ul_kbps <= 8640)) {
        eps_qos->bitRates.guarBitRateForUL = ((gbr_ul_kbps - 576) / 64) + 128;
      } else if (gbr_ul_kbps > 8640) {
        eps_qos->bitRates.guarBitRateForUL = 0xfe;
        eps_qos->bitRatesExtPresent        = 1;
        if ((gbr_ul_kbps >= 8600) && (gbr_ul_kbps <= 16000)) {
          eps_qos->bitRatesExt.guarBitRateForUL = (gbr_ul_kbps - 8600) / 100;
        } else if ((gbr_ul_kbps > 16000) && (gbr_ul_kbps <= 128000)) {
          eps_qos->bitRatesExt.guarBitRateForUL =
              ((gbr_ul_kbps - 16000) / 1000) + 74;
        } else if ((gbr_ul_kbps > 128000) && (gbr_ul_kbps <= 256000)) {
          eps_qos->bitRatesExt.guarBitRateForUL =
              ((gbr_ul_kbps - 128000) / 2000) + 186;
        }
      }
      if (gbr_dl_kbps == 0) {
        eps_qos->bitRates.guarBitRateForDL = 0xff;
      } else if ((gbr_dl_kbps > 0) && (gbr_dl_kbps <= 63)) {
        eps_qos->bitRates.guarBitRateForDL = gbr_dl_kbps;
      } else if ((gbr_dl_kbps > 63) && (gbr_dl_kbps <= 575)) {
        eps_qos->bitRates.guarBitRateForDL = ((gbr_dl_kbps - 64) / 8) + 64;
      } else if ((gbr_dl_kbps >= 575) && (gbr_dl_kbps <= 8640)) {
        eps_qos->bitRates.guarBitRateForDL = ((gbr_dl_kbps - 576) / 64) + 128;
      } else if (gbr_dl_kbps > 8640) {
        eps_qos->bitRates.guarBitRateForDL = 0xfe;
        eps_qos->bitRatesExtPresent        = 1;
        if ((gbr_dl_kbps >= 8600) && (gbr_dl_kbps <= 16000)) {
          eps_qos->bitRatesExt.guarBitRateForDL = (gbr_dl_kbps - 8600) / 100;
        } else if ((gbr_dl_kbps > 16000) && (gbr_dl_kbps <= 128000)) {
          eps_qos->bitRatesExt.guarBitRateForDL =
              ((gbr_dl_kbps - 16000) / 1000) + 74;
        } else if ((gbr_dl_kbps > 128000) && (gbr_dl_kbps <= 256000)) {
          eps_qos->bitRatesExt.guarBitRateForDL =
              ((gbr_dl_kbps - 128000) / 2000) + 186;
        }
      }
    }
    return RETURNok;
  }
  return RETURNerror;
}
