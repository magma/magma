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

  Source      amf_app_desc.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_config.h"
#include "amf_app_ue_context.h"
using namespace std;
#pragma once

namespace magma5g
{
    class amf_app_desc_t 
    {
        public:
        /* UE contexts */
        amf_ue_context_t amf_ue_contexts;

        long m5_statistic_timer_id;
        uint32_t m5_statistic_timer_period;
        amf_ue_ngap_id_t amf_app_ue_ngap_id_generator;

        /* ***************Statistics*************
        * number of registered UE,number of connected UE,
        * number of idle UE,number of PDU Sessions,
        * number of NG_U PDU session,number of PDN sessions
        */

        uint32_t nb_gnb_connected;
        uint32_t nb_ue_registered;
        uint32_t nb_ue_connected;
            
    };
}