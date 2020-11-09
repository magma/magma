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

  Source      amf_sap.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_fsm.h"
#include "amf_asDefs.h"
#include "amf_as.h"
#include "amf_sap.h"

using namespace std;

namespace magma5g
{
        amf_as amf_access; //Object to call respective functions
        int amf_sap_c::amf_sap_send(amf_sap_t *msg)
        {
            int rc = RETURNerror;
            amf_primitive_t primitive = msg->primitive;

            OAILOG_FUNC_IN(LOG_NAS_AMF);

            /*
            * Check the AMF-SAP primitive
            */
            if ((primitive > (amf_primitive_t) AMFREG_PRIMITIVE_MIN) && (primitive < (amf_primitive_t) AMFREG_PRIMITIVE_MAX)) 
            {
                    /*
                    * Forward to the AMFREG-SAP
                    * will handle for state update
                    */
                    //msg->u.amf_reg.primitive = primitive;
                    //TODO rc = amf_reg_send(&msg->u.amf_reg);//TODO -  NEED-RECHECK 
            } 
            //else if ((primitive > (amf_primitive_t) AMFSMF_PRIMITIVE_MIN) && (primitive < (amf_primitive_t) AMFSMF_PRIMITIVE_MAX)) 
            //{
                /*
                * Forward to the AMFESM-SAP
                */
                //msg->u.amf_smf.primitive = primitive;
                //rc = amf_smf_send(&msg->u.amf_smf);
            //} 
            else if ((primitive > (amf_primitive_t) AMFAS_PRIMITIVE_MIN) && (primitive < (amf_primitive_t) AMFAS_PRIMITIVE_MAX)) 
            {
                /*
                * Forward to the AMFAS-SAP
                */
                //msg->u.amf_as.primitive = primitive; //TODO -  NEED-RECHECK
                rc = amf_access.amf_as_send(&msg->u.amf_as);
            } 
            else if ((primitive > (amf_primitive_t) AMFCN_PRIMITIVE_MIN) && (primitive < (amf_primitive_t) AMFCN_PRIMITIVE_MAX)) 
            {
                /*
                * Forward to the AMFCN-SAP
                */
               // msg->u.amf_cn.primitive = primitive;
               // rc = amf_access.amf_cn_send(&msg->u.amf_cn); //TODO -  NEED-RECHECK
            } 
            else 
            {
                //OAILOG_WARNING( LOG_NAS_AMF, "AMF-SAP -   Out of range primitive (%d)\n", primitive);
            }

            OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
        }

        void amf_as_data_t::amf_as_set_security_data(amf_as_security_data_t* data, 
                              void* context, bool is_new,
                              bool is_ciphered) 
        {
           OAILOG_FUNC_IN(LOG_NAS_AMF);
        }
}
