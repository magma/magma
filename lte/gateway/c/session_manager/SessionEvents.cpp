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

#include <folly/json.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status.h>
#include <stdint.h>
#include <ostream>
#include <unordered_map>
#include <utility>
#include <vector>

#include "CreditKey.h"
#include "EnumToString.h"
#include "SessionCredit.h"
#include "SessionEvents.h"
#include "SessionState.h"
#include "Types.h"
#include "Utilities.h"
#include "includes/EventdClient.h"
#include "lte/protos/policydb.pb.h"
#include "lte/protos/session_manager.pb.h"
#include "magma_logging.h"
#include "orc8r/protos/common.pb.h"
#include "orc8r/protos/eventd.pb.h"

using magma::orc8r::Event;
using magma::orc8r::Void;

namespace {  // anonymous

const std::string SESSIOND_SERVICE_EV = "sessiond";
const std::string SESSION_CREATED_EV = "session_created";
const std::string SESSION_CREATE_FAILURE_EV = "session_create_failure";
const std::string SESSION_UPDATED_EV = "session_updated";
const std::string SESSION_UPDATE_FAILURE_EV = "session_update_failure";
const std::string SESSION_TERMINATED_EV = "session_terminated";

const std::string SESSION_ID = "session_id";
const std::string IMSI = "imsi";
const std::string IMEI = "imei";
const std::string USER_LOCATION = "user_location";
const std::string IP_ADDR = "ip_addr";
const std::string IPV6_ADDR = "ipv6_addr";
const std::string MAC_ADDR = "mac_addr";
const std::string MSISDN = "msisdn";
const std::string SPGW_IP = "spgw_ip";
const std::string APN = "apn";
const std::string PDP_START_TIME = "pdp_start_time";
const std::string PDP_END_TIME = "pdp_end_time";
const std::string DURATION = "duration";
const std::string FAILURE_REASON = "failure_reason";
const std::string RECORD_SEQUENCE_NUMBER = "record_sequence_number";
const std::string CAUSE_FOR_RECORD_CLOSING = "cause_for_rec_closing";
const std::string CHARGING_CHARACTERISTICS = "charging_characteristics";

const std::string SERVICE_DATA = "list_of_service_data";
const std::string RATING_GROUP = "rating_group";
const std::string SERVICE_IDENTIFIER = "service_identifier";
const std::string DATA_UPLINK = "data_volume_downlink";
const std::string DATA_DOWNLINK = "data_volume_uplink";
const std::string TIME_OF_FIRST_USAGE = "time_of_first_usage";
const std::string TIME_OF_LAST_USAGE = "time_of_last_usage";
const std::string SERVICE_CONDITION_CHANGE = "service_condition_change";
const std::string SERVICE_UPDATES = "list_of_updates";
const std::string UPDATE_REASON = "update_reason";
const std::string MONITORING_KEY = "monitoring_key";

const std::string TOTAL_TX = "total_tx";
const std::string TOTAL_RX = "total_rx";
const std::string CHARGING_TX = "charging_tx";
const std::string CHARGING_RX = "charging_rx";
const std::string MONITORING_TX = "monitoring_tx";
const std::string MONITORING_RX = "monitoring_rx";

enum CauseForRecordClosing {  // TS 132 298
  NORMAL_RELEASE = 0,
};
enum ServiceConditionChange {  // TS 132 298
  SERVICE_STOP = 1 << 9,
};
}  // namespace

namespace magma {
namespace lte {

EventsReporterImpl::EventsReporterImpl(EventdClient& eventd_client)
    : eventd_client_(eventd_client) {}

void EventsReporterImpl::session_created(
    const std::string& imsi, const std::string& session_id,
    const SessionConfig& session_context,
    const std::unique_ptr<SessionState>& session) {
  auto event = magma::orc8r::Event();
  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_CREATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[IP_ADDR] = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR] = session_context.common_context.ue_ipv6();
  event_value[MSISDN] = session_context.common_context.msisdn();
  event_value[APN] = session_context.common_context.apn();
  event_value[SESSION_ID] = session_id;
  event_value[PDP_START_TIME] = session->get_pdp_start_time();
  // LTE specific
  event_value[IMEI] = get_imei(session_context);
  event_value[SPGW_IP] = get_spgw_ipv4(session_context);
  event_value[USER_LOCATION] = get_user_location(session_context);
  event_value[CHARGING_CHARACTERISTICS] =
      get_charging_characteristics(session_context);
  // CWF specific
  event_value[MAC_ADDR] = get_mac_addr(session_context);

  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  eventd_client_.log_event(event, [=](Status status, Void v) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not log " << SESSION_CREATED_EV << " event "
                   << event_value_string
                   << ", Error Message: " << status.error_message();
    }
  });
}

void EventsReporterImpl::session_create_failure(
    const SessionConfig& session_context, const std::string& failure_reason) {
  auto event = magma::orc8r::Event();
  const std::string imsi = session_context.get_imsi();
  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_CREATE_FAILURE_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[APN] = session_context.common_context.apn();
  event_value[FAILURE_REASON] = failure_reason;
  event_value[MAC_ADDR] = get_mac_addr(session_context);

  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  eventd_client_.log_event(event, [=](Status status, Void v) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not log " << SESSION_CREATE_FAILURE_EV << " event "
                   << event_value_string
                   << ", Error Message: " << status.error_message();
    }
  });
}

void EventsReporterImpl::session_updated(const std::string& session_id,
                                         const SessionConfig& session_context,
                                         const UpdateRequests& update_request) {
  auto event = magma::orc8r::Event();
  const std::string imsi = session_context.get_imsi();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_UPDATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[SESSION_ID] = session_id;
  event_value[IP_ADDR] = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR] = session_context.common_context.ue_ipv6();
  event_value[APN] = session_context.common_context.apn();
  event_value[MAC_ADDR] = get_mac_addr(session_context);
  event_value[SERVICE_UPDATES] = get_update_summary(update_request);

  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  eventd_client_.log_event(event, [=](Status status, Void v) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not log " << SESSION_UPDATED_EV << " event "
                   << event_value_string
                   << ", Error Message: " << status.error_message();
    }
  });
}

void EventsReporterImpl::session_update_failure(
    const std::string& session_id, const SessionConfig& session_context,
    const UpdateRequests& failed_request, const std::string& failure_reason) {
  auto event = magma::orc8r::Event();
  const std::string imsi = session_context.get_imsi();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_UPDATE_FAILURE_EV);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[SESSION_ID] = session_id;
  event_value[IP_ADDR] = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR] = session_context.common_context.ue_ipv6();
  event_value[MAC_ADDR] = get_mac_addr(session_context);
  event_value[APN] = session_context.common_context.apn();
  event_value[FAILURE_REASON] = failure_reason;
  event_value[SERVICE_UPDATES] = get_update_summary(failed_request);

  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  eventd_client_.log_event(event, [=](Status status, Void v) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not log " << SESSION_UPDATE_FAILURE_EV << " event "
                   << event_value_string
                   << ", Error Message: " << status.error_message();
    }
  });
}

void EventsReporterImpl::session_terminated(
    const std::string& imsi, const std::unique_ptr<SessionState>& session) {
  auto event = magma::orc8r::Event();
  auto session_cfg = session->get_config();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_TERMINATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI] = imsi;
  event_value[IP_ADDR] = session_cfg.common_context.ue_ipv4();
  event_value[IPV6_ADDR] = session_cfg.common_context.ue_ipv6();
  event_value[MSISDN] = session_cfg.common_context.msisdn();
  event_value[APN] = session_cfg.common_context.apn();
  event_value[SESSION_ID] = session->get_session_id();

  TotalCreditUsage usage = session->get_total_credit_usage();
  event_value[TOTAL_TX] = usage.charging_tx + usage.monitoring_tx;
  event_value[TOTAL_RX] = usage.charging_rx + usage.monitoring_rx;
  event_value[CHARGING_TX] = usage.charging_tx;
  event_value[CHARGING_RX] = usage.charging_rx;
  event_value[MONITORING_TX] = usage.monitoring_tx;
  event_value[MONITORING_RX] = usage.monitoring_rx;
  const auto start_time = session->get_pdp_start_time();
  const auto end_time = session->get_pdp_end_time();
  event_value[PDP_START_TIME] = start_time;
  event_value[PDP_END_TIME] = end_time;
  // TODO these fields below should be handled by a CDR processor script
  event_value[DURATION] = end_time - start_time;
  event_value[CAUSE_FOR_RECORD_CLOSING] = int(NORMAL_RELEASE);
  event_value[RECORD_SEQUENCE_NUMBER] = 1;
  // LTE specific
  event_value[IMEI] = get_imei(session_cfg);
  event_value[SPGW_IP] = get_spgw_ipv4(session_cfg);
  event_value[USER_LOCATION] = get_user_location(session_cfg);
  event_value[CHARGING_CHARACTERISTICS] =
      get_charging_characteristics(session_cfg);
  // CWF specific
  event_value[MAC_ADDR] = get_mac_addr(session_cfg);

  // Add Gy tracked credits
  auto credit_summaries = session->get_charging_credit_summaries();
  folly::dynamic service_data_list = folly::dynamic::array;
  for (auto summary_pair : credit_summaries) {
    folly::dynamic service_data = folly::dynamic::object;
    auto summary = summary_pair.second;
    service_data[RATING_GROUP] = summary_pair.first.rating_group;
    if (summary_pair.first.service_identifier) {
      service_data[SERVICE_IDENTIFIER] = summary_pair.first.service_identifier;
    }
    service_data[DATA_UPLINK] = summary.usage.bytes_tx;
    service_data[DATA_DOWNLINK] = summary.usage.bytes_rx;
    service_data[TIME_OF_FIRST_USAGE] = summary.time_of_first_usage;
    service_data[TIME_OF_LAST_USAGE] = summary.time_of_last_usage;
    service_data[SERVICE_CONDITION_CHANGE] = int(SERVICE_STOP);
    service_data_list.push_back(service_data);
  }
  event_value[SERVICE_DATA] = service_data_list;

  std::string event_value_string = folly::toJson(event_value);
  event.set_value(event_value_string);

  eventd_client_.log_event(event, [=](Status status, Void v) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not log " << SESSION_TERMINATED_EV << " event "
                   << event_value_string
                   << ", Error Message: " << status.error_message();
    }
  });
}

folly::dynamic EventsReporterImpl::get_update_summary(
    const UpdateRequests& updates) {
  folly::dynamic update_array = folly::dynamic::array;
  for (const auto& charging : updates.charging_requests) {
    folly::dynamic data = folly::dynamic::object;
    data[RATING_GROUP] = charging.usage().charging_key();
    if (charging.usage().has_service_identifier()) {
      data[SERVICE_IDENTIFIER] = charging.usage().service_identifier().value();
    }
    data[UPDATE_REASON] = credit_update_type_to_str(charging.usage().type());
    update_array.push_back(data);
  }
  for (const auto& monitor : updates.monitor_requests) {
    folly::dynamic data = folly::dynamic::object;
    data[UPDATE_REASON] = event_trigger_to_str(monitor.event_trigger());
    if (monitor.has_update()) {
      data[MONITORING_KEY] = monitor.update().monitoring_key();
    }
    update_array.push_back(data);
  }
  return update_array;
}

std::string EventsReporterImpl::get_mac_addr(const SessionConfig& config) {
  // MacAddr is only relevant for WLAN
  const auto& rat_specific = config.rat_specific_context;
  std::string mac_addr = "";
  if (rat_specific.has_wlan_context()) {
    mac_addr = rat_specific.wlan_context().mac_addr();
  }
  return mac_addr;
}

std::string EventsReporterImpl::get_imei(const SessionConfig& config) {
  // IMEI is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string imei = "";
  if (rat_specific.has_lte_context()) {
    imei = rat_specific.lte_context().imei();
  }
  return imei;
}

std::string EventsReporterImpl::get_spgw_ipv4(const SessionConfig& config) {
  // SPGW_IPV4 is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string spgw_ipv4 = "";
  if (rat_specific.has_lte_context()) {
    spgw_ipv4 = rat_specific.lte_context().spgw_ipv4();
  }
  return spgw_ipv4;
}

std::string EventsReporterImpl::get_user_location(const SessionConfig& config) {
  // UserLocation is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string user_location = "";
  if (rat_specific.has_lte_context()) {
    user_location = rat_specific.lte_context().user_location();
  }
  // Return the HEX values in string
  return magma::bytes_to_hex(user_location);
}

std::string EventsReporterImpl::get_charging_characteristics(
    const SessionConfig& config) {
  // Charging Characteristics is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string charging_characteristics = "";
  if (rat_specific.has_lte_context()) {
    charging_characteristics =
        rat_specific.lte_context().charging_characteristics();
  }
  // Return the HEX values in string
  return charging_characteristics;
}

}  // namespace lte
}  // namespace magma
