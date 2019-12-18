// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <fcntl.h>
#include <sys/stat.h>
#include <unistd.h>
#include <cerrno>
#include <fstream>
#include <iostream>
#include <string>
#include <system_error>

#include <folly/GLog.h>

#include <devmand/utils/FileUtils.h>

namespace devmand {

// TODO this should really be async and handle errors but works for now.
void FileUtils::write(
    const std::string& filename,
    const std::string& contents) {
  std::ofstream out(filename);
  out << contents;
  out.close();
}

bool FileUtils::touch(const std::string& filename) {
  int fd =
      open(filename.c_str(), O_WRONLY | O_CREAT | O_NOCTTY | O_NONBLOCK, 0666);
  bool status{true};
  if (fd < 0 and errno != EISDIR) {
    status = false;
  } else if (utimensat(AT_FDCWD, filename.c_str(), nullptr, 0)) {
    status = false;
  }
  close(fd);
  return status;
}

std::string FileUtils::readContents(const std::string& filename) {
  std::string contents;
  std::ifstream in(filename, std::ios::in | std::ios::binary);
  if (in) {
    in.seekg(0, std::ios::end);
    contents.resize(static_cast<size_t>(in.tellg()));
    in.seekg(0, std::ios::beg);
    in.read(&contents[0], static_cast<std::streamsize>(contents.size()));
    in.close();
  }
  // TODO handle error
  return contents;
}

bool FileUtils::mkdir(const std::string& dir) {
  if (::mkdir(dir.c_str(), 0666) < 0 and
      errno != EEXIST) { // TODO allow mode to be passed
    LOG(ERROR) << "mkdir error "
               << std::error_code(errno, std::system_category());
    return false;
  }
  return true;
}

} // namespace devmand
