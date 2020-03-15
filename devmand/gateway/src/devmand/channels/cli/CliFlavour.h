// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/devices/cli/DeviceType.h>
#include <folly/Optional.h>
#include <chrono>
#include <memory>
#include <regex>

using devmand::channels::cli::sshsession::SessionAsync;
using devmand::devices::cli::DeviceType;
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

class CliFlavourParameters {
 public:
  const string newline;
  // regex that matches show config command
  const regex baseShowConfig;
  // match index that selects group containing just the command
  const unsigned int baseShowConfigIdx;
  const Optional<char> singleIndentChar;
  const string configSubsectionEnd;

  CliFlavourParameters(
      const string& _newline,
      const regex& _baseShowConfig,
      const unsigned int _baseShowConfigIdx,
      const Optional<char> _singleIndentChar,
      const string& _configSubsectionEnd)
      : newline(_newline),
        baseShowConfig(_baseShowConfig),
        baseShowConfigIdx(_baseShowConfigIdx),
        singleIndentChar(_singleIndentChar),
        configSubsectionEnd(_configSubsectionEnd) {}

  //  CliFlavourParameters& operator=(CliFlavourParameters&& x) = default;
  CliFlavourParameters& operator=(const CliFlavourParameters&) = default;
  CliFlavourParameters(CliFlavourParameters&&) = default;
  CliFlavourParameters(const CliFlavourParameters&) = default;
};

static CliFlavourParameters UBIQUITI_PARAMETERS{
    "\n",
    regex(R"(^((do )?sho?w? runn?i?n?g?-?c?o?n?f?i?g?).*)"),
    1,
    none,
    "exit"};

static CliFlavourParameters DEFAULT_PARAMETERS{
    "\n",
    regex(R"(^((do )?sho?w? runn?i?n?g?-?c?o?n?f?i?g?).*)"),
    1,
    ' ',
    "!"};

class CliFlavour {
 private:
  const shared_ptr<PromptResolver> resolver;
  const shared_ptr<CliInitializer> initializer;
  const shared_ptr<CliFlavourParameters> params;

 public:
  CliFlavour(
      unique_ptr<PromptResolver>&& _resolver,
      unique_ptr<CliInitializer>&& _initializer,
      shared_ptr<CliFlavourParameters> _params);

  static map<DeviceType, shared_ptr<CliFlavourParameters>>
  getHardcodedFlavours() {
    map<DeviceType, shared_ptr<CliFlavourParameters>> result;
    result.emplace(
        DeviceType::getDefaultInstance(),
        make_shared<CliFlavourParameters>(DEFAULT_PARAMETERS)); // TODO
    result.emplace(
        DeviceType(UBIQUITI, devmand::devices::cli::ANY_VERSION),
        make_shared<CliFlavourParameters>(UBIQUITI_PARAMETERS));
    return result;
  }

  static shared_ptr<CliFlavour> create(
      shared_ptr<CliFlavourParameters> parameters);

  static shared_ptr<CliFlavour> getDefaultInstance() {
    return create(
        make_shared<CliFlavourParameters>(DEFAULT_PARAMETERS)); // TODO
  }

  static shared_ptr<CliFlavour> getUbiquiti() {
    return create(
        make_shared<CliFlavourParameters>(UBIQUITI_PARAMETERS)); // TODO
  }

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
