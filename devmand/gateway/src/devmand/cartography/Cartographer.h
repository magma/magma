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

#include <memory>
#include <set>

#include <folly/io/async/EventBase.h>

#include <devmand/cartography/Method.h>

namespace devmand {

namespace cartography {

/*
 * This class is responsible for mapping devices on a network which is does
 * through a set of methods.
 */
class Cartographer {
 public:
  Cartographer(
      const AddHandler& addHandler,
      const DeleteHandler& deleteHandler);
  Cartographer() = delete;
  virtual ~Cartographer() = default;
  Cartographer(const Cartographer&) = delete;
  Cartographer& operator=(const Cartographer&) = delete;
  Cartographer(Cartographer&&) = delete;
  Cartographer& operator=(Cartographer&&) = delete;

 public:
  void addDeviceDiscoveryMethod(const std::shared_ptr<Method>& method);

 private:
  std::set<std::shared_ptr<Method>> methods;
  const AddHandler add;
  const DeleteHandler del;
};

} // namespace cartography
} // namespace devmand
