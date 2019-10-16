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
#define spgw_service
#define spgw_service_TASK_C

#include <stddef.h>

#include "assertions.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "log.h"
#include "mme_default_values.h"
#include "spgw_service.h"

static void *spgw_service_task(void *args)
{
  itti_mark_task_ready(TASK_SPGW_SERVICE);
  spgw_service_data_t *spgw_service_data = (spgw_service_data_t *) args;
  start_spgw_service(spgw_service_data->server_address);
  itti_exit_task();
  return NULL;
}

int spgw_service_init(void)
{
  OAILOG_DEBUG(LOG_SPGW_APP, "Initializing spgw_service task interface\n");
  spgw_service_data_t spgw_service_config;
  spgw_service_config.server_address = bfromcstr(SPGWSERVICE_SERVER_ADDRESS);

  if (
    itti_create_task(
      TASK_SPGW_SERVICE, &spgw_service_task, &spgw_service_config) < 0) {
    OAILOG_ALERT(LOG_SPGW_APP, "Initializing spgw_service: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}
