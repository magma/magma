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

  Source      amf_cn.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_app_ue_context_and_proc.h"
#include "amf_asDefs.h"
#include "amf_sap.h"
#include "amf_identity.h"

using namespace std;

namespace magma5g {

amf_identity_msg identity_msg;
int amf_cn_send(const amf_cn_t* msg) {
  int rc                       = RETURNerror;
  amf_cn_primitive_t primitive = msg->primitive;

  OAILOG_FUNC_IN(LOG_NAS_AMF);

  switch (primitive) {
    case AMFCN_AUTHENTICATION_PARAM_RES:
      // rc = amf_cn_authentication_res(msg->u.auth_res);
      break;
#if 0 
    case AMFCN_AUTHENTICATION_PARAM_FAIL:
      rc = amf_cn_authentication_fail(msg->u.auth_fail);
      break;

    case AMFCN_NW_INITIATED_DEREGISTRATION_UE:
      rc = amf_cn_nw_initiated_deregistration_ue(
          msg->u.amf_cn_nw_initiated_deregistration.ue_id,
          msg->u.amf_cn_nw_initiated_deregistration.deregistration_type);
      break;

    case AMFCN_DEACTIVATE_PDUSESSION_REQ:
      rc = amf_cn_deactivate_pdusession_req();

      break;
#endif
    case AMFCN_IDENTITY_PARAM_RES:
      // rc = identity_msg.amf_cn_identity_success();
      // Just return as nothing to do. Invocation of Authentication
      // taken care in identity complete routine.
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "AMFCN-SAP - successfully completed identification "
          "response retuning to invoke authentication request\n");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
      break;

    default:
      /*
       * Other primitives are forwarded to the Access Stratum
       */
      rc = RETURNerror;
      break;
  }

  if (rc != RETURNok) {
    // OAILOG_ERROR(
    //   LOG_NAS_AMF, "AMFCN-SAP - Failed to process primitive %s (%d)\n",
    //  amf_cn_primitive_str[primitive - _AMFCN_START - 1], primitive);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

}  // namespace magma5g
