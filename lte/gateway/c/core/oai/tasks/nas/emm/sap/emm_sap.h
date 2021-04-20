/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************

Source      emm_sap.h

Version     0.1

Date        2012/10/01

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines the EMM Service Access Points at which the EPS
        Mobility Management sublayer provides procedures for the
        control of security and mobility when the User Equipment
        is using the Evolved UTRA Network.

*****************************************************************************/
#ifndef FILE_EMM_SAP_SEEN
#define FILE_EMM_SAP_SEEN

#include "emm_data.h"
#include "emm_regDef.h"
#include "emm_esmDef.h"
#include "emm_asDef.h"
#include "emm_cnDef.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * EPS Mobility Management primitives
 * ----------------------------------
 * EMMREG-SAP provides registration services for location updating and
 * attach/detach procedures;
 * EMMESM-SAP provides interlayer services to the EPS Session Management
 * sublayer for service registration and activate/deactivate PDP context;
 * EMMAS-SAP provides services to the Access Stratum sublayer for NAS message
 * transfer;
 */
typedef enum {
  /* EMMREG-SAP */
  EMMREG_COMMON_PROC_REQ         = _EMMREG_COMMON_PROC_REQ,
  EMMREG_COMMON_PROC_CNF         = _EMMREG_COMMON_PROC_CNF,
  EMMREG_COMMON_PROC_REJ         = _EMMREG_COMMON_PROC_REJ,
  EMMREG_COMMON_PROC_ABORT       = _EMMREG_COMMON_PROC_ABORT,
  EMMREG_ATTACH_CNF              = _EMMREG_ATTACH_CNF,
  EMMREG_ATTACH_REJ              = _EMMREG_ATTACH_REJ,
  EMMREG_ATTACH_ABORT            = _EMMREG_ATTACH_ABORT,
  EMMREG_DETACH_INIT             = _EMMREG_DETACH_INIT,
  EMMREG_DETACH_REQ              = _EMMREG_DETACH_REQ,
  EMMREG_DETACH_FAILED           = _EMMREG_DETACH_FAILED,
  EMMREG_DETACH_CNF              = _EMMREG_DETACH_CNF,
  EMMREG_TAU_REQ                 = _EMMREG_TAU_REQ,
  EMMREG_TAU_CNF                 = _EMMREG_TAU_CNF,
  EMMREG_TAU_REJ                 = _EMMREG_TAU_REJ,
  EMMREG_SERVICE_REQ             = _EMMREG_SERVICE_REQ,
  EMMREG_SERVICE_CNF             = _EMMREG_SERVICE_CNF,
  EMMREG_SERVICE_REJ             = _EMMREG_SERVICE_REJ,
  EMMREG_LOWERLAYER_SUCCESS      = _EMMREG_LOWERLAYER_SUCCESS,
  EMMREG_LOWERLAYER_FAILURE      = _EMMREG_LOWERLAYER_FAILURE,
  EMMREG_LOWERLAYER_RELEASE      = _EMMREG_LOWERLAYER_RELEASE,
  EMMREG_LOWERLAYER_NON_DELIVERY = _EMMREG_LOWERLAYER_NON_DELIVERY,
  /* EMMESM-SAP */
  EMMESM_RELEASE_IND           = _EMMESM_RELEASE_IND,
  EMMESM_UNITDATA_REQ          = _EMMESM_UNITDATA_REQ,
  EMMESM_ACTIVATE_BEARER_REQ   = _EMMESM_ACTIVATE_BEARER_REQ,
  EMMESM_UNITDATA_IND          = _EMMESM_UNITDATA_IND,
  EMMESM_DEACTIVATE_BEARER_REQ = _EMMESM_DEACTIVATE_BEARER_REQ,
  /* EMMAS-SAP */
  EMMAS_SECURITY_REQ   = _EMMAS_SECURITY_REQ,
  EMMAS_SECURITY_IND   = _EMMAS_SECURITY_IND,
  EMMAS_SECURITY_RES   = _EMMAS_SECURITY_RES,
  EMMAS_SECURITY_REJ   = _EMMAS_SECURITY_REJ,
  EMMAS_ESTABLISH_REQ  = _EMMAS_ESTABLISH_REQ,
  EMMAS_ESTABLISH_CNF  = _EMMAS_ESTABLISH_CNF,
  EMMAS_ESTABLISH_REJ  = _EMMAS_ESTABLISH_REJ,
  EMMAS_RELEASE_REQ    = _EMMAS_RELEASE_REQ,
  EMMAS_RELEASE_IND    = _EMMAS_RELEASE_IND,
  EMMAS_ERAB_SETUP_REQ = _EMMAS_ERAB_SETUP_REQ,
  EMMAS_ERAB_SETUP_CNF = _EMMAS_ERAB_SETUP_CNF,
  EMMAS_ERAB_SETUP_REJ = _EMMAS_ERAB_SETUP_REJ,
  EMMAS_DATA_REQ       = _EMMAS_DATA_REQ,
  EMMAS_DATA_IND       = _EMMAS_DATA_IND,
  EMMAS_PAGE_IND       = _EMMAS_PAGE_IND,
  EMMAS_STATUS_IND     = _EMMAS_STATUS_IND,
  EMMAS_ERAB_REL_CMD   = _EMMAS_ERAB_REL_CMD,
  EMMAS_ERAB_REL_RSP   = _EMMAS_ERAB_REL_RSP,

  EMMCN_AUTHENTICATION_PARAM_RES      = _EMMCN_AUTHENTICATION_PARAM_RES,
  EMMCN_AUTHENTICATION_PARAM_FAIL     = _EMMCN_AUTHENTICATION_PARAM_FAIL,
  EMMCN_ULA_SUCCESS                   = _EMMCN_ULA_SUCCESS,
  EMMCN_CS_RESPONSE_SUCCESS           = _EMMCN_CS_RESPONSE_SUCCESS,
  EMMCN_ULA_OR_CSRSP_FAIL             = _EMMCN_ULA_OR_CSRSP_FAIL,
  EMMCN_ACTIVATE_DEDICATED_BEARER_REQ = _EMMCN_ACTIVATE_DEDICATED_BEARER_REQ,
  EMMCN_IMPLICIT_DETACH_UE            = _EMMCN_IMPLICIT_DETACH_UE,
  EMMCN_SMC_PROC_FAIL                 = _EMMCN_SMC_PROC_FAIL,
  EMMCN_NW_INITIATED_DETACH_UE        = _EMMCN_NW_INITIATED_DETACH_UE,
  EMMCN_CS_DOMAIN_LOCATION_UPDT_ACC   = _EMMCN_CS_DOMAIN_LOCATION_UPDT_ACC,
  EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL  = _EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL,
  EMMCN_CS_DOMAIN_MM_INFORMATION_REQ  = _EMMCN_CS_DOMAIN_MM_INFORMATION_REQ,
  EMMCN_DEACTIVATE_BEARER_REQ         = _EMMCN_DEACTIVATE_BEARER_REQ,
  EMMCN_PDN_DISCONNECT_RES            = _EMMCN_PDN_DISCONNECT_RES,
} emm_primitive_t;

/*
 * Minimal identifier for EMM-SAP primitives
 */
#define EMMREG_PRIMITIVE_MIN _EMMREG_START
#define EMMESM_PRIMITIVE_MIN _EMMESM_START
#define EMMAS_PRIMITIVE_MIN _EMMAS_START
#define EMMCN_PRIMITIVE_MIN _EMMCN_START

/*
 * Maximal identifier for EMM-SAP primitives
 */
#define EMMREG_PRIMITIVE_MAX _EMMREG_END
#define EMMESM_PRIMITIVE_MAX _EMMESM_END
#define EMMAS_PRIMITIVE_MAX _EMMAS_END
#define EMMCN_PRIMITIVE_MAX _EMMCN_END

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Structure of EPS Mobility Management primitive
 */
typedef struct emm_sap_s {
  emm_primitive_t primitive;
  union {
    emm_reg_t emm_reg; /* EMMREG-SAP primitives    */
    emm_esm_t emm_esm; /* EMMESM-SAP primitives    */
    emm_as_t emm_as;   /* EMMAS-SAP primitives     */
    emm_cn_t emm_cn;   /* EMMCN-SAP primitives     */
  } u;
} emm_sap_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

void emm_sap_initialize(void);

int emm_sap_send(emm_sap_t* msg);

#endif /* FILE_EMM_SAP_SEEN*/
