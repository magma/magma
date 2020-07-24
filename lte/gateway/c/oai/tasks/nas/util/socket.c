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
  Source   socket.c

  Version    0.1

  Date   2012/02/28

  Product    NAS stack

  Subsystem  Utilities

  Author   Frederic Maurel

  Description  Implements socket handlers

 *****************************************************************************/

#include "socket.h"
#include "nas/commonDef.h"

#include <stdlib.h>  // malloc, free, atoi
#include <string.h>  // memset
#include <unistd.h>  // close
#include <errno.h>   // EINTR
#include <sys/types.h>
#include <sys/socket.h>  // socket, setsockopt, connect, bind, recv, send
#include <netdb.h>       // getaddrinfo
#include "dynamic_memory_check.h"
/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* -------------------------------------
   Identifier of network socket endpoint
   -------------------------------------
    A network socket endpoint is defined with a type (client or server),
    a port number, and the name and address of the remote host it should
    be connected to. A socket file descriptor is created to handle the
    setup of the communication channel with the remote peer.
*/
struct socket_id_s {
  int type; /* connection type (client/server)  */
  int port; /* port number      */
#define SOCKET_HOSTNAME_SIZE 32
  char rhost[SOCKET_HOSTNAME_SIZE]; /* remote hostname    */
  struct sockaddr_storage addr;     /* remote address   */
  int fd;                           /* socket file descriptor */
};

/* Set socket option at the sockets API level (SOL_SOCKET) */
static int _socket_set_option(int sfd);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:  socket_udp_open()                                         **
 **                                                                        **
 ** Description: Initializes a communication endpoint for any IPv4 and     **
 **    IPv6 protocol used by a local host to send/receive UDP    **
 **    datagrams to/from a remote host in "connected" mode of    **
 **    operation.                                                **
 **                                                                        **
 ** Inputs:  type:    Type of the connection (client/server)     **
 **      host:    For client application, the remote host-   **
 **       name, or numerical IPv4 or IPv6 network    **
 **       address to connect to. NULL for server     **
 **       application.                               **
 **      port:    The port number                            **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   None                                                      **
 **      Return:  A pointer to the local endpoint identifier **
 **       that has been allocated for communication  **
 **       with the remote peer. NULL if the setup of **
 **       the communication endpoint failed.         **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
void* socket_udp_open(int type, const char* host, const char* port) {
  struct addrinfo socket_info;       /* endpoint information    */
  struct addrinfo *socket_addr, *sp; /* endpoint address    */
  int sfd;                           /* socket file descriptor  */

  /*
   * Parameters sanity check
   * -----------------------
   * The local connection endpoint shall be of type CLIENT or SERVER
   */
  if (host == NULL) {
    type = SOCKET_SERVER;
  } else if (type != SOCKET_CLIENT && type != SOCKET_SERVER) {
    return NULL;
  }

  /*
   * Initialize the endpoint address information
   * -------------------------------------------
   * The AI_PASSIVE flag allows "wildcard address" when hostname is not
   * provided (NULL). The wildcard address is used by applications (ty-
   * pically servers) that intend to accept connections on any of the
   * hosts's network addresses. If the hostname is not NULL, then the
   * AI_PASSIVE flag is ignored.
   * When the AI_V4MAPPED flag is set, and AF_INET6 is specified, and no
   * matching IPv6 addresses could be found, then IPv4-mapped IPv6 addresses
   * are returned by getaddrinfo in the list pointed to by result.
   */
  memset(&socket_info, 0, sizeof(struct addrinfo));
  socket_info.ai_socktype = SOCK_DGRAM;     /* Datagram socket   */
  socket_info.ai_flags    = AI_NUMERICSERV; /* numeric port number  */

  if (type != SOCKET_CLIENT) {
    /*
     * Setup socket address options at the server side
     */
    socket_info.ai_family = AF_INET6;    /* Accept either IPv4 or IPv6
                                          * connections     */
    socket_info.ai_flags |= AI_PASSIVE;  /* Use "wildcard address"  */
    socket_info.ai_flags |= AI_V4MAPPED; /* IPv4-mapped IPv6 address */
  } else {
    /*
     * Setup socket address options at the client side
     */
    socket_info.ai_family = AF_INET; /* Any address family   */
    //         socket_info.ai_flags |= AI_V4MAPPED; /* IPv4-mapped IPv6 address
    //         */
  }

  /*
   * getaddrinfo() returns a linked list of address structures:
   * - The network host may be multi-homed, accessible over multiple
   * protocols (e.g. both AF_INET and AF_INET6);
   * - The same service (port number) may be available from multiple
   * socket types (one SOCK_STREAM address and another SOCK_DGRAM address);
   */
  int rc = getaddrinfo(host, port, &socket_info, &socket_addr);

  if (rc != 0) {
    if (rc != EAI_SYSTEM) {
      errno = rc;
    }

    return NULL;
  }

  /*
   * Try each address until we successfully connect
   */
  for (sp = socket_addr; sp; sp = sp->ai_next) {
    /*
     * Create the socket endpoint for communication
     */
    sfd = socket(sp->ai_family, sp->ai_socktype, sp->ai_protocol);

    if (sfd < 0) {
      continue;
    }

    /*
     * Initiate a communication channel at the CLIENT side
     */
    if (type == SOCKET_CLIENT) {
      /*
       * Connect the socket to the remote server's address
       */
      if (connect(sfd, sp->ai_addr, sp->ai_addrlen) != -1) {
        break; /* Connection succeed */
      }
    }
    /*
     * Initiate a communication channel at the SERVER side
     */
    else {
      if (type == SOCKET_SERVER) {
        /*
         * Set socket options
         */
        if (_socket_set_option(sfd) != RETURNok) {
          continue;
        }

        /*
         * Bind the socket to the local server's address
         */
        if (bind(sfd, sp->ai_addr, sp->ai_addrlen) != -1) {
          break; /* Bind succeed */
        }
      }
    }

    close(sfd);
  }

  /*
   * Free the memory that was dynamically allocated for the linked list
   */
  freeaddrinfo(socket_addr);

  if (sp == NULL) {
    /*
     * Connect or bind failed
     */
    return NULL;
  }

  /*
   * The connection endpoint has been successfully setup
   */
  socket_id_t* sid = (socket_id_t*) malloc(sizeof(struct socket_id_s));

  if (sid) {
    sid->type = type;
    sid->port = atoi(port);
    sid->fd   = sfd;
  }

  return sid;
}

/****************************************************************************
 **                                                                        **
 ** Name:  socket_close()                                            **
 **                                                                        **
 ** Description: Cleanup the specified communication endpoint: Close the   **
 **    socket file descriptor and frees all the memory space     **
 **    allocated to handle the communication channel towards the **
 **    remote peer                                               **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   id:    The identifier of the connection endpoint  **
 **      Return:  None                                       **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
void socket_close(void* id) {
  if (id) {
    close(((socket_id_t*) id)->fd);
    free(id);
  }
}

/****************************************************************************
 **                                                                        **
 ** Name:  socket_recv()                                             **
 **                                                                        **
 ** Description: Receives data from a given communication endpoint to a    **
 **    receive buffer of specified length                        **
 **                                                                        **
 ** Inputs:  id:    The identifier of the connection endpoint  **
 **      length:  Length of the receive buffer               **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   id:    The identifier of the connection endpoint  **
 **      buffer:  The receive buffer                         **
 **      Return:  The number of bytes received when success; **
 **       RETURNerror otherwise.                     **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
ssize_t socket_recv(void* id, char* buffer, size_t length) {
  socket_id_t* sid = (socket_id_t*) (id);
  ssize_t rbytes   = -1;

  if (sid->type == SOCKET_CLIENT) {
    /*
     * Receive data from the connected socket
     */
    rbytes = recv(sid->fd, buffer, length, 0);
  } else if (sid->type == SOCKET_SERVER) {
    struct sockaddr_storage addr;
    socklen_t addrlen = sizeof(addr);

    /*
     * Receive data from the socket and retreive the remote host address
     */
    rbytes = recvfrom(
        sid->fd, buffer, length, 0, (struct sockaddr*) &addr, &addrlen);
    sid->addr = addr;
  }

  if (errno == EINTR) {
    /*
     * A signal was caught
     */
    return 0;
  } else if (rbytes < 0) {
    /*
     * Receive failed
     */
    return RETURNerror;
  }

  return rbytes;
}

/****************************************************************************
 **                                                                        **
 ** Name:  socket_send()                                             **
 **                                                                        **
 ** Description: Sends data to a given communication endpoint from a send  **
 **    buffer of specified length                                **
 **                                                                        **
 ** Inputs:  id:    The identifier of the connection endpoint  **
 **    buffer:  The send buffer                            **
 **      length:  Length of the send buffer                  **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   None                                                      **
 **      Return:  The number of bytes sent when success;     **
 **       RETURNerror otherwise.                     **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
ssize_t socket_send(const void* id, const char* buffer, size_t length) {
  const socket_id_t* sid = (socket_id_t*) (id);
  ssize_t sbytes         = -1;

  if (sid->type == SOCKET_CLIENT) {
    /*
     * Send data to the connected socket
     */
    sbytes = send(sid->fd, buffer, length, 0);
  } else if (sid->type == SOCKET_SERVER) {
    /*
     * Send data to the socket using the remote host address
     */
    sbytes = sendto(
        sid->fd, buffer, length, 0, (struct sockaddr*) &sid->addr,
        (socklen_t) sizeof(sid->addr));
  }

  if (errno == EINTR) {
    /*
     * A signal was caught
     */
    return 0;
  } else if (sbytes != length) {
    /*
     * Send failed
     */
    return RETURNerror;
  }

  return sbytes;
}

/****************************************************************************
 **                                                                        **
 ** Name:  socket_get_fd()                                           **
 **                                                                        **
 ** Description: Get the value of the socket file descriptor created for   **
 **    the given connection endpoint                             **
 **                                                                        **
 ** Inputs:  id:    The identifier of the connection endpoint  **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   None                                                      **
 **      Return:  The file descriptor of the socket created  **
 **       for this connection endpoint               **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
int socket_get_fd(const void* id) {
  if (id) {
    return ((socket_id_t*) id)->fd;
  }

  return RETURNerror;
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:  _socket_set_option()                                      **
 **                                                                        **
 ** Description: Set socket option at the sockets API level (SOL_SOCKET)   **
 **                                                                        **
 ** Inputs:  sfd:   Socket file descriptor                     **
 **      Others:  None                                       **
 **                                                                        **
 ** Outputs:   None                                                      **
 **      Return:  RETURNok, RETURNerror                      **
 **      Others:  None                                       **
 **                                                                        **
 ***************************************************************************/
static int _socket_set_option(int sfd) {
  int optval;

  /*
   * SO_REUSEADDR socket option:
   * * * * ---------------------------
   * * * * Allows the server to bind a socket to this port, unless
   * * * * there is an active listening socket already bound to the
   * * * * port. This is useful when recovering from a crash and the
   * * * * socket was not properly closed. The server can be restarted
   * * * * and it will simply open another socket on the same port and
   * * * * continue listening.
   */
  optval = true;

  if (setsockopt(sfd, SOL_SOCKET, SO_REUSEADDR, &optval, sizeof(optval)) < 0) {
    return RETURNerror;
  }

  /*
   * IPV6_V6ONLY socket option
   * * * * -------------------------
   * * * * When option is set to true, the socket is restricted to sending and
   * * * * receiving IPv6 packets only.
   * * * * When option is set to false, the socket can be used to send and
   * receive
   * * * * packets to and from an IPv6 address or an IPv4-mapped IPv6 address.
   */
  optval = false;

  if (setsockopt(sfd, IPPROTO_IPV6, IPV6_V6ONLY, &optval, sizeof(optval)) < 0) {
    return RETURNerror;
  }

  return RETURNok;
}
