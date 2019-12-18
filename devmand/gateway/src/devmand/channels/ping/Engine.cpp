// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <arpa/inet.h>
#include <fcntl.h>
#include <netdb.h>
#include <netinet/in.h>
#include <netinet/ip_icmp.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>

#include <algorithm>

#include <folly/GLog.h>

#include <devmand/channels/ping/Engine.h>
#include <devmand/utils/Time.h>

namespace devmand {
namespace channels {
namespace ping {

// Exhaustive switch statements, even for enums, don't guarantee a return.
// Just throw this if there's an unexpected enuenum value.
std::out_of_range DefaultSwitchError(
    "Reached end of exhaustive switch IPVersion.");

IcmpPacket::IcmpPacket(IPVersion ipv_)
    : packetType(PacketType::read), ipv(ipv_) {}

IcmpPacket::IcmpPacket(const folly::IPAddress& addr_, RequestId sequence)
    : packetType(PacketType::send),
      ipv(addr_.isV4() ? IPVersion::v4 : IPVersion::v6),
      addr(addr_) {
  switch (ipv) {
    case IPVersion::v4:
      hdrV4.type = ICMP_ECHO;
      hdrV4.un.echo.sequence = sequence;
      break;
    case IPVersion::v6:
      hdrV6.icmp6_type = ICMP6_ECHO_REQUEST;
      hdrV6.icmp6_seq = sequence;
      break;
    default:
      throw DefaultSwitchError;
  }
}

RequestId IcmpPacket::getSequence() const {
  switch (ipv) {
    case IPVersion::v4:
      return hdrV4.un.echo.sequence;
    case IPVersion::v6:
      return hdrV6.icmp6_seq;
    default:
      throw DefaultSwitchError;
  }
}

const folly::IPAddress& IcmpPacket::getAddr() const {
  return addr;
}

bool IcmpPacket::wasSuccess() {
  return success;
}

auto IcmpPacket::getType() {
  switch (ipv) {
    case IPVersion::v4:
      return hdrV4.type;
    case IPVersion::v6:
      return hdrV6.icmp6_type;
    default:
      throw DefaultSwitchError;
  }
}

bool IcmpPacket::isEchoReply() {
  switch (ipv) {
    case IPVersion::v4:
      return getType() == ICMP_ECHOREPLY;
    case IPVersion::v6:
      return getType() == ICMP6_ECHO_REPLY;
    default:
      throw DefaultSwitchError;
  }
}

auto IcmpPacket::getCode() {
  switch (ipv) {
    case IPVersion::v4:
      return hdrV4.code;
    case IPVersion::v6:
      return hdrV6.icmp6_code;
    default:
      throw DefaultSwitchError;
  }
}

const sockaddr_storage& IcmpPacket::getSrc() {
  return src;
}

auto IcmpPacket::send(int socket) const {
  assert(packetType == PacketType::send);
  sockaddr_storage dst;
  addr.toSockaddrStorage(&dst);
  switch (ipv) {
    case IPVersion::v4:
      return sendto(
          socket,
          &hdrV4,
          sizeof(hdrV4),
          0,
          reinterpret_cast<const sockaddr*>(&dst),
          sizeof(dst));
    case IPVersion::v6:
      return sendto(
          socket,
          &hdrV6,
          sizeof(hdrV6),
          0,
          reinterpret_cast<const sockaddr*>(&dst),
          sizeof(dst));
    default:
      throw DefaultSwitchError;
  }
}

void IcmpPacket::read(int socket) {
  assert(packetType == PacketType::read);
  switch (ipv) {
    case IPVersion::v4:
      success = recvfrom(
                    socket,
                    &hdrV4,
                    sizeof(hdrV4),
                    0,
                    reinterpret_cast<sockaddr*>(&src),
                    &srcLen) > 0;
      break;
    case IPVersion::v6:
      success = recvfrom(
                    socket,
                    &hdrV6,
                    sizeof(hdrV6),
                    0,
                    reinterpret_cast<sockaddr*>(&src),
                    &srcLen) > 0;
      break;
    default:
      throw DefaultSwitchError;
  }
}

IcmpPacket Engine::read() {
  IcmpPacket pkt(ipv);
  pkt.read(icmpSocket);
  return pkt;
}

Engine::Engine(
    folly::EventBase& _eventBase,
    IPVersion ipv_,
    const std::chrono::milliseconds& pingTimeout_,
    const std::chrono::milliseconds& timeoutFrequency_)
    : channels::Engine("Ping"),
      folly::EventHandler(&_eventBase),
      eventBase(_eventBase),
      ipv(ipv_),
      pingTimeout(pingTimeout_),
      timeoutFrequency(timeoutFrequency_) {
  switch (ipv) {
    case IPVersion::v4:
      icmpSocket = ::socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP);
      break;
    case IPVersion::v6:
      icmpSocket = ::socket(AF_INET6, SOCK_DGRAM, IPPROTO_ICMPV6);
      break;
  }
  if (icmpSocket < 0) {
    auto err = std::system_error(errno, std::generic_category());
    LOG(ERROR) << "Failed to open dgram socket: " << err.what();
    if (ipv_ == IPVersion::v6) {
      LOG(ERROR) << "ICMP over IPv6 will not work.";
      failedIpv6Socket = true;
    } else {
      throw err;
    }
  } else {
    if (fcntl(icmpSocket, F_SETFL, O_NONBLOCK) < 0) {
      throw std::system_error(errno, std::generic_category());
    }
    folly::EventHandler::changeHandlerFD(
        folly::NetworkSocket::fromFd(icmpSocket));
  }
  registerHandler(folly::EventHandler::READ | folly::EventHandler::PERSIST);
  start();
}

Engine::Engine(
    folly::EventBase& _eventBase,
    const std::chrono::milliseconds& pingTimeout_,
    const std::chrono::milliseconds& timeoutFrequency_)
    : Engine::Engine(
          _eventBase,
          IPVersion::v4,
          pingTimeout_,
          timeoutFrequency_) {}

Engine::~Engine() {
  unregisterHandler();
}

void Engine::start() {
  // TODO I implement a very simple type of timeout here where ever n
  // milliseconds we walk the pending requests and timeout ones that have
  // exceeded their time. This is neither the most efficient (we walk the entire
  // list) or the most accurate (we only guarentee and eventual timeout not a
  // precise timeout). Given this usecase I dont think either of these are that
  // important but I'm making this note so we know in the future we could
  // improve this by using one of the nice timeout queues that exist.
  eventBase.runInEventBaseThread([this]() {
    EventBaseUtils::scheduleEvery(
        eventBase, [this]() { timeout(); }, timeoutFrequency);
  });
}

void Engine::timeout() {
  sharedOutstandingRequests.withULockPtr([this](auto uOutstandingRequests) {
    auto outstandingRequests = uOutstandingRequests.moveFromUpgradeToWrite();
    // LOG(INFO) << "Processing ping timeouts";
    for (auto it = outstandingRequests->begin();
         it != outstandingRequests->end();) {
      utils::TimePoint end = utils::Time::now();
      if ((end - it->second.start) > pingTimeout) {
        LOG(ERROR) << "Ping request timed out";
        it->second.promise.setValue(0);
        it = outstandingRequests->erase(it);
      } else {
        ++it;
      }
    }
  });
}

folly::Future<Rtt> Engine::ping(const IcmpPacket& pkt) {
  if (failedIpv6Socket) {
    LOG(ERROR) << "Attempted ping on IPV6 where socket failed to open.";
    return folly::makeFuture<Rtt>(0);
  }
  incrementRequests();
  return sharedOutstandingRequests.withULockPtr([this, &pkt](
                                                    auto uOutstandingRequests) {
    auto outstandingRequests = uOutstandingRequests.moveFromUpgradeToWrite();
    auto request = outstandingRequests->emplace(
        std::piecewise_construct,
        std::forward_as_tuple(std::make_pair(pkt.getAddr(), pkt.getSequence())),
        std::forward_as_tuple(Request{}));
    if (request.second) {
      request.first->second.start = utils::Time::now();
      auto result = pkt.send(icmpSocket);
      if (result <= 0) {
        switch (result) {
          case EAGAIN: // case EWOULDBLOCK:
            // TODO if the ping fail because of a kernel buffer I'm not going to
            // implement retry logic as something is filling up the buffers. We
            // should probably alarm if this is the case.
            LOG(ERROR) << "Buffer full so ping failed";
            break;
          default:
            // TODO BOOTCAMP get errno string from syserror
            LOG(ERROR) << "Failed to send packet with errno " << errno;
            break;
        }
        outstandingRequests->erase(request.first);
        return folly::makeFuture<Rtt>(0);
      } else {
        auto ret = request.first->second.promise.getFuture();
        return ret;
      }
    } else {
      LOG(ERROR) << "ICMP Echo Id rollover with outstanding requests";
      return folly::makeFuture<Rtt>(0);
    }
  });
}

void Engine::handlerReady(uint16_t) noexcept {
  // TODO end time isn't really precise here as we don't have a kernel time
  // need to implement kernel timestamping
  utils::TimePoint end = utils::Time::now();
  for (IcmpPacket pkt = read(); pkt.wasSuccess(); pkt = read()) {
    bool processed{false};
    if (pkt.isEchoReply() and pkt.getCode() == 0) {
      sharedOutstandingRequests.withULockPtr([this, &end, &pkt, &processed](
                                                 auto uOutstandingRequests) {
        auto outstandingRequests =
            uOutstandingRequests.moveFromUpgradeToWrite();
        auto src = pkt.getSrc();
        auto request = outstandingRequests->find(std::make_pair(
            folly::IPAddress(reinterpret_cast<sockaddr*>(&src)),
            pkt.getSequence()));
        if (request != outstandingRequests->end()) {
          auto duration = std::chrono::duration_cast<std::chrono::microseconds>(
              end - request->second.start);
          LOG(INFO) << "Received ICMP response after " << duration.count()
                    << " microseconds";
          request->second.promise.setValue(duration.count());
          outstandingRequests->erase(request);
          processed = true;
        }
      });
    }

    if (not processed) {
      LOG(INFO) << "Packet received with ICMP type "
                << static_cast<int>(pkt.getType()) << " code "
                << static_cast<int>(pkt.getCode());
    }
  }
}

} // namespace ping
} // namespace channels
} // namespace devmand
