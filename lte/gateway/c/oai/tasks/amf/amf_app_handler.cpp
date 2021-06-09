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
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#include "directoryd.h"
#include "amf_config.h"
#include "dynamic_memory_check.h"
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
#include "n11_messages_types.h"

#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH)                                    \
  LENGTH += QUADLET - (LENGTH % QUADLET)

static int paging_t3513_handler(zloop_t* loop, int timer_id, void* arg);
namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;
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
        "Insertion of Hash entry failed for  "
        "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
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
        "Insertion of Hash entry failed for  "
        "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
        amf_ue_ngap_id);
  }

  ue_context_p->amf_teid_n11 = amf_teid_n11;

  if (guti_p) {
    if ((guti_p->guamfi.amf_set_id !=
         ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
        (guti_p->guamfi.amf_set_id !=
         ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
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
  uint8_t num_amf           = 0;  // Number of configured AMF in the AMF pool
  guti_p->m_tmsi            = s_tmsi_p->m_tmsi;
  guti_p->guamfi.amf_set_id = s_tmsi_p->amf_pointer;
  // Create GUTI by using PLMN Id and AMF-Group Id of serving AMF
  OAILOG_DEBUG(
      LOG_AMF_APP,
      "Construct GUTI using S-TMSI received form UE and AMG Group Id and "
      "PLMN "
      "id from AMF Conf: %u, %u \n",
      s_tmsi_p->m_tmsi, s_tmsi_p->amf_pointer);
  amf_config_read_lock(&amf_config_handler);
  /*
   * Check number of MMEs in the pool.
   * At present it is assumed that one AMF is supported in AMF pool but in
   * case there are more than one AMF configured then search the serving AMF
   * using AMF code. Assumption is that within one PLMN only one pool of AMF
   * will be configured
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
        (guti_p->guamfi.amf_set_id ==
         amf_config_handler.guamfi.guamfi[num_amf].amf_set_id)) {
      break;
    }
  }
  if (num_amf >= amf_config_handler.guamfi.nb) {
    OAILOG_DEBUG(LOG_AMF_APP, "No AMF serves this UE");
  } else {
    guti_p->guamfi.plmn = amf_config_handler.guamfi.guamfi[num_amf].plmn;
    guti_p->guamfi.amf_set_id =
        amf_config_handler.guamfi.guamfi[num_amf].amf_set_id;
    is_guti_valid = true;
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
    //  return amf_ue_context_exists_amf_ue_ngap_id(  //TODO -  NEED-RECHECK
    //    (amf_ue_ngap_id_t) amf_ue_ngap_id64);
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
    OAILOG_INFO(
        LOG_AMF_APP,
        "INITIAL UE Message: Valid amf_set_id and S-TMSI received ");
    guti.guamfi.plmn         = {0};
    guti.guamfi.amf_regionid = 0;
    guti.guamfi.amf_set_id   = 0;
    guti.guamfi.amf_pointer  = 0;
    guti.m_tmsi              = INVALID_M_TMSI;
    plmn.mcc_digit1          = initial_pP->tai.plmn.mcc_digit1;
    plmn.mcc_digit2          = initial_pP->tai.plmn.mcc_digit2;
    plmn.mcc_digit3          = initial_pP->tai.plmn.mcc_digit3;
    plmn.mnc_digit1          = initial_pP->tai.plmn.mnc_digit1;
    plmn.mnc_digit2          = initial_pP->tai.plmn.mnc_digit2;
    plmn.mnc_digit3          = initial_pP->tai.plmn.mnc_digit3;
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
    amf_insert_ue_context(
        ue_context_p->amf_ue_ngap_id, &amf_app_desc_p->amf_ue_contexts,
        ue_context_p);
    ue_context_p->sctp_assoc_id_key = initial_pP->sctp_assoc_id;
    ue_context_p->gnb_ue_ngap_id    = initial_pP->gnb_ue_ngap_id;

    // UEContextRequest
    ue_context_p->ue_context_request = initial_pP->ue_context_request;
    OAILOG_DEBUG(
        LOG_AMF_APP, "ue_context_requext received: %d\n ",
        ue_context_p->ue_context_request);

    notify_ngap_new_ue_amf_ngap_id_association(ue_context_p);
    s_tmsi_m5_t s_tmsi = {0};
    if (initial_pP->is_s_tmsi_valid) {
      s_tmsi = initial_pP->opt_s_tmsi;
    } else {
      s_tmsi.amf_pointer = 0;
      s_tmsi.m_tmsi      = INVALID_M_TMSI;
    }
    is_mm_ctx_new = true;

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
    amf_app_desc_t* amf_app_desc_p, bstring msg, amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  OAILOG_DEBUG(LOG_AMF_APP, " Received NAS UPLINK DATA from NGAP\n");
  if (msg) {
    amf_sap_t amf_sap;
    /*
     * Notify the AMF procedure call manager that data transfer
     * indication has been received from the Access-Stratum sublayer
     */
    amf_sap.primitive                    = AMFAS_ESTABLISH_REQ;
    amf_sap.u.amf_as.u.establish.ue_id   = ue_id;
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

/* Recieved the session created reponse message from SMF. Populate and Send
 * PDU Session Resource Setup Request message to gNB and  PDU Session
 * Establishment Accept Message to UE*/
void amf_app_handle_pdu_session_response(
    itti_n11_create_pdu_session_response_t* pdu_session_resp) {
  DLNASTransportMsg encode_msg;
  int amf_rc = RETURNerror;
  ue_m5gmm_context_s* ue_context;
  smf_context_t* smf_ctx;
  amf_smf_t amf_smf_msg;
  // TODO: hardcoded for now, addressed in the upcoming multi-UE PR
  uint32_t ue_id = 0;

  imsi64_t imsi64;
  IMSI_STRING_TO_IMSI64(pdu_session_resp->imsi, &imsi64);
  // Handle smf_context
  ue_context = lookup_ue_ctxt_by_imsi(imsi64);
  if (ue_context) {
    smf_ctx = amf_smf_context_exists_pdu_session_id(
        ue_context, pdu_session_resp->pdu_session_id);
    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "pdu session  not found for session_id = %u\n",
          pdu_session_resp->pdu_session_id);
      return;
    }
    ue_id = ue_context->amf_ue_ngap_id;
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP, "ue context not found for the imsi=%lu\n", imsi64);
    return;
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
      &(pdu_session_resp->qos_list), sizeof(qos_flow_request_list_t));
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
  amf_rc = pdu_session_resource_setup_request(ue_context, ue_id, smf_ctx);
  if (amf_rc != RETURNok) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "Failure in sending message to gNB for "
        "PDUSessionResourceSetupRequest\n");
    /* TODO: in future add negative case handling, send pdu reject
     * command to UE and release message to SMF
     */
  }
  /*Execute PDU establishement accept from AMF to gnodeb */
  pdu_state_handle_message(
      // ue_context->mm_state, STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT,
      REGISTERED_CONNECTED, STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT,
      // smf_ctx->pdu_session_state, ue_context, amf_smf_msg, NULL,
      CREATING, ue_context, amf_smf_msg, NULL, pdu_session_resp, ue_id);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_app_handle_pdu_session_accept()                           **
 **                                                                        **
 ** Description: Send the PDU establishment accept to gnodeb               **
 **                                                                        **
 ** Inputs:  pdu_session_resp:   pdusession response message               **
 **      ue_id:      ue identity                                           **
 **                                                                        **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int amf_app_handle_pdu_session_accept(
    itti_n11_create_pdu_session_response_t* pdu_session_resp, uint32_t ue_id) {
  nas5g_error_code_t rc = M5G_AS_SUCCESS;

  DLNASTransportMsg* encode_msg;
  amf_nas_message_t msg;
  uint32_t bytes = 0;
  uint32_t len;
  SmfMsg* smf_msg;
  bstring buffer;
  smf_context_t* smf_ctx;
  ue_m5gmm_context_s* ue_context;

  // Handle smf_context
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context) {
    smf_ctx = &(ue_context->amf_context.smf_context);
  } else {
    OAILOG_INFO(LOG_AMF_APP, "UE Context not found for UE ID: %d", ue_id);
  }

  // Message construction for PDU Establishment Accept
  msg.security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.security_protected.plain.amf.header.message_type = DLNASTRANSPORT;
  msg.header.security_header_type =
      SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.sequence_number =
      ue_context->amf_context._security.dl_count.seq_num;

  encode_msg = &msg.security_protected.plain.amf.msg.downlinknas5gtransport;
  smf_msg    = &encode_msg->payload_container.smf_msg;

  // AmfHeader
  encode_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  encode_msg->spare_half_octet.spare     = 0x00;
  encode_msg->sec_header_type.sec_hdr    = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  encode_msg->message_type.msg_type      = DLNASTRANSPORT;
  encode_msg->payload_container_type.iei = PAYLOAD_CONTAINER_TYPE;

  // SmfMsg
  encode_msg->payload_container_type.type_val = N1_SM_INFO;
  encode_msg->payload_container.iei           = PAYLOAD_CONTAINER;
  encode_msg->pdu_session_identity.iei        = PDU_SESSION_IDENTITY;
  encode_msg->pdu_session_identity.pdu_session_id =
      pdu_session_resp->pdu_session_id;

  smf_msg->header.extended_protocol_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->header.pdu_session_id = pdu_session_resp->pdu_session_id;
  smf_msg->header.message_type   = PDU_SESSION_ESTABLISHMENT_ACCEPT;
  // smf_msg->header.procedure_transaction_id = smf_ctx->smf_proc_data.pti.pti;
  smf_msg->header.procedure_transaction_id = 0x01;
  smf_msg->msg.pdu_session_estab_accept.extended_protocol_discriminator
      .extended_proto_discriminator = M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_identity.pdu_session_id =
      pdu_session_resp->pdu_session_id;
  smf_msg->msg.pdu_session_estab_accept.pti.pti = 0x01;
  // smf_ctx->smf_proc_data.pti.pti;
  smf_msg->msg.pdu_session_estab_accept.message_type.msg_type =
      PDU_SESSION_ESTABLISHMENT_ACCEPT;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_type.type_val = 1;
  // smf_msg->msg.pdu_session_estab_accept.pdu_session_type.type_val =
  // pdu_session_resp->pdu_session_type;
  smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val = 1;
  // smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val = SSC_MODE_ONE;

  memset(
      &(smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info), 0, 12);

  for (int i = 0; i < PDU_ADDR_IPV4_LEN; i++) {
    smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info[i] =
        pdu_session_resp->pdu_address.redirect_server_address[i];
  }
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
      pdu_session_resp->session_ambr.downlink_unit_type;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_unit =
      pdu_session_resp->session_ambr.uplink_unit_type;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_session_ambr =
      pdu_session_resp->session_ambr.downlink_units;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_session_ambr =
      pdu_session_resp->session_ambr.uplink_units;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.length = AMBR_LEN;

  //  encode_msg.payload_container.len = PDU_ESTAB_ACCPET_PAYLOAD_CONTAINER_LEN;
  //  len                              = PDU_ESTAB_ACCPET_NAS_PDU_LEN;
  //  buffer                           = bfromcstralloc(len, "\0");
  //  bytes = encode_msg.EncodeDLNASTransportMsg(&encode_msg, buffer->data,
  //  len);
  encode_msg->payload_container.len = 30;
  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF_TEST: start NAS encoding for PDU Session Establishment Accept\n");

  len = 41;  // originally 38 and 30

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);
  if (msg.header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    amf_msg_header* header = &msg.security_protected.plain.amf.header;
    /*
     * Expand size of protected NAS message
     */
    OAILOG_INFO(
        LOG_AMF_APP, "AMF_TEST:before adding sec header, length %d ", len);
    len += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    OAILOG_INFO(
        LOG_AMF_APP, "AMF_TEST:after adding sec header, length %d ", len);
    /*
     * Set header of plain NAS message
     */
    header->extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;
    header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }

  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(
      buffer->data, &msg, len, &ue_context->amf_context._security);
  if (bytes > 0) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "NAS encode success, sent PDU Establishment Accept to UE\n");
    buffer->slen = bytes;
    amf_app_handle_nas_dl_req(ue_id, buffer, rc);

  } else {
    bdestroy_wrapper(&buffer);
  }
  return rc;
}

/* Handling PDU Session Resource Setup Response sent from gNB*/
void amf_app_handle_resource_setup_response(
    itti_ngap_pdusessionresource_setup_rsp_t session_seup_resp) {
  amf_ue_ngap_id_t ue_id;

  ue_m5gmm_context_s* ue_context = nullptr;
  smf_context_t* smf_ctx         = nullptr;

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
  if (session_seup_resp.pduSessionResource_setup_list.no_of_items > 0) {
    ue_id = session_seup_resp.amf_ue_ngap_id;

    ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
    if (ue_context == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "ue context not found for the ue_id=%u\n", ue_id);
      return;
    }
    smf_ctx = amf_smf_context_exists_pdu_session_id(
        ue_context,
        session_seup_resp.pduSessionResource_setup_list.item[0].Pdu_Session_ID);
    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "pdu session  not found for session_id = %d\n",
          session_seup_resp.pduSessionResource_setup_list.item[0]
              .Pdu_Session_ID);
      return;
    }

    /* This is success case and we need not to send message to SMF
     * and drop the message here
     */
    OAILOG_DEBUG(
        LOG_AMF_APP,
        " this is success case and no need to hadle anything and drop "
        "the message\n");
    amf_ue_ngap_id_t ue_id;
    amf_smf_establish_t amf_smf_grpc_ies;
    ue_m5gmm_context_s* ue_context = nullptr;
    amf_context_t* amf_context     = nullptr;
    smf_context_t* smf_ctx         = nullptr;
    char imsi[IMSI_BCD_DIGITS_MAX + 1];

    ue_id = session_seup_resp.amf_ue_ngap_id;

    ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
    // Handling of ue context
    if (!ue_context) {
      OAILOG_INFO(LOG_AMF_APP, "UE Context not found for UE ID: %d", ue_id);
    }
    smf_ctx = &ue_context->amf_context.smf_context;
    OAILOG_DEBUG(LOG_AMF_APP, "filling gNB TEID info in smf context \n");
    // Store gNB ip and TEID in respective smf_context
    memset(
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0',
        sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr));
    memcpy(
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
        &session_seup_resp.pduSessionResource_setup_list.item[0]
             .PDU_Session_Resource_Setup_Response_Transfer.tunnel.gTP_TEID,
        4);
    OAILOG_DEBUG(LOG_AMF_APP, "filling gNB TEID info in gtp_ip_address \n");
    memcpy(
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr,
        &session_seup_resp.pduSessionResource_setup_list.item[0]
             .PDU_Session_Resource_Setup_Response_Transfer.tunnel
             .transportLayerAddress,
        4);  // time being 4 byte is copying.
    OAILOG_DEBUG(LOG_AMF_APP, "printing both teid and ip_address of gNB\n");
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "IP address %02x %02x %02x %02x  and TEID %02x "
        "%02x %02x %02x \n",
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[0],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[1],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[2],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[3],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid[0],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid[1],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid[0],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid[3]);
    // Incrementing the  pdu session version
    smf_ctx->pdu_session_version++;
    /*Copy respective gNB fields to amf_smf_establish_t compartible to gRPC
     * message*/
    memset(
        &amf_smf_grpc_ies.gnb_gtp_teid_ip_addr, '\0',
        sizeof(amf_smf_grpc_ies.gnb_gtp_teid_ip_addr));
    memset(
        &amf_smf_grpc_ies.gnb_gtp_teid, '\0',
        sizeof(amf_smf_grpc_ies.gnb_gtp_teid));
    memcpy(
        &amf_smf_grpc_ies.gnb_gtp_teid_ip_addr,
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
    memcpy(
        &amf_smf_grpc_ies.gnb_gtp_teid, &smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
        4);
    amf_smf_grpc_ies.pdu_session_id =
        session_seup_resp.pduSessionResource_setup_list.item[0].Pdu_Session_ID;
    // smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id;

    IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
    /* Prepare and send gNB setup response message to SMF through gRPC
     * 2nd time PDU session establish message
     */
    create_session_grpc_req_on_gnb_setup_rsp(
        &amf_smf_grpc_ies, imsi, smf_ctx->pdu_session_version);

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
  amf_ue_ngap_id_t ue_id;
  ue_m5gmm_context_s* ue_context = nullptr;
  smf_context_t* smf_ctx         = nullptr;
  ue_id                          = cm_idle_req.amf_ue_ngap_id;
  notify_ue_event notify_ue_event_type;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "AMF_APP: ue_context is NULL\n");
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
      amf_smf_notification_send(ue_id, ue_context, notify_ue_event_type);
      ue_context_release_command(
          ue_id, ue_context->gnb_ue_ngap_id, NGAP_NAS_NORMAL_RELEASE);

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

/* Routine to send ue context release command to NGAP after processing
 * ue context release request from NGAP. this command will change ue
 * state to idle.
 */
void ue_context_release_command(
    amf_ue_ngap_id_t amf_ue_ngap_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    Ngcause ng_cause) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  itti_ngap_ue_context_release_command_t* ctx_rel_cmd = nullptr;
  MessageDef* message_p                               = nullptr;

  OAILOG_INFO(
      LOG_AMF_APP,
      "preparing for context release command to NGAP "
      "for ue_id %d\n",
      amf_ue_ngap_id);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);
  ctx_rel_cmd = &message_p->ittiMsg.ngap_ue_context_release_command;
  memset(ctx_rel_cmd, 0, sizeof(itti_ngap_ue_context_release_command_t));
  // Filling the respective values of NGAP message
  ctx_rel_cmd->amf_ue_ngap_id = amf_ue_ngap_id;
  ctx_rel_cmd->gnb_ue_ngap_id = gnb_ue_ngap_id;
  ctx_rel_cmd->cause          = ng_cause;
  // Send message to NGAP task
  OAILOG_INFO(LOG_AMF_APP, "sent context release command to NGAP\n");
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  OAILOG_INFO(
      LOG_AMF_APP,
      "sent context release command to NGAP "
      "for ue_id %d\n",
      amf_ue_ngap_id);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

static int paging_t3513_handler(zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(LOG_AMF_APP, "Timer: In Paging handler\n");
  OAILOG_INFO(LOG_AMF_APP, "Timer: identification T3513 handler \n");
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id;
  ue_m5gmm_context_s* ue_context                 = nullptr;
  amf_context_t* amf_ctx                         = nullptr;
  paging_context_t* paging_ctx                   = nullptr;
  MessageDef* message_p                          = nullptr;
  itti_ngap_paging_request_t* ngap_paging_notify = nullptr;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "ue_context is NULL\n");
    return -1;
  }

  // Get Paging Context
  amf_ctx    = &ue_context->amf_context;
  paging_ctx = &ue_context->paging_context;

  paging_ctx->m5_paging_response_timer.id = NAS5G_TIMER_INACTIVE_ID;
  /*
   * Increment the retransmission counter
   */
  paging_ctx->paging_retx_count += 1;
  OAILOG_ERROR(
      LOG_AMF_APP, "Timer: Incrementing retransmission_count to %d\n",
      paging_ctx->paging_retx_count);

  if (paging_ctx->paging_retx_count < MAX_PAGING_RETRY_COUNT) {
    /*
     * ReSend Paging request message to the UE
     */
    OAILOG_INFO(LOG_AMF_APP, "AMF_APP: In Handler Starting PAGING Timer\n");
    paging_ctx->m5_paging_response_timer.id = start_timer(
        &amf_app_task_zmq_ctx, PAGING_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE,
        paging_t3513_handler, NULL);
    OAILOG_INFO(LOG_AMF_APP, "AMF_APP: After Starting PAGING Timer\n");
    // Fill the itti msg based on context info produced in amf core

    message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_PAGING_REQUEST);

    ngap_paging_notify = &message_p->ittiMsg.ngap_paging_request;
    memset(ngap_paging_notify, 0, sizeof(itti_ngap_paging_request_t));
    ngap_paging_notify->UEPagingIdentity.amf_set_id =
        amf_ctx->m5_guti.guamfi.amf_set_id;
    ngap_paging_notify->UEPagingIdentity.amf_pointer =
        amf_ctx->m5_guti.guamfi.amf_pointer;
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_APP: Filling NGAP structure for Downlink amf_ctx dec "
        "m_tmsi=%d",
        amf_ctx->m5_guti.m_tmsi);
    ngap_paging_notify->UEPagingIdentity.m_tmsi = amf_ctx->m5_guti.m_tmsi;
    OAILOG_INFO(
        LOG_AMF_APP, "AMF_APP: Filling NGAP structure for Downlink m_tmsi=%d",
        ngap_paging_notify->UEPagingIdentity.m_tmsi);
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit1 =
        amf_ctx->m5_guti.guamfi.plmn.mcc_digit1;
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit2 =
        amf_ctx->m5_guti.guamfi.plmn.mcc_digit2;
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit3 =
        amf_ctx->m5_guti.guamfi.plmn.mcc_digit3;
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit1 =
        amf_ctx->m5_guti.guamfi.plmn.mnc_digit1;
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit2 =
        amf_ctx->m5_guti.guamfi.plmn.mnc_digit2;
    ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit3 =
        amf_ctx->m5_guti.guamfi.plmn.mnc_digit3;
    ngap_paging_notify->TAIListForPaging.no_of_items     = 1;
    ngap_paging_notify->TAIListForPaging.tai_list[0].tac = 2;

    OAILOG_INFO(LOG_AMF_APP, "AMF_APP: sending downlink message to NGAP");
    rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
    OAILOG_ERROR(
        LOG_AMF_APP, "Timer: timer has expired Sending Paging request again\n");
    //    amf_paging_request(paging_ctx);
  } else {
    /*
     * Abort the Paging procedure
     */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Timer: Maximum retires done hence Abort the Paging Request "
        "procedure\n");
    return rc;
  }
  return rc;
}

// Doing Paging Request handling received from SMF in AMF CORE
// int amf_app_defs::amf_app_handle_notification_received(
int amf_app_handle_notification_received(
    itti_n11_received_notification_t* notification) {
  ue_m5gmm_context_s* ue_context                 = nullptr;
  amf_context_t* amf_ctx                         = nullptr;
  paging_context_t* paging_ctx                   = nullptr;
  MessageDef* message_p                          = nullptr;
  itti_ngap_paging_request_t* ngap_paging_notify = nullptr;
  int rc                                         = RETURNok;

  OAILOG_INFO(LOG_AMF_APP, "AMF_APP: PAGING NOTIFICATION received from SMF\n");
  imsi64_t imsi64;
  IMSI_STRING_TO_IMSI64(notification->imsi, &imsi64);

  OAILOG_INFO(
      LOG_AMF_APP, "AMF_APP: IMSI is %s %lu\n", notification->imsi, imsi64);
  // Handle smf_context
  ue_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "ue_context is NULL\n");
    return -1;
  }

  OAILOG_INFO(
      LOG_AMF_APP, "AMF_APP: IMSI is %d\n", notification->notify_ue_evnt);
  switch (notification->notify_ue_evnt) {
    case UE_PAGING_NOTIFY:
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: PAGING NOTIFICATION received\n");

      // Get Paging Context
      amf_ctx    = &ue_context->amf_context;
      paging_ctx = &ue_context->paging_context;

      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: Starting PAGING Timer\n");
      /* Start Paging Timer T3513 */
      paging_ctx->m5_paging_response_timer.id = start_timer(
          &amf_app_task_zmq_ctx, PAGING_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE,
          paging_t3513_handler, NULL);
      // Fill the itti msg based on context info produced in amf core
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: After Starting PAGING Timer\n");
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: Allocating memory ");
      message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_PAGING_REQUEST);

      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: ngap_paging_notify");

      ngap_paging_notify = &message_p->ittiMsg.ngap_paging_request;
      memset(ngap_paging_notify, 0, sizeof(itti_ngap_paging_request_t));
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: Filling NGAP structure for Downlink");
      ngap_paging_notify->UEPagingIdentity.amf_set_id =
          amf_ctx->m5_guti.guamfi.amf_set_id;
      ngap_paging_notify->UEPagingIdentity.amf_pointer =
          amf_ctx->m5_guti.guamfi.amf_pointer;
      ngap_paging_notify->UEPagingIdentity.m_tmsi = amf_ctx->m5_guti.m_tmsi;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit1 =
          amf_ctx->m5_guti.guamfi.plmn.mcc_digit1;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit2 =
          amf_ctx->m5_guti.guamfi.plmn.mcc_digit2;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mcc_digit3 =
          amf_ctx->m5_guti.guamfi.plmn.mcc_digit3;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit1 =
          amf_ctx->m5_guti.guamfi.plmn.mnc_digit1;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit2 =
          amf_ctx->m5_guti.guamfi.plmn.mnc_digit2;
      ngap_paging_notify->TAIListForPaging.tai_list[0].plmn.mnc_digit3 =
          amf_ctx->m5_guti.guamfi.plmn.mnc_digit3;
      ngap_paging_notify->TAIListForPaging.no_of_items     = 1;
      ngap_paging_notify->TAIListForPaging.tai_list[0].tac = 2;
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP: sending downlink message to NGAP");
      rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
      break;

    case UE_SERVICE_REQUEST_ON_PAGING:
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_APP: SERVICE ACCEPT NOTIFICATION received\n");
      // TODO: Service Accept code to be implemented in upcoming PR
    default:
      OAILOG_INFO(LOG_AMF_APP, "AMF_APP : default case nothing to do\n");
      break;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
}  // namespace magma5g
