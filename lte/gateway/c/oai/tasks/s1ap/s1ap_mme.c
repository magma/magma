/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file s1ap_mme.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#if HAVE_CONFIG_H
#include "config.h"
#endif

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>
#include <netinet/in.h>

#include "bstrlib.h"
#include "hashtable.h"
#include "log.h"
#include "msc.h"
#include "assertions.h"
#include "mme_app_statistics.h"
#include "s1ap_mme.h"
#include "s1ap_mme_decoder.h"
#include "s1ap_mme_handlers.h"
#include "s1ap_state.h"
#include "s1ap_ies_defs.h"
#include "s1ap_mme_nas_procedures.h"
#include "s1ap_mme_itti_messaging.h"
#include "service303.h"
#include "dynamic_memory_check.h"
#include "mme_config.h"
#include "timer.h"
#include "itti_free_defined_msg.h"
#include "S1ap-TimeToWait.h"
#include "asn_internal.h"
#include "common_defs.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_messages_types.h"
#include "mme_default_values.h"
#include "s1ap_messages_types.h"
#include "sctp_messages_types.h"
#include "timer_messages_types.h"

#if S1AP_DEBUG_LIST
#define eNB_LIST_OUT(x, args...)                                               \
  OAILOG_DEBUG(LOG_S1AP, "[eNB]%*s" x "\n", 4 * indent, "", ##args)
#define UE_LIST_OUT(x, args...)                                                \
  OAILOG_DEBUG(LOG_S1AP, "[UE] %*s" x "\n", 4 * indent, "", ##args)
#else
#define eNB_LIST_OUT(x, args...)
#define UE_LIST_OUT(x, args...)
#endif

bool hss_associated = false;
uint32_t nb_enb_associated = 0;

hash_table_ts_t g_s1ap_enb_coll = {
  .mutex = PTHREAD_MUTEX_INITIALIZER,
  0}; // contains eNB_description_s, key is eNB_description_s.enb_id (uint32_t);
hash_table_ts_t g_s1ap_mme_id2assoc_id_coll = {
  .mutex = PTHREAD_MUTEX_INITIALIZER,
  0}; // contains sctp association id, key is mme_ue_s1ap_id;

static int indent = 0;
void *s1ap_mme_thread(void *args);

//------------------------------------------------------------------------------
static int s1ap_send_init_sctp(void)
{
  // Create and alloc new message
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_S1AP, SCTP_INIT_MSG);
  message_p->ittiMsg.sctpInit.port = S1AP_PORT_NUMBER;
  message_p->ittiMsg.sctpInit.ppid = S1AP_SCTP_PPID;
  message_p->ittiMsg.sctpInit.ipv4 = 1;
  message_p->ittiMsg.sctpInit.ipv6 = 0;
  message_p->ittiMsg.sctpInit.nb_ipv4_addr = 1;
  message_p->ittiMsg.sctpInit.ipv4_address[0].s_addr =
    mme_config.ipv4.s1_mme.s_addr;
  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * * * * Disable it for now.
   */
  message_p->ittiMsg.sctpInit.nb_ipv6_addr = 0;
  message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  return itti_send_msg_to_task(TASK_SCTP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
void *s1ap_mme_thread(__attribute__((unused)) void *args)
{
  itti_mark_task_ready(TASK_S1AP);

  while (1) {
    MessageDef *received_message_p = NULL;
    MessagesIds message_id = MESSAGES_ID_MAX;
    /*
     * Trying to fetch a message from the message queue.
     * * * * If the queue is empty, this function will block till a
     * * * * message is sent to the task.
     */
    itti_receive_msg(TASK_S1AP, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case ACTIVATE_MESSAGE: {
        hss_associated = true;
      } break;

      case MESSAGE_TEST:
        OAILOG_DEBUG(LOG_S1AP, "Received MESSAGE_TEST\n");
        break;

      case SCTP_DATA_IND: {
        /*
         * New message received from SCTP layer.
         * * * * Decode and handle it.
         */
        s1ap_message message = {0};

        /*
         * Invoke S1AP message decoder
         */
        if (
          s1ap_mme_decode_pdu(
            &message, SCTP_DATA_IND(received_message_p).payload, &message_id) <
          0) {
          // TODO: Notify eNB of failure with right cause
          OAILOG_ERROR(LOG_S1AP, "Failed to decode new buffer\n");
        } else {
          s1ap_mme_handle_message(
            SCTP_DATA_IND(received_message_p).assoc_id,
            SCTP_DATA_IND(received_message_p).stream,
            &message);
        }

        if (message_id != MESSAGES_ID_MAX) {
          s1ap_free_mme_decode_pdu(&message, message_id);
        }

        /*
         * Free received PDU array
         */
        bdestroy_wrapper(&SCTP_DATA_IND(received_message_p).payload);
      } break;

      case SCTP_DATA_CNF:
        s1ap_mme_itti_nas_downlink_cnf(
          SCTP_DATA_CNF(received_message_p).mme_ue_s1ap_id,
          SCTP_DATA_CNF(received_message_p).is_success);
        break;
        /*
       * SCTP layer notifies S1AP of disconnection of a peer.
       */
      case SCTP_CLOSE_ASSOCIATION: {
        s1ap_handle_sctp_disconnection(
          SCTP_CLOSE_ASSOCIATION(received_message_p).assoc_id,
          SCTP_CLOSE_ASSOCIATION(received_message_p).reset);
      } break;

      case SCTP_NEW_ASSOCIATION: {
        increment_counter("mme_new_association", 1, NO_LABELS);
        if (s1ap_handle_new_association(
              &received_message_p->ittiMsg.sctp_new_peer)) {
          increment_counter("mme_new_association", 1, 1, "result", "failure");
        } else {
          increment_counter("mme_new_association", 1, 1, "result", "success");
        }
      } break;

      case S1AP_NAS_DL_DATA_REQ: {
        /*
         * New message received from NAS task.
         * * * * This corresponds to a S1AP downlink nas transport message.
         */
        s1ap_generate_downlink_nas_transport(
          S1AP_NAS_DL_DATA_REQ(received_message_p).enb_ue_s1ap_id,
          S1AP_NAS_DL_DATA_REQ(received_message_p).mme_ue_s1ap_id,
          &S1AP_NAS_DL_DATA_REQ(received_message_p).nas_msg);
      } break;

      case S1AP_E_RAB_SETUP_REQ: {
        s1ap_generate_s1ap_e_rab_setup_req(
          &S1AP_E_RAB_SETUP_REQ(received_message_p));
      } break;

      // From MME_APP task
      case S1AP_UE_CONTEXT_RELEASE_COMMAND: {
        s1ap_handle_ue_context_release_command(
          &received_message_p->ittiMsg.s1ap_ue_context_release_command);
      } break;

      case MME_APP_CONNECTION_ESTABLISHMENT_CNF: {
        s1ap_handle_conn_est_cnf(
          &MME_APP_CONNECTION_ESTABLISHMENT_CNF(received_message_p));
      } break;

      case MME_APP_S1AP_MME_UE_ID_NOTIFICATION: {
        s1ap_handle_mme_ue_id_notification(
          &MME_APP_S1AP_MME_UE_ID_NOTIFICATION(received_message_p));
      } break;

      case S1AP_ENB_INITIATED_RESET_ACK: {
        s1ap_handle_enb_initiated_reset_ack(
          &S1AP_ENB_INITIATED_RESET_ACK(received_message_p));
      } break;

      case S1AP_PAGING_REQUEST: {
        if (
          s1ap_handle_paging_request(
            &S1AP_PAGING_REQUEST(received_message_p)) != RETURNok) {
          OAILOG_ERROR(LOG_S1AP, "Failed to send paging message\n");
        }
      } break;
      case S1AP_UE_CONTEXT_MODIFICATION_REQUEST: {
        s1ap_handle_ue_context_mod_req(
          &received_message_p->ittiMsg.s1ap_ue_context_mod_request);
      } break;

      case TIMER_HAS_EXPIRED: {
        if (!timer_exists(
              received_message_p->ittiMsg.timer_has_expired.timer_id)) {
          break;
        }
        ue_description_t *ue_ref_p = NULL;
        enb_description_t *enb_ref_p = NULL;
        if (received_message_p->ittiMsg.timer_has_expired.arg != NULL) {
          // check whether timer is related to eNB procedure or UE procedure
          s1ap_timer_arg_t timer_arg =
            *((s1ap_timer_arg_t *) (received_message_p->ittiMsg
                                      .timer_has_expired.arg));
          if (timer_arg.timer_class == S1AP_UE_TIMER) {
            mme_ue_s1ap_id_t mme_ue_s1ap_id = timer_arg.instance_id;
            if (
              (ue_ref_p = s1ap_is_ue_mme_id_in_list(mme_ue_s1ap_id)) == NULL) {
              OAILOG_WARNING(
                LOG_S1AP,
                "Timer expired but no assoicated UE context for UE id %d\n",
                mme_ue_s1ap_id);
              timer_handle_expired(
                received_message_p->ittiMsg.timer_has_expired.timer_id);
              break;
            }
            if (
              received_message_p->ittiMsg.timer_has_expired.timer_id ==
              ue_ref_p->s1ap_ue_context_rel_timer.id) {
              // UE context release complete timer expiry handler
              OAILOG_INFO(
                LOG_S1AP,
                "ue_context_release_command_timer_expired for UE id %d\n",
                mme_ue_s1ap_id);
              increment_counter(
                "ue_context_release_command_timer_expired", 1, NO_LABELS);
              s1ap_mme_handle_ue_context_rel_comp_timer_expiry(ue_ref_p);
            }
          } else if (timer_arg.timer_class == S1AP_ENB_TIMER) {
            sctp_assoc_id_t assoc_id = timer_arg.instance_id;
            if ((enb_ref_p = s1ap_is_enb_assoc_id_in_list(assoc_id)) == NULL) {
              OAILOG_WARNING(
                LOG_S1AP,
                "Timer expired but no assoicated eNB context for eNB assoc_id "
                "%d\n",
                assoc_id);
              timer_handle_expired(
                received_message_p->ittiMsg.timer_has_expired.timer_id);
              break;
            }
            if (
              received_message_p->ittiMsg.timer_has_expired.timer_id ==
              enb_ref_p->s1ap_enb_assoc_clean_up_timer.id) {
              OAILOG_INFO(
                LOG_S1AP,
                "enb_sctp_shutdown_ue_clean_up_timer_expired for enb assoc_id "
                "%d\n",
                assoc_id);
              increment_counter(
                "enb_sctp_shutdown_ue_clean_up_timer_expired", 1, NO_LABELS);
              s1ap_enb_assoc_clean_up_timer_expiry(enb_ref_p);
            }
          } else {
            OAILOG_WARNING(
              LOG_S1AP,
              " S1AP Timer expired with invalid timer class  %u \n",
              timer_arg.timer_class);
          }
        }
        timer_handle_expired(
          received_message_p->ittiMsg.timer_has_expired.timer_id);

        /* TODO - Commenting out below function as it is not used as of now.
         * Need to handle it when we support other timers in S1AP
         */

        //s1ap_handle_timer_expiry (&received_message_p->ittiMsg.timer_has_expired);
      } break;

      case TERMINATE_MESSAGE: {
        s1ap_mme_exit();
        itti_free_msg_content(received_message_p);
        itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
        OAI_FPRINTF_INFO("TASK_S1AP terminated\n");
        itti_exit_task();
      } break;

      default: {
        OAILOG_ERROR(
          LOG_S1AP,
          "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int s1ap_mme_init(void)
{
  OAILOG_DEBUG(LOG_S1AP, "Initializing S1AP interface\n");

  if (get_asn1c_environment_version() < ASN1_MINIMUM_VERSION) {
    OAILOG_ERROR(
      LOG_S1AP,
      "ASN1C version %d fount, expecting at least %d\n",
      get_asn1c_environment_version(),
      ASN1_MINIMUM_VERSION);
    return RETURNerror;
  } else {
    OAILOG_DEBUG(
      LOG_S1AP, "ASN1C version %d\n", get_asn1c_environment_version());
  }

  OAILOG_DEBUG(LOG_S1AP, "S1AP Release v10.5\n");
  // 16 entries for n eNB.
  bstring bs1 = bfromcstr("s1ap_eNB_coll");
  hash_table_ts_t *h = hashtable_ts_init(
    &g_s1ap_enb_coll, mme_config.max_enbs, NULL, free_wrapper, bs1);
  bdestroy_wrapper(&bs1);
  if (!h) return RETURNerror;

  bstring bs2 = bfromcstr("s1ap_mme_id2assoc_id_coll");
  h = hashtable_ts_init(
    &g_s1ap_mme_id2assoc_id_coll,
    mme_config.max_ues,
    NULL,
    hash_free_int_func,
    bs2);
  bdestroy_wrapper(&bs2);
  if (!h) return RETURNerror;

  if (itti_create_task(TASK_S1AP, &s1ap_mme_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_S1AP, "Error while creating S1AP task\n");
    return RETURNerror;
  }

  if (s1ap_send_init_sctp() < 0) {
    OAILOG_ERROR(LOG_S1AP, "Error while sendind SCTP_INIT_MSG to SCTP \n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_S1AP, "Initializing S1AP interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
void s1ap_mme_exit(void)
{
  OAILOG_DEBUG(LOG_S1AP, "Cleaning S1AP\n");
  if (hashtable_ts_destroy(&g_s1ap_enb_coll) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying s1 eNB hash table");
  }
  if (hashtable_ts_destroy(&g_s1ap_mme_id2assoc_id_coll) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying assoc_id hash table");
  }
  OAILOG_DEBUG(LOG_S1AP, "Cleaning S1AP: DONE\n");
}

//------------------------------------------------------------------------------
void s1ap_dump_enb_list(void)
{
  hashtable_ts_apply_callback_on_elements(
    &g_s1ap_enb_coll, s1ap_dump_enb_hash_cb, NULL, NULL);
}

//------------------------------------------------------------------------------
bool s1ap_dump_enb_hash_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const eNB_void,
  void __attribute__((unused)) * unused_parameterP,
  void __attribute__((unused)) * *unused_resultP)
{
  const enb_description_t *const enb_ref = (const enb_description_t *) eNB_void;
  if (enb_ref == NULL) {
    return false;
  }
  s1ap_dump_enb(enb_ref);
  return false;
}

//------------------------------------------------------------------------------
void s1ap_dump_enb(const enb_description_t *const enb_ref)
{
#ifdef S1AP_DEBUG_LIST
  //Reset indentation
  indent = 0;

  if (enb_ref == NULL) {
    return;
  }

  eNB_LIST_OUT("");
  eNB_LIST_OUT(
    "eNB name:          %s",
    enb_ref->enb_name == NULL ? "not present" : enb_ref->enb_name);
  eNB_LIST_OUT("eNB ID:            %07x", enb_ref->enb_id);
  eNB_LIST_OUT("SCTP assoc id:     %d", enb_ref->sctp_assoc_id);
  eNB_LIST_OUT("SCTP instreams:    %d", enb_ref->instreams);
  eNB_LIST_OUT("SCTP outstreams:   %d", enb_ref->outstreams);
  eNB_LIST_OUT("UE attache to eNB: %d", enb_ref->nb_ue_associated);
  indent++;
  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t *const) & enb_ref->ue_coll,
    s1ap_dump_ue_hash_cb,
    NULL,
    NULL);
  indent--;
  eNB_LIST_OUT("");
#else
  s1ap_dump_ue(NULL);
#endif
}

//------------------------------------------------------------------------------
bool s1ap_dump_ue_hash_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const ue_void,
  void __attribute__((unused)) * unused_parameterP,
  void __attribute__((unused)) * *unused_resultP)
{
  ue_description_t *ue_ref = (ue_description_t *) ue_void;
  if (ue_ref == NULL) {
    return false;
  }
  s1ap_dump_ue(ue_ref);
  return false;
}
//------------------------------------------------------------------------------
void s1ap_dump_ue(const ue_description_t *const ue_ref)
{
#ifdef S1AP_DEBUG_LIST

  if (ue_ref == NULL) return;

  UE_LIST_OUT("eNB UE s1ap id:   0x%06x", ue_ref->enb_ue_s1ap_id);
  UE_LIST_OUT("MME UE s1ap id:   0x%08x", ue_ref->mme_ue_s1ap_id);
  UE_LIST_OUT("SCTP stream recv: 0x%04x", ue_ref->sctp_stream_recv);
  UE_LIST_OUT("SCTP stream send: 0x%04x", ue_ref->sctp_stream_send);
#endif
}

//------------------------------------------------------------------------------
bool s1ap_enb_compare_by_enb_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  const uint32_t *const enb_id_p = (const uint32_t *const) parameterP;
  enb_description_t *enb_ref = (enb_description_t *) elementP;
  if (*enb_id_p == enb_ref->enb_id) {
    *resultP = elementP;
    return true;
  }
  return false;
}
//------------------------------------------------------------------------------
enb_description_t *s1ap_is_enb_id_in_list(const uint32_t enb_id)
{
  enb_description_t *enb_ref = NULL;
  uint32_t *enb_id_p = (uint32_t *) &enb_id;
  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t *const) & g_s1ap_enb_coll,
    s1ap_enb_compare_by_enb_id_cb,
    (void *) enb_id_p,
    (void **) &enb_ref);
  return enb_ref;
}

//------------------------------------------------------------------------------
enb_description_t *s1ap_is_enb_assoc_id_in_list(
  const sctp_assoc_id_t sctp_assoc_id)
{
  enb_description_t *enb_ref = NULL;
  hashtable_ts_get(
    &g_s1ap_enb_coll, (const hash_key_t) sctp_assoc_id, (void **) &enb_ref);
  return enb_ref;
}

//------------------------------------------------------------------------------
// TODO: (amar) unused function check with OAI.
bool s1ap_ue_compare_by_enb_ue_s1ap_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *elementP,
  void *parameterP,
  void **resultP)
{
  enb_ue_s1ap_id_t *enb_ue_s1ap_id = (enb_ue_s1ap_id_t *) parameterP;
  ue_description_t *ue_ref = (ue_description_t *) elementP;
  if (*enb_ue_s1ap_id == ue_ref->enb_ue_s1ap_id) {
    *resultP = elementP;
    return true;
  }
  return false;
}
//------------------------------------------------------------------------------
ue_description_t *s1ap_is_ue_enb_id_in_list(
  enb_description_t *enb_ref,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id)
{
  ue_description_t *ue_ref = NULL;
  hashtable_ts_get(
    (hash_table_ts_t *const) & enb_ref->ue_coll,
    (const hash_key_t) enb_ue_s1ap_id,
    (void **) &ue_ref);
  return ue_ref;
}

//------------------------------------------------------------------------------
bool s1ap_ue_compare_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  mme_ue_s1ap_id_t *mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t *) parameterP;
  ue_description_t *ue_ref = (ue_description_t *) elementP;
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id) {
    *resultP = elementP;
    OAILOG_TRACE(
      LOG_S1AP,
      "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      ue_ref,
      ue_ref->mme_ue_s1ap_id);
    return true;
  }
  return false;
}

//------------------------------------------------------------------------------
bool s1ap_enb_find_ue_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  enb_description_t *enb_ref = (enb_description_t *) elementP;

  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t *const) & enb_ref->ue_coll,
    s1ap_ue_compare_by_mme_ue_id_cb,
    parameterP,
    resultP);
  if (*resultP) {
    OAILOG_TRACE(
      LOG_S1AP,
      "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      *resultP,
      ((ue_description_t *) (*resultP))->mme_ue_s1ap_id);
    return true;
  }
  return false;
}
//------------------------------------------------------------------------------
bool s1ap_ue_compare_by_s11_sgw_teid_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  s11_teid_t *s11_sgw_teid_p = (s11_teid_t *) parameterP;
  ue_description_t *ue_ref = (ue_description_t *) elementP;
  if (*s11_sgw_teid_p == ue_ref->s11_sgw_teid) {
    *resultP = elementP;
    return true;
  }
  return false;
}

//------------------------------------------------------------------------------
bool s1ap_enb_find_ue_by_s11_sgw_teid_cb(
  __attribute((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  enb_description_t *enb_ref = (enb_description_t *) elementP;

  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t *const) & enb_ref->ue_coll,
    s1ap_ue_compare_by_s11_sgw_teid_cb,
    parameterP,
    resultP);
  if (*resultP) {
    return true;
  }
  return false;
}

//------------------------------------------------------------------------------
ue_description_t *s1ap_is_ue_mme_id_in_list(
  const mme_ue_s1ap_id_t mme_ue_s1ap_id)
{
  ue_description_t *ue_ref = NULL;
  mme_ue_s1ap_id_t *mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t *) &mme_ue_s1ap_id;

  hashtable_ts_apply_callback_on_elements(
    &g_s1ap_enb_coll,
    s1ap_enb_find_ue_by_mme_ue_id_cb,
    (void *) mme_ue_s1ap_id_p,
    (void **) &ue_ref);
  OAILOG_TRACE(LOG_S1AP, "Return ue_ref %p \n", ue_ref);
  return ue_ref;
}

//------------------------------------------------------------------------------
// TODO(amar) unused function check with OAI.
ue_description_t *s1ap_is_s11_sgw_teid_in_list(const s11_teid_t teid)
{
  ue_description_t *ue_ref = NULL;
  s11_teid_t *teid_id_p = (s11_teid_t *) &teid;

  hashtable_ts_apply_callback_on_elements(
    &g_s1ap_enb_coll,
    s1ap_enb_find_ue_by_s11_sgw_teid_cb,
    (void *) teid_id_p,
    (void **) &ue_ref);
  return ue_ref;
}

//------------------------------------------------------------------------------
void s1ap_notified_new_ue_mme_s1ap_id_association(
  const sctp_assoc_id_t sctp_assoc_id,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id,
  const mme_ue_s1ap_id_t mme_ue_s1ap_id)
{
  enb_description_t *enb_ref = s1ap_is_enb_assoc_id_in_list(sctp_assoc_id);
  if (enb_ref) {
    ue_description_t *ue_ref =
      s1ap_is_ue_enb_id_in_list(enb_ref, enb_ue_s1ap_id);
    if (ue_ref) {
      ue_ref->mme_ue_s1ap_id = mme_ue_s1ap_id;
      hashtable_rc_t h_rc = hashtable_ts_insert(
        &g_s1ap_mme_id2assoc_id_coll,
        (const hash_key_t) mme_ue_s1ap_id,
        (void *) (uintptr_t) sctp_assoc_id);
      OAILOG_DEBUG(
        LOG_S1AP,
        "Associated  sctp_assoc_id %d, enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ":%s \n",
        sctp_assoc_id,
        enb_ue_s1ap_id,
        mme_ue_s1ap_id,
        hashtable_rc_code2string(h_rc));
      return;
    }
    OAILOG_DEBUG(
      LOG_S1AP,
      "Could not find  ue  with enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT "\n",
      enb_ue_s1ap_id);
    return;
  }
  OAILOG_DEBUG(
    LOG_S1AP, "Could not find  eNB with sctp_assoc_id %d \n", sctp_assoc_id);
}

//------------------------------------------------------------------------------
enb_description_t *s1ap_new_enb(void)
{
  enb_description_t *enb_ref = NULL;

  enb_ref = calloc(1, sizeof(enb_description_t));
  /*
   * Something bad happened during malloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify eNB with a cause like Hardware Failure.
   */
  DevAssert(enb_ref != NULL);
  // Update number of eNB associated
  nb_enb_associated++;
  bstring bs = bfromcstr("s1ap_ue_coll");
  hashtable_ts_init(
    &enb_ref->ue_coll, mme_config.max_ues, NULL, free_wrapper, bs);
  bdestroy_wrapper(&bs);
  enb_ref->nb_ue_associated = 0;
  enb_ref->s1ap_enb_assoc_clean_up_timer.sec = S1ap_TimeToWait_v20s;
  enb_ref->s1ap_enb_assoc_clean_up_timer.id = S1AP_TIMER_INACTIVE_ID;
  return enb_ref;
}

//------------------------------------------------------------------------------
ue_description_t *s1ap_new_ue(
  const sctp_assoc_id_t sctp_assoc_id,
  enb_ue_s1ap_id_t enb_ue_s1ap_id)
{
  enb_description_t *enb_ref = NULL;
  ue_description_t *ue_ref = NULL;

  enb_ref = s1ap_is_enb_assoc_id_in_list(sctp_assoc_id);
  DevAssert(enb_ref != NULL);
  ue_ref = calloc(1, sizeof(ue_description_t));
  /*
   * Something bad happened during malloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify eNB with a cause like Hardware Failure.
   */
  DevAssert(ue_ref != NULL);
  ue_ref->enb = enb_ref;
  ue_ref->enb_ue_s1ap_id = enb_ue_s1ap_id;

  hashtable_rc_t hashrc = hashtable_ts_insert(
    &enb_ref->ue_coll, (const hash_key_t) enb_ue_s1ap_id, (void *) ue_ref);
  if (HASH_TABLE_OK != hashrc) {
    OAILOG_ERROR(
      LOG_S1AP,
      "Could not insert UE descr in ue_coll: %s\n",
      hashtable_rc_code2string(hashrc));
    free_wrapper((void **) &ue_ref);
    return NULL;
  }
  MSC_LOG_EVENT(
    MSC_S1AP_MME,
    " Associating ue  (enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT ") to eNB %s",
    ue_ref->mme_ue_s1ap_id,
    enb_ref->enb_name);
  // Increment number of UE
  enb_ref->nb_ue_associated++;
  return ue_ref;
}

//------------------------------------------------------------------------------
void s1ap_remove_ue(ue_description_t *ue_ref)
{
  enb_description_t *enb_ref = NULL;

  /*
   * NULL reference...
   */
  if (ue_ref == NULL) return;

  mme_ue_s1ap_id_t mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id;
  enb_ref = ue_ref->enb;
  /*
   * Updating number of UE
   */
  DevAssert(enb_ref->nb_ue_associated > 0);
  enb_ref->nb_ue_associated--;

  /*
   * Remove any attached timer
   */
  // Stop UE Context Release Complete timer,if running
  if (ue_ref->s1ap_ue_context_rel_timer.id != S1AP_TIMER_INACTIVE_ID) {
    if (timer_remove(ue_ref->s1ap_ue_context_rel_timer.id, NULL)) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to stop s1ap ue context release complete timer for UE id  %d "
        "\n",
        ue_ref->mme_ue_s1ap_id);
    }
    ue_ref->s1ap_ue_context_rel_timer.id = S1AP_TIMER_INACTIVE_ID;
  }
  //     s1ap_timer_remove_ue(ue_ref->mme_ue_s1ap_id);
  OAILOG_TRACE(
    LOG_S1AP,
    "Removing UE enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT
    " mme_ue_s1ap_id:" MME_UE_S1AP_ID_FMT " in eNB id : %d\n",
    ue_ref->enb_ue_s1ap_id,
    ue_ref->mme_ue_s1ap_id,
    enb_ref->enb_id);

  ue_ref->s1_ue_state = S1AP_UE_INVALID_STATE;
  hashtable_ts_free(&enb_ref->ue_coll, ue_ref->enb_ue_s1ap_id);
  hashtable_ts_free(&g_s1ap_mme_id2assoc_id_coll, mme_ue_s1ap_id);
  if (!enb_ref->nb_ue_associated) {
    if (enb_ref->s1_state == S1AP_RESETING) {
      OAILOG_INFO(LOG_S1AP, "Moving eNB state to S1AP_INIT \n");
      enb_ref->s1_state = S1AP_INIT;
      update_mme_app_stats_connected_enb_sub();
    } else if (enb_ref->s1_state == S1AP_SHUTDOWN) {
      OAILOG_INFO(LOG_S1AP, "Deleting eNB \n");
      s1ap_remove_enb(enb_ref);
      update_mme_app_stats_connected_enb_sub();
    }
  }
}

//------------------------------------------------------------------------------
void s1ap_remove_enb(enb_description_t *enb_ref)
{
  if (enb_ref == NULL) {
    return;
  }
  // Stop associated UEs clean_up timer,if running
  if (enb_ref->s1ap_enb_assoc_clean_up_timer.id != S1AP_TIMER_INACTIVE_ID) {
    if (timer_remove(enb_ref->s1ap_enb_assoc_clean_up_timer.id, NULL)) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to stop wait_for_ue_cleanup timer for eNB association id  %u "
        "\n",
        enb_ref->sctp_assoc_id);
    } else {
      OAILOG_INFO(
        LOG_MME_APP,
        "Stopped wait_for_ue_cleanup timer for eNB association id  %u \n",
        enb_ref->sctp_assoc_id);
    }
    enb_ref->s1ap_enb_assoc_clean_up_timer.id = S1AP_TIMER_INACTIVE_ID;
  }
  enb_ref->s1_state = S1AP_INIT;
  hashtable_ts_destroy(&enb_ref->ue_coll);
  hashtable_ts_free(&g_s1ap_enb_coll, enb_ref->sctp_assoc_id);
  nb_enb_associated--;
}
