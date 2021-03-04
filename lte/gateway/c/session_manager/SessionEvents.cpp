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

#include "SessionEvents.h"

using magma::orc8r::Event;
using magma::orc8r::Void;

namespace {  // anonymous

const std::string SESSIOND_SERVICE_EV       = "sessiond";
const std::string SESSION_CREATED_EV        = "session_created";
const std::string SESSION_CREATE_FAILURE_EV = "session_create_failure";
const std::string SESSION_UPDATED_EV        = "session_updated";
const std::string SESSION_UPDATE_FAILURE_EV = "session_update_failure";
const std::string SESSION_TERMINATED_EV     = "session_terminated";

const std::string SESSION_ID     = "session_id";
const std::string IMSI           = "imsi";
const std::string IMEI           = "imei";
const std::string IP_ADDR        = "ip_addr";
const std::string IPV6_ADDR      = "ipv6_addr";
const std::string MAC_ADDR       = "mac_addr";
const std::string MSISDN         = "msisdn";
const std::string SPGW_IP        = "spgw_ip";
const std::string APN            = "apn";
const std::string PDP_START_TIME = "pdp_start_time";
const std::string PDP_END_TIME   = "pdp_end_time";
const std::string FAILURE_REASON = "failure_reason";

const std::string TOTAL_TX      = "total_tx";
const std::string TOTAL_RX      = "total_rx";
const std::string CHARGING_TX   = "charging_tx";
const std::string CHARGING_RX   = "charging_rx";
const std::string MONITORING_TX = "monitoring_tx";
const std::string MONITORING_RX = "monitoring_rx";

}  // namespace

namespace magma {
namespace lte {

EventsReporterImpl::EventsReporterImpl(AsyncEventdClient& eventd_client)
    : eventd_client_(eventd_client) {}

void EventsReporterImpl::session_created(
    const std::string& imsi, const std::string& session_id,
    const SessionConfig& session_context,
    const std::unique_ptr<SessionState>& session) {
  auto event = magma::orc8r::Event();
  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_CREATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value  = folly::dynamic::object;
  event_value[IMSI]           = imsi;
  event_value[IP_ADDR]        = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR]      = session_context.common_context.ue_ipv6();
  event_value[MSISDN]         = session_context.common_context.msisdn();
  event_value[APN]            = session_context.common_context.apn();
  event_value[SESSION_ID]     = session_id;
  event_value[PDP_START_TIME] = session->get_pdp_start_time();
  // LTE specific
  event_value[IMEI]    = get_imei(session_context);
  event_value[SPGW_IP] = get_spgw_ipv4(session_context);
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
    const std::string& imsi, const SessionConfig& session_context,
    const std::string& failure_reason) {
  auto event = magma::orc8r::Event();
  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_CREATE_FAILURE_EV);
  event.set_tag(imsi);

  folly::dynamic event_value  = folly::dynamic::object;
  event_value[IMSI]           = imsi;
  event_value[APN]            = session_context.common_context.apn();
  event_value[FAILURE_REASON] = failure_reason;
  event_value[MAC_ADDR]       = get_mac_addr(session_context);

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

void EventsReporterImpl::session_updated(
    const std::string& imsi, const std::string& session_id,
    const SessionConfig& session_context) {
  auto event = magma::orc8r::Event();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_UPDATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI]          = imsi;
  event_value[SESSION_ID]    = session_id;
  event_value[IP_ADDR]       = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR]     = session_context.common_context.ue_ipv6();
  event_value[APN]           = session_context.common_context.apn();
  event_value[MAC_ADDR]      = get_mac_addr(session_context);

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
    const std::string& imsi, const std::string& session_id,
    const SessionConfig& session_context, const std::string& failure_reason) {
  auto event = magma::orc8r::Event();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_UPDATE_FAILURE_EV);

  folly::dynamic event_value  = folly::dynamic::object;
  event_value[IMSI]           = imsi;
  event_value[SESSION_ID]     = session_id;
  event_value[IP_ADDR]        = session_context.common_context.ue_ipv4();
  event_value[IPV6_ADDR]      = session_context.common_context.ue_ipv6();
  event_value[MAC_ADDR]       = get_mac_addr(session_context);
  event_value[APN]            = session_context.common_context.apn();
  event_value[FAILURE_REASON] = failure_reason;

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
  auto event       = magma::orc8r::Event();
  auto session_cfg = session->get_config();

  event.set_stream_name(SESSIOND_SERVICE_EV);
  event.set_event_type(SESSION_TERMINATED_EV);
  event.set_tag(imsi);

  folly::dynamic event_value = folly::dynamic::object;
  event_value[IMSI]          = imsi;
  event_value[IP_ADDR]       = session_cfg.common_context.ue_ipv4();
  event_value[IPV6_ADDR]     = session_cfg.common_context.ue_ipv6();
  event_value[MSISDN]        = session_cfg.common_context.msisdn();
  event_value[APN]           = session_cfg.common_context.apn();
  event_value[SESSION_ID]    = session->get_session_id();

  SessionCredit::TotalCreditUsage usage = session->get_total_credit_usage();
  event_value[TOTAL_TX]       = usage.charging_tx + usage.monitoring_tx;
  event_value[TOTAL_RX]       = usage.charging_rx + usage.monitoring_rx;
  event_value[CHARGING_TX]    = usage.charging_tx;
  event_value[CHARGING_RX]    = usage.charging_rx;
  event_value[MONITORING_TX]  = usage.monitoring_tx;
  event_value[MONITORING_RX]  = usage.monitoring_rx;
  event_value[PDP_START_TIME] = session->get_pdp_start_time();
  event_value[PDP_END_TIME]   = session->get_pdp_end_time();
  // LTE specific
  event_value[IMEI]    = get_imei(session_cfg);
  event_value[SPGW_IP] = get_spgw_ipv4(session_cfg);
  // CWF specific
  event_value[MAC_ADDR] = get_mac_addr(session_cfg);

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

std::string EventsReporterImpl::get_mac_addr(const SessionConfig& config) {
  // MacAddr is only relevant for WLAN
  const auto& rat_specific = config.rat_specific_context;
  std::string mac_addr     = "";
  if (rat_specific.has_wlan_context()) {
    mac_addr = rat_specific.wlan_context().mac_addr();
  }
  return mac_addr;
}

std::string EventsReporterImpl::get_imei(const SessionConfig& config) {
  // IMEI is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string imei         = "";
  if (rat_specific.has_lte_context()) {
    imei = rat_specific.lte_context().imei();
  }
  return imei;
}

std::string EventsReporterImpl::get_spgw_ipv4(const SessionConfig& config) {
  // SPGW_IPV4 is only relevant for LTE
  const auto& rat_specific = config.rat_specific_context;
  std::string spgw_ipv4    = "";
  if (rat_specific.has_lte_context()) {
    spgw_ipv4 = rat_specific.lte_context().spgw_ipv4();
  }
  return spgw_ipv4;
}

}  // namespace lte
}  // namespace magma
