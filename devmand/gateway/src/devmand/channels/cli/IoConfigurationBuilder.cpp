// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG

#include <magma_logging.h>

#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/channels/cli/KeepaliveCli.h>
#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/channels/cli/QueuedCli.h>
#include <devmand/channels/cli/ReadCachingCli.h>
#include <devmand/channels/cli/ReconnectingCli.h>
#include <devmand/channels/cli/SshSession.h>
#include <devmand/channels/cli/SshSocketReader.h>
#include <devmand/channels/cli/TimeoutTrackingCli.h>
#include <folly/Singleton.h>

namespace devmand {
namespace channels {
namespace cli {

using devmand::channels::cli::IoConfigurationBuilder;
using devmand::channels::cli::SshSocketReader;
using devmand::channels::cli::sshsession::readCallback;
using devmand::channels::cli::sshsession::SshSession;

using folly::EvictingCacheMap;

IoConfigurationBuilder::IoConfigurationBuilder(
    shared_ptr<ConnectionParameters> _connectionParams)
    : connectionParameters(_connectionParams) {}

IoConfigurationBuilder::IoConfigurationBuilder(
    const DeviceConfig& deviceConfig,
    channels::cli::Engine& engine) {
  const std::map<std::string, std::string>& plaintextCliKv =
      deviceConfig.channelConfigs.at("cli").kvPairs;

  connectionParameters = makeConnectionParameters(
      deviceConfig.id,
      deviceConfig.ip,
      plaintextCliKv.at("username"),
      plaintextCliKv.at("password"),
      loadConfigValue(plaintextCliKv, "flavour", ""),
      folly::to<int>(plaintextCliKv.at("port")),
      toSeconds(loadConfigValue(
          plaintextCliKv, configKeepAliveIntervalSeconds, "60")),
      toSeconds(loadConfigValue(
          plaintextCliKv, configMaxCommandTimeoutSeconds, "60")),
      toSeconds(
          loadConfigValue(plaintextCliKv, reconnectingQuietPeriodConfig, "5")),
      std::stol(
          loadConfigValue(plaintextCliKv, sshConnectionTimeoutConfig, "30")),
      engine.getTimekeeper(),
      engine.getExecutor(Engine::executorRequestType::sshCli),
      engine.getExecutor(Engine::executorRequestType::paCli),
      engine.getExecutor(Engine::executorRequestType::rcCli),
      engine.getExecutor(Engine::executorRequestType::ttCli),
      engine.getExecutor(Engine::executorRequestType::qCli),
      engine.getExecutor(Engine::executorRequestType::rCli),
      engine.getExecutor(Engine::executorRequestType::kaCli));
}

IoConfigurationBuilder::~IoConfigurationBuilder() {
  MLOG(MDEBUG) << "~IoConfigurationBuilder";
}

shared_ptr<Cli> IoConfigurationBuilder::createAll(
    shared_ptr<CliCache> commandCache) {
  return createAllUsingFactory(commandCache);
}

chrono::seconds IoConfigurationBuilder::toSeconds(const string& value) {
  return chrono::seconds(folly::to<int>(value));
}

shared_ptr<Cli> IoConfigurationBuilder::createAllUsingFactory(
    shared_ptr<CliCache> commandCache) {
  function<SemiFuture<shared_ptr<Cli>>()> cliFactory =
      [params = connectionParameters, commandCache]() {
        return createPromptAwareCli(params).thenValue(
            [params, commandCache](shared_ptr<Cli> sshCli) -> shared_ptr<Cli> {
              MLOG(MDEBUG) << "[" << params->id << "] "
                           << "Creating cli layers rcclli, ttcli, qcli";
              // create caching cli
              const shared_ptr<ReadCachingCli>& rccli =
                  std::make_shared<ReadCachingCli>(
                      params->id, sshCli, commandCache, params->rcExecutor);
              // create timeout tracker
              shared_ptr<TimeoutTrackingCli> ttcli = TimeoutTrackingCli::make(
                  params->id,
                  rccli,
                  params->timekeeper,
                  params->ttExecutor,
                  params->cmdTimeout);
              // create Queued cli
              shared_ptr<QueuedCli> qcli =
                  QueuedCli::make(params->id, ttcli, params->qExecutor);
              return qcli;
            });
      };

  // create reconnecting cli that uses cliFactory to establish ssh connection
  shared_ptr<ReconnectingCli> rcli = ReconnectingCli::make(
      connectionParameters->id,
      connectionParameters->rExecutor,
      move(cliFactory),
      connectionParameters->timekeeper,
      connectionParameters->reconnectingQuietPeriod);
  // create keepalive cli
  shared_ptr<KeepaliveCli> kaCli = KeepaliveCli::make(
      connectionParameters->id,
      rcli,
      connectionParameters->kaExecutor,
      connectionParameters->timekeeper,
      connectionParameters->kaTimeout);
  return kaCli;
}

Future<shared_ptr<Cli>> IoConfigurationBuilder::configurePromptAwareCli(
    shared_ptr<PromptAwareCli> cli,
    shared_ptr<SshSessionAsync> session,
    shared_ptr<ConnectionParameters> params,
    shared_ptr<Executor> executor) {
  MLOG(MDEBUG) << "[" << params->id << "] "
               << "Initializing cli";
  // initialize CLI
  return cli->initializeCli(params->password)
      .via(executor.get())
      .thenValue([cli, params](auto) {
        // resolve prompt needs to happen
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "Resolving prompt";
        return cli->resolvePrompt();
      })
      .via(executor.get())
      .thenValue([cli, params, session](auto) {
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "Creating async data reader";
        event* sessionEvent = SshSocketReader::getInstance().addSshReader(
            readCallback, session->getSshFd(), session.get());
        session->setEvent(sessionEvent);
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "SSH layer configured";
        return makeFuture(cli);
      });
}

Future<shared_ptr<Cli>> IoConfigurationBuilder::createPromptAwareCli(
    shared_ptr<ConnectionParameters> params) {
  MLOG(MDEBUG) << "Creating CLI ssh device for " << params->id << " ("
               << params->ip << ":" << params->port << ")";

  // create session
  std::shared_ptr<SshSessionAsync> session = std::make_shared<SshSessionAsync>(
      params->id, params->sshExecutor, params->timekeeper);
  // open SSH connection

  MLOG(MDEBUG) << "[" << params->id << "] "
               << "Opening shell";

  return session
      ->openShell(
          params->ip,
          params->port,
          params->username,
          params->password,
          params->sshConnectionTimeout)
      .thenValue([params, session](auto) {
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "Creating shell";
        // create CLI
        shared_ptr<PromptAwareCli> cli = PromptAwareCli::make(
            params->id,
            session,
            params->flavour,
            params->paExecutor,
            params->timekeeper);
        return configurePromptAwareCli(
            cli, session, params, params->paExecutor);
      });
}

shared_ptr<IoConfigurationBuilder::ConnectionParameters>
IoConfigurationBuilder::makeConnectionParameters(
    string id,
    string hostname,
    string username,
    string password,
    string flavour,
    int port,
    chrono::seconds kaTimeout,
    chrono::seconds cmdTimeout,
    chrono::seconds reconnectingQuietPeriod,
    long sshConnectionTimeout,
    shared_ptr<CliThreadWheelTimekeeper> timekeeper,
    shared_ptr<Executor> sshExecutor,
    shared_ptr<Executor> paExecutor,
    shared_ptr<Executor> rcExecutor,
    shared_ptr<Executor> ttExecutor,
    shared_ptr<Executor> qExecutor,
    shared_ptr<Executor> rExecutor,
    shared_ptr<Executor> kaExecutor) {
  shared_ptr<IoConfigurationBuilder::ConnectionParameters>
      connectionParameters =
          make_shared<IoConfigurationBuilder::ConnectionParameters>();
  connectionParameters->id = id;
  connectionParameters->ip = hostname;
  connectionParameters->username = username;
  connectionParameters->password = password;
  connectionParameters->port = port;
  connectionParameters->flavour = CliFlavour::create(flavour);
  connectionParameters->kaTimeout = kaTimeout;
  connectionParameters->cmdTimeout = cmdTimeout;
  connectionParameters->reconnectingQuietPeriod = reconnectingQuietPeriod;
  connectionParameters->sshConnectionTimeout = sshConnectionTimeout;
  connectionParameters->timekeeper = timekeeper;
  connectionParameters->sshExecutor = sshExecutor;
  connectionParameters->paExecutor = paExecutor;
  connectionParameters->rcExecutor = rcExecutor;
  connectionParameters->ttExecutor = ttExecutor;
  connectionParameters->qExecutor = qExecutor;
  connectionParameters->rExecutor = rExecutor;
  connectionParameters->kaExecutor = kaExecutor;

  return connectionParameters;
}

string IoConfigurationBuilder::loadConfigValue(
    const std::map<std::string, std::string>& config,
    const string& key,
    const string& defaultValue) {
  return config.find(key) != config.end() ? config.at(key) : defaultValue;
}

} // namespace cli
} // namespace channels
} // namespace devmand
