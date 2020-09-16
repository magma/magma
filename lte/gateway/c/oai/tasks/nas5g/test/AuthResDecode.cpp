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

/*** Using this stub code we are going to test Decoding Functionality of Authentication Response Message ***/

#include <iostream>
#include <cstring>
#include "AuthenticationResponse.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g
{
   int Decode(void)
   {
      int ret = 0;
      uint8_t buffer[] = {0x7E, 0x00, 0x57};
      int len = 10;
      AuthenticationResponseMsg AuthRes;
      
      //Decoding Authentication Response Message
      MLOG(MDEBUG) << " ---Authentication response Message---\n";
      ret = AuthRes.DecodeAuthenticationResponseMsg (&AuthRes, buffer, len);

      //Printing Decoded Authentication Response Message
      MLOG(MDEBUG) << " ---Decoded Message---\n";
      MLOG(MDEBUG) << " Extended Protocol Discriminator :" << dec << int(AuthRes.extendedprotocoldiscriminator.extendedprotodiscriminator);
      MLOG(MDEBUG) << " Spare Half Octet : " << dec << int(AuthRes.sparehalfoctet.spare);
      MLOG(MDEBUG) << " Security Header Type : " << dec << int(AuthRes.securityheadertype.securityhdr);
      MLOG(MDEBUG) << " Message Type : 0x" << hex << int(AuthRes.messagetype.msgtype);

      return 0;
   }
}  

//Main Function to call test Decode function
int main(void)
{
   int ret;
   ret = magma5g::Decode();
   return 0;
   }
