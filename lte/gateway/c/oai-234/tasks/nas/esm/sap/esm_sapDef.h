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
Source      esm_sapDef.h

Version     0.1

Date        2012/11/21

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines the ESM Service Access Point that provides EPS
        bearer context handling and resources allocation procedures.

*****************************************************************************/

#ifndef __ESM_SAPDEF_H__
#define __ESM_SAPDEF_H__

#include "bstrlib.h"
#include "emm_data.h"
#include "3gpp_29.274.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * EPS Session Management primitives
 * ---------------------------------
 */
typedef enum esm_primitive_s {
  ESM_START = 0,
  /* Procedures related to EPS bearer contexts (initiated by the network) */
  ESM_DEFAULT_EPS_BEARER_CONTEXT_ACTIVATE_REQ,
  ESM_DEFAULT_EPS_BEARER_CONTEXT_ACTIVATE_CNF,
  ESM_DEFAULT_EPS_BEARER_CONTEXT_ACTIVATE_REJ,
  ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_REQ,
  ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_CNF,
  ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_REJ,
  ESM_EPS_BEARER_CONTEXT_MODIFY_REQ,
  ESM_EPS_BEARER_CONTEXT_MODIFY_CNF,
  ESM_EPS_BEARER_CONTEXT_MODIFY_REJ,
  ESM_EPS_BEARER_CONTEXT_DEACTIVATE_REQ,
  ESM_EPS_BEARER_CONTEXT_DEACTIVATE_CNF,
  /* Transaction related procedures (initiated by the UE) */
  ESM_PDN_CONNECTIVITY_REQ,
  ESM_PDN_CONNECTIVITY_REJ,
  ESM_PDN_DISCONNECT_REQ,
  ESM_PDN_DISCONNECT_REJ,
  ESM_BEARER_RESOURCE_ALLOCATE_REQ,
  ESM_BEARER_RESOURCE_ALLOCATE_REJ,
  ESM_BEARER_RESOURCE_MODIFY_REQ,
  ESM_BEARER_RESOURCE_MODIFY_REJ,
  /* ESM data indication ("raw" ESM message) */
  ESM_UNITDATA_IND,
  ESM_END
} esm_primitive_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Error code returned upon processing ESM-SAP primitive
 */
typedef enum {
  ESM_SAP_SUCCESS = 1, /* ESM-SAP primitive succeed           */
  ESM_SAP_DISCARDED,   /* ESM-SAP primitive failed, the caller should
                        * ignore the error                */
  ESM_SAP_FAILED       /* ESM-SAP primitive failed, the caller should
                        * take specific action and state transition may
                        * occurs                  */
} esm_sap_error_t;

/*
 * ESM primitive for activate EPS default bearer context procedure
 * ---------------------------------------------------------------
 */
typedef struct esm_activate_eps_default_bearer_context_s {
} esm_activate_eps_default_bearer_context_t;

/*
 * ESM primitive for PDN connectivity procedure
 * --------------------------------------------
 */
typedef struct esm_pdn_connectivity_s {
  pdn_cid_t cid;     /* PDN connection local identifier      */
  int is_defined;    /* Indicates whether a PDN context has been defined
                      * for the specified APN            */
  int pdn_type;      /* PDN address type (IPv4, IPv6, IPv4v6)    */
  const char* apn;   /* PDN's Access Point Name          */
  bool is_emergency; /* Indicates whether the PDN context has been
                      * defined to establish connection for emergency
                      * bearer services              */
} esm_pdn_connectivity_t;

/*
 * ESM primitive for PDN disconnect procedure
 * ------------------------------------------
 */
typedef struct esm_pdn_disconnect_s {
  pdn_cid_t cid; /* PDN connection local identifier      */
} esm_pdn_disconnect_t;

/*
 * ESM primitive for deactivate EPS bearer context procedure
 * ---------------------------------------------------------
 */
typedef struct esm_eps_bearer_context_deactivate_s {
#define ESM_SAP_ALL_EBI 0xff
  uint32_t no_of_bearers;
  ebi_t ebi[BEARERS_PER_UE]; /* EPS bearer identity of the EPS bearer context
                              * to be deactivated                */
  bool is_pcrf_initiated;
} esm_eps_bearer_context_deactivate_t;

/*
 * ESM primitive for activate dedicated EPS bearer context procedure
 * ---------------------------------------------------------
 */
typedef struct esm_eps_dedicated_bearer_context_activate_s {
  pdn_cid_t cid; /* PDN connection local identifier      */
  ebi_t ebi; /* EPS bearer identity of the EPS bearer context to be activated */
  ebi_t linked_ebi;
  qci_t qci;
  bitrate_t gbr_ul;
  bitrate_t gbr_dl;
  bitrate_t mbr_ul;
  bitrate_t mbr_dl;
  traffic_flow_template_t* tft;
  protocol_configuration_options_t* pco;
  fteid_t sgw_fteid;
} esm_eps_dedicated_bearer_context_activate_t;

/*
 * ------------------------------
 * Structure of ESM-SAP primitive
 * ------------------------------
 */
typedef union {
  esm_pdn_connectivity_t pdn_connect;
  esm_pdn_disconnect_t pdn_disconnect;
  esm_eps_bearer_context_deactivate_t eps_bearer_context_deactivate;
  esm_eps_dedicated_bearer_context_activate_t
      eps_dedicated_bearer_context_activate;
} esm_sap_data_t;

struct emm_context_s;

typedef struct esm_sap_s {
  esm_primitive_t primitive; /* ESM-SAP primitive to process     */
  bool is_standalone;        /* Indicates whether the ESM message handled
                              * within this primitive has to be sent/received
                              * standalone or together within an EMM related
                              * message              */
  struct emm_context_s* ctx; /* UE MM context                   */
  unsigned int ue_id;        /* Local UE identifier             */
  esm_sap_error_t err;       /* ESM-SAP error code               */
  const_bstring recv;        /* Encoded ESM message received     */
  bstring send;              /* Encoded ESM message to be sent   */
  esm_sap_data_t data;       /* ESM message data parameters      */
} esm_sap_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* __ESM_SAPDEF_H__*/
