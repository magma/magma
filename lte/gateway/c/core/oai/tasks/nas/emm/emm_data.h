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
Source      emm_data.h

Version     0.1

Date        2012/10/18

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines internal private data handled by EPS Mobility
        Management sublayer.

*****************************************************************************/
#ifndef FILE_EMM_DATA_SEEN
#define FILE_EMM_DATA_SEEN

#include <sys/types.h>
#include "queue.h"
#include "hashtable.h"
#include "obj_hashtable.h"
#include "nas/securityDef.h"
#include "TrackingAreaIdentityList.h"
#include "emm_fsm.h"
#include "nas_timer.h"
#include "nas_procedures.h"
#include "3gpp_24.301.h"
#include "3gpp_24.008.h"
#include "AdditionalUpdateType.h"
#include "UeNetworkCapability.h"
#include "EpsBearerContextStatus.h"
#include "EpsNetworkFeatureSupport.h"
#include "MobileStationClassmark2.h"
#include "esm_data.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* EPS NAS security context structure
 * EPS NAS security context: This context consists of K ASME with the associated
 * key set identifier, the UE security capabilities, and the uplink and downlink
 * NAS COUNT values. In particular, separate pairs of NAS COUNT values are used
 * for each EPS NAS security contexts, respectively. The distinction between
 * native and mapped EPS security contexts also applies to EPS NAS security
 * contexts. The EPS NAS security context is called 'full' if it additionally
 * contains the keys K NASint and K NASenc and the identifiers of the selected
 * NAS integrity and encryption algorithms.*/
typedef struct emm_security_context_s {
  emm_sc_type_t sc_type; /* Type of security context        */
  /* state of security context is implicit due to its storage location
   * (current/non-current)*/
#define EKSI_MAX_VALUE 6
  ksi_t eksi; /* NAS key set identifier for E-UTRAN      */
#define EMM_SECURITY_VECTOR_INDEX_INVALID (-1)
  int vector_index;                     /* Pointer on vector */
  uint8_t knas_enc[AUTH_KNAS_ENC_SIZE]; /* NAS cyphering key               */
  uint8_t knas_int[AUTH_KNAS_INT_SIZE]; /* NAS integrity key               */

  struct count_s {
    uint32_t spare : 8;
    uint32_t overflow : 16;
    uint32_t seq_num : 8;
  } dl_count, ul_count, kenb_ul_count; /* Downlink and uplink count params */
  struct {
    uint8_t eps_encryption;  /* algorithm used for ciphering            */
    uint8_t eps_integrity;   /* algorithm used for integrity protection */
    uint8_t umts_encryption; /* algorithm used for ciphering            */
    uint8_t umts_integrity;  /* algorithm used for integrity protection */
    uint8_t gprs_encryption; /* algorithm used for ciphering            */
    bool umts_present;
    bool gprs_present;
  } capability; /* UE network capability           */
  struct {
    uint8_t encryption : 4; /* algorithm used for ciphering           */
    uint8_t integrity : 4;  /* algorithm used for integrity protection */
  } selected_algorithms;    /* MME selected algorithms                */

  // Requirement MME24.301R10_4.4.4.3_2 (DETACH REQUEST (if sent before security
  // has been activated);)
  uint8_t activated;
  uint8_t direction_encode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
  uint8_t direction_decode;  // SECU_DIRECTION_DOWNLINK, SECU_DIRECTION_UPLINK
  // security keys for HO
  uint8_t next_hop[AUTH_NEXT_HOP_SIZE]; /* Next HOP security parameter */
  uint8_t next_hop_chaining_count;      /* Next Hop Chaining Count */
} emm_security_context_t;

/*
 * --------------------------------------------------------------------------
 *  EMM internal data handled by EPS Mobility Management sublayer in the MME
 * --------------------------------------------------------------------------
 */
struct emm_common_data_s;

typedef enum { SUCCESS, FAILURE } sgs_loc_updt_status_t;

typedef struct csfb_params_s {
#define MOBILE_IDENTITY (1 << 0)
#define LAI_CSFB (1 << 1)
#define ADD_UPDATE_TYPE (1 << 2)
  uint8_t presencemask;
  // LAI
  location_area_identification_t lai;
  // CSFB-New TMSI allocated
  bool newTmsiAllocated;
  // CSFB-Mobile Id
  mobile_identity_t mobileid;
  sgs_loc_updt_status_t sgs_loc_updt_status;
  uint8_t additional_updt_res;
  bstring esm_data;
  uint8_t tau_active_flag : 1;
} csfb_params_t;

typedef struct volte_params_s {
#define VOICE_DOMAIN_PREF_UE_USAGE_SETTING (1 << 0)
  uint8_t presencemask;
  voice_domain_preference_and_ue_usage_setting_t
      voice_domain_preference_and_ue_usage_setting;
} volte_params_t;

struct emm_attach_request_ies_s;
typedef struct new_attach_info_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id; /* mme_ue_s1ap_id for which Attach Request
                                      is received */
  bool
      is_mm_ctx_new; /* Attach request is received in S1ap initial UE message */
  struct emm_attach_request_ies_s* ies;
} new_attach_info_t;
/*
 * Structure of the EMM context established by the network for a particular UE
 * ---------------------------------------------------------------------------
 */
typedef struct emm_context_s {
  bool is_dynamic;   /* Dynamically allocated context indicator         */
  bool is_attached;  /* Attachment indicator                            */
  bool is_emergency; /* Emergency bearer services indicator             */
  bool is_initial_identity_imsi;  // If the IMSI was used for identification in
                                  // the initial NAS message
  bool is_guti_based_attach;
  /*
   * attach_type has type emm_proc_attach_type_t.
   *
   * Here, it is un-typedef'ed as uint8_t to avoid circular dependency issues.
   */
  uint8_t attach_type; /* EPS/Combined/etc. */
  additional_update_type_t additional_update_type;
  uint8_t tau_updt_type; /*TAU Update type - Normal Update, Periodic,
                                           combined TAU,combined TAU with IMSI*/

  uint num_attach_request; /* Num attach request received               */

  emm_procedures_t* emm_procedures;

  // this bitmask is here because we wanted to avoid modifying the EmmCommon
  // interface
  uint32_t common_proc_mask; /* bitmask, see significance of bits below */
#define EMM_CTXT_COMMON_PROC_GUTI ((uint32_t) 1 << 0)
#define EMM_CTXT_COMMON_PROC_AUTH ((uint32_t) 1 << 1)
#define EMM_CTXT_COMMON_PROC_SMC ((uint32_t) 1 << 2)
#define EMM_CTXT_COMMON_PROC_IDENT ((uint32_t) 1 << 3)
#define EMM_CTXT_COMMON_PROC_INFO ((uint32_t) 1 << 4)

  uint32_t member_present_mask; /* bitmask, see significance of bits below */
  uint32_t member_valid_mask;   /* bitmask, see significance of bits below */
#define EMM_CTXT_MEMBER_IMSI ((uint32_t) 1 << 0)
#define EMM_CTXT_MEMBER_IMEI ((uint32_t) 1 << 1)
#define EMM_CTXT_MEMBER_IMEI_SV ((uint32_t) 1 << 2)
#define EMM_CTXT_MEMBER_OLD_GUTI ((uint32_t) 1 << 3)
#define EMM_CTXT_MEMBER_GUTI ((uint32_t) 1 << 4)
#define EMM_CTXT_MEMBER_TAI_LIST ((uint32_t) 1 << 5)
#define EMM_CTXT_MEMBER_LVR_TAI ((uint32_t) 1 << 6)
#define EMM_CTXT_MEMBER_AUTH_VECTORS ((uint32_t) 1 << 7)
#define EMM_CTXT_MEMBER_SECURITY ((uint32_t) 1 << 8)
#define EMM_CTXT_MEMBER_NON_CURRENT_SECURITY ((uint32_t) 1 << 9)
#define EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE ((uint32_t) 1 << 10)
#define EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE ((uint32_t) 1 << 11)
#define EMM_CTXT_MEMBER_CURRENT_DRX_PARAMETER ((uint32_t) 1 << 12)
#define EMM_CTXT_MEMBER_PENDING_DRX_PARAMETER ((uint32_t) 1 << 13)
#define EMM_CTXT_MEMBER_EPS_BEARER_CONTEXT_STATUS ((uint32_t) 1 << 14)
#define EMM_CTXT_MEMBER_MOB_STATION_CLSMARK2 ((uint32_t) 1 << 15)
#define EMM_CTXT_MEMBER_UE_ADDITIONAL_SECURITY_CAPABILITY ((uint32_t) 1 << 16)

#define EMM_CTXT_MEMBER_AUTH_VECTOR0 ((uint32_t) 1 << 26)
  //#define           EMM_CTXT_MEMBER_AUTH_VECTOR1                 ((uint32_t)1
  //<< 27)  // reserved bit for AUTH VECTOR #define EMM_CTXT_MEMBER_AUTH_VECTOR2
  //((uint32_t)1 << 28)  // reserved bit for AUTH VECTOR #define
  // EMM_CTXT_MEMBER_AUTH_VECTOR3                 ((uint32_t)1 << 29)  //
  // reserved bit for AUTH VECTOR #define           EMM_CTXT_MEMBER_AUTH_VECTOR4
  //((uint32_t)1 << 30)  // reserved bit for AUTH VECTOR #define
  // EMM_CTXT_MEMBER_AUTH_VECTOR5                 ((uint32_t)1 << 31)  //
  // reserved bit for AUTH VECTOR

#define EMM_CTXT_MEMBER_SET_BIT(eMmCtXtMemBeRmAsK, bIt)                        \
  do {                                                                         \
    (eMmCtXtMemBeRmAsK) |= bIt;                                                \
  } while (0)
#define EMM_CTXT_MEMBER_CLEAR_BIT(eMmCtXtMemBeRmAsK, bIt)                      \
  do {                                                                         \
    (eMmCtXtMemBeRmAsK) &= ~bIt;                                               \
  } while (0)

  // imsi present mean we know it but was not checked with identity proc, or was
  // not provided in initial message
  imsi_t _imsi;     /* The IMSI provided by the UE or the MME, set valid when
                       identification returns IMSI */
  imsi64_t _imsi64; /* The IMSI provided by the UE or the MME, set valid when
                       identification returns IMSI */
  imsi64_t saved_imsi64; /* Useful for 5.4.2.7.c */
  imei_t _imei;          /* The IMEI provided by the UE                     */
  imeisv_t _imeisv;      /* The IMEISV provided by the UE                   */
  // bool                   _guti_is_new; /* The GUTI assigned to the UE is new
  // */
  guti_t _guti;         /* The GUTI assigned to the UE                     */
  guti_t _old_guti;     /* The old GUTI (GUTI REALLOCATION)                */
  tai_list_t _tai_list; /* TACs the the UE is registered to                */
  tai_t _lvr_tai;
  tai_t originating_tai;

  MobileStationClassmark2
      _mob_st_clsMark2; /* Mobile station classmark2 provided by the UE */
  ksi_t ksi;            /*key set identifier  */
  ue_network_capability_t _ue_network_capability;
  ms_network_capability_t _ms_network_capability;
  ue_additional_security_capability_t ue_additional_security_capability;
  drx_parameter_t _drx_parameter;

  int remaining_vectors;                       // remaining unused vectors
  auth_vector_t _vector[MAX_EPS_AUTH_VECTORS]; /* EPS authentication vector */
  emm_security_context_t
      _security; /* Current EPS security context: The security context which has
                    been activated most recently. Note that a current EPS
                                                          security context
                    originating from either a mapped or native EPS security
                    context may exist simultaneously with a native non-current
                    EPS security context.*/

  // Requirement MME24.301R10_4.4.2.1_2
  emm_security_context_t
      _non_current_security; /* Non-current EPS security context: A native EPS
                                security context that is not the current one. A
                                non-current EPS security context may be stored
                                along with a current EPS security context in the
                                UE and the MME. A non-current EPS security
                                context does not contain an EPS AS security
                                context. A non-current EPS security context is
                                either of type 'full
                                                          native' or of type
                                'partial native'.     */

  int emm_cause; /* EMM failure cause code                          */

  emm_fsm_state_t _emm_fsm_state;

  nas_timer_t T3422; /* Detach timer         */
  void* t3422_arg;

  struct esm_context_s esm_ctx;

  ue_network_capability_t
      tau_ue_network_capability; /* stored TAU Request IE Requirement
                                    MME24.301R10_5.5.3.2.4_4*/
  ms_network_capability_t
      tau_ms_network_capability;          /* stored TAU Request IE Requirement
                                             MME24.301R10_5.5.3.2.4_4*/
  drx_parameter_t _current_drx_parameter; /* stored TAU Request IE Requirement
                                             MME24.301R10_5.5.3.2.4_4*/
  eps_bearer_context_status_t
      _eps_bearer_context_status; /* stored TAU Request IE Requirement
                                     MME24.301R10_5.5.3.2.4_5*/
  eps_network_feature_support_t _eps_network_feature_support;

  // TODO: DO BETTER  WITH BELOW
  bstring esm_msg; /* ESM message contained within the initial request*/
#define EMM_CN_SAP_BUFFER_SIZE 4096

#define IS_EMM_CTXT_PRESENT_IMSI(eMmCtXtPtR)                                   \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_IMSI))
#define IS_EMM_CTXT_PRESENT_IMEI(eMmCtXtPtR)                                   \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_IMEI))
#define IS_EMM_CTXT_PRESENT_IMEISV(eMmCtXtPtR)                                 \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_IMEI_SV))
#define IS_EMM_CTXT_PRESENT_OLD_GUTI(eMmCtXtPtR)                               \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_OLD_GUTI))
#define IS_EMM_CTXT_PRESENT_GUTI(eMmCtXtPtR)                                   \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_GUTI))
#define IS_EMM_CTXT_PRESENT_TAI_LIST(eMmCtXtPtR)                               \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_TAI_LIST))
#define IS_EMM_CTXT_PRESENT_LVR_TAI(eMmCtXtPtR)                                \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_LVR_TAI))
#define IS_EMM_CTXT_PRESENT_AUTH_VECTORS(eMmCtXtPtR)                           \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_AUTH_VECTORS))
#define IS_EMM_CTXT_PRESENT_SECURITY(eMmCtXtPtR)                               \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_SECURITY))
#define IS_EMM_CTXT_PRESENT_NON_CURRENT_SECURITY(eMmCtXtPtR)                   \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_NON_CURRENT_SECURITY))
#define IS_EMM_CTXT_PRESENT_UE_NETWORK_CAPABILITY(eMmCtXtPtR)                  \
  (!!((eMmCtXtPtR)->member_present_mask &                                      \
      EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE))
#define IS_EMM_CTXT_PRESENT_MS_NETWORK_CAPABILITY(eMmCtXtPtR)                  \
  (!!((eMmCtXtPtR)->member_present_mask &                                      \
      EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE))
#define IS_EMM_CTXT_PRESENT_UE_ADDITIONAL_SECURITY_CAPABILITY(eMmCtXtPtR)      \
  (!!((eMmCtXtPtR)->member_present_mask &                                      \
      EMM_CTXT_MEMBER_UE_ADDITIONAL_SECURITY_CAPABILITY))

#define IS_EMM_CTXT_PRESENT_AUTH_VECTOR(eMmCtXtPtR, KsI)                       \
  (!!((eMmCtXtPtR)->member_present_mask &                                      \
      ((EMM_CTXT_MEMBER_AUTH_VECTOR0) << KsI)))
#define IS_EMM_CTXT_PRESENT_MOB_STATION_CLSMARK2(eMmCtXtPtR)                   \
  (!!((eMmCtXtPtR)->member_present_mask & EMM_CTXT_MEMBER_MOB_STATION_CLSMARK2))

#define IS_EMM_CTXT_VALID_IMSI(eMmCtXtPtR)                                     \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_IMSI))
#define IS_EMM_CTXT_VALID_IMEI(eMmCtXtPtR)                                     \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_IMEI))
#define IS_EMM_CTXT_VALID_IMEISV(eMmCtXtPtR)                                   \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_IMEI_SV))
#define IS_EMM_CTXT_VALID_OLD_GUTI(eMmCtXtPtR)                                 \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_OLD_GUTI))
#define IS_EMM_CTXT_VALID_GUTI(eMmCtXtPtR)                                     \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_GUTI))
#define IS_EMM_CTXT_VALID_TAI_LIST(eMmCtXtPtR)                                 \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_TAI_LIST))
#define IS_EMM_CTXT_VALID_LVR_TAI(eMmCtXtPtR)                                  \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_LVR_TAI))
#define IS_EMM_CTXT_VALID_AUTH_VECTORS(eMmCtXtPtR)                             \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_AUTH_VECTORS))
#define IS_EMM_CTXT_VALID_SECURITY(eMmCtXtPtR)                                 \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_SECURITY))
#define IS_EMM_CTXT_VALID_NON_CURRENT_SECURITY(eMmCtXtPtR)                     \
  (!!((eMmCtXtPtR)->member_valid_mask & EMM_CTXT_MEMBER_NON_CURRENT_SECURITY))
#define IS_EMM_CTXT_VALID_UE_NETWORK_CAPABILITY(eMmCtXtPtR)                    \
  (!!((eMmCtXtPtR)->member_valid_mask &                                        \
      EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE))
#define IS_EMM_CTXT_VALID_MS_NETWORK_CAPABILITY(eMmCtXtPtR)                    \
  (!!((eMmCtXtPtR)->member_valid_mask &                                        \
      EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE))

#define IS_EMM_CTXT_VALID_AUTH_VECTOR(eMmCtXtPtR, KsI)                         \
  (!!((eMmCtXtPtR)->member_valid_mask &                                        \
      ((EMM_CTXT_MEMBER_AUTH_VECTOR0) << KsI)))

  // CSFB related parameters
  csfb_params_t csfbparams;
  // VOLTE parameters
  volte_params_t volte_params;
  bool is_imsi_only_detach;
  /* Set the flag if pcrf initiated bearer deact and UE is in Idle state
   *if this flag is set after receving service req, send detach
   */
  bool nw_init_bearer_deactv;
  new_attach_info_t* new_attach_info;
  bool initiate_identity_after_smc;
} emm_context_t;

/*
 * Structure of the EMM data
 * -------------------------
 */
typedef struct emm_data_s {
  /*
   * MME configuration
   * -----------------
   */
  mme_api_emm_config_t conf;
  /*
   * EMM contexts
   * ------------
   */
  // TODO LG REMOVE hash_table_ts_t             *ctx_coll_ue_id; // key is emm
  // ue id, data is struct emm_context_s
  // TODO LG REMOVE hash_table_uint64_ts_t      *ctx_coll_imsi;  // key is
  // imsi_t, data is emm ue id (unsigned int)
  // TODO LG REMOVE obj_hash_table_uint64_t     *ctx_coll_guti;  // key is guti,
  // data is emm ue id (unsigned int)
} emm_data_t;

typedef struct {
  unsigned int ue_id;
#define DETACH_REQ_COUNTER_MAX 5
  unsigned int retransmission_count;
  uint8_t detach_type;
} nw_detach_data_t;

mme_ue_s1ap_id_t emm_ctx_get_new_ue_id(const emm_context_t* const ctxt)
    __attribute__((nonnull));

void emm_ctx_unmark_common_procedure_running(
    emm_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_set_attribute_present(
    emm_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_attribute_present(
    emm_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_set_attribute_valid(
    emm_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_attribute_valid(
    emm_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_guti(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_guti(emm_context_t* const ctxt, guti_t* guti)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_guti(emm_context_t* const ctxt, guti_t* guti)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_old_guti(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_old_guti(emm_context_t* const ctxt, guti_t* guti)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_old_guti(emm_context_t* const ctxt, guti_t* guti)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_imsi(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_imsi(emm_context_t* const ctxt, imsi_t* imsi, imsi64_t imsi64)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_imsi(
    emm_context_t* const ctxt, imsi_t* imsi, imsi64_t imsi64)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_imei(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_imei(emm_context_t* const ctxt, imei_t* imei)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_imei(emm_context_t* const ctxt, imei_t* imei)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_imeisv(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_imeisv(emm_context_t* const ctxt, imeisv_t* imeisv)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_imeisv(emm_context_t* const ctxt, imeisv_t* imeisv)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_lvr_tai(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_lvr_tai(emm_context_t* const ctxt, tai_t* lvr_tai)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_lvr_tai(emm_context_t* const ctxt, tai_t* lvr_tai)
    __attribute__((nonnull)) __attribute__((flatten));

void emm_ctx_clear_auth_vectors(emm_context_t* const ctxt)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_auth_vector(emm_context_t* const ctxt, ksi_t eksi)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_security(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_security_type(emm_context_t* const ctxt, emm_sc_type_t sc_type)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_security_eksi(emm_context_t* const ctxt, ksi_t eksi)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_security_vector_index(emm_context_t* const ctxt)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_security_vector_index(
    emm_context_t* const ctxt, int vector_index) __attribute__((nonnull))
__attribute__((flatten));

void emm_ctx_clear_non_current_security(emm_context_t* const ctxt)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_non_current_security_vector_index(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_non_current_security_vector_index(
    emm_context_t* const ctxt, int vector_index) __attribute__((nonnull));

void emm_ctx_clear_ue_nw_cap(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_ue_nw_cap(
    emm_context_t* const ctxt,
    const ue_network_capability_t* const ue_nw_cap_ie) __attribute__((nonnull));
void emm_ctx_set_valid_ue_nw_cap(
    emm_context_t* const ctxt,
    const ue_network_capability_t* const ue_nw_cap_ie) __attribute__((nonnull));

void emm_ctx_clear_ms_nw_cap(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_ms_nw_cap(
    emm_context_t* const ctxt,
    const ms_network_capability_t* const ms_nw_cap_ie);
void emm_ctx_set_valid_ms_nw_cap(
    emm_context_t* const ctxt,
    const ms_network_capability_t* const ms_nw_cap_ie);
void emm_ctx_clear_mobile_station_clsMark2(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_mobile_station_clsMark2(
    emm_context_t* const ctxt, MobileStationClassmark2* mob_st_clsMark2)
    __attribute__((nonnull));

void emm_ctx_clear_drx_parameter(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_drx_parameter(emm_context_t* const ctxt, drx_parameter_t* drx)
    __attribute__((nonnull));
void emm_ctx_set_valid_drx_parameter(
    emm_context_t* const ctxt, drx_parameter_t* drx) __attribute__((nonnull));

void emm_ctx_clear_ue_additional_security_capability(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_ue_additional_security_capability(
    emm_context_t* const ctxt, ue_additional_security_capability_t* drx)
    __attribute__((nonnull));

void free_emm_ctx_memory(
    emm_context_t* const ctxt, const mme_ue_s1ap_id_t ue_id);

struct emm_context_s* emm_context_get(
    emm_data_t* emm_data, const mme_ue_s1ap_id_t ue_id);
struct emm_context_s* emm_context_get_by_imsi(
    emm_data_t* emm_data, imsi64_t imsi64);
struct emm_context_s* emm_context_get_by_guti(
    emm_data_t* emm_data, guti_t* guti);

int emm_context_upsert_imsi(emm_data_t* emm_data, struct emm_context_s* elm)
    __attribute__((nonnull));

void nas_start_T3450(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3450,
    time_out_t time_out_cb, void* timer_callback_args);
void nas_start_T3460(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3460,
    time_out_t time_out_cb, void* timer_callback_args);
void nas_start_T3470(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3470,
    time_out_t time_out_cb, void* timer_callback_args);
void nas_stop_T3450(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3450,
    void* timer_callback_args);
void nas_stop_T3460(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3460,
    void* timer_callback_args);
void nas_stop_T3470(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3470,
    void* timer_callback_args);
void nas_start_Ts6a_auth_info(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const Ts6a_auth_info,
    time_out_t time_out_cb, void* timer_callback_args);
void nas_stop_Ts6a_auth_info(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const Ts6a_auth_info,
    void* timer_callback_args);
void emm_init_context(
    struct emm_context_s* const emm_ctx, const bool init_esm_ctxt)
    __attribute__((nonnull));
void emm_context_free(struct emm_context_s* const emm_ctx)
    __attribute__((nonnull));
void emm_context_free_content(struct emm_context_s* const emm_ctx)
    __attribute__((nonnull));
void emm_context_free_content_except_key_fields(
    struct emm_context_s* const emm_ctx) __attribute__((nonnull));
void emm_context_dump(
    const struct emm_context_s* const elm_pP, const uint8_t indent_spaces,
    bstring bstr_dump) __attribute__((nonnull));

/*
 *  Detach Proc: Timer handler
 */
void detach_t3422_handler(void*, imsi64_t* imsi64);
/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/*
 * --------------------------------------------------------------------------
 *      EPS mobility management data (used within EMM only)
 * --------------------------------------------------------------------------
 */
extern emm_data_t _emm_data;

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* FILE_EMM_DATA_SEEN*/
