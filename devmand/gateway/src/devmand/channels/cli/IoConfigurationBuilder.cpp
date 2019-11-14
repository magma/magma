// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/channels/cli/QueuedCli.h>
#include <devmand/channels/cli/ReadCachingCli.h>
#include <devmand/channels/cli/SshSession.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/channels/cli/SshSocketReader.h>
#include <folly/Singleton.h>
#include <folly/executors/IOThreadPoolExecutor.h>

namespace devmand {
namespace channels {
namespace cli {

using devmand::channels::cli::IoConfigurationBuilder;
using devmand::channels::cli::SshSocketReader;
using devmand::channels::cli::sshsession::readCallback;
using devmand::channels::cli::sshsession::SshSession;
using devmand::channels::cli::sshsession::SshSessionAsync;
using folly::EvictingCacheMap;
using folly::IOThreadPoolExecutor;
using std::make_shared;
using std::string;

// TODO executor?
shared_ptr<IOThreadPoolExecutor> executor =
    std::make_shared<IOThreadPoolExecutor>(10);

IoConfigurationBuilder::IoConfigurationBuilder(
    const DeviceConfig& _deviceConfig)
    : deviceConfig(_deviceConfig) {}

shared_ptr<Cli> IoConfigurationBuilder::getIo(
    shared_ptr<CliCache> commandCache) {
  MLOG(MDEBUG) << "Creating CLI ssh device for " << deviceConfig.id
               << " (host: " << deviceConfig.ip << ")";

  const auto& plaintextCliKv = deviceConfig.channelConfigs.at("cli").kvPairs;
  // crate session
  const std::shared_ptr<SshSessionAsync>& session =
      std::make_shared<SshSessionAsync>(deviceConfig.id, executor);
  // opening SSH connection
  session
      ->openShell(
          deviceConfig.ip,
          std::stoi(plaintextCliKv.at("port")),
          plaintextCliKv.at("username"),
          plaintextCliKv.at("password"))
      .get();

  shared_ptr<CliFlavour> cl =
      plaintextCliKv.find("flavour") != plaintextCliKv.end()
      ? CliFlavour::create(plaintextCliKv.at("flavour"))
      : CliFlavour::create("");

  // create CLI - how to create a CLI stack?
  const shared_ptr<PromptAwareCli>& cli =
      std::make_shared<PromptAwareCli>(session, cl);

  // initialize CLI
  cli->initializeCli();
  // resolve prompt needs to happen
  cli->resolvePrompt();
  // create async data reader
  event* sessionEvent = SshSocketReader::getInstance().addSshReader(
      readCallback, session->getSshFd(), session.get());
  session->setEvent(sessionEvent);

  // create caching cli
  const shared_ptr<ReadCachingCli>& ccli =
      std::make_shared<ReadCachingCli>(deviceConfig.id, cli, commandCache);

  // create Queued cli
  return std::make_shared<QueuedCli>(deviceConfig.id, ccli, executor);
}

} // namespace cli
} // namespace channels
} // namespace devmand