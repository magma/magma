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
#include "bstrlib.h"
using namespace std;

namespace magma5g {
// AuthenticationResponseParameter IE Class
class AuthenticationResponseParameterMsg {
 public:
#define AUTHENTICATION_RESPONSE_PARAMETER_MIN_LEN 18
  uint8_t iei;
  uint8_t length;
  // bstring response_parameter;
  unsigned char response_parameter[16];

  AuthenticationResponseParameterMsg();
  ~AuthenticationResponseParameterMsg();
  int EncodeAuthenticationResponseParameterMsg(
      AuthenticationResponseParameterMsg* response_parameter, uint8_t iei,
      uint8_t* buffer, uint32_t len);
  int DecodeAuthenticationResponseParameterMsg(
      AuthenticationResponseParameterMsg* response_parameter, uint8_t iei,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
