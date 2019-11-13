/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "Serializers.h"
#include <orc8r/protos/redis.pb.h>

using magma::orc8r::RedisState;
using google::protobuf::Message;
namespace magma {

std::function<bool(
  const Message&,
  std::string&,
  uint64_t&)> get_proto_serializer() {
  return [](const Message& message, std::string& str_out, uint64_t& version) -> bool {
    auto can_parse = message.SerializeToString(&str_out);
    if (!can_parse) {
      return false;
    }
    auto redis_state = RedisState();
    redis_state.set_version(version);
    redis_state.set_serialized_msg(str_out);
    return redis_state.SerializeToString(&str_out);
  };
}

std::function<bool(const std::string&, Message&)> get_proto_deserializer() {
  return [](const std::string& str, Message& msg_out) -> bool {
    RedisState redis_state;
    auto can_parse = redis_state.ParseFromString(str);
    if (!can_parse) {
      return false;
    }
    return msg_out.ParseFromString(redis_state.serialized_msg());
  };
}
}
