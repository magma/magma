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

/*! \file s1ap_mme.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "s1ap_mme.h"

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
#include "assertions.h"
#include "s1ap_mme_decoder.h"
#include "s1ap_mme_handlers.h"
#include "s1ap_mme_nas_procedures.h"
#include "s1ap_mme_itti_messaging.h"
#include "service303.h"
#include "service303_message_utils.h"
#include "dynamic_memory_check.h"
#include "mme_config.h"
#include "timer.h"
#include "itti_free_defined_msg.h"
#include "S1ap_TimeToWait.h"
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
  (LOG_S1AP, "[eNB]%*s" x "\n", 4 * indent, "", ##args)
#define UE_LIST_OUT(x, args...)                                                \
  OAILOG_DEBUG(LOG_S1AP, "[UE] %*s" x "\n", 4 * indent, "", ##args)
#else
#define eNB_LIST_OUT(x, args...)
#define UE_LIST_OUT(x, args...)
#endif

bool s1ap_dump_ue_hash_cb(
    const hash_key_t keyP, void* const ue_void, void* parameter,
    void** unused_res);
static void start_stats_timer(void);
static int handle_stats_timer(zloop_t* loop, int id, void* arg);
static long epc_stats_timer_id;

bool hss_associated = false;
static int indent   = 0;
task_zmq_ctx_t s1ap_task_zmq_ctx;

bool s1ap_congestion_control_enabled = true;
long s1ap_last_msg_latency           = 0;
long s1ap_zmq_th                     = LONG_MAX;

//------------------------------------------------------------------------------
static int s1ap_send_init_sctp(void) {
  // Create and alloc new message
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_S1AP, SCTP_INIT_MSG);
  message_p->ittiMsg.sctpInit.port         = S1AP_PORT_NUMBER;
  message_p->ittiMsg.sctpInit.ppid         = S1AP_SCTP_PPID;
  message_p->ittiMsg.sctpInit.ipv4         = 1;
  message_p->ittiMsg.sctpInit.ipv6         = 0;
  message_p->ittiMsg.sctpInit.nb_ipv4_addr = 1;
  message_p->ittiMsg.sctpInit.ipv4_address[0].s_addr =
      mme_config.ip.s1_mme_v4.s_addr;

  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * * * * Disable it for now.
   */
  message_p->ittiMsg.sctpInit.nb_ipv6_addr    = 0;
  message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_SCTP, message_p);
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  s1ap_state_t* state;
  MessageDef* received_message_p = receive_msg(reader);
  imsi64_t imsi64                = itti_get_associated_imsi(received_message_p);
  state                          = get_s1ap_state(false);
  AssertFatal(state != NULL, "failed to retrieve s1ap state (was null)");

  bool is_task_state_same = false;
  bool is_ue_state_same   = false;

  s1ap_last_msg_latency = ITTI_MSG_LATENCY(received_message_p);  // microseconds

  OAILOG_DEBUG(LOG_S1AP, "S1AP ZMQ latency: %ld.", s1ap_last_msg_latency);

  switch (ITTI_MSG_ID(received_message_p)) {
    case ACTIVATE_MESSAGE: {
      is_task_state_same = true;  // does not modify state
      is_ue_state_same   = true;
      hss_associated     = true;
    } break;

    case MESSAGE_TEST:
      is_task_state_same = true;  // does not modify state
      is_ue_state_same   = true;
      OAILOG_DEBUG(LOG_S1AP, "Received MESSAGE_TEST\n");
      break;

    case SCTP_DATA_IND: {
      /*
       * New message received from SCTP layer.
       * * * * Decode and handle it.
       */
      S1ap_S1AP_PDU_t pdu = {0};

      // Invoke S1AP message decoder
      if (s1ap_mme_decode_pdu(&pdu, SCTP_DATA_IND(received_message_p).payload) <
          0) {
        // TODO: Notify eNB of failure with right cause
        OAILOG_ERROR(LOG_S1AP, "Failed to decode new buffer\n");
      } else {
        s1ap_mme_handle_message(
            state, SCTP_DATA_IND(received_message_p).assoc_id,
            SCTP_DATA_IND(received_message_p).stream, &pdu);
      }

      // Free received PDU array
      ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
      bdestroy_wrapper(&SCTP_DATA_IND(received_message_p).payload);
    } break;

    case SCTP_DATA_CNF:
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_mme_itti_nas_downlink_cnf(
          SCTP_DATA_CNF(received_message_p).agw_ue_xap_id,
          SCTP_DATA_CNF(received_message_p).is_success);
      break;
    // SCTP layer notifies S1AP of disconnection of a peer.
    case SCTP_CLOSE_ASSOCIATION: {
      is_ue_state_same = true;
      s1ap_handle_sctp_disconnection(
          state, SCTP_CLOSE_ASSOCIATION(received_message_p).assoc_id,
          SCTP_CLOSE_ASSOCIATION(received_message_p).reset);
    } break;

    case SCTP_NEW_ASSOCIATION: {
      is_ue_state_same = true;
      increment_counter("mme_new_association", 1, NO_LABELS);
      if (s1ap_handle_new_association(
              state, &received_message_p->ittiMsg.sctp_new_peer)) {
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
          state, S1AP_NAS_DL_DATA_REQ(received_message_p).enb_ue_s1ap_id,
          S1AP_NAS_DL_DATA_REQ(received_message_p).mme_ue_s1ap_id,
          &S1AP_NAS_DL_DATA_REQ(received_message_p).nas_msg, imsi64,
          &is_task_state_same);
    } break;

    case S1AP_E_RAB_SETUP_REQ: {
      is_task_state_same = true;  // the following handler does not modify state
      s1ap_generate_s1ap_e_rab_setup_req(
          state, &S1AP_E_RAB_SETUP_REQ(received_message_p));
    } break;

    case S1AP_E_RAB_MODIFICATION_CNF: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_mme_generate_erab_modification_confirm(
          state, &received_message_p->ittiMsg.s1ap_e_rab_modification_cnf);
    } break;

    // From MME_APP task
    case S1AP_UE_CONTEXT_RELEASE_COMMAND: {
      is_ue_state_same = true;
      s1ap_handle_ue_context_release_command(
          state, &received_message_p->ittiMsg.s1ap_ue_context_release_command,
          imsi64);
    } break;

    case MME_APP_CONNECTION_ESTABLISHMENT_CNF: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_handle_conn_est_cnf(
          state, &MME_APP_CONNECTION_ESTABLISHMENT_CNF(received_message_p));
    } break;

    case MME_APP_S1AP_MME_UE_ID_NOTIFICATION: {
      s1ap_handle_mme_ue_id_notification(
          state, &MME_APP_S1AP_MME_UE_ID_NOTIFICATION(received_message_p));
    } break;

    case S1AP_ENB_INITIATED_RESET_ACK: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_handle_enb_initiated_reset_ack(
          &S1AP_ENB_INITIATED_RESET_ACK(received_message_p), imsi64);
    } break;

    case S1AP_PAGING_REQUEST: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      if (s1ap_handle_paging_request(
              state, &S1AP_PAGING_REQUEST(received_message_p), imsi64) !=
          RETURNok) {
        OAILOG_ERROR(LOG_S1AP, "Failed to send paging message\n");
      }
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_REQUEST: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_handle_ue_context_mod_req(
          state, &received_message_p->ittiMsg.s1ap_ue_context_mod_request,
          imsi64);
    } break;

    case S1AP_E_RAB_REL_CMD: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_generate_s1ap_e_rab_rel_cmd(
          state, &S1AP_E_RAB_REL_CMD(received_message_p));
    } break;

    case S1AP_PATH_SWITCH_REQUEST_ACK: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_handle_path_switch_req_ack(
          state, &received_message_p->ittiMsg.s1ap_path_switch_request_ack,
          imsi64);
    } break;

    case S1AP_PATH_SWITCH_REQUEST_FAILURE: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same   = true;
      s1ap_handle_path_switch_req_failure(
          state, &received_message_p->ittiMsg.s1ap_path_switch_request_failure,
          imsi64);
    } break;

    case MME_APP_HANDOVER_REQUEST: {
      s1ap_mme_handle_handover_request(
          state, &MME_APP_HANDOVER_REQUEST(received_message_p));
    } break;

    case MME_APP_HANDOVER_COMMAND: {
      s1ap_mme_handle_handover_command(
          state, &MME_APP_HANDOVER_COMMAND(received_message_p));
    } break;

    case TIMER_HAS_EXPIRED: {
      if (!timer_exists(
              received_message_p->ittiMsg.timer_has_expired.timer_id)) {
        break;
      }
      ue_description_t* ue_ref_p = NULL;
      if (received_message_p->ittiMsg.timer_has_expired.arg != NULL) {
        // check whether timer is related to eNB procedure or UE procedure
        s1ap_timer_arg_t timer_arg =
            *((s1ap_timer_arg_t*) (received_message_p->ittiMsg.timer_has_expired
                                       .arg));
        if (timer_arg.timer_class == S1AP_UE_TIMER) {
          mme_ue_s1ap_id_t mme_ue_s1ap_id = timer_arg.instance_id;
          if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
            OAILOG_WARNING_UE(
                LOG_S1AP, imsi64,
                "Timer expired but no assoicated UE context for UE id %d\n",
                mme_ue_s1ap_id);
            timer_handle_expired(
                received_message_p->ittiMsg.timer_has_expired.timer_id);
            is_ue_state_same = true;
            break;
          }
          if (received_message_p->ittiMsg.timer_has_expired.timer_id ==
              ue_ref_p->s1ap_ue_context_rel_timer.id) {
            // UE context release complete timer expiry handler
            OAILOG_WARNING_UE(
                LOG_S1AP, imsi64,
                "ue_context_release_command_timer_expired for UE id %d\n",
                mme_ue_s1ap_id);
            increment_counter(
                "ue_context_release_command_timer_expired", 1, NO_LABELS);
            s1ap_mme_handle_ue_context_rel_comp_timer_expiry(state, ue_ref_p);
          }
        } else {
          is_ue_state_same = true;
          OAILOG_WARNING_UE(
              LOG_S1AP, imsi64,
              "S1AP Timer expired with invalid timer class %u \n",
              timer_arg.timer_class);
        }
      }
      timer_handle_expired(
          received_message_p->ittiMsg.timer_has_expired.timer_id);

    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      s1ap_mme_exit();
    } break;

    default: {
      OAILOG_ERROR(
          LOG_S1AP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  if (!is_task_state_same) {
    put_s1ap_state();
  }
  if (!is_ue_state_same) {
    put_s1ap_imsi_map();
    put_s1ap_ue_state(imsi64);
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* s1ap_mme_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_S1AP);
  init_task_context(
      TASK_S1AP, (task_id_t[]){TASK_MME_APP, TASK_SCTP, TASK_SERVICE303}, 3,
      handle_message, &s1ap_task_zmq_ctx);

  if (s1ap_send_init_sctp() < 0) {
    OAILOG_ERROR(LOG_S1AP, "Error while sendind SCTP_INIT_MSG to SCTP \n");
  }
  start_stats_timer();

  zloop_start(s1ap_task_zmq_ctx.event_loop);
  s1ap_mme_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int s1ap_mme_init(const mme_config_t* mme_config_p) {
  OAILOG_DEBUG(LOG_S1AP, "Initializing S1AP interface\n");

  if (get_asn1c_environment_version() < ASN1_MINIMUM_VERSION) {
    OAILOG_ERROR(
        LOG_S1AP, "ASN1C version %d found, expecting at least %d\n",
        get_asn1c_environment_version(), ASN1_MINIMUM_VERSION);
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_S1AP, "ASN1C version %d\n", get_asn1c_environment_version());

  s1ap_congestion_control_enabled = mme_config_p->enable_congestion_control;
  s1ap_zmq_th                     = mme_config_p->s1ap_zmq_th;

  if (s1ap_state_init(
          mme_config_p->max_ues, mme_config_p->max_enbs,
          mme_config_p->use_stateless) < 0) {
    OAILOG_ERROR(LOG_S1AP, "Error while initing S1AP state\n");
    return RETURNerror;
  }

  if (itti_create_task(TASK_S1AP, &s1ap_mme_thread, NULL) == RETURNerror) {
    OAILOG_ERROR(LOG_S1AP, "Error while creating S1AP task\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_S1AP, "Initializing S1AP interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
void s1ap_mme_exit(void) {
  OAILOG_DEBUG(LOG_S1AP, "Cleaning S1AP\n");
  stop_timer(&s1ap_task_zmq_ctx, epc_stats_timer_id);

  put_s1ap_state();
  put_s1ap_imsi_map();

  s1ap_state_exit();

  destroy_task_context(&s1ap_task_zmq_ctx);

  OAILOG_DEBUG(LOG_S1AP, "Cleaning S1AP: DONE\n");
  OAI_FPRINTF_INFO("TASK_S1AP terminated\n");
  pthread_exit(NULL);
}

//------------------------------------------------------------------------------
void s1ap_dump_enb(const enb_description_t* const enb_ref) {
#ifdef S1AP_DEBUG_LIST
  // Reset indentation
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
  eNB_LIST_OUT("UEs attached to eNB: %d", enb_ref->nb_ue_associated);
  indent++;
  sctp_assoc_id_t sctp_assoc_id = enb_ref->sctp_assoc_id;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements(
      (hash_table_ts_t* const) state_ue_ht, s1ap_dump_ue_hash_cb,
      &sctp_assoc_id, NULL);
  indent--;
  eNB_LIST_OUT("");
#else
  s1ap_dump_ue(NULL);
#endif
}

//------------------------------------------------------------------------------
bool s1ap_dump_ue_hash_cb(
    __attribute__((unused)) const hash_key_t keyP, void* const ue_void,
    void* parameter, void __attribute__((unused)) * *unused_resultP) {
  ue_description_t* ue_ref       = (ue_description_t*) ue_void;
  sctp_assoc_id_t* sctp_assoc_id = (sctp_assoc_id_t*) parameter;
  if (ue_ref == NULL) {
    return false;
  }

  if (ue_ref->sctp_assoc_id == *sctp_assoc_id) {
    s1ap_dump_ue(ue_ref);
  }
  return false;
}

//------------------------------------------------------------------------------
void s1ap_dump_ue(const ue_description_t* const ue_ref) {
#ifdef S1AP_DEBUG_LIST

  if (ue_ref == NULL) return;

  UE_LIST_OUT("eNB UE s1ap id:   0x%06x", ue_ref->enb_ue_s1ap_id);
  UE_LIST_OUT("MME UE s1ap id:   0x%08x", ue_ref->mme_ue_s1ap_id);
  UE_LIST_OUT("SCTP stream recv: 0x%04x", ue_ref->sctp_stream_recv);
  UE_LIST_OUT("SCTP stream send: 0x%04x", ue_ref->sctp_stream_send);
#endif
}

//------------------------------------------------------------------------------
enb_description_t* s1ap_new_enb(s1ap_state_t* state) {
  enb_description_t* enb_ref = NULL;

  enb_ref = calloc(1, sizeof(enb_description_t));
  /*
   * Something bad happened during malloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify eNB with a cause like Hardware Failure.
   */
  DevAssert(enb_ref != NULL);
  bstring bs = bfromcstr("s1ap_ue_coll");
  hashtable_uint64_ts_init(&enb_ref->ue_id_coll, mme_config.max_ues, NULL, bs);
  bdestroy_wrapper(&bs);
  enb_ref->nb_ue_associated = 0;
  return enb_ref;
}

//------------------------------------------------------------------------------
ue_description_t* s1ap_new_ue(
    s1ap_state_t* state, const sctp_assoc_id_t sctp_assoc_id,
    enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  enb_description_t* enb_ref = NULL;
  ue_description_t* ue_ref   = NULL;

  enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
  DevAssert(enb_ref != NULL);
  ue_ref = calloc(1, sizeof(ue_description_t));
  /*
   * Something bad happened during malloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify eNB with a cause like Hardware Failure.
   */
  DevAssert(ue_ref != NULL);
  ue_ref->sctp_assoc_id  = sctp_assoc_id;
  ue_ref->enb_ue_s1ap_id = enb_ue_s1ap_id;
  ue_ref->comp_s1ap_id =
      S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id);

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_rc_t hashrc        = hashtable_ts_insert(
      state_ue_ht, (const hash_key_t) ue_ref->comp_s1ap_id, (void*) ue_ref);

  if (HASH_TABLE_OK != hashrc) {
    OAILOG_ERROR(
        LOG_S1AP, "Could not insert UE descr in ue_coll: %s\n",
        hashtable_rc_code2string(hashrc));
    free_wrapper((void**) &ue_ref);
    return NULL;
  }
  // Increment number of UE
  enb_ref->nb_ue_associated++;
  OAILOG_DEBUG(
      LOG_S1AP, "Num ue associated: %d on assoc id:%d",
      enb_ref->nb_ue_associated, sctp_assoc_id);
  return ue_ref;
}

//------------------------------------------------------------------------------
void s1ap_remove_ue(s1ap_state_t* state, ue_description_t* ue_ref) {
  enb_description_t* enb_ref = NULL;

  // NULL reference...
  if (ue_ref == NULL) return;

  mme_ue_s1ap_id_t mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id;
  enb_ref = s1ap_state_get_enb(state, ue_ref->sctp_assoc_id);
  DevAssert(enb_ref->nb_ue_associated > 0);
  // Updating number of UE
  enb_ref->nb_ue_associated--;

  // Stop UE Context Release Complete timer,if running
  if (ue_ref->s1ap_ue_context_rel_timer.id != S1AP_TIMER_INACTIVE_ID) {
    if (timer_remove(ue_ref->s1ap_ue_context_rel_timer.id, NULL)) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to stop s1ap ue context release complete timer, UE id: %d\n",
          ue_ref->mme_ue_s1ap_id);
    }
    ue_ref->s1ap_ue_context_rel_timer.id = S1AP_TIMER_INACTIVE_ID;
  }
  OAILOG_TRACE(
      LOG_S1AP,
      "Removing UE enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT
      " mme_ue_s1ap_id:" MME_UE_S1AP_ID_FMT " in eNB id : %d\n",
      ue_ref->enb_ue_s1ap_id, ue_ref->mme_ue_s1ap_id, enb_ref->enb_id);

  ue_ref->s1_ue_state = S1AP_UE_INVALID_STATE;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_free(state_ue_ht, ue_ref->comp_s1ap_id);
  hashtable_ts_free(&state->mmeid2associd, mme_ue_s1ap_id);
  hashtable_uint64_ts_remove(&enb_ref->ue_id_coll, mme_ue_s1ap_id);

  imsi64_t imsi64                = INVALID_IMSI64;
  s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      s1ap_imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);
  delete_s1ap_ue_state(imsi64);
  hashtable_uint64_ts_remove(
      s1ap_imsi_map->mme_ue_id_imsi_htbl, mme_ue_s1ap_id);

  OAILOG_DEBUG(
      LOG_S1AP, "Num UEs associated %u num ue_id_coll %zu",
      enb_ref->nb_ue_associated, enb_ref->ue_id_coll.num_elements);
  if (!enb_ref->nb_ue_associated) {
    if (enb_ref->s1_state == S1AP_RESETING) {
      OAILOG_INFO(LOG_S1AP, "Moving eNB state to S1AP_INIT \n");
      enb_ref->s1_state = S1AP_INIT;
      set_gauge("s1_connection", 0, 1, "enb_name", enb_ref->enb_name);
      state->num_enbs--;
    } else if (enb_ref->s1_state == S1AP_SHUTDOWN) {
      OAILOG_INFO(LOG_S1AP, "Deleting eNB \n");
      set_gauge("s1_connection", 0, 1, "enb_name", enb_ref->enb_name);
      s1ap_remove_enb(state, enb_ref);
    }
  }
}

//------------------------------------------------------------------------------
void s1ap_remove_enb(s1ap_state_t* state, enb_description_t* enb_ref) {
  if (enb_ref == NULL) {
    return;
  }
  enb_ref->s1_state = S1AP_INIT;
  hashtable_uint64_ts_destroy(&enb_ref->ue_id_coll);
  hashtable_ts_free(&state->enbs, enb_ref->sctp_assoc_id);
  state->num_enbs--;
}

static int handle_stats_timer(zloop_t* loop, int id, void* arg) {
  s1ap_state_t* s1ap_state_p = get_s1ap_state(false);
  application_s1ap_stats_msg_t stats_msg;
  stats_msg.nb_enb_connected = s1ap_state_p->num_enbs;
  return send_s1ap_stats_to_service303(
      &s1ap_task_zmq_ctx, TASK_S1AP, &stats_msg);
}

static void start_stats_timer(void) {
  epc_stats_timer_id = start_timer(
      &s1ap_task_zmq_ctx, EPC_STATS_TIMER_MSEC, TIMER_REPEAT_FOREVER,
      handle_stats_timer, NULL);
}
