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
/*! \file S11_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
MESSAGE_DEF(
    S11_CREATE_SESSION_REQUEST, itti_s11_create_session_request_t,
    s11_create_session_request)
MESSAGE_DEF(
    S11_CREATE_SESSION_RESPONSE, itti_s11_create_session_response_t,
    s11_create_session_response)
MESSAGE_DEF(
    S11_CREATE_BEARER_REQUEST, itti_s11_create_bearer_request_t,
    s11_create_bearer_request)
MESSAGE_DEF(
    S11_CREATE_BEARER_RESPONSE, itti_s11_create_bearer_response_t,
    s11_create_bearer_response)
MESSAGE_DEF(
    S11_MODIFY_BEARER_REQUEST, itti_s11_modify_bearer_request_t,
    s11_modify_bearer_request)
MESSAGE_DEF(
    S11_MODIFY_BEARER_RESPONSE, itti_s11_modify_bearer_response_t,
    s11_modify_bearer_response)
MESSAGE_DEF(
    S11_DELETE_SESSION_REQUEST, itti_s11_delete_session_request_t,
    s11_delete_session_request)
MESSAGE_DEF(
    S11_DELETE_SESSION_RESPONSE, itti_s11_delete_session_response_t,
    s11_delete_session_response)
MESSAGE_DEF(
    S11_DELETE_BEARER_COMMAND, itti_s11_delete_bearer_command_t,
    s11_delete_bearer_command)
MESSAGE_DEF(
    S11_RELEASE_ACCESS_BEARERS_REQUEST,
    itti_s11_release_access_bearers_request_t,
    s11_release_access_bearers_request)
MESSAGE_DEF(
    S11_RELEASE_ACCESS_BEARERS_RESPONSE,
    itti_s11_release_access_bearers_response_t,
    s11_release_access_bearers_response)
MESSAGE_DEF(S11_PAGING_REQUEST, itti_s11_paging_request_t, s11_paging_request)
MESSAGE_DEF(
    S11_SUSPEND_NOTIFICATION, itti_s11_suspend_notification_t,
    s11_suspend_notification)
MESSAGE_DEF(
    S11_SUSPEND_ACKNOWLEDGE, itti_s11_suspend_acknowledge_t,
    s11_suspend_acknowledge)
MESSAGE_DEF(
    S11_MODIFY_UE_AMBR_REQUEST, itti_s11_modify_ue_ambr_request_t,
    s11_modify_ue_ambr_request)
MESSAGE_DEF(
    S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST,
    itti_s11_nw_init_actv_bearer_request_t, s11_nw_init_actv_bearer_request)
MESSAGE_DEF(
    S11_NW_INITIATED_ACTIVATE_BEARER_RESP, itti_s11_nw_init_actv_bearer_rsp_t,
    s11_nw_init_actv_bearer_rsp)
MESSAGE_DEF(
    S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST,
    itti_s11_nw_init_deactv_bearer_request_t, s11_nw_init_deactv_bearer_request)
MESSAGE_DEF(
    S11_NW_INITIATED_DEACTIVATE_BEARER_RESP,
    itti_s11_nw_init_deactv_bearer_rsp_t, s11_nw_init_deactv_bearer_rsp)
MESSAGE_DEF(
    S11_DOWNLINK_DATA_NOTIFICATION, itti_s11_downlink_data_notification_t,
    s11_downlink_data_notification)
MESSAGE_DEF(
    S11_DOWNLINK_DATA_NOTIFICATION_ACKNOWLEDGE,
    itti_s11_downlink_data_notification_acknowledge_t,
    s11_downlink_data_notification_acknowledge)
