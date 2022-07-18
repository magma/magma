#include <string.h>
#include <gtest/gtest.h>

#include "feg/protos/s6a_proxy.pb.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_utility_funs.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/S6aClient.hpp"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/itti_msg_to_proto_msg.hpp"
#include "lte/protos/mconfig/mconfigs.pb.h"
#include "orc8r/gateway/c/common/config/MConfigLoader.hpp"
#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"
extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
}

namespace magma{
    using namespace feg;


        TEST(ULA2SubDataTest, TestULAallFields){

            //create UpdateLocationAnswer object
            auto ula_object = magma::feg::UpdateLocationAnswer();

            auto apn_config = ula_object.add_apn();
            apn_config->set_context_id(1);
            apn_config->set_service_selection("abc");
            auto qos_profile_msg = apn_config->mutable_qos_profile();
            qos_profile_msg->set_class_id(1);
            qos_profile_msg->set_priority_level(2);
            qos_profile_msg->set_preemption_capability(true);
            qos_profile_msg->set_preemption_vulnerability(true);

            apn_config->set_allocated_qos_profile(qos_profile_msg);

            auto ambr_msg = apn_config->mutable_ambr();
            ambr_msg->set_max_bandwidth_ul(100000);
            ambr_msg->set_max_bandwidth_dl(200000);
            ambr_msg->set_unit((magma::feg::UpdateLocationAnswer_AggregatedMaximumBitrate_BitrateUnitsAMBR)BPS);

            apn_config->set_allocated_ambr(ambr_msg);

            apn_config->set_pdn((magma::feg::UpdateLocationAnswer_APNConfiguration_PDNType)IPV4);

            auto resource_msg = apn_config->mutable_resource();
            resource_msg->set_apn_name("apn_name");
            resource_msg->set_gateway_ip("0.0.0.0");
            resource_msg->set_gateway_mac("A:B:C:D");
            resource_msg->set_vlan_id(123);

            apn_config->set_allocated_resource(resource_msg);

            apn_config->add_served_party_ip_address("123.123.123.123");



            // initialize subscriberData object

            magma::lte::SubscriberData sub_data = magma::lte::SubscriberData();
            magma::lte::SubscriberID sub_id = magma::lte::SubscriberID();
            sub_id.set_id("IMSI123123123");
            sub_id.set_type(magma::lte::SubscriberID::IMSI);
            sub_data.set_allocated_sid(&sub_id);

            // call data converting function

            S6aClient::convert_ula_to_subscriber_data(ula_object, sub_data);

            // test equality for each of the fields in the subscriberdata object

            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).context_id(), 1);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).service_selection(), "abc");
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().class_id(), 1);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().priority_level(), 2);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().preemption_capability(), true);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).qos_profile().preemption_vulnerability(), true);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().max_bandwidth_dl(), 200000);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().max_bandwidth_ul(), 100000);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).ambr().br_unit(), BPS);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).pdn(), IPV4);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().apn_name(), "apn_name");
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().gateway_ip(), "0.0.0.0");
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().gateway_mac(), "A:B:C:D");
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).resource().vlan_id(), 123);
            EXPECT_EQ(sub_data.non_3gpp().apn_config(0).assigned_static_ip(), "123.123.123.123");





    }

}
