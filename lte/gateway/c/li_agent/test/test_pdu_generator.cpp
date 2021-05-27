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

using testing::DoAll;
using ::testing::Test;
using ::testing::Return;
using ::testing::SetArgReferee;

namespace magma {
namespace lte {

class PDUGeneratorTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    auto proxy_connector_p = std::make_unique<MockProxyConnector>();
    proxy_connector        = proxy_connector_p.get();

    auto mobilityd_client_p = std::make_unique<MockMobilitydClient>();
    mobilityd_client        = mobilityd_client_p.get();

    pkt_generator = std::make_unique<PDUGenerator>(
        PKT_DST_MAC, PKT_SRC_MAC, 2, 4, std::move(proxy_connector_p),
        std::move(mobilityd_client_p));
  }

  MockProxyConnector* proxy_connector;
  MockMobilitydClient* mobilityd_client;
  std::unique_ptr<PDUGenerator> pkt_generator;
};

TEST_F(PDUGeneratorTest, test_pdu_generator) {
  struct pcap_pkthdr* phdr =
      (struct pcap_pkthdr*) malloc(sizeof(struct pcap_pkthdr));
  phdr->len       = sizeof(struct ether_header) + sizeof(struct ip);
  phdr->ts.tv_sec = 92;
  u_char* pdata   = reinterpret_cast<u_char*>(
      malloc(sizeof(struct ether_header) + sizeof(struct ip)));
  struct ether_header* ethernetHeader = (struct ether_header*) pdata;
  ethernetHeader->ether_type          = htons(ETHERTYPE_IP);

  struct ip* ipHeader     = (struct ip*) (pdata + sizeof(struct ether_header));
  ipHeader->ip_src.s_addr = 3232235522;
  ipHeader->ip_dst.s_addr = 3232235521;

  EXPECT_CALL(
    *mobilityd_client, GetSubscriberIDFromIP(_, _))
	.Times(1)
        .WillOnce(testing::DoAll(testing::SetArgPointee<1>("imsi1234"), testing::Return(0)));
	
  EXPECT_CALL(
    *proxy_connector, send_data(testing::_, testing::_))
	.Times(1)
	.WillOnce(testing::Return(1));

  pkt_generator->process_packet(phdr, pdata);

  free(pdata);
  free(phdr);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace lte
}  // namespace magma
