/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <iostream>

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

#define SESSIOND_SERVICE "sessiond"
#define SESSION_PROXY_SERVICE "session_proxy"
#define SESSIOND_VERSION "1.0"
#define MIN_USAGE_REPORTING_THRESHOLD 0.4
#define MAX_USAGE_REPORTING_THRESHOLD 1.1
#define DEFAULT_USAGE_REPORTING_THRESHOLD 0.8

#ifdef DEBUG
extern "C" void __gcov_flush(void);
#endif

static magma::mconfig::SessionD get_default_mconfig()
{
  magma::mconfig::SessionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  mconfig.set_relay_enabled(false);
  return mconfig;
}

static magma::mconfig::SessionD load_mconfig()
{
  magma::mconfig::SessionD mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(SESSIOND_SERVICE, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for sessiond, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

static void run_bare_service303()
{
  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
  server.Start();
  server.WaitForShutdown(); // blocks forever
  server.Stop();
}

static bool sessiond_enabled(const magma::mconfig::SessionD &mconfig)
{
  return mconfig.relay_enabled();
}

static const std::shared_ptr<grpc::Channel> get_controller_channel(
  const YAML::Node &config)
{
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
    addr, grpc::InsecureChannelCredentials(), grpc::ChannelArguments {});
}

static uint32_t get_log_verbosity(const YAML::Node &config)
{
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

int main(int argc, char *argv[])
{
#ifdef DEBUG
  __gcov_flush();
#endif

  MLOG(MINFO) << "Starting Session Manager";
  magma::init_logging(argv[0]);
  auto mconfig = load_mconfig();
  auto config =
    magma::ServiceConfigLoader {}.load_service_config(SESSIOND_SERVICE);
  magma::set_verbosity(get_log_verbosity(config));

  if (!sessiond_enabled(mconfig)) {
    MLOG(MINFO) << "Credit control disabled, local enforcer not running";
    run_bare_service303();
    return 0;
  }

  folly::EventBase *evb = folly::EventBaseManager::get()->getEventBase();

  // prep rule manager and rule update loop
  auto rule_store = std::make_shared<magma::StaticRuleStore>();
  magma::PolicyLoader policy_loader;
  std::thread policy_loader_thread([&]() {
    policy_loader.start_loop(
      [&](std::vector<magma::PolicyRule> rules) {
        rule_store->sync_rules(rules);
      },
      config["rule_update_inteval_sec"].as<uint32_t>());
    policy_loader.stop();
  });

  auto pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
  std::thread rule_manager_thread([&]() {
    MLOG(MINFO) << "Started pipelined response thread";
    pipelined_client->rpc_response_loop();
  });

  std::shared_ptr<magma::AsyncSpgwServiceClient> spgw_client;
  std::shared_ptr<aaa::AsyncAAAClient> aaa_client;

  std::thread optional_client_thread;
  if (config["support_carrier_wifi"].as<bool>()) {
    aaa_client = std::make_shared<aaa::AsyncAAAClient>();
    optional_client_thread = std::thread([&]() {
      MLOG(MINFO) << "Started AAA Client response thread";
      aaa_client->rpc_response_loop();
    });

    spgw_client = nullptr;
  } else {
    aaa_client = nullptr;

    spgw_client = std::make_shared<magma::AsyncSpgwServiceClient>();
    optional_client_thread = std::thread([&]() {
      MLOG(MINFO) << "Started SPGW response thread";
      spgw_client->rpc_response_loop();
    });
  }

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

  auto reporter = std::make_shared<magma::SessionCloudReporterImpl>(
    evb, get_controller_channel(config));
  std::thread reporter_thread([&]() {
    MLOG(MINFO) << "Started reporter thread";
    reporter->rpc_response_loop();
  });

  auto monitor = magma::LocalEnforcer(
    reporter,
    rule_store,
    pipelined_client,
    spgw_client,
    aaa_client,
    config["session_force_termination_timeout_ms"].as<long>());

  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
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

  std::thread local_thread([&]() {
    MLOG(MINFO) << "Started local service thread";
    local_service.wait_for_requests(); // block here instead of on server
    local_service.stop();              // stop queue after server shuts down
  });
  std::thread proxy_thread([&]() {
    MLOG(MINFO) << "Started proxy service thread";
    proxy_service.wait_for_requests(); // block here instead of on server
    proxy_service.stop();              // stop queue after server shuts down
  });

  // Block on main monitor (to keep evb in this thread)
  monitor.attachEventBase(evb);
  monitor.start();
  server.Stop();

  reporter_thread.join();
  local_thread.join();
  proxy_thread.join();
  rule_manager_thread.join();
  policy_loader_thread.join();
  optional_client_thread.join();

  return 0;
}
