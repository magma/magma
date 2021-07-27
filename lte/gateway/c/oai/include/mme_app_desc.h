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

/*! \file mme_app_desc.h */

#pragma once

#include <stdint.h>
#include <pthread.h>

#include "mme_app_ue_context.h"

typedef struct mme_app_desc_s {
  /* UE contexts */
  mme_ue_context_t mme_ue_contexts;

  long statistic_timer_id;
  uint32_t statistic_timer_period;
  mme_ue_s1ap_id_t mme_app_ue_s1ap_id_generator;

  /* ***************Statistics*************
   * number of attached UE,number of connected UE,
   * number of idle UE,number of default bearers,
   * number of S1_U bearers,number of PDN sessions
   */

  uint32_t nb_ue_attached;
  uint32_t nb_ue_connected;
  uint32_t nb_default_eps_bearers;
  uint32_t nb_s1u_bearers;

  /* ***************Changes in Statistics**************/

  uint32_t nb_ue_attached_since_last_stat;
  uint32_t nb_ue_detached_since_last_stat;
  uint32_t nb_ue_connected_since_last_stat;
  uint32_t nb_ue_disconnected_since_last_stat;
  uint32_t nb_eps_bearers_established_since_last_stat;
  uint32_t nb_eps_bearers_released_since_last_stat;
  uint32_t nb_s1u_bearers_released_since_last_stat;
  uint32_t nb_s1u_bearers_established_since_last_stat;
  uint32_t nb_ue_managed;
  uint32_t nb_ue_idle;
  uint32_t nb_bearers_managed;
  uint32_t nb_ue_since_last_stat;
  uint32_t nb_bearers_since_last_stat;
} mme_app_desc_t;
