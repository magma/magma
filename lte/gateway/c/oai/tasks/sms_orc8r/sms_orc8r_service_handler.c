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

#include <string.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "log.h"
#include "common_defs.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "sms_orc8r_messages.h"
#include "sgs_messages_types.h"

int handle_sms_orc8r_downlink_unitdata(
    const itti_sgsap_downlink_unitdata_t* sgs_dl_unitdata_p) {
  int rc = RETURNok;

  MessageDef* message_p                              = NULL;
  itti_sgsap_downlink_unitdata_t* sgs_dl_unit_data_p = NULL;

  message_p = itti_alloc_new_message(TASK_SMS_ORC8R, SGSAP_DOWNLINK_UNITDATA);
  AssertFatal(message_p, "itti_alloc_new_message Failed");
  sgs_dl_unit_data_p = &message_p->ittiMsg.sgsap_downlink_unitdata;
  memset((void*) sgs_dl_unit_data_p, 0, sizeof(itti_sgsap_downlink_unitdata_t));

  memcpy(
      sgs_dl_unit_data_p, sgs_dl_unitdata_p,
      sizeof(itti_sgsap_downlink_unitdata_t));
  // send it to NAS module for further processing
  OAILOG_DEBUG(
      LOG_SMS_ORC8R, "Received SMO Downlink UnitData message %s from Orc8r\n",
      sgs_dl_unit_data_p->nas_msg_container->data);
  rc = send_msg_to_task(&sms_orc8r_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
