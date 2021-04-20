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

#include "includes/ServiceRegistrySingleton.h"
#include <gtest/gtest.h>

using magma::ServiceRegistrySingleton;
using ::testing::Test;
namespace magma {

// This test relies on the values set in service_registry.yml and
// control_proxy.yml in /etc/magma/.
// If these tests start failing all of a sudden, see if those values have
// changed recently.

TEST(TestServiceRegistry, test_get_remote_channel_args) {
  auto args = ServiceRegistrySingleton::Instance()->GetCreateGrpcChannelArgs(
      "state", "cloud");
  // The default config has `proxy_cloud_connections` set
  EXPECT_EQ(args.ip, "127.0.0.1");
  EXPECT_EQ(args.port, "8443");  // proxy for local services
  EXPECT_EQ(args.authority, "state-controller.magma.test");
}

TEST(TestServiceRegistry, test_get_local_channel_args) {
  auto args = ServiceRegistrySingleton::Instance()->GetCreateGrpcChannelArgs(
      "mobilityd", "local");
  // These values should match service_registry.yml
  EXPECT_EQ(args.ip, "127.0.0.1");
  EXPECT_EQ(args.port, "60051");
  EXPECT_EQ(args.authority, "mobilityd.local");
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma
