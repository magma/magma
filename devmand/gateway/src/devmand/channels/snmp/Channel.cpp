// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <iostream>
#include <stdexcept>

#include <folly/IPAddress.h>
#include <folly/ScopeGuard.h>

#include <devmand/channels/snmp/Channel.h>
#include <devmand/channels/snmp/Engine.h>

namespace devmand {
namespace channels {
namespace snmp {

folly::Future<std::string> Channel::asFutureString(
    folly::Future<channels::snmp::Response>&& future) {
  return std::move(future).thenValue(
      [](auto result) { return result.value.asString(); });
}

std::string Channel::toStatus(const std::string& v) {
  if (v == "1") {
    return "UP";
  } else if (v == "2") {
    return "DOWN";
  } else {
    return v;
  }
}

folly::Future<int> Channel::asFutureInt(
    folly::Future<channels::snmp::Response>&& future) {
  return std::move(future).thenValue([](auto result) {
    return result.value.isInt() ? result.value.asInt() : 0;
  });
}

static inline int parseVersion(const Version& version) {
  if (version == "v1") {
    return SNMP_VERSION_1;
  } else if (version == "v2c") {
    return SNMP_VERSION_2c;
  } else if (version == "v3") {
    return SNMP_VERSION_3;
  } else {
    throw std::runtime_error("unknown version");
  }
}

static inline int parseSecurityLevel(const SecurityLevel& lvl) {
  if (lvl == "noAuth") {
    return SNMP_SEC_LEVEL_NOAUTH;
  } else if (lvl == "authNoPriv") {
    return SNMP_SEC_LEVEL_AUTHNOPRIV;
  } else if (lvl == "authPriv") {
    return SNMP_SEC_LEVEL_AUTHPRIV;
  } else {
    throw std::runtime_error("unknown version");
  }
}

Channel::Channel(
    Engine& engine_,
    const Peer& peer_,
    const Community& community,
    const Version& version,
    const std::string& passphrase,
    const std::string& securityName,
    const SecurityLevel& securityLevel,
    oid proto[])
    : engine(engine_), peer(peer_) {
  snmp_session sessionIn;

  snmp_sess_init(&sessionIn);
  sessionIn.version = parseVersion(version);

  if (sessionIn.version == SNMP_VERSION_3) {
    // TODO v3 is currently just a shell. Needs more work to really work
    // and handle all the config options.
    sessionIn.securityName = const_cast<char*>(securityName.data());
    sessionIn.securityNameLen = securityName.length();
    sessionIn.securityLevel = parseSecurityLevel(securityLevel);
    sessionIn.securityAuthProto = proto;
    // sessionIn.securityAuthProtoLen = sizeof(proto) / sizeof(oid);
    sessionIn.securityAuthKeyLen = USM_AUTH_KU_LEN;
    if (generate_Ku(
            sessionIn.securityAuthProto,
            static_cast<u_int>(sessionIn.securityAuthProtoLen),
            reinterpret_cast<u_char*>(const_cast<char*>(passphrase.data())),
            passphrase.length(),
            sessionIn.securityAuthKey,
            &sessionIn.securityAuthKeyLen) != SNMPERR_SUCCESS) {
      throw std::runtime_error("generate_Ku error");
    }
  } else {
    sessionIn.peername = const_cast<char*>(peer.data());
    sessionIn.community =
        reinterpret_cast<u_char*>(const_cast<char*>(community.c_str()));
    sessionIn.community_len = community.length();
  }

  SOCK_STARTUP;
  session = snmp_open(&sessionIn);
  if (session == nullptr) {
    snmp_perror("snmp_perror snmp_open");
    throw std::runtime_error("snmp_open error");
  }
}

Channel::~Channel() {
  if (session != nullptr and snmp_close(session) == 0) {
    snmp_perror("snmp_perror snmp_close");
  }
  SOCK_CLEANUP;
}

void Channel::handleAsyncResponse(
    Request* request,
    int operation,
    int requestId,
    snmp_pdu* response) {
  switch (operation) {
    case NETSNMP_CALLBACK_OP_RECEIVED_MESSAGE: {
      Response pr = processResponse(request, *response);
      if (not pr.isError()) {
        request->responsePromise.setValue(pr);
      } else {
        request->responsePromise.setException(Exception(pr.value.asString()));
      }
      break;
    }
    case NETSNMP_CALLBACK_OP_TIMED_OUT:
      request->responsePromise.setException(Exception(folly::sformat(
          "snmp timeout {} on oid {}", peer, request->oid.toString())));
      break;
    default:
      request->responsePromise.setException(Exception("snmp error"));
      break;
  }

  if (outstandingRequests.erase(requestId) != 1) {
    LOG(ERROR) << "Request " << requestId << " not found!";
  }
}

folly::Future<Response> Channel::asyncSend(
    Pdu& pdu,
    OutstandingRequests::iterator& result) {
  auto handler = [](int operation,
                    snmp_session* sess,
                    int requestId,
                    snmp_pdu* response,
                    void* magic) -> int {
    (void)sess;
    auto* request = reinterpret_cast<Request*>(magic);
    request->channel->handleAsyncResponse(
        request, operation, requestId, response);
    return 1;
  };

  engine.incrementRequests();
  if (snmp_async_send(session, pdu.get(), handler, &result->second)) {
    pdu.release();
    return result->second.responsePromise.getFuture();
  } else {
    outstandingRequests.erase(result);
    return folly::makeFuture<Response>(ErrorResponse("snmp_async_send error"));
  }
}

folly::Future<Response> Channel::asyncGet(const Oid& _oid) {
  Pdu pdu(SNMP_MSG_GET, _oid);
  auto result = outstandingRequests.emplace(
      std::piecewise_construct,
      std::forward_as_tuple(pdu.get()->reqid),
      std::forward_as_tuple(Request{this, _oid}));
  if (result.second) {
    return asyncSend(pdu, result.first);
  } else {
    throw std::runtime_error("emplace failed");
  }
}

folly::Future<Response> Channel::asyncGetNext(const Oid& _oid) {
  Pdu pdu(SNMP_MSG_GETNEXT, _oid);
  auto result = outstandingRequests.emplace(
      std::piecewise_construct,
      std::forward_as_tuple(pdu.get()->reqid),
      std::forward_as_tuple(Request{this, _oid}));
  if (result.second) {
    return asyncSend(pdu, result.first);
  } else {
    throw std::runtime_error("emplace failed");
  }
}

Oid Channel::firstOid(netsnmp_variable_list* vars) {
  for (; vars != nullptr; vars = vars->next_variable) {
    return Oid(vars->name, vars->name_length);
  }
  return Oid::error;
}

Response Channel::processVars(netsnmp_variable_list* vars) {
  // TODO handle multiple vars returned
  for (; vars != nullptr; vars = vars->next_variable) {
    Oid oid(vars->name, vars->name_length);
    switch (vars->type) {
      case SNMP_ENDOFMIBVIEW:
      case SNMP_NOSUCHOBJECT:
      case SNMP_NOSUCHINSTANCE:
        return Response(oid, nullptr);
      case ASN_OCTET_STR:
        return Response(
            oid,
            folly::dynamic(std::string(
                reinterpret_cast<char*>(vars->val.string), vars->val_len)));
      case ASN_IPADDRESS:
        return Response(
            oid,
            folly::IPAddress::fromBinary(
                folly::ByteRange(vars->val.string, vars->val_len))
                .str());
      case ASN_COUNTER64: {
        uint64_t val =
            (static_cast<uint64_t>((*vars->val.counter64).high) << 32) |
            vars->val.counter64->low;
        return Response(oid, val);
      }
      case ASN_INTEGER:
      case ASN_UNSIGNED:
      case ASN_TIMETICKS: // TODO note loss of type
      case ASN_COUNTER:
        return Response(oid, *vars->val.integer);
      case ASN_BOOLEAN:
        return Response(oid, static_cast<bool>(*vars->val.integer));
      case ASN_OBJECT_ID:
        return Response(
            oid,
            Oid(vars->val.objid, vars->val_len / sizeof(::oid)).toString());
      case ASN_BIT_STR:
      case ASN_NULL:
      default:
        return ErrorResponse(
            folly::sformat("snmp code path undefined of type {}", vars->type));
    }
  }

  return ErrorResponse("error");
}

Response Channel::processResponse(Request* request, snmp_pdu& response) {
  switch (response.errstat) {
    case SNMP_ERR_NOERROR:
      return processVars(response.variables);
    case SNMPERR_TIMEOUT:
      return ErrorResponse(folly::sformat(
          "snmp timeout {} on oid {}", peer, request->oid.toString()));
    case SNMP_ERR_NOSUCHNAME:
      for (netsnmp_variable_list* vars = response.variables; vars != nullptr;
           vars = vars->next_variable) {
        Oid oid(vars->name, vars->name_length);
        switch (vars->type) {
          case SNMP_ENDOFMIBVIEW:
          case SNMP_NOSUCHOBJECT:
          case SNMP_NOSUCHINSTANCE:
          case ASN_NULL:
            return Response(oid, nullptr);
          default:
            break;
        }
      }
      [[fallthrough]] default
          : return ErrorResponse(
                std::string("snmp packet error ") +
                snmp_errstring(static_cast<int>(response.errstat)) +
                " for oid " + firstOid(response.variables).toString() +
                " errno " + folly::to<std::string>(response.errstat));
  }
}

folly::Future<Responses> Channel::walk(const channels::snmp::Oid& tree) {
  return walk(tree, tree, {});
}

folly::Future<Responses> Channel::walk(
    const channels::snmp::Oid& tree,
    const channels::snmp::Oid& current,
    Responses responses) {
  return asyncGetNext(current).thenValue(
      [this, tree, current, responses = std::move(responses)](
          auto response) mutable {
        // TODO handle error
        if (not response.value.isNull() and response.oid.isDescendant(tree)) {
          responses.emplace_back(response);
          return this->walk(tree, response.oid, std::move(responses));
        } else {
          return folly::makeFuture<Responses>(std::move(responses));
        }
      });
}

} // namespace snmp
} // namespace channels
} // namespace devmand
