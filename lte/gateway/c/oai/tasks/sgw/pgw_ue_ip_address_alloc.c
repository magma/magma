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

/*! \file pgw_ue_ip_address_alloc.c
 * \brief
 * \author
 * \company
 * \email:
 */

#include "spgw_state.h"
#include "log.h"
#include "sgw.h"
#include "pgw_ue_ip_address_alloc.h"

// Allocate UE address pool
void pgw_ip_address_pool_init(spgw_state_t* spgw_state) {
  struct conf_ipv4_list_elm_s* conf_ipv4_p = NULL;
  struct ipv4_list_elm_s* ipv4_p           = NULL;

  STAILQ_INIT(&spgw_state->ipv4_list_free);
  STAILQ_INIT(&spgw_state->ipv4_list_allocated);
  STAILQ_FOREACH(
      conf_ipv4_p, &spgw_config.pgw_config.ipv4_pool_list, ipv4_entries) {
    ipv4_p              = calloc(1, sizeof(struct ipv4_list_elm_s));
    ipv4_p->addr.s_addr = conf_ipv4_p->addr.s_addr;
    STAILQ_INSERT_TAIL(&spgw_state->ipv4_list_free, ipv4_p, ipv4_entries);
    OAILOG_INFO(
        LOG_SPGW_APP, "Loaded IPv4 UE address in pool: %s\n",
        inet_ntoa(conf_ipv4_p->addr));
  }

  while ((conf_ipv4_p = STAILQ_FIRST(&spgw_config.pgw_config.ipv4_pool_list))) {
    STAILQ_REMOVE_HEAD(&spgw_config.pgw_config.ipv4_pool_list, ipv4_entries);
    free_wrapper((void**) &conf_ipv4_p);
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

// Allocate UE IP address from configured pool of UE IP address
int pgw_allocate_ue_ipv4_address(
    spgw_state_t* spgw_state, struct in_addr* addr_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  struct ipv4_list_elm_s* ipv4_p = NULL;

  if (STAILQ_EMPTY(&spgw_state->ipv4_list_free)) {
    addr_p->s_addr = INADDR_ANY;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  ipv4_p = STAILQ_FIRST(&spgw_state->ipv4_list_free);
  STAILQ_REMOVE(
      &spgw_state->ipv4_list_free, ipv4_p, ipv4_list_elm_s, ipv4_entries);
  STAILQ_INSERT_TAIL(&spgw_state->ipv4_list_allocated, ipv4_p, ipv4_entries);
  addr_p->s_addr = ipv4_p->addr.s_addr;
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

// Release the allocated UE IP address
int pgw_locally_release_ue_ipv4_address(
    spgw_state_t* spgw_state, const struct in_addr* const addr_pP) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  struct ipv4_list_elm_s* ipv4_p = NULL;

  STAILQ_FOREACH(ipv4_p, &spgw_state->ipv4_list_allocated, ipv4_entries) {
    if (ipv4_p->addr.s_addr == addr_pP->s_addr) {
      STAILQ_REMOVE(
          &spgw_state->ipv4_list_allocated, ipv4_p, ipv4_list_elm_s,
          ipv4_entries);
      STAILQ_INSERT_HEAD(&spgw_state->ipv4_list_free, ipv4_p, ipv4_entries);
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
    }
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
}

