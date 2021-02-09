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

#include "EventBaseWrapper.h"
#include "magma_logging.h"
#include "Utilities.h"
#include "date.h"

namespace magma {
using namespace date;

EventBaseWrapper::EventBaseWrapper(folly::EventBase* evb) : evb_(evb) {}

void EventBaseWrapper::loopForever() {
  evb_->loopForever();
};

void EventBaseWrapper::terminateLoopSoon() {
  evb_->terminateLoopSoon();
};

void EventBaseWrapper::runAfterDelay(
    folly::Function<void()> func, int32_t delayMs) {
  auto insert = std::chrono::system_clock::now();
  log_push(insert);
  evb_->runAfterDelay(
      [this, func = std::move(func), insert]() mutable {
        auto start = std::chrono::system_clock::now();
        func();
        auto end = std::chrono::system_clock::now();
        log_pop(insert, start, end);
      },
      delayMs);
};

void EventBaseWrapper::runInEventBaseThread(folly::Cob&& cob) {
  auto insert = std::chrono::system_clock::now();
  log_push(insert);
  evb_->runInEventBaseThread([this, func = std::move(cob), insert]() mutable {
    auto start = std::chrono::system_clock::now();
    func();
    auto end = std::chrono::system_clock::now();
    log_pop(insert, start, end);
  });
}

void EventBaseWrapper::log_push(std::chrono::system_clock::time_point now) {
  event_count++;
  MLOG(MINFO) << "[" << event_count << "] ==> Inserting into EventQueue at "
              << now;
};

void EventBaseWrapper::log_pop(
    std::chrono::system_clock::time_point insert,
    std::chrono::system_clock::time_point start,
    std::chrono::system_clock::time_point end) {
  event_count--;
  MLOG(MINFO) << "[" << event_count << "] "
              << "<== Popping from EventQueue.";
  MLOG(MINFO) << "* inserted: " << insert;
  MLOG(MINFO) << "* started execution: " << start;
  MLOG(MINFO) << "* finished execution: " << end;
  MLOG(MINFO) << "* duration waiting: " << (start - insert);
  MLOG(MINFO) << "* duration executing: " << (end - start);
};
};  // namespace magma