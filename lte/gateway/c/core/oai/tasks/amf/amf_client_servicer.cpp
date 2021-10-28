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

using magma5g::AsyncM5GAuthenticationServiceClient;

namespace magma5g {

status_code_e AMFClientServicerBase::amf_send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
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
