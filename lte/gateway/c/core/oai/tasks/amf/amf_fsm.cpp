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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include <iostream>
#include <cstring>
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/amf_as_message.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_smfDefs.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_app_statistics.hpp"

namespace magma5g {

// UE state matrix holding UE states,State events and PDU session states
ue_state_transition_t ue_state_matrix[UE_STATE_MAX][STATE_EVENT_MAX]
                                     [SESSION_MAX];

/*
 * Create a  list (array) of handlers
 *
 */
static UE_Handlers_t UE_handlers[] = {
    // NAME         // Handler
    {"Common_procedure_Initiated_step1",
     reinterpret_cast<void (*)(void)>(&amf_registration_run_procedure)},
    {"Common_procedure_Initiated_step2",
     reinterpret_cast<void (*)(void)>(&amf_registration_success_security_cb)},
    {"Register_complete",
     reinterpret_cast<void (*)(void)>(&amf_proc_registration_complete)},
    {"Deregister_Initiated",
     reinterpret_cast<void (*)(void)>(&amf_app_handle_deregistration_req)},
    {"Deregister_Completed",
     reinterpret_cast<void (*)(void)>(&amf_app_handle_deregistration_req)},
    {"Idle_mode_procedure",
     reinterpret_cast<void (*)(void)>(&amf_idle_mode_procedure)},
    {"PDU_Creating",
     reinterpret_cast<void (*)(void)>(&amf_smf_initiate_pdu_session_creation)},
    {"PDU_Created",
     reinterpret_cast<void (*)(void)>(&amf_app_handle_pdu_session_accept)},
    {"PDU_Release",
     reinterpret_cast<void (*)(void)>(&release_session_gprc_req)},
    {"PDU_Session_Modification_Request",
     reinterpret_cast<void (*)(void)>(
         &amf_app_pdu_session_modification_request)},
    {"PDU_Session_Modification_Complete",
     reinterpret_cast<void (*)(void)>(
         &amf_app_pdu_session_modification_complete)},
    {"PDU_Session_Modification_Reject",
     reinterpret_cast<void (*)(void)>(
         &amf_app_pdu_session_modification_command_reject)}};

/*
 * Update ue_state_matrix
 *
 * @param current ue state,state event,current PDU session state,next UE
 * state,next PDU session state, Function handler
 * @return null
 */
void Update_ue_state_matrix(m5gmm_state_t cur_state, int event,
                            SMSessionFSMState session_state,
                            m5gmm_state_t next_state,
                            SMSessionFSMState next_sess_state,
                            const char* func) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  uint8_t cnt = 0;
  for (cnt = 0; cnt < sizeof(UE_handlers) / sizeof(UE_handlers[0]); cnt++) {
    if (0 == strcmp(UE_handlers[cnt].name, func)) {
      ue_state_matrix[cur_state][event][session_state].handler =
          UE_handlers[cnt];
      ue_state_matrix[cur_state][event][session_state].next_sess_state =
          next_sess_state;
      ue_state_matrix[cur_state][event][session_state].next_state = next_state;
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/*
 * Create ue_state_matrix
 *
 * @param void
 * @return null
 */
void create_state_matrix() {
  /*Update_state_matrix holding
   * Current UE State
   * STATE Event
   * Current PDU Session State
   * next UE State
   * next PDU Session State
   * Function handler name holding in UE_handlers list
   */
  /* UE state Transitions */
  OAILOG_FUNC_IN(LOG_AMF_APP);
  Update_ue_state_matrix(DEREGISTERED, STATE_EVENT_REG_REQUEST, SESSION_NULL,
                         COMMON_PROCEDURE_INITIATED1, SESSION_NULL,
                         "Common_procedure_Initiated_step1");

  Update_ue_state_matrix(COMMON_PROCEDURE_INITIATED1,
                         STATE_EVENT_SEC_MODE_COMPLETE, SESSION_NULL,
                         COMMON_PROCEDURE_INITIATED2, SESSION_NULL,
                         "Common_procedure_Initiated_step2");
  Update_ue_state_matrix(COMMON_PROCEDURE_INITIATED2, STATE_EVENT_REG_COMPLETE,
                         SESSION_NULL, REGISTERED_CONNECTED, SESSION_NULL,
                         "Register_complete");

  Update_ue_state_matrix(REGISTERED_CONNECTED, STATE_EVENT_DEREGISTER,
                         SESSION_NULL, DEREGISTERED_INITIATED, SESSION_NULL,
                         "Deregister_Initiated");

  Update_ue_state_matrix(REGISTERED_IDLE, STATE_EVENT_DEREGISTER, SESSION_NULL,
                         DEREGISTERED_INITIATED, SESSION_NULL,
                         "Deregister_Initiated");

  Update_ue_state_matrix(DEREGISTERED_INITIATED, STATE_EVENT_DEREGISTER,
                         SESSION_NULL, DEREGISTERED, SESSION_NULL,
                         "Deregister_Completed");

  Update_ue_state_matrix(DEREGISTERED, STATE_EVENT_DEREGISTER, SESSION_NULL,
                         DEREGISTERED, SESSION_NULL, "Deregister_Completed");

  Update_ue_state_matrix(REGISTERED_CONNECTED, STATE_EVENT_CONTEXT_RELEASE,
                         SESSION_NULL, REGISTERED_IDLE, SESSION_NULL,
                         "Idle_mode_procedure");

  Update_ue_state_matrix(REGISTERED_IDLE, STATE_EVENT_REG_REQUEST, SESSION_NULL,
                         COMMON_PROCEDURE_INITIATED1, SESSION_NULL,
                         "Common_procedure_Initiated_step1");

  /* PDU session State Transitions*/
  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_ESTABLISHMENT_REQUEST, SESSION_NULL,
                         REGISTERED_CONNECTED, CREATING, "PDU_Creating");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT, CREATING,
                         REGISTERED_CONNECTED, ACTIVE, "PDU_Created");
  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_RELEASE_COMPLETE, ACTIVE,
                         REGISTERED_CONNECTED, RELEASED, "PDU_Release");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_RELEASE_COMPLETE, SESSION_NULL,
                         DEREGISTERED, SESSION_NULL, "PDU_Release");

  Update_ue_state_matrix(DEREGISTERED, STATE_PDU_SESSION_RELEASE_COMPLETE,
                         SESSION_NULL, DEREGISTERED, SESSION_NULL,
                         "PDU_Release");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_RELEASE_COMPLETE, INACTIVE,
                         REGISTERED_CONNECTED, RELEASED, "PDU_Release");
  Update_ue_state_matrix(REGISTERED_IDLE, STATE_PDU_SESSION_RELEASE_COMPLETE,
                         INACTIVE, REGISTERED_IDLE, RELEASED, "PDU_Release");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_MODIFICATION_REQUEST, ACTIVE,
                         REGISTERED_CONNECTED, SESSION_MODIFICATION,
                         "PDU_Session_Modification_Request");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_MODIFICATION_COMPLETE,
                         SESSION_MODIFICATION, REGISTERED_CONNECTED, ACTIVE,
                         "PDU_Session_Modification_Complete");

  Update_ue_state_matrix(REGISTERED_CONNECTED,
                         STATE_PDU_SESSION_MODIFICATION_COMMAND_REJECT,
                         SESSION_MODIFICATION, REGISTERED_CONNECTED, ACTIVE,
                         "PDU_Session_Modification_Reject");
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/*
 * Handle Deregister->CPI1 and CPI1->CPI2 UE state Transitions
 *
 * @param current UE state,state event,current pdu session
 * state,ue_context,amf_context
 * @return int for success or failure
 */
status_code_e ue_state_handle_message_initial(
    m5gmm_state_t cur_state, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s* ue_m5gmm_context, amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (ue_state_matrix[cur_state][event][session_state].handler.func) {
    OAILOG_INFO(
        LOG_NAS_AMF,
        "[%s] Event triggered, UE switches from [%s] to [%s] and PDU session "
        ":[%s]",
        get_state_event_string((state_events)event).c_str(),
        get_ue_state_string(cur_state).c_str(),
        get_ue_state_string(
            ue_state_matrix[cur_state][event][session_state].next_state)
            .c_str(),
        get_session_state_string(session_state).c_str());

    ue_m5gmm_context->mm_state =
        ue_state_matrix[cur_state][event][session_state].next_state;

    OAILOG_FUNC_RETURN(
        LOG_AMF_APP,
        reinterpret_cast<status_code_e (*)(amf_context_t*)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            amf_context));
  } else {
    OAILOG_ERROR(LOG_NAS_AMF, "FSM %s: No Proper Handler Found\n", __func__);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
}

/*
 * Handle CPI2->Register_connected UE state Transitions
 *
 * @param current UE state,state event,current pdu session
 * state,ue_context,ue_id,smf_msg,amf_cause,decode_status
 * @return int for success or failure
 */
status_code_e ue_state_handle_message_reg_conn(
    m5gmm_state_t cur_state, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s* ue_m5gmm_context, amf_ue_ngap_id_t ue_id,
    bstring smf_msg_pP, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (ue_state_matrix[cur_state][event][session_state].handler.func) {
    ue_m5gmm_context->mm_state =
        ue_state_matrix[cur_state][event][session_state].next_state;
    OAILOG_FUNC_RETURN(
        LOG_AMF_APP,
        reinterpret_cast<status_code_e (*)(
            amf_ue_ngap_id_t, bstring, int,
            const amf_nas_message_decode_status_t)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            ue_id, smf_msg_pP, amf_cause, decode_status));
  } else {
    OAILOG_ERROR(LOG_NAS_AMF, "FSM %s: No Proper Handler Found\n", __func__);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
}

/*
 * Handle Register_Connected->Deregister_Initiated and
 * Deregister_initiated->Deregister UE state Transitions
 *
 * @param current UE state,state event,current pdu session
 * state,ue_context,ue_id
 * @return int for success or failure
 */
status_code_e ue_state_handle_message_dereg(
    m5gmm_state_t cur_state, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s* ue_m5gmm_context, amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (ue_state_matrix[cur_state][event][session_state].handler.func) {
    ue_m5gmm_context->mm_state =
        ue_state_matrix[cur_state][event][session_state].next_state;

    OAILOG_FUNC_RETURN(
        LOG_AMF_APP,
        reinterpret_cast<status_code_e (*)(amf_ue_ngap_id_t)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            ue_id));
  } else {
    OAILOG_ERROR(LOG_NAS_AMF, "FSM %s: No Proper Handler Found\n", __func__);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
}

/*
 * Handle NULL->Creating,Creating->Created,Created->Release PDU state
 * Transitions
 *
 * @param current UE state,state event,current pdu session
 * state,ue_context,amf_smf_msg,imsi,pdu_session_resp,ue_id
 * @return int for success or failure
 */
status_code_e pdu_state_handle_message(
    m5gmm_state_t cur_state, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s* ue_m5gmm_context, amf_smf_t amf_smf_msg, char* imsi,
    itti_n11_create_pdu_session_response_t* pdu_session_resp, uint32_t ue_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  std::shared_ptr<smf_context_t> smf_ctx =
      amf_get_smf_context_by_pdu_session_id(ue_m5gmm_context,
                                            amf_smf_msg.pdu_session_id);

  if (ue_state_matrix[cur_state][event][session_state].handler.func) {
    OAILOG_INFO(
        LOG_NAS_AMF,
        "[%s] Event triggered, PDU session state switches from [%s] to [%s] "
        "and "
        "UE state :[%s]",
        get_state_event_string((state_events)event).c_str(),
        get_session_state_string(session_state).c_str(),
        get_session_state_string(
            ue_state_matrix[cur_state][event][session_state].next_sess_state)
            .c_str(),
        get_ue_state_string(cur_state).c_str());

    switch (event) {
      case STATE_PDU_SESSION_ESTABLISHMENT_REQUEST:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        return reinterpret_cast<status_code_e (*)(amf_smf_establish_t*, char*)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            &amf_smf_msg.u.establish, imsi);
        break;
      case STATE_PDU_SESSION_RELEASE_COMPLETE:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        update_amf_app_stats_pdusessions_ue_sub();
        return reinterpret_cast<status_code_e (*)(amf_smf_release_t*, char*)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            &amf_smf_msg.u.release, imsi);
        break;
      case STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        update_amf_app_stats_pdusessions_ue_add();

        return reinterpret_cast<status_code_e (*)(
            itti_n11_create_pdu_session_response_t*, uint32_t)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            pdu_session_resp, ue_id);
        break;
      case STATE_PDU_SESSION_MODIFICATION_REQUEST:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        return reinterpret_cast<status_code_e (*)(
            itti_n11_create_pdu_session_response_t*, uint32_t)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            pdu_session_resp, ue_id);
        break;
      case STATE_PDU_SESSION_MODIFICATION_COMPLETE:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        return reinterpret_cast<status_code_e (*)(amf_smf_establish_t*, char*)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            &amf_smf_msg.u.establish, imsi);
        break;
      case STATE_PDU_SESSION_MODIFICATION_COMMAND_REJECT:
        smf_ctx->pdu_session_state =
            ue_state_matrix[cur_state][event][session_state].next_sess_state;
        return reinterpret_cast<status_code_e (*)(amf_smf_establish_t*, char*)>(
            ue_state_matrix[cur_state][event][session_state].handler.func)(
            &amf_smf_msg.u.establish, imsi);
        break;
      default:
        OAILOG_ERROR(LOG_NAS_AMF, "FSM %s: No Proper Handler Found\n",
                     __func__);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}
std::string get_state_event_string(state_events event) {
  std::string eventStr;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  switch (event) {
    case STATE_EVENT_REG_REQUEST:
      eventStr = "STATE_EVENT_REG_REQUEST";
      break;
    case STATE_EVENT_SEC_MODE_COMPLETE:
      eventStr = "STATE_EVENT_SEC_MODE_COMPLETE";
      break;
    case STATE_EVENT_REG_COMPLETE:
      eventStr = "STATE_EVENT_REG_COMPLETE";
      break;
    case STATE_EVENT_DEREGISTER:
      eventStr = "STATE_EVENT_DEREGISTER";
      break;
    case STATE_PDU_SESSION_ESTABLISHMENT_REQUEST:
      eventStr = "STATE_PDU_SESSION_ESTABLISHMENT_REQUEST";
      break;
    case STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT:
      eventStr = "STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT";
      break;
    case STATE_PDU_SESSION_RELEASE_COMPLETE:
      eventStr = "STATE_PDU_SESSION_RELEASE_COMPLETE";
      break;
    case STATE_EVENT_CONTEXT_RELEASE:
      eventStr = "STATE_EVENT_CONTEXT_RELEASE";
      break;
    case STATE_PDU_SESSION_MODIFICATION_REQUEST:
      eventStr = "STATE_PDU_SESSION_MODIFICATION_REQUEST";
      break;
    case STATE_PDU_SESSION_MODIFICATION_COMPLETE:
      eventStr = "STATE_PDU_SESSION_MODIFICATION_COMPLETE";
      break;
    case STATE_PDU_SESSION_MODIFICATION_COMMAND_REJECT:
      eventStr = "STATE_PDU_SESSION_MODIFICATION_COMMAND_REJECT";
      break;
    default:
      eventStr = "UNKNOWN_EVENT";
      break;
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, eventStr);
}

std::string get_session_state_string(SMSessionFSMState s) {
  std::string sessStateStr;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  switch (s) {
    case SESSION_NULL:
      sessStateStr = "SESSION_NULL";
      break;
    case CREATING:
      sessStateStr = "CREATING";
      break;
    case CREATE:
      sessStateStr = "CREATE";
      break;
    case ACTIVE:
      sessStateStr = "ACTIVE";
      break;
    case INACTIVE:
      sessStateStr = "INACTIVE";
      break;
    case PENDING_RELEASE:
      sessStateStr = "PENDING_RELEASE";
      break;
    case RELEASED:
      sessStateStr = "RELEASED";
      break;
    case SESSION_MODIFICATION:
      sessStateStr = "MODIFICATION";
      break;
    default:
      sessStateStr = "UNKNOWN_SESSION_STATE";
      break;
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, sessStateStr);
}

std::string get_ue_state_string(m5gmm_state_t ueState) {
  std::string ueStateStr;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  switch (ueState) {
    case DEREGISTERED:
      ueStateStr = "DEREGISTERED";
      break;
    case REGISTERED_IDLE:
      ueStateStr = "REGISTERED_IDLE";
      break;
    case REGISTERED_CONNECTED:
      ueStateStr = "REGISTERED_CONNECTED";
      break;
    case DEREGISTERED_INITIATED:
      ueStateStr = "DEREGISTERED_INITIATED";
      break;
    case PENDING_RELEASE_RESPONSE:
      ueStateStr = "PENDING_RELEASE_RESPONSE";
      break;
    case COMMON_PROCEDURE_INITIATED1:
      ueStateStr = "COMMON_PROCEDURE_INITIATED1";
      break;
    case COMMON_PROCEDURE_INITIATED2:
      ueStateStr = "COMMON_PROCEDURE_INITIATED2";
      break;
    default:
      ueStateStr = "UE_UNKNOWN_STATE";
      break;
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, ueStateStr);
}
}  // namespace magma5g
