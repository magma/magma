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
#include <cstdint>
#include <string.h>
#include "M5GNetworkFeatureSupport.h"
#include "M5GCommonDefs.h"

namespace magma5g {
NetworkFeatureSupportMsg::NetworkFeatureSupportMsg(){};

NetworkFeatureSupportMsg::~NetworkFeatureSupportMsg(){};

int NetworkFeatureSupportMsg::EncodeNetworkFeatureSupportMsg(
    NetworkFeatureSupportMsg* feat_support, uint8_t iei, uint8_t* buffer,
    uint32_t len) {

  uint8_t encoded = 0;
  int i           = 0;

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) feat_support->iei);
    ENCODE_U8(buffer, iei, encoded);
  }

  ENCODE_U8(buffer + encoded, feat_support->len, encoded); 
  ENCODE_U8(buffer + encoded, feat_support->feature_list, encoded);
  ENCODE_U8(buffer + encoded, feat_support->spare, encoded); 

  return (encoded);
}

int NetworkFeatureSupportMsg::DecodeNetworkFeatureSupportMsg(
    NetworkFeatureSupportMsg* feat_support, uint8_t iei, uint8_t* buffer,
    uint32_t len) {

    // Decoding not supported
    return (RETURNerror);
}

} // namespace magma5g

