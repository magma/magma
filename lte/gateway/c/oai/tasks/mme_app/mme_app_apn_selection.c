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

/*! \file mme_app_apn_selection.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>

#include "bstrlib.h"

#include "dynamic_memory_check.h"
#include "log.h"
#include "msc.h"
#include "assertions.h"
#include "conversions.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "mme_app_extern.h"
#include "mme_app_ue_context.h"
#include "common_defs.h"
#include "mme_app_apn_selection.h"

//------------------------------------------------------------------------------
struct apn_configuration_s *mme_app_select_apn(
  ue_mm_context_t *const ue_context,
  const_bstring const ue_selected_apn)
{
  context_identifier_t default_context_identifier =
    ue_context->apn_config_profile.context_identifier;
  int index;

  for (index = 0; index < ue_context->apn_config_profile.nb_apns; index++) {
    if (!ue_selected_apn) {
      /*
       * OK we got our default APN
       */
      if (
        ue_context->apn_config_profile.apn_configuration[index]
          .context_identifier == default_context_identifier) {
        OAILOG_DEBUG(
          LOG_MME_APP,
          "Selected APN %s for UE " IMSI_64_FMT "\n",
          ue_context->apn_config_profile.apn_configuration[index]
            .service_selection,
          ue_context->emm_context._imsi64);
        return &ue_context->apn_config_profile.apn_configuration[index];
      }
    } else {
      /*
       * OK we got the UE selected APN
       */
      if (
        biseqcaselessblk(
          ue_selected_apn,
          ue_context->apn_config_profile.apn_configuration[index]
            .service_selection,
          strlen(ue_context->apn_config_profile.apn_configuration[index]
                   .service_selection)) == 1) {
        OAILOG_DEBUG(
          LOG_MME_APP,
          "Selected APN %s for UE " IMSI_64_FMT "\n",
          ue_context->apn_config_profile.apn_configuration[index]
            .service_selection,
          ue_context->emm_context._imsi64);
        return &ue_context->apn_config_profile.apn_configuration[index];
      }
    }
  }

  return NULL;
}

//------------------------------------------------------------------------------
struct apn_configuration_s *mme_app_get_apn_config(
  ue_mm_context_t *const ue_context,
  const context_identifier_t context_identifier)
{
  int index;

  for (index = 0; index < ue_context->apn_config_profile.nb_apns; index++) {
    if (
      ue_context->apn_config_profile.apn_configuration[index]
        .context_identifier == context_identifier) {
      return &ue_context->apn_config_profile.apn_configuration[index];
    }
  }
  return NULL;
}
