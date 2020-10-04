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

  Source      amf_as.h

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
    class amf_as: public AMFMsg
    {
        public:

        void amf_as_initialize(void);

        int amf_as_send(amf_as_t *msg);

        int amf_as_send_ng(const amf_as_t *msg);

        static AMGMsg* amf_as_set_header(nas_message_t* msg, const amf_as_security_data_t* security);
       


    };
    enum nas_error_code_t
    {
        AS_SUCCESS = 1,          /* Success code, transaction is going on    */
        AS_TERMINATED_NAS,       /* Transaction terminated by NAS        */
        AS_TERMINATED_AS,        /* Transaction terminated by AS         */
        AS_NON_DELIVERED_DUE_HO, /* Failure code                 */
        AS_FAILURE               /* Failure code, stand also for lower
                                  * layer failure AS_LOWER_LAYER_FAILURE */
    };
} 