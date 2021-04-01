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
/*! \file mme_app_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
// WARNING: Do not include this header directly. Use intertask_interface.h
// instead.

MESSAGE_DEF(
    MME_APP_CONNECTION_ESTABLISHMENT_CNF,
    itti_mme_app_connection_establishment_cnf_t,
    mme_app_connection_establishment_cnf)
MESSAGE_DEF(
    MME_APP_INITIAL_CONTEXT_SETUP_RSP, itti_mme_app_initial_context_setup_rsp_t,
    mme_app_initial_context_setup_rsp)
MESSAGE_DEF(
    MME_APP_INITIAL_CONTEXT_SETUP_FAILURE,
    itti_mme_app_initial_context_setup_failure_t,
    mme_app_initial_context_setup_failure)
MESSAGE_DEF(
    MME_APP_DELETE_SESSION_RSP, itti_mme_app_delete_session_rsp_t,
    mme_app_delete_session_rsp)
MESSAGE_DEF(
    MME_APP_S1AP_MME_UE_ID_NOTIFICATION,
    itti_mme_app_s1ap_mme_ue_id_notification_t,
    mme_app_s1ap_mme_ue_id_notification)
MESSAGE_DEF(
    MME_APP_UPLINK_DATA_IND, itti_mme_app_ul_data_ind_t, mme_app_ul_data_ind)
MESSAGE_DEF(
    MME_APP_DOWNLINK_DATA_CNF, itti_mme_app_dl_data_cnf_t, mme_app_dl_data_cnf)
MESSAGE_DEF(
    MME_APP_DOWNLINK_DATA_REJ, itti_mme_app_dl_data_rej_t, mme_app_dl_data_rej)
MESSAGE_DEF(
    MME_APP_HANDOVER_REQUEST, itti_mme_app_handover_request_t,
    mme_app_handover_request)
MESSAGE_DEF(
    MME_APP_HANDOVER_COMMAND, itti_mme_app_handover_command_t,
    mme_app_handover_command)
