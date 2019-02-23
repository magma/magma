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

#include <orc8r/protos/logging_service.grpc.pb.h>

#include "scribe_rpc_client.h"

#include "GRPCReceiver.h"

using grpc::Status;

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using google::protobuf::RepeatedPtrField;
using magma::orc8r::LoggingService;
using magma::orc8r::LogRequest;
using magma::orc8r::LoggerDestination;


namespace magma {
using namespace orc8r;
/*
 * gRPC client for LoggingService
 */
class LoggingServiceClient : public GRPCReceiver{
 public:
  /**
   * Log one scribe entry to the given category on scribe. API for C.
   *
   * @param category: category name of the scribe category to log to.
   * @param time: a timestamp associated with the logentry.
   * @param int_params[]: an array of scribe_int_param_t, where each
   * scribe_int_param_t contains a str name, and a int value.
   * @param int_params_len: length of the above array.
   * @param str_params[]: an array of scribe_string_param_t, where each
     scribe_string_param_t contains a str name, and a str value.
   * @param str_params_len: length of the above array.
   * @param sampling_rate: a float between 0 and 1 indicating the desired
   * samplingRate of the log. The ScribeClient will throw a die with value in
   * [0, 1) and drop the attempt to log the entry if the result of the die is
   * larger than the samplingRate.
   * @param callback: callback function is called when LogToScribe returns
   */
  static int log_to_scribe(
      char const *category,
      int time,
      scribe_int_param_t *int_params,
      int int_params_len,
      scribe_string_param_t *str_params,
      int str_params_len,
      float sampling_rate,
      std::function<void(Status, Void)> callback);

  /**
   * API for C++. Log on scribe entry to the given category on scribe.
   * @param category category name of the scribe category to log to.
   * @param time a timestamp associated with the logentry.
   * @param int_params a map of int parameters with string keys
   * @param str_params a map of string parameters with string keys
   * @param sampling_rate a float between 0 and 1 indicating the desired
   * sampling_rate of the log. The ScribeClient will throw a die with value in
   * [0, 1) and drop the attempt to log the entry if the result of the die is
   * larger than the sampling_rate.
   * @param callback callback function is called when LogToScribe returns
   */
  static void log_to_scribe(
      std::string category,
      time_t time,
      std::map<std::string, int> int_params,
      std::map<std::string, std::string> str_params,
      float sampling_rate,
      std::function<void(Status, Void)> callback);

 public:
  LoggingServiceClient(LoggingServiceClient const&) = delete;
  void operator=(LoggingServiceClient const&) = delete;

 private:
  explicit LoggingServiceClient();
  static LoggingServiceClient& get_instance();
  std::shared_ptr<LoggingService::Stub> stub_;
  bool shouldLog(float samplingRate);
  void initializeClient();
  static const uint32_t RESPONSE_TIMEOUT = 3; // seconds
};

} // namespace magma
