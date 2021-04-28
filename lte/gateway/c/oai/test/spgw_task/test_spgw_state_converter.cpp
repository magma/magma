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

#include "spgw_state_converter.h"
#include "spgw_state_manager.h"
#include "sgw_defs.h"
#include "state_creators.h"

namespace magma {
namespace lte {

class SPGWStateConverterTest : public ::testing::Test {
  virtual void SetUp() {
    spgw_config_t config;
    SpgwStateManager::getInstance().init(false, &config);
  }

  virtual void TearDown() { SpgwStateManager::getInstance().free_state(); }
};

TEST_F(SPGWStateConverterTest, TestSPGWStateConversion) {
  std::vector<spgw_state_t> original_states{
      make_spgw_state(teid_t(0), 4, 8, 12),
      make_spgw_state(teid_t(100), 500, 900, 1300),
      make_spgw_state(teid_t(2000), 6000, 1000, 1400),
      make_spgw_state(teid_t(9), 0, 0, 32),
  };

  for (spgw_state_t& initial_state : original_states) {
    oai::SpgwState proto_state;
    spgw_state_t final_state;

    SpgwStateConverter::state_to_proto(&initial_state, &proto_state);
    SpgwStateConverter::proto_to_state(proto_state, &final_state);

    EXPECT_EQ(initial_state.tunnel_id, final_state.tunnel_id);
    EXPECT_EQ(initial_state.gtpv1u_teid, final_state.gtpv1u_teid);

    gtpv1u_data_t initial_gtp_data = initial_state.gtpv1u_data;
    gtpv1u_data_t final_gtp_data   = final_state.gtpv1u_data;
    EXPECT_EQ(initial_gtp_data.fd0, final_gtp_data.fd0);
    EXPECT_EQ(initial_gtp_data.fd1u, final_gtp_data.fd1u);
  }
}

TEST_F(SPGWStateConverterTest, TestEmptySPGWUeConversion) {
  spgw_ue_context_t original_state, final_state;
  oai::SpgwUeContext proto_state;
  LIST_INIT(&original_state.sgw_s11_teid_list);

  SpgwStateConverter::ue_to_proto(&original_state, &proto_state);
  SpgwStateConverter::proto_to_ue(proto_state, &final_state);

  EXPECT_TRUE(LIST_EMPTY(&(original_state.sgw_s11_teid_list)));
  EXPECT_TRUE(LIST_EMPTY(&(final_state.sgw_s11_teid_list)));
}

// TODO add a state conversion test for UE context

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace lte
}  // namespace magma
