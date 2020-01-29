// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/channels/cli/SshSessionAsync.h>

using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::Command;
using devmand::channels::cli::PromptAwareCli;
using devmand::channels::cli::PromptResolver;
using std::string;
using namespace folly;

namespace devmand {
namespace channels {
namespace cli {

SemiFuture<Unit> PromptAwareCli::resolvePrompt() {
  return sharedCliFlavour->getResolver()
      ->resolvePrompt(
          sharedSession, sharedCliFlavour->getNewline(), sharedTimekeeper)
      .thenValue([params = promptAwareParameters](string _prompt) {
        boost::mutex::scoped_lock scoped_lock(params->promptMutex);
        params->prompt = _prompt;
      });
}

SemiFuture<Unit> PromptAwareCli::initializeCli(const string secret) {
  return sharedCliFlavour->getInitializer()->initialize(sharedSession, secret);
}

SemiFuture<std::string> PromptAwareCli::executeRead(const ReadCommand cmd) {
  const string& command = cmd.raw();

  return sharedSession->write(command)
      .semi()
      .via(sharedExecutor.get())
      .thenValue([params = promptAwareParameters, command, cmd](...) {
        if (auto _session = params->session.lock()) {
          MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                       << ") written command";
          SemiFuture<string> result = _session->readUntilOutput(command).semi();
          MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                       << ") obtained future readUntilOutput";
          return move(result);
        } else {
          return makeSemiFuture<std::string>(
              DisconnectedException("SSH session expired"));
        }
      })
      .thenValue([params = promptAwareParameters, cmd](const string& output) {
        if (auto _session = params->session.lock()) {
          if (auto _cliFlavour = params->cliFlavour.lock()) {
            if (auto _executor = params->executor.lock()) {
              return _session->write(_cliFlavour->getNewline())
                  .via(_executor.get())
                  .thenValue([params, output, cmd](...) {
                    MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                                 << ") written newline";
                    return output;
                  });
            }
          }
        }
        return makeFuture<std::string>(
            DisconnectedException("SSH session expired"));
      })
      .thenValue([params = promptAwareParameters, cmd](const string& output) {
        if (auto _session = params->session.lock()) {
          if (auto _executor = params->executor.lock()) {
            return _session->readUntilOutput(params->prompt)
                .via(_executor.get())
                .thenValue([id = params->id, output, cmd](
                               const string& readUntilOutput) {
                  // this might never run, do not capture params
                  MLOG(MDEBUG) << "[" << id << "] (" << cmd
                               << ") readUntilOutput - read result";
                  return output + readUntilOutput;
                })
                .semi();
          }
        }
        return makeSemiFuture<std::string>(
            DisconnectedException("SSH session expired"));
      })
      .semi();
}

PromptAwareCli::PromptAwareCli(
    string _id,
    shared_ptr<SessionAsync> _session,
    shared_ptr<CliFlavour> _cliFlavour,
    shared_ptr<Executor> _executor,
    shared_ptr<Timekeeper> _timekeeper)
    : id(_id),
      sharedSession(_session),
      sharedCliFlavour(_cliFlavour),
      sharedExecutor(_executor),
      sharedTimekeeper(_timekeeper) {
  promptAwareParameters = std::make_shared<PromptAwareParameters>(
      _id, _session, _cliFlavour, _executor, _timekeeper);
}

SemiFuture<Unit> PromptAwareCli::destroy() {
  string _id = id;
  // TODO cancel timekeeper futures
  return sharedSession->destroy();
}

PromptAwareCli::~PromptAwareCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~PromptAwareCli: started";
  destroy().get();
  MLOG(MDEBUG) << "[" << id << "] "
               << "~PromptAwareCli: done";
}

SemiFuture<std::string> PromptAwareCli::executeWrite(const WriteCommand cmd) {
  const string& command = cmd.raw();
  return sharedSession->write(command)
      .via(sharedExecutor.get())
      .thenValue([params = promptAwareParameters, command, cmd](...) {
        if (auto _session = params->session.lock()) {
          MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                       << ") written command";
          SemiFuture<string> result = _session->readUntilOutput(command).semi();
          MLOG(MDEBUG) << "[" << params->id << "] (" << cmd
                       << ") obtained future readUntilOutput";
          return move(result);
        } else {
          throw DisconnectedException("SSH session expired");
        }
      })
      .thenValue(
          [params = promptAwareParameters, command, cmd](const string& output) {
            if (auto _session = params->session.lock()) {
              if (auto _executor = params->executor.lock()) {
                if (auto _cliFlavour = params->cliFlavour.lock()) {
                  return _session->write(_cliFlavour->getNewline())
                      .via(_executor.get())
                      .thenValue([id = params->id, output, command, cmd](...) {
                        MLOG(MDEBUG)
                            << "[" << id << "] (" << cmd << ") written newline";
                        return output + command;
                      })
                      .semi();
                }
              }
            }
            throw DisconnectedException("SSH session expired");
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
