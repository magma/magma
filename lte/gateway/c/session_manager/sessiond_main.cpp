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
#include "SessionReporter.h"
#include "MagmaService.h"
#include "RedisStoreClient.h"
#include "RestartHandler.h"
#include "ServiceRegistrySingleton.h"
#include "PolicyLoader.h"
#include "MConfigLoader.h"
#include "magma_logging.h"
#include "SessionCredit.h"
#include "SessionStore.h"

#define SESSIOND_SERVICE "sessiond"
#define SESSION_PROXY_SERVICE "session_proxy"
#define POLICYDB_SERVICE "policydb"
#define SESSIOND_VERSION "1.0"
#define MIN_USAGE_REPORTING_THRESHOLD 0.4
#define MAX_USAGE_REPORTING_THRESHOLD 1.1
#define DEFAULT_USAGE_REPORTING_THRESHOLD 0.8
#define DEFAULT_QUOTA_EXHAUSTION_TERMINATION_MS 30000 // 30sec

#ifdef DEBUG
extern "C" void __gcov_flush(void);
#endif

static magma::mconfig::SessionD get_default_mconfig()
{
  magma::mconfig::SessionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  mconfig.set_relay_enabled(false);
  auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
  wallet_config->set_terminate_on_exhaust(false);
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

static const std::shared_ptr<grpc::Channel> get_controller_channel(
  const YAML::Node &config, const bool relay_enabled)
{
  if (relay_enabled) {
    MLOG(MINFO) << "Using proxied sessiond controller";
    return magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      SESSION_PROXY_SERVICE, magma::ServiceRegistrySingleton::CLOUD);
  } else {
    MLOG(MINFO) << "Using policydb controller";
    return magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      POLICYDB_SERVICE, magma::ServiceRegistrySingleton::LOCAL);
  }
}

static uint32_t get_log_verbosity(const YAML::Node &config, magma::mconfig::SessionD mconfig)
{
  if (!config["log_level"].IsDefined()) {
    if(mconfig.log_level() < 0 || mconfig.log_level() > 4){
      return MINFO;
    }
    return mconfig.log_level();
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
    magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
  magma::set_verbosity(get_log_verbosity(config, mconfig));

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

  auto directoryd_client = std::make_shared<magma::AsyncDirectorydClient>();
  std::thread directoryd_thread([&]() {
    MLOG(MINFO) << "Started pipelined response thread";
    directoryd_client->rpc_response_loop();
  });

  auto& eventd_client = magma::AsyncEventdClient::getInstance();
  std::thread eventd_thread([&]() {
    MLOG(MINFO) << "Started eventd response thread";
    eventd_client.rpc_response_loop();
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

  magma::SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED =
   config["terminate_service_when_quota_exhausted"].as<bool>();

  auto controller_channel = get_controller_channel(config,
    mconfig.relay_enabled());
  auto reporter = std::make_shared<magma::SessionReporterImpl>(
    evb, controller_channel);
  std::thread reporter_thread([&]() {
    MLOG(MINFO) << "Started reporter thread";
    reporter->rpc_response_loop();
  });

  // [CWF-ONLY]
  long quota_exhaust_termination_on_init_ms;
  if (config["cwf_quota_exhaustion_termination_on_init_ms"].IsDefined()) {
    quota_exhaust_termination_on_init_ms =
      config["cwf_quota_exhaustion_termination_on_init_ms"].as<long>();
  } else {
    quota_exhaust_termination_on_init_ms =
      DEFAULT_QUOTA_EXHAUSTION_TERMINATION_MS;
  }

  magma::SessionStore* session_store;
  bool is_stateless = config["support_stateless"].IsDefined()
      && config["support_stateless"].as<bool>();
  if (is_stateless) {
    auto store_client = std::make_shared<magma::lte::RedisStoreClient>(
        std::make_shared<cpp_redis::client>(),
        config["sessions_table"].as<std::string>(),
        rule_store);
    bool connected;
    do {
      MLOG(MINFO) << "Attempting to connect to Redis";
      connected = store_client->try_redis_connect();
      std::this_thread::sleep_for(std::chrono::milliseconds(100));
    } while (!connected);
    session_store = new magma::SessionStore(rule_store, store_client);
    MLOG(MINFO) << "Successfully connected to Redis";
  } else {
    MLOG(MINFO) << "Session store in memory";
    session_store = new magma::SessionStore(rule_store);
  }
  auto monitor = std::make_shared<magma::LocalEnforcer>(
    reporter,
    rule_store,
    *session_store,
    pipelined_client,
    directoryd_client,
    eventd_client,
    spgw_client,
    aaa_client,
    config["session_force_termination_timeout_ms"].as<long>(),
    quota_exhaust_termination_on_init_ms,
    mconfig);

  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
  auto local_handler = std::make_unique<magma::LocalSessionManagerHandlerImpl>(
    monitor, reporter.get(), directoryd_client, *session_store);
  auto proxy_handler =std::make_unique<magma::SessionProxyResponderHandlerImpl>(
    monitor, *session_store);

  auto restart_handler = std::make_shared<magma::sessiond::RestartHandler>(
      directoryd_client, aaa_client, monitor, reporter.get(), *session_store);
  std::thread restart_handler_thread([&]() {
    MLOG(MINFO) << "Started sessiond restart handler thread";
    if (!is_stateless) {
      restart_handler->cleanup_previous_sessions();
    } else if (config["support_carrier_wifi"].as<bool>()) {
      restart_handler->setup_aaa_sessions();
    }
  });

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
  monitor->attachEventBase(evb);
  monitor->sync_sessions_on_restart(time(NULL));
  monitor->start();
  server.Stop();

  reporter_thread.join();
  local_thread.join();
  proxy_thread.join();
  rule_manager_thread.join();
  directoryd_thread.join();
  restart_handler_thread.join();
  policy_loader_thread.join();
  optional_client_thread.join();

  return 0;
}
