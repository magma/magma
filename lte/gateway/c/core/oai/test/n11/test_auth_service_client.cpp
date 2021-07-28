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
#include <glog/logging.h>

#include <grpcpp/impl/codegen/status.h>
#include "lte/protos/subscriberauth.pb.h"
#include "lte/protos/subscriberauth.grpc.pb.h"

#include "amf_client_proto_msg_to_itti_msg.h"
#include "amf_app_messages_types.h"
#include "M5GAuthenticationServiceClient.h"

using ::testing::Test;

task_zmq_ctx_t grpc_service_task_zmq_ctx;

namespace magma {
namespace lte {

TEST(
    test_convert_proto_msg_to_itti_m5g_auth_info_ans,
    convert_proto_msg_to_itti_m5g_auth_info_ans) {
  magma::lte::M5GAuthenticationInformationAnswer response;
  itti_amf_subs_auth_info_ans_t amf_app_subs_auth_info_resp_p;

  auto* authVector1 = response.add_m5gauth_vectors();
  authVector1->set_rand("rand1");
  authVector1->set_xres_star("xres_star1");
  authVector1->set_autn("autn1");
  authVector1->set_kseaf("kseaf1");

  std::cout << "m5gauth_vectors_size :" << response.m5gauth_vectors_size();
  magma5g::convert_proto_msg_to_itti_m5g_auth_info_ans(
      response, &amf_app_subs_auth_info_resp_p);

  EXPECT_TRUE(
      response.m5gauth_vectors_size() ==
      amf_app_subs_auth_info_resp_p.auth_info.nb_of_vectors)
      << "size of response.m5gauth_vectors_size : "
      << response.m5gauth_vectors_size()
      << " amf_app_subs_auth_info_resp_p         : "
      << amf_app_subs_auth_info_resp_p.auth_info.nb_of_vectors;
}

TEST(test_get_subs_auth_info, get_subs_auth_info) {
  std::string imsi = "901700000000001";
  std::string snni = "5G:mnc070.mcc901.3gppnetwork.org";
  M5GAuthenticationInformationRequest req =
      magma5g::create_subs_auth_request(imsi, snni);

  EXPECT_TRUE(imsi == req.user_name());
  EXPECT_TRUE(snni == req.serving_network_name());
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace lte
}  // namespace magma
