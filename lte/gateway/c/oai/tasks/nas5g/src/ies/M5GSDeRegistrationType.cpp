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
#include "5GSDeRegistrationType.h"
#include "CommonDefs.h"
using namespace std;
namespace magma5g
{
   M5GSDeRegistrationTypeMsg::M5GSDeRegistrationTypeMsg()
   {
   };

   M5GSDeRegistrationTypeMsg::~M5GSDeRegistrationTypeMsg()
   {
   };

   int M5GSDeRegistrationTypeMsg::DecodeM5GSDeRegistrationTypeMsg(M5GSDeRegistrationTypeMsg *deregistrationtype, uint8_t iei, uint8_t *buffer, uint32_t len) 
   {
      uint8_t decoded = 0;

      deregistrationtype->switchoff = (*(buffer + decoded) >> 3) & 0x01;
      deregistrationtype->reregistrationrequired = (*(buffer + decoded) >> 2) & 0x01;
      deregistrationtype->accesstype = *(buffer + decoded) & 0x03;
      decoded++;
      MLOG(MDEBUG) << "DecodeM5GSDe-RegistrationType : \n   switchoff = " << hex << int(deregistrationtype->switchoff) << endl;
      MLOG(MDEBUG) << "   reregistrationrequired = " << hex << int(deregistrationtype->reregistrationrequired) << endl;
      MLOG(MDEBUG) << "   accesstype = " << hex << int(deregistrationtype->accesstype) << endl;
      return (decoded);
   };


   int M5GSDeRegistrationTypeMsg::EncodeM5GSDeRegistrationTypeMsg(M5GSDeRegistrationTypeMsg *deregistrationtype, uint8_t iei, uint8_t * buffer, uint32_t len)
   {
      uint8_t encoded = 0;

      *(buffer + encoded) = 0x00 | ((deregistrationtype->switchoff << 3) & 0x08) |
                            ((deregistrationtype->reregistrationrequired << 2) & 0x04) |
                            (deregistrationtype->accesstype & 0x03);
      encoded++;
      MLOG(MDEBUG) << "In EncodeM5GSDeRegistrationTypeMsg___: DeRegistrationType= " << hex << int(*(buffer + encoded)) << endl;
      return (encoded);
   };
}

