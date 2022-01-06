/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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

#include "lte/gateway/c/core/oai/include/mme_events.h"

#include <cstdlib>
#include <iostream>
#include <folly/Format.h>
#include <folly/json.h>
#include <folly/dynamic.h>
#include <grpcpp/support/status.h>

#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/common_types.h"

#include "orc8r/protos/common.pb.h"
#include "orc8r/protos/eventd.pb.h"

#include "lte/gateway/c/core/oai/lib/event_client/EventClientAPI.h"

using grpc::Status;
using magma::lte::init_eventd_client;
using magma::lte::log_event;
using magma::orc8r::Event;
using magma::orc8r::Void;

namespace {
constexpr char MME_STREAM_NAME[] = "mme";
constexpr char ATTACH_REQUEST[]  = "attach_request";
constexpr char ATTACH_ACCEPT[]   = "attach_accept";
constexpr char ATTACH_REJECT[]   = "attach_reject";
constexpr char ATTACH_COMPLETE[] = "attach_complete";
constexpr char ATTACH_SUCCESS[]  = "attach_success";
constexpr char ATTACH_FAILURE[]  = "attach_failure";
constexpr char DETACH_REQUEST[]  = "detach_request";
constexpr char DETACH_ACCEPT[]   = "detach_accept";
constexpr char DETACH_IMPLICIT[] = "detach_implicit";
constexpr char DETACH_SUCCESS[]  = "detach_success";
constexpr char DETACH_FAILURE[]  = "detach_failure";
constexpr char INITIAL_CONTEXT_SETUP_REQUEST[] =
    "initial_context_setup_request";
constexpr char INITIAL_CONTEXT_SETUP_RESPONSE[] =
    "initial_context_setup_response";
constexpr char INITIAL_CONTEXT_SETUP_FAILURE[] =
    "initial_context_setup_failure";
constexpr char UE_CONTEXT_RELEASE_REQUEST[]  = "ue_context_release_request";
constexpr char UE_CONTEXT_RELEASE_COMMAND[]  = "ue_context_release_command";
constexpr char UE_CONTEXT_RELEASE_COMPLETE[] = "ue_context_release_complete";
constexpr char AUTHENTICATION_INFORMATION_REQUEST[] = "authentication_information_request";
constexpr char AUTHENTICATION_INFORMATION_ANSWER[] = "authentication_information_answer";
constexpr char S1_SETUP_SUCCESS[]            = "s1_setup_success";
}  // namespace

void event_client_init(void) {
  init_eventd_client();
}

/**
 * Helper function to log event by sending RPC call to eventd service
 * @param event_value
 * @param event_type
 * @param stream_name
 * @param event_tag
 * @return response code
 */
static int report_event(
    folly::dynamic& event_value, const std::string& event_type,
    const std::string& stream_name, const std::string& event_tag) {
  Event event_request = Event();
  event_request.set_event_type(event_type);
  event_request.set_stream_name(stream_name);

  std::string event_value_string = folly::toJson(event_value);
  event_request.set_value(event_value_string);
  event_request.set_tag(event_tag);
  return log_event(event_request);
}

int attach_request_event(
    imsi64_t imsi64, const guti_t guti, const char* imei, const char* mme_id,
    const char* enb_id, const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";  // TODO(andreilee): Fix this
  event_value["imei"]        = imei;
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, ATTACH_REQUEST, MME_STREAM_NAME, imsi_str);
}

int attach_accept_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, ATTACH_ACCEPT, MME_STREAM_NAME, imsi_str);
}

int attach_reject_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["cause"]       = cause;

  return report_event(event_value, ATTACH_REJECT, MME_STREAM_NAME, imsi_str);
}

int attach_complete_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, ATTACH_COMPLETE, MME_STREAM_NAME, imsi_str);
}

int attach_success_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, ATTACH_SUCCESS, MME_STREAM_NAME, imsi_str);
}

int attach_failure_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* cause, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, ATTACH_FAILURE, MME_STREAM_NAME, imsi_str);
}

int detach_request_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* source) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["source"]      = source;

  return report_event(event_value, DETACH_REQUEST, MME_STREAM_NAME, imsi_str);
}

int detach_accept_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* source) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["source"]      = source;

  return report_event(event_value, DETACH_ACCEPT, MME_STREAM_NAME, imsi_str);
}

int detach_implicit_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(event_value, DETACH_IMPLICIT, MME_STREAM_NAME, imsi_str);
}

int detach_success_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* action) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["action"]      = action;

  return report_event(event_value, DETACH_SUCCESS, MME_STREAM_NAME, imsi_str);
}

int detach_failure_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["cause"]       = cause;

  return report_event(event_value, DETACH_FAILURE, MME_STREAM_NAME, imsi_str);
}

int initial_context_setup_request_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, INITIAL_CONTEXT_SETUP_REQUEST, MME_STREAM_NAME, imsi_str);
}

int initial_context_setup_response_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, INITIAL_CONTEXT_SETUP_RESPONSE, MME_STREAM_NAME, imsi_str);
}

int initial_context_setup_failure_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["cause"]         = cause;

  return report_event(
      event_value, INITIAL_CONTEXT_SETUP_FAILURE, MME_STREAM_NAME, imsi_str);
}

int ue_context_release_request_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;
  event_value["cause"]       = cause;

  return report_event(
      event_value, UE_CONTEXT_RELEASE_REQUEST, MME_STREAM_NAME, imsi_str);
}
int ue_context_release_command_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, UE_CONTEXT_RELEASE_COMMAND, MME_STREAM_NAME, imsi_str);
}
int ue_context_release_complete_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, UE_CONTEXT_RELEASE_COMPLETE, MME_STREAM_NAME, imsi_str);
}

int authentication_information_request_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, AUTHENTICATION_INFORMATION_REQUEST, MME_STREAM_NAME, imsi_str);
}

int authentication_information_answer_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["guti"]        = "guti";
  event_value["mme_id"]      = mme_id;
  event_value["enb_id"]      = enb_id;
  event_value["enb_ip"]      = enb_ip;
  event_value["apn"]         = apn;

  return report_event(
      event_value, AUTHENTICATION_INFORMATION_ANSWER, MME_STREAM_NAME, imsi_str);
}

int s1_setup_success_event(const char* enb_name, uint32_t enb_id) {
  folly::dynamic event_value = folly::dynamic::object;

  if (enb_name) {
    event_value["enb_name"] = enb_name;
  } else {
    event_value["enb_name"] = "";
  }

  event_value["enb_id"] = enb_id;

  return report_event(
      event_value, S1_SETUP_SUCCESS, MME_STREAM_NAME,
      folly::to<std::string>(enb_id));
}
