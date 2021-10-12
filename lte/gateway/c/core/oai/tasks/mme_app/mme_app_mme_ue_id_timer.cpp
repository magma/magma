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

#include "mme_app_mme_ue_id_timer.h"

#include <bits/stdc++.h>

typedef std::set<std::pair<mme_ue_s1ap_id_t, long>> MmeUeIdTimerIdSet;

MmeUeIdTimerIdSet mme_ue_id_timer_id_set;

void initialize_mme_ue_id_timer_id_set() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_id_timer_id_set = MmeUeIdTimerIdSet{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void clear_mme_ue_id_timer_id_set() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (!mme_ue_id_timer_id_set.empty()) {
    mme_ue_id_timer_id_set.clear();
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_insert_mme_ue_id_timer_id(
    mme_ue_s1ap_id_t mme_ue_id, long timer_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto ret_pair =
      mme_ue_id_timer_id_set.insert(std::make_pair(mme_ue_id, timer_id));
  if (ret_pair.second) {
    OAILOG_DEBUG(
        LOG_MME_APP, "Inserting mme_ue_id %u, timer id: %ld entry \n",
        mme_ue_id, timer_id);
  } else {
    OAILOG_WARNING(
        LOG_MME_APP, "mme_ue_id %u entry for timer id: %ld already exists\n",
        mme_ue_id, timer_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

bool mme_app_is_mme_ue_id_timer_id_key_valid(
    mme_ue_s1ap_id_t mme_ue_id, long timer_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  std::pair<mme_ue_s1ap_id_t, long> timer_ue_id_key =
      std::make_pair(mme_ue_id, timer_id);
  auto itr = mme_ue_id_timer_id_set.find(timer_ue_id_key);
  if (itr == mme_ue_id_timer_id_set.end()) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "No entry found on mme_ue_id_timer_id_set for mme_ue_id %u, timer_id "
        "%ld key\n",
        mme_ue_id, timer_id);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP, " Found timer_id: %ld for mme_ue_id %u entry \n", timer_id,
        mme_ue_id);
    return true;
  }
  return false;
}

void mme_app_remove_mme_ue_id_timer_id(
    mme_ue_s1ap_id_t mme_ue_id, long timer_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  int erased_elements =
      mme_ue_id_timer_id_set.erase(std::make_pair(mme_ue_id, timer_id));
  if (!erased_elements) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "No entry found on mme_ue_id_timer_id_set for mme_ue_id %u, timer_id "
        "%ld key\n",
        mme_ue_id, timer_id);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP, "Deleting mme_ue_id %u, timer_id %ld entries: %u \n",
        mme_ue_id, timer_id, erased_elements);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
