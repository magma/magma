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

#ifndef FILE_SERVICE303_MESSAGE_UTILS
#define FILE_SERVICE303_MESSAGE_UTILS

#include <stdbool.h>

#include "intertask_interface.h"
#include "intertask_interface_types.h"

int send_app_health_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id, bool healthy);

int send_mme_app_stats_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id,
    application_mme_app_stats_msg_t* stats_msg);

int send_s1ap_stats_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id,
    application_s1ap_stats_msg_t* stats_msg);

#endif /* FILE_SERVICE303_MESSAGE_UTILS */
