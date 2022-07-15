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

/*! \file s11_causes.cpp
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define S11_CAUSES_C

#include "lte/gateway/c/core/oai/tasks/sgw/s11_causes.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"

static const SGWCauseMapping_t causes[] = {
    {LOCAL_DETACH, const_cast<char*>("Local detach"), 0, 0, 0, 0},
    {COMPLETE_DETACH, const_cast<char*>("Complete detach"), 0, 0, 0, 0},
    {RAT_CHANGE_3GPP_TO_NON_3GPP,
     const_cast<char*>("From 3GPP to non 3GPP RAT change"), 0, 0, 0, 0},
    {ISR_DEACTIVATION, const_cast<char*>("ISR deactivation"), 0, 0, 0, 0},
    {ERROR_IND_FROM_RNC_ENB_SGSN,
     const_cast<char*>("Error ind received from RNC/eNB/SGSN"), 0, 0, 0, 0},
    {IMSI_DETACH_ONLY, const_cast<char*>("IMSI detach only"), 0, 0, 0, 0},
    {REQUEST_ACCEPTED, const_cast<char*>("Request accepted"), 1, 1, 1, 1},
    {REQUEST_ACCEPTED_PARTIALLY,
     const_cast<char*>("Request accepted partially"), 1, 1, 1, 0},
    {NEW_PDN_TYPE_NW_PREF, const_cast<char*>("New PDN type network preference"),
     1, 0, 0, 0},
    {NEW_PDN_TYPE_SAB_ONLY,
     const_cast<char*>("New PDN type single address bearer only"), 1, 0, 0, 0},
    {CONTEXT_NOT_FOUND, const_cast<char*>("Context not found"), 0, 1, 1, 1},
    {INVALID_MESSAGE_FORMAT, const_cast<char*>("Invalid message format"), 1, 1,
     1, 1},
    {INVALID_LENGTH, const_cast<char*>("Invalid length"), 1, 1, 1, 0},
    {SERVICE_NOT_SUPPORTED, const_cast<char*>("Service not supported"), 0, 0, 1,
     0},
    {SYSTEM_FAILURE, const_cast<char*>("System failure"), 1, 1, 1, 0},
    {NO_RESOURCES_AVAILABLE, const_cast<char*>("No resources available"), 1, 0,
     0, 0},
    {MISSING_OR_UNKNOWN_APN, const_cast<char*>("Missing or unknown APN"), 1, 0,
     0, 0},
    {GRE_KEY_NOT_FOUND, const_cast<char*>("GRE KEY not found"), 1, 0, 0, 0},
    {DENIED_IN_RAT, const_cast<char*>("Denied in RAT"), 1, 1, 0, 0},
    {UE_NOT_RESPONDING, const_cast<char*>("UE not responding"), 0, 1, 0, 0},
    {SERVICE_DENIED, const_cast<char*>("Service Denied"), 0, 0, 0, 0},
    {UNABLE_TO_PAGE_UE, const_cast<char*>("Unable to page UE"), 0, 1, 0, 0},
    {NO_MEMORY_AVAILABLE, const_cast<char*>("No memory available"), 1, 1, 1, 0},
    {REQUEST_REJECTED, const_cast<char*>("Request rejected"), 1, 1, 1, 0},
    {INVALID_PEER, const_cast<char*>("Invalid peer"), 0, 0, 0, 1},
    {TEMP_REJECT_HO_IN_PROGRESS,
     const_cast<char*>("Temporarily rejected due to HO in progress"), 0, 0, 0,
     0},
    {M_PDN_APN_NOT_ALLOWED,
     const_cast<char*>("Multiple PDN for a given APN not allowed"), 1, 0, 0, 0},
    {0, NULL, 0, 0, 0, 0},
};

static int compare_cause_id(const void* m1, const void* m2) {
  struct SGWCauseMapping_e* scm1 = (struct SGWCauseMapping_e*)m1;
  struct SGWCauseMapping_e* scm2 = (struct SGWCauseMapping_e*)m2;

  return (scm1->value - scm2->value);
}

char* sgw_cause_2_string(uint8_t cause_value) {
  SGWCauseMapping_t *res, key;

  key.value = cause_value;
  res = reinterpret_cast<SGWCauseMapping_t*>(
      bsearch((const void*)&key, causes, sizeof(causes),
              sizeof(SGWCauseMapping_t), compare_cause_id));

  if (res == NULL) {
    return const_cast<char*>("Unknown cause");
  } else {
    return res->name;
  }
}
