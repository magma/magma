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

#include "lte/gateway/c/core/oai/tasks/grpc_service/SpgwServiceImpl.hpp"

namespace magma {
namespace lte {
void CheckFillIpv4(packet_filter_contents_t* pf_content, int exp_addr[],
                   int exp_mask[]) {
  for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; i++) {
    EXPECT_EQ(pf_content->ipv4remoteaddr[i].mask, exp_mask[i]);
    EXPECT_EQ(pf_content->ipv4remoteaddr[i].addr, exp_addr[i]);
    pf_content->ipv4remoteaddr[i].mask = (uint8_t)256;  // reset mask
    pf_content->ipv4remoteaddr[i].addr = (uint8_t)256;  // reset addr
  }
}

TEST(SPGWServiceImplTest, TestParseIpv4Network) {
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

TEST(SPGWServiceImplTest, TestFillIpv4) {
  SpgwServiceImpl test_service;
  packet_filter_contents_t pf_content;

  // Input "192.168.32.118/8", expected output 192.0.0.0 and 255.0.0.0.
  bool return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/8");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 0, 0, 0}.data(),
                std::array<int, 4>{255, 0, 0, 0}.data());

  // Input "192.168.32.118/16", expected output 192.168.0.0 and 255.255.0.0.
  return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/16");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 168, 0, 0}.data(),
                std::array<int, 4>{255, 255, 0, 0}.data());

  // Input "192.168.32.118/17", expected output 192.168.0.0 and 255.255.128.0.
  return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/17");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 168, 0, 0}.data(),
                std::array<int, 4>{255, 255, 128, 0}.data());

  // Input "192.168.32.118/24", expected output 192.168.32.0 and
  // 255.255.255.0.
  return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/24");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 168, 32, 0}.data(),
                std::array<int, 4>{255, 255, 255, 0}.data());

  // Input "192.168.32.118/26", expected output 192.168.32.64 and
  // 255.255.255.192.
  return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/26");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 168, 32, 64}.data(),
                std::array<int, 4>{255, 255, 255, 192}.data());

  // Input "192.168.32.118/32", expected output 192.168.32.118 and
  // 255.255.255.255.
  return_val = test_service.fillIpv4(&pf_content, "192.168.32.118/32");
  EXPECT_TRUE(return_val);
  CheckFillIpv4(&pf_content, std::array<int, 4>{192, 168, 32, 118}.data(),
                std::array<int, 4>{255, 255, 255, 255}.data());
}

}  // namespace lte
}  // namespace magma
