/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/test/cli/utils/Log.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::devices::cli;
using namespace folly;

class PluginRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

TEST_F(PluginRegistryTest, api) {
  PluginRegistry reg;
  static const DeviceType type = {"IOS", "15.*"};

  typedef class : public Plugin {
   public:
    DeviceType getDeviceType() const override {
      return type;
    };

    void provideReaders(ReaderRegistryBuilder& reg) const override {
      (void)reg;
    };
    void provideWriters(WriterRegistryBuilder& reg) const override {
      (void)reg;
    };
  } TestPlugin;

  reg.registerPlugin(make_shared<TestPlugin>());
  MLOG(MDEBUG) << reg;

  EXPECT_THROW(
      reg.getDeviceContext({"NONEXISTING", "VERSION"}),
      PluginRegistryException);

  shared_ptr<DeviceContext> deviceCtx = reg.getDeviceContext(type);

  DeviceType reportedType = deviceCtx->getDeviceType();
  ASSERT_EQ(type, reportedType);
}

} // namespace cli
} // namespace test
} // namespace devmand
