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
#ifdef __cplusplus
extern "C" {
#endif
#include "3gpp_24.007.h"
#include "3gpp_38.401.h"
#include "amf_as_message.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_asDefs.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_as.h"
#include "amf_fsm.h"
#include "amf_recv.h"
#include "nas5g_network.h"
//#include "amf_message.h"
#include "amf_app_defs.h"
#include "ngap_messages_types.h"
#include "amf_app_state_manager.h"

using namespace std;

namespace magma5g {
amf_app_defs amf_defs;
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily store
                              // context inserted to ht
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
amf_app_desc_t* get_amf_nas_state(bool read_from_redis);
// int amf_defs.amf_app_handle_nas_dl_req(const amf_ue_ngap_id_t ue_id,
int amf_app_defs::amf_app_handle_nas_dl_req(
    const amf_ue_ngap_id_t ue_id, bstring nas_msg,
    nas5g_error_code_t transaction_status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  MessageDef* message_p           = NULL;
  int rc                          = RETURNok;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  nas_network nas_nw;

  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

  // amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);//TODO -
  // NEED-RECHECK
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (!amf_app_desc_p) {
    OAILOG_CRITICAL(
        LOG_AMF_APP,
        "DOWNLINK NAS TRANSPORT. Failed to get global amf_app_desc context \n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  ue_m5gmm_context_s* ue_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);  // TODO -
                                                    // NEED-RECHECK
  if (ue_context) {
    gnb_ue_ngap_id = ue_context->gnb_ue_ngap_id;
  } else {
    // OAILOG_WARNING(
    //    LOG_AMF_APP,
    //    " DOWNLINK NAS TRANSPORT. Null UE Context for "
    //    "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "\n",
    //    ue_id);
    //    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);AMF_TEST
    ue_context =
        &ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily
                                   // store context inserted to ht
  }

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = ue_id;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg        = bstrcpy(nas_msg);
  nas_nw.bdestroy_wrapper(&nas_msg);

  message_p->ittiMsgHeader.imsi = 13744632839234567870;
  /*
   * Store the NGAP NAS DL DATA REQ in case of IMSI or combined 5GMM/IMSI
   * deregister in sgs context and send it after recieving the 5GS IMSI
   * deregister Ack from 5GS task.
   */
#if 0
  if (ue_context->amf_context.amf_procedures != NULL) {
    /* Send the NGAP NAS DL DATA REQ to NGAP */
    rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  } else {
    rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  }

  /*
   * Move the UE to ECM Connected State ,if not in connected state already
   * N2 Signaling connection gets established via first DL NAS Trasnport
   * message in some scenarios so check the state
   * first
   */
  if (ue_context->ecm_state != M5GCM_CONNECTED) {
    // OAILOG_DEBUG( LOG_AMF_APP, "AMF_APP:DOWNLINK NAS TRANSPORT. Establishing
    // N2 sig connection. " "AMF_ue_NGap_id = " AMF_UE_NGAP_ID_FMT "\t
    //""gnb_ue_ngap_id = " GNB_UE_NGAP_ID_FMT " \n", ue_id, gnb_ue_ngap_id);
    // TODO
    // amf_ue_context_update_ue_sig_connection_state(&amf_app_desc_p->amf_ue_contexts,
    // ue_context, M5GCM_CONNECTED);
  }

#endif  // TODO ue_context null because hash exits fails
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: sending downlink message to NGAP");
  rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  // Check the transaction status. And trigger the UE context release command
  // accrordingly.
  if (transaction_status != M5G_AS_SUCCESS) {
    // TODO -  NEED-RECHECK use file from ngap to find cause
    // ue_context->ue_context_rel_cause = NGAP_NAS_NORMAL_RELEASE;
    // Notify NGAP to send UE Context Release Command to gNB.
    // amf_app_itti_ue_context_release(ue_context,
    // ue_context->ue_context_rel_cause);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}
#if 0
/****************************************************************************
 **                                                                        **
 ** Name:        amf_send_dl_nas_transportmsg()                            **
 **                                                                        **
 ** Description: Builds Downlink Nas Transport message                     **
 **                                                                        **
 **              The Downlink Nas Transport message is sent by the         **
 **              network to the UE to transfer the data in DL              **
 **              This function is used to send DL NAS Transport message    **
 **              via S1AP DL NAS Transport message.                        **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int amf_send_dl_nas_transportmsg(
    const amf_as_data_t* msg, DLNASTransportMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_MAXIMUM_LENGTH;
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type = DOWNLINK_NAS_TRANSPORT;
  /*
   * Mandatory - Nas message container
   */
  size += NAS_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  amf_msg->payload_container = bstrcpy(msg->nas_msg);
  OAILOG_INFO(LOG_NAS_AMF, "AMFAS-SAP - Sending DL NAS - DL NAS5G Transport\n");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}
#endif
}  // namespace magma5g
