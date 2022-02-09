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
#include <memory>
#include <event2/event.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <condition_variable>
#include <event2/thread.h>
#include "lte/gateway/c/session_manager/EventBaseManager.h"
#include "lte/gateway/c/session_manager/SchedulableCallback.h"

namespace magma {

using namespace std::chrono_literals;

class EventBaseManagerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    evthread_use_pthreads();
    base = EventBasePtr(event_base_new());
    assert(base);
    em = std::make_unique<EventBaseManager>(base);
  }

  virtual void TearDown() { em->Terminate(); }

 protected:
  EventBasePtr base;
  std::unique_ptr<EventBaseManager> em;
};

TEST_F(EventBaseManagerTest, test_construction) {
  // Test the basic SetUp and TearDown functionality that constructs and
  // destructs the EventBaseManager
  sleep(1);
}

TEST_F(EventBaseManagerTest, test_immediate_event) {
  std::mutex happened_mutex;
  std::condition_variable happened;
  std::unique_lock<std::mutex> lk(happened_mutex);

  // Schedule an immediate callback that notifies the condition variable
  // There *could* be a race condition here if notify_one is called before we
  // start waiting
  SchedulableCallback::MakeSchedulableCallback(base, [&]() {
    std::cout << "Hello from inside the callback!" << std::endl;
    happened.notify_one();
  })->scheduleCallbackSoon();

  // Assert we saw the callback. TODO: use enum
  EXPECT_EQ(std::cv_status::no_timeout, happened.wait_for(lk, 500 * 100ms));
}

TEST_F(EventBaseManagerTest, test_delayed_event_timeout) {
  std::mutex happened_mutex;
  std::condition_variable happened;
  std::unique_lock<std::mutex> lk(happened_mutex);

  SchedulableCallback::MakeSchedulableCallback(base, [&]() {
    std::cout << "Hello from inside the callback!" << std::endl;
    happened.notify_one();
  })->scheduleCallbackWithDelay(2);
  // Assert that it timed out TODO: use enum
  EXPECT_EQ(std::cv_status::timeout, happened.wait_for(lk, 100ms));
  // sleep so that we reach the callback eventually
  sleep(2);
}

TEST_F(EventBaseManagerTest, test_delayed_event_success) {
  std::mutex happened_mutex;
  std::condition_variable happened;
  std::unique_lock<std::mutex> lk(happened_mutex);

  SchedulableCallback::MakeSchedulableCallback(base, [&]() {
    std::cout << "Hello from inside the callback!" << std::endl;
    happened.notify_one();
  })->scheduleCallbackWithDelay(1);

  // wait up to 4 seconds should be plenty
  // Assert that the callback was called TODO: use enum
  EXPECT_EQ(std::cv_status::no_timeout, happened.wait_for(lk, 40 * 100ms));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
