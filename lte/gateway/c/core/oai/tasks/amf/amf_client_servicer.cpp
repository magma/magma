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

#include "include/amf_client_servicer.h"
#include "M5GAuthenticationServiceClient.h"
#include "amf_common.h"
#include <memory>
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.h"
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.h"

using magma5g::AsyncM5GAuthenticationServiceClient;

namespace magma5g {

status_code_e AMFClientServicerBase::amf_send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  OAILOG_INFO(
      LOG_AMF_APP, "Sending msg to :[%s] id: [%d]-[%s]\n",
      itti_get_task_name(destination_task_id), ITTI_MSG_ID(message),
      ITTI_MSG_NAME(message));
  return (send_msg_to_task(task_zmq_ctx_p, destination_task_id, message));
}

bool AMFClientServicerBase::get_subs_auth_info(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    amf_ue_ngap_id_t ue_id) {
  return (AsyncM5GAuthenticationServiceClient::getInstance().get_subs_auth_info(
      imsi, imsi_length, snni, ue_id));
}

bool AMFClientServicerBase::get_subs_auth_info_resync(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    const void* resync_info, uint8_t resync_info_len, amf_ue_ngap_id_t ue_id) {
  return (
      AsyncM5GAuthenticationServiceClient::getInstance()
          .get_subs_auth_info_resync(
              imsi, imsi_length, snni, resync_info, resync_info_len, ue_id));
}

int AMFClientServicerBase::allocate_ipv4_address(
    const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
    uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
    uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len) {
  return AsyncM5GMobilityServiceClient::getInstance().allocate_ipv4_address(
      subscriber_id, apn, pdu_session_id, pti, AF_INET, gnb_gtp_teid,
      gnb_gtp_teid_ip_addr, gnb_gtp_teid_ip_addr_len);
}

int AMFClientServicerBase::release_ipv4_address(
    const char* subscriber_id, const char* apn, const struct in_addr* addr) {
  return AsyncM5GMobilityServiceClient::getInstance().release_ipv4_address(
      subscriber_id, apn, addr);
}

int AMFClientServicerBase::amf_smf_create_pdu_session_ipv4(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version,
    const ambr_t& state_ambr) {
  return AsyncSmfServiceClient::getInstance().amf_smf_create_pdu_session_ipv4(
      imsi, apn, pdu_session_id, pdu_session_type, gnb_gtp_teid, pti,
      gnb_gtp_teid_ip_addr, ipv4_addr, version, state_ambr);
}

bool AMFClientServicerBase::set_smf_session(SetSMSessionContext& request) {
  return AsyncSmfServiceClient::getInstance().set_smf_session(request);
}

AMFClientServicer& AMFClientServicer::getInstance() {
  static AMFClientServicer instance;
  return instance;
}

status_code_e amf_send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  return (magma5g::AMFClientServicer::getInstance().amf_send_msg_to_task(
      task_zmq_ctx_p, destination_task_id, message));
}

}  // namespace magma5g
