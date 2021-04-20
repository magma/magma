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

#include <assert.h>

#include "common_defs.h"
#include "log.h"
#include "emm_reg.h"
#include "emm_fsm.h"

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
 ** Name:    emm_reg_initialize()                                      **
 **                                                                        **
 ** Description: Initializes the EMMREG Service Access Point               **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    NONE                                       **
 **                                                                        **
 ***************************************************************************/
void emm_reg_initialize(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Initialize the EMM state machine
   */
  emm_fsm_initialize();
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_reg_send()                                            **
 **                                                                        **
 ** Description: Processes the EMMREG Service Access Point primitive       **
 **                                                                        **
 ** Inputs:  msg:       The EMMREG-SAP primitive to process        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_reg_send(emm_reg_t* const msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNok;

  /*
   * Check the EMM-SAP primitive
   */
  emm_reg_primitive_t primitive = msg->primitive;

  assert((primitive > _EMMREG_START) && (primitive < _EMMREG_END));
  /*
   * Execute the EMM procedure
   */
  rc = emm_fsm_process(msg);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
