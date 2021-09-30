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

#include "include/ngap_client_servicer.h"
#include <memory>
extern "C" {
#include "common_defs.h"
}

namespace magma5g {

NGAPClientServicer::NGAPClientServicer() {
  mock_feature = false;
}

NGAPClientServicer& NGAPClientServicer::getInstance() {
  static NGAPClientServicer instance;

  return instance;
}

void NGAPClientServicer::set_mock_feature(bool val) {
  mock_feature = val;
}

bool NGAPClientServicer::is_mock_feature_enabled() {
  return mock_feature;
}

status_code_e NGAPClientServicer::send_message_to_amf(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  OAILOG_DEBUG(
      LOG_NGAP, " Is Mock Enabled %d\n",
      NGAPClientServicer::getInstance().is_mock_feature_enabled());

  if (NGAPClientServicer::getInstance().is_mock_feature_enabled()) {
    free(message);
    return (RETURNok);
  }

  return (send_msg_to_task(task_zmq_ctx_p, destination_task_id, message));
}

}  // namespace magma5g
