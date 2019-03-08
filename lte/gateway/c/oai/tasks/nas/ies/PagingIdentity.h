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

#ifndef PAGING_IDENTITY_SEEN
#define PAGING_IDENTITY_SEEN

#include <stdint.h>

#define PAGING_IDENTITY_MINIMUM_LENGTH 2
#define PAGING_IDENTITY_MAXIMUM_LENGTH 2

typedef uint8_t paging_identity_t;

int encode_paging_identity(
  paging_identity_t *pagingidentity,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

int decode_paging_identity(
  paging_identity_t *pagingidentity,
  uint8_t iei,
  uint8_t *buffer,
  uint32_t len);

#endif /* PAGING IDENTITY_SEEN */
