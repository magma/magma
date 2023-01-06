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

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/common_types.h"

#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/state_converter.hpp"
#include "lte/protos/oai/std_3gpp_types.pb.h"
#include "lte/protos/oai/spgw_state.pb.h"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/include/pgw_types.h"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

namespace magma {
namespace lte {

/**
 * Class for SGW / PGW tasks state conversion helper functions.
 */
class SpgwStateConverter : StateConverter {
 public:
  /**
   * Main function to convert SPGW state to proto definition
   * @param spgw_state const pointer to spgw_state struct
   * @param spgw_proto SpgwState proto object to be written to
   * Memory is owned by the caller
   */
  static void state_to_proto(const spgw_state_t* spgw_state,
                             oai::SpgwState* spgw_proto);

  /**
   * Main function to convert SPGW proto to state definition
   * @param spgw_proto SpgwState proto object to be written to
   * @param spgw_state const pointer to spgw_state struct
   * Memory is owned by the caller
   */
  static void proto_to_state(const oai::SpgwState& proto,
                             spgw_state_t* spgw_state);

  static void ue_to_proto(const oai::SpgwUeContext* ue_state,
                          oai::SpgwUeContext* ue_proto);

  static void proto_to_ue(const oai::SpgwUeContext& ue_proto,
                          oai::SpgwUeContext* ue_context_p);

  /**
   * Converts traffic flow template struct to proto, memory is owned by the
   * caller
   * @param tft_state
   * @param tft_proto
   */
 private:
  SpgwStateConverter();
  ~SpgwStateConverter();

  /**
   * Converts spgw bearer context struct to proto, memory is owned by the caller
   * @param spgw_bearer_state
   * @param spgw_bearer_proto
   */
  static void gtpv1u_data_to_proto(const gtpv1u_data_t* gtp_data,
                                   oai::GTPV1uData* gtp_proto);

  /**
   * Converts proto to gtpv1u_data struct
   * @param gtp_proto
   * @param gtp_data
   */
  static void proto_to_gtpv1u_data(const oai::GTPV1uData& gtp_proto,
                                   gtpv1u_data_t* gtp_data);

  /**
   * Converts port range struct to proto, memory is owned by the caller
   * @param port_range
   * @param port_range_proto
   */
};
}  // namespace lte
}  // namespace magma
void proto_to_traffic_flow_template(
    const magma::lte::oai::TrafficFlowTemplate& tft_proto,
    traffic_flow_template_t* tft_state);
