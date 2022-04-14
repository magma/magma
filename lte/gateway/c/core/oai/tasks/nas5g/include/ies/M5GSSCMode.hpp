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

namespace magma5g {
// SSCMode IE Class
class SSCModeMsg {
 public:
  uint8_t iei : 4;
  uint32_t mode_val : 3;

  SSCModeMsg();
  ~SSCModeMsg();
  int EncodeSSCModeMsg(SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer,
                       uint32_t len);
  int DecodeSSCModeMsg(SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer,
                       uint32_t len);
  void copy(const SSCModeMsg& s) {
    iei = s.iei;
    mode_val = s.mode_val;
  }
  bool isEqual(const SSCModeMsg& s) {
    return ((iei == s.iei) && (mode_val == s.mode_val));
  }
};
}  // namespace magma5g
