/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <lte/protos/mconfig/mconfigs.pb.h>

#include <cstdlib>
#include <iostream>

#include "GrpcMagmaUtils.h"
#include "UpfMsgManageHandler.h"
#include "LocalEnforcer.h"
#include "magma_logging_init.h"
#include "includes/MagmaService.h"
#include "includes/MConfigLoader.h"
#include "OperationalStatesHandler.h"
#include "includes/PolicyLoader.h"
#include "RedisStoreClient.h"
#include "RestartHandler.h"
#include "SentryWrappers.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessionCredit.h"
#include "SessionManagerServer.h"
#include "SessionReporter.h"
#include "SessionStore.h"

#define SESSIOND_SERVICE "sessiond"
#define SESSION_PROXY_SERVICE "session_proxy"
#define POLICYDB_SERVICE "policydb"
#define SESSIOND_VERSION "1.0"
#define MIN_USAGE_REPORTING_THRESHOLD 0.4
#define MAX_USAGE_REPORTING_THRESHOLD 1.0
#define DEFAULT_USAGE_REPORTING_THRESHOLD 0.8
#define DEFAULT_QUOTA_EXHAUSTION_TERMINATION_MS 30000  // 30sec

#ifdef DEBUG
extern "C" void __gcov_flush(void);
#endif

static magma::mconfig::SessionD get_default_mconfig() {
  magma::mconfig::SessionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  mconfig.set_gx_gy_relay_enabled(false);
  auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
  wallet_config->set_terminate_on_exhaust(false);
  return mconfig;
}

static magma::mconfig::SessionD load_mconfig() {
  magma::mconfig::SessionD mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(SESSIOND_SERVICE, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for SessionD, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

static const std::shared_ptr<grpc::Channel> get_controller_channel(
    const YAML::Node& config, const bool gx_gy_relay_enabled) {
  if (gx_gy_relay_enabled) {
    MLOG(MINFO) << "Using proxied SessionD controller";
    return magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        SESSION_PROXY_SERVICE, magma::ServiceRegistrySingleton::CLOUD);
  } else {
    MLOG(MINFO) << "Using policydb controller";
    return magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        POLICYDB_SERVICE, magma::ServiceRegistrySingleton::LOCAL);
  }
}

static uint32_t get_log_verbosity(
    const YAML::Node& config, magma::mconfig::SessionD mconfig) {
  if (!config["log_level"].IsDefined()) {
    if (mconfig.log_level() < 0 || mconfig.log_level() > 4) {
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

void set_consts(const YAML::Node& config) {
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

  if (config["default_requested_units"].IsDefined()) {
    magma::SessionCredit::DEFAULT_REQUESTED_UNITS =
        config["default_requested_units"].as<uint64_t>();
  }

  if (config["send_access_timezone"].IsDefined()) {
    magma::LocalEnforcer::SEND_ACCESS_TIMEZONE =
        config["send_access_timezone"].as<bool>();
  }
  // default value for this config is true
  if (config["cleanup_all_dangling_flows"].IsDefined()) {
    magma::LocalEnforcer::CLEANUP_DANGLING_FLOWS =
        config["cleanup_all_dangling_flows"].as<bool>();
  }
  if (config["enable_ipfix"].IsDefined()) {
    magma::LocalEnforcer::SEND_IPFIX = config["enable_ipfix"].as<bool>();
  }

  // log all configs on startup
  MLOG(MINFO) << "==== Constants/Configs loaded from sessiond.yml ====";
  MLOG(MINFO) << "USAGE_REPORTING_THRESHOLD: "
              << magma::SessionCredit::USAGE_REPORTING_THRESHOLD;
  MLOG(MINFO) << "TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED: "
              << magma::SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED;
  MLOG(MINFO) << "DEFAULT_REQUESTED_UNITS: "
              << magma::SessionCredit::DEFAULT_REQUESTED_UNITS;
  MLOG(MINFO) << "SEND_ACCESS_TIMEZONE: "
              << magma::LocalEnforcer::SEND_ACCESS_TIMEZONE;
  MLOG(MINFO) << "CLEANUP_DANGLING_FLOWS: "
              << magma::LocalEnforcer::CLEANUP_DANGLING_FLOWS;
  MLOG(MINFO) << "SEND_IPFIX: " << magma::LocalEnforcer::SEND_IPFIX;
  MLOG(MINFO) << "==== Constants/Configs loaded from sessiond.yml ====";
}

magma::SessionStore* create_session_store(
    const YAML::Node& config,
    std::shared_ptr<magma::StaticRuleStore> rule_store,
    std::shared_ptr<magma::MeteringReporter> metering_reporter) {
  bool is_stateless = config["support_stateless"].IsDefined() &&
                      config["support_stateless"].as<bool>();
  if (is_stateless) {
    auto store_client = std::make_shared<magma::lte::RedisStoreClient>(
        std::make_shared<cpp_redis::client>(),
        config["sessions_table"].as<std::string>(), rule_store);
    bool connected;
    do {
      MLOG(MINFO) << "Attempting to connect to Redis";
      connected = store_client->try_redis_connect();
      std::this_thread::sleep_for(std::chrono::milliseconds(100));
    } while (!connected);
    MLOG(MINFO) << "Successfully connected to Redis";
    return new magma::SessionStore(rule_store, metering_reporter, store_client);
  } else {
    MLOG(MINFO) << "Session store in memory";
    return new magma::SessionStore(rule_store, metering_reporter);
  }
}

long get_quota_exhaust_termination_time(const YAML::Node& config) {
  long quota_exhaust_termination_on_init_ms;
  if (config["cwf_quota_exhaustion_termination_on_init_ms"].IsDefined()) {
    quota_exhaust_termination_on_init_ms =
        config["cwf_quota_exhaustion_termination_on_init_ms"].as<long>();
  } else {
    quota_exhaust_termination_on_init_ms =
        DEFAULT_QUOTA_EXHAUSTION_TERMINATION_MS;
  }
  return quota_exhaust_termination_on_init_ms;
}

int main(int argc, char* argv[]) {
#ifdef DEBUG
  __gcov_flush();
#endif

  magma::init_logging(argv[0]);

  auto mconfig = load_mconfig();
  auto config =
      magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
  magma::set_verbosity(get_log_verbosity(config, mconfig));

  if ((config["print_grpc_payload"].IsDefined())) {
    set_grpc_logging_level(config["print_grpc_payload"].as<bool>());
  }

  initialize_sentry();

  bool converged_access = false;
  // Check converged sessiond is enabled or not
  if ((config["converged_access"].IsDefined()) &&
      (config["converged_access"].as<bool>())) {
    converged_access = true;
  }
  MLOG(MINFO) << "Starting Session Manager";
  folly::EventBase* evb = folly::EventBaseManager::get()->getEventBase();

  // Start off a thread to periodically load policy definitions from Redis into
  // RuleStore
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
  std::thread pipelined_response_handling_thread([&]() {
    MLOG(MINFO) << "Started PipelineD response thread";
    pipelined_client->rpc_response_loop();
  });

  auto directoryd_client = std::make_shared<magma::AsyncDirectorydClient>();
  std::thread directoryd_response_handling_thread([&]() {
    MLOG(MINFO) << "Started DirectoryD response thread";
    directoryd_client->rpc_response_loop();
  });

  auto& eventd_client = magma::AsyncEventdClient::getInstance();
  auto events_reporter =
      std::make_shared<magma::lte::EventsReporterImpl>(eventd_client);
  std::thread eventd_response_handling_thread([&]() {
    MLOG(MINFO) << "Started EventD response thread";
    eventd_client.rpc_response_loop();
  });

  std::shared_ptr<magma::AsyncSpgwServiceClient> spgw_client;
  std::shared_ptr<aaa::AsyncAAAClient> aaa_client;
  std::shared_ptr<magma::AsyncAmfServiceClient> amf_srv_client;

  if (converged_access) {
    // AMF service client to handle response message
    amf_srv_client = std::make_shared<magma::AsyncAmfServiceClient>();
    spgw_client    = nullptr;
    aaa_client     = nullptr;
  }
  // Case on config, setup the appropriate client for the access component
  std::thread access_response_handling_thread;
  if (config["support_carrier_wifi"].as<bool>()) {
    aaa_client                      = std::make_shared<aaa::AsyncAAAClient>();
    access_response_handling_thread = std::thread([&]() {
      MLOG(MINFO) << "Started AAA Client response thread";
      aaa_client->rpc_response_loop();
    });
    spgw_client                     = nullptr;
    amf_srv_client                  = nullptr;
  } else {
    spgw_client = std::make_shared<magma::AsyncSpgwServiceClient>();
    access_response_handling_thread = std::thread([&]() {
      MLOG(MINFO) << "Started SPGW response thread";
      spgw_client->rpc_response_loop();
    });
    aaa_client                      = nullptr;
  }

  // Setup SessionReporter which talks to the policy component
  // (FeG+PCRF/PolicyDB).
  bool gx_gy_relay_enabled = mconfig.gx_gy_relay_enabled();
  auto reporter            = std::make_shared<magma::SessionReporterImpl>(
      evb, get_controller_channel(config, gx_gy_relay_enabled));
  std::thread policy_response_handler([&]() {
    MLOG(MINFO) << "Started reporter thread";
    reporter->rpc_response_loop();
  });

  // Case on stateless config, setup the appropriate store client
  auto metering_reporter = std::make_shared<magma::MeteringReporter>();
  magma::SessionStore* session_store =
      create_session_store(config, rule_store, metering_reporter);
  // service restart clears the UE metering metrics, so we need to offset
  // metering_reporter with existing usage
  session_store->initialize_metering_counter();

  // Some setup work for the SessionCredit class
  set_consts(config);
  // Initialize the main logical component of SessionD
  auto local_enforcer = std::make_shared<magma::LocalEnforcer>(
      reporter, rule_store, *session_store, pipelined_client, events_reporter,
      spgw_client, aaa_client,
      config["session_force_termination_timeout_ms"].as<long>(),
      get_quota_exhaust_termination_time(config), mconfig);

  // RestartHandler will cleanup sessions from previous SessionD run. We do not
  // care about the return value of this thread.
  auto restart_handler = std::make_shared<magma::sessiond::RestartHandler>(
      directoryd_client, aaa_client, local_enforcer, reporter.get(),
      *session_store);
  std::thread restart_handler_thread([&]() {
    MLOG(MINFO) << "Started SessionD restart handler thread";
    bool is_stateless = config["support_stateless"].IsDefined() &&
                        config["support_stateless"].as<bool>();
    if (!is_stateless) {
      restart_handler->cleanup_previous_sessions();
    } else if (config["support_carrier_wifi"].as<bool>()) {
      restart_handler->setup_aaa_sessions();
    }
  });

  // Setup threads to serve as GRPC servers for the LocalSessionManagerHandler
  // and the SessionProxyHandler (RARs)
  auto local_handler = std::make_unique<magma::LocalSessionManagerHandlerImpl>(
      local_enforcer, reporter.get(), directoryd_client, events_reporter,
      *session_store);
  auto proxy_handler =
      std::make_shared<magma::SessionProxyResponderHandlerImpl>(
          local_enforcer, *session_store);
  magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);
  magma::LocalSessionManagerAsyncService local_service(
      server.GetNewCompletionQueue(), std::move(local_handler));
  magma::SessionProxyResponderAsyncService proxy_service(
      server.GetNewCompletionQueue(), proxy_handler);
  server.AddServiceToServer(&local_service);
  MLOG(MINFO) << "Added LocalSessionManagerAsyncService to service's server";
  server.AddServiceToServer(&proxy_service);
  MLOG(MINFO) << "Added SessionProxyResponderAsyncService to service's server";

  // Register state polling callback
  server.SetOperationalStatesCallback([evb, session_store]() {
    std::promise<magma::OpState> result;
    std::future<magma::OpState> future = result.get_future();
    evb->runInEventBaseThread([session_store, &result, &future]() {
      set_sentry_transaction("GetOperationalStates");
      result.set_value(magma::get_operational_states(session_store));
    });
    return future.get();
  });

  magma::AmfPduSessionSmContextAsyncService* conv_set_message_service = nullptr;
  magma::SetInterfaceForUserPlaneAsyncService* conv_upf_message_service =
      nullptr;
  if (converged_access) {
    // Initialize the main thread of session management by folly event to handle
    // logical component of 5G of SessionD
    extern std::shared_ptr<magma::SessionStateEnforcer> conv_session_enforcer;
    conv_session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pipelined_client, amf_srv_client, mconfig,
        config["session_force_termination_timeout_ms"].as<long>());
    // 5G related async msg handler service framework creation
    auto conv_set_message_handler =
        std::make_unique<magma::SetMessageManagerHandler>(
            conv_session_enforcer, *session_store);
    MLOG(MINFO) << "Initialized SetMessageManagerHandler";
    // 5G specific services to handle set messages from AMF and mme
    conv_set_message_service = new magma::AmfPduSessionSmContextAsyncService(
        server.GetNewCompletionQueue(), std::move(conv_set_message_handler));
    // 5G related services
    server.AddServiceToServer(conv_set_message_service);
    MLOG(MINFO)
        << "Added SessionProxyResponderAsyncService to service's server";

    // 5G related upf  async service framework creation
    auto conv_upf_message_handler =
        std::make_unique<magma::UpfMsgManageHandler>(
            conv_session_enforcer, *session_store);
    // 5G  upf converged service to handler set message from UPF
    conv_upf_message_service = new magma::SetInterfaceForUserPlaneAsyncService(
        server.GetNewCompletionQueue(), std::move(conv_upf_message_handler));
    MLOG(MINFO) << "SetInterfaceForUserPlaneAsyncService ";
    server.AddServiceToServer(conv_upf_message_service);
    MLOG(MINFO) << "Add converged UPF message service";
    // 5G related SessionStateEnforcer main thread start to handled session
    // state
    conv_session_enforcer->attachEventBase(evb);
  }

  // For FWA always handle abort session
  magma::AbortSessionResponderAsyncService* abort_session_service = nullptr;
  if (!config["support_carrier_wifi"].as<bool>()) {
    abort_session_service = new magma::AbortSessionResponderAsyncService(
        server.GetNewCompletionQueue(), proxy_handler);
    server.AddServiceToServer(abort_session_service);
  }

  server.Start();

  // 5G set message handling thread from access.
  std::thread access_common_message_thread([&]() {
    // conv_set_message_service is initialized only if it is converged_access
    if (converged_access) {
      MLOG(MDEBUG) << "Started access message thread";
      conv_set_message_service
          ->wait_for_requests();         // block here instead of on server
      conv_set_message_service->stop();  // stop queue after server shutsdown
    }
  });
  std::thread conv_upf_message_thread([&]() {
    if (converged_access) {
      MLOG(MINFO) << "Started upf message thread";
      conv_upf_message_service
          ->wait_for_requests();         // block here instead of on server
      conv_upf_message_service->stop();  // stop queue after server shutsdown
    }
  });
  // session_enforcer->sync_sessions_on_restart(time(NULL));//not part of drop-1
  // MVC
  std::thread local_thread([&]() {
    MLOG(MINFO) << "Started local service thread";
    local_service.wait_for_requests();  // block here instead of on server
    local_service.stop();               // stop queue after server shuts down
  });
  std::thread proxy_thread([&]() {
    MLOG(MINFO) << "Started proxy service thread";
    proxy_service.wait_for_requests();  // block here instead of on server
    proxy_service.stop();               // stop queue after server shuts down
  });

  // Only start abort session handler for non-CWF deployments
  std::thread abort_session_thread;
  if (abort_session_service != nullptr) {
    abort_session_thread = std::thread([&]() {
      MLOG(MINFO) << "Started abort session service thread";
      // block here instead of on server
      abort_session_service->wait_for_requests();
      abort_session_service->stop();  // stop queue after server shuts down
    });
  }

  // Block on main local_enforcer (to keep evb in this thread)
  local_enforcer->attachEventBase(evb);
  MLOG(MDEBUG) << "local enforcer Attached EventBase to evb";
  local_enforcer->sync_sessions_on_restart(time(NULL));
  MLOG(MDEBUG) << "Synced session on restart";
  evb->loopForever();
  MLOG(MINFO) << "Stoping.. session manager GRPC server";

  server.Stop();

  // Clean up threads & resources
  policy_response_handler.join();
  local_thread.join();
  proxy_thread.join();
  pipelined_response_handling_thread.join();
  directoryd_response_handling_thread.join();
  restart_handler_thread.join();
  policy_loader_thread.join();
  if (abort_session_service != nullptr) {
    abort_session_thread.join();
    free(abort_session_service);
  }
  access_response_handling_thread.join();
  if (converged_access) {
    // 5G related thread join
    access_common_message_thread.join();
    conv_upf_message_thread.join();
    free(conv_set_message_service);
    free(conv_upf_message_service);
  }
  delete session_store;

  shutdown_sentry();
  return 0;
}
