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
#include "directoryd.h"
#include "amf_config.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include "conversions.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_asDefs.h"
#include "amf_sap.h"
#include "amf_recv.h"
#include "amf_app_state_manager.h"
#include "M5gNasMessage.h"
#include "dynamic_memory_check.h"

namespace magma5g {
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO: This has been taken care in new PR with
                              // multi UE feature
amf_config_t amf_config_handler;

//----------------------------------------------------------------------------
static void amf_directoryd_report_location(uint64_t imsi, uint8_t imsi_len) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  directoryd_report_location(imsi_str);
  OAILOG_DEBUG_UE(LOG_AMF_APP, imsi, " Reported UE location to directoryd\n");
}

//------------------------------------------------------------------------------
void amf_ue_context_update_coll_keys(
    amf_ue_context_t* const amf_ue_context_p,
    ue_m5gmm_context_s* const ue_context_p,
    const gnb_ngap_id_key_t gnb_ngap_id_key,
    const amf_ue_ngap_id_t amf_ue_ngap_id, const imsi64_t imsi,
    const teid_t amf_teid_n11, const guti_m5_t* const guti_p) {
  hashtable_rc_t h_rc                 = HASH_TABLE_OK;
  hash_table_ts_t* amf_state_ue_id_ht = get_amf_ue_state();
  OAILOG_FUNC_IN(LOG_AMF_APP);
  OAILOG_TRACE(
      LOG_AMF_APP,
      "Existing ue context, old_gnb_ue_ngap_id_key %ld "
      "old_amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "old_IMSI " IMSI_64_FMT
      "old_GUTI " GUTI_FMT "\n",
      ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id,
      ue_context_p->amf_context.imsi64,
      GUTI_ARG_M5G(&ue_context_p->amf_context._guti));

  if ((gnb_ngap_id_key != INVALID_GNB_UE_NGAP_ID_KEY) &&
      (ue_context_p->gnb_ngap_id_key != gnb_ngap_id_key)) {
    h_rc = hashtable_uint64_ts_remove(
        amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
        (const hash_key_t) ue_context_p->gnb_ngap_id_key);
    h_rc = hashtable_uint64_ts_insert(
        amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
        (const hash_key_t) gnb_ngap_id_key, amf_ue_ngap_id);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_ERROR_UE(
          LOG_AMF_APP, imsi,
          "Error could not update this ue context %p "
          "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT
          "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " %s\n",
          ue_context_p, ue_context_p->gnb_ue_ngap_id,
          ue_context_p->amf_ue_ngap_id, hashtable_rc_code2string(h_rc));
    }
    ue_context_p->gnb_ngap_id_key = gnb_ngap_id_key;
  }

  if (amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    if (ue_context_p->amf_ue_ngap_id != amf_ue_ngap_id) {
      h_rc = hashtable_ts_remove(
          amf_state_ue_id_ht, (const hash_key_t) ue_context_p->amf_ue_ngap_id,
          (void**) &ue_context_p);
      h_rc = hashtable_ts_insert(
          amf_state_ue_id_ht, (const hash_key_t) amf_ue_ngap_id,
          (void*) ue_context_p);

      if (HASH_TABLE_OK != h_rc) {
        // TODO: this method is deprecated and will be removed once the AMF's
        // context is migrated to map in the upcoming multi-UE PR
        OAILOG_ERROR(
            LOG_AMF_APP,
            "Insertion of Hash entry failed for  "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
            amf_ue_ngap_id);
      }
      ue_context_p->amf_ue_ngap_id = amf_ue_ngap_id;
    }
  } else {
    // TODO: this method is deprecated and will be removed once the AMF's
    // context is migrated to map in the upcoming multi-UE PR
    OAILOG_ERROR(
        LOG_AMF_APP, "Invalid  amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
        amf_ue_ngap_id);
  }

  h_rc = hashtable_uint64_ts_remove(
      amf_ue_context_p->imsi_amf_ue_id_htbl,
      (const hash_key_t) ue_context_p->amf_context.imsi64);

  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    h_rc = hashtable_uint64_ts_insert(
        amf_ue_context_p->imsi_amf_ue_id_htbl, (const hash_key_t) imsi,
        amf_ue_ngap_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }

  if (HASH_TABLE_OK != h_rc) {
    // TODO: this method is deprecated and will be removed once the AMF's
    // context is migrated to map in the upcoming multi-UE PR
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Insertion of Hash entry failed for  amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
            PRIX32 " \n",
        amf_ue_ngap_id);
  }

  amf_directoryd_report_location(
      ue_context_p->amf_context.imsi64, ue_context_p->amf_context.imsi.length);
  h_rc = hashtable_uint64_ts_remove(
      amf_ue_context_p->tun11_ue_context_htbl,
      (const hash_key_t) ue_context_p->amf_teid_n11);

  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    h_rc = hashtable_uint64_ts_insert(
        amf_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) amf_teid_n11, (uint64_t) amf_ue_ngap_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }

  if (HASH_TABLE_OK != h_rc) {
    // TODO: this method is deprecated and will be removed once the AMF's
    // context is migrated to map in the upcoming multi-UE PR
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Insertion of Hash entry failed for  amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
            PRIX32 " \n",
        amf_ue_ngap_id);
  }

  ue_context_p->amf_teid_n11 = amf_teid_n11;

  if (guti_p) {
    if ((guti_p->guamfi.amf_code !=
         ue_context_p->amf_context.m5_guti.guamfi.amf_code) ||
        (guti_p->guamfi.amf_gid !=
         ue_context_p->amf_context.m5_guti.guamfi.amf_gid) ||
        (guti_p->m_tmsi != ue_context_p->amf_context.m5_guti.m_tmsi) ||
        (guti_p->guamfi.plmn.mcc_digit1 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit1) ||
        (guti_p->guamfi.plmn.mcc_digit2 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
        (guti_p->guamfi.plmn.mcc_digit3 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3) ||
        (ue_context_p->amf_ue_ngap_id != amf_ue_ngap_id)) {
      h_rc = obj_hashtable_uint64_ts_remove(
          amf_ue_context_p->guti_ue_context_htbl,
          &ue_context_p->amf_context.m5_guti, sizeof(*guti_p));
      if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
        h_rc = obj_hashtable_uint64_ts_insert(
            amf_ue_context_p->guti_ue_context_htbl, (const void* const) guti_p,
            sizeof(*guti_p), (uint64_t) amf_ue_ngap_id);
      } else {
        h_rc = HASH_TABLE_KEY_NOT_EXISTS;
      }
      if (HASH_TABLE_OK != h_rc) {
        // TODO: this method is deprecated and will be removed once the AMF's
        // context is migrated to map in the upcoming multi-UE PR
        OAILOG_ERROR(
            LOG_AMF_APP,
            "Insertion of Hash entry failed for  "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
            amf_ue_ngap_id);
      }
      ue_context_p->amf_context.m5_guti = *guti_p;
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//----------------------------------------------------------------------------------------------
/* This is deprecated function and removed in upcoming PRs related to
 * Service request and Periodic Reg updating.*/
static bool amf_app_construct_guti(
    const plmn_t* const plmn_p, const s_tmsi_m5_t* const s_tmsi_p,
    guti_m5_t* const guti_p) {
  /*
   * This is a helper function to construct GUTI from S-TMSI. It uses PLMN id
   * and AMF Group Id of the serving AMF for this purpose.
   *
   */
  bool is_guti_valid =
      false;  // Set to true if serving AMF is found and GUTI is constructed
  uint8_t num_amf         = 0;  // Number of configured AMF in the AMF pool
  guti_p->m_tmsi          = s_tmsi_p->m_tmsi;
  guti_p->guamfi.amf_code = s_tmsi_p->amf_code;
  // Create GUTI by using PLMN Id and AMF-Group Id of serving AMF
  OAILOG_DEBUG(
      LOG_AMF_APP,
      "Construct GUTI using S-TMSI received form UE and AMG Group Id and PLMN "
      "id from AMF Conf: %u, %u \n",
      s_tmsi_p->m_tmsi, s_tmsi_p->amf_code);
  amf_config_read_lock(&amf_config_handler);
  /*
   * Check number of MMEs in the pool.
   * At present it is assumed that one AMF is supported in AMF pool but in case
   * there are more than one AMF configured then search the serving AMF using
   * AMF code. Assumption is that within one PLMN only one pool of AMF will be
   * configured
   */
  if (amf_config_handler.guamfi.nb > 1) {
    OAILOG_DEBUG(LOG_AMF_APP, "More than one AMFs are configured.");
  }
  for (num_amf = 0; num_amf < amf_config_handler.guamfi.nb; num_amf++) {
    /*Verify that the AMF code within S-TMSI is same as what is configured in
     * AMF conf*/
    if ((plmn_p->mcc_digit2 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mcc_digit2) &&
        (plmn_p->mcc_digit1 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mcc_digit1) &&
        (plmn_p->mnc_digit3 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mnc_digit3) &&
        (plmn_p->mcc_digit3 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mcc_digit3) &&
        (plmn_p->mnc_digit2 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mnc_digit2) &&
        (plmn_p->mnc_digit1 ==
         amf_config_handler.guamfi.guamfi[num_amf].plmn.mnc_digit1) &&
        (guti_p->guamfi.amf_code ==
         amf_config_handler.guamfi.guamfi[num_amf].amf_code)) {
      break;
    }
  }
  if (num_amf >= amf_config_handler.guamfi.nb) {
    OAILOG_DEBUG(LOG_AMF_APP, "No AMF serves this UE");
  } else {
    guti_p->guamfi.plmn    = amf_config_handler.guamfi.guamfi[num_amf].plmn;
    guti_p->guamfi.amf_gid = amf_config_handler.guamfi.guamfi[num_amf].amf_gid;
    is_guti_valid          = true;
  }
  amf_config_unlock(&amf_config_handler);
  return is_guti_valid;
}

//------------------------------------------------------------------------------
// Get existing GUTI details
ue_m5gmm_context_s* amf_ue_context_exists_guti(
    amf_ue_context_t* const amf_ue_context_p, const guti_m5_t* const guti_p) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t amf_ue_ngap_id64 = 0;
  h_rc                      = obj_hashtable_uint64_ts_get(
      amf_ue_context_p->guti_ue_context_htbl, (const void*) guti_p,
      sizeof(*guti_p), &amf_ue_ngap_id64);

  if (HASH_TABLE_OK == h_rc) {
    // TODO: this method is deprecated and will be removed once the AMF's
    // context is migrated to map in the upcoming multi-UE PR
  } else {
    OAILOG_WARNING(LOG_AMF_APP, " No GUTI hashtable for GUTI ");
  }

  return NULL;
}

//-----------------------------------------------------------------------------------------
/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_intial_ue_message()                                **
 **                                                                        **
 ** Description: Processes Initial UE message                              **
 **                                                                        **
 ** Inputs:  amf_app_desc_p:    amf application descriptors                **
 **      initial_pP:      ngap initial ue message structure                **
 **                                                                        **
 **      Return:    imsi value                                             **
 **                                                                        **
 ***************************************************************************/
imsi64_t amf_app_handle_initial_ue_message(
    amf_app_desc_t* amf_app_desc_p,
    itti_ngap_initial_ue_message_t* const initial_pP) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ue_m5gmm_context_s* ue_context_p  = NULL;
  bool is_guti_valid                = false;
  bool is_mm_ctx_new                = false;
  gnb_ngap_id_key_t gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  imsi64_t imsi64                   = INVALID_IMSI64;
  guti_m5_t guti;
  plmn_t plmn;

  if (initial_pP->amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "AMF UE NGAP Id (" AMF_UE_NGAP_ID_FMT ") is already assigned\n",
        initial_pP->amf_ue_ngap_id);
  }

  // Check if there is any existing UE context using S-TMSI/GUTI
  if (initial_pP->is_s_tmsi_valid) {
    /* This check is not used in this PR and code got changed in upcoming PRs
     * hence not-used functions are take out
     */
    OAILOG_DEBUG(
        LOG_AMF_APP, "INITIAL UE Message: Valid amf_code and S-TMSI received ");
    guti.guamfi.plmn        = {0};
    guti.guamfi.amf_gid     = 0;
    guti.guamfi.amf_code    = 0;
    guti.guamfi.amf_Pointer = 0;
    guti.m_tmsi             = INVALID_M_TMSI;
    plmn.mcc_digit1         = initial_pP->tai.plmn.mcc_digit1;
    plmn.mcc_digit2         = initial_pP->tai.plmn.mcc_digit2;
    plmn.mcc_digit3         = initial_pP->tai.plmn.mcc_digit3;
    plmn.mnc_digit1         = initial_pP->tai.plmn.mnc_digit1;
    plmn.mnc_digit2         = initial_pP->tai.plmn.mnc_digit2;
    plmn.mnc_digit3         = initial_pP->tai.plmn.mnc_digit3;
    is_guti_valid =
        amf_app_construct_guti(&plmn, &(initial_pP->opt_s_tmsi), &guti);
    // create a new ue context if nothing is found
    if (is_guti_valid) {
      ue_context_p =
          amf_ue_context_exists_guti(&amf_app_desc_p->amf_ue_contexts, &guti);
      if (ue_context_p) {
        initial_pP->amf_ue_ngap_id = ue_context_p->amf_ue_ngap_id;
        if (ue_context_p->gnb_ngap_id_key != INVALID_GNB_UE_NGAP_ID_KEY) {
          /*
           * Ideally this should never happen. When UE moves to IDLE,
           * this key is set to INVALID.
           * Note - This can happen if eNB detects RLF late and by that time
           * UE sends Initial NAS message via new RRC connection.
           * However if this key is valid, remove the key from the hashtable.
           */
          OAILOG_ERROR(
              LOG_AMF_APP,
              "AMF_APP_INITAIL_UE_MESSAGE: gnb_ngap_id_key %ld has "
              "valid value \n",
              ue_context_p->gnb_ngap_id_key);
          amf_app_ue_context_release(
              ue_context_p, ue_context_p->ue_context_rel_cause);
          hashtable_uint64_ts_remove(
              amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl,
              (const hash_key_t) ue_context_p->gnb_ngap_id_key);
          ue_context_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
        }
        // Update AMF UE context with new gnb_ue_ngap_id
        ue_context_p->gnb_ue_ngap_id = initial_pP->gnb_ue_ngap_id;
        amf_ue_context_update_coll_keys(
            &amf_app_desc_p->amf_ue_contexts, ue_context_p, gnb_ngap_id_key,
            ue_context_p->amf_ue_ngap_id, ue_context_p->amf_context.imsi64,
            ue_context_p->amf_teid_n11, &guti);
        imsi64 = ue_context_p->amf_context.imsi64;
      }
    } else {
      // TODO This piece of code got changed in upcoming PRs with feature
      // like Service Req and Periodic Reg Updating.
    }
  } else {
    OAILOG_DEBUG(
        LOG_AMF_APP, "AMF_APP_INITIAL_UE_MESSAGE from NGAP,without S-TMSI. \n");
  }
  // create a new ue context if nothing is found
  if (ue_context_p == NULL) {
    OAILOG_DEBUG(LOG_AMF_APP, " UE context doesn't exist -> create one\n");
    if (!(ue_context_p = amf_create_new_ue_context())) {
      OAILOG_INFO(LOG_AMF_APP, "Failed to create context \n");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
    }
    // Allocate new amf_ue_ngap_id
    ue_context_p->amf_ue_ngap_id = amf_app_ctx_get_new_ue_id(
        &amf_app_desc_p->amf_app_ue_ngap_id_generator);
    if (ue_context_p->amf_ue_ngap_id == INVALID_AMF_UE_NGAP_ID) {
      OAILOG_CRITICAL(
          LOG_AMF_APP,
          "AMF_APP_INITIAL_UE_MESSAGE. AMF_UE_NGAP_ID allocation Failed.\n");
      amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
    }
    AMF_APP_GNB_NGAP_ID_KEY(
        ue_context_p->gnb_ngap_id_key, initial_pP->gnb_id,
        initial_pP->gnb_ue_ngap_id);
    amf_insert_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);

    notify_ngap_new_ue_amf_ngap_id_association(ue_context_p);
    s_tmsi_m5_t s_tmsi = {0};
    if (initial_pP->is_s_tmsi_valid) {
      s_tmsi = initial_pP->opt_s_tmsi;
    } else {
      s_tmsi.amf_code = 0;
      s_tmsi.m_tmsi   = INVALID_M_TMSI;
    }

    OAILOG_DEBUG(
        LOG_AMF_APP,
        " Sending NAS Establishment Indication to NAS for ue_id = "
        "(%d)\n",
        ue_context_p->amf_ue_ngap_id);
    nas_proc_establish_ind(
        ue_context_p->amf_ue_ngap_id, is_mm_ctx_new, initial_pP->tai,
        initial_pP->ecgi, initial_pP->m5g_rrc_establishment_cause, s_tmsi,
        initial_pP->nas);
  }
  return RETURNok;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_uplink_nas_message()                               **
 **                                                                        **
 ** Description: Handle uplink nas message                                 **
 **                                                                        **
 ** Inputs:  amf_app_desc_p:    amf application descriptors                **
 **      msg:      nstring msg                                             **
 **                                                                        **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int amf_app_handle_uplink_nas_message(
    amf_app_desc_t* amf_app_desc_p, bstring msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  OAILOG_DEBUG(LOG_AMF_APP, " Received NAS UPLINK DATA from NGAP\n");
  if (msg) {
    amf_sap_t amf_sap;
    /*
     * Notify the AMF procedure call manager that data transfer
     * indication has been received from the Access-Stratum sublayer
     */
    amf_sap.primitive = AMFAS_ESTABLISH_REQ;
    // TODO: hardcoded for now, addressed in the upcoming multi-UE PR
    amf_sap.u.amf_as.u.establish.ue_id   = 1;
    amf_sap.u.amf_as.u.establish.nas_msg = msg;
    msg                                  = NULL;
    rc                                   = amf_sap_send(&amf_sap);
  } else {
    OAILOG_WARNING(
        LOG_NAS, "Received NAS message in uplink is NULL for ue_id = (%u)\n",
        amf_app_desc_p->amf_app_ue_ngap_id_generator);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/* Recieved the session created reponse message from SMF. Populate and Send PDU
 * Session Resource Setup Request message to gNB and  PDU Session Establishment
 * Accept Message to UE*/
void amf_app_handle_pdu_session_response(
    itti_n11_create_pdu_session_response_t* pdu_session_resp) {
  DLNASTransportMsg encode_msg;
  SmfMsg* smf_msg;
  bstring buffer;
  uint32_t len;
  nas5g_error_code_t rc = M5G_AS_SUCCESS;
  int amf_rc            = RETURNerror;
  ue_m5gmm_context_s* ue_context;
  smf_context_t* smf_ctx;
  uint32_t bytes = 0;
  // TODO: hardcoded for now, addressed in the upcoming multi-UE PR
  uint32_t ue_id = 1;

  // Handle smf_context
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context) {
    smf_ctx = &(ue_context->amf_context.smf_context);
  } else {
    ue_context = &ue_m5gmm_global_context;
    smf_ctx    = &ue_m5gmm_global_context.amf_context
                   .smf_context;  // TODO: This has been taken care in new PR
                                  // with multi UE feature
  }
  smf_ctx->dl_session_ambr = pdu_session_resp->session_ambr.downlink_units;
  smf_ctx->dl_ambr_unit    = pdu_session_resp->session_ambr.downlink_unit_type;
  smf_ctx->ul_session_ambr = pdu_session_resp->session_ambr.uplink_units;
  smf_ctx->ul_ambr_unit    = pdu_session_resp->session_ambr.uplink_unit_type;
  /* Message construction for PDUSessionResourceSetupRequest
   * to gNB with UPF TEID info
   */
  memcpy(
      &(smf_ctx->pdu_resource_setup_req
            .pdu_session_resource_setup_request_transfer
            .qos_flow_setup_request_list),
      &(pdu_session_resp->qos_list), sizeof(qos_flow_request_list));
  memcpy(
      smf_ctx->gtp_tunnel_id.upf_gtp_teid_ip_addr,
      pdu_session_resp->upf_endpoint.end_ipv4_addr,
      sizeof(smf_ctx->gtp_tunnel_id.upf_gtp_teid_ip_addr));
  memcpy(
      smf_ctx->gtp_tunnel_id.upf_gtp_teid, pdu_session_resp->upf_endpoint.teid,
      sizeof(smf_ctx->gtp_tunnel_id.upf_gtp_teid));

  OAILOG_DEBUG(
      LOG_AMF_APP,
      "Sending message to gNB for PDUSessionResourceSetupRequest\n");
  amf_rc = pdu_session_resource_setup_request(ue_context, ue_id);
  if (amf_rc != RETURNok) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "Failure in sending message to gNB for "
        "PDUSessionResourceSetupRequest\n");
    /* TODO: in future add negative case handling, send pdu reject
     * command to UE and release message to SMF
     */
  }
  smf_msg = &encode_msg.payload_container.smf_msg;

  // Message construction for PDU Establishment Accept
  // AmfHeader
  encode_msg.extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  encode_msg.spare_half_octet.spare     = 0x00;
  encode_msg.sec_header_type.sec_hdr    = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  encode_msg.message_type.msg_type      = DLNASTRANSPORT;
  encode_msg.payload_container_type.iei = PAYLOAD_CONTAINER_TYPE;

  // SmfMsg
  encode_msg.payload_container_type.type_val = N1_SM_INFO;
  encode_msg.payload_container.iei           = PAYLOAD_CONTAINER;
  encode_msg.pdu_session_identity.iei        = PDU_SESSION_IDENTITY;
  encode_msg.pdu_session_identity.pdu_session_id =
      pdu_session_resp->pdu_session_id;
  smf_msg->header.extended_protocol_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->header.pdu_session_id           = pdu_session_resp->pdu_session_id;
  smf_msg->header.message_type             = PDU_SESSION_ESTABLISHMENT_ACCEPT;
  smf_msg->header.procedure_transaction_id = smf_ctx->smf_proc_data.pti.pti;
  smf_msg->msg.pdu_session_estab_accept.extended_protocol_discriminator
      .extended_proto_discriminator = M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_identity.pdu_session_id =
      pdu_session_resp->pdu_session_id;
  smf_msg->msg.pdu_session_estab_accept.pti.pti =
      smf_ctx->smf_proc_data.pti.pti;
  smf_msg->msg.pdu_session_estab_accept.message_type.msg_type =
      PDU_SESSION_ESTABLISHMENT_ACCEPT;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_type.type_val =
      pdu_session_resp->pdu_session_type;
  smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val = SSC_MODE_ONE;
  memset(
      smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info, '\0',
      sizeof(smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info));
  memcpy(
      smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info,
      pdu_session_resp->pdu_address.redirect_server_address, PDU_ADDR_IPV4_LEN);
  smf_msg->msg.pdu_session_estab_accept.pdu_address.type_val = PDU_ADDR_TYPE;

  /* QOSrules are hardcoded as it is not exchanged in AMF-SMF
   * gRPC calls as of now, handled in upcoming PR
   * TODO: get the rules for the session from SMF and use it here
   */
  smf_msg->msg.pdu_session_estab_accept.qos_rules.length = 0x9;
  QOSRule qos_rule;
  qos_rule.qos_rule_id         = 0x1;
  qos_rule.len                 = 0x6;
  qos_rule.rule_oper_code      = 0x1;
  qos_rule.dqr_bit             = 0x1;
  qos_rule.no_of_pkt_filters   = 0x1;
  qos_rule.qos_rule_precedence = 0xff;
  qos_rule.spare               = 0x0;
  qos_rule.segregation         = 0x0;
  qos_rule.qfi                 = 0x6;
  NewQOSRulePktFilter new_qos_rule_pkt_filter;
  new_qos_rule_pkt_filter.spare          = 0x0;
  new_qos_rule_pkt_filter.pkt_filter_dir = 0x3;
  new_qos_rule_pkt_filter.pkt_filter_id  = 0x1;
  new_qos_rule_pkt_filter.len            = 0x1;
  uint8_t contents                       = 0x1;
  memcpy(
      new_qos_rule_pkt_filter.contents, &contents, new_qos_rule_pkt_filter.len);
  memcpy(
      qos_rule.new_qos_rule_pkt_filter, &new_qos_rule_pkt_filter,
      1 * sizeof(NewQOSRulePktFilter));
  memcpy(
      smf_msg->msg.pdu_session_estab_accept.qos_rules.qos_rule, &qos_rule,
      1 * sizeof(QOSRule));
  smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_unit =
      smf_ctx->dl_ambr_unit;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_unit =
      smf_ctx->ul_ambr_unit;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_session_ambr =
      smf_ctx->dl_session_ambr;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_session_ambr =
      smf_ctx->ul_session_ambr;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.length = AMBR_LEN;
  encode_msg.payload_container.len = PDU_ESTAB_ACCPET_PAYLOAD_CONTAINER_LEN;
  len                              = PDU_ESTAB_ACCPET_NAS_PDU_LEN;
  buffer                           = bfromcstralloc(len, "\0");
  bytes = encode_msg.EncodeDLNASTransportMsg(&encode_msg, buffer->data, len);
  if (bytes > 0) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "NAS encode success, sent PDU Establishment Accept to UE\n");
    buffer->slen =
        bytes +
        3;  // TODO fix the buffer length returned in NAS encode function
    amf_app_handle_nas_dl_req(ue_id, buffer, rc);

  } else {
    bdestroy_wrapper(&buffer);
  }
}

/* Handling PDU Session Resource Setup Response sent from gNB*/
void amf_app_handle_resource_setup_response(
    itti_ngap_pdusessionresource_setup_rsp_t session_setup_resp) {
  /* Check if failure message is not NULL and if NULL,
   * it is successful message from gNB.
   * Nothing to in this case. If failure message comes from gNB
   * AMF need to report this failed message to SMF
   *
   * NOTE: only handling success part not failure part
   * will be handled later
   */
  OAILOG_DEBUG(
      LOG_AMF_APP, " handling uplink PDU session setup response message\n");
  if (session_setup_resp.pduSessionResource_setup_list.no_of_items > 0) {
    /* This is success case and we need not to send message to SMF
     * and drop the message here
     */
    OAILOG_DEBUG(
        LOG_AMF_APP,
        " this is success case and no need to hadle anything and drop "
        "the message\n");
  amf_ue_ngap_id_t    ue_id;
  amf_smf_establish_t amf_smf_grpc_ies;
  ue_m5gmm_context_s* ue_context  = nullptr;
  amf_context_t*      amf_context = nullptr;
  smf_context_t*      smf_ctx     = nullptr;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  ue_id = session_setup_resp.amf_ue_ngap_id;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // Handling of ue context
  if(!ue_context){
     ue_context = &ue_m5gmm_global_context;
  }
  smf_ctx = &ue_context->amf_context.smf_context;
  OAILOG_DEBUG(LOG_AMF_APP, "filling gNB TEID info in smf context \n");
  //Store gNB ip and TEID in respective smf_context
  memset(&smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0',
          sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr));
  memcpy(&smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
        &session_setup_resp.pduSessionResource_setup_list.item[0]
        .PDU_Session_Resource_Setup_Response_Transfer.tunnel.gTP_TEID, 4);
  OAILOG_DEBUG(LOG_AMF_APP, "filling gNB TEID info in gtp_ip_address \n");
  memcpy(&smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr,
        &session_setup_resp.pduSessionResource_setup_list.item[0]
        .PDU_Session_Resource_Setup_Response_Transfer.tunnel
        .transportLayerAddress, 4); //time being 4 byte is copying.
  OAILOG_DEBUG(LOG_AMF_APP, "printing both teid and ip_address of gNB\n");
  OAILOG_DEBUG(LOG_AMF_APP, "IP address %02x %02x %02x %02x  and TEID %02x %02x %02x %02x \n",
         smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[0], smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[1],
         smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[2], smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[3],
         smf_ctx->gtp_tunnel_id.gnb_gtp_teid[0], smf_ctx->gtp_tunnel_id.gnb_gtp_teid[1],
         smf_ctx->gtp_tunnel_id.gnb_gtp_teid[0], smf_ctx->gtp_tunnel_id.gnb_gtp_teid[3]);
  //Incrementing the  pdu session version
  smf_ctx->pdu_session_version++;
  /*Copy respective gNB fields to amf_smf_establish_t compartible to gRPC message*/
  memset(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr, '\0', sizeof(amf_smf_grpc_ies.gnb_gtp_teid_ip_addr));
  memset(&amf_smf_grpc_ies.gnb_gtp_teid, '\0', sizeof(amf_smf_grpc_ies.gnb_gtp_teid));
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr, &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid, &smf_ctx->gtp_tunnel_id.gnb_gtp_teid, 4);
  amf_smf_grpc_ies.pdu_session_id = smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id;
  IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
  /* Prepare and send gNB setup response message to SMF through gRPC
   * 2nd time PDU session establish message
   */
  create_session_grpc_req_on_gnb_setup_rsp(&amf_smf_grpc_ies,
                  imsi, smf_ctx->pdu_session_version);

   } else {
    // TODO: implement failure message from gNB. messagge to send to SMF
    OAILOG_DEBUG(
        LOG_AMF_APP, " Failure message not handled and dropping the message\n");
  }
}

/* Handling Resource Release Response from gNB */
void amf_app_handle_resource_release_response(
    itti_ngap_pdusessionresource_rel_rsp_t session_rel_resp) {
  /*
   * Release request always should be successful.
   * This response message will be dropped here as nothing to do.
   * as pdu_session_resource_release_response_transfer is
   * optional as per 38.413 - 9.3.4.2.1
   */
  OAILOG_DEBUG(
      LOG_AMF_APP, " handling uplink PDU session release response message\n");
  if (session_rel_resp.pduSessionResourceReleasedRspList.no_of_items > 0) {
    /* This is success case and we need not to send message to SMF
     * and drop the message here
     */
    OAILOG_DEBUG(
        LOG_AMF_APP,
        " this is success case of release response and no need to "
        "hadle anything and drop the message\n");
  } else {
    // TODO implement failure message from gNB. messagge to send to SMF
    OAILOG_DEBUG(
        LOG_AMF_APP, " Failure message not handled and dropping the message\n");
  }
}

/* This function gets invoked based on the message NGAP_UE_CONTEXT_RELEASE_REQ
 * from gNB/NGAP in UL for handling CM-idle state of UE/IMSI/SUPI.
 * Action logic:
 * - Fetch AMF context, match the no of PDU sessions in message with no of
 *   PDU sessions in AMF_context and cause NGAP_RADIO_NR_GENERATED_REASON
 *   it means gNB RRC-Inactive triggered and UE state must be changed from
 *   CM-connected to CM-Idle state.
 *   Then send message to SMF to change all respective PDU session state
 *   to inactive state.
 * - Retrive the required field of UE, like IMSI and fill gRPC notification
 *   proto structure.
 * - In AMF move the UE/IMSI state to CM-idle
 *   Go over all PDU sessions and change the state to in-active.
 * */
void amf_app_handle_cm_idle_on_ue_context_release(
    itti_ngap_ue_context_release_req_t cm_idle_req) {
  OAILOG_DEBUG(
      LOG_AMF_APP, " Handling UL UE context release for CM-idle for ue id %d\n",
      cm_idle_req.amf_ue_ngap_id);
  /* Currently only one PDU session is considered.
   * for multiple PDU session context (smf_context_t) will be part of vector
   * and no. of PDU sessions can be derived from this vector and compared
   * with NGAP message in future.
   * Now only need to check the cause and proceed further.
   * note: check if UE in connected state else already in idle state
   * nothing to do.
   */
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id;
  ue_m5gmm_context_s* ue_context = nullptr;
  amf_context_t* amf_context     = nullptr;
  smf_context_t* smf_ctx         = nullptr;
  ue_id                          = cm_idle_req.amf_ue_ngap_id;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // TODO: This has been taken care in new PR
  // with multi UE feature
  if (!ue_context) {
    ue_context = &ue_m5gmm_global_context;
  }
  // if UE on REGISTERED_IDLE, so no need to do anyting
  if (ue_context->mm_state == REGISTERED_CONNECTED) {
    // UE in connected state and need to check if cause is proper
    if (cm_idle_req.relCause == NGAP_RADIO_NR_GENERATED_REASON) {
      // Change the respective UE/PDU session state to idle/inactive.
      ue_context->mm_state == REGISTERED_IDLE;
      // Handling of smf_context as vector
      // TODO: This has been taken care in new PR
      // with multi UE feature
      smf_ctx                    = &ue_context->amf_context.smf_context;
      smf_ctx->pdu_session_state = INACTIVE;

      // construct the proto structure and send message to SMF
      rc = amf_smf_notification_send(ue_id, ue_context);

    } else {
      OAILOG_DEBUG(
          LOG_AMF_APP,
          " UE in REGISTERED_CONNECTED state, but cause from NGAP"
          " is wrong for UE ID %d and return\n",
          cm_idle_req.amf_ue_ngap_id);
      return;
    }
  } else {
    /* TODO: Single or multiple PDU session state change notification
     * should be taken care here. amf_smf_notification_send will be used
     * with one more parameter as boolean for idle mode or single PDU
     * session state change. Currently nothing to do
     */
    OAILOG_DEBUG(
        LOG_AMF_APP,
        " UE in REGISTERED_IDLE or CM-idle state, nothing to do"
        " for UE ID %d\n",
        cm_idle_req.amf_ue_ngap_id);
    return;
  }
}

}  // namespace magma5g
