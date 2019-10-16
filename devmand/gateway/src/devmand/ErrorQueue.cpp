// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/ErrorQueue.h>

namespace devmand {

void ErrorQueue::add(std::string&& error) {
  errors.emplace_back(std::forward<std::string>(error));
  // on max size, discard oldest error
  if (errors.size() > maxSize) {
    errors.pop_front();
  }
}

folly::dynamic ErrorQueue::get() {
  auto ret = folly::dynamic::array();
  for (auto& error : errors) {
    ret.push_back(error);
  }
  return ret;
}

} // namespace devmand
