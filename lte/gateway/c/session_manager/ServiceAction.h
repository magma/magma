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

#include <memory>

#include <lte/protos/session_manager.grpc.pb.h>
#include <experimental/optional>

#include "CreditKey.h"
#include "Types.h"

namespace magma {
using namespace lte;
using std::experimental::optional;

enum ServiceActionType {
  CONTINUE_SERVICE = 0,
  TERMINATE_SERVICE = 1,
  ACTIVATE_SERVICE = 2,
  REDIRECT = 3,
  RESTRICT_ACCESS = 4,
};

/**
 * ServiceAction is the base class for any action that needs to be taken on
 * a subscriber's service. This could be terminate, redirect data, or just
 * continue.
 */
class ServiceAction {
 public:
  ServiceAction(ServiceActionType action_type) : action_type_(action_type) {}

  ServiceActionType get_type() const { return action_type_; }

  ServiceAction& set_imsi(const std::string& imsi) {
    imsi_ = std::make_unique<std::string>(imsi);
    return *this;
  }

  ServiceAction& set_session_id(const std::string& session_id) {
    session_id_ = std::make_unique<std::string>(session_id);
    return *this;
  }

  ServiceAction& set_ip_addr(const std::string& ip_addr) {
    ip_addr_ = std::make_unique<std::string>(ip_addr);
    return *this;
  }

  ServiceAction& set_ipv6_addr(const std::string& ipv6_addr) {
    ipv6_addr_ = std::make_unique<std::string>(ipv6_addr);
    return *this;
  }

  ServiceAction& set_teids(const Teids& teids) {
    teids_ = std::make_unique<Teids>(teids);
    return *this;
  }

  ServiceAction& set_credit_key(const CreditKey& credit_key) {
    credit_key_ = credit_key;
    return *this;
  }

  ServiceAction& set_ambr(const optional<AggregatedMaximumBitrate> ambr) {
    ambr_ = ambr;
    return *this;
  }

  ServiceAction& set_msisdn(const std::string& msisdn) {
    msisdn_ = std::make_unique<std::string>(msisdn);
    return *this;
  }

  /**
   * get_imsi returns the associated IMSI for the action, or throws a nullptr
   * exception if there is none stored
   */
  const std::string& get_imsi() const { return *imsi_; }

  /**
   * get_imsi returns the associated IMSI for the action, or throws a nullptr
   * exception if there is none stored
   */
  const std::string& get_session_id() const { return *session_id_; }

  /**
   * get_ip_addr returns the associated subscriber's ip_addr for the action,
   * or throws a nullptr exception if there is none stored
   */
  const std::string& get_ip_addr() const { return *ip_addr_; }

  const std::string& get_ipv6_addr() const { return *ipv6_addr_; }

  const Teids& get_teids() const { return *teids_; }

  const CreditKey& get_credit_key() const { return credit_key_; }

  const optional<AggregatedMaximumBitrate> get_ambr() const { return ambr_; }

  const std::string& get_msisdn() const { return *msisdn_; }

  // RulesToProcess
  RulesToProcess get_gx_rules_to_install() const { return gx_to_install_; }
  RulesToProcess* get_mutable_gx_rules_to_install() { return &gx_to_install_; }

  RulesToProcess get_gy_rules_to_install() const { return gy_to_install_; }
  RulesToProcess* get_mutable_gy_rules_to_install() { return &gy_to_install_; }

 private:
  ServiceActionType action_type_;
  std::unique_ptr<std::string> imsi_;
  std::unique_ptr<std::string> session_id_;
  std::unique_ptr<std::string> ip_addr_;
  std::unique_ptr<std::string> ipv6_addr_;
  std::unique_ptr<Teids> teids_;
  std::unique_ptr<std::string> msisdn_;
  CreditKey credit_key_;
  optional<AggregatedMaximumBitrate> ambr_;
  RulesToProcess gx_to_install_;
  RulesToProcess gy_to_install_;
};

}  // namespace magma
