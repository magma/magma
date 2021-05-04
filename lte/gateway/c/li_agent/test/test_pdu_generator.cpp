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
#include <chrono>
#include <thread>

#include <netinet/ip.h>
#include <net/ethernet.h>

#include "Consts.h"
#include "PDUGenerator.h"
#include "LIAgentdMocks.h"

#include <gtest/gtest.h>

using ::testing::Test;

namespace magma {

class PDUGeneratorTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    auto proxy_connector_p   = std::make_unique<MockProxyConnector>();
    proxy_connector          = proxy_connector_p.get();
    auto directoryd_client_p = std::make_unique<MockDirectorydClient>();
    directoryd_client        = directoryd_client_p.get();

    pkt_generator = std::make_unique<PDUGenerator>(
        std::move(proxy_connector_p), std::move(directoryd_client_p),
        PKT_DST_MAC, PKT_SRC_MAC);
  }

  MockProxyConnector* proxy_connector;
  MockDirectorydClient* directoryd_client;
  std::unique_ptr<PDUGenerator> pkt_generator;
};

TEST_F(PDUGeneratorTest, test_pdu_generator) {
  struct pcap_pkthdr* phdr =
      (struct pcap_pkthdr*) malloc(sizeof(struct pcap_pkthdr));
  phdr->len       = sizeof(struct ether_header) + sizeof(struct ip);
  phdr->ts.tv_sec = 92;
  u_char* pdata =
      reinterpret_cast<u_char*>(malloc(sizeof(struct ether_header) + sizeof(struct ip)));
  struct ether_header* ethernetHeader = (struct ether_header*) pdata;
  ethernetHeader->ether_type          = htons(ETHERTYPE_IP);
  struct ip* ipHeader     = (struct ip*) (pdata + sizeof(struct ether_header));
  ipHeader->ip_src.s_addr = 123;
  ipHeader->ip_dst.s_addr = 1222787743;
  pkt_generator->send_packet(phdr, pdata);

  // TODO(koolzz): For some reason these are not properly caught, fix...
  //  EXPECT_CALL(
  //      *directoryd_client, get_directoryd_xid_field(testing::_, testing::_));
  //  EXPECT_CALL(
  //      *directoryd_client, get_directoryd_xid_field(testing::_, testing::_));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
