/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include "ServiceRegistrySingleton.h"
#include <gtest/gtest.h>

using magma::ServiceRegistrySingleton;
using ::testing::Test;

// Tests the GetGrpcChannel
TEST(TestServiceRegistry, TestGetGrpcCloudChannelArgs) {
  auto args = ServiceRegistrySingleton::Instance()->GetCreateGrpcChannelArgs(
      "logger", "cloud");
  EXPECT_EQ(args.ip, "127.0.0.1");
  EXPECT_EQ(args.port, "8443");
  EXPECT_EQ(args.authority, "logger-controller.magma.test");
}

TEST(TestServiceRegistry, TestGetGrpcLocalChannelArgs) {
  auto args = ServiceRegistrySingleton::Instance()->GetCreateGrpcChannelArgs(
      "mobilityd", "local");
  EXPECT_EQ(args.ip, "127.0.0.1");
  EXPECT_EQ(args.port, "60051");
  EXPECT_EQ(args.authority, "mobilityd.local");
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
