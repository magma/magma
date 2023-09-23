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

/*! \file sctp_itti_messaging.hpp
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#pragma once

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif

#define S1AP 18
#define NGAP 60

extern task_zmq_ctx_t sctp_task_zmq_ctx;

status_code_e sctp_itti_send_lower_layer_conf(
    task_id_t origin_task_id, sctp_ppid_t ppid, sctp_assoc_id_t assoc_id,
    sctp_stream_id_t stream, uint32_t agw_ue_xap_id, bool is_success);

status_code_e sctp_itti_send_new_association(sctp_ppid_t ppid,
                                             sctp_assoc_id_t assoc_id,
                                             sctp_stream_id_t instreams,
                                             sctp_stream_id_t outstreams,
                                             STOLEN_REF bstring* ran_cp_ipaddr);

status_code_e sctp_itti_send_new_message_ind(STOLEN_REF bstring* payload,
                                             sctp_ppid_t ppid,
                                             sctp_assoc_id_t assoc_id,
                                             sctp_stream_id_t stream);

status_code_e sctp_itti_send_com_down_ind(sctp_ppid_t ppid,
                                          sctp_assoc_id_t assoc_id, bool reset);
