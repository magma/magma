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

#include <folly/json.h>
#include <gtest/gtest.h>

#include <devmand/devices/demo/Device.h>

namespace devmand {
namespace test {

class DemoDeviceTest : public ::testing::Test {
 public:
  DemoDeviceTest() = default;
  virtual ~DemoDeviceTest() = default;
  DemoDeviceTest(const DemoDeviceTest&) = delete;
  DemoDeviceTest& operator=(const DemoDeviceTest&) = delete;
  DemoDeviceTest(DemoDeviceTest&&) = delete;
  DemoDeviceTest& operator=(DemoDeviceTest&&) = delete;
};

TEST_F(DemoDeviceTest, jsonSample) {
  folly::dynamic data = devices::demo::Device::getDemoDatastore();
  std::cerr << folly::toJson(data) << std::endl;
}

} // namespace test
} // namespace devmand
