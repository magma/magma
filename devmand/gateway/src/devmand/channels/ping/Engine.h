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

#include <map>
#include <string>

#include <netinet/icmp6.h>
#include <netinet/ip_icmp.h>

#include <folly/futures/Future.h>
#include <folly/io/async/AsyncSocket.h>

#include <devmand/channels/Engine.h>
#include <devmand/utils/EventBaseUtils.h>
#include <devmand/utils/Time.h>

namespace devmand {
namespace channels {
namespace ping {

extern std::out_of_range DefaultSwitchError;

using Rtt = uint64_t;

struct Request {
  utils::TimePoint start;
  folly::Promise<Rtt> promise;
};

enum IPVersion { v4, v6 };
enum PacketType { send, read };

using RequestId = uint16_t;

using OutstandingRequests =
    std::map<std::pair<folly::IPAddress, RequestId>, Request>;

// The data structure for icmp headers.
// Generalized for v4 or v6.
class IcmpPacket final {
 public:
  IcmpPacket(IPVersion ipv); // read ping
  IcmpPacket(const folly::IPAddress& addr, RequestId sequence); // send ping
  IcmpPacket() = delete;
  ~IcmpPacket() = default;
  IcmpPacket(const IcmpPacket&) = delete;
  IcmpPacket& operator=(const IcmpPacket&) = delete;
  IcmpPacket(IcmpPacket&&) = default;
  IcmpPacket& operator=(IcmpPacket&&) = default;

 public:
  // getters for private members
  RequestId getSequence() const;
  const folly::IPAddress& getAddr() const;
  bool wasSuccess();
  bool isEchoReply();
  auto getType();
  auto getCode();
  const sockaddr_storage& getSrc();

 public:
  auto send(int socket) const;
  void read(int socket);

 private:
  PacketType packetType;
  IPVersion ipv;
  bool success{false};
  folly::IPAddress addr;
  icmphdr hdrV4{};
  icmp6_hdr hdrV6{};
  sockaddr_storage src{};
  socklen_t srcLen{sizeof(sockaddr_storage)};
};

class Engine : public channels::Engine, public folly::EventHandler {
 public:
  Engine(
      folly::EventBase& _eventBase,
      IPVersion ipv_,
      const std::chrono::milliseconds& pingTimeout_ =
          std::chrono::milliseconds(5000),
      const std::chrono::milliseconds& timeoutFrequency_ =
          std::chrono::milliseconds(10000));
  Engine(
      folly::EventBase& _eventBase,
      const std::chrono::milliseconds& pingTimeout_ =
          std::chrono::milliseconds(5000),
      const std::chrono::milliseconds& timeoutFrequency_ =
          std::chrono::milliseconds(10000));
  Engine() = delete;
  ~Engine() override;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  folly::Future<Rtt> ping(const IcmpPacket& pkt);

  // NOTE this must be called after the event base is running.
  void start();

 private:
  virtual void handlerReady(uint16_t events) noexcept override;
  void timeout();
  IcmpPacket read();

 private:
  folly::EventBase& eventBase;
  folly::Synchronized<OutstandingRequests> sharedOutstandingRequests;
  int icmpSocket{-1};
  IPVersion ipv;
  bool failedIpv6Socket{false};
  std::chrono::milliseconds pingTimeout;
  std::chrono::milliseconds timeoutFrequency;
};

} // namespace ping
} // namespace channels
} // namespace devmand
