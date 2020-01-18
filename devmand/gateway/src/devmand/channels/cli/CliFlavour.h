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

using devmand::channels::cli::sshsession::SessionAsync;
using folly::Future;
using folly::SemiFuture;
using folly::Timekeeper;
using folly::Unit;
using std::shared_ptr;
using std::string;

namespace devmand {
namespace channels {
namespace cli {

static const char* const UBIQUITI = "ubiquiti";

static const std::chrono::milliseconds delayDelta =
    std::chrono::milliseconds(100);

using std::shared_ptr;
using std::string;

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
  Future<folly::Optional<string>> resolvePromptAsync(
      shared_ptr<SessionAsync> session,
      const string& newline,
      std::chrono::milliseconds delay,
      shared_ptr<Timekeeper> timekeeper);
  Future<string> resolvePrompt(
      shared_ptr<SessionAsync> session,
      const string& newline,
      std::chrono::milliseconds delay,
      shared_ptr<Timekeeper> timekeeper);

 public:
  DefaultPromptResolver() = default;

  Future<string> resolvePrompt(
      shared_ptr<SessionAsync> session,
      const string& newline,
      shared_ptr<Timekeeper> timekeeper);
  void removeEmptyStrings(std::vector<string>& split) const;
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
  shared_ptr<folly::Timekeeper> timekeeper;

 public:
  static shared_ptr<CliFlavour> create(string flavour);

  std::unique_ptr<PromptResolver> resolver;
  std::unique_ptr<CliInitializer> initializer;
  string newline;

  CliFlavour(
      std::unique_ptr<PromptResolver>&& _resolver,
      std::unique_ptr<CliInitializer>&& _initializer,
      string _newline = "\n");
};

} // namespace cli
} // namespace channels
} // namespace devmand
