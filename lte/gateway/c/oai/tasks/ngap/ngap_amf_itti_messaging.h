/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/****************************************************************************
  Source      ngap_amf_itti_messaging.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <czmq.h>

#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_38.401.h"
#include "3gpp_38.413.h"
#include "Ngap_Cause.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface.h"
#include "ngap_state.h"

task_zmq_ctx_t ngap_task_zmq_ctx;

int ngap_amf_itti_send_sctp_request(
    STOLEN_REF bstring* payload, const uint32_t sctp_assoc_id_t,
    const sctp_stream_id_t stream, const amf_ue_ngap_id_t ue_id);

int ngap_amf_itti_nas_uplink_ind(
    const amf_ue_ngap_id_t ue_id, STOLEN_REF bstring* payload,
    const tai_t const* tai, const ecgi_t const* cgi);

int ngap_amf_itti_nas_downlink_cnf(
    const amf_ue_ngap_id_t ue_id, const bool is_success);

void ngap_amf_itti_ngap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t enb_id,
    const gnb_ue_ngap_id_t gnb_ue_ngap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t const* tai,
    const ecgi_t const* ecgi, const long rrc_cause,
    const s_tmsi_m5_t const* opt_s_tmsi, const csg_id_t const* opt_csg_id,
    const guamfi_t const* opt_guamfi,
    const void const* opt_cell_access_mode,          /* unused*/
    const void const* opt_cell_gw_transport_address, /* unused*/
    const void const* opt_relay_node_indicator       /* unused*/
);

void ngap_amf_itti_nas_non_delivery_ind(
    const amf_ue_ngap_id_t ue_id, uint8_t* const nas_msg,
    const size_t nas_msg_length, const Ngap_Cause_t* const cause,
    const imsi64_t imsi64);

