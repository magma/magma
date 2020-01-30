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
#include <errno.h>
#include <stdbool.h>
#include <string.h>

#if HAVE_CONFIG_H
#include "config.h"
#endif

#include "dynamic_memory_check.h"
#include "assertions.h"
#include "log.h"
#include "mme_config.h"
#include "shared_ts_log.h"

#include "intertask_interface_init.h"
#include "sctp_primitives_server.h"
#include "s1ap_mme.h"
#include "mme_app_extern.h"
/* FreeDiameter headers for support of S6A interface */
#include "s6a_defs.h"
#include "sgs_defs.h"
#include "oai_mme.h"
#include "pid_file.h"
#include "service303_message_utils.h"
#include "mme_app_embedded_spgw.h"
#include "bstrlib.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "pgw_defs.h"
#include "service303.h"
#include "sgw_defs.h"
#include "shared_ts_log.h"
#include "spgw_config.h"
#include "grpc_service.h"

int main(int argc, char *argv[])
{
  char *pid_file_name;

  CHECK_INIT_RETURN(OAILOG_INIT(
    MME_CONFIG_STRING_MME_CONFIG, OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(shared_log_init(MAX_LOG_PROTOS));
  CHECK_INIT_RETURN(itti_init(
    TASK_MAX,
    THREAD_MAX,
    MESSAGES_ID_MAX,
    tasks_info,
    messages_info,
    NULL,
    NULL));

  /*
   * Parse the command line for options and set the mme_config accordingly.
   */
#if EMBEDDED_SGW
  CHECK_INIT_RETURN(mme_config_embedded_spgw_parse_opt_line(
    argc, argv, &mme_config, &spgw_config));
#else
  CHECK_INIT_RETURN(mme_config_parse_opt_line(argc, argv, &mme_config));
#endif

  pid_file_name = get_pid_file_name(mme_config.pid_dir);

  if (!pid_file_lock(pid_file_name)) {
    exit(-EDEADLK);
  }
  free_wrapper((void **) &pid_file_name);

  /*
   * Calling each layer init function
   */
  // Intialize loggers and configured log levels.
  OAILOG_LOG_CONFIGURE(&mme_config.log_config);
  CHECK_INIT_RETURN(service303_init(&(mme_config.service303_config)));

  // Service started, but not healthy yet
  send_app_health_to_service303(TASK_MME_APP, false);

  CHECK_INIT_RETURN(mme_app_init(&mme_config));
  CHECK_INIT_RETURN(sctp_init(&mme_config));
#if EMBEDDED_SGW
  CHECK_INIT_RETURN(sgw_init(&spgw_config, mme_config.use_stateless));
  CHECK_INIT_RETURN(pgw_init(&spgw_config));
#else
  CHECK_INIT_RETURN(s11_mme_init(&mme_config));
#endif
  CHECK_INIT_RETURN(s1ap_mme_init(&mme_config));
  CHECK_INIT_RETURN(s6a_init(&mme_config));

  //Create SGS Task only if non_eps_service_control is not set to OFF
  char *non_eps_service_control = bdata(mme_config.non_eps_service_control);
  if (
    !(strcmp(non_eps_service_control, "SMS")) ||
    !(strcmp(non_eps_service_control, "CSFB_SMS"))) {
    CHECK_INIT_RETURN(sgs_init(&mme_config));
    OAILOG_DEBUG(LOG_MME_APP, "SGS Task initialized\n");
  }
  CHECK_INIT_RETURN(grpc_service_init());
  OAILOG_DEBUG(LOG_MME_APP, "MME app initialization complete\n");

#if EMBEDDED_SGW
  /*
   * Display the configuration
   */
  mme_config_display(&mme_config);
  spgw_config_display(&spgw_config);
#endif

  /*
   * Handle signals here
   */
  itti_wait_tasks_end();
#if EMBEDDED_SGW
  free_spgw_config(&spgw_config);
#endif
  pid_file_unlock();

  return 0;
}
