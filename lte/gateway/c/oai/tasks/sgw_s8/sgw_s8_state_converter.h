/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "assertions.h"
#include "common_types.h"
#include "hashtable.h"

#ifdef __cplusplus
}
#endif

#include "state_converter.h"
#include "lte/protos/oai/sgw_state.pb.h"
#include "spgw_types.h"

namespace magma {
namespace lte {

// Class for SGW_S8 tasks state conversion helper functions.
class SgwStateConverter : StateConverter {
 public:
  /**
   * Main function to convert SGW state to proto definition
   * @param sgw_state const pointer to sgw_state struct
   * @param sgw_proto SgwState proto object to be written to
   * Memory is owned by the caller
   */
  static void state_to_proto(
      const sgw_state_t* sgw_state, oai::SgwState* sgw_proto);

  /**
   * Main function to convert SGW proto to state definition
   * @param sgw_proto SgwState proto object to be written to
   * @param sgw_state const pointer to sgw_state struct
   * Memory is owned by the caller
   */
  static void proto_to_state(
      const oai::SgwState& proto, sgw_state_t* sgw_state);

  static void ue_to_proto(
      const spgw_ue_context_t* ue_state, oai::SgwUeContext* ue_proto);

  static void proto_to_ue(
      const oai::SgwUeContext& ue_proto, spgw_ue_context_t* ue_context_p);

 private:
  SgwStateConverter();
  ~SgwStateConverter();
};
}  // namespace lte
}  // namespace magma
