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


#include <stdint.h>
#include <string.h>
#include <iostream>
#include <string>

#include "proto_msg_to_itti_msg.h"
#include "3gpp_33.401.h"
#include "common_types.h"
#include "feg/protos/s6a_proxy.pb.h"
#include "security_types.h"

extern "C" {
}

namespace magma {

void convert_proto_msg_to_itti_s6a_auth_info_ans(
  AuthenticationInformationAnswer msg,
  s6a_auth_info_ans_t *itti_msg)
{
  if (msg.eutran_vectors_size() > MAX_EPS_AUTH_VECTORS) {
    std::cout << "[ERROR] Number of eutran auth vectors received is:"
                 << msg.eutran_vectors_size() << std::endl;
    return;
  }
  itti_msg->auth_info.nb_of_vectors = msg.eutran_vectors_size();
  uint8_t idx = 0;
  while (idx < itti_msg->auth_info.nb_of_vectors) {
    auto eutran_vector = msg.eutran_vectors(idx);
    eutran_vector_t *itti_eutran_vector =
      &(itti_msg->auth_info.eutran_vector[idx]);
    if (eutran_vector.rand().length() <= RAND_LENGTH_OCTETS) {
      memcpy(
        itti_eutran_vector->rand,
        eutran_vector.rand().c_str(),
        eutran_vector.rand().length());
    }
    uint8_t xres_len = 0;
    xres_len = eutran_vector.xres().length();
    if ((xres_len > XRES_LENGTH_MIN) && (xres_len <= XRES_LENGTH_MAX)) {
      itti_eutran_vector->xres.size = eutran_vector.xres().length();
      memcpy(
        itti_eutran_vector->xres.data, eutran_vector.xres().c_str(), xres_len);
    } else {
      std::cout << "[ERROR] Invalid xres length " << xres_len << std::endl;
      return;
    }
    if (eutran_vector.autn().length() == AUTN_LENGTH_OCTETS) {
      memcpy(
        itti_eutran_vector->autn,
        eutran_vector.autn().c_str(),
        eutran_vector.autn().length());
    } else {
      std::cout << "[ERROR] Invalid AUTN length " << eutran_vector.autn().length()
                   << std::endl;
      return;
    }
    if (eutran_vector.kasme().length() == KASME_LENGTH_OCTETS) {
      memcpy(
        itti_eutran_vector->kasme,
        eutran_vector.kasme().c_str(),
        eutran_vector.kasme().length());
    } else {
      std::cout << "[ERROR] Invalid KASME length " << eutran_vector.kasme().length()
                   << std::endl;
      return;
    }
    ++idx;
  }
  return;
}

void convert_proto_msg_to_itti_s6a_update_location_ans(
  UpdateLocationAnswer msg,
  s6a_update_location_ans_t *itti_msg)
{
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
      itti_msg->subscription_data.msisdn,
      msg.msisdn().c_str(),
      msg.msisdn().length());
    itti_msg->subscription_data.msisdn_length = msg.msisdn().length();
  }
  itti_msg->subscription_data.subscriber_status = SS_SERVICE_GRANTED;
  itti_msg->subscription_data.access_restriction =
    ARD_HO_TO_NON_3GPP_NOT_ALLOWED;

  if (msg.network_access_mode()
      == UpdateLocationAnswer_NetworkAccessMode_PACKET_AND_CIRCUIT) {
    itti_msg->subscription_data.access_mode = NAM_PACKET_AND_CIRCUIT;
  } else if (msg.network_access_mode()
      == UpdateLocationAnswer_NetworkAccessMode_RESERVED) {
    itti_msg->subscription_data.access_mode = NAM_RESERVED;
  } else {
    itti_msg->subscription_data.access_mode = NAM_ONLY_PACKET;
  }

#define SUBSCRIBER_PERIODIC_RAU_TAU_TIMER_VAL 10
  itti_msg->subscription_data.rau_tau_timer =
    SUBSCRIBER_PERIODIC_RAU_TAU_TIMER_VAL;

  // apn configuration
  itti_msg->subscription_data.apn_config_profile.nb_apns = msg.apn_size();
  uint8_t idx = 0;
  while (idx < msg.apn_size() && idx < MAX_APN_PER_UE) {
    auto apn = msg.apn(idx);
    struct apn_configuration_s *itti_msg_apn =
      &(itti_msg->subscription_data.apn_config_profile.apn_configuration[idx]);

    itti_msg_apn->context_identifier = apn.context_id();
    itti_msg_apn->pdn_type = (pdn_type_t) apn.pdn();
    auto service_sel = apn.service_selection();
    if (service_sel.length() > APN_MAX_LENGTH) {
      itti_msg_apn->service_selection_length = APN_MAX_LENGTH;
    } else {
      itti_msg_apn->service_selection_length = service_sel.length();
    }
    memcpy(
      itti_msg_apn->service_selection,
      service_sel.c_str(),
      itti_msg_apn->service_selection_length);

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

    //apn ambr
    itti_msg_apn->ambr.br_ul = apn.ambr().max_bandwidth_ul();
    itti_msg_apn->ambr.br_dl = apn.ambr().max_bandwidth_dl();
    ++idx;
  }
  return;
}

} // namespace magma
