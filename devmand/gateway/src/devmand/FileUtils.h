// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

namespace devmand {

class FileUtils final {
 public:
  FileUtils() = delete;
  ~FileUtils() = delete;
  FileUtils(const FileUtils&) = delete;
  FileUtils& operator=(const FileUtils&) = delete;
  FileUtils(FileUtils&&) = delete;
  FileUtils& operator=(FileUtils&&) = delete;

 public:
  static bool touch(const std::string& filename);
  static void write(const std::string& filename, const std::string& contents);
  static std::string readContents(const std::string& filename);
  static bool mkdir(const std::string& dir);
};

} // namespace devmand
