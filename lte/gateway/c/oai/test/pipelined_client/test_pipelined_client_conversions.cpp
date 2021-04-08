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

#include "proto_converters.h"

namespace magma {
namespace lte {

class PipelineDClientTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}
};

TEST_F(PipelineDClientTest, TestMakeRequestIPv4) {
  struct in_addr enb_ipv4_addr;
  struct in_addr ue_ipv4_addr;
  uint32_t in_teid = 0;
  uint32_t out_teid = 1;
  struct ip_flow_dl flow_dl;
  uint32_t ue_state = 3;
  UESessionSet request = make_update_request_ipv4(enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, flow_dl, ue_state);
  
  // Run all assertions!
  EXPECT_EQ(request.in_teid(), in_teid);

}


int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace lte
}  // namespace magma
