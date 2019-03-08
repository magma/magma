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

/*! \file sctp_common.c
    \brief MME SCTP related common procedures
    \author Sebastien ROUX, Lionel GAUTHIER
    \date 2013
    \version 1.0
    @ingroup _sctp
*/

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <netinet/sctp.h>
#include <stdint.h>

#include "sctp_common.h"
#include "log.h"

/* Pre-bind socket options configuration.
   See http://linux.die.net/man/7/sctp for more informations on these options.
*/
//------------------------------------------------------------------------------
int sctp_set_init_opt(
  const int sd,
  const sctp_stream_id_t instreams,
  const sctp_stream_id_t outstreams,
  const uint16_t max_attempts,
  const uint16_t init_timeout)
{
  int on = 1;
  struct sctp_initmsg init = {0};

  memset((void *) &init, 0, sizeof(struct sctp_initmsg));
  /*
   * Request a number of streams
   */
  init.sinit_num_ostreams = outstreams;
  init.sinit_max_instreams = instreams;
  init.sinit_max_attempts = max_attempts;
  init.sinit_max_init_timeo = init_timeout;

  if (
    setsockopt(
      sd, IPPROTO_SCTP, SCTP_INITMSG, &init, sizeof(struct sctp_initmsg)) < 0) {
    OAILOG_ERROR(LOG_SCTP, "setsockopt: %d:%s\n", errno, strerror(errno));
    close(sd);
    return -1;
  }

  /*
   * Allow socket reuse
   */
  if (setsockopt(sd, SOL_SOCKET, SO_REUSEADDR, &on, sizeof(on)) < 0) {
    OAILOG_ERROR(
      LOG_SCTP,
      "setsockopt SO_REUSEADDR failed (%d:%s)\n",
      errno,
      strerror(errno));
    close(sd);
    return -1;
  }

  return 0;
}

//------------------------------------------------------------------------------
int sctp_get_sockinfo(
  int sock,
  sctp_stream_id_t *instream,
  sctp_stream_id_t *outstream,
  sctp_assoc_id_t *assoc_id)
{
  socklen_t i = 0;
  struct sctp_status status = {0};

  if (socket <= 0) {
    return -1;
  }

  memset(&status, 0, sizeof(struct sctp_status));
  i = sizeof(struct sctp_status);

  if (getsockopt(sock, IPPROTO_SCTP, SCTP_STATUS, &status, &i) < 0) {
    OAILOG_ERROR(
      LOG_SCTP, "Getsockopt SCTP_STATUS failed: %s\n", strerror(errno));
    return -1;
  }

  OAILOG_DEBUG(LOG_SCTP, "----------------------\n");
  OAILOG_DEBUG(LOG_SCTP, "SCTP Status:\n");
  OAILOG_DEBUG(LOG_SCTP, "assoc id .....: %u\n", status.sstat_assoc_id);
  OAILOG_DEBUG(LOG_SCTP, "state ........: %d\n", status.sstat_state);
  OAILOG_DEBUG(LOG_SCTP, "instrms ......: %u\n", status.sstat_instrms);
  OAILOG_DEBUG(LOG_SCTP, "outstrms .....: %u\n", status.sstat_outstrms);
  OAILOG_DEBUG(
    LOG_SCTP, "fragmentation : %u\n", status.sstat_fragmentation_point);
  OAILOG_DEBUG(LOG_SCTP, "pending data .: %u\n", status.sstat_penddata);
  OAILOG_DEBUG(LOG_SCTP, "unack data ...: %u\n", status.sstat_unackdata);
  OAILOG_DEBUG(LOG_SCTP, "rwnd .........: %u\n", status.sstat_rwnd);
  OAILOG_DEBUG(LOG_SCTP, "peer info     :\n");
  OAILOG_DEBUG(
    LOG_SCTP, "    state ....: %u\n", status.sstat_primary.spinfo_state);
  OAILOG_DEBUG(
    LOG_SCTP, "    cwnd .....: %u\n", status.sstat_primary.spinfo_cwnd);
  OAILOG_DEBUG(
    LOG_SCTP, "    srtt .....: %u\n", status.sstat_primary.spinfo_srtt);
  OAILOG_DEBUG(
    LOG_SCTP, "    rto ......: %u\n", status.sstat_primary.spinfo_rto);
  OAILOG_DEBUG(
    LOG_SCTP, "    mtu ......: %u\n", status.sstat_primary.spinfo_mtu);
  OAILOG_DEBUG(LOG_SCTP, "----------------------\n");

  if (instream != NULL) {
    *instream = status.sstat_instrms;
  }

  if (outstream != NULL) {
    *outstream = status.sstat_outstrms;
  }

  if (assoc_id != NULL) {
    *assoc_id = status.sstat_assoc_id;
  }

  return 0;
}

//------------------------------------------------------------------------------
int sctp_get_peeraddresses(
  int sock,
  struct sockaddr **remote_addr,
  int *nb_remote_addresses)
{
  int nb = 0, j = 0;
  struct sockaddr *temp_addr_p = NULL;

  if ((nb = sctp_getpaddrs(sock, -1, &temp_addr_p)) <= 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to retrieve peer addresses\n");
    return -1;
  }

  OAILOG_DEBUG(LOG_SCTP, "----------------------\n");
  OAILOG_DEBUG(LOG_SCTP, "Peer addresses:\n");

  for (j = 0; j < nb; j++) {
    if (temp_addr_p[j].sa_family == AF_INET) {
      char address[INET_ADDRSTRLEN] = {0};
      struct sockaddr_in *addr = NULL;

      addr = (struct sockaddr_in *) &temp_addr_p[j];

      if (
        inet_ntop(AF_INET, &addr->sin_addr, address, sizeof(address)) != NULL) {
        OAILOG_DEBUG(LOG_SCTP, "    - [%s]\n", address);
      }
    } else {
      struct sockaddr_in6 *addr = NULL;
      char address[INET6_ADDRSTRLEN] = {0};

      addr = (struct sockaddr_in6 *) &temp_addr_p[j];

      if (
        inet_ntop(
          AF_INET6, &addr->sin6_addr.s6_addr, address, sizeof(address)) !=
        NULL) {
        OAILOG_DEBUG(LOG_SCTP, "    - [%s]\n", address);
      }
    }
  }

  OAILOG_DEBUG(LOG_SCTP, "----------------------\n");

  if (remote_addr != NULL && nb_remote_addresses != NULL) {
    *nb_remote_addresses = nb;
    *remote_addr = temp_addr_p;
  } else {
    /*
     * We can destroy buffer
     */
    sctp_freepaddrs((struct sockaddr *) temp_addr_p);
  }

  return 0;
}

//------------------------------------------------------------------------------
int sctp_get_localaddresses(
  int sock,
  struct sockaddr **local_addr,
  int *nb_local_addresses)
{
  int nb = 0, j = 0;
  struct sockaddr *temp_addr_p = NULL;

  if ((nb = sctp_getladdrs(sock, -1, &temp_addr_p)) <= 0) {
    OAILOG_ERROR(LOG_SCTP, "Failed to retrieve local addresses\n");
    return -1;
  }

  if (temp_addr_p) {
    OAILOG_DEBUG(LOG_SCTP, "----------------------\n");
    OAILOG_DEBUG(LOG_SCTP, "Local addresses:\n");

    for (j = 0; j < nb; j++) {
      if (temp_addr_p[j].sa_family == AF_INET) {
        char address[INET_ADDRSTRLEN] = {0};
        struct sockaddr_in *addr = NULL;

        addr = (struct sockaddr_in *) &temp_addr_p[j];

        if (
          inet_ntop(AF_INET, &addr->sin_addr, address, sizeof(address)) !=
          NULL) {
          OAILOG_DEBUG(LOG_SCTP, "    - [%s]\n", address);
        }
      } else if (temp_addr_p[j].sa_family == AF_INET6) {
        struct sockaddr_in6 *addr = NULL;
        char address[INET6_ADDRSTRLEN] = {0};

        addr = (struct sockaddr_in6 *) &temp_addr_p[j];

        if (
          inet_ntop(
            AF_INET6, &addr->sin6_addr.s6_addr, address, sizeof(address)) !=
          NULL) {
          OAILOG_DEBUG(LOG_SCTP, "    - [%s]\n", address);
        }
      } else {
        OAILOG_DEBUG(
          LOG_SCTP,
          "    - unhandled address family %d\n",
          temp_addr_p[j].sa_family);
      }
    }
    OAILOG_DEBUG(LOG_SCTP, "----------------------\n");

    if (local_addr != NULL && nb_local_addresses != NULL) {
      *nb_local_addresses = nb;
      *local_addr = temp_addr_p;
    } else {
      /*
       * We can destroy buffer
       */
      sctp_freeladdrs((struct sockaddr *) temp_addr_p);
    }
  }

  return 0;
}
