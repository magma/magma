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

#include <arpa/inet.h>
#include <folly/GLog.h>
#include <random>

#include <devmand/channels/ping/Channel.h>

namespace devmand {
namespace channels {
namespace ping {

Channel::Channel(Engine& engine_, folly::IPAddress target_)
    : engine(engine_), target(target_), sequence(genRandomRequestId()) {}

folly::Future<Rtt> Channel::ping() {
  auto pkt = IcmpPacket(target, getSequence());
  LOG(INFO) << "Sending ping to " << target.str() << " with sequence "
            << pkt.getSequence();
  return engine.ping(pkt);
}

RequestId Channel::genRandomRequestId() {
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<uint16_t> dis(0, UINT16_MAX);
  return dis(gen);
}

RequestId Channel::getSequence() {
  return ++sequence;
}

} // namespace ping
} // namespace channels
} // namespace devmand
