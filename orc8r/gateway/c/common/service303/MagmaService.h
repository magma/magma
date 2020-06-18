/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <grpc++/grpc++.h>
#include <orc8r/protos/service303.grpc.pb.h>
#include <chrono>

#include "MetricsRegistry.h"
#include "MetricsSingleton.h"

using grpc::ServerContext;
using grpc::Status;
using grpc::Server;
using magma::orc8r::Service303;
using magma::orc8r::ServiceInfo;
using magma::orc8r::State;
using magma::orc8r::Void;
using magma::service303::MetricsSingleton;

namespace magma { namespace service303 {

using ServiceInfoMeta = std::map<std::string,std::string>;
using States = std::list<std::map<std::string, std::string>>;
using ServiceInfoCallback = std::function<ServiceInfoMeta()>;
using ConfigReloadCallback = std::function<bool()>;
using OperationalStatesCallback = std::function<States()>;

/**
 * MagmaService provides the framework for all Magma services.
 * This class also implements the Service303 interface for external
 * entities to interact with the service.
 */
class MagmaService final : public Service303::Service {
  public:
    MagmaService(const std::string& name, const std::string& version);

    /**
     * Starts the service info gRPC Service
     */
    void Start();

    /**
     * Stops the service info gRPC Service
     */
    void Stop();

    /**
     * Blocks until the gRPC Service finshes shutdown
     */
    void WaitForShutdown();

    /**
     * Add an additional service to the grpc server before starting
     *
     * @param service: pointer to service to add
     */
    void AddServiceToServer(grpc::Service* service);

    /**
     * Return a new completion queue for handling async services
     */
    std::unique_ptr<grpc::ServerCompletionQueue> GetNewCompletionQueue();

    /**
     * Sets the callback to generate the meta field for service info
     */
    void SetServiceInfoCallback(ServiceInfoCallback callback);

    /**
     * Unsets the callback to generate the meta field for service info
     */
    void ClearServiceInfoCallback();

    /**
     * Sets the callback to request a config reload from a service
     */
    void SetConfigReloadCallback(ConfigReloadCallback callback);

    /**
     * Unsets the callback to request a config reload from a service
     */
    void ClearConfigReloadCallback();

    /**
     * Sets the callback to generate the operational states for a service
     */
    void SetOperationalStatesCallback(OperationalStatesCallback callback);

    /**
     * Unsets the callback to generate the operational states for a service
     */
    void ClearOperationalStatesCallback();


    /*
    * Returns the service info (name, version, state, etc.)
    *
    * @param context: the grpc Server context
    * @param request: void request param
    * @param response (out): the ServiceInfo response
    * @return grpc Status instance
    */
    Status GetServiceInfo(
        ServerContext* context,
        const Void* request,
        ServiceInfo* response) override;

    /*
    * Handles request to stop the service
    *
    * @param context: the grpc Server context
    * @param request: void request param
    * @param response (out): void response param
    * @return grpc Status instance
    */
    Status StopService(
        ServerContext* context,
        const Void* request,
        Void* response) override;

    /*
     * Collects timeseries samples from prometheus client interface on this
     * process
     *
     * @param context: the grpc Server context
     * @param request: void request param
     * @param response (out): container of all collected metrics
     * @return grpc Status instance
     */
    Status GetMetrics(
        ServerContext* context,
        const Void* request,
        MetricsContainer* response) override;

    /*
     * Sets the log verbosity to print to syslog at runtime
     *
     * @param context: the grpc Server context
     * @param request: log level (FATAL, ERROR, etc)
     * @param response (out): Void
     * @return grpc Status instance
     */
    Status SetLogLevel(
         ServerContext* context,
         const LogLevelMessage* request,
         Void* response) override;

    /*
     * Handles request to reload the service config file
     *
     * @param context: the grpc Server context
     * @param request: Void
     * @param response (out): reload config result (SUCCESS/FAILURE/UNSUPPORTED)
     * @return grpc Status instance
     */
    Status ReloadServiceConfig(
        ServerContext *context,
        const Void *request,
        ReloadConfigResponse *response) override;

    /*
     * Returns the  operational states of devices managed by this service.
     *
     * @param context: the grpc Server context
     * @param request: Void
     * @param response (out): a list of states
     * @return grpc Status instance
     */
    Status GetOperationalStates(
        ServerContext *context,
        const Void *request,
        GetOperationalStatesResponse *response) override;

    /*
     * Simple setter function to set the new application health
     *
     * @param newState: the new application health you want to set
     *   One of: APP_UNKNOWN, APP_HEALTHY, APP_UNHEALTHY
     */
    void setApplicationHealth(ServiceInfo::ApplicationHealth newHealth);

  private:
    /*
     * Helper function to set the process_start_time_seconds in metricsd
     */
    void setMetricsStartTime();

    /*
     * Helper function to set all shared metrics among all services, like
     * uptime and memory usage
     */
    void setSharedMetrics();

    /*
     * Helper function to set the process_cpu_seconds_total in metrics
     */
    void setMetricsUptime();

    /*
     * Helper function to set process physical memory and virtual memory
     */
    void setMemoryUsage();

  private:
    const std::string name_;
    const std::string version_;
    const std::chrono::steady_clock::time_point start_time_;
    const std::chrono::system_clock::time_point wall_start_time_;
    ServiceInfo::ApplicationHealth health_;
    std::unique_ptr<Server> server_;
    grpc::ServerBuilder builder_;
    ServiceInfoCallback service_info_callback_;
    ConfigReloadCallback config_reload_callback_;
    OperationalStatesCallback operational_states_callback_;
};

}} // namespace magma::service303
