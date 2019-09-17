// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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

folly::Future<InterfaceIndicies> IfMib::handleNextInterfaceIndex(
    channels::snmp::Channel& channel,
    int numInterfacesRemaining,
    InterfaceIndicies indicies,
    channels::snmp::Oid marker) {
  // TODO I need to look into cancellable futures in the case of the device
  // lifetime ending.
  if (numInterfacesRemaining != 0) {
    return channel.asyncGetNext(marker).thenValue(
        [&channel, numInterfacesRemaining, indicies = std::move(indicies)](
            auto response) mutable {
          indicies.emplace_back(response.value.asInt());
          return IfMib::handleNextInterfaceIndex(
              channel,
              numInterfacesRemaining - 1,
              std::move(indicies),
              response.oid);
        });
  }
  return folly::makeFuture<InterfaceIndicies>(std::move(indicies));
}

folly::Future<InterfaceIndicies> IfMib::getInterfaceIndicies(
    channels::snmp::Channel& channel) {
  return getNumberOfInterfaces(channel).thenValue(
      [&channel](auto numberOfInterfaces) {
        channels::snmp::Oid marker(".1.3.6.1.2.1.2.1.0");
        InterfaceIndicies indicies;
        return IfMib::handleNextInterfaceIndex(
            channel, numberOfInterfaces, std::move(indicies), marker);
      });
}

folly::Future<InterfaceStatuses> IfMib::getInterfaceStatuses(
    channels::snmp::Channel& channel) {
  return getInterfaceIndicies(channel).thenValue([&channel](auto indicies) {
    std::vector<folly::Future<InterfaceStatus>> statuses;
    for (auto index : indicies) {
      std::string o{".1.3.6.1.2.1.2.2.1.8."};
      o += folly::to<std::string>(index);
      statuses.emplace_back(
          channel.asyncGet(channels::snmp::Oid(o))
              .thenValue([index](auto result) {
                return InterfaceStatus{
                    index, Channel::toStatus(result.value.asString())};
              }));
    }
    return folly::collect(std::move(statuses));
  });
}

folly::Future<InterfaceNames> IfMib::getInterfaceNames(
    channels::snmp::Channel& channel) {
  return getInterfaceIndicies(channel).thenValue([&channel](auto indicies) {
    std::vector<folly::Future<InterfaceName>> statuses;
    for (auto index : indicies) {
      std::string o{".1.3.6.1.2.1.2.2.1.2."};
      o += folly::to<std::string>(index);
      statuses.emplace_back(
          channel.asyncGet(channels::snmp::Oid(o))
              .thenValue([index](auto result) {
                return InterfaceName{index, result.value.asString()};
              }));
    }
    return folly::collect(std::move(statuses));
  });
}

} // namespace snmp
} // namespace channels
} // namespace devmand
