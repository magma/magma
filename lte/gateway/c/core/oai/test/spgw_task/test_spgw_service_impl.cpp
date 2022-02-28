/**
 * Copyright 2022 The Magma Authors.
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

#include <gtest/gtest.h>

#include "lte/gateway/c/core/oai/tasks/grpc_service/SpgwServiceImpl.h"

namespace magma {
namespace lte {

class SPGWServiceImplTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}
};

TEST_F(SPGWServiceImplTest, TestSpgwServiceImpl) {
  SpgwServiceImpl test_service;

  // parseIpv4Network calls with valid ip and optional valid subnet masks
  ipv4_network_t result = test_service.parseIpv4Network("255.255.255.0/24");
  EXPECT_EQ(result.addr_hbo, 4294967040);
  EXPECT_EQ(result.mask_len, 24);
  EXPECT_EQ(result.success, true);

  result = test_service.parseIpv4Network("192.168.255.0");
  EXPECT_EQ(result.addr_hbo, 3232300800);
  EXPECT_EQ(result.mask_len, 32);
  EXPECT_EQ(result.success, true);

  result = test_service.parseIpv4Network("192.168.0.0/0");
  EXPECT_EQ(result.addr_hbo, 3232235520);
  EXPECT_EQ(result.mask_len, 0);
  EXPECT_EQ(result.success, true);

  result = test_service.parseIpv4Network("0.0.0.0/0");
  EXPECT_EQ(result.addr_hbo, 0);
  EXPECT_EQ(result.mask_len, 0);
  EXPECT_EQ(result.success, true);

  result = test_service.parseIpv4Network("192.168.0.0/32");
  EXPECT_EQ(result.addr_hbo, 3232235520);
  EXPECT_EQ(result.mask_len, 32);
  EXPECT_EQ(result.success, true);

  // parseIpv4Network calls with invalid combinations of ip and subnet masks
  std::vector<std::string> false_test_ips = {"0.0.0.0//0",
                                             "0.0.0.0\\0",
                                             "0/0.0.0.0",
                                             "5",
                                             "3.5",
                                             "abc",
                                             "192.168.0.0/33",
                                             "192.168.0.0/a",
                                             "192.168.0.0/-1",
                                             "192.468.0.0",
                                             "192.-168.0.0",
                                             "192.168.5"};

  for (std::string false_test_ip : false_test_ips) {
    result = test_service.parseIpv4Network(false_test_ip);
    EXPECT_EQ(result.success, false);
  }
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace lte
}  // namespace magma
