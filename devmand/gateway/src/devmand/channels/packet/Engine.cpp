// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <arpa/inet.h>
#include <linux/if_packet.h>
#include <net/ethernet.h>
#include <net/if.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>
#include <iomanip>
#include <sstream>
#include <system_error>

#include <iostream>
#include <stdexcept>

#include <folly/GLog.h>

#include <devmand/channels/packet/Engine.h>

namespace devmand {
namespace channels {
namespace packet {

Engine::Engine(const std::string& interfaceName) : channels::Engine("Packet") {
  LOG(INFO) << "Listening on interface " << interfaceName;

  fd = ::socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
  if (fd < 0) {
    throw std::system_error(errno, std::generic_category());
  }

  if (interfaceName.length() >= IFNAMSIZ) {
    throw std::range_error("interface name too long");
  }

  struct sockaddr_ll sock_address {};

  sock_address.sll_family = AF_PACKET;
  sock_address.sll_protocol = htons(ETH_P_ALL);
  sock_address.sll_ifindex =
      static_cast<int>(if_nametoindex(interfaceName.c_str()));
  if (bind(
          fd,
          reinterpret_cast<struct sockaddr*>(&sock_address),
          sizeof(sock_address)) < 0) {
    throw std::system_error(errno, std::generic_category());
  }
}

Engine::~Engine() {
  if (fd > 0 and ::close(fd) != 0) {
    LOG(ERROR) << "Failed to close " << __PRETTY_FUNCTION__;
  }
}

void Engine::handleIncomingPacket() {
  static constexpr unsigned int mtu{65536};
  unsigned char buffer[mtu] = {};

  ssize_t msgLength = ::recvfrom(fd, buffer, mtu, 0, nullptr, nullptr);
  if (msgLength < 0) {
    throw std::system_error(errno, std::generic_category());
  }

  std::stringstream ss;
  ss << std::hex << std::setw(2) << std::setfill('0');

  for (int i = 0; i < msgLength; ++i) {
    ss << static_cast<unsigned int>(buffer[i]);
  }

  LOG(INFO) << "received: " << ss.str();
}

} // namespace packet
} // namespace channels
} // namespace devmand
