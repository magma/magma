/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#ifndef FILE_EMM_HEARDERS_SEEN
#define FILE_EMM_HEADERS_SEEN

/*TODO: This file has temporary function declarations to
 * resolve undefined references. Delete
 * this file after moving all the files to c++
 * GH issue: https://github.com/magma/magma/issues/13096
 */
#include <sys/types.h>

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/nas/securityDef.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/queue.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/lib/hashtable/obj_hashtable.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_fsm.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsBearerContextStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsNetworkFeatureSupport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MobileStationClassmark2.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeNetworkCapability.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/
void nas_start_T3422(const mme_ue_s1ap_id_t ue_id,
                     struct nas_timer_s* const T3422, time_out_t time_out_cb);
void nas_stop_T3460(const mme_ue_s1ap_id_t ue_id,
                    struct nas_timer_s* const T3460);
void nas_stop_T3470(const mme_ue_s1ap_id_t ue_id,
                    struct nas_timer_s* const T3470);
void nas_stop_T3422(const imsi64_t imsi64, struct nas_timer_s* const T3422);

struct emm_context_s* emm_context_get(emm_data_t* emm_data,
                                      const mme_ue_s1ap_id_t ue_id);
void free_emm_detach_request_ies(emm_detach_request_ies_t** const ies);
void free_emm_tau_request_ies(emm_tau_request_ies_t** const ies);
void free_emm_attach_request_ies(emm_attach_request_ies_t** const params);
void nas_start_T3450(const mme_ue_s1ap_id_t ue_id,
                     struct nas_timer_s* const T3450, time_out_t time_out_cb);
void nas_start_T3460(const mme_ue_s1ap_id_t ue_id,
                     struct nas_timer_s* const T3460, time_out_t time_out_cb);
void nas_start_T3470(const mme_ue_s1ap_id_t ue_id,
                     struct nas_timer_s* const T3470, time_out_t time_out_cb);
status_code_e emm_context_upsert_imsi(emm_data_t* emm_data,
                                      struct emm_context_s* elm)
    __attribute__((nonnull));
bool is_nas_common_procedure_identification_running(
    const struct emm_context_s* const ctxt);

nas_emm_ident_proc_t* get_nas_common_procedure_identification(
    const struct emm_context_s* const ctxt);
bool is_nas_attach_reject_sent(const nas_emm_attach_proc_t* const attach_proc);
bool is_nas_specific_procedure_attach_running(
    const struct emm_context_s* const ctxt);
void emm_ctx_set_mobile_station_clsMark2(
    emm_context_t* const ctxt, MobileStationClassmark2* mob_st_clsMark2)
    __attribute__((nonnull));
void emm_ctx_set_ue_additional_security_capability(
    emm_context_t* const ctxt, ue_additional_security_capability_t* drx)
    __attribute__((nonnull));
void emm_ctx_set_guti(emm_context_t* const ctxt, guti_t* guti)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_attribute_valid(emm_context_t* const ctxt,
                                 const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_ms_nw_cap(
    emm_context_t* const ctxt,
    const ms_network_capability_t* const ms_nw_cap_ie);
void emm_ctx_set_valid_imsi(emm_context_t* const ctxt, imsi_t* imsi,
                            imsi64_t imsi64) __attribute__((nonnull))
__attribute__((flatten));
bool is_nas_specific_procedure_attach_running(
    const struct emm_context_s* const ctxt);
struct emm_context_s* emm_context_get_by_imsi(emm_data_t* emm_data,
                                              imsi64_t imsi64);
void emm_ctx_clear_ms_nw_cap(emm_context_t* const ctxt)
    __attribute__((nonnull));
void emm_ctx_set_valid_drx_parameter(emm_context_t* const ctxt,
                                     drx_parameter_t* drx);
status_code_e emm_proc_emm_information(ue_mm_context_t* emm_ctx);
nas_emm_attach_proc_t* nas_new_attach_procedure(
    struct emm_context_s* const emm_context);
emm_fsm_state_t emm_fsm_get_state(
    const struct emm_context_s* const emm_context);
void emm_ctx_set_valid_lvr_tai(emm_context_t* const ctxt, tai_t* lvr_tai)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_valid_imei(emm_context_t* const ctxt, imei_t* imei)
    __attribute__((nonnull)) __attribute__((flatten));
void nas_start_Ts6a_auth_info(const mme_ue_s1ap_id_t ue_id,
                              struct nas_timer_s* const Ts6a_auth_info,
                              time_out_t time_out_cb);
void emm_ctx_set_attribute_present(emm_context_t* const ctxt,
                                   const int attribute_bit_pos)
    __attribute__((nonnull));
void emm_ctx_clear_auth_vectors(emm_context_t* const ctxt)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_security_eksi(emm_context_t* const ctxt, ksi_t eksi)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_old_guti(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_clear_imsi(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_clear_imei(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_clear_security(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_clear_non_current_security(emm_context_t* const ctxt)
    __attribute__((nonnull)) __attribute__((flatten));
nas_emm_ident_proc_t* nas_new_identification_procedure(
    struct emm_context_s* const emm_context);
void emm_ctx_clear_guti(emm_context_t* const ctxt) __attribute__((nonnull))
__attribute__((flatten));
void emm_ctx_set_valid_imeisv(emm_context_t* const ctxt, imeisv_t* imeisv)
    __attribute__((nonnull)) __attribute__((flatten));
void emm_ctx_set_security_vector_index(emm_context_t* const ctxt,
                                       int vector_index)
    __attribute__((nonnull)) __attribute__((flatten));
nas_emm_smc_proc_t* get_nas_common_procedure_smc(
    const struct emm_context_s* const ctxt);
nas_emm_smc_proc_t* nas_new_smc_procedure(
    struct emm_context_s* const emm_context);
nas_emm_tau_proc_t* get_nas_specific_procedure_tau(
    const struct emm_context_s* const ctxt);
nas_emm_tau_proc_t* nas_new_tau_procedure(
    struct emm_context_s* const emm_context);
nas_auth_info_proc_t* get_nas_cn_procedure_auth_info(
    const struct emm_context_s* const ctxt);
int nas_message_encode(unsigned char* buffer, const nas_message_t* const msg,
                       size_t length, void* security);
void nas_emm_procedure_register_emm_message(mme_ue_s1ap_id_t ue_id,
                                            const uint64_t puid,
                                            bstring nas_msg);
int nas_message_encrypt(const unsigned char* inbuf, unsigned char* outbuf,
                        const nas_message_security_header_t* header,
                        size_t length, void* security);
int nas_message_decrypt(const unsigned char* const inbuf,
                        unsigned char* const outbuf,
                        nas_message_security_header_t* header, size_t length,
                        void* security, nas_message_decode_status_t* status);

int nas_message_decode(const unsigned char* const buffer, nas_message_t* msg,
                       size_t length, void* security,
                       nas_message_decode_status_t* status);
status_code_e emm_proc_status_ind(mme_ue_s1ap_id_t ue_id,
                                  emm_cause_t emm_cause);
status_code_e emm_proc_status(mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause);
void set_callbacks_for_attach_proc(nas_emm_attach_proc_t* attach_proc);
void free_emm_tau_request_ies(emm_tau_request_ies_t** const ies);
void set_callbacks_for_auth_proc(nas_emm_auth_proc_t* auth_proc);
void free_emm_detach_request_ies(emm_detach_request_ies_t** const ies);
status_code_e emm_proc_emm_information(ue_mm_context_t* emm_ctx);
void set_callbacks_for_auth_info_proc(nas_auth_info_proc_t* auth_info_proc);
void set_notif_callbacks_for_auth_proc(nas_emm_auth_proc_t* auth_proc);
void set_callbacks_for_smc_proc(nas_emm_smc_proc_t* smc_proc);
void set_notif_callbacks_for_smc_proc(nas_emm_smc_proc_t* smc_proc);

#ifdef __cplusplus
extern "C" {
#endif
void free_emm_ctx_memory(emm_context_t* const ctxt,
                         const mme_ue_s1ap_id_t ue_id);

void clear_emm_ctxt(emm_context_t* emm_ctx);
void emm_ctx_set_valid_ue_nw_cap(
    emm_context_t* const ctxt,
    const ue_network_capability_t* const ue_nw_cap_ie) __attribute__((nonnull));

bool is_nas_attach_accept_sent(const nas_emm_attach_proc_t* const attach_proc);

void nas_stop_T3450(const mme_ue_s1ap_id_t ue_id,
                    struct nas_timer_s* const T3450);
void emm_init_context(struct emm_context_s* const emm_ctx,
                      const bool init_esm_ctxt) __attribute__((nonnull));

void emm_ctx_clear_ue_nw_cap(emm_context_t* const ctxt)
    __attribute__((nonnull));
#ifdef __cplusplus
}
#endif

bool is_nas_common_procedure_authentication_running(
    const struct emm_context_s* const ctxt);
status_code_e nas_timer_init(void);
void nas_timer_cleanup(void);
void emm_ctx_set_security_type(emm_context_t* const ctxt, emm_sc_type_t sc_type)
    __attribute__((nonnull)) __attribute__((flatten));
#endif /* FILE_EMM_HEADERS_SEEN*/
