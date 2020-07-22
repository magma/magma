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

/*! \file s11_causes.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define S11_CAUSES_C

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "s11_causes.h"
#include "3gpp_29.274.h"

static const SGWCauseMapping_t causes[] = {
    {LOCAL_DETACH, "Local detach", 0, 0, 0, 0},
    {COMPLETE_DETACH, "Complete detach", 0, 0, 0, 0},
    {RAT_CHANGE_3GPP_TO_NON_3GPP, "From 3GPP to non 3GPP RAT change", 0, 0, 0,
     0},
    {ISR_DEACTIVATION, "ISR deactivation", 0, 0, 0, 0},
    {ERROR_IND_FROM_RNC_ENB_SGSN, "Error ind received from RNC/eNB/SGSN", 0, 0,
     0, 0},
    {IMSI_DETACH_ONLY, "IMSI detach only", 0, 0, 0, 0},
    {REQUEST_ACCEPTED, "Request accepted", 1, 1, 1, 1},
    {REQUEST_ACCEPTED_PARTIALLY, "Request accepted partially", 1, 1, 1, 0},
    {NEW_PDN_TYPE_NW_PREF, "New PDN type network preference", 1, 0, 0, 0},
    {NEW_PDN_TYPE_SAB_ONLY, "New PDN type single address bearer only", 1, 0, 0,
     0},
    {CONTEXT_NOT_FOUND, "Context not found", 0, 1, 1, 1},
    {INVALID_MESSAGE_FORMAT, "Invalid message format", 1, 1, 1, 1},
    {INVALID_LENGTH, "Invalid length", 1, 1, 1, 0},
    {SERVICE_NOT_SUPPORTED, "Service not supported", 0, 0, 1, 0},
    {SYSTEM_FAILURE, "System failure", 1, 1, 1, 0},
    {NO_RESOURCES_AVAILABLE, "No resources available", 1, 0, 0, 0},
    {MISSING_OR_UNKNOWN_APN, "Missing or unknown APN", 1, 0, 0, 0},
    {GRE_KEY_NOT_FOUND, "GRE KEY not found", 1, 0, 0, 0},
    {DENIED_IN_RAT, "Denied in RAT", 1, 1, 0, 0},
    {UE_NOT_RESPONDING, "UE not responding", 0, 1, 0, 0},
    {SERVICE_DENIED, "Service Denied", 0, 0, 0, 0},
    {UNABLE_TO_PAGE_UE, "Unable to page UE", 0, 1, 0, 0},
    {NO_MEMORY_AVAILABLE, "No memory available", 1, 1, 1, 0},
    {REQUEST_REJECTED, "Request rejected", 1, 1, 1, 0},
    {INVALID_PEER, "Invalid peer", 0, 0, 0, 1},
    {TEMP_REJECT_HO_IN_PROGRESS, "Temporarily rejected due to HO in progress",
     0, 0, 0, 0},
    {M_PDN_APN_NOT_ALLOWED, "Multiple PDN for a given APN not allowed", 1, 0, 0,
     0},
    {0, NULL, 0, 0, 0, 0},
};

static int compare_cause_id(const void* m1, const void* m2) {
  struct SGWCauseMapping_e* scm1 = (struct SGWCauseMapping_e*) m1;
  struct SGWCauseMapping_e* scm2 = (struct SGWCauseMapping_e*) m2;

  return (scm1->value - scm2->value);
}

char* sgw_cause_2_string(uint8_t cause_value) {
  SGWCauseMapping_t *res, key;

  key.value = cause_value;
  res       = bsearch(
      &key, causes, sizeof(causes), sizeof(SGWCauseMapping_t),
      compare_cause_id);

  if (res == NULL) {
    return "Unknown cause";
  } else {
    return res->name;
  }
}
