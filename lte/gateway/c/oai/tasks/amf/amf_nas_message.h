/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************

  Source      amf_nas_message.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>

using namespace std;
#pragma once
namespace mangam5g
{
////////////////////typecast require///////////////////////////////
    class nas_message_decode_status_t 
    {
        public:
        uint8_t integrity_protected_message : 1;
        uint8_t ciphered_message : 1;
        uint8_t mac_matched : 1;
        uint8_t security_context_available : 1;
        int emm_cause;
    };

}