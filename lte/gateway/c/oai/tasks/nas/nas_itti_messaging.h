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

/*! \file nas_itti_messaging.h
   \brief
   \author  Sebastien ROUX, Lionel GAUTHIER
   \date
   \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_ITTI_MESSAGING_SEEN
#define FILE_NAS_ITTI_MESSAGING_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "nas_message.h"
#include "nas/as_message.h"
#include "esm_proc.h"
#include "MobileStationClassmark2.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "common_types.h"

int nas_itti_dl_data_req(
  const mme_ue_s1ap_id_t ue_idP,
  bstring nas_msgP,
  nas_error_code_t transaction_status);

void nas_itti_sgsap_downlink_unitdata(
  const char *const imsi,
  uint8_t imsi_len,
  const_bstring nas_msg);

void s6a_auth_info_rsp_timer_expiry_handler(void *args);

#endif /* FILE_NAS_ITTI_MESSAGING_SEEN */
