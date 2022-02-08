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
#include <event2/event.h>
#include <glog/logging.h>
#include <gtest/gtest.h>
#include <condition_variable>
#include <event2/thread.h>
#include "lte/gateway/c/session_manager/EventManager.h"

namespace magma {

using namespace std::chrono_literals;

TEST(EventManagerTest, test_construction) {
  evthread_use_pthreads();
  struct event_base* base = event_base_new();
  assert(base);
  std::cout << "TEST: Starting em..." << std::endl;
  EventManager em(base);
  sleep(1);
  std::cout << "TEST: Terminating em..." << std::endl;
  em.Terminate();
  std::cout << "TEST: Should have terminated em..." << std::endl;
  event_base_free(base);
}

void it_happened(evutil_socket_t sig, short events, void* user_data) {
  std::cout << "TEST: WUT! it_happened!!!" << std::endl;
  std::condition_variable* happened =
      static_cast<std::condition_variable*>(user_data);
  happened->notify_one();
}

TEST(EventManagerTest, test_immediate_event) {
  std::cout << "TEST: Starting test!" << std::endl;
  evthread_use_pthreads();
  struct event_base* base = event_base_new();
  assert(base);
  EventManager em(base);

  // Create conditional for notify and then capture in a lambda.
  std::mutex happened_mutex;
  std::condition_variable happened;
  em.AddEvent(&it_happened, (void*)&happened);

  std::cout << "TEST: WOOHOO A" << std::endl;

  // Pass into event request system.
  std::unique_lock<std::mutex> lk(happened_mutex);
  happened.wait_for(lk, 500 * 100ms);

  std::cout << "TEST: WOOHOO B" << std::endl;

  em.Terminate();

  std::cout << "TEST: WOOHOO C" << std::endl;
  event_base_free(base);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
