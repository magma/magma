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
// PDUSessionType Class
class PDUSessionTypeMsg {
 public:
  uint8_t iei : 4;
  uint32_t type_val : 3;

  PDUSessionTypeMsg();
  ~PDUSessionTypeMsg();
  int EncodePDUSessionTypeMsg(
      PDUSessionTypeMsg* pdu_session_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodePDUSessionTypeMsg(
      PDUSessionTypeMsg* pdu_session_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
