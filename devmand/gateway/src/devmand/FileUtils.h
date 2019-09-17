// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
