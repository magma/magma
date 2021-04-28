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

#include "mme_app_state_manager.h"
#include "nas_state_converter.h"
#include "emm_data.h"
#include "state_creators.h"
#include "timer.h"

namespace magma {
namespace lte {

class NASStateConverterTest : public ::testing::Test {
 virtual void SetUp() {
    mme_config_t config;
    CHECK_INIT_RETURN(timer_init());

    MmeNasStateManager::getInstance().initialize_state(&config);
  }

  virtual void TearDown() { MmeNasStateManager::getInstance().free_state(); }
  };

TEST_F(NASStateConverterTest, TestEMMStateConversion) {
  emm_context_t original_state, final_state;
  oai::EmmContext proto_state;
  NasStateConverter::emm_context_to_proto(&original_state, &proto_state);
  NasStateConverter::proto_to_emm_context(proto_state, &final_state);
}


int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace lte
}  // namespace magma
