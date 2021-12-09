/*
   Copyright 2020 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#include <sstream>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowDescriptor.h"

namespace magma5g {
M5GQosFlowDescription::M5GQosFlowDescription() {}
M5GQosFlowDescription::~M5GQosFlowDescription() {}

int M5GQosFlowDescription::EncodeM5GQosFlowDescription(
    M5GQosFlowDescription* qosFlowDesc, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  MLOG(MDEBUG) << " EncodeQosFlowDescriptor : ";
  *(buffer + encoded) = qosFlowDesc->qfi;
  MLOG(MDEBUG) << "QFI = 0x" << std::hex << int(*(buffer + encoded));
  encoded++;

  *(buffer + encoded) = qosFlowDesc->operationCode;
  MLOG(MDEBUG) << "OperationCode = 0x" << std::hex << int(*(buffer + encoded));
  encoded++;

  *(buffer + encoded) = qosFlowDesc->Ebit << 6;
  *(buffer + encoded) |= qosFlowDesc->numOfParams & 0x3f;
  MLOG(MDEBUG) << "NumOfParams = 0x" << std::hex << int(*(buffer + encoded));
  encoded++;
  for (uint8_t i = 0; i < qosFlowDesc->numOfParams; i++) {
    M5GQosFlowParam* qosParams = &qosFlowDesc->paramList[i];
    encoded +=
        qosParams->EncodeM5GQosFlowParam(qosParams, (buffer + encoded), len);
  }

  return encoded;
}

int M5GQosFlowDescription::DecodeM5GQosFlowDescription(
    M5GQosFlowDescription* qosFlowDesc, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  MLOG(MDEBUG) << " EncodeQosFlowDescriptor : ";
  qosFlowDesc->qfi = (*(buffer + decoded)) & 0x3F;
  MLOG(MDEBUG) << "QFI = 0x" << std::hex << int(qosFlowDesc->qfi);
  decoded++;
  qosFlowDesc->operationCode = (*(buffer + decoded) & 0xE0);
  MLOG(MDEBUG) << "OperationCode = 0x" << std::hex
               << int(qosFlowDesc->operationCode);
  decoded++;
  qosFlowDesc->numOfParams = ((*(buffer + decoded) & 0x3F));
  MLOG(MDEBUG) << "NumOfParams = 0x" << std::hex
               << int(qosFlowDesc->numOfParams);
  qosFlowDesc->Ebit = ((*(buffer + decoded) & 0x40) >> 6);
  MLOG(MDEBUG) << "Ebit = 0x" << std::hex << int(qosFlowDesc->Ebit);
  decoded++;

  for (uint8_t i = 0; i < qosFlowDesc->numOfParams; i++) {
    M5GQosFlowParam* qosParams = &qosFlowDesc->paramList[i];
    decoded +=
        qosParams->DecodeM5GQosFlowParam(qosParams, (buffer + decoded), len);
  }

  return decoded;
}
}  // namespace magma5g
