
**
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

  Source      amf_proc.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "bstrlib.h"
#include "amf_securityDef.h"
#include "amf_fsm.h"
using namespace std;
#pragma once
typedef uint8_t ksi_t;
#define IS_AMF_CTXT_PRESENT_SECURITY(aMfCtXtPtR)                               \
  (!!((aMfCtXtPtR)->member_present_mask & AMF_CTXT_MEMBER_SECURITY))

 #define IS_AMF_CTXT_VALID_AUTH_VECTORS(aMfCtXtPtR)                             \
  (!!((aMfCtXtPtR)->member_valid_mask & AMF_CTXT_MEMBER_AUTH_VECTORS))


namespace magma5g
{

  class amf_procedures_t {
    nas_amf_specific_proc_t* amf_specific_proc;
    LIST_HEAD(nas_amf_common_procedures_head_s, nas_amf_common_procedure_s)
    amf_common_procs;
    LIST_HEAD(nas_cn_procedures_head_s, nas_cn_procedure_s)
    cn_procs;  // triggered by AMF
    nas_amf_con_mngt_proc_t* amf_con_mngt_proc;

      int nas_proc_mess_sign_next_location;  // next index in array
     #define MAX_NAS_PROC_MESS_SIGN 3
      nas_proc_mess_sign_t nas_proc_mess_sign[MAX_NAS_PROC_MESS_SIGN];
  } ;

    /*
    * Structure of the AMF context established by the network for a particular UE
    * ---------------------------------------------------------------------------
    */
    class amf_context_t
    {
        public:
        bool is_dynamic;   /* Dynamically allocated context indicator         */
        bool is_registered;  /* Registration indicator                            */
        //bool is_emergency; /* Emergency bearer services indicator             */
        bool is_initial_identity_imsi; // If the IMSI was used for identification in the initial NAS message
        bool is_guti_based_registered;
        /*
        * registration_type has type amf_proc_registration_type_t.
        *
        * Here, it is un-typedef'ed as uint8_t to avoid circular dependency issues.
        */
        uint8_t m5gsregistrationtype; 
        amf_procedures_t* amf_procedures;

        uint num_registration_request; /* Num registration request received               */
         
       // imsi present mean we know it but was not checked with identity proc, or was not provided in initial message
        imsi_t _imsi; /* The IMSI provided by the UE or the AMF, set valid when identification returns IMSI */
        imsi64_t _imsi64; /* The IMSI provided by the UE or the AMF, set valid when identification returns IMSI */
        imsi64_t saved_imsi64; /* Useful for 5.4.2.7.c */
        imei_t _imei;          /* The IMEI provided by the UE                     */
        imeisv_t _imeisv;      /* The IMEISV provided by the UE                   */
        //bool                   _guti_is_new; /* The GUTI assigned to the UE is new              */
        guti_m5_t _m5_guti;         /* The GUTI assigned to the UE                     */
        //guti_t m5_old_guti;     /* The old GUTI (GUTI REALLOCATION)                */
        //tai_list_t _tai_list;   /* TACs the the UE is registered to                */
        //tai_t _lvr_tai;
        //tai_t originating_tai;

       
        ksi_t ksi;          /*key set identifier  */
        //ue_network_capability_t _ue_network_capability; // will be use in perodic registration
        //ms_network_capability_t _ms_network_capability;
        drx_parameter_t _drx_parameter;

        int remaining_vectors; // remaining unused vectors
        auth_vector_t _vector[MAX_EPS_AUTH_VECTORS]; /* 5GMM authentication vector                            */
        //amf_security_context_t _security; /* Current 5GMM security context: The security context which has been activated most recently. Note that a current 5GMM
                                                                //security context originating from either a mapped or native 5GMM security context may exist simultaneously with a native
                                                               // non-current 5GMM security context.*/

        // Requirement MME24.301R10_4.4.2.1_2
        //amf_security_context_t  _non_current_security; /* Non-current 5GMM security context: A native 5GMM security context that is not the current one. A non-current 5GMM
                                                               /* security context may be stored along with a current 5GMM security context in the UE and the MME. A non-current 5GMM
                                                                security context does not contain an 5GMM AS security context. A non-current 5GMM security context is either of type 'full
                                                                native' or of type 'partial native'.     */

        int amf_cause; /* EMM failure cause code                          */

        amf_fsm_state_t anf_fsm_state;

        //nas_timer_t T3422; /* Deregister timer         */
        void *t3422_arg;

        //struct smf_context_s smf_ctx; //smf contents

        drx_parameter_t _current_drx_parameter; /* stored TAU Request IE Requirement AMF24.501R15_5.5.3.2.4_4*/
       
        // TODO: DO BETTER  WITH BELOW
        std::string smf_msg; /* SMF message contained within the initial request*/
        bool is_imsi_only_detach;
        
    };
    
    class count_s 
    {
      public:
      uint32_t spare : 8;
      uint32_t overflow : 16;
      uint32_t seq_num : 8;
    };  /* Downlink and uplink count params */
    class capability 
    {
      public:
      uint8_t m5gs_encryption;  /* algorithm used for ciphering            */
      uint8_t m5gs_integrity;   /* algorithm used for integrity protection */
      uint8_t umts_encryption; /* algorithm used for ciphering            */
      uint8_t umts_integrity;  /* algorithm used for integrity protection */
      uint8_t gprs_encryption; /* algorithm used for ciphering            */
      bool umts_present;
      bool gprs_present;
    } ; /* UE network capability           */
    enum amf_sc_type_t {
      SECURITY_CTX_TYPE_NOT_AVAILABLE = 0,
      SECURITY_CTX_TYPE_PARTIAL_NATIVE,
      SECURITY_CTX_TYPE_FULL_NATIVE,
      SECURITY_CTX_TYPE_MAPPED  // UNUSED
    } ;
    class selected_algorithms
    {
      public:
      uint8_t encryption : 4; /* algorithm used for ciphering           */
      uint8_t integrity : 4;  /* algorithm used for integrity protection */
    } ;    /* MME selected algorithms                */
    class amf_security_context_t : public count_s , public capability , public selected_algorithms
    {
      public:
      amf_sc_type_t sc_type; /* Type of security context        */
      /* state of security context is implicit due to its storage location
      * (current/non-current)*/
    #define EKSI_MAX_VALUE 6
      ksi_t eksi; /* NAS key set identifier for E-UTRAN      */
    #define EMM_SECURITY_VECTOR_INDEX_INVALID (-1)
      int vector_index;                     /* Pointer on vector */
      uint8_t knas_enc[AUTH_KNAS_ENC_SIZE]; /* NAS cyphering key               */
      uint8_t knas_int[AUTH_KNAS_INT_SIZE]; /* NAS integrity key               */ 

      // Requirement MME24.301R10_4.4.4.3_2 (DETACH REQUEST (if sent before security
      // has been activated);)
      uint8_t activated;
      uint8_t direction_encode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
      uint8_t direction_decode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
      // security keys for HO
      uint8_t next_hop[AUTH_NEXT_HOP_SIZE]; /* Next HOP security parameter */
      uint8_t next_hop_chaining_count;      /* Next Hop Chaining Count */
    };
    void amf_ctx_set_valid_imsi(amf_context_t* const ctxt, imsi_t* imsi, const imsi64_t imsi64)__attribute__((nonnull)) __attribute__((flatten));
    void amf_ctx_set_attribute_valid(amf_context_t* const ctxt, const int attribute_bit_pos)__attribute__((nonnull)) __attribute__((flatten));
   // int  amf_context_upsert_imsi(amf_data_t* amf_data,  amf_context_t* elm)__attribute__((nonnull));

}