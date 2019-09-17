// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <map>
#include <string>

namespace devmand {

using UnifiedView = std::map<std::string, std::string>;
using SharedUnifiedView = folly::Synchronized<UnifiedView>;

} // namespace devmand
