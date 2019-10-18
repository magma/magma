// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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

class Engine;

using Exception = std::runtime_error;

using OutstandingRequests = std::unordered_map<int, Request>;

class Channel final : public channels::Channel {
 public:
  Channel(
      Engine& engine_,
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
  folly::Future<Responses> walk(const channels::snmp::Oid& tree);

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

  folly::Future<Responses> walk(
      const channels::snmp::Oid& tree,
      const channels::snmp::Oid& current,
      Responses responses);

  Response processResponse(Request* request, snmp_pdu& response);
  Response processVars(netsnmp_variable_list* vars);
  static Oid firstOid(netsnmp_variable_list* vars);

 private:
  Engine& engine;
  const Peer peer;
  snmp_session* session{nullptr};
  OutstandingRequests outstandingRequests;
};

} // namespace snmp
} // namespace channels
} // namespace devmand
