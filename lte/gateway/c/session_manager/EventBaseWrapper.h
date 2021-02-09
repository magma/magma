/**
 * Copyright 2020 The Magma Authors.
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
#pragma once

#include <folly/io/async/EventBaseManager.h>

namespace magma {

class EventBaseWrapper {
 public:
  EventBaseWrapper();

  EventBaseWrapper(folly::EventBase* evb);

  void loopForever();

  void terminateLoopSoon();

  void runAfterDelay(folly::Function<void()> func, int32_t delayMs);

  void runInEventBaseThread(folly::Cob&& cob);

 private:
  folly::EventBase* evb_;
  int event_count = 0;
};

}  // namespace magma
