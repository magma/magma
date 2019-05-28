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

#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "s6a_messages_types.h"

int delete_subscriber_request(const char *imsi, const uint imsi_len)
{
  // send it to MME module for further processing
  MessageDef *message_p = NULL;
  s6a_cancel_location_req_t *s6a_cancel_location_req_p = NULL;
  message_p = itti_alloc_new_message(TASK_S6A, S6A_CANCEL_LOCATION_REQ);
  s6a_cancel_location_req_p = &message_p->ittiMsg.s6a_cancel_location_req;
  memcpy(s6a_cancel_location_req_p->imsi, imsi, imsi_len);
  s6a_cancel_location_req_p->imsi[imsi_len] = '\0';
  s6a_cancel_location_req_p->imsi_length = imsi_len;
  s6a_cancel_location_req_p->cancellation_type = SUBSCRIPTION_WITHDRAWL;
  int ret = 0;
  ret = itti_send_msg_to_task(TASK_MME_APP, message_p);
  return ret;
}

void handle_reset_request(void)
{
  // send it to MME module for further processing
  MessageDef *message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S6A, S6A_RESET_REQ);
  // TBD - To add support for partial reset
  itti_send_msg_to_task(TASK_MME_APP, message_p);
  return;
}
