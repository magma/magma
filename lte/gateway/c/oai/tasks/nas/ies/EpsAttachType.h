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

#ifndef EPS_ATTACH_TYPE_SEEN
#define EPS_ATTACH_TYPE_SEEN
#include <stdint.h>

#define EPS_ATTACH_TYPE_MINIMUM_LENGTH 1
#define EPS_ATTACH_TYPE_MAXIMUM_LENGTH 1

#define EPS_ATTACH_TYPE_EPS 0b0001
#define EPS_ATTACH_TYPE_COMBINED_EPS_IMSI 0b0010
#define EPS_ATTACH_TYPE_EMERGENCY 0b0110
#define EPS_ATTACH_TYPE_RESERVED 0b0111

typedef uint8_t eps_attach_type_t;

int encode_eps_attach_type(
    eps_attach_type_t* epsattachtype, uint8_t iei, uint8_t* buffer,
    uint32_t len);

uint8_t encode_u8_eps_attach_type(eps_attach_type_t* epsattachtype);

int decode_eps_attach_type(
    eps_attach_type_t* epsattachtype, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_u8_eps_attach_type(
    eps_attach_type_t* epsattachtype, uint8_t iei, uint8_t value, uint32_t len);

#endif /* EPS_ATTACH_TYPE_SEEN */
