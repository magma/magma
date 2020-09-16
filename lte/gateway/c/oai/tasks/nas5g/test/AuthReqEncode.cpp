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

/*** Using this stub code we are going to test Encoding Functionality of Authentication Request Message ***/

#include <iostream>
#include <iomanip>
#include <cstring>
#include "AuthenticationRequest.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g
{
   int Encode(void)
   {
      int ret = 0;
      uint8_t buffer[10] = {};
      int len = 10;

      AuthenticationRequestMsg AuthReq;
      AuthReq.extendedprotocoldiscriminator.extendedprotodiscriminator = 126;
      AuthReq.securityheadertype.securityhdr = 0;
      AuthReq.sparehalfoctet.spare = 0;
      AuthReq.messagetype.msgtype = 0x56;
      AuthReq.naskeysetidentifier.tsc = 0;
      AuthReq.naskeysetidentifier.naskeysetidentifier = 0;
      uint8_t abba_buff[] = {0x71, 0x00, 0x0d, 0x01};
      AuthReq.abba.contents.assign((const char *)abba_buff, 4);
      
      //Encoding the Authentication Message
      MLOG(MDEBUG) << "\n\n---Encoding Authentication request Message---\n\n";
      ret = AuthReq.EncodeAuthenticationRequestMsg(&AuthReq, buffer, len);

      //Printing Encoded Message
      MLOG(MDEBUG) << " ENCODED MESSAGE : " << setfill('0') << hex << int(buffer[0]) << hex << setw(2) << int(buffer[1]) << hex << int(buffer[2]) << hex << setw(2) << int(buffer[3])<< hex << setw(2) << int(buffer[4])<< hex << setw(2) << int(buffer[5])<< hex << setw(2) << int(buffer[6])<< hex << setw(2) << int(buffer[7])<< hex << setw(2) << int(buffer[8])<< "\n\n";

      return 0;
   }
}  

//Main Function to call Test Encode Function
int main(void)
{
   int ret;
   ret = magma5g::Encode();
   return 0;
   }
