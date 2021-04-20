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

#include "bytes_to_ie.h"

void bytes_to_lai(const char* bytes, lai_t* lai) {
  /*plmn_t plmn = {
    .mcc_digit2 = bytes[0] >> 4,
    .mcc_digit1 = bytes[0] & 0x0F,
    .mnc_digit3 = bytes[1] >> 4,
    .mcc_digit3 = bytes[1] & 0x0F,
    .mnc_digit2 = bytes[2] >> 4,
    .mnc_digit1 = bytes[2] & 0x0F
  };*/

  lai->mccdigit2 = bytes[0] >> 4, lai->mccdigit1 = bytes[0] & 0x0F,
  lai->mncdigit3 = bytes[1] >> 4, lai->mccdigit3 = bytes[1] & 0x0F,
  lai->mncdigit2 = bytes[2] >> 4, lai->mncdigit1 = bytes[2] & 0x0F;

  unsigned char lac_1 = bytes[3];
  unsigned char lac_2 = bytes[4];
  uint16_t lac        = (lac_1 << 8) | lac_2;

  //  lai->plmn = plmn;
  lai->lac = lac;

  return;
}

void bytes_to_tmsi(const char* bytes, tmsi_t* tmsi) {
  unsigned char tmsi_1 = bytes[0];
  unsigned char tmsi_2 = bytes[1];
  unsigned char tmsi_3 = bytes[2];
  unsigned char tmsi_4 = bytes[3];

  *tmsi = (((uint32_t) tmsi_1) << 24) | (((uint32_t) tmsi_2) << 16) |
          (((uint32_t) tmsi_3) << 8) | ((uint32_t) tmsi_4);

  return;
}

void bytes_to_mobile_identity(
    const char* bytes, uint8_t mobile_identity_len, bool is_imsi,
    MobileIdentity_t* mobile_identity) {
  uint8_t typeofidentity;

  if (is_imsi) {
    typeofidentity = MOBILE_IDENTITY_IMSI;
    for (uint8_t i = 0; i < mobile_identity_len; ++i) {
      mobile_identity->u.imsi[i] = bytes[i];
    }
  } else {
    typeofidentity = MOBILE_IDENTITY_TMSI;
    for (uint8_t i = 0; i < mobile_identity_len; ++i) {
      mobile_identity->u.tmsi[i] = bytes[i];
    }
  }

  mobile_identity->length         = mobile_identity_len;
  mobile_identity->typeofidentity = typeofidentity;

  return;
}
