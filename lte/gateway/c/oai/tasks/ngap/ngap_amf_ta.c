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
/****************************************************************************
  Source      ngap_amf_ta.c
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include <stdio.h>
#include <stdint.h>

#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "amf_config.h"
#include "ngap_amf_ta.h"
#include "Ngap_BroadcastPLMNList.h"
#include "Ngap_BroadcastPLMNItem.h"
#include "Ngap_PLMNIdentity.h"
#include "Ngap_SupportedTAItem.h"
#include "Ngap_TAC.h"
#include "TrackingAreaIdentity.h"
#include "ngap_types.h"

char buf_plmn[3];

static int ngap_amf_compare_plmn(const Ngap_PLMNIdentity_t* const plmn) {
  int i            = 0;
  uint16_t mcc     = 0;
  uint16_t mnc     = 0;
  uint16_t mnc_len = 0;

  OAILOG_DEBUG(LOG_NGAP, " :%s, %d", plmn->buf, (int) plmn->size);
  OAILOG_DEBUG(
      LOG_NGAP, " :%02x %02x %02x ", plmn->buf[0], plmn->buf[1], plmn->buf[2]);

  memcpy(buf_plmn, plmn->buf, plmn->size);

  DevAssert(plmn != NULL);
  TBCD_TO_MCC_MNC(plmn, mcc, mnc, mnc_len);
  amf_config_read_lock(&amf_config);

  for (i = 0; i < amf_config.served_tai.nb_tai; i++) {
    OAILOG_TRACE(
        LOG_NGAP,
        "Comparing plmn_mcc %d/%d, plmn_mnc %d/%d plmn_mnc_len %d/%d\n",
        amf_config.served_tai.plmn_mcc[i], mcc,
        amf_config.served_tai.plmn_mnc[i], mnc,
        amf_config.served_tai.plmn_mnc_len[i], mnc_len);

    if ((amf_config.served_tai.plmn_mcc[i] == mcc) &&
        (amf_config.served_tai.plmn_mnc[i] == mnc) &&
        (amf_config.served_tai.plmn_mnc_len[i] == mnc_len))
      /*
       * There is a matching plmn
       */
      return TA_LIST_AT_LEAST_ONE_MATCH;
  }

  amf_config_unlock(&amf_config);
  return TA_LIST_NO_MATCH;
}

/* @brief compare a list of broadcasted plmns against the AMF configured.
 o*/
static int ngap_amf_compare_plmns(Ngap_BroadcastPLMNList_t* b_plmns) {
  int i                   = 0;
  int matching_occurrence = 0;
  DevAssert(b_plmns != NULL);

  for (i = 0; i < b_plmns->list.count; i++) {
    if (ngap_amf_compare_plmn(&b_plmns->list.array[i]->pLMNIdentity) ==
        TA_LIST_AT_LEAST_ONE_MATCH)
      matching_occurrence++;
    // TBD will work on match case
  }

  if (matching_occurrence == 0)
    return TA_LIST_NO_MATCH;
  else if (matching_occurrence == b_plmns->list.count - 1)
    return TA_LIST_COMPLETE_MATCH;
  else
    return TA_LIST_AT_LEAST_ONE_MATCH;
}

/* @brief compare a TAC
 */
static int ngap_amf_compare_tac(const Ngap_TAC_t* tac) {
  int i              = 0;
  uint16_t tac_value = 0;

  DevAssert(tac != NULL);
  OCTET_STRING_TO_TAC(tac, tac_value);
  amf_config_read_lock(&amf_config);

  for (i = 0; i < amf_config.served_tai.nb_tai; i++) {
    OAILOG_TRACE(
        LOG_NGAP, "Comparing config tac %d, received tac = %d\n",
        amf_config.served_tai.tac[i], tac_value);

    if (amf_config.served_tai.tac[i] == tac_value)
      return TA_LIST_AT_LEAST_ONE_MATCH;
  }

  amf_config_unlock(&amf_config);
  return TA_LIST_NO_MATCH;
}

/* @brief compare a given ta list against the one provided by amf configuration.
   @param ta_list
   @return - TA_LIST_UNKNOWN_PLMN if at least one TAC match and no PLMN match
           - TA_LIST_UNKNOWN_TAC if at least one PLMN match and no TAC match
           - TA_LIST_RET_OK if both tac and plmn match at least one element
*/
int ngap_amf_compare_ta_lists(Ngap_SupportedTAList_t* ta_list) {
  int i;
  int tac_ret, bplmn_ret;

  DevAssert(ta_list != NULL);

  /*
   * Parse every item in the list and try to find matching parameters
   */
  for (i = 0; i < ta_list->list.count; i++) {
    Ngap_SupportedTAItem_t* ta;

    ta = ta_list->list.array[i];
    DevAssert(ta != NULL);
    tac_ret   = ngap_amf_compare_tac(&ta->tAC);
    bplmn_ret = ngap_amf_compare_plmns(&ta->broadcastPLMNList);

    if (tac_ret == TA_LIST_NO_MATCH && bplmn_ret == TA_LIST_NO_MATCH) {
      return TA_LIST_UNKNOWN_PLMN + TA_LIST_UNKNOWN_TAC;
    } else {
      if (tac_ret > TA_LIST_NO_MATCH && bplmn_ret == TA_LIST_NO_MATCH) {
        return TA_LIST_UNKNOWN_PLMN;
      } else if (tac_ret == TA_LIST_NO_MATCH && bplmn_ret > TA_LIST_NO_MATCH) {
        return TA_LIST_UNKNOWN_TAC;
      }
    }
  }

  return TA_LIST_RET_OK;
}

/* @brief compare PLMNs
 */
static int ngap_paging_compare_plmns(
    plmn_t* gnb_bplmns, uint8_t gnb_plmn_count,
    const paging_tai_list_t* p_tai_list) {
  int plmn_idx, p_plmn_idx;

  for (plmn_idx = 0; plmn_idx < gnb_plmn_count; plmn_idx++) {
    plmn_t* gnb_plmn = NULL;
    gnb_plmn         = &gnb_bplmns[plmn_idx];
    if (gnb_plmn == NULL) {
      OAILOG_ERROR(LOG_NGAP, "PLMN Information not found in eNB tai list\n");
      return false;
    }

    for (p_plmn_idx = 0; p_plmn_idx < (p_tai_list->numoftac + 1);
         p_plmn_idx++) {
      tai_t p_plmn;
      p_plmn = p_tai_list->tai_list[p_plmn_idx];

      if ((gnb_plmn->mcc_digit1 == p_plmn.plmn.mcc_digit1) &&
          (gnb_plmn->mcc_digit2 == p_plmn.plmn.mcc_digit2) &&
          (gnb_plmn->mcc_digit3 == p_plmn.plmn.mcc_digit3) &&
          (gnb_plmn->mnc_digit1 == p_plmn.plmn.mnc_digit1) &&
          (gnb_plmn->mnc_digit2 == p_plmn.plmn.mnc_digit2) &&
          (gnb_plmn->mnc_digit3 == p_plmn.plmn.mnc_digit3)) {
        return true;
      }
    }
  }
  return false;
}

/* @brief compare a TAC
 */
static int ngap_paging_compare_tac(
    uint8_t gnb_tac, const paging_tai_list_t* p_tai_list) {
  for (int p_tac_count = 0; p_tac_count < (p_tai_list->numoftac + 1);
       p_tac_count++) {
    if (gnb_tac == p_tai_list->tai_list[p_tac_count].tac) {
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
int ngap_paging_compare_ta_lists(
    m5g_supported_ta_list_t* gnb_ta_list, const paging_tai_list_t* p_tai_list,
    uint8_t p_tai_list_count) {
  bool tac_ret = false, bplmn_ret = false;
  int gnb_tai_count, p_list_count;

  for (gnb_tai_count = 0; gnb_tai_count < gnb_ta_list->list_count;
       gnb_tai_count++) {
    m5g_supported_tai_items_t* gnb_tai_item = NULL;
    gnb_tai_item = &gnb_ta_list->supported_tai_items[gnb_tai_count];
    if (gnb_tai_item == NULL) {
      OAILOG_ERROR(LOG_NGAP, "TAI Item not found in eNB TA List\n");
      return false;
    }
    for (p_list_count = 0; p_list_count < p_tai_list_count; p_list_count++) {
      const paging_tai_list_t* tai = NULL;
      tai                          = &p_tai_list[p_list_count];
      if (tai == NULL) {
        OAILOG_ERROR(LOG_NGAP, "Paging TAI list not found\n");
        return false;
      }

      tac_ret = ngap_paging_compare_tac(gnb_tai_item->tac, tai);
      if (tac_ret != true) {
        return false;
      } else {
        bplmn_ret = ngap_paging_compare_plmns(
            gnb_tai_item->bplmns, gnb_tai_item->bplmnlist_count, tai);
      }
      // Returns TRUE only if both TAC and PLMN matches
      if (tac_ret && bplmn_ret) {
        return true;
      }
    }
  }
  return false;
}
