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
#ifdef __cplusplus
};
#endif
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/include/amf_as_message.h"

namespace magma5g {

#define M5GSMobileIdentityMsg_GUTI_LENGTH 11

// Checks PLMNs equality
#define PLMN_ARE_EQUAL(p1, p2) \
  ((MCCS_ARE_EQUAL((p1), (p2))) && (MNCS_ARE_EQUAL((p1), (p2))))

// Checks MCC equality
#define MCCS_ARE_EQUAL(n1, n2)             \
  (((n1).mcc_digit1 == (n2).mcc_digit1) && \
   ((n1).mcc_digit2 == (n2).mcc_digit2) && \
   ((n1).mcc_digit3 == (n2).mcc_digit3))

// Checks Mobile Network Code equality
#define MNCS_ARE_EQUAL(n1, n2)             \
  (((n1).mnc_digit1 == (n2).mnc_digit1) && \
   ((n1).mnc_digit2 == (n2).mnc_digit2) && \
   ((n1).mnc_digit3 == (n2).mnc_digit3))

// AMF_AS Service Access point primitive
status_code_e amf_as_send(amf_as_t* msg);

// Builds NAS message according to the given AMFAS Service Access Point
// primitive
status_code_e amf_as_send_ng(const amf_as_t* msg);

status_code_e initial_context_setup_request(amf_ue_ngap_id_t ue_id,
                                            amf_context_t* amf_ctx,
                                            bstring nas_msg);

// For _AMFAS_DATA_REQ primitive
uint16_t amf_as_data_req(const amf_as_data_t* msg,
                         m5g_dl_info_transfer_req_t* as_msg);

// For _AMFAS_ESTABLISH_CNF primitive
uint16_t amf_as_establish_cnf(const amf_as_establish_t* establish,
                              nas5g_establish_rsp_t* nas_establish_rsp);

enum nas_error_code_t {
  AS_SUCCESS = 1,          /* Success code, transaction is going on    */
  AS_TERMINATED_NAS,       /* Transaction terminated by NAS        */
  AS_TERMINATED_AS,        /* Transaction terminated by AS         */
  AS_NON_DELIVERED_DUE_HO, /* Failure code                 */
  AS_FAILURE               /* Failure code, stand also for lower
                            * layer failure AS_LOWER_LAYER_FAILURE */
};

/* Cause of RRC connection establishment */
typedef enum as_m5gcause_s {
  AS_CAUSE_UNKNOWN = 0,
} as_m5gcause_t;

}  // namespace magma5g
