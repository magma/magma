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

#include "M5GAmfMock.h"
#include "amf_app_messages_types.h"
#include "M5GAuthenticationServiceClient.h"
#include "Consts.h"
#include "include/amf_client_servicer.h"

#include <gtest/gtest.h>
#include <glog/logging.h>

#include <grpcpp/impl/codegen/status.h>

using ::testing::Test;

namespace magma5g {

class AMFClientServerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    subscriberdb_client =
        std::make_shared<MockM5GAuthenticationServiceClient>();
    amf_client = std::make_unique<AmfClientServicer>(subscriberdb_client);
  }
  virtual void TearDown() {}

 protected:
  std::unique_ptr<AmfClientServicer> amf_client;
  std::shared_ptr<MockM5GAuthenticationServiceClient> subscriberdb_client;
};

TEST_F(AMFClientServerTest, test_get_subs_auth_info) {
  EXPECT_CALL(
      *subscriberdb_client,
      get_subs_auth_info(testing::_, testing::_, testing::_, testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));

  auto succeeded = amf_client->get_subscriber_authentication_info(
      IMSI1, IMSI_LENGTH, SNNI_STR, UE_IDENTITY);

  EXPECT_TRUE(succeeded);
}

TEST_F(AMFClientServerTest, test_get_subs_auth_info_resync) {
  unsigned char* data =
      (unsigned char*) malloc(sizeof(unsigned char) * RESYNC_INFO_LEN);

  EXPECT_CALL(
      *subscriberdb_client, get_subs_auth_info_resync(
                                testing::_, testing::_, testing::_, testing::_,
                                testing::_, testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));

  auto succeeded = amf_client->get_subscriber_authentication_info_resync(
      IMSI1, IMSI_LENGTH, SNNI_STR, data, RESYNC_INFO_LEN, UE_IDENTITY);

  EXPECT_TRUE(succeeded);
  free(data);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma5g
