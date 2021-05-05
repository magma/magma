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
#pragma once

#include "Ngap_Cause.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_ProtocolIE-Field.h"
#include "bstrlib.h"

#ifdef __cplusplus
extern "C" {
#endif

#include "ngap_amf_encoder.h"

#ifdef __cplusplus
}
#endif

#define NGAP_FIND_PROTOCOLIE_BY_ID(IE_TYPE, ie, container, IE_ID)              \
  do {                                                                         \
    IE_TYPE** ptr;                                                             \
    ie = NULL;                                                                 \
    for (ptr = container->protocolIEs.list.array;                              \
         ptr < &container->protocolIEs.list                                    \
                    .array[container->protocolIEs.list.count];                 \
         ptr++) {                                                              \
      if ((*ptr)->id == IE_ID) {                                               \
        ie = *ptr;                                                             \
        break;                                                                 \
      }                                                                        \
    }                                                                          \
  } while (0)

// Base test function
int ngap_ng_setup_failure_stream(
    const Ngap_Cause_PR cause_type, const long cause_value, bstring& stream);

int ngap_ng_setup_failure_pdu(
    const Ngap_Cause_PR cause_type, const long cause_value,
    Ngap_NGAP_PDU_t& encode_pdu);

bool ng_setup_failure_decode(const_bstring const raw, Ngap_NGAP_PDU_t* pdu);
