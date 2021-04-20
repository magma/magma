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

/*! \file s1ap_mme_ta.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S1AP_MME_TA_SEEN
#define FILE_S1AP_MME_TA_SEEN

#include "S1ap_SupportedTAs.h"
#include "TrackingAreaIdentity.h"
#include "s1ap_types.h"

enum {
  TA_LIST_UNKNOWN_TAC        = -2,
  TA_LIST_UNKNOWN_PLMN       = -1,
  TA_LIST_RET_OK             = 0,
  TA_LIST_NO_MATCH           = 0x1,
  TA_LIST_AT_LEAST_ONE_MATCH = 0x2,
  TA_LIST_COMPLETE_MATCH     = 0x3,
};

int s1ap_mme_compare_ta_lists(S1ap_SupportedTAs_t* ta_list);
int s1ap_paging_compare_ta_lists(
    supported_ta_list_t* enb_ta_list, const paging_tai_list_t* p_tai_list,
    uint8_t p_tai_list_count);

#endif /* FILE_S1AP_MME_TA_SEEN */
