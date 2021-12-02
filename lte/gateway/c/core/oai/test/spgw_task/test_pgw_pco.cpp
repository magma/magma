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

#include <gtest/gtest.h>
#include <string.h>

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

class SPGWPcoTest : public ::testing::Test {
  virtual void SetUp() {
    spgw_config_init(&spgw_config);
    spgw_config.pgw_config.ipv4.default_dns.s_addr = DEFAULT_DNS_PRIMARY_HEX;
    spgw_config.pgw_config.ipv4.default_dns_sec.s_addr =
        DEFAULT_DNS_SECONDARY_HEX;
  }

  virtual void TearDown() { free_spgw_config(&spgw_config); }

 protected:
  char test_dns_primary[4]   = DEFAULT_DNS_PRIMARY_ARRAY;
  char test_dns_secondary[4] = DEFAULT_DNS_SECONDARY_ARRAY;

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
}  // namespace lte
}  // namespace magma
