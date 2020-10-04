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

  Source      amf_data.cpp

  Version     0.1

  Date        2020/07/28

  Product     AMF 

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "amf_data.h"
#include "log.h"

using namespace std;
#pragma once
typedef uint32_t amf_ue_ngap_id_t;
#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH) LENGTH += QUADLET - (LENGTH % QUADLET)

namespace magma5g
{
    void amf_ctx_set_attribute_valid(amf_context_t* const ctxt, const int attribute_bit_pos) 
    {
        ctxt->member_present_mask |= attribute_bit_pos;
        ctxt->member_valid_mask |= attribute_bit_pos;
    }
    /* Set IMSI, mark it as valid */
    void amf_ctx_set_valid_imsi(amf_context_t* const ctxt, imsi_t* imsi, const imsi64_t imsi64) {
    ctxt->_imsi   = *imsi;
    ctxt->_imsi64 = imsi64;
    amf_ctx_set_attribute_valid(ctxt, AMF_CTXT_MEMBER_IMSI);
    #if DEBUG_IS_ON
    char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
    IMSI64_TO_STRING(ctxt->_imsi64, imsi_str, ctxt->_imsi.length);
    OAILOG_DEBUG(LOG_NAS_AMF, "ue_id=" AMF_UE_NGAP_ID_FMT " set IMSI %s (valid)\n",
        (PARENT_STRUCT(ctxt, struct ue_m5gmm_context_s, amf_context))
            ->amf_ue_ngap_id,
        imsi_str);
    #endif
    //TODO 
    //amf_api_notify_imsi((PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id, imsi64);
    }
}//namespace magma5g