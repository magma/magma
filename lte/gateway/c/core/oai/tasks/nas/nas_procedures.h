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

/*! \file nas_procedures.h
   \brief
   \author  Lionel GAUTHIER
   \date 2017
   \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_PROCEDURES_SEEN
#define FILE_NAS_PROCEDURES_SEEN
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "3gpp_24.008.h"
#include "3gpp_23.003.h"
#include "3gpp_33.401.h"
#include "3gpp_36.401.h"
#include "bstrlib.h"
#include "common_types.h"
#include "emm_fsm.h"
#include "nas_timer.h"
#include "queue.h"
#include "nas/securityDef.h"
#include "security_types.h"

struct emm_context_s;
struct nas_base_proc_s;
struct nas_emm_proc_s;

struct emm_context_s;
struct nas_base_proc_s;
struct nas_emm_proc_s;

typedef int (*success_cb_t)(struct emm_context_s*);
typedef int (*failure_cb_t)(struct emm_context_s*);
typedef int (*proc_abort_t)(struct emm_context_s*, struct nas_base_proc_s*);

typedef int (*pdu_in_resp_t)(
    struct emm_context_s*,
    void* arg);  // can be RESPONSE, COMPLETE, ACCEPT
typedef int (*pdu_in_rej_t)(struct emm_context_s*, void* arg);  // REJECT.
typedef int (*pdu_out_rej_t)(
    struct emm_context_s*, struct nas_base_proc_s*);  // REJECT.
typedef void (*time_out_t)(void* arg, imsi64_t* imsi64);

typedef int (*sdu_out_delivered_t)(
    struct emm_context_s*, struct nas_emm_proc_s*);
typedef int (*sdu_out_not_delivered_t)(
    struct emm_context_s*, struct nas_emm_proc_s*);
typedef int (*sdu_out_not_delivered_ho_t)(
    struct emm_context_s*, struct nas_emm_proc_s*);

typedef enum {
  NAS_PROC_TYPE_NONE = 0,
  NAS_PROC_TYPE_EMM,
  NAS_PROC_TYPE_ESM,
  NAS_PROC_TYPE_CN,
} nas_base_proc_type_t;

typedef struct nas_base_proc_s {
  success_cb_t success_notif;
  failure_cb_t failure_notif;
  proc_abort_t abort;

  // PDU interface
  // pdu_in_resp_t           resp_in;
  pdu_in_rej_t fail_in;
  pdu_out_rej_t fail_out;
  time_out_t time_out;
  nas_base_proc_type_t type;  // EMM, ESM, CN
  uint64_t nas_puid;          // procedure unique identifier for internal use

  struct nas_base_proc_s* parent;
  struct nas_base_proc_s* child;
} nas_base_proc_t;

////////////////////////////////////////////////////////////////////////////////
// EMM procedures
////////////////////////////////////////////////////////////////////////////////
//------------------------------------------------------------------------------

typedef enum {
  NAS_EMM_PROC_TYPE_NONE = 0,
  NAS_EMM_PROC_TYPE_SPECIFIC,
  NAS_EMM_PROC_TYPE_COMMON,
  NAS_EMM_PROC_TYPE_CONN_MNGT,
} nas_emm_proc_type_t;

// EMM Specific procedures
typedef struct nas_emm_proc_s {
  nas_base_proc_t base_proc;
  nas_emm_proc_type_t type;  // specific, common, connection management
  // SDU interface
  sdu_out_delivered_t delivered;
  sdu_out_not_delivered_t not_delivered;
  sdu_out_not_delivered_ho_t not_delivered_ho;

  emm_fsm_state_t previous_emm_fsm_state;
} nas_emm_proc_t;

typedef enum {
  EMM_SPEC_PROC_TYPE_NONE = 0,
  EMM_SPEC_PROC_TYPE_ATTACH,
  EMM_SPEC_PROC_TYPE_DETACH,
  EMM_SPEC_PROC_TYPE_TAU,
} emm_specific_proc_type_t;

// EMM Specific procedures
typedef struct nas_emm_specific_proc_s {
  nas_emm_proc_t emm_proc;
  emm_specific_proc_type_t type;
} nas_emm_specific_proc_t;

struct emm_attach_request_ies_s;

typedef struct nas_emm_attach_proc_s {
  nas_emm_specific_proc_t emm_spec_proc;
  struct nas_timer_s T3450;  // EMM message retransmission timer
#define ATTACH_COUNTER_MAX 5
  int attach_accept_sent;
  bool attach_reject_sent;
  bool attach_complete_received;
  guti_t guti;
  bstring
      esm_msg_out;  // ESM message to be sent within the Attach Accept message
  struct emm_attach_request_ies_s* ies;
  mme_ue_s1ap_id_t ue_id;
  ksi_t ksi;
  int emm_cause;
} nas_emm_attach_proc_t;

struct emm_detach_request_ies_s;

typedef struct nas_emm_detach_proc_s {
  nas_emm_specific_proc_t emm_spec_proc;
  struct emm_detach_request_ies_s* ies;
} nas_emm_detach_proc_t;

struct emm_tau_request_ies_s;

typedef struct nas_emm_tau_proc_s {
  nas_emm_specific_proc_t emm_spec_proc;
  struct nas_timer_s T3450;  // EMM message retransmission timer
#define TAU_COUNTER_MAX 5
  unsigned int retransmission_count; /* Retransmission counter   */
  bstring
      esm_msg_out;  // ESM message to be sent within the Attach Accept message
  struct emm_tau_request_ies_s* ies;
  mme_ue_s1ap_id_t ue_id;
  int emm_cause;
} nas_emm_tau_proc_t;

//------------------------------------------------------------------------------

typedef enum {
  EMM_COMM_PROC_NONE = 0,
  EMM_COMM_PROC_GUTI,
  EMM_COMM_PROC_AUTH,
  EMM_COMM_PROC_SMC,
  EMM_COMM_PROC_IDENT,
  EMM_COMM_PROC_INFO,
} emm_common_proc_type_t;

// EMM Common procedures
typedef struct nas_emm_common_proc_s {
  nas_emm_proc_t emm_proc;
  emm_common_proc_type_t type;
} nas_emm_common_proc_t;

typedef struct nas_emm_guti_proc_s {
  nas_emm_common_proc_t emm_com_proc;
} nas_emm_guti_proc_t;

typedef struct nas_emm_ident_proc_s {
  nas_emm_common_proc_t emm_com_proc;
  struct nas_timer_s T3470; /* Identification timer         */
#define IDENTIFICATION_COUNTER_MAX 5
  unsigned int retransmission_count;
  mme_ue_s1ap_id_t ue_id;
  bool is_cause_is_attach;  //  could also be done by seeking parent procedure
  identity_type2_t identity_type;
} nas_emm_ident_proc_t;

typedef struct nas_emm_auth_proc_s {
  nas_emm_common_proc_t emm_com_proc;
  struct nas_timer_s T3460; /* Authentication timer         */
#define AUTHENTICATION_COUNTER_MAX 5
  unsigned int retransmission_count;
#define EMM_AUTHENTICATION_SYNC_FAILURE_MAX 2
  unsigned int
      sync_fail_count; /* counter of successive AUTHENTICATION FAILURE messages
                          from the UE with EMM cause #21 "synch failure" */
  unsigned int mac_fail_count;
  mme_ue_s1ap_id_t ue_id;
  bool is_cause_is_attach;  //  could also be done by seeking parent procedure
  ksi_t ksi;
  uint8_t rand[AUTH_RAND_SIZE]; /* Random challenge number  */
  uint8_t autn[AUTH_AUTN_SIZE]; /* Authentication token     */
  imsi_t* unchecked_imsi;
  int emm_cause;
} nas_emm_auth_proc_t;

typedef struct nas_emm_smc_proc_s {
  nas_emm_common_proc_t emm_com_proc;
  struct nas_timer_s T3460; /* Authentication timer         */
#define SECURITY_COUNTER_MAX 5
  mme_ue_s1ap_id_t ue_id;
  unsigned int retransmission_count; /* Retransmission counter    */
  int ksi;                           /* NAS key set identifier                */
  int eea;                           /* Replayed EPS encryption algorithms    */
  int eia;                           /* Replayed EPS integrity algorithms     */
  int ucs2;                          /* Replayed Alphabet                     */
  int uea;                           /* Replayed UMTS encryption algorithms   */
  int uia;                           /* Replayed UMTS integrity algorithms    */
  int gea;                           /* Replayed G encryption algorithms      */
  bool umts_present;
  bool gprs_present;
  int selected_eea;        /* Selected EPS encryption algorithms    */
  int selected_eia;        /* Selected EPS integrity algorithms     */
  int saved_selected_eea;  /* Previous selected EPS encryption algorithms    */
  int saved_selected_eia;  /* Previous selected EPS integrity algorithms     */
  int saved_eksi;          /* Previous ksi     */
  uint16_t saved_overflow; /* Previous dl_count overflow     */
  uint8_t saved_seq_num;   /* Previous dl_count seq_num     */
  emm_sc_type_t saved_sc_type;
  bool notify_failure; /* Indicates whether the identification
                        * procedure failure shall be notified
                        * to the ongoing EMM procedure */
  bool is_new;         /* new security context for SMC header type */
  bool imeisv_request;
  bool replayed_ue_add_sec_cap_present;
  uint16_t _5g_ea; /* Replayed 5GS encryption algorithms */
  uint16_t _5g_ia; /* Replayed 5GS integrity algorithms */
} nas_emm_smc_proc_t;

typedef struct nas_emm_info_proc_s {
  nas_emm_common_proc_t emm_com_proc;
} nas_emm_info_proc_t;

//------------------------------------------------------------------------------
typedef enum {
  EMM_CON_MNGT_PROC_NONE = 0,
  EMM_CON_MNGT_PROC_PAGING,
  EMM_CON_MNGT_PROC_SERVICE_REQUEST
} emm_con_mngt_proc_type_t;

// EMM Connection management procedures
typedef struct nas_emm_con_mngt_proc_s {
  nas_emm_proc_t emm_proc;
  emm_con_mngt_proc_type_t type;
} nas_emm_con_mngt_proc_t;

typedef struct nas_sr_proc_s {
  nas_emm_con_mngt_proc_t con_mngt_proc;
  mme_ue_s1ap_id_t ue_id;
  int emm_cause;
} nas_sr_proc_t;

////////////////////////////////////////////////////////////////////////////////
// ESM procedures
////////////////////////////////////////////////////////////////////////////////
typedef enum {
  ESM_PROC_NONE = 0,
  ESM_PROC_EPS_BEARER_CONTEXT,
  ESM_PROC_TRANSACTION
} esm_proc_type_t;

typedef struct nas_esm_proc_s {
  nas_base_proc_t base_proc;
  esm_proc_type_t type;
} nas_esm_proc_t;

typedef enum {
  ESM_BEARER_CTX_PROC_NONE = 0,
  ESM_PROC_DEFAULT_EPS_BEARER_CTXT_ACTIVATION,
  ESM_PROC_DEDICATED_EPS_BEARER_CTXT_ACTIVATION,
  ESM_PROC_EPS_BEARER_CTXT_MODIFICATION,
  ESM_PROC_EPS_BEARER_CTXT_DEACTIVATION
} esm_bearer_ctx_proc_type_t;

typedef struct nas_esm_bearer_ctx_proc_s {
  nas_esm_proc_t esm_proc;
  esm_bearer_ctx_proc_type_t type;
} nas_esm_bearer_ctx_proc_t;

typedef enum {
  ESM_TRANSACTION_PROC_NONE = 0,
  ESM_PROC_TRANSACTION_PDN_CONNECTIVITY,
  ESM_PROC_TRANSACTION_PDN_DISCONNECT,
  ESM_PROC_TRANSACTION_BEARER_RESOURCE_ALLOCATION,
  ESM_PROC_TRANSACTION_BEARER_RESOURCE_MODIFICATION,
} esm_transaction_proc_type_t;

typedef struct nas_esm_transaction_proc_s {
  nas_esm_proc_t esm_proc;
  esm_transaction_proc_type_t type;
} nas_esm_transaction_proc_t;

////////////////////////////////////////////////////////////////////////////////
// CN procedures
////////////////////////////////////////////////////////////////////////////////
typedef enum {
  CN_PROC_NONE = 0,
  CN_PROC_AUTH_INFO,
} cn_proc_type_t;

typedef struct nas_cn_proc_s {
  nas_base_proc_t base_proc;
  cn_proc_type_t type;
} nas_cn_proc_t;

typedef struct nas_auth_info_proc_s {
  nas_cn_proc_t cn_proc;
  success_cb_t success_notif;
  failure_cb_t failure_notif;
  bool request_sent;
  uint8_t nb_vectors;
  eutran_vector_t* vector[MAX_EPS_AUTH_VECTORS];
  int nas_cause;
#define S6A_AIR_RESPONSE_TIMER 3  // In seconds
  struct nas_timer_s timer_s6a;
  mme_ue_s1ap_id_t ue_id;
  bool resync;  // Indicates whether the authentication information is requested
                // due to sync failure
} nas_auth_info_proc_t;

////////////////////////////////////////////////////////////////////////////////
// Mix them (kind of class hierarchy--)
////////////////////////////////////////////////////////////////////////////////
typedef union {
  nas_base_proc_t base_proc;
  nas_emm_proc_t emm_proc;
  nas_emm_specific_proc_t emm_specific_proc;
  nas_emm_attach_proc_t emm_attach_proc;
  nas_emm_detach_proc_t emm_detach_proc;
  nas_emm_tau_proc_t emm_tau_proc;
  nas_emm_common_proc_t emm_common_proc;
  nas_emm_guti_proc_t emm_guti_proc;
  nas_emm_ident_proc_t emm_ident_proc;
  nas_emm_auth_proc_t emm_auth_proc;
  nas_emm_smc_proc_t emm_smc_proc;
  nas_emm_info_proc_t emm_info_proc;
  nas_emm_con_mngt_proc_t emm_con_mngt_proc;
  nas_esm_proc_t esm_proc;
  nas_cn_proc_t cn_proc;
  nas_auth_info_proc_t auth_info;
} nas_proc_t;

typedef struct nas_emm_common_procedure_s {
  nas_emm_common_proc_t* proc;
  LIST_ENTRY(nas_emm_common_procedure_s) entries;
} nas_emm_common_procedure_t;

typedef struct nas_cn_procedure_s {
  nas_cn_proc_t* proc;
  LIST_ENTRY(nas_cn_procedure_s) entries;
} nas_cn_procedure_t;

typedef struct nas_proc_mess_sign_s {
  uint64_t puid;
#define NAS_MSG_DIGEST_SIZE 16
  uint8_t digest[NAS_MSG_DIGEST_SIZE];
  size_t digest_length;
  size_t nas_msg_length;
} nas_proc_mess_sign_t;

typedef struct emm_procedures_s {
  nas_emm_specific_proc_t* emm_specific_proc;
  LIST_HEAD(nas_emm_common_procedures_head_s, nas_emm_common_procedure_s)
  emm_common_procs;
  LIST_HEAD(nas_cn_procedures_head_s, nas_cn_procedure_s)
  cn_procs;  // triggered by EMM
  nas_emm_con_mngt_proc_t* emm_con_mngt_proc;

  int nas_proc_mess_sign_next_location;  // next index in array
#define MAX_NAS_PROC_MESS_SIGN 3
  nas_proc_mess_sign_t nas_proc_mess_sign[MAX_NAS_PROC_MESS_SIGN];
} emm_procedures_t;

bool is_nas_common_procedure_guti_realloc_running(
    const struct emm_context_s* const ctxt);
bool is_nas_common_procedure_authentication_running(
    const struct emm_context_s* const ctxt);
bool is_nas_common_procedure_smc_running(
    const struct emm_context_s* const ctxt);
bool is_nas_common_procedure_identification_running(
    const struct emm_context_s* const ctxt);

nas_emm_guti_proc_t* get_nas_common_procedure_guti_realloc(
    const struct emm_context_s* const ctxt);
nas_emm_ident_proc_t* get_nas_common_procedure_identification(
    const struct emm_context_s* const ctxt);
nas_emm_smc_proc_t* get_nas_common_procedure_smc(
    const struct emm_context_s* const ctxt);
nas_emm_auth_proc_t* get_nas_common_procedure_authentication(
    const struct emm_context_s* const ctxt);

nas_auth_info_proc_t* get_nas_cn_procedure_auth_info(
    const struct emm_context_s* const ctxt);

nas_sr_proc_t* get_nas_con_mngt_procedure_service_request(
    const struct emm_context_s* const ctxt);

bool is_nas_specific_procedure_attach_running(
    const struct emm_context_s* const ctxt);
bool is_nas_specific_procedure_detach_running(
    const struct emm_context_s* const ctxt);
bool is_nas_specific_procedure_tau_running(
    const struct emm_context_s* const ctxt);

nas_emm_attach_proc_t* get_nas_specific_procedure_attach(
    const struct emm_context_s* const ctxt);
nas_emm_detach_proc_t* get_nas_specific_procedure_detach(
    const struct emm_context_s* const ctxt);
nas_emm_tau_proc_t* get_nas_specific_procedure_tau(
    const struct emm_context_s* const ctxt);

bool is_nas_attach_accept_sent(const nas_emm_attach_proc_t* const attach_proc);
bool is_nas_attach_reject_sent(const nas_emm_attach_proc_t* const attach_proc);
bool is_nas_attach_complete_received(
    const nas_emm_attach_proc_t* const attach_proc);

int nas_unlink_procedures(
    nas_base_proc_t* const parent_proc, nas_base_proc_t* const child_proc);

void nas_delete_all_emm_procedures(struct emm_context_s* const emm_context);
void nas_delete_common_procedure(
    struct emm_context_s* const emm_context, nas_emm_common_proc_t** proc);
void nas_delete_attach_procedure(struct emm_context_s* const emm_context);
void nas_delete_tau_procedure(struct emm_context_s* emm_context);
void nas_delete_detach_procedure(struct emm_context_s* emm_context);

void nas_delete_cn_procedure(
    struct emm_context_s* emm_context, nas_cn_proc_t* cn_proc);

nas_emm_attach_proc_t* nas_new_attach_procedure(
    struct emm_context_s* const emm_context);
nas_emm_tau_proc_t* nas_new_tau_procedure(
    struct emm_context_s* const emm_context);
nas_sr_proc_t* nas_new_service_request_procedure(
    struct emm_context_s* const emm_context);
nas_emm_ident_proc_t* nas_new_identification_procedure(
    struct emm_context_s* const emm_context);
nas_emm_auth_proc_t* nas_new_authentication_procedure(
    struct emm_context_s* const emm_context);
nas_emm_smc_proc_t* nas_new_smc_procedure(
    struct emm_context_s* const emm_context);
nas_auth_info_proc_t* nas_new_cn_auth_info_procedure(
    struct emm_context_s* const emm_context);

void nas_digest_msg(
    const unsigned char* const msg, const size_t msg_len, char* const digest,
    /*INOUT*/ size_t* const digest_length);
void nas_emm_procedure_register_emm_message(
    mme_ue_s1ap_id_t ue_id, const uint64_t puid, bstring nas_msg);
nas_emm_proc_t* nas_emm_find_procedure_by_msg_digest(
    struct emm_context_s* const emm_context, const char* const digest,
    const size_t digest_bytes, const size_t msg_size);

#endif
