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

#ifndef SGS_MESSAGES_H_
#define SGS_MESSAGES_H_

#include "intertask_interface.h"

extern task_zmq_ctx_t sgs_task_zmq_ctx;

int sgs_send_eps_detach_indication(
    itti_sgsap_eps_detach_ind_t* sgs_eps_detach_ind_p);

int sgs_send_imsi_detach_indication(
    itti_sgsap_imsi_detach_ind_t* sgs_imsi_detach_ind_p);

int sgs_send_tmsi_reallocation_complete(
    itti_sgsap_tmsi_reallocation_comp_t* sgs_tmsi_realloc_comp_p);

int sgs_send_service_request(
    itti_sgsap_service_request_t* const sgs_service_request_p);

int sgs_send_paging_reject(
    itti_sgsap_paging_reject_t* const sgs_paging_reject_p);

int sgs_send_ue_unreachable(
    itti_sgsap_ue_unreachable_t* const sgs_ue_unreachable_p);

#endif
