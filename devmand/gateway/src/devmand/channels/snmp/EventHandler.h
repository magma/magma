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

#include <folly/io/async/EventHandler.h>

namespace devmand {
namespace channels {
namespace snmp {

class Engine;

class EventHandler final : public folly::EventHandler {
 public:
  EventHandler(Engine& engine_, int fd_);
  EventHandler() = delete;
  ~EventHandler() override;
  EventHandler(const EventHandler&) = delete;
  EventHandler& operator=(const EventHandler&) = delete;
  EventHandler(EventHandler&& e) = delete;
  EventHandler& operator=(EventHandler&&) = delete;

 public:
  int getFd() const;

 private:
  void handlerReady(uint16_t events) noexcept override;
  void read();

 private:
  Engine& engine;
  int fd{-1};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
