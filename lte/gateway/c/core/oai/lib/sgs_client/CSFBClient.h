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

#pragma once

#include <gmp.h>
#include <grpc++/grpc++.h>
#include <stdint.h>
#include <string>
#include <functional>
#include <memory>

#include "includes/GRPCReceiver.h"
#include "feg/protos/csfb.grpc.pb.h"
#include "sgs_messages_types.h"

extern "C" {
#include "intertask_interface.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
}

namespace magma {
using namespace orc8r;
using namespace feg;

/**
 * CSFBClient is the main client for sending message to FeG
 * FeG will forward the message to MSC then respond instantly with Void
 */
class CSFBClient : public GRPCReceiver {
 public:
  /**
   * Send SGsAP-ALERT-ACK
   */
  static void alert_ack(
      const itti_sgsap_alert_ack_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * Send SGsAP-ALERT-REJECT
   */
  static void alert_reject(
      const itti_sgsap_alert_reject_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * Send SGsAP-LOCATION-UPDATE-REQUEST
   */
  static void location_update_request(
      const itti_sgsap_location_update_req_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * Send SGsAP-TMSI-REALLOCATION-COMPLETE
   */
  static void tmsi_reallocation_complete(
      const itti_sgsap_tmsi_reallocation_comp_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-EPS-DETACH-INDICATION
   */
  static void eps_detach_indication(
      const itti_sgsap_eps_detach_ind_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-IMSI-DETACH-INDICATION
   */
  static void imsi_detach_indication(
      const itti_sgsap_imsi_detach_ind_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-PAGING-REJECT
   */
  static void paging_reject(
      const itti_sgsap_paging_reject_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-SERVICE-REQUEST
   */
  static void service_request(
      const itti_sgsap_service_request_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-UE-ACTIVITY-INDICATION
   */
  static void ue_activity_indication(
      const itti_sgsap_ue_activity_ind_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-UE-UNREACHABLE
   */
  static void ue_unreachable(
      const itti_sgsap_ue_unreachable_t* msg,
      std::function<void(grpc::Status, Void)> callback);

  /**
   * SGsAP-UPLINK-UNITDATA
   */
  static void send_uplink_unitdata(
      const itti_sgsap_uplink_unitdata_t* msg,
      std::function<void(grpc::Status, Void)> callback);

 public:
  CSFBClient(CSFBClient const&) = delete;
  void operator=(CSFBClient const&) = delete;

 private:
  CSFBClient();
  static CSFBClient& get_instance();
  std::unique_ptr<CSFBFedGWService::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 3;  // seconds
};

}  // namespace magma
