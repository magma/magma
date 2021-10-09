/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include "mme_app_mme_ue_id_timer_id.h"

#include <bits/stdc++.h>

using namespace std;

typedef unordered_map<uint32_t, long> MmeUeIdTimerIdMap;

MmeUeIdTimerIdMap mme_ue_id_timer_id_map;

void initialize_mme_ue_id_timer_id_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_id_timer_id_map = MmeUeIdTimerIdMap{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_upsert_mme_ue_id_timer_id(mme_ue_s1ap_id_t mme_ue_id, long timer_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = mme_ue_id_timer_id_map.find(mme_ue_id);
  if (itr == mme_ue_id_timer_id_map.end()) {
    mme_ue_id_timer_id_map[mme_ue_id] = timer_id;
    OAILOG_DEBUG(
        LOG_MME_APP, "Inserting mme_ue_id %u entry for timer id: %lu \n", mme_ue_id,
        timer_id);
  } else {
      OAILOG_WARNING(LOG_MME_APP, "Replacing current timer id: %lu with new timer id: %lu", mme_ue_id_timer_id_map[mme_ue_id], timer_id);
    mme_ue_id_timer_id_map[mme_ue_id] = timer_id;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

long mme_app_get_timer_id_from_mme_ue_id(mme_ue_s1ap_id_t mme_ue_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = mme_ue_id_timer_id_map.find(mme_ue_id);
  if (itr == mme_ue_id_timer_id_map.end()) {
    OAILOG_WARNING(
        LOG_MME_APP, "No timer_id found on mme_ue_id_timer_id_map for mme_ue_id %u\n", mme_ue_id);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP, " Found timer_id: %ld for mme_ue_id %u \n", itr->second, mme_ue_id);
    return itr->second;
  }
  return NAS_TIMER_INACTIVE_ID;
}

void mme_app_remove_mme_ue_id_timer_id(mme_ue_s1ap_id_t mme_ue_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = mme_ue_id_timer_id_map.find(mme_ue_id);
  if (itr == mme_ue_id_timer_id_map.end()) {
    OAILOG_ERROR(
        LOG_MME_APP, "No timer_id found on mme_ue_id_timer_id_map for mme_ue_id %u\n", mme_ue_id);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP, "Deleting timer_id entry: %lu \n", itr->second);
    mme_ue_id_timer_id_map.erase(itr);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
