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

#include "assertions.h"
#include "log.h"
#include "common_defs.h"
#include "emm_main.h"
#include "mme_config.h"
#include "mme_api.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

emm_data_t _emm_data;

/****************************************************************************
 **                                                                        **
 ** Name:    emm_main_initialize()                                     **
 **                                                                        **
 ** Description: Initializes EMM internal data                             **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    _emm_data                                  **
 **                                                                        **
 ***************************************************************************/
void emm_main_initialize(const mme_config_t* mme_config_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Retreive MME supported configuration data
   */
  memset(&_emm_data.conf, 0, sizeof(_emm_data.conf));
  if (mme_api_get_emm_config(&_emm_data.conf, mme_config_p) != RETURNok) {
    Fatal(
        "EMM-MAIN  - Failed to get all required MME config data, "
        "check the error logs for missing configs");
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_main_cleanup()                                        **
 **                                                                        **
 ** Description: Performs the EPS Mobility Management clean up procedure   **
 **                                                                        **
 ** Inputs:  None                                                      **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **          Return:    None                                       **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void emm_main_cleanup(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
