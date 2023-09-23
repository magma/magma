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

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/message_utils/service303_message_utils.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/mme_init.hpp"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_decoder.hpp"
#include "S1ap_TimeToWait.h"
#include "asn_internal.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/mme_default_values.h"
#include "lte/gateway/c/core/oai/include/mme_app_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/mme_init.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/include/sctp_messages_types.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_itti_messaging.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_nas_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_timer.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

bool hss_associated = false;
namespace magma {
namespace lte {

static void start_stats_timer(void);
static int handle_stats_timer(zloop_t* loop, int id, void* arg);
static long epc_stats_timer_id;
static size_t epc_stats_timer_sec = 60;

task_zmq_ctx_t s1ap_task_zmq_ctx;

bool s1ap_congestion_control_enabled = true;
long s1ap_last_msg_latency = 0;
long s1ap_zmq_th = LONG_MAX;

static void s1ap_mme_exit(void);
//------------------------------------------------------------------------------
static int s1ap_send_init_sctp(void) {
  // Create and alloc new message
  MessageDef* message_p = NULL;

  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_S1AP, SCTP_INIT_MSG);
  message_p->ittiMsg.sctpInit.port = S1AP_PORT_NUMBER;
  message_p->ittiMsg.sctpInit.ppid = S1AP_SCTP_PPID;
  message_p->ittiMsg.sctpInit.ipv6 = mme_config.ip.s1_ipv6_enabled;

  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * Only allow multi homing when IPv6 is enabled.
   */
  message_p->ittiMsg.sctpInit.ipv4 = 1;
  message_p->ittiMsg.sctpInit.nb_ipv4_addr = 1;
  message_p->ittiMsg.sctpInit.ipv4_address[0].s_addr =
      mme_config.ip.s1_mme_v4.s_addr;

  if (message_p->ittiMsg.sctpInit.ipv6) {
    message_p->ittiMsg.sctpInit.nb_ipv6_addr = 1;
    message_p->ittiMsg.sctpInit.ipv6_address[0] = mme_config.ip.s1_mme_v6;
  } else {
    message_p->ittiMsg.sctpInit.nb_ipv6_addr = 0;
    message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  }

  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_SCTP, message_p);
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);
  imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
  oai::S1apState* state = get_s1ap_state(false);
  AssertFatal(state != NULL, "failed to retrieve s1ap state (was null)");

  bool is_task_state_same = false;
  bool is_ue_state_same = false;

  s1ap_last_msg_latency = ITTI_MSG_LATENCY(received_message_p);  // microseconds

  OAILOG_DEBUG(LOG_S1AP, "S1AP ZMQ latency: %ld.", s1ap_last_msg_latency);

  switch (ITTI_MSG_ID(received_message_p)) {
    case ACTIVATE_MESSAGE: {
      is_task_state_same = true;  // does not modify state
      is_ue_state_same = true;
      hss_associated = true;
    } break;

    case MESSAGE_TEST:
      is_task_state_same = true;  // does not modify state
      is_ue_state_same = true;
      OAILOG_DEBUG(LOG_S1AP, "Received MESSAGE_TEST\n");
      break;

    case SCTP_DATA_IND: {
      /*
       * New message received from SCTP layer.
       * * * * Decode and handle it.
       */
      S1ap_S1AP_PDU_t pdu = {S1ap_S1AP_PDU_PR_NOTHING, {0}};

      // Invoke S1AP message decoder
      if (s1ap_mme_decode_pdu(&pdu, SCTP_DATA_IND(received_message_p).payload) <
          0) {
        // TODO: Notify eNB of failure with right cause
        OAILOG_ERROR(LOG_S1AP, "Failed to decode new buffer\n");
      } else {
        s1ap_mme_handle_message(state,
                                SCTP_DATA_IND(received_message_p).assoc_id,
                                SCTP_DATA_IND(received_message_p).stream, &pdu);
      }

      // Free received PDU array
      ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
      bdestroy_wrapper(&SCTP_DATA_IND(received_message_p).payload);
    } break;

    case SCTP_DATA_CNF:
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
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
      is_ue_state_same = true;
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
      is_task_state_same =
          false;  // the following handler does not modify state
      is_ue_state_same = false;
      s1ap_handle_conn_est_cnf(
          state, &MME_APP_CONNECTION_ESTABLISHMENT_CNF(received_message_p),
          imsi64);
    } break;

    case MME_APP_S1AP_MME_UE_ID_NOTIFICATION: {
      s1ap_handle_mme_ue_id_notification(
          state, &MME_APP_S1AP_MME_UE_ID_NOTIFICATION(received_message_p));
    } break;

    case S1AP_ENB_INITIATED_RESET_ACK: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      s1ap_handle_enb_initiated_reset_ack(
          &S1AP_ENB_INITIATED_RESET_ACK(received_message_p), imsi64);
    } break;

    case S1AP_PAGING_REQUEST: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      if (s1ap_handle_paging_request(state,
                                     &S1AP_PAGING_REQUEST(received_message_p),
                                     imsi64) != RETURNok) {
        OAILOG_ERROR(LOG_S1AP, "Failed to send paging message\n");
      }
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_REQUEST: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      s1ap_handle_ue_context_mod_req(
          state, &received_message_p->ittiMsg.s1ap_ue_context_mod_request,
          imsi64);
    } break;

    case S1AP_E_RAB_REL_CMD: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      s1ap_generate_s1ap_e_rab_rel_cmd(state,
                                       &S1AP_E_RAB_REL_CMD(received_message_p));
    } break;

    case S1AP_PATH_SWITCH_REQUEST_ACK: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      s1ap_handle_path_switch_req_ack(
          state, &received_message_p->ittiMsg.s1ap_path_switch_request_ack,
          imsi64);
    } break;

    case S1AP_PATH_SWITCH_REQUEST_FAILURE: {
      is_task_state_same = true;  // the following handler does not modify state
      is_ue_state_same = true;
      s1ap_handle_path_switch_req_failure(
          &received_message_p->ittiMsg.s1ap_path_switch_request_failure,
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

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      s1ap_mme_exit();
    } break;

    default: {
      OAILOG_DEBUG(LOG_S1AP, "Unknown message ID %d:%s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
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
  const task_id_t peer_task_ids[] = {TASK_MME_APP, TASK_SCTP, TASK_SERVICE303};
  init_task_context(TASK_S1AP, peer_task_ids, 3, handle_message,
                    &s1ap_task_zmq_ctx);

  if (s1ap_send_init_sctp() < 0) {
    OAILOG_ERROR(LOG_S1AP, "Error while sendind SCTP_INIT_MSG to SCTP \n");
  }
  start_stats_timer();

  zloop_start(s1ap_task_zmq_ctx.event_loop);
  AssertFatal(0,
              "Asserting as s1ap_mme_thread should not be exiting on its own! "
              "This is likely due to a timer handler function returning -1 "
              "(RETURNerror) on one of the conditions.");
  return NULL;
}

//------------------------------------------------------------------------------
// TODO(rsarwad): remove extern C, when complete functionality of MME is
// migrated to cpp
extern "C" status_code_e s1ap_mme_init(const mme_config_t* mme_config_p) {
  OAILOG_DEBUG(LOG_S1AP, "Initializing S1AP interface\n");
  OAILOG_DEBUG(LOG_S1AP, "ASN1C version %d\n", get_asn1c_environment_version());

  s1ap_congestion_control_enabled = mme_config_p->enable_congestion_control;
  s1ap_zmq_th = mme_config_p->s1ap_zmq_th;

  // Initialize global stats timer
  epc_stats_timer_sec = (size_t)mme_config_p->stats_timer_sec;

  if (s1ap_state_init(mme_config_p->use_stateless) < 0) {
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
static void s1ap_mme_exit(void) {
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
void s1ap_new_enb(oai::EnbDescription* enb_ref) {
  magma::proto_map_uint32_uint64_t ue_id_coll;

  if (enb_ref == nullptr) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Received invalid pointer for structure, oai::EnbDescription ");
    return;
  }
  ue_id_coll.map = enb_ref->mutable_ue_id_map();
  ue_id_coll.set_name("s1ap_ue_coll");
  enb_ref->set_nb_ue_associated(0);
  return;
}

//------------------------------------------------------------------------------
oai::UeDescription* s1ap_new_ue(oai::EnbDescription* enb_ref,
                                const sctp_assoc_id_t sctp_assoc_id,
                                enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  oai::UeDescription* ue_ref = nullptr;

  if (!enb_ref) {
    OAILOG_ERROR(LOG_S1AP, "Invalid enb context for assoc_id: %u",
                 sctp_assoc_id);
    return nullptr;
  }

  ue_ref = new oai::UeDescription();
  /*
   * Something bad happened during memory allocation...
   * * * * May be we are running out of memory.
   * * * * TODO: Notify eNB with a cause like Hardware Failure.
   */
  if (ue_ref == nullptr) {
    OAILOG_ERROR(LOG_S1AP,
                 "Failed to allocate memory for protobuf object UeDescription");
    return nullptr;
  }
  ue_ref->set_sctp_assoc_id(sctp_assoc_id);
  ue_ref->set_enb_ue_s1ap_id(enb_ue_s1ap_id);
  ue_ref->set_comp_s1ap_id(
      S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id));

  map_uint64_ue_description_t* s1ap_ue_state = get_s1ap_ue_state();
  if (s1ap_ue_state == nullptr) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    return nullptr;
  }
  magma::proto_map_rc_t rc =
      s1ap_ue_state->insert(ue_ref->comp_s1ap_id(), ue_ref);

  if (rc != magma::PROTO_MAP_OK) {
    OAILOG_ERROR(LOG_S1AP, "Could not insert UE descr in ue_coll: %s\n",
                 magma::map_rc_code2string(rc));
    free_cpp_wrapper(reinterpret_cast<void**>(&ue_ref));
    return nullptr;
  }
  // Increment number of UE
  enb_ref->set_nb_ue_associated((enb_ref->nb_ue_associated() + 1));
  OAILOG_DEBUG(LOG_S1AP, "Num ue associated: %d on assoc id:%d",
               enb_ref->nb_ue_associated(), sctp_assoc_id);
  return ue_ref;
}

//------------------------------------------------------------------------------
void s1ap_remove_ue(oai::S1apState* state, oai::UeDescription* ue_ref) {
  oai::EnbDescription enb_ref;
  // NULL reference...
  if (ue_ref == nullptr) return;

  mme_ue_s1ap_id_t mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id();
  if ((s1ap_state_get_enb(state, ue_ref->sctp_assoc_id(), &enb_ref)) !=
      PROTO_MAP_OK) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get enb association for assoc_id: %u",
                 ue_ref->sctp_assoc_id());
    return;
  }
  DevAssert(enb_ref.nb_ue_associated() > 0);
  // Updating number of UE
  enb_ref.set_nb_ue_associated((enb_ref.nb_ue_associated() - 1));
  OAILOG_TRACE(LOG_S1AP,
               "Removing UE enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT
               " mme_ue_s1ap_id:" MME_UE_S1AP_ID_FMT " in eNB id : %d\n",
               ue_ref->enb_ue_s1ap_id(), ue_ref->mme_ue_s1ap_id(),
               enb_ref.enb_id);

  ue_ref->set_s1ap_ue_state(oai::S1AP_UE_INVALID_STATE);
  if (ue_ref->s1ap_ue_context_rel_timer().id() != S1AP_TIMER_INACTIVE_ID) {
    s1ap_stop_timer(ue_ref->s1ap_ue_context_rel_timer().id());
    ue_ref->mutable_s1ap_ue_context_rel_timer()->set_id(S1AP_TIMER_INACTIVE_ID);
  }

  map_uint64_ue_description_t* s1ap_ue_state = get_s1ap_ue_state();
  if (s1ap_ue_state == nullptr) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    OAILOG_FUNC_OUT(LOG_S1AP);
  }
  s1ap_ue_state->remove(ue_ref->comp_s1ap_id());
  proto_map_uint32_uint32_t mmeid2associd_map;
  mmeid2associd_map.map = state->mutable_mmeid2associd();
  mmeid2associd_map.remove(mme_ue_s1ap_id);

  magma::proto_map_uint32_uint64_t ue_id_coll;
  ue_id_coll.map = enb_ref.mutable_ue_id_map();
  ue_id_coll.remove(mme_ue_s1ap_id);

  imsi64_t imsi64 = INVALID_IMSI64;
  magma::proto_map_uint32_uint64_t ueid_imsi_map;
  get_s1ap_ueid_imsi_map(&ueid_imsi_map);
  ueid_imsi_map.get(mme_ue_s1ap_id, &imsi64);
  delete_s1ap_ue_state(imsi64);
  ueid_imsi_map.remove(mme_ue_s1ap_id);

  OAILOG_DEBUG(LOG_S1AP, "Num UEs associated %u num elements in ue_id_coll %lu",
               enb_ref.nb_ue_associated(), ue_id_coll.size());
  if (!enb_ref.nb_ue_associated()) {
    if (enb_ref.s1_enb_state() == oai::S1AP_RESETING) {
      OAILOG_INFO(LOG_S1AP, "Moving eNB state to S1AP_INIT \n");
      enb_ref.set_s1_state(oai::S1AP_INIT);
      set_gauge("s1_connection", 0, 1, "enb_name", enb_ref.enb_name().c_str());
      state->set_num_enbs(state->num_enbs() - 1);
    } else if (enb_ref.s1_enb_state() == oai::S1AP_SHUTDOWN) {
      OAILOG_INFO(LOG_S1AP, "Deleting eNB \n");
      set_gauge("s1_connection", 0, 1, "enb_name", enb_ref.enb_name().c_str());
      s1ap_remove_enb(state, &enb_ref);
    }
  }
  s1ap_state_update_enb_map(state, enb_ref.sctp_assoc_id(), &enb_ref);
}

//------------------------------------------------------------------------------
void s1ap_remove_enb(oai::S1apState* state, oai::EnbDescription* enb_ref) {
  if (enb_ref == nullptr) {
    return;
  }
  magma::proto_map_uint32_uint64_t ue_id_coll;
  proto_map_uint32_enb_description_t enb_map;
  enb_ref->set_s1_state(oai::S1AP_INIT);

  ue_id_coll.map = enb_ref->mutable_ue_id_map();
  ue_id_coll.clear();
  OAILOG_INFO(LOG_S1AP, "Deleting eNB on assoc_id :%u\n",
              enb_ref->sctp_assoc_id());
  enb_map.map = state->mutable_enbs();
  enb_map.remove(enb_ref->sctp_assoc_id());
  state->set_num_enbs(state->num_enbs() - 1);
}

static int handle_stats_timer(zloop_t* loop, int id, void* arg) {
  oai::S1apState* s1ap_state_p = get_s1ap_state(false);
  application_s1ap_stats_msg_t stats_msg;
  stats_msg.nb_enb_connected = s1ap_state_p->num_enbs();
  stats_msg.nb_s1ap_last_msg_latency = s1ap_last_msg_latency;
  return send_s1ap_stats_to_service303(&s1ap_task_zmq_ctx, TASK_S1AP,
                                       &stats_msg);
}

static void start_stats_timer(void) {
  epc_stats_timer_id =
      start_timer(&s1ap_task_zmq_ctx, 1000 * epc_stats_timer_sec,
                  TIMER_REPEAT_FOREVER, handle_stats_timer, NULL);
}

}  // namespace lte
}  // namespace magma
