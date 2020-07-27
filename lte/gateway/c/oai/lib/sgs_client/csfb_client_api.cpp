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

#include <grpcpp/impl/codegen/status.h>
#include <iostream>
#include <string>

#include "csfb_client_api.h"
#include "CSFBClient.h"
#include "orc8r/protos/common.pb.h"

void empty_callback(grpc::Status status, magma::Void void_response) {
  return;
}

void send_alert_ack(const itti_sgsap_alert_ack_t* msg) {
  magma::CSFBClient::alert_ack(msg, empty_callback);
}

void send_alert_reject(const itti_sgsap_alert_reject_t* msg) {
  magma::CSFBClient::alert_reject(msg, empty_callback);
}

void send_location_update_request(const itti_sgsap_location_update_req_t* msg) {
  std::cout << "[DEBUG] Sending LOCATION_UDPATE_REQUEST with IMSI: "
            << std::string(msg->imsi) << std::endl;
  magma::CSFBClient::location_update_request(
      msg, [imsiStr = std::string(msg->imsi)](
               grpc::Status status, magma::Void void_response) {
        if (status.ok()) {
          std::cout
              << "[DEBUG] Successfully sent LOCATION_UDPATE_REQUEST with IMSI: "
              << imsiStr << std::endl;
        } else {
          std::cout
              << "[ERROR] Failed to send LOCATION_UDPATE_REQUEST with IMSI: "
              << imsiStr << "; Status: " << status.error_message() << std::endl;
        }
        return;
      });
}

void send_tmsi_reallocation_complete(
    const itti_sgsap_tmsi_reallocation_comp_t* msg) {
  std::cout << "[DEBUG] Sending TMSI_REALLOCATION_COMPLETE with IMSI: "
            << std::string(msg->imsi) << std::endl;
  magma::CSFBClient::tmsi_reallocation_complete(
      msg, [imsiStr = std::string(msg->imsi)](
               grpc::Status status, magma::Void void_response) {
        if (status.ok()) {
          std::cout
              << "[DEBUG] "
              << "Successfully sent TMSI_REALLOCATION_COMPLETE with IMSI: "
              << imsiStr << std::endl;
        } else {
          std::cout
              << "[ERROR] Failed to send TMSI_REALLOCATION_COMPLETE with IMSI: "
              << imsiStr << "; Status: " << status.error_message() << std::endl;
        }
        return;
      });
}

void send_eps_detach_indication(const itti_sgsap_eps_detach_ind_t* msg) {
  magma::CSFBClient::eps_detach_indication(msg, empty_callback);
}

void send_imsi_detach_indication(const itti_sgsap_imsi_detach_ind_t* msg) {
  magma::CSFBClient::imsi_detach_indication(msg, empty_callback);
}

void send_paging_reject(const itti_sgsap_paging_reject_t* msg) {
  magma::CSFBClient::paging_reject(msg, empty_callback);
}

void send_service_request(const itti_sgsap_service_request_t* msg) {
  std::cout << "[DEBUG] Sending SERVICE REQUEST with IMSI: "
            << std::string(msg->imsi) << std::endl;
  magma::CSFBClient::service_request(msg, empty_callback);
}

void send_ue_activity_indication(const itti_sgsap_ue_activity_ind_t* msg) {
  magma::CSFBClient::ue_activity_indication(msg, empty_callback);
}

void send_ue_unreachable(const itti_sgsap_ue_unreachable_t* msg) {
  magma::CSFBClient::ue_unreachable(msg, empty_callback);
}

void send_uplink_unitdata(const itti_sgsap_uplink_unitdata_t* msg) {
  magma::CSFBClient::send_uplink_unitdata(msg, empty_callback);
}
