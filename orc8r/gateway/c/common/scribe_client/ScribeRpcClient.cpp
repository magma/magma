/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <string>
#include <thread>

#include "ScribeClient.h"
#include "scribe_rpc_client_for_cpp.h"
#include "ServiceRegistrySingleton.h"

using magma::LoggingServiceClient;
using magma::ServiceRegistrySingleton;
using magma::LogEntry;

static void log_to_scribe_done(const grpc::Status& status);

int log_to_scribe(
    char const *category,
    scribe_int_param_t *int_params,
    int int_params_len,
    scribe_string_param_t *str_params,
    int str_params_len) {
  time_t t = time(NULL);
  return log_to_scribe_with_time_and_sampling_rate(category,
                                                   t,
                                                   int_params,
                                                   int_params_len,
                                                   str_params,
                                                   str_params_len,
                                                   1);
}

int log_to_scribe_with_sampling_rate(
    char const *category,
    scribe_int_param_t *int_params,
    int int_params_len,
    scribe_string_param_t *str_params,
    int str_params_len,
    float sampling_rate) {
  time_t t = time(NULL);
  return log_to_scribe_with_time_and_sampling_rate(category,
                                                   t,
                                                   int_params,
                                                   int_params_len,
                                                   str_params,
                                                   str_params_len,
                                                   sampling_rate);
}

int log_to_scribe_with_time_and_sampling_rate(
    char const *category,
    int time,
    scribe_int_param_t *int_params,
    int int_params_len,
    scribe_string_param_t *str_params,
    int str_params_len,
    float sampling_rate) {
  int status = LoggingServiceClient::log_to_scribe(
      category,
      time,
      int_params,
      int_params_len,
      str_params,
      str_params_len,
      sampling_rate,
      [](grpc::Status g_status, magma::Void response) {
        log_to_scribe_done(g_status);
      });
  return status;
}

void log_to_scribe_with_time_and_sampling_rate(
    std::string category,
    time_t time,
    std::map<std::string, int> int_params,
    std::map<std::string, std::string> str_params,
    float sampling_rate) {
  return LoggingServiceClient::log_to_scribe(
      category,
      time,
      int_params,
      str_params,
      sampling_rate,
      [](grpc::Status status, magma::Void response) {
        log_to_scribe_done(status);
      });
}

void log_to_scribe(
    std::string category,
    std::map<std::string, int> int_params,
    std::map<std::string, std::string> str_params) {
  time_t t = time(NULL);
  return log_to_scribe_with_time_and_sampling_rate(
      category, t, int_params, str_params, 1);
}

static void log_to_scribe_done(const grpc::Status& status) {
  if (!status.ok()) {
    std::cerr << "log_to_scribe fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
  }
}
