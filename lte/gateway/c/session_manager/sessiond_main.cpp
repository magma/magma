/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <iostream>
#include <string>

#include <lte/protos/mconfig/mconfigs.pb.h>

#include "SessionManagerServer.h"
#include "LocalEnforcer.h"
#include "CloudReporter.h"
#include "MagmaService.h"
#include "ServiceRegistrySingleton.h"
#include "PolicyLoader.h"
#include "MConfigLoader.h"
#include "magma_logging.h"
#include "SessionCredit.h"

#ifdef DEBUG
extern "C" void __gcov_flush(void);
#endif

namespace {

const char *SESSIOND_SERVICE = "sessiond";
const char *SESSION_PROXY_SERVICE = "session_proxy";
const char *SESSIOND_VERSION = "1.0";
const double MIN_USAGE_REPORTING_THRESHOLD = 0.4;
const double MAX_USAGE_REPORTING_THRESHOLD = 1.1;
const double DEFAULT_USAGE_REPORTING_THRESHOLD = 0.8;

magma::mconfig::SessionD get_default_mconfig() {
  magma::mconfig::SessionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  mconfig.set_relay_enabled(false);
  return mconfig;
}

magma::mconfig::SessionD load_mconfig() {
  magma::mconfig::SessionD mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(SESSIOND_SERVICE, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for sessiond, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

void run_bare_service303() {
  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
  server.Start();
  server.WaitForShutdown(); // blocks forever
  server.Stop();
}

bool sessiond_enabled(const magma::mconfig::SessionD &mconfig) {
  return mconfig.relay_enabled();
}

std::shared_ptr<grpc::Channel> get_controller_channel(
    const YAML::Node &config) {
  if (
      !config["use_proxied_controller"].IsDefined() ||
          config["use_proxied_controller"].as<bool>()) {
    MLOG(MINFO) << "Using proxied sessiond controller";
    return magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        SESSION_PROXY_SERVICE, magma::ServiceRegistrySingleton::CLOUD);
  }
  auto port = config["local_controller_port"].as<std::string>();
  auto addr = "127.0.0.1:" + port;
  MLOG(MINFO) << "Using local address " << addr << " for controller";
  return grpc::CreateCustomChannel(
      addr, grpc::InsecureChannelCredentials(), grpc::ChannelArguments{});
}

uint32_t get_log_verbosity(const YAML::Node &config) {
  if (!config["log_level"].IsDefined()) {
    return MINFO;
  }
  std::string log_level = config["log_level"].as<std::string>();
  if (log_level == "DEBUG") {
    return MDEBUG;
  } else if (log_level == "INFO") {
    return MINFO;
  } else if (log_level == "WARNING") {
    return MWARNING;
  } else if (log_level == "ERROR") {
    return MERROR;
  } else if (log_level == "FATAL") {
    return MFATAL;
  } else {
    MLOG(MINFO) << "Invalid log level in config: "
                << config["log_level"].as<std::string>();
    return MINFO;
  }
}

void init_rules(const YAML::Node &config,
                std::vector<std::thread> &threads,
                std::shared_ptr<magma::StaticRuleStore> &rule_store) {
  magma::PolicyLoader policy_loader{};
  threads.emplace_back([&]() {
    policy_loader.start_loop(
        [&](const std::vector<magma::PolicyRule> &rules) {
          rule_store->sync_rules(rules);
        },
        config["rule_update_inteval_sec"].as<uint32_t>());
    policy_loader.stop();
  });
}

void init_pipelined_client(std::vector<std::thread> &threads,
    std::shared_ptr<magma::AsyncPipelinedClient> &pipelined_client) {
  pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
  threads.emplace_back([&]() {
    MLOG(MINFO) << "Started pipelined response thread";
    pipelined_client->rpc_response_loop();
  });
}

void init_cwf(std::vector<std::thread> &threads,
              std::shared_ptr<aaa::AsyncAAAClient> &aaa_client) {
  aaa_client = std::make_shared<aaa::AsyncAAAClient>();
  threads.emplace_back(([&]() {
    MLOG(MINFO) << "Started AAA Client response thread";
    aaa_client->rpc_response_loop();
  }));
}

void init_lte(std::vector<std::thread> &threads,
              std::shared_ptr<magma::AsyncSpgwServiceClient> &spgw_client) {
  spgw_client = std::make_shared<magma::AsyncSpgwServiceClient>();
  threads.emplace_back([&]() {
    MLOG(MINFO) << "Started SPGW response thread";
    spgw_client->rpc_response_loop();
  });
}

void init_reporting(const YAML::Node &config,
    folly::EventBase *evb,
    std::vector<std::thread>& threads,
    std::shared_ptr<magma::SessionCloudReporterImpl> &reporter) {
  auto reporting_threshold = config["usage_reporting_threshold"].as<float>();
  if (reporting_threshold <= MIN_USAGE_REPORTING_THRESHOLD ||
      reporting_threshold >= MAX_USAGE_REPORTING_THRESHOLD) {
    MLOG(MWARNING) << "Usage reporting threshold should be between "
                   << MIN_USAGE_REPORTING_THRESHOLD << " and "
                   << MAX_USAGE_REPORTING_THRESHOLD << ", apply default value: "
                   << DEFAULT_USAGE_REPORTING_THRESHOLD;
    reporting_threshold = DEFAULT_USAGE_REPORTING_THRESHOLD;
  }
  magma::SessionCredit::USAGE_REPORTING_THRESHOLD = reporting_threshold;
  magma::SessionCredit::EXTRA_QUOTA_MARGIN =
      config["extra_quota_margin"].as<uint64_t>();
  magma::SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED =
      config["terminate_service_when_quota_exhausted"].as<bool>();

  reporter = std::make_shared<magma::SessionCloudReporterImpl>(
      evb, get_controller_channel(config));
  threads.emplace_back([&]() {
    MLOG(MINFO) << "Started reporter thread";
    reporter->rpc_response_loop();
  });
}

void init_async_server(std::vector<std::thread> &threads,
    magma::LocalEnforcer &monitor,
    std::shared_ptr<magma::SessionCloudReporterImpl> &reporter,
    magma::service303::MagmaService &server) {
  auto local_handler = std::make_unique<magma::LocalSessionManagerHandlerImpl>(
      &monitor, reporter.get());
  auto proxy_handler =
      std::make_unique<magma::SessionProxyResponderHandlerImpl>(&monitor);

  magma::LocalSessionManagerAsyncService local_service(
      server.GetNewCompletionQueue(), std::move(local_handler));
  magma::SessionProxyResponderAsyncService proxy_service(
      server.GetNewCompletionQueue(), std::move(proxy_handler));
  server.AddServiceToServer(&local_service);
  server.AddServiceToServer(&proxy_service);
  server.Start();

  threads.emplace_back([&]() {
    MLOG(MINFO) << "Started local service thread";
    local_service.wait_for_requests(); // block here instead of on server
    local_service.stop();              // stop queue after server shuts down
  });
  threads.emplace_back([&]() {
    MLOG(MINFO) << "Started proxy service thread";
    proxy_service.wait_for_requests(); // block here instead of on server
    proxy_service.stop();              // stop queue after server shuts down
  });
}
}  // anonymous

/* Sessiond is the control plane for the PCEF.
 * Sessiond consists of the following.
 * An async GRPC server that implements three services
 * 1 - LocalSessionManager that handles within the gateway call (e.g. MME)
 * 2 - Proxy responder for servicing requests from the FeG
 * 3 - Service303
 * Service303 runs in the main thread, a completion queue for the other services
 * run on dedicated threads of their own.
 * Additionally the process has:
 * - An async client to the pipelined enforcement service for pushing rules
 * to the PCEF
 * - An async client to the MME/AAA to perform network initiated detaches
 * - An async client to communicate towards the cloud for PCRF and OCS
 *   related updates.
 *   TODO: Service needs a signal handler.
 */
int main(int argc, char *argv[]) {
#ifdef DEBUG
  __gcov_flush();
#endif

  MLOG(MINFO) << "Starting Session Manager";
  magma::init_logging(argv[0]);
  auto mconfig = load_mconfig();
  auto config =
      magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
  magma::set_verbosity(get_log_verbosity(config));

  if (!sessiond_enabled(mconfig)) {
    MLOG(MINFO) << "Credit control disabled, local enforcer not running";
    run_bare_service303();
    return 0;
  }

  std::vector<std::thread> threads;

  // prep rule manager and rule update loop
  std::shared_ptr<magma::StaticRuleStore> rule_store;
  init_rules(config, threads, rule_store);

  std::shared_ptr<magma::AsyncPipelinedClient> pipelined_client;
  init_pipelined_client(threads, pipelined_client);

  // init connection to AAA/MME.
  // TODO: Ideally this should be abstracted behind a single client interface
  std::shared_ptr<magma::AsyncSpgwServiceClient> spgw_client;
  std::shared_ptr<aaa::AsyncAAAClient> aaa_client;
  if (config["support_carrier_wifi"].as<bool>()) {
    init_cwf(threads, aaa_client);
  } else {
    init_lte(threads, spgw_client);
  }

  // init reporting
  std::shared_ptr<magma::SessionCloudReporterImpl> reporter;
  folly::EventBase *evb = folly::EventBaseManager::get()->getEventBase();
  init_reporting(config, evb, threads, reporter);


  // Create the local enforcer and
  auto monitor = magma::LocalEnforcer(
      reporter,
      rule_store,
      pipelined_client,
      spgw_client,
      aaa_client,
      config["session_force_termination_timeout_ms"].as<long>());

  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
  init_async_server(threads, monitor, reporter, server);

  // Block on main monitor (to keep evb in this thread)
  monitor.attachEventBase(evb);
  monitor.start();
  server.Stop();

  for (auto &thread : threads) {
    thread.join();
  }

  return 0;
}
