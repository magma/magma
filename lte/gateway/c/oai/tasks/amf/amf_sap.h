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

  Source      amf_sap.h

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
#define MIN_GUMMEI 1
#define MAX_GUMMEI 5
typedef uint64_t imsi64_t;
typedef uint32_t amf_ue_ngap_id_t;
#include "amf_app_desc.h"
namespace magma5g
{
    class amf_sap_c
    {
        public:

        void amf_sap_initialize(void);

        int amf_sap_send(amf_sap_t* msg);


    };
    /****************************************************************************/
    /*********************  G L O B A L    C O N S T A N T S  *******************/
    /****************************************************************************/

    /*
    * 5GMM Mobility Management primitives
    * ----------------------------------
    * AMFREG-SAP provides registration services for location updating and
    * registration/deregistration procedures;
    * AMFSMF-SAP provides interlayer services to the 5GMM Session Management
    * NF for service registration and activate/deactivate PDU session context;
    * AMFAS-SAP provides services to the Access Stratum sublayer for NAS message
    * transfer;
    */
    enum amf_primitive_t
    {
      /* AMFREG-SAP */
      AMFREG_COMMON_PROC_REQ = _AMFREG_COMMON_PROC_REQ,
      AMFREG_COMMON_PROC_CNF = _AMFREG_COMMON_PROC_CNF,
      AMFREG_COMMON_PROC_REJ = _AMFREG_COMMON_PROC_REJ,
      AMFREG_COMMON_PROC_ABORT = _AMFREG_COMMON_PROC_ABORT,
      AMFREG_REGISTRATION_CNF = _AMFREG_REGISTRATION_CNF,
      AMFREG_REGISTRATION_REJ = _AMFREG_REGISTRATION_REJ,
      AMFREG_REGISTRATION_ABORT = _AMFREG_REGISTRATION_ABORT,
      AMFREG_DETACH_INIT = _AMFREG_DETACH_INIT,
      AMFREG_DETACH_REQ = _AMFREG_DETACH_REQ,
      AMFREG_DETACH_FAILED = _AMFREG_DETACH_FAILED,
      AMFREG_DETACH_CNF = _AMFREG_DETACH_CNF,
      AMFREG_TAU_REQ = _AMFREG_TAU_REQ,
      AMFREG_TAU_CNF = _AMFREG_TAU_CNF,
      AMFREG_TAU_REJ = _AMFREG_TAU_REJ,
      AMFREG_SERVICE_REQ = _AMFREG_SERVICE_REQ,
      AMFREG_SERVICE_CNF = _AMFREG_SERVICE_CNF,
      AMFREG_SERVICE_REJ = _AMFREG_SERVICE_REJ,
      AMFREG_LOWERLAYER_SUCCESS = _AMFREG_LOWERLAYER_SUCCESS,
      AMFREG_LOWERLAYER_FAILURE = _AMFREG_LOWERLAYER_FAILURE,
      AMFREG_LOWERLAYER_RELEASE = _AMFREG_LOWERLAYER_RELEASE,
      AMFREG_LOWERLAYER_NON_DELIVERY = _AMFREG_LOWERLAYER_NON_DELIVERY,
      /* AMFSMF-SAP */
      AMFSMF_RELEASE_IND = _AMFSMF_RELEASE_IND,
      AMFSMF_UNITDATA_REQ = _AMFSMF_UNITDATA_REQ,
      AMFSMF_ACTIVATE_BEARER_REQ = _AMFSMF_ACTIVATE_BEARER_REQ,
      AMFSMF_UNITDATA_IND = _AMFSMF_UNITDATA_IND,
      AMFSMF_DEACTIVATE_BEARER_REQ = _AMFSMF_DEACTIVATE_BEARER_REQ,
      /* AMFAS-SAP */
      AMFAS_SECURITY_REQ = _AMFAS_SECURITY_REQ,
      AMFAS_SECURITY_IND = _AMFAS_SECURITY_IND,
      AMFAS_SECURITY_RES = _AMFAS_SECURITY_RES,
      AMFAS_SECURITY_REJ = _AMFAS_SECURITY_REJ,
      AMFAS_ESTABLISH_REQ = _AMFAS_ESTABLISH_REQ,
      AMFAS_ESTABLISH_CNF = _AMFAS_ESTABLISH_CNF,
      AMFAS_ESTABLISH_REJ = _AMFAS_ESTABLISH_REJ,
      AMFAS_RELEASE_REQ = _AMFAS_RELEASE_REQ,
      AMFAS_RELEASE_IND = _AMFAS_RELEASE_IND,
      AMFAS_ERAB_SETUP_REQ = _AMFAS_ERAB_SETUP_REQ,
      AMFAS_ERAB_SETUP_CNF = _AMFAS_ERAB_SETUP_CNF,
      AMFAS_ERAB_SETUP_REJ = _AMFAS_ERAB_SETUP_REJ,
      AMFAS_DATA_REQ = _AMFAS_DATA_REQ,
      AMFAS_DATA_IND = _AMFAS_DATA_IND,
      AMFAS_PAGE_IND = _AMFAS_PAGE_IND,
      AMFAS_STATUS_IND = _AMFAS_STATUS_IND,
      AMFAS_ERAB_REL_CMD = _AMFAS_ERAB_REL_CMD,
      AMFAS_ERAB_REL_RSP = _AMFAS_ERAB_REL_RSP,

      AMFCN_AUTHENTICATION_PARAM_RES = _AMFCN_AUTHENTICATION_PARAM_RES,
      AMFCN_AUTHENTICATION_PARAM_FAIL = _AMFCN_AUTHENTICATION_PARAM_FAIL,
      AMFCN_ULA_SUCCESS = _AMFCN_ULA_SUCCESS,
      AMFCN_CS_RESPONSE_SUCCESS = _AMFCN_CS_RESPONSE_SUCCESS,
      AMFCN_ULA_OR_CSRSP_FAIL = _AMFCN_ULA_OR_CSRSP_FAIL,
      AMFCN_ACTIVATE_DEDICATED_BEARER_REQ = _AMFCN_ACTIVATE_DEDICATED_BEARER_REQ,
      AMFCN_IMPLICIT_DETACH_UE = _AMFCN_IMPLICIT_DETACH_UE,
      AMFCN_SMC_PROC_FAIL = _AMFCN_SMC_PROC_FAIL,
      AMFCN_NW_INITIATED_DETACH_UE = _AMFCN_NW_INITIATED_DETACH_UE,
      AMFCN_CS_DOMAIN_LOCATION_UPDT_ACC = _AMFCN_CS_DOMAIN_LOCATION_UPDT_ACC,
      AMFCN_CS_DOMAIN_LOCATION_UPDT_FAIL = _AMFCN_CS_DOMAIN_LOCATION_UPDT_FAIL,
      AMFCN_CS_DOMAIN_MM_INFORMATION_REQ = _AMFCN_CS_DOMAIN_MM_INFORMATION_REQ,
      AMFCN_DEACTIVATE_BEARER_REQ = _AMFCN_DEACTIVATE_BEARER_REQ,
      AMFCN_PDN_DISCONNECT_RES = _AMFCN_PDN_DISCONNECT_RES,
    };

    /*
    * Minimal identifier for AMF-SAP primitives
    */
    #define AMFREG_PRIMITIVE_MIN _AMFREG_START
    #define AMFSMF_PRIMITIVE_MIN _AMFSMF_START
    #define AMFAS_PRIMITIVE_MIN _AMFAS_START
    #define AMFCN_PRIMITIVE_MIN _AMFCN_START

    /*
    * Maximal identifier for AMF-SAP primitives
    */
    #define AMFREG_PRIMITIVE_MAX _AMFREG_END
    #define AMFSMF_PRIMITIVE_MAX _AMFSMF_END
    #define AMFAS_PRIMITIVE_MAX _AMFAS_END
    #define AMFCN_PRIMITIVE_MAX _AMFCN_END

    /****************************************************************************/
    /************************  G L O B A L    T Y P E S  ************************/
    /****************************************************************************/

    /*
    * Structure of 5GMM Mobility Management primitive
    */
    class amf_sap_t
    {
      public:
      amf_sap_t();
      amf_primitive_t primitive;
      union c
      {
        amf_reg_t amf_reg; /* AMFREG-SAP primitives    */
        amf_smf_t amf_smf; /* AMFSMF-SAP primitives    */
        amf_as_t amf_as;   /* AMFAS-SAP primitives     */
        amf_cn_t amf_cn;   /* AMFCN-SAP primitives     */
      };
      
    };
    

} 