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

/*! \file mme_app_sgw_selection.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <inttypes.h>
#include <netinet/in.h>

#include "bstrlib.h"

#include "log.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "TrackingAreaIdentity.h"
#include "mme_config.h"
#include "mme_app_sgw_selection.h"
#include "mme_app_edns_emulation.h"

//------------------------------------------------------------------------------
void mme_app_select_sgw(
  const tai_t *const tai,
  struct in_addr *const sgw_in_addr)
{
  // see in 3GPP TS 29.303 version 10.5.0 Release 10:
  // 5.2 Procedures for Discovering and Selecting a SGW
  // ....

  // do it the simplest way for now
  //
  char tmp[8];
  bstring application_unique_string = bfromcstr("tac-lb");
  if (0 < snprintf(tmp, 8, "%02x", tai->tac & 0x00FF)) {
    bcatcstr(application_unique_string, tmp);
  } else {
    sgw_in_addr->s_addr = 0;
    return;
  }
  bcatcstr(application_unique_string, ".tac-hb");
  if (0 < snprintf(tmp, 8, "%02x", tai->tac >> 8)) {
    bcatcstr(application_unique_string, tmp);
  } else {
    goto lookup_error;
  }
  bcatcstr(application_unique_string, ".tac.epc.mnc");
  uint16_t mnc = (tai->mnc_digit1 * 10) + tai->mnc_digit2;
  if (10 > tai->mnc_digit3) {
    mnc = (mnc * 10) + tai->mnc_digit3;
  }
  if (0 < snprintf(tmp, 8, "%03u", mnc)) {
    bcatcstr(application_unique_string, tmp);
  } else {
    goto lookup_error;
  }
  bcatcstr(application_unique_string, ".mcc");
  if (
    0 <
    snprintf(
      tmp, 8, "%u%u%u", tai->mcc_digit1, tai->mcc_digit2, tai->mcc_digit3)) {
    bcatcstr(application_unique_string, tmp);
  } else {
    goto lookup_error;
  }
  bcatcstr(application_unique_string, ".3gppnetwork.org");

  struct in_addr *entry = mme_app_edns_get_sgw_entry(application_unique_string);

  if (entry) {
    sgw_in_addr->s_addr = entry->s_addr;
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    "SGW lookup %s returned %s\n",
    application_unique_string->data,
    inet_ntoa(*sgw_in_addr));
  bdestroy_wrapper(&application_unique_string);
  return;

lookup_error:
  OAILOG_WARNING(
    LOG_MME_APP, "Failed SGW lookup for TAI " TAI_FMT "\n", TAI_ARG(tai));
  sgw_in_addr->s_addr = 0;
  bdestroy_wrapper(&application_unique_string);
  return;
}
