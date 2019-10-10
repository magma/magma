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

#include <folly/GLog.h>

#include <devmand/channels/ping/Channel.h>
#include <devmand/utils/Time.h>

namespace devmand {
namespace channels {
namespace ping {

Channel::Channel(folly::EventBase& _eventBase, folly::IPAddress target_)
    : folly::EventHandler(&_eventBase), eventBase(_eventBase), target(target_) {
  // TODO this socket needs to be moved to an engine...
  icmpSocket = ::socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP);
  if (icmpSocket < 0) {
    throw std::system_error(errno, std::generic_category());
  }

  if (fcntl(icmpSocket, F_SETFL, O_NONBLOCK) < 0) {
    throw std::system_error(errno, std::generic_category());
  }

  folly::EventHandler::changeHandlerFD(
      folly::NetworkSocket::fromFd(icmpSocket));

  registerHandler(folly::EventHandler::READ | folly::EventHandler::PERSIST);
}

folly::Future<Rtt> Channel::ping() {
  icmphdr hdr = makeIcmpPacket();

  LOG(INFO) << "Sending ping to " << target.str() << " with sequence "
            << hdr.un.echo.sequence;

  // TODO BOOTCAMP this handles ipv4 only we should support ipv6 as well.
  sockaddr_in destination;
  if (inet_pton(AF_INET, target.str().c_str(), &destination.sin_addr) != 1) {
    LOG(ERROR) << "Invalid IPv4 Address " << target.str();
    return folly::makeFuture<Rtt>(0);
  }
  destination.sin_family = AF_INET;

  auto request = outstandingRequests.emplace(
      std::piecewise_construct,
      std::forward_as_tuple(hdr.un.echo.sequence),
      std::forward_as_tuple(Request{}));
  if (request.second) {
    request.first->second.start = utils::Time::now();
    auto result = sendto(
        icmpSocket,
        &hdr,
        sizeof(hdr),
        0,
        reinterpret_cast<sockaddr*>(&destination),
        sizeof(destination));
    if (result <= 0) {
      switch (result) {
        case EAGAIN: // case EWOULDBLOCK:
          // TODO if the ping fail because of a kernel buffer I'm not going to
          // implement retry logic as something is filling up the buffers. We
          // should probably alarm if this is the case.
          LOG(ERROR) << "Buffer full so ping failed to " << target.str();
          break;
        default:
          // TODO BOOTCAMP get errno string from syserror
          LOG(ERROR) << "Failed to send packet with errno " << errno;
          break;
      }
      outstandingRequests.erase(request.first);
      return folly::makeFuture<Rtt>(0);
    } else {
      return request.first->second.promise.getFuture();
    }
  } else {
    LOG(ERROR) << "ICMP Echo Id rollover with outstanding requests";
    return folly::makeFuture<Rtt>(0);
  }
  // TODO implement timeout
}

RequestId Channel::getSequence() {
  return ++sequence;
}

icmphdr Channel::makeIcmpPacket() {
  icmphdr hdr{};
  hdr.type = ICMP_ECHO;
  // hdr.un.echo.id = 0;
  hdr.un.echo.sequence = getSequence();
  // hdr.checksum = 0; // Let the kernel fill in the checksum
  return hdr;
}

void Channel::handlerReady(uint16_t) noexcept {
  // TODO end time isn't really precise here as we don't have a kernel time
  // need to implement kernel timestamping
  utils::TimePoint end = utils::Time::now();
  icmphdr hdr;
  while (recv(icmpSocket, &hdr, sizeof(hdr), 0) > 0) {
    LOG(INFO) << "Packet received with ICMP type " << static_cast<int>(hdr.type)
              << " code " << static_cast<int>(hdr.code);

    if (hdr.type == 0 and hdr.code == 0) {
      auto request = outstandingRequests.find(hdr.un.echo.sequence);
      if (request != outstandingRequests.end()) {
        auto duration = std::chrono::duration_cast<std::chrono::microseconds>(
            end - request->second.start);
        LOG(INFO) << "Received ICMP response after " << duration.count()
                  << " microseconds";
        request->second.promise.setValue(duration.count());
      }
    }
  }
}

} // namespace ping
} // namespace channels
} // namespace devmand
