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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowParam.h"

namespace magma5g {
M5GQosFlowParam::M5GQosFlowParam() {}
M5GQosFlowParam::~M5GQosFlowParam() {}

int M5GQosFlowParam::EncodeM5GQosFlowParam(M5GQosFlowParam* param,
                                           uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  MLOG(MDEBUG) << " EncodeQosFlowParamemeter : ";
  switch (param->iei) {
    case param_id_5qi: {
      *(buffer + encoded) = param->iei;
      MLOG(MDEBUG) << "5Qi iei = 0x" << std::hex << int(*(buffer + encoded));
      encoded++;

      *(buffer + encoded) = param->length;
      MLOG(MDEBUG) << "5Qi length= 0x" << std::hex << int(*(buffer + encoded));
      encoded++;

      *(buffer + encoded) = param->element;
      MLOG(MDEBUG) << "5Qi  = 0x" << std::hex << int(*(buffer + encoded));
      encoded++;
    } break;
    case param_id_mfbr_uplink:
    case param_id_mfbr_downlink:
    case param_id_gfbr_uplink:
    case param_id_gfbr_downlink: {
      *(buffer + encoded) = param->iei;
      MLOG(MDEBUG) << "mfbr_uplink iei = 0x" << std::hex
                   << int(*(buffer + encoded));
      encoded++;

      *(buffer + encoded) = param->length;
      MLOG(MDEBUG) << "iei length= 0x" << std::hex << int(*(buffer + encoded));
      encoded++;

      *(buffer + encoded) = param->units;
      MLOG(MDEBUG) << "mfbr_uplink units  = 0x" << std::hex
                   << int(*(buffer + encoded));
      encoded++;

      *(buffer + encoded) = (param->element & 0xff00) >> 8;
      MLOG(MDEBUG) << "Element Octet1  = 0x" << std::hex
                   << int(*(buffer + encoded));
      encoded++;
      *(buffer + encoded) = param->element & (0x00ff);
      MLOG(MDEBUG) << "Element Octet2 = 0x" << std::hex
                   << int(*(buffer + encoded));
      encoded++;
    } break;
    default: {
    }
  }
  return encoded;
}

int M5GQosFlowParam::DecodeM5GQosFlowParam(M5GQosFlowParam* param,
                                           uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;

  MLOG(MDEBUG) << " DecodeM5GQosFlowParam : ";
  param->iei = *(buffer + decoded);
  decoded++;

  switch (param->iei) {
    case param_id_5qi: {
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->iei);
      param->length = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->length);
      param->element = *(buffer + decoded);
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->element);
      decoded++;
    } break;
    case param_id_gfbr_uplink: {
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->iei);
      param->length = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->length);
      param->units = *(buffer + decoded);
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->units);
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded) << 8;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->element);
      decoded++;
    } break;
    case param_id_gfbr_downlink: {
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->iei);
      param->length = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->length);
      param->units = *(buffer + decoded);
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->units);
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded) << 8;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->element);
      decoded++;
    } break;
    case param_id_mfbr_uplink: {
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->iei);
      param->length = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->length);
      param->units = *(buffer + decoded);
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->units);
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded) << 8;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->element);
      decoded++;
    } break;
    case param_id_mfbr_downlink: {
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->iei);
      param->length = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->length);
      param->units = *(buffer + decoded);
      decoded++;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->units);
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded) << 8;
      MLOG(MDEBUG) << "IEI = 0x" << std::hex << int(param->element);
      decoded++;
    } break;
    default: {
    } break;
  }

  return decoded;
}

}  // namespace magma5g
