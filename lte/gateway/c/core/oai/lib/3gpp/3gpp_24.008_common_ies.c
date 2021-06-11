/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file 3gpp_24.008_common_ies.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdbool.h>
#include <stdint.h>

#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "common_defs.h"
#include "TLVDecoder.h"
#include "TLVEncoder.h"
#include "log.h"

//******************************************************************************
// 10.5.1 Common information elements
//******************************************************************************
//------------------------------------------------------------------------------
// 10.5.1.2 Ciphering Key Sequence Number
//------------------------------------------------------------------------------
int decode_ciphering_key_sequence_number_ie(
    ciphering_key_sequence_number_t* cipheringkeysequencenumber,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, CIPHERING_KEY_SEQUENCE_NUMBER_IE_MAX_LENGTH, len);

  if (iei_present) {
    CHECK_IEI_DECODER((*buffer & 0xf0), C_CIPHERING_KEY_SEQUENCE_NUMBER_IEI);
  }

  *cipheringkeysequencenumber = *buffer & 0x7;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_ciphering_key_sequence_number_ie(
    ciphering_key_sequence_number_t* cipheringkeysequencenumber,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint8_t encoded = 0;

  /*
   * Checking length and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, CIPHERING_KEY_SEQUENCE_NUMBER_IE_MAX_LENGTH, len);
  *(buffer + encoded) = 0x00;
  if (iei_present) {
    *(buffer + encoded) = (C_CIPHERING_KEY_SEQUENCE_NUMBER_IEI & 0xf0);
  }
  *(buffer + encoded) |= (*cipheringkeysequencenumber & 0x7);
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.1.3 Location Area Identification
//------------------------------------------------------------------------------

int decode_location_area_identification_ie(
    location_area_identification_t* locationareaidentification,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(C_LOCATION_AREA_IDENTIFICATION_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH - 1), len);
  }

  locationareaidentification->mccdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mccdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  locationareaidentification->mncdigit3 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mccdigit3 = *(buffer + decoded) & 0xf;
  decoded++;
  locationareaidentification->mncdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mncdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  // IES_DECODE_U16(locationareaidentification->lac, *(buffer + decoded));
  IES_DECODE_U16(buffer, decoded, locationareaidentification->lac);
  return decoded;
}

//------------------------------------------------------------------------------
int encode_location_area_identification_ie(
    location_area_identification_t* locationareaidentification,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH, len);
    *buffer = C_LOCATION_AREA_IDENTIFICATION_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH - 1), len);
  }

  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mccdigit2 & 0xf) << 4) |
                        (locationareaidentification->mccdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mncdigit3 & 0xf) << 4) |
                        (locationareaidentification->mccdigit3 & 0xf);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mncdigit2 & 0xf) << 4) |
                        (locationareaidentification->mncdigit1 & 0xf);
  encoded++;
  IES_ENCODE_U16(buffer, encoded, locationareaidentification->lac);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.1.4 Mobile Identity
//------------------------------------------------------------------------------
static int decode_imsi_mobile_identity(
    imsi_mobile_identity_t* imsi, uint8_t* buffer, const uint32_t len);
static int decode_imei_mobile_identity(
    imei_mobile_identity_t* imei, uint8_t* buffer, const uint32_t len);
static int decode_imeisv_mobile_identity(
    imeisv_mobile_identity_t* imeisv, uint8_t* buffer, const uint32_t len);
static int decode_tmsi_mobile_identity(
    tmsi_mobile_identity_t* tmsi, uint8_t* buffer, const uint32_t len);
static int decode_tmgi_mobile_identity(
    tmgi_mobile_identity_t* tmgi, uint8_t* buffer, const uint32_t len);
static int decode_no_mobile_identity(
    no_mobile_identity_t* no_id, uint8_t* buffer, const uint32_t len);
static int encode_imsi_mobile_identity(
    imsi_mobile_identity_t* imsi, uint8_t* buffer, const uint32_t len);
static int encode_imei_mobile_identity(
    imei_mobile_identity_t* imei, uint8_t* buffer, const uint32_t len);
static int encode_imeisv_mobile_identity(
    imeisv_mobile_identity_t* imeisv, uint8_t* buffer, const uint32_t len);
static int encode_tmsi_mobile_identity(
    tmsi_mobile_identity_t* tmsi, uint8_t* buffer, const uint32_t len);
static int encode_tmgi_mobile_identity(
    tmgi_mobile_identity_t* tmgi, uint8_t* buffer, const uint32_t len);
static int encode_no_mobile_identity(
    no_mobile_identity_t* no_id, uint8_t* buffer, const uint32_t len);

//------------------------------------------------------------------------------
int decode_mobile_identity_ie(
    mobile_identity_t* mobileidentity, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  int decoded_rc = TLV_VALUE_DOESNT_MATCH;
  int decoded    = 0;
  uint8_t ielen  = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, MOBILE_IDENTITY_IE_MIN_LENGTH, len);
    CHECK_IEI_DECODER(C_MOBILE_IDENTITY_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (MOBILE_IDENTITY_IE_MIN_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  uint8_t typeofidentity = *(buffer + decoded) & 0x7;

  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if (typeofidentity == MOBILE_IDENTITY_IMSI) {
    decoded_rc =
        decode_imsi_mobile_identity(&mobileidentity->imsi, buffer, ielen);
  } else if (typeofidentity == MOBILE_IDENTITY_IMEI) {
    decoded_rc = decode_imei_mobile_identity(
        &mobileidentity->imei, buffer + decoded, ielen);
  } else if (typeofidentity == MOBILE_IDENTITY_IMEISV) {
    decoded_rc = decode_imeisv_mobile_identity(
        &mobileidentity->imeisv, buffer + decoded, ielen);
  } else if (typeofidentity == MOBILE_IDENTITY_TMSI) {
    decoded_rc = decode_tmsi_mobile_identity(
        &mobileidentity->tmsi, buffer + decoded, ielen);
  } else if (typeofidentity == MOBILE_IDENTITY_TMGI) {
    decoded_rc = decode_tmgi_mobile_identity(
        &mobileidentity->tmgi, buffer + decoded, ielen);
  } else if (typeofidentity == MOBILE_IDENTITY_NOT_AVAILABLE) {
    decoded_rc = decode_no_mobile_identity(
        &mobileidentity->no_id, buffer + decoded, ielen);
  } else {
    return TLV_VALUE_DOESNT_MATCH;
  }
  if (decoded_rc < 0) {
    return decoded_rc;
  }
  return (decoded + decoded_rc);
}

//------------------------------------------------------------------------------
int encode_mobile_identity_ie(
    mobile_identity_t* mobileidentity, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  uint8_t* lenPtr;
  int encoded_rc   = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, MOBILE_IDENTITY_IE_MIN_LENGTH, len);
    *buffer = C_MOBILE_IDENTITY_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (MOBILE_IDENTITY_IE_MIN_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if (mobileidentity->no_id.typeofidentity != MOBILE_IDENTITY_NOT_AVAILABLE) {
    if (mobileidentity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      encoded_rc = encode_imsi_mobile_identity(
          &mobileidentity->imsi, buffer + encoded, len - encoded);
    } else if (mobileidentity->imei.typeofidentity == MOBILE_IDENTITY_IMEI) {
      encoded_rc = encode_imei_mobile_identity(
          &mobileidentity->imei, buffer + encoded, len - encoded);
    } else if (
        mobileidentity->imeisv.typeofidentity == MOBILE_IDENTITY_IMEISV) {
      encoded_rc = encode_imeisv_mobile_identity(
          &mobileidentity->imeisv, buffer + encoded, len - encoded);
    } else if (mobileidentity->tmsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      encoded_rc = encode_tmsi_mobile_identity(
          &mobileidentity->tmsi, buffer + encoded, len - encoded);
    } else if (mobileidentity->tmgi.typeofidentity == MOBILE_IDENTITY_TMGI) {
      encoded_rc = encode_tmgi_mobile_identity(
          &mobileidentity->tmgi, buffer + encoded, len - encoded);
    }

    if (encoded_rc > 0) {
      *lenPtr = encoded + encoded_rc - 1 - ((iei_present) ? 1 : 0);
    }
  } else {
    encoded_rc = encode_no_mobile_identity(
        &mobileidentity->no_id, buffer + encoded, len - encoded);

    if (encoded_rc > 0) {
      *lenPtr = MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH;
    }
  }

  if (encoded_rc < 0) {
    return encoded_rc;
  }

  return (encoded + encoded_rc);
}

//------------------------------------------------------------------------------
static int decode_imsi_mobile_identity(
    imsi_mobile_identity_t* imsi, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded   = 0;
  uint8_t ielen = 0;

  /* Pointing buffer to IE length field, to include the ieLen byte*/
  ielen = *(buffer + decoded);
  decoded++;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MOBILE_IDENTITY_IE_IMSI_LENGTH, len);
  imsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (imsi->typeofidentity != MOBILE_IDENTITY_IMSI) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  imsi->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  imsi->digit1  = (*(buffer + decoded) >> 4) & 0xf;
  imsi->numOfValidImsiDigits++;
  decoded++;
  if (decoded <= ielen) {
    imsi->digit2 = *(buffer + decoded) & 0xf;
    imsi->digit3 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit4 = *(buffer + decoded) & 0xf;
    imsi->digit5 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit5 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit6 = *(buffer + decoded) & 0xf;
    imsi->digit7 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit7 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit8 = *(buffer + decoded) & 0xf;
    imsi->digit9 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit9 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit10 = *(buffer + decoded) & 0xf;
    imsi->digit11 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit11 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit12 = *(buffer + decoded) & 0xf;
    imsi->digit13 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit13 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit14 = *(buffer + decoded) & 0xf;
    imsi->digit15 = (*(buffer + decoded) >> 4) & 0xf;
    if ((MOBILE_IDENTITY_EVEN == imsi->oddeven) && (imsi->digit15 != 0x0f)) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
    }
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  /*
   * IMSI is coded using BCD coding. If the number of identity digits is
   * even then bits 5 to 8 of the last octet shall be filled with an end
   * mark coded as "1111".
   */
  if ((imsi->oddeven == MOBILE_IDENTITY_EVEN) && (imsi->digit15 != 0x0f)) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded--; /*ielen is already included*/
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_imei_mobile_identity(
    imei_mobile_identity_t* imei, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MOBILE_IDENTITY_IE_IMEI_LENGTH, len);
  imei->typeofidentity = *(buffer + decoded) & 0x7;

  if (imei->typeofidentity != MOBILE_IDENTITY_IMEI) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  imei->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  imei->tac1    = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->tac2 = *(buffer + decoded) & 0xf;
  imei->tac3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->tac4 = *(buffer + decoded) & 0xf;
  imei->tac5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->tac6 = *(buffer + decoded) & 0xf;
  imei->tac7 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->tac8 = *(buffer + decoded) & 0xf;
  imei->snr1 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->snr2 = *(buffer + decoded) & 0xf;
  imei->snr3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->snr4 = *(buffer + decoded) & 0xf;
  imei->snr5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imei->snr6 = *(buffer + decoded) & 0xf;
  imei->cdsd = (*(buffer + decoded) >> 4) & 0xf;

  /*
   * IMEI is coded using BCD coding. If the number of identity digits is
   * even then bits 5 to 8 of the last octet shall be filled with an end
   * mark coded as "1111".
   */
  if ((imei->oddeven == MOBILE_IDENTITY_EVEN) && (imei->cdsd != 0x0f)) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_imeisv_mobile_identity(
    imeisv_mobile_identity_t* imeisv, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MOBILE_IDENTITY_IE_IMEISV_LENGTH, len);
  imeisv->typeofidentity = *(buffer + decoded) & 0x7;

  if (imeisv->typeofidentity != MOBILE_IDENTITY_IMEISV) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  imeisv->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  imeisv->tac1    = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->tac2 = *(buffer + decoded) & 0xf;
  imeisv->tac3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->tac4 = *(buffer + decoded) & 0xf;
  imeisv->tac5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->tac6 = *(buffer + decoded) & 0xf;
  imeisv->tac7 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->tac8 = *(buffer + decoded) & 0xf;
  imeisv->snr1 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->snr2 = *(buffer + decoded) & 0xf;
  imeisv->snr3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->snr4 = *(buffer + decoded) & 0xf;
  imeisv->snr5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->snr6 = *(buffer + decoded) & 0xf;
  imeisv->svn1 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  imeisv->svn2 = *(buffer + decoded) & 0xf;
  imeisv->last = (*(buffer + decoded) >> 4) & 0xf;

  /*
   * IMEISV is coded using BCD coding. If the number of identity digits is
   * even then bits 5 to 8 of the last octet shall be filled with an end
   * mark coded as "1111".
   */
  if ((imeisv->oddeven == MOBILE_IDENTITY_EVEN) && (imeisv->last != 0x0f)) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_tmsi_mobile_identity(
    tmsi_mobile_identity_t* tmsi, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MOBILE_IDENTITY_IE_TMSI_LENGTH, len);
  tmsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (tmsi->typeofidentity != MOBILE_IDENTITY_TMSI) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  tmsi->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  tmsi->f       = (*(buffer + decoded) >> 4) & 0xf;

  /*
   * If the mobile identity is the TMSI/P-TMSI/M-TMSI then bits 5 to 8
   * of octet 3 are coded as "1111".
   */
  if (tmsi->f != 0xf) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  tmsi->tmsi[0] = *(buffer + decoded);
  decoded++;
  tmsi->tmsi[1] = *(buffer + decoded);
  decoded++;
  tmsi->tmsi[2] = *(buffer + decoded);
  decoded++;
  tmsi->tmsi[3] = *(buffer + decoded);
  decoded++;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int decode_tmgi_mobile_identity(
    tmgi_mobile_identity_t* tmgi, uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MOBILE_IDENTITY_IE_TMGI_LENGTH, len);
  tmgi->spare = (*(buffer + decoded) >> 6) & 0x2;

  /*
   * Spare bits are coded with 0s
   */
  if (tmgi->spare != 0) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  tmgi->mbmssessionidindication = (*(buffer + decoded) >> 5) & 0x1;
  tmgi->mccmncindication        = (*(buffer + decoded) >> 4) & 0x1;
  tmgi->oddeven                 = (*(buffer + decoded) >> 3) & 0x1;
  tmgi->typeofidentity          = *(buffer + decoded) & 0x7;

  if (tmgi->typeofidentity != MOBILE_IDENTITY_TMGI) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  // IES_DECODE_U24(tmgi->mbmsserviceid, *(buffer + decoded));
  IES_DECODE_U24(buffer, decoded, tmgi->mbmsserviceid);
  tmgi->mccdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  tmgi->mccdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  tmgi->mncdigit3 = (*(buffer + decoded) >> 4) & 0xf;
  tmgi->mccdigit3 = *(buffer + decoded) & 0xf;
  decoded++;
  tmgi->mncdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  tmgi->mncdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  tmgi->mbmssessionid = *(buffer + decoded);
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
static int decode_no_mobile_identity(
    no_mobile_identity_t* no_id, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, 1, len);
  no_id->typeofidentity = *(buffer + decoded) & 0x7;

  if (no_id->typeofidentity != MOBILE_IDENTITY_NOT_AVAILABLE) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, TLV_VALUE_DOESNT_MATCH);
  }

  no_id->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  no_id->digit1  = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  if (len > 1) {
    no_id->digit2 = *(buffer + decoded) & 0xf;
    no_id->digit3 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    if (len > 2) {
      no_id->digit4 = *(buffer + decoded) & 0xf;
      no_id->digit5 = (*(buffer + decoded) >> 4) & 0xf;
      decoded++;
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded);
}

//------------------------------------------------------------------------------
static int encode_imsi_mobile_identity(
    imsi_mobile_identity_t* imsi, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | (imsi->digit1 << 4) | (imsi->oddeven << 3) |
                        (imsi->typeofidentity);
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit3 << 4) | imsi->digit2;
  encoded++;

  if (imsi->digit4 != 0xf) {
    *(buffer + encoded) = 0x00 | (imsi->digit5 << 4) | imsi->digit4;
  } else {
    *(buffer + encoded) = 0xf0 | imsi->digit4;
  }
  encoded++;
  if (imsi->digit6 != 0xf) {
    if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
      *(buffer + encoded) = 0x00 | (imsi->digit7 << 4) | imsi->digit6;
    } else {
      *(buffer + encoded) = 0xf0 | imsi->digit6;
    }
    encoded++;
    if (imsi->digit8 != 0xf) {
      if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
        *(buffer + encoded) = 0x00 | (imsi->digit9 << 4) | imsi->digit8;
      } else {
        *(buffer + encoded) = 0xf0 | imsi->digit8;
      }
      encoded++;
      if (imsi->digit10 != 0xf) {
        if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
          *(buffer + encoded) = 0x00 | (imsi->digit11 << 4) | imsi->digit10;
        } else {
          *(buffer + encoded) = 0xf0 | imsi->digit10;
        }
        encoded++;
        if (imsi->digit12 != 0xf) {
          if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
            *(buffer + encoded) = 0x00 | (imsi->digit13 << 4) | imsi->digit12;
          } else {
            *(buffer + encoded) = 0xf0 | imsi->digit12;
          }
          encoded++;
          if (imsi->digit14 != 0xf) {
            if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
              *(buffer + encoded) = 0x00 | (imsi->digit15 << 4) | imsi->digit14;
            } else {
              *(buffer + encoded) = 0xf0 | imsi->digit14;
            }
            encoded++;
          }
        }
      }
    }
  }
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_imei_mobile_identity(
    imei_mobile_identity_t* imei, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_IDENTITY_IE_IMEI_LENGTH, len);
  *(buffer + encoded) =
      (imei->tac1 << 4) | (imei->oddeven << 3) | (imei->typeofidentity);
  encoded++;
  *(buffer + encoded) = (imei->tac3 << 4) | imei->tac2;
  encoded++;
  *(buffer + encoded) = (imei->tac5 << 4) | imei->tac4;
  encoded++;
  *(buffer + encoded) = (imei->tac7 << 4) | imei->tac6;
  encoded++;
  *(buffer + encoded) = (imei->snr1 << 4) | imei->tac8;
  encoded++;
  *(buffer + encoded) = (imei->snr3 << 4) | imei->snr2;
  encoded++;
  *(buffer + encoded) = (imei->snr5 << 4) | imei->snr4;
  encoded++;

  if (imei->oddeven != MOBILE_IDENTITY_EVEN) {
    *(buffer + encoded) = (imei->cdsd << 4) | imei->snr6;
  } else {
    *(buffer + encoded) = 0xf0 | imei->snr6;
  }

  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_imeisv_mobile_identity(
    imeisv_mobile_identity_t* imeisv, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_IDENTITY_IE_IMEISV_LENGTH, len);
  *(buffer + encoded) =
      (imeisv->tac1 << 4) | (imeisv->oddeven << 3) | (imeisv->typeofidentity);
  encoded++;
  *(buffer + encoded) = (imeisv->tac3 << 4) | imeisv->tac2;
  encoded++;
  *(buffer + encoded) = (imeisv->tac5 << 4) | imeisv->tac4;
  encoded++;
  *(buffer + encoded) = (imeisv->tac7 << 4) | imeisv->tac6;
  encoded++;
  *(buffer + encoded) = (imeisv->snr1 << 4) | imeisv->tac8;
  encoded++;
  *(buffer + encoded) = (imeisv->snr3 << 4) | imeisv->snr2;
  encoded++;
  *(buffer + encoded) = (imeisv->snr5 << 4) | imeisv->snr4;
  encoded++;
  *(buffer + encoded) = (imeisv->svn1 << 4) | imeisv->snr6;
  encoded++;

  if (imeisv->oddeven != MOBILE_IDENTITY_EVEN) {
    *(buffer + encoded) = imeisv->last | imeisv->svn2;
  } else {
    *(buffer + encoded) = 0xf0 | imeisv->svn2;
  }

  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_tmsi_mobile_identity(
    tmsi_mobile_identity_t* tmsi, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_IDENTITY_IE_TMSI_LENGTH, len);
  *(buffer + encoded) = 0xf0 | (tmsi->oddeven << 3) | (tmsi->typeofidentity);
  encoded++;
  *(buffer + encoded) = tmsi->tmsi[0];
  encoded++;
  *(buffer + encoded) = tmsi->tmsi[1];
  encoded++;
  *(buffer + encoded) = tmsi->tmsi[2];
  encoded++;
  *(buffer + encoded) = tmsi->tmsi[3];
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_tmgi_mobile_identity(
    tmgi_mobile_identity_t* tmgi, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_IDENTITY_IE_TMGI_LENGTH, len);
  *(buffer + encoded) = ((tmgi->mbmssessionidindication & 0x1) << 5) |
                        ((tmgi->mccmncindication & 0x1) << 4) |
                        ((tmgi->oddeven & 0x1) << 3) |
                        (tmgi->typeofidentity & 0x7);
  encoded++;
  IES_ENCODE_U24(buffer, encoded, tmgi->mbmsserviceid);
  *(buffer + encoded) =
      ((tmgi->mccdigit2 & 0xf) << 4) | (tmgi->mccdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) =
      ((tmgi->mncdigit3 & 0xf) << 4) | (tmgi->mccdigit3 & 0xf);
  encoded++;
  *(buffer + encoded) =
      ((tmgi->mncdigit2 & 0xf) << 4) | (tmgi->mncdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) = tmgi->mbmssessionid;
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
static int encode_no_mobile_identity(
    no_mobile_identity_t* no_id, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, 1, len);
  *(buffer + encoded) = no_id->typeofidentity;
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.1.6 Mobile Station Classmark 2
//------------------------------------------------------------------------------
int decode_mobile_station_classmark_2_ie(
    mobile_station_classmark2_t* mobilestationclassmark2,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, MOBILE_STATION_CLASSMARK_2_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(C_MOBILE_STATION_CLASSMARK_2_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (MOBILE_STATION_CLASSMARK_2_IE_MAX_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  mobilestationclassmark2->revisionlevel     = (*(buffer + decoded) >> 5) & 0x3;
  mobilestationclassmark2->esind             = (*(buffer + decoded) >> 4) & 0x1;
  mobilestationclassmark2->a51               = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->rfpowercapability = *(buffer + decoded) & 0x7;
  decoded++;
  mobilestationclassmark2->pscapability      = (*(buffer + decoded) >> 6) & 0x1;
  mobilestationclassmark2->ssscreenindicator = (*(buffer + decoded) >> 4) & 0x3;
  mobilestationclassmark2->smcapability      = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->vbs               = (*(buffer + decoded) >> 2) & 0x1;
  mobilestationclassmark2->vgcs              = (*(buffer + decoded) >> 1) & 0x1;
  mobilestationclassmark2->fc                = *(buffer + decoded) & 0x1;
  decoded++;
  mobilestationclassmark2->cm3      = (*(buffer + decoded) >> 7) & 0x1;
  mobilestationclassmark2->lcsvacap = (*(buffer + decoded) >> 5) & 0x1;
  mobilestationclassmark2->ucs2     = (*(buffer + decoded) >> 4) & 0x1;
  mobilestationclassmark2->solsa    = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->cmsp     = (*(buffer + decoded) >> 2) & 0x1;
  mobilestationclassmark2->a53      = (*(buffer + decoded) >> 1) & 0x1;
  mobilestationclassmark2->a52      = *(buffer + decoded) & 0x1;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_mobile_station_classmark_2_ie(
    mobile_station_classmark2_t* mobilestationclassmark2,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, MOBILE_STATION_CLASSMARK_2_IE_MAX_LENGTH, len);
    *buffer = C_MOBILE_STATION_CLASSMARK_2_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (MOBILE_STATION_CLASSMARK_2_IE_MAX_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((mobilestationclassmark2->revisionlevel & 0x3) << 5) |
                        ((mobilestationclassmark2->esind & 0x1) << 4) |
                        ((mobilestationclassmark2->a51 & 0x1) << 3) |
                        (mobilestationclassmark2->rfpowercapability & 0x7);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((mobilestationclassmark2->pscapability & 0x1) << 6) |
      ((mobilestationclassmark2->ssscreenindicator & 0x3) << 4) |
      ((mobilestationclassmark2->smcapability & 0x1) << 3) |
      ((mobilestationclassmark2->vbs & 0x1) << 2) |
      ((mobilestationclassmark2->vgcs & 0x1) << 1) |
      (mobilestationclassmark2->fc & 0x1);
  encoded++;
  *(buffer + encoded) = 0x00 | ((mobilestationclassmark2->cm3 & 0x1) << 7) |
                        ((mobilestationclassmark2->lcsvacap & 0x1) << 5) |
                        ((mobilestationclassmark2->ucs2 & 0x1) << 4) |
                        ((mobilestationclassmark2->solsa & 0x1) << 3) |
                        ((mobilestationclassmark2->cmsp & 0x1) << 2) |
                        ((mobilestationclassmark2->a53 & 0x1) << 1) |
                        (mobilestationclassmark2->a52 & 0x1);
  encoded++;
  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.1.7 Mobile Station Classmark 3
//------------------------------------------------------------------------------

int decode_mobile_station_classmark_3_ie(
    mobile_station_classmark3_t* mobilestationclassmark3,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  // Temporary fix so that we decode other IEs required for CSFB
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei_present > 0) {
    CHECK_IEI_DECODER(C_MOBILE_STATION_CLASSMARK_3_IEI, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;

  decoded += ielen;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_mobile_station_classmark_3_ie(
    mobile_station_classmark3_t* mobilestationclassmark3,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  return 0;
}

//------------------------------------------------------------------------------
// 10.5.1.13 PLMN list
//------------------------------------------------------------------------------
int decode_plmn_list_ie(
    plmn_list_t* plmnlist, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;
  uint8_t i     = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (plmnlist->num_plmn * 3 + 2), len);
    CHECK_IEI_DECODER(C_PLMN_LIST_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (plmnlist->num_plmn * 3 + 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  plmnlist->num_plmn = 0;
  while (decoded < len) {
    plmnlist->plmn[i].mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
    plmnlist->plmn[i].mcc_digit1 = *(buffer + decoded) & 0xf;
    decoded++;
    plmnlist->plmn[i].mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
    plmnlist->plmn[i].mcc_digit3 = *(buffer + decoded) & 0xf;
    decoded++;
    plmnlist->plmn[i].mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
    plmnlist->plmn[i].mnc_digit1 = *(buffer + decoded) & 0xf;
    decoded++;
    plmnlist->num_plmn += 1;
    i += 1;
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_plmn_list_ie(
    plmn_list_t* plmnlist, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (2 + (plmnlist->num_plmn * 3)), len);
    *buffer = C_PLMN_LIST_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (1 + (plmnlist->num_plmn * 3)), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;
  for (int i = 0; i < plmnlist->num_plmn; i++) {
    *(buffer + encoded) = 0x00 | ((plmnlist->plmn[i].mcc_digit2 & 0xf) << 4) |
                          (plmnlist->plmn[i].mcc_digit1 & 0xf);
    encoded++;
    *(buffer + encoded) = 0x00 | ((plmnlist->plmn[i].mnc_digit3 & 0xf) << 4) |
                          (plmnlist->plmn[i].mcc_digit3 & 0xf);
    encoded++;
    *(buffer + encoded) = 0x00 | ((plmnlist->plmn[i].mnc_digit2 & 0xf) << 4) |
                          (plmnlist->plmn[i].mnc_digit1 & 0xf);
    encoded++;
  }
  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.1.15 MS network feature support
//------------------------------------------------------------------------------
int decode_ms_network_feature_support_ie(
    ms_network_feature_support_t* msnetworkfeaturesupport,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, MS_NETWORK_FEATURE_SUPPORT_IE_MAX_LENGTH, len);
  if (iei_present) {
    CHECK_IEI_DECODER(C_MS_NETWORK_FEATURE_SUPPORT_IEI, (*buffer & 0xc0));
  }
  msnetworkfeaturesupport->spare_bits = (*(buffer + decoded) >> 3) & 0x7;
  msnetworkfeaturesupport->extended_periodic_timers = *(buffer + decoded) & 0x1;
  decoded++;

  return decoded;
}

//------------------------------------------------------------------------------
int encode_ms_network_feature_support_ie(
    ms_network_feature_support_t* msnetworkfeaturesupport,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;
  /* Checking IEI and pointer */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MS_NETWORK_FEATURE_SUPPORT_IE_MAX_LENGTH, len);

  *(buffer + encoded) =
      0x00 | ((msnetworkfeaturesupport->spare_bits & 0x7) << 3) |
      (msnetworkfeaturesupport->extended_periodic_timers & 0x1);
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.5.31 Network Resource Identifier Container
//------------------------------------------------------------------------------

int decode_network_resource_identifier_container_ie(
    network_resource_identifier_container_t* networkresourceidentifiercontainer,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  // Temporary fix so that we decode other IEs
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei_present > 0) {
    CHECK_IEI_DECODER(C_NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IEI, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;

  decoded += ielen;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_network_resource_identifier_container_ie(
    network_resource_identifier_container_t* networkresourceidentifiercontainer,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  return 0;
}
