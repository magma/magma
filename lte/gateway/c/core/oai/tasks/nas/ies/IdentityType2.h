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

#ifndef IDENTITY_TYPE_2_H_
#define IDENTITY_TYPE_2_H_
#include <stdint.h>

#define IDENTITY_TYPE_2_MINIMUM_LENGTH 1
#define IDENTITY_TYPE_2_MAXIMUM_LENGTH 1

#define IDENTITY_TYPE_2_IMSI 0b001
#define IDENTITY_TYPE_2_IMEI 0b010
#define IDENTITY_TYPE_2_IMEISV 0b011
#define IDENTITY_TYPE_2_TMSI 0b100
typedef uint8_t IdentityType2;

int encode_identity_type_2(
    IdentityType2* identitytype2, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_identity_type_2_xml(IdentityType2* identitytype2, uint8_t iei);

uint8_t encode_u8_identity_type_2(IdentityType2* identitytype2);

int decode_identity_type_2(
    IdentityType2* identitytype2, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_identity_type_2(
    IdentityType2* identitytype2, uint8_t iei, uint8_t value, uint32_t len);

#endif /* IDENTITY TYPE 2_H_ */
