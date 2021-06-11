/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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

/*! \file udp_primitives_server.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <arpa/inet.h>
#include <errno.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>

#include "bstrlib.h"

#include "assertions.h"
#include "conversions.h"
#include "dynamic_memory_check.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "log.h"
#include "queue.h"
#include "udp_primitives_server.h"

task_zmq_ctx_t udp_task_zmq_ctx;

struct udp_socket_desc_s {
  uint8_t buffer[4096];
  int sd; /* Socket descriptor to use */

  pthread_t listener_thread; /* Thread affected to recv */

  struct sockaddr local_addr; /* Local ipv4 or ipv6 address to use */
  uint16_t local_port;        /* Local port to use */

  task_id_t task_id; /* Task who has requested the new endpoint */
  STAILQ_ENTRY(udp_socket_desc_s) entries;
};

static STAILQ_HEAD(udp_socket_list_s, udp_socket_desc_s) udp_socket_list;
static pthread_mutex_t udp_socket_list_mutex = PTHREAD_MUTEX_INITIALIZER;

static void udp_server_receive_and_process(
    struct udp_socket_desc_s* udp_sock_pP);

/* @brief Retrieve the descriptor associated with the task_id
 */
static struct udp_socket_desc_s* udp_server_get_socket_desc(
    task_id_t task_id, uint16_t local_port, uint16_t peer_port, int sa_family) {
  struct udp_socket_desc_s* udp_sock_p = NULL;

  OAILOG_DEBUG(LOG_UDP, "Looking for task %d\n", task_id);
  STAILQ_FOREACH(udp_sock_p, &udp_socket_list, entries) {
    if (udp_sock_p->task_id == task_id &&
        udp_sock_p->local_addr.sa_family == sa_family) {
      if (local_port) {
        if (udp_sock_p->local_port == local_port) {
          OAILOG_DEBUG(LOG_UDP, "Found matching local port %d. \n", local_port);
          break;
        } else {
          continue;
        }
      }
      /** Reach here only if no local port. */
      if (udp_sock_p->local_port != peer_port) {
        /** Must be high port. */
        OAILOG_DEBUG(
            LOG_UDP, "Found matching port with high port %d. \n",
            udp_sock_p->local_port);
        break;
      } else {
        /** If no local port is given, the destination will be 2123. That cannot
         * be the local port. */
        continue;
      }
    }
  }
  return udp_sock_p;
}

static struct udp_socket_desc_s* udp_server_get_socket_desc_by_sd(int sdP) {
  struct udp_socket_desc_s* udp_sock_p = NULL;

  OAILOG_DEBUG(LOG_UDP, "Looking for sd %d\n", sdP);
  STAILQ_FOREACH(udp_sock_p, &udp_socket_list, entries) {
    if (udp_sock_p->sd == sdP) {
      break;
    }
  }
  return udp_sock_p;
}

static void udp_server_receive_and_process(
    struct udp_socket_desc_s* udp_sock_pP) {
  OAILOG_DEBUG(
      LOG_UDP, "Inserting new descriptor for task %d, sd %d\n",
      udp_sock_pP->task_id, udp_sock_pP->sd);
  {
    int bytes_received = 0;
    socklen_t from_len;
    struct sockaddr_in addr;
    struct sockaddr_in6 addr6;

    bool ipv6 = udp_sock_pP->local_addr.sa_family == AF_INET6;

    if (ipv6) {
      from_len = (socklen_t) sizeof(struct sockaddr_in6);
      memset(&addr6, 0, sizeof(struct sockaddr_in6));
    } else {
      from_len = (socklen_t) sizeof(struct sockaddr_in);
      memset(&addr, 0, sizeof(struct sockaddr_in));
    }

    if ((bytes_received = recvfrom(
             udp_sock_pP->sd, udp_sock_pP->buffer, sizeof(udp_sock_pP->buffer),
             0, ipv6 ? (struct sockaddr*) &addr6 : (struct sockaddr*) &addr,
             &from_len)) <= 0) {
      OAILOG_ERROR(LOG_UDP, "Recvfrom failed %s\n", strerror(errno));
      // break;
    } else {
      MessageDef* message_p = NULL;
      udp_data_ind_t* udp_data_ind_p;
      AssertFatal(
          sizeof(udp_sock_pP->buffer) >= bytes_received, "UDP BUFFER OVERFLOW");
      message_p = itti_alloc_new_message(TASK_UDP, UDP_DATA_IND);
      DevAssert(message_p != NULL);
      udp_data_ind_p = &message_p->ittiMsg.udp_data_ind;
      memcpy(udp_data_ind_p->msgBuf, udp_sock_pP->buffer, bytes_received);

      udp_data_ind_p->buffer_length = bytes_received;
      udp_data_ind_p->local_port    = udp_sock_pP->local_port;
      udp_data_ind_p->peer_port =
          ipv6 ? htons(addr6.sin6_port) : htons(addr.sin_port);
      memcpy(
          (void*) &udp_data_ind_p->sock_addr,
          (ipv6) ? (void*) &addr6 : (void*) &addr,
          (ipv6) ? sizeof(struct sockaddr_in6) : sizeof(struct sockaddr_in));

      OAILOG_DEBUG(
          LOG_UDP, "Msg of length %d received from %s:%u\n", bytes_received,
          (!ipv6) ? inet_ntoa(addr.sin_addr) : "TODO_IPV6",
          ntohs(addr.sin_port));

      if (send_msg_to_task(&udp_task_zmq_ctx, udp_sock_pP->task_id, message_p) <
          0) {
        OAILOG_DEBUG(
            LOG_UDP, "Failed to send message %d to task %d\n", UDP_DATA_IND,
            udp_sock_pP->task_id);
        // break;
      }
    }
  }
}

static int udp_socket_handler(zloop_t* loop, zmq_pollitem_t* item, void* arg) {
  struct udp_socket_desc_s* udp_sock_p = NULL;

  pthread_mutex_lock(&udp_socket_list_mutex);
  udp_sock_p = udp_server_get_socket_desc_by_sd(item->fd);

  if (udp_sock_p != NULL) {
    udp_server_receive_and_process(udp_sock_p);
  } else {
    OAILOG_ERROR(
        LOG_UDP, "Failed to retrieve the udp socket descriptor %d", item->fd);
  }

  pthread_mutex_unlock(&udp_socket_list_mutex);

  return 0;
}

//------------------------------------------------------------------------------
static int udp_server_create_socket_v4(
    uint16_t port, struct in_addr* address, task_id_t task_id) {
  struct sockaddr_in addr;
  int sd;
  struct udp_socket_desc_s* socket_desc_p = NULL;

  // OAILOG_INFO();
  /*
   * Create UDP socket
   */
  if ((sd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
    /*
     * Socket creation has failed...
     */
    OAILOG_ERROR(
        LOG_UDP, "IPv4 socket creation failed (%s)\n", strerror(errno));
    return sd;
  }

  memset(&addr, 0, sizeof(struct sockaddr_in));
  addr.sin_family = AF_INET;
  addr.sin_port   = htons(port);
  addr.sin_addr   = *address;

  char ipv4[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, (void*) &addr.sin_addr, ipv4, INET_ADDRSTRLEN);
  OAILOG_DEBUG(
      LOG_UDP,
      "Creating new listen socket on IPv4 address %s and port %" PRIu16 "\n",
      ipv4, port);

  if (bind(sd, (struct sockaddr*) &addr, sizeof(struct sockaddr_in)) < 0) {
    /*
     * Bind failed
     */
    OAILOG_ERROR(
        LOG_UDP,
        "Socket bind failed (%s) for IPv4 address %s and port %" PRIu16 "\n",
        strerror(errno), ipv4, port);
    close(sd);
    return -1;
  }
  struct sockaddr_in addr_check;
  socklen_t len = sizeof(addr_check);
  if (getsockname(sd, (struct sockaddr*) &addr_check, &len) != -1)
    OAILOG_DEBUG(
        LOG_UDP, "Listened on port %" PRIu16 " for IPv4. \n",
        ntohs(addr_check.sin_port));

  /*
   * Add the socket to list of fd monitored by ITTI
   */
  /*
   * Mark the socket as non-blocking
   */
  if (fcntl(sd, F_SETFL, O_NONBLOCK) < 0) {
    OAILOG_ERROR(
        LOG_UDP, "fcntl F_SETFL O_NONBLOCK failed for IPv4: %s\n",
        strerror(errno));
    close(sd);
    return -1;
  }

  socket_desc_p = calloc(1, sizeof(struct udp_socket_desc_s));
  DevAssert(socket_desc_p != NULL);
  socket_desc_p->sd                                            = sd;
  ((struct sockaddr_in*) &socket_desc_p->local_addr)->sin_addr = *address;
  socket_desc_p->local_addr.sa_family                          = AF_INET;
  socket_desc_p->local_port = ntohs(addr_check.sin_port);
  socket_desc_p->task_id    = task_id;
  OAILOG_DEBUG(
      LOG_UDP, "(IPv4) Inserting new descriptor for task %d, sd %d\n",
      socket_desc_p->task_id, socket_desc_p->sd);
  pthread_mutex_lock(&udp_socket_list_mutex);
  STAILQ_INSERT_TAIL(&udp_socket_list, socket_desc_p, entries);
  pthread_mutex_unlock(&udp_socket_list_mutex);

  zmq_pollitem_t item = {0, sd, ZMQ_POLLIN, 0};
  zloop_poller(udp_task_zmq_ctx.event_loop, &item, udp_socket_handler, NULL);

  return sd;
}

//------------------------------------------------------------------------------
static int udp_server_create_socket_v6(
    uint16_t port, struct in6_addr* address, task_id_t task_id) {
  struct sockaddr_in6 addr;
  int sd;
  struct udp_socket_desc_s* socket_desc_p = NULL;

  /*
   * Create UDP socket
   */
  if ((sd = socket(AF_INET6, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
    /*
     * Socket creation has failed...
     */
    OAILOG_ERROR(
        LOG_UDP, "IPv6 socket creation failed (%s)\n", strerror(errno));
    return sd;
  }

  memset(&addr, 0, sizeof(struct sockaddr_in6));
  addr.sin6_family = AF_INET6;
  addr.sin6_port   = htons(port);
  addr.sin6_addr   = *address;

  char ipv6[INET6_ADDRSTRLEN];
  inet_ntop(AF_INET6, (void*) &addr, ipv6, INET6_ADDRSTRLEN);
  OAILOG_DEBUG(
      LOG_UDP,
      "Creating new listen socket on IPv6 address %s and port %" PRIu16 "\n",
      ipv6, port);

  if (bind(sd, (struct sockaddr*) &addr, sizeof(struct sockaddr_in6)) < 0) {
    /*
     * Bind failed
     */
    //    OAILOG_ERROR (LOG_UDP, "Socket bind failed (%s) for address %s and
    //    port %" PRIu16 "\n", strerror (errno), ipv4, port);
    close(sd);
    return -1;
  }
  struct sockaddr_in addr_check;
  socklen_t len = sizeof(addr_check);
  if (getsockname(sd, (struct sockaddr*) &addr_check, &len) != -1)
    OAILOG_DEBUG(
        LOG_UDP, "Listened on port %" PRIu16 " for IPv6. \n",
        ntohs(addr_check.sin_port));

  /*
   * Add the socket to list of fd monitored by ITTI
   */
  /*
   * Mark the socket as non-blocking
   */
  if (fcntl(sd, F_SETFL, O_NONBLOCK) < 0) {
    OAILOG_ERROR(
        LOG_UDP, "fcntl F_SETFL O_NONBLOCK failed: %s\n", strerror(errno));
    close(sd);
    return -1;
  }

  socket_desc_p = calloc(1, sizeof(struct udp_socket_desc_s));
  DevAssert(socket_desc_p != NULL);
  socket_desc_p->sd                                              = sd;
  ((struct sockaddr_in6*) &socket_desc_p->local_addr)->sin6_addr = *address;
  socket_desc_p->local_addr.sa_family                            = AF_INET6;

  //  ((struct sockaddr_in6*)&socket_desc_p->local_addr)->sin6_family = AF_INET;
  socket_desc_p->local_port = ntohs(addr_check.sin_port);
  socket_desc_p->task_id    = task_id;
  OAILOG_DEBUG(
      LOG_UDP, "(IPv6) Inserting new descriptor for task %d, sd %d\n",
      socket_desc_p->task_id, socket_desc_p->sd);
  pthread_mutex_lock(&udp_socket_list_mutex);
  STAILQ_INSERT_TAIL(&udp_socket_list, socket_desc_p, entries);
  pthread_mutex_unlock(&udp_socket_list_mutex);

  zmq_pollitem_t item = {0, sd, ZMQ_POLLIN, 0};
  zloop_poller(udp_task_zmq_ctx.event_loop, &item, udp_socket_handler, NULL);

  return sd;
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST: {
      OAI_FPRINTF_INFO("TASK_UDP received MESSAGE_TEST\n");
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      udp_exit();
    } break;

    case UDP_INIT: {
      udp_init_t* udp_init_p = &received_message_p->ittiMsg.udp_init;
      if (udp_init_p->in_addr)
        udp_server_create_socket_v4(
            udp_init_p->port, udp_init_p->in_addr,
            ITTI_MSG_ORIGIN_ID(received_message_p));
      if (udp_init_p->in6_addr)
        udp_server_create_socket_v6(
            udp_init_p->port, udp_init_p->in6_addr,
            ITTI_MSG_ORIGIN_ID(received_message_p));
    } break;

    case UDP_DATA_REQ: {
      int udp_sd = -1;
      ssize_t bytes_written;
      struct udp_socket_desc_s* udp_sock_p = NULL;
      udp_data_req_t* udp_data_req_p;
      struct sockaddr_in peer_addr;
      struct sockaddr_in6 peer_addr6;

      udp_data_req_p = &received_message_p->ittiMsg.udp_data_req;
      // UDP_DEBUG("-- UDP_DATA_REQ
      // -----------------------------------------------------\n%s :\n",
      //        __FUNCTION__);
      // udp_print_hex_octets(&udp_data_req_p->buffer[udp_data_req_p->buffer_offset],
      //        udp_data_req_p->buffer_length);
      if (udp_data_req_p->peer_address->sa_family == AF_INET) {
        memset(&peer_addr, 0, sizeof(struct sockaddr_in));
        peer_addr.sin_family = AF_INET;
        peer_addr.sin_port   = htons(udp_data_req_p->peer_port);
        peer_addr.sin_addr =
            ((struct sockaddr_in*) udp_data_req_p->peer_address)->sin_addr;
        pthread_mutex_lock(&udp_socket_list_mutex);
        udp_sock_p = udp_server_get_socket_desc(
            ITTI_MSG_ORIGIN_ID(received_message_p), udp_data_req_p->local_port,
            udp_data_req_p->peer_port, AF_INET);

        if (udp_sock_p == NULL) {
          OAILOG_ERROR(
              LOG_UDP,
              "Failed to retrieve the udp socket descriptor for IPv4 "
              "associated with task %d\n",
              ITTI_MSG_ORIGIN_ID(received_message_p));
          pthread_mutex_unlock(&udp_socket_list_mutex);
          // no free udp_data_req_p->buffer, statically allocated
          break;
        }

        udp_sd = udp_sock_p->sd;
        pthread_mutex_unlock(&udp_socket_list_mutex);
        OAILOG_DEBUG(
            LOG_UDP,
            "[%d] Sending message of size %u to " IN_ADDR_FMT " and port %u\n",
            udp_sd, udp_data_req_p->buffer_length,
            PRI_IN_ADDR(
                ((struct sockaddr_in*) udp_data_req_p->peer_address)->sin_addr),
            udp_data_req_p->peer_port);
        bytes_written = sendto(
            udp_sd, &udp_data_req_p->buffer[udp_data_req_p->buffer_offset],
            udp_data_req_p->buffer_length, 0, (struct sockaddr*) &peer_addr,
            sizeof(struct sockaddr_in));
        // no free udp_data_req_p->buffer, statically allocated

        if (bytes_written != udp_data_req_p->buffer_length) {
          OAILOG_ERROR(
              LOG_UDP,
              "There was an error while writing to socket "
              "(%d:%s)\n",
              errno, strerror(errno));
        }
      } else if (udp_data_req_p->peer_address->sa_family == AF_INET6) {
        memset(&peer_addr6, 0, sizeof(struct sockaddr_in6));
        peer_addr6.sin6_family = AF_INET6;
        peer_addr6.sin6_port   = htons(udp_data_req_p->peer_port);
        peer_addr6.sin6_addr =
            ((struct sockaddr_in6*) udp_data_req_p->peer_address)->sin6_addr;
        pthread_mutex_lock(&udp_socket_list_mutex);
        udp_sock_p = udp_server_get_socket_desc(
            ITTI_MSG_ORIGIN_ID(received_message_p), udp_data_req_p->local_port,
            udp_data_req_p->peer_port, AF_INET6);

        if (udp_sock_p == NULL) {
          OAILOG_ERROR(
              LOG_UDP,
              "Failed to retrieve the udp socket descriptor for IPv6 "
              "associated with task %d\n",
              ITTI_MSG_ORIGIN_ID(received_message_p));
          pthread_mutex_unlock(&udp_socket_list_mutex);
          // no free udp_data_req_p->buffer, statically allocated
          break;
        }

        udp_sd = udp_sock_p->sd;
        pthread_mutex_unlock(&udp_socket_list_mutex);
        bytes_written = sendto(
            udp_sd, &udp_data_req_p->buffer[udp_data_req_p->buffer_offset],
            udp_data_req_p->buffer_length, 0, (struct sockaddr*) &peer_addr6,
            sizeof(struct sockaddr_in6));
        // no free udp_data_req_p->buffer, statically allocated

        if (bytes_written != udp_data_req_p->buffer_length) {
          OAILOG_ERROR(
              LOG_UDP,
              "There was an error while writing to socket "
              "(%d:%s)\n",
              errno, strerror(errno));
        }
      }

      else {
        OAILOG_DEBUG(LOG_UDP, "Unknown address type");
      }

    } break;

    default: {
      OAILOG_DEBUG(
          LOG_UDP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* udp_thread(void* args) {
  itti_mark_task_ready(TASK_UDP);
  init_task_context(
      TASK_UDP, (task_id_t[]){TASK_MME_APP, TASK_S11}, 2, handle_message,
      &udp_task_zmq_ctx);

  zloop_start(udp_task_zmq_ctx.event_loop);
  udp_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int udp_init(void) {
  OAILOG_DEBUG(LOG_UDP, "Initializing UDP task interface\n");
  STAILQ_INIT(&udp_socket_list);

  if (itti_create_task(TASK_UDP, &udp_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_UDP, "udp pthread_create (%s)\n", strerror(errno));
    return -1;
  }

  OAILOG_DEBUG(LOG_UDP, "Initializing UDP task interface: DONE\n");
  return 0;
}

//------------------------------------------------------------------------------
void udp_exit(void) {
  struct udp_socket_desc_s* socket_desc_p = NULL;
  while ((socket_desc_p = STAILQ_FIRST(&udp_socket_list))) {
    close(socket_desc_p->sd);
    pthread_mutex_destroy(&udp_socket_list_mutex);
    STAILQ_REMOVE_HEAD(&udp_socket_list, entries);
    free_wrapper((void**) &socket_desc_p);
  }

  destroy_task_context(&udp_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_UDP terminated\n");
  pthread_exit(NULL);
}
