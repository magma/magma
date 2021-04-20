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

  Source      mme_app_sgs_la_updated.c

  Version

  Date

  Product    MME app

  Subsystem  SGS (an interface between MME and MSC/VLR) state machine handling

  Author

  Description Implements the SGS procedures executed
        when the SGS state in SGS-LA_UPDATE_REQUESTED.

*****************************************************************************/

#include "common_defs.h"
#include "log.h"
#include "mme_app_sgs_fsm.h"
#include "common_types.h"

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
 ** Name:    sgs_la_update_requested_handler() **
 **                                                                        **
 ** Description: Handles the behaviour of the UE in MME while the          **
 **              SGS is in SGS-LA_UPDATE_REQUEST state.                    **
 **                                                                        **
 ** Inputs:  sgs_evt:   The received SGS event                             **
 **                                                                        **
 ** Outputs:                                                               **
 **          Return:    RETURNok, RETURNerror                              **
 **                                                                        **
 ***************************************************************************/
int sgs_la_update_requested_handler(const sgs_fsm_t* evt) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  int rc = RETURNerror;

  if (sgs_fsm_get_status(evt->ue_id, evt->ctx) != SGS_LA_UPDATE_REQUESTED) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS not in the SGS_LA_UPDATE_REQUESTED state for UE "
        "Id: " MME_UE_S1AP_ID_FMT "\n",
        evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  switch (evt->primitive) {
    case _SGS_LOCATION_UPDATE_ACCEPT: {
      rc = sgs_fsm_la_updt_req_loc_updt_acc(evt);
    } break;

    case _SGS_LOCATION_UPDATE_REJECT: {
      rc = sgs_fsm_la_updt_req_loc_updt_rej(evt);
    } break;

    case _SGS_PAGING_REQUEST: {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Handle paging request in SGS_LA_UPDATE_REQUESTED state for ue-id :"
          "" MME_UE_S1AP_ID_FMT " \n",
          evt->ue_id);
      rc = RETURNok;
    } break;

    case _SGS_EPS_DETACH_IND:
      /*
       * SGS EPS Detach procedure successful
       * enter state SGS-NULL.
       */
      rc = sgs_fsm_set_status(evt->ue_id, evt->ctx, SGS_NULL);
      break;

    case _SGS_IMSI_DETACH_IND:
      /*
       * SGS IMSI Detach procedure successful
       * enter state SGS-NULL.
       */
      rc = sgs_fsm_set_status(evt->ue_id, evt->ctx, SGS_NULL);
      break;

    case _SGS_RESET_INDICATION: {
      /* No handling required, if Reset indication received in
       * La-Update-Requested state */
      OAILOG_DEBUG(
          LOG_MME_APP,
          " Received Reset Indication while SGS context is in "
          "La-Update-Requested state for ue_id"
          " :%d \n",
          evt->ue_id);
      rc = RETURNok;
    } break;

    default:
      OAILOG_ERROR(
          LOG_MME_APP, "SGS-FSM   - Primitive is not valid (%d)\n",
          evt->primitive);
      break;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
