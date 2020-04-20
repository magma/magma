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
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "MobileIdentity.h"

static int decode_imsi_mobile_identity(
    ImsiMobileIdentity_t* imsi, uint8_t* buffer);
static int decode_imei_mobile_identity(
    ImeiMobileIdentity_t* imei, uint8_t* buffer);
static int decode_imeisv_mobile_identity(
    ImeisvMobileIdentity_t* imeisv, uint8_t* buffer);
static int decode_tmsi_mobile_identity(
    TmsiMobileIdentity_t* tmsi, uint8_t* buffer);
static int decode_tmgi_mobile_identity(
    TmgiMobileIdentity_t* tmgi, uint8_t* buffer);
static int decode_no_mobile_identity(
    NoMobileIdentity_t* no_id, uint8_t* buffer);

static int encode_imsi_mobile_identity(
    ImsiMobileIdentity_t* imsi, uint8_t* buffer);
static int encode_imei_mobile_identity(
    ImeiMobileIdentity_t* imei, uint8_t* buffer);
static int encode_imeisv_mobile_identity(
    ImeisvMobileIdentity_t* imeisv, uint8_t* buffer);
static int encode_tmsi_mobile_identity(
    TmsiMobileIdentity_t* tmsi, uint8_t* buffer);
static int encode_tmgi_mobile_identity(
    TmgiMobileIdentity_t* tmgi, uint8_t* buffer);
static int encode_no_mobile_identity(
    NoMobileIdentity_t* no_id, uint8_t* buffer);

int decode_mobile_identity(
    MobileIdentity* mobileidentity, uint8_t iei, uint8_t* buffer,
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
  uint8_t typeofidentity = *(buffer + decoded) & 0x7;

  if (typeofidentity != MOBILE_IDENTITY_NOT_AVAILABLE) {
    CHECK_LENGTH_DECODER(len - decoded, ielen);

    if (typeofidentity == MOBILE_IDENTITY_IMSI) {
      decoded_rc = decode_imsi_mobile_identity(&mobileidentity->imsi, buffer);
    } else if (typeofidentity == MOBILE_IDENTITY_IMEI) {
      decoded_rc =
          decode_imei_mobile_identity(&mobileidentity->imei, buffer + decoded);
    } else if (typeofidentity == MOBILE_IDENTITY_IMEISV) {
      decoded_rc = decode_imeisv_mobile_identity(
          &mobileidentity->imeisv, buffer + decoded);
    } else if (typeofidentity == MOBILE_IDENTITY_TMSI) {
      decoded_rc =
          decode_tmsi_mobile_identity(&mobileidentity->tmsi, buffer + decoded);
    } else if (typeofidentity == MOBILE_IDENTITY_TMGI) {
      decoded_rc =
          decode_tmgi_mobile_identity(&mobileidentity->tmgi, buffer + decoded);
    }
  } else if (ielen == MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH) {
    decoded_rc =
        decode_no_mobile_identity(&mobileidentity->no_id, buffer + decoded);
  }

  if (decoded_rc < 0) {
    return decoded_rc;
  }
#if NAS_DEBUG
  dump_mobile_identity_xml(mobileidentity, iei);
#endif
  return (decoded + decoded_rc);
}

int encode_mobile_identity(
    MobileIdentity* mobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  int encoded_rc   = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_IDENTITY_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_mobile_identity_xml(mobileidentity, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if (mobileidentity->no_id.typeofidentity != MOBILE_IDENTITY_NOT_AVAILABLE) {
    if (mobileidentity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      encoded_rc =
          encode_imsi_mobile_identity(&mobileidentity->imsi, buffer + encoded);
    } else if (mobileidentity->imei.typeofidentity == MOBILE_IDENTITY_IMEI) {
      encoded_rc =
          encode_imei_mobile_identity(&mobileidentity->imei, buffer + encoded);
    } else if (
        mobileidentity->imeisv.typeofidentity == MOBILE_IDENTITY_IMEISV) {
      encoded_rc = encode_imeisv_mobile_identity(
          &mobileidentity->imeisv, buffer + encoded);
    } else if (mobileidentity->tmsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      encoded_rc =
          encode_tmsi_mobile_identity(&mobileidentity->tmsi, buffer + encoded);
    } else if (mobileidentity->tmgi.typeofidentity == MOBILE_IDENTITY_TMGI) {
      encoded_rc =
          encode_tmgi_mobile_identity(&mobileidentity->tmgi, buffer + encoded);
    }

    if (encoded_rc > 0) {
      *lenPtr = encoded + encoded_rc - 1 - ((iei > 0) ? 1 : 0);
    }
  } else {
    encoded_rc =
        encode_no_mobile_identity(&mobileidentity->no_id, buffer + encoded);

    if (encoded_rc > 0) {
      *lenPtr = MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH;
    }
  }

  if (encoded_rc < 0) {
    return encoded_rc;
  }

  return (encoded + encoded_rc);
}

void dump_mobile_identity_xml(MobileIdentity* mobileidentity, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Mobile Identity>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  if (mobileidentity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
    ImsiMobileIdentity_t* imsi = &mobileidentity->imsi;

    OAILOG_DEBUG(LOG_NAS, "    <odd even>%u</odd even>\n", imsi->oddeven);
    OAILOG_DEBUG(LOG_NAS, "    <Type of identity>IMSI</Type of identity>\n");
    OAILOG_DEBUG(LOG_NAS, "    <digit1>%u</digit1>\n", imsi->digit1);
    OAILOG_DEBUG(LOG_NAS, "    <digit2>%u</digit2>\n", imsi->digit2);
    OAILOG_DEBUG(LOG_NAS, "    <digit3>%u</digit3>\n", imsi->digit3);
    OAILOG_DEBUG(LOG_NAS, "    <digit4>%u</digit4>\n", imsi->digit4);
    OAILOG_DEBUG(LOG_NAS, "    <digit5>%u</digit5>\n", imsi->digit5);
    OAILOG_DEBUG(LOG_NAS, "    <digit6>%u</digit6>\n", imsi->digit6);
    OAILOG_DEBUG(LOG_NAS, "    <digit7>%u</digit7>\n", imsi->digit7);
    OAILOG_DEBUG(LOG_NAS, "    <digit8>%u</digit8>\n", imsi->digit8);
    OAILOG_DEBUG(LOG_NAS, "    <digit9>%u</digit9>\n", imsi->digit9);
    OAILOG_DEBUG(LOG_NAS, "    <digit10>%u</digit10>\n", imsi->digit10);
    OAILOG_DEBUG(LOG_NAS, "    <digit11>%u</digit11>\n", imsi->digit11);
    OAILOG_DEBUG(LOG_NAS, "    <digit12>%u</digit12>\n", imsi->digit12);
    OAILOG_DEBUG(LOG_NAS, "    <digit13>%u</digit13>\n", imsi->digit13);
    OAILOG_DEBUG(LOG_NAS, "    <digit14>%u</digit14>\n", imsi->digit14);
    OAILOG_DEBUG(LOG_NAS, "    <digit15>%u</digit15>\n", imsi->digit15);
  } else if (mobileidentity->imei.typeofidentity == MOBILE_IDENTITY_IMEI) {
    ImeiMobileIdentity_t* imei = &mobileidentity->imei;

    OAILOG_DEBUG(LOG_NAS, "    <odd even>%u</odd even>\n", imei->oddeven);
    OAILOG_DEBUG(LOG_NAS, "    <Type of identity>IMEI</Type of identity>\n");
    OAILOG_DEBUG(LOG_NAS, "    <tac1>%u</tac1>\n", imei->tac1);
    OAILOG_DEBUG(LOG_NAS, "    <tac2>%u</tac2>\n", imei->tac2);
    OAILOG_DEBUG(LOG_NAS, "    <tac3>%u</tac3>\n", imei->tac3);
    OAILOG_DEBUG(LOG_NAS, "    <tac4>%u</tac4>\n", imei->tac4);
    OAILOG_DEBUG(LOG_NAS, "    <tac5>%u</tac5>\n", imei->tac5);
    OAILOG_DEBUG(LOG_NAS, "    <tac6>%u</tac6>\n", imei->tac6);
    OAILOG_DEBUG(LOG_NAS, "    <tac7>%u</tac7>\n", imei->tac7);
    OAILOG_DEBUG(LOG_NAS, "    <tac8>%u</tac8>\n", imei->tac8);
    OAILOG_DEBUG(LOG_NAS, "    <snr1>%u</snr1>\n", imei->snr1);
    OAILOG_DEBUG(LOG_NAS, "    <snr2>%u</snr2>\n", imei->snr2);
    OAILOG_DEBUG(LOG_NAS, "    <snr3>%u</snr3>\n", imei->snr3);
    OAILOG_DEBUG(LOG_NAS, "    <snr4>%u</snr4>\n", imei->snr4);
    OAILOG_DEBUG(LOG_NAS, "    <snr5>%u</snr5>\n", imei->snr5);
    OAILOG_DEBUG(LOG_NAS, "    <snr6>%u</snr6>\n", imei->snr6);
    OAILOG_DEBUG(LOG_NAS, "    <cdsd>%u</cdsd>\n", imei->cdsd);
  } else if (mobileidentity->imeisv.typeofidentity == MOBILE_IDENTITY_IMEISV) {
    ImeisvMobileIdentity_t* imeisv = &mobileidentity->imeisv;

    OAILOG_DEBUG(LOG_NAS, "    <odd even>%u</odd even>\n", imeisv->oddeven);
    OAILOG_DEBUG(LOG_NAS, "    <Type of identity>IMEISV</Type of identity>\n");
    OAILOG_DEBUG(LOG_NAS, "    <tac1>%u</tac1>\n", imeisv->tac1);
    OAILOG_DEBUG(LOG_NAS, "    <tac2>%u</tac2>\n", imeisv->tac2);
    OAILOG_DEBUG(LOG_NAS, "    <tac3>%u</tac3>\n", imeisv->tac3);
    OAILOG_DEBUG(LOG_NAS, "    <tac4>%u</tac4>\n", imeisv->tac4);
    OAILOG_DEBUG(LOG_NAS, "    <tac5>%u</tac5>\n", imeisv->tac5);
    OAILOG_DEBUG(LOG_NAS, "    <tac6>%u</tac6>\n", imeisv->tac6);
    OAILOG_DEBUG(LOG_NAS, "    <tac7>%u</tac7>\n", imeisv->tac7);
    OAILOG_DEBUG(LOG_NAS, "    <tac8>%u</tac8>\n", imeisv->tac8);
    OAILOG_DEBUG(LOG_NAS, "    <snr1>%u</snr1>\n", imeisv->snr1);
    OAILOG_DEBUG(LOG_NAS, "    <snr2>%u</snr2>\n", imeisv->snr2);
    OAILOG_DEBUG(LOG_NAS, "    <snr3>%u</snr3>\n", imeisv->snr3);
    OAILOG_DEBUG(LOG_NAS, "    <snr4>%u</snr4>\n", imeisv->snr4);
    OAILOG_DEBUG(LOG_NAS, "    <snr5>%u</snr5>\n", imeisv->snr5);
    OAILOG_DEBUG(LOG_NAS, "    <snr6>%u</snr6>\n", imeisv->snr6);
    OAILOG_DEBUG(LOG_NAS, "    <svn1>%u</svn1>\n", imeisv->svn1);
    OAILOG_DEBUG(LOG_NAS, "    <svn2>%u</svn2>\n", imeisv->svn2);
  } else if (mobileidentity->tmsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
    TmsiMobileIdentity_t* tmsi = &mobileidentity->tmsi;

    OAILOG_DEBUG(LOG_NAS, "    <odd even>%u</odd even>\n", tmsi->oddeven);
    OAILOG_DEBUG(LOG_NAS, "    <Type of identity>TMSI</Type of identity>\n");
    OAILOG_DEBUG(LOG_NAS, "    <digit1>%u</digit1>\n", tmsi->digit1);
    OAILOG_DEBUG(LOG_NAS, "    <digit2>%u</digit2>\n", tmsi->digit2);
    OAILOG_DEBUG(LOG_NAS, "    <digit3>%u</digit3>\n", tmsi->digit3);
    OAILOG_DEBUG(LOG_NAS, "    <digit4>%u</digit4>\n", tmsi->digit4);
    OAILOG_DEBUG(LOG_NAS, "    <digit5>%u</digit5>\n", tmsi->digit5);
    OAILOG_DEBUG(LOG_NAS, "    <digit6>%u</digit6>\n", tmsi->digit6);
    OAILOG_DEBUG(LOG_NAS, "    <digit7>%u</digit7>\n", tmsi->digit7);
    OAILOG_DEBUG(LOG_NAS, "    <digit8>%u</digit8>\n", tmsi->digit8);
    OAILOG_DEBUG(LOG_NAS, "    <digit9>%u</digit9>\n", tmsi->digit9);
    OAILOG_DEBUG(LOG_NAS, "    <digit10>%u</digit10>\n", tmsi->digit10);
    OAILOG_DEBUG(LOG_NAS, "    <digit11>%u</digit11>\n", tmsi->digit11);
    OAILOG_DEBUG(LOG_NAS, "    <digit12>%u</digit12>\n", tmsi->digit12);
    OAILOG_DEBUG(LOG_NAS, "    <digit13>%u</digit13>\n", tmsi->digit13);
    OAILOG_DEBUG(LOG_NAS, "    <digit14>%u</digit14>\n", tmsi->digit14);
    OAILOG_DEBUG(LOG_NAS, "    <digit15>%u</digit15>\n", tmsi->digit15);
  } else if (mobileidentity->tmgi.typeofidentity == MOBILE_IDENTITY_TMGI) {
    TmgiMobileIdentity_t* tmgi = &mobileidentity->tmgi;

    OAILOG_DEBUG(
        LOG_NAS,
        "    <MBMS session ID indication>%u</MBMS session ID indication>\n",
        tmgi->mbmssessionidindication);
    OAILOG_DEBUG(
        LOG_NAS, "    <MCC MNC indication>%u</MCC MNC indication>\n",
        tmgi->mccmncindication);
    OAILOG_DEBUG(LOG_NAS, "    <Odd even>%u</Odd even>\n", tmgi->oddeven);
    OAILOG_DEBUG(LOG_NAS, "    <Type of identity>TMGI</Type of identity>\n");
    OAILOG_DEBUG(
        LOG_NAS, "    <MBMS service ID>%u</MBMS service ID>\n",
        tmgi->mbmsserviceid);
    OAILOG_DEBUG(
        LOG_NAS, "    <MCC digit 2>%u</MCC digit 2>\n", tmgi->mccdigit2);
    OAILOG_DEBUG(
        LOG_NAS, "    <MCC digit 1>%u</MCC digit 1>\n", tmgi->mccdigit1);
    OAILOG_DEBUG(
        LOG_NAS, "    <MNC digit 3>%u</MNC digit 3>\n", tmgi->mncdigit3);
    OAILOG_DEBUG(
        LOG_NAS, "    <MCC digit 3>%u</MCC digit 3>\n", tmgi->mccdigit3);
    OAILOG_DEBUG(
        LOG_NAS, "    <MNC digit 2>%u</MNC digit 2>\n", tmgi->mncdigit2);
    OAILOG_DEBUG(
        LOG_NAS, "    <MNC digit 1>%u</MNC digit 1>\n", tmgi->mncdigit1);
    OAILOG_DEBUG(
        LOG_NAS, "    <MBMS session ID>%u</MBMS session ID>\n",
        tmgi->mbmssessionid);
  } else {
    OAILOG_DEBUG(
        LOG_NAS, "    Wrong type of mobile identity (%u)\n",
        mobileidentity->imsi.typeofidentity);
  }

  OAILOG_DEBUG(LOG_NAS, "</Mobile Identity>\n");
}

static int decode_imsi_mobile_identity(
    ImsiMobileIdentity_t* imsi, uint8_t* buffer) {
  int decoded   = 0;
  uint8_t ielen = 0;

  ielen = *(buffer + decoded); /* Pointing buffer to IE length field, to include
                                  the ieLen byte*/
  decoded++;
  imsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (imsi->typeofidentity != MOBILE_IDENTITY_IMSI) {
    return (TLV_VALUE_DOESNT_MATCH);
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
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit6 = *(buffer + decoded) & 0xf;
    imsi->digit7 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit8 = *(buffer + decoded) & 0xf;
    imsi->digit9 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit10 = *(buffer + decoded) & 0xf;
    imsi->digit11 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit12 = *(buffer + decoded) & 0xf;
    imsi->digit13 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }
  if (decoded <= ielen) {
    imsi->digit14 = *(buffer + decoded) & 0xf;
    imsi->digit15 = (*(buffer + decoded) >> 4) & 0xf;
    decoded++;
    imsi->numOfValidImsiDigits += 2;
  }

  if (imsi->oddeven == false) {
    imsi->numOfValidImsiDigits--; /* For even number of digits*/
  }
  /*
   * IMSI is coded using BCD coding. If the number of identity digits is
   * even then bits 5 to 8 of the last octet shall be filled with an end
   * mark coded as "1111".
   */
  if ((imsi->oddeven == MOBILE_IDENTITY_EVEN) && (imsi->digit15 != 0x0f)) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded--; /*ielen is already included*/
  return decoded;
}

static int decode_imei_mobile_identity(
    ImeiMobileIdentity_t* imei, uint8_t* buffer) {
  int decoded = 0;

  imei->typeofidentity = *(buffer + decoded) & 0x7;

  if (imei->typeofidentity != MOBILE_IDENTITY_IMEI) {
    return (TLV_VALUE_DOESNT_MATCH);
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
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  return decoded;
}

static int decode_imeisv_mobile_identity(
    ImeisvMobileIdentity_t* imeisv, uint8_t* buffer) {
  int decoded = 0;

  imeisv->typeofidentity = *(buffer + decoded) & 0x7;

  if (imeisv->typeofidentity != MOBILE_IDENTITY_IMEISV) {
    return (TLV_VALUE_DOESNT_MATCH);
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
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  return decoded;
}

static int decode_tmsi_mobile_identity(
    TmsiMobileIdentity_t* tmsi, uint8_t* buffer) {
  int decoded = 0;

  tmsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (tmsi->typeofidentity != MOBILE_IDENTITY_TMSI) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  tmsi->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  tmsi->digit1  = (*(buffer + decoded) >> 4) & 0xf;

  /*
   * If the mobile identity is the TMSI/P-TMSI/M-TMSI then bits 5 to 8
   * of octet 3 are coded as "1111".
   */
  if (tmsi->digit1 != 0xf) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  tmsi->digit2 = *(buffer + decoded) & 0xf;
  tmsi->digit3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  tmsi->digit4 = *(buffer + decoded) & 0xf;
  tmsi->digit5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  tmsi->digit6 = *(buffer + decoded) & 0xf;
  tmsi->digit7 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  tmsi->digit8 = *(buffer + decoded) & 0xf;
  tmsi->digit9 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  return decoded;
}

static int decode_tmgi_mobile_identity(
    TmgiMobileIdentity_t* tmgi, uint8_t* buffer) {
  int decoded = 0;

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

static int decode_no_mobile_identity(
    NoMobileIdentity_t* no_id, uint8_t* buffer) {
  int decoded = 0;

  no_id->typeofidentity = *(buffer + decoded) & 0x7;

  if (no_id->typeofidentity != MOBILE_IDENTITY_NOT_AVAILABLE) {
    return (TLV_VALUE_DOESNT_MATCH);
  }

  no_id->oddeven = (*(buffer + decoded) >> 3) & 0x1;
  no_id->digit1  = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit2 = *(buffer + decoded) & 0xf;
  no_id->digit3 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit4 = *(buffer + decoded) & 0xf;
  no_id->digit5 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit6 = *(buffer + decoded) & 0xf;
  no_id->digit7 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit8 = *(buffer + decoded) & 0xf;
  no_id->digit9 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit10 = *(buffer + decoded) & 0xf;
  no_id->digit11 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit12 = *(buffer + decoded) & 0xf;
  no_id->digit13 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  no_id->digit14 = *(buffer + decoded) & 0xf;
  no_id->digit15 = (*(buffer + decoded) >> 4) & 0xf;
  decoded++;
  return decoded;
}

static int encode_imsi_mobile_identity(
    ImsiMobileIdentity_t* imsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | (imsi->digit1 << 4) | (imsi->oddeven << 3) |
                        (imsi->typeofidentity);
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit3 << 4) | imsi->digit2;
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit5 << 4) | imsi->digit4;
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit7 << 4) | imsi->digit6;
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit9 << 4) | imsi->digit8;
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit11 << 4) | imsi->digit10;
  encoded++;
  *(buffer + encoded) = 0x00 | (imsi->digit13 << 4) | imsi->digit12;
  encoded++;

  if (imsi->oddeven != MOBILE_IDENTITY_EVEN) {
    *(buffer + encoded) = 0x00 | (imsi->digit15 << 4) | imsi->digit14;
  } else {
    *(buffer + encoded) = 0xf0 | imsi->digit14;
  }

  encoded++;
  return encoded;
}

static int encode_imei_mobile_identity(
    ImeiMobileIdentity_t* imei, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) =
      0x00 | (imei->tac1 << 4) | (imei->oddeven << 3) | (imei->typeofidentity);
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->tac3 << 4) | imei->tac2;
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->tac5 << 4) | imei->tac4;
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->tac7 << 4) | imei->tac6;
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->snr1 << 4) | imei->tac8;
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->snr3 << 4) | imei->snr2;
  encoded++;
  *(buffer + encoded) = 0x00 | (imei->snr5 << 4) | imei->snr4;
  encoded++;

  if (imei->oddeven != MOBILE_IDENTITY_EVEN) {
    *(buffer + encoded) = 0x00 | (imei->cdsd << 4) | imei->snr6;
  } else {
    *(buffer + encoded) = 0xf0 | imei->snr6;
  }

  encoded++;
  return encoded;
}

static int encode_imeisv_mobile_identity(
    ImeisvMobileIdentity_t* imeisv, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | (imeisv->tac1 << 4) | (imeisv->oddeven << 3) |
                        (imeisv->typeofidentity);
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->tac3 << 4) | imeisv->tac2;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->tac5 << 4) | imeisv->tac4;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->tac7 << 4) | imeisv->tac6;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->snr1 << 4) | imeisv->tac8;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->snr3 << 4) | imeisv->snr2;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->snr5 << 4) | imeisv->snr4;
  encoded++;
  *(buffer + encoded) = 0x00 | (imeisv->svn1 << 4) | imeisv->snr6;
  encoded++;

  if (imeisv->oddeven != MOBILE_IDENTITY_EVEN) {
    *(buffer + encoded) = imeisv->last | imeisv->svn2;
  } else {
    *(buffer + encoded) = 0xf0 | imeisv->svn2;
  }

  encoded++;
  return encoded;
}

static int encode_tmsi_mobile_identity(
    TmsiMobileIdentity_t* tmsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0xf0 | (tmsi->oddeven << 3) | (tmsi->typeofidentity);
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit3 << 4) | tmsi->digit2;
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit5 << 4) | tmsi->digit4;
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit7 << 4) | tmsi->digit6;
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit9 << 4) | tmsi->digit8;
  encoded++;
  /*Below code is not required as TMSI will be 4 bytes which is 9 digits*/
  /*  *(buffer + encoded) = 0x00 | (tmsi->digit11 << 4) | tmsi->digit10;
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit13 << 4) | tmsi->digit12;
  encoded++;
  *(buffer + encoded) = 0x00 | (tmsi->digit15 << 4) | tmsi->digit14;
  encoded++;*/
  return encoded;
}

static int encode_tmgi_mobile_identity(
    TmgiMobileIdentity_t* tmgi, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | ((tmgi->mbmssessionidindication & 0x1) << 5) |
                        ((tmgi->mccmncindication & 0x1) << 4) |
                        ((tmgi->oddeven & 0x1) << 3) |
                        (tmgi->typeofidentity & 0x7);
  encoded++;
  IES_ENCODE_U24(buffer, encoded, tmgi->mbmsserviceid);
  *(buffer + encoded) =
      0x00 | ((tmgi->mccdigit2 & 0xf) << 4) | (tmgi->mccdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((tmgi->mncdigit3 & 0xf) << 4) | (tmgi->mccdigit3 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((tmgi->mncdigit2 & 0xf) << 4) | (tmgi->mncdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) = tmgi->mbmssessionid;
  encoded++;
  return encoded;
}

static int encode_no_mobile_identity(
    NoMobileIdentity_t* no_id, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = no_id->typeofidentity;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  *(buffer + encoded) = 0x00;
  encoded++;
  return encoded;
}
