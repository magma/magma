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

#pragma once

#include <sys/types.h>
#include "messages_types.h"

/*
 * Sends a SPGW_NW_INITIATED_ACTIVATE_BEARER_REQ message to SPGW.
 */
int send_activate_bearer_request_itti(
  itti_spgw_nw_init_actv_bearer_request_t* itti_msg);
/*
 * Sends a PGW_NW_INITIATED_DEACTIVATE_BEARER_REQ message to PGW.
 */
void send_deactivate_bearer_request_itti(
  itti_pgw_nw_init_deactv_bearer_request_t *itti_msg);
