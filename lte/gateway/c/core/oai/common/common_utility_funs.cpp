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

#include <stdint.h>
#include <iostream>
#include "common_utility_funs.h"
#include "lte/protos/mconfig/mconfigs.pb.h"

// Extract MCC and MNC from the imsi received and match with
// configuration
extern "C" int match_fed_mode_map(const char* imsi, log_proto_t module) {
  OAILOG_FUNC_IN(module);
  imsi64_t imsi64;
  IMSI_STRING_TO_IMSI64(imsi, &imsi64);
  uint8_t mcc_d1 = imsi[0] - '0';
  uint8_t mcc_d2 = imsi[1] - '0';
  uint8_t mcc_d3 = imsi[2] - '0';
  uint8_t mnc_d1 = imsi[3] - '0';
  uint8_t mnc_d2 = imsi[4] - '0';
  uint8_t mnc_d3 = imsi[5] - '0';
  if ((mcc_d1 < 0 || mcc_d1 > 9) || (mcc_d2 < 0 || mcc_d2 > 9) ||
      (mcc_d3 < 0 || mcc_d3 > 9) || (mnc_d1 < 0 || mnc_d1 > 9) ||
      (mnc_d2 < 0 || mnc_d2 > 9) || (mnc_d3 < 0 || mnc_d3 > 9)) {
    OAILOG_ERROR_UE(module, imsi64, "MCC/MNC is not a decimal digit \n");
    OAILOG_FUNC_RETURN(module, RETURNerror);
  }
  for (uint8_t itr = 0; itr < mme_config.mode_map_config.num; itr++) {
    if (((mcc_d1 == mme_config.mode_map_config.mode_map[itr].plmn.mcc_digit1) &&
         (mcc_d2 == mme_config.mode_map_config.mode_map[itr].plmn.mcc_digit2) &&
         (mcc_d3 == mme_config.mode_map_config.mode_map[itr].plmn.mcc_digit3) &&
         (mnc_d1 == mme_config.mode_map_config.mode_map[itr].plmn.mnc_digit1) &&
         (mnc_d2 ==
          mme_config.mode_map_config.mode_map[itr].plmn.mnc_digit2))) {
      if (mme_config.mode_map_config.mode_map[itr].plmn.mnc_digit3 != 0xf) {
        if (mnc_d3 !=
            mme_config.mode_map_config.mode_map[itr].plmn.mnc_digit3) {
          continue;
        }
      }
      OAILOG_FUNC_RETURN(module, mme_config.mode_map_config.mode_map[itr].mode);
    }
  }
  // If the plmn is not configured set the default mode as hss + spgw_task.
  OAILOG_INFO_UE(
      module, imsi64,
      "PLMN is not configured. Selecting default mode: SPGW_SUBSCRIBER \n");
  OAILOG_FUNC_RETURN(
      module, magma::mconfig::ModeMapItem_FederatedMode_SPGW_SUBSCRIBER);
}

// Verify that tac is included in registered subscription areas
extern "C" int verify_service_area_restriction(
    tac_t tac, const regional_subscription_t* reg_sub, uint8_t num_reg_sub) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  tac_list_per_sac_t* tac_list = NULL;
  for (uint8_t itr = 0; itr < num_reg_sub; itr++) {
    hashtable_rc_t htbl = obj_hashtable_get(
        mme_config.sac_to_tacs_map.sac_to_tacs_map_htbl, reg_sub[itr].zone_code,
        ZONE_CODE_LEN, (void**) &tac_list);
    if (htbl == HASH_TABLE_OK) {
      for (uint8_t idx = 0; idx < tac_list->num_tac_entries; idx++) {
        if (tac_list->tacs[idx] == tac) {
          OAILOG_FUNC_RETURN(LOG_COMMON, RETURNok);
        }
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_COMMON, RETURNerror);
}