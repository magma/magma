/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <functional>
#include <google/protobuf/message.h>

using google::protobuf::Message;
namespace magma {
// This file defines some common serializer methods

/**
 * Serialize a protobuf message into the standard string form
 */
std::function<bool(const Message&, std::string&, uint64_t&)>
get_proto_serializer();

/**
 * Deserialize a string into a protobuf message
 */
std::function<bool(const std::string&, Message&)>
get_proto_deserializer();

}
