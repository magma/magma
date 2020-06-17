/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "mme_events.h"

#include <cstdlib>
#include <iostream>
#include <folly/Format.h>
#include <folly/json.h>
#include <folly/dynamic.h>
#include <grpcpp/support/status.h>

#include "conversions.h"
#include "common_types.h"
#include "orc8r/protos/common.pb.h"
#include "orc8r/protos/eventd.pb.h"

#include "EventClientAPI.h"

using grpc::Status;
using magma::orc8r::Event;
using magma::orc8r::Void;
using magma::lte::log_event;
using magma::lte::init_eventd_client;
using magma::orc8r::Void;

namespace {
constexpr char MME_STREAM_NAME[]   = "mme";
constexpr char ATTACH_SUCCESSFUL[] = "attach_successful";
constexpr char DETACH_SUCCESSFUL[] = "detach_successful";
}  // namespace

int event_client_init(void) {
  init_eventd_client();
}

int attach_successful(imsi64_t imsi64) {
  Event eventRequest = Event();
  eventRequest.set_event_type(ATTACH_SUCCESSFUL);
  eventRequest.set_stream_name(MME_STREAM_NAME);

  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;

  std::string event_value_string = folly::toJson(event_value);
  eventRequest.set_value(event_value_string);
  eventRequest.set_tag(imsi_str);

  int rc = log_event(eventRequest);
  return rc;
}

int detach_successful(imsi64_t imsi64, const char* action) {
  Event eventRequest = Event();
  eventRequest.set_event_type(DETACH_SUCCESSFUL);
  eventRequest.set_stream_name(MME_STREAM_NAME);

  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["action"]      = action;

  std::string event_value_string = folly::toJson(event_value);
  eventRequest.set_value(event_value_string);
  eventRequest.set_tag(imsi_str);

  int rc = log_event(eventRequest);
  return rc;
}
