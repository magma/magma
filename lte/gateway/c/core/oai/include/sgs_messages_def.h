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
MESSAGE_DEF(
    SGSAP_LOCATION_UPDATE_REQ, itti_sgsap_location_update_req_t,
    sgsap_location_update_req)
MESSAGE_DEF(
    SGSAP_LOCATION_UPDATE_ACC, itti_sgsap_location_update_acc_t,
    sgsap_location_update_acc)
MESSAGE_DEF(
    SGSAP_LOCATION_UPDATE_REJ, itti_sgsap_location_update_rej_t,
    sgsap_location_update_rej)
MESSAGE_DEF(SGSAP_STATUS, itti_sgsap_status_t, sgsap_status)
MESSAGE_DEF(
    SGSAP_EPS_DETACH_IND, itti_sgsap_eps_detach_ind_t, sgsap_eps_detach_ind)
MESSAGE_DEF(
    SGSAP_EPS_DETACH_ACK, itti_sgsap_eps_detach_ack_t, sgsap_eps_detach_ack)
MESSAGE_DEF(
    SGSAP_IMSI_DETACH_IND, itti_sgsap_imsi_detach_ind_t, sgsap_imsi_detach_ind)
MESSAGE_DEF(
    SGSAP_IMSI_DETACH_ACK, itti_sgsap_imsi_detach_ack_t, sgsap_imsi_detach_ack)
MESSAGE_DEF(
    SGSAP_TMSI_REALLOC_COMP, itti_sgsap_tmsi_reallocation_comp_t,
    sgsap_tmsi_realloc_comp)
MESSAGE_DEF(
    SGSAP_VLR_RESET_INDICATION, itti_sgsap_vlr_reset_indication_t,
    sgsap_vlr_reset_indication)
MESSAGE_DEF(
    SGSAP_VLR_RESET_ACK, itti_sgsap_vlr_reset_ack_t, sgsap_vlr_reset_ack)
MESSAGE_DEF(
    SGSAP_PAGING_REQUEST, itti_sgsap_paging_request_t, sgsap_paging_request)
MESSAGE_DEF(
    SGSAP_SERVICE_REQUEST, itti_sgsap_service_request_t, sgsap_service_request)
MESSAGE_DEF(
    SGSAP_PAGING_REJECT, itti_sgsap_paging_reject_t, sgsap_paging_reject)
MESSAGE_DEF(
    SGSAP_UE_UNREACHABLE, itti_sgsap_ue_unreachable_t, sgsap_ue_unreachable)
MESSAGE_DEF(
    SGSAP_UPLINK_UNITDATA, itti_sgsap_uplink_unitdata_t, sgsap_uplink_unitdata)
MESSAGE_DEF(
    SGSAP_DOWNLINK_UNITDATA, itti_sgsap_downlink_unitdata_t,
    sgsap_downlink_unitdata)
MESSAGE_DEF(SGSAP_RELEASE_REQ, itti_sgsap_release_req_t, sgsap_release_req)
MESSAGE_DEF(
    SGSAP_ALERT_REQUEST, itti_sgsap_alert_request_t, sgsap_alert_request)
MESSAGE_DEF(SGSAP_ALERT_ACK, itti_sgsap_alert_ack_t, sgsap_alert_ack)
MESSAGE_DEF(SGSAP_ALERT_REJECT, itti_sgsap_alert_reject_t, sgsap_alert_reject)
MESSAGE_DEF(
    SGSAP_UE_ACTIVITY_IND, itti_sgsap_ue_activity_ind_t, sgsap_ue_activity_ind)
MESSAGE_DEF(
    SGSAP_MM_INFORMATION_REQ, itti_sgsap_mm_information_req_t,
    sgsap_mm_information_req)
MESSAGE_DEF(
    SGSAP_SERVICE_ABORT_REQ, itti_sgsap_service_abort_req_t,
    sgsap_service_abort_req)
