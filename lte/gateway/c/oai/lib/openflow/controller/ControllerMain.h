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

#ifdef __cplusplus
extern "C" {
#endif

#define CONTROLLER_ADDR "127.0.0.1"
#define CONTROLLER_PORT 6654
#define NUM_WORKERS 2

int start_of_controller(bool persist_state);

int stop_of_controller(void);

int openflow_controller_add_gtp_tunnel(
  struct in_addr ue,
  struct in_addr enb,
  uint32_t i_tei,
  uint32_t o_tei,
  const char* imsi,
  struct ipv4flow_dl* flow_dl,
  uint32_t flow_precedence_dl);

int openflow_controller_del_gtp_tunnel(
    struct in_addr ue,
    uint32_t i_tei,
    struct ipv4flow_dl *flow_dl);

int openflow_controller_discard_data_on_tunnel(
  struct in_addr ue,
  uint32_t i_tei,
  struct ipv4flow_dl *flow_dl);

int openflow_controller_forward_data_on_tunnel(
  struct in_addr ue,
  uint32_t i_tei,
  struct ipv4flow_dl *flow_dl,
  uint32_t flow_precedence_dl);

int openflow_controller_add_paging_rule(struct in_addr ue_ip);

int openflow_controller_delete_paging_rule(struct in_addr ue_ip);

#ifdef __cplusplus
}
#endif
