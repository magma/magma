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

  Source      nas_proc.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "../3gpp/3gpp_23.003.h"
#include "amf_common_defs.h"
#include "amf_data.h"
using namespace std;
#pragma once
////////////////////typecast require///////////////////////////////
typedef int (*success_cb_t)( amf_context_s *);
typedef int (*failure_cb_t)( amf_context_s *);
typedef int (*proc_abort_t)( amf_context_s *,  nas_base_proc_s *);

typedef int (*pdu_in_resp_t)( amf_context_s *, void *arg); // can be RESPONSE, COMPLETE, ACCEPT
typedef int (*pdu_in_rej_t)( amf_context_s *, void *arg); // REJECT.
typedef int (*pdu_out_rej_t)( amf_context_s *, nas_base_proc_s *); // REJECT.
typedef void (*time_out_t)(void *arg);

typedef int (*sdu_out_delivered_t)(amf_context_s *,  nas_amf_proc_s *);
typedef int (*sdu_out_not_delivered_t)( amf_context_s *,  nas_amf_proc_s *);
typedef int (*sdu_out_not_delivered_ho_t)( amf_context_s *,  nas_amf_proc_s *);


/////////////////////////////////////////////////////////////////////
namespace magma5g
{
    class nas_proc:amf_sap_t
    {

        public:
        /*
        * --------------------------------------------------------------------------
        *          NAS procedures triggered by the user
        * --------------------------------------------------------------------------
        */
       int nas_proc_establish_ind(const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
                                  const tai_t originating_tai, const ecgi_t ecgi, 
                                  const as_m5gcause_t as_cause,const s_tmsi_t s_tmsi, 
                                  std::string *msg);
      nas_amf_registration_proc_t* get_nas_specific_procedure_registration(const struct amf_context_s* const ctxt);
    };
    class nas_base_proc_t
     {
       public:
       success_cb_t success_notif;
       failure_cb_t failure_notif;
       proc_abort_t abort;

        // PDU interface
        //pdu_in_resp_t           resp_in;
       pdu_in_rej_t fail_in;
       pdu_out_rej_t fail_out;
       time_out_t time_out;
       nas_base_proc_type_t type; // AMF, SMF, CN
       uint64_t nas_puid;         // procedure unique identifier for internal use

        struct nas_base_proc_s *parent;
        struct nas_base_proc_s *child;
      };

      enum nas_amf_proc_type_t
      {
        NAS_AMF_PROC_TYPE_NONE = 0,
        NAS_AMF_PROC_TYPE_SPECIFIC,
        NAS_AMF_PROC_TYPE_COMMON,
        NAS_AMF_PROC_TYPE_CONN_MNGT,
      } ;
    // AMF Specific procedures
     class nas_amf_proc_t
     {
       public:
        nas_base_proc_t base_proc;
        nas_amf_proc_type_t type; // specific, common, connection management
        // SDU interface
        sdu_out_delivered_t delivered;
        sdu_out_not_delivered_t not_delivered;
        sdu_out_not_delivered_ho_t not_delivered_ho;

        amf_fsm_state_t previous_amf_fsm_state;
    };
      enum nas_base_proc_type_t
      {
        NAS_PROC_TYPE_NONE = 0,
        NAS_PROC_TYPE_AMF,
        NAS_PROC_TYPE_SMF,
        NAS_PROC_TYPE_CN,
      };

     enum amf_specific_proc_type_t
     {
        AMF_SPEC_PROC_TYPE_NONE = 0,
        AMF_SPEC_PROC_TYPE_ATTACH,
        AMF_SPEC_PROC_TYPE_DETACH,
        AMF_SPEC_PROC_TYPE_TAU,
     } ;

    // AMF Specific procedures
    class nas_amf_specific_proc_t
    {
      public:
      nas_amf_proc_t amf_proc;
      amf_specific_proc_type_t type;
    } ;
    class nas_amf_registration_proc_t
    {
      public:
      nas_amf_specific_proc_t amf_spec_proc;
      //struct nas_timer_s T3450; // AMF message retransmission timer
      //#define REGISTRATION_COUNTER_MAX 5
      int registration_accept_sent;
      bool registration_reject_sent;
      bool registration_complete_received;
      guti_t guti;
      std::string amf_msg_out; // SMF message to be sent within the Registration Accept message
      class amf_registration_request_ies_t *ies;
      amf_ue_ngap_id_t ue_id;
      ksi_t ksi;
      int amf_cause;
    }; 
    class identification:public amf_context_t
    {
      public:
      static const char* amf_identity_type_str[] = {"NOT AVAILABLE", "IMSI", "IMEI", "IMEISV", "TMSI"};
      int amf_proc_identification(amf_context_t* const amf_context, nas_amf_proc_t* const amf_proc,
                                  const identity_type2_t type, success_cb_t success, failure_cb_t failure);
      int amf_proc_identification_complete( const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, 
                                      imei_t* const imei,imeisv_t* const imeisv, uint32_t* const tmsi);


    };
    class nas_amf_auth_proc_t 
    {
      public:
      nas_amf_common_proc_t amf_com_proc;
      struct nas_timer_s T3460; /* Authentication timer         */
      #define AUTHENTICATION_COUNTER_MAX 5
      unsigned int retransmission_count;
      #define EMM_AUTHENTICATION_SYNC_FAILURE_MAX 2
      unsigned int  sync_fail_count; /* counter of successive AUTHENTICATION FAILURE messages
                                    from the UE with AMF cause #21 "synch failure" */
      unsigned int mac_fail_count;
      amf_ue_ngap_id_t ue_id;
      bool is_cause_is_registered;  //  could also be done by seeking parent procedure
      ksi_t ksi;
      uint8_t rand[AUTH_RAND_SIZE]; /* Random challenge number  */
      uint8_t autn[AUTH_AUTN_SIZE]; /* Authentication token     */
      imsi_t* unchecked_imsi;
      int amf_cause;
    };
    class nas_5g_auth_info_proc_t {
      public:
      nas_cn_proc_t cn_proc;
      success_cb_t success_notif;
      failure_cb_t failure_notif;
      bool request_sent;
      uint8_t nb_vectors;
      m5g_vector_t* vector[MAX_5G_AUTH_VECTORS];
      int nas_cause;
      amf_ue_ngap_id_t ue_id;
      bool resync;  // Indicates whether the authentication information is requested
                    // due to sync failure
    } ;
    class authentication: public amf_context_t
    {
      public:
      int amf_proc_authentication_ksi(amf_context_s* amf_context, nas_amf_specific_proc_t* const amf_specific_proc,
                                     ksi_t ksi,const uint8_t* const rand, const uint8_t* const autn, 
                                     success_cb_t success, failure_cb_t failure);

        int amf_proc_authentication( amf_context_s* amf_context, nas_amf_specific_proc_t* const amf_specific_proc,
                                     success_cb_t success, failure_cb_t failure);

        int amf_proc_authentication_failure(amf_ue_ngap_id_t ue_id, int amf_cause, const_bstring auts);

        int amf_proc_authentication_complete(amf_ue_ngap_id_t ue_id, authentication_response_msg* msg, 
                                            int amf_cause, const_bstring const res);

        int amf_registration_security(amf_context_s* amf_context);

        void set_notif_callbacks_for_5g_auth_proc(nas_amf_auth_proc_t* auth_proc);
        void set_callbacks_for_5g_auth_proc(nas_amf_auth_proc_t* auth_proc);
        void set_callbacks_for_5g_auth_info_proc(nas_5g_auth_info_proc_t* auth_info_proc);

    };
    // 5G CN Specific procedures
    typedef struct nas_amf_common_procedure_s {
      nas_amf_common_proc_t* proc;
      LIST_ENTRY(nas_amf_common_procedure_s) entries;
    } nas_amf_common_procedure_t;
      typedef struct nas_amf_proc_s {
      nas_base_proc_t base_proc;
      nas_amf_proc_type_t type;  // specific, common, connection management
      // SDU interface
      sdu_out_delivered_t delivered;
      sdu_out_not_delivered_t not_delivered;
      sdu_out_not_delivered_ho_t not_delivered_ho;

      amf_fsm_state_t previous_amf_fsm_state;
    } nas_amf_proc_t;
    typedef struct nas_amf_common_proc_s {
    nas_amf_proc_t emm_proc;
    amf_common_proc_type_t type;
    } nas_amf_common_proc_t;
    typedef struct nas_amf_ident_proc_s {
      nas_amf_common_proc_t amf_com_proc;
      struct nas_timer_s T3470; /* Identification timer         */
      #define IDENTIFICATION_COUNTER_MAX 5
      unsigned int retransmission_count;
      amf_ue_ngap_id_t ue_id;
      bool is_cause_is_registered;  //  could also be done by seeking parent procedure
      identity_type2_t identity_type;
    } nas_amf_ident_proc_t;
}