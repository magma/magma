/**
 * Copyright 2026 The Magma Authors.
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

#include <cassert>
#include <cstdlib>
#include <cstring>

#include "lte/gateway/c/core/common/dynamic_memory_check.h"

namespace {

bool destructor_called = false;

class DestructionTracked {
 public:
  DestructionTracked() : payload_(new int[64]) {}

  ~DestructionTracked() {
    delete[] payload_;
    destructor_called = true;
  }

 private:
  int* payload_;
};

}  // namespace

int main() {
  void* allocation = std::malloc(256);
  assert(allocation != nullptr);
  std::memset(allocation, 0xA5, 256);

  free_wrapper(&allocation);
  assert(allocation == nullptr);

  // A pointer-to-pointer with a null pointee must be a safe no-op. Repeating
  // the call also verifies that free_wrapper leaves the caller's state null.
  free_wrapper(&allocation);
  assert(allocation == nullptr);

  // Typed deletion must invoke the concrete destructor and release the
  // object's owned allocation.
  DestructionTracked* tracked = new DestructionTracked();
  free_cpp_wrapper(&tracked);
  if (tracked != nullptr || !destructor_called) {
    return 1;
  }

  return 0;
}
