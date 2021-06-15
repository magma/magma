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

#pragma once

typedef enum ngap_Cause_PR {
  ngap_Cause_PR_NOTHING = 0,
  ngap_Cause_PR_radioNetwork,
  ngap_Cause_PR_transport,
  ngap_Cause_PR_nas,
  ngap_Cause_PR_protocol,
  ngap_Cause_PR_misc,
} ngap_Cause_PR;

typedef enum ngap_CauseRadioNetwork {
  ngap_CRN_unspecified = 0,
  ngap_CRN_txnrelocoverall_expiry,
  ngap_CRN_successful_handover,
  ngap_CRN_release_due_to_ngran_generated_reason,
  ngap_CRN_release_due_to_5gc_generated_reason,
  ngap_CRN_handover_canceled,
  ngap_CRN_partial_handover,
  ngap_CRN_ho_failure_in_target_5GC_ngran_node_or_target_system,
  ngap_CRN_ho_target_not_allowed,
  ngap_CRN_tngrelocoverall_expiry,
  ngap_CRN_tngrelocprep_expiry,
  ngap_CRN_cell_not_available,
  ngap_CRN_unknown_targetID,
  ngap_CRN_no_radio_resources_available_in_target_cell,
  ngap_CRN_unknown_local_UE_ngap_ID,
  ngap_CRN_inconsistent_remote_UE_ngap_ID,
  ngap_CRN_handover_desirable_for_radio_reason,
  ngap_CRN_time_critical_handover,
  ngap_CRN_resource_optimisation_handover,
  ngap_CRN_reduce_load_in_serving_cell,
  ngap_CRN_user_inactivity,
  ngap_CRN_radio_connection_with_ue_lost,
  ngap_CRN_radio_resources_not_available,
  ngap_CRN_invalid_qos_combination,
  ngap_CRN_failure_in_radio_interface_procedure,
  ngap_CRN_interaction_with_other_procedure,
  ngap_CRN_unknown_PDU_session_ID,
  ngap_CRN_unknown_qos_flow_ID,
  ngap_CRN_multiple_PDU_session_ID_instances,
  ngap_CRN_multiple_qos_flow_ID_instances,
  ngap_CRN_encryption_and_or_integrity_protection_algorithms_not_supported,
  ngap_CRN_ng_intra_system_handover_triggered,
  ngap_CRN_ng_inter_system_handover_triggered,
  ngap_CRN_xn_handover_triggered,
  ngap_CRN_not_supported_5QI_value,
  ngap_CRN_ue_context_transfer,
  ngap_CRN_ims_voice_eps_fallback_or_rat_fallback_triggered,
  ngap_CRN_up_integrity_protection_not_possible,
  ngap_CRN_up_confidentiality_protection_not_possible,
  ngap_CRN_slice_not_supported,
  ngap_CRN_ue_in_rrc_inactive_state_not_reachable,
  ngap_CRN_redirection,
  ngap_CRN_resources_not_available_for_the_slice,
  ngap_CRN_ue_max_integrity_protected_data_rate_reason,
  ngap_CRN_release_due_to_cn_detected_mobility
} e_ngap_CauseRadioNetwork;

// TODO for future use
typedef enum ngap_CauseTransport {
  ngap_CT_transport_resource_unavailable = 0,
  ngap_CT_unspecified
} e_ngap_CauseTransport;

typedef enum ngap_CauseNas {
  ngap_CauseNas_normal_release = 0,
  ngap_CauseNas_authentication_failure,
  ngap_CauseNas_deregister,
  ngap_CauseNas_unspecified
} e_ngap_CauseNAS;

// TODO for future use
typedef enum ngap_CauseProtocol {
  ngap_CP_transfer_syntax_error,
  ngap_CP_abstract_syntax_error_reject,
  ngap_CP_abstract_syntax_error_ignore_and_notify,
  ngap_CP_message_not_compatible_with_receiver_state,
  ngap_CP_semantic_error,
  ngap_CP_abstract_syntax_error_falsely_constructed_message
} e_ngap_CauseProtocol;

// TODO for future use
typedef enum ngap_CauseMisc {
  ngap_CM_control_processing_overload = 0,
  ngap_CM_not_enough_user_plane_processing_resources,
  ngap_CM_hardware_failure,
  ngap_CM_om_intervention,
  ngap_CM_unknown_PLMN,
  ngap_CM_unspecified
} e_ngap_CauseMisc;

// The purpose of the Cause IE is to indicate the reason for a particular event
// for the ngap protocol.
typedef struct ngap_Cause {
  ngap_Cause_PR present;
  union ngap_Cause_u {
    e_ngap_CauseRadioNetwork radioNetwork;  // Radio Network Layer Cause
    e_ngap_CauseTransport transport;        // Transport Layer Cause
    e_ngap_CauseNAS nas;                    // NAS Cause
    e_ngap_CauseProtocol protocol;          // Protocol Cause
    e_ngap_CauseMisc misc;                  // Miscellaneous Cause
  } ngapCause_u;
} ngap_Cause_t;
