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

#ifndef SECURITY_HEADER_TYPE_SEEN
#define SECURITY_HEADER_TYPE_SEEN

#define SECURITY_HEADER_TYPE_MINIMUM_LENGTH 1
#define SECURITY_HEADER_TYPE_MAXIMUM_LENGTH 1

typedef uint8_t security_header_type_t;

int encode_security_header_type(
  security_header_type_t *securityheadertype,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

int decode_security_header_type(
  security_header_type_t *securityheadertype,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

#endif /* SECURITY HEADER TYPE_SEEN */
