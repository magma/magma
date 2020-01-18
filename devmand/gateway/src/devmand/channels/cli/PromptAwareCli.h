// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/futures/Future.h>

namespace devmand {
namespace channels {
namespace cli {

using boost::mutex;
using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::PromptResolver;
using devmand::channels::cli::sshsession::SessionAsync;
using folly::Executor;
using folly::SemiFuture;
using folly::Unit;
using std::shared_ptr;
using std::string;
using std::weak_ptr;

class PromptAwareCli : public Cli {
 private:
  string id;

  shared_ptr<SessionAsync> sharedSession;
  shared_ptr<CliFlavour> sharedCliFlavour;
  shared_ptr<Executor> sharedExecutor;
  shared_ptr<Timekeeper> sharedTimekeeper;

  struct PromptAwareParameters {
    string id;
    weak_ptr<SessionAsync> session;
    weak_ptr<CliFlavour> cliFlavour;
    weak_ptr<Executor> executor;
    weak_ptr<Timekeeper> timekeeper;
    mutex promptMutex;
    string prompt;

    PromptAwareParameters(
        const string& id,
        const shared_ptr<SessionAsync>& session,
        const shared_ptr<CliFlavour>& cliFlavour,
        const shared_ptr<Executor>& executor,
        const shared_ptr<Timekeeper>& timekeeper);
  };

  shared_ptr<PromptAwareParameters> promptAwareParameters;

 public:
  PromptAwareCli(
      string id,
      shared_ptr<SessionAsync> session,
      shared_ptr<CliFlavour> cliFlavour,
      shared_ptr<Executor> executor,
      shared_ptr<Timekeeper> timekeeper);
  static shared_ptr<PromptAwareCli> make(
      string id,
      shared_ptr<SessionAsync> session,
      shared_ptr<CliFlavour> cliFlavour,
      shared_ptr<Executor> executor,
      shared_ptr<Timekeeper> timekeeper);

  SemiFuture<Unit> destroy() override;

  ~PromptAwareCli();

  SemiFuture<Unit> resolvePrompt();
  SemiFuture<Unit> initializeCli(const string secret);
  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd);
  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd);
};

} // namespace cli
} // namespace channels
} // namespace devmand
