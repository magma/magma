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

#include "orc8r/gateway/c/common/service303/includes/MagmaService.h"

using ::testing::Test;

namespace magma {
namespace service303 {

const std::string SERVICE_NAME = "test_service";
const std::string SERVICE_VERSION = "0.0.0";
const std::string META_KEY = "key";
const std::string META_VALUE = "value";

TEST(test_magma_service, test_GetServiceInfo) {
  MagmaService magma_service(SERVICE_NAME, SERVICE_VERSION);
  ServiceInfo response;

  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_EQ(response.name(), SERVICE_NAME);
  EXPECT_EQ(response.version(), SERVICE_VERSION);
  EXPECT_EQ(response.state(), ServiceInfo::ALIVE);
  EXPECT_EQ(response.health(), ServiceInfo::APP_UNKNOWN);
  auto start_time_1 = response.start_time_secs();
  EXPECT_TRUE(response.status().meta().empty());

  response.Clear();

  magma_service.setApplicationHealth(ServiceInfo::APP_HEALTHY);
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_EQ(response.name(), SERVICE_NAME);
  EXPECT_EQ(response.version(), SERVICE_VERSION);
  EXPECT_EQ(response.state(), ServiceInfo::ALIVE);
  EXPECT_EQ(response.health(), ServiceInfo::APP_HEALTHY);
  auto start_time_2 = response.start_time_secs();
  EXPECT_TRUE(response.status().meta().empty());

  EXPECT_EQ(start_time_1, start_time_2);
}

ServiceInfoMeta test_callback() {
  return ServiceInfoMeta{{META_KEY, META_VALUE}};
}

TEST(test_magma_service, test_GetServiceInfo_with_callback) {
  MagmaService magma_service(SERVICE_NAME, SERVICE_VERSION);
  ServiceInfo response;

  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_TRUE(response.status().meta().empty());

  response.Clear();

  magma_service.SetServiceInfoCallback(test_callback);
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  auto meta = response.status().meta();
  EXPECT_FALSE(meta.empty());
  EXPECT_EQ(meta.size(), 1);
  EXPECT_EQ(meta[META_KEY], META_VALUE);

  response.Clear();

  magma_service.ClearServiceInfoCallback();
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_TRUE(response.status().meta().empty());
}

bool reload_succeeded() { return true; }

bool reload_failed() { return false; }

TEST(test_magma_service, test_ReloadServiceConfig) {
  MagmaService magma_service(SERVICE_NAME, SERVICE_VERSION);
  ReloadConfigResponse response;

  magma_service.ReloadServiceConfig(nullptr, nullptr, &response);
  EXPECT_EQ(response.result(), ReloadConfigResponse::RELOAD_UNSUPPORTED);

  response.Clear();

  magma_service.SetConfigReloadCallback(reload_succeeded);
  magma_service.ReloadServiceConfig(nullptr, nullptr, &response);
  EXPECT_EQ(response.result(), ReloadConfigResponse::RELOAD_SUCCESS);

  response.Clear();

  magma_service.SetConfigReloadCallback(reload_failed);
  magma_service.ReloadServiceConfig(nullptr, nullptr, &response);
  EXPECT_EQ(response.result(), ReloadConfigResponse::RELOAD_FAILURE);

  response.Clear();

  magma_service.ClearConfigReloadCallback();
  magma_service.ReloadServiceConfig(nullptr, nullptr, &response);
  EXPECT_EQ(response.result(), ReloadConfigResponse::RELOAD_UNSUPPORTED);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace service303
}  // namespace magma
