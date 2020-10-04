/****************************************************************************
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
 ****************************************************************************/
/*****************************************************************************

  Source      amf_app_msg.cpp

  Version     0.1

  Date        2020/09/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include "common_types.h"
#include "intertask_interface.h"
#include "amf_app_msg.h"

extern task_zmq_ctx_t mme_app_task_zmq_ctx;
using namespace std;
namespace magmam5g
{
    
        /****************************************************************************
         **                                                                        **
        ** name:    amf_app_ue_context_release()                             **
        **                                                                        **
        ** description: Send itti mesage to ngap task to send UE Context Release  **
        **              Request                                                   **
        **                                                                        **
        ** inputs:  ue_context_p: Pointer to UE context amf_casue: failed cause   **
        **                                                                        **
        ***************************************************************************/
        void amf_app_ue_context_release(class ue_m5gmm_context_s* ue_context_p, enum ngcause cause) 
        {
            MessageDef* message_p;

            OAILOG_FUNC_IN(LOG_AMF_APP);
            message_p = itti_alloc_new_message(TASK_AMG_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);//ngap_message_def
            if (message_p == NULL) {
                OAILOG_ERROR_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
                    "Failed to allocate memory for NGAP_UE_CONTEXT_RELEASE_COMMAND \n");
                OAILOG_FUNC_OUT(LOG_AMF_APP);
            }

            OAILOG_INFO_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
                "Sending UE Context Release Cmd to ngap for (ue_id = %u)\n"
                "UE Context Release Cause = (%d)\n",ue_context_p->amf_ue_ngap_id, cause);

            NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id = ue_context_p->amf_ue_ngap_id;
            NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id = ue_context_p->gnb_ue_ngap_id;
            NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause = cause;

            message_p->ittiMsgHeader.imsi = ue_context_p->amf_context._imsi64;
            send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
            OAILOG_FUNC_OUT(LOG_AMF_APP);
        }

    
} // namespace magmam5g


