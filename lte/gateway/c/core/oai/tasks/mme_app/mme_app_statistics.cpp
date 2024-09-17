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

#include "lte/gateway/c/core/oai/include/mme_app_statistics.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"

/*********************************** Utility Functions to update
 * Statistics**************************************/

/*****************************************************/
static inline int get_max(int num1, int num2) {
  return (num1 > num2 ? num1 : num2);
}
// Number of Connected UEs
void update_mme_app_stats_connected_ue_add(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  (mme_app_desc_p->nb_ue_connected)++;
  mme_app_desc_p->nb_ue_idle = get_max(
      mme_app_desc_p->nb_ue_attached - mme_app_desc_p->nb_ue_connected, 0);
  return;
}
void update_mme_app_stats_connected_ue_sub(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (mme_app_desc_p->nb_ue_connected != 0) (mme_app_desc_p->nb_ue_connected)--;
  mme_app_desc_p->nb_ue_idle = get_max(
      mme_app_desc_p->nb_ue_attached - mme_app_desc_p->nb_ue_connected, 0);
  return;
}

/*****************************************************/
// Number of S1U Bearers
void update_mme_app_stats_s1u_bearer_add(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  (mme_app_desc_p->nb_s1u_bearers)++;
  return;
}
void update_mme_app_stats_s1u_bearer_sub(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (mme_app_desc_p->nb_s1u_bearers != 0) (mme_app_desc_p->nb_s1u_bearers)--;
  return;
}

/*****************************************************/
// Number of Default EPS Bearers
void update_mme_app_stats_default_bearer_add(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  (mme_app_desc_p->nb_default_eps_bearers)++;
  return;
}
void update_mme_app_stats_default_bearer_sub(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (mme_app_desc_p->nb_default_eps_bearers != 0)
    (mme_app_desc_p->nb_default_eps_bearers)--;
  return;
}

/*****************************************************/
// Number of Attached UEs
void update_mme_app_stats_attached_ue_add(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  (mme_app_desc_p->nb_ue_attached)++;
  return;
}
void update_mme_app_stats_attached_ue_sub(void) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (mme_app_desc_p->nb_ue_attached != 0) (mme_app_desc_p->nb_ue_attached)--;
  mme_app_desc_p->nb_ue_idle = get_max(
      mme_app_desc_p->nb_ue_attached - mme_app_desc_p->nb_ue_connected, 0);
  return;
}
/*****************************************************/
