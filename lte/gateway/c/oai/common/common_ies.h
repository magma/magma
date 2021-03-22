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

#ifndef FILE_COMMON_IES_TYPES_SEEN
#define FILE_COMMON_IES_TYPES_SEEN

#include <stdint.h>

#include "3gpp_23.003.h"

#define TMSI_SIZE 4
#define MAX_IMEISV_SIZE 16
#define MAX_IMEI_SIZE 15
#define MAX_MME_NAME_LENGTH 255
#define MAX_VLR_NAME_LENGTH 255
#define SGS_ASSOC_ACTIVE 1
#define SGS_ASSOC_INACTIVE 0

typedef uint8_t TimeZone;
typedef uint8_t UeEmmMode;

typedef struct MobileIdentity_s {
#define MOBILE_IDENTITY_IMSI 0b001
#define MOBILE_IDENTITY_TMSI 0b100
  uint8_t length;          // Length of IMSI or TMSI
  uint8_t typeofidentity;  // IMSI or TMSI
  union {
    char imsi[IMSI_BCD_DIGITS_MAX + 1];
    char tmsi[TMSI_SIZE + 1];
  } u;
} MobileIdentity_t;

typedef struct MobileStationClassmark2_s {
  uint8_t revisionlevel : 2;
  uint8_t esind : 1;
  uint8_t a51 : 1;
  uint8_t rfpowercapability : 3;
  uint8_t pscapability : 1;
  uint8_t ssscreenindicator : 2;
  uint8_t smcapability : 1;
  uint8_t vbs : 1;
  uint8_t vgcs : 1;
  uint8_t fc : 1;
  uint8_t cm3 : 1;
  uint8_t lcsvacap : 1;
  uint8_t ucs2 : 1;
  uint8_t solsa : 1;
  uint8_t cmsp : 1;
  uint8_t a53 : 1;
  uint8_t a52 : 1;
} MobileStationClassmark2_t;

typedef enum additional_updt_s {
  MME_APP_NO_ADDITIONAL_INFO = 0,
  MME_APP_SMS_ONLY           = 1
} additional_updt_t;

typedef enum additional_updt_result_s {
  ADDITONAL_UPDT_RES_NO_ADDITIONAL_INFO = 0,
  ADDITONAL_UPDT_RESCSFB_NOT_PREFERRED  = 1,
  ADDITONAL_UPDT_RES_SMS_ONLY           = 2,
} additional_updt_result_t;

typedef enum ongoing_procedure_s {
  COMBINED_ATTACH,
  COMBINED_TAU
} ongoing_procedure_t;
#endif /* FILE_COMMON_IES_TYPES_SEEN */
