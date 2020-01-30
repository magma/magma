/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */
//WARNING: Do not include this header directly. Use intertask_interface.h instead.

/*! \file s1ap_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

/* Messages for S1AP logging */
MESSAGE_DEF(
  S1AP_UPLINK_NAS_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_uplink_nas_log)
MESSAGE_DEF(
  S1AP_UE_CAPABILITY_IND_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_ue_capability_ind_log)
MESSAGE_DEF(
  S1AP_INITIAL_CONTEXT_SETUP_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_initial_context_setup_log)
MESSAGE_DEF(
  S1AP_NAS_NON_DELIVERY_IND_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_nas_non_delivery_ind_log)
MESSAGE_DEF(
  S1AP_DOWNLINK_NAS_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_downlink_nas_log)
MESSAGE_DEF(
  S1AP_S1_SETUP_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_s1_setup_log)
MESSAGE_DEF(
  S1AP_INITIAL_UE_MESSAGE_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_initial_ue_message_log)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_REQ_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_ue_context_release_req_log)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_COMMAND_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_ue_context_release_command_log)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_ue_context_release_log)
MESSAGE_DEF(
  S1AP_ENB_RESET_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_enb_reset_log)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_MODIFICATION_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_ue_context_modification_log)
MESSAGE_DEF(
  S1AP_UE_CAPABILITIES_IND,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_cap_ind_t,
  s1ap_ue_cap_ind)
MESSAGE_DEF(
  S1AP_ENB_DEREGISTERED_IND,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_eNB_deregistered_ind_t,
  s1ap_eNB_deregistered_ind)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_REQ,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_release_req_t,
  s1ap_ue_context_release_req)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_COMMAND,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_release_command_t,
  s1ap_ue_context_release_command)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_RELEASE_COMPLETE,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_release_complete_t,
  s1ap_ue_context_release_complete)
MESSAGE_DEF(
  S1AP_NAS_DL_DATA_REQ,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_nas_dl_data_req_t,
  s1ap_nas_dl_data_req)
MESSAGE_DEF(
  S1AP_INITIAL_UE_MESSAGE,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_initial_ue_message_t,
  s1ap_initial_ue_message)
MESSAGE_DEF(
  S1AP_E_RAB_SETUP_REQ,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_e_rab_setup_req_t,
  s1ap_e_rab_setup_req)
MESSAGE_DEF(
  S1AP_E_RAB_SETUP_RSP,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_e_rab_setup_rsp_t,
  s1ap_e_rab_setup_rsp)
MESSAGE_DEF(
  S1AP_ENB_INITIATED_RESET_REQ,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_enb_initiated_reset_req_t,
  s1ap_enb_initiated_reset_req)
MESSAGE_DEF(
  S1AP_ENB_INITIATED_RESET_ACK,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_enb_initiated_reset_ack_t,
  s1ap_enb_initiated_reset_ack)
MESSAGE_DEF(
  S1AP_PAGING_REQUEST,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_paging_request_t,
  s1ap_paging_request)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_mod_req_t,
  s1ap_ue_context_mod_request)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_MODIFICATION_RESPONSE,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_mod_resp_t,
  s1ap_ue_context_mod_response)
MESSAGE_DEF(
  S1AP_UE_CONTEXT_MODIFICATION_FAILURE,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_ue_context_mod_resp_fail_t,
  s1ap_ue_context_mod_failure)
MESSAGE_DEF(
  S1AP_E_RAB_REL_CMD,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_e_rab_rel_cmd_t,
  s1ap_e_rab_rel_cmd)
MESSAGE_DEF(
  S1AP_E_RAB_REL_RSP,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_e_rab_rel_rsp_t,
  s1ap_e_rab_rel_rsp)
MESSAGE_DEF(
  S1AP_ENB_CONFIGURATION_TRANSFER_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_enb_configuration_transfer_log)
MESSAGE_DEF(
  S1AP_PATH_SWITCH_REQUEST_LOG,
  MESSAGE_PRIORITY_MED,
  IttiMsgText,
  s1ap_path_switch_request_log)
MESSAGE_DEF(
  S1AP_PATH_SWITCH_REQUEST,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_path_switch_request_t,
  s1ap_path_switch_request)
MESSAGE_DEF(
  S1AP_PATH_SWITCH_REQUEST_ACK,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_path_switch_request_ack_t,
  s1ap_path_switch_request_ack)
MESSAGE_DEF(
  S1AP_PATH_SWITCH_REQUEST_FAILURE,
  MESSAGE_PRIORITY_MED,
  itti_s1ap_path_switch_request_failure_t,
  s1ap_path_switch_request_failure)
