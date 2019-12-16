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

folly::SemiFuture<std::string> Channel::executeRead(const ReadCommand cmd) {
  MLOG(MDEBUG) << "[" << id << "] "
               << "Executing command and reading: "
               << "\"" << cmd << "\"";

  return cli->executeRead(cmd);
}

folly::SemiFuture<std::string> Channel::executeWrite(const WriteCommand cmd) {
  MLOG(MDEBUG) << "[" << id << "]"
               << "Executing command"
               << "\"" << cmd << "\"";

  return cli->executeWrite(cmd);
}

} // namespace cli
} // namespace channels
} // namespace devmand
