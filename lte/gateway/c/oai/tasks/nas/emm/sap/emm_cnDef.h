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

Source      emm_cnDef.h

Version     0.1

Date        2013/12/05

Product     NAS stack

Subsystem   EPS Core Network

Author      Sebastien Roux, Lionel GAUTHIER

Description

*****************************************************************************/

#ifndef FILE_EMM_CNDEF_SEEN
#define FILE_EMM_CNDEF_SEEN

#include "intertask_interface_types.h"

#include "nas/as_message.h"
#include "common_ies.h"
#include "LocationAreaIdentification.h"

typedef enum emmcn_primitive_s {
  _EMMCN_START = 400,
  _EMMCN_AUTHENTICATION_PARAM_RES,
  _EMMCN_AUTHENTICATION_PARAM_FAIL,
  _EMMCN_ULA_SUCCESS,
  _EMMCN_CS_RESPONSE_SUCCESS,
  _EMMCN_ULA_OR_CSRSP_FAIL,
  _EMMCN_ACTIVATE_DEDICATED_BEARER_REQ,
  _EMMCN_IMPLICIT_DETACH_UE,
  _EMMCN_SMC_PROC_FAIL,
  _EMMCN_NW_INITIATED_DETACH_UE,
  _EMMCN_CS_DOMAIN_LOCATION_UPDT_ACC,
  _EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL,
  _EMMCN_CS_DOMAIN_MM_INFORMATION_REQ,
  _EMMCN_DEACTIVATE_BEARER_REQ,  // LG
  _EMMCN_PDN_DISCONNECT_RES,
  _EMMCN_END
} emm_cn_primitive_t;

typedef enum pdn_conn_rsp_cause_e {
  CAUSE_OK                             = 16,
  CAUSE_CONTEXT_NOT_FOUND              = 64,
  CAUSE_INVALID_MESSAGE_FORMAT         = 65,
  CAUSE_SERVICE_NOT_SUPPORTED          = 68,
  CAUSE_SYSTEM_FAILURE                 = 72,
  CAUSE_NO_RESOURCES_AVAILABLE         = 73,
  CAUSE_ALL_DYNAMIC_ADDRESSES_OCCUPIED = 84
} pdn_conn_rsp_cause_t;

typedef struct emm_cn_auth_res_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;

  /* For future use: nb of vectors provided */
  uint8_t nb_vectors;

  /* Consider only one E-UTRAN vector for the moment... */
  eutran_vector_t* vector[MAX_EPS_AUTH_VECTORS];
} emm_cn_auth_res_t;

typedef struct emm_cn_auth_fail_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;

  /* S6A mapped to NAS cause */
  nas_cause_t cause;
} emm_cn_auth_fail_t;

typedef struct emm_cn_ula_success_s {
  mme_ue_s1ap_id_t ue_id;  // nas ref
} emm_cn_ula_success_t;

/* emm_cn_ula_or_csrsp_fail_s is used for handling failed
 * Location update procedure and Create session procedure
 */
typedef struct emm_cn_ula_or_csrsp_fail_s {
  mme_ue_s1ap_id_t ue_id;
  int pti;
  pdn_conn_rsp_cause_t cause;
} emm_cn_ula_or_csrsp_fail_t;

typedef struct emm_cn_cs_response_success_s {
  pdn_cid_t pdn_cid;
  /* Identity of the procedure transaction executed to
   * activate the PDN connection enty
   */
  proc_tid_t pti;
  network_qos_t qos;
  protocol_configuration_options_t pco;
  bstring pdn_addr;
  int pdn_type;
  int request_type;
  mme_ue_s1ap_id_t ue_id;
  ambr_t ambr;
  ambr_t apn_ambr;
  unsigned ebi : 4;
  /* QoS */
  qci_t qci;
  priority_level_t prio_level;
  pre_emption_vulnerability_t pre_emp_vulnerability;
  pre_emption_capability_t pre_emp_capability;

  /* S-GW TEID and IP address for user-plane */
  fteid_t sgw_s1u_fteid;
} emm_cn_cs_response_success_t;

typedef struct emm_cn_activate_dedicated_bearer_req_s {
  mme_ue_s1ap_id_t ue_id;
  pdn_cid_t cid;
  ebi_t ebi;
  ebi_t linked_ebi;
  bearer_qos_t bearer_qos;
  traffic_flow_template_t* tft;
  protocol_configuration_options_t* pco;
  fteid_t sgw_fteid;
} emm_cn_activate_dedicated_bearer_req_t;

typedef struct emm_cn_deactivate_dedicated_bearer_req_s {
  uint32_t no_of_bearers;
  ebi_t ebi[BEARERS_PER_UE];  // EPS Bearer ID
  mme_ue_s1ap_id_t ue_id;
} emm_cn_deactivate_dedicated_bearer_req_t;

typedef struct emm_cn_pdn_disconnect_rsp_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  ebi_t lbi;  // Default EPS Bearer ID
} emm_cn_pdn_disconnect_rsp_t;

typedef struct emm_cn_deregister_ue_s {
  uint32_t ue_id;
} emm_cn_deregister_ue_t;

typedef struct emm_cn_implicit_detach_ue_s {
  uint32_t ue_id;
} emm_cn_implicit_detach_ue_t;

typedef struct emm_cn_smc_fail_s {
  mme_ue_s1ap_id_t ue_id;
  nas_cause_t emm_cause;
} emm_cn_smc_fail_t;

typedef struct emm_cn_nw_initiated_detach_ue_s {
  uint32_t ue_id;
#define HSS_INITIATED_EPS_DETACH 0x00
#define SGS_INITIATED_IMSI_DETACH 0x01
#define MME_INITIATED_EPS_DETACH 0x02
  uint8_t detach_type;
} emm_cn_nw_initiated_detach_ue_t;

typedef struct emm_cn_cs_domain_location_updt_fail_s {
#define LAI (1 << 0)
  uint8_t presencemask;
  mme_ue_s1ap_id_t ue_id;
  int reject_cause;
  lai_t laicsfb;
} emm_cn_cs_domain_location_updt_fail_t;

typedef itti_sgsap_mm_information_req_t emm_cn_cs_domain_mm_information_req_t;

typedef struct emm_mme_ul_s {
  emm_cn_primitive_t primitive;
  union {
    emm_cn_auth_res_t* auth_res;
    emm_cn_auth_fail_t* auth_fail;
    emm_cn_deregister_ue_t deregister;
    emm_cn_ula_success_t* emm_cn_ula_success;
    emm_cn_cs_response_success_t* emm_cn_cs_response_success;
    emm_cn_ula_or_csrsp_fail_t* emm_cn_ula_or_csrsp_fail;
    emm_cn_activate_dedicated_bearer_req_t* activate_dedicated_bearer_req;
    emm_cn_deactivate_dedicated_bearer_req_t* deactivate_dedicated_bearer_req;
    emm_cn_implicit_detach_ue_t emm_cn_implicit_detach;
    emm_cn_smc_fail_t* smc_fail;
    emm_cn_nw_initiated_detach_ue_t emm_cn_nw_initiated_detach;
    emm_cn_cs_domain_location_updt_fail_t emm_cn_cs_domain_location_updt_fail;
    emm_cn_cs_domain_mm_information_req_t* emm_cn_cs_domain_mm_information_req;
    emm_cn_pdn_disconnect_rsp_t* emm_cn_pdn_disconnect_rsp;
  } u;
} emm_cn_t;

#endif /* FILE_EMM_CNDEF_SEEN */
