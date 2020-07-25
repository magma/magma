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

#include <devmand/channels/Channel.h>
#include <devmand/channels/ping/Engine.h>
#include <gtest/gtest_prod.h>
namespace devmand {
namespace test {
class PingChannelTest_checkSequenceIdGeneration_Test;
}
} // namespace devmand

namespace devmand {
namespace channels {
namespace ping {

class Channel : public channels::Channel {
 public:
  Channel(Engine& engine, folly::IPAddress target_);
  Channel() = delete;
  ~Channel() override = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  folly::Future<Rtt> ping();

 private:
  friend devmand::test::PingChannelTest_checkSequenceIdGeneration_Test;
  // this was not working?
  // FRIEND_TEST(PingChannelTest, checkSequenceIdGeneration);
  RequestId getSequence();
  icmphdr makeIcmpPacket();
  RequestId genRandomRequestId();

 private:
  Engine& engine;
  folly::IPAddress target;
  RequestId sequence;
};

} // namespace ping
} // namespace channels
} // namespace devmand
