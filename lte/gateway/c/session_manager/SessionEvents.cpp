/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SessionEvents.h"

using magma::orc8r::Event;
using magma::orc8r::Void;

namespace magma {
namespace session_events {

#define SESSIOND_SERVICE "sessiond"
#define SESSION_CREATED "session_created"
#define SESSION_TERMINATED "session_terminated"
#define SESSION_ID "session_id"
#define IMSI "imsi"
#define IP_ADDR "ip_addr"

#define CHARGING_ "charging_"
#define MONITORING_ "monitoring_"
#define TX "tx"
#define RX "rx"
#define CHARGING_TX CHARGING_ TX
#define CHARGING_RX CHARGING_ RX
#define MONITORING_TX MONITORING_ TX
#define MONITORING_RX MONITORING_ RX


void session_created(
    AsyncEventdClient& client,
    const std::string& imsi,
    const std::string& session_id) {
  auto event = Event();
  event.set_stream_name(SESSIOND_SERVICE);
  event.set_event_type(SESSION_CREATED);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[SESSION_ID] = session_id;
  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  client.log_event(
      event, [=](Status status, Void v) {
      if (!status.ok()) {
      MLOG(MERROR)
      << "Could not log " << SESSION_CREATED << " event " << event_value_string
      << ", Error Message: " << status.error_message();
      }
      });
}

void session_terminated(
    AsyncEventdClient& client,
    const std::unique_ptr<SessionState>& session) {
  auto event = Event();
  SessionState::SessionInfo session_info;
  session->get_session_info(session_info);

  event.set_stream_name(SESSIOND_SERVICE);
  event.set_event_type(SESSION_TERMINATED);
  event.set_tag(session_info.imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = session_info.imsi;
  event_value[IP_ADDR] = session_info.ip_addr;
  event_value[SESSION_ID] = session->get_session_id();
  SessionState::TotalCreditUsage usage = session->get_total_credit_usage();
  event_value[CHARGING_TX] = usage.charging_tx;
  event_value[CHARGING_RX] = usage.charging_rx;
  event_value[MONITORING_TX] = usage.monitoring_tx;
  event_value[MONITORING_RX] = usage.monitoring_rx;
  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  client.log_event(
      event, [=](Status status, Void v) {
      if (!status.ok()) {
      MLOG(MERROR)
      << "Could not log "<< SESSION_TERMINATED << " event " << event_value_string
      << ", Error Message: " << status.error_message();
      }
      });
}

}  // namespace session_events
}  // namespace magma
