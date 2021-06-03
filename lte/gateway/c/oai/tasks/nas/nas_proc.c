/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************
  Source      nas_proc.c

  Version     0.1

  Date        2012/09/20

  Product     NAS stack

  Subsystem   NAS main process

  Author      Frederic Maurel, Lionel GAUTHIER

  Description NAS procedure call manager

*****************************************************************************/
#include <stdbool.h>
#include <stdint.h>
#include <string.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "nas_proc.h"
#include "emm_proc.h"
#include "emm_main.h"
#include "emm_sap.h"
#include "esm_main.h"
#include "s6a_defs.h"
#include "mme_app_ue_context.h"
#include "3gpp_24.008.h"
#include "3gpp_33.401.h"
#include "DetachRequest.h"
#include "MobileIdentity.h"
#include "common_types.h"
#include "emm_asDef.h"
#include "emm_data.h"
#include "hashtable.h"
#include "mme_api.h"
#include "mme_app_state.h"
#include "nas_procedures.h"
#include "service303.h"
#include "sgs_messages_types.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

static nas_cause_t s6a_error_2_nas_cause(uint32_t s6a_error, int experimental);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_initialize()                                     **
 **                                                                        **
 ** Description:                                                           **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void nas_proc_initialize(const mme_config_t* mme_config_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Initialize the EMM procedure manager
   */
  emm_main_initialize(mme_config_p);
  /*
   * Initialize the ESM procedure manager
   */
  esm_main_initialize();
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_cleanup()                                        **
 **                                                                        **
 ** Description: Performs clean up procedure before the system is shutdown **
 **                                                                        **
 ** Inputs:  None                                                      **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **          Return:    None                                       **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void nas_proc_cleanup(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Perform the EPS Mobility Manager's clean up procedure
   */
  emm_main_cleanup();
  /*
   * Perform the EPS Session Manager's clean up procedure
   */
  esm_main_cleanup();
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/*
   --------------------------------------------------------------------------
            NAS procedures triggered by the user
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_establish_ind()                                  **
 **                                                                        **
 ** Description: Processes the NAS signalling connection establishment     **
 **      indication message received from the network              **
 **                                                                        **
 ** Inputs:  ueid:      UE identifier                              **
 **      tac:       The code of the tracking area the initia-  **
 **             ting UE belongs to                         **
 **      data:      The initial NAS message transferred within  **
 **             the message                                **
 **      len:       The length of the initial NAS message      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_proc_establish_ind(
    const mme_ue_s1ap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi, const as_cause_t as_cause,
    const s_tmsi_t s_tmsi, STOLEN_REF bstring* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (msg) {
    emm_sap_t emm_sap = {0};

    /*
     * Notify the EMM procedure call manager that NAS signalling
     * connection establishment indication message has been received
     * from the Access-Stratum sublayer
     */

    emm_sap.primitive                          = EMMAS_ESTABLISH_REQ;
    emm_sap.u.emm_as.u.establish.ue_id         = ue_id;
    emm_sap.u.emm_as.u.establish.is_initial    = true;
    emm_sap.u.emm_as.u.establish.is_mm_ctx_new = is_mm_ctx_new;

    emm_sap.u.emm_as.u.establish.nas_msg = *msg;
    *msg                                 = NULL;
    emm_sap.u.emm_as.u.establish.tai     = &originating_tai;
    // emm_sap.u.emm_as.u.establish.plmn_id            = &originating_tai.plmn;
    // emm_sap.u.emm_as.u.establish.tac                = originating_tai.tac;
    emm_sap.u.emm_as.u.establish.ecgi = ecgi;

    rc = emm_sap_send(&emm_sap);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_dl_transfer_cnf()                                **
 **                                                                        **
 ** Description: Processes the downlink data transfer confirm message re-  **
 **      ceived from the network while NAS message has been succes-**
 **      sfully delivered to the NAS sublayer on the receiver side.**
 **                                                                        **
 ** Inputs:  ueid:      UE identifier                              **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_proc_dl_transfer_cnf(
    const uint32_t ue_id, const nas_error_code_t status,
    bstring* STOLEN_REF nas_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap_t emm_sap = {0};
  int rc            = RETURNok;

  /*
   * Notify the EMM procedure call manager that downlink NAS message
   * has been successfully delivered to the NAS sublayer on the
   * receiver side
   */
  emm_sap.primitive = EMMAS_DATA_IND;
  if (AS_SUCCESS == status) {
    emm_sap.u.emm_as.u.data.delivered = EMM_AS_DATA_DELIVERED_TRUE;
  } else {
    emm_sap.u.emm_as.u.data.delivered =
        EMM_AS_DATA_DELIVERED_LOWER_LAYER_FAILURE;
  }
  emm_sap.u.emm_as.u.data.ue_id = ue_id;
  if (*nas_msg) {
    emm_sap.u.emm_as.u.data.nas_msg = *nas_msg;
    *nas_msg                        = NULL;
  }
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_dl_transfer_rej()                                **
 **                                                                        **
 ** Description: Processes the downlink data transfer confirm message re-  **
 **      ceived from the network while NAS message has not been    **
 **      delivered to the NAS sublayer on the receiver side.       **
 **                                                                        **
 ** Inputs:  ueid:      UE identifier                              **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_proc_dl_transfer_rej(
    const uint32_t ue_id, const nas_error_code_t status,
    bstring* STOLEN_REF nas_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap_t emm_sap = {0};
  int rc            = RETURNok;

  /*
   * Notify the EMM procedure call manager that transmission
   * failure of downlink NAS message indication has been received
   * from lower layers
   */
  emm_sap.primitive             = EMMAS_DATA_IND;
  emm_sap.u.emm_as.u.data.ue_id = ue_id;
  if (AS_SUCCESS == status) {
    emm_sap.u.emm_as.u.data.delivered = EMM_AS_DATA_DELIVERED_TRUE;
  } else if (AS_NON_DELIVERED_DUE_HO == status) {
    emm_sap.u.emm_as.u.data.delivered =
        EMM_AS_DATA_DELIVERED_LOWER_LAYER_NON_DELIVERY_INDICATION_DUE_TO_HO;
  } else {
    emm_sap.u.emm_as.u.data.delivered =
        EMM_AS_DATA_DELIVERED_LOWER_LAYER_FAILURE;
  }
  emm_sap.u.emm_as.u.data.delivered = status;
  emm_sap.u.emm_as.u.data.nas_msg   = NULL;
  if (*nas_msg) {
    emm_sap.u.emm_as.u.data.nas_msg = *nas_msg;
    *nas_msg                        = NULL;
  }
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_ul_transfer_ind()                                **
 **                                                                        **
 ** Description: Processes uplink data transfer indication message recei-  **
 **      ved from the network                                      **
 **                                                                        **
 ** Inputs:  ueid:      UE identifier                              **
 **      data:      The transferred NAS message                 **
 **      len:       The length of the NAS message              **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int nas_proc_ul_transfer_ind(
    const mme_ue_s1ap_id_t ue_id, const tai_t originating_tai, const ecgi_t cgi,
    STOLEN_REF bstring* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  OAILOG_INFO(
      LOG_NAS,
      "Received NAS UPLINK DATA IND from S1AP for ue_id = " MME_UE_S1AP_ID_FMT
      "\n",
      ue_id);
  if (msg) {
    emm_sap_t emm_sap = {0};

    /*
     * Notify the EMM procedure call manager that data transfer
     * indication has been received from the Access-Stratum sublayer
     */
    emm_sap.primitive                 = EMMAS_DATA_IND;
    emm_sap.u.emm_as.u.data.ue_id     = ue_id;
    emm_sap.u.emm_as.u.data.delivered = true;
    emm_sap.u.emm_as.u.data.nas_msg   = *msg;
    *msg                              = NULL;
    emm_sap.u.emm_as.u.data.tai       = &originating_tai;
    // emm_sap.u.emm_as.u.data.plmn_id   = &originating_tai.plmn;
    // emm_sap.u.emm_as.u.data.tac       = originating_tai.tac;
    emm_sap.u.emm_as.u.data.ecgi = cgi;
    rc                           = emm_sap_send(&emm_sap);
  } else {
    OAILOG_WARNING(
        LOG_NAS, "Received NAS message in uplink is NULL for ue_id = (%u)\n",
        ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//-----------------------------------------------------------------------------
int nas_proc_authentication_info_answer(
    mme_app_desc_t* mme_app_desc_p, s6a_auth_info_ans_t* aia) {
  imsi64_t imsi64                  = INVALID_IMSI64;
  int rc                           = RETURNerror;
  emm_context_t* emm_ctxt_p        = NULL;
  ue_mm_context_t* ue_mm_context_p = NULL;
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  DevAssert(aia);
  IMSI_STRING_TO_IMSI64((char*) aia->imsi, &imsi64);

  OAILOG_DEBUG(LOG_NAS_EMM, "Handling imsi " IMSI_64_FMT "\n", imsi64);

  ue_mm_context_p = mme_ue_context_exists_imsi(
      &mme_app_desc_p->mme_ue_contexts, (const hash_key_t) imsi64);
  if (ue_mm_context_p) {
    emm_ctxt_p = &ue_mm_context_p->emm_context;
  }

  if (!(emm_ctxt_p)) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  mme_ue_s1ap_id_t mme_ue_s1ap_id = ue_mm_context_p->mme_ue_s1ap_id;
  OAILOG_INFO(
      LOG_NAS_EMM,
      "Received Authentication Information Answer from S6A for"
      " ue_id = " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
  if ((aia->result.present == S6A_RESULT_BASE) &&
      (aia->result.choice.base == DIAMETER_SUCCESS)) {
    /*
     * Check that list is not empty and contain at most MAX_EPS_AUTH_VECTORS
     * elements
     */
    DevCheck(
        aia->auth_info.nb_of_vectors <= MAX_EPS_AUTH_VECTORS,
        aia->auth_info.nb_of_vectors, MAX_EPS_AUTH_VECTORS, 0);
    DevCheck(
        aia->auth_info.nb_of_vectors > 0, aia->auth_info.nb_of_vectors, 1, 0);

    OAILOG_DEBUG(
        LOG_NAS_EMM, "INFORMING NAS ABOUT AUTH RESP SUCCESS got %u vector(s)\n",
        aia->auth_info.nb_of_vectors);
    rc = nas_proc_auth_param_res(
        mme_ue_s1ap_id, aia->auth_info.nb_of_vectors,
        aia->auth_info.eutran_vector);
  } else {
    OAILOG_ERROR(LOG_NAS_EMM, "INFORMING NAS ABOUT AUTH RESP ERROR CODE\n");
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause",
        "auth_info_failure_from_hss");
    /*
     * Inform NAS layer with the right failure
     */
    if (aia->result.present == S6A_RESULT_BASE) {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Auth info Rsp failure for imsi " IMSI_64_FMT
          ", base_error_code %d \n",
          imsi64, aia->result.choice.base);
      rc = nas_proc_auth_param_fail(
          mme_ue_s1ap_id, s6a_error_2_nas_cause(aia->result.choice.base, 0));
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Auth info Rsp failure for imsi " IMSI_64_FMT
          ", experimental_error_code %d \n",
          imsi64, aia->result.choice.experimental);
      rc = nas_proc_auth_param_fail(
          mme_ue_s1ap_id,
          s6a_error_2_nas_cause(aia->result.choice.experimental, 1));
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_auth_param_res(
    mme_ue_s1ap_id_t ue_id, uint8_t nb_vectors, eutran_vector_t* vectors) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                            = RETURNerror;
  emm_sap_t emm_sap                 = {0};
  emm_cn_auth_res_t emm_cn_auth_res = {0};

  emm_cn_auth_res.ue_id      = ue_id;
  emm_cn_auth_res.nb_vectors = nb_vectors;
  for (int i = 0; i < nb_vectors; i++) {
    emm_cn_auth_res.vector[i] = &vectors[i];
  }

  emm_sap.primitive           = EMMCN_AUTHENTICATION_PARAM_RES;
  emm_sap.u.emm_cn.u.auth_res = &emm_cn_auth_res;
  rc                          = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_auth_param_fail(mme_ue_s1ap_id_t ue_id, nas_cause_t cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                              = RETURNerror;
  emm_sap_t emm_sap                   = {0};
  emm_cn_auth_fail_t emm_cn_auth_fail = {0};

  emm_cn_auth_fail.cause = cause;
  emm_cn_auth_fail.ue_id = ue_id;

  emm_sap.primitive            = EMMCN_AUTHENTICATION_PARAM_FAIL;
  emm_sap.u.emm_cn.u.auth_fail = &emm_cn_auth_fail;
  rc                           = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_ula_success(mme_ue_s1ap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                                  = RETURNerror;
  emm_sap_t emm_sap                       = {0};
  emm_cn_ula_success_t emm_cn_ula_success = {0};
  emm_cn_ula_success.ue_id                = ue_id;
  emm_sap.primitive                       = EMMCN_ULA_SUCCESS;
  emm_sap.u.emm_cn.u.emm_cn_ula_success   = &emm_cn_ula_success;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "Received S6a-Update Location Answer Success for ue_id "
      "= " MME_UE_S1AP_ID_FMT "\n",
      emm_cn_ula_success.ue_id);
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_cs_respose_success(
    emm_cn_cs_response_success_t* cs_response_success) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  emm_sap.primitive                             = EMMCN_CS_RESPONSE_SUCCESS;
  emm_sap.u.emm_cn.u.emm_cn_cs_response_success = cs_response_success;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "Handle Create Session Response Success at NAS for ue_id "
      "= " MME_UE_S1AP_ID_FMT "\n",
      cs_response_success->ue_id);
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_ula_or_csrsp_fail(emm_cn_ula_or_csrsp_fail_t* ula_or_csrsp_fail) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  emm_sap.primitive                           = EMMCN_ULA_OR_CSRSP_FAIL;
  emm_sap.u.emm_cn.u.emm_cn_ula_or_csrsp_fail = ula_or_csrsp_fail;
  rc                                          = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_create_dedicated_bearer(
    emm_cn_activate_dedicated_bearer_req_t* emm_cn_activate) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};
  emm_sap.primitive = _EMMCN_ACTIVATE_DEDICATED_BEARER_REQ;
  emm_sap.u.emm_cn.u.activate_dedicated_bearer_req = emm_cn_activate;
  rc                                               = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_delete_dedicated_bearer(
    emm_cn_deactivate_dedicated_bearer_req_t* emm_cn_deactivate) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};
  emm_sap.primitive = _EMMCN_DEACTIVATE_BEARER_REQ;
  emm_sap.u.emm_cn.u.deactivate_dedicated_bearer_req = emm_cn_deactivate;
  rc                                                 = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_implicit_detach_ue_ind(mme_ue_s1ap_id_t ue_id) {
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap.primitive                               = EMMCN_IMPLICIT_DETACH_UE;
  emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = ue_id;
  rc                                              = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_nw_initiated_detach_ue_request(
    emm_cn_nw_initiated_detach_ue_t* const nw_initiated_detach_p) {
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap.primitive = EMMCN_NW_INITIATED_DETACH_UE;
  emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.ue_id =
      nw_initiated_detach_p->ue_id;

  if ((nw_initiated_detach_p->detach_type == HSS_INITIATED_EPS_DETACH) ||
      (nw_initiated_detach_p->detach_type == MME_INITIATED_EPS_DETACH)) {
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.detach_type =
        NW_DETACH_TYPE_RE_ATTACH_NOT_REQUIRED;
  } else if (nw_initiated_detach_p->detach_type == SGS_INITIATED_IMSI_DETACH) {
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.detach_type =
        NW_DETACH_TYPE_IMSI_DETACH;
  }
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_downlink_unitdata(itti_sgsap_downlink_unitdata_t* dl_unitdata) {
  imsi64_t imsi64       = INVALID_IMSI64;
  int rc                = RETURNerror;
  emm_context_t* ctxt   = NULL;
  emm_sap_t emm_sap     = {0};
  emm_as_data_t* emm_as = &emm_sap.u.emm_as.u.data;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  DevAssert(dl_unitdata);

  IMSI_STRING_TO_IMSI64(dl_unitdata->imsi, &imsi64);

  OAILOG_DEBUG(LOG_NAS_EMM, "Handling imsi " IMSI_64_FMT "\n", imsi64);

  ctxt = emm_context_get_by_imsi(&_emm_data, imsi64);

  if (!(ctxt)) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  emm_as->nas_info = EMM_AS_NAS_DL_NAS_TRANSPORT;
  emm_as->nas_msg  = bstrcpy(dl_unitdata->nas_msg_container);
  /*
   * Set the UE identifier
   */
  emm_as->ue_id =
      PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context)->mme_ue_s1ap_id;
  /*
   * Setup EPS NAS security data
   */
  emm_as_set_security_data(&emm_as->sctx, &ctxt->_security, false, true);
  /*
   * Notify EMM-AS SAP that Downlink Nas transport message has to be sent to the
   * ue
   */
  emm_sap.primitive = EMMAS_DATA_REQ;
  rc                = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

int encode_mobileid_imsi_tmsi(
    MobileIdentity* out, MobileIdentity in, uint8_t typeofidentity)

{
  if (typeofidentity == MOBILE_IDENTITY_IMSI) {
    out->imsi.digit1               = in.imsi.digit1;
    out->imsi.oddeven              = in.imsi.oddeven;
    out->imsi.typeofidentity       = in.imsi.typeofidentity;
    out->imsi.digit2               = in.imsi.digit2;
    out->imsi.digit3               = in.imsi.digit3;
    out->imsi.digit4               = in.imsi.digit4;
    out->imsi.digit5               = in.imsi.digit5;
    out->imsi.digit6               = in.imsi.digit6;
    out->imsi.digit7               = in.imsi.digit7;
    out->imsi.digit8               = in.imsi.digit8;
    out->imsi.digit9               = in.imsi.digit9;
    out->imsi.digit10              = in.imsi.digit10;
    out->imsi.digit11              = in.imsi.digit11;
    out->imsi.digit12              = in.imsi.digit12;
    out->imsi.digit13              = in.imsi.digit13;
    out->imsi.digit14              = in.imsi.digit14;
    out->imsi.digit15              = in.imsi.digit15;
    out->imsi.numOfValidImsiDigits = in.imsi.numOfValidImsiDigits;
  } else if (typeofidentity == MOBILE_IDENTITY_TMSI) {
    out->tmsi.digit1               = in.tmsi.digit1;
    out->tmsi.oddeven              = in.tmsi.oddeven;
    out->tmsi.typeofidentity       = in.tmsi.typeofidentity;
    out->tmsi.digit2               = in.tmsi.digit2;
    out->tmsi.digit3               = in.tmsi.digit3;
    out->tmsi.digit4               = in.tmsi.digit4;
    out->tmsi.digit5               = in.tmsi.digit5;
    out->tmsi.digit6               = in.tmsi.digit6;
    out->tmsi.digit7               = in.tmsi.digit7;
    out->tmsi.digit8               = in.tmsi.digit8;
    out->tmsi.digit9               = in.tmsi.digit9;
    out->tmsi.digit10              = in.tmsi.digit10;
    out->tmsi.digit11              = in.tmsi.digit11;
    out->tmsi.digit12              = in.tmsi.digit12;
    out->tmsi.digit13              = in.tmsi.digit13;
    out->tmsi.digit14              = in.tmsi.digit14;
    out->tmsi.digit15              = in.tmsi.digit15;
    out->tmsi.numOfValidImsiDigits = in.tmsi.numOfValidImsiDigits;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

//------------------------------------------------------------------------------
int nas_proc_cs_domain_location_updt_fail(
    SgsRejectCause_t cause, lai_t* lai, mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  int rc                                                      = RETURNerror;
  emm_sap_t emm_sap                                           = {0};
  emm_cn_cs_domain_location_updt_fail_t cs_location_updt_fail = {0};

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap.primitive = EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL;
  cs_location_updt_fail =
      emm_sap.u.emm_cn.u.emm_cn_cs_domain_location_updt_fail;

  cs_location_updt_fail.ue_id = mme_ue_s1ap_id;
  // LAI
  if (lai) {
    memcpy(&(cs_location_updt_fail.laicsfb), lai, sizeof(lai_t));
    cs_location_updt_fail.presencemask = LAI;
  }
  // SGS Reject Cause
  cs_location_updt_fail.reject_cause = map_sgs_emm_cause(cause);

  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int nas_proc_sgs_release_req(itti_sgsap_release_req_t* sgs_release_req) {
  imsi64_t imsi64     = INVALID_IMSI64;
  int rc              = RETURNerror;
  emm_context_t* ctxt = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  DevAssert(sgs_release_req);

  IMSI_STRING_TO_IMSI64(sgs_release_req->imsi, &imsi64);

  OAILOG_DEBUG(LOG_NAS_EMM, "Handling imsi " IMSI_64_FMT "\n", imsi64);

  ctxt = emm_context_get_by_imsi(&_emm_data, imsi64);

  if (!(ctxt)) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * As per spec 29.118 section 5.11.4
   * Check the SGS cause recieved in SGS Release Request
   * if sgs cause is "IMSI unknown" or "IMSI detached for non-EPS services"
   * set the "VLR-Reliable" MM context variable to "false"
   * MME requests the UE to re-attach for non-EPS services
   */
  if ((sgs_release_req->opt_cause == SGS_CAUSE_IMSI_UNKNOWN) ||
      (sgs_release_req->opt_cause ==
       SGS_CAUSE_IMSI_DETACHED_FOR_NONEPS_SERVICE)) {
    // NAS trigger UE to re-attach for non-EPS services.
    mme_ue_s1ap_id_t ue_id =
        PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id;
    // update the ue context vlr_reliable flag to false
    mme_ue_context_update_ue_sgs_vlr_reliable(ue_id, false);
    emm_sap_t emm_sap = {0};
    emm_sap.primitive = EMMCN_NW_INITIATED_DETACH_UE;
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.ue_id = ue_id;
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.detach_type =
        NW_DETACH_TYPE_IMSI_DETACH;
    rc = emm_sap_send(&emm_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_cs_service_notification()                            **
 **                                                                        **
 ** Description: Processes CS Paging Request message from MSC/VLR          **
 **              over SGs interface                                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE identifier                                     **
 **          paging_id   Indicates the identity used for                   **
 **                      paging non-eps services                           **
 **          cli         Calling Line Identification                       **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int nas_proc_cs_service_notification(
    mme_ue_s1ap_id_t ue_id, uint8_t paging_id, bstring cli) {
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap.primitive                = EMMAS_DATA_REQ;
  emm_sap.u.emm_as.u.data.nas_info = EMM_AS_NAS_DATA_CS_SERVICE_NOTIFICATION;
  emm_sap.u.emm_as.u.data.ue_id    = ue_id;
  emm_sap.u.emm_as.u.data.nas_msg  = NULL; /*No Esm container*/
  emm_sap.u.emm_as.u.data.paging_identity = paging_id;
  bassign(emm_sap.u.emm_as.u.data.cli, cli);
  rc = emm_sap_send(&emm_sap);
  if (emm_sap.u.emm_as.u.data.cli) {
    bdestroy(emm_sap.u.emm_as.u.data.cli);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
static nas_cause_t s6a_error_2_nas_cause(uint32_t s6a_error, int experimental) {
  if (experimental == 0) {
    /*
     * Base protocol errors
     */
    switch (s6a_error) {
        /*
         * 3002
         */
      case ER_DIAMETER_UNABLE_TO_DELIVER: /* Fall through */

        /*
         * 3003
         */
      case ER_DIAMETER_REALM_NOT_SERVED: /* Fall through */

        /*
         * 5003
         */
      case ER_DIAMETER_AUTHORIZATION_REJECTED:
        return NAS_CAUSE_IMSI_UNKNOWN_IN_HSS;

        /*
         * 5012
         */
      case ER_DIAMETER_UNABLE_TO_COMPLY: /* Fall through */

        /*
         * 5004
         */
      case ER_DIAMETER_INVALID_AVP_VALUE: /* Fall through */

        /*
         * Any other permanent errors from the diameter base protocol
         */
      default:
        break;
    }
  } else {
    switch (s6a_error) {
        /*
         * 4181
         */
      case DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE:
        return NAS_CAUSE_ILLEGAL_UE;

        /*
         * 5001
         */
      case DIAMETER_ERROR_USER_UNKNOWN:
        return NAS_CAUSE_EPS_SERVICES_AND_NON_EPS_SERVICES_NOT_ALLOWED;

        /*
         * TODO: distinguish GPRS_DATA_SUBSCRIPTION
         */
        /*
         * 5420
         */
      case DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION:
        return NAS_CAUSE_NO_SUITABLE_CELLS_IN_TRACKING_AREA;

        /*
         * 5421
         */
      case DIAMETER_ERROR_RAT_NOT_ALLOWED:
        /*
         * One of the following parameter can be sent depending on
         * operator preference:
         * ROAMING_NOT_ALLOWED_IN_THIS_TRACKING_AREA
         * TRACKING_AREA_NOT_ALLOWED
         * NO_SUITABLE_CELLS_IN_TRACKING_AREA
         */
        return NAS_CAUSE_TRACKING_AREA_NOT_ALLOWED;

        /*
         * 5004 without error diagnostic
         */
      case DIAMETER_ERROR_ROAMING_NOT_ALLOWED:
        return NAS_CAUSE_PLMN_NOT_ALLOWED;

        /*
         * TODO: 5004 with error diagnostic of ODB_HPLMN_APN or
         * ODB_VPLMN_APN
         */
        /*
         * TODO: 5004 with error diagnostic of ODB_ALL_APN
         */
      default:
        break;
    }
  }
  return NAS_CAUSE_NETWORK_FAILURE;
}

// Handle CS domain MM-Information request from MSC/VLR

int nas_proc_cs_domain_mm_information_request(
    itti_sgsap_mm_information_req_t* const mm_information_req_pP) {
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap.primitive = EMMCN_CS_DOMAIN_MM_INFORMATION_REQ;
  emm_sap.u.emm_cn.u.emm_cn_cs_domain_mm_information_req =
      mm_information_req_pP;
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    nas_proc_pdn_disconnect_rsp                                   **
 **                                                                        **
 ** Description: Processes _pdn_disconnect_rsp received from MME APP       **
 **                                                                        **
 ** Inputs:                                                                **
 **      emm_cn_pdn_disconnect_rsp : The received message from MME APP     **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int nas_proc_pdn_disconnect_rsp(
    emm_cn_pdn_disconnect_rsp_t* emm_cn_pdn_disconnect_rsp) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc            = RETURNerror;
  emm_sap_t emm_sap = {0};

  OAILOG_DEBUG(
      LOG_NAS_EMM, "Received pdn_disconnect_rsp for ue id %u\n",
      emm_cn_pdn_disconnect_rsp->ue_id);
  emm_sap.primitive                            = EMMCN_PDN_DISCONNECT_RES;
  emm_sap.u.emm_cn.u.emm_cn_pdn_disconnect_rsp = emm_cn_pdn_disconnect_rsp;
  rc                                           = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
