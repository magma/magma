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

#include <arpa/inet.h>
#include <gtest/gtest.h>
#include <netinet/in.h>
#include <string.h>
#include <cstdint>
#include <string>

extern "C" {
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/rfc_1332.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_pco.h"
}

namespace magma {
namespace lte {

#define PCO_PI_IPCP_LEN 16

#define DEFAULT_DNS_PRIMARY_HEX 0x8080808  // 8.8.8.8
#define DEFAULT_DNS_PRIMARY_ARRAY                                              \
  { 0x08, 0x08, 0x08, 0x08 }
#define DEFAULT_DNS_SECONDARY_HEX 0x4040808  // 8.8.4.4
#define DEFAULT_DNS_SECONDARY_ARRAY                                            \
  { 0x08, 0x08, 0x04, 0x04 }
#define DEFAULT_DNS_IPV6 "2001:4860:4860:0:0:0:0:8888"

#define DEFAULT_PCSCF_IPV4 0x96171bac  // "172.27.23.150"
#define DEFAULT_PCSCF_IPV4_ARRAY                                               \
  { 0xac, 0x1b, 0x17, 0x96 }
#define DEFAULT_PCSCF_IPV6 "2a12:577:9941:f99c:0002:0001:c731:f114"

#define DEFAULT_MTU 1400
#define PCO_IDS_EXPECT_EQ(ipcp, dns, ipnas, dhcp, mtu)                         \
  do {                                                                         \
    EXPECT_EQ(pco_ids.pi_ipcp, ipcp);                                          \
    EXPECT_EQ(pco_ids.ci_dns_server_ipv4_address_request, dns);                \
    EXPECT_EQ(pco_ids.ci_ip_address_allocation_via_nas_signalling, ipnas);     \
    EXPECT_EQ(pco_ids.ci_ipv4_address_allocation_via_dhcpv4, dhcp);            \
    EXPECT_EQ(pco_ids.ci_ipv4_link_mtu_request, mtu);                          \
  } while (0)

class SPGWPcoTest : public ::testing::Test {
  virtual void SetUp() {
    spgw_config_init(&spgw_config);
    spgw_config.pgw_config.ipv4.default_dns.s_addr = DEFAULT_DNS_PRIMARY_HEX;
    spgw_config.pgw_config.ipv4.default_dns_sec.s_addr =
        DEFAULT_DNS_SECONDARY_HEX;
    spgw_config.pgw_config.ue_mtu = DEFAULT_MTU;

    inet_pton(
        AF_INET6, test_dns_ipv6.c_str(),
        &spgw_config.pgw_config.ipv6.dns_ipv6_addr);
  }

  virtual void TearDown() { free_spgw_config(&spgw_config); }

 protected:
  const char test_dns_primary[4]   = DEFAULT_DNS_PRIMARY_ARRAY;
  const char test_dns_secondary[4] = DEFAULT_DNS_SECONDARY_ARRAY;
  const std::string test_dns_ipv6  = DEFAULT_DNS_IPV6;
  const char test_mtu[2]           = {DEFAULT_MTU >> 8, DEFAULT_MTU & 0xFF};
  const uint8_t test_pcscf_ipv4_addr[4]  = DEFAULT_PCSCF_IPV4_ARRAY;
  const std::string test_pcscf_ipv6_addr = DEFAULT_PCSCF_IPV6;

  void fill_ipcp(
      pco_protocol_or_container_id_t* poc_id, char* primary_dns,
      char* secondary_dns) {
    poc_id->id     = PCO_PI_IPCP;
    poc_id->length = PCO_PI_IPCP_LEN;

    int poc_idx = 0;
    char poc_content[PCO_PI_IPCP_LEN];

    poc_content[0]  = 0x01;  // Code = 01 , i.e. Config Request
    poc_content[1]  = 0x00;  // Identifier = 00
    poc_content[2]  = 0x00;  // Length = 0x0010 , i.e. 16
    poc_content[3]  = 0x10;
    poc_content[4]  = 0x81;  // Option: 81 for primary DNS IP addr
    poc_content[5]  = 0x06;  // length = 6
    poc_content[6]  = primary_dns[0];
    poc_content[7]  = primary_dns[1];
    poc_content[8]  = primary_dns[2];
    poc_content[9]  = primary_dns[3];
    poc_content[10] = 0x83;  // Option: 83 for secondary DNS IP addr
    poc_content[11] = 0x06;  // length = 6
    poc_content[12] = secondary_dns[0];
    poc_content[13] = secondary_dns[1];
    poc_content[14] = secondary_dns[2];
    poc_content[15] = secondary_dns[3];

    poc_id->contents = blk2bstr(poc_content, PCO_PI_IPCP_LEN);
  }

  void clear_pco(protocol_configuration_options_t* pco) {
    for (int i = 0; i < pco->num_protocol_or_container_id; i++) {
      bdestroy_wrapper(&pco->protocol_or_container_ids[i].contents);
    }
  }

  bool are_force_push_pcos_valid(
      protocol_configuration_options_t* pco_resp, bool expected_dns,
      bool expected_mtu) {
    bool has_dns = false;
    bool has_mtu = false;
    for (int i = 0; i < pco_resp->num_protocol_or_container_id; i++) {
      switch (pco_resp->protocol_or_container_ids[i].id) {
        case PCO_CI_DNS_SERVER_IPV4_ADDRESS:
          if (memcmp(
                  pco_resp->protocol_or_container_ids[i].contents->data,
                  test_dns_primary, sizeof(test_dns_primary)) == 0) {
            has_dns = true;
          }
          break;

        case PCO_CI_IPV4_LINK_MTU:
          if (memcmp(
                  pco_resp->protocol_or_container_ids[i].contents->data,
                  test_mtu, sizeof(test_mtu)) == 0) {
            has_mtu = true;
          }
          break;
      }
    }
    if ((has_dns == expected_dns) && (has_mtu == expected_mtu)) {
      return true;
    }
    return false;
  }
};

TEST_F(SPGWPcoTest, TestIPCPWithNoDNS) {
  status_code_e return_code                 = RETURNerror;
  protocol_configuration_options_t pco_resp = {};
  pco_protocol_or_container_id_t poc_id     = {};

  char no_dns[4] = {0x00, 0x00, 0x00, 0x00};

  fill_ipcp(&poc_id, no_dns, no_dns);

  return_code = pgw_process_pco_request_ipcp(&pco_resp, &poc_id);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  // compare the values in pco_resp with those in the poc_id
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].id, PCO_PI_IPCP);
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].length, PCO_PI_IPCP_LEN);
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[1],
      poc_id.contents->data[1]);  // Identifier is same as poc_id

  // check that return code is NACK
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[0],
      IPCP_CODE_CONFIGURE_NACK);

  // check that DNS addresses are filled correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 6,
          test_dns_primary, sizeof(test_dns_primary)),
      0);

  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 12,
          test_dns_secondary, sizeof(test_dns_secondary)),
      0);

  bdestroy_wrapper(&poc_id.contents);
  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestIPCPWithRandomDNS) {
  status_code_e return_code                 = RETURNerror;
  protocol_configuration_options_t pco_resp = {};
  pco_protocol_or_container_id_t poc_id     = {};

  char primary_dns[4]   = {0x01, 0x02, 0x03, 0x04};  // 1.2.3.4
  char secondary_dns[4] = {0x05, 0x06, 0x07, 0x08};  // 5.6.7.8

  fill_ipcp(&poc_id, primary_dns, secondary_dns);

  return_code = pgw_process_pco_request_ipcp(&pco_resp, &poc_id);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  // compare the values in pco_resp with those in the poc_id
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].id, PCO_PI_IPCP);
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].length, PCO_PI_IPCP_LEN);
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[1],
      poc_id.contents->data[1]);  // Identifier is same as poc_id

  // check that return code is NACK
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[0],
      IPCP_CODE_CONFIGURE_NACK);

  // check that DNS addresses are filled correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 6,
          test_dns_primary, sizeof(test_dns_primary)),
      0);
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 12,
          test_dns_secondary, sizeof(test_dns_secondary)),
      0);

  bdestroy_wrapper(&poc_id.contents);
  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestIPCPWithMatchingDNS) {
  status_code_e return_code                 = RETURNerror;
  protocol_configuration_options_t pco_resp = {};
  pco_protocol_or_container_id_t poc_id     = {};

  char primary_dns[4]   = DEFAULT_DNS_PRIMARY_ARRAY;
  char secondary_dns[4] = DEFAULT_DNS_SECONDARY_ARRAY;

  fill_ipcp(&poc_id, primary_dns, secondary_dns);

  return_code = pgw_process_pco_request_ipcp(&pco_resp, &poc_id);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  // compare the values in pco_resp with those in the poc_id
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].id, PCO_PI_IPCP);
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].length, PCO_PI_IPCP_LEN);
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[1],
      poc_id.contents->data[1]);  // Identifier is same as poc_id

  // check that return code is ACK
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[0],
      IPCP_CODE_CONFIGURE_ACK);

  // check that DNS addresses are filled correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 6,
          test_dns_primary, sizeof(test_dns_primary)),
      0);

  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 12,
          test_dns_secondary, sizeof(test_dns_secondary)),
      0);

  bdestroy_wrapper(&poc_id.contents);
  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestIpv4DnsServerRequest) {
  status_code_e return_code                 = RETURNerror;
  protocol_configuration_options_t pco_resp = {};

  return_code = pgw_process_pco_dns_server_request(&pco_resp, NULL);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].id, PCO_CI_DNS_SERVER_IPV4_ADDRESS);

  // check that DNS is assigned correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data,
          test_dns_primary, sizeof(test_dns_primary)),
      0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestLinkMtuRequest) {
  status_code_e return_code = RETURNerror;

  protocol_configuration_options_t pco_resp = {};

  return_code = pgw_process_pco_link_mtu_request(&pco_resp, NULL);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].id, PCO_CI_IPV4_LINK_MTU);

  // check that MTU was assigned correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data, test_mtu,
          sizeof(test_mtu)),
      0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestIpv6DNS) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id = 1;
  pco_req.protocol_or_container_ids[0].id =
      PCO_CI_DNS_SERVER_IPV6_ADDRESS_REQUEST;

  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].id, PCO_CI_DNS_SERVER_IPV6_ADDRESS);

  // check that Ipv6 DNS is assigned correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data,
          spgw_config.pgw_config.ipv6.dns_ipv6_addr.s6_addr,
          sizeof(struct in6_addr)),
      0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestPcscfIpv4) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id    = 1;
  pco_req.protocol_or_container_ids[0].id = PCO_CI_P_CSCF_IPV4_ADDRESS_REQUEST;

  // process PCO for PCSCF without initializing SPGW config
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 0);

  // Initialize SPGW config
  spgw_config.pgw_config.pcscf.ipv4_addr.s_addr = DEFAULT_PCSCF_IPV4;

  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].id, PCO_CI_P_CSCF_IPV4_ADDRESS);

  // check that Ipv4 PCSCF is assigned correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data,
          test_pcscf_ipv4_addr, sizeof(test_pcscf_ipv4_addr)),
      0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestPcscfIpv6) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id    = 1;
  pco_req.protocol_or_container_ids[0].id = PCO_CI_P_CSCF_IPV6_ADDRESS_REQUEST;

  // process PCO for PCSCF without initializing SPGW config
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 0);

  // Initialize SPGW config
  inet_pton(
      AF_INET6, test_pcscf_ipv6_addr.c_str(),
      &spgw_config.pgw_config.pcscf.ipv6_addr);

  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].id, PCO_CI_P_CSCF_IPV6_ADDRESS);

  // check that Ipv4 PCSCF is assigned correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data,
          spgw_config.pgw_config.pcscf.ipv6_addr.s6_addr,
          sizeof(struct in6_addr)),
      0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestNasSignallingIPCP) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  char no_dns[4] = {0x00, 0x00, 0x00, 0x00};

  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id = 2;

  pco_req.protocol_or_container_ids[0].id =
      PCO_CI_IP_ADDRESS_ALLOCATION_VIA_NAS_SIGNALLING;

  fill_ipcp(&pco_req.protocol_or_container_ids[1], no_dns, no_dns);

  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);

  EXPECT_EQ(return_code, RETURNok);

  // Only one container is added in pco_resp for IPCP
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);

  PCO_IDS_EXPECT_EQ(1, 0, 1, 0, 0);

  // compare the values in pco_resp with those in the poc_id
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].id, PCO_PI_IPCP);
  EXPECT_EQ(pco_resp.protocol_or_container_ids[0].length, PCO_PI_IPCP_LEN);
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[1],
      pco_req.protocol_or_container_ids[1]
          .contents->data[1]);  // Identifier is same as poc_id

  // check that return code is NACK
  EXPECT_EQ(
      pco_resp.protocol_or_container_ids[0].contents->data[0],
      IPCP_CODE_CONFIGURE_NACK);

  // check that DNS addresses are filled correctly
  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 6,
          test_dns_primary, sizeof(test_dns_primary)),
      0);

  EXPECT_EQ(
      memcmp(
          pco_resp.protocol_or_container_ids[0].contents->data + 12,
          test_dns_secondary, sizeof(test_dns_secondary)),
      0);

  bdestroy_wrapper(&pco_req.protocol_or_container_ids[1].contents);
  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestForcePush) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  // pco request without poc_ids
  pco_req.num_protocol_or_container_id = 0;
  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;

  // Disable PCO force push
  spgw_config.pgw_config.force_push_pco = false;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 0);

  // Enable PCO force push
  spgw_config.pgw_config.force_push_pco = true;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 2);
  // Check that both MTU and DNS PCOs are added in PCO response
  EXPECT_TRUE(are_force_push_pcos_valid(&pco_resp, true, true));
  // check that only NAS signalling flag is set to true in pco_ids
  PCO_IDS_EXPECT_EQ(0, 0, 1, 0, 0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestDNSReqForcePush) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  // pco request with DNS request
  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id = 1;
  pco_req.protocol_or_container_ids[0].id =
      PCO_CI_DNS_SERVER_IPV4_ADDRESS_REQUEST;

  // Disable PCO force push
  spgw_config.pgw_config.force_push_pco = false;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);
  // Check that only DNS PCO is added in PCO response
  EXPECT_TRUE(are_force_push_pcos_valid(&pco_resp, true, false));
  // Check that only DNS flag is set in pco_ids
  PCO_IDS_EXPECT_EQ(0, 1, 0, 0, 0);
  clear_pco(&pco_resp);

  // Enable PCO force push
  spgw_config.pgw_config.force_push_pco = true;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 2);
  // Check that both MTU and DNS PCOs are added in PCO response
  EXPECT_TRUE(are_force_push_pcos_valid(&pco_resp, true, true));
  // Check that DNS and NAS flag is set in pco_ids
  PCO_IDS_EXPECT_EQ(0, 1, 1, 0, 0);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestMTUReqForcePush) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  // pco request with MTU request
  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  pco_req.num_protocol_or_container_id    = 1;
  pco_req.protocol_or_container_ids[0].id = PCO_CI_IPV4_LINK_MTU;

  // Disable PCO force push
  spgw_config.pgw_config.force_push_pco = false;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 1);
  // Check that only MTU PCO is added in PCO response
  EXPECT_TRUE(are_force_push_pcos_valid(&pco_resp, false, true));
  // Check that only MTU flag is set in pco_ids
  PCO_IDS_EXPECT_EQ(0, 0, 0, 0, 1);
  clear_pco(&pco_resp);

  // Enable PCO force push
  spgw_config.pgw_config.force_push_pco = true;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.num_protocol_or_container_id, 2);
  // Check that both MTU and DNS PCOs are added in PCO response
  EXPECT_TRUE(are_force_push_pcos_valid(&pco_resp, true, true));
  // Check that MTU and NAS flag is set in pco_ids
  PCO_IDS_EXPECT_EQ(0, 0, 1, 0, 1);

  clear_pco(&pco_resp);
}

TEST_F(SPGWPcoTest, TestPcoRequestConfigurationProtocol) {
  status_code_e return_code                    = RETURNerror;
  protocol_configuration_options_t pco_req     = {};
  protocol_configuration_options_t pco_resp    = {};
  protocol_configuration_options_ids_t pco_ids = {};

  // pco request without poc_ids
  pco_req.num_protocol_or_container_id = 0;

  // PCO request with random configuration protocol
  pco_req.configuration_protocol =
      1 +
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.configuration_protocol, 0);
  PCO_IDS_EXPECT_EQ(0, 0, 0, 0, 0);

  // PCO request with configuration protocol set
  pco_req.configuration_protocol =
      PCO_CONFIGURATION_PROTOCOL_PPP_FOR_USE_WITH_IP_PDP_TYPE_OR_IP_PDN_TYPE;
  return_code = pgw_process_pco_request(&pco_req, &pco_resp, &pco_ids);
  EXPECT_EQ(return_code, RETURNok);
  EXPECT_EQ(pco_resp.configuration_protocol, pco_req.configuration_protocol);
  PCO_IDS_EXPECT_EQ(0, 0, 0, 0, 0);
  clear_pco(&pco_resp);
}

}  // namespace lte
}  // namespace magma
