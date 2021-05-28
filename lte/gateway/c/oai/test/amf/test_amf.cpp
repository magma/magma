/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "util_nas5g_pkt.h"
#include "dynamic_memory_check.h"
#include <gtest/gtest.h>
#include <thread>

using ::testing::Test;

namespace magma5g{

uint8_t NAS5GPktSnapShot::reg_req_buffer[38] =
          {0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x09,
           0xf1, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
           0x00, 0x00, 0x10, 0x10, 0x01, 0x00, 0x2e, 0x04,
           0xf0, 0xf0, 0xf0, 0xf0, 0x2f, 0x05, 0x04, 0x01,
           0x00, 0x00, 0x01, 0x53, 0x01, 0x00};

TEST(test_amf_nas5g_pkt_process, test_amf_ue_register_req_msg) {

   NAS5GPktSnapShot nas5g_pkt_snap;
   RegistrationRequestMsg reg_request;
   bool decode_res = false;

   uint32_t len = nas5g_pkt_snap.get_reg_req_buffer_len();

   memset(&reg_request, 0, sizeof(RegistrationRequestMsg)); 
   decode_res = decode_registration_request_msg(&reg_request,
                                                nas5g_pkt_snap.reg_req_buffer, len);
   EXPECT_TRUE(decode_res == true);

   EXPECT_TRUE(reg_request.extended_protocol_discriminator.
	        extended_proto_discriminator ==
	       M5G_MOBILITY_MANAGEMENT_MESSAGES);

   //Type is registration Request
   EXPECT_TRUE(reg_request.message_type.msg_type == REG_REQUEST); 

   //Registraiton Type is Initial Registration
   EXPECT_TRUE(reg_request.m5gs_reg_type.type_val == 1); 

   //5GS Mobile Identity SUPI FORMAT
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.type_of_identity == M5GSMobileIdentityMsg_SUCI_IMSI);
   
   //5GS Mobile mms digit2
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.mcc_digit1 == 0x09);
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.mcc_digit2 == 0x00);
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.mcc_digit3 == 0x01);
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.mnc_digit1 == 0x07);
   EXPECT_TRUE(reg_request.m5gs_mobile_identity.mobile_identity.
	       imsi.mcc_digit2 == 0x0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma5g


