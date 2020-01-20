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

#ifndef FILE_S1AP_MME_HANDLERS_SEEN
#define FILE_S1AP_MME_HANDLERS_SEEN
#include <stdbool.h>

#include "s1ap_ies_defs.h"
#include "s1ap_mme.h"
#include "intertask_interface.h"
#include "S1ap-Cause.h"
#include "common_types.h"
#include "s1ap_messages_types.h"
#include "sctp_messages_types.h"

struct s1ap_message_s;

#define MAX_NUM_PARTIAL_S1_CONN_RESET 256

const char *s1_enb_state2str(enum mme_s1_enb_state_s state);
const char *s1ap_direction2str(uint8_t dir);

/** \brief Handle decoded incoming messages from SCTP
 * \param assoc_id SCTP association ID
 * \param stream Stream number
 * \param message_p The message decoded by the ASN1C decoder
 * @returns int
 **/
int s1ap_mme_handle_message(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_ue_cap_indication(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

/** \brief Handle an S1 Setup request message.
 * Typically add the eNB in the list of served eNB if not present, simply reset
 * UEs association otherwise. S1SetupResponse message is sent in case of success or
 * S1SetupFailure if the MME cannot accept the configuration received.
 * \param assoc_id SCTP association ID
 * \param stream Stream number
 * \param message_p The message decoded by the ASN1C decoder
 * @returns int
 **/
int s1ap_mme_handle_s1_setup_request(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_path_switch_request(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_ue_context_release_request(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_handle_ue_context_release_command(
  s1ap_state_t *state,
  const itti_s1ap_ue_context_release_command_t
    *const ue_context_release_command_pP);

int s1ap_mme_handle_ue_context_release_complete(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_handle_ue_context_mod_req(
  s1ap_state_t *state,
  const itti_s1ap_ue_context_mod_req_t *const ue_context_mod_req_pP);

int s1ap_mme_handle_initial_context_setup_failure(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_initial_context_setup_response(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_handle_sctp_disconnection(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  bool reset);

int s1ap_handle_new_association(
  s1ap_state_t *state,
  sctp_new_peer_t *sctp_new_peer_p);

int s1ap_mme_set_cause(
  S1ap_Cause_t *cause_p,
  const S1ap_Cause_PR cause_type,
  const long cause_value);

int s1ap_mme_generate_s1_setup_failure(
  const sctp_assoc_id_t assoc_id,
  const S1ap_Cause_PR cause_type,
  const long cause_value,
  const long time_to_wait);

int s1ap_mme_handle_erab_setup_response(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

int s1ap_mme_handle_erab_setup_failure(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

void s1ap_mme_handle_ue_context_rel_comp_timer_expiry(
  s1ap_state_t *state,
  ue_description_t *ue_ref_p);

void s1ap_mme_release_ue_context(
  s1ap_state_t *state,
  ue_description_t *ue_ref_p);

int s1ap_mme_handle_error_ind_message(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

int s1ap_mme_handle_enb_reset(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

int s1ap_handle_enb_initiated_reset_ack(
  const itti_s1ap_enb_initiated_reset_ack_t *const enb_reset_ack_p);

void s1ap_enb_assoc_clean_up_timer_expiry(
  s1ap_state_t *state,
  enb_description_t *enb_ref_p);

int s1ap_handle_paging_request(
  s1ap_state_t* state,
  const itti_s1ap_paging_request_t* paging_request);

int s1ap_mme_handle_ue_context_modification_response(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_ue_context_modification_failure(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_mme_handle_erab_rel_response(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message);

int s1ap_mme_handle_enb_configuration_transfer(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message_p);

int s1ap_handle_path_switch_req_ack(
  s1ap_state_t *state,
  const itti_s1ap_path_switch_request_ack_t *path_switch_req_ack_p);

int s1ap_handle_path_switch_req_failure(
  s1ap_state_t *state,
  const itti_s1ap_path_switch_request_failure_t *path_switch_req_failure_p);

#endif /* FILE_S1AP_MME_HANDLERS_SEEN */
