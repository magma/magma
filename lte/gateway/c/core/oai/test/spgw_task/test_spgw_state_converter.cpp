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
#include <string.h>
#include <gtest/gtest.h>
#include <cstdio>
#include <cstdlib>

#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_converter.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/state_creators.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.hpp"
#include "lte/protos/oai/mme_nas_state.pb.h"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/lib/message_utils/ie_to_bytes.h"
}

namespace magma {
namespace lte {

TEST(SPGWStateConverterTest, TestSPGWStateConversion) {
  std::vector<spgw_state_t> original_states{
      make_spgw_state(4, 8, 12),
      make_spgw_state(500, 900, 1300),
      make_spgw_state(6000, 1000, 1400),
      make_spgw_state(0, 0, 32),
  };

  for (spgw_state_t& initial_state : original_states) {
    oai::SpgwState proto_state;
    spgw_state_t final_state;

    SpgwStateConverter::state_to_proto(&initial_state, &proto_state);
    SpgwStateConverter::proto_to_state(proto_state, &final_state);

    EXPECT_EQ(initial_state.gtpv1u_teid, final_state.gtpv1u_teid);

    gtpv1u_data_t initial_gtp_data = initial_state.gtpv1u_data;
    gtpv1u_data_t final_gtp_data = final_state.gtpv1u_data;
    EXPECT_EQ(initial_gtp_data.fd0, final_gtp_data.fd0);
    EXPECT_EQ(initial_gtp_data.fd1u, final_gtp_data.fd1u);
  }
}

TEST(SPGWStateConverterTest, TestUEContextConversion) {
  // Init SPGW hashtable
  spgw_config_t* spgw_config_p =
      (spgw_config_t*)calloc(1, sizeof(spgw_config_t));
  spgw_config_init(spgw_config_p);
  spgw_state_init(false, spgw_config_p);

  // Create some UE contexts, populate them with TEIDs
  std::vector<s_plus_p_gw_eps_bearer_context_information_t*> bearer_ctxs{
      make_bearer_context(10, 100),
      make_bearer_context(20, 200),
      make_bearer_context(30, 300),
  };

  // Repeat for each UE context
  for (s_plus_p_gw_eps_bearer_context_information_t* want : bearer_ctxs) {
    auto imsi = want->sgw_eps_bearer_context_information.imsi64;

    // Convert to proto and back
    auto initial_ctx = spgw_create_or_get_ue_context(imsi);
    auto want_teid = LIST_FIRST(&initial_ctx->sgw_s11_teid_list)->sgw_s11_teid;

    oai::SpgwUeContext proto_ctx;
    spgw_ue_context_t* final_ctx = new spgw_ue_context_t();
    SpgwStateConverter::ue_to_proto(initial_ctx, &proto_ctx);
    SpgwStateConverter::proto_to_ue(proto_ctx, final_ctx);

    // Ensure underlying bearer contexts match
    auto got = sgw_cm_get_spgw_context(want_teid);
    EXPECT_NE(got, nullptr);

    // Ensure PGW context matches
    auto want_pgw = want->pgw_eps_bearer_context_information;
    auto got_pgw = got->pgw_eps_bearer_context_information;

    EXPECT_TRUE(!memcmp(&want_pgw.imsi, &got_pgw.imsi, sizeof(want_pgw.imsi)));

    EXPECT_EQ(want_pgw.imsi_unauthenticated_indicator,
              got_pgw.imsi_unauthenticated_indicator);
    EXPECT_EQ(std::string(want_pgw.msisdn), std::string(got_pgw.msisdn));

    // Ensure SGW context matches
    auto want_sgw = want->sgw_eps_bearer_context_information;
    auto got_sgw = got->sgw_eps_bearer_context_information;

    EXPECT_TRUE(!memcmp(&want_sgw.imsi, &got_sgw.imsi, sizeof(want_sgw.imsi)));

    EXPECT_EQ(want_sgw.imsi64, got_sgw.imsi64);
    EXPECT_EQ(want_sgw.imsi_unauthenticated_indicator,
              got_sgw.imsi_unauthenticated_indicator);
    EXPECT_EQ(std::string(want_sgw.msisdn), std::string(got_sgw.msisdn));
    EXPECT_EQ(want_sgw.mme_teid_S11, got_sgw.mme_teid_S11);
    EXPECT_EQ(want_sgw.s_gw_teid_S11_S4, got_sgw.s_gw_teid_S11_S4);

    EXPECT_TRUE(!memcmp(&want_sgw.mme_ip_address_S11,
                        &got_sgw.mme_ip_address_S11,
                        sizeof(want_sgw.mme_ip_address_S11)));

    EXPECT_TRUE(!memcmp(&want_sgw.s_gw_ip_address_S11_S4,
                        &got_sgw.s_gw_ip_address_S11_S4,
                        sizeof(want_sgw.s_gw_ip_address_S11_S4)));

    EXPECT_TRUE(!memcmp(&want_sgw.last_known_cell_Id,
                        &got_sgw.last_known_cell_Id,
                        sizeof(want_sgw.last_known_cell_Id)));
    sgw_free_ue_context(reinterpret_cast<void**>(&final_ctx));
  }
}
}  // namespace lte
}  // namespace magma
