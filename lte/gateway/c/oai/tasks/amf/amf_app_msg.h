/****************************************************************************
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
 ****************************************************************************/
/*****************************************************************************

  Source      amf_app_msg.h

  Version     0.1

  Date        2020/09/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include "amf_config.h"
using namespace std;
#pragma once

namespace magmam5g
{
    class amf_app_msg
    {
        public:
        void amf_app_ue_context_release(class ue_m5gmm_context_s* ue_context_p, enum ngcause cause) ;
    };
}//magmam5g