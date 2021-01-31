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

  Source      nas_proc.cpp

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

#include <sstream>

#include "amf_fsm.h"
#include "amf_asDefs.h"
//#include "amf_nas5g_proc.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_as.h"
#include "amf_sap.h"
//#include "nas_procedures.h"
#include "nas5g_network.h"
//#include "amf_data.h"
//#include "log.h"
using namespace std;

extern amf_config_t amf_config;
namespace magma5g {
amf_sap_c amf_sap_op;
nas_network nas_networks_proc;
AmfMsg amf_msg_obj;
int nas_proc::nas_proc_establish_ind(
    const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi,
    const m5g_rrc_establishment_cause_t as_cause, const s_tmsi_m5_t s_tmsi,
    bstring msg) {
  amf_sap_t amf_sap;
  uint32_t rc = RETURNerror;
  if (msg) {
    /*
     * Notify the AMF procedure call manager that NAS signalling
     * connection establishment indication message has been received
     * from the Access-Stratum sublayer
     */
    amf_sap.primitive = AMFAS_ESTABLISH_REQ;
    amf_sap.u.amf_as.primitive =
        _AMFAS_ESTABLISH_REQ;  // TODO verify with Sanjay, Sandeep
    amf_sap.u.amf_as.u.establish.ue_id          = ue_id;
    amf_sap.u.amf_as.u.establish.is_initial     = true;
    amf_sap.u.amf_as.u.establish.is_amf_ctx_new = is_mm_ctx_new;

    amf_sap.u.amf_as.u.establish.nas_msg = msg;
    // TODO -  NEED-RECHECK
    // amf_sap.u.amf_as.u.establish.tai = &originating_tai;
    // amf_sap.u.amf_as.u.establish.plmn_id            = &originating_tai.plmn;
    // amf_sap.u.amf_as.u.establish.tac                = originating_tai.tac;
    amf_sap.u.amf_as.u.establish.ecgi = ecgi;

    rc = amf_sap_op.amf_sap_send(&amf_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

amf_procedures_t* _nas_new_amf_procedures(amf_context_t* const amf_context) {
  amf_procedures_t* amf_procedures = new amf_procedures_t;
  LIST_INIT(&amf_procedures->amf_common_procs);
  return amf_procedures;
}

//-----------------------------------------------------------------------------
nas_5g_auth_info_proc_t* nas_new_5gcn_auth_info_procedure(
    amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = _nas_new_amf_procedures(amf_context);
  }

  nas_5g_auth_info_proc_t* auth_info_proc = new nas_5g_auth_info_proc_t;
  auth_info_proc->cn_proc.base_proc.nas_puid =
      __sync_fetch_and_add(&nas_puid, 1);
  // auth_info_proc->cn_proc.base_proc.type = NAS_PROC_TYPE_CN;//TODO -
  // NEED-RECHECK
  auth_info_proc->cn_proc.type = CN5G_PROC_AUTH_INFO;

  nas5g_cn_procedure_t* wrapper = new nas5g_cn_procedure_t;
  if (wrapper) {
    wrapper->proc = &auth_info_proc->cn_proc;
    // LIST_INSERT_HEAD(&amf_context->amf_procedures->cn_procs, wrapper,
    // entries);//TODO -  NEED-RECHECK commented as list in class. recheck
    // OAILOG_TRACE(LOG_NAS_AMF, "New CN5G_PROC_AUTH_INFO\n");
    return auth_info_proc;
  } else {
    nas_networks_proc.free_wrapper((void**) &auth_info_proc);
  }
  return NULL;
}
//---------------------------------------------------------------------------------------------
nas_amf_registration_proc_t* nas_proc::get_nas_specific_procedure_registration(
    const amf_context_t* ctxt) {
  OAILOG_INFO(
      LOG_AMF_APP, "AMF-TEST: in get_nas_specific_procedure_registration\n");
  OAILOG_INFO(LOG_AMF_APP, "AMF-TEST:ctxt:%p\n", ctxt);
  OAILOG_INFO(
      LOG_AMF_APP, "AMF-TEST:amf_procedures:%p\n", ctxt->amf_procedures);
  OAILOG_INFO(
      LOG_AMF_APP, "AMF-TEST:amf_specific_proc:%p\n",
      ctxt->amf_procedures->amf_specific_proc);
  OAILOG_INFO(
      LOG_AMF_APP, "AMF-TEST:amf_specific_proc->type:%d\n",
      ctxt->amf_procedures->amf_specific_proc->type);

  if ((ctxt) && (ctxt->amf_procedures) &&
      (ctxt->amf_procedures->amf_specific_proc) &&
      ((AMF_SPEC_PROC_TYPE_REGISTRATION ==
        ctxt->amf_procedures->amf_specific_proc->type))) {
    return (nas_amf_registration_proc_t*)
        ctxt->amf_procedures->amf_specific_proc;
  }

  return NULL;
}
//------------------------------------------------------------------------------
bool nas_proc::is_nas_specific_procedure_registration_running(
    const amf_context_t* ctxt) {
  if ((ctxt) && (ctxt->amf_procedures) &&
      (ctxt->amf_procedures->amf_specific_proc) &&
      ((AMF_SPEC_PROC_TYPE_REGISTRATION ==
        ctxt->amf_procedures->amf_specific_proc->type)))
    return true;
  return false;
}
//-----------------------------------------------------------------------------
#if 0
amf_procedures_t* _nas_new_amf_procedures(amf_context_t* amf_context) {
  amf_procedures_t* amf_procedures;
  // amf_procedures = new (amf_context->amf_procedures);
  amf_procedures = new amf_procedures_t;
  LIST_INIT(&amf_procedures->amf_common_procs);
  return amf_procedures;
}
#endif
//-----------------------------------------------------------------------------
int nas_proc::nas5g_message_decode(
    unsigned char* buffer, amf_nas_message_t* nas_msg, int length,
    amf_security_context_t* amf_security_context,
    amf_nas_message_decode_status_t* decode_status) {
  OAILOG_FUNC_IN(LOG_NAS5G);
  amf_security_context_t* amf_security = amf_security_context;
  int bytes                            = 0;
  uint32_t mac                         = 0;
  uint16_t short_mac                   = 0;
  int size                             = 0;
  bool is_sr                           = false;
  uint8_t sequence_number              = 0;
  uint8_t temp_sequence_number         = 0;
  AmfMsgHeader* msg_header             = nullptr;
  AmfMsg* msg_amf                      = nullptr;
  /*
   * Decode the header
   */
  // OAILOG_STREAM_HEX( OAILOG_LEVEL_DEBUG, LOG_NAS5G, "Incoming NAS message: ",
  // buffer, length); size = amf_msg_obj.AmfMsgDecodeHeaderMsg(&nas_msg->header,
  // &buffer, length);
  msg_header = (AmfMsgHeader*) &nas_msg->header;
  size       = amf_msg_obj.AmfMsgDecodeHeaderMsg(msg_header, buffer, length);
  // OAILOG_DEBUG(LOG_NAS5G, "nas_message_header_decode returned size %d\n",
  // size);
  if (size < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS5G, TLV_BUFFER_TOO_SHORT);
  }
  // TODO:  else part for nas5g_message_protected_decode
  /*
   * Decode plain NAS message
   */
  //  msg_amf = (AmfMsg*) &nas_msg;
  msg_amf = (AmfMsg*) &nas_msg->plain.amf;
  bytes   = amf_msg_obj.M5gNasMessageDecodeMsg(msg_amf, buffer, length);
  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}
//------------------------------------------------------------------------------------

nas_amf_ident_proc_t* nas5g_new_identification_procedure(
    amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: From nas5g_new_identification_procedure allocating for "
        "amf_procedures\n");
    amf_context->amf_procedures = _nas_new_amf_procedures(amf_context);
  }
  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF_TEST: From nas5g_new_identification_procedure amf_procedures:%p\n",
      amf_context->amf_procedures);
  nas_amf_ident_proc_t* ident_proc = new nas_amf_ident_proc_t;

  ident_proc->amf_com_proc.amf_proc.base_proc.nas_puid =
      __sync_fetch_and_add(&nas_puid, 1);
  ident_proc->amf_com_proc.amf_proc.type = NAS_AMF_PROC_TYPE_COMMON;
  // ident_proc->amf_com_proc.amf_proc.base_proc.type = NAS_PROC_TYPE_AMF;//TODO
  // -  NEED-RECHECK ident_proc->amf_com_proc.type                    =
  // AMF_COMM_PROC_IDENT;

  ident_proc->T3570.sec = amf_config.nas_config.t3570_sec;
  ident_proc->T3570.id  = NAS5G_TIMER_INACTIVE_ID;

  nas_amf_common_procedure_t* wrapper = new nas_amf_common_procedure_t;
  if (wrapper) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: From nas5g_new_identification_procedure amf_procedures:%p\n",
        amf_context->amf_procedures);
    wrapper->proc = &ident_proc->amf_com_proc;
    LIST_INSERT_HEAD(
        &amf_context->amf_procedures->amf_common_procs, wrapper, entries);
    // OAILOG_TRACE(LOG_NAS_AMF, "New AMF_COMM_PROC_IDENT\n");
    return ident_proc;
  } else {
    nas_networks_proc.free_wrapper((void**) &ident_proc);
  }
  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF_TEST: From nas5g_new_identification_procedure amf_procedures:%p\n",
      amf_context->amf_procedures);
  return ident_proc;
}
/*
   --------------------------------------------------------------------------
                AMF NAS specific local functions
   --------------------------------------------------------------------------
*/

/*
 * Description: Sends IDENTITY REQUEST message and start timer T3470.
 *
 * Inputs:  args:      handler parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    None
 *      Others:    T3470
 */
static int amf_identification_request(nas_amf_ident_proc_t* const proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  amf_sap_t amf_sap;  //             = {0};
  int rc                 = RETURNok;
  amf_context_t* amf_ctx = NULL;
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: Sending AS IDENTITY_REQUEST\n");
  // ue_m5gmm_context_s* ue_5gmm_context =
  // amf_ue_context_exists_amf_ue_ngap_id(proc->ue_id);
  ue_m5gmm_context_s* ue_5gmm_context;  // =
  if (ue_5gmm_context) {
    amf_ctx = &ue_5gmm_context->amf_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * Notify AMF-AS SAP that Identity Request message has to be sent
   * to the UE
   */
  amf_sap.primitive = AMFAS_SECURITY_REQ;
  // amf_sap.primitive = AMFAS_DATA_REQ;
  amf_sap.u.amf_as.u.security.puid =
      proc->amf_com_proc.amf_proc.base_proc.nas_puid;
  // amf_sap.u.amf_as.u.security.guti        = NULL;
  amf_sap.u.amf_as.u.security.ue_id                    = proc->ue_id;
  amf_sap.u.amf_as.u.security.msg_type                 = AMF_AS_MSG_TYPE_IDENT;
  amf_sap.u.amf_as.u.security.ident_type               = proc->identity_type;
  amf_sap.u.amf_as.u.security.sctx.is_knas_int_present = true;
  amf_sap.u.amf_as.u.security.sctx.is_knas_enc_present = true;
  amf_sap.u.amf_as.u.security.sctx.is_new =
      true;  // TODO AMF-TEST, handle the bool values in
             // amf_as_set_security_data
  /*
   * Setup 5G CN NAS security data
   */
  // TODO
  // amf_as_set_security_data(&amf_sap.u.amf_as.u.security.sctx,
  // &amf_ctx->_security, false, true);
  rc = amf_sap_op.amf_sap_send(&amf_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3470 timer
     */
    // TODO
    // nas_start_T3470(proc->ue_id,
    // &proc->T3470,proc->amf_com_proc.amf_proc.base_proc.time_out, (void*)
    // amf_ctx);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
//-------------------------------------------------------------------------------------
int identification::amf_proc_identification(
    amf_context_t* const amf_context, nas_amf_proc_t* const amf_proc,
    const identity_type2_t type, success_cb_t success, failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  identification amf_identity;  // to access respective functions

  amf_context->amf_fsm_state =
      AMF_REGISTERED;  // TODO AMF-TEST amf_context->amf_fsm_state is set to 0;
                       // bypass
  if ((amf_context) && ((AMF_DEREGISTERED == amf_context->amf_fsm_state) ||
                        (AMF_REGISTERED == amf_context->amf_fsm_state))) {
    amf_ue_ngap_id_t ue_id =
        PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
            ->amf_ue_ngap_id;

    // OAILOG_INFO(LOG_NAS_AMF,"AMF-PROC  - Initiate identification type = %s
    // (%d), ctx = %p\n",
    //   amf_identity.amf_identity_type_str[type], type, amf_context);

    nas_amf_ident_proc_t* ident_proc =
        nas5g_new_identification_procedure(amf_context);
    if (ident_proc) {
      if (amf_proc) {
        if ((NAS_AMF_PROC_TYPE_SPECIFIC == amf_proc->type) &&
            (AMF_SPEC_PROC_TYPE_REGISTRATION ==
             ((nas_amf_specific_proc_t*) amf_proc)->type)) {
          ident_proc->is_cause_is_registered = true;
        }
      }
      ident_proc->identity_type        = type;
      ident_proc->retransmission_count = 0;
      ident_proc->ue_id                = ue_id;
      //((nas5g_base_proc_t*) ident_proc)->parent = (nas5g_base_proc_t*)
      // amf_proc;
      ident_proc->amf_com_proc.amf_proc.delivered = NULL;
      // TODO - RECHECK later
      // ident_proc->amf_com_proc.amf_proc.previous_amf_fsm_state =
      // amf_fsm_get_state(amf_context);
      // ident_proc->amf_com_proc.amf_proc.not_delivered
      // =identification_Nll_failure;
      // ident_proc->amf_com_proc.amf_proc.not_delivered_ho =
      // _identification_non_delivered_ho;
      ident_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
      ident_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
      // ident_proc->amf_com_proc.amf_proc.base_proc.abort =
      // _identification_abort;
      ident_proc->amf_com_proc.amf_proc.base_proc.fail_in =
          NULL;  // only response
      // ident_proc->amf_com_proc.amf_proc.base_proc.time_out =
      // _identification_t3470_handler;
    }

    rc = amf_identification_request(ident_proc);

    if (rc != RETURNerror) {
      /*
       * Notify 5G CN that common procedure has been initiated
       */
      amf_sap_t amf_sap;  // = {0};

      amf_sap.primitive       = AMFREG_COMMON_PROC_REQ;
      amf_sap.u.amf_reg.ue_id = ue_id;
      amf_sap.u.amf_reg.ctx   = amf_context;
      rc                      = amf_sap_op.amf_sap_send(&amf_sap);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
}  // namespace magma5g
