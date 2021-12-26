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

#include "orc8r/gateway/c/common/service303/includes/MagmaService.h"

#include <grpcpp/impl/codegen/completion_queue.h>  // for ServerCompletionQueue
#include <grpcpp/security/server_credentials.h>    // for InsecureServerCred...
#include <orc8r/protos/common.pb.h>                // for FATAL, LogLevel
#include <orc8r/protos/metricsd.pb.h>              // for MetricsContainer
#include <orc8r/protos/service303.pb.h>            // for ServiceInfo, Reloa...
#include <prometheus/registry.h>                   // for Registry
#include <stdarg.h>                                // for va_list
#include <google/protobuf/stubs/common.h>
#include <metrics.pb.h>
#include <chrono>       // for seconds, duration
#include <csignal>      // for raise, SIGTERM
#include <ctime>        // for time
#include <string>       // for string, stoi
#include <type_traits>  // for enable_if<>::type
#include <utility>      // for move
#include <vector>       // for vector
#include <algorithm>

#include "orc8r/gateway/c/common/service303/includes/MetricsSingleton.h"  // for MetricsSingleton
#include "orc8r/gateway/c/common/service303/ProcFileUtils.h"  // for ProcFileUtils::mem...
#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"  // for ServiceRegistrySin...
#include "orc8r/gateway/c/common/logging/magma_logging_init.h"  // for set_verbosity

namespace grpc {
class ServerContext;
}
namespace grpc {
class Service;
}

using grpc::InsecureServerCredentials;
using grpc::Status;
using io::prometheus::client::MetricFamily;
using magma::orc8r::GetOperationalStatesResponse;
using magma::orc8r::ReloadConfigResponse;
using magma::orc8r::Service303;
using magma::orc8r::ServiceInfo;
using magma::orc8r::State;
using magma::orc8r::Void;
using magma::service303::MagmaService;
using magma::service303::MetricsSingleton;
using namespace std::chrono;

MagmaService::MagmaService(const std::string& name, const std::string& version)
    : name_(name),
      version_(version),
      start_time_(steady_clock::now()),
      wall_start_time_(system_clock::now()),
      health_(ServiceInfo::APP_UNKNOWN),
      service_info_callback_(nullptr),
      config_reload_callback_(nullptr),
      operational_states_callback_(nullptr) {}

void MagmaService::AddServiceToServer(grpc::Service* service) {
  builder_.RegisterService(service);
}

std::unique_ptr<grpc::ServerCompletionQueue>
MagmaService::GetNewCompletionQueue() {
  return std::move(builder_.AddCompletionQueue());
}

void MagmaService::Start() {
  setMetricsStartTime();
  builder_.RegisterService(this);
  std::string service_addr =
      magma::ServiceRegistrySingleton::Instance()->GetServiceAddrString(name_);
  builder_.AddListeningPort(service_addr, grpc::InsecureServerCredentials());
  server_ = builder_.BuildAndStart();
}

void MagmaService::WaitForShutdown() {
  server_->Wait();  // Blocking call
}

void MagmaService::Stop() { server_->Shutdown(); }

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
    OperationalStatesCallback callback) {
  operational_states_callback_ = callback;
}

void MagmaService::ClearOperationalStatesCallback() {
  operational_states_callback_ = nullptr;
}

Status MagmaService::GetServiceInfo(__attribute__((unused))
                                    ServerContext* context,
                                    __attribute__((unused)) const Void* request,
                                    ServiceInfo* response) {
  auto start_time_secs =
      time_point_cast<seconds>(wall_start_time_).time_since_epoch().count();

  auto meta = (service_info_callback_ != nullptr)
                  ? service_info_callback_()
                  : std::map<std::string, std::string>();

  response->set_name(name_);
  response->set_version(version_);
  response->set_state(ServiceInfo::ALIVE);
  response->set_health(health_);
  response->set_start_time_secs(start_time_secs);
  response->mutable_status()->mutable_meta()->insert(meta.begin(), meta.end());

  return Status::OK;
}

Status MagmaService::StopService(__attribute__((unused)) ServerContext* context,
                                 __attribute__((unused)) const Void* request,
                                 __attribute__((unused)) Void* response) {
  std::raise(SIGTERM);
  return Status::OK;
}

Status MagmaService::GetMetrics(__attribute__((unused)) ServerContext* context,
                                __attribute__((unused)) const Void* request,
                                MetricsContainer* response) {
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

Status MagmaService::SetLogLevel(__attribute__((unused)) ServerContext* context,
                                 const LogLevelMessage* request,
                                 __attribute__((unused)) Void* response) {
  // log level FATAL is minimum verbosity and maximum level
  auto verbosity = LogLevel::FATAL - request->level();
  set_verbosity(verbosity);
  return Status::OK;
}

Status MagmaService::ReloadServiceConfig(__attribute__((unused))
                                         ServerContext* context,
                                         __attribute__((unused))
                                         const Void* request,
                                         ReloadConfigResponse* response) {
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
    __attribute__((unused)) ServerContext* context,
    __attribute__((unused)) const Void* request,
    GetOperationalStatesResponse* response) {
  auto op_states = (operational_states_callback_ != nullptr)
                       ? operational_states_callback_()
                       : std::list<std::map<std::string, std::string>>();

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
                                        (double)std::time(nullptr), 0, ap);
}

void MagmaService::setMetricsUptime() {
  va_list ap;
  // Use monotonic time for uptime to avoid clock skew
  steady_clock::time_point t2 = steady_clock::now();
  duration<double> time_span =
      duration_cast<duration<double>>(t2 - start_time_);
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
