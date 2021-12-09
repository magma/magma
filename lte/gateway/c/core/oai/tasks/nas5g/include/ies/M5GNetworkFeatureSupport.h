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
#pragma once
#include <sstream>
#include <cstdint>
using namespace std;
namespace magma5g {
#define NAS5G_NETWORK_FEATURE_IE 0x21
#define NAS5G_FEATURE_SUPPORT_MIN_LENGTH 2
#define NAS5G_IMS_VOICE_OVER_PS_SESSION_INDICATOR 1
class NetworkFeatureSupportMsg {
 public:
  uint8_t iei;
  uint8_t len;
  uint8_t feature_list;
  uint8_t spare;
  NetworkFeatureSupportMsg();
  ~NetworkFeatureSupportMsg();
  int EncodeNetworkFeatureSupportMsg(
      NetworkFeatureSupportMsg* feat_support, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeNetworkFeatureSupportMsg(
      NetworkFeatureSupportMsg* feat_support, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
