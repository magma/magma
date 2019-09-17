// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <string>
#include <unordered_map>

#include <folly/dynamic.h>
#include <folly/futures/Future.h>

#include <devmand/channels/Channel.h>
#include <devmand/channels/snmp/Pdu.h>
#include <devmand/channels/snmp/Request.h>
#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

using Exception = std::runtime_error;

using OutstandingRequests = std::unordered_map<int, Request>;

class Channel final : public channels::Channel {
 public:
  Channel(
      const Peer& peer_,
      const Community& community,
      const Version& version,
      const std::string& passphrase = "",
      const std::string& securityName = "",
      const SecurityLevel& securityLevel = "",
      oid proto[] = {});

  Channel() = delete;
  ~Channel() override;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  folly::Future<Response> asyncGet(const Oid& _oid);
  folly::Future<Response> asyncGetNext(const Oid& _oid);

  static folly::Future<std::string> asFutureString(
      folly::Future<channels::snmp::Response>&& future);
  static std::string toStatus(const std::string& v);
  static folly::Future<int> asFutureInt(
      folly::Future<channels::snmp::Response>&& future);

 private:
  void handleAsyncResponse(
      Request* request,
      int operation,
      int requestId,
      snmp_pdu* response);

  folly::Future<Response> asyncSend(
      Pdu& pdu,
      OutstandingRequests::iterator& result);

  Response processResponse(snmp_pdu& response);
  Response processVars(netsnmp_variable_list* vars);
  static Oid firstOid(netsnmp_variable_list* vars);

 private:
  const Peer peer;
  snmp_session* session{nullptr};
  OutstandingRequests outstandingRequests;
};

} // namespace snmp
} // namespace channels
} // namespace devmand
