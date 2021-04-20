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

/*! \file 3gpp_24.008_cc_ies.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdbool.h>
#include <stdint.h>

#include "3gpp_24.008.h"
#include "TLVDecoder.h"
#include "TLVEncoder.h"

//******************************************************************************
// 10.5.4 Call control information elements.
//******************************************************************************

//------------------------------------------------------------------------------
// 10.5.4.32 Supported codec list
//------------------------------------------------------------------------------
int decode_supported_codec_list(
    supported_codec_list_t* supportedcodeclist, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  int decode_result = 0;
  int decoded       = 0;
  uint8_t ielen     = 0;

  if (iei_present) {
    CHECK_IEI_DECODER(CC_SUPPORTED_CODEC_LIST_IE, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  if ((decode_result = decode_bstring(
           supportedcodeclist, ielen, buffer + decoded, len - decoded)) < 0) {
    return decode_result;
  } else {
    decoded += decode_result;
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_supported_codec_list(
    supported_codec_list_t* supportedcodeclist, const bool iei_present,
    uint8_t* buffer, const uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded       = 0;
  uint32_t encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SUPPORTED_CODEC_LIST_IE_MIN_LENGTH, len);

  if (iei_present) {
    *buffer = CC_SUPPORTED_CODEC_LIST_IE;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  if ((encode_result = encode_bstring(
           *supportedcodeclist, buffer + encoded, len - encoded)) < 0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei_present) ? 1 : 0);
  return encoded;
}
