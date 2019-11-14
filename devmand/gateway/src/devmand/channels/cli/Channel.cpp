// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Channel.h>

namespace devmand {
namespace channels {
namespace cli {

Channel::Channel(
    string _id,
    const std::shared_ptr<devmand::channels::cli::Cli> _cli)
    : id(_id), cli(_cli) {}

Channel::~Channel() {}

folly::Future<std::string> Channel::executeAndRead(const Command& cmd) {
  MLOG(MDEBUG) << "[" << id << "] "
               << "Executing command and reading: "
               << "\"" << cmd << "\"";

  return cli->executeAndRead(cmd);
}

folly::Future<std::string> Channel::execute(const Command& cmd) {
  MLOG(MDEBUG) << "[" << id << "]"
               << "Executing command"
               << "\"" << cmd << "\"";

  return cli->execute(cmd);
}

} // namespace cli
} // namespace channels
} // namespace devmand
