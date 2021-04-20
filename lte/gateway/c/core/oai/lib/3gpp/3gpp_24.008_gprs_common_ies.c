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

/*! \file 3gpp_24.008_gprs_common_ies.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdint.h>

#include "3gpp_24.008.h"
#include "TLVDecoder.h"
#include "TLVEncoder.h"

//******************************************************************************
// 10.5.7 GPRS Common information elements
//******************************************************************************

//------------------------------------------------------------------------------
// 10.5.7.3 GPRS Timer
//------------------------------------------------------------------------------

static const long gprs_timer_unit[] = {2, 60, 360, 60, 60, 60, 60, 0};

//------------------------------------------------------------------------------
int decode_gprs_timer_ie(
    gprs_timer_t* gprstimer, uint8_t iei, uint8_t* buffer, const uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, GPRS_TIMER_IE_MIN_LENGTH, len);
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  } else {
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, GPRS_TIMER_IE_MIN_LENGTH - 1, len);
  }

  gprstimer->unit       = (*(buffer + decoded) >> 5) & 0x7;
  gprstimer->timervalue = *(buffer + decoded) & 0x1f;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_gprs_timer_ie(
    gprs_timer_t* gprstimer, uint8_t iei, uint8_t* buffer, const uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, GPRS_TIMER_IE_MIN_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) =
      0x00 | ((gprstimer->unit & 0x7) << 5) | (gprstimer->timervalue & 0x1f);
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------

long gprs_timer_value(gprs_timer_t* gprstimer) {
  return (gprstimer->timervalue * gprs_timer_unit[gprstimer->unit]);
}
