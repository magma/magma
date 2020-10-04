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

  Source      amf_ap_transport.cpp

  Version     0.1

  Date        2020/07/28

  Product     AMF stack

  Subsystem   from AMF to NGAP

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "amf_as.h"
#include "nas_network.h"
#include "amf_recv.h"
#include "amf_message.h"
#include "amf_common_defs.h"
#include "../3gpp/3gpp_24.007.h"
#include "../3gpp/3gpp_38.401.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "tasks_def.h"
using namespace std;
#pragma once


namespace magma5g
{
    int amf_app_defs::amf_app_handle_nas_dl_req(const amf_ue_ngap_id_t ue_id, bstring nas_msg, nas_error_code_t transaction_status)
    {
        OAILOG_FUNC_IN(LOG_AMF_APP);
        MessageDef* message_p = NULL;
        int rc = RETURNok;
        gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;

        message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

        amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
        if (!amf_app_desc_p) 
        {
            OAILOG_CRITICAL(
            LOG_AMF_APP,
            "DOWNLINK NAS TRANSPORT. Failed to get global amf_app_desc context \n");
            OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }
        ue_m5gmm_context_s *ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
        if (ue_context) 
        {
            gnb_ue_ngap_id = ue_context->gnb_ue_ngap_id;
        } 
        else 
        {
            OAILOG_WARNING(
            LOG_AMF_APP,
            " DOWNLINK NAS TRANSPORT. Null UE Context for "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "\n",
            ue_id);
            OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }

        NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
        NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = ue_id;
        NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(nas_msg);
        bdestroy_wrapper(&nas_msg);

        message_p->ittiMsgHeader.imsi = ue_context->amf_context._imsi64;
        /*
        * Store the NGAP NAS DL DATA REQ in case of IMSI or combined 5GMM/IMSI deregister in sgs context
        * and send it after recieving the 5GS IMSI deregister Ack from 5GS task.
        */
        if (ue_context->sgs_context != NULL) 
        {
            /* Send the NGAP NAS DL DATA REQ to NGAP */
            rc = send_msg_to_task(&amf_app_task_zmq_ctx,TASK_NGAP, message_p);
            
        } 
        else 
        {
            rc = send_msg_to_task(&amf_app_task_zmq_ctx,TASK_NGAP, message_p);
            
        }

        /*
        * Move the UE to ECM Connected State ,if not in connected state already
        * N2 Signaling connection gets established via first DL NAS Trasnport
        * message in some scenarios so check the state
        * first
        */
        if (ue_context->ecm_state != ECM_CONNECTED) 
        {
            OAILOG_DEBUG( LOG_AMF_APP, "AMF_APP:DOWNLINK NAS TRANSPORT. Establishing N2 sig connection. "
            "AMF_ue_NGap_id = " AMF_UE_NGAP_ID_FMT "\t ""gnb_ue_ngap_id = " GNB_UE_NGAP_ID_FMT " \n",
            ue_id, gnb_ue_ngap_id);
            //TODO
            //amf_ue_context_update_ue_sig_connection_state(&amf_app_desc_p->amf_ue_contexts, ue_context, ECM_CONNECTED);
        }

        // Check the transaction status. And trigger the UE context release command accrordingly.
        if (transaction_status != AS_SUCCESS) 
        {
            ue_context->ue_context_rel_cause = NGAP_NAS_NORMAL_RELEASE;
            // Notify NGAP to send UE Context Release Command to gNB.
            amf_app_itti_ue_context_release(ue_context, ue_context->ue_context_rel_cause);
        }
        OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
    }


}