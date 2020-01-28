// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Command.h>
#include <folly/futures/Future.h>
#include <folly/futures/Promise.h>
#include <magma_logging.h>

#include <chrono>
#include <thread>
#include <vector>

namespace devmand {
namespace channels {
namespace cli {
class Cli {
 public:
  Cli() = default;
  virtual ~Cli() = default;
  Cli(const Cli&) = delete;
  Cli& operator=(const Cli&) = delete;
  Cli(Cli&&) = delete;
  Cli& operator=(Cli&&) = delete;

 public:
  virtual folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) = 0;

  virtual folly::SemiFuture<std::string> executeWrite(
      const WriteCommand cmd) = 0;

  //! Destruct object asynchronously. Must be idempotent.
  //! Thread safety: it is expected that this method will only be called
  //! from the same thread as actual destructor.
  //! Destructor should call destroy().get() to block until all resources are
  //! cleaned. Calling destroy() on outer CLI layer should call destroy() on
  //! underlying layer immediately to start closing ssh session.
  virtual folly::SemiFuture<folly::Unit> destroy() = 0;
};

class CliException : public std::runtime_error {
 public:
  CliException(string msg) : std::runtime_error(msg) {}
};

class DisconnectedException : public CliException {
 public:
  DisconnectedException(string msg = "Not connected") : CliException(msg) {}
};

class CommandExecutionException : public CliException {
 public:
  CommandExecutionException(string msg = "Command execution failed")
      : CliException(msg) {}
};

class CommandTimeoutException : public CommandExecutionException {
 public:
  CommandTimeoutException(string msg = "Command execution timed out")
      : CommandExecutionException(msg) {}
};

} // namespace cli
} // namespace channels
} // namespace devmand
