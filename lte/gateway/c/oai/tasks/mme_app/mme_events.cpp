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
using magma::lte::init_eventd_client;
using magma::lte::log_event;
using magma::orc8r::Event;
using magma::orc8r::Void;

namespace {
constexpr char MME_STREAM_NAME[]  = "mme";
constexpr char ATTACH_SUCCESS[]   = "attach_success";
constexpr char DETACH_SUCCESS[]   = "detach_success";
constexpr char S1_SETUP_SUCCESS[] = "s1_setup_success";
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

int attach_success_event(imsi64_t imsi64) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;

  return report_event(event_value, ATTACH_SUCCESS, MME_STREAM_NAME, imsi_str);
}

int detach_success_event(imsi64_t imsi64, const char* action) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["action"]      = action;

  return report_event(event_value, DETACH_SUCCESS, MME_STREAM_NAME, imsi_str);
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
