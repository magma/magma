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

#include <chrono>
#include <thread>

#define EXPECT_BECOMES_TRUE(exp)                                     \
  do {                                                               \
    constexpr std::chrono::seconds _maxExpectsWait{10};              \
    auto _expectsWait = _maxExpectsWait;                             \
    while ((not(exp)) and _expectsWait != std::chrono::seconds(0)) { \
      _expectsWait -= std::chrono::seconds(1);                       \
      std::this_thread::sleep_for(std::chrono::seconds(1));          \
    }                                                                \
    EXPECT_TRUE(exp);                                                \
  } while (false);
