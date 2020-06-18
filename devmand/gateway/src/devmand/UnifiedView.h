// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <map>
#include <string>

#include <folly/dynamic.h>

#include <devmand/devices/Id.h>

namespace devmand {

// TODO convert this to ydk?
using YangModelBundle = folly::dynamic;
using UnifiedView = std::map<devices::Id, YangModelBundle>;
using SharedUnifiedView = folly::Synchronized<UnifiedView>;

} // namespace devmand
