/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under 
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.  
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "DetachAccept.h"

#ifndef ENABLE_CODE_COVERAGE
int decode_detach_accept(
  detach_accept_msg *detach_accept,
  uint8_t *buffer,
  uint32_t len)
{
  uint32_t decoded = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
    buffer, DETACH_ACCEPT_MINIMUM_LENGTH, len);
  /*
   * Decoding mandatory fields
   */
  return decoded;
}
#endif /* ENABLE_CODE_COVERAGE */

int encode_detach_accept(
  detach_accept_msg *detach_accept,
  uint8_t *buffer,
  uint32_t len)
{
  int encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
    buffer, DETACH_ACCEPT_MINIMUM_LENGTH, len);
  return encoded;
}
