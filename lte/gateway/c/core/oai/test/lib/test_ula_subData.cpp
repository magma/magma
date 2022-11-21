/*
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
#include <string.h>
#include <gtest/gtest.h>

#include "feg/protos/s6a_proxy.pb.h"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/S6aClient.hpp"

#include "lte/gateway/c/core/oai/include/mme_config.hpp"

namespace magma {

class ULA2SubDataTest : public ::testing::Test {
  virtual void SetUp() { mme_config_init(&mme_config); }

  virtual void TearDown() { free_mme_config(&mme_config); }
};

void create_ula_object(
    magma::feg::UpdateLocationAnswer *ula_object,
    google::protobuf::uint32 context_id, std::string service_selection,
    google::protobuf::int32 class_id, google::protobuf::uint32 priority_level,
    bool preemption_capability, bool preemption_vulnerability,

    std::string apn_name, std::string gateway_mac, std::string gateway_ip,
    google::protobuf::uint32 vlan_id,

    google::protobuf::uint32 max_bandwidth_ul,
    google::protobuf::uint32 max_bandwidth_dl, apn_ambr_bitrate_unit_t unit,

    pdu_session_type_e pdn,

    std::string served_party_ip_address) {
  // create UpdateLocationAnswer object

  magma::feg::UpdateLocationAnswer_APNConfiguration *apn_config =
      ula_object->add_apn();
  apn_config->set_context_id(context_id);
  apn_config->set_service_selection(service_selection);
  auto qos_profile_msg = apn_config->mutable_qos_profile();
  qos_profile_msg->set_class_id(class_id);
  qos_profile_msg->set_priority_level(priority_level);
  qos_profile_msg->set_preemption_capability(preemption_capability);
  qos_profile_msg->set_preemption_vulnerability(preemption_vulnerability);

  auto ambr_msg = apn_config->mutable_ambr();
  ambr_msg->set_max_bandwidth_ul(max_bandwidth_ul);
  ambr_msg->set_max_bandwidth_dl(max_bandwidth_dl);
  ambr_msg->set_unit(
      (magma::feg::
           UpdateLocationAnswer_AggregatedMaximumBitrate_BitrateUnitsAMBR)unit);

  apn_config->set_pdn(
      (magma::feg::UpdateLocationAnswer_APNConfiguration_PDNType)pdn);

  auto resource_msg = apn_config->mutable_resource();
  resource_msg->set_apn_name(apn_name);
  resource_msg->set_gateway_ip(gateway_ip);
  resource_msg->set_gateway_mac(gateway_mac);
  resource_msg->set_vlan_id(vlan_id);

  apn_config->add_served_party_ip_address(served_party_ip_address);
}
void initSubscriber(magma::lte::SubscriberData *sub_data) {
  // initialize subscriberData object
  auto sub_id = sub_data->mutable_sid();
  sub_id->set_id("IMSI123123123");
  sub_id->set_type(magma::lte::SubscriberID::IMSI);
}

TEST(ULA2SubDataTest, TestULAallFields) {
  magma::feg::UpdateLocationAnswer ula_object =
      magma::feg::UpdateLocationAnswer();
  magma::lte::SubscriberData sub_data = magma::lte::SubscriberData();

  google::protobuf::uint32 context_id = 1;
  std::string service_selection = "abc";

  google::protobuf::int32 class_id = 1;
  google::protobuf::uint32 priority_level = 2;
  bool preemption_capability = true;
  bool preemption_vulnerability = false;

  std::string apn_name = "apn_name";
  std::string gateway_mac = "A:B:C:D";
  std::string gateway_ip = "0.0.0.0";
  google::protobuf::uint32 vlan_id = 123;

  google::protobuf::uint32 max_bandwidth_ul = 200000;
  google::protobuf::uint32 max_bandwidth_dl = 100000;
  apn_ambr_bitrate_unit_t unit = BPS;

  pdu_session_type_e pdn = IPV4;

  std::string served_party_ip_address = "123.123.123.123";

  create_ula_object(&ula_object, context_id, service_selection, class_id,
                    priority_level, preemption_capability,
                    preemption_vulnerability, apn_name, gateway_mac, gateway_ip,
                    vlan_id, max_bandwidth_ul, max_bandwidth_dl, unit, pdn,
                    served_party_ip_address);
  initSubscriber(&sub_data);

  // call data converting function

  S6aClient::convert_ula_to_subscriber_data(ula_object, &sub_data);

  // test equality for each of the fields in the subscriberdata object

  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).context_id(), context_id);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).service_selection(),
            service_selection);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().class_id(),
            class_id);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().priority_level(),
            priority_level);
  EXPECT_EQ(
      sub_data.non_3gpp().apn_config(0).qos_profile().preemption_capability(),
      preemption_capability);
  EXPECT_EQ(sub_data.non_3gpp()
                .apn_config(0)
                .qos_profile()
                .preemption_vulnerability(),
            preemption_vulnerability);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().max_bandwidth_dl(),
            max_bandwidth_dl);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().max_bandwidth_ul(),
            max_bandwidth_ul);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().br_unit(), unit);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).pdn(),
            (magma::lte::APNConfiguration_PDNType)pdn);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().apn_name(), apn_name);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().gateway_ip(),
            gateway_ip);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().gateway_mac(),
            gateway_mac);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().vlan_id(), vlan_id);
  EXPECT_EQ(sub_data.non_3gpp().apn_config(0).assigned_static_ip(),
            served_party_ip_address);
}
}  // namespace magma
