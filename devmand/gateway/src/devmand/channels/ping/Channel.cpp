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

namespace devmand {
namespace channels {
namespace ping {

Channel::Channel(folly::EventBase& _eventBase, folly::IPAddress target_)
    : eventBase(_eventBase), target(target_) {
  icmpSocket = ::socket(AF_INET, SOCK_DGRAM, IPPROTO_ICMP);
  if (icmpSocket < 0) {
    throw std::system_error(errno, std::generic_category());
  }

  /*
  if (fcntl(icmpSocket, F_SETFL, O_NONBLOCK) < 0) {
    throw std::system_error(errno, std::generic_category());
  }
  */
}

// TODO first pass sync
folly::Future<Rtt> Channel::ping() {
  auto hdr = makeIcmpPacket();

  LOG(INFO) << "Sending ping to " << target.str() << " with sequence "
            << hdr.un.echo.sequence;

  // TODO BOOTCAMP this handles ipv4 only we should support ipv6 as well.
  sockaddr_in destination;
  if (inet_pton(AF_INET, target.str().c_str(), &destination.sin_addr) != 1) {
    LOG(ERROR) << "Invalid IPv4 Address " << target.str();
    return folly::makeFuture<Rtt>(0);
  }
  destination.sin_family = AF_INET;

  auto result = sendto(
      icmpSocket,
      &hdr,
      sizeof(hdr),
      0,
      reinterpret_cast<sockaddr*>(&destination),
      sizeof(destination));
  if (result <= 0) {
    switch (result) {
      case EAGAIN:
        // case EWOULDBLOCK:
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
    return folly::makeFuture<Rtt>(0);
  }

  sockaddr_in retAddr;
  unsigned int addrLen = static_cast<unsigned int>(sizeof(retAddr));
  if (recvfrom(
          icmpSocket,
          &hdr,
          sizeof(hdr),
          0,
          reinterpret_cast<sockaddr*>(&retAddr),
          &addrLen) <= 0) {
    LOG(ERROR) << "Packet receive failed!";
    return folly::makeFuture<Rtt>(0);
  } else {
    if (not(hdr.type == 0 and hdr.code == 0)) {
      LOG(ERROR) << "Packet received with ICMP type "
                 << static_cast<int>(hdr.type) << " code "
                 << static_cast<int>(hdr.code);
    } else {
      LOG(INFO) << "got!";
    }
  }

  // TODO setup promise

  return folly::makeFuture<Rtt>(0);
}

uint16_t Channel::getSequence() {
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

} // namespace ping
} // namespace channels
} // namespace devmand
