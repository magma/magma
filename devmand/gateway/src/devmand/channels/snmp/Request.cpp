// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/Request.h>

namespace devmand {
namespace channels {
namespace snmp {

Request::Request(Channel* channel_, Oid oid_) : channel(channel_), oid(oid_) {}

} // namespace snmp
} // namespace channels
} // namespace devmand
