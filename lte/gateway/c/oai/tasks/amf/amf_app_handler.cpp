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

  Source      amf_app_handler.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_config.h"
//#include "amf_app_ue_context.h"
#include "nas_proc.h"
#include "3gpp_23.003.h"
#include "amf_app_msg.h"
using namespace std;

namespace magma5g
{
    class amf_app_handler:public amf_app_ue_context: public amf_app_desc_t
    {
        public:
        
    };
//----------------------------------------------------------------------------
static void amf_directoryd_report_location(uint64_t imsi, uint8_t imsi_len) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  directoryd_report_location(imsi_str);
  OAILOG_INFO_UE(LOG_AMF_APP, imsi, "Reported UE location to directoryd\n");
}
//------------------------------------------------------------------------------
void amf_ue_context_update_coll_keys(
    amf_ue_context_t* const amf_ue_context_p,
    ue_m5gmm_context_s* const ue_context_p,
    const gnb_ngap_id_key_t gnb_ngap_id_key,
    const amf_ue_ngap_id_t amf_ue_ngap_id, const imsi64_t imsi,
    const n11_teid_t amf_teid_n11,
    const guti_m5_t* const guti_p)  //  never NULL, if none put &ue_context_p->guti
{
  hashtable_rc_t h_rc                 = HASH_TABLE_OK;
  hash_table_ts_t* amf_state_ue_id_ht = get_amf_ue_state();
  OAILOG_FUNC_IN(LOG_AMF_APP);

  OAILOG_TRACE(LOG_AMF_APP,"Update ue context.old_gnb_ue_ngap_id_key %ld ue "
      "context.old_amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
      " ue context.old_IMSI " IMSI_64_FMT " ue context.old_GUTI " GUTI_FMT "\n",
      ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id,
      ue_context_p->amf_context._imsi64,
      GUTI_ARG(&ue_context_p->amf_context._guti));

  OAILOG_TRACE(LOG_AMF_APP,"Update ue context %p updated_gnb_ue_ngap_id_key %ld "
      "updated_amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " updated_IMSI " IMSI_64_FMT
      " updated_GUTI " GUTI_FMT "\n",
      ue_context_p, gnb_ngap_id_key, amf_ue_ngap_id, imsi, GUTI_ARG(guti_p));

  if ((INVALID_GNB_UE_NGAP_ID_KEY != gnb_ngap_id_key) &&
      (ue_context_p->gnb_ngap_id_key != gnb_ngap_id_key)) {
    // new insertion of gnb_ue_ngap_id_key,
    h_rc = hashtable_uint64_ts_remove(amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,(const hash_key_t) ue_context_p->gnb_ngap_id_key);
    h_rc = hashtable_uint64_ts_insert(amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,(const hash_key_t) gnb_ngap_id_key, amf_ue_ngap_id);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_ERROR_UE(LOG_AMF_APP, imsi,"Error could not update this ue context %p "
          "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " %s\n",
          ue_context_p, ue_context_p->gnb_ue_ngap_id,ue_context_p->amf_ue_ngap_id, hashtable_rc_code2string(h_rc));
    }
    ue_context_p->gnb_ngap_id_key = gnb_ngap_id_key;
  } else {
    OAILOG_DEBUG_UE(LOG_AMF_APP, imsi, "Did not update gnb_ngap_id_key %ld in ue context %p "
        "gnb_ue_ngap_ue_id "GENB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "\n",
        gnb_ngap_id_key, ue_context_p, ue_context_p->gnb_ue_ngap_id,ue_context_p->amf_ue_ngap_id);
  }

  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    if (ue_context_p->amf_ue_ngap_id != amf_ue_ngap_id) {
      // new insertion of amf_ue_ngap_id, not a change in the id
      h_rc = hashtable_ts_remove(amf_state_ue_id_ht, (const hash_key_t) ue_context_p->amf_ue_ngap_id,(void**) &ue_context_p);
      h_rc = hashtable_ts_insert(amf_state_ue_id_ht, (const hash_key_t) amf_ue_ngap_id,(void*) ue_context_p);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR_UE(LOG_AMF_APP, imsi,"Error could not update this ue context %p "
            "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " %s\n",
            ue_context_p, ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id, hashtable_rc_code2string(h_rc));
      }
      ue_context_p->amf_ue_ngap_id = amf_ue_ngap_id;
    }
  } else {
    OAILOG_DEBUG_UE(LOG_AMF_APP, imsi,  "Did not update hashtable  for ue context %p "
        "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " imsi " IMSI_64_FMT " \n",
        ue_context_p, ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id, imsi);
  }

  h_rc = hashtable_uint64_ts_remove(amf_ue_context_p->imsi_amf_ue_id_htbl,
      (const hash_key_t) ue_context_p->amf_context._imsi64);
  if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
    h_rc = hashtable_uint64_ts_insert(amf_ue_context_p->imsi_amf_ue_id_htbl, (const hash_key_t) imsi, amf_ue_ngap_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }
  if (HASH_TABLE_OK != h_rc) {
    OAILOG_ERROR_UE(LOG_AMF_APP, imsi, "Error could not update this ue context %p "
        "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " imsi " IMSI_64_FMT ": %s\n",
        ue_context_p, ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id, imsi, hashtable_rc_code2string(h_rc));
  }
     amf_directoryd_report_location(ue_context_p->amf_context._imsi64, ue_context_p->amf_context._imsi.length);

  h_rc = hashtable_uint64_ts_remove(amf_ue_context_p->tun11_ue_context_htbl,
      (const hash_key_t) ue_context_p->amf_teid_n11);
  if (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) {
    h_rc = hashtable_uint64_ts_insert( amf_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) amf_teid_n11, (uint64_t) amf_ue_ngap_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }

  if (HASH_TABLE_OK != h_rc) {
    OAILOG_ERROR_UE(LOG_AMF_APP, imsi, "Error could not update this ue context %p "
        "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " amf_teid_n11 " TEID_FMT
        " : %s\n",
        ue_context_p, ue_context_p->gnb_ue_ngap_id,
        ue_context_p->amf_ue_ngap_id, amf_teid_n11,
        hashtable_rc_code2string(h_rc));
  }
  ue_context_p->amf_teid_n11 = amf_teid_n11;

  if (guti_p) {
    if ((guti_p->guamfi.amf_code !=
         ue_context_p->amf_context._m5_guti.guamfi.amf_code) ||
        (guti_p->guamfi.amf_gid !=
         ue_context_p->amf_context._m5_guti.guamfi.amf_gid) ||
        (guti_p->m_tmsi != ue_context_p->amf_context._m5_guti.m_tmsi) ||
        (guti_p->guamfi.plmn.mcc_digit1 !=
         ue_context_p->amf_context._m5_guti.guamfi.plmn.mcc_digit1) ||
        (guti_p->guamfi.plmn.mcc_digit2 !=
         ue_context_p->amf_context._m5_guti.guamfi.plmn.mcc_digit2) ||
        (guti_p->guamfi.plmn.mcc_digit3 !=
         ue_context_p->amf_context._m5_guti.guamfi.plmn.mcc_digit3) ||
        (ue_context_p->amf_ue_ngap_id != amf_ue_ngap_id)) {
      // may check guti_p with a kind of instanceof()?
      h_rc = obj_hashtable_uint64_ts_remove(
          amf_ue_context_p->guti_ue_context_htbl,
          &ue_context_p->amf_context._m5_guti, sizeof(*guti_p));
      if (INVALID_AMF_UE_NGAP_ID != amf_ue_ngap_id) {
        h_rc = obj_hashtable_uint64_ts_insert(
            amf_ue_context_p->guti_ue_context_htbl, (const void* const) guti_p,
            sizeof(*guti_p), (uint64_t) amf_ue_ngap_id);
      } else {
        h_rc = HASH_TABLE_KEY_NOT_EXISTS;
      }

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR_UE(LOG_AMF_APP, imsi, "Error could not update this ue context %p "
            "gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " guti " GUTI_FMT " %s\n",
            ue_context_p, ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id, GUTI_ARG(guti_p),
            hashtable_rc_code2string(h_rc));
      }
      ue_context_p->amf_context._m5_guti = *guti_p;
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
//----------------------------------------------------------------------------------------------
    static bool amf_app_construct_guti(const plmn_t* const plmn_p, const s_tmsi_m5_t* const s_tmsi_p,
            guti_m5_t* const guti_p) {
        /*
        * This is a helper function to construct GUTI from S-TMSI. It uses PLMN id
        * and AMF Group Id of the serving AMF for this purpose.
        *
        */

        bool is_guti_valid =  false;  // Set to true if serving AMF is found and GUTI is constructed
        uint8_t num_amf         = 0;  // Number of configured AMF in the AMF pool
        guti_p->m_tmsi          = s_tmsi_p->m_tmsi;
        guti_p-> guamfi_t.amf_code = s_tmsi_p->amf_code;
        // Create GUTI by using PLMN Id and AMF-Group Id of serving AMF
        OAILOG_DEBUG(LOG_AMF_APP,"Construct GUTI using S-TMSI received form UE and AMG Group Id and PLMN "
            "id "
            "from AMF Conf: %u, %u \n",
            s_tmsi_p->m_tmsi, s_tmsi_p->amf_code);
        amf_config_read_lock(&amf_config);
        /*
        * Check number of MMEs in the pool.
        * At present it is assumed that one AMF is supported in AMF pool but in case
        * there are more than one AMF configured then search the serving AMF using
        * AMF code. Assumption is that within one PLMN only one pool of AMF will be
        * configured
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
                (guti_p->guamfi.amf_code ==
                amf_config.guamfi.guamfi[num_amf].amf_code)) {
            break;
            }
        }
        if (num_amf >= amf_config.guamfi.nb) {
            OAILOG_DEBUG(LOG_AMF_APP, "No AMF serves this UE");
        } else {
            guti_p->guamfi.plmn    = amf_config.guamfi.guamfi[num_amf].plmn;
            guti_p->guamfi.amf_gid = amf_config.guamfi.guamfi[num_amf].amf_gid;
            is_guti_valid          = true;
        }
        amf_config_unlock(&amf_config);
        return is_guti_valid;
    }
    //------------------------------------------------------------------------------
    ue_m5gmm_context_s* amf_ue_context_exists_guti(amf_ue_context_t* const amf_ue_context_p, const guti_m5_t* const guti_p) 
    {

        hashtable_rc_t h_rc       = HASH_TABLE_OK;
        uint64_t amf_ue_ngap_id64 = 0;

        h_rc = obj_hashtable_uint64_ts_get(
            amf_ue_context_p->guti_ue_context_htbl, (const void*) guti_p,
            sizeof(*guti_p), &amf_ue_ngap_id64);

        if (HASH_TABLE_OK == h_rc) {
            return amf_ue_context_exists_amf_ue_ngap_id(
                (amf_ue_ngap_id_t) amf_ue_ngap_id64);
        } else {
            OAILOG_WARNING(LOG_AMF_APP, " No GUTI hashtable for GUTI ");
        }

        return NULL;
    }

//------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------
           
    imsi64_t amf_app_defs::amf_app_handle_initial_ue_message(amf_app_desc_t *amf_app_desc_p, itti_ngap_initial_ue_message_t *const initial_pP)
    {
         OAILOG_FUNC_IN(LOG_AMF_APP);
        class ue_m5gmm_context_s* ue_context_p = NULL;
        bool is_guti_valid                   = false;
        bool is_mm_ctx_new                   = false;
        gnb_ngap_id_key_t gnb_ngap_id_key    = INVALID_GNB_UE_NGAP_ID_KEY;
        imsi64_t imsi64                      = INVALID_IMSI64;
        OAILOG_INFO(LOG_AMF_APP, "Received AMF_APP_INITIAL_UE_MESSAGE from NGAP\n");

    if (initial_pP-> amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
        OAILOG_ERROR(LOG_AMF_APP,"AMF UE NGAP Id (" AMF_UE_NGAP_ID_FMT ") is already assigned\n",
            initial_pP->amf_ue_ngap_id);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
    }

    // Check if there is any existing UE context using S-TMSI/GUTI
    if (initial_pP->is_s_tmsi_valid) {
        OAILOG_DEBUG(LOG_AMF_APP,"INITIAL UE Message: Valid amf_code %u and S-TMSI %u received from "
           "gNB.\n",initial_pP->opt_s_tmsi._code, initial_pP->opt_s_tmsi.m_tmsi);
        guti_m5_t guti = {.guamfi.plmn     = {0},
                    .guamfi.amf_gid  = 0,
                    .guamfi.amf_code = 0,
                    .guamfi.amf_Pointer=0,
                    .m_tmsi          = INVALID_M_TMSI};
        plmn_t plmn = {.mcc_digit1 = initial_pP->tai.mcc_digit1,
                    .mcc_digit2 = initial_pP->tai.mcc_digit2,
                    .mcc_digit3 = initial_pP->tai.mcc_digit3,
                    .mnc_digit1 = initial_pP->tai.mnc_digit1,
                    .mnc_digit2 = initial_pP->tai.mnc_digit2,
                    .mnc_digit3 = initial_pP->tai.mnc_digit3};
        is_guti_valid = amf_app_construct_guti(&plmn, &(initial_pP->opt_s_tmsi), &guti);
        // create a new ue context if nothing is found
        if (is_guti_valid) 
        {
            ue_context_p = amf_ue_context_exists_guti(&amf_app_desc_p->amf_ue_contexts, &guti);
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

                OAILOG_ERROR(LOG_AMF_APP,"AMF_APP_INITAIL_UE_MESSAGE: gnb_ngap_id_key %ld has "
                    "valid value \n",ue_context_p->gnb_ngap_id_key);
                // Inform ngap for local cleanup of gnb_ue_ngap_id from ue context
                ue_context_p->ue_context_rel_cause = NGAP_INVALID_GNB_ID;
                OAILOG_ERROR(LOG_AMF_APP," Sending UE Context Release to NGAP for ue_id =(%u)\n",
                    ue_context_p->amf_ue_ngap_id);
                amf_app_ue_context_release(ue_context_p, ue_context_p->ue_context_rel_cause);
                hashtable_uint64_ts_remove(amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl,
                    (const hash_key_t) ue_context_p->gnb_ngap_id_key);
                ue_context_p->gnb_ngap_id_key      = INVALID_GNB_UE_NGAP_ID_KEY;
                ue_context_p->ue_context_rel_cause = NGAP_INVALID_CAUSE;
                }
        // Update AMF UE context with new gnb_ue_ngap_id
        ue_context_p->gnb_ue_ngap_id = initial_pP->gnb_ue_ngap_id;
        // regenerate the gnb_ngap_id_key as gnb_ue_ngap_id is changed.
        AMF_APP_GNB_NGAP_ID_KEY(gnb_ngap_id_key, initial_pP->gnb_id, initial_pP->gnb_ue_ngap_id);
        // Update gnb_ngap_id_key in hashtable
        amf_ue_context_update_coll_keys(
            &amf_app_desc_p->amf_ue_contexts, ue_context_p, gnb_ngap_id_key,
            ue_context_p->amf_ue_ngap_id, ue_context_p->amf_context._imsi64,
            ue_context_p->amf_teid_n11, &guti);
        imsi64 = ue_context_p->amf_context._imsi64;
        // Check if paging timer exists for UE and remove
        if (ue_context_p->paging_response_timer.id !=
            MME_APP_TIMER_INACTIVE_ID) {
          nas_itti_timer_arg_t* timer_argP = NULL;
          if (timer_remove(ue_context_p->paging_response_timer.id, (void**) &timer_argP)) {
            OAILOG_ERROR_UE( LOG_AMF_APP, imsi64, "Failed to stop paging response timer for UE id %d\n",
                ue_context_p->amf_ue_ngap_id);
          }
          if (timer_argP) {
            free_wrapper((void**) &timer_argP);
          }
          ue_context_p->paging_response_timer.id = AMF_APP_TIMER_INACTIVE_ID;
        }
      } else {
        OAILOG_DEBUG(LOG_AMF_APP, "No UE context found for AMF code %u and S-TMSI %u\n",
            initial_pP->opt_s_tmsi.amf_code, initial_pP->opt_s_tmsi.m_tmsi);
      }
    } else {
      OAILOG_DEBUG(LOG_AMF_APP, "No AMF is configured with AMF code %u received in S-TMSI %u from "
          "UE.\n", initial_pP->opt_s_tmsi.amf_code, initial_pP->opt_s_tmsi.m_tmsi);
    }
  } else {
    OAILOG_DEBUG(LOG_AMF_APP, "AMF_APP_INITIAL_UE_MESSAGE from NGAP,without S-TMSI. \n");
  }
  // create a new ue context if nothing is found
  if (!(ue_context_p)) {
    OAILOG_DEBUG(LOG_AMF_APP, "UE context doesn't exist -> create one\n");
    if (!(ue_context_p=amf_create_new_ue_context();
        // Allocate new amf_ue_ngap_id
        ue_context_p->amf_ue_ngap_id = amf_app_ue_context::amf_app_ctx_get_new_ue_id(&amf_app_desc_p->amf_app_ue_ngap_id_generator);
        if (ue_context_p->amf_ue_ngap_id == INVALID_AMF_UE_NGAP_ID) 
        {
            OAILOG_CRITICAL(LOG_AMF_APP,"AMF_APP_INITIAL_UE_MESSAGE. AMF_UE_NGAP_ID allocation Failed.\n");
            amf_app_ue_context::amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
            OAILOG_FUNC_RETURN(LOG_AMF_APP, imsi64);
        }
        amf_app_ue_context::amf_insert_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
        
        amf_app_ue_context::notify_ngap_new_ue_amf_ngap_id_association(ue_m5gmm_context_s *ue_context_p);
        s_tmsi_m5_t s_tmsi = {0};
        if (initial_pP->is_s_tmsi_valid) {
            s_tmsi = initial_pP->opt_s_tmsi;
        } else {
            s_tmsi.amf_code = 0;
            s_tmsi.m_tmsi   = INVALID_M_TMSI;
        }
        OAILOG_INFO_UE( LOG_AMF_APP, ue_context_p->amf_context._imsi64,
            "INITIAL_UE_MESSAGE RCVD \n" "amf_ue_ngap_id  = %d\n" "gnb_ue_ngap_id  = %d\n",
            ue_context_p->amf_ue_ngap_id, ue_context_p->gnb_ue_ngap_id);
        OAILOG_DEBUG(LOG_AMF_APP, "Is S-TMSI Valid - (%d)\n", initial_pP->is_s_tmsi_valid);

        OAILOG_INFO_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
            "Sending NAS Establishment Indication to NAS for ue_id = (%d)\n",
            ue_context_p->amf_ue_ngap_id);

        amf_ue_ngap_id_t ue_id = ue_context_p->amf_ue_ngap_id;
        int nas_proc::nas_proc_establish_ind(ue_context_p->amf_ue_ngap_id, is_mm_ctx_new, 
                                             initial_pP->tai,initial_pP->ecgi, 
                                             initial_pP->m5g_rrc_establishment_cause, 
                                             s_tmsi,&initial_pP->nas);


    }

}