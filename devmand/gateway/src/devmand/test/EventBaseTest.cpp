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

#include <devmand/test/EventBaseTest.h>

namespace devmand {
namespace test {

EventBaseTest::EventBaseTest() {
  start();
}

EventBaseTest::~EventBaseTest() {
  if (started) {
    stop();
  }
}
void EventBaseTest::start() {
  assert(not started);
  eventBaseThread =
      std::async(std::launch::async, [this] { eventBase.loopForever(); });
  started = true;
}

void EventBaseTest::stop() {
  assert(started);
  eventBase.terminateLoopSoon();
  eventBaseThread.wait();
  started = false;
}

} // namespace test
} // namespace devmand
