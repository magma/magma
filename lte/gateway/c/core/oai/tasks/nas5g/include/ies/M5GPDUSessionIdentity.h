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
// PDUSessionIdentity IE Class
class PDUSessionIdentityMsg {
 public:
  uint8_t iei;
  uint8_t pdu_session_id;

  PDUSessionIdentityMsg();
  ~PDUSessionIdentityMsg();
  int EncodePDUSessionIdentityMsg(
      PDUSessionIdentityMsg* pdu_session_id, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodePDUSessionIdentityMsg(
      PDUSessionIdentityMsg* pdu_session_id, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
