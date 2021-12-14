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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"
#include "include/amf_client_servicer.h"

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
int amf_start_registration_proc_authentication(
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
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }
  amf_context->amf_procedures->amf_specific_proc =
      reinterpret_cast<nas_amf_specific_proc_t*>(
          new nas_amf_registration_proc_t());

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
    char imei_str[MAX_IMEISV_SIZE];
    IMEI_TO_STRING(ies->imei, imei_str, MAX_IMEISV_SIZE);
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (IMEI = %s ) \n",
        ue_id, imei_str);
  }

  ue_m5gmm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_m5gmm_context == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "ue context not found for the"
        "ue_id=" AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  ue_m5gmm_context->amf_context.amf_procedures = NULL;
  ue_m5gmm_context->amf_context.is_dynamic     = false;
  ue_m5gmm_context->amf_ue_ngap_id             = ue_id;

  if (!(is_nas_specific_procedure_registration_running(
          &ue_m5gmm_context->amf_context))) {
    amf_proc_create_procedure_registration_request(ue_m5gmm_context, ies);
  } else {
    /* Update the GUTI */
    if (ies->guti) {
      nas_amf_registration_proc_t* registration_proc =
          get_nas_specific_procedure_registration(
              &(ue_m5gmm_context->amf_context));

      registration_proc->ies = ies;
    }
  }

  /* If in a connected state REGISTRATION_REQUEST is received
   * Just respond with plan response.
   * This can happen in periodic registration case.
   */
  if (ue_m5gmm_context->mm_state == REGISTERED_CONNECTED) {
    amf_registration_run_procedure(&ue_m5gmm_context->amf_context);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

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

  if (amf_ctx) {
    if (is_nas_specific_procedure_registration_running(amf_ctx)) {
      nas_amf_registration_proc_t* registration_proc =
          reinterpret_cast<nas_amf_registration_proc_t*>(
              amf_ctx->amf_procedures->amf_specific_proc);
      registration_proc->amf_cause = amf_cause;
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
      nas_amf_registration_proc_t no_registration_proc = {0};
      no_registration_proc.ue_id                       = ue_id;
      no_registration_proc.amf_cause                   = amf_cause;
      no_registration_proc.amf_msg_out                 = NULL;
      rc = amf_registration_reject(amf_ctx, &no_registration_proc);
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
  int rc            = RETURNerror;
  amf_sap_t amf_sap = {};
  nas_amf_registration_proc_t* registration_proc =
      (nas_amf_registration_proc_t*) nas_base_proc;
  OAILOG_WARNING(
      LOG_AMF_APP, "AMF-PROC  - AMF Registration procedure not accepted ");
  /*
   * Notify AMF-AS SAP that Registration Reject message has to be sent
   * onto the network
   */
  amf_sap.primitive                      = AMFAS_ESTABLISH_REJ;
  amf_sap.u.amf_as.u.establish.ue_id     = registration_proc->ue_id;
  amf_sap.u.amf_as.u.establish.amf_cause = registration_proc->amf_cause;
  amf_sap.u.amf_as.u.establish.nas_info  = AMF_AS_NAS_INFO_REGISTERED;

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
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGISTRATION_REJECT message\n");
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
    OAILOG_WARNING(
        LOG_AMF_APP, " Registration_proc null, from %s\n", __FUNCTION__);
  }
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
      if (amf_context->is_initial_identity_imsi == true) {
        if (registration_proc->ies->decode_status.mac_matched == 0) {
          /* IMSI is known but mac-mismatch start the authentication process */
          amf_ctx_clear_auth_vectors(amf_context);

          rc = amf_start_registration_proc_authentication(
              amf_context, registration_proc);
        } else {
          /* IMSI is known and Mac is matching */
          amf_registration(amf_context);
        }
      } else {
        /* If its first time GUTI Identify the IMSI */
        rc = amf_proc_identification(
            amf_context, (nas_amf_proc_t*) registration_proc,
            IDENTITY_TYPE_2_IMSI, amf_registration_success_identification_cb,
            amf_registration_failure_identification_cb);
      }
    } else {
      OAILOG_ERROR(LOG_NAS_AMF, "Unsupported Identifier type! \n");
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
    registration_proc->amf_cause        = amf_context->amf_cause;
    amf_sap_t amf_sap                   = {};
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
 ** Description: Performs the registration signaling                      **
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
    registration_proc->T3550.id             = -1;
    registration_proc->retransmission_count = 0;
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
    amf_sap_t amf_sap = {};
    nas_amf_registration_proc_t* registration_proc =
        get_nas_specific_procedure_registration(amf_context);
    ue_m5gmm_context_s* ue_m5gmm_context_p =
        PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context);
    amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;

    if (registration_proc) {
      registration_proc->T3550.id = NAS5G_TIMER_INACTIVE_ID;

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

      m5gmm_state_t state =
          PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)->mm_state;

      /*
       * Notify AMF-AS SAP that Registaration Accept message
       * if this is a re-transmit or periodic registration
       */
      if ((registration_proc->registration_accept_sent) ||
          ((registration_proc->ies->m5gsregistrationtype ==
            AMF_REGISTRATION_TYPE_PERIODIC_UPDATING) &&
           (state == REGISTERED_CONNECTED))) {
        amf_sap.primitive                = AMFAS_DATA_REQ;
        amf_sap.u.amf_as.u.data.ue_id    = ue_id;
        amf_sap.u.amf_as.u.data.nas_info = AMF_AS_NAS_DATA_REGISTRATION_ACCEPT;
        amf_sap.u.amf_as.u.data.guti     = new (guti_m5_t)();
        *(amf_sap.u.amf_as.u.data.guti)  = amf_context->m5_guti;
      } else {
        /*
         * Notify AMF-AS SAP that Registaration Accept message together
         * with an Activate Pdu session Context Request message has to
         * be sent to the UE
         */
        amf_sap.primitive                     = AMFAS_ESTABLISH_CNF;
        amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
        amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_REGISTERED;

        /* GUTI have already updated in amf_context during Identification
         * response complete, now assign to amf_sap
         */
        amf_sap.u.amf_as.u.establish.guti = amf_context->m5_guti;
      }

      rc = amf_sap_send(&amf_sap);
      if (rc == RETURNok) {
        registration_proc->registration_accept_sent++;
      }
      /*
       * Start T3550 timer
       */
      registration_proc->T3550.id = amf_app_start_timer(
          REGISTRATION_ACCEPT_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE,
          registration_accept_t3550_handler, registration_proc->ue_id);
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "Timer: Registration_accept timer T3550 with id  %lu "
          "Started for ue id: " AMF_UE_NGAP_ID_FMT,
          registration_proc->T3550.id, registration_proc->ue_id);
    }

    // s6a update location request
    int rc =
        amf_send_n11_update_location_req(ue_m5gmm_context_p->amf_ue_ngap_id);
    if (rc == RETURNerror) {
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: n11_update_location_req failure\n");
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

static int registration_accept_t3550_handler(
    zloop_t* loop, int timer_id, void* arg) {
  amf_context_t* amf_ctx                         = NULL;
  ue_m5gmm_context_s* ue_amf_context             = NULL;
  nas_amf_registration_proc_t* registration_proc = NULL;
  amf_ue_ngap_id_t ue_id                         = 0;
  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3550: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }
  /*
   * Get the UE context
   */
  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context == NULL) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "ue context not found for the ue_id=" AMF_UE_NGAP_ID_FMT "\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ctx = &ue_amf_context->amf_context;

  registration_proc =
      (nas_amf_registration_proc_t*)
          ue_amf_context->amf_context.amf_procedures->amf_specific_proc;

  if (registration_proc) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "T3550: timer id: %lu expired for"
        "ue_id=" AMF_UE_NGAP_ID_FMT "\n",
        registration_proc->T3550.id, registration_proc->ue_id);

    registration_proc->retransmission_count += 1;
    if (registration_proc->retransmission_count < REGISTRATION_COUNTER_MAX) {
      /* Send entity Registration accept message to the UE */

      OAILOG_WARNING(
          LOG_AMF_APP,
          "T3550: timer has expired retransmitting registration accept\n");
      amf_send_registration_accept(amf_ctx);
    } else {
      /* Abort the registration procedure */
      OAILOG_ERROR(
          LOG_AMF_APP,
          "T3550: Maximum retires:%d, for registration accept done hence Abort "
          "the registration "
          "procedure\n",
          registration_proc->retransmission_count);
      // To abort the registration procedure
      amf_proc_registration_abort(amf_ctx, ue_amf_context);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_registration_complete()                              **
 **                                                                        **
 ** Description: Completion of Registration Procedure                      **
 **                                                                        **
 ** Inputs:  amf_ctx:     UE Related AMF context                           **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_proc_registration_complete(amf_context_t* amf_ctx) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  nas_amf_registration_proc_t* registration_proc = NULL;
  int rc                                         = RETURNerror;
  amf_sap_t amf_sap                              = {};
  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_ctx, struct ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;

  if (amf_ctx) {
    if (is_nas_specific_procedure_registration_running(amf_ctx)) {
      registration_proc = (nas_amf_registration_proc_t*)
                              amf_ctx->amf_procedures->amf_specific_proc;

      amf_app_stop_timer(registration_proc->T3550.id);
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "Timer: after stop registration timer T3550 with id = %lu\n",
          registration_proc->T3550.id);
      registration_proc->T3550.id = NAS5G_TIMER_INACTIVE_ID;

      /*
       * Upon receiving an REGISTRATION COMPLETE message, the AMF shall enter
       * state AMF-REGISTERED and consider the GUTI sent in the REGISTRATION
       * ACCEPT message as valid.
       */
      amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_GUTI);
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_AMF,
        "UE Context not found for "
        "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  /*
   * Set the network registration indicator
   */
  amf_ctx->is_registered = true;

  /*
   * Notify AMF that registration procedure has successfully completed
   */
  amf_sap.primitive                   = AMFREG_REGISTRATION_CNF;
  amf_sap.u.amf_reg.ue_id             = ue_id;
  amf_sap.u.amf_reg.ctx               = amf_ctx;
  amf_sap.u.amf_reg.notify            = true;
  amf_sap.u.amf_reg.free_proc         = true;
  amf_sap.u.amf_reg.u.registered.proc = registration_proc;
  rc                                  = amf_sap_send(&amf_sap);
  if (rc == RETURNok) {
    /*
     * Send AMF Information after handling Registration Complete message
     * TODO this logic will handled in future when PDU Session Establish
     * resquest comes along with Initial Registration request.
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
  int rc                               = RETURNerror;
  ue_m5gmm_context_s* ue_m5gmm_context = NULL;

  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "AMFAS-SAP - received registration complete message for ue_id "
      "=" AMF_UE_NGAP_ID_FMT "\n",
      ue_id);

  ue_m5gmm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_m5gmm_context == NULL) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  /* FSM to move from COMMON_PROCEDURE_INITIATED2 -> REGISTERED_CONNECTED */
  ue_state_handle_message_initial(
      COMMON_PROCEDURE_INITIATED2, STATE_EVENT_REG_COMPLETE, SESSION_NULL,
      ue_m5gmm_context, &ue_m5gmm_context->amf_context);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
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
  int rc                 = RETURNerror;
  amf_sap_t amf_sap      = {};
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
  // PDU session establishment with initial registration
  amf_primitive_t primitive          = msg->primitive;
  amf_reg_t* evt                     = &msg->u.amf_reg;
  amf_context_t* amf_ctx             = msg->u.amf_reg.ctx;
  ue_m5gmm_context_s* ue_amf_context = NULL;

  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(evt->ue_id);

  if (!ue_amf_context) {
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "Ue context not found for the ue id" AMF_UE_NGAP_ID_FMT "\n",
        evt->ue_id);
    rc = RETURNerror;
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  if (evt && amf_ctx) {
    switch (primitive) {
      case AMFREG_REGISTRATION_CNF: {
        if (evt->free_proc) {
          amf_delete_registration_proc(amf_ctx);
        }

        /* Update the state */
        ue_amf_context->mm_state = REGISTERED_CONNECTED;
        OAILOG_DEBUG(
            LOG_NAS_AMF, "UE current state is %u\n", ue_amf_context->mm_state);
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

    amf_delete_child_procedures(amf_ctx, (nas5g_base_proc_t*) proc);
    delete_wrapper(&proc);

    amf_ctx->amf_procedures->amf_specific_proc = nullptr;

    nas_amf_procedure_gc(amf_ctx);
  }
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
    delete_wrapper(&(*ies)->imsi);
  }

  if ((*ies)->guti) {
    delete_wrapper(&(*ies)->guti);
  }

  if ((*ies)->imei) {
    delete_wrapper(&(*ies)->imei);
  }

  if ((*ies)->drx_parameter) {
    delete_wrapper(&(*ies)->drx_parameter);
  }

  if ((*ies)->last_visited_registered_tai) {
    delete_wrapper(&(*ies)->last_visited_registered_tai);
  }

  delete_wrapper(ies);
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

static void delete_common_proc_by_type(nas_amf_common_proc_t* proc) {
  if (proc) {
    switch (proc->type) {
      case AMF_COMM_PROC_AUTH: {
        delete (reinterpret_cast<nas5g_amf_auth_proc_t*>(proc));
      } break;
      case AMF_COMM_PROC_SMC: {
        delete (reinterpret_cast<nas_amf_smc_proc_t*>(proc));
      } break;
      case AMF_COMM_PROC_IDENT: {
        delete (reinterpret_cast<nas_amf_ident_proc_t*>(proc));
      } break;
      default: {}
    }
  }
}

/***********************************************************************
 ** Name:    amf_delete_common_procedure()                            **
 **                                                                   **
 ** Description: deletes the nas registration specific common         **
 **              procedures                                           **
 **                                                                   **
 ** Inputs:  proc: nas amf common proc                                **
 **                                                                   **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/
void amf_delete_common_procedure(
    amf_context_t* amf_ctx, nas_amf_common_proc_t** proc) {
  if (proc && *proc) {
    switch ((*proc)->type) {
      case AMF_COMM_PROC_AUTH: {
      } break;
      case AMF_COMM_PROC_SMC: {
      } break;
      case AMF_COMM_PROC_IDENT: {
      } break;
      default: {}
    }
  }

  // remove proc from list
  if (amf_ctx->amf_procedures) {
    nas_amf_common_procedure_t* p1 =
        LIST_FIRST(&amf_ctx->amf_procedures->amf_common_procs);
    nas_amf_common_procedure_t* p2 = NULL;

    // 2 methods: this one, the other: use parent struct macro and LIST_REMOVE
    // without searching matching element in the list
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      if (p1->proc == (nas_amf_common_proc_t*) (*proc)) {
        LIST_REMOVE(p1, entries);
        delete_common_proc_by_type(p1->proc);
        delete (p1);
        return;
      }
      p1 = p2;
    }
    nas_amf_procedure_gc(amf_ctx);
  }

  return;
}
/****************************************************************************
**                                                                        **
** Name:    amf_proc_registration_abort()                                 **
**                                                                        **
** Description: Abort the ongoing registration procedure                  **
**              for timer failure cases                                   **
**                                                                        **
** Inputs:  amf_ctx  : AMF context                                        **
**          ue_amf_context: UE context                                    **
**                                                                        **
**      Others:    None                                                   **
**                                                                        **
** Outputs:                                                               **
**      Return:    RETURNok, RETURNerror                                  **
**      Others:    None                                                   **
**                                                                        **
***************************************************************************/
int amf_proc_registration_abort(
    amf_context_t* amf_ctx, struct ue_m5gmm_context_s* ue_amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int rc = RETURNerror;
  if ((ue_amf_context) && (amf_ctx)) {
    MessageDef* message_p = nullptr;
    message_p =
        itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);

    if (message_p == NULL) {
      OAILOG_ERROR(LOG_AMF_APP, "message is NULL");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
    }
    memset(
        &NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p), 0,
        sizeof(itti_ngap_ue_context_release_command_t));

    NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id =
        ue_amf_context->amf_ue_ngap_id;
    NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id =
        ue_amf_context->gnb_ue_ngap_id;
    NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause =
        (Ngcause) ngap_CauseNas_deregister;
    message_p->ittiMsgHeader.imsi = ue_amf_context->amf_context.imsi64;
    send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
    amf_delete_registration_proc(amf_ctx);
    amf_free_ue_context(ue_amf_context);
    rc = RETURNok;
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}
/***************************************************************************
**                                                                        **
** Name:    get_decrypt_imsi_suci_extension()                             **
**                                                                        **
** Description: Invokes .get_decrypt_imsi_info                            **
**              to fetch decrypted imsi                                   **
**                                                                        **
**                                                                        **
***************************************************************************/
int get_decrypt_imsi_suci_extension(
    amf_context_t* amf_context, uint8_t ue_pubkey_identifier,
    const std::string& ue_pubkey, const std::string& ciphertext,
    const std::string& mac_tag) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[subscriberdb] for ue: [" AMF_UE_NGAP_ID_FMT
      "] decrypt-imsi\n",
      ue_id);

  AMFClientServicer::getInstance().get_decrypt_imsi_info(
      ue_pubkey_identifier, ue_pubkey, ciphertext, mac_tag, ue_id);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

}  // namespace magma5g
