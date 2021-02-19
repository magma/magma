
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
/*****************************************************************************

  Source      amf_data.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_DATA_SEEN
#define AMF_DATA_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"
#include "3gpp_24.501.h"
#ifdef __cplusplus
};
#endif
#include "amf_securityDef.h"
using namespace std;
typedef uint8_t ksi_t;
#define AMF_CTXT_MEMBER_AUTH_VECTORS ((uint32_t) 1 << 7)
#define AMF_CTXT_MEMBER_SECURITY ((uint32_t) 1 << 8)
#define AMF_CTXT_MEMBER_GUTI ((uint32_t) 1 << 4)
#define AMF_CTXT_MEMBER_OLD_GUTI ((uint32_t) 1 << 3)
#define AMF_CTXT_MEMBER_AUTH_VECTOR0 ((uint32_t) 1 << 26)

#define IS_AMF_CTXT_PRESENT_SECURITY(aMfCtXtPtR)                               \
  (!!((aMfCtXtPtR)->member_present_mask & AMF_CTXT_MEMBER_SECURITY))

#define IS_AMF_CTXT_VALID_AUTH_VECTORS(aMfCtXtPtR)                             \
  (!!((aMfCtXtPtR)->member_valid_mask & AMF_CTXT_MEMBER_AUTH_VECTORS))

#define IS_AMF_CTXT_VALID_AUTH_VECTOR(aMfCtXtPtR, KsI)                         \
  (!!((aMfCtXtPtR)->member_valid_mask &                                        \
      ((AMF_CTXT_MEMBER_AUTH_VECTOR0) << KsI)))

#define AUTS_LENGTH 14
#define RAND_LENGTH_BITS (128)
#define RAND_LENGTH_OCTETS (RAND_LENGTH_BITS / 8)
#define M5G_IMSI_BCD_DIGITS_MAX 15

namespace magma5g {

/* class count_s
 {
   public:
   uint32_t spare : 8;
   uint32_t overflow : 16;
   uint32_t seq_num : 8;
 }; */ /* Downlink and uplink count params */
class capability {
 public:
  uint8_t m5gs_encryption; /* algorithm used for ciphering            */
  uint8_t m5gs_integrity;  /* algorithm used for integrity protection */
  uint8_t umts_encryption; /* algorithm used for ciphering            */
  uint8_t umts_integrity;  /* algorithm used for integrity protection */
  uint8_t gprs_encryption; /* algorithm used for ciphering            */
  bool umts_present;
  bool gprs_present;
}; /* UE network capability           */

#if 0
    typedef enum amf_sc_type_s {
      SECURITY_CTX_TYPE_NOT_AVAILABLE = 0,
      SECURITY_CTX_TYPE_PARTIAL_NATIVE,
      SECURITY_CTX_TYPE_FULL_NATIVE,
      SECURITY_CTX_TYPE_MAPPED  // UNUSED
    } amf_sc_type_t;
#endif

typedef struct selected_algorithms_s {
 public:
  uint8_t encryption : 4; /* algorithm used for ciphering           */
  uint8_t integrity : 4;  /* algorithm used for integrity protection */
} selected_algorithms_t;  /* AMF selected algorithms                */
typedef struct count_s {
  uint32_t spare : 8;
  uint32_t overflow : 16;
  uint32_t seq_num : 8;
} count_t; /* Downlink and uplink count params */

class amf_security_context_t : public count_s, public capability {
 public:
  amf_sc_type_t sc_type; /* Type of security context        */
  /* state of security context is implicit due to its storage location
   * (current/non-current)*/
#define EKSI_MAX_VALUE 6
  ksi_t eksi; /* NAS key set identifier for E-UTRAN      */
#define AMF_SECURITY_VECTOR_INDEX_INVALID (-1)
  int vector_index;                     /* Pointer on vector */
  uint8_t knas_enc[AUTH_KNAS_ENC_SIZE]; /* NAS cyphering key               */
  uint8_t knas_int[AUTH_KNAS_INT_SIZE]; /* NAS integrity key               */
  count_t dl_count;
  count_t ul_count;
  selected_algorithms_t selected_algorithms;
  count_t kenb_ul_count;
  // Requirement AMF24.501R15_4.4.4.3_2 (DEREGISTRATION REQUEST (if sent before
  // security has been activated);)
  uint8_t activated;
  uint8_t direction_encode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
  uint8_t direction_decode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
  // security keys for HO
  uint8_t next_hop[AUTH_NEXT_HOP_SIZE]; /* Next HOP security parameter */
  uint8_t next_hop_chaining_count;      /* Next Hop Chaining Count */

  // void amf_ctx_clear_security(amf_context_t*  ctxt)
  // __attribute__((nonnull))__attribute__((flatten));
};

// TODO -  NEED-RECHECK
typedef struct n6_auth_info_req_s {
  char imsi[M5G_IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;
  // plmn_t visited_plmn;
  /* Number of vectors to retrieve from HSS, should be equal to one */
  uint8_t nb_of_vectors;

  /* Bit to indicate that USIM has requested a re-synchronization of SQN */
  unsigned re_synchronization : 1;
  /* AUTS to provide to AUC.
   * Only present and interpreted if re_synchronization == 1.
   */
  uint8_t resync_param[RAND_LENGTH_OCTETS + AUTS_LENGTH];
} n6_auth_info_req_t;

typedef enum {
  AMF_IMEISV_NOT_REQUESTED = 0,
  AMF_IMEISV_REQUESTED     = 1
} amf_imeisv_req_type_t;

// void amf_ctx_set_security_eksi(amf_context_t* ctxt, ksi_t eksi);

}  // namespace magma5g
#endif
