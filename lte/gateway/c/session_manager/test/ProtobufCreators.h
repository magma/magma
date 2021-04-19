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
#pragma once

#include <lte/protos/mconfig/mconfigs.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <string>
#include <vector>

#include "DiameterCodes.h"
#include "StoredState.h"

namespace magma {
using namespace lte;

CommonSessionContext build_common_context(
    const std::string& imsi, const std::string& ue_ipv4,
    const std::string& ue_ipv6, const Teids teids, const std::string& apn,
    const std::string& msisdn, const RATType rat_type);

LTESessionContext build_lte_context(
    const std::string& spgw_ipv4, const std::string& imei,
    const std::string& plmn_id, const std::string& imsi_plmn_id,
    const std::string& user_location, uint32_t bearer_id,
    QosInformationRequest* qos_info);

WLANSessionContext build_wlan_context(
    const std::string& mac_addr, const std::string& radius_session_id);

RuleSet create_rule_set(
    const bool apply_subscriber_wide, const std::string& apn,
    std::vector<std::string> static_rules,
    std::vector<PolicyRule> dynamic_rules);

void create_rule_record(
    const std::string& imsi, const std::string& rule_id, uint64_t bytes_rx,
    uint64_t bytes_tx, RuleRecord* rule_record);

void create_rule_record(
    const std::string& imsi, const uint32_t teid, const std::string& rule_id,
    uint64_t bytes_rx, uint64_t bytes_tx, RuleRecord* rule_record);

void create_rule_record(
    const std::string& imsi, const uint32_t teid, const std::string& rule_id,
    uint64_t bytes_rx, uint64_t bytes_tx, uint64_t dropped_rx,
    uint64_t dropped_tx, RuleRecord* rule_record);

void create_charging_credit(
    uint64_t volume, bool is_final, ChargingCredit* credit);

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, CreditLimitType limit_type,
    CreditUpdateResponse* response);

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, CreditUpdateResponse* response);

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, ChargingCredit_FinalAction action,
    std::string redirect_server, std::string restrict_rule,
    CreditUpdateResponse* response);

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t total_volume, uint64_t tx_volume,
    uint64_t rx_volume, bool is_final, CreditUpdateResponse* response);

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, bool is_final,
    CreditUpdateResponse* response);

void create_credit_update_response_with_error(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, bool is_final, DiameterResultCode resultCode,
    CreditUpdateResponse* response);

void create_credit_update_response_with_error(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, bool is_final, DiameterResultCode resultCode,
    ChargingCredit_FinalAction action, std::string redirect_server,
    std::string restrict_rule, CreditUpdateResponse* response);

void create_charging_credit(
    uint64_t total_volume, uint64_t tx_volume, uint64_t rx_volume,
    bool is_final, ChargingCredit* credit);

void create_charging_credit(
    uint64_t volume, ChargingCredit_FinalAction action,
    std::string redirect_server, std::string restrict_rule,
    ChargingCredit* credit);

void create_monitor_credit(
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    UsageMonitoringCredit* response);

void create_monitor_credit(
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume, UsageMonitoringCredit* credit);

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume,
    UsageMonitoringUpdateResponse* response);

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    UsageMonitoringUpdateResponse* response);

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    UsageMonitoringUpdateResponse* response);

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    UsageMonitoringUpdateResponse* response);

void create_usage_update(
    const std::string& imsi, uint32_t charging_key, uint64_t bytes_rx,
    uint64_t bytes_tx, CreditUsage::UpdateType type, CreditUsageUpdate* update);

void create_update_session_request(
    std::string imsi, std::string session_id, uint32_t ckey, std::string mkey,
    CreditUsage::UpdateType type, uint64_t bytes_rx, uint64_t bytes_tx,
    UpdateSessionRequest* usr);

void create_usage_monitoring_update_request(
    const std::string& imsi, std::string monitoring_key, uint64_t bytes_rx,
    uint64_t bytes_tx, UsageMonitoringUpdateRequest* update);

void create_policy_reauth_request(
    const std::string& session_id, const std::string& imsi,
    const std::vector<std::string>& rules_to_remove,
    const std::vector<StaticRuleInstall>& rules_to_install,
    const std::vector<DynamicRuleInstall>& dynamic_rules_to_install,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    const std::vector<UsageMonitoringCredit>& usage_monitoring_credits,
    PolicyReAuthRequest* request);

void create_tgpp_context(
    const std::string& gx_dest_host, const std::string& gy_dest_host,
    TgppContext* context);

void create_subscriber_quota_update(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state, SubscriberQuotaUpdate* update);

void create_session_create_response(
    const std::string& imsi, const std::string session_id,
    const std::string& monitoring_key, std::vector<std::string>& static_rules,
    CreateSessionResponse* response);

PolicyRule create_policy_rule(
    const std::string& rule_id, const std::string& m_key, const uint32_t rg);

PolicyRule create_policy_rule_with_qos(
    const std::string& rule_id, const std::string& m_key, const uint32_t rg,
    const int qci);

void create_granted_units(
    uint64_t* total, uint64_t* tx, uint64_t* rx, GrantedUnits* gsu);

magma::mconfig::SessionD get_default_mconfig();

PolicyBearerBindingRequest create_policy_bearer_bind_req(
    const std::string& imsi, const uint32_t linked_bearer_id,
    const std::string& rule_id, const uint32_t bearer_id,
    const uint32_t agw_teid, const uint32_t enb_teid);

UpdateTunnelIdsRequest create_update_tunnel_ids_request(
    const std::string& imsi, const uint32_t bearer_id, const Teids teids);

UpdateTunnelIdsRequest create_update_tunnel_ids_request(
    const std::string& imsi, const uint32_t bearer_id, const uint32_t agw_teid,
    const uint32_t enb_teid);

StaticRuleInstall create_static_rule_install(const std::string& rule_id);

DynamicRuleInstall create_dynamic_rule_install(const PolicyRule& rule);

}  // namespace magma
