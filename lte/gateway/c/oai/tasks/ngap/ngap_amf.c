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
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#if HAVE_CONFIG_H
#include "config.h"
#endif

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
#include "ngap_amf_decoder.h"
#include "ngap_amf_handlers.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_itti_messaging.h"

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

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  ngap_state_t* state = NULL;
  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);

  imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
  state           = get_ngap_state(false);
  AssertFatal(state != NULL, "failed to retrieve ngap state (was null)");

  switch (ITTI_MSG_ID(received_message_p)) {

    case NGAP_PDUSESSION_RESOURCE_SETUP_REQ: {
      ngap_generate_ngap_pdusession_resource_setup_req(
          state, &NGAP_PDUSESSION_RESOURCE_SETUP_REQ(received_message_p));
    } break;


    default: {
      OAILOG_ERROR(
          LOG_NGAP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  put_ngap_state();
  put_ngap_imsi_map();
  put_ngap_ue_state(imsi64);
  itti_free_msg_content(received_message_p);
  zframe_destroy(&msg_frame);
  return 0;
}




