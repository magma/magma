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

#include <google/protobuf/message.h>
#include <gtest/gtest.h>

#include <sstream>

#include "orc8r/gateway/c/common/config/includes/MConfigLoader.h"
#include "lte/protos/mconfig/mconfigs.pb.h"

namespace {

static constexpr const char* SERVICE_NAME_MME = "mme";

const char* healthy_mconfig =
    R"proto({
              "configs_by_key": {
                "mobilityd": {
                  "@type": "type.googleapis.com/magma.mconfig.MobilityD",
                  "logLevel": "INFO",
                  "ipBlock": "192.168.128.0/24",
                  "ipv6Block": "fdee:5:6c::/48",
                  "ipv6PrefixAllocationType": "RANDOM",
                  "ip_allocator_type": "IP_POOL",
                  "static_ip_enabled": true,
                  "multi_apn_ip_alloc": true
                },
                "mme": {
                  "@type": "type.googleapis.com/magma.mconfig.MME",
                  "mmeCode": 1,
                  "mmeGid": 1,
                  "mmeRelativeCapacity": 11,
                  "logLevel": "INFO",
                  "mcc": "001",
                  "mnc": "01",
                  "tac": 1,
                  "enableDnsCaching": false,
                  "relayEnabled": false,
                  "hssRelayEnabled": false,
                  "csfbMcc": "001",
                  "dnsPrimary": "8.8.8.8",
                  "dnsSecondary": "8.8.4.4",
                  "ipv6PCscfAddress": "2a12:577:9941:f99c:0002:0001:c731:f114",
                  "ipv6DnsAddress": "2001:4860:4860:0:0:0:0:8888",
                  "enable5gFeatures": false
                }
              }
            })proto";

TEST(MConfigLoader, FailsEmptyStream) {
  magma::mconfig::MME mconfig;
  std::istringstream config_stream("");
  EXPECT_FALSE(
      magma::load_service_mconfig(SERVICE_NAME_MME, &config_stream, &mconfig));
}

TEST(MConfigLoader, HealthyConfigLoads) {
  magma::mconfig::MME mconfig;
  std::istringstream config_stream(healthy_mconfig);
  EXPECT_TRUE(
      magma::load_service_mconfig(SERVICE_NAME_MME, &config_stream, &mconfig));
  EXPECT_EQ(mconfig.tac(), 1);
  EXPECT_EQ(mconfig.ipv6_p_cscf_address(),
            "2a12:577:9941:f99c:0002:0001:c731:f114");
}

TEST(MConfigLoader, MissingServiceNameFails) {
  magma::mconfig::MME mconfig;
  std::istringstream config_stream(healthy_mconfig);
  EXPECT_FALSE(magma::load_service_mconfig("MISSING_SERVICE_NAME",
                                           &config_stream, &mconfig));
}

}  // namespace
