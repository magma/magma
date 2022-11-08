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
#include <sstream>
#include <cstdint>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowParam.hpp"

namespace magma5g {
M5GQosFlowParam::M5GQosFlowParam() {}
M5GQosFlowParam::~M5GQosFlowParam() {}

int M5GQosFlowParam::EncodeM5GQosFlowParam(M5GQosFlowParam* param,
                                           uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "EncodeQosFlowParamemeter: ");
  switch (param->iei) {
    case param_id_5qi: {
      *(buffer + encoded) = param->iei;
      OAILOG_DEBUG(LOG_NAS5G, "5Qi iei = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->length;
      OAILOG_DEBUG(LOG_NAS5G, "5Qi length = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->element;
      OAILOG_DEBUG(LOG_NAS5G, "5Qi = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
    } break;
    case param_id_mfbr_uplink: {
      *(buffer + encoded) = param->iei;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink iei = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->length;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->units;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink units = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = (param->element & 0Xff00) >> 8;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink element octet1 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
      *(buffer + encoded) = param->element & (0X00ff);
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink element octet2 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
    } break;
    case param_id_mfbr_downlink: {
      *(buffer + encoded) = param->iei;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink iei = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->length;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->units;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink units = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = (param->element & 0Xff00) >> 8;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink element octet1 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
      *(buffer + encoded) = param->element & (0X00ff);
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink element octet2 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
    } break;
    case param_id_gfbr_uplink: {
      *(buffer + encoded) = param->iei;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink iei = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->length;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->units;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink units = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = (param->element & 0Xff00) >> 8;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink element octet1 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
      *(buffer + encoded) = param->element & (0X00ff);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink element octet2 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
    } break;
    case param_id_gfbr_downlink: {
      *(buffer + encoded) = param->iei;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink iei = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->length;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = param->units;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink units = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;

      *(buffer + encoded) = (param->element & 0Xff00) >> 8;
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink element octet1 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
      *(buffer + encoded) = param->element & (0X00ff);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink element octet2 = 0X%x",
                   static_cast<int>(*(buffer + encoded)));
      encoded++;
    } break;
    default: {
    } break;
  }
  return encoded;
}

int M5GQosFlowParam::DecodeM5GQosFlowParam(M5GQosFlowParam* param,
                                           uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "DecodeM5GQosFlowParam :");
  param->iei = *(buffer + decoded);
  decoded++;

  switch (param->iei) {
    case param_id_5qi: {
      OAILOG_DEBUG(LOG_NAS5G, "5Qi iei = 0X%x", static_cast<int>(param->iei));
      param->length = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "5Qi length = 0X%x",
                   static_cast<int>(param->length));
      param->element = *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "5Qi = 0X%x", static_cast<int>(param->element));
      decoded++;
    } break;
    case param_id_gfbr_uplink: {
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink iei = 0X%x",
                   static_cast<int>(param->iei));
      param->length = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(param->length));
      param->units = *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink units = 0X%x",
                   static_cast<int>(param->units));
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_uplink element = 0X%x",
                   static_cast<int>(param->element));
      decoded++;
    } break;
    case param_id_gfbr_downlink: {
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink iei = 0X%x",
                   static_cast<int>(param->iei));
      param->length = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(param->length));
      param->units = *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink units = 0X%x",
                   static_cast<int>(param->units));
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "gfbr_downlink element = 0X%x",
                   static_cast<int>(param->element));
      decoded++;
    } break;
    case param_id_mfbr_uplink: {
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink iei = 0X%x",
                   static_cast<int>(param->iei));
      param->length = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(param->length));
      param->units = *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink units = 0X%x",
                   static_cast<int>(param->units));
      decoded++;
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_uplink element = 0X%x",
                   static_cast<int>(param->element));
      decoded++;
    } break;
    case param_id_mfbr_downlink: {
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink iei = 0X%x",
                   static_cast<int>(param->iei));
      param->length = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "iei length = 0X%x",
                   static_cast<int>(param->length));
      param->units = *(buffer + decoded);
      decoded++;
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink units = 0X%x",
                   static_cast<int>(param->units));
      param->element = *(buffer + decoded);
      decoded++;
      param->element |= *(buffer + decoded);
      OAILOG_DEBUG(LOG_NAS5G, "mfbr_downlink element = 0X%x",
                   static_cast<int>(param->element));
      decoded++;
    } break;
    default: {
    } break;
  }

  return decoded;
}

// Convert mfbr and gfbr values in proper format
void M5GQosFlowParam::mfbr_gbr_convert(
    magma5g::M5GQosFlowParam* flow_des_paramList, uint64_t element) {
  int count = 0;
  while (element > UINT16_MAX) {
    element /= 1024;
    if (count == 0) {
      count = M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_KBPS;
    } else {
      count += M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_COUNT;
    }
  }
  if (count == 0) {
    flow_des_paramList->element = (uint16_t)element / 1024;
    flow_des_paramList->units = M5G_QOS_FLOW_PARAM_BIT_RATE_UNITS_KBPS;
  } else {
    flow_des_paramList->element = (uint16_t)element;
    flow_des_paramList->units = count;
  }
}

}  // namespace magma5g
