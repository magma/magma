/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file nas_itti_messaging.h
   \brief
   \author  Sebastien ROUX, Lionel GAUTHIER
   \date
   \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_ITTI_MESSAGING_SEEN
#define FILE_NAS_ITTI_MESSAGING_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "nas_message.h"
#include "nas/as_message.h"
#include "esm_proc.h"
#include "MobileStationClassmark2.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "common_types.h"

int nas_itti_dl_data_req(
  const mme_ue_s1ap_id_t ue_idP,
  bstring nas_msgP,
  nas_error_code_t transaction_status);

int nas_itti_erab_setup_req(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  const bitrate_t mbr_dl,
  const bitrate_t mbr_ul,
  const bitrate_t gbr_dl,
  const bitrate_t gbr_ul,
  bstring nas_msg);

void nas_itti_pdn_config_req(
  int ptiP,
  unsigned int ue_idP,
  const imsi_t *const imsi_pP,
  esm_proc_data_t *proc_data_pP,
  esm_proc_pdn_request_t request_typeP);

void nas_itti_pdn_connectivity_req(
  int ptiP,
  const mme_ue_s1ap_id_t ue_idP,
  const pdn_cid_t pdn_cidP,
  const imsi_t *const imsi_pP,
  imeisv_t imeisv,
  esm_proc_data_t *proc_data_pP,
  esm_proc_pdn_request_t request_typeP);

void nas_itti_auth_info_req(
  const mme_ue_s1ap_id_t ue_idP,
  const imsi_t *const imsiP,
  const bool is_initial_reqP,
  plmn_t *const visited_plmnP,
  const uint8_t num_vectorsP,
  const_bstring const auts_pP);

void nas_itti_establish_rej(
  const mme_ue_s1ap_id_t ue_idP,
  const imsi_t *const imsi_pP,
  uint8_t initial_reqP);

void nas_itti_establish_cnf(
  const mme_ue_s1ap_id_t ue_idP,
  const nas_error_code_t error_codeP,
  bstring msgP,
  const uint16_t selected_encryption_algorithmP,
  const uint16_t selected_integrity_algorithmP,
  const uint8_t csfb_response,
  const uint8_t presencemask,
  const uint8_t service_type);

void nas_itti_detach_req(const mme_ue_s1ap_id_t ue_idP);

void nas_itti_dedicated_eps_bearer_complete(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP);

void nas_itti_extended_service_req(
  const mme_ue_s1ap_id_t ue_id,
  const uint8_t servicetype,
  uint8_t csfb_response);

void nas_itti_dedicated_eps_bearer_reject(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP);

void nas_itti_sgsap_uplink_unitdata(
  const char *const imsi,
  uint8_t imsi_len,
  bstring nas_msg,
  imeisv_t *imeisv,
  MobileStationClassmark2 *mobilestationclassmark2,
  tai_t *tai,
  ecgi_t *ecgi);

void nas_itti_sgsap_downlink_unitdata(
  const char *const imsi,
  uint8_t imsi_len,
  const_bstring nas_msg);

void nas_itti_sgsap_tmsi_reallocation_comp(
  const char *imsi,
  const unsigned int imsi_len);

void nas_itti_cs_domain_location_update_req(
  unsigned int ue_idP,
  uint8_t msg_type);

void nas_itti_tau_complete(unsigned int ue_idP);

void nas_itti_sgs_detach_req(const uint32_t ue_idP, const uint8_t detach_type);

void nas_itti_sgsap_tmsi_reallocation_comp(
  const char *imsi,
  const unsigned int imsi_len);

void s6a_auth_info_rsp_timer_expiry_handler(void *args);

void nas_itti_sgsap_ue_activity_ind(
  const char *imsi,
  const unsigned int imsi_len);

void nas_itti_deactivate_eps_bearer_context(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP,
  bool delete_default_bearer,
  teid_t s_gw_teid_s11_s4);

int nas_itti_erab_rel_cmd(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  bstring nas_msg);

void nas_itti_dedicated_eps_bearer_deactivation_reject(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP,
  bool delete_default_bearer,
  teid_t s_gw_teid_s11_s4);

#endif /* FILE_NAS_ITTI_MESSAGING_SEEN */
