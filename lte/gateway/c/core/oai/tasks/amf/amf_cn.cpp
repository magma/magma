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

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"

namespace magma5g {

//------------------------------------------------------------------------------
static int amf_cn_authentication_res(amf_cn_auth_res_t* const msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  amf_context_t* amf_ctx = NULL;
  int rc                 = RETURNerror;

  /*
   * We received security vector from HSS. Try to setup security with UE
   */
  ue_m5gmm_context_s* ue_m5gmm_context =
      amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

  if (ue_m5gmm_context) {
    amf_ctx = &ue_m5gmm_context->amf_context;
    nas5g_auth_info_proc_t* auth_info_proc =
        get_nas5g_cn_procedure_auth_info(amf_ctx);

    if (auth_info_proc) {
      for (int i = 0; i < msg->nb_vectors; i++) {
        auth_info_proc->vector[i] = msg->vector[i];
        msg->vector[i]            = NULL;
      }
      auth_info_proc->nb_vectors = msg->nb_vectors;

      rc = amf_authentication_proc_success(amf_ctx);
    } else {
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "EMM-PROC  - "
          "Failed to find Auth_info procedure associated to UE "
          "ID: " AMF_UE_NGAP_ID_FMT,
          msg->ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
int amf_cn_send(const amf_cn_t* msg) {
  int rc                       = RETURNerror;
  amf_cn_primitive_t primitive = msg->primitive;

  OAILOG_FUNC_IN(LOG_NAS_AMF);

  switch (primitive) {
    case _AMFCN_AUTHENTICATION_PARAM_RES:
      rc = amf_cn_authentication_res(msg->u.auth_res);
      break;

    default:
      /*
       * Other primitives are forwarded to the Access Stratum
       */
      rc = RETURNerror;
      break;
  }

  if (rc != RETURNok) {
    OAILOG_ERROR(LOG_NAS_AMF, "AMF-SAP - Failed to process primitive \n");
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
}  // namespace magma5g
