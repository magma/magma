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

#include <string.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "log.h"
#include "3gpp_24.007.h"
#include "3gpp_24.301.h"
#include "esm_send.h"
#include "esm_cause.h"
#include "ApnAggregateMaximumBitRate.h"
#include "PdnAddress.h"
#include "common_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Functions executed by the MME to send ESM messages
   --------------------------------------------------------------------------
*/
int esm_send_esm_information_request(
    pti_t pti, ebi_t ebi, esm_information_request_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity            = ebi;
  msg->messagetype                  = ESM_INFORMATION_REQUEST;
  msg->proceduretransactionidentity = pti;
  OAILOG_NOTICE(
      LOG_NAS_ESM,
      "ESM-SAP   - Send ESM_INFORMATION_REQUEST message (pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}
/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_status()                                         **
 **                                                                        **
 ** Description: Builds ESM status message                                 **
 **                                                                        **
 **      The ESM status message is sent by the network or the UE   **
 **      to pass information on the status of the indicated EPS    **
 **      bearer context and report certain error conditions.       **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      esm_cause: ESM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_status(pti_t pti, ebi_t ebi, esm_status_msg* msg, int esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity            = ebi;
  msg->messagetype                  = ESM_STATUS;
  msg->proceduretransactionidentity = pti;
  /*
   * Mandatory - ESM cause code
   */
  msg->esmcause = esm_cause;
  OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Send ESM Status message (pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/*
   --------------------------------------------------------------------------
   Functions executed by the MME to send ESM message to the UE
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_pdn_connectivity_reject()                        **
 **                                                                        **
 ** Description: Builds PDN Connectivity Reject message                    **
 **                                                                        **
 **      The PDN connectivity reject message is sent by the net-   **
 **      work to the UE to reject establishment of a PDN connec-   **
 **      tion.                                                     **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      esm_cause: ESM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_pdn_connectivity_reject(
    pti_t pti, pdn_connectivity_reject_msg* msg, int esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity            = EPS_BEARER_IDENTITY_UNASSIGNED;
  msg->messagetype                  = PDN_CONNECTIVITY_REJECT;
  msg->proceduretransactionidentity = pti;
  /*
   * Mandatory - ESM cause code
   */
  msg->esmcause = esm_cause;
  /*
   * Optional IEs
   */
  msg->presencemask = 0;
  OAILOG_DEBUG(
      LOG_NAS_ESM,
      "ESM-SAP   - Send PDN Connectivity Reject message "
      "(pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_pdn_disconnect_reject()                          **
 **                                                                        **
 ** Description: Builds PDN Disconnect Reject message                      **
 **                                                                        **
 **      The PDN disconnect reject message is sent by the network  **
 **      to the UE to reject release of a PDN connection.          **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      esm_cause: ESM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_pdn_disconnect_reject(
    pti_t pti, pdn_disconnect_reject_msg* msg, int esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity            = EPS_BEARER_IDENTITY_UNASSIGNED;
  msg->messagetype                  = PDN_DISCONNECT_REJECT;
  msg->proceduretransactionidentity = pti;
  /*
   * Mandatory - ESM cause code
   */
  msg->esmcause = esm_cause;
  /*
   * Optional IEs
   */
  msg->presencemask = 0;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-SAP   - Send PDN Disconnect Reject message "
      "(pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_activate_default_eps_bearer_context_request()    **
 **                                                                        **
 ** Description: Builds Activate Default EPS Bearer Context Request        **
 **      message                                                   **
 **                                                                        **
 **      The activate default EPS bearer context request message   **
 **      is sent by the network to the UE to request activation of **
 **      a default EPS bearer context.                             **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      qos:       Subscribed EPS quality of service          **
 **      apn:       Access Point Name in used                  **
 **      pdn_addr:  PDN IPv4 address and/or IPv6 suffix        **
 **      esm_cause: ESM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_activate_default_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, activate_default_eps_bearer_context_request_msg* msg,
    pdn_context_t* pdn_context_p, const protocol_configuration_options_t* pco,
    int pdn_type, bstring pdn_addr, const EpsQualityOfService* qos,
    int esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-SAP   - Send Activate Default EPS Bearer Context Request message\n");
  bstring apn = pdn_context_p->apn_subscribed;
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity     = ebi;
  msg->messagetype           = ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST;
  msg->proceduretransactionidentity = pti;
  OAILOG_DEBUG(
      LOG_NAS_ESM, "Activate default EPS bearer context (pti=%d, ebi=%d) \n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  /*
   * Mandatory - EPS QoS
   */
  msg->epsqos = *qos;
  OAILOG_DEBUG(LOG_NAS_ESM, "ESM-SAP   - epsqos  qci:  %u\n", qos->qci);

  if (qos->bitRatesPresent) {
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  maxBitRateForUL:  %u\n",
        qos->bitRates.maxBitRateForUL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  maxBitRateForDL:  %u\n",
        qos->bitRates.maxBitRateForDL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  guarBitRateForUL: %u\n",
        qos->bitRates.guarBitRateForUL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  guarBitRateForDL: %u\n",
        qos->bitRates.guarBitRateForDL);
  } else {
    OAILOG_DEBUG(LOG_NAS_ESM, "ESM-SAP   - epsqos  no bit rates defined\n");
  }

  if (qos->bitRatesExtPresent) {
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  maxBitRateForUL  Ext: %u\n",
        qos->bitRatesExt.maxBitRateForUL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  maxBitRateForDL  Ext: %u\n",
        qos->bitRatesExt.maxBitRateForDL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  guarBitRateForUL Ext: %u\n",
        qos->bitRatesExt.guarBitRateForUL);
    OAILOG_DEBUG(
        LOG_NAS_ESM, "ESM-SAP   - epsqos  guarBitRateForDL Ext: %u\n",
        qos->bitRatesExt.guarBitRateForDL);
  } else {
    OAILOG_DEBUG(LOG_NAS_ESM, "ESM-SAP   - epsqos  no bit rates ext defined\n");
  }

  if (apn == NULL) {
    OAILOG_WARNING(LOG_NAS_ESM, "ESM-SAP   - apn is NULL!\n");
  } else {
    OAILOG_DEBUG(LOG_NAS_ESM, "ESM-SAP   - apn is %s\n", bdata(apn));
  }
  /*
   * Mandatory - Access Point Name
   */
  msg->accesspointname = apn;
  /*
   * Mandatory - PDN address
   */
  OAILOG_DEBUG(LOG_NAS_ESM, "ESM-SAP   - pdn_type is %u\n", pdn_type);
  msg->pdnaddress.pdntypevalue = pdn_type;
  OAILOG_STREAM_HEX(
      OAILOG_LEVEL_DEBUG, LOG_NAS_ESM, "ESM-SAP   - pdn_addr is ",
      bdata(pdn_addr), blength(pdn_addr));
  msg->pdnaddress.pdnaddressinformation = pdn_addr;
  /*
   * Optional - ESM cause code
   */
  msg->presencemask = 0;

  if (esm_cause != ESM_CAUSE_SUCCESS) {
    msg->presencemask |=
        ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_ESM_CAUSE_PRESENT;
    msg->esmcause = esm_cause;
  }

  if (pco->num_protocol_or_container_id) {
    msg->presencemask |=
        ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;
    copy_protocol_configuration_options(
        &msg->protocolconfigurationoptions, pco);
  }
  //#pragma message  "TEST LG FORCE APN-AMBR"
  OAILOG_DEBUG(
      LOG_NAS_ESM, "ESM-SAP   - FORCE APN-AMBR DL %lu UL %lu\n",
      pdn_context_p->subscribed_apn_ambr.br_dl,
      pdn_context_p->subscribed_apn_ambr.br_ul);
  msg->presencemask |=
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_PRESENT;
  bit_rate_value_to_eps_qos(
      &msg->apnambr, pdn_context_p->subscribed_apn_ambr.br_dl,
      pdn_context_p->subscribed_apn_ambr.br_ul);

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-SAP   - Send Activate Default EPS Bearer Context "
      "Request message (pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_activate_dedicated_eps_bearer_context_request()  **
 **                                                                        **
 ** Description: Builds Activate Dedicated EPS Bearer Context Request      **
 **      message                                                   **
 **                                                                        **
 **      The activate dedicated EPS bearer context request message **
 **      is sent by the network to the UE to request activation of **
 **      a dedicated EPS bearer context associated with the same   **
 **      PDN address(es) and APN as an already active default EPS  **
 **      bearer context.                                           **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      linked_ebi:    EPS bearer identity of the default bearer  **
 **             associated with the EPS dedicated bearer   **
 **             to be activated                            **
 **      qos:       EPS quality of service                     **
 **      tft:       Traffic flow template                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_activate_dedicated_eps_bearer_context_request(
    pti_t pti, ebi_t ebi,
    activate_dedicated_eps_bearer_context_request_msg* msg, ebi_t linked_ebi,
    const EpsQualityOfService* qos, traffic_flow_template_t* tft,
    protocol_configuration_options_t* pco) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity     = ebi;
  msg->messagetype           = ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST;
  msg->proceduretransactionidentity = pti;
  msg->linkedepsbeareridentity      = linked_ebi;
  /*
   * Mandatory - EPS QoS
   */
  msg->epsqos = *qos;
  /*
   * Mandatory - traffic flow template
   */
  if (tft) {
    memcpy(&msg->tft, tft, sizeof(traffic_flow_template_t));
  }

  /*
   * Optional
   */
  msg->presencemask = 0;
  if (pco) {
    memcpy(
        &msg->protocolconfigurationoptions, pco,
        sizeof(protocol_configuration_options_t));
    msg->presencemask |=
        ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_IEI;
  }
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-SAP   - Send Activate Dedicated EPS Bearer Context "
      "Request message (pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_send_deactivate_eps_bearer_context_request()          **
 **                                                                        **
 ** Description: Builds Deactivate EPS Bearer Context Request message      **
 **                                                                        **
 **      The deactivate EPS bearer context request message is sent **
 **      by the network to request deactivation of an active EPS   **
 **      bearer context.                                           **
 **                                                                        **
 ** Inputs:  pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      esm_cause: ESM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The ESM message to be sent                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_send_deactivate_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, deactivate_eps_bearer_context_request_msg* msg,
    int esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  /*
   * Mandatory - ESM message header
   */
  msg->protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;
  msg->epsbeareridentity            = ebi;
  msg->messagetype                  = DEACTIVATE_EPS_BEARER_CONTEXT_REQUEST;
  msg->proceduretransactionidentity = pti;
  /*
   * Mandatory - ESM cause code
   */
  msg->esmcause = esm_cause;
  /*
   * Optional IEs
   */
  msg->presencemask = 0;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-SAP   - Send Deactivate EPS Bearer Context Request "
      "message (pti=%d, ebi=%d)\n",
      msg->proceduretransactionidentity, msg->epsbeareridentity);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
