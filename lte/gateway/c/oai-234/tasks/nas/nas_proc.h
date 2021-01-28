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
Source      nas_proc.h

Version     0.1

Date        2012/09/20

Product     NAS stack

Subsystem   NAS main process

Author      Frederic Maurel, Lionel GAUTHIER

Description NAS procedure call manager

*****************************************************************************/
#ifndef FILE_NAS_PROC_SEEN
#define FILE_NAS_PROC_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "common_defs.h"
#include "mme_config.h"
#include "emm_cnDef.h"
#include "nas/commonDef.h"
#include "nas/networkDef.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"
#include "nas/as_message.h"
#include "bstrlib.h"
#include "mme_app_defs.h"
#include "s6a_messages_types.h"
#include "security_types.h"
#include "sgs_messages_types.h"

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

int nas_proc_establish_ind(
    const mme_ue_s1ap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi, const as_cause_t as_cause,
    const s_tmsi_t s_tmsi, STOLEN_REF bstring* msg);

int nas_proc_dl_transfer_cnf(
    const mme_ue_s1ap_id_t ueid, const nas_error_code_t status,
    STOLEN_REF bstring* nas_msg);
int nas_proc_dl_transfer_rej(
    const mme_ue_s1ap_id_t ueid, const nas_error_code_t status,
    STOLEN_REF bstring* nas_msg);
int nas_proc_ul_transfer_ind(
    const mme_ue_s1ap_id_t ueid, const tai_t originating_tai, const ecgi_t cgi,
    STOLEN_REF bstring* msg);

/*
 * --------------------------------------------------------------------------
 *      NAS procedures triggered by the mme applicative layer
 * --------------------------------------------------------------------------
 */
int nas_proc_authentication_info_answer(
    mme_app_desc_t* mme_app_desc_p, s6a_auth_info_ans_t* ans);
int nas_proc_auth_param_res(
    mme_ue_s1ap_id_t ue_id, uint8_t nb_vectors, eutran_vector_t* vectors);
int nas_proc_auth_param_fail(mme_ue_s1ap_id_t ue_id, nas_cause_t cause);
int nas_proc_ula_success(mme_ue_s1ap_id_t ue_id);
int nas_proc_cs_respose_success(
    emm_cn_cs_response_success_t* nas_cs_response_success);
int nas_proc_ula_or_csrsp_fail(emm_cn_ula_or_csrsp_fail_t* ula_or_csrsp_fail);
int nas_proc_create_dedicated_bearer(
    emm_cn_activate_dedicated_bearer_req_t* emm_cn_activate);
int nas_proc_signalling_connection_rel_ind(mme_ue_s1ap_id_t ue_id);
int nas_proc_implicit_detach_ue_ind(mme_ue_s1ap_id_t ue_id);
int nas_proc_smc_fail(emm_cn_smc_fail_t* emm_cn_smc_fail);
int nas_proc_nw_initiated_detach_ue_request(
    emm_cn_nw_initiated_detach_ue_t* const nw_initiated_detach_p);
int nas_proc_cs_domain_location_updt_fail(
    SgsRejectCause_t cause, lai_t* lai, mme_ue_s1ap_id_t mme_ue_s1ap_id);
int nas_proc_downlink_unitdata(itti_sgsap_downlink_unitdata_t* dl_unitdata);
int nas_proc_sgs_release_req(itti_sgsap_release_req_t* sgs_rel);
int nas_proc_cs_domain_mm_information_request(
    itti_sgsap_mm_information_req_t* const mm_information_req_pP);
int nas_proc_cs_service_notification(
    mme_ue_s1ap_id_t ue_id, uint8_t paging_id, bstring cli);
int nas_proc_delete_dedicated_bearer(
    emm_cn_deactivate_dedicated_bearer_req_t* emm_cn_deactivate);
int nas_proc_pdn_disconnect_rsp(
    emm_cn_pdn_disconnect_rsp_t* emm_cn_pdn_disconnect_rsp);
#endif /* FILE_NAS_PROC_SEEN*/
