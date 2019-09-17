// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <list>
#include <string>

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
  // TODO make sync
  std::list<std::string> errors;
  unsigned int maxSize{0};
};

} // namespace devmand
