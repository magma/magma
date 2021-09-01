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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

extern "C" {
#include "emm_data.h"
#include "emm_proc.h"
#include "esm_proc.h"
#include "nas_message.h"
#include "nas_procedures.h"
#include "timer.h"
}

#include <sstream>
#include "lte/protos/oai/nas_state.pb.h"
#include "lte/protos/oai/spgw_state.pb.h"
#include "state_converter.h"

using magma::lte::oai::TrafficFlowTemplate;
/******************************************************************************
 * This is a helper class to encapsulate all functions for converting in-memory
 * state of MME and NAS task to/from proto for persisting UE state in data
 * store. The class does not have any member variables. This class does
 * dynamically allocate memory for EMM procedures. All the allocated memory is
 * cleared by the MmeNasStateManager
 ******************************************************************************/

namespace magma {
namespace lte {

class NasStateConverter : StateConverter {
 public:
  // Constructor
  NasStateConverter();

  // Destructor
  ~NasStateConverter();

  /**
   * Serialize struct emm_context_t to oai::EmmContext proto message, the memory
   * for EMM procedures in EMM context is dynamically allocated by NAS task and
   * NasStateConvertor, and freed by MmeNasStateManager
   */
  static void emm_context_to_proto(
      const emm_context_t* state_emm_context,
      oai::EmmContext* emm_context_proto);

  /**
   * Deserialize oai::EmmContext proto message to struct emm_context_t, the
   * memory for EMM procedures in EMM context is allocated by NasStateConvertor,
   * and freed by MmeNasStateManager
   */
  static void proto_to_emm_context(
      const oai::EmmContext& emm_context_proto,
      emm_context_t* state_emm_context);

  template<typename NodeType>
  static void identity_tuple_to_proto(
      const NodeType* state_identity, oai::IdentityTuple* identity_proto,
      int size) {
    identity_proto->set_value(state_identity->u.value, size);
    identity_proto->set_num_digits(state_identity->length);
  }

  template<typename NodeType>
  static void proto_to_identity_tuple(
      const oai::IdentityTuple& identity_proto, NodeType* state_identity,
      int size) {
    memcpy(
        (void*) &state_identity->u.value, identity_proto.value().data(), size);
    state_identity->length = identity_proto.num_digits();
  }

  // TODO: To be moved to base state converter class
  static void proto_to_guti(const oai::Guti& guti_proto, guti_t* state_guti);

  static void proto_to_ecgi(const oai::Ecgi& ecgi_proto, ecgi_t* state_ecgi);

  static void tai_list_to_proto(
      const tai_list_t* state_tai_list, oai::TaiList* tai_list_proto);

  static void proto_to_tai_list(
      const oai::TaiList& tai_list_proto, tai_list_t* state_tai_list);

  static void tai_to_proto(const tai_t* state_tai, oai::Tai* tai_proto);

  static void proto_to_tai(const oai::Tai& tai_proto, tai_t* state_tai);

  static void esm_ebr_context_to_proto(
      const esm_ebr_context_t& state_esm_ebr_context,
      oai::EsmEbrContext* esm_ebr_context_proto);

  static void proto_to_esm_ebr_context(
      const oai::EsmEbrContext& esm_ebr_context_proto,
      esm_ebr_context_t* state_esm_ebr_context);

  static void protocol_configuration_options_to_proto(
      const protocol_configuration_options_t&
          state_protocol_configuration_options,
      oai::ProtocolConfigurationOptions* protocol_configuration_options_proto);

  static void proto_to_protocol_configuration_options(
      const oai::ProtocolConfigurationOptions&
          protocol_configuration_options_proto,
      protocol_configuration_options_t* state_protocol_configuration_options);

  static void esm_ebr_timer_data_to_proto(
      const esm_ebr_timer_data_t& state_esm_ebr_timer_data,
      oai::EsmEbrTimerData* proto_esm_ebr_timer_data);

  static void proto_to_esm_ebr_timer_data(
      const oai::EsmEbrTimerData& proto_esm_ebr_timer_data,
      esm_ebr_timer_data_t** state_esm_ebr_timer_data);

 private:
  static void partial_tai_list_to_proto(
      const partial_tai_list_t* state_partial_tai_list,
      oai::PartialTaiList* partial_tai_list_proto);

  static void proto_to_partial_tai_list(
      const oai::PartialTaiList& partial_tai_list_proto,
      partial_tai_list_t* state_partial_tai_list);

  static void ue_network_capability_to_proto(
      const ue_network_capability_t* state_ue_network_capability,
      oai::UeNetworkCapability* ue_network_capability_proto);

  static void proto_to_ue_network_capability(
      const oai::UeNetworkCapability& ue_network_capability_proto,
      ue_network_capability_t* state_ue_network_capability);

  static void classmark2_to_proto(
      const MobileStationClassmark2* state_MobileStationClassmark,
      oai::MobileStaClassmark2* mobile_station_classmark2_proto);

  static void proto_to_classmark2(
      const oai::MobileStaClassmark2& mobile_sta_classmark2_proto,
      MobileStationClassmark2* state_MobileStationClassmar);

  static void voice_preference_to_proto(
      const voice_domain_preference_and_ue_usage_setting_t*
          state_voice_domain_preference_and_ue_usage_setting,
      oai::VoicePreference* voice_preference_proto);

  static void proto_to_voice_preference(
      const oai::VoicePreference& voice_preference_proto,
      voice_domain_preference_and_ue_usage_setting_t*
          state_voice_domain_preference_and_ue_usage_setting);

  static void ue_additional_security_capability_to_proto(
      const ue_additional_security_capability_t*
          state_ue_additional_security_capability,
      oai::UeAdditionalSecurityCapability*
          ue_additional_security_capability_proto);

  static void proto_to_ue_additional_security_capability(
      const oai::UeAdditionalSecurityCapability&
          ue_additional_security_capability_proto,
      ue_additional_security_capability_t*
          state_ue_additional_security_capability);

  static void ambr_to_proto(const ambr_t& state_ambr, oai::Ambr* ambr_proto);

  static void proto_to_ambr(const oai::Ambr& ambr_proto, ambr_t* state_ambr);

  static void bearer_qos_to_proto(
      const bearer_qos_t& state_bearer_qos, oai::BearerQos* bearer_qos_proto);

  static void proto_to_bearer_qos(
      const oai::BearerQos& bearer_qos_proto, bearer_qos_t* state_bearer_qos);

  static void pco_protocol_or_container_id_to_proto(
      const protocol_configuration_options_t&
          state_protocol_configuration_options,
      oai::ProtocolConfigurationOptions* protocol_configuration_options_proto);

  static void proto_to_pco_protocol_or_container_id(
      const oai::ProtocolConfigurationOptions&
          protocol_configuration_options_proto,
      protocol_configuration_options_t* state_protocol_configuration_options);

  static void esm_proc_data_to_proto(
      const esm_proc_data_t* state_esm_proc_data,
      oai::EsmProcData* esm_proc_data_proto);

  static void proto_to_esm_proc_data(
      const oai::EsmProcData& esm_proc_data_proto,
      esm_proc_data_t* state_esm_proc_data);

  static void esm_context_to_proto(
      const esm_context_t* state_esm_context,
      oai::EsmContext* esm_context_proto);

  static void proto_to_esm_context(
      const oai::EsmContext& esm_context_proto,
      esm_context_t* state_esm_context);

  static void nas_message_decode_status_to_proto(
      const nas_message_decode_status_t* state_nas_message_decode_status,
      oai::NasMsgDecodeStatus* nas_msg_decode_status_proto);

  static void proto_to_nas_message_decode_status(
      const oai::NasMsgDecodeStatus& nas_msg_decode_status_proto,
      nas_message_decode_status_t* state_nas_message_decode_status);

  static void emm_attach_request_ies_to_proto(
      const emm_attach_request_ies_t* state_emm_attach_request_ies,
      oai::AttachRequestIes* attach_request_ies_proto);

  static void proto_to_emm_attach_request_ies(
      const oai::AttachRequestIes& attach_request_ies_proto,
      emm_attach_request_ies_t* state_emm_attach_request_ies);

  static void nas_attach_proc_to_proto(
      const nas_emm_attach_proc_t* state_nas_attach_proc,
      oai::AttachProc* attach_proc_proto);

  static void proto_to_nas_emm_attach_proc(
      const oai::AttachProc& attach_proc_proto,
      nas_emm_attach_proc_t* state_nas_emm_attach_proc);

  static void emm_detach_request_ies_to_proto(
      const emm_detach_request_ies_t* state_emm_detach_request_ies,
      oai::DetachRequestIes* detach_request_ies_proto);

  static void proto_to_emm_detach_request_ies(
      const oai::DetachRequestIes& detach_request_ies_proto,
      emm_detach_request_ies_t* state_emm_detach_request_ies);

  static void emm_tau_request_ies_to_proto(
      const emm_tau_request_ies_t* state_emm_tau_request_ies,
      oai::TauRequestIes* tau_request_ies_proto);

  static void proto_to_emm_tau_request_ies(
      const oai::TauRequestIes& tau_request_ies_proto,
      emm_tau_request_ies_t* state_emm_tau_request_ies);

  static void nas_emm_tau_proc_to_proto(
      const nas_emm_tau_proc_t* state_nas_emm_tau_proc,
      oai::NasTauProc* nas_tau_proc_proto);

  static void proto_to_nas_emm_tau_proc(
      const oai::NasTauProc& nas_tau_proc_proto,
      nas_emm_tau_proc_t* state_nas_emm_tau_proc);

  static void nas_emm_auth_proc_to_proto(
      const nas_emm_auth_proc_t* state_nas_emm_auth_proc,
      oai::AuthProc* auth_proc_proto);

  static void proto_to_nas_emm_auth_proc(
      const oai::AuthProc& auth_proc_proto,
      nas_emm_auth_proc_t* state_nas_emm_auth_proc);

  static void nas_emm_smc_proc_to_proto(
      const nas_emm_smc_proc_t* state_nas_emm_smc_proc,
      oai::SmcProc* smc_proc_proto);

  static void proto_to_nas_emm_smc_proc(
      const oai::SmcProc& smc_proc_proto,
      nas_emm_smc_proc_t* state_nas_emm_smc_proc);

  static void nas_proc_mess_sign_to_proto(
      const nas_proc_mess_sign_t* state_nas_proc_mess_sign,
      oai::NasProcMessSign* nas_proc_mess_sign_proto);

  static void proto_to_nas_proc_mess_sign(
      const oai::NasProcMessSign& nas_proc_mess_sign_proto,
      nas_proc_mess_sign_t* state_nas_proc_mess_sign);

  static void nas_base_proc_to_proto(
      const nas_base_proc_t* base_proc_p, oai::NasBaseProc* base_proc_proto);

  static void proto_to_nas_base_proc(
      const oai::NasBaseProc& nas_base_proc_proto,
      nas_base_proc_t* state_nas_base_proc);

  static void emm_proc_to_proto(
      const nas_emm_proc_t* emm_proc_p, oai::NasEmmProc* emm_proc_proto);

  static void proto_to_nas_emm_proc(
      const oai::NasEmmProc& nas_emm_proc_proto,
      nas_emm_proc_t* state_nas_emm_proc);

  static void emm_specific_proc_to_proto(
      const nas_emm_specific_proc_t* state_emm_specific_proc,
      oai::NasEmmProcWithType* emm_proc_with_type);

  static void proto_to_emm_specific_proc(
      const oai::NasEmmProcWithType& proto_emm_proc_with_type,
      emm_procedures_t* state_emm_procedures);

  static void emm_common_proc_to_proto(
      const emm_procedures_t* state_emm_procedures,
      oai::EmmProcedures* emm_procedures_proto);

  static void insert_proc_into_emm_common_procs(
      emm_procedures_t* state_emm_procedures,
      nas_emm_common_proc_t* emm_com_proc);

  static void proto_to_emm_common_proc(
      const oai::EmmProcedures& emm_procedures_proto,
      emm_context_t* state_emm_context);

  static void eutran_vectors_to_proto(
      eutran_vector_t** state_eutran_vector_array, uint8_t num_vectors,
      oai::AuthInfoProc* auth_info_proc_proto);

  static void proto_to_eutran_vectors(
      const oai::AuthInfoProc& auth_info_proc_proto,
      nas_auth_info_proc_t* state_nas_auth_info_proc);

  static void nas_auth_info_proc_to_proto(
      nas_auth_info_proc_t* state_nas_auth_info_proc,
      oai::AuthInfoProc* auth_info_proc_proto);

  static void proto_to_nas_auth_info_proc(
      const oai::AuthInfoProc& auth_info_proc_proto,
      nas_auth_info_proc_t* state_nas_auth_info_proc);

  static void nas_cn_procs_to_proto(
      const emm_procedures_t* state_emm_procedures,
      oai::EmmProcedures* emm_procedures_proto);

  static void insert_proc_into_cn_procs(
      emm_procedures_t* state_emm_procedures, nas_cn_proc_t* cn_proc);

  static void proto_to_nas_cn_proc(
      const oai::EmmProcedures& emm_procedures_proto,
      emm_procedures_t* state_emm_procedures);

  static void mess_sign_array_to_proto(
      const emm_procedures_t* state_emm_procedures,
      oai::EmmProcedures* emm_procedures_proto);

  static void proto_to_mess_sign_array(
      const oai::EmmProcedures& emm_procedures_proto,
      emm_procedures_t* state_emm_procedures);

  static void emm_procedures_to_proto(
      const emm_procedures_t* state_emm_procedures,
      oai::EmmProcedures* emm_procedures_proto);

  static void proto_to_emm_procedures(
      const oai::EmmProcedures& emm_procedures_proto,
      emm_context_t* state_emm_context);

  static void auth_vectors_to_proto(
      const auth_vector_t* state_auth_vector_array, int num_vectors,
      oai::EmmContext* emm_context_proto);

  static int proto_to_auth_vectors(
      const oai::EmmContext& emm_context_proto,
      auth_vector_t* state_auth_vector);

  static void emm_security_context_to_proto(
      const emm_security_context_t* state_emm_security_context,
      oai::EmmSecurityContext* emm_security_context_proto);

  static void proto_to_emm_security_context(
      const oai::EmmSecurityContext& emm_security_context_proto,
      emm_security_context_t* state_emm_security_context);

  static void nw_detach_data_to_proto(
      nw_detach_data_t* detach_timer_arg,
      oai::NwDetachData* detach_timer_arg_proto);

  static void proto_to_nw_detach_data(
      const oai::NwDetachData& detach_timer_arg_proto,
      nw_detach_data_t** detach_timer_arg);

  static void new_attach_info_to_proto(
      const new_attach_info_t* state_new_attach_info,
      oai::NewAttachInfo* proto_new_attach_info);

  static void proto_to_new_attach_info(
      const oai::NewAttachInfo& proto_new_attach_info,
      new_attach_info_t* state_new_attach_info);
};

}  // namespace lte
}  // namespace magma
