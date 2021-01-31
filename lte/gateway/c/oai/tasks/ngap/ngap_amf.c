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

/*! \file ngap_amf.c
 */

#if HAVE_CONFIG_H
#include "config.h"
#endif

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <unistd.h>
#include <pthread.h>
#include <netinet/in.h>

#include "log.h"
#include "bstrlib.h"
#include "hashtable.h"
#include "assertions.h"
#include "ngap_amf_decoder.h"
#include "ngap_amf_handlers.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_itti_messaging.h"

#include "service303.h"
#include "dynamic_memory_check.h"
#include "amf_config.h"
#include "mme_config.h"
#include "amf_default_values.h"
//#include "ngap_timer.h"
#include "timer.h"
//#include "ngap_itti_free_defined_msg.h"
#include "itti_free_defined_msg.h"
#include "Ngap_TimeToWait.h"
//#include "TimeToWait.h"

#include "asn_internal.h"
#include "sctp_messages_types.h"
//#include "sctp_messages_def.h"

#include "common_defs.h"
//#include "ngap_intertask_interface.h"
#include "intertask_interface.h"
//#include "ngap_intertask_interface_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "amf_app_messages_types.h"
#include "amf_default_values.h"

#include "ngap_messages_types.h"
#include "timer_messages_types.h"
#include "ngap_amf.h"

#if NGAP_DEBUG_LIST
#define gNB_LIST_OUT(x, args...)                                               \
  (LOG_NGAP, "[gNB]%*s" x "\n", 4 * indent, "", ##args)
#define UE_LIST_OUT(x, args...)                                                \
  OAILOG_DEBUG(LOG_NGAP, "[UE] %*s" x "\n", 4 * indent, "", ##args)
#else
#define gNB_LIST_OUT(x, args...)
#define UE_LIST_OUT(x, args...)
#endif

/*
bool ngap_dump_enb_hash_cb(
  const hash_key_t keyP,
  void* const enb_void,
  void* unused_param,
  void** unused_res);
bool ngap_dump_ue_hash_cb(
  const hash_key_t keyP,
  void* const ue_void,
  void* parameter,
  void** unused_res);
*/

// extern int ngap_amf_decode_pdu(ngap_message* message,  const bstring const
// raw,  MessagesIds* messages_id) __attribute__((warn_unus    ed_result));

amf_config_t amf_config;
// bool hss_associated = false;
// static int indent = 0;
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
      mme_config.ip.s1_mme_v4.s_addr;  // Need change in future

  /*
   * SR WARNING: ipv6 multi-homing fails sometimes for localhost.
   * * * * Disable it for now.
   */
  message_p->ittiMsg.sctpInit.nb_ipv6_addr    = 0;
  message_p->ittiMsg.sctpInit.ipv6_address[0] = in6addr_loopback;
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_SCTP, message_p);
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  ngap_state_t* state = NULL;
  // MessagesIds message_id = MESSAGES_ID_MAX;

  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);

  imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
  OAILOG_INFO(
      LOG_NGAP,
      "NGAP-TEST : inside handler received message with imsi %lu  and "
      "message "
      "type %d \n",
      imsi64, ITTI_MSG_ID(received_message_p));

  state = get_ngap_state(false);  // TODO Need change
  AssertFatal(
      state != NULL,
      "failed to retrieve ngap state (was null)");  // Need change

  switch (ITTI_MSG_ID(received_message_p)) {
#if 0 
   case ACTIVATE_MESSAGE: {
      hss_associated = true;
    } break;

    case MESSAGE_TEST:
      OAILOG_DEBUG(LOG_NGAP, "Received MESSAGE_TEST\n");
      break;
#endif
    case SCTP_DATA_IND: {
      /*
       * New message received from SCTP layer.
       * * * * Decode and handle it.
       */

      // Invoke NGAP message decoder
#if 1 /*decoder from lib*/
      Ngap_NGAP_PDU_t pdu = {0};

      if (ngap_amf_decode_pdu(
              &pdu, SCTP_DATA_IND(received_message_p).payload)) {
        // TODO: Notify gNB of failure with right cause
        OAILOG_ERROR(LOG_NGAP, " Failed to decode new buffer\n");

      } else {
        OAILOG_INFO(LOG_NGAP, "decode new buffer done\n");
        ngap_amf_handle_message(
            state, SCTP_DATA_IND(received_message_p).assoc_id,
            SCTP_DATA_IND(received_message_p).stream, &pdu);
        OAILOG_INFO(LOG_NGAP, " decode new buffer done\n");
      }

      // Free received PDU array
      bdestroy_wrapper(&SCTP_DATA_IND(received_message_p).payload);

#endif /*decoder from lib*/
    } break;

    case SCTP_NEW_ASSOCIATION: {
      OAILOG_DEBUG(LOG_NGAP, "#####SCTP_NEW_ASSOCIATION\n");
      increment_counter("amf_new_association", 1, NO_LABELS);
      if (ngap_handle_new_association(
              state, &received_message_p->ittiMsg.sctp_new_peer)) {
        increment_counter("amf_new_association", 1, 1, "result", "failure");
      } else {
        increment_counter("amf_new_association", 1, 1, "result", "success");
      }
    } break;

    case NGAP_PDUSESSION_RESOURCE_SETUP_REQ: {
      ngap_generate_ngap_pdusession_resource_setup_req(
          state, &NGAP_PDUSESSION_RESOURCE_SETUP_REQ(received_message_p));
    } break;

    case NGAP_PDUSESSIONRESOURCE_REL_REQ: {
      ngap_generate_ngap_pdusession_resource_rel_cmd(
          state, &NGAP_PDUSESSIONRESOURCE_REL_REQ(received_message_p));
    } break;

    case AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION: {
      OAILOG_ERROR(LOG_NGAP, " ##### AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION\n");
      ngap_handle_amf_ue_id_notification(
          state, &AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION(received_message_p));
    } break;

    case NGAP_NAS_DL_DATA_REQ: {  // packets from NAS
                                  /*
                                   * New message received from NAS task.
                                   * * * * This corresponds to a NGAP downlink nas transport message.
                                   */
      OAILOG_ERROR(LOG_NGAP, ":%s, NGAP_NAS_DL_DATA_REQ\n", __func__);
      ngap_generate_downlink_nas_transport(
          state, NGAP_NAS_DL_DATA_REQ(received_message_p).gnb_ue_ngap_id,
          NGAP_NAS_DL_DATA_REQ(received_message_p).amf_ue_ngap_id,
          &NGAP_NAS_DL_DATA_REQ(received_message_p).nas_msg, imsi64);

#if 0
      bstring payload;
      static int count = 1;
      int length = 0; 

      /*initial_ue_msgi: 71-6 */
      char* new_data =
          "\x00\x0f\x40\x3d\x00\x00\x04\x00\x55\x00\x02\x00\x01\x00\x26\x00\x18"
          "\x17\x7e\x00\x41\x71\x00\x0d\x01\x13\x00\x14\xf0\xff\x00\x00\x01\x00"
          "\x00\x00\xf1\x2e\x02\x80\x40\x00\x79\x00\x0f\x40\x13\x40\x01\x00\x00"
          "\x00\x5c\xd0\x13\x40\x01\x00\xa0\x00\x00\x5a\x40\x01\x18";

      switch (count) {
        case 0: {
          /*initial_ue_msgi: 71-6 */
          new_data =
              "\x00\x0f\x40\x3d\x00\x00\x04\x00\x55\x00\x02\x00\x01\x00\x26\x00"
              "\x18\x17\x7e\x00\x41\x71\x00\x0d\x01\x13\x00\x14\xf0\xff\x00\x00"
              "\x01\x00\x00\x00\xf1\x2e\x02\x80\x40\x00\x79\x00\x0f\x40\x13\x40"
              "\x01\x00\x00\x00\x5c\xd0\x13\x40\x01\x00\xa0\x00\x00\x5a\x40\x01"
              "\x18";

        } break;

        case 1: {
          /*UplnkNas: identity response : 75-6*/
	 length =69;
          new_data =
              "\x00\x2e\x40\x40\x00\x00\x04\x00\x0a\x00\x02\x00\x01\x00\x55\x00"
              "\x02\x00\x01\x00\x26\x00\x1a\x19\x7e\x02\xd3\x40\xbc\x55\x01\x7e"
              "\x00\x5c\x00\x0d\x01\x13\x00\x14\xf0\xff\x00\x00\x01\x00\x00\x00"
              "\xf1\x00\x79\x40\x0f\x40\x13\x40\x01\x00\x00\x00\x5c\xd0\x13\x40"
              "\x01\x00\xa0\x00";

        } break;

        case 2: {
          /*UplnkNas: Auth response : 75-6*/
	 length = 64;
          new_data =
              "\x00\x2e\x40\x3c\x00\x00\x04\x00\x0a\x00\x02\x00\x01\x00\x55\x00"
              "\x02\x00\x01\x00\x26\x00\x16\x15\x7e\x00\x57\x2d\x10\x25\xe8\x7b"
              "\x06\x52\xc3\xc6\x3b\x36\x82\x8b\x54\x51\x7e\xbf\x15\x00\x79\x40"
              "\x0f\x40\x13\x40\x01\x00\x00\x00\x5c\xd0\x13\x40\x01\x00\xa0"
              "\x00";

        } break;

        case 3: {
          /*UplnkNas: Security mode Complete : 75-6*/
			length = 65;
          new_data =
              "\x00\x2e\x40\x3d\x00\x00\x04\x00\x0a\x00\x02\x00\x01\x00\x55\x00"
              "\x02\x00\x01\x00\x26\x00\x17\x16\x7e\x04\xd1\xf8\x69\x38\x00\x7e"
              "\x00\x5e\x77\x00\x09\x05\x00\x00\x00\x00\x00\x00\x00\xf0\x00\x79"
              "\x40\x0f\x40\x13\x40\x01\x00\x00\x00\x5c\xd0\x13\x40\x01\x00\xa0"
              "\x00";

        } break;

    /*    case 4: {
         // Initial Context Setup Response: 19 
			length=19;
          new_data =
              "\x20\x0e\x00\x0f\x00\x00\x02\x00\x0a\x40\x02\x00\x01\x00\x55\x40"
              "\x02\x00\x01";
        } break;
*/
        case 4: {
          /*UplnkNas: Registration Complete : 75-6*/
			length=53;
          new_data =
              "\x00\x2e\x40\x31\x00\x00\x04\x00\x0a\x00\x02\x00\x01\x00\x55\x00"
              "\x02\x00\x01\x00\x26\x00\x0b\x0a\x7e\x02\xad\xb6\xc7\xc9\x02\x7e"
              "\x00\x43\x00\x79\x40\x0f\x40\x13\x40\x01\x00\x00\x00\x5c\xd0\x13"
              "\x40\x01\x00\xa0\x00";
        } break;

        default:
          OAILOG_ERROR(LOG_SCTP, "default case\n");
      }

      count++;
      
      OAILOG_ERROR(LOG_SCTP,"bef blk2bstr, count: %d\n", count);
     // payload = blk2bstr((unsigned char*) new_data, 75 - 6);
      payload = blk2bstr((unsigned char*) new_data, length);
      if (payload == NULL) {
        OAILOG_ERROR(LOG_SCTP, "failed to allocate bstr for SendUl\n");
      }
      OAILOG_ERROR(LOG_SCTP, "after blk2bstr:%p\n", payload);
      /**/
      Ngap_NGAP_PDU_t pdu = {0};
      // if (ngap_amf_decode_pdu( &pdu,
      // SCTP_DATA_IND(received_message_p).payload)) {
      if (ngap_amf_decode_pdu(&pdu, payload)) {
        // TODO: Notify gNB of failure with right cause
        OAILOG_ERROR(LOG_NGAP, "########Failed to decode new buffer\n");

      } else {
        OAILOG_ERROR(LOG_NGAP, "#######decode new buffer DONE\n");
        ngap_amf_handle_message(
            state, SCTP_DATA_IND(received_message_p).assoc_id,
            SCTP_DATA_IND(received_message_p).stream, &pdu);
      }

      // Free received PDU array
      bdestroy_wrapper(&payload);
      /**/
     if(count == 5) {
             sleep(30);
	     length = 67;
	     new_data = "\x00\x2e\x40\x3f\x00\x00\x04\x00\x0a\x00\x02\x00\x04\x00\x55\x00\x02\x00\x04\x00\x26\x00\x19\x18\x7e\x01\x47\x17\x33\xd4\x07\x7e\x00\x45\x01\x00\x0b\xf2\x13\x00\x14\x44\x33\x12\x00\x00\x00\x01\x00\x79\x40\x0f\x40\x13\x40\x01\x00\x00\x00\x5c\xd0\x13\x40\x01\x00\xa0\x00";

      OAILOG_ERROR(LOG_SCTP, "ACL_TAG bef blk2bstr, count: %d\n", count);
     // payload = blk2bstr((unsigned char*) new_data, 75 - 6);
      payload = blk2bstr((unsigned char*) new_data, length);
      if (payload == NULL) {
        OAILOG_ERROR(LOG_SCTP, "failed to allocate bstr for SendUl\n");
      }
      OAILOG_ERROR(LOG_SCTP, "ACL_TAG after blk2bstr:%p\n", payload);
      /**/
      Ngap_NGAP_PDU_t pdu = {0};
      // if (ngap_amf_decode_pdu( &pdu,
      // SCTP_DATA_IND(received_message_p).payload)) {
      if (ngap_amf_decode_pdu(&pdu, payload)) {
        // TODO: Notify gNB of failure with right cause
        OAILOG_ERROR(LOG_NGAP, "########Failed to decode new buffer\n");

      } else {
        OAILOG_ERROR(LOG_NGAP, "#######decode new buffer DONE\n");
        ngap_amf_handle_message(
            state, SCTP_DATA_IND(received_message_p).assoc_id,
            SCTP_DATA_IND(received_message_p).stream, &pdu);
      }

      // Free received PDU array
      bdestroy_wrapper(&payload);
      /**/
     }
#endif

    } break;

#if 0 /* TODO later*/

    // SCTP layer notifies NGAP of disconnection of a peer.
    case SCTP_CLOSE_ASSOCIATION: {
      ngap_handle_sctp_disconnection(
          state, SCTP_CLOSE_ASSOCIATION(received_message_p).assoc_id,
          SCTP_CLOSE_ASSOCIATION(received_message_p).reset);
    } break;

    */
    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      zframe_destroy(&msg_frame);
      ngap_amf_exit(); // Need change
    } break;

#endif

    default: {
      OAILOG_ERROR(
          LOG_NGAP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  put_ngap_state();
  put_ngap_imsi_map();
  // put_ngap_ue_state(imsi64);
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
    OAILOG_ERROR(
        LOG_NGAP, " ######Error while sending SCTP_INIT_MSG to SCTP \n");
  } else
    OAILOG_ERROR(LOG_NGAP, "####### sending SCTP_INIT_MSG to SCTP \n");

  zloop_start(ngap_task_zmq_ctx.event_loop);
  ngap_amf_exit();
  return NULL;
}
// ACL_TAG temp to test remove later
void guamfi_config_init(guamfi_config_t* guamfi_conf) {
  guamfi_conf->nb                 = 1;
  guamfi_conf->guamfi[0].amf_code = AMFC;
  guamfi_conf->guamfi[0].amf_gid  = AMFGID;
  //  guamfi_conf->guamfi[0].amf_Pointer     = AMFPOINTER;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 0x0F;
}

//------------------------------------------------------------------------------
int ngap_amf_init(/*const*/ amf_config_t* amf_config_p) {
  OAILOG_ERROR(LOG_NGAP, "Initializing NGAP interface\n");

  // ACL_TAG temp to test remove later
  /*init amfconfig*/
#if 1
  amf_config_t* config = amf_config_p;

  memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  config->config_file                    = NULL;
  config->max_gnbs                       = 2;
  config->max_ues                        = 2;
  config->unauthenticated_imsi_supported = 0;
  config->relative_capacity              = RELATIVE_CAPACITY;
  config->amf_statistic_timer            = AMF_STATISTIC_TIMER_S;

  guamfi_config_init(&config->guamfi);
#endif
#if 0  // TODO Need change  ##############################
  if (get_asn1c_environment_version() < ASN1_MINIMUM_VERSION) {
    OAILOG_ERROR(
      LOG_NGAP,
      "ASN1C version %d found, expecting at least %d\n",
      get_asn1c_environment_version(),
      ASN1_MINIMUM_VERSION);
    return RETURNerror;
  }

  //OAILOG_DEBUG(LOG_NGAP, "ASN1C version %d\n", get_asn1c_environment_version());
  // OAILOG_DEBUG(LOG_NGAP, "NGAP Release v10.5\n"); // Need change

#endif
  /*set the amf_config temp*/

  OAILOG_ERROR(LOG_NGAP, "####cobraranu");
  if (ngap_state_init(
          amf_config_p->max_ues, amf_config_p->max_gnbs,
          amf_config_p->use_stateless) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Error while initing NGAP state\n");
    OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
    return RETURNerror;
  }

  OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
  if (itti_create_task(TASK_NGAP, &ngap_amf_thread, NULL) == RETURNerror) {
    OAILOG_ERROR(LOG_NGAP, " Error while creating NGAP task\n");
    return RETURNerror;
  }

  OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
  OAILOG_DEBUG(LOG_NGAP, " Initializing NGAP interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
void ngap_amf_exit(void) {
  OAILOG_DEBUG(LOG_NGAP, "Cleaning NGAP\n");

  destroy_task_context(&ngap_task_zmq_ctx);

  put_ngap_imsi_map();

  ngap_state_exit();  // Need change

  OAILOG_DEBUG(LOG_NGAP, "Cleaning NGAP: DONE\n");
  OAI_FPRINTF_INFO("TASK_NGAP terminated\n");
  pthread_exit(NULL);
}

/*
//------------------------------------------------------------------------------
void ngap_dump_gnb_list(ngap_state_t* state)
{
  hashtable_ts_apply_callback_on_elements(
    &state->gnbs, ngap_dump_gnb_hash_cb, NULL, NULL);
}

//------------------------------------------------------------------------------
bool ngap_dump_gnb_hash_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void* const gNB_void,
  void __attribute__((unused)) * unused_parameterP,
  void __attribute__((unused)) * *unused_resultP)
{
  const gnb_description_t* const enb_ref = (const enb_description_t*) eNB_void;
  if (enb_ref == NULL) {
    return false;
  }
  ngap_dump_enb(enb_ref);
  return false;
}

//------------------------------------------------------------------------------
void ngap_dump_enb(const enb_description_t* const enb_ref)
{
#ifdef NGAP_DEBUG_LIST
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
  sctp_assoc_id_t sctp_assoc_id = enb_ref->sctp_assoc_id;

  hash_table_ts_t* state_ue_ht = get_ngap_ue_state();
  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t* const) state_ue_ht,
    ngap_dump_ue_hash_cb,
    &sctp_assoc_id,
    NULL);
  indent--;
  eNB_LIST_OUT("");
#else
  ngap_dump_ue(NULL);
#endif
}

//------------------------------------------------------------------------------
bool ngap_dump_ue_hash_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void* const ue_void,
  void* parameter,
  void __attribute__((unused)) * *unused_resultP)
{
  ue_description_t* ue_ref = (ue_description_t*) ue_void;
  sctp_assoc_id_t* sctp_assoc_id = (sctp_assoc_id_t *) parameter;
  if (ue_ref == NULL) {
    return false;
  }

  if(ue_ref->sctp_assoc_id == *sctp_assoc_id) {
    ngap_dump_ue(ue_ref);
  }
  return false;
}

//------------------------------------------------------------------------------
void ngap_dump_ue(const ue_description_t* const ue_ref)
{
#ifdef NGAP_DEBUG_LIST

  if (ue_ref == NULL) return;

  UE_LIST_OUT("eNB UE ngap id:   0x%06x", ue_ref->enb_ue_ngap_id);
  UE_LIST_OUT("MME UE ngap id:   0x%08x", ue_ref->amf_ue_ngap_id);
  UE_LIST_OUT("SCTP stream recv: 0x%04x", ue_ref->sctp_stream_recv);
  UE_LIST_OUT("SCTP stream send: 0x%04x", ue_ref->sctp_stream_send);
#endif
}
*/

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
  hashtable_uint64_ts_init(&gnb_ref->ue_id_coll, amf_config.max_ues, NULL, bs);
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

/*
//------------------------------------------------------------------------------
void ngap_remove_ue(n1ap_state_t* state, ue_description_t* ue_ref) {
  enb_description_t* enb_ref = NULL;

  // NULL reference...
  if (ue_ref == NULL) return;

  amf_ue_ngap_id_t amf_ue_ngap_id = ue_ref->amf_ue_ngap_id;
  enb_ref = n1ap_state_get_enb(state, ue_ref->sctp_assoc_id);
  DevAssert(enb_ref->nb_ue_associated > 0);
  // Updating number of UE
  enb_ref->nb_ue_associated--;

  // Stop UE Context Release Complete timer,if running
  if (ue_ref->ngap_ue_context_rel_timer.id != NGAP_TIMER_INACTIVE_ID) {
    if (timer_remove(ue_ref->ngap_ue_context_rel_timer.id, NULL)) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to stop ngap ue context release complete timer, UE id: %d\n",
          ue_ref->amf_ue_ngap_id);
    }
    ue_ref->ngap_ue_context_rel_timer.id = NGAP_TIMER_INACTIVE_ID;
  }
  OAILOG_TRACE(
      LOG_NGAP,
      "Removing UE enb_ue_ngap_id: " ENB_UE_NGAP_ID_FMT
      " amf_ue_ngap_id:" MME_UE_NGAP_ID_FMT " in eNB id : %d\n",
      ue_ref->enb_ue_ngap_id, ue_ref->amf_ue_ngap_id, enb_ref->enb_id);

  ue_ref->n1_ue_state = NGAP_UE_INVALID_STATE;

  hash_table_ts_t* state_ue_ht = get_ngap_ue_state();
  hashtable_ts_free(state_ue_ht, ue_ref->comp_ngap_id);
  hashtable_ts_free(&state->amfid2associd, amf_ue_ngap_id);
  hashtable_uint64_ts_free(&enb_ref->ue_id_coll, amf_ue_ngap_id);

  imsi64_t imsi64                = INVALID_IMSI64;
  ngap_imsi_map_t* ngap_imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      ngap_imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) amf_ue_ngap_id,
      &imsi64);
  delete_ngap_ue_state(imsi64);

  if (!enb_ref->nb_ue_associated) {
    if (enb_ref->ng_state == NGAP_RESETING) {
      OAILOG_INFO(LOG_NGAP, "Moving eNB state to NGAP_INIT \n");
      enb_ref->ng_state = NGAP_INIT;
      set_gauge("ng_connection", 0, 1, "enb_name", enb_ref->enb_name);
      update_amf_app_stats_connected_enb_sub();
    } else if (enb_ref->ng_state == NGAP_SHUTDOWN) {
      OAILOG_INFO(LOG_NGAP, "Deleting eNB \n");
      set_gauge("ng_connection", 0, 1, "enb_name", enb_ref->enb_name);
      ngap_remove_enb(state, enb_ref);
      update_amf_app_stats_connected_enb_sub();
    }
  }
}

*/

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
