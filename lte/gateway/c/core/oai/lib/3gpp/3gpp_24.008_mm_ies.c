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

/*! \file 3gpp_24.008_mm_ies.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdbool.h>
#include <stdint.h>

#include "bstrlib.h"
#include "common_defs.h"
#include "assertions.h"
#include "3gpp_24.008.h"
#include "TLVDecoder.h"
#include "TLVEncoder.h"
#include "log.h"

//******************************************************************************
// 10.5.3 Mobility management information elements.
//******************************************************************************
//------------------------------------------------------------------------------
// 10.5.3.1 Authentication parameter RAND
//------------------------------------------------------------------------------
int decode_authentication_parameter_rand_ie(
    authentication_parameter_rand_t* authenticationparameterrand,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS);
  int decoded   = 0;
  uint8_t ielen = 16;
  int decode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(MM_AUTHENTICATION_PARAMETER_RAND_IEI, *buffer);
    decoded++;
  } else {
    // Type 3
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH - 1), len);
  }

  if ((decode_result = decode_bstring(
           authenticationparameterrand, ielen, buffer + decoded,
           len - decoded)) < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, decode_result);
  } else
    decoded += decode_result;

  OAILOG_FUNC_RETURN(LOG_NAS, decoded);
}

//------------------------------------------------------------------------------
int encode_authentication_parameter_rand_ie(
    authentication_parameter_rand_t authenticationparameterrand,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint32_t encode_result;
  uint32_t encoded = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH, len);
    *buffer = MM_AUTHENTICATION_PARAMETER_RAND_IEI;
    encoded++;
  } else {
    // Type 4
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH - 2), len);
  }

  if ((encode_result = encode_bstring(
           authenticationparameterrand, buffer + encoded, len - encoded)) < 0) {
    return encode_result;
  } else {
    encoded += encode_result;
  }

  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.1.1 Authentication Parameter AUTN (UMTS and EPS authentication
// challenge)
//------------------------------------------------------------------------------
int decode_authentication_parameter_autn_ie(
    authentication_parameter_autn_t* authenticationparameterautn,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS);
  int decoded   = 0;
  uint8_t ielen = 0;
  int decode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(MM_AUTHENTICATION_PARAMETER_AUTN_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if ((decode_result = decode_bstring(
           authenticationparameterautn, ielen, buffer + decoded,
           len - decoded)) < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, decode_result);
  } else
    decoded += decode_result;

  OAILOG_FUNC_RETURN(LOG_NAS, decoded);
}

//------------------------------------------------------------------------------
int encode_authentication_parameter_autn_ie(
    authentication_parameter_autn_t authenticationparameterautn,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  int encode_result;
  uint32_t encoded = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH, len);
    *buffer = MM_AUTHENTICATION_PARAMETER_AUTN_IEI;  // ???
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if ((encode_result = encode_bstring(
           authenticationparameterautn, buffer + encoded, len - encoded)) < 0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.2 Authentication Response parameter
//------------------------------------------------------------------------------
int decode_authentication_response_parameter_ie(
    authentication_response_parameter_t* authenticationresponseparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS);
  int decoded   = 0;
  uint8_t ielen = 0;
  int decode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(MM_AUTHENTICATION_RESPONSE_PARAMETER_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if ((decode_result = decode_bstring(
           authenticationresponseparameter, ielen, buffer + decoded,
           len - decoded)) < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, decode_result);
  } else
    decoded += decode_result;

  OAILOG_FUNC_RETURN(LOG_NAS, decoded);
}

//------------------------------------------------------------------------------
int encode_authentication_response_parameter_ie(
    authentication_response_parameter_t authenticationresponseparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  int encode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH, len);
    *buffer = MM_AUTHENTICATION_RESPONSE_PARAMETER_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if ((encode_result = encode_bstring(
           authenticationresponseparameter, buffer + encoded, len - encoded)) <
      0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.2.2 Authentication Failure parameter (UMTS and EPS authentication
// challenge)
//------------------------------------------------------------------------------
int decode_authentication_failure_parameter_ie(
    authentication_failure_parameter_t* authenticationfailureparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS);
  int decoded   = 0;
  uint8_t ielen = 0;
  int decode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, AUTHENTICATION_FAILURE_PARAMETER_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(MM_AUTHENTICATION_FAILURE_PARAMETER_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (AUTHENTICATION_FAILURE_PARAMETER_IE_MAX_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if ((decode_result = decode_bstring(
           authenticationfailureparameter, ielen, buffer + decoded,
           len - decoded)) < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS, decode_result);
  } else
    decoded += decode_result;

  OAILOG_FUNC_RETURN(LOG_NAS, decoded);
}

//------------------------------------------------------------------------------
int encode_authentication_failure_parameter_ie(
    authentication_failure_parameter_t authenticationfailureparameter,
    const bool iei_present, uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  int encode_result;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, AUTHENTICATION_FAILURE_PARAMETER_IE_MAX_LENGTH, len);
    *buffer = MM_AUTHENTICATION_FAILURE_PARAMETER_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (AUTHENTICATION_FAILURE_PARAMETER_IE_MAX_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if ((encode_result = encode_bstring(
           authenticationfailureparameter, buffer + encoded, len - encoded)) <
      0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.5a Network Name
//------------------------------------------------------------------------------
int decode_network_name_ie(
    network_name_t* networkname, const uint8_t iei, uint8_t* buffer,
    const uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS);
  int decoded   = 0;
  uint8_t ielen = 0;
  int decode_result;

  if (iei > 0) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, NETWORK_NAME_IE_MIN_LENGTH, len);
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (NETWORK_NAME_IE_MIN_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if (((*buffer >> 7) & 0x1) != 1) {
    errorCodeDecoder = TLV_VALUE_DOESNT_MATCH;
    return TLV_VALUE_DOESNT_MATCH;
  }

  networkname->codingscheme                 = (*(buffer + decoded) >> 5) & 0x7;
  networkname->addci                        = (*(buffer + decoded) >> 4) & 0x1;
  networkname->numberofsparebitsinlastoctet = (*(buffer + decoded) >> 1) & 0x7;

  if ((decode_result = decode_bstring(
           &networkname->textstring, ielen, buffer + decoded, len - decoded)) <
      0) {
    OAILOG_FUNC_RETURN(LOG_NAS, decode_result);
  } else
    decoded += decode_result;
  OAILOG_FUNC_RETURN(LOG_NAS, decoded);
}

//------------------------------------------------------------------------------
int encode_network_name_ie(
    network_name_t* networkname, const uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  int encode_result;

  if (iei > 0) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer,
        NETWORK_NAME_IE_MIN_LENGTH + blength(networkname->textstring) - 1, len);
    *buffer = iei;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer,
        NETWORK_NAME_IE_MIN_LENGTH + blength(networkname->textstring) - 2, len);
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 | (1 << 7) |
                        ((networkname->codingscheme & 0x7) << 4) |
                        ((networkname->addci & 0x1) << 3) |
                        (networkname->numberofsparebitsinlastoctet & 0x7);
  encoded++;

  if ((encode_result = encode_bstring(
           networkname->textstring, buffer + encoded, len - encoded)) < 0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.8 Time Zone
//------------------------------------------------------------------------------
int encode_time_zone(
    time_zone_t* timezone, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, TIME_ZONE_IE_MAX_LENGTH, len);

  if (iei_present) {
    *buffer = MM_TIME_ZONE_IEI;
    encoded++;
  }

  *(buffer + encoded) = *timezone;
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
int decode_time_zone(
    time_zone_t* timezone, const bool iei_present, uint8_t* buffer,
    const uint32_t len) {
  int decoded = 0;

  if (iei_present) {
    CHECK_IEI_DECODER(MM_TIME_ZONE_IEI, *buffer);
    decoded++;
  }

  *timezone = *(buffer + decoded);
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
// 10.5.3.9 Time Zone and Time
//------------------------------------------------------------------------------
int encode_time_zone_and_time(
    time_zone_and_time_t* timezoneandtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TIME_ZONE_AND_TIME_MAX_LENGTH, len);

  if (iei_present) {
    *buffer = MM_TIME_ZONE_AND_TIME_IEI;
    encoded++;
  }

  *(buffer + encoded) = timezoneandtime->year;
  encoded++;
  *(buffer + encoded) = timezoneandtime->month;
  encoded++;
  *(buffer + encoded) = timezoneandtime->day;
  encoded++;
  *(buffer + encoded) = timezoneandtime->hour;
  encoded++;
  *(buffer + encoded) = timezoneandtime->minute;
  encoded++;
  *(buffer + encoded) = timezoneandtime->second;
  encoded++;
  *(buffer + encoded) = timezoneandtime->timezone;
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
int decode_time_zone_and_time(
    time_zone_and_time_t* timezoneandtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  if (iei_present) {
    CHECK_IEI_DECODER(MM_TIME_ZONE_AND_TIME_IEI, *buffer);
    decoded++;
  }

  timezoneandtime->year = *(buffer + decoded);
  decoded++;
  timezoneandtime->month = *(buffer + decoded);
  decoded++;
  timezoneandtime->day = *(buffer + decoded);
  decoded++;
  timezoneandtime->hour = *(buffer + decoded);
  decoded++;
  timezoneandtime->minute = *(buffer + decoded);
  decoded++;
  timezoneandtime->second = *(buffer + decoded);
  decoded++;
  timezoneandtime->timezone = *(buffer + decoded);
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
// 10.5.3.12 Daylight Saving Time
//------------------------------------------------------------------------------
int decode_daylight_saving_time_ie(
    daylight_saving_time_t* daylightsavingtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH, len);
    CHECK_IEI_DECODER(MM_DAYLIGHT_SAVING_TIME_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  *daylightsavingtime = *buffer & 0x3;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_daylight_saving_time_ie(
    daylight_saving_time_t* daylightsavingtime, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH, len);
    *buffer = MM_DAYLIGHT_SAVING_TIME_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 | (*daylightsavingtime & 0x3);
  encoded++;
  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}

//------------------------------------------------------------------------------
// 10.5.3.13 Emergency Number List
//------------------------------------------------------------------------------
int decode_emergency_number_list_ie(
    emergency_number_list_t* emergencynumberlist, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  int decoded                = 0;
  uint8_t ielen              = 0;
  emergency_number_list_t* e = emergencynumberlist;

  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, EMERGENCY_NUMBER_LIST_IE_MIN_LENGTH, len);
    CHECK_IEI_DECODER(MM_EMERGENCY_NUMBER_LIST_IEI, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, (EMERGENCY_NUMBER_LIST_IE_MIN_LENGTH - 1), len);
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  e->lengthofemergencynumberinformation = *(buffer + decoded);
  decoded++;
  emergencynumberlist->emergencyservicecategoryvalue =
      *(buffer + decoded) & 0x1f;
  decoded++;
  for (int i = 0; i < e->lengthofemergencynumberinformation - 1; i++) {
    e->number_digit[i] = *(buffer + decoded);
    decoded++;
  }
  for (int i = e->lengthofemergencynumberinformation - 1;
       i < EMERGENCY_NUMBER_MAX_DIGITS; i++) {
    e->number_digit[i] = 0xFF;
  }
  Fatal("TODO emergency_number_list_t->next");

  return decoded;
}

//------------------------------------------------------------------------------
int encode_emergency_number_list_ie(
    emergency_number_list_t* emergencynumberlist, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded           = 0;
  emergency_number_list_t* e = emergencynumberlist;

  Fatal("TODO Implement encode_emergency_number_list_ie");
  if (iei_present) {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, EMERGENCY_NUMBER_LIST_IE_MIN_LENGTH, len);
    *buffer = MM_EMERGENCY_NUMBER_LIST_IEI;
    encoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
        buffer, (EMERGENCY_NUMBER_LIST_IE_MIN_LENGTH - 1), len);
  }

  lenPtr = (buffer + encoded);
  encoded++;
  while (e) {
    *(buffer + encoded) =
        emergencynumberlist->lengthofemergencynumberinformation;
    encoded++;
    *(buffer + encoded) =
        0x00 | (emergencynumberlist->emergencyservicecategoryvalue & 0x1f);
    encoded++;
    for (int i = 0; i < EMERGENCY_NUMBER_MAX_DIGITS; i++) {
      if (e->number_digit[i] < 10) {
        *(buffer + encoded) = e->number_digit[i];
      } else {
        break;
      }
      if (e->number_digit[i] < 10) {
        *(buffer + encoded) |= (e->number_digit[i] << 4);
      } else {
        *(buffer + encoded) |= 0xF0;
        encoded++;
        break;
      }
      encoded++;
    }
    e = e->next;
  }
  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}
