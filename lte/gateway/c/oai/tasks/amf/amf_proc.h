**
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

  Source      amf_proc.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "../bstr/bstrlib.h"
using namespace std;
#pragma once

namespace magma5g
{
    class amf_registration_request_ies_t: public registration_request_msg
    {
        // need to put registration ies.


    };

    enum amf_proc_registration_type_t
    {
        AMF_REGISTRATION_TYPE_INITIAL = 0, 
        AMF_REGISTRATION_TYPE_MOBILITY_UPDATING,
        AMF_REGISTRATION_TYPE_PERODIC_UPDATING,
        AMF_REGISTRATION_TYPE_EMERGENCY,
        AMF_REGISTRATION_TYPE_RESERVED
    } ;

//int amf_proc_registration_request( amf_ue_ngap_id_t ue_id, const bool ctx_is_new, amf_registration_request_ies_t *const params);

int amf_registration_reject(amf_context_t *amf_context, struct nas_base_proc_s *nas_base_proc);

int amf_proc_registration_reject(amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause);

int amf_proc_registration_complete(amf_ue_ngap_id_t ue_id, const_bstring amf_msg_pP, int amf_cause, const nas_message_decode_status_t status);

}


/*
0 0 1 initial registration
0 1 0 mobility registration updating
0 1 1 periodic registration updating
1 0 0 emergency registration
*/