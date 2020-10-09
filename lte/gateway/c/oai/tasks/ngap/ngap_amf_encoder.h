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

#include <stdint.h>
#include "Ngap_NGAP-PDU.h"
#ifndef FILE_NGAP_AMF_ENCODER_SEEN
#define FILE_NGAP_AMF_ENCODER_SEEN

int ngap_amf_encode_pdu(
    Ngap_NGAP_PDU_t* message, uint8_t** buffer, uint32_t* len)
    __attribute__((warn_unused_result));

#endif /* FILE_NGAP_AMF_ENCODER_SEEN */
