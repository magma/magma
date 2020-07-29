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

#include <assert.h>

namespace devmand {
namespace utils {

template <class T>
class LifetimeTracker {
 public:
  LifetimeTracker() {
    ++allocations;
  }

  virtual ~LifetimeTracker() {
    ++deallocations;
  }

  LifetimeTracker(const LifetimeTracker&) = default;
  LifetimeTracker& operator=(const LifetimeTracker&) = default;
  LifetimeTracker(LifetimeTracker&&) = delete;
  LifetimeTracker& operator=(LifetimeTracker&&) = default;

  static unsigned int getAllocations() {
    return allocations;
  }

  static unsigned int getDeallocations() {
    return deallocations;
  }

  static unsigned int getLivingCount() {
    assert(allocations > deallocations);

    return allocations - deallocations;
  }

 private:
  static unsigned int allocations;
  static unsigned int deallocations;
};

template <class T>
unsigned int LifetimeTracker<T>::allocations;
template <class T>
unsigned int LifetimeTracker<T>::deallocations;

} // namespace utils
} // namespace devmand
