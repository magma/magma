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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowParam.hpp"

namespace magma5g {

class M5GQosFlowDescription {
 public:
  uint8_t iei;
  uint16_t length;

#define MAX_QOS_FLOW_PARAMS_LIST 7
#define QOS_FLOW_DESC_BUF_LEN_MAX 4096

  uint8_t qfi;
  uint8_t operationCode;
  uint8_t numOfParams : 6;
  uint8_t Ebit : 2;
  M5GQosFlowParam paramList[MAX_QOS_FLOW_PARAMS_LIST];

  int EncodeM5GQosFlowDescription(M5GQosFlowDescription* qosFlowDesc,
                                  uint8_t* buffer, uint32_t len);
  int DecodeM5GQosFlowDescription(M5GQosFlowDescription* qosFlowDesc,
                                  uint8_t* buffer, uint32_t len);

  M5GQosFlowDescription();
  ~M5GQosFlowDescription();
};
}  // namespace magma5g
