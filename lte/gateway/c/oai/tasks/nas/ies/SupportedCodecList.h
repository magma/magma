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

#ifndef SUPPORTED_CODEC_LIST_H_
#define SUPPORTED_CODEC_LIST_H_
#include <stdint.h>

#define SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION                       \
  2 /*Taking consideration of GERAN and UTRAN*/
#define SUPPORTED_CODEC_LIST_MINIMUM_LENGTH 5
#define SUPPORTED_CODEC_LIST_MAXIMUM_LENGTH                                    \
  ((4 * SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION) + 1)

typedef struct SupportedCodecList_tag {
  uint8_t systemidentification;
  uint8_t lengthofbitmap;
  uint16_t codecbitmap;
} SupportedCodecList[SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION];

int encode_supported_codec_list(
  SupportedCodecList *supportedcodeclist,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

int decode_supported_codec_list(
  SupportedCodecList *supportedcodeclist,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

void dump_supported_codec_list_xml(
  SupportedCodecList *supportedcodeclist,
  uint8_t iei);

#endif /* SUPPORTED CODEC LIST_H_ */
