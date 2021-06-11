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

/*****************************************************************************

  Source    network_api.c

  Version   0.1

  Date    2012/03/01

  Product   NAS stack

  Subsystem Application Programming Interface

  Author    Frederic Maurel

  Description Implements the API used by the NAS layer to send/receive
    message to/from the network sublayer

*****************************************************************************/

#include "network_api.h"
#include "common_types.h"
#include "log.h"
#include "socket.h"

#include "nas/as_message.h"

#include <string.h>  // strerror, memset
#include <netdb.h>   // gai_strerror
#include <errno.h>   // errno
#include <unistd.h>  // gethostname

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* -------------------
   Connection endpoint
   -------------------
    The connection endpoint is used to send/receive data to/from the
    network sublayer. Its definition depends on the underlaying mechanism
    chosen to communicate (network socket, I/O terminal device, async FIFO).
    A connection endpoint is handled using an identifier, and functions
    used to retreive the file descriptor actually allocated by the system,
    to receive data, to send data, and to perform clean up when connection
    is shut down.
    Only one single end to end connection with the network sublayer is
    managed at a time.
*/
static struct {
  /*
   * Connection endpoint reference
   */
  void* endpoint;
  /*
   * Connection endpoint handlers
   */
  void* (*open)(int, const char*, const char*);
  int (*getfd)(const void*);
  ssize_t (*recv)(void*, char*, size_t);
  ssize_t (*send)(const void*, const char*, size_t);
  void (*close)(void*);
} _network_api_id;

#define NETWORK_API_OPEN(a, b, c) _network_api_id.open(a, b, c)
#define NETWORK_API_GETFD() _network_api_id.getfd(_network_api_id.endpoint)
#define NETWORK_API_RECV(a, b)                                                 \
  _network_api_id.recv(_network_api_id.endpoint, a, b)
#define NETWORK_API_SEND(a, b)                                                 \
  _network_api_id.send(_network_api_id.endpoint, a, b)
#define NETWORK_API_CLOSE() _network_api_id.close(_network_api_id.endpoint)

/*
   The buffer used to receive data from the network sublayer
*/
#define NETWORK_API_RECV_BUFFER_SIZE 4096
static char _network_api_recv_buffer[NETWORK_API_RECV_BUFFER_SIZE];

/*
   The buffer used to send data to the network sublayer
*/
#define NETWORK_API_SEND_BUFFER_SIZE NETWORK_API_RECV_BUFFER_SIZE
static char _network_api_send_buffer[NETWORK_API_SEND_BUFFER_SIZE];

/*
   The decoded data received from the network sublayer
*/
static as_message_t _as_data = {}; /* Access Stratum message     */

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_initialize()                                  **
 **                                                                        **
 ** Description: Initializes the network API from which the NAS layer      **
 **    will send/receive messages to/from the network sublayer   **
 **                                                                        **
 ** Inputs:  host:    The name of the host from which the net-   **
 **       work sublayer will connect                 **
 **      port:    The local port number                      **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   Return:  RETURNerror, RETURNok                      **
 **      Others:  _network_api_id                            **
 **                                                                        **
 ***************************************************************************/
int network_api_initialize(const char* host, const char* port) {
  OAILOG_FUNC_IN(LOG_NAS);
  /*
   * Initialize network socket handlers
   */
  _network_api_id.open  = socket_udp_open;
  _network_api_id.getfd = socket_get_fd;
  _network_api_id.recv  = socket_recv;
  _network_api_id.send  = socket_send;
  _network_api_id.close = socket_close;
  /*
   * Initialize UDP communication channel with the network layer
   */
#if NAS_MME
  _network_api_id.endpoint = NETWORK_API_OPEN(SOCKET_SERVER, NULL, port);
#endif

  if (_network_api_id.endpoint == NULL) {
    const char* error = ((errno < 0) ? gai_strerror(errno) : strerror(errno));

    OAILOG_ERROR(
        LOG_NAS, "NET-API   - Failed to open connection endpoint, %s", error);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  gethostname(_network_api_send_buffer, NETWORK_API_SEND_BUFFER_SIZE);
  OAILOG_INFO(
      LOG_NAS, "NET-API   - Network's UDP socket %d is BOUND to %s/%s",
      network_api_get_fd(), _network_api_send_buffer, port);
  OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_get_fd()                                      **
 **                                                                        **
 ** Description: Get the file descriptor of the connection endpoint used   **
 **    to send/receive messages to/from the network sublayer     **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:  _network_api_id                            **
 **                                                                        **
 ** Outputs:   Return:  The file descriptor of the connection end- **
 **       point used by the network sublayer         **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int network_api_get_fd(void) {
  OAILOG_FUNC_IN(LOG_NAS);
  OAILOG_FUNC_RETURN(LOG_NAS, NETWORK_API_GETFD());
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_get_data()                                    **
 **                                                                        **
 ** Description: Get a generic pointer to the network data structure.      **
 **    Casting to the proper type is necessary before its usage. **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:  _as_data                                   **
 **                                                                        **
 ** Outputs:   Return:  A generic pointer to the network data      **
 **       structure                                  **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
const void* network_api_get_data(void) {
  OAILOG_FUNC_IN(LOG_NAS);
  OAILOG_FUNC_RETURN(LOG_NAS, (void*) (&_as_data));
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_read_data()                                   **
 **                                                                        **
 ** Description: Read data received from the network sublayer              **
 **                                                                        **
 ** Inputs:  fd:    File descriptor of the connection endpoint **
 **       from which data have been received         **
 **      Others:  _network_api_id                            **
 **                                                                        **
 ** Outputs:   Return:  The number of bytes read when success;     **
 **       RETURNerror otherwise                      **
 **      Others:  _network_api_recv_buffer, _network_api_id  **
 **                                                                        **
 ***************************************************************************/
int network_api_read_data(int fd) {
  OAILOG_FUNC_IN(LOG_NAS);
  int rbytes;

  /*
   * Sanity check
   */
  int sfd = network_api_get_fd();

  if (fd != sfd) {
    OAILOG_ERROR(
        LOG_NAS,
        "NET-API   - Endpoint %d is not the one created for communication with "
        "the network sublayer (%d)",
        fd, sfd);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  memset(_network_api_recv_buffer, 0, NETWORK_API_RECV_BUFFER_SIZE);
  /*
   * Receive data from the network sublayer
   */
  rbytes =
      NETWORK_API_RECV(_network_api_recv_buffer, NETWORK_API_RECV_BUFFER_SIZE);

  if (rbytes == RETURNerror) {
    OAILOG_ERROR(LOG_NAS, "NET-API   - recv() failed, %s", strerror(errno));
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  } else if (rbytes == 0) {
    OAILOG_WARNING(LOG_NAS, "NET-API   - A signal was caught");
  } else {
    OAILOG_INFO(
        LOG_NAS,
        "NET-API   - %d bytes received from the network "
        "sublayer",
        rbytes);
    OAILOG_STREAM_HEX(
        OAILOG_LEVEL_DEBUG, LOG_NAS, "", _network_api_recv_buffer, rbytes);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, rbytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_send_data()                                   **
 **                                                                        **
 ** Description: Send data to the network sublayer                         **
 **                                                                        **
 ** Inputs:  fd:    File descriptor of the connection endpoint **
 **       to which data have to be sent              **
 **      length:  Number of bytes to send                    **
 **      Others:  _network_api_send_buffer, _network_api_id  **
 **                                                                        **
 ** Outputs:   Return:  The number of bytes sent when success;     **
 **       RETURNerror otherwise                      **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int network_api_send_data(int fd, size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  int sbytes;

  /*
   * Sanity check
   */
  int sfd = network_api_get_fd();

  if (fd != sfd) {
    OAILOG_ERROR(
        LOG_NAS,
        "NET-API   - Endpoint %d is not the one created for communication with "
        "the network sublayer (%d)",
        fd, sfd);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  /*
   * Send data to the network sublayer
   */
  sbytes = NETWORK_API_SEND(_network_api_send_buffer, length);

  if (sbytes == RETURNerror) {
    OAILOG_ERROR(LOG_NAS, "NET-API   - send() failed, %s", strerror(errno));
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  } else if (sbytes == 0) {
    OAILOG_WARNING(LOG_NAS, "NET-API   - A signal was caught");
  } else {
    OAILOG_INFO(
        LOG_NAS, "NET-API   - %d bytes sent to the network sublayer", sbytes);
    OAILOG_STREAM_HEX(
        OAILOG_LEVEL_DEBUG, LOG_NAS, "", _network_api_recv_buffer, sbytes);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, sbytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_close()                                       **
 **                                                                        **
 ** Description: Clean the network API from which the NAS layer sent/recei-**
 **    ved messages to/from the network sublayer                 **
 **                                                                        **
 ** Inputs:  fd:    File descriptor of the connection endpoint **
 **       allocated by the system to communicate     **
 **       with the network sublayer                  **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   Return:  None                                       **
 **      Others:  _network_api_id                            **
 **                                                                        **
 ***************************************************************************/
void network_api_close(int fd) {
  OAILOG_FUNC_IN(LOG_NAS);
  /*
   * Sanity check
   */
  int sfd = network_api_get_fd();

  if (fd != sfd) {
    OAILOG_ERROR(
        LOG_NAS,
        "NET-API   - Endpoint %d is not the one created for communication with "
        "the network sublayer (%d)",
        fd, sfd);
    OAILOG_FUNC_OUT(LOG_NAS);
    return;
  }

  /*
   * Cleanup the connection endpoint
   */
  NETWORK_API_CLOSE();
  _network_api_id.endpoint = NULL;
  OAILOG_FUNC_OUT(LOG_NAS);
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_decode_data()                                 **
 **                                                                        **
 ** Description: Decodes the message received from the network and fills   **
 **    the corresponding network data structure.                 **
 **                                                                        **
 ** Inputs:  length:  Number of bytes to decode                  **
 **      Others:  _network_api_recv_buffer                   **
 **                                                                        **
 ** Outputs:   Return:  The identifier of the received AS          **
 **       message when successfully decoded;         **
 **       RETURNerror otherwise                      **
 **      Others:  _as_data                                   **
 **                                                                        **
 ***************************************************************************/
int network_api_decode_data(size_t length) {
  OAILOG_FUNC_IN(LOG_NAS);
  /*
   * Decode the Access Stratum message received from the network
   */
  int as_id = as_message_decode(_network_api_recv_buffer, &_as_data, length);

  if (as_id != RETURNerror) {
    OAILOG_INFO(
        LOG_NAS, "NET-API   - AS message id=0x%x successfully decoded", as_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, as_id);
}

/****************************************************************************
 **                                                                        **
 ** Name:  network_api_encode_data()                                 **
 **                                                                        **
 ** Description: Encodes Access Stratum message to be sent to the network  **
 **                                                                        **
 ** Inputs   data:    Generic pointer to the AS data to encode   **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   Return:  The number of characters that have been    **
 **       successfully encoded;                      **
 **       RETURNerror otherwise.                     **
 **      Others:  _network_api_send_buffer                   **
 **                                                                        **
 ***************************************************************************/
int network_api_encode_data(void* data) {
  OAILOG_FUNC_IN(LOG_NAS);
  /*
   * Encode the Access Stratum  message
   */
  int bytes = as_message_encode(
      _network_api_send_buffer, (as_message_t*) (data),
      NETWORK_API_SEND_BUFFER_SIZE);

  if (bytes != RETURNerror) {
    OAILOG_INFO(LOG_NAS, "NET-API   - %d bytes encoded", bytes);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:  as_message_send()                                         **
 **                                                                        **
 ** Description: Service provided to the EPS Mobility Management protocol  **
 **    at the EMMAS Access Point (EMMAS-SAP) to send AS messages **
 **    to the Access Stratum sublayer.                           **
 **                                                                        **
 ** Inputs:  as_msg:  The AS message to send                     **
 **      Others:  _network_api_send_buffer, _network_api_id  **
 **                                                                        **
 ** Outputs:   Return:  The number of bytes sent when success;     **
 **       RETURNerror Otherwise                      **
 **      Others:  _network_api_send_buffer                   **
 **                                                                        **
 ***************************************************************************/
int as_message_send(as_message_t* as_msg) {
  int bytes;

  OAILOG_FUNC_IN(LOG_NAS);
  OAILOG_INFO(
      LOG_NAS,
      "NET-API   - Send message 0x%.4x to the Access Stratum "
      "layer",
      as_msg->msg_id);
  /*
   * Encode the AS message
   */
  bytes = network_api_encode_data(as_msg);

  if (bytes > 0) {
    /*
     * Get the network file descriptor
     */
    int fd = network_api_get_fd();

    if (fd != RETURNerror) {
      /*
       * Send the AS message to the network
       */
      bytes = network_api_send_data(fd, bytes);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS, bytes);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
