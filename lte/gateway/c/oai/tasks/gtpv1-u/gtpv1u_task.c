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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
/*! \file gtpv1u_task.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>

#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "gtpv1u.h"
#include "gtpv1u_sgw_defs.h"
#include "pgw_ue_ip_address_alloc.h"
#include "intertask_interface_types.h"
#include "pgw_config.h"
#include "spgw_config.h"

const struct gtp_tunnel_ops* gtp_tunnel_ops;

static void* gtpv1u_thread(void* args)
{
  itti_mark_task_ready(TASK_GTPV1_U);

  gtpv1u_data_t* gtpv1u_data = (gtpv1u_data_t*) args;

  while (1) {
    /*
     * Trying to fetch a message from the message queue.
     * * * * If the queue is empty, this function will block till a
     * * * * message is sent to the task.
     */
    MessageDef* received_message_p = NULL;

    itti_receive_msg(TASK_GTPV1_U, &received_message_p);
    DevAssert(received_message_p != NULL);

    switch (ITTI_MSG_ID(received_message_p)) {
      case TERMINATE_MESSAGE: gtpv1u_exit(gtpv1u_data); break;

      default: {
        OAILOG_ERROR(
          LOG_GTPV1U,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    // TODO: Add state write in case of using libgptnl

    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int gtpv1u_init(
  spgw_state_t* spgw_state_p,
  spgw_config_t* spgw_config,
  bool persist_state)
{
  int rv = 0;
  struct in_addr netaddr;
  uint32_t netmask = 0;

  OAILOG_DEBUG(LOG_GTPV1U, "Initializing GTPV1U interface\n");

  // Init gtp_tunnel_ops
#if ENABLE_OPENFLOW
  OAILOG_DEBUG(LOG_GTPV1U, "Initializing gtp_tunnel_ops_openflow\n");
  gtp_tunnel_ops = gtp_tunnel_ops_init_openflow();
#else
  OAILOG_DEBUG(LOG_GTPV1U, "Initializing gtp_tunnel_ops_libgtpnl\n");
  gtp_tunnel_ops = gtp_tunnel_ops_init_libgtpnl();
#endif

  if (gtp_tunnel_ops == NULL) {
    OAILOG_CRITICAL(LOG_GTPV1U, "ERROR in initializing gtp_tunnel_ops\n");
    return -1;
  }

  // Reset GTP tunnel states
  rv = gtp_tunnel_ops->reset();
  if (rv != 0) {
    OAILOG_CRITICAL(LOG_GTPV1U, "ERROR clean existing gtp states.\n");
    return -1;
  }

  rv = get_ip_block(&netaddr, &netmask);
  if (rv != 0) {
    OAILOG_CRITICAL(
      LOG_GTPV1U, "ERROR in getting assigned IP block from mobilityd\n");
    return -1;
  }

  // Init GTP device, using the same MTU as SGi.
  gtp_tunnel_ops->init(
    &netaddr,
    netmask,
    spgw_config->pgw_config.ipv4.mtu_SGI,
    &spgw_state_p->sgw_state.gtpv1u_data.fd0,
    &spgw_state_p->sgw_state.gtpv1u_data.fd1u,
    persist_state);

  // END-GTP quick integration only for evaluation purpose

  if (
    itti_create_task(
      TASK_GTPV1_U, &gtpv1u_thread, &spgw_state_p->sgw_state.gtpv1u_data) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "gtpv1u phtread_create: %s", strerror(errno));
    gtp_tunnel_ops->uninit();
    return -1;
  }

  OAILOG_DEBUG(LOG_GTPV1U, "Initializing GTPV1U interface: DONE\n");
  return 0;
}

//------------------------------------------------------------------------------
void gtpv1u_exit(gtpv1u_data_t* const gtpv1u_data)
{
  gtp_tunnel_ops->uninit();
  itti_exit_task();
}
