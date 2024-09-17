/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/amf/include/amf_app_statistics.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"

// Utility Functions to update Statistics
namespace magma5g {
static inline int get_max(int num1, int num2) {
  return (num1 > num2 ? num1 : num2);
}
/*****************************************************/
// Number of Connected UEs
void update_amf_app_stats_connected_ue_add(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  (amf_app_desc_p->nb_ue_connected)++;
  amf_app_desc_p->nb_ue_idle = get_max(
      amf_app_desc_p->nb_ue_attached - amf_app_desc_p->nb_ue_connected, 0);
  return;
}
void update_amf_app_stats_connected_ue_sub(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  if (amf_app_desc_p->nb_ue_connected != 0) {
    (amf_app_desc_p->nb_ue_connected)--;
  }
  amf_app_desc_p->nb_ue_idle = get_max(
      amf_app_desc_p->nb_ue_attached - amf_app_desc_p->nb_ue_connected, 0);
  return;
}

/*****************************************************/
// Number of Registered UEs
void update_amf_app_stats_registered_ue_add(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  (amf_app_desc_p->nb_ue_attached)++;
  return;
}
void update_amf_app_stats_registered_ue_sub(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  if (amf_app_desc_p->nb_ue_attached != 0) (amf_app_desc_p->nb_ue_attached)--;
  amf_app_desc_p->nb_ue_idle = get_max(
      amf_app_desc_p->nb_ue_attached - amf_app_desc_p->nb_ue_connected, 0);
  return;
}

/*****************************************************/
// Number of Pdusessions UEs
void update_amf_app_stats_pdusessions_ue_add(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  (amf_app_desc_p->nb_pdu_sessions)++;
  return;
}
void update_amf_app_stats_pdusessions_ue_sub(void) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (amf_app_desc_p == NULL) {
    return;
  }
  if (amf_app_desc_p->nb_pdu_sessions != 0) (amf_app_desc_p->nb_pdu_sessions)--;
  return;
}

/*****************************************************/
}  // namespace magma5g
