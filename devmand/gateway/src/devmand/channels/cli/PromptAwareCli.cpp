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

SemiFuture<Unit> PromptAwareCli::resolvePrompt() {
  return promptAwareParameters->cliFlavour->resolver
      ->resolvePrompt(
          promptAwareParameters->session,
          promptAwareParameters->cliFlavour->newline,
          promptAwareParameters->timekeeper)
      .thenValue([params = promptAwareParameters](string _prompt) {
        params->prompt = _prompt;
      });
}

SemiFuture<Unit> PromptAwareCli::initializeCli(const string secret) {
  return promptAwareParameters->cliFlavour->initializer->initialize(
      promptAwareParameters->session, secret);
}

folly::SemiFuture<std::string> PromptAwareCli::executeRead(
    const ReadCommand cmd) {
  const string& command = cmd.raw();

  return promptAwareParameters->session->write(command)
      .semi()
      .via(promptAwareParameters->executor.get())
      .thenValue([params = promptAwareParameters, command, cmd](...) {
        MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                     << ") written command";
        SemiFuture<string> result =
            params->session->readUntilOutput(command).semi();
        MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                     << ") obtained future readUntilOutput";
        return move(result);
      })
      .semi()
      .via(promptAwareParameters->executor.get())
      .thenValue([params = promptAwareParameters, cmd](const string& output) {
        return params->session->write(params->cliFlavour->newline)
            .semi()
            .via(params->executor.get())
            .thenValue([params, output, cmd](...) {
              MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                           << ") written newline";
              return output;
            })
            .semi();
      })
      .semi()
      .via(promptAwareParameters->executor.get())
      .thenValue([params = promptAwareParameters, cmd](const string& output) {
        return params->session->readUntilOutput(params->prompt)
            .semi()
            .via(params->executor.get())
            .thenValue(
                [id = params->id, output, cmd](const string& readUntilOutput) {
                  // this might never run, do not capture params
                  MLOG(MDEBUG) << "[" << id << "] (" << cmd
                               << ") readUntilOutput - read result";
                  return output + readUntilOutput;
                })
            .semi();
      })
      .semi();
}

PromptAwareCli::PromptAwareCli(
    string id,
    shared_ptr<SessionAsync> _session,
    shared_ptr<CliFlavour> _cliFlavour,
    shared_ptr<Executor> _executor,
    shared_ptr<Timekeeper> _timekeeper) {
  promptAwareParameters = std::make_shared<PromptAwareParameters>(
      id, _session, _cliFlavour, _executor, _timekeeper);
}

PromptAwareCli::~PromptAwareCli() {
  string id = promptAwareParameters->id;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~PromptAwareCli started";
  while (promptAwareParameters.use_count() > 1) {
    MLOG(MDEBUG) << "[" << id << "] "
                 << "~PromptAwareCli sleeping";
    std::this_thread::sleep_for(std::chrono::seconds(1));
  }

  // Closing session explicitly to finish executing future and disconnect before
  // releasing the rest of state
  promptAwareParameters->session = nullptr;
  promptAwareParameters = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~PromptAwareCli done";
}

folly::SemiFuture<std::string> PromptAwareCli::executeWrite(
    const WriteCommand cmd) {
  const string& command = cmd.raw();
  return promptAwareParameters->session->write(command)
      .semi()
      .via(promptAwareParameters->executor.get())
      .thenValue([params = promptAwareParameters, command, cmd](...) {
        MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                     << ") written command";
        SemiFuture<string> result =
            params->session->readUntilOutput(command).semi();
        MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                     << ") obtained future readUntilOutput";
        return move(result);
      })
      .thenValue([params = promptAwareParameters, command, cmd](
                     const string& output) {
        return params->session->write(params->cliFlavour->newline)
            .semi()
            .via(params->executor.get())
            .thenValue([id = params->id, output, command, cmd](...) {
              MLOG(MDEBUG) << "[" << id << "] (" << cmd << ") written newline";
              return output + command;
            })
            .semi();
      })
      .semi();
}

shared_ptr<PromptAwareCli> PromptAwareCli::make(
    string id,
    shared_ptr<SessionAsync> session,
    shared_ptr<CliFlavour> cliFlavour,
    shared_ptr<Executor> executor,
    shared_ptr<Timekeeper> timekeeper) {
  return std::make_shared<PromptAwareCli>(
      id, session, cliFlavour, executor, timekeeper);
}

PromptAwareCli::PromptAwareParameters::PromptAwareParameters(
    const string& _id,
    const shared_ptr<SessionAsync>& _session,
    const shared_ptr<CliFlavour>& _cliFlavour,
    const shared_ptr<Executor>& _executor,
    const shared_ptr<Timekeeper>& _timekeeper)
    : id(_id),
      session(_session),
      cliFlavour(_cliFlavour),
      executor(_executor),
      timekeeper(_timekeeper) {}
} // namespace cli
} // namespace channels
} // namespace devmand
