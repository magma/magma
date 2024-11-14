/*
 * Copyright 2022 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNetworkFeatureSupport.hpp"

namespace magma5g {
NetworkFeatureSupportMsg::NetworkFeatureSupportMsg() {}
NetworkFeatureSupportMsg::~NetworkFeatureSupportMsg() {}

int NetworkFeatureSupportMsg::DecodeNetworkFeatureSupportMsg(uint8_t iei,
                                                             uint8_t* buffer,
                                                             uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    this->iei = *(buffer + decoded++);
  }
  this->len = *(buffer + decoded++);

  this->MPSI = (*(buffer + decoded) >> 7) & 0x01;
  this->IWK_N26 = (*(buffer + decoded) >> 6) & 0x01;
  this->EMF = (*(buffer + decoded) >> 4) & 0x03;
  this->EMC = (*(buffer + decoded) >> 2) & 0x03;
  this->IMS_VoPS_N3GPP = (*(buffer + decoded) >> 1) & 0x01;
  this->IMS_VoPS_3GPP = *(buffer + decoded) & 0x01;
  decoded++;

  this->MCSI = (*(buffer + decoded) >> 1) & 0x01;
  this->EMCN3 = (*(buffer + decoded)) & 0x01;
  decoded++;

  return decoded;
}

int NetworkFeatureSupportMsg::EncodeNetworkFeatureSupportMsg(uint8_t iei,
                                                             uint8_t* buffer,
                                                             uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, NETWORK_FEATURE_MINIMUM_LENGTH,
                                       len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)this->iei);
    *(buffer + encoded++) = iei;
  }

  *(buffer + encoded++) = this->len;
  *(buffer + encoded++) =
      0x00 | ((this->MPSI & 0x01) << 7) | ((this->IWK_N26 & 0x01) << 6) |
      ((this->EMF & 0x03) << 4) | ((this->EMC & 0x03) << 2) |
      ((this->IMS_VoPS_N3GPP & 0x01) << 1) | (this->IMS_VoPS_3GPP & 0x01);
  *(buffer + encoded++) =
      0x00 | ((this->MCSI & 0x01) << 1) | (this->EMCN3 & 0x01);
  return encoded;
}

}  // namespace magma5g
