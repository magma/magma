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

  Source      amf_msgdef.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_message.h"
namespace std;
/* Header length boundaries of 5GS Mobility Management messages  */
#define AMF_HEADER_LENGTH sizeof(amf_msg_header)
#define AMF_HEADER_MINIMUM_LENGTH AMF_HEADER_LENGTH
#define AMF_HEADER_MAXIMUM_LENGTH AMF_HEADER_LENGTH
namespace magma5g
{
	class amf_msg_header 
	{
		public:
		uint8_t extended_protocol_discriminator;
		uint8_t security_header_type ;
		uint8_t message_type;
	};
	/* union of plain NAS message */
	union nas_message_plain_t 
	{
  		 AMFMsg amf; /* 5GMM Mobility Management messages */
  		 SMFMsg smf; /* 5GS Session Management messages  */
	};

	/* class of security protected NAS message */
	class nas_message_security_protected_t 
	{
		public:
		amf_msg_header header;
		nas_message_plain_t plain;
	};
	union nas_message_t 
	{
		amf_msg_header header;
		nas_message_security_protected_t security_protected;
		nas_message_plain_t plain;
	} ;

}//namespace magma5gs