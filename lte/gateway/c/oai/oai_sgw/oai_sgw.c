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

/*! \file oai_sgw.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#if HAVE_CONFIG_H
#include "config.h"
#endif

#include "dynamic_memory_check.h"
#include "assertions.h"
#include "log.h"
#include "intertask_interface_init.h"
#include "spgw_config.h"
#include "udp_primitives_server.h"
#include "s11_sgw.h"
#include "sgw_defs.h"
#include "pgw_defs.h"
#include "oai_sgw.h"
#include "pid_file.h"
#include "service303.h"
#include "service303_message_utils.h"
#include "daemonize.h"

int main(int argc, char *argv[])
{
  char *pid_file_name;

#if DAEMONIZE
  daemon_start();
#endif /* DAEMONIZE */

  pid_file_name = get_pid_file_name(NULL);

  if (!pid_file_lock(pid_file_name)) {
#if DAEMONIZE
    daemon_stop();
#endif /* DAEMONIZE */
    exit(-EDEADLK);
  }
  free_wrapper((void **) &pid_file_name);

  CHECK_INIT_RETURN(OAILOG_INIT(
    SGW_CONFIG_STRING_SGW_CONFIG, OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(shared_log_init(MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(async_system_init());
  /*
   * Parse the command line for options and set the mme_config accordingly.
   */
  CHECK_INIT_RETURN(spgw_config_parse_opt_line(argc, argv, &spgw_config));
  /*
   * Calling each layer init function
   */
  CHECK_INIT_RETURN(itti_init(
    TASK_MAX,
    THREAD_MAX,
    MESSAGES_ID_MAX,
    tasks_info,
    messages_info,
    NULL,
    NULL));
  OAILOG_LOG_CONFIGURE(&spgw_config.sgw_config.log_config);

  CHECK_INIT_RETURN(service303_init(&(spgw_config.service303_config)));

  // Initialize grpc clients
  grpc_start_spgw_client_receivers();

  // Tell service303 that spgw is unhealthy
  send_app_health_to_service303(TASK_SPGW_APP, false);

  CHECK_INIT_RETURN(udp_init());
  CHECK_INIT_RETURN(s11_sgw_init(&spgw_config.sgw_config));
  //CHECK_INIT_RETURN (gtpv1u_init (&spgw_config));
  CHECK_INIT_RETURN(sgw_init(&spgw_config));
  CHECK_INIT_RETURN(pgw_init(&spgw_config));
  // Finished setup, notify service303 that spgw is healthy
  send_app_health_to_service303(TASK_SPGW_APP, true);
  /*
   * Handle signals here
   */
  itti_wait_tasks_end();
  pid_file_unlock();
  return 0;
}
