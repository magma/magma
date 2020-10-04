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

  Source      amf_recv.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_data.h"
using namespace std;
#pragma once

namespace magma5g
{
    class amf_procedure_handler
    {

        public:
        int amf_handle_registration_request(msg->ue_id,&originating_tai,&msg->ecgi,&amf_msg->registration_request,msg->is_initial,msg->is_mm_ctx_new,amf_cause,&decode_status);
    };
    class amf_registration_procedure: public amf_context_t
    {
      public:
      static int amf_proc_registration_request(amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,amf_registration_request_ies_t* const ies);
      static int amf_registration_run_procedure(&ue_mm_context->amf_context);
      int amf_registration_success_identification_cb(amf_context);
      static int amf_registration(amf_context_t *amf_context);
      static int amf_proc_registration_reject(const amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause);

    }

}