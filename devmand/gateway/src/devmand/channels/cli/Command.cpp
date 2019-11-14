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

Command::Command(const std::string _command) : command(_command) {}

// Factory methods
Command Command::makeReadCommand(const std::string& cmd) {
  return Command(cmd);
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
        commands.emplace_back(Command::makeReadCommand(str.substr(prev, pos - prev)));
        prev = pos + 1;
    }

    commands.emplace_back(Command::makeReadCommand(str.substr(prev)));
    return commands;
}


} // namespace cli
} // namespace channels
} // namespace devmand