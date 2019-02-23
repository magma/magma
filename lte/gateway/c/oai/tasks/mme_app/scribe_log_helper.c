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

#include "conversions.h"
#include "common_types.h"
#include "log.h"
#include "scribe_rpc_client.h"

void log_ue_state_to_scribe(
  imsi64_t imsi,
  uint8_t imsi_len,
  const char *ue_status)
{
  char imsi_str[16];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  scribe_string_param_t str_params[] = {
    {"ue_status", ue_status},
    {"imsi", imsi_str},
  };
  char const *category = "perfpipe_magma_ue_stats";
  int status = log_to_scribe(category, NULL, 0, str_params, 2);
  if (status != 0) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to log to scribe  category %s, log status: %d \n",
      category,
      status);
  }
}
