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
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/directoryd/directoryd.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_pdu_session_configs.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_session_manager_pco.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/include/n11_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
#include "lte/gateway/c/core/oai/include/map.h"

extern amf_config_t amf_config;
extern amf_config_t amf_config;
namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;

//------------------------------------------------------------------------------
void amf_ue_context_update_coll_keys(
    amf_ue_context_t* const amf_ue_context_p, ue_m5gmm_context_s* ue_context_p,
    const gnb_ngap_id_key_t gnb_ngap_id_key,
    const amf_ue_ngap_id_t amf_ue_ngap_id, const imsi64_t imsi,
    const teid_t amf_teid_n11, const guti_m5_t* const guti_p) {
  hashtable_rc_t h_rc  = HASH_TABLE_OK;
  magma::map_rc_t m_rc = magma::MAP_OK;
  // TODO: Migrate the amf_state_ue_id_ht to the map implementation.
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
    m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.remove(
        ue_context_p->gnb_ngap_id_key);
    m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.insert(
        gnb_ngap_id_key, amf_ue_ngap_id);

    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR_UE(
          LOG_AMF_APP, imsi,
          "Error could not update this ue context %p "
          "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT
          "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " %s\n",
          ue_context_p, ue_context_p->gnb_ue_ngap_id,
          ue_context_p->amf_ue_ngap_id, map_rc_code2string(m_rc).c_str());
    }
    ue_context_p->gnb_ngap_id_key = gnb_ngap_id_key;
  }

  if (amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    if (ue_context_p->amf_ue_ngap_id != amf_ue_ngap_id) {
      h_rc = hashtable_ts_remove(
          amf_state_ue_id_ht, (const hash_key_t) ue_context_p->amf_ue_ngap_id,
          reinterpret_cast<void**>(&ue_context_p));
      h_rc = hashtable_ts_insert(
          amf_state_ue_id_ht, (const hash_key_t) amf_ue_ngap_id,
          reinterpret_cast<void*>(ue_context_p));

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR(
            LOG_AMF_APP,
            "Insertion of Hash entry failed for  "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
            amf_ue_ngap_id);
      }

      ue_context_p->amf_ue_ngap_id = amf_ue_ngap_id;
    }
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP, "Invalid  amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
        amf_ue_ngap_id);
  }

  m_rc = amf_ue_context_p->imsi_amf_ue_id_htbl.remove(
      ue_context_p->amf_context.imsi64);

  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    m_rc = amf_ue_context_p->imsi_amf_ue_id_htbl.insert(imsi, amf_ue_ngap_id);
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Insertion of Hash entry failed for  "
        "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
        amf_ue_ngap_id);
  }

  m_rc = amf_ue_context_p->tun11_ue_context_htbl.remove(
      ue_context_p->amf_teid_n11);

  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    m_rc = amf_ue_context_p->tun11_ue_context_htbl.insert(
        amf_teid_n11, amf_ue_ngap_id);
  } else {
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
        (guti_p->guamfi.amf_regionid !=
         ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
        (guti_p->m_tmsi != ue_context_p->amf_context.m5_guti.m_tmsi) ||
        (guti_p->guamfi.plmn.mcc_digit1 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit1) ||
        (guti_p->guamfi.plmn.mcc_digit2 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
        (guti_p->guamfi.plmn.mcc_digit3 !=
         ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3) ||
        (ue_context_p->amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID)) {
      m_rc = amf_ue_context_p->guti_ue_context_htbl.remove(*guti_p);
      if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
        m_rc = amf_ue_context_p->guti_ue_context_htbl.insert(
            *guti_p, amf_ue_ngap_id);
      } else {
        OAILOG_ERROR(
            LOG_AMF_APP,
            "Insertion of Hash entry failed for  "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT PRIX32 " \n",
            amf_ue_ngap_id);
      }
      ue_context_p->amf_context.m5_guti = *guti_p;
    }
  }
}

/* Insert guti into guti_ue_context_table */
void amf_ue_context_on_new_guti(
    ue_m5gmm_context_t* const ue_context_p, const guti_m5_t* const guti_p) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  if (ue_context_p) {
    amf_ue_context_update_coll_keys(
        &amf_app_desc_p->amf_ue_contexts, ue_context_p,
        ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id,
        ue_context_p->amf_context.imsi64, ue_context_p->amf_teid_n11, guti_p);
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
  uint8_t num_amf             = 0;  // Number of configured AMF in the AMF pool
  guti_p->m_tmsi              = s_tmsi_p->m_tmsi;
  guti_p->guamfi.amf_set_id   = s_tmsi_p->amf_set_id;
  guti_p->guamfi.amf_pointer  = s_tmsi_p->amf_pointer;
  guti_p->guamfi.amf_regionid = amf_config.guamfi.guamfi[0].amf_regionid;
  // Create GUTI by using PLMN Id and AMF-Group Id of serving AMF
  OAILOG_DEBUG(
      LOG_AMF_APP,
      "Construct GUTI using S-TMSI received form UE and AMG set Id and pointer"
      "PLMN "
      "id from AMF Conf: %0x, %u %u\n",
      s_tmsi_p->m_tmsi, s_tmsi_p->amf_set_id, s_tmsi_p->amf_pointer);
  amf_config_read_lock(&amf_config);

  /*
   * Check number of MMEs in the pool.
   * At present it is assumed that one AMF is supported in AMF pool but in
   * case there are more than one AMF configured then search the serving AMF
   * using AMF code. Assumption is that within one PLMN only one pool of AMF
   * will be configured
   */
  if (amf_config.guamfi.nb > 1) {
    OAILOG_DEBUG(LOG_AMF_APP, "More than one AMFs are configured.");
  }
  for (num_amf = 0; num_amf < amf_config.guamfi.nb; num_amf++) {
    /*Verify that the AMF code within S-TMSI is same as what is configured in
     * AMF conf*/
    if ((plmn_p->mcc_digit2 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mcc_digit2) &&
        (plmn_p->mcc_digit1 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mcc_digit1) &&
        (plmn_p->mnc_digit3 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mnc_digit3) &&
        (plmn_p->mcc_digit3 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mcc_digit3) &&
        (plmn_p->mnc_digit2 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mnc_digit2) &&
        (plmn_p->mnc_digit1 ==
         amf_config.guamfi.guamfi[num_amf].plmn.mnc_digit1) &&
        (guti_p->guamfi.amf_set_id ==
         amf_config.guamfi.guamfi[num_amf].amf_set_id)) {
      break;
    }
  }
  if (num_amf >= amf_config.guamfi.nb) {
    OAILOG_DEBUG(LOG_AMF_APP, "No AMF serves this UE");
  } else {
    guti_p->guamfi.plmn        = amf_config.guamfi.guamfi[num_amf].plmn;
    guti_p->guamfi.amf_set_id  = amf_config.guamfi.guamfi[num_amf].amf_set_id;
    guti_p->guamfi.amf_pointer = amf_config.guamfi.guamfi[num_amf].amf_pointer;
    guti_p->guamfi.amf_regionid =
        amf_config.guamfi.guamfi[num_amf].amf_regionid;
    is_guti_valid = true;
  }

  amf_config_unlock(&amf_config);
  return is_guti_valid;
}

//------------------------------------------------------------------------------
// Get existing GUTI details
ue_m5gmm_context_s* amf_ue_context_exists_guti(
    amf_ue_context_t* const amf_ue_context_p, const guti_m5_t* const guti_p) {
  magma::map_rc_t m_rc           = magma::MAP_OK;
  uint64_t amf_ue_ngap_id64      = 0;
  ue_m5gmm_context_t* ue_context = NULL;

  m_rc = amf_ue_context_p->guti_ue_context_htbl.get(*guti_p, &amf_ue_ngap_id64);

  if (m_rc == magma::MAP_OK) {
    ue_context = amf_ue_context_exists_amf_ue_ngap_id(
        (amf_ue_ngap_id_t) amf_ue_ngap_id64);
    if (ue_context) {
      return ue_context;
    }
  }

  if (!ue_context) {
    OAILOG_WARNING(
        LOG_AMF_APP, "No GUTI hashtable for GUTI Hash %x", guti_p->m_tmsi);
    ue_context = ue_context_loopkup_by_guti(guti_p->m_tmsi);
    if (ue_context) {
      return ue_context;
    }
  }

  return NULL;
}

//-----------------------------------------------------------------------------------------
/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_initial_ue_message()                                **
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
  guti_m5_t guti                    = {0};
  plmn_t plmn                       = {0};
  s_tmsi_m5_t s_tmsi                = {0};
  amf_ue_ngap_id_t amf_ue_ngap_id   = INVALID_AMF_UE_NGAP_ID;

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
          ue_context_release_command(
              ue_context_p->amf_ue_ngap_id, ue_context_p->gnb_ue_ngap_id,
              NGAP_NAS_NORMAL_RELEASE);
          amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.remove(
              ue_context_p->gnb_ngap_id_key);
          ue_context_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
        }

        /* remove amf_ngap_ud_id entry from ue context */
        amf_remove_ue_context(ue_context_p);

        // Update AMF UE context with new gnb_ue_ngap_id
        ue_context_p->gnb_ue_ngap_id = initial_pP->gnb_ue_ngap_id;

        AMF_APP_GNB_NGAP_ID_KEY(
            gnb_ngap_id_key, initial_pP->gnb_id, initial_pP->gnb_ue_ngap_id);

        // generate new amf_ngap_ue_id
        amf_ue_ngap_id = amf_app_ctx_get_new_ue_id(
            &amf_app_desc_p->amf_app_ue_ngap_id_generator);

        amf_ue_context_update_coll_keys(
            &amf_app_desc_p->amf_ue_contexts, ue_context_p, gnb_ngap_id_key,
            amf_ue_ngap_id, ue_context_p->amf_context.imsi64,
            ue_context_p->amf_teid_n11, &guti);

        amf_insert_ue_context(ue_context_p->amf_ue_ngap_id, ue_context_p);
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

  /* Five_G_TMSI not configured */
  if (ue_context_p == NULL) {
    /* Check if Context can be found by GNB UE ID */
    ue_context_p = ue_context_lookup_by_gnb_ue_id(initial_pP->gnb_ue_ngap_id);

    /* Make sure its with same connection */
    if (ue_context_p &&
        (initial_pP->sctp_assoc_id != ue_context_p->sctp_assoc_id_key)) {
      ue_context_p = NULL;
    }
  }

  /*
   * UE Context already present in AMF. This can happen during PERIODIC
   * Registration. Steps for Periodic Registration:
   *   1. Context Setup and Guti established.
   *   2. UEContext Release Sequence
   *   3. InitialUEContextSetup Request from GNB with type as periodic
   *      registration
   */
  if (ue_context_p && ue_context_p->amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    /* If NGAP is not aware of ue_id or fiveGTmsi Is received
     * send the ue_id notification.
     */
    if ((initial_pP->amf_ue_ngap_id == INVALID_AMF_UE_NGAP_ID) ||
        (initial_pP->is_s_tmsi_valid)) {
      /* Sync data between AMF and NGAP */
      if (initial_pP->sctp_assoc_id == ue_context_p->sctp_assoc_id_key) {
        notify_ngap_new_ue_amf_ngap_id_association(ue_context_p);
      }
    }
  }

  // create a new ue context if nothing is found
  if (ue_context_p == NULL) {
    if (!(ue_context_p = amf_create_new_ue_context())) {
      OAILOG_ERROR(LOG_AMF_APP, "Failed to create ue_m5gmm_context for ue \n");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
    }

    // Allocate new amf_ue_ngap_id
    ue_context_p->amf_ue_ngap_id = amf_app_ctx_get_new_ue_id(
        &amf_app_desc_p->amf_app_ue_ngap_id_generator);

    if (ue_context_p->amf_ue_ngap_id == INVALID_AMF_UE_NGAP_ID) {
      OAILOG_ERROR(LOG_AMF_APP, "amf_ue_ngap_id allocation failed.\n");
      amf_remove_ue_context(ue_context_p);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
    }

    OAILOG_DEBUG(
        LOG_AMF_APP,
        "Creating new ue_m5gmm_context: [%p]"
        "for amf_ue_ngap_id: [" AMF_UE_NGAP_ID_FMT "]",
        ue_context_p, ue_context_p->amf_ue_ngap_id);

    AMF_APP_GNB_NGAP_ID_KEY(
        ue_context_p->gnb_ngap_id_key, initial_pP->gnb_id,
        initial_pP->gnb_ue_ngap_id);
    amf_insert_ue_context(ue_context_p->amf_ue_ngap_id, ue_context_p);
  }
  ue_context_p->sctp_assoc_id_key = initial_pP->sctp_assoc_id;
  ue_context_p->gnb_ue_ngap_id    = initial_pP->gnb_ue_ngap_id;

  // UEContextRequest
  ue_context_p->ue_context_request = initial_pP->ue_context_request;
  OAILOG_DEBUG(
      LOG_AMF_APP, "UE context request received: %d\n ",
      ue_context_p->ue_context_request);

  notify_ngap_new_ue_amf_ngap_id_association(ue_context_p);
  if (initial_pP->is_s_tmsi_valid) {
    s_tmsi = initial_pP->opt_s_tmsi;
  } else {
    s_tmsi.amf_pointer = 0;
    s_tmsi.m_tmsi      = INVALID_M_TMSI;
    OAILOG_DEBUG(
        LOG_AMF_APP,
        " Sending nas establishment indication to nas for ue_id = "
        "(" AMF_UE_NGAP_ID_FMT ")",
        ue_context_p->amf_ue_ngap_id);
  }
  is_mm_ctx_new = true;

  OAILOG_DEBUG(
      LOG_AMF_APP,
      " Sending NAS Establishment Indication to NAS for ue_id = "
      "(" AMF_UE_NGAP_ID_FMT ")",
      ue_context_p->amf_ue_ngap_id);

  amf_ue_ngap_id_t ue_id = ue_context_p->amf_ue_ngap_id;

  nas_proc_establish_ind(
      ue_context_p->amf_ue_ngap_id, is_mm_ctx_new, initial_pP->tai,
      initial_pP->ecgi, initial_pP->m5g_rrc_establishment_cause, s_tmsi,
      initial_pP->nas);

  initial_pP->nas = NULL;

  /* In case duplicate attach handling, ue_context_p might be removed
   * Before accessing ue_context_p, we shall validate whether UE context
   * exists or not
   */
  if (INVALID_AMF_UE_NGAP_ID != ue_id) {
    hash_table_ts_t* amf_state_ue_id_ht = get_amf_ue_state();
    if (hashtable_ts_is_key_exists(
            amf_state_ue_id_ht, (const hash_key_t) ue_id) == HASH_TABLE_OK) {
      imsi64 = ue_context_p->amf_context.imsi64;
    }
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
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
    amf_app_desc_t* amf_app_desc_p, bstring msg, amf_ue_ngap_id_t ue_id,
    const tai_t originating_tai) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  int rc = RETURNerror;
  if (msg) {
    amf_sap_t amf_sap = {};
    /*
     * Notify the AMF procedure call manager that data transfer
     * indication has been received from the Access-Stratum sublayer
     */
    amf_sap.primitive                    = AMFAS_ESTABLISH_REQ;
    amf_sap.u.amf_as.u.establish.ue_id   = ue_id;
    amf_sap.u.amf_as.u.establish.nas_msg = msg;
    amf_sap.u.amf_as.u.establish.tai     = originating_tai;
    msg                                  = NULL;
    rc                                   = amf_sap_send(&amf_sap);
  } else {
    OAILOG_WARNING(
        LOG_NAS,
        "Received NAS message in uplink is NULL for for UE "
        "ID: " AMF_UE_NGAP_ID_FMT,
        amf_app_desc_p->amf_app_ue_ngap_id_generator);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

static void get_ambr_unit(
    uint8_t apn_ambr_unit, uint32_t apn_session_ambr, uint8_t* calc_ambr_unit,
    uint16_t* calc_session_ambr) {
  *calc_ambr_unit = 0;

  while (apn_session_ambr) {
    *calc_ambr_unit += 1;

    apn_session_ambr >>= 2;
  }

  *calc_session_ambr = PDU_SESSION_DEFAULT_AMBR;
}

/* Received the session created response message from SMF. Populate and Send
 * PDU Session Resource Setup Request message to gNB and  PDU Session
 * Establishment Accept Message to UE*/
int amf_app_handle_pdu_session_response(
    itti_n11_create_pdu_session_response_t* pdu_session_resp) {
  DLNASTransportMsg encode_msg;
  ue_m5gmm_context_s* ue_context;
  std::shared_ptr<smf_context_t> smf_ctx;
  amf_smf_t amf_smf_msg;
  // TODO: hardcoded for now, addressed in the upcoming multi-UE PR
  uint64_t ue_id = 0;
  int rc         = RETURNerror;

  imsi64_t imsi64;
  IMSI_STRING_TO_IMSI64(pdu_session_resp->imsi, &imsi64);
  // Handle smf_context
  ue_context = lookup_ue_ctxt_by_imsi(imsi64);
  if (ue_context) {
    smf_ctx = amf_get_smf_context_by_pdu_session_id(
        ue_context, pdu_session_resp->pdu_session_id);
    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "pdu session  not found for session_id = %u\n",
          pdu_session_resp->pdu_session_id);
      return RETURNerror;
    }
    ue_id = ue_context->amf_ue_ngap_id;
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP, "ue context not found for the imsi=%lu\n", imsi64);
    return RETURNerror;
  }

  get_ambr_unit(
      pdu_session_resp->session_ambr.downlink_unit_type,
      pdu_session_resp->session_ambr.downlink_units, &smf_ctx->dl_ambr_unit,
      &smf_ctx->dl_session_ambr);

  get_ambr_unit(
      pdu_session_resp->session_ambr.uplink_unit_type,
      pdu_session_resp->session_ambr.uplink_units, &smf_ctx->ul_ambr_unit,
      &smf_ctx->ul_session_ambr);

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

  smf_ctx->n_active_pdus += 1;

  if (REGISTERED_IDLE == ue_context->mm_state) {
    // pdu session state
    smf_ctx->pdu_session_state            = ACTIVE;
    amf_sap_t amf_sap                     = {};
    amf_sap.primitive                     = AMFAS_ESTABLISH_CNF;
    amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
    amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_SR;

    amf_sap.u.amf_as.u.establish.pdu_session_status_ie =
        (AMF_AS_PDU_SESSION_STATUS | AMF_AS_PDU_SESSION_REACTIVATION_STATUS);
    amf_sap.u.amf_as.u.establish.pdu_session_status =
        (1 << smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id);
    amf_sap.u.amf_as.u.establish.pdu_session_reactivation_status =
        (1 << smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id);
    amf_sap.u.amf_as.u.establish.guti = ue_context->amf_context.m5_guti;
    rc                                = amf_sap_send(&amf_sap);
    if (RETURNok == rc) {
      ue_context->mm_state = REGISTERED_CONNECTED;
    }
  } else {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "Sending message to gNB for PDUSessionResourceSetupRequest "
        "**n_active_pdus=%d **\n",
        smf_ctx->n_active_pdus);

    amf_smf_msg.pdu_session_id = pdu_session_resp->pdu_session_id;
    /*Execute PDU establishement accept from AMF to gnodeb */
    pdu_state_handle_message(
        // ue_context->mm_state, STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT,
        REGISTERED_CONNECTED, STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT,
        // smf_ctx->pdu_session_state, ue_context, amf_smf_msg, NULL,
        CREATING, ue_context, amf_smf_msg, NULL, pdu_session_resp, ue_id);
    rc = RETURNok;
  }

  return rc;
}
/****************************************************************************
 **                                                                        **
 ** Name:    convert_ambr()                                                **
 **                                                                        **
 ** Description: Converts the session ambr format from                     **
 **  one defined in create_pdu_session_response to one defined in          **
 **  pdu_session_estab_accept message                                      **
 **                                                                        **
 ** Inputs:  pdu_ambr_response_unit, pdu_ambr_response_value               **
 **          ambr_unit, ambr_value                                         **
 **                                                                        **
 **      Return:   void                                                    **
 **                                                                        **
 ***************************************************************************/
void convert_ambr(
    const uint32_t* pdu_ambr_response_unit,
    const uint32_t* pdu_ambr_response_value, uint8_t* ambr_unit,
    uint16_t* ambr_value) {
  int count                             = 1;
  uint32_t temp_pdu_ambr_response_value = *pdu_ambr_response_value;

  // minimum rate unit is KBPS
  if (*pdu_ambr_response_unit == BPS &&
      temp_pdu_ambr_response_value / 1000 == 0) {
    // Values less than 1Kbps are defaulted to 1Kbps
    *ambr_value = static_cast<uint16_t>(1);
    *ambr_unit  = static_cast<uint8_t>(
        magma5g::M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // Kbps
    return;
  }

  if (*pdu_ambr_response_unit == BPS) {
    temp_pdu_ambr_response_value /= 1000;
  }

  while (temp_pdu_ambr_response_value >= AMBR_UNIT_CONVERT_THRESHOLD) {
    temp_pdu_ambr_response_value = (temp_pdu_ambr_response_value / 1000);
    count++;
  }

  switch (count) {
    case 1:
      *ambr_unit = static_cast<uint8_t>(
          magma5g::M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // Kbps
      break;
    case 2:
      *ambr_unit = static_cast<uint8_t>(
          magma5g::M5GSessionAmbrUnit::MULTIPLES_1MBPS);  // Mbps
      break;
    case 3:
      *ambr_unit = static_cast<uint8_t>(
          magma5g::M5GSessionAmbrUnit::MULTIPLES_1GBPS);  // Gbps
      break;
  }
  *ambr_value = static_cast<uint16_t>(temp_pdu_ambr_response_value);
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
    itti_n11_create_pdu_session_response_t* pdu_session_resp, uint64_t ue_id) {
  nas5g_error_code_t rc = M5G_AS_SUCCESS;

  DLNASTransportMsg* encode_msg;
  amf_nas_message_t msg = {};
  uint32_t bytes        = 0;
  uint32_t len          = 0;
  SmfMsg* smf_msg       = nullptr;
  bstring buffer;
  // smf_ctx declared and set but not used, commented to cleanup warnings
  std::shared_ptr<smf_context_t> smf_ctx;
  ue_m5gmm_context_s* ue_context                   = nullptr;
  protocol_configuration_options_t* msg_accept_pco = nullptr;

  // Handle smf_context
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP, "ue context not found for the ue_id:" AMF_UE_NGAP_ID_FMT,
        ue_id);
    return M5G_AS_FAILURE;
  }

  smf_ctx = amf_get_smf_context_by_pdu_session_id(
      ue_context, pdu_session_resp->pdu_session_id);
  if (!smf_ctx) {
    OAILOG_ERROR(
        LOG_AMF_APP, "Smf context is not exist UE ID:" AMF_UE_NGAP_ID_FMT,
        ue_id);
    return M5G_AS_FAILURE;
  }
  // updating session state
  smf_ctx->pdu_session_state = ACTIVE;

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
  encode_msg->payload_container_type.iei = 0;
  // encode_msg->payload_container_type.iei = PAYLOAD_CONTAINER_TYPE;

  // SmfMsg
  encode_msg->payload_container_type.type_val = N1_SM_INFO;
  encode_msg->payload_container.iei           = PAYLOAD_CONTAINER;
  encode_msg->pdu_session_identity.iei =
      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2);
  encode_msg->pdu_session_identity.pdu_session_id =
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
  smf_msg->msg.pdu_session_estab_accept.pdu_session_type.type_val = 1;
  smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val =
      (pdu_session_resp->selected_ssc_mode + 1);

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
  qos_rule.qfi =
      smf_ctx->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item.qos_flow_identifier;
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
  convert_ambr(
      &pdu_session_resp->session_ambr.downlink_unit_type,
      &pdu_session_resp->session_ambr.downlink_units,
      &smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_unit,
      &smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_session_ambr);

  convert_ambr(
      &pdu_session_resp->session_ambr.uplink_unit_type,
      &pdu_session_resp->session_ambr.uplink_units,
      &smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_unit,
      &smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_session_ambr);
  smf_msg->msg.pdu_session_estab_accept.session_ambr.length = AMBR_LEN;

  msg_accept_pco =
      &(smf_msg->msg.pdu_session_estab_accept.protocolconfigurationoptions.pco);

  auto pco_len = sm_process_pco_request(&(smf_ctx->pco), msg_accept_pco);

  /* NSSAI
  --------------------------------------
  Parameters | IEI | Length | SST | SD |
  --------------------------------------
  Size       | 1   | 1      | 1   | 3  |
  -------------------------------------- */
  smf_msg->msg.pdu_session_estab_accept.nssai.iei =
      static_cast<uint8_t>(M5GIei::S_NSSAI);
  uint32_t buf_len = 0;
  if (smf_ctx->sst) {
    if (smf_ctx->sd[0]) {
      smf_msg->msg.pdu_session_estab_accept.nssai.len = SST_LENGTH + SD_LENGTH;
      smf_msg->msg.pdu_session_estab_accept.nssai.sst = smf_ctx->sst;
      memcpy(smf_msg->msg.pdu_session_estab_accept.nssai.sd, smf_ctx->sd, 3);
    } else {
      smf_msg->msg.pdu_session_estab_accept.nssai.len = SST_LENGTH;
      smf_msg->msg.pdu_session_estab_accept.nssai.sst = smf_ctx->sst;
    }
    buf_len = smf_msg->msg.pdu_session_estab_accept.nssai.len + 2;
  }

  /* DNN
  -------------------------------------
  Parameters | IEI | Length | DNN     |
  -------------------------------------
  Size       | 1   | 1      | 1 - 100 |
  ------------------------------------- */
  smf_msg->msg.pdu_session_estab_accept.dnn.iei =
      static_cast<uint8_t>(M5GIei::DNN);
  smf_msg->msg.pdu_session_estab_accept.dnn.len = smf_ctx->dnn.length() + 1;
  smf_ctx->dnn.copy(
      reinterpret_cast<char*>(smf_msg->msg.pdu_session_estab_accept.dnn.dnn),
      smf_ctx->dnn.length());
  buf_len += smf_msg->msg.pdu_session_estab_accept.dnn.len + 2;

  encode_msg->payload_container.len =
      PDU_ESTAB_ACCPET_PAYLOAD_CONTAINER_LEN + pco_len + buf_len;
  len = PDU_ESTAB_ACCEPT_NAS_PDU_LEN + pco_len + buf_len;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);
  if (msg.header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    amf_msg_header* header = &msg.security_protected.plain.amf.header;
    /*
     * Expand size of protected NAS message
     */
    len += NAS_MESSAGE_SECURITY_HEADER_SIZE;
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
    buffer->slen = bytes;
    if (ue_context->ue_context_request) {
      pdu_session_resource_setup_request(ue_context, ue_id, smf_ctx, buffer);
    } else {
      initial_context_setup_request(ue_id, &ue_context->amf_context, buffer);
    }

  } else {
    OAILOG_WARNING(LOG_AMF_APP, "NAS encode failed \n");
    bdestroy_wrapper(&buffer);
  }

  /* Clean up the PCO contents */
  sm_free_protocol_configuration_options(&msg_accept_pco);

  return rc;
}

/* Handling PDU Session Resource Setup Response sent from gNB*/
void amf_app_handle_resource_setup_response(
    itti_ngap_pdusessionresource_setup_rsp_t session_seup_resp) {
  amf_ue_ngap_id_t ue_id;

  ue_m5gmm_context_s* ue_context = nullptr;
  std::shared_ptr<smf_context_t> smf_ctx;

  /* Check if failure message is not NULL and if NULL,
   * it is successful message from gNB.
   * Nothing to in this case. If failure message comes from gNB
   * AMF need to report this failed message to SMF
   *
   * NOTE: only handling success part not failure part
   * will be handled later
   */
  OAILOG_DEBUG(
      LOG_AMF_APP, "Handling uplink PDU session setup response message\n");

  if (session_seup_resp.pduSessionResource_setup_list.no_of_items > 0) {
    ue_id = session_seup_resp.amf_ue_ngap_id;

    ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
    if (ue_context == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "UE context not found for the ue_id = " AMF_UE_NGAP_ID_FMT, ue_id);
      return;
    }

    smf_ctx = amf_get_smf_context_by_pdu_session_id(
        ue_context,
        session_seup_resp.pduSessionResource_setup_list.item[0].Pdu_Session_ID);
    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "PDU session  not found for session_id = %lu\n",
          session_seup_resp.pduSessionResource_setup_list.item[0]
              .Pdu_Session_ID);
      return;
    }

    /* This is success case and we need not to send message to SMF
     * and drop the message here
     */
    amf_ue_ngap_id_t ue_id;
    amf_smf_establish_t amf_smf_grpc_ies;
    ue_m5gmm_context_s* ue_context = nullptr;
    char imsi[IMSI_BCD_DIGITS_MAX + 1];

    ue_id = session_seup_resp.amf_ue_ngap_id;

    ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
    // Handling of ue context
    if (!ue_context) {
      OAILOG_ERROR(
          LOG_AMF_APP, "UE context not found for UE ID: " AMF_UE_NGAP_ID_FMT,
          ue_id);
    }

    // Store gNB ip and TEID in respective smf_context
    memset(
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0',
        sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr));

    smf_ctx->gtp_tunnel_id
        .gnb_gtp_teid = htonl(*(reinterpret_cast<unsigned int*>(
        session_seup_resp.pduSessionResource_setup_list.item[0]
            .PDU_Session_Resource_Setup_Response_Transfer.tunnel.gTP_TEID)));
    memcpy(
        &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr,
        &session_seup_resp.pduSessionResource_setup_list.item[0]
             .PDU_Session_Resource_Setup_Response_Transfer.tunnel
             .transportLayerAddress,
        4);  // time being 4 byte is copying.
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "gnb_gtp_teid_ipaddr: [%02x %02x %02x %02x]  and gnb_gtp_teid "
        "[" GNB_GTP_TEID_FMT "]\n",
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[0],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[1],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[2],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr[3],
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid);
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
 * - Retrieve the required field of UE, like IMSI and fill gRPC notification
 *   proto structure.
 * - In AMF move the UE/IMSI state to CM-idle
 *   Go over all PDU sessions and change the state to in-active.
 * */
void amf_app_handle_cm_idle_on_ue_context_release(
    itti_ngap_ue_context_release_req_t cm_idle_req) {
  int rc = RETURNerror;
  OAILOG_DEBUG(
      LOG_AMF_APP,
      " Handling UL UE context release for CM-idle for UE "
      "ID: " AMF_UE_NGAP_ID_FMT,
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
  ue_id                          = cm_idle_req.amf_ue_ngap_id;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context) {
    return;
  }

  // if UE on REGISTERED_IDLE, so no need to do anyting
  if ((ue_context->mm_state == REGISTERED_CONNECTED) ||
      (ue_context->mm_state == DEREGISTERED)) {
    // UE in connected state and need to check if cause is proper
    if (cm_idle_req.relCause == NGAP_RADIO_NR_GENERATED_REASON) {
      rc = ue_state_handle_message_initial(
          ue_context->mm_state, STATE_EVENT_CONTEXT_RELEASE, SESSION_NULL,
          ue_context, &ue_context->amf_context);

      if (rc != RETURNok) {
        OAILOG_WARNING(LOG_AMF_APP, "Failed transitioning to idle mode\n");
      }

      ue_context_release_command(
          ue_id, ue_context->gnb_ue_ngap_id, NGAP_USER_INACTIVITY);
      ue_context->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;

    } else {
      OAILOG_WARNING(
          LOG_AMF_APP,
          " UE in registered_connected state, but cause from NGAP"
          " is wrong for UE ID: " AMF_UE_NGAP_ID_FMT " and return",
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
        " for UE ID: " AMF_UE_NGAP_ID_FMT,
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

  OAILOG_DEBUG(
      LOG_AMF_APP,
      "preparing for context release command to NGAP "
      "for ue_id " AMF_UE_NGAP_ID_FMT,
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
  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

static int paging_t3513_handler(zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(LOG_AMF_APP, "T3513: In Paging handler\n");
  int rc                                         = RETURNerror;
  amf_ue_ngap_id_t ue_id                         = 0;
  ue_m5gmm_context_s* ue_context                 = nullptr;
  amf_context_t* amf_ctx                         = nullptr;
  paging_context_t* paging_ctx                   = nullptr;
  MessageDef* message_p                          = nullptr;
  itti_ngap_paging_request_t* ngap_paging_notify = nullptr;

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3513: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "ue_context is NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  // Get Paging Context
  amf_ctx    = &ue_context->amf_context;
  paging_ctx = &ue_context->paging_context;

  paging_ctx->m5_paging_response_timer.id = NAS5G_TIMER_INACTIVE_ID;

  if (paging_ctx->paging_retx_count < MAX_PAGING_RETRY_COUNT) {
    paging_ctx->paging_retx_count += 1;

    OAILOG_DEBUG(
        LOG_AMF_APP,
        "T3513: timer has expired for UE ID: " AMF_UE_NGAP_ID_FMT
        " with timer id: %d, Sending Paging request again",
        ue_id, timer_id);
    /*
     * Increment the retransmission counter
     */
    OAILOG_DEBUG(
        LOG_AMF_APP, "T3513: Incrementing retransmission_count to %d\n",
        paging_ctx->paging_retx_count);

    /*
     * ReSend Paging request message to the UE
     */

    // Fill the itti msg based on context info produced in amf core

    message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_PAGING_REQUEST);

    ngap_paging_notify = &message_p->ittiMsg.ngap_paging_request;
    memset(ngap_paging_notify, 0, sizeof(itti_ngap_paging_request_t));
    ngap_paging_notify->UEPagingIdentity.amf_set_id =
        amf_ctx->m5_guti.guamfi.amf_set_id;
    ngap_paging_notify->UEPagingIdentity.amf_pointer =
        amf_ctx->m5_guti.guamfi.amf_pointer;
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "T3513: Filling NGAP structure for Downlink amf_ctx dec "
        "m_tmsi=%d",
        amf_ctx->m5_guti.m_tmsi);
    ngap_paging_notify->UEPagingIdentity.m_tmsi = amf_ctx->m5_guti.m_tmsi;
    OAILOG_INFO(
        LOG_AMF_APP, "T3513: Filling NGAP structure for Downlink m_tmsi=%d",
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

    OAILOG_INFO(LOG_AMF_APP, "T3513: sending downlink message to NGAP");
    rc = send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
    if (rc != RETURNok)
      OAILOG_ERROR(LOG_AMF_APP, "Could not send msg to task\n");
    //    amf_paging_request(paging_ctx);
  } else {
    /*
     * Abort the Paging procedure
     */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3513: Maximum retires done hence Abort the Paging Request "
        "procedure\n");
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
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

  imsi64_t imsi64;
  IMSI_STRING_TO_IMSI64(notification->imsi, &imsi64);

  OAILOG_DEBUG(LOG_AMF_APP, "IMSI is %s %lu\n", notification->imsi, imsi64);
  // Handle smf_context
  ue_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (!ue_context) {
    OAILOG_ERROR(LOG_AMF_APP, "UE context is null\n");
    return -1;
  }

  switch (notification->notify_ue_evnt) {
    case UE_PAGING_NOTIFY:
      OAILOG_DEBUG(LOG_AMF_APP, "Paging notification received\n");

      // Get Paging Context
      amf_ctx    = &ue_context->amf_context;
      paging_ctx = &ue_context->paging_context;

      OAILOG_INFO(
          LOG_AMF_APP,
          "T3513: Starting PAGING Timer for UE ID: " AMF_UE_NGAP_ID_FMT,
          ue_context->amf_ue_ngap_id);
      paging_ctx->paging_retx_count = 0;
      /* Start Paging Timer T3513 */
      paging_ctx->m5_paging_response_timer.id = amf_app_start_timer(
          PAGING_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE, paging_t3513_handler,
          ue_context->amf_ue_ngap_id);
      // Fill the itti msg based on context info produced in amf core
      OAILOG_INFO(
          LOG_AMF_APP,
          "T3513: Starting PAGING Timer for UE ID: " AMF_UE_NGAP_ID_FMT
          " and timer id: %ld",
          ue_context->amf_ue_ngap_id, paging_ctx->m5_paging_response_timer.id);

      message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_PAGING_REQUEST);

      ngap_paging_notify = &message_p->ittiMsg.ngap_paging_request;
      memset(ngap_paging_notify, 0, sizeof(itti_ngap_paging_request_t));
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
      OAILOG_DEBUG(LOG_AMF_APP, "Service Accept notification received\n");
      // TODO: Service Accept code to be implemented in upcoming PR
    default:
      OAILOG_DEBUG(LOG_AMF_APP, "default case nothing to do\n");
      break;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

void amf_app_handle_initial_context_setup_rsp(
    amf_app_desc_t* amf_app_desc_p,
    itti_amf_app_initial_context_setup_rsp_t* initial_context_rsp) {
  ue_m5gmm_context_s* ue_context = NULL;
  std::shared_ptr<smf_context_t> smf_context;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  Ngap_PDUSession_Resource_Setup_Response_List_t* pdu_list =
      &initial_context_rsp->PDU_Session_Resource_Setup_Response_Transfer;

  // Handle smf_context
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(initial_context_rsp->ue_id);

  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP, " Ue context not found for the id " AMF_UE_NGAP_ID_FMT,
        initial_context_rsp->ue_id);
    return;
  }

  if (pdu_list->no_of_items) {
    for (uint32_t index = 0; index < pdu_list->no_of_items; index++) {
      smf_context = amf_get_smf_context_by_pdu_session_id(
          ue_context, pdu_list->item[index].Pdu_Session_ID);
      if (smf_context == NULL) {
        OAILOG_ERROR(
            LOG_AMF_APP, "pdu session  not found for session_id = %ld\n",
            pdu_list->item[index].Pdu_Session_ID);
      } else {
        amf_smf_establish_t amf_smf_grpc_ies;

        // gnb tunnel info

        smf_context->gtp_tunnel_id.gnb_gtp_teid =
            htonl(*(reinterpret_cast<unsigned int*>(
                pdu_list->item[index]
                    .PDU_Session_Resource_Setup_Response_Transfer.tunnel
                    .gTP_TEID)));
        memcpy(
            smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr,
            pdu_list->item[index]
                .PDU_Session_Resource_Setup_Response_Transfer.tunnel
                .transportLayerAddress,
            4);

        OAILOG_DEBUG(
            LOG_AMF_APP,
            "IP address %02x %02x %02x %02x  and TEID" GNB_GTP_TEID_FMT "\n",
            smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr[0],
            smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr[1],
            smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr[2],
            smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr[3],
            smf_context->gtp_tunnel_id.gnb_gtp_teid);

        smf_context->pdu_session_version++;
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
            &smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
        memcpy(
            &amf_smf_grpc_ies.gnb_gtp_teid,
            &smf_context->gtp_tunnel_id.gnb_gtp_teid, 4);
        amf_smf_grpc_ies.pdu_session_id = pdu_list->item[index].Pdu_Session_ID;

        IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
        /* Prepare and send gNB setup response message to SMF through gRPC
         * 2nd time PDU session establish message
         */
        create_session_grpc_req_on_gnb_setup_rsp(
            &amf_smf_grpc_ies, imsi, smf_context->pdu_session_version);
      }
    }
  }

  // update UE state
  ue_context->mm_state = REGISTERED_CONNECTED;
}
}  // namespace magma5g
