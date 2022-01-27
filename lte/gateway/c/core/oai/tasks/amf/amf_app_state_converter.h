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
#include "lte/gateway/c/core/oai/include/state_converter.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
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
  static void state_to_proto(
      const amf_app_desc_t* amf_nas_state_p, MmeNasState* state_proto);

  // Deserialize amf_app_desc_t from oai::MmeNasState proto
  static void proto_to_state(
      const MmeNasState& state_proto, amf_app_desc_t* amf_nas_state_p);

  static void ue_to_proto(
      const ue_m5gmm_context_t* ue_ctxt, UeContext* ue_ctxt_proto);

  static void proto_to_ue(
      const UeContext& ue_ctxt_proto, ue_m5gmm_context_t* ue_ctxt);

  static void ue_m5gmm_context_to_proto(
      const ue_m5gmm_context_t* state_ue_m5gmm_context,
      UeContext* ue_context_proto);

  static void proto_to_ue_m5gmm_context(
      const UeContext& ue_context_proto,
      ue_m5gmm_context_t* state_ue_m5gmm_context);

  static std::string amf_app_convert_guti_m5_to_string(const guti_m5_t& guti);
  static void amf_app_convert_string_to_guti_m5(
      const std::string& guti_str, guti_m5_t* guti_m5_p);

  static void amf_security_context_to_proto(
      const amf_security_context_t* state_amf_security_context,
      EmmSecurityContext* emm_security_context_proto);

  static void proto_to_amf_security_context(
      const EmmSecurityContext& emm_security_context_proto,
      amf_security_context_t* state_amf_security_context);

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
  static void amf_context_to_proto(
      const amf_context_t* amf_ctx, EmmContext* emm_context_proto);

  /**
   * Deserialize oai::EmmContext proto message to struct amf_context_t, the
   * memory for EMM procedures in AMF context is allocated by AmfStateConvertor,
   * and freed by AmfStateManager
   */
  static void proto_to_amf_context(
      const EmmContext& emm_context_proto, amf_context_t* amf_ctx);

  template<typename NodeType>
  static void identity_tuple_to_proto(
      const NodeType* state_identity,
      magma::lte::oai::IdentityTuple* identity_proto, int size) {
    identity_proto->set_value(state_identity->u.value, size);
    identity_proto->set_num_digits(state_identity->length);
  }
  template<typename NodeType>
  static void proto_to_identity_tuple(
      const magma::lte::oai::IdentityTuple& identity_proto,
      NodeType* state_identity, int size) {
    memcpy(
        reinterpret_cast<void*>(&state_identity->u.value),
        identity_proto.value().data(), size);
    state_identity->length = identity_proto.num_digits();
  }

  static void tai_to_proto(
      const tai_t* state_tai, magma::lte::oai::Tai* tai_proto);

  static void proto_to_tai(
      const magma::lte::oai::Tai& tai_proto, tai_t* state_tai);
};
}  // namespace magma5g
