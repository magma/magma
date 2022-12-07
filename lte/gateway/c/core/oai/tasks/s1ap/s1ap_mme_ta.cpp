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

/*! \file s1ap_mme_ta.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_ta.hpp"

#include <stdio.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "S1ap_BPLMNs.h"
#include "S1ap_PLMNidentity.h"
#include "S1ap_SupportedTAs-Item.h"
#include "S1ap_TAC.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.hpp"

static int s1ap_mme_compare_plmn(const S1ap_PLMNidentity_t* const plmn) {
  int i = 0;
  uint16_t mcc = 0;
  uint16_t mnc = 0;
  uint16_t mnc_len = 0;

  DevAssert(plmn != NULL);
  TBCD_TO_MCC_MNC(plmn, mcc, mnc, mnc_len);
  mme_config_read_lock(&mme_config);

  for (i = 0; i < mme_config.served_tai.nb_tai; i++) {
    OAILOG_TRACE(
        LOG_S1AP,
        "Comparing plmn_mcc %d/%d, plmn_mnc %d/%d plmn_mnc_len %d/%d\n",
        mme_config.served_tai.plmn_mcc[i], mcc,
        mme_config.served_tai.plmn_mnc[i], mnc,
        mme_config.served_tai.plmn_mnc_len[i], mnc_len);

    if ((mme_config.served_tai.plmn_mcc[i] == mcc) &&
        (mme_config.served_tai.plmn_mnc[i] == mnc) &&
        (mme_config.served_tai.plmn_mnc_len[i] == mnc_len))
      /*
       * There is a matching plmn
       */
      return TA_LIST_AT_LEAST_ONE_MATCH;
  }

  mme_config_unlock(&mme_config);
  return TA_LIST_NO_MATCH;
}

/* @brief compare a list of broadcasted plmns against the MME configured.
 */
static int s1ap_mme_compare_plmns(S1ap_BPLMNs_t* b_plmns) {
  int i = 0;
  int matching_occurence = 0;

  DevAssert(b_plmns != NULL);

  for (i = 0; i < b_plmns->list.count; i++) {
    if (s1ap_mme_compare_plmn(b_plmns->list.array[i]) ==
        TA_LIST_AT_LEAST_ONE_MATCH)
      matching_occurence++;
  }

  if (matching_occurence == 0)
    return TA_LIST_NO_MATCH;
  else if (matching_occurence == b_plmns->list.count - 1)
    return TA_LIST_COMPLETE_MATCH;
  else
    return TA_LIST_AT_LEAST_ONE_MATCH;
}

/* @brief compare a TAC
 */
static int s1ap_mme_compare_tac(const S1ap_TAC_t* const tac) {
  int i = 0;
  uint16_t tac_value = 0;

  DevAssert(tac != NULL);
  OCTET_STRING_TO_TAC(tac, tac_value);
  mme_config_read_lock(&mme_config);

  for (i = 0; i < mme_config.served_tai.nb_tai; i++) {
    OAILOG_TRACE(LOG_S1AP, "Comparing config tac %d, received tac = %d\n",
                 mme_config.served_tai.tac[i], tac_value);

    if (mme_config.served_tai.tac[i] == tac_value)
      return TA_LIST_AT_LEAST_ONE_MATCH;
  }

  mme_config_unlock(&mme_config);
  return TA_LIST_NO_MATCH;
}

/* @brief compare a given ta list against the one provided by mme configuration.
   @param ta_list
   @return - TA_LIST_UNKNOWN_PLMN if at least one TAC match and no PLMN match
           - TA_LIST_UNKNOWN_TAC if at least one PLMN match and no TAC match
           - TA_LIST_RET_OK if both tac and plmn match at least one element
*/
int s1ap_mme_compare_ta_lists(S1ap_SupportedTAs_t* ta_list) {
  int i;
  int tac_ret, bplmn_ret;
  int tac_matches = 0, bplmn_matches = 0;

  DevAssert(ta_list != NULL);

  /*
   * Parse every item in the list and try to find matching parameters
   */
  for (i = 0; i < ta_list->list.count; i++) {
    S1ap_SupportedTAs_Item_t* ta;

    ta = ta_list->list.array[i];
    DevAssert(ta != NULL);
    tac_ret = s1ap_mme_compare_tac(&ta->tAC);
    bplmn_ret = s1ap_mme_compare_plmns(&ta->broadcastPLMNs);

    if (tac_ret == TA_LIST_NO_MATCH && bplmn_ret == TA_LIST_NO_MATCH) {
      continue;
    }

    if (tac_ret > TA_LIST_NO_MATCH && bplmn_ret == TA_LIST_NO_MATCH) {
      tac_matches++;
    } else if (tac_ret == TA_LIST_NO_MATCH && bplmn_ret > TA_LIST_NO_MATCH) {
      bplmn_matches++;
    } else if (tac_ret > TA_LIST_NO_MATCH && bplmn_ret > TA_LIST_NO_MATCH) {
      return TA_LIST_RET_OK;
    }
  }

  if (tac_matches > 0) {
    return TA_LIST_UNKNOWN_PLMN;
  } else if (bplmn_matches > 0) {
    return TA_LIST_UNKNOWN_TAC;
  }

  return TA_LIST_UNKNOWN_TAC + TA_LIST_UNKNOWN_PLMN;
}

/* @brief compare PLMNs
 */
static int s1ap_paging_compare_plmns(
    const magma::lte::oai::SupportedTaiItems& enb_tai_item,
    uint8_t enb_plmn_count, const paging_tai_list_t* p_tai_list) {
  int plmn_idx, p_plmn_idx;

  if (!(enb_tai_item.bplmns_size())) {
    OAILOG_ERROR(LOG_S1AP, "PLMN Information not found in eNB tai list\n");
    return false;
  }
  for (plmn_idx = 0; plmn_idx < enb_plmn_count; plmn_idx++) {
    plmn_t enb_plmn;
    char plmn_array[6] = {0};
    memcpy(plmn_array, enb_tai_item.bplmns(plmn_idx).c_str(),
           sizeof(plmn_array));
    COPY_PLMN_FROM_CHAR_ARRAY_FMT(enb_plmn, plmn_array);

    for (p_plmn_idx = 0; p_plmn_idx < (p_tai_list->numoftac + 1);
         p_plmn_idx++) {
      tai_t p_plmn;
      p_plmn = p_tai_list->tai_list[p_plmn_idx];

      if (IS_PLMN_EQUAL(enb_plmn, p_plmn.plmn)) {
        return true;
      }
    }
  }
  return false;
}

/* @brief compare a TAC
 */
static int s1ap_paging_compare_tac(uint8_t enb_tac,
                                   const paging_tai_list_t* p_tai_list) {
  for (int p_tac_count = 0; p_tac_count < (p_tai_list->numoftac + 1);
       p_tac_count++) {
    if (enb_tac == p_tai_list->tai_list[p_tac_count].tac) {
      return true;
    }
  }
  return false;
}

/* @brief compare given tai list against the one stored in eNB structure.
   @param ta_list, paging_request, p_tai_list_count
   @return - tai_matching=0 if both TAC and PLMN does not match with list of
   ENBs
           - tai_matching=1 if both TAC and PLMN matches with list of ENBs
*/
int s1ap_paging_compare_ta_lists(magma::lte::oai::SupportedTaList& enb_ta_list,
                                 const paging_tai_list_t* p_tai_list,
                                 uint8_t p_tai_list_count) {
  bool tac_ret = false, bplmn_ret = false;
  int enb_tai_count, p_list_count;

  for (enb_tai_count = 0; enb_tai_count < enb_ta_list.list_count();
       enb_tai_count++) {
    if (!(enb_ta_list.supported_tai_items_size())) {
      OAILOG_ERROR(LOG_S1AP, "TAI Item not found in eNB TA List\n");
      return false;
    }
    magma::lte::oai::SupportedTaiItems enb_tai_item =
        enb_ta_list.supported_tai_items(enb_tai_count);
    for (p_list_count = 0; p_list_count < p_tai_list_count; p_list_count++) {
      const paging_tai_list_t* tai = NULL;
      tai = &p_tai_list[p_list_count];
      if (tai == NULL) {
        OAILOG_ERROR(LOG_S1AP, "Paging TAI list not found\n");
        return false;
      }
      tac_ret = s1ap_paging_compare_tac(enb_tai_item.tac(), tai);
      if (tac_ret != true) {
        return false;
      } else {
        bplmn_ret = s1ap_paging_compare_plmns(enb_tai_item,
                                              enb_tai_item.bplmns_size(), tai);
      }
      // Returns TRUE only if both TAC and PLMN matches
      if (tac_ret && bplmn_ret) {
        return true;
      }
    }
  }
  return false;
}
