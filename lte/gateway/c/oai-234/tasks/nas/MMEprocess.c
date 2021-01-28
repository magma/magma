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
  Source     MMEprocess.c

  Version        0.1

  Date       2013/02/26

  Product        NAS stack

  Subsystem  MME NAS main process

  Author     Frederic Maurel

  Description    Implements the Non-Access Stratum protocol for Evolved Packet
  system (EPS) running at the Network side.

 *****************************************************************************/

#include "nas/commonDef.h"
#include "log.h"
#include "nas_timer.h"

#include "network_api.h"
#include "nas_network.h"
#include "nas_parser.h"

#include <stdlib.h>  // exit
#include <poll.h>    // poll
#include <string.h>  // memset
#include <signal.h>  // sigaction
#include <pthread.h>

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

#define NAS_SLEEP_TIMEOUT 1000 /* 1 second */

static void* _nas_network_mngr(void*);

static int _nas_set_signal_handler(int signal, void(handler)(int));
static void _nas_signal_handler(int signal);

static void _nas_clean(int net_fd);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************/
int main(int argc, const char* argv[]) {
  /*
   * Get the command line options
   */
  if (nas_parser_get_options(argc, argv) != RETURNok) {
    nas_parser_print_usage(FIRMWARE_VERSION);
    exit(EXIT_FAILURE);
  }

  /*
   * Initialize logging trace utility
   */
  nas_log_init(nas_parser_get_trace_level());
  const char* nhost = nas_parser_get_network_host();
  const char* nport = nas_parser_get_network_port();

  OAILOG_TRACE(
      INFO, "MME-MAIN  - %s -nhost %s -nport %s -trace 0x%x", argv[0], nhost,
      nport, nas_parser_get_trace_level());

  /*
   * Initialize the Network interface
   */
  if (network_api_initialize(nhost, nport) != RETURNok) {
    OAILOG_TRACE(ERROR, "MME-MAIN  - network_api_initialize() failed");
    exit(EXIT_FAILURE);
  }

  int network_fd = network_api_get_fd();

  /*
   * Initialize the NAS contexts
   */
  nas_network_initialize();
  /*
   * Initialize NAS timer handlers
   */
  nas_timer_init();
  /*
   * Set up signal handlers
   */
  (void) _nas_set_signal_handler(SIGINT, _nas_signal_handler);
  (void) _nas_set_signal_handler(SIGTERM, _nas_signal_handler);
  pthread_attr_t attr;

  pthread_attr_init(&attr);
  pthread_attr_setscope(&attr, PTHREAD_SCOPE_SYSTEM);
  pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
  /*
   * Start thread use to manage the network connection endpoint
   */
  pthread_t network_mngr;

  if (pthread_create(&network_mngr, &attr, _nas_network_mngr, &network_fd) !=
      0) {
    OAILOG_TRACE(
        ERROR,
        "MME-MAIN  - "
        "Failed to create the network management thread");
    network_api_close(network_fd);
    exit(EXIT_FAILURE);
  }

  pthread_attr_destroy(&attr);

  /*
   * Suspend execution of the main process until the network connection
   * endpoint is still active
   */
  while (network_fd != -1) {
    poll(NULL, 0, NAS_SLEEP_TIMEOUT);
    network_fd = network_api_get_fd();
  }

  /*
   * Termination cleanup
   */
  _nas_clean(network_fd);
  OAILOG_TRACE(WARNING, "MME-MAIN  - NAS main process exited");
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    _nas_network_mngr()                                       **
 **                                                                        **
 ** Description: Manages the connection endpoint use to communicate with   **
 **      the network sublayer                                      **
 **                                                                        **
 ** Inputs:  fd:        The descriptor of the network connection   **
 **             endpoint                                   **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     Return:    None                                       **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static void* _nas_network_mngr(void* args) {
  OAILOG_FUNC_IN;
  int ret_code;
  int network_message_id;
  int bytes;
  int* fd = (int*) args;

  OAILOG_TRACE(
      INFO, "MME-MAIN  - Network connection manager started (%d)", *fd);

  /*
   * Network receiving loop
   */
  while (true) {
    /*
     * Read the network data message
     */
    bytes = network_api_read_data(*fd);

    if (bytes == RETURNerror) {
      /*
       * Failed to read data from the network sublayer;
       * * * * exit from the receiving loop
       */
      OAILOG_TRACE(
          ERROR,
          "MME-MAIN  - "
          "Failed to read data from the network sublayer");
      break;
    }

    if (bytes == 0) {
      /*
       * A signal was caught before any data were available
       */
      continue;
    }

    /*
     * Decode the network data message
     */
    network_message_id = network_api_decode_data(bytes);

    if (network_message_id == RETURNerror) {
      /*
       * Failed to decode data read from the network sublayer
       */
      continue;
    }

    /*
     * Process the network data message
     */
    ret_code =
        nas_network_process_data(network_message_id, network_api_get_data());

    if (ret_code != RETURNok) {
      /*
       * The network data message has not been successfully
       * * * * processed
       */
      OAILOG_TRACE(
          WARNING,
          "MME-MAIN  - "
          "The network procedure call 0x%x failed",
          network_message_id);
    }
  }

  /*
   * Close the connection to the network sublayer
   */
  OAILOG_TRACE(
      WARNING,
      "MME-MAIN  - "
      "The network connection endpoint manager exited");
  OAILOG_FUNC_RETURN(NULL);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _nas_set_signal_handler()                                 **
 **                                                                        **
 ** Description: Set up a signal handler                                   **
 **                                                                        **
 ** Inputs:  signal:    Signal number                              **
 **          handler:   Signal handler                             **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     Return:    RETURNerror, RETURNok                      **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int _nas_set_signal_handler(int signal, void(handler)(int)) {
  OAILOG_FUNC_IN;
  struct sigaction act;

  /*
   * Initialize signal set
   */
  (void) memset(&act, 0, sizeof(act));
  (void) sigfillset(&act.sa_mask);
  (void) sigdelset(&act.sa_mask, SIGHUP);
  (void) sigdelset(&act.sa_mask, SIGINT);
  (void) sigdelset(&act.sa_mask, SIGTERM);
  (void) sigdelset(&act.sa_mask, SIGILL);
  (void) sigdelset(&act.sa_mask, SIGTRAP);
  (void) sigdelset(&act.sa_mask, SIGIOT);
#ifndef LINUX
  (void) sigdelset(&act.sa_mask, SIGEMT);
#endif
  (void) sigdelset(&act.sa_mask, SIGFPE);
  (void) sigdelset(&act.sa_mask, SIGBUS);
  (void) sigdelset(&act.sa_mask, SIGSEGV);
  (void) sigdelset(&act.sa_mask, SIGSYS);
  /*
   * Initialize signal handler
   */
  act.sa_handler = handler;

  if (sigaction(signal, &act, 0) < 0) {
    return RETURNerror;
  }

  OAILOG_TRACE(
      INFO, "MME-MAIN  - Handler successfully set for signal %d", signal);
  OAILOG_FUNC_RETURN(RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _nas_signal_handler()                                     **
 **                                                                        **
 ** Description: Signal handler                                            **
 **                                                                        **
 ** Inputs:  signal:    Signal number                              **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     Return:    None                                       **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static void _nas_signal_handler(int signal) {
  OAILOG_FUNC_IN;
  OAILOG_TRACE(WARNING, "MME-MAIN  - Signal %d received", signal);
  _nas_clean(network_api_get_fd());
  exit(EXIT_SUCCESS);
  OAILOG_FUNC_OUT;
}

/****************************************************************************
 **                                                                        **
 ** Name:    _nas_clean()                                              **
 **                                                                        **
 ** Description: Performs termination cleanup                              **
 **                                                                        **
 ** Inputs:  net_fd:    Network's connection file descriptor       **
 **          Others:    None                                       **
 **                                                                        **
 ** Outputs:     Return:    None                                       **
 **          Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static void _nas_clean(int net_fd) {
  OAILOG_FUNC_IN;
  OAILOG_TRACE(INFO, "MME-MAIN  - Perform EMM and ESM cleanup");
  nas_network_cleanup();
  OAILOG_TRACE(INFO, "MME-MAIN  - Closing network connection %d", net_fd);
  network_api_close(net_fd);
  OAILOG_FUNC_OUT;
}
