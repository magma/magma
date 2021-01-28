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

#include "log.h"
#include "common_defs.h"
#include "emm_sap.h"
#include "emm_reg.h"
#include "emm_esm.h"
#include "emm_as.h"
#include "emm_cn.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_sap_initialize()                                      **
 **                                                                        **
 ** Description: Initializes the EMM Service Access Points                 **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    NONE                                       **
 **                                                                        **
 ***************************************************************************/
void emm_sap_initialize(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_reg_initialize();
  emm_esm_initialize();
  emm_as_initialize();
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_sap_send()                                            **
 **                                                                        **
 ** Description: Processes the EMM Service Access Point primitive          **
 **                                                                        **
 ** Inputs:  msg:       The EMM-SAP primitive to process           **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The EMM-SAP primitive to process           **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_sap_send(emm_sap_t* msg) {
  int rc                    = RETURNerror;
  emm_primitive_t primitive = msg->primitive;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  /*
   * Check the EMM-SAP primitive
   */
  if ((primitive > (emm_primitive_t) EMMREG_PRIMITIVE_MIN) &&
      (primitive < (emm_primitive_t) EMMREG_PRIMITIVE_MAX)) {
    /*
     * Forward to the EMMREG-SAP
     */
    msg->u.emm_reg.primitive = primitive;
    rc                       = emm_reg_send(&msg->u.emm_reg);
  } else if (
      (primitive > (emm_primitive_t) EMMESM_PRIMITIVE_MIN) &&
      (primitive < (emm_primitive_t) EMMESM_PRIMITIVE_MAX)) {
    /*
     * Forward to the EMMESM-SAP
     */
    msg->u.emm_esm.primitive = primitive;
    rc                       = emm_esm_send(&msg->u.emm_esm);
  } else if (
      (primitive > (emm_primitive_t) EMMAS_PRIMITIVE_MIN) &&
      (primitive < (emm_primitive_t) EMMAS_PRIMITIVE_MAX)) {
    /*
     * Forward to the EMMAS-SAP
     */
    msg->u.emm_as.primitive = primitive;
    rc                      = emm_as_send(&msg->u.emm_as);
  } else if (
      (primitive > (emm_primitive_t) EMMCN_PRIMITIVE_MIN) &&
      (primitive < (emm_primitive_t) EMMCN_PRIMITIVE_MAX)) {
    /*
     * Forward to the EMMCN-SAP
     */
    msg->u.emm_cn.primitive = primitive;
    rc                      = emm_cn_send(&msg->u.emm_cn);
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM, "EMM-SAP -   Out of range primitive (%d)\n", primitive);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
