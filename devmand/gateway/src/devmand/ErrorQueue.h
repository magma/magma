// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <list>
#include <string>

#include <folly/Synchronized.h>
#include <folly/dynamic.h>

namespace devmand {

// TODO add timestamps
class ErrorQueue final {
 public:
  ErrorQueue(unsigned int maxSize_ = 10) : maxSize(maxSize_) {}
  ~ErrorQueue() = default;
  ErrorQueue(const ErrorQueue&) = delete;
  ErrorQueue& operator=(const ErrorQueue&) = delete;
  ErrorQueue(ErrorQueue&&) = delete;
  ErrorQueue& operator=(ErrorQueue&&) = delete;

 public:
  void add(std::string&& error);
  folly::dynamic get();

 private:
  folly::Synchronized<std::list<std::string>> errors;
  unsigned int maxSize{0};
};

} // namespace devmand
