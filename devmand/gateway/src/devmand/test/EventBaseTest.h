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

#include <future>

#include <gtest/gtest.h>

#include <folly/io/async/EventBase.h>

namespace devmand {
namespace test {

class EventBaseTest : public ::testing::Test {
 public:
  EventBaseTest();
  virtual ~EventBaseTest();
  EventBaseTest(const EventBaseTest&) = delete;
  EventBaseTest& operator=(const EventBaseTest&) = delete;
  EventBaseTest(EventBaseTest&&) = delete;
  EventBaseTest& operator=(EventBaseTest&&) = delete;

 protected:
  void start();
  void stop();

 protected:
  folly::EventBase eventBase;

 private:
  bool started{false};
  std::future<void> eventBaseThread;
};

} // namespace test
} // namespace devmand
