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

TEST(
    test_convert_proto_msg_to_itti_m5g_auth_info_ans,
    convert_proto_msg_to_itti_m5g_auth_info_ans) {
  magma::lte::M5GAuthenticationInformationAnswer response;
  itti_amf_subs_auth_info_ans_t amf_app_subs_auth_info_resp_p;

  std::string rand("rand1");
  std::string xres_star("xres_star1");
  std::string autn("authenticationtn");
  std::string kseaf("SecurityAnchorFunctionAMFKeyOf22");

  auto* authVector1 = response.add_m5gauth_vectors();
  authVector1->set_rand(rand);
  authVector1->set_xres_star(xres_star);
  authVector1->set_autn(autn);
  authVector1->set_kseaf(kseaf);

  magma5g::convert_proto_msg_to_itti_m5g_auth_info_ans(
      response, &amf_app_subs_auth_info_resp_p);

  // build expected itti_amf_subs_auth_info_ans_t
  itti_amf_subs_auth_info_ans_t expect_auth_info;
  expect_auth_info.auth_info.nb_of_vectors = 1;
  m5gauth_vector_t& expected_m5gauth_vector =
      expect_auth_info.auth_info.m5gauth_vector[0];
  memcpy(expected_m5gauth_vector.rand, rand.c_str(), rand.length());
  expected_m5gauth_vector.xres_star.size = xres_star.length();
  memcpy(
      expected_m5gauth_vector.xres_star.data, xres_star.c_str(),
      xres_star.length());
  memcpy(expected_m5gauth_vector.autn, autn.c_str(), autn.length());
  memcpy(expected_m5gauth_vector.kseaf, kseaf.c_str(), kseaf.length());

  // check generated & expected
  m5gauth_vector_t& generated_m5gauth_vector =
      amf_app_subs_auth_info_resp_p.auth_info.m5gauth_vector[0];

  EXPECT_TRUE(
      expect_auth_info.auth_info.nb_of_vectors ==
      amf_app_subs_auth_info_resp_p.auth_info.nb_of_vectors);
  EXPECT_TRUE(
      0 == memcmp(
               expected_m5gauth_vector.rand, generated_m5gauth_vector.rand,
               rand.length()));
  EXPECT_TRUE(
      expected_m5gauth_vector.xres_star.size ==
      generated_m5gauth_vector.xres_star.size);
  EXPECT_TRUE(
      0 == memcmp(
               expected_m5gauth_vector.xres_star.data,
               generated_m5gauth_vector.xres_star.data, xres_star.length()));
  EXPECT_TRUE(
      0 == memcmp(
               expected_m5gauth_vector.autn, generated_m5gauth_vector.autn,
               autn.length()));
  EXPECT_TRUE(
      0 == memcmp(
               expected_m5gauth_vector.kseaf, generated_m5gauth_vector.kseaf,
               kseaf.length()));
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
