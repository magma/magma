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

#ifndef FILE_SECU_DEFS_SEEN
#define FILE_SECU_DEFS_SEEN

#include <stdint.h>

#include "security_types.h"

#define SECU_DIRECTION_UPLINK 0
#define SECU_DIRECTION_DOWNLINK 1

void kdf(
    const uint8_t* key, const unsigned key_len, uint8_t* s,
    const unsigned s_len, uint8_t* out, const unsigned out_len);

int derive_keNB(
    const uint8_t* kasme_32, const uint32_t nas_count, uint8_t* keNB);

int derive_key_nas(
    algorithm_type_dist_t nas_alg_type, uint8_t nas_enc_alg_id,
    const uint8_t* kasme_32, uint8_t* knas);

int derive_NH(
    const uint8_t* kasme_32, const uint8_t* syncInput, uint8_t* next_hop,
    uint8_t* next_hop_chaining_count);

#define derive_key_nas_enc(aLGiD, kASME, kNAS)                                 \
  derive_key_nas(NAS_ENC_ALG, aLGiD, kASME, kNAS)

#define derive_key_nas_int(aLGiD, kASME, kNAS)                                 \
  derive_key_nas(NAS_INT_ALG, aLGiD, kASME, kNAS)

#define derive_key_rrc_enc(aLGiD, kASME, kNAS)                                 \
  derive_key_nas(RRC_ENC_ALG, aLGiD, kASME, kNAS)

#define derive_key_rrc_int(aLGiD, kASME, kNAS)                                 \
  derive_key_nas(RRC_INT_ALG, aLGiD, kASME, kNAS)

#define derive_key_up_enc(aLGiD, kASME, kNAS)                                  \
  derive_key_nas(UP_ENC_ALG, aLGiD, kASME, kNAS)

#define derive_key_up_int(aLGiD, kASME, kNAS)                                  \
  derive_key_nas(UP_INT_ALG, aLGiD, kASME, kNAS)

#define SECU_DIRECTION_UPLINK 0
#define SECU_DIRECTION_DOWNLINK 1

typedef struct {
  uint8_t* key;
  uint32_t key_length;
  uint32_t count;
  uint8_t bearer;
  uint8_t direction;
  uint8_t* message;
  /* length in bits */
  uint32_t blength;
} nas_stream_cipher_t;

int nas_stream_encrypt_eea1(
    nas_stream_cipher_t* const stream_cipher, uint8_t* const out);

int nas_stream_encrypt_eia1(
    nas_stream_cipher_t* const stream_cipher, uint8_t const out[4]);

int nas_stream_encrypt_eea2(
    nas_stream_cipher_t* const stream_cipher, uint8_t* const out);

int nas_stream_encrypt_eia2(
    nas_stream_cipher_t* const stream_cipher, uint8_t const out[4]);

#endif /* FILE_SECU_DEFS_SEEN */
