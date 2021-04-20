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
#include "3gpp_23.003.h"
#include "3gpp_29.274.h"
#include "common_types.h"

#define S8_CREATE_SESSION_RSP(mSGpTR) (mSGpTR)->ittiMsg.s8_create_session_rsp
#define S8_DELETE_SESSION_RSP(mSGpTR) (mSGpTR)->ittiMsg.s8_delete_session_rsp

typedef struct s8_bearer_context_s {
  ebi_t eps_bearer_id;
  bearer_qos_t qos;
  fteid_t pgw_s8_up;
  uint32_t charging_id;
} s8_bearer_context_t;

typedef struct s8_create_session_response_s {
  uint8_t imsi_length;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  pdn_type_t pdn_type;
  paa_t paa;
  teid_t context_teid;  // SGW_S11_teid, created per PDN
  ebi_t eps_bearer_id;
  s8_bearer_context_t bearer_context[BEARERS_PER_UE];
  uint8_t apn_restriction_value;
  fteid_t pgw_s8_cp_teid;
  uint32_t cause;
  protocol_configuration_options_t pco;
} s8_create_session_response_t;

typedef struct s8_delete_session_response_s {
  teid_t context_teid;  // SGW_S11_teid, created per PDN
  uint32_t cause;
} s8_delete_session_response_t;
