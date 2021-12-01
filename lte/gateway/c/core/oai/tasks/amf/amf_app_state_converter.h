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

/******************************************************************************
 * This is a helper class to encapsulate all functions for converting in-memory
 * state of AMF and NAS task to/from proto for persisting UE state in data
 * store. The class does not have any member variables. The class does not
 * allocate any memory, but calls NAS state converter, which dynamically
 * allocates memory for EMM procedures. All the allocated memory is cleared by
 * the caller class AmfNasStateManager
 ******************************************************************************/

namespace magma5g {
class AmfNasStateConverter : public magma::lte::StateConverter {
 public:
  // Constructor
  AmfNasStateConverter();

  // Destructor
  ~AmfNasStateConverter();

  // Serialize amf_app_desc_t to oai::MmeNasState proto
  static void state_to_proto(
      const amf_app_desc_t* amf_nas_state_p,
      magma::lte::oai::MmeNasState* state_proto);

  // Deserialize amf_app_desc_t from oai::MmeNasState proto
  static void proto_to_state(
      const magma::lte::oai::MmeNasState& state_proto,
      amf_app_desc_t* amf_nas_state_p);

  static void ue_to_proto(
      const ue_m5gmm_context_t* ue_ctxt,
      magma::lte::oai::UeContext* ue_ctxt_proto);

  static void proto_to_ue(
      const magma::lte::oai::UeContext& ue_ctxt_proto,
      ue_m5gmm_context_t* ue_ctxt);

  // Note: declare these helper functions as private after testing
  static std::string amf_app_convert_guti_m5_to_string(guti_m5_t guti);
  static void amf_app_convert_string_to_guti_m5(
      guti_m5_t* guti_m5_p, const std::string& guti_str);

 private:
  /***********************************************************
   *                 Map <-> Proto
   * Functions to serialize/deserialize in-memory maps
   * for AMF task. Only AMF task inserts/removes elements in
   * the maps, so these calls are thread-safe.
   * We only need to lock the UE context structure as it can
   * also be accessed by the NAS task. If map is empty
   * the proto field is also empty
   ***********************************************************/

  static void map_uint64_uint64_to_proto(
      map_uint64_uint64_t map,
      google::protobuf::Map<uint64_t, uint64_t>* proto_map);

  static void proto_to_map_uint64_uint64(
      const google::protobuf::Map<uint64_t, uint64_t>& proto_map,
      map_uint64_uint64_t* map);

  static void map_guti_uint64_to_proto(
      map_guti_m5_uint64_t guti_map,
      google::protobuf::Map<std::string, uint64_t>* proto_map);

  static void proto_to_guti_map(
      const google::protobuf::Map<std::string, uint64_t>& proto_map,
      map_guti_m5_uint64_t* guti_map);
};
}  // namespace magma5g
