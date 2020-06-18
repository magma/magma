// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <set>
#include <string>
#include <tuple>

#include <folly/Format.h>
#include <folly/GLog.h>

namespace devmand {

namespace utils {

class ConfigGenerator final {
 public:
  ConfigGenerator(
      const std::string& configFile_,
      const std::string& fileTemplate_ = "{}");
  ConfigGenerator() = delete;
  ~ConfigGenerator() = default;
  ConfigGenerator(const ConfigGenerator&) = delete;
  ConfigGenerator& operator=(const ConfigGenerator&) = delete;
  ConfigGenerator(ConfigGenerator&&) = delete;
  ConfigGenerator& operator=(ConfigGenerator&&) = delete;

 public:
  template <class... Args>
  void add(const std::string& templateS, Args&&... args) {
    if (entries.emplace(folly::sformat(templateS, std::forward<Args>(args)...))
            .second) {
      rewrite();
    } else {
      LOG(ERROR) << "Failed to add entry";
    }
  }

  template <class... Args>
  void remove(Args&&... args) {
    if (entries.erase(folly::sformat(std::forward<Args>(args)...)) != 1) {
      LOG(ERROR) << "Failed to delete entry ";
    } else {
      rewrite();
    }
  }

 private:
  void rewrite();

 private:
  std::string configFile;
  std::string fileTemplate;
  std::set<std::string> entries;
};

} // namespace utils

} // namespace devmand
