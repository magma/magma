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

  Source      amf_app_ue_context.cpp

  Version     0.1

  Date        2020/09/21

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include "amf_app_ue_context.h"
#include "amf_data.h"

using namespace std;

namespace magma5g
{
    amf_app_ue_context::amf_app_ctx_get_new_ue_id(amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p))
    {
       amf_ue_ngap_id_t tmp = 0;
       tmp = __sync_fetch_and_add(amf_app_ue_ngap_id_generator_p, 1);
       return tmp;
    }
 //------------------------------------------------------------------------------
 void amf_app_ue_context::notify_ngap_new_ue_amf_ngap_id_association(class ue_m5gmm_context_s* ue_context_p) 
 {

    MessageDef* message_p                                      = NULL;
    itti_amf_app_ngap_amf_ue_id_notification_t* notification_p = NULL;

    OAILOG_FUNC_IN(LOG_AMF_APP);
    if (ue_context_p == NULL) {
        OAILOG_ERROR(LOG_AMF_APP, " NULL UE context pointer!\n");
        OAILOG_FUNC_OUT(LOG_AMF_APP);
    }
    message_p = itti_alloc_new_message(TASK_AMF_APP, AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION);
    notification_p = &message_p->ittiMsg.amf_app_ngap_amf_ue_id_notification;
    memset(notification_p, 0, sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
    notification_p->gnb_ue_ngap_id = ue_context_p->gnb_ue_ngap_id;
    notification_p->amf_ue_ngap_id = ue_context_p->amf_ue_ngap_id;
    notification_p->sctp_assoc_id  = ue_context_p->sctp_assoc_id_key;

    OAILOG_DEBUG_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
        " Sent AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION to NGAP for (ue_id = %u)\n",notification_p->amf_ue_ngap_id);

    send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
    OAILOG_FUNC_OUT(LOG_AMF_APP);
}
//-------------------------------------------------------------------------------------------------
    amf_app_ue_context::amf_insert_ue_context(amf_ue_context_t *const amf_ue_context_p, 
                                              const ue_m5gmm_context_s* const ue_context_p)
    {
        // TODO: replace hastable with fooly lib .. 
        hashtable_rc_t h_rc                 = HASH_TABLE_OK;
       // hash_table_ts_t* amf_state_ue_id_ht = get_amf_ue_state();

        OAILOG_FUNC_IN(LOG_AMF_APP);
        if (amf_ue_context_p == NULL) {
            OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
            OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }
        if (ue_context_p == NULL) {
            OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
            OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
        }

        // filled GNB UE NGAP ID
        h_rc = hashtable_uint64_ts_is_key_exists(
           amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
           (const hash_key_t) ue_context_p->gnb_ngap_id_key);
       
        h_rc = hashtable_uint64_ts_insert(amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
         (const hash_key_t) ue_context_p->gnb_ngap_id_key,
        ue_context_p->amf_ue_ngap_id);

        if (INVALID_AMF_UE_NGAP_ID != ue_context_p->amf_ue_ngap_id) {
           // filled IMSI
            if (ue_context_p->amf_context_t._imsi64) {
           

            _directoryd_report_location(
                ue_context_p->amf_context_t._imsi64,
                ue_context_p->amf_context_t._imsi.length);
            }

          // filled guti
            if ((0 != ue_context_p->amf_context_t.m5_guti.guamfi.amf_code) ||
                (0 != ue_context_p->amf_context_t.m5_guti.guamfi.amf_gid) ||
                (0 != ue_context_p->amf_context_t.m5_guti.m_tmsi) ||
                (0 != ue_context_p->amf_context_t.m5_guti.guamfi.plmn.mcc_digit1) ||  // MCC 000 does not exist in ITU table
                (0 != ue_context_p->amf_context_t.m5_guti.guamfi.plmn.mcc_digit2) ||
                (0 != ue_context_p->amf_context_t.m5_guti.guamfi.plmn.mcc_digit3)) {
            h_rc = obj_hashtable_uint64_ts_insert(
                amf_ue_context_p->guti_ue_context_htbl,
                (const void* const) & ue_context_p->amf_context_t.m5_guti,
                sizeof(ue_context_p->amf_context_t.m5_guti),
                ue_context_p->amf_ue_ngap_id);

            if (HASH_TABLE_OK != h_rc) {
                OAILOG_WARNING(
                    LOG_AMF_APP,
                    "Error could not register this ue context %p "
                    "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " guti " GUTI_FMT "\n",
                    ue_context_p, ue_context_p->amf_ue_ngap_id,
                    GUTI_ARG(&ue_context_p->amf_context_t.m5_guti));
                OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
            }
            }
        }

        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);

    }
    //------------------------------------------------------------------------------
// warning: lock the UE context
    ue_m5gmm_context_s* amf_create_new_ue_context(void) 
    {
        ue_m5gmm_context_s* new_p = new(ue_m5gmm_context_s);
        if (!new_p) {
            OAILOG_ERROR(LOG_AMF_APP, "Failed to allocate memory for UE context \n");
            return NULL;
        }

        new_p->amf_ue_ngap_id  = INVALID_AMF_UE_NGAP_ID;
        new_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
        // TODO amf_init_context(&new_p->amf_context, true);

        // Initialize timers to INVALID IDs
        new_p->m5_mobile_reachability_timer.id = AMF_APP_TIMER_INACTIVE_ID;
        new_p->m5_implicit_detach_timer.id     = AMF_APP_TIMER_INACTIVE_ID;

        new_p->m5_initial_context_setup_rsp_timer = (struct amf_app_timer_t){
            AMF_APP_TIMER_INACTIVE_ID, AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
        new_p->m5_paging_response_timer = (struct amf_app_timer_t){
            AMF_APP_TIMER_INACTIVE_ID, AMF_APP_PAGING_RESPONSE_TIMER_VALUE};
        new_p->m5_ulr_response_timer = (struct amf_app_timer_t){
            AMF_APP_TIMER_INACTIVE_ID, AMF_APP_ULR_RESPONSE_TIMER_VALUE};
        new_p->m5_ue_context_modification_timer = (struct amf_app_timer_t){
            AMF_APP_TIMER_INACTIVE_ID, AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};

        new_p->ue_context_rel_cause = NGAP_INVALID_CAUSE;
        
        return new_p;
    }

//------------------------------------------------------------------------------
    amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p) 
    {
        amf_ue_ngap_id_t tmp = 0;
        tmp = __sync_fetch_and_add(amf_app_ue_ngap_id_generator_p, 1);
        return tmp;
    }
//-------------------------------------------------------------------------------
//------------------------------------------------------------------------------
    void amf_app_ue_context::amf_remove_ue_context(amf_ue_context_t* const amf_ue_context_p, class ue_m5gmm_context_s* ue_context_p) 
    {
        OAILOG_FUNC_IN(LOG_AMF_APP);
        hashtable_rc_t hash_rc              = HASH_TABLE_OK;
        hash_table_ts_t* amf_state_ue_id_ht;// TODO = get_amf_ue_state();

        if (!amf_ue_context_p) {
            OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
            OAILOG_FUNC_OUT(LOG_AMF_APP);
        }
        if (!ue_context_p) {
            OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
            OAILOG_FUNC_OUT(LOG_AMF_APP);
        }

        // Release emm and esm context
       //TODO below 
       //  delete_amf_ue_state(ue_context_p->amf_context._imsi64);
       // _clear_amf_ctxt(&ue_context_p->amf_context);
       // amf_app_ue_context_free_content(ue_context_p);
        // IMSI
        if (ue_context_p->emm_context._imsi64) {
            hash_rc = hashtable_uint64_ts_remove(amf_ue_context_p->imsi_amf_ue_id_htbl,
                (const hash_key_t) ue_context_p->emm_context._imsi64);
            if (HASH_TABLE_OK != hash_rc) {
            OAILOG_ERROR_UE(LOG_AMF_APP, ue_context_p->emm_context._imsi64,
                "UE context not found!\n" " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
                " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " not in IMSI collection\n",
                ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
            }
        }

        // gNB UE NGAP UE ID
        hash_rc = hashtable_uint64_ts_remove( amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
            (const hash_key_t) ue_context_p->gnb_ngap_id_key);
        if (HASH_TABLE_OK != hash_rc)
            OAILOG_ERROR(LOG_AMF_APP, "UE context not found!\n"
                " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
                " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
                ", GNB_UE_NGAP_ID not in GNB_UE_NGAP_ID collection",
                ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);

        // filled N11 tun id
        if (ue_context_p->amf_teid_n11) {
            hash_rc = hashtable_uint64_ts_remove( amf_ue_context_p->tun11_ue_context_htbl,
                (const hash_key_t) ue_context_p->amf_teid_s11);
            if (HASH_TABLE_OK != hash_rc)
            OAILOG_ERROR(LOG_AMF_APP, "UE Context not found!\n"
                " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT ", AMF S11 TEID  " TEID_FMT
                "  not in N11 collection\n", ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id,
                ue_context_p->amf_teid_n11);
        }
        // filled guti
        if ((ue_context_p->amf_context._m5_guti.guamfi.amf_code) ||
            (ue_context_p->emm_context._m5_guti.guamfi.amf_gid) ||
            (ue_context_p->emm_context._m5_guti.m_tmsi) ||
            (ue_context_p->emm_context._m5_guti.guamfi.plmn.mcc_digit1) ||
            (ue_context_p->emm_context._m5_guti.guamfi.plmn.mcc_digit2) ||
            (ue_context_p->emm_context._m5_guti.guamfi.plmn.mcc_digit3)) 
            {  // MCC 000 does not exist in ITU table
            hash_rc = obj_hashtable_uint64_ts_remove( amf_ue_context_p->guti_ue_context_htbl,
                (const void* const) & ue_context_p->emm_context._guti, sizeof(ue_context_p->emm_context._guti));
            if (HASH_TABLE_OK != hash_rc)
            OAILOG_ERROR(LOG_AMF_APP,"UE Context not found!\n"
                " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
                ", GUTI  not in GUTI collection\n", ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
        }

        // filled NAS UE ID/ AMF UE NGAP ID
        if (INVALID_AMF_UE_NGAP_ID != ue_context_p->amf_ue_ngap_id) {
            hash_rc = hashtable_ts_remove(
                amf_state_ue_id_ht, (const hash_key_t) ue_context_p->amf_ue_ngap_id,
                (void**) &ue_context_p);
            if (HASH_TABLE_OK != hash_rc)
                OAILOG_ERROR(LOG_AMF_APP, "UE context not found!\n"
                "  gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT ", amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
                " not in AMF UE NGAP ID collection",ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
        }

        //amf_directoryd_remove_location( ue_context_p->amf_context._imsi64,ue_context_p->amf_context._imsi.length);
        free_wrapper((void**) &ue_context_p);
        OAILOG_FUNC_OUT(LOG_AMF_APP);
    }   


}//namespace magma5g
