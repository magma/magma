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

#include <stdint.h>
#include <string.h>
#include <iostream>
#include <string>

#include "proto_msg_to_itti_msg.h"
#include "3gpp_33.401.h"
#include "common_types.h"
#include "feg/protos/s6a_proxy.pb.h"
#include "security_types.h"

extern "C" {}

namespace magma {

void copy_charging_characteristics(
    charging_characteristics_t* target, const char* proto_c_str, int length) {
  if (length > CHARGING_CHARACTERISTICS_LENGTH) {
    length = CHARGING_CHARACTERISTICS_LENGTH;
  }
  if (length) memcpy(target->value, proto_c_str, length);
  target->value[length] = '\0';
  target->length        = length;
}

void convert_proto_msg_to_itti_s6a_auth_info_ans(
    AuthenticationInformationAnswer msg, s6a_auth_info_ans_t* itti_msg) {
  if (msg.eutran_vectors_size() > MAX_EPS_AUTH_VECTORS) {
    std::cout << "[ERROR] Number of eutran auth vectors received is:"
              << msg.eutran_vectors_size() << std::endl;
    return;
  }
  itti_msg->auth_info.nb_of_vectors = msg.eutran_vectors_size();
  uint8_t idx                       = 0;
  while (idx < itti_msg->auth_info.nb_of_vectors) {
    auto eutran_vector = msg.eutran_vectors(idx);
    eutran_vector_t* itti_eutran_vector =
        &(itti_msg->auth_info.eutran_vector[idx]);
    if (eutran_vector.rand().length() <= RAND_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->rand, eutran_vector.rand().c_str(),
          eutran_vector.rand().length());
    }
    uint8_t xres_len = 0;
    xres_len         = eutran_vector.xres().length();
    if ((xres_len > XRES_LENGTH_MIN) && (xres_len <= XRES_LENGTH_MAX)) {
      itti_eutran_vector->xres.size = eutran_vector.xres().length();
      memcpy(
          itti_eutran_vector->xres.data, eutran_vector.xres().c_str(),
          xres_len);
    } else {
      std::cout << "[ERROR] Invalid xres length " << xres_len << std::endl;
      return;
    }
    if (eutran_vector.autn().length() == AUTN_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->autn, eutran_vector.autn().c_str(),
          eutran_vector.autn().length());
    } else {
      std::cout << "[ERROR] Invalid AUTN length "
                << eutran_vector.autn().length() << std::endl;
      return;
    }
    if (eutran_vector.kasme().length() == KASME_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->kasme, eutran_vector.kasme().c_str(),
          eutran_vector.kasme().length());
    } else {
      std::cout << "[ERROR] Invalid KASME length "
                << eutran_vector.kasme().length() << std::endl;
      return;
    }
    ++idx;
  }
  return;
}

void convert_proto_msg_to_itti_s6a_update_location_ans(
    UpdateLocationAnswer msg, s6a_update_location_ans_t* itti_msg) {
  itti_msg->subscription_data.apn_config_profile.context_identifier =
      msg.default_context_id();
  itti_msg->subscription_data.subscribed_ambr.br_ul =
      msg.total_ambr().max_bandwidth_ul();
  itti_msg->subscription_data.subscribed_ambr.br_dl =
      msg.total_ambr().max_bandwidth_dl();
  if (msg.all_apns_included()) {
    itti_msg->subscription_data.apn_config_profile.all_apn_conf_ind =
        MODIFIED_ADDED_APN_CONFIGURATIONS_INCLUDED;
  } else {
    itti_msg->subscription_data.apn_config_profile.all_apn_conf_ind =
        ALL_APN_CONFIGURATIONS_INCLUDED;
  }
  if (msg.msisdn().length() <= (MSISDN_LENGTH + 1)) {
    memcpy(
        itti_msg->subscription_data.msisdn, msg.msisdn().c_str(),
        msg.msisdn().length());
    itti_msg->subscription_data.msisdn_length = msg.msisdn().length();
  }

  copy_charging_characteristics(
      &itti_msg->subscription_data.default_charging_characteristics,
      msg.default_charging_characteristics().c_str(),
      msg.default_charging_characteristics().length());

  itti_msg->subscription_data.subscriber_status = SS_SERVICE_GRANTED;
  itti_msg->subscription_data.access_restriction =
      ARD_HO_TO_NON_3GPP_NOT_ALLOWED;

  if (msg.network_access_mode() ==
      UpdateLocationAnswer_NetworkAccessMode_PACKET_AND_CIRCUIT) {
    itti_msg->subscription_data.access_mode = NAM_PACKET_AND_CIRCUIT;
  } else if (
      msg.network_access_mode() ==
      UpdateLocationAnswer_NetworkAccessMode_RESERVED) {
    itti_msg->subscription_data.access_mode = NAM_RESERVED;
  } else {
    itti_msg->subscription_data.access_mode = NAM_ONLY_PACKET;
  }

  // Regional subscription zone codes
  itti_msg->subscription_data.num_zcs =
      ((msg.regional_subscription_zone_code_size() > MAX_REGIONAL_SUB) ?
           MAX_REGIONAL_SUB :
           (msg.regional_subscription_zone_code_size()));
  uint8_t reg_sub_idx = 0;
  for (uint8_t itr = 0; itr < itti_msg->subscription_data.num_zcs; itr++) {
    auto regional_subscription_zone_code =
        msg.regional_subscription_zone_code(itr);
    // Copy only if the zonecode is 2 octets
    if (regional_subscription_zone_code.length() == ZONE_CODE_LEN) {
      memcpy(
          itti_msg->subscription_data.reg_sub[reg_sub_idx].zone_code,
          regional_subscription_zone_code.c_str(), ZONE_CODE_LEN);
      ++reg_sub_idx;
    } else {
      std::cout << "[WARNING] Invalid zonecode length received, ignoring"
                << regional_subscription_zone_code.length() << std::endl;
    }
  }
  itti_msg->subscription_data.num_zcs = reg_sub_idx;

#define SUBSCRIBER_PERIODIC_RAU_TAU_TIMER_VAL 10
  itti_msg->subscription_data.rau_tau_timer =
      SUBSCRIBER_PERIODIC_RAU_TAU_TIMER_VAL;

  // apn configuration
  uint8_t nb_apns = 0;
  if (msg.apn_size() > MAX_APN_PER_UE) {
    std::cout << "[WARNING] The number of APNs configured in subscriber data ("
              << msg.apn_size() << ") is larger than MME limit of "
              << MAX_APN_PER_UE << ". Truncating the list to this MME limit."
              << std::endl;
    nb_apns = MAX_APN_PER_UE;
  } else {
    nb_apns = msg.apn_size();
  }
  itti_msg->subscription_data.apn_config_profile.nb_apns = nb_apns;
  for (uint8_t idx = 0; idx < nb_apns; ++idx) {
    auto apn                                 = msg.apn(idx);
    struct apn_configuration_s* itti_msg_apn = &(
        itti_msg->subscription_data.apn_config_profile.apn_configuration[idx]);

    itti_msg_apn->context_identifier = apn.context_id();
    itti_msg_apn->pdn_type           = (pdn_type_t) apn.pdn();

    auto service_sel = apn.service_selection();
    if (service_sel.length() > APN_MAX_LENGTH) {
      itti_msg_apn->service_selection_length = APN_MAX_LENGTH;
    } else {
      itti_msg_apn->service_selection_length = service_sel.length();
    }
    memcpy(
        itti_msg_apn->service_selection, service_sel.c_str(),
        itti_msg_apn->service_selection_length);

    copy_charging_characteristics(
        &itti_msg_apn->charging_characteristics,
        apn.charging_characteristics().c_str(),
        apn.charging_characteristics().length());

    // Qos profile
    itti_msg_apn->subscribed_qos.qci = (qci_t) apn.qos_profile().class_id();
    itti_msg_apn->subscribed_qos.allocation_retention_priority.priority_level =
        apn.qos_profile().priority_level();
    itti_msg_apn->subscribed_qos.allocation_retention_priority
        .pre_emp_vulnerability = (pre_emption_vulnerability_t) apn.qos_profile()
                                     .preemption_vulnerability();
    itti_msg_apn->subscribed_qos.allocation_retention_priority
        .pre_emp_capability =
        (pre_emption_capability_t) apn.qos_profile().preemption_capability();

    // apn ambr
    itti_msg_apn->ambr.br_ul   = apn.ambr().max_bandwidth_ul();
    itti_msg_apn->ambr.br_dl   = apn.ambr().max_bandwidth_dl();
    itti_msg_apn->ambr.br_unit = (apn_ambr_bitrate_unit_t) apn.ambr().unit();
  }

  return;
}

}  // namespace magma
