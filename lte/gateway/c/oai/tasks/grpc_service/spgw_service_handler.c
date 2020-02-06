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
#include <string.h>
#include <sys/types.h>
#include <conversions.h>

#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "sgw_messages_types.h"

int send_activate_bearer_request_itti(
  itti_pgw_nw_init_actv_bearer_request_t *itti_msg)
{
  OAILOG_DEBUG(LOG_SPGW_APP, "Sending pgw_nw_init_actv_bearer_request\n");
  MessageDef *message_p = itti_alloc_new_message(
    TASK_GRPC_SERVICE, PGW_NW_INITIATED_ACTIVATE_BEARER_REQ);
  message_p->ittiMsg.pgw_nw_init_actv_bearer_request = *itti_msg;

  IMSI_STRING_TO_IMSI64((char*) itti_msg->imsi, &message_p->ittiMsgHeader.imsi);

  return itti_send_msg_to_task(TASK_PGW_APP, INSTANCE_DEFAULT, message_p);
}

int send_deactivate_bearer_request_itti(
  itti_pgw_nw_init_deactv_bearer_request_t* itti_msg)
{
  OAILOG_DEBUG(LOG_SPGW_APP, "Sending pgw_nw_init_deactv_bearer_request\n");
  MessageDef* message_p = itti_alloc_new_message(
    TASK_GRPC_SERVICE, PGW_NW_INITIATED_DEACTIVATE_BEARER_REQ);
  message_p->ittiMsg.pgw_nw_init_deactv_bearer_request = *itti_msg;

  IMSI_STRING_TO_IMSI64((char*) itti_msg->imsi, &message_p->ittiMsgHeader.imsi);

  return itti_send_msg_to_task(TASK_PGW_APP, INSTANCE_DEFAULT, message_p);
}
