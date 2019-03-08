/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under 
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.  
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include <stdio.h>
#include <stdint.h>

#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "mme_config.h"
#include "s1ap_mme_ta.h"
#include "S1ap-BPLMNs.h"
#include "S1ap-PLMNidentity.h"
#include "S1ap-SupportedTAs-Item.h"
#include "S1ap-TAC.h"

static int s1ap_mme_compare_plmn(const S1ap_PLMNidentity_t *const plmn)
{
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
      mme_config.served_tai.plmn_mcc[i],
      mcc,
      mme_config.served_tai.plmn_mnc[i],
      mnc,
      mme_config.served_tai.plmn_mnc_len[i],
      mnc_len);

    if (
      (mme_config.served_tai.plmn_mcc[i] == mcc) &&
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
static int s1ap_mme_compare_plmns(S1ap_BPLMNs_t *b_plmns)
{
  int i = 0;
  int matching_occurence = 0;

  DevAssert(b_plmns != NULL);

  for (i = 0; i < b_plmns->list.count; i++) {
    if (
      s1ap_mme_compare_plmn(b_plmns->list.array[i]) ==
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
static int s1ap_mme_compare_tac(const S1ap_TAC_t *const tac)
{
  int i = 0;
  uint16_t tac_value = 0;

  DevAssert(tac != NULL);
  OCTET_STRING_TO_TAC(tac, tac_value);
  mme_config_read_lock(&mme_config);

  for (i = 0; i < mme_config.served_tai.nb_tai; i++) {
    OAILOG_TRACE(
      LOG_S1AP,
      "Comparing config tac %d, received tac = %d\n",
      mme_config.served_tai.tac[i],
      tac_value);

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
int s1ap_mme_compare_ta_lists(S1ap_SupportedTAs_t *ta_list)
{
  int i;
  int tac_ret, bplmn_ret;

  DevAssert(ta_list != NULL);

  /*
   * Parse every item in the list and try to find matching parameters
   */
  for (i = 0; i < ta_list->list.count; i++) {
    S1ap_SupportedTAs_Item_t *ta;

    ta = ta_list->list.array[i];
    DevAssert(ta != NULL);
    tac_ret = s1ap_mme_compare_tac(&ta->tAC);
    bplmn_ret = s1ap_mme_compare_plmns(&ta->broadcastPLMNs);

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
