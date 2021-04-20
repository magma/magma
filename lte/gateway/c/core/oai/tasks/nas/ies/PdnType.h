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

#ifndef PDN_TYPE_SEEN
#define PDN_TYPE_SEEN

#include <stdint.h>

#include "common_types.h"

#define PDN_TYPE_MINIMUM_LENGTH 1
#define PDN_TYPE_MAXIMUM_LENGTH 1

#define PDN_TYPE_IPV4 0b001
#define PDN_TYPE_IPV6 0b010
#define PDN_TYPE_IPV4V6 0b011
#define PDN_TYPE_UNUSED 0b100

// defined in common_types.h
// typedef uint8_t pdn_type_t;

int encode_pdn_type(
    pdn_type_t* pdntype, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_pdn_type(pdn_type_t* pdntype);

int decode_pdn_type(
    pdn_type_t* pdntype, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_pdn_type(
    pdn_type_t* pdntype, uint8_t iei, uint8_t value, uint32_t len);

#endif /* PDN TYPE_SEEN */
