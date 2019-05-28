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

/*! \file sctp_primitives_server.c
    \brief Main server primitives
    \author Sebastien ROUX, Lionel GAUTHIER
    \date 2013
    \version 1.0
    @ingroup _sctp
*/

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/sctp.h>
#include <arpa/inet.h>
#include <sys/select.h>

#include "dynamic_memory_check.h"
#include "common_defs.h"
#include "assertions.h"
#include "log.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "sctp_primitives_server.h"
#include "sctp_common.h"
#include "sctp_itti_messaging.h"
#include "service303.h"
#include "bstrlib.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "mme_default_values.h"
#include "sctp_messages_types.h"

#define SCTP_RC_ERROR -1
#define SCTP_RC_NORMAL_READ 0
#define SCTP_RC_DISCONNECT 1

typedef struct sctp_association_s {
  struct sctp_association_s *next_assoc; ///< Next association in the list
  struct sctp_association_s
    *previous_assoc; ///< Previous association in the list
  int sd;            ///< Socket descriptor
  uint32_t ppid;     ///< Payload protocol Identifier
  uint16_t
    instreams; ///< Number of input streams negociated for this connection
  uint16_t
    outstreams; ///< Number of output strams negotiated for this connection
  sctp_assoc_id_t assoc_id; ///< SCTP association id for the connection
  uint32_t messages_recv;   ///< Number of messages received on this connection
  uint32_t messages_sent;   ///< Number of messages sent on this connection

  struct sockaddr *peer_addresses; ///< A list of peer addresses
  int nb_peer_addresses;
} sctp_association_t;

typedef struct sctp_descriptor_s {
  // List of connected peers
  struct sctp_association_s *available_connections_head;
  struct sctp_association_s *available_connections_tail;

  uint32_t number_of_connections;
  uint16_t nb_instreams;
  uint16_t nb_outstreams;
} sctp_descriptor_t;

typedef struct sctp_arg_s {
  int sd;
  uint32_t ppid;
} sctp_arg_t;

static sctp_descriptor_t sctp_desc;

// Thread used to handle sctp messages
static pthread_t assoc_thread;

// LOCAL FUNCTIONS prototypes
void *sctp_receiver_thread(void *args_p);
static int sctp_send_msg(
  sctp_assoc_id_t sctp_assoc_id,
  uint16_t stream,
  STOLEN_REF bstring *payload);

// Association list related local functions prototypes
static sctp_association_t *sctp_is_assoc_in_list(sctp_assoc_id_t assoc_id);
static sctp_association_t *sctp_add_new_peer(void);
static int handle_assoc_change(
  int sd,
  uint32_t ppid,
  struct sctp_assoc_change *assoc_change);
static int sctp_handle_com_down(sctp_assoc_id_t assoc_id);
static int sctp_handle_reset(const sctp_assoc_id_t assoc_id);
static void sctp_dump_list(void);
static void sctp_exit(void);

//------------------------------------------------------------------------------
static sctp_association_t *sctp_add_new_peer(void)
{
  sctp_association_t *new_sctp_descriptor =
    calloc(1, sizeof(sctp_association_t));

  if (new_sctp_descriptor == NULL) {
    OAILOG_ERROR(
      LOG_SCTP,
      "Failed to allocate memory for new peer (%s:%d)\n",
      __FILE__,
      __LINE__);
    return NULL;
  }

  new_sctp_descriptor->next_assoc = NULL;
  new_sctp_descriptor->previous_assoc = NULL;

  if (sctp_desc.available_connections_tail == NULL) {
    sctp_desc.available_connections_head = new_sctp_descriptor;
    sctp_desc.available_connections_tail = sctp_desc.available_connections_head;
  } else {
    new_sctp_descriptor->previous_assoc = sctp_desc.available_connections_tail;
    sctp_desc.available_connections_tail->next_assoc = new_sctp_descriptor;
    sctp_desc.available_connections_tail = new_sctp_descriptor;
  }

  sctp_desc.number_of_connections++;
  sctp_dump_list();
  return new_sctp_descriptor;
}

//------------------------------------------------------------------------------
static sctp_association_t *sctp_is_assoc_in_list(sctp_assoc_id_t assoc_id)
{
  sctp_association_t *assoc_desc = NULL;

  if (assoc_id < 0) {
    return NULL;
  }

  for (assoc_desc = sctp_desc.available_connections_head; assoc_desc;
       assoc_desc = assoc_desc->next_assoc) {
    if (assoc_desc->assoc_id == assoc_id) {
      break;
    }
  }

  return assoc_desc;
}

//------------------------------------------------------------------------------
static int sctp_remove_assoc_from_list(sctp_assoc_id_t assoc_id)
{
  sctp_association_t *assoc_desc = NULL;

  /*
   * Association not in the list
   */
  if ((assoc_desc = sctp_is_assoc_in_list(assoc_id)) == NULL) {
    return -1;
  }

  if (assoc_desc->next_assoc == NULL) {
    if (assoc_desc->previous_assoc == NULL) {
      /*
       * Head and tail
       */
      sctp_desc.available_connections_head =
        sctp_desc.available_connections_tail = NULL;
    } else {
      /*
       * Not head but tail
       */
      sctp_desc.available_connections_tail = assoc_desc->previous_assoc;
      assoc_desc->previous_assoc->next_assoc = NULL;
    }
  } else {
    if (assoc_desc->previous_assoc == NULL) {
      /*
       * Head but not tail
       */
      sctp_desc.available_connections_head = assoc_desc->next_assoc;
      assoc_desc->next_assoc->previous_assoc = NULL;
    } else {
      /*
       * Not head and not tail
       */
      assoc_desc->previous_assoc->next_assoc = assoc_desc->next_assoc;
      assoc_desc->next_assoc->previous_assoc = assoc_desc->previous_assoc;
    }
  }

  if (assoc_desc->peer_addresses) {
    int rv = sctp_freepaddrs(assoc_desc->peer_addresses);
    if (rv)
      OAILOG_DEBUG(
        LOG_SCTP, "sctp_freepaddrs(%p) failed\n", assoc_desc->peer_addresses);
  }
  free_wrapper((void **) &assoc_desc);
  sctp_desc.number_of_connections--;
  return 0;
}

//------------------------------------------------------------------------------
static void sctp_dump_assoc(sctp_association_t *sctp_assoc_p)
{
#if SCTP_DUMP_LIST
  int i;

  if (sctp_assoc_p == NULL) {
    return;
  }

  OAILOG_DEBUG(LOG_SCTP, "sd           : %d\n", sctp_assoc_p->sd);
  OAILOG_DEBUG(LOG_SCTP, "input streams: %d\n", sctp_assoc_p->instreams);
  OAILOG_DEBUG(LOG_SCTP, "out streams  : %d\n", sctp_assoc_p->outstreams);
  OAILOG_DEBUG(LOG_SCTP, "assoc_id     : %d\n", sctp_assoc_p->assoc_id);
  OAILOG_DEBUG(LOG_SCTP, "peer address :\n");

  for (i = 0; i < sctp_assoc_p->nb_peer_addresses; i++) {
    char address[40];

    memset(address, 0, sizeof(address));

    if (
      inet_ntop(
        sctp_assoc_p->peer_addresses[i].sa_family,
        sctp_assoc_p->peer_addresses[i].sa_data,
        address,
        sizeof(address)) != NULL) {
      OAILOG_DEBUG(LOG_SCTP, "    - [%s]\n", address);
    }
  }

#else
  sctp_assoc_p = sctp_assoc_p;
#endif
}

//------------------------------------------------------------------------------
static void sctp_dump_list(void)
{
#if SCTP_DUMP_LIST
  sctp_association_t *sctp_assoc_p = sctp_desc.available_connections_head;
  OAILOG_DEBUG(
    LOG_SCTP,
    "SCTP list contains %d associations\n",
    sctp_desc.number_of_connections);

  while (sctp_assoc_p != NULL) {
    sctp_dump_assoc(sctp_assoc_p);
    sctp_assoc_p = sctp_assoc_p->next_assoc;
  }

#else
  sctp_dump_assoc(NULL);
#endif
}

//------------------------------------------------------------------------------
static int sctp_send_msg(
  sctp_assoc_id_t sctp_assoc_id,
  uint16_t stream,
  STOLEN_REF bstring *payload)
{
  sctp_association_t *assoc_desc = NULL;

  DevAssert(*payload);

  if ((assoc_desc = sctp_is_assoc_in_list(sctp_assoc_id)) == NULL) {
    OAILOG_DEBUG(
      LOG_SCTP,
      "This assoc id has not been fount in list (%d)\n",
      sctp_assoc_id);
    return -1;
  }

  if (assoc_desc->sd == -1) {
    /*
     * The socket is invalid may be closed.
     */
    OAILOG_DEBUG(
      LOG_SCTP,
      "The socket is invalid may be closed (assoc id %d)\n",
      sctp_assoc_id);
    return -1;
  }

  OAILOG_DEBUG(
    LOG_SCTP,
    "[%d][%d] Sending buffer %p of %d bytes on stream %d with ppid %d\n",
    assoc_desc->sd,
    sctp_assoc_id,
    bdata(*payload),
    blength(*payload),
    stream,
    assoc_desc->ppid);

  /*
   * Send message_p on specified stream of the sd association
   */
  if (
    sctp_sendmsg(
      assoc_desc->sd,
      (const void *) bdata(*payload),
      (size_t) blength(*payload),
      NULL,
      0,
      htonl(assoc_desc->ppid),
      0,
      stream,
      0,
      0) < 0) {
    bdestroy(*payload);
    OAILOG_ERROR(LOG_SCTP, "send: %s:%d\n", strerror(errno), errno);
    return -1;
  }
  OAILOG_DEBUG(
    LOG_SCTP,
    "Successfully sent %d bytes on stream %d\n",
    blength(*payload),
    stream);
  bdestroy(*payload);
  *payload = NULL;

  assoc_desc->messages_sent++;
  return 0;
}

//------------------------------------------------------------------------------
static int sctp_create_new_listener(sctp_init_t *init_p)
{
  struct sctp_event_subscribe event = {0};
  struct sockaddr *addr = NULL;
  sctp_arg_t *sctp_arg_p = NULL;
  uint16_t i = 0, j = 0;
  int sd = 0;
  int used_addresses = 0;

  DevAssert(init_p != NULL);

  if (init_p->ipv4 == 0 && init_p->ipv6 == 0) {
    OAILOG_ERROR(
      LOG_SCTP,
      "Illegal IP configuration upper layer should request at"
      "least ipv4 and/or ipv6 config\n");
    return -1;
  }

  if ((used_addresses = init_p->nb_ipv4_addr + init_p->nb_ipv6_addr) == 0) {
    OAILOG_WARNING(LOG_SCTP, "No address provided...\n");
    return -1;
  }

  addr = calloc((size_t) used_addresses, sizeof(struct sockaddr));
  OAILOG_DEBUG(
    LOG_SCTP, "Creating new listen socket on port %u with\n", init_p->port);

  if (init_p->ipv4 == 1) {
    struct sockaddr_in *ip4_addr;

    OAILOG_DEBUG(LOG_SCTP, "ipv4 addresses:\n");

    for (i = 0; i < init_p->nb_ipv4_addr; i++) {
      ip4_addr = (struct sockaddr_in *) &addr[i];
      ip4_addr->sin_family = AF_INET;
      ip4_addr->sin_port = htons(init_p->port);
      ip4_addr->sin_addr.s_addr = init_p->ipv4_address[i].s_addr;
      char ipv4[INET_ADDRSTRLEN];
      inet_ntop(
        AF_INET, (void *) &ip4_addr->sin_addr.s_addr, ipv4, INET_ADDRSTRLEN);
      OAILOG_DEBUG(LOG_SCTP, "\t- %s\n", ipv4);
    }
  }

  if (init_p->ipv6 == 1) {
    struct sockaddr_in6 *ip6_addr = NULL;

    OAILOG_DEBUG(LOG_SCTP, "ipv6 addresses:\n");

    for (j = 0; j < init_p->nb_ipv6_addr; j++) {
      char ipv6[INET6_ADDRSTRLEN];
      inet_ntop(
        AF_INET6, (void *) &init_p->ipv6_address[j], ipv6, INET6_ADDRSTRLEN);
      OAILOG_DEBUG(LOG_SCTP, "\t- %s\n", ipv6);
      ip6_addr = (struct sockaddr_in6 *) &addr[i + j];
      ip6_addr->sin6_family = AF_INET6;
      ip6_addr->sin6_port = htons(init_p->port);

      ip6_addr->sin6_addr = init_p->ipv6_address[j];
    }
  }

  if ((sd = socket(AF_INET6, SOCK_STREAM, IPPROTO_SCTP)) < 0) {
    OAILOG_ERROR(LOG_SCTP, "socket: %s:%d\n", strerror(errno), errno);
    return -1;
  }

  memset((void *) &event, 0, sizeof(struct sctp_event_subscribe));
  event.sctp_association_event = 1;
  event.sctp_shutdown_event = 1;
  event.sctp_data_io_event = 1;

  if (
    setsockopt(
      sd,
      IPPROTO_SCTP,
      SCTP_EVENTS,
      &event,
      sizeof(struct sctp_event_subscribe)) < 0) {
    OAILOG_ERROR(LOG_SCTP, "setsockopt: %s:%d\n", strerror(errno), errno);
    return -1;
  }

  /*
   * Some pre-bind socket configuration
   */
  if (
    sctp_set_init_opt(
      sd, sctp_desc.nb_instreams, sctp_desc.nb_outstreams, 0, 0) < 0) {
    goto err;
  }

  if (sctp_bindx(sd, addr, used_addresses, SCTP_BINDX_ADD_ADDR) != 0) {
    OAILOG_ERROR(LOG_SCTP, "sctp_bindx: %s:%d\n", strerror(errno), errno);
    goto err;
  }

  if (listen(sd, 5) < 0) {
    OAILOG_ERROR(LOG_SCTP, "listen: %s:%d\n", strerror(errno), errno);
    goto err;
  }

  if ((sctp_arg_p = malloc(sizeof(sctp_arg_t))) == NULL) {
    goto err;
  }

  sctp_arg_p->sd = sd;
  sctp_arg_p->ppid = init_p->ppid;

  if (
    pthread_create(
      &assoc_thread, NULL, &sctp_receiver_thread, (void *) sctp_arg_p) < 0) {
    OAILOG_ERROR(LOG_SCTP, "pthread_create: %s:%d\n", strerror(errno), errno);
    return -1;
  }

  free_wrapper((void **) &addr);
  return sd;
err:

  if (addr) {
    free_wrapper((void **) &addr);
  }

  if (sd != -1) {
    close(sd);
    sd = -1;
  }

  return -1;
}

//------------------------------------------------------------------------------
static inline int sctp_read_from_socket(int sd, uint32_t ppid)
{
  int flags = 0, n;
  socklen_t from_len = 0;
  struct sctp_sndrcvinfo sinfo = {0};
  struct sockaddr_in6 addr = {0};
  uint8_t buffer[SCTP_RECV_BUFFER_SIZE];

  if (sd < 0) {
    return -1;
  }

  memset((void *) &addr, 0, sizeof(struct sockaddr_in6));
  from_len = (socklen_t) sizeof(struct sockaddr_in6);
  memset((void *) &sinfo, 0, sizeof(struct sctp_sndrcvinfo));
  n = sctp_recvmsg(
    sd,
    (void *) buffer,
    SCTP_RECV_BUFFER_SIZE,
    (struct sockaddr *) &addr,
    &from_len,
    &sinfo,
    &flags);

  if (n < 0) {
    OAILOG_DEBUG(LOG_SCTP, "An error occured during read\n");
    OAILOG_ERROR(LOG_SCTP, "sctp_recvmsg: %s:%d\n", strerror(errno), errno);
    return SCTP_RC_ERROR;
  }

  if (flags & MSG_NOTIFICATION) {
    union sctp_notification *snp = (union sctp_notification *) buffer;

    switch (snp->sn_header.sn_type) {
      case SCTP_SHUTDOWN_EVENT: {
        OAILOG_DEBUG(LOG_SCTP, "SCTP_SHUTDOWN_EVENT received\n");
        return sctp_handle_com_down(
          (sctp_assoc_id_t) snp->sn_shutdown_event.sse_assoc_id);
      }
      case SCTP_ASSOC_CHANGE: {
        OAILOG_DEBUG(LOG_SCTP, "SCTP association change event received\n");
        return handle_assoc_change(sd, ppid, &snp->sn_assoc_change);
      }
      default: {
        OAILOG_WARNING(
          LOG_SCTP, "Unhandled notification type %u\n", snp->sn_header.sn_type);
        break;
      }
    }
  } else {
    /*
     * Data payload received
     */
    sctp_association_t *association;

    if (
      (association = sctp_is_assoc_in_list(
         (sctp_assoc_id_t) sinfo.sinfo_assoc_id)) == NULL) {
      // TODO: handle this case
      return SCTP_RC_ERROR;
    }

    association->messages_recv++;

    if (ntohl(sinfo.sinfo_ppid) != association->ppid) {
      /*
       * Mismatch in Payload Protocol Identifier,
       * * * * may be we received unsollicited traffic from stack other than S1AP.
       */
      OAILOG_ERROR(
        LOG_SCTP,
        "Received data from peer with unsollicited PPID %d, expecting %d\n",
        ntohl(sinfo.sinfo_ppid),
        association->ppid);
      return SCTP_RC_ERROR;
    }

    OAILOG_DEBUG(
      LOG_SCTP,
      "[%d][%d] Msg of length %d received from port %u, on stream %d, PPID "
      "%d\n",
      sinfo.sinfo_assoc_id,
      sd,
      n,
      ntohs(addr.sin6_port),
      sinfo.sinfo_stream,
      ntohl(sinfo.sinfo_ppid));
    bstring payload = blk2bstr(buffer, n);
    sctp_itti_send_new_message_ind(
      &payload,
      (sctp_assoc_id_t) sinfo.sinfo_assoc_id,
      sinfo.sinfo_stream,
      association->instreams,
      association->outstreams);
  }

  return SCTP_RC_NORMAL_READ;
}

//------------------------------------------------------------------------------
static int sctp_handle_com_down(sctp_assoc_id_t assoc_id)
{
  OAILOG_DEBUG(
    LOG_SCTP, "Sending close connection for assoc_id %u\n", assoc_id);

  if (sctp_itti_send_com_down_ind(assoc_id, false) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to send message to TASK_S1AP\n");
  }

  if (sctp_remove_assoc_from_list(assoc_id) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to find client in list\n");
  }

  return SCTP_RC_DISCONNECT;
}

static int sctp_handle_reset(const sctp_assoc_id_t assoc_id)
{
  OAILOG_DEBUG(LOG_SCTP, "Handling sctp reset\n");

  if (sctp_itti_send_com_down_ind(assoc_id, true) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to send release message to TASK_S1AP\n");
    return SCTP_RC_ERROR;
  }
  sctp_association_t *assoc = sctp_is_assoc_in_list(assoc_id);
  DevAssert(assoc != NULL);

  return SCTP_RC_NORMAL_READ;
}

//------------------------------------------------------------------------------
void *sctp_receiver_thread(void *args_p)
{
  sctp_arg_t sctp_arg_p;

  /*
   * maximum file descriptor number
   */
  int fdmax, clientsock, i;

  /*
   * master file descriptor list
   */
  fd_set master;

  /*
   * temp file descriptor list for select()
   */
  fd_set read_fds;

  if (args_p == NULL) {
    pthread_exit(NULL);
  }

  memcpy(&sctp_arg_p, args_p, sizeof sctp_arg_p);
  free_wrapper(&args_p);

  /*
   * clear the master and temp sets
   */
  FD_ZERO(&master);
  FD_ZERO(&read_fds);
  FD_SET(sctp_arg_p.sd, &master);
  fdmax = sctp_arg_p.sd; /* so far, it's this one */

  while (1) {
    memcpy(&read_fds, &master, sizeof(master));

    if (select(fdmax + 1, &read_fds, NULL, NULL, NULL) == -1) {
      OAILOG_ERROR(
        LOG_SCTP, "[%d] Select() error: %s\n", sctp_arg_p.sd, strerror(errno));
      free_wrapper((void **) &args_p);
      close(sctp_arg_p.sd);
      args_p = NULL;
      pthread_exit(NULL);
    }

    for (i = 0; i <= fdmax; i++) {
      if (FD_ISSET(i, &read_fds)) {
        if (i == sctp_arg_p.sd) {
          /*
           * There is data to read on listener socket. This means we have to accept
           * * * * the connection.
           */
          if ((clientsock = accept(sctp_arg_p.sd, NULL, NULL)) < 0) {
            OAILOG_ERROR(
              LOG_SCTP,
              "[%d] accept: %s:%d\n",
              sctp_arg_p.sd,
              strerror(errno),
              errno);
            free_wrapper((void **) &args_p);
            close(sctp_arg_p.sd);
            args_p = NULL;
            pthread_exit(NULL);
          } else {
            FD_SET(clientsock, &master); /* add to master set */

            if (clientsock > fdmax) {
              /*
               * keep track of the maximum
               */
              fdmax = clientsock;
            }
          }
        } else {
          int ret;

          /*
           * Read from socket
           */
          ret = sctp_read_from_socket(i, sctp_arg_p.ppid);

          /*
           * When the socket is disconnected we have to update
           * * * * the fd_set.
           */
          if (ret == SCTP_RC_DISCONNECT) {
            /*
             * Remove the socket from the FD set and update the max sd
             */
            FD_CLR(i, &master);

            if (i == fdmax) {
              while (FD_ISSET(fdmax, &master) == false) fdmax -= 1;
            }
          }
        }
      }
    }
  }

  return NULL;
}

//------------------------------------------------------------------------------
static void *sctp_intertask_interface(__attribute__((unused)) void *args_p)
{
  int sctp_sd = -1;
  itti_mark_task_ready(TASK_SCTP);

  while (1) {
    MessageDef *received_message_p = NULL;

    itti_receive_msg(TASK_SCTP, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case SCTP_INIT_MSG: {
        OAILOG_DEBUG(LOG_SCTP, "Received SCTP_INIT_MSG\n");

        /*
         * We received a new connection request
         */
        if (
          (sctp_sd = sctp_create_new_listener(
             &received_message_p->ittiMsg.sctpInit)) < 0) {
          /*
           * SCTP socket creation or bind failed...
           * Die as this MME is not going to be useful.
           */
          AssertFatal(false, "Failed to create new SCTP listener\n");
        }
        MessageDef *message_p = NULL;
        message_p =
          itti_alloc_new_message(TASK_S1AP, SCTP_MME_SERVER_INITIALIZED);
        SCTP_MME_SERVER_INITIALIZED(message_p).successful = true;
        AssertFatal(message_p != NULL, "itti_alloc_new_message Failed");
        itti_send_msg_to_task(TASK_MME_APP, message_p);
      } break;

      case SCTP_CLOSE_ASSOCIATION: {
      } break;

      case SCTP_DATA_REQ: {
        if (
          sctp_send_msg(
            SCTP_DATA_REQ(received_message_p).assoc_id,
            SCTP_DATA_REQ(received_message_p).stream,
            &SCTP_DATA_REQ(received_message_p).payload) < 0) {
          sctp_itti_send_lower_layer_conf(
            received_message_p->ittiMsgHeader.originTaskId,
            SCTP_DATA_REQ(received_message_p).assoc_id,
            SCTP_DATA_REQ(received_message_p).stream,
            SCTP_DATA_REQ(received_message_p).mme_ue_s1ap_id,
            false);
        } /* NO NEED FOR CONFIRM success yet else {
          if (INVALID_MME_UE_S1AP_ID != SCTP_DATA_REQ (received_message_p).mme_ue_s1ap_id) {
            sctp_itti_send_lower_layer_conf(received_message_p->ittiMsgHeader.originTaskId,
                SCTP_DATA_REQ (received_message_p).assoc_id,
                SCTP_DATA_REQ (received_message_p).stream,
                SCTP_DATA_REQ (received_message_p).mme_ue_s1ap_id,
                true);
          }
        }*/
      } break;

      case MESSAGE_TEST: {
        OAI_FPRINTF_INFO("TASK_SCTP received MESSAGE_TEST\n");
      } break;

      case TERMINATE_MESSAGE: {
        close(sctp_sd);
        sctp_exit();
        itti_free_msg_content(received_message_p);
        itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
        itti_exit_task();
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_SCTP,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
// Function adds a new association and sends a new association notification message.
sctp_association_t *add_new_association(
  int sd,
  uint32_t ppid,
  struct sctp_assoc_change *sctp_assoc_changed)
{
  sctp_association_t *new_association = NULL;
  if ((new_association = sctp_add_new_peer()) == NULL) {
    OAILOG_ERROR(LOG_SCTP, "Failed to allocate new sctp peer \n");
    return NULL;
  }

  new_association->sd = sd;
  new_association->ppid = ppid;
  new_association->instreams = sctp_assoc_changed->sac_inbound_streams;
  new_association->outstreams = sctp_assoc_changed->sac_outbound_streams;
  new_association->assoc_id =
    (sctp_assoc_id_t) sctp_assoc_changed->sac_assoc_id;
  sctp_get_localaddresses(sd, NULL, NULL);
  sctp_get_peeraddresses(
    sd, &new_association->peer_addresses, &new_association->nb_peer_addresses);

  if (
    sctp_itti_send_new_association(
      new_association->assoc_id,
      new_association->instreams,
      new_association->outstreams) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to send message to S1AP\n");
    return NULL;
  }
  return new_association;
}

//------------------------------------------------------------------------------
// Handle association change events.

int handle_assoc_change(
  int sd,
  uint32_t ppid,
  struct sctp_assoc_change *sctp_assoc_changed)
{
  int rc = SCTP_RC_NORMAL_READ;
  switch (sctp_assoc_changed->sac_state) {
    case SCTP_COMM_UP: {
      if (add_new_association(sd, ppid, sctp_assoc_changed) == NULL) {
        rc = SCTP_RC_ERROR;
      }
      break;
    }
    case SCTP_RESTART: {
      DevAssert(
        sctp_is_assoc_in_list(
          (sctp_assoc_id_t) sctp_assoc_changed->sac_assoc_id) != NULL);
      /* Don't remove the sctp assoc from the list of associations, just send remove the s1ap state */
      rc =
        sctp_handle_reset((sctp_assoc_id_t) sctp_assoc_changed->sac_assoc_id);
      increment_counter("sctp_reset", 1, NO_LABELS);
      break;
    }
    case SCTP_COMM_LOST:
    case SCTP_SHUTDOWN_COMP:
    case SCTP_CANT_STR_ASSOC: {
      DevAssert(
        sctp_is_assoc_in_list(
          (sctp_assoc_id_t) sctp_assoc_changed->sac_assoc_id) != NULL);
      rc = sctp_handle_com_down(
        (sctp_assoc_id_t) sctp_assoc_changed->sac_assoc_id);
      increment_counter("sctp_shutdown", 1, NO_LABELS);
      break;
    }
    default:
      OAILOG_DEBUG(
        LOG_SCTP,
        "Logging unhandled sctp message %u\n",
        sctp_assoc_changed->sac_state);
      break;
  }
  return rc;
}

//------------------------------------------------------------------------------
int sctp_init(const mme_config_t *mme_config_p)
{
  OAILOG_DEBUG(LOG_SCTP, "Initializing SCTP task interface\n");
  memset(&sctp_desc, 0, sizeof(sctp_descriptor_t));
  /*
   * Number of streams from configuration
   */
  sctp_desc.nb_instreams = mme_config_p->sctp_config.in_streams;
  sctp_desc.nb_outstreams = mme_config_p->sctp_config.out_streams;

  if (itti_create_task(TASK_SCTP, &sctp_intertask_interface, NULL) < 0) {
    OAILOG_ERROR(LOG_SCTP, "create task failed\n");
    OAILOG_DEBUG(LOG_SCTP, "Initializing SCTP task interface: FAILED\n");
    return -1;
  }

  OAILOG_DEBUG(LOG_SCTP, "Initializing SCTP task interface: DONE\n");
  return 0;
}

//------------------------------------------------------------------------------
static void sctp_exit(void)
{
  int rv = pthread_cancel(assoc_thread);
  pthread_join(assoc_thread, NULL);
  if (rv)
    OAILOG_DEBUG(
      LOG_SCTP,
      "pthread_cancel(%08lX) failed: %d:%s\n",
      assoc_thread,
      rv,
      strerror(rv));
  ;

  sctp_association_t *sctp_assoc_p = sctp_desc.available_connections_head;
  sctp_association_t *next_sctp_assoc_p = sctp_desc.available_connections_head;

  while (next_sctp_assoc_p) {
    next_sctp_assoc_p = sctp_assoc_p->next_assoc;
    if (sctp_assoc_p->peer_addresses) {
      close(sctp_assoc_p->sd);
      rv = sctp_freepaddrs(sctp_assoc_p->peer_addresses);
      if (rv)
        OAILOG_DEBUG(
          LOG_SCTP,
          "sctp_freepaddrs(%p) failed\n",
          sctp_assoc_p->peer_addresses);
    }
    free_wrapper((void **) &sctp_assoc_p);
    sctp_desc.number_of_connections--;
  }
  OAI_FPRINTF_INFO("TASK_SCTP terminated\n");
}
