/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf.h"

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>
#include <netinet/in.h>

#include "Ngap_TimeToWait.h"
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/amf_default_values.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_decoder.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_itti_messaging.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_nas_procedures.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

#include "asn_internal.h"
#include "lte/gateway/c/core/oai/include/sctp_messages_types.hpp"

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/amf_default_values.h"
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"

#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"

task_zmq_ctx_t ngap_task_zmq_ctx;

uint64_t ngap_last_msg_latency = 0;

static int ngap_send_init_sctp(void) {
  // Create and alloc new message
  MessageDef* message_p = NULL;

  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_NGAP, SCTP_INIT_MSG);
  message_p->ittiMsg.sctpInit.port = NGAP_PORT_NUMBER;
  message_p->ittiMsg.sctpInit.ppid = NGAP_SCTP_PPID;
  message_p->ittiMsg.sctpInit.ipv4 = 1;
  message_p->ittiMsg.sctpInit.ipv6 = 0;
  message_p->ittiMsg.sctpInit.nb_ipv4_addr = 1;
  message_p->ittiMsg.sctpInit.ipv4_address[0].s_addr =
      mme_config.ip.s1_mme_v4.s_addr;

  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * * * * Disable it for now.
   */
  message_p->ittiMsg.sctpInit.nb_ipv6_addr = 0;
  message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  return ngap_send_msg_to_task(&ngap_task_zmq_ctx, TASK_SCTP, message_p);
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  ngap_state_t* state = NULL;
  MessageDef* received_message_p = receive_msg(reader);

  imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
  state = get_ngap_state(false);
  AssertFatal(state != NULL, "failed to retrieve ngap state (was null)");

  OAILOG_INFO(LOG_NGAP, "Received msg from :[%s] id:[%d] name:[%s]\n",
              ITTI_MSG_ORIGIN_NAME(received_message_p),
              ITTI_MSG_ID(received_message_p),
              ITTI_MSG_NAME(received_message_p));

  bool is_task_state_same = false;
  bool is_ue_state_same = false;

  ngap_last_msg_latency = ITTI_MSG_LATENCY(received_message_p);  // microseconds

  OAILOG_DEBUG(LOG_NGAP, "NGAP ZMQ latency: %ld.", ngap_last_msg_latency);

  switch (ITTI_MSG_ID(received_message_p)) {
    case SCTP_DATA_IND: {
      /*
       * New message received from SCTP layer.
       * * * * Decode and handle it.
       */

      // Invoke NGAP message decoder
      Ngap_NGAP_PDU_t pdu = {0};

      if (ngap_amf_decode_pdu(&pdu,
                              SCTP_DATA_IND(received_message_p).payload)) {
        // TODO: Notify gNB of failure with right cause
        OAILOG_ERROR(LOG_NGAP, "Failed to decode new buffer\n");

      } else {
        ngap_amf_handle_message(state,
                                SCTP_DATA_IND(received_message_p).assoc_id,
                                SCTP_DATA_IND(received_message_p).stream, &pdu);
      }

      ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &pdu);

      // Free received PDU array
      bdestroy_wrapper(&SCTP_DATA_IND(received_message_p).payload);

    } break;

    // SCTP layer notifies NGAP of disconnection of a peer.
    case SCTP_CLOSE_ASSOCIATION: {
      is_ue_state_same = true;
      ngap_handle_sctp_disconnection(
          state, SCTP_CLOSE_ASSOCIATION(received_message_p).assoc_id,
          SCTP_CLOSE_ASSOCIATION(received_message_p).reset);
    } break;

    case SCTP_NEW_ASSOCIATION: {
      is_ue_state_same = true;
      increment_counter("amf_new_association", 1, NO_LABELS);
      if (ngap_handle_new_association(
              state, &received_message_p->ittiMsg.sctp_new_peer)) {
        increment_counter("amf_new_association", 1, 1, "result", "failure");
      } else {
        increment_counter("amf_new_association", 1, 1, "result", "success");
      }
    } break;

    case NGAP_NAS_DL_DATA_REQ: {  // packets from NAS

      /*
       * New message received from NAS task.
       * * * * This corresponds to a NGAP downlink nas transport message.
       */
      ngap_generate_downlink_nas_transport(
          state, NGAP_NAS_DL_DATA_REQ(received_message_p).gnb_ue_ngap_id,
          NGAP_NAS_DL_DATA_REQ(received_message_p).amf_ue_ngap_id,
          &NGAP_NAS_DL_DATA_REQ(received_message_p).nas_msg, imsi64);
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      ngap_amf_exit();
    } break;

    case AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION: {
      ngap_handle_amf_ue_id_notification(
          state, &AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION(received_message_p));
    } break;

    // From AMF_APP task
    case NGAP_UE_CONTEXT_RELEASE_COMMAND: {
      is_ue_state_same = true;
      ngap_handle_ue_context_release_command(
          state, &received_message_p->ittiMsg.ngap_ue_context_release_command,
          imsi64);
    } break;

    case NGAP_INITIAL_CONTEXT_SETUP_REQ: {
      ngap_handle_conn_est_cnf(
          state, &NGAP_INITIAL_CONTEXT_SETUP_REQ(received_message_p));
    } break;

    case NGAP_PDUSESSION_RESOURCE_SETUP_REQ: {
      ngap_generate_ngap_pdusession_resource_setup_req(
          state, &NGAP_PDUSESSION_RESOURCE_SETUP_REQ(received_message_p));
    } break;

    case NGAP_PDU_SESSION_RESOURCE_MODIFY_REQ: {
      ngap_generate_ngap_pdusession_resource_modify_req(
          state, &NGAP_PDU_SESSION_RESOURCE_MODIFY_REQ(received_message_p));
    } break;

    case NGAP_PDUSESSIONRESOURCE_REL_REQ: {
      ngap_generate_ngap_pdusession_resource_rel_cmd(
          state, &NGAP_PDUSESSIONRESOURCE_REL_REQ(received_message_p));
    } break;

    case NGAP_GNB_INITIATED_RESET_ACK: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      ngap_handle_gnb_initiated_reset_ack(
          &NGAP_GNB_INITIATED_RESET_ACK(received_message_p));
    } break;

    case NGAP_PAGING_REQUEST: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      if (ngap_handle_paging_request(state,
                                     &NGAP_PAGING_REQUEST(received_message_p),
                                     imsi64) != RETURNok) {
        OAILOG_ERROR(LOG_NGAP, "Failed to send paging message\n");
      }
    } break;

    default: {
      OAILOG_ERROR(LOG_NGAP, "Unknown message ID %d:%s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
    } break;
  }

  if (!is_task_state_same) {
    put_ngap_state();
  }
  if (!is_ue_state_same) {
    put_ngap_imsi_map();
    put_ngap_ue_state(imsi64);
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* ngap_amf_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_NGAP);
  init_task_context(TASK_NGAP, (task_id_t[]){TASK_AMF_APP, TASK_SCTP}, 2,
                    handle_message, &ngap_task_zmq_ctx);

  if (ngap_send_init_sctp() < 0) {
    OAILOG_ERROR(LOG_NGAP, "Error while sending SCTP_INIT_MSG to SCTP \n");
  } else {
    OAILOG_DEBUG(LOG_NGAP, " Sending SCTP_INIT_MSG to SCTP \n");
  }
  zloop_start(ngap_task_zmq_ctx.event_loop);
  AssertFatal(0,
              "Asserting as ngap_amg_thread should not be exiting on its own! "
              "This is likely due to a timer handler function returning -1 "
              "(RETURNerror) on one of the conditions.");
  return NULL;
}

//------------------------------------------------------------------------------
status_code_e ngap_amf_init(const amf_config_t* amf_config_p) {
  OAILOG_DEBUG(LOG_NGAP, "Initializing NGAP interface\n");

  if (NULL == amf_config_p) {
    return RETURNerror;
  }

  if (ngap_state_init(amf_config_p->max_ues, amf_config_p->max_gnbs,
                      amf_config_p->use_stateless) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Error while initing NGAP state\n");
    return RETURNerror;
  }

  if (itti_create_task(TASK_NGAP, &ngap_amf_thread, NULL) == RETURNerror) {
    OAILOG_ERROR(LOG_NGAP, "Error while creating NGAP task\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_NGAP, "Initializing NGAP interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
void ngap_amf_exit(void) {
  OAILOG_DEBUG(LOG_NGAP, "Cleaning NGAP\n");

  destroy_task_context(&ngap_task_zmq_ctx);

  put_ngap_imsi_map();
  ngap_state_exit();

  OAILOG_DEBUG(LOG_NGAP, "Cleaning NGAP: DONE\n");
  OAI_FPRINTF_INFO("TASK_NGAP terminated\n");
  pthread_exit(NULL);
}

//------------------------------------------------------------------------------
gnb_description_t* ngap_new_gnb(ngap_state_t* state) {
  gnb_description_t* gnb_ref = NULL;

  gnb_ref = calloc(1, sizeof(gnb_description_t));
  /*
   * Something bad happened during calloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify gNB with a cause like Hardware Failure.
   */
  DevAssert(gnb_ref != NULL);
  // Update number of gNB associated
  state->num_gnbs++;
  bstring bs = bfromcstr("ngap_ue_coll");
  // Need change in below line in future after amf_config
  hashtable_uint64_ts_init(&gnb_ref->ue_id_coll, amf_config.max_ues, NULL,
                           bs);  // Need change
  bdestroy_wrapper(&bs);
  gnb_ref->nb_ue_associated = 0;
  return gnb_ref;
}

//------------------------------------------------------------------------------
m5g_ue_description_t* ngap_new_ue(ngap_state_t* state,
                                  const sctp_assoc_id_t sctp_assoc_id,
                                  gnb_ue_ngap_id_t gnb_ue_ngap_id) {
  gnb_description_t* gnb_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  gnb_ref = ngap_state_get_gnb(state, sctp_assoc_id);
  DevAssert(gnb_ref != NULL);
  ue_ref = calloc(1, sizeof(m5g_ue_description_t));

  /*
   * Something bad happened during malloc...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify gNB with a cause like Hardware Failure.
   */

  DevAssert(ue_ref != NULL);
  ue_ref->sctp_assoc_id = sctp_assoc_id;
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;
  ue_ref->comp_ngap_id = ngap_get_comp_ngap_id(sctp_assoc_id, gnb_ue_ngap_id);

  hash_table_ts_t* state_ue_ht = get_ngap_ue_state();
  hashtable_rc_t hashrc = hashtable_ts_insert(
      state_ue_ht, (const hash_key_t)ue_ref->comp_ngap_id, (void*)ue_ref);

  if (HASH_TABLE_OK != hashrc) {
    OAILOG_ERROR(LOG_NGAP, "Could not insert UE descr in ue_coll: %s\n",
                 hashtable_rc_code2string(hashrc));
    free_wrapper((void**)&ue_ref);
    return NULL;
  }
  // Increment number of UE
  gnb_ref->nb_ue_associated++;
  return ue_ref;
}

//------------------------------------------------------------------------------
void ngap_remove_ue(ngap_state_t* state, m5g_ue_description_t* ue_ref) {
  gnb_description_t* gNB_ref = NULL;

  // NULL reference...
  if (ue_ref == NULL) return;

  amf_ue_ngap_id_t amf_ue_ngap_id = ue_ref->amf_ue_ngap_id;
  gNB_ref = ngap_state_get_gnb(state, ue_ref->sctp_assoc_id);

  // Updating number of UE
  gNB_ref->nb_ue_associated--;

  ue_ref->ng_ue_state = NGAP_UE_INVALID_STATE;

  hash_table_ts_t* state_ue_ht = get_ngap_ue_state();
  hashtable_ts_free(state_ue_ht, ue_ref->comp_ngap_id);
  hashtable_ts_free(&state->amfid2associd, amf_ue_ngap_id);
  hashtable_uint64_ts_remove(&gNB_ref->ue_id_coll, amf_ue_ngap_id);

  imsi64_t imsi64 = INVALID_IMSI64;
  ngap_imsi_map_t* ngap_imsi_map = get_ngap_imsi_map();

  hashtable_uint64_ts_get(ngap_imsi_map->amf_ue_id_imsi_htbl,
                          (const hash_key_t)amf_ue_ngap_id, &imsi64);

  delete_ngap_ue_state(imsi64);

  if (!gNB_ref->nb_ue_associated) {
    if (gNB_ref->ng_state == NGAP_RESETING) {
      gNB_ref->ng_state = NGAP_INIT;
      set_gauge("ngap_connection", 0, 1, "gnb_name", gNB_ref->gnb_name);
      // update_mme_app_stats_connected_enb_sub();
    } else if (gNB_ref->ng_state == NGAP_SHUTDOWN) {
      OAILOG_INFO(LOG_NGAP, "Deleting gNB \n");
      // set_gauge("s1_connection", 0, 1, "enb_name", enb_ref->enb_name);
      ngap_remove_gnb(state, gNB_ref);
    }
  }
}

//------------------------------------------------------------------------------
void ngap_remove_gnb(ngap_state_t* state, gnb_description_t* gnb_ref) {
  if (gnb_ref == NULL) {
    return;
  }
  gnb_ref->ng_state = NGAP_INIT;
  hashtable_uint64_ts_destroy(&gnb_ref->ue_id_coll);
  hashtable_ts_free(&state->gnbs, gnb_ref->sctp_assoc_id);
  state->num_gnbs--;
}
