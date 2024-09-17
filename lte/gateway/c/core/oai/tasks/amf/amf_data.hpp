/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once
#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"
#ifdef __cplusplus
};
#endif
#include "lte/gateway/c/core/oai/include/amf_securityDef.h"

typedef uint8_t ksi_t;
#define AMF_CTXT_MEMBER_AUTH_VECTORS ((uint32_t)1 << 7)
#define AMF_CTXT_MEMBER_SECURITY ((uint32_t)1 << 8)
#define AMF_CTXT_MEMBER_GUTI ((uint32_t)1 << 4)
#define AMF_CTXT_MEMBER_AUTH_VECTOR0 ((uint32_t)1 << 26)
#define IS_AMF_CTXT_PRESENT_SECURITY(aMfCtXtPtR) \
  (!!((aMfCtXtPtR)->member_present_mask & AMF_CTXT_MEMBER_SECURITY))
#define IS_AMF_CTXT_VALID_AUTH_VECTORS(aMfCtXtPtR) \
  (!!((aMfCtXtPtR)->member_valid_mask & AMF_CTXT_MEMBER_AUTH_VECTORS))
#define IS_AMF_CTXT_VALID_AUTH_VECTOR(aMfCtXtPtR, KsI) \
  (!!((aMfCtXtPtR)->member_valid_mask &                \
      ((AMF_CTXT_MEMBER_AUTH_VECTOR0) << KsI)))
#define AUTS_LENGTH 14
#define RAND_LENGTH_BITS (128)
#define RAND_LENGTH_OCTETS (RAND_LENGTH_BITS / 8)
#define M5G_IMSI_BCD_DIGITS_MAX 15

// Encryption and Integrity algorithms used
namespace magma5g {
typedef struct selected_algorithms_s {
  uint8_t encryption : 4; /* algorithm used for ciphering           */
  uint8_t integrity : 4;  /* algorithm used for integrity protection */
} selected_algorithms_t;  /* AMF selected algorithms                */

typedef struct count_s {
  uint32_t overflow : 16;
  uint32_t seq_num : 8;
} count_t; /* Downlink and uplink count params */

// AMF security context
typedef struct amf_security_context_s {
  amf_sc_type_t sc_type; /* Type of security context        */
#define EKSI_MAX_VALUE 6
  ksi_t eksi; /* NAS key set identifier for E-UTRAN      */
#define AMF_SECURITY_VECTOR_INDEX_INVALID (-1)
  int vector_index;                     /* Pointer on vector */
  uint8_t knas_enc[AUTH_KNAS_ENC_SIZE]; /* NAS cyphering key               */
  uint8_t knas_int[AUTH_KNAS_INT_SIZE]; /* NAS integrity key               */
  uint8_t kamf[AUTH_KAMF_SIZE];         /* AMF key               */
  uint8_t kgnb[AUTH_KGNB_SIZE];         /* GNB key               */
  count_t dl_count;
  count_t ul_count;
  selected_algorithms_t selected_algorithms;
  count_t kenb_ul_count;
  // Requirement AMF24.501R15_4.4.4.3_2 (DEREGISTRATION REQUEST (if sent before
  // security has been activated);)
  uint8_t direction_encode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
  uint8_t direction_decode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
} amf_security_context_t;

// Authentication request information
typedef struct n6_auth_info_req_s {
  char imsi[M5G_IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;
  // Number of vectors to retrieve from HSS, should be equal to one
  uint8_t nb_of_vectors;
  // Bit to indicate that USIM has requested a re-synchronization of SQN
  unsigned re_synchronization : 1;
  /* AUTS to provide to AUC.
   * Only present and interpreted if re_synchronization == 1.
   */
  uint8_t resync_param[RAND_LENGTH_OCTETS + AUTS_LENGTH];
} n6_auth_info_req_t;

typedef enum {
  AMF_IMEISV_NOT_REQUESTED = 0,
  AMF_IMEISV_REQUESTED = 1
} amf_imeisv_req_type_t;

}  // namespace magma5g
