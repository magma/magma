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

#include <iostream>
#include <sstream>
#include <cstdint>
#include <cstring>
#include "PayloadContainerType.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
PayloadContainerTypeMsg::PayloadContainerTypeMsg(){};
PayloadContainerTypeMsg::~PayloadContainerTypeMsg(){};

// Decode PayloadContainerType IE
int PayloadContainerTypeMsg::DecodePayloadContainerTypeMsg(
    PayloadContainerTypeMsg* payloadcontainertype, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  payloadcontainertype->typeval = (*buffer & 0x0f);
  decoded++;
  MLOG(MDEBUG) << "DecodePayloadContainerTypeMsg__: typeval = " << hex
               << int(payloadcontainertype->typeval) << endl;

  return (decoded);
};

// Encode PayloadContainerType IE
int PayloadContainerTypeMsg::EncodePayloadContainerTypeMsg(
    PayloadContainerTypeMsg* payloadcontainertype, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  *buffer = payloadcontainertype->typeval & 0x0f;
  MLOG(MDEBUG) << "DecodePayloadContainerTypeMsg__: typeval = " << hex
               << int(*buffer) << endl;
  encoded++;

  return (encoded);
};
}  // namespace magma5g
