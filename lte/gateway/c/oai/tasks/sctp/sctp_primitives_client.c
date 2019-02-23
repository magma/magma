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

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <assert.h>
#include <poll.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
//WARNING: sctp requires libsctp-dev and -lsctp as linker option
#include <netinet/sctp.h>
#include <arpa/inet.h>

#include "assertions.h"
#include "sctp_eNB_defs.h"
#include "sctp_common.h"
#include "sctp_primitives_client.h"
#include "mme_default_values.h"
#include "log.h"

/* Send buffer to SCTP association */
int sctp_send_msg(
  sctp_data_t *sctp_data_p,
  constuint16_t ppid,
  constsctp_stream_id_t stream,
  const uint8_t *buffer,
  const size_t length)
{
  DevAssert(buffer != NULL);
  DevAssert(sctp_data_p != NULL);

  /*
   * Send message on specified stream of the sd association
   * * * * NOTE: PPID should be defined in network order
   */
  if (
    sctp_sendmsg(
      sctp_data_p->sd,
      (const void *) buffer,
      length,
      NULL,
      0,
      htonl(ppid),
      0,
      stream,
      0,
      0) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Sctp_sendmsg failed: %s\n", strerror(errno));
    return -1;
  }

  OAILOG_DEBUG(
    LOG_SCTP,
    "Successfully sent %d bytes to port %d on stream %d\n",
    length,
    sctp_data_p->remote_port,
    stream);
  return 0;
}

static int sctp_handle_notifications(union sctp_notification *snp)
{
  if (SCTP_SHUTDOWN_EVENT == snp->sn_header.sn_type) {
    /*
     * Client deconnection
     */
    OAILOG_DEBUG(LOG_SCTP, "Notification received: server deconnected\n");
  } else if (SCTP_ASSOC_CHANGE == snp->sn_header.sn_type) {
    /*
     * Association has changed
     */
    OAILOG_DEBUG(
      LOG_SCTP, "Notification received: server association changed\n");
  } else {
    OAILOG_DEBUG(
      LOG_SCTP, "Notification received: %d TODO\n", snp->sn_header.sn_type);
  }

  /*
   * TODO: handle more notif here
   */
  return 0;
}

int sctp_run(sctp_data_t *sctp_data_p)
{
  int ret, repeat = 1;
  int total_size = 0;
  int sd;
  struct pollfd fds;

  DevAssert(sctp_data_p != NULL);
  sd = sctp_data_p->sd;
  memset(&fds, 0, sizeof(struct pollfd));
  fds.fd = sd;
  fds.events = POLLIN;

  while (repeat == 1) {
    ret = poll(&fds, 1, 0);

    if (ret < 0) {
      OAILOG_ERROR(
        LOG_SCTP,
        "[SD %d] Poll has failed (%d:%s)\n",
        sd,
        errno,
        strerror(errno));
      return errno;
    } else if (ret == 0) {
      /*
       * No data to read, just leave the loop
       */
      OAILOG_DEBUG(LOG_SCTP, "[SD %d] Poll: no data available\n", sd);
      repeat = 0;
    } else {
      /*
       * Socket has some data to read
       */
      uint8_t buffer[SCTP_RECV_BUFFER_SIZE];
      int flags = 0;
      int n;
      struct sockaddr_in addr;
      struct sctp_sndrcvinfo sinfo;
      socklen_t from_len;

      memset((void *) &addr, 0, sizeof(struct sockaddr_in));
      from_len = (socklen_t) sizeof(struct sockaddr_in);
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
        /*
         * Other peer is deconnected
         */
        OAILOG_ERROR(
          LOG_SCTP,
          "[SD %d] An error occured during read (%d:%s)\n",
          sd,
          errno,
          strerror(errno));
        return 0;
      }

      if (flags & MSG_NOTIFICATION) {
        sctp_handle_notifications((union sctp_notification *) buffer);
      } else {
        struct sctp_queue_item_s *new_item_p;

        new_item_p = calloc(1, sizeof(struct sctp_queue_item_s));
        /*
         * Normal data received
         */
        OAILOG_DEBUG(
          LOG_SCTP,
          "[SD %d] Msg of length %d received from %s:%u on "
          "stream %d, PPID %d, assoc_id %d\n",
          sd,
          n,
          inet_ntoa(addr.sin_addr),
          ntohs(addr.sin_port),
          sinfo.sinfo_stream,
          sinfo.sinfo_ppid,
          sinfo.sinfo_assoc_id);
        new_item_p->local_stream = sinfo.sinfo_stream;
        new_item_p->remote_port = ntohs(addr.sin_port);
        new_item_p->remote_addr = addr.sin_addr.s_addr;
        new_item_p->ppid = sinfo.sinfo_ppid;
        new_item_p->assoc_id = sinfo.sinfo_assoc_id;
        new_item_p->length = n;
        new_item_p->buffer = malloc(sizeof(uint8_t) * n);
        memcpy(new_item_p->buffer, buffer, n);
        /*
         * Insert the new packet at the tail of queue.
         */
        TAILQ_INSERT_TAIL(&sctp_data_p->sctp_queue, new_item_p, entry);
        /*
         * Update queue related data
         */
        sctp_data_p->queue_size += n;
        sctp_data_p->queue_length++;
        total_size += n;
      }
    }
  }

  return total_size;
}

int sctp_connect_to_remote_host(
  char *local_ip_addr[],
  int nb_local_addr,
  char *remote_ip_addr,
  constuint16_t port,
  int socket_type,
  sctp_data_t *sctp_data_p)
{
  int sd = -1;
  socklen_t i = 0;
  struct sctp_initmsg init;
  struct sctp_event_subscribe events;
  struct sockaddr *bindx_add_addr;

  DevAssert(sctp_data_p != NULL);
  DevAssert(remote_ip_addr != NULL);
  DevAssert(local_ip_addr != NULL);
  DevCheck((socket_type == SOCK_STREAM), socket_type, 0, 0);
  OAILOG_DEBUG(LOG_SCTP, "Creating socket type %d\n", socket_type);

  /*
   * Create new socket
   */
  if ((sd = socket(AF_INET6, SOCK_STREAM, IPPROTO_SCTP)) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Socket creation failed: %s\n", strerror(errno));
    return -1;
  }

  /*
   * Bind to provided IP adresses
   */
  bindx_add_addr = calloc(nb_local_addr, sizeof(struct sockaddr));

  for (i = 0; i < nb_local_addr; i++) {
    if (
      inet_pton(
        AF_INET,
        local_ip_addr[i],
        &((struct sockaddr_in *) &bindx_add_addr[i])->sin_addr.s_addr) != 1) {
      if (
        inet_pton(
          AF_INET6,
          local_ip_addr[i],
          &((struct sockaddr_in6 *) &bindx_add_addr[i])->sin6_addr.s6_addr) !=
        1) {
        continue;
      } else {
        ((struct sockaddr_in6 *) &bindx_add_addr[i])->sin6_port = 0;
        bindx_add_addr[i].sa_family = AF_INET6;
      }
    } else {
      ((struct sockaddr_in *) &bindx_add_addr[i])->sin_port = 0;
      bindx_add_addr[i].sa_family = AF_INET;
    }
  }

  if (sctp_bindx(sd, bindx_add_addr, nb_local_addr, SCTP_BINDX_ADD_ADDR) < 0) {
    OAILOG_ERROR(LOG_SCTP, "Socket bind failed: %s\n", strerror(errno));
    return -1;
  }

  memset((void *) &init, 0, sizeof(struct sctp_initmsg));
  /*
   * Request a number of in/out streams
   */
  init.sinit_num_ostreams = SCTP_OUT_STREAMS;
  init.sinit_max_instreams = SCTP_IN_STREAMS;
  init.sinit_max_attempts = SCTP_MAX_ATTEMPTS;
  OAILOG_DEBUG(
    LOG_SCTP,
    "Requesting (%d %d) (in out) streams\n",
    init.sinit_num_ostreams,
    init.sinit_max_instreams);

  if (
    setsockopt(
      sd,
      IPPROTO_SCTP,
      SCTP_INITMSG,
      &init,
      (socklen_t) sizeof(struct sctp_initmsg)) < 0) {
    OAILOG_ERROR(
      LOG_SCTP,
      "Setsockopt IPPROTO_SCTP_INITMSG failed: %s\n",
      strerror(errno));
    return -1;
  }

  /*
   * Subscribe to all events
   */
  memset((void *) &events, 1, sizeof(struct sctp_event_subscribe));

  if (
    setsockopt(
      sd,
      IPPROTO_SCTP,
      SCTP_EVENTS,
      &events,
      sizeof(struct sctp_event_subscribe)) < 0) {
    OAILOG_ERROR(
      LOG_SCTP, "Setsockopt IPPROTO_SCTP_EVENTS failed: %s\n", strerror(errno));
    return -1;
  }

  /*
   * SOCK_STREAM socket type requires an explicit connect to the remote host
   * * * * address and port.
   * * * * Only use IPv4 for the first connection attempt
   */
  {
    struct sockaddr_in addr;

    memset(&addr, 0, sizeof(struct sockaddr_in));

    if (inet_pton(AF_INET, remote_ip_addr, &addr.sin_addr.s_addr) != 1) {
      OAILOG_ERROR(
        LOG_SCTP,
        "Failed to convert ip address %s to network type\n",
        remote_ip_addr);
      goto err;
    }

    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    OAILOG_DEBUG(
      LOG_SCTP,
      "[%d] Sending explicit connect to %s:%u\n",
      sd,
      remote_ip_addr,
      port);

    /*
     * Connect to remote host and port
     */
    if (sctp_connectx(sd, (struct sockaddr *) &addr, 1, NULL) < 0) {
      OAILOG_ERROR(
        LOG_SCTP,
        "Connect to %s:%u failed: %s\n",
        remote_ip_addr,
        port,
        strerror(errno));
      goto err;
    }
  }
  /*
   * Get SCTP status
   */
  sctp_get_sockinfo(
    sd,
    &sctp_data_p->instreams,
    &sctp_data_p->outstreams,
    &sctp_data_p->assoc_id);
  sctp_data_p->sd = sd;
  sctp_get_peeraddresses(
    sd, &sctp_data_p->remote_ip_addresses, &sctp_data_p->nb_remote_addresses);
  sctp_get_localaddresses(sd, NULL, NULL);
  TAILQ_INIT(&sctp_data_p->sctp_queue);
  return sd;
err:

  if (sd != 0) {
    close(sd);
  }

  return -1;
}
