// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
