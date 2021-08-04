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
/****************************************************************************
  Source      ngap_amf.c
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>
#include <netinet/in.h>

#include "log.h"
#include "bstrlib.h"
#include "hashtable.h"
#include "assertions.h"
#include "ngap_state.h"
#include "service303.h"
#include "dynamic_memory_check.h"
#include "amf_config.h"
#include "mme_config.h"
#include "amf_default_values.h"
#include "timer.h"
#include "itti_free_defined_msg.h"
#include "Ngap_TimeToWait.h"

#include "asn_internal.h"
#include "sctp_messages_types.h"

#include "common_defs.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "amf_app_messages_types.h"
#include "amf_default_values.h"

#include "ngap_messages_types.h"
#include "timer_messages_types.h"
#include "ngap_amf.h"

amf_config_t amf_config;
task_zmq_ctx_t ngap_task_zmq_ctx;

static int ngap_send_init_sctp(void) {
  // Create and alloc new message
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NGAP, SCTP_INIT_MSG);
  message_p->ittiMsg.sctpInit.port         = NGAP_PORT_NUMBER;
  message_p->ittiMsg.sctpInit.ppid         = NGAP_SCTP_PPID;
  message_p->ittiMsg.sctpInit.ipv4         = 1;
  message_p->ittiMsg.sctpInit.ipv6         = 0;
  message_p->ittiMsg.sctpInit.nb_ipv4_addr = 1;
  message_p->ittiMsg.sctpInit.ipv4_address[0].s_addr =
      mme_config.ip.s1_mme_v4.s_addr;  // TODO Need change

  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * * * * Disable it for now.
   */
  message_p->ittiMsg.sctpInit.nb_ipv6_addr    = 0;
  message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_SCTP, message_p);
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);

  switch (ITTI_MSG_ID(received_message_p)) {
    case SCTP_DATA_IND: {
      /*
       * New message received from SCTP layer.
       * * * * Decode and handle it.
       */

      // Invoke NGAP message decoder
      // TODO
    } break;

    case SCTP_NEW_ASSOCIATION: {
      // TODO
    } break;

    case NGAP_NAS_DL_DATA_REQ: {  // packets from NAS
                                  /*
                                   * New message received from NAS task.
                                   */
                                  // TODO
    } break;

    default: {
      OAILOG_ERROR(
          LOG_NGAP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  // put_ngap_state(); //TODO : implement
  put_ngap_imsi_map();
  // put_ngap_ue_state(imsi64); //TODO
  itti_free_msg_content(received_message_p);
  zframe_destroy(&msg_frame);
  return 0;
}

//------------------------------------------------------------------------------
static void* ngap_amf_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_NGAP);
  init_task_context(
      TASK_NGAP, (task_id_t[]){TASK_AMF_APP, TASK_SCTP}, 2, handle_message,
      &ngap_task_zmq_ctx);

  if (ngap_send_init_sctp() < 0) {
    OAILOG_ERROR(LOG_NGAP, "Error while sending SCTP_INIT_MSG to SCTP \n");
  } else {
    OAILOG_INFO(LOG_NGAP, " sending SCTP_INIT_MSG to SCTP \n");
  }
  zloop_start(ngap_task_zmq_ctx.event_loop);
  AssertFatal(
      0, "Asserting as ngap_amf_thread should not be exiting on its own!");
  return NULL;
}

//------------------------------------------------------------------------------
int ngap_amf_init(const amf_config_t* amf_config_p) {
  OAILOG_DEBUG(LOG_NGAP, "Initializing NGAP interface\n");

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

  // put_ngap_imsi_map();//TODO

  // n1ap_state_exit(); //TODO

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
  hashtable_uint64_ts_init(
      &gnb_ref->ue_id_coll, amf_config.max_ues, NULL, bs);  // Need change
  bdestroy_wrapper(&bs);
  gnb_ref->nb_ue_associated = 0;
  return gnb_ref;
}

//------------------------------------------------------------------------------
m5g_ue_description_t* ngap_new_ue(
    ngap_state_t* state, const sctp_assoc_id_t sctp_assoc_id,
    gnb_ue_ngap_id_t gnb_ue_ngap_id) {
  gnb_description_t* gnb_ref   = NULL;
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
  ue_ref->sctp_assoc_id  = sctp_assoc_id;
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;
  ue_ref->comp_ngap_id   = ngap_get_comp_ngap_id(sctp_assoc_id, gnb_ue_ngap_id);

  hash_table_ts_t* state_ue_ht = get_ngap_ue_state();
  hashtable_rc_t hashrc        = hashtable_ts_insert(
      state_ue_ht, (const hash_key_t) ue_ref->comp_ngap_id, (void*) ue_ref);

  if (HASH_TABLE_OK != hashrc) {
    OAILOG_ERROR(
        LOG_NGAP, "Could not insert UE descr in ue_coll: %s\n",
        hashtable_rc_code2string(hashrc));
    free_wrapper((void**) &ue_ref);
    return NULL;
  }
  // Increment number of UE
  gnb_ref->nb_ue_associated++;
  return ue_ref;
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
