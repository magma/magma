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

#include <list>
#include <memory>
#include <string>

#include <folly/io/async/EventHandler.h>

#include <devmand/channels/Engine.h>
#include <devmand/channels/snmp/EventHandler.h>

namespace devmand {
namespace channels {
namespace snmp {

class Engine final : public channels::Engine {
 public:
  Engine(folly::EventBase& eventBase_, const std::string& appName);

  Engine() = delete;
  ~Engine() override = default;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  folly::EventBase& getEventBase();
  void sync();

 private:
  void timeout();

  void enableDebug();

 private:
  folly::EventBase& eventBase;
  std::list<std::unique_ptr<EventHandler>> handlers;
};

} // namespace snmp
} // namespace channels
} // namespace devmand
