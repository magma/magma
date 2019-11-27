// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/futures/Future.h>

namespace devmand {
namespace channels {
namespace cli {

using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::PromptResolver;
using devmand::channels::cli::sshsession::SshSessionAsync;
using std::shared_ptr;
using std::string;

class PromptAwareCli : public Cli {
 private:
  shared_ptr<SshSessionAsync> session;
  shared_ptr<CliFlavour> cliFlavour;

 public:
  PromptAwareCli(
      shared_ptr<SshSessionAsync> session,
      shared_ptr<CliFlavour> cliFlavour);

  string prompt;

 public:
  void init(
      const string hostname,
      const int port,
      const string username,
      const string password);
  void resolvePrompt();
  void initializeCli();
  folly::Future<std::string> executeAndRead(const Command& cmd);

  folly::Future<std::string> execute(const Command& cmd);
};

} // namespace cli
} // namespace channels
} // namespace devmand
