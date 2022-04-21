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

#include <sstream>
#include "lte/gateway/c/core/oai/include/state_converter.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"
#include "lte/gateway/c/core/oai/include/map.h"
#include "lte/protos/oai/mme_nas_state.pb.h"
#include "lte/protos/oai/nas_state.pb.h"

/******************************************************************************
 * This is a helper class to encapsulate all functions for converting in-memory
 * state of AMF and NAS task to/from proto for persisting UE state in data
 * store. The class does not have any member variables. The class does not
 * allocate any memory, but calls NAS state converter, which dynamically
 * allocates memory for EMM procedures. All the allocated memory is cleared by
 * the caller class AmfNasStateManager
 ******************************************************************************/

using magma::lte::oai::EmmContext;
using magma::lte::oai::EmmSecurityContext;
using magma::lte::oai::MmeNasState;
using magma::lte::oai::UeContext;
namespace magma5g {
class AmfNasStateConverter : public magma::lte::StateConverter {
 public:
  // Constructor
  AmfNasStateConverter();

  // Destructor
  ~AmfNasStateConverter();

  // Serialize amf_app_desc_t to oai::MmeNasState proto
  static void state_to_proto(const amf_app_desc_t* amf_nas_state_p,
                             MmeNasState* state_proto);

  // Deserialize amf_app_desc_t from oai::MmeNasState proto
  static void proto_to_state(const MmeNasState& state_proto,
                             amf_app_desc_t* amf_nas_state_p);

  static void ue_to_proto(const ue_m5gmm_context_t* ue_ctxt,
                          UeContext* ue_ctxt_proto);

  static void proto_to_ue(const UeContext& ue_ctxt_proto,
                          ue_m5gmm_context_t* ue_ctxt);

  static void ue_m5gmm_context_to_proto(
      const ue_m5gmm_context_t* state_ue_m5gmm_context,
      UeContext* ue_context_proto);

  static void proto_to_ue_m5gmm_context(
      const UeContext& ue_context_proto,
      ue_m5gmm_context_t* state_ue_m5gmm_context);

  static void smf_context_to_proto(
      const smf_context_t* state_smf_context,
      magma::lte::oai::SmfContext* smf_context_proto);

  static void proto_to_smf_context(
      const magma::lte::oai::SmfContext& smf_context_proto,
      smf_context_t* state_smf_context);

  static void smf_proc_data_to_proto(
      const smf_proc_data_t* state_smf_proc_data,
      magma::lte::oai::Smf_Proc_Data* smf_proc_data_proto);

  static void proto_to_smf_proc_data(
      const magma::lte::oai::Smf_Proc_Data& smf_proc_data_proto,
      smf_proc_data_t* state_smf_proc_data);

  static void s_nssai_to_proto(const s_nssai_t* state_s_nssai,
                               magma::lte::oai::SNssai* snassi_proto);

  static void proto_to_s_nssai(const magma::lte::oai::SNssai& snassi_proto,
                               s_nssai_t* state_s_nssai);

  static void protocol_configuration_options_to_proto(
      const protocol_configuration_options_t&
          state_protocol_configuration_options,
      magma::lte::oai::ProtocolConfigurationOptions*
          protocol_configuration_options_proto);

  static void proto_to_protocol_configuration_options(
      const magma::lte::oai::ProtocolConfigurationOptions&
          protocol_configuration_options_proto,
      protocol_configuration_options_t* state_protocol_configuration_options);

  static void pco_protocol_or_container_id_to_proto(
      const protocol_configuration_options_t&
          state_protocol_configuration_options,
      magma::lte::oai::ProtocolConfigurationOptions*
          protocol_configuration_options_proto);

  static void proto_to_pco_protocol_or_container_id(
      const magma::lte::oai::ProtocolConfigurationOptions&
          protocol_configuration_options_proto,
      protocol_configuration_options_t* state_protocol_configuration_options);

  static void session_ambr_to_proto(const session_ambr_t& state_session_ambr,
                                    magma::lte::oai::Ambr* ambr_proto);

  static void proto_to_session_ambr(const magma::lte::oai::Ambr& ambr_proto,
                                    session_ambr_t* state_ambr);

  static void qos_flow_setup_request_item_to_proto(
      const qos_flow_setup_request_item& state_qos_flow_request_item,
      magma::lte::oai::M5GQosFlowItem* qos_flow_item_proto);

  static void proto_to_qos_flow_setup_request_item(
      const magma::lte::oai::M5GQosFlowItem& qos_flow_item_proto,
      qos_flow_setup_request_item* state_qos_flow_request_item);

  static void qos_flow_level_parameters_to_proto(
      const qos_flow_level_qos_parameters& state_qos_flow_parameters,
      magma::lte::oai::QosFlowParameters* qos_flow_parameters_proto);

  static void proto_to_qos_flow_level_parameters(
      const magma::lte::oai::QosFlowParameters& qos_flow_parameters_proto,
      qos_flow_level_qos_parameters* state_qos_flow_parameters);

  static std::string amf_app_convert_guti_m5_to_string(const guti_m5_t& guti);
  static void amf_app_convert_string_to_guti_m5(const std::string& guti_str,
                                                guti_m5_t* guti_m5_p);

  static void amf_security_context_to_proto(
      const amf_security_context_t* state_amf_security_context,
      EmmSecurityContext* emm_security_context_proto);

  static void proto_to_amf_security_context(
      const EmmSecurityContext& emm_security_context_proto,
      amf_security_context_t* state_amf_security_context);

  static void guti_m5_to_proto(const guti_m5_t& state_guti_m5,
                               magma::lte::oai::Guti_m5* guti_m5_proto);
  static void proto_to_guti_m5(const magma::lte::oai::Guti_m5& guti_m5_proto,
                               guti_m5_t* state_guti_m5);

  static void smf_context_map_to_proto(
      const std::unordered_map<uint8_t, std::shared_ptr<smf_context_t>>&
          smf_ctxt_map,
      google::protobuf::Map<uint32_t, magma::lte::oai::SmfContext>* proto_map);

  static void proto_to_smf_context_map(
      const google::protobuf::Map<uint32_t, magma::lte::oai::SmfContext>&
          proto_map,
      std::unordered_map<uint8_t, std::shared_ptr<smf_context_t>>*
          smf_ctxt_map);

  static void plmn_to_chars(const plmn_t& state_plmn, char* plmn_array);

  static void chars_to_plmn(const char* plmn_array, plmn_t* state_plmn);

  /***********************************************************
   *                 Map <-> Proto
   * Functions to serialize/deserialize in-memory maps
   * for AMF task. Only AMF task inserts/removes elements in
   * the maps, so these calls are thread-safe.
   * We only need to lock the UE context structure as it can
   * also be accessed by the NAS task. If map is empty
   * the proto field is also empty
   ***********************************************************/

  static void map_guti_uint64_to_proto(
      const map_guti_m5_uint64_t guti_map,
      google::protobuf::Map<std::string, uint64_t>* proto_map);

  static void proto_to_guti_map(
      const google::protobuf::Map<std::string, uint64_t>& proto_map,
      map_guti_m5_uint64_t* guti_map);

  /**
   * Serialize struct amf_context_t to oai::EmmContext proto message, the memory
   * for AMF procedures in AMF context is dynamically allocated by AMF task and
   * AmfStateConvertor, and freed by AmfStateManager
   */
  static void amf_context_to_proto(const amf_context_t* amf_ctx,
                                   EmmContext* emm_context_proto);

  /**
   * Deserialize oai::EmmContext proto message to struct amf_context_t, the
   * memory for EMM procedures in AMF context is allocated by AmfStateConvertor,
   * and freed by AmfStateManager
   */
  static void proto_to_amf_context(const EmmContext& emm_context_proto,
                                   amf_context_t* amf_ctx);

  template <typename NodeType>
  static void identity_tuple_to_proto(
      const NodeType* state_identity,
      magma::lte::oai::IdentityTuple* identity_proto, int size) {
    identity_proto->set_value(state_identity->u.value, size);
    identity_proto->set_num_digits(state_identity->length);
  }
  template <typename NodeType>
  static void proto_to_identity_tuple(
      const magma::lte::oai::IdentityTuple& identity_proto,
      NodeType* state_identity, int size) {
    memcpy(reinterpret_cast<void*>(&state_identity->u.value),
           identity_proto.value().data(), size);
    state_identity->length = identity_proto.num_digits();
  }

  static void tai_to_proto(const tai_t* state_tai,
                           magma::lte::oai::Tai* tai_proto);

  static void proto_to_tai(const magma::lte::oai::Tai& tai_proto,
                           tai_t* state_tai);
};
}  // namespace magma5g
