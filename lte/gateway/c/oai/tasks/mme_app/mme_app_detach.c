/*******************************************************************************
    OpenAirInterface
    Copyright(c) 1999 - 2014 Eurecom

    OpenAirInterface is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.


    OpenAirInterface is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with OpenAirInterface.The full GNU General Public License is
   included in this distribution in the file called "COPYING". If not,
   see <http://www.gnu.org/licenses/>.

  Contact Information
  OpenAirInterface Admin: openair_admin@eurecom.fr
  OpenAirInterface Tech : openair_tech@eurecom.fr
  OpenAirInterface Dev  : openair4g-devel@eurecom.fr

  Address      : Eurecom, Compus SophiaTech 450, route des chappes, 06451 Biot, France.

 *******************************************************************************/

#include <stdio.h>
#include <stdbool.h>

#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "gcc_diag.h"
#include "mme_config.h"
#include "mme_app_ue_context.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_defs.h"
#include "3gpp_24.007.h"
#include "3gpp_29.274.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "nas_messages_types.h"
#include "s11_messages_types.h"
#include "s1ap_messages_types.h"
#include "service303.h"
#include "sgw_ie_defs.h"

#if EMBEDDED_SGW
#define TASK_SPGW TASK_SPGW_APP
#else
#define TASK_SPGW TASK_S11
#endif

//------------------------------------------------------------------------------
void mme_app_send_delete_session_request(
  struct ue_mm_context_s *const ue_context_p,
  const ebi_t ebi,
  const pdn_cid_t cid)
{
  MessageDef *message_p = NULL;
  OAILOG_FUNC_IN(LOG_MME_APP);
  message_p = itti_alloc_new_message(TASK_MME_APP, S11_DELETE_SESSION_REQUEST);
  AssertFatal(message_p, "itti_alloc_new_message Failed");
  S11_DELETE_SESSION_REQUEST(message_p).local_teid = ue_context_p->mme_teid_s11;
  S11_DELETE_SESSION_REQUEST(message_p).teid =
    ue_context_p->pdn_contexts[cid]->s_gw_teid_s11_s4;
  S11_DELETE_SESSION_REQUEST(message_p).lbi = ebi; //default bearer

  /* clang-format off */
  OAI_GCC_DIAG_OFF(pointer-to-int-cast);
  /* clang-format on */
  S11_DELETE_SESSION_REQUEST(message_p).sender_fteid_for_cp.teid =
    (teid_t) ue_context_p;
  OAI_GCC_DIAG_ON(pointer - to - int - cast);
  S11_DELETE_SESSION_REQUEST(message_p).sender_fteid_for_cp.interface_type =
    S11_MME_GTP_C;
  mme_config_read_lock(&mme_config);
  S11_DELETE_SESSION_REQUEST(message_p).sender_fteid_for_cp.ipv4_address =
    mme_config.ipv4.s11;
  mme_config_unlock(&mme_config);
  S11_DELETE_SESSION_REQUEST(message_p).sender_fteid_for_cp.ipv4 = 1;
  S11_DELETE_SESSION_REQUEST(message_p).indication_flags.oi = 1;

  /*
   * S11 stack specific parameter. Not used in standalone epc mode
   */
  S11_DELETE_SESSION_REQUEST(message_p).trxn = NULL;
  mme_config_read_lock(&mme_config);
  S11_DELETE_SESSION_REQUEST(message_p).peer_ip =
    ue_context_p->pdn_contexts[cid]->s_gw_address_s11_s4.address.ipv4_address;
  mme_config_unlock(&mme_config);

  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);
  increment_counter("mme_spgw_delete_session_req", 1, NO_LABELS);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_detach_req(mme_app_desc_t *mme_app_desc_p,
    const itti_nas_detach_req_t *const detach_req_p)
{
  struct ue_mm_context_s *ue_context = NULL;

  DevAssert(detach_req_p != NULL);
  ue_context = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc_p->mme_ue_contexts, detach_req_p->ue_id);
  if (ue_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "UE context doesn't exist -> Nothing to do :-) \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if ((!ue_context->mme_teid_s11) && (!ue_context->nb_active_pdn_contexts)) {
    /* No Session.
     * If UE is already in idle state, skip asking eNB to release UE context and
     * just clean up locally.
     */
    if (ECM_IDLE == ue_context->ecm_state) {
      ue_context->ue_context_rel_cause = S1AP_IMPLICIT_CONTEXT_RELEASE;
      // Notify S1AP to release S1AP UE context locally.
      mme_app_itti_ue_context_release(
        ue_context, ue_context->ue_context_rel_cause);
      // Free MME UE Context
      mme_notify_ue_context_released(&mme_app_desc_p->mme_ue_contexts,
          ue_context);
      // Send PUR,before removal of ue contexts
      if (
        (ue_context->send_ue_purge_request == true) &&
        (ue_context->hss_initiated_detach == false)) {
        mme_app_send_s6a_purge_ue_req(mme_app_desc_p, ue_context);
      }
      mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context);
    } else {
      if (ue_context->ue_context_rel_cause == S1AP_INVALID_CAUSE) {
        ue_context->ue_context_rel_cause = S1AP_NAS_DETACH;
      }
      // Notify S1AP to send UE Context Release Command to eNB.
      mme_app_itti_ue_context_release(
        ue_context, ue_context->ue_context_rel_cause);
      if (ue_context->ue_context_rel_cause == S1AP_SCTP_SHUTDOWN_OR_RESET) {
        // Just cleanup the MME APP state associated with s1.
        mme_ue_context_update_ue_sig_connection_state(
          &mme_app_desc_p->mme_ue_contexts, ue_context, ECM_IDLE);
        // Free MME UE Context
        mme_notify_ue_context_released(
          &mme_app_desc_p->mme_ue_contexts, ue_context);
        // Send PUR,before removal of ue contexts
        if (
          (ue_context->send_ue_purge_request == true) &&
          (ue_context->hss_initiated_detach == false)) {
          mme_app_send_s6a_purge_ue_req(mme_app_desc_p, ue_context);
        }
        mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context);
      } else {
        ue_context->ue_context_rel_cause = S1AP_INVALID_CAUSE;
        unlock_ue_contexts(ue_context);
      }
    }
  } else {
    for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
      if (ue_context->pdn_contexts[i]) {
        // Send a DELETE_SESSION_REQUEST message to the SGW
        mme_app_send_delete_session_request(
          ue_context, ue_context->pdn_contexts[i]->default_ebi, i);
      }
    }
    unlock_ue_contexts(ue_context);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
