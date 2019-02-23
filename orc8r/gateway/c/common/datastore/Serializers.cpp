/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "Serializers.h"

using google::protobuf::Message;
namespace magma {

std::function<bool(const Message&, std::string&)> get_proto_serializer() {
  return [](const Message& message, std::string& str_out) -> bool {
    return message.SerializeToString(&str_out);
  };
}

std::function<bool(const std::string&, Message&)> get_proto_deserializer() {
  return [](const std::string& str, Message& msg_out) -> bool {
    return msg_out.ParseFromString(str);
  };
}
}
