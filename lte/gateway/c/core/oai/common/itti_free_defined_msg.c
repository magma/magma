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

/*! \file itti_free_defined_msg.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdlib.h>

#include "dynamic_memory_check.h"
#include "assertions.h"
#include "3gpp_24.008.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "async_system_messages_types.h"
#include "ip_forward_messages_types.h"
#include "s11_messages_types.h"
#include "sctp_messages_types.h"

//------------------------------------------------------------------------------
void itti_free_msg_content(MessageDef* const message_p) {
  switch (ITTI_MSG_ID(message_p)) {
    case ASYNC_SYSTEM_COMMAND: {
      if (ASYNC_SYSTEM_COMMAND(message_p).system_command) {
        bdestroy_wrapper(&ASYNC_SYSTEM_COMMAND(message_p).system_command);
      }
    } break;

    case GTPV1U_CREATE_TUNNEL_REQ:
    case GTPV1U_CREATE_TUNNEL_RESP:
    case GTPV1U_UPDATE_TUNNEL_REQ:
    case GTPV1U_UPDATE_TUNNEL_RESP:
    case GTPV1U_DELETE_TUNNEL_REQ:
    case GTPV1U_DELETE_TUNNEL_RESP:
      // DO nothing
      break;

    case GTPV1U_TUNNEL_DATA_IND:
    case GTPV1U_TUNNEL_DATA_REQ:
      // UNUSED actually
      break;

    case SGI_CREATE_ENDPOINT_REQUEST:
      break;

    case SGI_CREATE_ENDPOINT_RESPONSE: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.sgi_create_end_point_response.pco);
    } break;

    case SGI_UPDATE_ENDPOINT_REQUEST:
    case SGI_UPDATE_ENDPOINT_RESPONSE:
    case SGI_DELETE_ENDPOINT_REQUEST:
    case SGI_DELETE_ENDPOINT_RESPONSE:
      // DO nothing
      break;

    case MME_APP_CONNECTION_ESTABLISHMENT_CNF: {
      itti_mme_app_connection_establishment_cnf_t mme_app_est_cnf =
          message_p->ittiMsg.mme_app_connection_establishment_cnf;
      for (uint8_t index = 0; index < BEARERS_PER_UE; index++) {
        bdestroy_wrapper(&mme_app_est_cnf.nas_pdu[index]);
      }
      for (uint8_t index = 0; index < mme_app_est_cnf.no_of_e_rabs; index++) {
        bdestroy_wrapper(&(mme_app_est_cnf.transport_layer_address[index]));
      }
      bdestroy_wrapper(&mme_app_est_cnf.ue_radio_capability);
    } break;

    case MME_APP_INITIAL_CONTEXT_SETUP_RSP:
      break;

    case MME_APP_DELETE_SESSION_RSP:
      // DO nothing
      break;

    case MME_APP_UPLINK_DATA_IND:
      bdestroy_wrapper(&message_p->ittiMsg.mme_app_ul_data_ind.nas_msg);
      break;

    case MME_APP_HANDOVER_REQUEST:
      bdestroy_wrapper(
          &message_p->ittiMsg.mme_app_handover_request.src_tgt_container);
      break;
    case MME_APP_HANDOVER_COMMAND:
      bdestroy_wrapper(
          &message_p->ittiMsg.mme_app_handover_command.tgt_src_container);
      break;

    case S11_CREATE_SESSION_REQUEST: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.s11_create_session_request.pco);
    } break;

    case S11_CREATE_SESSION_RESPONSE: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.s11_create_session_response.pco);
    } break;

    case S11_CREATE_BEARER_REQUEST: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.s11_create_bearer_request.pco);
    } break;

    case S11_CREATE_BEARER_RESPONSE: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.s11_create_bearer_response.pco);
    } break;

    case S11_MODIFY_BEARER_REQUEST:
    case S11_MODIFY_BEARER_RESPONSE:
    case S11_DELETE_SESSION_REQUEST:
      // DO nothing (trxn)
      break;

    case S11_DELETE_SESSION_RESPONSE: {
      clear_protocol_configuration_options(
          &message_p->ittiMsg.s11_delete_session_response.pco);
    } break;

    case S11_RELEASE_ACCESS_BEARERS_REQUEST:
    case S11_RELEASE_ACCESS_BEARERS_RESPONSE:
      // DO nothing (trxn)
      break;
    case S11_PAGING_REQUEST: {
      if (message_p->ittiMsg.s11_paging_request.imsi) {
        free_wrapper((void**) &message_p->ittiMsg.s11_paging_request.imsi);
      }
    } break;

    case S1AP_ENB_INITIATED_RESET_ACK:
      free_wrapper((void**) &message_p->ittiMsg.s1ap_enb_initiated_reset_ack
                       .ue_to_reset_list);
      break;
    case S1AP_UE_CAPABILITIES_IND:
    case S1AP_ENB_DEREGISTERED_IND:
    case S1AP_UE_CONTEXT_RELEASE_REQ:
    case S1AP_UE_CONTEXT_RELEASE_COMMAND:
    case S1AP_UE_CONTEXT_RELEASE_COMPLETE:
      // DO nothing
      break;

    case S1AP_ENB_INITIATED_RESET_REQ:
      // Do Nothing
      // No need to free ue_to_reset_list in "S1AP_ENB_INITIATED_RESET_REQ"
      // because it is re-used in another ITTI message
      break;
    case S1AP_E_RAB_REL_CMD:
      bdestroy_wrapper(&message_p->ittiMsg.s1ap_e_rab_rel_cmd.nas_pdu);
      break;
    case S1AP_E_RAB_SETUP_REQ:
      bdestroy_wrapper(
          &message_p->ittiMsg.s1ap_e_rab_setup_req.e_rab_to_be_setup_list
               .item[0]
               .transport_layer_address);
      bdestroy_wrapper(
          &message_p->ittiMsg.s1ap_e_rab_setup_req.e_rab_to_be_setup_list
               .item[0]
               .nas_pdu);
      break;
    case S1AP_NAS_DL_DATA_REQ:
      bdestroy_wrapper(&message_p->ittiMsg.s1ap_nas_dl_data_req.nas_msg);
      break;
    case S1AP_HANDOVER_REQUIRED:
      bdestroy_wrapper(
          &message_p->ittiMsg.s1ap_handover_required.src_tgt_container);
      break;
    case S1AP_HANDOVER_REQUEST_ACK:
      bdestroy_wrapper(
          &message_p->ittiMsg.s1ap_handover_request_ack.tgt_src_container);
      break;
    case S6A_UPDATE_LOCATION_REQ:
    case S6A_UPDATE_LOCATION_ANS:
    case S6A_AUTH_INFO_REQ:
    case S6A_AUTH_INFO_ANS:
      // DO nothing
      break;

    case S11_NW_INITIATED_DEACTIVATE_BEARER_RESP:
      free_wrapper(
          (void**) &message_p->ittiMsg.s11_nw_init_deactv_bearer_rsp.lbi);
      break;

    case SCTP_INIT_MSG:
      // DO nothing (ipv6_address statically allocated)
      break;

    case SCTP_DATA_REQ:
      bdestroy_wrapper(&message_p->ittiMsg.sctp_data_req.payload);
      break;

    case SCTP_DATA_IND:
      bdestroy_wrapper(&message_p->ittiMsg.sctp_data_ind.payload);
      break;

    case SCTP_NEW_ASSOCIATION:
      bdestroy_wrapper(&message_p->ittiMsg.sctp_new_peer.ran_cp_ipaddr);
      break;

    case SCTP_DATA_CNF:
    case SCTP_CLOSE_ASSOCIATION:
      // DO nothing
      break;
    default:;
  }
}
