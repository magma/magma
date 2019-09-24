// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/Channel.h>
#include <devmand/devices/mikrotik/Mib.h>

namespace devmand {
namespace devices {
namespace mikrotik {

folly::Future<std::string> Mib::getBaseMac(channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.17.1.1.0")));
}

folly::Future<std::string> Mib::getSerialNumber(
    channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.7.3.0")));
}

folly::Future<std::string> Mib::getFirmwareVersion(
    channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.7.4.0")));
}

folly::Future<std::string> Mib::getModel(channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
             channel.asyncGet(
                 channels::snmp::Oid(".1.3.6.1.2.1.47.1.1.1.1.2.65536")))
      .thenValue([](auto v) {
        static const std::string delim = ") on ";
        auto pos = v.find(delim);
        if (pos != std::string::npos) {
          return v.erase(0, pos + delim.length());
        } else {
          return std::string("unable to parse model");
        }
      });
}

folly::Future<std::string> Mib::getUpTime(channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.1.3.0")));
}

folly::Future<std::string> Mib::getLongtitude(
    channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.12.2")));
}

folly::Future<std::string> Mib::getLatitude(channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.12.3")));
}

folly::Future<std::string> Mib::getAltitude(channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.12.4")));
}

folly::Future<std::string> Mib::getIpv4Address(
    channels::snmp::Channel& channel) {
  return channels::snmp::Channel::asFutureString(
      channel.asyncGetNext(channels::snmp::Oid(".1.3.6.1.2.1.4.20.1.1")));
}

folly::Future<std::string> Mib::getIpv6Address(
    channels::snmp::Channel& channel) {
  // TODO we need to confirm this device has an ipv6 address and error if not
  // e.g. if getnext returns something other than the requested mib.
  return channels::snmp::Channel::asFutureString(
      channel.asyncGetNext(channels::snmp::Oid(".1.3.6.1.2.1.55.1.8.1.1")));
}

} // namespace mikrotik
} // namespace devices
} // namespace devmand
