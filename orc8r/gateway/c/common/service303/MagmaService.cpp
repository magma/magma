/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <string>
#include <csignal>
#include <ctime>
#include <chrono>
#include <ratio>

#include <iostream>

#include <orc8r/protos/service303.grpc.pb.h>
#include <orc8r/protos/service303.pb.h>
#include <orc8r/protos/common.pb.h>

#include "MagmaService.h"
#include "MetricsRegistry.h"
#include "MetricsSingleton.h"
#include "ProcFileUtils.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Channel;
using grpc::ServerContext;
using grpc::Status;
using grpc::ServerBuilder;
using grpc::InsecureServerCredentials;
using grpc::Server;
using magma::orc8r::Service303;
using magma::orc8r::ServiceInfo;
using magma::orc8r::State;
using magma::orc8r::GetOperationalStatesResponse;
using magma::orc8r::ReloadConfigResponse;
using magma::orc8r::Void;
using magma::service303::MetricsSingleton;
using magma::service303::MagmaService;
using io::prometheus::client::MetricFamily;
using namespace std::chrono;

MagmaService::MagmaService(const std::string& name, const std::string& version)
    : name_(name), version_(version), health_(ServiceInfo::APP_UNKNOWN),
      start_time_(steady_clock::now()), wall_start_time_(system_clock::now()),
      service_info_callback_(nullptr), config_reload_callback_(nullptr),
      operational_states_callback_(nullptr)
      {}

void MagmaService::AddServiceToServer(grpc::Service *service) {
  builder_.RegisterService(service);
}

std::unique_ptr<grpc::ServerCompletionQueue>
MagmaService::GetNewCompletionQueue() {
  return std::move(builder_.AddCompletionQueue());
}

void MagmaService::Start() {
    setMetricsStartTime();
    builder_.RegisterService(this);
    std::string service_addr = magma::ServiceRegistrySingleton::Instance()
      ->GetServiceAddrString(name_);
    builder_.AddListeningPort(service_addr,
                              grpc::InsecureServerCredentials());
    server_ = builder_.BuildAndStart();
}

void MagmaService::WaitForShutdown() {
  server_->Wait(); // Blocking call
}

void MagmaService::Stop() {
  server_->Shutdown();
}

void MagmaService::SetServiceInfoCallback(ServiceInfoCallback callback) {
  service_info_callback_ = callback;
}

void MagmaService::ClearServiceInfoCallback() {
  service_info_callback_ = nullptr;
}

void MagmaService::SetConfigReloadCallback(ConfigReloadCallback callback) {
  config_reload_callback_ = callback;
}

void MagmaService::ClearConfigReloadCallback() {
  config_reload_callback_ = nullptr;
}

void MagmaService::SetOperationalStatesCallback(
  OperationalStatesCallback callback
) {
  operational_states_callback_ = callback;
}

void MagmaService::ClearOperationalStatesCallback() {
  operational_states_callback_ = nullptr;
}

Status MagmaService::GetServiceInfo(
    ServerContext* context, const Void* request, ServiceInfo* response) {
  auto start_time_secs =
    time_point_cast<seconds>(wall_start_time_).time_since_epoch().count();

  auto meta = (service_info_callback_ != nullptr) ?
    service_info_callback_() :
    std::map<std::string, std::string>();

  response->set_name(name_);
  response->set_version(version_);
  response->set_state(ServiceInfo::ALIVE);
  response->set_health(health_);
  response->set_start_time_secs(start_time_secs);
  response->mutable_status()->mutable_meta()->insert(meta.begin(), meta.end());

  return Status::OK;
}

Status MagmaService::StopService(
    ServerContext* context, const Void* request, Void* response) {
  std::raise(SIGTERM);
  return Status::OK;
}

Status MagmaService::GetMetrics(
    ServerContext* context, const Void* request, MetricsContainer* response) {
  // Set all common metrics
  setSharedMetrics();

  MetricsSingleton& instance = MetricsSingleton::Instance();
  const std::vector<MetricFamily>& collected = instance.registry_->Collect();
  for (auto it = collected.begin(); it != collected.end(); it++) {
    MetricFamily* family = response->add_family();
    family->CopyFrom(*it);
  }
  return Status::OK;
}

Status MagmaService::SetLogLevel(
    ServerContext* context,
    const LogLevelMessage* request,
    Void* response) {
  // log level FATAL is minimum verbosity and maximum level
  auto verbosity = LogLevel::FATAL - request->level();
  set_verbosity(verbosity);
  return Status::OK;
}

Status MagmaService::ReloadServiceConfig(
    ServerContext *context,
    const Void *request,
    ReloadConfigResponse *response) {
  if (config_reload_callback_ != nullptr) {
    if (config_reload_callback_()) {
      response->set_result(ReloadConfigResponse::RELOAD_SUCCESS);
    } else {
      response->set_result(ReloadConfigResponse::RELOAD_FAILURE);
    }
  } else {
    response->set_result(ReloadConfigResponse::RELOAD_UNSUPPORTED);
  }

  return Status::OK;
}

Status MagmaService::GetOperationalStates(
    ServerContext *context,
    const Void *request,
    GetOperationalStatesResponse *response) {
    auto op_states = (operational_states_callback_ != nullptr) ?
      operational_states_callback_() :
      std::list<std::map<std::string, std::string>>();

      for (auto op_state : op_states) {
        State* state = response->add_states();
        state->set_type(op_state["type"]);
        state->set_deviceid(op_state["device_id"]);
        state->set_value(op_state["value"]);
        state->set_version(stoi(version_));
      }

    return Status::OK;
}

void MagmaService::setSharedMetrics() {
  setMetricsUptime();
  setMemoryUsage();
}

void MagmaService::setApplicationHealth(
    ServiceInfo::ApplicationHealth newHealth) {
  health_ = newHealth;
}

void MagmaService::setMetricsStartTime() {
  va_list ap;
  // Use standard time to get start time
  MetricsSingleton::Instance().SetGauge("process_start_time_seconds",
    (double) std::time(nullptr), 0, ap);
}

void MagmaService::setMetricsUptime() {
  va_list ap;
  // Use monotonic time for uptime to avoid clock skew
  steady_clock::time_point t2 = steady_clock::now();
  duration<double> time_span = duration_cast<duration<double>>(
    t2 - start_time_);
  double uptime = time_span.count();
  MetricsSingleton::Instance().SetGauge("process_cpu_seconds_total", uptime, 0,
    ap);
}

void MagmaService::setMemoryUsage() {
  va_list ap;
  const ProcFileUtils::memory_info_t mem_info = ProcFileUtils::getMemoryInfo();
  MetricsSingleton::Instance().SetGauge("process_virtual_memory_bytes",
    mem_info.virtual_mem, 0, ap);
  MetricsSingleton::Instance().SetGauge("process_resident_memory_bytes",
    mem_info.physical_mem, 0, ap);
}
