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

namespace magma5g {
/**
 * This class is single place holder for all client related services.
 * For instance : subscriberdb, sessiond, mobilityd
 */

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
};

class AMFClientServicer : public AMFClientServicerBase {
 public:
  static AMFClientServicer& getInstance();

  AMFClientServicer(AMFClientServicer const&) = delete;
  void operator=(AMFClientServicer const&) = delete;

#if MME_UNIT_TEST
  status_code_e amf_send_msg_to_task(
      task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
      MessageDef* message_p) override {
    OAILOG_DEBUG(LOG_AMF_APP, " Mock is Enabled \n");

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
#endif /* MME_UNIT_TEST */

 private:
  AMFClientServicer(){};
};

}  // namespace magma5g
