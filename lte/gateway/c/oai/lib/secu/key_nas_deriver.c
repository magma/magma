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

#include <stdint.h>
#include <string.h>

#include "security_types.h"
#include "secu_defs.h"

/*!
   @brief Derive the kNASenc from kasme and perform truncate on the generated
   key to reduce his size to 128 bits. Definition of the derivation function can
   be found in 3GPP TS.33401 #A.7
   @param[in] nas_alg_type NAS algorithm distinguisher
   @param[in] nas_enc_alg_id NAS encryption/integrity algorithm identifier.
   Possible values are:
        - 0 for EIA0 algorithm (Null Integrity Protection algorithm)
        - 1 for 128-EIA1 SNOW 3G
        - 2 for 128-EIA2 AES
   @param[in] kasme Key for MME as provided by AUC
   @param[out] knas Pointer to reference where output of KDF will be stored.
   NOTE: knas is dynamically allocated by the KDF function
*/
int derive_key_nas(
    algorithm_type_dist_t nas_alg_type, uint8_t nas_enc_alg_id,
    const uint8_t* kasme_32, uint8_t* knas) {
  uint8_t s[7]    = {0};
  uint8_t out[32] = {0};

  /*
   * FC
   */
  s[0] = FC_ALG_KEY_DER;
  /*
   * P0 = algorithm type distinguisher
   */
  s[1] = (uint8_t)(nas_alg_type & 0xFF);
  /*
   * L0 = length(P0) = 1
   */
  s[2] = 0x00;
  s[3] = 0x01;
  /*
   * P1
   */
  s[4] = nas_enc_alg_id;
  /*
   * L1 = length(P1) = 1
   */
  s[5] = 0x00;
  s[6] = 0x01;
  // OAILOG_TRACE (LOG_NAS, "FC %d nas_alg_type distinguisher %d
  // nas_enc_alg_identity %d\n", FC_ALG_KEY_DER, nas_alg_type, nas_enc_alg_id);
  // OAILOG_STREAM_HEX(OAILOG_LEVEL_TRACE, LOG_NAS, "s:", s, 7);
  // OAILOG_STREAM_HEX(OAILOG_LEVEL_TRACE, LOG_NAS, "kasme_32:", kasme_32, 32);
  kdf(kasme_32, 32, &s[0], 7, &out[0], 32);
  memcpy(knas, &out[31 - 16 + 1], 16);
  return 0;
}
