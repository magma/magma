/*
   Copyright 2022 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#pragma once
#include <sstream>
#include <cstdint>

namespace magma5g {

class M5GQosFlowParam {
 public:
  typedef enum qosflow_param_id_type {
    param_id_5qi = 1,
    param_id_gfbr_uplink,
    param_id_gfbr_downlink,
    param_id_mfbr_uplink,
    param_id_mfbr_downlink,
    param_id_avg_window,
    param_id_qos_flow_identity,
  } qos_flow_param_id_type_t;

  uint8_t iei;

#define M5G_QOS_FLOW_PARAM_BIT_RATE_LEN 3
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_KBPS 1
  uint8_t length;
  uint8_t units;
  uint16_t element;

  M5GQosFlowParam();
  ~M5GQosFlowParam();

  int EncodeM5GQosFlowParam(M5GQosFlowParam* param, uint8_t* buffer,
                            uint32_t len);
  int DecodeM5GQosFlowParam(M5GQosFlowParam* param, uint8_t* buffer,
                            uint32_t len);
};
}  // namespace magma5g
