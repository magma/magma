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
#define SST_LENGTH 1
#define SD_LENGTH 3
#define NSSAI_MSG_IE_MIN_LEN 2
namespace magma5g {
class NSSAIMsg {
 public:
  const int NSSAI_VAL_MAX = 74;
  const int NSSAI_MIN_LENGTH = 4;
  uint8_t iei;
  uint8_t len;
  uint8_t sst;
  uint8_t sd[SD_LENGTH];
  uint8_t hplmn_sst;
  uint8_t hplmn_sd[SD_LENGTH];

  NSSAIMsg();
  ~NSSAIMsg();
  int EncodeNSSAIMsg(NSSAIMsg* nssai, uint8_t iei, uint8_t* buffer,
                     uint32_t len);
  int DecodeNSSAIMsg(NSSAIMsg* nssai, uint8_t iei, uint8_t* buffer,
                     uint32_t len);
};

class NSSAIMsgList {
 public:
  uint8_t iei;
  uint8_t len;
  NSSAIMsg nssai;

  NSSAIMsgList();
  ~NSSAIMsgList();
  int EncodeNSSAIMsgList(NSSAIMsgList* allowed_nssai, uint8_t iei,
                         uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
