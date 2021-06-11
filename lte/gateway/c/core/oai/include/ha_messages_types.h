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

#pragma once

#include <stdint.h>

#include "3gpp_23.003.h"

#define AGW_OFFLOAD_REQ(mSGpTR) (mSGpTR)->ittiMsg.ha_agw_offload_req

// ALL: offload all idle and connected UEs
// ANY: offload any UE
// ANY_CONNECTED: offload any one of ECM_CONNECTED UEs
// ANY_IDLE: offload any one of ECM_IDLE UEs
typedef enum { ALL, ANY, ANY_CONNECTED, ANY_IDLE } offload_type_t;

// The imsi and eNB_id fields are used as filters.
// A UE that satisfy any of these filters will be offloaded.
typedef struct ha_agw_offload_req_s {
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;

  uint32_t eNB_id;
  offload_type_t enb_offload_type;
} ha_agw_offload_req_t;
