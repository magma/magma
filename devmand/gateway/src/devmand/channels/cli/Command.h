// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/algorithm/string/replace.hpp>
#include <folly/futures/Future.h>
#include <atomic>
#include <iostream>

using std::string;
using std::vector;

namespace devmand {
namespace channels {
namespace cli {
/*
 * Command struct encapsulating a string to be executed on a device.
 */
static std::atomic<int> commandCounter;

class Command {
 protected:
  explicit Command(std::string _command, int idx, bool skipCache);
  string command;
  int idx;
  bool skipCache_;

 public:
  Command() = delete;
  bool isMultiCommand();
  vector<Command> splitMultiCommand();

  string raw() const {
    return command;
  }

  bool skipCache() const {
    return skipCache_;
  }

  int getIdx() const {
    return idx;
  }

  friend std::ostream& operator<<(std::ostream& _stream, Command const& c) {
    auto rawCmd = c.raw();
    boost::replace_all(rawCmd, "\n", "\\n");
    boost::replace_all(rawCmd, "\r", "\\r");
    boost::replace_all(rawCmd, "\t", "\\t");
    _stream << rawCmd;
    return _stream;
  }
};

class WriteCommand : public Command {
 public:
  static WriteCommand create(const std::string& cmd, bool skipCache = false);
  static WriteCommand create(const Command& cmd);

  WriteCommand(const WriteCommand& wc);
  WriteCommand& operator=(const WriteCommand& other);

 private:
  WriteCommand(const string& command, int idx, bool skipCache);
};

class ReadCommand : public Command {
 public:
  static ReadCommand create(const std::string& cmd, bool skipCache = false);
  static ReadCommand create(const Command& cmd);

  ReadCommand& operator=(const ReadCommand& other);
  ReadCommand(const ReadCommand& rc);

 private:
  ReadCommand(const string& command, int idx, bool skipCache);
};

} // namespace cli
} // namespace channels
} // namespace devmand
