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

#include "lte/gateway/c/core/oai/include/mme_events.hpp"

#include <cstdlib>
#include <iostream>
#include <nlohmann/json.hpp>
#include <grpcpp/support/status.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/conversions.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/common/common_types.h"

#include "orc8r/protos/common.pb.h"
#include "orc8r/protos/eventd.pb.h"

#include "lte/gateway/c/core/oai/lib/event_client/EventClientAPI.hpp"

using grpc::Status;
using magma::lte::init_eventd_client;
using magma::lte::log_event;
using magma::orc8r::Event;
using magma::orc8r::Void;

namespace {
constexpr char MME_STREAM_NAME[] = "mme";
constexpr char ATTACH_SUCCESS[] = "attach_success";
constexpr char DETACH_SUCCESS[] = "detach_success";
constexpr char S1_SETUP_SUCCESS[] = "s1_setup_success";
constexpr char ATTACH_REJECT[] = "attach_reject";
}  // namespace

void event_client_init(void) { init_eventd_client(); }

/**
 * Helper function to log event by sending RPC call to eventd service
 * @param event_value
 * @param event_type
 * @param stream_name
 * @param event_tag
 * @return response code
 */
static int report_event(const nlohmann::json& event_value,
                        const std::string& event_type,
                        const std::string& stream_name,
                        const std::string& event_tag) {
  Event event_request = Event();
  event_request.set_event_type(event_type);
  event_request.set_stream_name(stream_name);

  std::string event_value_string = event_value.dump();
  event_request.set_value(event_value_string);
  event_request.set_tag(event_tag);
  return log_event(event_request);
}

int attach_success_event(imsi64_t imsi64) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);

  nlohmann::json event_value;
  event_value["imsi"] = imsi_str;

  return report_event(event_value, ATTACH_SUCCESS, MME_STREAM_NAME, imsi_str);
}

int detach_success_event(imsi64_t imsi64, const char* action) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);

  nlohmann::json event_value;
  event_value["imsi"] = imsi_str;
  event_value["action"] = action;

  return report_event(event_value, DETACH_SUCCESS, MME_STREAM_NAME, imsi_str);
}

int s1_setup_success_event(const char* enb_name, uint32_t enb_id) {
  nlohmann::json event_value;

  if (enb_name) {
    event_value["enb_name"] = enb_name;
  } else {
    event_value["enb_name"] = "";
  }

  event_value["enb_id"] = enb_id;

  return report_event(event_value, S1_SETUP_SUCCESS, MME_STREAM_NAME,
                      std::to_string(enb_id));
}

int attach_reject_event(imsi64_t imsi64) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);

  nlohmann::json event_value;
  event_value["imsi"] = imsi_str;

  return report_event(event_value, ATTACH_REJECT, MME_STREAM_NAME, imsi_str);
}
