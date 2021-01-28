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

  Source      mme_app_sgs_associated.c

  Version

  Date

  Product    MME app

  Subsystem  SGS (an interface between MME and MSC/VLR) state machine handling

  Author

  Description Implements the SGS procedures executed
        when the SGS state in SGS-associated.

*****************************************************************************/

#include "common_defs.h"
#include "log.h"
#include "mme_app_sgs_fsm.h"
#include "mme_app_defs.h"

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
 ** Name:    sgs_associated_handler                                        **
 **                                                                        **
 ** Description: Handles the SGSAP messages for UE while the               **
 **              SGS is in SGS-Associated state.                           **
 **                                                                        **
 ** Inputs:  sgs_evt:   The received SGS event                             **
 **                                                                        **
 ** Outputs:                                                               **
 **          Return:    RETURNok, RETURNerror                              **
 **                                                                        **
 ***************************************************************************/
int sgs_associated_handler(const sgs_fsm_t* evt) {
  int rc = RETURNerror;
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (sgs_fsm_get_status(evt->ue_id, evt->ctx) != SGS_ASSOCIATED) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS not in the SGS_Associated state, UE Id: " MME_UE_S1AP_ID_FMT "\n",
        evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  switch (evt->primitive) {
    case _SGS_LOCATION_UPDATE_ACCEPT: {
      rc = sgs_fsm_associated_loc_updt_acc(evt);
    } break;

    case _SGS_LOCATION_UPDATE_REJECT: {
      rc = sgs_fsm_associated_loc_updt_rej(evt);
    } break;

    case _SGS_PAGING_REQUEST: {
      rc = sgs_handle_associated_paging_request(evt);
    } break;

    case _SGS_EPS_DETACH_IND: {
      /*
       * SGS EPS Detach procedure successful
       * enter state SGS-NULL.
       */
      rc = sgs_fsm_set_status(evt->ue_id, evt->ctx, SGS_NULL);
    } break;

    case _SGS_SERVICE_ABORT_REQUEST: {
      rc = sgs_fsm_associated_service_abort_request(evt);
    } break;

    case _SGS_IMSI_DETACH_IND: {
      /*
       * SGS IMSI Detach procedure successful
       * enter state SGS-NULL.
       */
      rc = sgs_fsm_set_status(evt->ue_id, evt->ctx, SGS_NULL);
    } break;

    case _SGS_RESET_INDICATION: {
      rc = sgs_fsm_associated_reset_indication(evt);
    } break;

    default: {
      OAILOG_ERROR(
          LOG_MME_APP, "SGS-FSM   - Primitive is not valid (%d)\n",
          evt->primitive);
    } break;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
