// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/Optional.h>
#include <chrono>
#include <memory>
#include <regex>

using devmand::channels::cli::sshsession::SessionAsync;
using namespace std;
using namespace folly;

namespace devmand {
namespace channels {
namespace cli {

static const char* const UBIQUITI = "ubiquiti";

static const chrono::milliseconds delayDelta = chrono::milliseconds(100);

class PromptResolver {
 public:
  PromptResolver() = default;
  virtual Future<string> resolvePrompt(
      shared_ptr<SessionAsync> session,
      const string& newline,
      shared_ptr<Timekeeper> timekeeper) = 0;
  virtual ~PromptResolver() = default;
};

class DefaultPromptResolver : public PromptResolver {
 private:
  Future<Optional<string>> resolvePromptAsync(
      shared_ptr<SessionAsync> session,
      const string& newline,
      chrono::milliseconds delay,
      shared_ptr<Timekeeper> timekeeper);
  Future<string> resolvePrompt(
      shared_ptr<SessionAsync> session,
      const string& newline,
      chrono::milliseconds delay,
      shared_ptr<Timekeeper> timekeeper);

 public:
  DefaultPromptResolver() = default;

  Future<string> resolvePrompt(
      shared_ptr<SessionAsync> session,
      const string& newline,
      shared_ptr<Timekeeper> timekeeper);
  void removeEmptyStrings(vector<string>& split) const;
};

class CliInitializer {
 public:
  virtual ~CliInitializer() = default;
  virtual SemiFuture<Unit> initialize(
      shared_ptr<SessionAsync> session,
      string secret) = 0;
};

class EmptyInitializer : public CliInitializer {
 public:
  SemiFuture<Unit> initialize(shared_ptr<SessionAsync> session, string secret)
      override;
  ~EmptyInitializer() override = default;
};

class UbiquitiInitializer : public CliInitializer {
 public:
  SemiFuture<Unit> initialize(shared_ptr<SessionAsync> session, string secret)
      override;
  ~UbiquitiInitializer() override = default;
};

class CliFlavour {
 private:
  const shared_ptr<PromptResolver> resolver;
  const shared_ptr<CliInitializer> initializer;
  const string newline;
  // regex that matches show config command
  const regex baseShowConfig;
  // match index that selects group containing just the command
  const unsigned int baseShowConfigIdx;
  const Optional<char> singleIndentChar;
  const string configSubsectionEnd;

 public:
  static shared_ptr<CliFlavour> create(string flavour);

  CliFlavour(
      unique_ptr<PromptResolver>&& _resolver,
      unique_ptr<CliInitializer>&& _initializer,
      string _newline,
      regex baseShowConfig,
      unsigned int baseShowConfigIdx,
      Optional<char> singleIndentChar,
      string configSubsectionEnd);

  shared_ptr<PromptResolver> getResolver();

  shared_ptr<CliInitializer> getInitializer();

  string getNewline();

  /*
   * If supplied command is command to show running config or other command
   * where caching is supported, return index where the base command ends. E.g.
   * 'sh run interface 0/14' should return index 6:
   *             ^
   * Otherwise return none value.
   */
  Optional<size_t> getBaseShowConfigIdx(const string cmd) const;

  Optional<char> getSingleIndentChar();

  string getConfigSubsectionEnd();

  vector<string> splitSubcommands(string subcommands);
};

} // namespace cli
} // namespace channels
} // namespace devmand
