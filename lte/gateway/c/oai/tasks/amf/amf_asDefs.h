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

  Source      amf_asDefs.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
using namespace std;
#pragma once
namespace magmam5g
{
    ////////////////////typecast require///////////////////////////////
    /****************************************************************************/
    /*********************  G L O B A L    C O N S T A N T S  *******************/
    /****************************************************************************/

    /*
    * AMFAS-SAP primitives
    */
     enum amf_as_primitive_t {
    _AMFAS_START = 200,
    _AMFAS_SECURITY_REQ,   /* AMF->AS: Security request          */
    _AMFAS_SECURITY_IND,   /* AS->AMF: Security indication       */
    _AMFAS_SECURITY_RES,   /* AMF->AS: Security response         */
    _AMFAS_SECURITY_REJ,   /* AMF->AS: Security reject           */
    _AMFAS_ESTABLISH_REQ,  /* AMF->AS: Connection establish request  */
    _AMFAS_ESTABLISH_CNF,  /* AS->AMF: Connection establish confirm  */
    _AMFAS_ESTABLISH_REJ,  /* AS->AMF: Connection establish reject   */
    _AMFAS_RELEASE_REQ,    /* AMF->AS: Connection release request    */
    _AMFAS_RELEASE_IND,    /* AS->AMF: Connection release indication */
    _AMFAS_ERAB_SETUP_REQ, /* AMF->AS: ERAB setup request  */
    _AMFAS_ERAB_SETUP_CNF, /* AS->AMF  */
    _AMFAS_ERAB_SETUP_REJ, /* AS->AMF  */
    _AMFAS_DATA_REQ,       /* AMF->AS: Data transfer request     */
    _AMFAS_DATA_IND,       /* AS->AMF: Data transfer indication      */
    _AMFAS_PAGE_IND,       /* AS->AMF: Paging data indication        */
    _AMFAS_STATUS_IND,     /* AS->AMF: Status indication         */
    _AMFAS_ERAB_REL_CMD,   /* AMF->AS: ERAB Release Cmd  */
    _AMFAS_ERAB_REL_RSP,   /* AMF->AS: ERAB Release Rsp  */
    _AMFAS_END
    };
    /*
    * AMFS primitive for data transfer
    * ---------------------------------
    */
    class amf_as_data_t 
    {
      public:
      amf_ue_ngap_id_t ue_id;       /* UE lower layer identifier        */
      amf_as_M5GS_identity_t m5gs_id; /* UE's M5G mobile identity         */
      const guti_m5_t* guti;           /* GUTI temporary mobile identity   */
      //const guti_t* new_guti;       /* New GUTI, if re-allocated        */
     // amf_as_security_data_t sctx;  /* M5G NAS security context         */
      //uint8_t encryption : 4;       /* Ciphering algorithm              */
      //uint8_t integrity : 4;        /* Integrity protection algorithm   */
      const plmn_t* plmn_id;        /* Identifier of the selected PLMN  */
      ecgi_t ecgi;      /* NR CGI This information element is used to globally
                          identify a cell */
      //const tai_t* tai; /* Code of the first tracking area identity the UE is registered to          */
      //tai_list_t tai_list;                  /* Valid field if num tai > 0 */
      
      bool switch_off;                      /* true if the UE is switched off   */
      uint8_t type;                         /* Network deregister type          */
    #define AMF_AS_DATA_DELIVERED_LOWER_LAYER_FAILURE 0
    #define AMF_AS_DATA_DELIVERED_TRUE 1
    #define AMF_AS_DATA_DELIVERED_LOWER_LAYER_NON_DELIVERY_INDICATION_DUE_TO_HO 2
    uint8_t delivered;                    /* Data message delivery indicator  */
    #define AMF_AS_NAS_DATA_REGISTRATION 0x01     /* REGISTRATION Complete          */
    #define AMF_AS_NAS_DATA_DEREGISTRATION_REQ 0x02 /* DEREGISTRATION Request           */
    #define AMF_AS_NAS_DATA_TAU 0x03        /* TAU    REGISTRATION            */
    #define AMF_AS_NAS_DATA_REGISTRATION_ACCEPT 0x04 /* REGISTRATION Accept            */
    #define AMF_AS_NAS_AMF_INFORMATION 0x05    /* Emm information          */
    #define AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT 0x06 /* DEREGISTRATION Accept            */
    #define AMF_AS_NAS_DATA_CS_SERVICE_NOTIFICATION  0x07      /* CS Service Notification  */
    #define AMF_AS_NAS_DATA_INFO_SR 0x08     /* Service Reject in DL NAS */
    #define AMF_AS_NAS_DL_NAS_TRANSPORT 0x09 /* Downlink Nas Transport */
      uint8_t nas_info; /* Type of NAS information to transfer  */
      std::string nas_msg;  /* NAS message to be transferred     */
      std::string full_network_name;
      std::string short_network_name;
      std::uint8_t daylight_saving_time;
      //LAI_t* location_area_identification; /* Location area identification */
      mobile_identity_t* ms_identity; /* MS identity This IE may be included to assign or unassign
                          a new TMSI to a UE during a combined TA/LA update. */
      std::uint8_t* additional_update_result;   /* TAU Additional update result   */
      std::uint32_t* amf_cause;                 /* EMM failure cause code        */
      std::string cli;             /* Calling Line Identification  */
    };
}