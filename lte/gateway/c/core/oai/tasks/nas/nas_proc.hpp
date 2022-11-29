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
Source      nas_proc.hpp

Version     0.1

Date        2012/09/20

Product     NAS stack

Subsystem   NAS main process

Author      Frederic Maurel, Lionel GAUTHIER

Description NAS procedure call manager

*****************************************************************************/

#pragma once

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/security_types.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/nas/as_message.h"
#include "lte/gateway/c/core/oai/include/nas/commonDef.h"
#include "lte/gateway/c/core/oai/include/nas/networkDef.h"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/sgs_messages_types.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_cnDef.hpp"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

void nas_proc_initialize(const mme_config_t* mme_config_p);

void nas_proc_cleanup(void);

/*
 * --------------------------------------------------------------------------
 *          NAS procedures triggered by the user
 * --------------------------------------------------------------------------
 */

/*
 * --------------------------------------------------------------------------
 *      NAS procedures triggered by the network
 * --------------------------------------------------------------------------
 */

status_code_e nas_proc_establish_ind(
    const mme_ue_s1ap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi, const as_cause_t as_cause,
    const s_tmsi_t s_tmsi, STOLEN_REF bstring* msg);

status_code_e nas_proc_dl_transfer_cnf(const mme_ue_s1ap_id_t ueid,
                                       const nas_error_code_t status,
                                       STOLEN_REF bstring* nas_msg);
status_code_e nas_proc_dl_transfer_rej(const mme_ue_s1ap_id_t ueid,
                                       const nas_error_code_t status,
                                       STOLEN_REF bstring* nas_msg);
status_code_e nas_proc_ul_transfer_ind(const mme_ue_s1ap_id_t ueid,
                                       const tai_t originating_tai,
                                       const ecgi_t cgi,
                                       STOLEN_REF bstring* msg);

/*
 * --------------------------------------------------------------------------
 *      NAS procedures triggered by the mme applicative layer
 * --------------------------------------------------------------------------
 */
status_code_e nas_proc_authentication_info_answer(
    mme_app_desc_t* mme_app_desc_p, s6a_auth_info_ans_t* ans);
status_code_e nas_proc_downlink_unitdata(
    itti_sgsap_downlink_unitdata_t* dl_unitdata);
status_code_e nas_proc_sgs_release_req(itti_sgsap_release_req_t* sgs_rel);
status_code_e nas_proc_cs_domain_mm_information_request(
    itti_sgsap_mm_information_req_t* const mm_information_req_pP);
status_code_e nas_proc_cs_respose_success(
    emm_cn_cs_response_success_t* nas_cs_response_success);
status_code_e nas_proc_pdn_disconnect_rsp(
    emm_cn_pdn_disconnect_rsp_t* emm_cn_pdn_disconnect_rsp);
status_code_e nas_proc_ula_or_csrsp_fail(
    emm_cn_ula_or_csrsp_fail_t* ula_or_csrsp_fail);
status_code_e nas_proc_create_dedicated_bearer(
    emm_cn_activate_dedicated_bearer_req_t* emm_cn_activate);
status_code_e nas_proc_implicit_detach_ue_ind(mme_ue_s1ap_id_t ue_id);
status_code_e nas_proc_delete_dedicated_bearer(
    emm_cn_deactivate_dedicated_bearer_req_t* emm_cn_deactivate);
status_code_e nas_proc_nw_initiated_detach_ue_request(
    emm_cn_nw_initiated_detach_ue_t* const nw_initiated_detach_p);
status_code_e nas_proc_ula_success(mme_ue_s1ap_id_t ue_id);
status_code_e nas_proc_cs_domain_location_updt_fail(
    SgsRejectCause_t cause, lai_t* lai, mme_ue_s1ap_id_t mme_ue_s1ap_id);
status_code_e nas_proc_cs_service_notification(mme_ue_s1ap_id_t ue_id,
                                               uint8_t paging_id, bstring cli);
status_code_e nas_proc_auth_param_res(mme_ue_s1ap_id_t ue_id,
                                      uint8_t nb_vectors,
                                      eutran_vector_t* vectors);
status_code_e nas_proc_auth_param_fail(mme_ue_s1ap_id_t ue_id,
                                       nas_cause_t cause);
int nas_proc_signalling_connection_rel_ind(mme_ue_s1ap_id_t ue_id);
int nas_proc_smc_fail(emm_cn_smc_fail_t* emm_cn_smc_fail);

void mme_ue_context_update_ue_sgs_vlr_reliable(mme_ue_s1ap_id_t mme_ue_s1ap_id,
                                               bool vlr_reliable);
