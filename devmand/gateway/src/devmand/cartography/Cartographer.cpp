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

#include <devmand/cartography/Cartographer.h>

namespace devmand {
namespace cartography {

Cartographer::Cartographer(
    const AddHandler& addHandler,
    const DeleteHandler& deleteHandler)
    : add(addHandler), del(deleteHandler) {
  assert(add != nullptr);
  assert(del != nullptr);
}

void Cartographer::addDeviceDiscoveryMethod(
    const std::shared_ptr<Method>& method) {
  assert(method != nullptr);
  auto result = methods.emplace(method);
  if (result.second) {
    method->setHandlers(add, del);
    method->enable();
  } else {
    LOG(ERROR) << "Failed to add device discovery method";
    throw std::runtime_error("Failed to add device discovery method.");
  }
}

} // namespace cartography
} // namespace devmand
