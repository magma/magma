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
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_MBPS 6
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_GBPS 11
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_TBPS 16
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_PBPS 21
#define M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_COUNT \
  5  // Considering only 1KBPS,1MBPS,1GBPS,1TBPS,1PBPS as per the specs.
  uint8_t length;
  uint8_t units;
  uint16_t element;

  M5GQosFlowParam();
  ~M5GQosFlowParam();

  void mfbr_gbr_convert(magma5g::M5GQosFlowParam* flow_des_paramList,
                        uint64_t element);
  int EncodeM5GQosFlowParam(M5GQosFlowParam* param, uint8_t* buffer,
                            uint32_t len);
  int DecodeM5GQosFlowParam(M5GQosFlowParam* param, uint8_t* buffer,
                            uint32_t len);
};
}  // namespace magma5g
