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
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_decoder.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_encoder.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_nas_procedures.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_types.h"

#define NGAP_TEST_PDU_FETCH_AMF_SET_ID_FROM_PDU(aSN, Amf_Set_Id) \
  DevCheck((aSN).size == 2, (aSN).size, 0, 0);                   \
  DevCheck((aSN).bits_unused == 6, (aSN).bits_unused, 6, 0);     \
  Amf_Set_Id = (aSN.buf[0] << 2) + ((aSN.buf[1] >> 6) & 0x03);

#define NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(IE_TYPE, ie, container, IE_ID) \
  do {                                                                     \
    IE_TYPE** ptr;                                                         \
    ie = NULL;                                                             \
    for (ptr = container->protocolIEs.list.array;                          \
         ptr < &container->protocolIEs.list                                \
                    .array[container->protocolIEs.list.count];             \
         ptr++) {                                                          \
      if ((*ptr)->id == IE_ID) {                                           \
        ie = *ptr;                                                         \
        break;                                                             \
      }                                                                    \
    }                                                                      \
  } while (0)

// Base test function
int ngap_ng_setup_failure_stream(const Ngap_Cause_PR cause_type,
                                 const long cause_value, bstring& stream);

int ngap_ng_setup_failure_pdu(const Ngap_Cause_PR cause_type,
                              const long cause_value,
                              Ngap_NGAP_PDU_t& encode_pdu);

bool ng_setup_failure_decode(const_bstring const raw, Ngap_NGAP_PDU_t* pdu);

bool ngap_initiate_ue_message(bstring& stream);

bool generator_ngap_pdusession_resource_setup_req(bstring& stream);

bool generator_itti_ngap_pdusession_resource_setup_req(bstring& stream);

bool generator_ngap_pdusession_resource_rel_cmd_stream(bstring& stream);

bool generate_guti_ngap_pdu(Ngap_NGAP_PDU_t* pdu);

bool generate_ngap_request_msg(Ngap_NGAP_PDU_t* pdu);

bool validate_ngap_setup_request(Ngap_NGAP_PDU_t* pdu);

bool validate_handle_initial_ue_message(gnb_description_t* gNB_ref,
                                        m5g_ue_description_t* ue_ref,
                                        Ngap_NGAP_PDU_t* pdu);

status_code_e send_ngap_gnb_reset_ack();
