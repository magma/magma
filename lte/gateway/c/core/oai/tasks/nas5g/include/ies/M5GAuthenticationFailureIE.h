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
// M5GMMCause IE Class
class M5GAuthenticationFailureIE {
 public:
  /* Optional Parameters */
  uint8_t iei;
  bstring authentication_failure_info;
  uint32_t presencemask;

  M5GAuthenticationFailureIE();
  ~M5GAuthenticationFailureIE();
  int DecodeM5GAuthenticationFailureIE(
      M5GAuthenticationFailureIE* m5g_auth_failure_ie, uint8_t iei,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
