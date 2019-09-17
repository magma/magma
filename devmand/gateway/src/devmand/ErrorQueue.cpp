// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <devmand/ErrorQueue.h>

namespace devmand {

void ErrorQueue::add(std::string&& error) {
  errors.emplace_front(std::forward<std::string>(error));

  if (maxSize > errors.size()) {
    errors.pop_back();
  }
}

folly::dynamic ErrorQueue::get() {
  folly::dynamic ret = folly::dynamic::array;
  for (auto& error : errors) {
    ret.push_back(error);
  }
  return ret;
}

} // namespace devmand
