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

#include <devmand/error/ErrorQueue.h>

namespace devmand {

void ErrorQueue::add(std::string&& error) {
  errors.withWLock([this, &error](auto& queue) {
    queue.emplace_back(std::forward<std::string>(error));
    // on max size, discard oldest error
    if (queue.size() > maxSize) {
      queue.pop_front();
    }
  });
}

folly::dynamic ErrorQueue::get() {
  auto ret = folly::dynamic::array();
  // NOTE: this is a shared read lock but if you modify "get()" to clear
  // the errors, you'll need a write lock.
  errors.withRLock([this, &ret](auto& queue) {
    for (auto& error : queue) {
      ret.push_back(error);
    }
  });
  return ret;
}

} // namespace devmand
