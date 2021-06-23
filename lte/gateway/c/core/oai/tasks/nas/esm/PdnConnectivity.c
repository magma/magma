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

#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <sys/socket.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "log.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "3gpp_36.401.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "esm_data.h"
#include "esm_cause.h"
#include "esm_pt.h"
#include "mme_api.h"
#include "emm_sap.h"
#include "mme_app_pdn_context.h"
#include "3gpp_24.301.h"
#include "EsmCause.h"
#include "common_defs.h"
#include "emm_data.h"
#include "emm_esmDef.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
    Internal data handled by the PDN connectivity procedure in the MME
   --------------------------------------------------------------------------
*/
/*
   PDN connection handlers
*/
static pdn_cid_t pdn_connectivity_create(
    emm_context_t* emm_context, const proc_tid_t pti, const pdn_cid_t pdn_cid,
    const context_identifier_t context_identifier, const_bstring const apn,
    pdn_type_t pdn_type, const_bstring const pdn_addr,
    protocol_configuration_options_t* const pco, const bool is_emergency,
    esm_cause_t* esm_cause);

proc_tid_t pdn_connectivity_delete(
    emm_context_t* emm_context, pdn_cid_t pdn_cid);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
        PDN connectivity procedure executed by the MME
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_pdn_connectivity_request()                       **
 **                                                                        **
 ** Description: Performs PDN connectivity procedure requested by the UE.  **
 **                                                                        **
 **              3GPP TS 24.301, section 6.5.1.3                           **
 **      Upon receipt of the PDN CONNECTIVITY REQUEST message, the **
 **      MME checks if connectivity with the requested PDN can be  **
 **      established. If no requested  APN  is provided  the  MME  **
 **      shall use the default APN as the  requested  APN if the   **
 **      request type is different from "emergency", or the APN    **
 **      configured for emergency bearer services if the request   **
 **      type is "emergency".                                      **
 **      If connectivity with the requested PDN is accepted by the **
 **      network, the MME shall initiate the default EPS bearer    **
 **      context activation procedure.                             **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Identifies the PDN connectivity procedure  **
 **             requested by the UE                        **
 **      request_type:  Type of the PDN request                    **
 **      pdn_type:  PDN type value (IPv4, IPv6, IPv4v6)        **
 **      apn:       Requested Access Point Name                **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     apn:       Default Access Point Name                  **
 **      pdn_addr:  Assigned IPv4 address and/or IPv6 suffix   **
 **      esm_qos:   EPS bearer level QoS parameters            **
 **      esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    The identifier of the PDN connection if    **
 **             successfully created;                      **
 **             RETURNerror otherwise.                     **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_connectivity_request(
    emm_context_t* emm_context, const proc_tid_t pti, const pdn_cid_t pdn_cid,
    const context_identifier_t context_identifier,
    const esm_proc_pdn_request_t request_type, const_bstring const apn,
    pdn_type_t pdn_type, const_bstring const pdn_addr,
    bearer_qos_t* default_qos, protocol_configuration_options_t* const pco,
    esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNok;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN connectivity requested by the UE "
      "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);
  OAILOG_DEBUG(
      LOG_NAS_ESM,
      "PDN connectivity request :pti= %d PDN type = %s, APN = %s pdn addr = %s "
      "pdn id %d for (ue_id = %u)\n",
      pti,
      (pdn_type == ESM_PDN_TYPE_IPV4) ?
          "IPv4" :
          (pdn_type == ESM_PDN_TYPE_IPV6) ? "IPv6" : "IPv4v6",
      (apn) ? (char*) bdata(apn) : "null",
      (pdn_addr) ? (char*) bdata(pdn_addr) : "null", pdn_cid, ue_id);

  /*
   * Check network IP capabilities
   */
  OAILOG_DEBUG(
      LOG_NAS_ESM,
      "ESM-PROC  - _esm_data.conf.features %08x for (ue_id = %u)\n",
      _esm_data.conf.features, ue_id);

  int is_emergency = (request_type == ESM_PDN_REQUEST_EMERGENCY);
  OAILOG_DEBUG(
      LOG_NAS_ESM, "ESM-PROC  - is_emergency = (%d) for (ue_id = %u)\n",
      is_emergency, ue_id);

  /*
   * Create new PDN connection
   */
  rc = pdn_connectivity_create(
      emm_context, pti, pdn_cid, context_identifier, apn, pdn_type, pdn_addr,
      pco, is_emergency, esm_cause);

  if (rc < 0) {
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - Failed to create PDN connection for "
        "ue_id " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    *esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_pdn_connectivity_reject()                        **
 **                                                                        **
 ** Description: Performs PDN connectivity procedure not accepted by the   **
 **      network.                                                  **
 **                                                                        **
 **              3GPP TS 24.301, section 6.5.1.4                           **
 **      If connectivity with the requested PDN cannot be accepted **
 **      by the network, the MME shall send a PDN CONNECTIVITY RE- **
 **      JECT message to the UE.                                   **
 **                                                                        **
 ** Inputs:  is_standalone: Indicates whether the PDN connectivity     **
 **             procedure was initiated as part of the at- **
 **             tach procedure                             **
 **      ue_id:      UE lower layer identifier                  **
 **      ebi:       Not used                                   **
 **      msg:       Encoded PDN connectivity reject message to **
 **             be sent                                    **
 **      ue_triggered:  Not used                                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_connectivity_reject(
    bool is_standalone, emm_context_t* emm_context, ebi_t ebi,
    STOLEN_REF bstring* msg, bool ue_triggered) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_WARNING(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN connectivity not accepted by the "
      "network (ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);

  if (is_standalone) {
    emm_sap_t emm_sap = {0};

    /*
     * Notity EMM that ESM PDU has to be forwarded to lower layers
     */
    emm_sap.primitive            = EMMESM_UNITDATA_REQ;
    emm_sap.u.emm_esm.ue_id      = ue_id;
    emm_sap.u.emm_esm.ctx        = emm_context;
    emm_sap.u.emm_esm.u.data.msg = *msg;
    rc                           = emm_sap_send(&emm_sap);
  }

  /*
   * If the PDN connectivity procedure initiated as part of the initial
   * * * * attach procedure has failed, an error is returned to notify EMM that
   * * * * the ESM sublayer did not accept UE requested PDN connectivity
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:        esm_proc_pdn_connectivity_failure()                       **
 **                                                                        **
 ** Description: Performs PDN connectivity procedure upon receiving noti-  **
 **              fication from the EPS Mobility Management sublayer that   **
 **              EMM procedure that initiated PDN connectivity activation  **
 **              locally failed.                                           **
 **                                                                        **
 **              The MME releases the PDN connection entry allocated when  **
 **              the PDN connectivity procedure was requested by the UE.   **
 **                                                                        **
 **         Inputs:  ue_id:      UE local identifier                        **
 **                  pdn_cid:       Identifier of the PDN connection to be **
 **                             released                                   **
 **                  Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **                  Return:    RETURNok, RETURNerror                      **
 **                  Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_connectivity_failure(
    emm_context_t* emm_context, pdn_cid_t pdn_cid) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  proc_tid_t pti = ESM_PT_UNASSIGNED;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_WARNING(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN connectivity failure (ue_id=" MME_UE_S1AP_ID_FMT
      ", pdn_cid=%d)\n",
      ue_id, pdn_cid);
  /*
   * Delete the PDN connection entry
   */
  pti = pdn_connectivity_delete(emm_context, pdn_cid);

  if (pti != ESM_PT_UNASSIGNED) {
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/*
  ---------------------------------------------------------------------------
                PDN connection handlers
  ---------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:        _pdn_connectivity_create()                                **
 **                                                                        **
 ** Description: Creates a new PDN connection entry for the specified UE   **
 **                                                                        **
 ** Inputs:          ue_id:      UE local identifier                        **
 **                  ctx:       UE context                                 **
 **                  pti:       Procedure transaction identity             **
 **                  apn:       Access Point Name of the PDN connection    **
 **                  pdn_type:  PDN type (IPv4, IPv6, IPv4v6)              **
 **                  pdn_addr:  Network allocated PDN IPv4 or IPv6 address **
 **              is_emergency:  true if the PDN connection has to be esta- **
 **                             blished for emergency bearer services      **
 **                  Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     None                                                      **
 **                  Return:    The identifier of the PDN connection if    **
 **                             successfully created; -1 otherwise.        **
 **                  Others:    _esm_data                                  **
 **                                                                        **
 ***************************************************************************/
static int pdn_connectivity_create(
    emm_context_t* emm_context, const proc_tid_t pti, const pdn_cid_t pdn_cid,
    const context_identifier_t context_identifier, const_bstring const apn,
    pdn_type_t pdn_type, const_bstring const pdn_addr,
    protocol_configuration_options_t* const pco, const bool is_emergency,
    esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);

  OAILOG_DEBUG(
      LOG_NAS_ESM,
      "ESM-PROC  - Create new PDN connection (pti=%d), APN = %s, pdn_type = "
      "%d, IP address = %s "
      "PDN id %d (ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      pti, bdata(apn), pdn_type,
      (pdn_type == ESM_PDN_TYPE_IPV4) ?
          esm_data_get_ipv4_addr(pdn_addr) :
          (pdn_type == ESM_PDN_TYPE_IPV6) ? esm_data_get_ipv6_addr(pdn_addr) :
                                            esm_data_get_ipv4v6_addr(pdn_addr),
      pdn_cid, ue_mm_context->mme_ue_s1ap_id);

  if (!ue_mm_context->pdn_contexts[pdn_cid]) {
    /*
     * Create new PDN connection
     */
    pdn_context_t* pdn_context =
        mme_app_create_pdn_context(ue_mm_context, pdn_cid, context_identifier);

    if (pdn_context) {
      /*
       * Increment the number of PDN connections
       */
      ue_mm_context->nb_active_pdn_contexts += 1;
      /*
       * Set the procedure transaction identity
       */
      pdn_context->esm_data.pti = pti;
      /*
       * Set the emergency bearer services indicator
       */
      pdn_context->esm_data.is_emergency = is_emergency;

      // Set the esm cause
      pdn_context->esm_data.esm_cause = *esm_cause;

      if (pco) {
        if (!pdn_context->pco) {
          pdn_context->pco =
              calloc(1, sizeof(protocol_configuration_options_t));
        } else {
          clear_protocol_configuration_options(pdn_context->pco);
        }
        copy_protocol_configuration_options(pdn_context->pco, pco);
      } else {
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - PCO is NULL for ue_id " MME_UE_S1AP_ID_FMT "\n",
            ue_mm_context->mme_ue_s1ap_id);
      }

      /*
       * Setup the IP address allocated by the network
       */
      OAILOG_DEBUG(
          LOG_NAS_ESM, "PDN TYPE = %d for (ue_id = %u)\n", pdn_type,
          ue_mm_context->mme_ue_s1ap_id);
      pdn_context->pdn_type = pdn_type;
      if (pdn_addr) {
        pdn_context->paa.pdn_type = pdn_type;
        switch (pdn_type) {
          case IPv4:
            IPV4_STR_ADDR_TO_INADDR(
                (const char*) pdn_addr->data, pdn_context->paa.ipv4_address,
                "BAD IPv4 ADDRESS FORMAT FOR PAA!\n");
            break;
          case IPv6:
            AssertFatal(
                1 == inet_pton(
                         AF_INET6, (const char*) pdn_addr->data,
                         &pdn_context->paa.ipv6_address),
                "BAD IPv6 ADDRESS FORMAT FOR PAA!\n");
            break;
          // TODO Handle static IPv4v6 addr allocation
          case IPv4_AND_v6:
            Fatal("TODO Implement pdn_connectivity_create IPv4_AND_v6 \n");
            break;
          case IPv4_OR_v6:
            Fatal("TODO Implement pdn_connectivity_create IPv4_OR_v6 \n");
            break;
          default:;
        }
      }
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
    }

    OAILOG_ERROR(
        LOG_NAS_ESM,
        "ESM-PROC  - Failed to create new PDN connection (pdn_cid=%d) for "
        "(ue_id = " MME_UE_S1AP_ID_FMT ")\n",
        pdn_cid, ue_mm_context->mme_ue_s1ap_id);
  } else {
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - PDN connection already exist (pdn_cid=%d) for (ue_id "
        "= " MME_UE_S1AP_ID_FMT ")\n",
        pdn_cid, ue_mm_context->mme_ue_s1ap_id);
    // already created
    pdn_context_t* pdn_context = ue_mm_context->pdn_contexts[pdn_cid];

    if (pdn_context) {
      // QUICK WORKAROUND, TODO seriously

      /*
       * Set the procedure transaction identity
       */
      pdn_context->esm_data.pti          = pti;
      pdn_context->esm_data.is_emergency = is_emergency;
      if (pco) {
        if (!pdn_context->pco) {
          pdn_context->pco =
              calloc(1, sizeof(protocol_configuration_options_t));
        } else {
          clear_protocol_configuration_options(pdn_context->pco);
        }
        copy_protocol_configuration_options(pdn_context->pco, pco);
      } else {
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - PCO is NULL for ue_id " MME_UE_S1AP_ID_FMT "\n",
            ue_mm_context->mme_ue_s1ap_id);
      }
      pdn_context->pdn_type = pdn_type;
      if (pdn_addr) {
        pdn_context->paa.pdn_type = pdn_type;
        switch (pdn_type) {
          case IPv4:
            IPV4_STR_ADDR_TO_INADDR(
                (const char*) pdn_addr->data, pdn_context->paa.ipv4_address,
                "BAD IPv4 ADDRESS FORMAT FOR PAA!\n");
            break;
          case IPv6:
            AssertFatal(
                1 == inet_pton(
                         AF_INET6, (const char*) pdn_addr->data,
                         &pdn_context->paa.ipv6_address),
                "BAD IPv6 ADDRESS FORMAT FOR PAA!\n");
            break;
          case IPv4_AND_v6:
            Fatal("TODO Implement pdn_connectivity_create IPv4_AND_v6 \n");
            break;
          case IPv4_OR_v6:
            Fatal("TODO Implement pdn_connectivity_create IPv4_OR_v6 \n");
            break;
          default:;
        }
      }
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:        _pdn_connectivity_delete()                                **
 **                                                                        **
 ** Description: Deletes PDN connection to the specified UE associated to  **
 **              PDN connection entry with given identifier                **
 **                                                                        **
 ** Inputs:          ue_id:     UE local identifier                        **
 **                  pdn_cid:   Identifier of the PDN connection to be     **
 **                             released                                   **
 **                  Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     None                                                      **
 **                  Return:    The identity of the procedure transaction  **
 **                             assigned to the PDN connection when suc-   **
 **                             cessfully released;                        **
 **                             UNASSIGNED value otherwise.                **
 **                  Others:    _esm_data                                  **
 **                                                                        **
 ***************************************************************************/
proc_tid_t pdn_connectivity_delete(
    emm_context_t* emm_context, pdn_cid_t pdn_cid) {
  proc_tid_t pti = ESM_PT_UNASSIGNED;

  if (!emm_context) {
    return pti;
  }
  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);

  if (pdn_cid < MAX_APN_PER_UE) {
    if (!ue_mm_context->pdn_contexts[pdn_cid]) {
      OAILOG_ERROR(
          LOG_NAS_ESM, "ESM-PROC  - PDN connection has not been allocated\n");
    } else if (ue_mm_context->pdn_contexts[pdn_cid]->is_active) {
      OAILOG_ERROR(LOG_NAS_ESM, "ESM-PROC  - PDN connection is active\n");
    } else {
      /*
       * Get the identity of the procedure transaction that created
       * the PDN connection
       */
      pti = ue_mm_context->pdn_contexts[pdn_cid]->esm_data.pti;
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_ESM, "ESM-PROC  - PDN connection identifier is not valid\n");
  }
  if (pti != ESM_PT_UNASSIGNED) {
    // Release allocated PDN connection data
    bdestroy_wrapper(&ue_mm_context->pdn_contexts[pdn_cid]->apn_in_use);
    bdestroy_wrapper(&ue_mm_context->pdn_contexts[pdn_cid]->apn_oi_replacement);
    bdestroy_wrapper(&ue_mm_context->pdn_contexts[pdn_cid]->apn_subscribed);
    memset(
        &ue_mm_context->pdn_contexts[pdn_cid]->esm_data, 0,
        sizeof(ue_mm_context->pdn_contexts[pdn_cid]->esm_data));

    // Free protocol configuration options and its contents
    if (ue_mm_context->pdn_contexts[pdn_cid]->pco) {
      free_protocol_configuration_options(
          &ue_mm_context->pdn_contexts[pdn_cid]->pco);
    }
    // TODO Think about free ue_mm_context->pdn_contexts[pdn_cid]
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - PDN connection %d released for ue id " MME_UE_S1AP_ID_FMT
        "\n",
        pdn_cid, ue_mm_context->mme_ue_s1ap_id);
  }
  // Return the procedure transaction identity
  return (pti);
}
