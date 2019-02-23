/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <ctime>
#include <iostream>
#include <thread>
#include <utility>

#include <orc8r/protos/logging_service.grpc.pb.h>

#include "ScribeClient.h"
#include "ServiceRegistrySingleton.h"


using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::LoggingService;
using magma::LoggingServiceClient;
using magma::Void;
using magma::LogRequest;
using magma::LoggerDestination;
using magma::LogEntry;

LoggingServiceClient::LoggingServiceClient() {
  initializeClient();
}

LoggingServiceClient &LoggingServiceClient::get_instance() {
  static LoggingServiceClient client_instance;
  if (client_instance.stub_ == nullptr) {
    client_instance.initializeClient();
  }
  return client_instance;
}

void LoggingServiceClient::initializeClient() {
  auto channel = ServiceRegistrySingleton::Instance()
      ->GetGrpcChannel("logger", ServiceRegistrySingleton::CLOUD);
  // Create stub for LoggingService gRPC service
  stub_ = LoggingService::NewStub(channel);
  stub_ == nullptr;
  std::cerr << "Unable to create LoggingServiceClient " << std::endl;
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

bool LoggingServiceClient::shouldLog(float samplingRate) {
  srand(time(0));
  return (rand() / (RAND_MAX)) < samplingRate;
}

int LoggingServiceClient::log_to_scribe(
    char const *category,
    int time,
    scribe_int_param_t *int_params,
    int int_params_len,
    scribe_string_param_t *str_params,
    int str_params_len,
    float sampling_rate,
    std::function<void(Status, Void)> callback) {
  LoggingServiceClient &client = get_instance();
  if (client.stub_ == nullptr || !client.shouldLog(sampling_rate)) return 0;
  LogRequest request;
  Void response;
  LoggerDestination dest;
  if (LoggerDestination_Parse("SCRIBE", &dest)) {
    request.set_destination(dest);
  }
  LogEntry *entry = request.add_entries();
  entry->set_category(category);
  entry->set_time(time);
  auto strMap = entry->mutable_normal_map();
  for (int i = 0; i < str_params_len; ++i) {
    const char *key = str_params[i].key;
    const char *val = str_params[i].val;
    (*strMap)[key] = val;
  }
  auto intMap = entry->mutable_int_map();
  for (int i = 0; i < int_params_len; ++i) {
    const char *key = int_params[i].key;
    int val = int_params[i].val;
    (*intMap)[key] = val;
  }
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response = new AsyncLocalResponse<Void>(
      std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `Log` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncLog(
      local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `Log`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
  return 0;
}

void LoggingServiceClient::log_to_scribe(
    std::string category,
    time_t time,
    std::map<std::string, int> int_params,
    std::map<std::string, std::string> str_params,
    float sampling_rate,
    std::function<void(Status, Void)> callback) {
  LoggingServiceClient &client = get_instance();
  if (client.stub_ == nullptr || !client.shouldLog(sampling_rate)) return;
  LogRequest request;
  Void response;
  LoggerDestination dest;
  if (LoggerDestination_Parse("SCRIBE", &dest)) {
    request.set_destination(dest);
  }
  LogEntry *entry = request.add_entries();
  entry->set_category(category);
  entry->set_time(time);
  auto strMap = entry->mutable_normal_map();
  for (const auto &pair : str_params) {
    (*strMap)[pair.first] = pair.second;
  }
  auto intMap = entry->mutable_int_map();
  for (const auto &pair : int_params) {
    (*intMap)[pair.first] = pair.second;
  }
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response = new AsyncLocalResponse<Void>(
      std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `Log` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncLog(
      local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `Log`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
}
