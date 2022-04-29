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

#include <gmock/gmock.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/status_code_enum.h>
#include <gtest/gtest.h>
#include <lte/protos/subscriberdb.pb.h>
#include <net/ethernet.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <pcap.h>
#include <stdlib.h>
#include <sys/time.h>
#include <sys/types.h>
#include <limits>
#include <memory>
#include <string>
#include <utility>

#include "lte/gateway/c/li_agent/src/PDUGenerator.hpp"
#include "lte/gateway/c/li_agent/src/test/Consts.hpp"
#include "lte/gateway/c/li_agent/src/test/LIAgentdMocks.hpp"

namespace magma {
namespace lte {

class PDUGeneratorTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    std::string target_id = "IMSI12345";
    std::string task_id = "29f28e1c-f230-486a-a860-f5a784ab9177";
    auto mconfig = create_liagentd_mconfig(task_id, target_id);

    int sync_time = std::numeric_limits<int>::max();  // Prevent sync

    auto proxy_connector_p = std::make_unique<MockProxyConnector>();
    proxy_connector = proxy_connector_p.get();

    auto mobilityd_client_p = std::make_unique<MockMobilitydClient>();
    mobilityd_client = mobilityd_client_p.get();

    pkt_generator = std::make_unique<PDUGenerator>(
        PKT_DST_MAC, PKT_SRC_MAC, sync_time, sync_time,
        std::move(proxy_connector_p), std::move(mobilityd_client_p), mconfig);
  }

  MockProxyConnector* proxy_connector;
  MockMobilitydClient* mobilityd_client;
  std::unique_ptr<PDUGenerator> pkt_generator;
};

TEST_F(PDUGeneratorTest, test_pdu_generator) {
  struct pcap_pkthdr* phdr =
      (struct pcap_pkthdr*)malloc(sizeof(struct pcap_pkthdr));
  phdr->len = sizeof(struct ether_header) + sizeof(struct ip);
  phdr->ts.tv_sec = 56;
  u_char* pdata = reinterpret_cast<u_char*>(
      malloc(sizeof(struct ether_header) + sizeof(struct ip)));
  struct ether_header* ethernetHeader = (struct ether_header*)pdata;
  ethernetHeader->ether_type = htons(ETHERTYPE_IP);

  struct ip* ipHeader = (struct ip*)(pdata + sizeof(struct ether_header));
  ipHeader->ip_src.s_addr = 3232235522;
  ipHeader->ip_dst.s_addr = 3232235521;

  EXPECT_CALL(*proxy_connector, send_data(testing::_, testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));

  SubscriberID response;
  response.set_id("12345");
  EXPECT_CALL(*mobilityd_client,
              get_subscriber_id_from_ip(testing::_, testing::_))
      .WillRepeatedly(testing::InvokeArgument<1>(grpc::Status::OK, response));

  auto succeeded = pkt_generator->process_packet(phdr, pdata);
  EXPECT_TRUE(succeeded);
  free(pdata);
  free(phdr);
}

TEST_F(PDUGeneratorTest, test_generator_unknown_subscriber) {
  struct pcap_pkthdr* phdr =
      (struct pcap_pkthdr*)malloc(sizeof(struct pcap_pkthdr));
  phdr->len = sizeof(struct ether_header) + sizeof(struct ip);
  phdr->ts.tv_sec = 56;
  u_char* pdata = reinterpret_cast<u_char*>(
      malloc(sizeof(struct ether_header) + sizeof(struct ip)));
  struct ether_header* ethernetHeader = (struct ether_header*)pdata;
  ethernetHeader->ether_type = htons(ETHERTYPE_IP);

  struct ip* ipHeader = (struct ip*)(pdata + sizeof(struct ether_header));
  ipHeader->ip_src.s_addr = 3232235522;

  SubscriberID response;
  EXPECT_CALL(*mobilityd_client,
              get_subscriber_id_from_ip(testing::_, testing::_))
      .WillRepeatedly(testing::InvokeArgument<1>(
          grpc::Status(grpc::DEADLINE_EXCEEDED, "timeout"), response));

  auto succeeded = pkt_generator->process_packet(phdr, pdata);
  EXPECT_FALSE(succeeded);
  free(pdata);
  free(phdr);
}

TEST_F(PDUGeneratorTest, test_generator_non_ip_packet) {
  struct pcap_pkthdr* phdr =
      (struct pcap_pkthdr*)malloc(sizeof(struct pcap_pkthdr));
  phdr->len = sizeof(struct ether_header);
  u_char* pdata =
      reinterpret_cast<u_char*>(malloc(sizeof(struct ether_header)));
  struct ether_header* ethernetHeader = (struct ether_header*)pdata;
  ethernetHeader->ether_type = htons(ETHERTYPE_ARP);

  auto succeeded = pkt_generator->process_packet(phdr, pdata);
  EXPECT_FALSE(succeeded);

  free(pdata);
  free(phdr);
}
}  // namespace lte
}  // namespace magma
