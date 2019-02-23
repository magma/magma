/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/** @brief Intertask Interface Signal Dumper
   Allows users to connect their itti_analyzer to this process and dump
   signals exchanged between tasks.
   @author Sebastien Roux <sebastien.roux@eurecom.fr>
*/

#define _GNU_SOURCE // required for pthread_setname_np()
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <error.h>
#include <sched.h>

#include <sys/ioctl.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <sys/types.h>
#include <arpa/inet.h>

#include <sys/eventfd.h>

#include "assertions.h"
#include "liblfds710.h"

#include "itti_types.h"
#include "intertask_interface.h"
#include "intertask_interface_dump.h"
#include "dynamic_memory_check.h"

#if OAI_EMU
#include "vcd_signal_dumper.h"
#endif

static const int itti_dump_debug = 0; // 0x8 | 0x4 | 0x2;

#define ITTI_DUMP_DEBUG(m, x, args...)                                         \
  do {                                                                         \
    if ((m) &itti_dump_debug) fprintf(stdout, "[ITTI_DUMP][D]" x, ##args);     \
  } while (0)
#define ITTI_DUMP_ERROR(x, args...)                                            \
  do {                                                                         \
    fprintf(stdout, "[ITTI_DUMP][E]" x, ##args);                               \
  } while (0)

typedef struct itti_dump_queue_item_s {
  MessageDef *data;
  uint32_t data_size;
  uint32_t message_number;
  uint32_t message_type;
  uint32_t message_size;
} itti_dump_queue_item_t;

typedef struct {
  int sd;
  uint32_t last_message_number;
} itti_client_desc_t;

typedef struct itti_desc_s {
  /*
   * Asynchronous thread that write to file/accept new clients
   */
  pthread_t itti_acceptor_thread;
  pthread_attr_t attr;

  /*
   * List of messages to dump.
   * * * NOTE: we limit the size of this queue to retain only the last exchanged
   * * * messages. The size can be increased by setting up the ITTI_QUEUE_MAX_ELEMENTS
   * * * in mme_default_values.h or by putting a custom in the configuration file.
   */
  struct lfds611_ringbuffer_state *itti_message_queue;

  int nb_connected;

  /*
   * Event fd used to notify new messages (semaphore)
   */
  int event_fd;

  int itti_listen_socket;

  itti_client_desc_t itti_clients[ITTI_DUMP_MAX_CON];
} itti_desc_t;

typedef struct {
  itti_socket_header_t socket_header;

  itti_signal_header_t signal_header;

  /*
   * Message payload is added here, this struct is used as an header
   */
} itti_dump_message_t;

typedef struct {
  itti_socket_header_t socket_header;
} itti_statistic_message_t;

static const itti_message_types_t itti_dump_xml_definition_end =
  ITTI_DUMP_XML_DEFINITION_END;
static const itti_message_types_t itti_dump_message_type_end =
  ITTI_DUMP_MESSAGE_TYPE_END;

static itti_desc_t itti_dump_queue;
static FILE *dump_file = NULL;
static int itti_dump_running = 1;

static volatile uint32_t pending_messages = 0;

/*------------------------------------------------------------------------------*/
static int itti_dump_send_message(int sd, itti_dump_queue_item_t *message)
{
  itti_dump_message_t *new_message;
  ssize_t bytes_sent = 0, total_sent = 0;
  uint8_t *data_ptr;

  /*
   * Allocate memory for message header and payload
   */
  size_t size = sizeof(itti_dump_message_t) + message->data_size +
                sizeof(itti_message_types_t);

  AssertFatal(sd > 0, "Socket descriptor (%d) is invalid!\n", sd);
  AssertFatal(message != NULL, "Message is NULL!\n");
  new_message = malloc(size);
  AssertFatal(new_message != NULL, "New message allocation failed!\n");
  /*
   * Preparing the header
   */
  new_message->socket_header.message_size = size;
  new_message->socket_header.message_type = ITTI_DUMP_MESSAGE_TYPE;
  /*
   * Adds message number in unsigned decimal ASCII format
   */
  snprintf(
    new_message->signal_header.message_number_char,
    sizeof(new_message->signal_header.message_number_char),
    MESSAGE_NUMBER_CHAR_FORMAT,
    message->message_number);
  new_message->signal_header.message_number_char
    [sizeof(new_message->signal_header.message_number_char) - 1] = '\n';
  /*
   * Appends message payload
   */
  memcpy(&new_message[1], message->data, message->data_size);
  memcpy(
    ((void *) &new_message[1]) + message->data_size,
    &itti_dump_message_type_end,
    sizeof(itti_message_types_t));
  data_ptr = (uint8_t *) &new_message[0];

  do {
    bytes_sent = send(sd, &data_ptr[total_sent], size - total_sent, 0);

    if (bytes_sent < 0) {
      ITTI_DUMP_ERROR(
        "[%d] Failed to send %zu bytes to socket (%d:%s)\n",
        sd,
        size,
        errno,
        strerror(errno));
      free_wrapper((void **) &new_message);
      return -1;
    }

    total_sent += bytes_sent;
  } while (total_sent != size);

  free_wrapper((void **) &new_message);
  return total_sent;
}

static int itti_dump_fwrite_message(itti_dump_queue_item_t *message)
{
  itti_dump_message_t new_message_header;

  if ((dump_file != NULL) && (message != NULL)) {
    new_message_header.socket_header.message_size =
      message->message_size + sizeof(itti_dump_message_t) +
      sizeof(itti_message_types_t);
    new_message_header.socket_header.message_type = message->message_type;
    snprintf(
      new_message_header.signal_header.message_number_char,
      sizeof(new_message_header.signal_header.message_number_char),
      MESSAGE_NUMBER_CHAR_FORMAT,
      message->message_number);
    new_message_header.signal_header.message_number_char
      [sizeof(new_message_header.signal_header.message_number_char) - 1] = '\n';
    fwrite(&new_message_header, sizeof(itti_dump_message_t), 1, dump_file);
    fwrite(message->data, message->data_size, 1, dump_file);
    fwrite(
      &itti_dump_message_type_end, sizeof(itti_message_types_t), 1, dump_file);
    fflush(dump_file);
    return (1);
  }

  return (0);
}

static int itti_dump_send_xml_definition(
  const int sd,
  const char *message_definition_xml,
  const uint32_t message_definition_xml_length)
{
  itti_socket_header_t *itti_dump_message;

  /*
   * Allocate memory for message header and payload
   */
  size_t itti_dump_message_size;
  ssize_t bytes_sent = 0, total_sent = 0;
  uint8_t *data_ptr;

  AssertFatal(sd > 0, "Socket descriptor (%d) is invalid!\n", sd);
  AssertFatal(
    message_definition_xml != NULL, "Message definition XML is NULL!\n");
  itti_dump_message_size = sizeof(itti_socket_header_t) +
                           message_definition_xml_length +
                           sizeof(itti_message_types_t);
  itti_dump_message = calloc(1, itti_dump_message_size);
  ITTI_DUMP_DEBUG(
    0x2,
    "[%d] Sending XML definition message of size %zu to observer peer\n",
    sd,
    itti_dump_message_size);
  itti_dump_message->message_size = itti_dump_message_size;
  itti_dump_message->message_type = ITTI_DUMP_XML_DEFINITION;
  /*
   * Copying message definition
   */
  memcpy(
    &itti_dump_message[1],
    message_definition_xml,
    message_definition_xml_length);
  memcpy(
    ((void *) &itti_dump_message[1]) + message_definition_xml_length,
    &itti_dump_xml_definition_end,
    sizeof(itti_message_types_t));
  data_ptr = (uint8_t *) &itti_dump_message[0];

  do {
    bytes_sent =
      send(sd, &data_ptr[total_sent], itti_dump_message_size - total_sent, 0);

    if (bytes_sent < 0) {
      ITTI_DUMP_ERROR(
        "[%d] Failed to send %zu bytes to socket (%d:%s)\n",
        sd,
        itti_dump_message_size,
        errno,
        strerror(errno));
      free_wrapper((void **) &itti_dump_message);
      return -1;
    }

    total_sent += bytes_sent;
  } while (total_sent != itti_dump_message_size);

  free_wrapper((void **) &itti_dump_message);
  return 0;
}

static void itti_dump_user_data_delete_function(
  void *user_data,
  void *user_state)
{
  (void) user_state; // UNUSED

  if (user_data != NULL) {
    itti_dump_queue_item_t *item;
    task_id_t task_id;
    int result;

    item = (itti_dump_queue_item_t *) user_data;

    if (item->data != NULL) {
      task_id = ITTI_MSG_ORIGIN_ID(item->data);
      result = itti_free(task_id, item->data);
      AssertFatal(
        result == EXIT_SUCCESS, "Failed to free memory (%d)!\n", result);
    } else {
      task_id = TASK_UNKNOWN;
    }

    result = itti_free(task_id, item);
    AssertFatal(
      result == EXIT_SUCCESS, "Failed to free memory (%d)!\n", result);
  }
}

static int itti_dump_enqueue_message(
  itti_dump_queue_item_t *new,
  uint32_t message_size,
  uint32_t message_type)
{
  struct lfds611_freelist_element *new_queue_element = NULL;
  int overwrite_flag;

  AssertFatal(new != NULL, "Message to queue is NULL!\n");
#if OAI_EMU
  VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
    VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE, VCD_FUNCTION_IN);
#endif
  new->message_type = message_type;
  new->message_size = message_size;
  ITTI_DUMP_DEBUG(
    0x1, " itti_dump_enqueue_message: lfds611_ringbuffer_get_write_element\n");
  new_queue_element = lfds611_ringbuffer_get_write_element(
    itti_dump_queue.itti_message_queue, &new_queue_element, &overwrite_flag);

  if (overwrite_flag != 0) {
    // no free element available: overwrite a non read one => data loss!
    void *old = NULL;

    lfds611_freelist_get_user_data_from_element(new_queue_element, &old);
    ITTI_DUMP_DEBUG(
      0x4,
      " overwrite_flag set, freeing old data %p %p\n",
      new_queue_element,
      old);
    itti_dump_user_data_delete_function(old, NULL);
  }

  lfds611_freelist_set_user_data_in_element(new_queue_element, new);
  lfds611_ringbuffer_put_write_element(
    itti_dump_queue.itti_message_queue, new_queue_element);

  if (overwrite_flag == 0) {
    {
      ssize_t write_ret;
      eventfd_t sem_counter = 1;

      /*
       * Call to write for an event fd must be of 8 bytes
       */
      write_ret =
        write(itti_dump_queue.event_fd, &sem_counter, sizeof(sem_counter));
      AssertFatal(
        write_ret == sizeof(sem_counter),
        "Write to dump event failed (%d/%d)!\n",
        (int) write_ret,
        (int) sizeof(sem_counter));
    }
    // add one to pending_messages, atomically
    __sync_fetch_and_add(&pending_messages, 1);
  }

  ITTI_DUMP_DEBUG(
    0x2,
    " Added element to queue %p %p, pending %u, type %u\n",
    new_queue_element,
    new,
    pending_messages,
    message_type);
#if OAI_EMU
  VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
    VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE, VCD_FUNCTION_OUT);
#endif
  return 0;
}

static void itti_dump_socket_exit(void)
{
  close(itti_dump_queue.event_fd);
  close(itti_dump_queue.itti_listen_socket);
  /*
   * Leave the thread as we detected end signal
   */
  pthread_exit(NULL);
}

static int itti_dump_flush_ring_buffer(int flush_all)
{
  struct lfds611_freelist_element *element = NULL;
  void *user_data;
  int j;
  int consumer;

  /*
   * Check if there is a least one consumer
   */
  consumer = 0;

  if (dump_file != NULL) {
    consumer = 1;
  } else {
    for (j = 0; j < ITTI_DUMP_MAX_CON; j++) {
      if (itti_dump_queue.itti_clients[j].sd > 0) {
        consumer = 1;
        break;
      }
    }
  }

  if (consumer > 0) {
    do {
      /*
       * Acquire the ring element
       */
      lfds611_ringbuffer_get_read_element(
        itti_dump_queue.itti_message_queue, &element);
      // subtract one from pending_messages, atomically
      __sync_fetch_and_sub(&pending_messages, 1);

      if (element == NULL) {
        if (flush_all != 0) {
          flush_all = 0;
        } else {
          AssertFatal(0, "Dump event with no data!\n");
        }
      } else {
        /*
         * Retrieve user part of the message
         */
        lfds611_freelist_get_user_data_from_element(element, &user_data);
        ITTI_DUMP_DEBUG(
          0x2,
          " removed element from queue %p %p, pending %u\n",
          element,
          user_data,
          pending_messages);

        if (
          ((itti_dump_queue_item_t *) user_data)->message_type ==
          ITTI_DUMP_EXIT_SIGNAL) {
          lfds611_ringbuffer_put_read_element(
            itti_dump_queue.itti_message_queue, element);
          itti_dump_socket_exit();
        }

        /*
         * Write message to file
         */
        itti_dump_fwrite_message((itti_dump_queue_item_t *) user_data);

        /*
         * Send message to remote analyzer
         */
        for (j = 0; j < ITTI_DUMP_MAX_CON; j++) {
          if (itti_dump_queue.itti_clients[j].sd > 0) {
            itti_dump_send_message(
              itti_dump_queue.itti_clients[j].sd,
              (itti_dump_queue_item_t *) user_data);
          }
        }

        itti_dump_user_data_delete_function(user_data, NULL);
        lfds611_freelist_set_user_data_in_element(element, NULL);
        /*
         * We have finished with this element, reinsert it in the ring buffer
         */
        lfds611_ringbuffer_put_read_element(
          itti_dump_queue.itti_message_queue, element);
      }
    } while (flush_all);
  }

  return (consumer);
}

static int itti_dump_handle_new_connection(
  int sd,
  const char *xml_definition,
  uint32_t xml_definition_length)
{
  if (itti_dump_queue.nb_connected < ITTI_DUMP_MAX_CON) {
    uint8_t i;

    for (i = 0; i < ITTI_DUMP_MAX_CON; i++) {
      /*
       * Let's find a place to store the new client
       */
      if (itti_dump_queue.itti_clients[i].sd == -1) {
        break;
      }
    }

    ITTI_DUMP_DEBUG(0x2, " Found place to store new connection: %d\n", i);
    AssertFatal(
      i < ITTI_DUMP_MAX_CON,
      "No more connection available (%d/%d) for socked %d!\n",
      i,
      ITTI_DUMP_MAX_CON,
      sd);
    ITTI_DUMP_DEBUG(0x2, " Socket %d accepted\n", sd);

    /*
     * Send the XML message definition
     */
    if (
      itti_dump_send_xml_definition(sd, xml_definition, xml_definition_length) <
      0) {
      AssertError(0, {}, "Failed to send XML definition!\n");
      close(sd);
      return -1;
    }

    itti_dump_queue.itti_clients[i].sd = sd;
    itti_dump_queue.nb_connected++;
  } else {
    ITTI_DUMP_DEBUG(0x2, " Socket %d rejected\n", sd);
    /*
     * We have reached max number of users connected...
     * * * Reject the connection.
     */
    close(sd);
    return -1;
  }

  return 0;
}

static void *itti_dump_socket(void *arg_p)
{
  uint32_t message_definition_xml_length;
  char *message_definition_xml;
  int rc;
  int itti_listen_socket, max_sd;
  int on = 1;
  fd_set read_set, working_set;
  struct sockaddr_in servaddr; /* socket address structure */
  struct timeval *timeout_p = NULL;

  ITTI_DUMP_DEBUG(0x2, " Creating TCP dump socket on port %u\n", ITTI_PORT);
  message_definition_xml = (char *) arg_p;
  AssertFatal(
    message_definition_xml != NULL, "Message definition XML is NULL!\n");
  message_definition_xml_length = strlen(message_definition_xml) + 1;

  if ((itti_listen_socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)) < 0) {
    ITTI_DUMP_ERROR(" ocket creation failed (%d:%s)\n", errno, strerror(errno));
    pthread_exit(NULL);
  }

  /*
   * Allow socket reuse
   */
  rc = setsockopt(
    itti_listen_socket, SOL_SOCKET, SO_REUSEADDR, (char *) &on, sizeof(on));

  if (rc < 0) {
    ITTI_DUMP_ERROR(
      " setsockopt SO_REUSEADDR failed (%d:%s)\n", errno, strerror(errno));
    close(itti_listen_socket);
    pthread_exit(NULL);
  }

  /*
   * Set socket to be non-blocking.
   * * * NOTE: sockets accepted will inherit this option.
   */
  rc = ioctl(itti_listen_socket, FIONBIO, (char *) &on);

  if (rc < 0) {
    ITTI_DUMP_ERROR(
      " ioctl FIONBIO (non-blocking) failed (%d:%s)\n", errno, strerror(errno));
    close(itti_listen_socket);
    pthread_exit(NULL);
  }

  memset(&servaddr, 0, sizeof(servaddr));
  servaddr.sin_family = AF_INET;
  servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
  servaddr.sin_port = htons(ITTI_PORT);

  if (
    bind(itti_listen_socket, (struct sockaddr *) &servaddr, sizeof(servaddr)) <
    0) {
    ITTI_DUMP_ERROR(" Bind failed (%d:%s)\n", errno, strerror(errno));
    pthread_exit(NULL);
  }

  if (listen(itti_listen_socket, 5) < 0) {
    ITTI_DUMP_ERROR(" Listen failed (%d:%s)\n", errno, strerror(errno));
    pthread_exit(NULL);
  }

  FD_ZERO(&read_set);
  /*
   * Add the listener
   */
  FD_SET(itti_listen_socket, &read_set);
  /*
   * Add the event fd
   */
  FD_SET(itti_dump_queue.event_fd, &read_set);
  /*
   * Max of both sd
   */
  max_sd = itti_listen_socket > itti_dump_queue.event_fd ?
             itti_listen_socket :
             itti_dump_queue.event_fd;

  itti_dump_queue.itti_listen_socket = itti_listen_socket;

  /*
   * Loop waiting for incoming connects or for incoming data
   * * * on any of the connected sockets.
   */
  while (1) {
    int desc_ready;
    int client_socket = -1;
    int i;

    memcpy(&working_set, &read_set, sizeof(read_set));
    timeout_p = NULL;
    /*
     * No timeout: select blocks till a new event has to be handled
     * * * on sd's.
     */
    rc = select(max_sd + 1, &working_set, NULL, NULL, timeout_p);

    if (rc < 0) {
      ITTI_DUMP_ERROR(" select failed (%d:%s)\n", errno, strerror(errno));
      pthread_exit(NULL);
    } else if (rc == 0) {
      /*
       * Timeout
       */
      if (itti_dump_flush_ring_buffer(1) == 0) {
        if (itti_dump_running) {
          ITTI_DUMP_DEBUG(0x4, " No messages consumers, waiting ...\n");
          usleep(100 * 1000);
        } else {
          itti_dump_socket_exit();
        }
      }
    }

    desc_ready = rc;

    for (i = 0; i <= max_sd && desc_ready > 0; i++) {
      if (FD_ISSET(i, &working_set)) {
        desc_ready -= 1;

        if (i == itti_dump_queue.event_fd) {
          /*
           * Notification of new element to dump from other tasks
           */
          eventfd_t sem_counter;
          ssize_t read_ret;

          /*
           * Read will always return 1 for kernel versions > 2.6.30
           */
          read_ret =
            read(itti_dump_queue.event_fd, &sem_counter, sizeof(sem_counter));

          if (read_ret < 0) {
            ITTI_DUMP_ERROR(
              " Failed read for semaphore: %s\n", strerror(errno));
            pthread_exit(NULL);
          }

          AssertFatal(
            read_ret == sizeof(sem_counter),
            "Failed to read from dump event FD (%d/%d)!\n",
            (int) read_ret,
            (int) sizeof(sem_counter));

          if (itti_dump_flush_ring_buffer(0) == 0) {
            if (itti_dump_running) {
              ITTI_DUMP_DEBUG(0x4, " No messages consumers, waiting ...\n");
              usleep(100 * 1000);
              {
                ssize_t write_ret;

                sem_counter = 1;
                /*
                 * Call to write for an event fd must be of 8 bytes
                 */
                write_ret = write(
                  itti_dump_queue.event_fd, &sem_counter, sizeof(sem_counter));
                AssertFatal(
                  write_ret == sizeof(sem_counter),
                  "Failed to write to dump event FD (%d/%d)!\n",
                  (int) write_ret,
                  (int) sem_counter);
              }
            } else {
              itti_dump_socket_exit();
            }
          } else {
            ITTI_DUMP_DEBUG(0x1, " Write element to file\n");
          }
        } else if (i == itti_listen_socket) {
          do {
            client_socket = accept(itti_listen_socket, NULL, NULL);

            if (client_socket < 0) {
              if (errno == EWOULDBLOCK || errno == EAGAIN) {
                /*
                 * No more new connection
                 */
                ITTI_DUMP_DEBUG(0x2, " No more new connection\n");
                continue;
              } else {
                ITTI_DUMP_ERROR(
                  " accept failed (%d:%s)\n", errno, strerror(errno));
                pthread_exit(NULL);
              }
            }

            if (
              itti_dump_handle_new_connection(
                client_socket,
                message_definition_xml,
                message_definition_xml_length) == 0) {
              /*
               * The socket has been accepted.
               * * * We have to update the set to include this new sd.
               */
              FD_SET(client_socket, &read_set);

              if (client_socket > max_sd) max_sd = client_socket;
            }
          } while (client_socket != -1);
        } else {
          /*
           * For now the MME itti dumper should not receive data
           * * * other than connection oriented (CLOSE).
           */
          uint8_t j;

          ITTI_DUMP_DEBUG(0x2, " Socket %d disconnected\n", i);
          /*
           * Close the socket and update info related to this connection
           */
          close(i);

          for (j = 0; j < ITTI_DUMP_MAX_CON; j++) {
            if (itti_dump_queue.itti_clients[j].sd == i) break;
          }

          /*
           * In case we don't find the matching sd in list of known
           * * * connections -> assert.
           */
          AssertFatal(
            j < ITTI_DUMP_MAX_CON,
            "Connection index not found (%d/%d) for socked %d!\n",
            j,
            ITTI_DUMP_MAX_CON,
            i);
          /*
           * Re-initialize the socket to -1 so we can accept new
           * * * incoming connections.
           */
          itti_dump_queue.itti_clients[j].sd = -1;
          itti_dump_queue.itti_clients[j].last_message_number = 0;
          itti_dump_queue.nb_connected--;
          /*
           * Remove the socket from the FD set and update the max sd
           */
          FD_CLR(i, &read_set);

          if (i == max_sd) {
            if (itti_dump_queue.nb_connected == 0) {
              /*
               * No more new connection max_sd = itti_listen_socket
               */
              max_sd = itti_listen_socket;
            } else {
              while (FD_ISSET(max_sd, &read_set) == 0) {
                max_sd -= 1;
              }
            }
          }
        }
      }
    }
  }

  return NULL;
}

/*------------------------------------------------------------------------------*/
int itti_dump_queue_message(
  task_id_t sender_task,
  message_number_t message_number,
  MessageDef *message_p,
  const char *message_name,
  const uint32_t message_size)
{
  if (itti_dump_running) {
    itti_dump_queue_item_t *new;

    AssertFatal(message_name != NULL, "Message name is NULL!\n");
    AssertFatal(message_p != NULL, "Message is NULL!\n");
#if OAI_EMU
    VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
      VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE_malloc,
      VCD_FUNCTION_IN);
#endif
    new = itti_malloc(sender_task, TASK_MAX, sizeof(itti_dump_queue_item_t));
#if OAI_EMU
    VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
      VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE_malloc,
      VCD_FUNCTION_OUT);
#endif
#if OAI_EMU
    VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
      VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE_malloc,
      VCD_FUNCTION_IN);
#endif
    new->data = itti_malloc(sender_task, TASK_MAX, message_size);
#if OAI_EMU
    VCD_SIGNAL_DUMPER_DUMP_FUNCTION_BY_NAME(
      VCD_SIGNAL_DUMPER_FUNCTIONS_ITTI_DUMP_ENQUEUE_MESSAGE_malloc,
      VCD_FUNCTION_OUT);
#endif
    memcpy(new->data, message_p, message_size);
    new->data_size = message_size;
    new->message_number = message_number;
    itti_dump_enqueue_message(new, message_size, ITTI_DUMP_MESSAGE_TYPE);
  }

  return 0;
}

/* This function should be called by each thread that will use the ring buffer */
void itti_dump_thread_use_ring_buffer(void)
{
  lfds611_ringbuffer_use(itti_dump_queue.itti_message_queue);
}

int itti_dump_init(
  const char *const messages_definition_xml,
  const char *const dump_file_name)
{
  int i, ret;
  struct sched_param scheduler_param;

  scheduler_param.sched_priority = sched_get_priority_min(SCHED_FIFO) + 1;

  if (dump_file_name != NULL) {
    dump_file = fopen(dump_file_name, "wb");

    if (dump_file == NULL) {
      ITTI_DUMP_ERROR(
        " can not open dump file \"%s\" (%d:%s)\n",
        dump_file_name,
        errno,
        strerror(errno));
    } else {
      /*
       * Output the XML to file
       */
      uint32_t message_size = strlen(messages_definition_xml) + 1;
      itti_socket_header_t header;

      header.message_size = sizeof(itti_socket_header_t) + message_size +
                            sizeof(itti_message_types_t);
      header.message_type = ITTI_DUMP_XML_DEFINITION;
      fwrite(&header, sizeof(itti_socket_header_t), 1, dump_file);
      fwrite(messages_definition_xml, message_size, 1, dump_file);
      fwrite(
        &itti_dump_xml_definition_end,
        sizeof(itti_message_types_t),
        1,
        dump_file);
      fflush(dump_file);
    }
  }

  memset(&itti_dump_queue, 0, sizeof(itti_desc_t));
  ITTI_DUMP_DEBUG(
    0x2,
    " Creating new ring buffer for itti dump of %u elements\n",
    ITTI_QUEUE_MAX_ELEMENTS);

  if (
    lfds611_ringbuffer_new(
      &itti_dump_queue.itti_message_queue,
      ITTI_QUEUE_MAX_ELEMENTS,
      NULL,
      NULL) != 1) {
    /*
     * Always assert on this condition
     */
    AssertFatal(0, " Failed to create ring buffer!\n");
  }

  itti_dump_queue.event_fd = eventfd(0, EFD_SEMAPHORE);

  if (itti_dump_queue.event_fd == -1) {
    /*
     * Always assert on this condition
     */
    AssertFatal(0, "eventfd failed: %s!\n", strerror(errno));
  }
  itti_dump_queue.nb_connected = 0;

  for (i = 0; i < ITTI_DUMP_MAX_CON; i++) {
    itti_dump_queue.itti_clients[i].sd = -1;
    itti_dump_queue.itti_clients[i].last_message_number = 0;
  }

  /*
   * initialized with default attributes
   */
  ret = pthread_attr_init(&itti_dump_queue.attr);

  if (ret < 0) {
    AssertFatal(
      0, "pthread_attr_init failed (%d:%s)!\n", errno, strerror(errno));
  }

  ret = pthread_attr_setschedpolicy(&itti_dump_queue.attr, SCHED_FIFO);

  if (ret < 0) {
    AssertFatal(
      0,
      "pthread_attr_setschedpolicy (SCHED_IDLE) failed (%d:%s)!\n",
      errno,
      strerror(errno));
  }

  ret = pthread_attr_setschedparam(&itti_dump_queue.attr, &scheduler_param);

  if (ret < 0) {
    AssertFatal(
      0,
      "pthread_attr_setschedparam failed (%d:%s)!\n",
      errno,
      strerror(errno));
  }

  ret = pthread_create(
    &itti_dump_queue.itti_acceptor_thread,
    &itti_dump_queue.attr,
    &itti_dump_socket,
    (void *) messages_definition_xml);

  if (ret < 0) {
    AssertFatal(0, "pthread_create failed (%d:%s)!\n", errno, strerror(errno));
  }

  pthread_setname_np(itti_dump_queue.itti_acceptor_thread, "ITTI acceptor");
  return 0;
}

void itti_dump_exit(void)
{
  void *arg;
  itti_dump_queue_item_t *new;

  new = itti_malloc(TASK_UNKNOWN, TASK_UNKNOWN, sizeof(itti_dump_queue_item_t));
  memset(new, 0, sizeof(itti_dump_queue_item_t));
  /*
   * Set a flag to stop recording message
   */
  itti_dump_running = 0;
  /*
   * Send the exit signal to other thread
   */
  itti_dump_enqueue_message(new, 0, ITTI_DUMP_EXIT_SIGNAL);
  ITTI_DUMP_DEBUG(0x2, " waiting for dumper thread to finish\n");
  /*
   * wait for the thread to terminate
   */
  pthread_join(itti_dump_queue.itti_acceptor_thread, &arg);
  ITTI_DUMP_DEBUG(0x2, " dumper thread correctly exited\n");

  if (dump_file != NULL) {
    /*
     * Synchronise file and then close it
     */
    fclose(dump_file);
    dump_file = NULL;
  }

  if (itti_dump_queue.itti_message_queue) {
    lfds611_ringbuffer_delete(
      itti_dump_queue.itti_message_queue,
      itti_dump_user_data_delete_function,
      NULL);
  }
}
