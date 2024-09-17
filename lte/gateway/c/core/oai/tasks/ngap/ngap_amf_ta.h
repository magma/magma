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
  Source      ngap_amf_ta.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include "Ngap_SupportedTAList.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_types.h"

enum {
  TA_LIST_UNKNOWN_TAC = -2,
  TA_LIST_UNKNOWN_PLMN = -1,
  TA_LIST_RET_OK = 0,
  TA_LIST_NO_MATCH = 0x1,
  TA_LIST_AT_LEAST_ONE_MATCH = 0x2,
  TA_LIST_COMPLETE_MATCH = 0x3,
};

enum {
  NGAP_SLICE_CMP_MATCH_FAILED_SD = -101,
  NGAP_SLICE_CMP_MATCH_FAILED_SST = -100,
  NGAP_SLICE_CMP_MATCH_SUCCESS = 0
};

typedef enum s_nssai_sst_s {
  _SST_eMBB = 1,
  _SST_URLLC = 2,
  _SST_mMTC = 3,
} s_nssai_sst_t;

int ngap_amf_compare_ta_lists(Ngap_SupportedTAList_t* ta_list);
int ngap_paging_compare_ta_lists(m5g_supported_ta_list_t* enb_ta_list,
                                 const paging_tai_list_t* p_tai_list,
                                 uint8_t p_tai_list_count);
