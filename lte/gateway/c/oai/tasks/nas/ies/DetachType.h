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

#ifndef DETACH_TYPE_SEEN
#define DETACH_TYPE_SEEN

#include <stdbool.h>
#include <stdint.h>

#define DETACH_TYPE_MINIMUM_LENGTH 1
#define DETACH_TYPE_MAXIMUM_LENGTH 1

typedef struct detach_type_s {
#define DETACH_TYPE_NORMAL_DETACH 0
#define DETACH_TYPE_SWITCH_OFF 1
  bool switchoff;
#define DETACH_TYPE_EPS 0b001
#define DETACH_TYPE_IMSI 0b010
#define DETACH_TYPE_EPS_IMSI 0b011
#define DETACH_TYPE_RESERVED_1 0b110
#define DETACH_TYPE_RESERVED_2 0b111
  uint8_t typeofdetach;
} detach_type_t;

int encode_detach_type(
    detach_type_t* detachtype, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_detach_type(detach_type_t* detachtype);

int decode_detach_type(
    detach_type_t* detachtype, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_detach_type(
    detach_type_t* detachtype, uint8_t iei, uint8_t value, uint32_t len);

#endif /* DETACH_TYPE_SEEN */
