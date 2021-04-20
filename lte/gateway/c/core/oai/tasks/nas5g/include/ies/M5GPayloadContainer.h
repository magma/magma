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
#include <SmfMessage.h>

using namespace std;
namespace magma5g {
class PayloadContainerMsg {
 public:
#define PAYLOAD_CONTAINER_CONTENTS_MAX_LEN 65534
  uint8_t iei;
  uint32_t len;
  uint8_t contents[PAYLOAD_CONTAINER_CONTENTS_MAX_LEN];
  SmfMsg smf_msg;

  PayloadContainerMsg();
  ~PayloadContainerMsg();
  int EncodePayloadContainerMsg(
      PayloadContainerMsg* payload_container, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodePayloadContainerMsg(
      PayloadContainerMsg* payload_container, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
