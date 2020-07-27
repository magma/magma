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

#ifndef GUTI_TYPE_SEEN
#define GUTI_TYPE_SEEN

#include <stdbool.h>
#include <stdint.h>

#define GUTI_TYPE_MINIMUM_LENGTH 1
#define GUTI_TYPE_MAXIMUM_LENGTH 1

#define GUTI_NATIVE 0
#define GUTI_MAPPED 1

typedef bool guti_type_t;

int encode_guti_type(
    guti_type_t* gutitype, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_guti_type(guti_type_t* gutitype);

int decode_guti_type(
    guti_type_t* gutitype, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_guti_type(
    guti_type_t* gutitype, uint8_t iei, uint8_t value, uint32_t len);

#endif /* GUTI_TYPE_SEEN */
