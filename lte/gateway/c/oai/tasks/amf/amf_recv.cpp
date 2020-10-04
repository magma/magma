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

  Source      amf_recv.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_proc.h"
#include "5GSMobileIdentity.h"
#include "registration_request.h"
#include "amf_recv.h"
#include "amf_common_defs.h"
#include "amf_sap.h"
using namespace std;
#pragma once
#define AMF_CAUSE_SUCCESS (-1)
namespace magma5g
{
        int amf_procedure_handler::amf_handle_registration_request
        (const amf_ue_ngap_id_t ue_id, const tai_t *const originating_tai,
        const ecgi_t *const originating_ecgi,registration_request_msg *const msg,
        const bool is_initial,const bool is_mm_ctx_new,int *const amf_cause,
        const nas_message_decode_status_t *decode_status)
        {
            int rc = RETURNok;

            /*
            * Message checking
            */
            if (msg->uenetworkcapability.spare != 0b000) 
            {
                /*
                * Spare bits shall be coded as zero
                */
                *amf_cause = AMF_CAUSE_PROTOCOL_ERROR;                
            }
            /*
            * Handle message checking error
            */
            if (*amf_cause != AMF_CAUSE_SUCCESS) 
            {
                
               rc = amf_proc_registration_reject(ue_id, *amf_cause);
                              
            }

            amf_registration_request_ies_t *params = new(amf_registration_request_ies_t);
            /*
            * Message processing
            */
            /*
            * Get the 5GS Registration type
            */
            params->m5gsregistrationtype=AMF_REGISTRATION_TYPE_RESERVED;
            if (msg->m5gsregistrationtype == AMF_REGISTRATION_TYPE_INITIAL) 
            {
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;

            } 
            else if (msg->m5gsregistrationtype == AMF_REGISTRATION_TYPE_MOBILITY_UPDATING) 
            {                
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_MOBILITY_UPDATING;
            } 
            else if (msg->m5gsregistrationtype == AMF_REGISTRATION_TYPE_PERODIC_UPDATING) 
            {
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_PERODIC_UPDATING;
            } 
            else if (msg->m5gsregistrationtype == AMF_REGISTRATION_TYPE_EMERGENCY) 
            {
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_EMERGENCY;
            } 
            else if (msg->m5gsregistrationtype == AMF_REGISTRATION_TYPE_RESERVED) 
            {
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
            } 
            else 
            {
                params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
            }

            /*
            * Get the AMF mobile identity
            */
		    if (msg->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti.typeofidentity == M5GS_Mobile_Identity_M5G_GUTI) 
            {
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti = new(guti_t); // need to define guti_t like below in 3gpp 23003.h
                 /*!< \brief  Globally Unique MME Identity gummei_t gummei;            */
                   /*!< \brief  M-Temporary Mobile Subscriber Identity tmsi_t m_tmsi;  */

                /* below need to update after header file of 5gsmobileidenty.h completion*/
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mcc_digit1 = msg->oldgutiorimsi.guti.mcc_digit1;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mcc_digit2 = msg->oldgutiorimsi.guti.mcc_digit2;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mcc_digit3 = msg->oldgutiorimsi.guti.mcc_digit3;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mnc_digit1 = msg->oldgutiorimsi.guti.mnc_digit1;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mnc_digit2 = msg->oldgutiorimsi.guti.mnc_digit2;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.plmn.mnc_digit3 = msg->oldgutiorimsi.guti.mnc_digit3;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.mme_gid = msg->oldgutiorimsi.guti.mme_group_id;
                params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti->gummei.mme_code = msg->oldgutiorimsi.guti.mme_code;
                params->guti->m_tmsi = msg->oldgutiorimsi.guti.m_tmsi;
                            }
            else if (msg->m5gsmobileidentity.m5gs_mobile_identity_t.imsi.typeofidentity == M5GS_Mobile_Identity_IMSI) 
            {
                /*
                * Get the IMSI
                */
                
                params->m5gsmobileidentity.m5gs_mobile_identity_t.imsi = new(imsi_t); // need to define imsi_t
                params->imsi->u.num.digit1 = msg->oldgutiorimsi.imsi.identity_digit1;
                params->imsi->u.num.digit2 = msg->oldgutiorimsi.imsi.identity_digit2;
                params->imsi->u.num.digit3 = msg->oldgutiorimsi.imsi.identity_digit3;
                params->imsi->u.num.digit4 = msg->oldgutiorimsi.imsi.identity_digit4;
                params->imsi->u.num.digit5 = msg->oldgutiorimsi.imsi.identity_digit5;
                params->imsi->u.num.digit6 = msg->oldgutiorimsi.imsi.identity_digit6;
                params->imsi->u.num.digit7 = msg->oldgutiorimsi.imsi.identity_digit7;
                params->imsi->u.num.digit8 = msg->oldgutiorimsi.imsi.identity_digit8;
                params->imsi->u.num.digit9 = msg->oldgutiorimsi.imsi.identity_digit9;
                params->imsi->u.num.digit10 = msg->oldgutiorimsi.imsi.identity_digit10;
                params->imsi->u.num.digit11 = msg->oldgutiorimsi.imsi.identity_digit11;
                params->imsi->u.num.digit12 = msg->oldgutiorimsi.imsi.identity_digit12;
                params->imsi->u.num.digit13 = msg->oldgutiorimsi.imsi.identity_digit13;
                params->imsi->u.num.digit14 = msg->oldgutiorimsi.imsi.identity_digit14;
                params->imsi->u.num.digit15 = msg->oldgutiorimsi.imsi.identity_digit15;
                params->imsi->u.num.parity = 0x0f;
                params->imsi->length = msg->oldgutiorimsi.imsi.num_digits;

                // need to update from 23.003
            }
            else if (msg->m5gsmobileidentity.m5gs_mobile_identity_t.imei.typeofidentity == M5GS_Mobile_Identity_IMEI) 
            {
                //assign IMEI value
            }
            else if (msg->m5gsmobileidentity.m5gs_mobile_identity_t.m5gstmsi.typeofidentity == M5GS_Mobile_Identity_M5GS_TMSI) 
            {
                //assign m5gstmsi value
            }
            else if (msg->m5gsmobileidentity.m5gs_mobile_identity_t.imeisv.typeofidentity == M5GS_Mobile_Identity_IMEISV) 
            {
                //assign imeisv value  REGISTRATION_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT
            }       
            /*
            * Execute the requested UE registration procedure
            */
            rc = amf_proc_registration_request(ue_id, is_mm_ctx_new, params);
        }

    
}