/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */
// WARNING: Do not include this header directly. Use intertask_interface.h
// instead.

/*! \file s1ap_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

MESSAGE_DEF(S1AP_UE_CAPABILITIES_IND, itti_s1ap_ue_cap_ind_t, s1ap_ue_cap_ind)
MESSAGE_DEF(
    S1AP_ENB_DEREGISTERED_IND, itti_s1ap_eNB_deregistered_ind_t,
    s1ap_eNB_deregistered_ind)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_RELEASE_REQ, itti_s1ap_ue_context_release_req_t,
    s1ap_ue_context_release_req)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_RELEASE_COMMAND, itti_s1ap_ue_context_release_command_t,
    s1ap_ue_context_release_command)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_RELEASE_COMPLETE, itti_s1ap_ue_context_release_complete_t,
    s1ap_ue_context_release_complete)
MESSAGE_DEF(
    S1AP_NAS_DL_DATA_REQ, itti_s1ap_nas_dl_data_req_t, s1ap_nas_dl_data_req)
MESSAGE_DEF(
    S1AP_INITIAL_UE_MESSAGE, itti_s1ap_initial_ue_message_t,
    s1ap_initial_ue_message)
MESSAGE_DEF(
    S1AP_E_RAB_SETUP_REQ, itti_s1ap_e_rab_setup_req_t, s1ap_e_rab_setup_req)
MESSAGE_DEF(
    S1AP_E_RAB_SETUP_RSP, itti_s1ap_e_rab_setup_rsp_t, s1ap_e_rab_setup_rsp)
MESSAGE_DEF(
    S1AP_ENB_INITIATED_RESET_REQ, itti_s1ap_enb_initiated_reset_req_t,
    s1ap_enb_initiated_reset_req)
MESSAGE_DEF(
    S1AP_ENB_INITIATED_RESET_ACK, itti_s1ap_enb_initiated_reset_ack_t,
    s1ap_enb_initiated_reset_ack)
MESSAGE_DEF(
    S1AP_PAGING_REQUEST, itti_s1ap_paging_request_t, s1ap_paging_request)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST, itti_s1ap_ue_context_mod_req_t,
    s1ap_ue_context_mod_request)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_MODIFICATION_RESPONSE, itti_s1ap_ue_context_mod_resp_t,
    s1ap_ue_context_mod_response)
MESSAGE_DEF(
    S1AP_UE_CONTEXT_MODIFICATION_FAILURE, itti_s1ap_ue_context_mod_resp_fail_t,
    s1ap_ue_context_mod_failure)
MESSAGE_DEF(S1AP_E_RAB_REL_CMD, itti_s1ap_e_rab_rel_cmd_t, s1ap_e_rab_rel_cmd)
MESSAGE_DEF(S1AP_E_RAB_REL_RSP, itti_s1ap_e_rab_rel_rsp_t, s1ap_e_rab_rel_rsp)
MESSAGE_DEF(
    S1AP_PATH_SWITCH_REQUEST, itti_s1ap_path_switch_request_t,
    s1ap_path_switch_request)
MESSAGE_DEF(
    S1AP_PATH_SWITCH_REQUEST_ACK, itti_s1ap_path_switch_request_ack_t,
    s1ap_path_switch_request_ack)
MESSAGE_DEF(
    S1AP_PATH_SWITCH_REQUEST_FAILURE, itti_s1ap_path_switch_request_failure_t,
    s1ap_path_switch_request_failure)
MESSAGE_DEF(
    S1AP_E_RAB_MODIFICATION_IND, itti_s1ap_e_rab_modification_ind_t,
    s1ap_e_rab_modification_ind)
MESSAGE_DEF(
    S1AP_E_RAB_MODIFICATION_CNF, itti_s1ap_e_rab_modification_cnf_t,
    s1ap_e_rab_modification_cnf)
MESSAGE_DEF(
    S1AP_REMOVE_STALE_UE_CONTEXT, itti_s1ap_remove_stale_ue_context_t,
    s1ap_remove_stale_ue_context)
MESSAGE_DEF(
    S1AP_HANDOVER_REQUIRED, itti_s1ap_handover_required_t,
    s1ap_handover_required)
MESSAGE_DEF(
    S1AP_HANDOVER_REQUEST_ACK, itti_s1ap_handover_request_ack_t,
    s1ap_handover_request_ack)
MESSAGE_DEF(
    S1AP_HANDOVER_NOTIFY, itti_s1ap_handover_notify_t, s1ap_handover_notify)
