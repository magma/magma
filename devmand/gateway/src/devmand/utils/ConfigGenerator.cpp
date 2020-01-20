// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <algorithm>
#include <experimental/filesystem>

#include <folly/Format.h>
#include <folly/GLog.h>

#include <devmand/utils/ConfigGenerator.h>
#include <devmand/utils/FileUtils.h>

#include <boost/algorithm/string/join.hpp>

namespace devmand {

namespace utils {

ConfigGenerator::ConfigGenerator(
    const std::string& configFile_,
    const std::string& fileTemplate_)
    : configFile(configFile_), fileTemplate(fileTemplate_) {
  FileUtils::mkdir(
      std::experimental::filesystem::path(configFile).parent_path());
}

void ConfigGenerator::rewrite() {
  // TODO make this async
  FileUtils::write(
      configFile,
      folly::sformat(fileTemplate, boost::algorithm::join(entries, "")));
}

} // namespace utils

} // namespace devmand
