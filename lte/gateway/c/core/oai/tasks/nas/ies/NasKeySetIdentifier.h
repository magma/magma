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

#ifndef NAS_KEY_SET_IDENTIFIER_SEEN
#define NAS_KEY_SET_IDENTIFIER_SEEN

#include <stdint.h>

#define NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH 1
#define NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH 1

typedef struct NasKeySetIdentifier_tag {
#define NAS_KEY_SET_IDENTIFIER_NATIVE 0
#define NAS_KEY_SET_IDENTIFIER_MAPPED 1
  uint8_t tsc : 1;
#define NAS_KEY_SET_IDENTIFIER_NOT_AVAILABLE 0b111
  uint8_t naskeysetidentifier : 3;
} NasKeySetIdentifier;

int encode_nas_key_set_identifier(
    NasKeySetIdentifier* naskeysetidentifier, uint8_t iei, uint8_t* buffer,
    uint32_t len);

uint8_t encode_u8_nas_key_set_identifier(
    NasKeySetIdentifier* naskeysetidentifier);

int decode_nas_key_set_identifier(
    NasKeySetIdentifier* naskeysetidentifier, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_u8_nas_key_set_identifier(
    NasKeySetIdentifier* naskeysetidentifier, uint8_t iei, uint8_t value,
    uint32_t len);

#endif /* NAS KEY SET IDENTIFIER_SEEN */
