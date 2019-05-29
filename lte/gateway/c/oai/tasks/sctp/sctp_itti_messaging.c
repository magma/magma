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

/*! \file sctp_itti_messaging.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <string.h>
#include <stdbool.h>

#include "intertask_interface.h"
#include "sctp_itti_messaging.h"
#include "itti_types.h"
#include "sctp_messages_types.h"

//------------------------------------------------------------------------------
int sctp_itti_send_lower_layer_conf(
  const task_id_t origin_task_id,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  const uint32_t mme_ue_s1ap_id,
  const bool is_success)
{
  MessageDef *message_p = NULL;
  sctp_data_cnf_t *sctp_p = NULL;

  message_p = itti_alloc_new_message(TASK_SCTP, SCTP_DATA_CNF);
  sctp_p = &message_p->ittiMsg.sctp_data_cnf;
  sctp_p->assoc_id = assoc_id;
  sctp_p->stream = stream;
  sctp_p->mme_ue_s1ap_id = mme_ue_s1ap_id;
  sctp_p->is_success = is_success;
  return itti_send_msg_to_task(origin_task_id, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int sctp_itti_send_new_association(
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t instreams,
  const sctp_stream_id_t outstreams)
{
  MessageDef *message_p = NULL;
  sctp_new_peer_t *sctp_new_peer_p = NULL;

  message_p = itti_alloc_new_message(TASK_SCTP, SCTP_NEW_ASSOCIATION);
  sctp_new_peer_p = &message_p->ittiMsg.sctp_new_peer;
  sctp_new_peer_p->assoc_id = assoc_id;
  sctp_new_peer_p->instreams = instreams;
  sctp_new_peer_p->outstreams = outstreams;
  return itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int sctp_itti_send_new_message_ind(
  STOLEN_REF bstring *payload,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  const sctp_stream_id_t instreams,
  const sctp_stream_id_t outstreams)
{
  MessageDef *message_p = itti_alloc_new_message(TASK_SCTP, SCTP_DATA_IND);
  if (message_p) {
    SCTP_DATA_IND(message_p).payload = *payload;
    STOLEN_REF *payload = NULL;
    SCTP_DATA_IND(message_p).stream = stream;
    SCTP_DATA_IND(message_p).assoc_id = assoc_id;
    SCTP_DATA_IND(message_p).instreams = instreams;
    SCTP_DATA_IND(message_p).outstreams = outstreams;
    return itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
  }
  return RETURNerror;
}

//------------------------------------------------------------------------------
int sctp_itti_send_com_down_ind(const sctp_assoc_id_t assoc_id, bool reset)
{
  MessageDef *message_p = NULL;
  sctp_close_association_t *sctp_close_association_p = NULL;

  message_p = itti_alloc_new_message(TASK_SCTP, SCTP_CLOSE_ASSOCIATION);
  sctp_close_association_p = &message_p->ittiMsg.sctp_close_association;
  sctp_close_association_p->assoc_id = assoc_id;
  sctp_close_association_p->reset = reset;
  return itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
}
