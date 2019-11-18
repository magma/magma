// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/executors/IOThreadPoolExecutor.h>

using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::Command;
using devmand::channels::cli::PromptAwareCli;
using devmand::channels::cli::PromptResolver;
using std::string;

namespace devmand {
namespace channels {
namespace cli {

void PromptAwareCli::resolvePrompt() {
  this->prompt =
      cliFlavour->resolver->resolvePrompt(session, cliFlavour->newline);
}

void PromptAwareCli::initializeCli() {
  cliFlavour->initializer->initialize(session);
}

folly::Future<string> PromptAwareCli::executeAndRead(const Command& cmd) {
  const string& command = cmd.toString();

  return session->write(command)
      .thenValue([=](...) { return session->readUntilOutput(command); })
      .thenValue([=](const string& output) {
        auto returnOutputParameter = [output](...) { return output; };
        return session->write(cliFlavour->newline)
            .thenValue(returnOutputParameter);
      })
      .thenValue([=](const string& output) {
        auto concatOutputParameter = [output](const string& readUntilOutput) {
          return output + readUntilOutput;
        };
        return session->readUntilOutput(prompt).thenValue(
            concatOutputParameter);
      });
}

PromptAwareCli::PromptAwareCli(
    shared_ptr<SshSessionAsync> _session,
    shared_ptr<CliFlavour> _cliFlavour)
    : session(_session), cliFlavour(_cliFlavour) {}

void PromptAwareCli::init( // TODO remove
    const string hostname,
    const int port,
    const string username,
    const string password) {
  session->openShell(hostname, port, username, password).get();
}

folly::Future<std::string> PromptAwareCli::execute(const Command& cmd) {
  const string& command = cmd.toString();
  return session->write(command)
      .thenValue([=](...) { return session->readUntilOutput(command); })
      .thenValue([=](const string& output) {
        return session->write(cliFlavour->newline)
            .thenValue([output, command](...) { return output + command; });
      });
}
} // namespace cli
} // namespace channels
} // namespace devmand