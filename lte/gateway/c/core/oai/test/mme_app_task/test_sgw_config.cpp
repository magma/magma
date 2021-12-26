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
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include <libconfig.h>

extern "C" {
#include "lte/gateway/c/core/oai/include/sgw_config.h"
}

namespace magma {
namespace lte {

const char* sEmptyConfig =
    R"libconfig(S-GW :
{
};)libconfig";

const char* sHealthyConfig =
    R"libconfig(S-GW :
{
    NETWORK_INTERFACES :
    {
        # S-GW binded interface for S11 communication (GTPV2-C), if none selected the ITTI message interface is used
        SGW_INTERFACE_NAME_FOR_S11              = "lo";
        SGW_IPV4_ADDRESS_FOR_S11                = "127.0.0.1/8";

        # S-GW binded interface for S1-U communication (GTPV1-U) can be ethernet interface, virtual ethernet interface, we don't advise wireless interfaces
        SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP    = "eth1";
        SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP      = "192.168.60.142/24";
        SGW_IPV4_PORT_FOR_S1U_S12_S4_UP         = 2152;                         # INTEGER, port number, PREFER NOT CHANGE UNLESS YOU KNOW WHAT YOU ARE DOING

        # S-GW binded interface for S5 or S8 communication, not implemented, so leave it to none
        SGW_INTERFACE_NAME_FOR_S5_S8_UP         = "eth0";         # STRING, interface name
        SGW_IPV4_ADDRESS_FOR_S5_S8_UP           = "10.0.2.15/24";                 # STRING, CIDR

        SGW_IPV6_ADDRESS_FOR_S1U_S12_S4_UP      = "2001:db8::1234:5678";
        SGW_IPV6_PORT_FOR_S1U_S12_S4_UP         = 2152;
        SGW_S1_IPV6_ENABLED = "True";
    };

    INTERTASK_INTERFACE :
    {
        # max queue size per task
        ITTI_QUEUE_SIZE            = 2000000;                                   # INTEGER
    };

    OVS :
    {
      BRIDGE_NAME                          = "gtp_br0";
      GTP_PORT_NUM                         = 32768;
      MTR_PORT_NUM                         = 15577;
      INTERNAL_SAMPLING_PORT_NUM           = 15578;
      INTERNAL_SAMPLING_FWD_TBL_NUM        = 201;
      UPLINK_PORT_NUM                      = 2;
      UPLINK_MAC                           = "ff:ff:ff:ff:ff:ff";
      MULTI_TUNNEL                         = "True";
      GTP_ECHO                             = "True";
      GTP_CHECKSUM                         = "False";
      AGW_L3_TUNNEL                        = "False";
      PIPELINED_CONFIG_ENABLED             = "False";
    };
};)libconfig";

TEST(SGWConfigTest, TestParseHealthyConfig) {
  sgw_config_t sgw_config = {0};
  EXPECT_EQ(sgw_config_parse_string(sHealthyConfig, &sgw_config), 0);
  EXPECT_EQ(std::string(bdata(sgw_config.ovs_config.bridge_name)), "gtp_br0");
  EXPECT_EQ(std::string(bdata(sgw_config.ipv6.if_name_S1u_S12_S4_up)), "eth1");
  ASSERT_EQ(sgw_config.ipv6.s1_ipv6_enabled, true);
  free_sgw_config(&sgw_config);
}

TEST(SGWConfigTest, TestParseHealthyConfigDisplay) {
  sgw_config_t sgw_config = {0};
  EXPECT_EQ(sgw_config_parse_string(sHealthyConfig, &sgw_config), 0);
  sgw_config_display(&sgw_config);
  free_sgw_config(&sgw_config);
}

}  // namespace lte
}  // namespace magma
