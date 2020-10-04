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

  Source      amf_app_defs.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_config.h"
#include "amf_as.h"
using namespace std;
#pragma once

namespace magma5g
{
    class amf_app_defs:amf_app_ue_context
    {
        public:
        imsi64_t amf_app_handle_initial_ue_message(amf_app_desc_t *amf_app_desc_p,itti_ngap_initial_ue_message_t *const conn_est_ind_pP);
        int amf_app_handle_nas_dl_req(const amf_ue_ngap_id_t ue_id, bstring nas_msg,nas_error_code_t transaction_status);
       

    };
}