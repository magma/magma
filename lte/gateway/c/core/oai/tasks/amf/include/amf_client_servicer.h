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

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "dynamic_memory_check.h"
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#ifdef __cplusplus
}
#endif

#include <memory>
#include <string>
#include "lte/protos/session_manager.grpc.pb.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/gateway/c/core/oai/include/map.h"

using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContext;
using magma::lte::SmContextVoid;

namespace magma5g {
/**
 * This class is single place holder for all client related services.
 * For instance : subscriberdb, sessiond, mobilityd
 */
typedef magma::map_s<std::string, std::string> map_string_string_t;

class AMFClientServicerBase {
 public:
  virtual status_code_e amf_send_msg_to_task(
      task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
      MessageDef* message);

  virtual bool get_subs_auth_info(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      amf_ue_ngap_id_t ue_id);

  virtual bool get_subs_auth_info_resync(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      const void* resync_info, uint8_t resync_info_len, amf_ue_ngap_id_t ue_id);

  virtual int allocate_ipv4_address(
      const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
      uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
      uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len);

  virtual int release_ipv4_address(
      const char* subscriber_id, const char* apn, const struct in_addr* addr);

  virtual int amf_smf_create_pdu_session_ipv4(
      char* imsi, uint8_t* apn, uint32_t pdu_session_id,
      uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
      uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version,
      const ambr_t& state_ambr);

  virtual bool set_smf_session(SetSMSessionContext& request);
  virtual bool get_decrypt_imsi_info(
      const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
      const std::string& ciphertext, const std::string& mac_tag,
      amf_ue_ngap_id_t ue_id);
};

class AMFClientServicer : public AMFClientServicerBase {
 public:
  std::vector<MessagesIds>
      msgtype_stack;  // stack maintains type of msgs sent to ngap
  static AMFClientServicer& getInstance();

  AMFClientServicer(AMFClientServicer const&) = delete;
  void operator=(AMFClientServicer const&) = delete;

  map_string_string_t map_tableKey_protoStr;

#if MME_UNIT_TEST
  status_code_e amf_send_msg_to_task(
      task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
      MessageDef* message_p) override {
    OAILOG_DEBUG(LOG_AMF_APP, " Mock is Enabled \n");
    msgtype_stack.push_back(ITTI_MSG_ID(message_p));
    itti_free_msg_content(message_p);
    free(message_p);
    return RETURNok;
  }
  bool get_subs_auth_info(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      amf_ue_ngap_id_t ue_id) override {
    return true;
  }

  bool get_subs_auth_info_resync(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      const void* resync_info, uint8_t resync_info_len,
      amf_ue_ngap_id_t ue_id) override {
    return true;
  }

  int allocate_ipv4_address(
      const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
      uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
      uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len) {
    return RETURNok;
  }

  int release_ipv4_address(
      const char* subscriber_id, const char* apn, const struct in_addr* addr) {
    return RETURNok;
  }

  int amf_smf_create_pdu_session_ipv4(
      char* imsi, uint8_t* apn, uint32_t pdu_session_id,
      uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
      uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version,
      const ambr_t& state_ambr) {
    return RETURNok;
  }

  bool set_smf_session(SetSMSessionContext& request) { return true; }
  bool get_decrypt_imsi_info(
      const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
      const std::string& ciphertext, const std::string& mac_tag,
      amf_ue_ngap_id_t ue_id) override {
    return true;
  }
#endif /* MME_UNIT_TEST */

 private:
  AMFClientServicer(){};
};

}  // namespace magma5g
