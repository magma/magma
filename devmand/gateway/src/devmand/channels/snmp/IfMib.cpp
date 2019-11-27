// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/IfMib.h>

namespace devmand {
namespace channels {
namespace snmp {

folly::Future<std::string> IfMib::getSystemName(
    channels::snmp::Channel& channel) {
  return Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.1.5.0")));
}

folly::Future<Contact> IfMib::getSystemContact(
    channels::snmp::Channel& channel) {
  return Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.1.4.0")));
}

folly::Future<Location> IfMib::getSystemLocation(
    channels::snmp::Channel& channel) {
  return Channel::asFutureString(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.1.6.0")));
}

folly::Future<int> IfMib::getNumberOfInterfaces(
    channels::snmp::Channel& channel) {
  return Channel::asFutureInt(
      channel.asyncGet(channels::snmp::Oid(".1.3.6.1.2.1.2.1.0")));
}

folly::Future<InterfaceIndicies> IfMib::getInterfaceIndicies(
    channels::snmp::Channel& channel) {
  return channel.walk(channels::snmp::Oid(".1.3.6.1.2.1.2.2.1.1"))
      .thenValue([](auto responses) {
        InterfaceIndicies indices;
        for (auto res : responses) {
          indices.push_back((int)res.value.asInt());
        }
        return indices;
      });
}

folly::Future<InterfacePairs> IfMib::getInterfaceField(
    channels::snmp::Channel& channel,
    const std::string& oid,
    const std::function<std::string(std::string)>& formatter) {
  return getInterfaceIndicies(channel).thenValue(
      [&channel, oid, formatter](auto indicies) {
        std::vector<folly::Future<InterfacePair>> pairs;
        for (auto index : indicies) {
          pairs.emplace_back(
              channel
                  .asyncGet(
                      channels::snmp::Oid(oid + folly::to<std::string>(index)))
                  .thenValue([index, formatter](auto result) {
                    auto val = result.value.asString();
                    return InterfacePair{
                        index, formatter == nullptr ? val : formatter(val)};
                  }));
        }
        return folly::collect(std::move(pairs));
      });
}

folly::Future<InterfacePairs> IfMib::getInterfaceField(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices,
    const std::string& oid,
    const std::function<std::string(std::string)>& formatter) {
  std::vector<folly::Future<InterfacePair>> pairs;
  for (auto index : indices) {
    pairs.emplace_back(
        channel
            .asyncGet(channels::snmp::Oid(oid + folly::to<std::string>(index)))
            .thenValue([index, formatter](auto result) {
              auto val = result.value.asString();
              return InterfacePair{index,
                                   formatter == nullptr ? val : formatter(val)};
            }));
  }
  return folly::collect(std::move(pairs));
}

folly::Future<InterfacePairs> IfMib::getInterfaceOperStatuses(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(
      channel, indices, ".1.3.6.1.2.1.2.2.1.8.", Channel::toStatus);
}

folly::Future<InterfacePairs> IfMib::getInterfaceAdminStatuses(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(
      channel, indices, ".1.3.6.1.2.1.2.2.1.7.", Channel::toStatus);
}

folly::Future<InterfacePairs> IfMib::getInterfaceNames(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(channel, indices, ".1.3.6.1.2.1.2.2.1.2.");
}

folly::Future<InterfacePairs> IfMib::getInterfaceMtus(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(channel, indices, ".1.3.6.1.2.1.2.2.1.4.");
}

folly::Future<InterfacePairs> IfMib::getInterfaceTypes(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(channel, indices, ".1.3.6.1.2.1.2.2.1.3.");
}

folly::Future<InterfacePairs> IfMib::getInterfaceDescriptions(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(channel, indices, ".1.3.6.1.2.1.2.2.1.2.");
}

folly::Future<InterfacePairs> IfMib::getInterfaceLastChange(
    channels::snmp::Channel& channel,
    const InterfaceIndicies& indices) {
  return getInterfaceField(channel, indices, ".1.3.6.1.2.1.2.2.1.9.");
}

} // namespace snmp
} // namespace channels
} // namespace devmand
