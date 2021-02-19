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

  Source      amf_as.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_AS_SEEN
#define AMF_AS_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"

#ifdef __cplusplus
};
#endif
#include "amf_app_ue_context_and_proc.h"  // "amf_data.h" included in it
#include "amf_asDefs.h"
#include "amf_as_message.h"
//#include "amf_data.h"
using namespace std;
namespace magma5g {
#if 1
class amf_as : public AMFMsg {
 public:
  void amf_as_initialize(void);

  int amf_as_send(amf_as_t* msg);

  int amf_as_send_ng(const amf_as_t* msg);

  static AMFMsg* amf_as_set_header(
      amf_nas_message_t* msg, const amf_as_security_data_t* security);
};

class amf_as_dl_message {
 public:
  // For _AMFAS_DATA_REQ primitive
  uint16_t amf_as_data_req(
      const amf_as_data_t* msg, m5g_dl_info_transfer_req_t* as_msg);
  // For _AMFAS_ESTABLISH_CNF primitive
  uint16_t amf_as_establish_cnf(
      const amf_as_establish_t* establish,
      nas5g_establish_rsp_t* nas_establish_rsp);
  // For _AMFAS_ESTABLISH_REJ primitive
  uint16_t amf_as_establish_rej(
      const amf_as_establish_t* establish,
      nas5g_establish_rsp_t* nas_establish_rsp);
  // For _AMFAS_SECURITY_REQ primitive
  // TODO-RECHECK commented till get implemented
  // uint16_t amf_as_security_req(amf_as_security_t security,
  // m5g_dl_info_transfer_req_t dl_transfer_req);
};

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
#if 0
  AS_CAUSE_EMERGENCY = EMERGENCY,
  AS_CAUSE_HIGH_PRIO = HIGH_PRIORITY_ACCESS,
  AS_CAUSE_MT_ACCESS = MT_ACCESS,
  AS_CAUSE_MO_SIGNAL = MO_SIGNALLING,
  AS_CAUSE_MO_DATA   = MO_DATA,
  AS_CAUSE_V1020     = DELAY_TOLERANT_ACCESS_V1020
#endif
} as_m5gcause_t;
#endif

}  // namespace magma5g
#endif
