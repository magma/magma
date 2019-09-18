// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <string>

#include <devmand/channels/snmp/Response.h>

namespace devmand {
namespace channels {
namespace snmp {

// name or address of peer (may include transport specifier and/or port number)
using Peer = std::string;
using Community = std::string;
using Version = std::string;
using SecurityLevel = std::string;

} // namespace snmp
} // namespace channels
} // namespace devmand
