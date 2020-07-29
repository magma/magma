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

#include <string>

#include <devmand/channels/Engine.h>

namespace devmand {
namespace channels {
namespace packet {

/*
 * This represents a packet core which processes raw packets seen on the wire.
 */
class Engine final : public channels::Engine {
 public:
  Engine(const std::string& interfaceName);

  Engine() = delete;
  ~Engine() override;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  void handleIncomingPacket();

 private:
  // A file descriptor for a raw socket.
  int fd{-1};
};

} // namespace packet
} // namespace channels
} // namespace devmand
