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

  Source      amf_app_ue_context.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_config.h"
#include "bstrlib.h"
#include "amf_common_defs.h"
#include "5GMMCapability.h"
#include "amf_data.h"
extern c
{
    #include "hashtable.h"
    //#include "bstrlib.h"
    #include "../3gpp/3gpp_23.003.h"
};
using namespace std;
#pragma once

namespace magma5g
{
    class amf_ue_context_t 
    {
        public:

        /* hash_table_uint64_ts_t is defined in lib/hastable*/
        hash_table_uint64_ts_t* imsi_amf_ue_id_htbl;   // data is amf_ue_ngap_id_t
        hash_table_uint64_ts_t* tun11_ue_context_htbl; // data is amf_ue_ngap_id_t
        hash_table_uint64_ts_t* gnb_ue_ngap_id_ue_context_htbl;              // data is amf_ue_ngap_id_t
        obj_hash_table_uint64_t* guti_ue_context_htbl; // data is amf_ue_ngap_id_t
    };
    
    
    enum mm_state_t
    {
        UE_UNREGISTERED = 0,
        UE_REGISTERED,
    };
    enum ecm_state_t
    {
        ECM_IDLE = 0,
        ECM_CONNECTED,
    };
    class ue_m5gmm_context_s :public amf_context_t
    {
        public:
        /* msisdn: The basic MSISDN of the UE. The presence is dictated by its storage
        *         in the HSS, set by S6A UPDATE LOCATION ANSWER
        */
        bstring msisdn;

        enum Ngcause ue_context_rel_cause; //define require for Ngcause in NGAP module
        mm_state_t mm_state;
        ecm_state_t ecm_state;

        /* Last known 5G cell, set by nas_registration_req_t */
        ecgi_t e_utran_cgi;// from 3gpp 23.003

        /* cell_age: Time elapsed since the last 5G Cell Global Identity was
        *           acquired. Set by nas_auth_param_req_t
        */
        time_t cell_age; 
        /* TODO: add csg_id */
        /* TODO: add csg_membership */
        /* TODO Access mode: Access mode of last known ECGI when the UE was active */

        /* apn_config_profile: set by UPDATE LOCATION ANSWER */
        apn_config_profile_t apn_config_profile;

        /* access_restriction_data: The access restriction subscription information.
        *           set by UPDATE LOCATION ANSWER
        */
        ard_t access_restriction_data;

        
        bstring apn_oi_replacement;
        teid_t amf_teid_n11;
        /* SCTP assoc id */
        sctp_assoc_id_t sctp_assoc_id_key;

        /* gNB UE NGAP ID,  Unique identity the UE within gNodeB */
        gnb_ue_ngap_id_t gnb_ue_ngap_id : 24;
        
        gnb_ngap_id_key_t gnb_ngap_id_key;

        /* AMF UE NGAP ID, Unique identity the UE within AMF */
        amf_ue_ngap_id_t amf_ue_ngap_id;

        /* Subscribed UE-AMBR: The Maximum Aggregated uplink and downlink MBR values
        *           to be shared across all Non-GBR bearers according to the
        *           subscription of the user. Set by SMF UPDATE LOCATION ANSWER
        */
        ambr_t subscribed_ue_ambr;
        /* used_ue_ambr: The currently used Maximum Aggregated uplink and downlink
        *           MBR values to be shared across all Non-GBR bearers.
        *           Set by S6A UPDATE LOCATION ANSWER
        */
        ambr_t used_ue_ambr;
        /* rau_tau_timer: Indicates a subscribed Periodic RAU/TAU Timer value
        *           Set by S6A UPDATE LOCATION ANSWER
        */
        rau_tau_timer_t rau_tau_timer;

        int nb_active_pdn_contexts;
        pdn_context_t* pdn_contexts[MAX_APN_PER_UE];

        amf_context_t amf_context;
        bearer_context_t* bearer_contexts[BEARERS_PER_UE];

        /* ue_radio_capability: Store the radio capabilities as received in
        *           S1AP UE capability indication message
        */
        bstring ue_radio_capability;

        /* mobile_reachability_timer: Start when UE moves to idle state.
        *             Stop when UE moves to connected state
        */
        struct amf_app_timer_t m5_mobile_reachability_timer;
        /* implicit_detach_timer: Start at the expiry of Mobile Reachability timer.
        * Stop when UE moves to connected state
        */
        struct amf_app_timer_t m5_implicit_detach_timer;
        /* Initial Context Setup Procedure Guard timer */
        struct amf_app_timer_t m5_initial_context_setup_rsp_timer;
        /* UE Context Modification Procedure Guard timer */
        struct amf_app_timer_t m5_ue_context_modification_timer;
        /* Timer for retrying paging messages */
        struct amf_app_timer_t m5_paging_response_timer;
        /* send_ue_purge_request: If true AMF shall send - Purge Req to
        * delete contexts at HSS
        */
        bool send_ue_purge_request;

        bool hss_initiated_detach;
        bool location_info_confirmed_in_hss;
        /* S6a- update location request guard timer */
        struct amf_app_timer_t m5_ulr_response_timer;
        
        uint8_t attach_type;
        lai_t lai;
        int cs_fallback_indicator;
        uint8_t sgs_detach_type;
        /* granted_service_t: informs the granted service to UE */
        granted_service_t m5_granted_service;
        /*  paging_proceeding_flag (PPF) shall set to true, when UE moves to connected.
        * Indicates that paging procedure can be prooceeded,
        * Is set to false, due to "Inactivity of UE including lack of periodic TAU"
        */
        bool ppf;

        #define SUBSCRIPTION_UNKNOWN false
        #define SUBSCRIPTION_KNOWN true
        bool subscription_known;
        ambr_t used_ambr;
        subscriber_status_t subscriber_status;
        network_access_mode_t network_access_mode;

        bool path_switch_req;
        LIST_HEAD(s11_procedures_s, mme_app_s11_proc_s) * s11_procedures;
    };
                /** @class ue_m5gmm_context_s
             *  @brief Useful parameters to know in AMF application layer. They are set
             * according to 3GPP TS.23.518 #6.1.6.2.25
             */
    class ue_mm_context
    {
        public:

        /* msisdn: The basic MSISDN of the UE. The presence is dictated by its storage
        *         in the UDM, set by N8 UPDATE LOCATION ANSWER
        */
        std::string imsi;
        bool supi_UnauthInd;
       // std::string gpsiList[] array(Gpsi);
        std::string Pei pei;
        uint64_t udmGroupId ; //NfGroupId
        uint64_t ausfGroupId ;//NfGroupId;
        std::string routingIndicator;
        //std::auto groupList[] array(GroupId);
        std::string drxParameter;
        std::string subRfsp;
        uint32_t subRfsp;//RfspIndex type;
        uint32_t usedRfsp;//RfspIndex type;
        ambr_t subUeAmbr ;
        bool smsSupport;
        std::string smsfId //NfInstanceId type
        std::string seafData //SeafData which will come while AUSF communication for AUTN.
        class M5GMM_Capability_msg m5gMmCapability //5GMmCapability
        std::string pcfId; //NfInstanceId
        std::string pcfAmPolicyUri; //Uri
        std::auto amPolicyReqTriggerList ;//array(PolicyReqTrigger)
        std::string pcfUePolicyUri ;//Uri
        std::auto uePolicyReqTriggerList; //array(PolicyReqTrigger)
        std::string hpcfId; //NfInstanceId
        std::string restrictedRatList; //array(RatType)
        std::string forbiddenAreaList; //array(Area)
        std::string serviceAreaRestriction ;//ServiceAreaRestriction
        std::string restrictedCnList; //array(CoreNetworkType)
        std::string eventSubscriptionList ;//array(AmfEventSubscription)
        std::string mmContextList ;//array(MmContext)
        std::string sessionContextList; //array(PduSessionContext)
        std::string traceData; //TraceData
         /* SCTP assoc id */
        sctp_assoc_id_t sctp_assoc_id_key;

        /* gNB UE NGAP ID,  Unique identity the UE within gNodeB */
        gnb_ue_ngap_id_t gnb_ue_ngap_id : 24;
       
        gnb_ngap_id_key_t gnb_ngap_id_key;

        /* MME UE S1AP ID, Unique identity the UE within MME */
        amf_ue_ngap_id_t amf_ue_ngap_id;
    };

    /*
    * Timer identifier returned when in inactive state (timer is stopped or has
    * failed to be started)
    */
    #define AMF_APP_TIMER_INACTIVE_ID (-1)

    #define AMF_APP_DELTA_T3412_REACHABILITY_TIMER 4            // in minutes
    #define AMF_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER 0  // in minutes

    #define AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE 2  // In seconds
    #define AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE 2    // In seconds
    #define AMF_APP_PAGING_RESPONSE_TIMER_VALUE 4            // In seconds
    #define AMF_APP_ULR_RESPONSE_TIMER_VALUE 3               // In seconds
    /* Timer structure */
    struct amf_app_timer_t {
    long id;  /* The timer identifier                 */
    long sec; /* The timer interval value in seconds  */
    };
    class amf_app_ue_context: public amf_ue_context_t :ue_m5gmm_context_s
    {
        public:
        // check & create state information. 
        int amf_insert_ue_context(amf_ue_context_t *const amf_ue_context, const class ue_m5gmm_context_s* const ue_context_p);
        amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p);
        // Notify NGAP about the mapping between amf_ue_ngap_id and
        // sctp assoc id + gnb_ue_ngap_id
        void notify_ngap_new_ue_amf_ngap_id_association(ue_context_p);
        void amf_remove_ue_context(amf_ue_context_t* const amf_ue_context, class ue_m5gmm_context_s* const ue_context_p);
    };
    
    ue_m5gmm_context_s* amf_create_new_ue_context(void);

}