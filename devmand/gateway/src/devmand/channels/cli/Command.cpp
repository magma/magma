// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Command.h>

namespace devmand {
namespace channels {
namespace cli {

using std::string;
using std::vector;

Command::Command(const std::string _command, int _idx, bool skipCache)
    : command(_command), idx(_idx), skipCache_(skipCache) {}

// Factory methods
ReadCommand ReadCommand::create(const std::string& cmd, bool skipCache) {
  return ReadCommand(cmd, commandCounter++, skipCache);
}
WriteCommand WriteCommand::create(const std::string& cmd, bool skipCache) {
  return WriteCommand(cmd, commandCounter++, skipCache);
}
static const char DELIMITER = '\n';

bool Command::isMultiCommand() {
  return command.find(DELIMITER) != std::string::npos;
}

vector<Command> Command::splitMultiCommand() {
  string str = command;
  vector<Command> commands;
  string::size_type pos = 0;
  string::size_type prev = 0;
  while ((pos = str.find(DELIMITER, prev)) != std::string::npos) {
    commands.emplace_back(WriteCommand::create(str.substr(prev, pos - prev)));
    prev = pos + 1;
  }

  commands.emplace_back(ReadCommand::create(str.substr(prev)));
  return commands;
}

ReadCommand::ReadCommand(const string& _command, int _idx, bool _skipCache)
    : Command(_command, _idx, _skipCache) {}

ReadCommand ReadCommand::create(const Command& cmd) {
  return create(cmd.raw(), cmd.skipCache());
}

ReadCommand& ReadCommand::operator=(const ReadCommand& other) {
  this->command = other.command;
  this->idx = other.idx;
  this->skipCache_ = other.skipCache_;
  return *this;
}

ReadCommand::ReadCommand(const ReadCommand& rc)
    : Command(rc.raw(), rc.idx, rc.skipCache()) {}

WriteCommand::WriteCommand(const string& _command, int _idx, bool _skipCache)
    : Command(_command, _idx, _skipCache) {}

WriteCommand WriteCommand::create(const Command& cmd) {
  return create(cmd.raw(), cmd.skipCache());
}

WriteCommand::WriteCommand(const WriteCommand& wc)
    : Command(wc.raw(), wc.idx, wc.skipCache()) {}

WriteCommand& WriteCommand::operator=(const WriteCommand& other) {
  this->command = other.command;
  this->idx = other.idx;
  this->skipCache_ = other.skipCache_;
  return *this;
}

} // namespace cli
} // namespace channels
} // namespace devmand
