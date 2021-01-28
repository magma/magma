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

  Source      Registration.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "3gpp_24.501.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif
#include "amf_data.h"
#include "amf_fsm.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_asDefs.h"
#include "amf_as.h"
#include "amf_sap.h"
#include "M5GSRegistrationResult.h"
#include "amf_recv.h"
#include "nas5g_network.h"

using namespace std;
#define M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH                                \
  1  // TODO temporary should be in nas5g
#define INVALID_IMSI64 (imsi64_t) 0
#define INVALID_AMF_UE_NGAP_ID 0x0

namespace magma5g {
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily store
                              // context inserted to ht
amf_sap_c amf_sap_reg;
nas_proc nas_procedure_reg;
amf_as_data_t amf_data_sec;
m5g_authentication m5g_auth;
identification amf_proc;
nas_network nas_nw;
nas_amf_smc_proc_t smc_proc;
static int amf_registration_failure_authentication_cb(
    amf_context_t* amf_context);
static int amf_registration_failure_identification_cb(
    amf_context_t* amf_context);
static int amf_start_registration_proc_security(
    amf_context_t* amf_context, nas_amf_registration_proc_t* registration_proc);
static int amf_registration(amf_context_t* amf_context);
static int amf_send_registration_accept(amf_context_t* amf_context);
static int amf_registration_failure_security_cb(amf_context_t* amf_context);

//------------------------------------------------------------------------------
// int amf_registration_procedure::amf_registration_success_authentication_cb( {
int amf_registration_procedure::amf_registration_success_authentication_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST Authentication procedure success and start  Security mode"
      "command procedures\n");
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_3__1);
    rc = amf_start_registration_proc_security(amf_context, registration_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
static int amf_start_registration_proc_authentication(
    amf_context_t* amf_context,
    nas_amf_registration_proc_t* registration_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  OAILOG_INFO(LOG_AMF_APP, "AMF-TEST: , from %s\n", __FUNCTION__);
  if ((amf_context) && (registration_proc)) {
    OAILOG_INFO(
        LOG_AMF_APP, "AMF-TEST: , calling amf_proc_authentication() from %s\n",
        __FUNCTION__);
    rc = m5g_auth.amf_proc_authentication(
        amf_context, &registration_proc->amf_spec_proc,
        amf_registration_procedure::amf_registration_success_authentication_cb,
        amf_registration_failure_authentication_cb);
  }
  // amf_registration_success_authentication_cb(amf_context);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

nas_amf_registration_proc_t* nas_new_registration_procedure(
    amf_context_t* amf_context) {
  if (!(amf_context->amf_procedures)) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: From nas_new_registration_procedure allocating for "
        "amf_procedures\n");
    amf_context->amf_procedures = _nas_new_amf_procedures(amf_context);
  }
  //  amf_context->amf_procedures->amf_specific_proc =
  //      new nas_amf_registration_proc_t ;// TODO AMF_TEST
  amf_context->amf_procedures->amf_specific_proc = new nas_amf_specific_proc_t;
  // amf_context->amf_procedures->amf_specific_proc = new
  // nas_amf_specific_proc_t;
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

  // TODO timer to be implemented later
  // proc->T3450.sec = amf_config.nas_config.t3450_sec;
  // proc->T3450.id  = NAS_TIMER_INACTIVE_ID;
  ue_m5gmm_global_context.amf_context.amf_procedures =
      amf_context->amf_procedures;  // TODO AMF_TEST global var to temporarily
                                    // store context inserted to ht
  OAILOG_TRACE(
      LOG_NAS_AMF, "New AMF_SPEC_PROC_TYPE_REGISTRATION initialized\n");
  return proc;
}

void amf_proc_create_procedure_registration_request(
    ue_m5gmm_context_s* ue_ctx, amf_registration_request_ies_t* ies) {
  nas_amf_registration_proc_t* reg_proc =
      nas_new_registration_procedure(&ue_ctx->amf_context);
  if ((reg_proc)) {
    reg_proc->ies   = ies;
    reg_proc->ue_id = ue_ctx->amf_ue_ngap_id;
// TODO callbacl to be implemeted later
#if 0
    ((nas_base_proc_t*) reg_proc)->abort    = NULL;
    ((nas_base_proc_t*) reg_proc)->fail_in  = NULL;  
    ((nas_base_proc_t*) reg_proc)->time_out = NULL;
    ((nas_base_proc_t*) reg_proc)->fail_out = NULL;
#endif
  }
}

int amf_registration_procedure::amf_proc_registration_request(
    amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    amf_registration_request_ies_t* ies) {
  int rc = RETURNerror;
  ue_m5gmm_context_s ue_ctx;
  amf_fsm_state_t fsm_state             = AMF_DEREGISTERED;
  bool clear_amf_ctxt                   = false;
  ue_m5gmm_context_s* ue_m5gmm_context  = NULL;
  ue_m5gmm_context_s* guti_ue_m5gmm_ctx = NULL;
  ue_m5gmm_context_s* imsi_ue_m5gmm_ctx = NULL;
  amf_context_t* new_amf_ctx            = NULL;
  imsi64_t imsi64                       = INVALID_IMSI64;
  amf_ue_ngap_id_t old_ue_id            = INVALID_AMF_UE_NGAP_ID;

  if (ies->imsi) {
    imsi64 = amf_imsi_to_imsi64(ies->imsi);
    OAILOG_INFO(
        LOG_AMF_APP,
        "During initial registration request "
        "SUPI as IMSI converted to imsi64 " IMSI_64_FMT " = ",
        imsi64);
  } else if (ies->guti) {
    // OAILOG_INFO(LOG_NAS_AMF,
    //"REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (GUTI = " GUTI_FMT ")
    //\n", ue_id,GUTI_ARG_M5G(ies->guti));
  } else if (ies->imei) {
    char imei_str[16];
    IMEI_TO_STRING(ies->imei, imei_str, 16);
    OAILOG_INFO(
        LOG_AMF_APP,
        "REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (IMEI = %s ) \n",
        ue_id, imei_str);
  }

  // OAILOG_DEBUG(
  //    LOG_NAS_AMF,
  //    "is_initial request = %u\n (ue_id=" AMF_UE_NGAP_ID_FMT
  //    ") \n(imsi = " IMSI_64_FMT ") \n",
  //    ies->is_initial, ue_id, imsi64);
  /*
   * Initialize the temporary UE context
   */
  // memset(&ue_ctx, 0, sizeof(ue_m5gmm_context_s)); TODO AMF_TEST, commented
  // due to crash after adding smf_context
  memset(&ue_ctx, 0, sizeof(ue_m5gmm_context_s));
  ue_ctx.amf_context.is_dynamic = false;
  ue_ctx.amf_ue_ngap_id         = ue_id;

  if (!(nas_procedure_reg.is_nas_specific_procedure_registration_running(
          &ue_ctx.amf_context))) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: From amf_proc_registration_request "
        "is_nas_specific_procedure_registration_running");
    amf_proc_create_procedure_registration_request(&ue_ctx, ies);
  }

  rc = amf_registration_procedure::amf_registration_run_procedure(
      &ue_ctx.amf_context);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
int amf_registration_procedure::amf_proc_registration_reject(
    amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  amf_context_t* amf_ctx;  //= amf_context_get(&amf_data, ue_id);//TODO -
                           // NEED-RECHECK need to declare the function
  nas_amf_registration_proc_t* registration_proc =
      (nas_amf_registration_proc_t*) (amf_ctx->amf_procedures
                                          ->amf_specific_proc);
  registration_proc->amf_cause = amf_cause;
  if (amf_ctx) {
    if (nas_procedure_reg.is_nas_specific_procedure_registration_running(
            amf_ctx)) {
      // TODO could be in callback of attach procedure triggered by
      // AMF REG__REJ
      rc = amf_registration_procedure::amf_registration_reject(
          amf_ctx, registration_proc);
      amf_sap_t amf_sap;  //               = {0};
      amf_sap.primitive                   = AMFREG_REGISTRATION_REJ;
      amf_sap.u.amf_reg.ue_id             = ue_id;
      amf_sap.u.amf_reg.ctx               = amf_ctx;
      amf_sap.u.amf_reg.notify            = false;
      amf_sap.u.amf_reg.free_proc         = true;
      amf_sap.u.amf_reg.u.registered.proc = registration_proc;
      rc                                  = amf_sap_reg.amf_sap_send(&amf_sap);
    } else {
      nas_amf_registration_proc_t no_registration_proc;  // = {0};
      no_registration_proc.ue_id       = ue_id;
      no_registration_proc.amf_cause   = amf_cause;
      no_registration_proc.amf_msg_out = {0};
      rc = amf_registration_procedure::amf_registration_reject(
          amf_ctx, registration_proc);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
int amf_registration_procedure::amf_registration_reject(
    amf_context_t* amf_context, nas_amf_registration_proc_t* nas_base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  amf_sap_t amf_sap;
  nas_amf_registration_proc_t* registration_proc =
      (nas_amf_registration_proc_t*) nas_base_proc;

  OAILOG_WARNING(
      LOG_AMF_APP, "AMF-PROC  - AMF Registration procedure not accepted ");
  //"by the network (ue_id=" AMF_UE_NGAP_ID_FMT ", cause=%d)\n",
  // registration_proc->ue_id, registration_proc->amf_cause);
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
    // OAILOG_ERROR(LOG_NAS_EMM, "AMF-PROC  - SMF message is missing\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  /*
   * Setup 5G CN NAS security data
   */
  if (amf_context) {  // TODO -  NEED-RECHECK
    amf_data_sec.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.establish.sctx, &amf_context->_security, false,
        false);
  } else {
    amf_data_sec.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.establish.sctx, NULL, false, false);
  }
  rc = amf_sap_reg.amf_sap_send(&amf_sap);
  increment_counter(
      "ue_Registration", 1, 1, "action", "Registration_reject_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
/*
 * --------------------------------------------------------------------------
 * Functions that may initiate AMF common procedures
 * --------------------------------------------------------------------------
 */

//------------------------------------------------------------------------------
int amf_registration_procedure::amf_registration_run_procedure(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);
#if 1
  if (registration_proc == NULL) {
    OAILOG_INFO(
        LOG_AMF_APP, "AMF_TEST: registration_proc NULL, from %s\n",
        __FUNCTION__);
  }
#endif  // AMF_TEST
  if (registration_proc) {
    if (registration_proc->ies->last_visited_registered_tai) {
      // amf_ctx_set_valid_lvr_tai(amf_context,
      // registration_proc->ies->last_visited_registered_tai);
      // amf_ctx_set_valid_ue_nw_cap(amf_context,
      // &registration_proc->ies->ue_network_capability);
    }

    // if (registration_proc->ies->ms_network_capability)
    // {
    // amf_ctx_set_valid_ms_nw_cap(amf_context,
    // registration_proc->ies->ms_network_capability);
    // }
    // amf_context->originating_tai = *registration_proc->ies->originating_tai;

    // temporary choice to clear security context if it exist
    // amf_ctx_clear_security(amf_context);

    if (registration_proc->ies->imsi) {
      if ((registration_proc->ies->decode_status.mac_matched) ||
          //!(registration_proc->ies->decode_status //TODO AMF_TEST value coming
          // up as 0
          (registration_proc->ies->decode_status.integrity_protected_message)) {
        // force authentication, even if not necessary
        imsi64_t imsi64 = amf_imsi_to_imsi64(registration_proc->ies->imsi);
        amf_ctx_set_valid_imsi(
            amf_context, registration_proc->ies->imsi, imsi64);
        // TODO amf_context_upsert_imsi(&_amf_data, amf_context);
        rc = amf_start_registration_proc_authentication(
            amf_context, registration_proc);
        if (rc != RETURNok) {
          OAILOG_ERROR(
              LOG_NAS_AMF,
              "Failed to start registration authentication procedure! \n");
        }
      } else {
        // force identification, even if not necessary
        rc = amf_proc.amf_proc_identification(
            amf_context, (nas_amf_proc_t*) registration_proc,
            IDENTITY_TYPE_2_IMSI,
            amf_registration_procedure::
                amf_registration_success_identification_cb,
            amf_registration_failure_identification_cb);
      }
    } else if (registration_proc->ies->guti) {
      /* TODO: Check the amf_supi_guti_map for presence of respective guti
       * especially compare with 4 octet of 5G-TMSI of GUTI from UE
       * if matches, skip identification request and send authentication
       * request. At the same time update GUTI with other random TMSI value
       */
      rc = amf_proc.amf_proc_identification(
          amf_context, (nas_amf_proc_t*) registration_proc,
          IDENTITY_TYPE_2_IMSI,
          amf_registration_procedure::
              amf_registration_success_identification_cb,
          amf_registration_failure_identification_cb);
    } else if (registration_proc->ies->imei) {
      // emergency allowed if go here, but have to be implemented...
      // AssertFatal(0, "TODO emergency");//TODO -  NEED-RECHECK
    }
  }
  // amf_registration_procedure::amf_registration_success_identification_cb(
  //    amf_context);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
int amf_registration_procedure::amf_registration_success_identification_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: Identification procedure success\n");
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF-TEST: ,registration_proc valid, start authentication, from %s\n",
        __FUNCTION__);
    rc = amf_start_registration_proc_authentication(
        amf_context, registration_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
static int amf_registration_failure_identification_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  // OAILOG_ERROR(LOG_NAS_AMF, "registration - Identification procedure
  // failed!\n");

  // AssertFatal(0, "Cannot happen...\n");//TODO -  NEED-RECHECK
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
static int amf_registration_failure_authentication_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  // OAILOG_ERROR(LOG_NAS_AMF, "REGISTRATION - Authentication procedure
  // failed!\n");
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    registration_proc->amf_cause = amf_context->amf_cause;

    amf_sap_t amf_sap;
    amf_sap.primitive                   = AMFREG_REGISTRATION_REJ;
    amf_sap.u.amf_reg.ue_id             = registration_proc->ue_id;
    amf_sap.u.amf_reg.ctx               = amf_context;
    amf_sap.u.amf_reg.notify            = true;
    amf_sap.u.amf_reg.free_proc         = true;
    amf_sap.u.amf_reg.u.registered.proc = registration_proc;
    // dont' care amf_sap.u.amf_reg.u.registration.is_emergency = false;
    rc = amf_sap_reg.amf_sap_send(&amf_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
int amf_registration_procedure::amf_registration_success_security_cb(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  // OAILOG_INFO(LOG_NAS_AMF, "REGISTRATION - Security procedure success!\n");
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    rc = amf_registration(amf_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
static int amf_start_registration_proc_security(
    amf_context_t* amf_context,
    nas_amf_registration_proc_t* registration_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  if ((amf_context) && (registration_proc)) {
    // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_3__1);
    amf_ue_ngap_id_t ue_id =
        PARENT_STRUCT(amf_context, struct ue_m5gmm_context_s, amf_context)
            ->amf_ue_ngap_id;
    /*
     * Create new NAS security context
     */
    smc_proc.amf_ctx_clear_security(amf_context);
    rc = amf_proc_security_mode_control(
        amf_context, &registration_proc->amf_spec_proc, registration_proc->ksi,
        amf_registration_procedure::amf_registration_success_security_cb,
        // amf_registration_success_security_cb,
        amf_registration_failure_security_cb);
    // amf_registration_success_security_cb(amf_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//------------------------------------------------------------------------------
static int amf_registration_failure_security_cb(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  // OAILOG_ERROR(LOG_NAS_AMF, "REGISTRATION - Security procedure failed!\n");
  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
    //_amf_registration_release(amf_context);//TODO -  NEED-RECHECK
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/*
 *
 * Name:        amf_registration_security()
 *
 * Description: Initiates security mode control AMF common procedure.
 *
 * Inputs:          args:      security argument parameters
 *                  Others:    None
 *
 * Outputs:     None
 *                  Return:    RETURNok, RETURNerror
 *                  Others:    _amf_data
 */
//------------------------------------------------------------------------------
int m5g_authentication::amf_registration_security(amf_context_t* amf_context) {
  // return m5g_auth.amf_registration_security(amf_context);
  return RETURNok;
}

/*
--------------------------------------------------------------------------
                AMF specific local functions
--------------------------------------------------------------------------
*/

/*
 *
 * Name:    amf_registration()
 *
 * Description: Performs the registration signalling procedure while a context
 *      exists for the incoming UE in the network.
 *
 *              3GPP TS 24.501, section 5.5.1.2.4
 *      Upon receiving the REGISTRATION REQUEST message, the AMF shall
 *      send an REGISTRATION ACCEPT message to the UE and start timer
 *      T3450.
 *
 * Inputs:  args:      registration argument parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    _amf_data
 *
 */
//------------------------------------------------------------------------------
int amf_registration_procedure::amf_registration(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_context, struct ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: "
      "ue_id=" AMF_UE_NGAP_ID_FMT
      "Start REGISTRATION_ACCEPT procedures for UE \n",
      ue_id);

  nas_amf_registration_proc_t* registration_proc =
      nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);

  if (registration_proc) {
#if 0  // TODO -  NEED-RECHECK
            if (registration_proc->ies->smf_msg) 
            {
                smf_sap_t smf_sap;
                smf_sap.primitive = SMF_UNITDATA_IND;
                smf_sap.is_standalone = false;
                smf_sap.ue_id = ue_id;
                smf_sap.ctx = amf_context;
                smf_sap.recv = registration_proc->ies->smf_msg;
                //rc = smf_sap_send(&smf_sap);
                if ((rc != RETURNerror) && (smf_sap.err == SMF_SAP_SUCCESS)) 
                {
                    rc = RETURNok;
                } 
                /*else if (smf_sap.err != SMF_SAP_DISCARDED) 
                {
                    
                    * Theregistration procedure failed due to an SMF procedure failure
                    */
                    registration_proc->amf_cause = AMF_CAUSE_SMF_FAILURE;

                    /*
                    * Setup the SMF message container to include pdu session Connectivity Reject
                    * message within the REGISTRATION Reject message
                    */
                   /* bdestroy_wrapper(&registration_proc->ies->smf_msg);
                    registration_proc->smf_msg_out = smf_sap.send;
                    OAILOG_ERROR(LOG_NAS_AMF, "Sending Registration Reject to UE ue_id = (%u), amf_cause = (%d)\n", ue_id, registration_proc->amf_cause);
                    rc = _amf_registration_reject(amf_context, &registration_proc->amf_spec_proc.amf_proc.base_proc);
                }*/
                else 
                {
                    /*
                    * SMF procedure failed and, received message has been discarded or
                    * Status message has been returned; ignore SMF procedure failure
                    */
                    OAILOG_WARNING(LOG_NAS_AMF,"Ignore SMF procedure failure &""received message has been discarded for ue_id = (%u)\n", ue_id);
                    rc = RETURNok;
                }
            } 
            else
#endif
    //{
    rc = RETURNok;
    rc = amf_registration_procedure::amf_send_registration_accept(amf_context);
    //}
  }

  if (rc != RETURNok) {
    /*
     * The registration procedure failed
     */
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "ue_id=" AMF_UE_NGAP_ID_FMT
        " AMF-PROC  - Failed to respond to Registration Request\n",
        ue_id);
    registration_proc->amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
    /*
     * Do not accept the UE to registration to the network
     */
    // OAILOG_ERROR(LOG_NAS_AMF,"Sending Registration Reject to UE ue_id = (%u),
    // amf_cause = (%d)\n", ue_id, registration_proc->amf_cause);
    // rc = _amf_registration_reject(amf_context,
    // &registration_proc->amf_spec_proc.amf_proc.base_proc);
    // increment_counter("ue_registration", 1, 2, "result", "failure", "cause",
    // "protocol_error");
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_send_registration_accept() **
 **                                                                        **
 ** Description: Sends REGISTRATION ACCEPT message and start timer T3450 **
 **                                                                        **
 ** Inputs:  data:      Registration accept retransmission data          **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3450                                      **
 **                                                                        **
 ***************************************************************************/
int amf_registration_procedure::amf_send_registration_accept(
    amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  // may be caused by timer not stopped when deleted context
  if (amf_context) {
    amf_sap_t amf_sap;
    nas_amf_registration_proc_t* registration_proc =
        nas_procedure_reg.get_nas_specific_procedure_registration(amf_context);
    ue_m5gmm_context_s* ue_m5gmm_context_p =
        PARENT_STRUCT(amf_context, class ue_m5gmm_context_s, amf_context);
    amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;

    if (registration_proc) {
      // TODO - NEED-RECHECK
      //_amf_registration_update(amf_context, registration_proc->ies);
      /*
       * The IMSI if provided by the UE
       */
      if (registration_proc->ies->imsi) {
        imsi64_t new_imsi64 = amf_imsi_to_imsi64(registration_proc->ies->imsi);
        if (new_imsi64 != amf_context->_imsi64) {
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
      amf_sap.u.amf_as.u.establish.guti = amf_context->_m5_guti;
      OAILOG_INFO(
          LOG_NAS_AMF,
          " AMF_TEST in %s assigned GUTI to amf_sap from amf_context\n",
          __FUNCTION__);
      OAILOG_INFO(
          LOG_NAS_AMF, " AMF_TEST Value of TMSI of GUTI %08" PRIx32 "\n",
          amf_sap.u.amf_as.u.establish.guti.m_tmsi);

      // NO_REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__3);
      // bdestroy_wrapper(&ue_m5gmm_context_s->ue_radio_capability);
      //----------------------------------------
      // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__4);
      // amf_ctx_set_attribute_valid(amf_context,
      // AMF_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
      //----------------------------------------
      if (registration_proc->ies->drx_parameter) {
        // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__5);
        // TODO - NEED-RECHECK
        // amf_ctx_set_valid_drx_parameter(amf_context,
        // registration_proc->ies->drx_parameter);
      }
      //----------------------------------------
      // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__9);
      // the set of amf_sap.u.amf_as.u.establish.new_guti is for including the
      // GUTI in the registration accept message
      // ONLY ONE MME NOW NO S10
      /*if (!IS_AMF_CTXT_PRESENT_GUTI(amf_context))
      {
          // Sure it is an unknown GUTI in this AMF
          guti_t old_guti = amf_context->_old_guti;
          guti_t guti = {.gummei.plmn = {0},
                      .gummei.mme_gid = 0,
                      .gummei.mme_code = 0,
                      .m_tmsi = INVALID_M_TMSI};
          clear_guti(&guti);

          rc = amf_api_new_guti(
          &amf_context->_imsi,
          &old_guti,
          &guti,
          &amf_context->originating_tai,
          &amf_context->_tai_list);
          if (RETURNok == rc) {
          amf_ctx_set_guti(amf_context, &guti);
          amf_ctx_set_attribute_valid(amf_context, AMF_CTXT_MEMBER_TAI_LIST);
          //----------------------------------------
          REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__6);
          REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__10);
          memcpy(
              &amf_sap.u.amf_as.u.establish.tai_list,
              &amf_context->_tai_list,
              sizeof(tai_list_t));
          }
          else
          {
              //OAILOG_ERROR(LOG_NAS_AMF,"Failed to assign amf api new guti for
      ue_id = %u\n",ue_id);
              //OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
          }
      }
      else
      {
          // Set the TAI attributes from the stored context for resends.
          memcpy(
          &amf_sap.u.amf_as.u.establish.tai_list,
          &amf_context->_tai_list,
          sizeof(tai_list_t));
      }*/
    }

    // amf_sap.u.amf_as.u.establish.eps_id.guti = &amf_context->_guti;

    /*if (!IS_AMF_CTXT_VALID_GUTI(amf_context) &&
    IS_AMF_CTXT_PRESENT_GUTI(amf_context) &&
    IS_AMF_CTXT_PRESENT_OLD_GUTI(amf_context))
    {
        /*
        * Implicit GUTI reallocation;
        * include the new assigned GUTI in the Registration Accept message
        */
    /*OAILOG_DEBUG(LOG_NAS_AMF, "ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  -
Implicit GUTI reallocation, include the new assigned " "GUTI in the Registration
Accept message\n", ue_id); amf_sap.u.amf_as.u.establish.new_guti = &amf->_guti;
} */
    /*else if (!IS_AMF_CTXT_VALID_GUTI(amf_context) &&
    IS_AMF_CTXT_PRESENT_GUTI(amf_context))
    {
        /*
        * include the new assigned GUTI in the Attach Accept message
        */
    /*OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - Include
the new assigned GUTI in the Registration Accept ""message\n", ue_id);
    amf_sap.u.amf_as.u.establish.new_guti = &amf_context->_guti;
} */
    // else
    {  // IS_AMF_CTXT_VALID_GUTI(ue_amf_context) is true
       // amf_sap.u.amf_as.u.establish.new_guti = NULL;
    }
    //----------------------------------------
    // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__14);
    // amf_sap.u.amf_as.u.establish.eps_network_feature_support =
    // &_amf_data.conf.eps_network_feature_support;

    /*
     * Delete any preexisting UE radio capabilities, pursuant to
     * GPP 24.5R15:5.5.1.2.4
     */
    // Note: this is safe from double-free errors because it sets to NULL
    // after freeing, which free treats as a no-op.
    // nas_nw.bdestroy_wrapper(&ue_m5gmm_context_p->ue_radio_capability);

    /*
    * Setup EPS NAS security data

    amf_as_set_security_data(&amf_sap.u.amf_as.u.establish.sctx,
    &amf_context->_security, false, true);
    amf_sap.u.amf_as.u.establish.encryption
    =amf_context->_security.selected_algorithms.encryption;
    amf_sap.u.amf_as.u.establish.integrity
    =amf_context->_security.selected_algorithms.integrity;
    OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  -
    encryption = 0x%X
    (0x%X)\n",ue_id,amf_sap.u.amf_as.u.establish.encryption,amf_context->_security.selected_algorithms.encryption);
    OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - integrity
    = 0x%X
    (0x%X)\n",ue_id,amf_sap.u.amf_as.u.establish.integrity,amf_context->_security.selected_algorithms.integrity);
*/
    /*
     * Get the activate default 5GMM PDu Session context request message to
     * transfer within the SMF container of the Registration accept message
     */
    amf_sap.u.amf_as.u.establish.nas_msg = registration_proc->amf_msg_out;
    OAILOG_TRACE(
        LOG_NAS_AMF,
        "ue_id=" AMF_UE_NGAP_ID_FMT
        " AMF-PROC  - nas_msg  src size = %d nas_msg  dst size = %d \n",
        ue_id, blength(registration_proc->amf_msg_out),
        blength(amf_sap.u.amf_as.u.establish.nas_msg));

    // Send T3402
    // amf_sap.u.amf_as.u.establish.t3402 = &amf_config.nas_config.t3402_min;

    // Encode CSFB parameters
    // _encode_csfb_parameters_attach_accept(amf_context,
    // &amf_sap.u.amf_as.u.establish);

    // REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__2);
    rc = amf_sap_reg.amf_sap_send(&amf_sap);

    if (RETURNerror != rc) {
      /*
       * Start T3450 timer
       */
      // nas_stop_T3450(registration_proc->ue_id, &registration_proc->T3450,
      // NULL);
      // nas_start_T3450(registration_proc->ue_id,&registration_proc->T3450,registration_proc->amf_spec_proc.amf_proc.base_proc.time_out,(void
      // *) amf_context);
    }
  } else {
    // OAILOG_WARNING(LOG_NAS_AMF, "ue_amf_context NULL\n");
  }
  // increment_counter("ue_registration", 1, 1, "action",
  // "registration_accept_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
/****************************************************************************
 **                                                                       **
 ** Name:    amf_send_registration_accept_dl_nas()                         **
 **                                                                        **
 ** Description: Builds Registration Accept message to be sent
 ** is NGAP : DL NAS Tx **
 **                                                                        **
 **      The registration Accept message is sent by the network to the     **
 **      UE to indicate that the corresponding attach request has          **
 **      been accepted.                                                    **
 **                                                                        **
 ** Inputs:  msg:       The AMFAS-SAP primitive to process                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_msg:   The AMF message to be sent                     **
 **      Return:    The size of the AMF message                            **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_registration_procedure::amf_send_registration_accept_dl_nas(
    const amf_as_data_t* msg, RegistrationAcceptMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_MAXIMUM_LENGTH;

  // Get the UE context
  amf_context_t* amf_ctx = amf_context_get(msg->ue_id);
  // DevAssert(amf_ctx);
  ue_m5gmm_context_s* ue_m5gmm_context_p =
      PARENT_STRUCT(amf_ctx, class ue_m5gmm_context_s, amf_context);
  amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;
  // DevAssert(msg->ue_id == ue_id);

  // OAILOG_INFO(LOG_AMF_APP, "AMFAS-SAP - Send Regisration Accept message\n");
  // OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP - size =
  // AMF_HEADER_MAXIMUM_LENGTH(%d)\n", size);
  /*
   * Mandatory - Message type
   */
  // amf_msg->messagetype = (uint8_t) REGISTRATION_ACCEPT;//TODO -  NEED-RECHECK
  /*
   * Mandatory - 5GS Registration result
   */
  size += M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH;  // TODO -  NEED-RECHECK
  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMFAS-SAP - size += AMF_REGISTRATION_RESULT_MAXIMUM_LENGTH(%d) (%d)\n",
      M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH, size);
  switch (amf_ctx->m5gsregistrationtype) {
    case AMF_REGISTRATION_TYPE_INITIAL:
      // amf_msg->m5gsregistrationresult = M5GS_REGISTRATION_RESULT_3GPP_ACCESS;
      // OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP -
      // M5GS_REGISTRATION_RESULT_3GPP_ACCESS\n");
      break;
    case AMF_REGISTRATION_TYPE_EMERGENCY:  // We should not reach here
      // OAILOG_ERROR(LOG_NAS_AMF,"AMFAS-SAP - M5GS emergency Registration,
      // currently unsupported\n"); OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);  //
      // TODO: fix once supported
      break;
  }
#if 0  // TODO -  NEED-RECHECK
     /*
    * Optional - Mobile Identity
    */
    if (msg->m5gsmobileidentity) {
        size += M5GS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
        amf_msg->presencemask |= REGISTRATION_ACCEPT_UE_IDENTITY_PRESENT;
        if (msg->msidentity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
        memcpy(
            &amf_msg->msidentity.imsi, &msg->ms_identity->imsi,
            sizeof(amf_msg->msidentity.imsi));
        } else if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
        memcpy(&amf_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
        sizeof(amf_msg->msidentity.tmsi));
        }
    }
#endif
  /*
  * Optional - Additional Update Result

  if (msg->additional_update_result) {
      size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
      amf_msg->presencemask |=
  REGISTRATION_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
      amf_msg->additionalupdateresult = SMS_ONLY;
  }
   */
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_registrationcomplete_response() **
 **                                                                        **
 ** Description: Processes registration Complete message **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      msg:       The received AMF message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                             **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int amf_procedure_handler::amf_handle_registrationcomplete_response(
    amf_ue_ngap_id_t ue_id, RegistrationCompleteMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc;

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMFAS-SAP - Received Registration Complete message for ue_id = (%u)\n",
      ue_id);
  /*
   * Execute the registration procedure completion
   */
  rc = amf_registration_procedure::amf_proc_registration_complete(
      ue_id, msg->smf_pdu, amf_cause, status);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

#if 1
//===========================================================================================================

int amf_registration_procedure::amf_proc_registration_complete(
    amf_ue_ngap_id_t ue_id, bstring smf_msg_pP, int amf_cause,
    const amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_amf_context             = NULL;
  nas_amf_registration_proc_t* registration_proc = NULL;
  int rc                                         = RETURNerror;
  amf_sap_t amf_sap;
  // smf_sap_t smf_sap ; //TODO -  NEED-RECHECK as PDU ses req comes in
  // different mesg
  amf_context_t* amf_ctx = NULL;

  /*
   * Get the UE context
   */
  //  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ue_amf_context =
      &ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily
                                 // store context inserted to ht

  if (ue_amf_context) {
    if (nas_procedure_reg.is_nas_specific_procedure_registration_running(
            &ue_amf_context->amf_context)) {
      registration_proc =
          (nas_amf_registration_proc_t*)
              ue_amf_context->amf_context.amf_procedures->amf_specific_proc;

      amf_ctx = &ue_amf_context->amf_context;
      /*
       * Upon receiving an REGISTRATION COMPLETE message, the AMF shall enter
       * state AMF-REGISTERED and consider the GUTI sent in the REGISTRATION
       * ACCEPT message as valid.
       */
      amf_ctx_set_attribute_valid(
          amf_ctx, AMF_CTXT_MEMBER_GUTI);  // TODO-RECHECK

/*currently by default Activate  Bearer Context Accept message was
 * sent in Registration complete Now, modified the code to send the message
 * received in Uplink/smfContainer.
 * third byte of smf message container is a message_type*/
#if 0
      switch (smf_msg_pP->data[2]) {
        case ACTIVATE_DEFAULT_PDU_SESSION_CONTEXT_ACCEPT:
          //smf_sap.primitive = SMF_DEFAULT_PDU_SESSION_CONTEXT_ACTIVATE_CNF;
          break;
        case ACTIVATE_DEFAULT_PDU_SESSION_CONTEXT_REJECT:
          //smf_sap.primitive = SMF_DEFAULT_PDU_SESSION_CONTEXT_ACTIVATE_REJ;
          break;
        default:
          OAILOG_ERROR(
              LOG_NAS_AMF, "Invalid SMF Message type, value = [%x] \n",
              smf_msg_pP->data[2]);
          break;
      }
      //smf_sap.is_standalone = false;
      //smf_sap.ue_id         = ue_id;
      //smf_sap.recv          = smf_msg_pP;
      //smf_sap.ctx           = &ue_amf_context->amf_context;
      rc                    = smf_sap_send(&smf_sap);
#endif
    } else {
      OAILOG_INFO(
          LOG_NAS_AMF,
          "UE " AMF_UE_NGAP_ID_FMT
          " REGISTRATION COMPLETE discarded (AMF procedure not found)\n",
          ue_id);
      bdestroy((bstring)(smf_msg_pP));
    }
  } else {
    OAILOG_WARNING(LOG_NAS_AMF, "UE Context not found..\n");
    OAILOG_INFO(
        LOG_NAS_AMF,
        "UE " AMF_UE_NGAP_ID_FMT
        " REGISTRATION COMPLETE discarded (context not found)\n",
        ue_id);
  }

  // if ((rc != RETURNerror) && (smf_sap.err == SMF_SAP_SUCCESS))
  rc = RETURNok;  // AMF_TEST
  if ((rc != RETURNerror)) {
    /*
     * Set the network registrationment indicator
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
    rc                                  = amf_sap_reg.amf_sap_send(&amf_sap);
    if (rc == RETURNok) {
      /*
       * Send AMF Information after handling Registration Complete message
       * */
      OAILOG_INFO(
          LOG_NAS_AMF, " Sending AMF INFORMATION for ue_id = (%u)\n", ue_id);
      amf_proc_amf_informtion(ue_amf_context);
      increment_counter(
          "ue_registration", 1, 1, "result", "registration_proc_successful");
      // registration_success_event(ue_amf_context->amf_context._imsi64);
    }
  }
#if 0
    else if (smf_sap.err != SMF_SAP_DISCARDED) {
    /*
     * Notify SMF that registration procedure failed
     */
    amf_sap.primitive               = AMFREG_REGISTRATION_REJ;
    amf_sap.u.amf_reg.ue_id         = ue_id;
    amf_sap.u.amf_reg.ctx           = &ue_amf_context->amf_context;
    amf_sap.u.amf_reg.notify        = true;
    amf_sap.u.amf_reg.free_proc     = true;
    amf_sap.u.amf_reg.u.registration.proc = registration_proc;
    rc                              = amf_sap_send(&amf_sap);
    }
#endif
  else {
    /*
     * SMF procedure failed and, received message has been discarded or
     * Status message has been returned; ignore SMF procedure failure
     */
    OAILOG_WARNING(
        LOG_NAS_AMF,
        "Ignore SMF procedure failure/received "
        "message has been discarded for"
        "ue_id = (%u)\n",
        ue_id);
    rc = RETURNok;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
//==========================================================
int amf_proc_amf_informtion(ue_m5gmm_context_s* ue_amf_ctx) {
  int rc = RETURNerror;
  amf_sap_t amf_sap;
  amf_sap_c amf_sap_reg;
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
#if 0
  unsigned char result[256] = {0};
  amf_as->daylight_saving_time = _amf_data.conf.daylight_saving_time;
  /*
   * Encode full_network_name with gsm 7 bit encoding
   * The encoding is done referring to 3gpp 24.008
   * (section: 10.5.3.5a)and 23.038
   */
  _amf_information_pack_gsm_7Bit(_amf_data.conf.full_network_name, result);
  amf_as->full_network_name = bfromcstr((const char*) result);
  /*
   * Encode short_network_name with gsm 7 bit encoding
   */
  memset(result, 0, sizeof(result));
  _amf_information_pack_gsm_7Bit(_amf_data.conf.short_network_name, result);
  amf_as->short_network_name = bfromcstr((const char*) result);
#endif  // TODO revisit later, as it is not part of demo scope

  /*
   * Setup EPS NAS security data
   */
  amf_as->amf_as_set_security_data(
      &amf_as->sctx, &amf_ctx->_security, false, true);
  /*
   * Notify AMF-AS SAP that TAU Accept message has to be sent to the network
   */
  amf_sap.primitive = AMFAS_DATA_REQ;
  rc                = amf_sap_reg.amf_sap_send(&amf_sap);

#if 0
  bdestroy(amf_as->full_network_name);
  bdestroy(amf_as->short_network_name);
#endif  // TODO revisit later, as it is not part of demo scope

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
int amf_reg_send(amf_reg_t* const msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;
#if 0
  /*
   * Check the AMF-SAP primitive
   */
  amf_reg_primitive_t primitive = msg->primitive;

  assert((primitive > _AMFREG_START) && (primitive < _AMFREG_END));
  /*
   * Execute the AMF procedure
   */
  rc = amf_fsm_process(msg);
#endif
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

//==========================================================

#endif
}  // namespace magma5g
