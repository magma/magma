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
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "3gpp_24.501.h"
#include "conversions.h"
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include "M5gNasMessage.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_authentication.h"
#include "amf_as.h"
#include "amf_sap.h"
#include "amf_recv.h"

#define M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH 1
#define INVALID_IMSI64 (imsi64_t) 0
#define INVALID_AMF_UE_NGAP_ID 0x0

namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;
amf_as_data_t amf_data_sec;
nas_amf_smc_proc_t smc_proc;
static int amf_registration_failure_authentication_cb(
    amf_context_t* amf_context);
static int amf_start_registration_proc_security(
    amf_context_t* amf_context, nas_amf_registration_proc_t* registration_proc);
static int amf_registration(amf_context_t* amf_context);
static int amf_registration_failure_security_cb(amf_context_t* amf_context);

static int amf_registration_reject(
    amf_context_t* amf_context, nas_amf_registration_proc_t* nas_base_proc);
static int registration_accept_t3550_handler(
    zloop_t* loop, int timer_id, void* arg);

/***************************************************************************
**                                                                        **
** Name:  amf_registration_success_authentication_cb()                    **
**                                                                        **
** Description: Callback for successful authentication                    **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_registration_success_authentication_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  OAILOG_DEBUG(LOG_NAS_AMF, " Authentication procedure is successful");
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    rc = amf_start_registration_proc_security(amf_context, registration_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_start_registration_proc_authentication()                    **
**                                                                        **
** Description:Validates amf_context and invokes authentication procedure **
**                                                                        **
**                                                                        **
***************************************************************************/
static int amf_start_registration_proc_authentication(
    amf_context_t* amf_context,
    nas_amf_registration_proc_t* registration_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  if ((amf_context) && (registration_proc)) {
    rc = amf_proc_authentication(
        amf_context, &registration_proc->amf_spec_proc,
        amf_registration_success_authentication_cb,
        amf_registration_failure_authentication_cb);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    nas_new_registration_procedure                                **
**                                                                        **
** Description: Allocate and initialize amf_procedures                    **
**                                                                        **
**                                                                        **
***************************************************************************/
nas_amf_registration_proc_t* nas_new_registration_procedure(
    ue_m5gmm_context_s* ue_ctxt) {
  amf_context_t* amf_context = &ue_ctxt->amf_context;

  if (!(amf_context->amf_procedures)) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "From nas_new_registration_procedure allocating for "
        "amf_procedures\n");
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }
  amf_context->amf_procedures->amf_specific_proc =
      new nas_amf_specific_proc_t();

  amf_context->amf_procedures->amf_specific_proc->amf_proc.base_proc.nas_puid =
      __sync_fetch_and_add(&nas_puid, 1);
  amf_context->amf_procedures->amf_specific_proc->amf_proc.base_proc.type =
      NAS_PROC_TYPE_AMF;
  amf_context->amf_procedures->amf_specific_proc->amf_proc.type =
      NAS_AMF_PROC_TYPE_CONN_MNGT;
  amf_context->amf_procedures->amf_specific_proc->type =
      AMF_SPEC_PROC_TYPE_REGISTRATION;

  nas_amf_registration_proc_t* proc =
      (nas_amf_registration_proc_t*)
          amf_context->amf_procedures->amf_specific_proc;
  proc->registration_accept_sent = 0;

  /* TIMERS_PLACE_HOLDER */

  OAILOG_TRACE(
      LOG_NAS_AMF, "New AMF_SPEC_PROC_TYPE_REGISTRATION initialized\n");
  return proc;
}

/***************************************************************************
**                                                                        **
** Name:    amf_proc_create_procedure_registration_request()              **
**                                                                        **
** Description: Create registration request procedure                     **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_proc_create_procedure_registration_request(
    ue_m5gmm_context_s* ue_ctx, amf_registration_request_ies_t* ies) {
  nas_amf_registration_proc_t* reg_proc =
      nas_new_registration_procedure(ue_ctx);
  if ((reg_proc)) {
    reg_proc->ies   = ies;
    reg_proc->ue_id = ue_ctx->amf_ue_ngap_id;
  }
}

/***************************************************************************
**                                                                        **
** Name:   amf_proc_registration_request                                  **
**                                                                        **
** Description: Handler for processing registration request               **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_proc_registration_request(
    amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    amf_registration_request_ies_t* ies) {
  int rc = RETURNerror;
  ue_m5gmm_context_s ue_ctx;
  imsi64_t imsi64                      = INVALID_IMSI64;
  ue_m5gmm_context_s* ue_m5gmm_context = NULL;
  if (ies->imsi) {
    imsi64 = amf_imsi_to_imsi64(ies->imsi);
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "During initial registration request "
        "SUPI as IMSI converted to imsi64 " IMSI_64_FMT " = ",
        imsi64);
  } else if (ies->imei) {
    char imei_str[16];
    IMEI_TO_STRING(ies->imei, imei_str, 16);
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (IMEI = %s ) \n",
        ue_id, imei_str);
  }

  ue_m5gmm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_m5gmm_context == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "ue context not found for the ue_id=%u\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  ue_m5gmm_context->amf_context.amf_procedures = NULL;
  ue_m5gmm_context->amf_context.is_dynamic     = false;
  ue_m5gmm_context->amf_ue_ngap_id             = ue_id;

  if (!(is_nas_specific_procedure_registration_running(
          &ue_m5gmm_context->amf_context))) {
    amf_proc_create_procedure_registration_request(ue_m5gmm_context, ies);
  }

  OAILOG_INFO(LOG_AMF_APP, "ue_m5gmm_context %p\n", ue_m5gmm_context);
  rc = ue_state_handle_message_initial(
      ue_m5gmm_context->mm_state, STATE_EVENT_REG_REQUEST, SESSION_NULL,
      ue_m5gmm_context, &ue_m5gmm_context->amf_context);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_proc_registration_reject                                  **
**                                                                        **
** Description:  Handler to trigger registration reject                   **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_proc_registration_reject(
    amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                 = RETURNerror;
  amf_context_t* amf_ctx = amf_context_get(ue_id);
  nas_amf_registration_proc_t* registration_proc =
      (nas_amf_registration_proc_t*) (amf_ctx->amf_procedures
                                          ->amf_specific_proc);
  registration_proc->amf_cause = amf_cause;
  if (amf_ctx) {
    if (is_nas_specific_procedure_registration_running(amf_ctx)) {
      rc = amf_registration_reject(amf_ctx, registration_proc);
      amf_sap_t amf_sap;
      amf_sap.primitive                   = AMFREG_REGISTRATION_REJ;
      amf_sap.u.amf_reg.ue_id             = ue_id;
      amf_sap.u.amf_reg.ctx               = amf_ctx;
      amf_sap.u.amf_reg.notify            = false;
      amf_sap.u.amf_reg.free_proc         = true;
      amf_sap.u.amf_reg.u.registered.proc = registration_proc;
      rc                                  = amf_sap_send(&amf_sap);
    } else {
      rc = amf_registration_reject(amf_ctx, registration_proc);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_registration_reject                                       **
**                                                                        **
** Description: Notify AS-SAP about Registration Reject message           **
**                                                                        **
**                                                                        **
***************************************************************************/
static int amf_registration_reject(
    amf_context_t* amf_context, nas_amf_registration_proc_t* nas_base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  amf_sap_t amf_sap;
  nas_amf_registration_proc_t* registration_proc =
      (nas_amf_registration_proc_t*) nas_base_proc;
  OAILOG_WARNING(
      LOG_AMF_APP, "AMF-PROC  - AMF Registration procedure not accepted ");
  /*
   * Notify AMF-AS SAP that Registration Reject message has to be sent
   * onto the network
   */
  amf_sap.primitive                      = AMFREG_REGISTRATION_REJ;
  amf_sap.u.amf_as.u.establish.ue_id     = registration_proc->ue_id;
  amf_sap.u.amf_as.u.establish.amf_cause = registration_proc->amf_cause;
  amf_sap.u.amf_as.u.establish.nas_info  = AMF_AS_NAS_INFO_REGISTERD;

  if (registration_proc->amf_cause != AMF_CAUSE_SMF_FAILURE) {
    amf_sap.u.amf_as.u.establish.nas_msg = NULL;
  } else if (registration_proc->amf_msg_out) {
    amf_sap.u.amf_as.u.establish.nas_msg = registration_proc->amf_msg_out;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  // Setup 5G CN NAS security data
  if (amf_context) {
    amf_data_sec.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.establish.sctx, &amf_context->_security, false,
        false);
  } else {
    amf_data_sec.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.establish.sctx, NULL, false, false);
  }
  rc = amf_sap_send(&amf_sap);
  increment_counter(
      "ue_Registration", 1, 1, "action", "Registration_reject_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_registration_run_procedure                                **
**                                                                        **
** Description: Functions that will initiate AMF common procedures        **
**                                                                        **
**                                                                        **
***************************************************************************/

int amf_registration_run_procedure(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);
  if (registration_proc == NULL) {
    OAILOG_DEBUG(
        LOG_AMF_APP, " registration_proc NULL, from %s\n", __FUNCTION__);
  }
  OAILOG_DEBUG(
      LOG_NAS_AMF, " decode_status.mac_matched :%d",
      registration_proc->ies->decode_status.mac_matched);
  OAILOG_DEBUG(
      LOG_NAS_AMF, " decode_status.integrity_protected_message :%d",
      registration_proc->ies->decode_status.integrity_protected_message);

  if (registration_proc) {
    if (registration_proc->ies->imsi) {
      /* If registratin ie is IMSI and if mac matched or
       * Intergrity type is not protected start authentication
       * procedure.
       */
      if ((registration_proc->ies->decode_status.mac_matched) ||
          !(registration_proc->ies->decode_status
                .integrity_protected_message)) {
        if (amf_context->reg_id_type != M5GSMobileIdentityMsg_SUCI_IMSI) {
          OAILOG_ERROR(LOG_AMF_APP, "ies and type mismatch \n");
          OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
        }

        // Convert recevied imsi to uint64
        imsi64_t imsi64 = amf_imsi_to_imsi64(registration_proc->ies->imsi);

        amf_ctx_set_valid_imsi(
            amf_context, registration_proc->ies->imsi, imsi64);

        rc = amf_start_registration_proc_authentication(
            amf_context, registration_proc);
        if (rc != RETURNok) {
          OAILOG_ERROR(
              LOG_NAS_AMF,
              "Failed to start registration authentication procedure! \n");
        }

      } else {
        // force identification, even if not necessary
        rc = amf_proc_identification(
            amf_context, (nas_amf_proc_t*) registration_proc,
            IDENTITY_TYPE_2_IMSI,

            amf_registration_success_identification_cb,
            amf_registration_failure_identification_cb);
      }
    } else if (registration_proc->ies->guti) {
      /* TODO: Currently we are not receving GUTI during intial
       * Registration procedure and in future this code can be used.
       */
      rc = amf_proc_identification(
          amf_context, (nas_amf_proc_t*) registration_proc,
          IDENTITY_TYPE_2_IMSI, amf_registration_success_identification_cb,
          amf_registration_failure_identification_cb);
    } else {
      OAILOG_ERROR(LOG_NAS_AMF, "Unsupported IE type! \n");
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_registration_success_identification_cb()                    **
**                                                                        **
** Description: Callback for successful identification                    **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_registration_success_identification_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  OAILOG_DEBUG(LOG_NAS_AMF, " Identification procedure success\n");
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    rc = amf_start_registration_proc_authentication(
        amf_context, registration_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_registration_failure_identification_cb()                    **
**                                                                        **
** Description: Callback for identification failure                       **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_registration_failure_identification_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  // TODO nagetive scenario will be taken care in future.
  int rc = RETURNerror;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_registration_failure_authentication_cb()                    **
**                                                                        **
** Description: Callback for authentication failure                       **
**                                                                        **
**                                                                        **
***************************************************************************/
static int amf_registration_failure_authentication_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    registration_proc->amf_cause = amf_context->amf_cause;
    amf_sap_t amf_sap;
    amf_sap.primitive                   = AMFREG_REGISTRATION_REJ;
    amf_sap.u.amf_reg.ue_id             = registration_proc->ue_id;
    amf_sap.u.amf_reg.ctx               = amf_context;
    amf_sap.u.amf_reg.notify            = true;
    amf_sap.u.amf_reg.free_proc         = true;
    amf_sap.u.amf_reg.u.registered.proc = registration_proc;
    rc                                  = amf_sap_send(&amf_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_registration_success_security_cb()                          **
**                                                                        **
** Description: Callback for successful security mode complete            **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_registration_success_security_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    rc = amf_registration(amf_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_start_registration_proc_security()                          **
**                                                                        **
** Description: Create new security context and initiate SMC procedures   **
**                                                                        **
**                                                                        **
***************************************************************************/
static int amf_start_registration_proc_security(
    amf_context_t* amf_context,
    nas_amf_registration_proc_t* registration_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  if ((amf_context) && (registration_proc)) {
    /*
     * Create new NAS security context
     */
    smc_proc.amf_ctx_clear_security(amf_context);
    rc = amf_proc_security_mode_control(
        amf_context, &registration_proc->amf_spec_proc, registration_proc->ksi,
        amf_registration_success_security_cb,
        amf_registration_failure_security_cb);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:  amf_registration_failure_security_cb                            **
**                                                                        **
** Description: Callback for security mode command failure                **
**                                                                        **
**                                                                        **
***************************************************************************/
static int amf_registration_failure_security_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  // TODO: In future implement as part of handling negative scenarios
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 ** Name:    amf_registration()                                            **
 **                                                                        **
 ** Description: Performs the registration signalling                      **
 **              procedure while a context  exists for                     **
 **	       	the incoming UE in the network.                            **
 **                                                                        **
 **              3GPP TS 24.501, section 5.5.1.2.4                         **
 **      Upon receiving the REGISTRATION REQUEST message, the AMF shall    **
 **      send an REGISTRATION ACCEPT message to the UE and start timer     **
 **      T3450.                                                            **
 **                                                                        **
 ****************************************************************************/
static int amf_registration(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_context, struct ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "ue_id= " AMF_UE_NGAP_ID_FMT
      "Start REGISTRATION_ACCEPT procedures for UE \n",
      ue_id);
  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    rc = amf_send_registration_accept(amf_context);
  }

  if (rc != RETURNok) {
    /*
     * The Registration procedure failed
     */
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "ue_id= " AMF_UE_NGAP_ID_FMT
        " AMF-PROC  - Failed to respond to registration request\n",
        ue_id);
    registration_proc->amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_send_registration_accept()                                **
 **                                                                        **
 ** Description: Sends REGISTRATION ACCEPT message and start timer T3550   **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
int amf_send_registration_accept(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  if (amf_context) {
    amf_sap_t amf_sap;
    nas_amf_registration_proc_t* registration_proc =
        get_nas_specific_procedure_registration(amf_context);
    ue_m5gmm_context_s* ue_m5gmm_context_p =
        PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context);
    amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;

    registration_proc->T3550.id = NAS5G_TIMER_INACTIVE_ID;
    if (registration_proc) {
      /*
       * The IMSI if provided by the UE
       */
      if (registration_proc->ies->imsi) {
        imsi64_t new_imsi64 = amf_imsi_to_imsi64(registration_proc->ies->imsi);
        if (new_imsi64 != amf_context->imsi64) {
          amf_ctx_set_valid_imsi(
              amf_context, registration_proc->ies->imsi, new_imsi64);
        }
      }
      /*
       * Notify AMF-AS SAP that Registaration Accept message together with an
       * Activate Pdu session Context Request message has to be sent to the UE
       */
      amf_sap.primitive = AMFAS_ESTABLISH_CNF;
      amf_sap.u.amf_as.u.establish.puid =
          registration_proc->amf_spec_proc.amf_proc.base_proc.nas_puid;
      amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
      amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_REGISTERD;

      /* GUTI have already updated in amf_context during Identification
       * response complete, now assign to amf_sap
       */
      amf_sap.u.amf_as.u.establish.guti = amf_context->m5_guti;
      amf_sap.u.amf_as.u.establish.guti.m_tmsi =
          htonl(amf_sap.u.amf_as.u.establish.guti.m_tmsi);
      rc = amf_sap_send(&amf_sap);
      /*
       * Start T3550 timer
       */
      OAILOG_INFO(LOG_AMF_APP, "Timer: registration_accept timer start\n");
      registration_proc->T3550.id = start_timer(
          &amf_app_task_zmq_ctx, REGISTRATION_ACCEPT_TIMER_EXPIRY_MSECS,
          TIMER_REPEAT_ONCE, registration_accept_t3550_handler,
          (void*) registration_proc->ue_id);
      OAILOG_INFO(
          LOG_AMF_APP,
          "Timer: Registration_accept timer T3550 with id  %d Started\n",
          registration_proc->T3550.id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

static int registration_accept_t3550_handler(
    zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(LOG_AMF_APP, "Timer: In registration_accept_t3550 handler\n");
#if 0 /* TIMER_CHANGES_REVIEW */
  amf_context_t* amf_ctx                         = NULL;
  ue_m5gmm_context_s* ue_amf_context             = NULL;
  nas_amf_registration_proc_t* registration_proc = NULL;
  amf_ue_ngap_id_t ue_id                         = 0;

  ue_id = *((amf_ue_ngap_id_t*) (arg));
  /*
   * Get the UE context
   */
  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: ue_context is NULL\n");
    return -1;
  }

  registration_proc =
      (nas_amf_registration_proc_t*)
          ue_amf_context->amf_context.amf_procedures->amf_specific_proc;

  OAILOG_INFO(
      LOG_AMF_APP, "Timer: In _registration_t3550_handler - T3550 id %p\n",
      registration_proc);
  if (registration_proc) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3550 timer   timer id %d ue id %d\n",
        registration_proc->T3550.id, registration_proc->ue_id);
    registration_proc->T3550.id = -1;

    registration_proc->retransmission_count += 1;
    OAILOG_ERROR(
        LOG_AMF_APP, "Timer: Incrementing retransmission_count to %d\n",
        registration_proc->retransmission_count);
    if (registration_proc->retransmission_count < REGISTRATION_COUNTER_MAX) {
      /* Send entity Registration request message to the UE */

      OAILOG_ERROR(
          LOG_AMF_APP,
          "Timer: timer has expired Sending Registration request again\n");
      // amf_registration_request(registration_proc); /* TO DO for negative
      // scenario */
    } else {
      /* Abort the registration procedure */

      OAILOG_ERROR(
          LOG_AMF_APP,
          "Timer: Maximum retires done hence Abort the registration "
          "procedure\n");
      return -1;
    }

    return 0;
  }

#endif /* TIMER_CHANGES_REVIEW */
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_send_registration_accept_dl_nas()                         **
 **                                                                        **
 ** Description: Builds Registration Accept message to be sent             **
 **              is NGAP : DL NAS Tx                                       **
 **                                                                        **
 **      The registration Accept message is sent by the network to the     **
 **      UE to indicate that the corresponding attach request has          **
 **      been accepted.                                                    **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
int amf_send_registration_accept_dl_nas(
    const amf_as_data_t* msg, RegistrationAcceptMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = REGISTRATION_ACCEPT_MINIMUM_LENGTH;
  amf_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  amf_msg->sec_header_type.sec_hdr        = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  amf_msg->message_type.msg_type          = REG_ACCEPT;
  amf_msg->m5gs_reg_result.sms_allowed    = NOT_ALLOWED;
  amf_msg->m5gs_reg_result.reg_result_val = M3GPP_ACCESS;

  /* 5GS mobile identity IE */
  size += MOBILE_IDENTITY_MAX_LENGTH;
  amf_msg->mobile_id.mobile_identity.guti.odd_even = EVEN_IDENTITY_DIGITS;
  amf_msg->mobile_id.iei                           = M5GS_MOBILE_IDENTITY;
  amf_msg->mobile_id.len = M5GSMobileIdentityMsg_GUTI_LENGTH;
  amf_msg->mobile_id.mobile_identity.guti.type_of_identity =
      M5GSMobileIdentityMsg_GUTI;
  amf_msg->mobile_id.mobile_identity.guti.mcc_digit1 =
      msg->guti->guamfi.plmn.mcc_digit1;
  amf_msg->mobile_id.mobile_identity.guti.mcc_digit2 =
      msg->guti->guamfi.plmn.mcc_digit2;
  amf_msg->mobile_id.mobile_identity.guti.mcc_digit3 =
      msg->guti->guamfi.plmn.mcc_digit3;
  amf_msg->mobile_id.mobile_identity.guti.mnc_digit1 =
      msg->guti->guamfi.plmn.mnc_digit1;
  amf_msg->mobile_id.mobile_identity.guti.mnc_digit2 =
      msg->guti->guamfi.plmn.mnc_digit2;
  amf_msg->mobile_id.mobile_identity.guti.mnc_digit3 =
      msg->guti->guamfi.plmn.mnc_digit3;
  uint8_t* offset;
  offset                                        = (uint8_t*) &msg->guti->m_tmsi;
  amf_msg->mobile_id.mobile_identity.guti.tmsi1 = *offset;
  offset++;
  amf_msg->mobile_id.mobile_identity.guti.tmsi2 = *offset;
  offset++;
  amf_msg->mobile_id.mobile_identity.guti.tmsi3 = *offset;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_registration_complete()                              **
 **                                                                        **
 ** Description: Completion of Registration Procedure                      **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      smf_msg_pP:       The received SMF message                        **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                                 **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_proc_registration_complete(
    amf_ue_ngap_id_t ue_id, bstring smf_msg_pP, int amf_cause,
    const amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_amf_context             = NULL;
  nas_amf_registration_proc_t* registration_proc = NULL;
  int rc                                         = RETURNerror;
  amf_sap_t amf_sap;
  amf_context_t* amf_ctx = NULL;

  /*
   * Get the UE context
   */
  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context) {
    if (is_nas_specific_procedure_registration_running(
            &ue_amf_context->amf_context)) {
      registration_proc =
          (nas_amf_registration_proc_t*)
              ue_amf_context->amf_context.amf_procedures->amf_specific_proc;
      amf_ctx = &ue_amf_context->amf_context;

      stop_timer(&amf_app_task_zmq_ctx, registration_proc->T3550.id);
      OAILOG_INFO(
          LOG_AMF_APP, "Timer: after stop registration timer with id = %d\n",
          registration_proc->T3550.id);
      registration_proc->T3550.id = NAS5G_TIMER_INACTIVE_ID;

      /*
       * Upon receiving an REGISTRATION COMPLETE message, the AMF shall enter
       * state AMF-REGISTERED and consider the GUTI sent in the REGISTRATION
       * ACCEPT message as valid.
       */
      amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_GUTI);
    } else {
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "UE " AMF_UE_NGAP_ID_FMT
          " REGISTRATION COMPLETE discarded (AMF procedure not found)\n",
          ue_id);
      bdestroy((bstring)(smf_msg_pP));
    }
  } else {
    OAILOG_WARNING(LOG_NAS_AMF, "UE Context not found..\n");
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "UE " AMF_UE_NGAP_ID_FMT
        " REGISTRATION COMPLETE discarded (context not found)\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }
  /*
   * Set the network registration indicator
   */
  ue_amf_context->amf_context.is_registered = true;
  /*
   * Notify AMF that registration procedure has successfully completed
   */
  amf_sap.primitive                   = AMFREG_REGISTRATION_CNF;
  amf_sap.u.amf_reg.ue_id             = ue_id;
  amf_sap.u.amf_reg.ctx               = &ue_amf_context->amf_context;
  amf_sap.u.amf_reg.notify            = true;
  amf_sap.u.amf_reg.free_proc         = true;
  amf_sap.u.amf_reg.u.registered.proc = registration_proc;
  rc                                  = amf_sap_send(&amf_sap);
  if (rc == RETURNok) {
    /*
     * Send AMF Information after handling Registration Complete message
     * TODO this logic will handled in future when PDU Session Establish
     * resquest comes along with Inital Registration request.
     */
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_registration_complete_response()                    **
 **                                                                        **
 ** Description: Processes registration Complete message                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      msg:       The received AMF message                               **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                                 **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_handle_registration_complete_response(
    amf_ue_ngap_id_t ue_id, RegistrationCompleteMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc;
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "AMFAS-SAP - received registration complete message for ue_id = (%u)\n",
      ue_id);
  /*
   * Execute the registration procedure completion
   */
  rc = amf_proc_registration_complete(ue_id, msg->smf_pdu, amf_cause, status);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_amf_information()                                     **
 **                                                                        **
 ** Description: Send AMF Information after handling                       **
 **              Registration Complete message                             **
 **                                                                        **
 ** Inputs:  ue_amf_ctx:   UE context                                      **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/

int amf_proc_amf_information(ue_m5gmm_context_s* ue_amf_ctx) {
  int rc = RETURNerror;
  amf_sap_t amf_sap;
  amf_as_data_t* amf_as  = &amf_sap.u.amf_as.u.data;
  amf_context_t* amf_ctx = &(ue_amf_ctx->amf_context);
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  /*
   * Setup NAS information message to transfer
   */
  amf_as->nas_info = AMF_AS_NAS_AMF_INFORMATION;
  amf_as->nas_msg  = "";
  /*
   * Set the UE identifier
   */
  amf_as->ue_id = ue_amf_ctx->amf_ue_ngap_id;
  /*
   * Setup EPS NAS security data
   */
  amf_as->amf_as_set_security_data(
      &amf_as->sctx, &amf_ctx->_security, false, true);
  /*
   * Notify AMF-AS SAP that Registration Accept message has to be sent to the
   * network
   */
  amf_sap.primitive = AMFAS_DATA_REQ;
  rc                = amf_sap_send(&amf_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***********************************************************************
 ** Name:    amf_reg_send()                                           **
 **                                                                   **
 ** Description: Processes the AMFREG Service Access Point primitive  **
 **                                                                   **
 ** Inputs:  msg:       The AMFREG-SAP primitive to process           **
 **      Others:    None                                              **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    RETURNok, RETURNerror                             **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
int amf_reg_send(amf_sap_t* const msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;
  // TODO in future it will be implemented based on request of
  // PDU session establishment with inital registration
  amf_primitive_t primitive          = msg->primitive;
  amf_reg_t* evt                     = &msg->u.amf_reg;
  amf_context_t* amf_ctx             = msg->u.amf_reg.ctx;
  ue_m5gmm_context_s* ue_amf_context = NULL;

  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(evt->ue_id);

  if (!ue_amf_context) {
    OAILOG_ERROR(
        LOG_NAS_AMF, "Ue context not found for the ue id %u\n", evt->ue_id);
    rc = RETURNerror;
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  if (evt && amf_ctx) {
    switch (primitive) {
      case AMFREG_REGISTRATION_CNF: {
        OAILOG_DEBUG(LOG_NAS_AMF, "AMFREG_REGISTRATION_CNF");
        if (evt->free_proc) {
          amf_delete_registration_proc(amf_ctx);
        }

        /* Update the state */
        // ue_amf_context->mm_state = REGISTERED_CONNECTED;
      } break;
      case AMFREG_COMMON_PROC_REJ: {
      }
      default: {}
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***********************************************************************
 ** Name:    amf_delete_registration_proc()                       **
 **                                                                   **
 ** Description: deletes the nas registration specific procedure      **
 **                                                                   **
 ** Inputs:  amf_ctx:       The Amf context to process                **
 **      Others:    None                                              **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void amf_delete_registration_proc(amf_context_t* amf_ctx) {
  nas_amf_registration_proc_t* proc =
      get_nas_specific_procedure_registration(amf_ctx);

  if (proc) {
    if (proc->ies) {
      amf_delete_registration_ies(&proc->ies);
    }
  }

  amf_delete_child_procedures(amf_ctx, (nas5g_base_proc_t*) proc);
}  // namespace magma5g

/***********************************************************************
 ** Name:    amf_delete_registration_ies()                            **
 **                                                                   **
 ** Description: deletes the nas registration specific ies            **
 **                                                                   **
 ** Inputs:  ies:   The registration to delete                        **
 **      Others:    None                                              **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void amf_delete_registration_ies(amf_registration_request_ies_t** ies) {
  if ((*ies)->imsi) {
    delete_wrapper((void**) &(*ies)->imsi);
  }

  if ((*ies)->guti) {
    delete_wrapper((void**) &(*ies)->guti);
  }

  if ((*ies)->imei) {
    delete_wrapper((void**) &(*ies)->imei);
  }

  if ((*ies)->drx_parameter) {
    delete_wrapper((void**) &(*ies)->drx_parameter);
  }

  if ((*ies)->last_visited_registered_tai) {
    delete_wrapper((void**) &(*ies)->last_visited_registered_tai);
  }
}

/***********************************************************************
 ** Name:    amf_delete_child_procedures()                            **
 **                                                                   **
 ** Description: deletes the nas registration specific child          **
 **              child procedures                                     **
 **                                                                   **
 ** Inputs:  amf_ctx:   The amf context                               **
 **          parent_proc: nas 5g base proc                            **
 **                                                                   **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void amf_delete_child_procedures(
    amf_context_t* amf_ctx, struct nas5g_base_proc_t* const parent_proc) {
  if (amf_ctx && amf_ctx->amf_procedures) {
    nas_amf_common_procedure_t* p1 =
        LIST_FIRST(&amf_ctx->amf_procedures->amf_common_procs);
    nas_amf_common_procedure_t* p2 = NULL;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      if (((nas5g_base_proc_t*) p1->proc)->parent == parent_proc) {
        amf_delete_common_procedure(amf_ctx, &p1->proc);
      }
      p1 = p2;
    }
  }
}

/***********************************************************************
 ** Name:    amf_delete_common_procedure()                            **
 **                                                                   **
 ** Description: deletes the nas registration specific common         **
 **              procedures                                           **
 **                                                                   **
 ** Inputs:  amf_ctx:   The amf context                               **
 **          proc: nas amf common proc                                **
 **                                                                   **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void amf_delete_common_procedure(
    amf_context_t* amf_ctx, nas_amf_common_proc_t** proc) {
  if (*proc) {
    /* delete proc content */
    switch ((*proc)->type) {
      case AMF_COMM_PROC_AUTH: {
      } break;
      case AMF_COMM_PROC_SMC: {
      } break;
      case AMF_COMM_PROC_IDENT: {
      } break;
      default: {}
    }

    // remove proc from list
    if (amf_ctx->amf_procedures) {
      nas_amf_common_procedure_t* p1 =
          LIST_FIRST(&amf_ctx->amf_procedures->amf_common_procs);
      nas_amf_common_procedure_t* p2 = NULL;
      while (p1) {
        p2 = LIST_NEXT(p1, entries);
        if (p1->proc == (nas_amf_common_proc_t*) (*proc)) {
          LIST_REMOVE(p1, entries);
          delete_wrapper((void**) &p1->proc);
          delete_wrapper((void**) &p1);
          return;
        }
        p1 = p2;
      }
    }
  }

  // if not found in list, free it anyway
  if (*proc) {
    delete_wrapper((void**) proc);
  }
}

/***********************************************************************
 ** Name:    delete_wrapper()                                         **
 **                                                                   **
 ** Description: deletes the memory                                   **
 **                                                                   **
 ** Inputs: ptr:   pointer to be freed                                **
 **                                                                   **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void delete_wrapper(void** ptr) {
  if (ptr) {
    delete (*ptr);
    *ptr = NULL;
  }
}

}  // namespace magma5g
