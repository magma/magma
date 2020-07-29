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

#pragma once

#define LOG_WITH_GLOG
#include <devmand/cartography/DeviceConfig.h>
#include <magma_logging.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

static const char* const ANY_VERSION = "*";

class DeviceType {
 private:
  string device;
  string version;

 public:
  DeviceType(const string& device, const string& version);

  static DeviceType getDefaultInstance() {
    return DeviceType("default", ANY_VERSION); // TODO
  }

  static DeviceType create(
      const devmand::cartography::DeviceConfig& deviceConfig);

  friend ostream& operator<<(ostream& os, const DeviceType& type);

  string str() const;

  bool operator==(const DeviceType& rhs) const;

  bool operator!=(const DeviceType& rhs) const;

  bool operator<(const DeviceType& rhs) const;

  bool operator>(const DeviceType& rhs) const;

  bool operator<=(const DeviceType& rhs) const;

  bool operator>=(const DeviceType& rhs) const;
};

} // namespace cli
} // namespace devices
} // namespace devmand
