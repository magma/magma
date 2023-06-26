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
#pragma once
#include <sstream>
#include <cstdint>
namespace magma5g {
class NetworkFeatureSupportMsg {
 public:
#define NETWORK_FEATURE_MINIMUM_LENGTH 3
#define NETWORK_FEATURE_MAXIMUM_LENGTH 5

  uint8_t iei;
  uint8_t len;
  uint8_t IMS_VoPS_3GPP : 1;
  uint8_t IMS_VoPS_N3GPP : 1;
  uint8_t EMC : 2;
  uint8_t EMF : 2;
  uint8_t IWK_N26 : 1;
  uint8_t MPSI : 1;
  uint8_t EMCN3 : 1;
  uint8_t MCSI : 1;
  NetworkFeatureSupportMsg();
  ~NetworkFeatureSupportMsg();

  int EncodeNetworkFeatureSupportMsg(NetworkFeatureSupportMsg* networkfeature,
                                     uint8_t iei, uint8_t* buffer,
                                     uint32_t len);

  int DecodeNetworkFeatureSupportMsg(NetworkFeatureSupportMsg* networkfeature,
                                     uint8_t iei, uint8_t* buffer,
                                     uint32_t len);
};
}  // namespace magma5g
