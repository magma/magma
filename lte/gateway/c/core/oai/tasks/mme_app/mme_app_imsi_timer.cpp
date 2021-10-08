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

#include "mme_app_imsi_timer_id.h"

#include <bits/stdc++.h>

using namespace std;

typedef unordered_map<uint64_t, long> ImsiTimerIdMap;

ImsiTimerIdMap imsi_timer_id_map;

void initialize_imsi_timer_id_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  imsi_timer_id_map = ImsiTimerIdMap{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_insert_imsi_timer_id(imsi64_t imsi64, long timer_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = imsi_timer_id_map.find(imsi64);
  if (itr == imsi_timer_id_map.end()) {
    imsi_timer_id_map[imsi64] = timer_id;
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64, "Inserting IMSI entry for timer id: %lu \n",
        timer_id);
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64, "Inserting IMSI for existing timer_id: %lu \n",
        timer_id);
    imsi_timer_id_map[imsi64] = timer_id;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

long mme_app_get_timer_id_from_imsi(imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = imsi_timer_id_map.find(imsi64);
  if (itr == imsi_timer_id_map.end()) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64, "No timer_id found on imsi_timer_id_map \n");
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64, " Found timer_id: %ld \n", itr->second);
    return itr->second;
  }
  return -1;
}

void mme_app_remove_imsi_timer_id(imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = imsi_timer_id_map.find(imsi64);
  if (itr == imsi_timer_id_map.end()) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64, "No timer_id found on imsi_timer_id_map \n");
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64, "Deleting timer_id entry: %lu \n", itr->second);
    imsi_timer_id_map.erase(itr);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
