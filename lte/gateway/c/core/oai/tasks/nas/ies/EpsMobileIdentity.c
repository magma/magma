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
#include <stdbool.h>

#include "log.h"
#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "EpsMobileIdentity.h"
#include "common_defs.h"

static int decode_guti_eps_mobile_identity(
    guti_eps_mobile_identity_t* guti, uint8_t* buffer);
static int decode_imsi_eps_mobile_identity(
    imsi_eps_mobile_identity_t* imsi, uint8_t* buffer, uint8_t ie_len);
static int decode_imei_eps_mobile_identity(
    imei_eps_mobile_identity_t* imei, uint8_t* buffer);

static int encode_guti_eps_mobile_identity(
    guti_eps_mobile_identity_t* guti, uint8_t* buffer);
static int encode_imsi_eps_mobile_identity(
    imsi_eps_mobile_identity_t* imsi, uint8_t* buffer);
static int encode_imei_eps_mobile_identity(
    imei_eps_mobile_identity_t* imei, uint8_t* buffer);

//------------------------------------------------------------------------------
int decode_eps_mobile_identity(
    eps_mobile_identity_t* epsmobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded_rc = TLV_VALUE_DOESNT_MATCH;
  int decoded    = 0;
  uint8_t ielen  = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  uint8_t typeofidentity = *(buffer + decoded) & 0x7;

  if (typeofidentity == EPS_MOBILE_IDENTITY_IMSI) {
    decoded_rc = decode_imsi_eps_mobile_identity(
        &epsmobileidentity->imsi, buffer, ielen);
  } else if (typeofidentity == EPS_MOBILE_IDENTITY_GUTI) {
    decoded_rc = decode_guti_eps_mobile_identity(
        &epsmobileidentity->guti, buffer + decoded);
  } else if (typeofidentity == EPS_MOBILE_IDENTITY_IMEI) {
    decoded_rc = decode_imei_eps_mobile_identity(
        &epsmobileidentity->imei, buffer + decoded);
  }

  if (decoded_rc < 0) {
    return decoded_rc;
  }
  return (decoded + decoded_rc);
}

//------------------------------------------------------------------------------
int encode_eps_mobile_identity(
    eps_mobile_identity_t* epsmobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  int encoded_rc   = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, EPS_MOBILE_IDENTITY_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if (epsmobileidentity->imsi.typeofidentity == EPS_MOBILE_IDENTITY_IMSI) {
    encoded_rc = encode_imsi_eps_mobile_identity(
        &epsmobileidentity->imsi, buffer + encoded);
  } else if (
      epsmobileidentity->guti.typeofidentity == EPS_MOBILE_IDENTITY_GUTI) {
    encoded_rc = encode_guti_eps_mobile_identity(
        &epsmobileidentity->guti, buffer + encoded);
  } else if (
      epsmobileidentity->imei.typeofidentity == EPS_MOBILE_IDENTITY_IMEI) {
    encoded_rc = encode_imei_eps_mobile_identity(
        &epsmobileidentity->imei, buffer + encoded);
  }

  if (encoded_rc < 0) {
    return encoded_rc;
  }

  *lenPtr = encoded + encoded_rc - 1 - ((iei > 0) ? 1 : 0);
  return (encoded + encoded_rc);
}

//------------------------------------------------------------------------------
static int decode_guti_eps_mobile_identity(
    guti_eps_mobile_identity_t* guti, uint8_t* buffer) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  guti->spare = (*(buffer + decoded) >> 4) & 0xf;

  /*
   * For the GUTI, bits 5 to 8 of octet 3 are coded as "1111"
   */
  if (guti->spare != 0xf) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  guti->oddeven        = (*(buffer + decoded) >> 3) & 0x1;
  guti->typeofidentity = *(buffer + decoded) & 0x7;

  if (guti->typeofidentity != EPS_MOBILE_IDENTITY_GUTI) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  guti->mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mcc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  guti->mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mcc_digit3 = *(buffer + decoded) & 0xf;
  decoded++;
  guti->mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mnc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  // IES_DECODE_U16(guti->mmegroupid, *(buffer + decoded));
  IES_DECODE_U16(buffer, decoded, guti->mme_group_id);
  guti->mme_code = *(buffer + decoded);
  decoded++;
  // IES_DECODE_U32(guti->mtmsi, *(buffer + decoded));
  IES_DECODE_U32(buffer, decoded, guti->m_tmsi);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_imsi_eps_mobile_identity(
    imsi_eps_mobile_identity_t* imsi, uint8_t* buffer, uint8_t len) {
  uint8_t decoded = 0;
  uint8_t ielen   = 0;

  ielen = *(buffer + decoded); /* Pointing buffer to IE length field, to include
                                  the ieLen byte*/
  decoded++;
  imsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (imsi->typeofidentity != EPS_MOBILE_IDENTITY_IMSI) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  imsi->oddeven         = (*(buffer + decoded) >> 3) & 0x1;
  imsi->identity_digit1 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->num_digits++;
  decoded++;
  if (decoded <= ielen) {
    imsi->identity_digit2 = *(buffer + decoded) & 0xf;
    imsi->identity_digit3 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit4 = *(buffer + decoded) & 0xf;
    imsi->identity_digit5 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit6 = *(buffer + decoded) & 0xf;
    imsi->identity_digit7 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit8 = *(buffer + decoded) & 0xf;
    imsi->identity_digit9 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit10 = *(buffer + decoded) & 0xf;
    imsi->identity_digit11 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit12 = *(buffer + decoded) & 0xf;
    imsi->identity_digit13 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->num_digits += 2;
  }
  if (decoded <= ielen) {
    imsi->identity_digit14 = *(buffer + decoded) & 0xf;
    imsi->identity_digit15 = (*(buffer + decoded) >> 4) & 0xf;
    imsi->num_digits += 2;
    decoded++;
  }

  if (imsi->oddeven == false) {
    imsi->num_digits--; /* For even number of digits*/
  }
  /*
   * IMSI is coded using BCD coding. If the number of identity digits is
   * even then bits 5 to 8 of the last octet shall be filled with an end
   * mark coded as "1111".
   */
  if ((imsi->oddeven == EPS_MOBILE_IDENTITY_EVEN) &&
      (imsi->identity_digit15 != 0x0f)) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded--; /*ielen is already included*/
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_imei_eps_mobile_identity(
    imei_eps_mobile_identity_t* imei, uint8_t* buffer) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  imei->typeofidentity = *(buffer + decoded) & 0x7;

  if (imei->typeofidentity != EPS_MOBILE_IDENTITY_IMEI) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  imei->oddeven         = (*(buffer + decoded) >> 3) & 0x1;
  imei->identity_digit1 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit2 = *(buffer + decoded) & 0xf;
  imei->identity_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit4 = *(buffer + decoded) & 0xf;
  imei->identity_digit5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit6 = *(buffer + decoded) & 0xf;
  imei->identity_digit7 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit8 = *(buffer + decoded) & 0xf;
  imei->identity_digit9 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit10 = *(buffer + decoded) & 0xf;
  imei->identity_digit11 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit12 = *(buffer + decoded) & 0xf;
  imei->identity_digit13 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->identity_digit14 = *(buffer + decoded) & 0xf;
  imei->identity_digit15 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int encode_guti_eps_mobile_identity(
    guti_eps_mobile_identity_t* guti, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) =
      0xf0 | ((guti->oddeven & 0x1) << 3) | (guti->typeofidentity & 0x7);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mcc_digit2 & 0xf) << 4) | (guti->mcc_digit1 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit3 & 0xf) << 4) | (guti->mcc_digit3 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit2 & 0xf) << 4) | (guti->mnc_digit1 & 0xf);
  encoded++;
  IES_ENCODE_U16(buffer, encoded, guti->mme_group_id);
  *(buffer + encoded) = guti->mme_code;
  encoded++;
  IES_ENCODE_U32(buffer, encoded, guti->m_tmsi);
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_imsi_eps_mobile_identity(
    imsi_eps_mobile_identity_t* imsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | (imsi->identity_digit1 << 4) |
                        (imsi->oddeven << 3) | (imsi->typeofidentity);
  encoded++;
  *(buffer + encoded) =
      0x00 | (imsi->identity_digit3 << 4) | imsi->identity_digit2;
  encoded++;
  // Quick fix, should do a loop, but try without modifying struct!
  if (imsi->num_digits > 3) {
    if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
      *(buffer + encoded) =
          0x00 | (imsi->identity_digit5 << 4) | imsi->identity_digit4;
    } else {
      *(buffer + encoded) = 0xf0 | imsi->identity_digit4;
    }
    encoded++;
    if (imsi->num_digits > 5) {
      if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
        *(buffer + encoded) =
            0x00 | (imsi->identity_digit7 << 4) | imsi->identity_digit6;
      } else {
        *(buffer + encoded) = 0xf0 | imsi->identity_digit6;
      }
      encoded++;
      if (imsi->num_digits > 7) {
        if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
          *(buffer + encoded) =
              0x00 | (imsi->identity_digit9 << 4) | imsi->identity_digit8;
        } else {
          *(buffer + encoded) = 0xf0 | imsi->identity_digit8;
        }
        encoded++;
        if (imsi->num_digits > 9) {
          if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
            *(buffer + encoded) =
                0x00 | (imsi->identity_digit11 << 4) | imsi->identity_digit10;
          } else {
            *(buffer + encoded) = 0xf0 | imsi->identity_digit10;
          }
          encoded++;
          if (imsi->num_digits > 11) {
            if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
              *(buffer + encoded) =
                  0x00 | (imsi->identity_digit13 << 4) | imsi->identity_digit12;
            } else {
              *(buffer + encoded) = 0xf0 | imsi->identity_digit12;
            }
            encoded++;
            if (imsi->num_digits > 13) {
              if (imsi->oddeven != EPS_MOBILE_IDENTITY_EVEN) {
                *(buffer + encoded) = 0x00 | (imsi->identity_digit15 << 4) |
                                      imsi->identity_digit14;
              } else {
                *(buffer + encoded) = 0xf0 | imsi->identity_digit14;
              }
              encoded++;
            }
          }
        }
      }
    }
  }

  return encoded;
}

//------------------------------------------------------------------------------
static int encode_imei_eps_mobile_identity(
    imei_eps_mobile_identity_t* imei, uint8_t* buffer) {
  uint32_t encoded = 0;

  // IMEI fixed length of 15 digits
  *(buffer + encoded) = 0x00 | (imei->identity_digit1 << 4) |
                        (imei->oddeven << 3) | (imei->typeofidentity);
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit3 << 4) | imei->identity_digit2;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit5 << 4) | imei->identity_digit4;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit7 << 4) | imei->identity_digit6;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit9 << 4) | imei->identity_digit8;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit11 << 4) | imei->identity_digit10;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit13 << 4) | imei->identity_digit12;
  encoded++;
  *(buffer + encoded) =
      0x00 | (imei->identity_digit15 << 4) | imei->identity_digit14;
  encoded++;
  return encoded;
}
