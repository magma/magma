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

/*! \file sctp_itti_messaging.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_SCTP_ITTI_MESSAGING_SEEN
#define FILE_SCTP_ITTI_MESSAGING_SEEN
#include <stdbool.h>
#include <stdint.h>

#include "common_defs.h"
#include "bstrlib.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"

#define S1AP 18
#define NGAP 60

extern task_zmq_ctx_t sctp_task_zmq_ctx;

int sctp_itti_send_lower_layer_conf(
    task_id_t origin_task_id, sctp_ppid_t ppid, sctp_assoc_id_t assoc_id,
    sctp_stream_id_t stream, uint32_t mme_ue_s1ap_id, bool is_success);

int sctp_itti_send_new_association(
    sctp_ppid_t ppid, sctp_assoc_id_t assoc_id, sctp_stream_id_t instreams,
    sctp_stream_id_t outstreams, STOLEN_REF bstring* ran_cp_ipaddr);

int sctp_itti_send_new_message_ind(
    STOLEN_REF bstring* payload, sctp_ppid_t ppid, sctp_assoc_id_t assoc_id,
    sctp_stream_id_t stream);

int sctp_itti_send_com_down_ind(
    sctp_ppid_t ppid, sctp_assoc_id_t assoc_id, bool reset);

#endif /* FILE_SCTP_ITTI_MESSAGING_SEEN */
